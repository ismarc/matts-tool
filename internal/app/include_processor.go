package app

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Fragment struct {
	content *yaml.Node
}

type IncludeProcessor struct {
	target interface{}
}

func (fragment *Fragment) UnmarshalYAML(value *yaml.Node) error {
	var err error

	fragment.content, err = processIncludes(value)
	return err
}

func (include *IncludeProcessor) UnmarshalYAML(value *yaml.Node) error {
	processed, err := processIncludes(value)

	if err != nil {
		return err
	}

	return processed.Decode(include.target)
}

func processIncludes(node *yaml.Node) (*yaml.Node, error) {
	// Skip `replace` nodes because they aren't v5 compatible
	if node.Value == "replace" {
		node.Value = ""
		return node, nil
	}
	if node.Tag == "!include" {
		if node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("!include on a non-scalar node")
		}

		data, err := loadFile(node.Value)
		if err != nil {
			return nil, err
		}

		var fragment Fragment
		err = yaml.Unmarshal(data, &fragment)
		return fragment.content, err
	}

	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var err error
		for i := range node.Content {
			node.Content[i], err = processIncludes(node.Content[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return node, nil
}
