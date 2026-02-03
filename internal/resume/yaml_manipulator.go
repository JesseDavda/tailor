package resume

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLDocument struct {
	root *yaml.Node
}

type PathSegment struct {
	key   string
	index int // -1 if not an array access
}

func ParseYAMLDocument(content string) (*YAMLDocument, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		return &YAMLDocument{root: root.Content[0]}, nil
	}

	return &YAMLDocument{root: &root}, nil
}

func (d *YAMLDocument) ToYAML() (string, error) {
	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(d.root); err != nil {
		return "", fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return "", fmt.Errorf("failed to close encoder: %w", err)
	}

	return buf.String(), nil
}

func parsePath(path string) ([]PathSegment, error) {
	var segments []PathSegment

	re := regexp.MustCompile(`([^.\[\]]+)(?:\[(\d+)\])?`)
	matches := re.FindAllStringSubmatch(path, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	for _, match := range matches {
		if len(match) < 2 || match[1] == "" {
			continue
		}

		segment := PathSegment{
			key:   match[1],
			index: -1,
		}

		if len(match) > 2 && match[2] != "" {
			idx, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, fmt.Errorf("invalid array index in path: %s", match[2])
			}
			segment.index = idx
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

func (d *YAMLDocument) GetNodeByPath(path string) (*yaml.Node, error) {
	segments, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	return d.traverseToNode(d.root, segments)
}

func (d *YAMLDocument) traverseToNode(current *yaml.Node, segments []PathSegment) (*yaml.Node, error) {
	if len(segments) == 0 {
		return current, nil
	}

	segment := segments[0]
	remaining := segments[1:]

	switch current.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(current.Content); i += 2 {
			keyNode := current.Content[i]
			valueNode := current.Content[i+1]

			if keyNode.Value == segment.key {
				if segment.index >= 0 {
					if valueNode.Kind != yaml.SequenceNode {
						return nil, fmt.Errorf("expected array at %s, got %v", segment.key, valueNode.Kind)
					}
					if segment.index >= len(valueNode.Content) {
						return nil, fmt.Errorf("array index %d out of bounds for %s (length: %d)",
							segment.index, segment.key, len(valueNode.Content))
					}
					return d.traverseToNode(valueNode.Content[segment.index], remaining)
				}
				return d.traverseToNode(valueNode, remaining)
			}
		}
		return nil, fmt.Errorf("key %s not found in mapping", segment.key)

	case yaml.SequenceNode:
		if segment.index < 0 {
			return nil, fmt.Errorf("array access requires an index")
		}
		if segment.index >= len(current.Content) {
			return nil, fmt.Errorf("array index %d out of bounds (length: %d)", segment.index, len(current.Content))
		}
		return d.traverseToNode(current.Content[segment.index], remaining)

	default:
		return nil, fmt.Errorf("unexpected node kind: %v", current.Kind)
	}
}

func (d *YAMLDocument) SetNodeValue(path string, value any) error {
	segments, err := parsePath(path)
	if err != nil {
		return err
	}

	if len(segments) == 0 {
		return fmt.Errorf("empty path")
	}

	parentSegments := segments[:len(segments)-1]
	lastSegment := segments[len(segments)-1]

	var parent *yaml.Node
	if len(parentSegments) == 0 {
		parent = d.root
	} else {
		parent, err = d.traverseToNode(d.root, parentSegments)
		if err != nil {
			return fmt.Errorf("failed to navigate to parent: %w", err)
		}
	}

	return d.setValueInNode(parent, lastSegment, value)
}

func (d *YAMLDocument) setValueInNode(parent *yaml.Node, segment PathSegment, value any) error {
	switch parent.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(parent.Content); i += 2 {
			keyNode := parent.Content[i]
			if keyNode.Value == segment.key {
				valueNode := parent.Content[i+1]

				if segment.index >= 0 {
					if valueNode.Kind != yaml.SequenceNode {
						return fmt.Errorf("expected array at %s", segment.key)
					}
					if segment.index >= len(valueNode.Content) {
						return fmt.Errorf("array index %d out of bounds", segment.index)
					}
					return d.updateNodeValue(valueNode.Content[segment.index], value)
				}

				return d.updateNodeValue(valueNode, value)
			}
		}
		return fmt.Errorf("key %s not found", segment.key)

	case yaml.SequenceNode:
		if segment.index < 0 {
			return fmt.Errorf("array access requires an index")
		}
		if segment.index >= len(parent.Content) {
			return fmt.Errorf("array index %d out of bounds", segment.index)
		}
		return d.updateNodeValue(parent.Content[segment.index], value)

	default:
		return fmt.Errorf("cannot set value in node of kind %v", parent.Kind)
	}
}

func (d *YAMLDocument) updateNodeValue(node *yaml.Node, value any) error {
	if node.Kind == yaml.ScalarNode {
		node.Value = fmt.Sprintf("%v", value)
		return nil
	}

	valueBytes, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	var newNode yaml.Node
	if err := yaml.Unmarshal(valueBytes, &newNode); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	if newNode.Kind == yaml.DocumentNode && len(newNode.Content) > 0 {
		*node = *newNode.Content[0]
	} else {
		*node = newNode
	}

	return nil
}

func (d *YAMLDocument) ReorderArray(path string, newOrder []int) error {
	node, err := d.GetNodeByPath(path)
	if err != nil {
		return err
	}

	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("path does not point to an array")
	}

	if len(newOrder) != len(node.Content) {
		return fmt.Errorf("newOrder length (%d) does not match array length (%d)",
			len(newOrder), len(node.Content))
	}

	oldContent := make([]*yaml.Node, len(node.Content))
	copy(oldContent, node.Content)

	newContent := make([]*yaml.Node, len(newOrder))
	for newIdx, oldIdx := range newOrder {
		if oldIdx < 0 || oldIdx >= len(oldContent) {
			return fmt.Errorf("invalid index %d in newOrder", oldIdx)
		}
		newContent[newIdx] = oldContent[oldIdx]
	}

	node.Content = newContent
	return nil
}

