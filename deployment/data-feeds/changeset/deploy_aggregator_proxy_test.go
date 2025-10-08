package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestAggregatorProxy(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(DeployCacheChangeset, types.DeployConfig{
			ChainsToDeploy: []uint64{selector},
			Labels:         []string{"data-feeds"},
		}),
		runtime.ChangesetTask(DeployAggregatorProxyChangeset, types.DeployAggregatorProxyConfig{
			ChainsToDeploy:   []uint64{selector},
			AccessController: []common.Address{common.HexToAddress("0x")},
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addrs, 2) // AggregatorProxy and DataFeedsCache
}
