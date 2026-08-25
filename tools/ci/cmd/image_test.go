package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
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
	assert.Contains(t, string(ghContent), "resolved_image=123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0")
}

func TestImageResolve_ParityWithBash(t *testing.T) {
	repoRoot, err := cmd.FindRepoRoot(context.Background())
	require.NoError(t, err)

	bashScript := filepath.Join(repoRoot, ".github", "scripts", "resolve-chainlink-image.sh")
	if _, err := os.Stat(bashScript); os.IsNotExist(err) {
		t.Skip("bash script does not exist, skipping parity test")
	}

	testCases := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "public ecr valid",
			env: map[string]string{
				"ECR_TYPE":                  "public",
				"CHAINLINK_IMAGE_REPO_PATH": "chainlink",
				"CHAINLINK_IMAGE_TAG":       "v2.1.0",
			},
		},
		{
			name: "sdlc ecr valid",
			env: map[string]string{
				"ECR_TYPE":                  "sdlc",
				"CHAINLINK_IMAGE_REPO_PATH": "chainlink-integration-tests",
				"CHAINLINK_IMAGE_TAG":       "v2.1.0",
				"AWS_ACCOUNT_NUMBER":        "123456789012",
				"AWS_REGION":                "us-west-2",
			},
		},
		{
			name: "case insensitive normalization with mixed tag",
			env: map[string]string{
				"ECR_TYPE":                  "PuBLic",
				"CHAINLINK_IMAGE_REPO_PATH": "ChAinLink",
				"CHAINLINK_IMAGE_TAG":       "V2.1.0-custom",
			},
		},
		{
			name: "whitespace padded inputs",
			env: map[string]string{
				"ECR_TYPE":                  "  SDLC  ",
				"CHAINLINK_IMAGE_REPO_PATH": "  chainlink-integration-tests  ",
				"CHAINLINK_IMAGE_TAG":       "  v2.1.0  ",
				"AWS_ACCOUNT_NUMBER":        "  123456789012  ",
				"AWS_REGION":                "  US-WEST-2  ",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run bash script
			bashCmd := exec.CommandContext(t.Context(), "bash", bashScript)
			bashCmd.Env = os.Environ()
			for k, v := range tc.env {
				bashCmd.Env = append(bashCmd.Env, k+"="+v)
			}
			bashOut, bashErr := bashCmd.CombinedOutput()
			require.NoError(t, bashErr, "bash script failed: %s", string(bashOut))

			// Run Go CLI
			rootCmd := cmd.NewRootCmd()
			var goOut bytes.Buffer
			rootCmd.SetOut(&goOut)
			rootCmd.SetErr(&goOut)
			rootCmd.SetArgs([]string{"image", "resolve"})

			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			err := rootCmd.ExecuteContext(context.Background())
			require.NoError(t, err)

			assert.Equal(t, strings.TrimSpace(string(bashOut)), strings.TrimSpace(goOut.String()))
		})
	}
}
