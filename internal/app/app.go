package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type policyLoader struct {
	filePath string
}

func (p policyLoader) loadFile(inputFile string) ([]byte, error) {
	var inFile io.Reader
	var err error

	switch inputFile {
	case "-":
		inFile = os.Stdin
	default:
		inFile, err = os.Open(filepath.Join(p.filePath, inputFile))
		if err != nil {
			panic(err)
		}
	}

	return ioutil.ReadAll(inFile)
}

func (p policyLoader) loadPolicyYaml(inData []byte) (result yaml.Node, err error) {
	processor := IncludeProcessor{&result, p}
	// err = yaml.Unmarshal(inData, &IncludeProcessor{&result, "foo"})
	err = yaml.Unmarshal(inData, &processor)
	return
}

// RunPolicy is the main entrypoint for the policy subcommand
func RunPolicy(inputPolicyFile string) {
	loader := policyLoader{filepath.Dir(inputPolicyFile)}

	data, err := loader.loadFile(filepath.Base(inputPolicyFile))
	if err != nil {
		panic(err.Error())
	}

	result, err := loader.loadPolicyYaml(data)
	out, err := yaml.Marshal(result)
	fmt.Printf("%+v\n", string(out))
}

// RunDB is the main entrypoint for the db subcommand
func RunDB(sourceConjurRC string, sourceVersion string, destinationConjurRC string, destinationVersion string, noAct bool) {
	source, err := loadApi(sourceConjurRC, sourceVersion)
	if err != nil {
		panic(err)
	}
	fmt.Printf("source: %s\n", source.GetConfig().ApplianceURL)
	destination, err := loadApi(destinationConjurRC, destinationVersion)
	if err != nil {
		panic(err)
	}
	fmt.Printf("destination: %s\n", destination.GetConfig().ApplianceURL)

	resources, err := loadResources(source)
	if err != nil {
		panic(err)
	}

	if !noAct {
		err = syncResources(resources, source, destination)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Would sync resources:\n")
		for _, resource := range resources {
			fmt.Printf("Id: %s\n", resource)
		}
	}
}
