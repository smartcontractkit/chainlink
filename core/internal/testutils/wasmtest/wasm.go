package wasmtest

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/wasmbuild"
)

var getRepoRoot = sync.OnceValues(func() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("could not get caller info")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find repo root (no go.mod found)")
		}
		dir = parent
	}
})

// GetTestBinary returns a WASM binary for the given package path relative to the repo root
// (e.g. core/capabilities/compute/test/simple/cmd). The binary is built lazily on first use
// and cached in .wasm-cache/ using a content-addressed filename.
func GetTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()

	repoRoot, err := getRepoRoot()
	require.NoError(tb, err, "failed to get repo root: %s", err)

	pkgDir := outputPath
	if !filepath.IsAbs(pkgDir) {
		pkgDir = filepath.Join(repoRoot, outputPath)
	}

	binary, err := wasmbuild.Compile(context.Background(), wasmbuild.Config{
		PkgDir:   pkgDir,
		RepoRoot: repoRoot,
		Compress: compress,
	})
	require.NoError(tb, err)

	return binary
}
