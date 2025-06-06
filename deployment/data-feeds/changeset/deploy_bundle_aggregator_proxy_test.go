package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestBundleAggregatorProxy(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	resp, err := commonChangesets.Apply(t, env, commonChangesets.Configure(
		DeployCacheChangeset,
		types.DeployConfig{
			ChainsToDeploy: []uint64{chainSelector},
			Labels:         []string{"data-feeds"},
		},
	), commonChangesets.Configure(
		DeployBundleAggregatorProxyChangeset,
		types.DeployBundleAggregatorProxyConfig{
			ChainsToDeploy: []uint64{chainSelector},
			Owners:         map[uint64]common.Address{chainSelector: common.HexToAddress("0x1234")},
			Labels:         []string{"data-feeds"},
			CacheLabel:     "data-feeds",
		},
	))

	require.NoError(t, err)
	require.NotNil(t, resp)

	addrs, err := resp.ExistingAddresses.AddressesForChain(chainSelector) //nolint:staticcheck // TODO: replace with DataStore when ready
	require.NoError(t, err)
	require.Len(t, addrs, 2) // BundleAggregatorProxy and DataFeedsCache
}
