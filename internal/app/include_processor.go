package app

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Fragment struct {
	content *yaml.Node
	loader  policyLoader
}

type IncludeProcessor struct {
	target interface{}
	loader policyLoader
}

func (fragment *Fragment) UnmarshalYAML(value *yaml.Node) error {
	var err error

	fragment.content, err = processIncludes(value, fragment.loader)
	return err
}

func (include *IncludeProcessor) UnmarshalYAML(value *yaml.Node) error {
	processed, err := processIncludes(value, include.loader)

	if err != nil {
		return err
	}

	return processed.Decode(include.target)
}

func processIncludes(node *yaml.Node, loader policyLoader) (*yaml.Node, error) {
	if node.Tag == "!include" {
		if node.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("!include on a non-scalar node")
		}

		data, err := loader.loadFile(node.Value)
		if err != nil {
			return nil, err
		}

		fragment := Fragment{loader: loader}
		err = yaml.Unmarshal(data, &fragment)
		return fragment.content, err
	}

	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var content []*yaml.Node
		for i := range node.Content {
			// Remove `replace` entries for v5 compatibility
			if node.Content[i].Value != "replace" {
				entry, err := processIncludes(node.Content[i], loader)
				if err != nil {
					return nil, err
				}
				content = append(content, entry)
			}
		}
		node.Content = content
	}
	return node, nil
}
