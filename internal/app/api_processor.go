package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

func loadApi(conjurrc string, version string) (conjur *conjurapi.Client, err error) {
	majorVersion := os.Getenv("CONJUR_MAJOR_VERSION")
	conjurVersion := os.Getenv("CONJUR_VERSION")
	os.Setenv("CONJUR_MAJOR_VERSION", version)
	os.Setenv("CONJUR_VERSION", version)
	defer os.Setenv("CONJUR_MAJOR_VERSION", majorVersion)
	defer os.Setenv("CONJUR_VERSION", conjurVersion)

	originalConjurRC := os.Getenv("CONJURRC")
	os.Setenv("CONJURRC", conjurrc)
	defer os.Setenv("CONJURRC", originalConjurRC)

	config, err := conjurapi.LoadConfig()
	if err != nil {
		return
	}

	conjur, err = conjurapi.NewClientFromEnvironment(config)
	if err != nil {
		return
	}

	return
}

func loadResources(conjur *conjurapi.Client) (result []string, err error) {
	batchSize := 25
	offset := 0
	resources, err := conjur.Resources(&conjurapi.ResourceFilter{Kind: "variable", Limit: batchSize, Offset: offset})

	for len(resources) == batchSize {
		if err != nil {
			return
		}

		for _, resource := range resources {
			result = append(result, resource["id"].(string))
		}

		offset += 25
		resources, err = conjur.Resources(&conjurapi.ResourceFilter{Kind: "variable", Limit: batchSize, Offset: offset})
	}

	return
}

func syncResources(resources []string, source *conjurapi.Client, destination *conjurapi.Client) (err error) {
	batchSize := 10
	account := source.GetConfig().Account
	variablePrefix := fmt.Sprintf("%s:variable:", account)
	resources = resources[len(resources)-10:]
	var variables []string
	for _, resource := range resources {
		if strings.HasPrefix(resource, variablePrefix) {
			variables = append(variables, strings.TrimPrefix(resource, variablePrefix))
		}
	}

	// Uncomment the following to only operate on the last 10 variables instead of all 600+
	// variables = variables[len(variables)-10:]

	for index := 0; index < len(variables); index += batchSize {
		end := index + batchSize
		batch := variables[index:end]
		data, err := source.RetrieveBatchSecrets(batch)
		if err != nil {
			return err
		}

		for variable, value := range data {
			addSecret(destination, variable, string(value))
		}
	}

	return
}

func addSecret(destination *conjurapi.Client, variable string, value string) {
	fmt.Printf("Would write secret: %s\n", variable)
	// Uncomment the following to write the secret to the destination.  Should theoretically work, has not been run
	// destination.AddSecret(variable, value)
}
