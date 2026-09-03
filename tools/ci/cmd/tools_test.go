package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/cmd"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/tools"
)

func TestToolsMatrixCmd_JSON(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"tools", "matrix", "--all", "--json"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var targets []tools.Target
	err = json.Unmarshal(out.Bytes(), &targets)
	require.NoError(t, err)
	assert.NotEmpty(t, targets)

	// Verify all known targets exist
	names := make([]string, len(targets))
	for i, tgt := range targets {
		names[i] = tgt.Name
	}
	assert.Contains(t, names, "tools/ci")
	assert.Contains(t, names, "tools/githooks")
	assert.Contains(t, names, "tools/test")
	assert.Contains(t, names, "tools/root")
}

func TestToolsMatrixCmd_GHAOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "github_output")
	require.NoError(t, os.WriteFile(outputPath, []byte{}, 0o600))

	t.Setenv("GITHUB_OUTPUT", outputPath)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"tools", "matrix", "--all"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "matrix")
	assert.Contains(t, string(content), "tools/ci")
}

func TestToolsMatrixCmd_FilterChangedFiles(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"tools", "matrix", "--changed-files", "tools/ci/cmd/version.go", "--json"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var targets []tools.Target
	err = json.Unmarshal(out.Bytes(), &targets)
	require.NoError(t, err)
	require.Len(t, targets, 1)
	assert.Equal(t, "tools/ci", targets[0].Name)
}
