package wasmtest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"

	"github.com/stretchr/testify/require"
)

var (
	binaryCache   = make(map[string][]byte)
	binaryCacheMu sync.RWMutex
	binaryGroup   singleflight.Group
)

func CreateTestBinary(outputPath string, compress bool, t *testing.T) []byte {
	cacheKey := fmt.Sprintf("%s-%t", outputPath, compress)

	binaryCacheMu.RLock()
	cached, ok := binaryCache[cacheKey]
	binaryCacheMu.RUnlock()
	if ok {
		return cached
	}

	v, err, _ := binaryGroup.Do(cacheKey, func() (interface{}, error) {
		// Recheck cache just in case
		binaryCacheMu.RLock()
		if cached, ok := binaryCache[cacheKey]; ok {
			binaryCacheMu.RUnlock()
			return cached, nil
		}
		binaryCacheMu.RUnlock()

		filePath := filepath.Join(t.TempDir(), uuid.New().String()+".wasm")
		cmd := exec.Command("go", "build", "-o", filePath, "github.com/smartcontractkit/chainlink/v2/"+outputPath) // #nosec
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

		binaryCacheMu.Lock()
		binaryCache[cacheKey] = binary
		binaryCacheMu.Unlock()

		return binary, nil
	})

	require.NoError(t, err)
	return v.([]byte)
}
