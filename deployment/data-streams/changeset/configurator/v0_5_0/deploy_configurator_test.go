package v0_5_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployConfigurator(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	out, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployConfiguratorChangeset,
			DeployConfiguratorConfig{
				ChainsToDeploy: []uint64{testutil.TestChain.Selector},
			},
		),
	)

	require.NoError(t, err)

	ab, err := out.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, testutil.TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, types.Configurator, tv.Type)
		}
	}
}
