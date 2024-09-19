package job_test

import (
	"encoding/json"
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
	rawSpec, err := factory.RawSpec(testutils.Context(t), binaryLocation)
	require.NoError(t, err)
	actual, err := factory.Spec(testutils.Context(t), logger.NullLogger, rawSpec, config)
	require.NoError(t, err)

	rawBinary, err := os.ReadFile(binaryLocation)
	require.NoError(t, err)
	expected, err := host.GetWorkflowSpec(&host.ModuleConfig{Logger: logger.NullLogger}, rawBinary, config)
	require.NoError(t, err)

	require.Equal(t, *expected, actual)
}

func createTestBinary(t *testing.T) string {
	const testBinaryLocation = "testdata/wasm/testmodule.wasm"

	cmd := exec.Command("go", "build", "-o", testBinaryLocation, "github.com/smartcontractkit/chainlink/v2/core/services/job/testdata/wasm")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	return testBinaryLocation
}
