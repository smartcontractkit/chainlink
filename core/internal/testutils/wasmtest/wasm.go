package wasmtest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"
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

var (
	binaryCache   = make(map[string][]byte)
	binaryCacheMu sync.RWMutex
)

// GetTestBinary fetches the pre-compiled WASM binary from the package's testdata/ directory.
// The binary MUST be generated with `go generate` before running tests.
// Example:
// //go:generate go run path/to/internal/testutils/wasmtest/generator/main.go -pkg core/target/package
func GetTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()

	cacheKey := fmt.Sprintf("%s:%t", outputPath, compress)
	binaryCacheMu.RLock()
	cached, ok := binaryCache[cacheKey]
	binaryCacheMu.RUnlock()
	if ok {
		res := make([]byte, len(cached))
		copy(res, cached)
		return res
	}

	binaryCacheMu.Lock()
	defer binaryCacheMu.Unlock()
	if cached, ok = binaryCache[cacheKey]; ok {
		res := make([]byte, len(cached))
		copy(res, cached)
		return res
	}

	repoRoot, err := getRepoRoot()
	require.NoError(tb, err, "failed to get repo root: %s", err)

	binary, err := readFixture(filepath.Join(repoRoot, outputPath), compress)
	require.NoError(tb, err)

	cachedCopy := make([]byte, len(binary))
	copy(cachedCopy, binary)
	binaryCache[cacheKey] = cachedCopy

	return binary
}

// readFixture reads the pre-compiled WASM fixture from pkgDir/testdata. A missing fixture is
// reported as an actionable error directing the caller to regenerate it.
func readFixture(pkgDir string, compress bool) ([]byte, error) {
	filePath := filepath.Join(pkgDir, "testdata", "output.wasm.br")
	binary, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("WASM fixture missing or out of date. Run 'go generate ./...' to rebuild it. Missing file: %s", filePath)
		}
		return nil, fmt.Errorf("read fixture failed: %w", err)
	}

	if !compress {
		var b bytes.Buffer
		bwr := brotli.NewReader(bytes.NewReader(binary))
		if _, err := io.Copy(&b, bwr); err != nil {
			return nil, fmt.Errorf("decompress fixture failed: %w", err)
		}
		return b.Bytes(), nil
	}

	return binary, nil
}
