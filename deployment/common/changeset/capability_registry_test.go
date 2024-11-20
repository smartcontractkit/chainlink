package changeset

import (
	"encoding/json"
	"os"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestHydrateCapabilityRegistry(t *testing.T) {
	b, err := os.ReadFile("testdata/capability_registry_view.json")
	require.NoError(t, err)
	require.NotEmpty(t, b)
	var capabilityRegistryView v1_0.CapabilityRegistryView
	require.NoError(t, json.Unmarshal(b, &capabilityRegistryView))

	chainID := chainsel.TEST_90000001.EvmChainID + 1
	cfg := HydrateConfig{ChainID: chainID}
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	hydrated, err := HydrateCapabilityRegistry(capabilityRegistryView, env, cfg)
	require.NoError(t, err)
	require.NotNil(t, hydrated)
}
