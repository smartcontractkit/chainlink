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

func TestEOFCmd_Help(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"end-of-file-fixer", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Ensure files end with a single newline")
	assert.Contains(t, buf.String(), "--check")
}

func TestEOFCmd_RunFile(t *testing.T) {
	t.Parallel()

	t.Run("fixes file without newline", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "test.go")
		require.NoError(t, os.WriteFile(target, []byte("package main"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"end-of-file-fixer", target})

		err := root.Execute()
		require.NoError(t, err)

		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "package main\n", string(content))
	})

	t.Run("fails in check mode when file needs fixing", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "test.go")
		require.NoError(t, os.WriteFile(target, []byte("package main"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"end-of-file-fixer", "--check", target})

		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "require end-of-file newline fixes")
	})
}
