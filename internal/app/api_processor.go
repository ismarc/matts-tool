package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

func loadAPI(conjurrc string, version string, netrc string) (conjur *conjurapi.Client, err error) {
	majorVersion := os.Getenv("CONJUR_MAJOR_VERSION")
	conjurVersion := os.Getenv("CONJUR_VERSION")
	os.Setenv("CONJUR_MAJOR_VERSION", version)
	os.Setenv("CONJUR_VERSION", version)
	defer os.Setenv("CONJUR_MAJOR_VERSION", majorVersion)
	defer os.Setenv("CONJUR_VERSION", conjurVersion)

	originalConjurRC := os.Getenv("CONJURRC")
	os.Setenv("CONJURRC", conjurrc)
	defer os.Setenv("CONJURRC", originalConjurRC)

	originalNetrc := os.Getenv("CONJUR_NETRC_PATH")
	os.Setenv("CONJUR_NETRC_PATH", netrc)
	defer os.Setenv("CONJUR_NETRC_PATH", originalNetrc)

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

func loadResources(conjur *conjurapi.Client, batchSize int) (result []string, err error) {
	offset := 0
	resources, err := conjur.Resources(&conjurapi.ResourceFilter{Kind: "variable", Limit: batchSize, Offset: offset})
	if err != nil {
		return
	}

	for _, resource := range resources {
		result = append(result, resource["id"].(string))
	}

	for len(resources) == batchSize {
		resources, err = conjur.Resources(&conjurapi.ResourceFilter{Kind: "variable", Limit: batchSize, Offset: offset})
		if err != nil {
			return
		}

		for _, resource := range resources {
			result = append(result, resource["id"].(string))
		}

		offset += 25
	}

	return
}

func syncResources(resources []string, source *conjurapi.Client, destination *conjurapi.Client, batchSize int, continueOnError bool) (errors map[string][]string) {
	errors = make(map[string][]string)
	account := source.GetConfig().Account
	variablePrefix := fmt.Sprintf("%s:variable:", account)

	var variables []string
	for _, resource := range resources {
		if strings.HasPrefix(resource, variablePrefix) {
			variables = append(variables, strings.TrimPrefix(resource, variablePrefix))
		}
	}

	for index := 0; index < len(variables); index += batchSize {
		end := index + batchSize
		if end > len(variables) {
			end = len(variables)
		}
		batch := variables[index:end]
		data, err := source.RetrieveBatchSecrets(batch)
		if err != nil {
			// An error in the batch means that at least one had an error response, not that they all had an error
			// Attempt each item in turn so variables aren't missed
			data = make(map[string][]byte)
			for _, entry := range batch {
				value, err := source.RetrieveSecret(entry)
				if err != nil {
					if errors[err.Error()] != nil {
						errors[err.Error()] = append(errors[err.Error()], entry)
					} else {
						errors[err.Error()] = []string{entry}
					}
					if !continueOnError {
						return
					}
				} else {
					data[entry] = value
				}
			}
		}

		for variable, value := range data {
			err := addSecret(destination, variable, string(value))
			if err != nil {
				if errors[err.Error()] != nil {
					errors[err.Error()] = append(errors[err.Error()], variable)
				} else {
					errors[err.Error()] = []string{variable}
				}
				if !continueOnError {
					return
				}
			}
		}
	}

	return
}

func addSecret(destination *conjurapi.Client, variable string, value string) error {
	// Ignore errors, because it could indicate a missing value which should be set
	destinationSecret, _ := destination.RetrieveSecret(variable)

	if string(destinationSecret) == value {
		fmt.Printf("Skipping variable: %s\n", variable)
		return nil
	} else {
		fmt.Printf("Writing variable: %s\n", variable)
		return destination.AddSecret(variable, value)
	}
}
