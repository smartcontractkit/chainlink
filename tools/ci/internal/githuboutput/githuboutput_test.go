package githuboutput

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendVar_WritesKeyValuePair(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	require.NoError(t, AppendVar("matrix", `[{"test_name":"Test1"}]`))
	require.NoError(t, AppendVar("version", "2.21.0"))

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "matrix=[{\"test_name\":\"Test1\"}]\nversion=2.21.0\n", string(content))
}

func TestAppendMultilineVar_UsesDelimitedSyntax(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	require.NoError(t, AppendMultilineVar("pr_body", "line one\nline two"))

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "pr_body<<EOF\nline one\nline two\nEOF\n", string(content))
}

func TestAppendVars_NoOpWhenEnvUnset(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")

	require.NoError(t, AppendVars(map[string]string{"a": "b"}))
}

func TestAppendToFile_CreatesFileAndAppends(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "summary.md")

	require.NoError(t, AppendToFile(path, "hello"))
	require.NoError(t, AppendToFile(path, " world"))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(content))
}
