package changeset

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
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

	chainSelector, err := chainsel.SelectorFromChainId(chainID)
	require.NoError(t, err)
	chain, ok := env.Chains[chainSelector]
	require.True(t, ok)
	capabilityRegistry, err := capabilities_registry.NewCapabilitiesRegistry(hydrated.Address, chain.Client)
	require.NoError(t, err)
	require.NotNil(t, capabilityRegistry)
	capView, err := v1_0.GenerateCapabilityRegistryView(capabilityRegistry)
	require.NoError(t, err)

	// Setting address/owner values to be the same in order to compare the views
	capView.Address = capabilityRegistryView.Address
	capView.Owner = capabilityRegistryView.Owner
	b1, err := capabilityRegistryView.MarshalJSON()
	require.NoError(t, err)
	b2, err := capView.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(b1), string(b2))
}
