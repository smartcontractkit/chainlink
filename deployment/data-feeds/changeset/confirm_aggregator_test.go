package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestConfirmAggregator(t *testing.T) {
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

	tx, _ := proxy.Contract.ProposeAggregator(env.Chains[chainSelector].DeployerKey, common.HexToAddress("0x123"))
	_, err := env.Chains[chainSelector].Confirm(tx)

	resp, err := changeset.ConfirmAggregatorChangeset(env, types.ProposeConfirmAggregatorConfig{
		ChainSelector: chainSelector,
		Proxy:         proxy.Address,
		NewAggregator: common.HexToAddress("0x123"),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}
