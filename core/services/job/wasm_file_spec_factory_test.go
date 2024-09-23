package job_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestWasmFileSpecFactory(t *testing.T) {
	binaryLocation := createTestBinary(t)
	config, err := json.Marshal(sdk.NewWorkflowParams{
		Owner: "owner",
		Name:  "name",
	})
	require.NoError(t, err)

	factory := job.WasmFileSpecFactory{}
	actual, actualSha, err := factory.Spec(testutils.Context(t), logger.NullLogger, binaryLocation, config)
	require.NoError(t, err)

	rawBinary, err := os.ReadFile(binaryLocation)
	require.NoError(t, err)
	expected, err := host.GetWorkflowSpec(&host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}, rawBinary, config)
	require.NoError(t, err)

	expectedSha := sha256.New()
	expectedSha.Write(rawBinary)
	expectedSha.Write(config)
	require.Equal(t, fmt.Sprintf("%x", expectedSha.Sum(nil)), actualSha)

	require.Equal(t, *expected, actual)
}

func createTestBinary(t *testing.T) string {
	const testBinaryLocation = "testdata/wasm/testmodule.wasm"

	cmd := exec.Command("go1.22.7", "build", "-o", testBinaryLocation, "github.com/smartcontractkit/chainlink/v2/core/services/job/testdata/wasm")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	return testBinaryLocation
}
