package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestProposeAggregator(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	cache, _ := changeset.DeployCache(env.Chains[chainSelector], []string{})
	proxy, _ := changeset.DeployAggregatorProxy(env.Chains[chainSelector], cache.Address, common.HexToAddress("0x"), []string{})

	resp, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			changeset.ProposeAggregatorChangeset,
			types.ProposeConfirmAggregatorConfig{
				ChainSelector:        chainSelector,
				ProxyAddress:         proxy.Address,
				NewAggregatorAddress: common.HexToAddress("0x123"),
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
