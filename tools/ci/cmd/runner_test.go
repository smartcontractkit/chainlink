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
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/runner"
)

func TestRunnerSpot_CLI_JSON(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"runner", "spot", "--event", "merge_group", "--json"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res runner.SpotResult
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "false", res.Spot)
	assert.Equal(t, "spot=false", res.SpotFlag)
	assert.False(t, res.Enabled)
	assert.True(t, res.IsMergeQueue)
}

func TestRunnerSpot_CLI_Human(t *testing.T) {
	t.Parallel()

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"runner", "spot", "--event", "push", "--ref", "refs/heads/develop"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "spot=co\n", out.String())
}

func TestRunnerSpot_CLI_GithubOutput(t *testing.T) {
	tmpDir := t.TempDir()
	ghOutput := filepath.Join(tmpDir, "gh_output")
	require.NoError(t, os.WriteFile(ghOutput, []byte{}, 0o600))
	t.Setenv("GITHUB_OUTPUT", ghOutput)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"runner", "spot", "--event", "pull_request", "--base-ref", "release/2.57.1"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	ghContent, err := os.ReadFile(ghOutput)
	require.NoError(t, err)
	assert.Contains(t, string(ghContent), "spot")
	assert.Contains(t, string(ghContent), "spot_flag")
	assert.Contains(t, string(ghContent), "is_release")
}

func TestRunnerSpot_CLI_EnvFallback(t *testing.T) {
	t.Setenv("GITHUB_EVENT_NAME", "merge_group")
	t.Setenv("GITHUB_REF", "refs/heads/gh-readonly-queue/develop/pr-123")

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"runner", "spot"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "spot=false\n", out.String())
}
