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

// Run is the main entrypoint for the policy-handler
func Run(inputPolicyFile string) {
	loader := policyLoader{filepath.Dir(inputPolicyFile)}

	data, err := loader.loadFile(filepath.Base(inputPolicyFile))
	if err != nil {
		panic(err.Error())
	}

	result, err := loader.loadPolicyYaml(data)
	out, err := yaml.Marshal(result)
	fmt.Printf("%+v\n", string(out))
}
