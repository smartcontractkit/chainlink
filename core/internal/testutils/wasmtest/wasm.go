package wasmtest

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func CreateTestBinary(outputPath string, compress bool, t *testing.T) []byte {
	filePath := filepath.Join(t.TempDir(), uuid.New().String()+".wasm")
	cmd := exec.Command("go", "build", "-o", filePath, "github.com/smartcontractkit/chainlink/v2/"+outputPath) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	binary, err := os.ReadFile(filePath)
	require.NoError(t, err)

	if !compress {
		return binary
	}

	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	_, err = bwr.Write(binary)
	require.NoError(t, err)
	require.NoError(t, bwr.Close())

	cb, err := io.ReadAll(&b)
	require.NoError(t, err)
	return cb
}
