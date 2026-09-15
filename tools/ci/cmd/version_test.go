package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/cmd"
)

func TestRootCmd_Help(t *testing.T) {
	t.Parallel()
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out.String(), "ci provides tooling for CI workflows")
}

func TestVersionCmd_Text(t *testing.T) {
	t.Parallel()
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out.String(), "ci version")
}

func TestVersionCmd_JSON(t *testing.T) {
	t.Parallel()
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"version", "--json"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var payload map[string]string
	err = json.Unmarshal(out.Bytes(), &payload)
	require.NoError(t, err)
	assert.NotEmpty(t, payload["version"])
}

func TestFindRepoRoot(t *testing.T) {
	t.Parallel()
	repoRoot, err := cmd.FindRepoRoot(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, repoRoot)
}
