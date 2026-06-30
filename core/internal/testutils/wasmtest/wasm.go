package wasmtest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var getRepoRoot = sync.OnceValues(func() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get caller info")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repo root (no go.mod found)")
		}
		dir = parent
	}
})

// GetTestBinary looks up the WASM binary from outputPath and optionally brotli-compresses it.
func GetTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()

	repoRoot, err := getRepoRoot()
	require.NoError(tb, err, "failed to get repo root: %s", err)

	pkgDir := filepath.Join(repoRoot, outputPath)

	hash, err := HashPackage(pkgDir)
	require.NoError(tb, err, "failed to hash package %s: %s", pkgDir, err)

	// Determine output filename
	cacheFile := fmt.Sprintf("output-%s.wasm", hash)
	if compress {
		cacheFile += ".br"
	}
	filePath := filepath.Join(pkgDir, "testdata", cacheFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		require.NoError(tb, err, "WASM fixture missing or out of date. Please run 'go generate ./...' in the project root to rebuild it. Missing file: %s", filePath)
	}

	// Read from cache
	binary, err := os.ReadFile(filePath)
	require.NoError(tb, err, "read cache file failed: %s", err)

	return binary
}

// HashPackage computes a SHA-256 hash of all .go files within the specified package directory.
// It is used to uniquely identify the state of the package's source code for WASM compilation caching.
func HashPackage(pkgDir string) (string, error) {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return "", err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	h := sha256.New()
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(pkgDir, f))
		if err != nil {
			return "", err
		}
		h.Write([]byte(f))
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}
