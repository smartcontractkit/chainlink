package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "package.json")
	require.NoError(t, os.WriteFile(pkgPath, []byte(`{"name": "chainlink", "version": "2.21.0"}`), 0600))

	ver, err := ReadVersionFromPackageJSON(pkgPath)
	require.NoError(t, err)
	require.Equal(t, "2.21.0", ver)
}

func TestFormatChangelog(t *testing.T) {
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "package.json")
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	outputFile := filepath.Join(tempDir, "github_output")

	require.NoError(t, os.WriteFile(pkgPath, []byte(`{"version": "2.21.0"}`), 0600))

	sampleChangelog := `# Changelog Chainlink Core

## 2.21.0

### Minor Changes

- [#added] [#nops] Add new OCR3 consensus feature
  extra details on consensus
- [#bugfix] Fix token expiration check
- Untagged change item

## 2.20.0

- Old version change
`
	require.NoError(t, os.WriteFile(changelogPath, []byte(sampleChangelog), 0600))

	t.Setenv("GITHUB_OUTPUT", outputFile)

	res, err := Format(changelogPath, pkgPath, true)
	require.NoError(t, err)

	require.Equal(t, "2.21.0", res.Version)
	require.Contains(t, res.NewChangelog, "## 2.21.0 - PREVIEW")
	require.Contains(t, res.NewChangelog, "## added")
	require.Contains(t, res.NewChangelog, "## bugfix")
	require.Contains(t, res.NewChangelog, "## untagged")
	require.Contains(t, res.NewChangelog, "## 2.20.0")

	// PR body check
	require.Contains(t, res.PRBody, "This PR is a preview of the changes")
	require.Contains(t, res.PRBody, "## added")

	// Verify CHANGELOG.md was updated
	diskContent, err := os.ReadFile(changelogPath)
	require.NoError(t, err)
	require.Equal(t, res.NewChangelog, string(diskContent))

	// Verify GITHUB_OUTPUT was written
	ghOutput, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Contains(t, string(ghOutput), "version=2.21.0")
	require.Contains(t, string(ghOutput), "pr_body<<EOF")
}

func TestFormatChangelog_Truncation(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	pkgPath := filepath.Join(tempDir, "package.json")
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")

	require.NoError(t, os.WriteFile(pkgPath, []byte(`{"version": "2.21.0"}`), 0600))

	// Create a huge changelog entry exceeding 65,000 characters
	hugeChange := "- [#added] " + strings.Repeat("A", 70000)
	sampleChangelog := "# Changelog Chainlink Core\n\n## 2.21.0\n\n" + hugeChange + "\n\n## 2.20.0\n"
	require.NoError(t, os.WriteFile(changelogPath, []byte(sampleChangelog), 0600))

	res, err := Format(changelogPath, pkgPath, false)
	require.NoError(t, err)

	require.Contains(t, res.PRBody, "The changelog content is too long for the PR description.")
}
