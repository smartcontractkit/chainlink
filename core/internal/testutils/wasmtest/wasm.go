package wasmtest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	binaryOnce   = make(map[string]func() ([]byte, error))
	binaryOnceMu sync.Mutex
)

// CreateTestBinary looks up the WASM binary from outputPath and optionally brotli-compresses it.
// Results are cached across test runs by saving the binary to the package directory with a hash.
// Within a single test process, results are further cached via sync.OnceValues.
func CreateTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()
	cacheKey := fmt.Sprintf("%s-%t", outputPath, compress)

	binaryOnceMu.Lock()
	once, ok := binaryOnce[cacheKey]
	if !ok {
		once = sync.OnceValues(func() ([]byte, error) {
			pkgPath := "github.com/smartcontractkit/chainlink/v2/" + outputPath

			// Find the absolute directory of the package to store the cached binary
			listCmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgPath)
			listCmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
			out, err := listCmd.Output()
			if err != nil {
				return nil, fmt.Errorf("failed to find package dir for %s: %w\n%s", pkgPath, err, string(out))
			}
			pkgDir := strings.TrimSpace(string(out))

			hash, err := HashPackage(pkgDir)
			if err != nil {
				return nil, fmt.Errorf("failed to hash package: %w", err)
			}

			// Determine output filename
			cacheFile := fmt.Sprintf("output-%s.wasm", hash)
			if compress {
				cacheFile += ".br"
			}
			testdataDir := filepath.Join(pkgDir, "testdata")
			filePath := filepath.Join(testdataDir, cacheFile)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("WASM fixture missing or out of date. Run 'go generate ...' to rebuild it. Missing file: %s", filePath)
			}

			// Read from cache
			binary, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("read cache file failed: %w", err)
			}
			return binary, nil
		})
		binaryOnce[cacheKey] = once
	}
	binaryOnceMu.Unlock()

	result, err := once()
	require.NoError(tb, err)
	return result
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
