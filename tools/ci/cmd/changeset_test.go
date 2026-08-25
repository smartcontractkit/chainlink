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
	assert.Contains(t, string(ghContent), "has_tags=true\n")
	assert.Contains(t, string(ghContent), "found_tags=#bugfix\n")
}

func TestChangesetCheckTags_ParityWithBash(t *testing.T) {
	repoRoot, err := cmd.FindRepoRoot(context.Background())
	require.NoError(t, err)

	bashScript := filepath.Join(repoRoot, ".github", "scripts", "check-changeset-tags.sh")
	if _, err := os.Stat(bashScript); os.IsNotExist(err) {
		t.Skip("bash script does not exist, skipping parity test")
	}

	tmpDir := t.TempDir()
	testCases := []struct {
		filename string
		content  string
	}{
		{
			filename: "single_tag_patch.md",
			content: `---
"chainlink": patch
---

Added bugfix description #bugfix
`,
		},
		{
			filename: "multiple_tags_minor.md",
			content: `---
"chainlink": minor
---

New feature #added with database update #db_update and #wip
`,
		},
		{
			filename: "no_tags_major.md",
			content: `---
"chainlink": major
---

No tags in this description
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, tc.filename)
			require.NoError(t, os.WriteFile(filePath, []byte(tc.content), 0o600))

			// Run bash script
			bashOutputFile := filepath.Join(tmpDir, tc.filename+"_bash_out")
			require.NoError(t, os.WriteFile(bashOutputFile, []byte{}, 0o600))

			bashCmd := exec.CommandContext(t.Context(), "bash", bashScript, filePath)
			bashCmd.Env = append(os.Environ(), "GITHUB_OUTPUT="+bashOutputFile)
			bashOut, bashErr := bashCmd.CombinedOutput()
			require.NoError(t, bashErr, "bash script failed: %s", string(bashOut))

			bashOutContent, err := os.ReadFile(bashOutputFile)
			require.NoError(t, err)

			// Run Go CLI
			goOutputFile := filepath.Join(tmpDir, tc.filename+"_go_out")
			require.NoError(t, os.WriteFile(goOutputFile, []byte{}, 0o600))

			rootCmd := cmd.NewRootCmd()
			var goOut bytes.Buffer
			rootCmd.SetOut(&goOut)
			rootCmd.SetErr(&goOut)
			rootCmd.SetArgs([]string{"changeset", "check-tags", filePath})

			t.Setenv("GITHUB_OUTPUT", goOutputFile)

			err = rootCmd.ExecuteContext(context.Background())
			require.NoError(t, err)

			goOutContent, err := os.ReadFile(goOutputFile)
			require.NoError(t, err)

			assert.Equal(t, strings.TrimSpace(string(bashOutContent)), strings.TrimSpace(string(goOutContent)))
		})
	}
}
