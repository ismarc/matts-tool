package app

import (
	"fmt"
	"regexp"

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

	fragment.content, err = processIncludes(value, fragment.loader, "")
	return err
}

func (include *IncludeProcessor) UnmarshalYAML(value *yaml.Node) error {
	processed, err := processIncludes(value, include.loader, "")

	if err != nil {
		return err
	}

	return processed.Decode(include.target)
}

var seenAnchors = make(map[string]bool)

func excludedTag(value string) bool {
	switch value {
	case
		"replace",
		"role_name",
		"record":
		return true
	default:
		return false
	}
	return false
}

func processIncludes(node *yaml.Node, loader policyLoader, incomingId string) (*yaml.Node, error) {
	// Get the id for the currently being processed tree
	if node.Kind == yaml.MappingNode {
		for i := range node.Content {
			if node.Content[i].Value == "id" {
				r1 := regexp.MustCompile(`\W`)
				incomingId = string(r1.ReplaceAll([]byte(node.Content[i+1].Value), []byte("_")))
			}
		}
	}

	if node.Anchor != "" {
		// Duplicate ids (because they can be nested) can have duplicate anchors
		// Work around this by incrementing each time it's seen
		anchor := node.Anchor
		for i := 0; seenAnchors[anchor]; i++ {
			anchor = fmt.Sprintf("%s_%d_%s", incomingId, i, node.Anchor)
		}
		node.Anchor = anchor
		seenAnchors[anchor] = true
	}
	// An alias has to follow an anchor, and refers to the most recent.  The related
	// anchor to the alias will be <id>_<highest counter>_<anchor>
	if node.Kind == yaml.AliasNode {
		anchor := node.Value
		maxIndex := 0
		for i := maxIndex; seenAnchors[anchor]; i++ {
			anchor = fmt.Sprintf("%s_%d_%s", incomingId, i, node.Value)
			maxIndex = i - 1
		}
		if maxIndex > 0 || seenAnchors[anchor] {
			node.Value = anchor
		}
	}

	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var content []*yaml.Node
		removeNode := false
		for i := range node.Content {
			// Remove excluded tags for v5 compatibility
			if excludedTag(node.Content[i].Value) {
				removeNode = true
				continue
			}
			// Remove the content of the replace node as well
			if removeNode {
				removeNode = false
				continue
			}
			if node.Content[i].Tag == "!include" && node.Content[i].Kind == yaml.ScalarNode {
				data, err := loader.loadFile(node.Content[i].Value)
				if err != nil {
					return nil, err
				}
				fragment := Fragment{loader: loader}
				err = yaml.Unmarshal(data, &fragment)
				content = append(content, fragment.content)
			} else {
				if node.Content[i].Tag == "!automatic-role" {
					if len(content) > 0 {
						content = content[:len(content)-1]
					}
				} else {
					entry, err := processIncludes(node.Content[i], loader, incomingId)
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
