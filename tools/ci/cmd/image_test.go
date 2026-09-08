package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/cmd"
)

func TestImageResolve_CLI_Public(t *testing.T) {
	t.Parallel()
	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"image", "resolve",
		"--ecr-type", "public",
		"--repo-path", "chainlink",
		"--tag", "v2.1.0",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "public.ecr.aws/chainlink:v2.1.0", strings.TrimSpace(out.String()))
}

func TestImageResolve_CLI_SDLC_Env(t *testing.T) {
	tmpDir := t.TempDir()
	ghOutput := filepath.Join(tmpDir, "gh_output")
	require.NoError(t, os.WriteFile(ghOutput, []byte{}, 0o600))
	t.Setenv("GITHUB_OUTPUT", ghOutput)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{
		"image", "resolve",
		"--ecr-type", "sdlc",
		"--repo-path", "chainlink-integration-tests",
		"--tag", "v2.1.0",
		"--aws-account", "123456789012",
		"--aws-region", "us-west-2",
		"--json",
	})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var payload map[string]string
	err = json.Unmarshal(out.Bytes(), &payload)
	require.NoError(t, err)
	assert.Equal(t, "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0", payload["resolved_image"])

	ghContent, err := os.ReadFile(ghOutput)
	require.NoError(t, err)
	assert.Contains(t, string(ghContent), "resolved_image<<_GitHubActionsFileCommandDelimeter_\n123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0\n_GitHubActionsFileCommandDelimeter_\n") //typos:ignore `Delimeter` should be `Delimiter` // From GitHub
}
