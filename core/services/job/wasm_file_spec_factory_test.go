package job_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestWasmFileSpecFactory(t *testing.T) {
	binaryLocation := createTestBinary(t)
	configLocation := "testdata/config.json"
	config, err := os.ReadFile(configLocation)
	require.NoError(t, err)

	rawBinary, err := os.ReadFile(binaryLocation)
	require.NoError(t, err)

	b := bytes.Buffer{}
	bwr := brotli.NewWriter(&b)
	_, err = bwr.Write(rawBinary)
	require.NoError(t, err)

	require.NoError(t, bwr.Close())

	t.Run("Raw binary", func(t *testing.T) {
		ctx := testutils.Context(t)
		factory := job.WasmFileSpecFactory{}
		actual, rawSpec, actualSha, err2 := factory.Spec(testutils.Context(t), binaryLocation, configLocation)
		require.NoError(t, err2)

		expected, err2 := host.GetWorkflowSpec(ctx, &host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}, rawBinary, config)
		require.NoError(t, err2)

		expectedSha := sha256.New()
		expectedSha.Write(rawBinary)
		expectedSha.Write(config)
		require.Equal(t, fmt.Sprintf("%x", expectedSha.Sum(nil)), actualSha)

		require.Equal(t, *expected, actual)

		assert.Equal(t, b.Bytes(), rawSpec)
	})

	t.Run("Compressed binary", func(t *testing.T) {
		ctx := testutils.Context(t)
		brLoc := strings.Replace(binaryLocation, ".wasm", ".br", 1)
		compressedBytes := b.Bytes()
		require.NoError(t, os.WriteFile(brLoc, compressedBytes, 0600))

		factory := job.WasmFileSpecFactory{}
		actual, rawSpec, actualSha, err2 := factory.Spec(testutils.Context(t), brLoc, configLocation)
		require.NoError(t, err2)

		expected, err2 := host.GetWorkflowSpec(ctx, &host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}, rawBinary, config)
		require.NoError(t, err2)

		expectedSha := sha256.New()
		expectedSha.Write(rawBinary)
		expectedSha.Write(config)
		require.Equal(t, fmt.Sprintf("%x", expectedSha.Sum(nil)), actualSha)

		require.Equal(t, *expected, actual)

		assert.Equal(t, b.Bytes(), rawSpec)
	})

	t.Run("Config", func(t *testing.T) {
		factory := job.WasmFileSpecFactory{}
		actual, err3 := factory.Config(testutils.Context(t), configLocation)
		require.NoError(t, err3)

		assert.Equal(t, config, actual)
	})
}

func createTestBinary(t *testing.T) string {
	const testBinaryLocation = "testdata/wasm/testmodule.wasm"

	cmd := exec.Command("go", "build", "-o", testBinaryLocation, "github.com/smartcontractkit/chainlink/v2/core/services/job/testdata/wasm")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	return testBinaryLocation
}
