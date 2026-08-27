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

func TestWhitespaceCmd_Help(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"whitespace-fixer", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Fix erroneous trailing whitespace in eligible code and text files")
	assert.Contains(t, buf.String(), "--check")
}

func TestWhitespaceCmd_RunFile(t *testing.T) {
	t.Parallel()

	t.Run("fixes trailing whitespace in file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "test.yaml")
		require.NoError(t, os.WriteFile(target, []byte("key: value   \n"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"whitespace-fixer", target})

		err := root.Execute()
		require.NoError(t, err)

		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "key: value\n", string(content))
	})

	t.Run("skips Go files", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "test.go")
		require.NoError(t, os.WriteFile(target, []byte("package main   \n"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"whitespace-fixer", target})

		err := root.Execute()
		require.NoError(t, err)

		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "package main   \n", string(content))
	})

	t.Run("fails in check mode when file has trailing whitespace", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "test.yaml")
		require.NoError(t, os.WriteFile(target, []byte("key: value   \n"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"whitespace-fixer", "--check", target})

		err := root.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "require whitespace fixes")
	})

	t.Run("preserves markdown two-space linebreaks by default", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "README.md")
		require.NoError(t, os.WriteFile(target, []byte("Line with break.  \nNext line.\n"), 0o600))

		root := cmd.NewRootCmd()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetErr(buf)
		root.SetArgs([]string{"whitespace-fixer", target})

		err := root.Execute()
		require.NoError(t, err)

		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "Line with break.  \nNext line.\n", string(content))
	})
}
