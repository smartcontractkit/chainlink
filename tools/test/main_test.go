package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncSkillsFromBaseFailsBeforeRemovingDestinationWhenSourceMissing(t *testing.T) {
	base := t.TempDir()
	dst := filepath.Join(base, skillsDst)
	require.NoError(t, os.MkdirAll(dst, 0o755))
	sentinel := filepath.Join(dst, "keep.md")
	require.NoError(t, os.WriteFile(sentinel, []byte("keep"), 0o600))

	err := syncSkillsFromBase(base)

	require.Error(t, err)
	b, readErr := os.ReadFile(sentinel)
	require.NoError(t, readErr)
	assert.Equal(t, "keep", string(b))
}

func TestSyncSkillsFromBasePreservesFileMode(t *testing.T) {
	base := t.TempDir()
	src := filepath.Join(base, skillsSrc, "example")
	require.NoError(t, os.MkdirAll(src, 0o755))
	sourceFile := filepath.Join(src, "run.sh")
	require.NoError(t, os.WriteFile(sourceFile, []byte("#!/bin/sh\n"), 0o600))
	require.NoError(t, os.Chmod(sourceFile, 0o750))

	require.NoError(t, syncSkillsFromBase(base))

	info, err := os.Stat(filepath.Join(base, skillsDst, "example", "run.sh"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o750), info.Mode().Perm())
}
