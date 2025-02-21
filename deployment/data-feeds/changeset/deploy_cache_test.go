package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployCache(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	resp, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			changeset.DeployCacheChangeset,
			types.DeployConfig{
				ChainsToDeploy: []uint64{chainSelector},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// DataFeedsCache should be deployed on chain 0
	addrs, err := resp.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// no DataFeedsCache deployed on chain 1
	require.NotEqual(t, chainSelector, env.AllChainSelectors()[1])
	oaddrs, err := resp.ExistingAddresses.AddressesForChain(env.AllChainSelectors()[1])
	require.NoError(t, err)
	require.Empty(t, oaddrs)
}
