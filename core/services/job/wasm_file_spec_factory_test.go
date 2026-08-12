package job_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestWasmFileSpecFactory(t *testing.T) {
	const pkgRelPath = "core/services/job/testdata/wasm"
	configLocation := "testdata/config.json"
	config, err := os.ReadFile(configLocation)
	require.NoError(t, err)

	rawBinary := wasmtest.GetTestBinary(t, pkgRelPath, false)
	compressedBinary := wasmtest.GetTestBinary(t, pkgRelPath, true)

	tmpDir := t.TempDir()
	rawBinaryPath := filepath.Join(tmpDir, "testmodule.wasm")
	compressedBinaryPath := filepath.Join(tmpDir, "testmodule.br")
	require.NoError(t, os.WriteFile(rawBinaryPath, rawBinary, 0o600))
	require.NoError(t, os.WriteFile(compressedBinaryPath, compressedBinary, 0o600))

	t.Run("Raw binary", func(t *testing.T) {
		ctx := t.Context()
		factory := job.WasmFileSpecFactory{}
		actual, rawSpec, actualSha, err2 := factory.Spec(t.Context(), rawBinaryPath, configLocation)
		require.NoError(t, err2)

		expected, err2 := host.GetWorkflowSpec(ctx, &host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}, rawBinary, config)
		require.NoError(t, err2)

		expectedSha := sha256.New()
		expectedSha.Write(rawBinary)
		expectedSha.Write(config)
		require.Equal(t, hex.EncodeToString(expectedSha.Sum(nil)), actualSha)

		require.Equal(t, *expected, actual)

		assert.Equal(t, compressedBinary, rawSpec)
	})

	t.Run("Compressed binary", func(t *testing.T) {
		ctx := t.Context()
		factory := job.WasmFileSpecFactory{}
		actual, rawSpec, actualSha, err2 := factory.Spec(t.Context(), compressedBinaryPath, configLocation)
		require.NoError(t, err2)

		expected, err2 := host.GetWorkflowSpec(ctx, &host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}, rawBinary, config)
		require.NoError(t, err2)

		expectedSha := sha256.New()
		expectedSha.Write(rawBinary)
		expectedSha.Write(config)
		require.Equal(t, hex.EncodeToString(expectedSha.Sum(nil)), actualSha)

		require.Equal(t, *expected, actual)

		assert.Equal(t, compressedBinary, rawSpec)
	})

	t.Run("Config", func(t *testing.T) {
		factory := job.WasmFileSpecFactory{}
		actual, err3 := factory.Config(t.Context(), configLocation)
		require.NoError(t, err3)

		assert.Equal(t, config, actual)
	})
}
