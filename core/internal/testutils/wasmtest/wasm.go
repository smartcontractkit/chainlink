package wasmtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/andybalholm/brotli"

	"github.com/stretchr/testify/require"
)

var (
	binaryOnce   = make(map[string]func() ([]byte, error))
	binaryOnceMu sync.Mutex
)

// CreateTestBinary compiles a WASM binary from outputPath and optionally brotli-compresses it.
// Results are cached for the lifetime of the test process keyed by (outputPath, compress)
// via sync.OnceValues — the build runs exactly once per key regardless of concurrency.
// This assumes source at outputPath does not change during a single `go test` invocation.
func CreateTestBinary(tb testing.TB, outputPath string, compress bool) []byte {
	tb.Helper()
	cacheKey := fmt.Sprintf("%s-%t", outputPath, compress)

	binaryOnceMu.Lock()
	once, ok := binaryOnce[cacheKey]
	if !ok {
		once = sync.OnceValues(func() ([]byte, error) {
			tmpDir, err := os.MkdirTemp("", "wasmtest-*")
			if err != nil {
				return nil, fmt.Errorf("create temp dir: %w", err)
			}
			defer os.RemoveAll(tmpDir)

			cmdCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			filePath := filepath.Join(tmpDir, "output.wasm")
			cmd := exec.CommandContext(cmdCtx, "go", "build", "-o", filePath, "github.com/smartcontractkit/chainlink/v2/"+outputPath) // #nosec
			cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

			output, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("build failed: %s %w", string(output), err)
			}

			binary, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("read file failed: %w", err)
			}

			if compress {
				var b bytes.Buffer
				bwr := brotli.NewWriter(&b)
				if _, err = bwr.Write(binary); err != nil {
					return nil, err
				}
				if err = bwr.Close(); err != nil {
					return nil, err
				}

				cb, err := io.ReadAll(&b)
				if err != nil {
					return nil, err
				}
				binary = cb
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
