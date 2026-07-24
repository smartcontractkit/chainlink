package infra

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// PatchChartValues injects a TOML configuration layer into the chart values
// YAML file for a specific node. The layer is named "30-cre" and is inserted
// under chainlink-node.instances.<nodeName>.configuration.
//
// If the layer already exists, its content is replaced. If the node's
// configuration list doesn't exist, it's created.
func PatchChartValues(yamlPath string, patches map[string]string) error {
	// patches maps nodeName -> TOML content to inject as the "30-cre" layer

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read chart values %s", yamlPath)
	}

	var root yaml.Node
	if err = yaml.Unmarshal(data, &root); err != nil {
		return errors.Wrapf(err, "failed to parse YAML %s", yamlPath)
	}

	// Ensure we have a mapping node at the root
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return fmt.Errorf("unexpected YAML structure in %s", yamlPath)
	}

	docMap := root.Content[0]
	if docMap.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping at root of %s", yamlPath)
	}

	// Navigate to chainlink-node.instances.<name>.configuration
	clNode := FindKey(docMap, "chainlink-node")
	if clNode == nil {
		return fmt.Errorf("chainlink-node not found in %s", yamlPath)
	}

	instances := FindKey(clNode, "instances")
	if instances == nil {
		return fmt.Errorf("chainlink-node.instances not found in %s", yamlPath)
	}

	for nodeName, tomlContent := range patches {
		nodeEntry := FindKey(instances, nodeName)
		if nodeEntry == nil {
			return fmt.Errorf("node %s not found in chainlink-node.instances in %s", nodeName, yamlPath)
		}

		configList := findOrCreateKey(nodeEntry, "configuration", yaml.SequenceNode)
		updateOrCreateLayer(configList, "30-cre", tomlContent)
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return errors.Wrap(err, "failed to marshal patched YAML")
	}

	output := string(out)

	// Prepend managed-layer warning if not already present
	headerComment := "# This file was patched by reconciler.\n" +
		"# The '30-cre' configuration layers are managed — do not edit them manually.\n"
	if !strings.HasPrefix(output, "# This file was patched") {
		output = headerComment + output
	}

	if err := os.WriteFile(yamlPath, []byte(output), 0600); err != nil {
		return errors.Wrapf(err, "failed to write patched YAML %s", yamlPath)
	}

	return nil
}

// FindKey finds a value node by key name in a mapping node.
func FindKey(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// findOrCreateKey finds a key in a mapping node, or creates it with the
// given kind if it doesn't exist. Returns the value node.
func findOrCreateKey(mapping *yaml.Node, key string, kind yaml.Kind) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}

	// Create new key-value pair
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
	valNode := &yaml.Node{Kind: kind}
	mapping.Content = append(mapping.Content, keyNode, valNode)
	return valNode
}

// updateOrCreateLayer updates an existing layer by name, or appends a new one.
// configuration is a sequence node containing mapping nodes with a single
// key (layer name) → value (TOML string).
func updateOrCreateLayer(configList *yaml.Node, layerName, tomlContent string) {
	if configList.Kind != yaml.SequenceNode {
		return
	}

	// Each item in the sequence is a mapping with one key (the layer name)
	for _, item := range configList.Content {
		if item.Kind == yaml.MappingNode && len(item.Content) >= 2 && item.Content[0].Value == layerName {
			// Update existing layer
			item.Content[1].Value = tomlContent
			item.Content[1].Style = yaml.LiteralStyle
			return
		}
	}

	// Append new layer
	layerMap := &yaml.Node{Kind: yaml.MappingNode}
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: layerName, Tag: "!!str"}
	valNode := &yaml.Node{Kind: yaml.ScalarNode, Value: tomlContent, Tag: "!!str", Style: yaml.LiteralStyle}
	layerMap.Content = append(layerMap.Content, keyNode, valNode)
	configList.Content = append(configList.Content, layerMap)
}
