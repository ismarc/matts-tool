package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func loadFile(inputFile string) ([]byte, error) {
	var inFile io.Reader
	var err error

	switch inputFile {
	case "-":
		inFile = os.Stdin
	default:
		inFile, err = os.Open(inputFile)
		if err != nil {
			panic(err)
		}
	}

	return ioutil.ReadAll(inFile)
}

func loadPolicyYaml(inData []byte) (result yaml.Node, err error) {
	var foo yaml.Node
	err = yaml.Unmarshal(inData, &IncludeProcessor{&foo})
	result = foo
	return
}

// Run is the main entrypoint for the policy-handler
func Run(inputPolicyFile string) {
	data, err := loadFile(inputPolicyFile)
	if err != nil {
		panic(err.Error())
	}

	result, err := loadPolicyYaml(data)
	out, err := yaml.Marshal(result)
	fmt.Printf("%+v\n", string(out))
}
