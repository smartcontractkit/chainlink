package tron_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestDeployCache(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TRON_DEVNET.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithTronContainer(t, []uint64{selector}),
	))
	require.NoError(t, err)

	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	err = rt.Exec(
		runtime.ChangesetTask(tron.DeployCacheChangeset, types.DeployTronConfig{
			ChainsToDeploy: []uint64{selector},
			Labels:         []string{"data-feeds"},
			Qualifier:      "tron",
			DeployOptions:  deployOptions,
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			selector,
			"DataFeedsCache",
			semver.MustParse("1.0.0"),
			"tron",
		))
	require.NoError(t, err)
	require.NotNil(t, addrs.Address)
	require.Equal(t, datastore.ContractType("DataFeedsCache"), addrs.Type)
	require.Equal(t, "tron", addrs.Qualifier)
}
