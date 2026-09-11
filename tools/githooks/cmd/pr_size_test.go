package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

//nolint:paralleltest // t.Chdir is process-global
func TestPRSizeCmd_Run(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "base.go"), []byte("package main\n"), 0o600))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "init")
	baseSHA := git(t, dir, "rev-parse", "HEAD")
	git(t, dir, "update-ref", "refs/remotes/origin/develop", baseSHA)
	git(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	require.NoError(t, os.WriteFile(filepath.Join(dir, "feature.go"), []byte("package main\n// new\n"), 0o600))
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-m", "feature")

	t.Chdir(dir)

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"pr-size", "--no-color"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "PR Diff Size")
	assert.Contains(t, buf.String(), "Classification: [ SMALL ] [ OK ]")
}
