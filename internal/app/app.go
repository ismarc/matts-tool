package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// APIConfig provides an interface for API related configuration options
type APIConfig struct {
	SourceConjurRC      string
	SourceNetRC         string
	SourceVersion       string
	DestinationConjurRC string
	DestinationNetRC    string
	DestinationVersion  string
	NoAct               bool
	ContinueOnError     bool
	ResourceBatchSize   int
	VariableBatchSize   int
}

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

func (p policyLoader) loadPolicyYaml(inData []byte, stripAnnotations bool) (result yaml.Node, err error) {
	processor := IncludeProcessor{&result, p, stripAnnotations}
	err = yaml.Unmarshal(inData, &processor)
	return
}

// RunPolicy is the main entrypoint for the policy subcommand
func RunPolicy(inputPolicyFile string, stripAnnotations bool) {
	loader := policyLoader{filepath.Dir(inputPolicyFile)}

	data, err := loader.loadFile(filepath.Base(inputPolicyFile))
	if err != nil {
		panic(err.Error())
	}

	result, err := loader.loadPolicyYaml(data, stripAnnotations)
	out, err := yaml.Marshal(result)
	fmt.Printf("%+v\n", string(out))
}

// RunAPI is the main entrypoint for the db subcommand
func RunAPI(config APIConfig) {
	source, err := loadAPI(config.SourceConjurRC, config.SourceVersion, config.SourceNetRC)
	if err != nil {
		panic(err)
	}
	fmt.Printf("source: %s\n", source.GetConfig().ApplianceURL)
	destination, err := loadAPI(config.DestinationConjurRC, config.DestinationVersion, config.DestinationNetRC)
	if err != nil {
		panic(err)
	}
	fmt.Printf("destination: %s\n", destination.GetConfig().ApplianceURL)

	resources, err := loadResources(source, config.ResourceBatchSize)
	if err != nil {
		panic(err)
	}

	if !config.NoAct {
		errors := syncResources(resources, source, destination, config.VariableBatchSize, config.ContinueOnError)
		if len(errors) != 0 {
			fmt.Printf("Errors received:\n")
			for k, v := range errors {
				fmt.Printf("Error: %s\n", k)
				for _, value := range v {
					if !strings.Contains(k, value) {
						fmt.Printf("  %s\n", value)
					}
				}
			}
		}
	} else {
		fmt.Printf("Would sync resources:\n")
		for _, resource := range resources {
			fmt.Printf("Id: %s\n", resource)
		}
	}
}
