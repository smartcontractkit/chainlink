package wasmtest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPackage(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Create a dummy .go file
	goFilePath := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(goFilePath, []byte("package main"), 0600)
	require.NoError(t, err)

	hash1, err := HashPackage(tmpDir)
	require.NoError(t, err)
	assert.NotEmpty(t, hash1)

	// Change file content
	err = os.WriteFile(goFilePath, []byte("package main\n\nfunc main() {}"), 0600)
	require.NoError(t, err)

	hash2, err := HashPackage(tmpDir)
	require.NoError(t, err)
	assert.NotEmpty(t, hash2)
	assert.NotEqual(t, hash1, hash2, "hash should change when file content changes")

	// Adding another file changes hash
	goFilePath2 := filepath.Join(tmpDir, "helper.go")
	err = os.WriteFile(goFilePath2, []byte("package main"), 0600)
	require.NoError(t, err)

	hash3, err := HashPackage(tmpDir)
	require.NoError(t, err)
	assert.NotEqual(t, hash2, hash3, "hash should change when new file is added")
}
