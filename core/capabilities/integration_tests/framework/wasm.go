package framework

import (
	"bytes"
	"encoding/base64"
	"os"
	"os/exec"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"
)

// GetCompressedWorkflowWasm compresses the workflow wasm file using brotli and returns the compressed binary and the base64 encoded compressed binary
func GetCompressedWorkflowWasm(t *testing.T, workflowWasmPath string) ([]byte, string) {
	uncompressedBinary, err := os.ReadFile(workflowWasmPath)
	require.NoError(t, err)

	var compressedBuffer bytes.Buffer
	writer := brotli.NewWriter(&compressedBuffer)
	_, err = writer.Write(uncompressedBinary)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)
	compressedBinary := compressedBuffer.Bytes()
	base64EncodedCompressedBinary := base64.StdEncoding.EncodeToString(compressedBinary)
	return compressedBinary, base64EncodedCompressedBinary
}

func CreateWasmBinary(t *testing.T, goFile string, wasmFile string) {
	cmd := exec.Command("go", "build", "-o", wasmFile, goFile) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}
