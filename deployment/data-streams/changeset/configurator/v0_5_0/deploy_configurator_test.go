package v0_5_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployConfigurator(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployConfiguratorChangeset,
			DeployConfiguratorConfig{
				ChainsToDeploy: []uint64{testutil.TestChain.Selector},
			},
		),
	)
	require.NoError(t, err)

	_, err = dsutil.GetContractAddress(e.DataStore.Addresses(), types.Configurator)
	require.NoError(t, err)
}
