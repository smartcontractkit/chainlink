package changeset

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

func TestHydrateCapabilityRegistry(t *testing.T) {
	b, err := os.ReadFile("testdata/capability_registry_view.json")
	require.NoError(t, err)
	require.NotEmpty(t, b)
	var capabilityRegistryView v1_0.CapabilityRegistryView
	require.NoError(t, json.Unmarshal(b, &capabilityRegistryView))

	chainSelector := chainsel.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{chainSelector}),
	)
	require.NoError(t, err)

	cfg := HydrateConfig{ChainSelector: chainSelector}
	hydrated, err := HydrateCapabilityRegistry(t, capabilityRegistryView, *env, cfg)
	require.NoError(t, err)
	require.NotNil(t, hydrated)
	hydratedCapView, err := v1_0.GenerateCapabilityRegistryView(hydrated)
	require.NoError(t, err)

	// Setting address/owner values to be the same in order to compare the views
	hydratedCapView.Address = capabilityRegistryView.Address
	hydratedCapView.Owner = capabilityRegistryView.Owner
	b1, err := capabilityRegistryView.MarshalJSON()
	require.NoError(t, err)
	b2, err := hydratedCapView.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(b1), string(b2))
}