func (d *YAMLDocument) AddArrayElement(path string, value any, position int) error {
	node, err := d.GetNodeByPath(path)
	if err != nil {
		return err
	}

	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("path does not point to an array")
	}

	valueBytes, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	var valueNode yaml.Node
	if err := yaml.Unmarshal(valueBytes, &valueNode); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	var nodeToInsert *yaml.Node
	if valueNode.Kind == yaml.DocumentNode && len(valueNode.Content) > 0 {
		nodeToInsert = valueNode.Content[0]
	} else {
		nodeToInsert = &valueNode
	}

	if position < 0 || position > len(node.Content) {
		position = len(node.Content)
	}

	newContent := make([]*yaml.Node, 0, len(node.Content)+1)
	newContent = append(newContent, node.Content[:position]...)
	newContent = append(newContent, nodeToInsert)
	newContent = append(newContent, node.Content[position:]...)

	node.Content = newContent
	return nil
}

func (d *YAMLDocument) RemoveArrayElement(path string, index int) error {
	node, err := d.GetNodeByPath(path)
	if err != nil {
		return err
	}

	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("path does not point to an array")
	}

	if index < 0 || index >= len(node.Content) {
		return fmt.Errorf("index %d out of bounds (length: %d)", index, len(node.Content))
	}

	newContent := make([]*yaml.Node, 0, len(node.Content)-1)
	newContent = append(newContent, node.Content[:index]...)
	newContent = append(newContent, node.Content[index+1:]...)

	node.Content = newContent
	return nil
}

func (d *YAMLDocument) GetArrayLength(path string) (int, error) {
	node, err := d.GetNodeByPath(path)
	if err != nil {
		return 0, err
	}

	if node.Kind != yaml.SequenceNode {
		return 0, fmt.Errorf("path does not point to an array")
	}

	return len(node.Content), nil
}

func (d *YAMLDocument) GetScalarValue(path string) (string, error) {
	node, err := d.GetNodeByPath(path)
	if err != nil {
		return "", err
	}

	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("path does not point to a scalar value")
	}

	return node.Value, nil
}
