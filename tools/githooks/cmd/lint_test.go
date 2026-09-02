package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

func TestRootCmd(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "githooks provides tooling for Git hooks")
	assert.Contains(t, buf.String(), "lint")
	assert.Contains(t, buf.String(), "test")
	assert.Contains(t, buf.String(), "tidy")
	assert.Contains(t, buf.String(), "generate")
	assert.Contains(t, buf.String(), "end-of-file-fixer")
	assert.Contains(t, buf.String(), "whitespace-fixer")
}

func TestLintCmdHelp(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"lint", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Discovers enclosing Go modules and packages for files changed since the merge-base with the default branch")
	assert.Contains(t, buf.String(), "--fix")
	assert.Contains(t, buf.String(), "--rev")
}

func TestTestCmdHelp(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"test", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Discovers affected Go test packages for changed/staged files")
	assert.Contains(t, buf.String(), "--short")
}

func TestTidyCmdHelp(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"tidy", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Run go mod tidy in parallel on all changed Go modules")
}

func TestGenerateCmdHelp(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"generate", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Run targeted code generators (proto, config docs) on changed files")
}
