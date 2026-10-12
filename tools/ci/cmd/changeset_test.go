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
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/changeset"
)

func TestChangesetCheckTags_CLI_JSON(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := `---
"chainlink": patch
---

Fixing issue #bugfix
`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"changeset", "check-tags", filePath, "--json"})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	var res changeset.Result
	err = json.Unmarshal(out.Bytes(), &res)
	require.NoError(t, err)
	assert.True(t, res.HasTags)
	assert.Equal(t, []string{"#bugfix"}, res.FoundTags)
}

func TestChangesetCheckTags_CLI_MultipleFiles(t *testing.T) {
	t.Setenv("CHANGESET_FILE_PATH", "a.md b.md")

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"changeset", "check-tags"})

	err := rootCmd.ExecuteContext(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiple changeset file paths")
}

func TestChangesetCheckTags_CLI_Human_NoOutputPollution(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := `---
"chainlink": patch
---

Fixing issue #bugfix
`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"changeset", "check-tags", filePath})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Found tag: #bugfix in "+filePath)
	assert.NotContains(t, out.String(), "has_tags=")
	assert.NotContains(t, out.String(), "found_tags=")
}

func TestChangesetCheckTags_CLI_Human_Env(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := `---
"chainlink": patch
---

Fixing issue #bugfix
`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	ghOutput := filepath.Join(tmpDir, "gh_output")
	require.NoError(t, os.WriteFile(ghOutput, []byte{}, 0o600))
	t.Setenv("GITHUB_OUTPUT", ghOutput)

	rootCmd := cmd.NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"changeset", "check-tags", filePath})

	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Found tag: #bugfix in "+filePath)

	ghContent, err := os.ReadFile(ghOutput)
	require.NoError(t, err)
	assert.Contains(t, string(ghContent), "has_tags<<_GitHubActionsFileCommandDelimeter_\ntrue\n_GitHubActionsFileCommandDelimeter_\n")      //typos:ignore `Delimeter` should be `Delimiter` // From GitHub
	assert.Contains(t, string(ghContent), "found_tags<<_GitHubActionsFileCommandDelimeter_\n#bugfix\n_GitHubActionsFileCommandDelimeter_\n") //typos:ignore `Delimeter` should be `Delimiter` // From GitHub
}
