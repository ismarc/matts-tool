#!/usr/bin/env bash

set -euo pipefail

# End to end test for datakey rotation feature of
#  __  __       _   _   _       _____           _
# |  \/  | __ _| |_| |_( )___  |_   _|__   ___ | |
# | |\/| |/ _` | __| __|// __|   | |/ _ \ / _ \| |
# | |  | | (_| | |_| |_  \__ \   | | (_) | (_) | |
# |_|  |_|\__,_|\__|\__| |___/   |_|\___/ \___/|_|

get_container_logs(){
  log_dir="rotation-e2e-logs"
  mkdir -p "${log_dir}"
  docker compose ps --format json \
    |jq -r  '.[]|.Name' \
    |while read -r container; do \
      docker logs "${container}" &> "${log_dir}/${container}.log"
    done

    echo "Logs written to ${log_dir}:"
    ls "${log_dir}"
}

deploy_conjur(){
    echo "🔑 Starting Quickstart Auto"
  docker compose down -v||:
  echo "✅ cleanup"
  export CONJUR_DATA_KEY="$(< data_key)"
  echo "✅ read data key: ${CONJUR_DATA_KEY}"
  docker compose up -d
  echo "✅ start containers"
  docker compose exec conjur conjurctl wait
  docker compose exec conjur conjurctl account create conjur > admin_data
  echo "✅ create account"
  admin_api_key="$(awk '/API key for admin:/{print $NF}' admin_data|tr -d '\r')"
  echo "✅ read api key: ${admin_api_key}"
  docker compose exec client conjur init --insecure --url http://conjur:80 --account conjur
  echo "✅ init cli"
  docker compose exec client conjur login -i admin -p "${admin_api_key}"
  echo "✅ login cli"
  echo '- !policy
  id: test
  body:
    - !variable test' > test.yml
  docker cp test.yml conjur_client:/test.yml
  echo "✅ Write test policy file"
  docker compose exec client conjur policy load --branch root --file /test.yml
  echo "✅ load test policy"
  docker compose exec client conjur list
  echo "✅ list objects"
  docker compose exec client conjur variable set -i test/test -v firstvalue
  docker compose exec client conjur variable set -i test/test -v t3st
  echo "✅ set test value test/test = t3st"
  docker compose exec client conjur user  change-password  -p "${admin_password}"
  echo "✅ set admin password: ${admin_password}"
  echo "🏁 Quickstart Auto Complete"
}

check_conjur(){
    if docker compose exec client conjur list |grep test/test; then
        echo "✅ list check passed, test/test found"
    else
        echo "❌ list check failed, test/test not found"
    fi

    if [[ "$(docker compose exec client conjur variable get -i test/test)" == "t3st" ]]; then
        echo "✅ value check passed, test/test = t3st"
    else
        echo "❌ value check failed, test/test != t3st"
    fi
}

run_checks_with_password_and_api_key(){
    echo "Logging in with admin api key: ${admin_api_key}"
    docker compose exec client conjur login -i admin -p "${admin_api_key}"
    check_conjur
    echo "Logging in with admin password: ${admin_password}"
    docker compose exec client conjur login -i admin -p "${admin_password}"
    check_conjur
}

check_tool(){
    tool="${1}"
    echo -n "${tool} "
    if command -v $tool >/dev/null; then
        echo "✅"
    else
        echo "❌ not found"
        exit 1
    fi
}

echo "🔵 Starting End-To-End test of Conjur Data Key rotation"

echo "🔵 Check for required tools"
check_tool git
check_tool docker
check_tool jq

if [[ "$(docker compose version |grep -o ' v[0-9]\+'|sed 's/v//')"  -le 1 ]]; then
    echo "❌ docker compose major version 2 or higher is required. V1 restarts dependencies when restarting a container which"
    echo " breaks the test as postgres gets recreated."
    exit 1
else
    echo "docker-compose v2 or higher ✅"
fi

# Trap won't work if the tools don't exist, so don't set it
# till after we know we have the tools.
trap get_container_logs EXIT

repoRoot="$(git rev-parse --show-toplevel)"
git submodule update --init
pushd "${repoRoot}/conjur-quickstart"

echo "🔵 Building Matt's Tool"
../scripts/build
cp ../output/mt .
chmod +x mt

echo "🔵 Generating Data Keys"
admin_password="abcABC123---"
docker compose run --no-deps --rm conjur data-key generate > data_key.in
docker compose run --no-deps --rm conjur data-key generate > data_key.out

echo "🔵 Deploying Conjur With Initial Data Key"
cp data_key.in data_key
deploy_conjur

echo "🔵 Running checks before rotation"
run_checks_with_password_and_api_key

echo "🔵 Rotating Data Key"
export IN_CONJUR_DATA_KEY="$(<data_key.in)"
export OUT_CONJUR_DATA_KEY="$(<data_key.out)"
export cs="postgresql://postgres:SuperSecretPg@localhost:8432"
./mt rotate-datakey --dsn "${cs}"

echo "🔵 Restarting Conjur with new data key"
export CONJUR_DATA_KEY="$(<data_key.out)"
docker compose up -d --force-recreate conjur
docker compose exec conjur conjurctl wait

running_data_key="$(docker compose exec conjur sh -c 'echo $CONJUR_DATA_KEY')"
# %? strips trailing newline, which causes this to fail on some versions of bash/docker
if [[ "${running_data_key}" != "$(<data_key.out)" ]]; then
    echo "Conjur container does not have the new data key after restart 🙁 abort"
    exit 1
else
    echo "✅ Verified Conjur has started with new data key"
fi

echo "🔵 Running Checks Post Rotation"
run_checks_with_password_and_api_key

echo "✅ DataKey Rotation E2E test complete"
