package infra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchChartValues_UpdateExistingLayer(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "dev.yaml")

	original := `chainlink-node:
  instances:
    node-0:
      configuration:
        - 01-network: |
            [[EVM]]
            ChainID = '1337'
`
	require.NoError(t, os.WriteFile(yamlPath, []byte(original), 0600))

	patches := map[string]string{
		"node-0": "[Capabilities.ExternalRegistry]\n  Address = '0xabc'\n",
	}
	require.NoError(t, PatchChartValues(yamlPath, patches))

	data, _ := os.ReadFile(yamlPath)
	content := string(data)
	require.Contains(t, content, "30-cre")
	require.Contains(t, content, "0xabc")
	require.Contains(t, content, "01-network") // original layer preserved
}

func TestPatchChartValues_Idempotent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "dev.yaml")

	original := `chainlink-node:
  instances:
    node-0:
      configuration:
        - 01-network: |
            [[EVM]]
            ChainID = '1337'
`
	require.NoError(t, os.WriteFile(yamlPath, []byte(original), 0600))

	patch := map[string]string{
		"node-0": "[Capabilities.ExternalRegistry]\n  Address = '0xabc'\n",
	}

	require.NoError(t, PatchChartValues(yamlPath, patch))
	require.NoError(t, PatchChartValues(yamlPath, patch)) // run again

	data, _ := os.ReadFile(yamlPath)
	content := string(data)
	// Should only have one 30-cre layer (count the YAML key, not the header comment)
	count := strings.Count(content, "- 30-cre:")
	require.Equal(t, 1, count, "expected exactly one 30-cre layer after double patch")
}

func TestReadLayerValue_MissingFile(t *testing.T) {
	t.Parallel()

	value, exists, err := ReadLayerValue(filepath.Join(t.TempDir(), "nope.yaml"), "node-0", "30-cre")
	require.NoError(t, err)
	require.False(t, exists)
	require.Empty(t, value)
}

func TestReadLayerValue_LayerAbsent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "dev.yaml")
	original := `chainlink-node:
  instances:
    node-0:
      configuration:
        - 01-network: |
            [[EVM]]
            ChainID = '1337'
`
	require.NoError(t, os.WriteFile(yamlPath, []byte(original), 0600))

	value, exists, err := ReadLayerValue(yamlPath, "node-0", "30-cre")
	require.NoError(t, err)
	require.False(t, exists)
	require.Empty(t, value)
}

func TestReadLayerValue_ReflectsPatchedContent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "dev.yaml")
	original := `chainlink-node:
  instances:
    node-0:
      configuration:
        - 01-network: |
            [[EVM]]
            ChainID = '1337'
`
	require.NoError(t, os.WriteFile(yamlPath, []byte(original), 0600))

	content := "[Capabilities.ExternalRegistry]\n  Address = '0xabc'\n"
	require.NoError(t, PatchChartValues(yamlPath, map[string]string{"node-0": content}))

	value, exists, err := ReadLayerValue(yamlPath, "node-0", "30-cre")
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, content, value)
}
