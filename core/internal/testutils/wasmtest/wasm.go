package wasmtest

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"
)

func CreateTestBinary(outputPath, path string, compress bool, t *testing.T) []byte {
	cmd := exec.Command("go", "build", "-o", path, fmt.Sprintf("github.com/smartcontractkit/chainlink/v2/%s", outputPath)) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	binary, err := os.ReadFile(path)
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
