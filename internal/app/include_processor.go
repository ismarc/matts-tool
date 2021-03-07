package app

import (
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
	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var content []*yaml.Node
		for i := range node.Content {
			// Remove `replace` entries for v5 compatibility
			if node.Content[i].Value != "replace" {
				if node.Content[i].Tag == "!include" && node.Content[i].Kind == yaml.ScalarNode {
					data, err := loader.loadFile(node.Content[i].Value)
					if err != nil {
						return nil, err
					}
					fragment := Fragment{loader: loader}
					err = yaml.Unmarshal(data, &fragment)
					content = append(content, fragment.content)
				} else {
					entry, err := processIncludes(node.Content[i], loader)
					if err != nil {
						return nil, err
					}
					content = append(content, entry)
				}
			}
		}
		node.Content = content
	}
	return node, nil
}
