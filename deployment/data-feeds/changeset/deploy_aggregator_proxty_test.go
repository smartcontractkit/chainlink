package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestAggregatorProxy(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	ab, _ := DeployCacheChangeset(env, types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"},
	})
	addresses, _ := ab.AddressBook.Addresses()
	env.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(addresses)

	resp, err := DeployAggregatorProxyChangeset(env, types.DeployAggregatorProxyConfig{
		ChainsToDeploy:   []uint64{chainSelector},
		AccessController: []common.Address{common.HexToAddress("0x")},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// AggregatorProxy should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// no AggregatorProxy deployed on chain 1
	require.NotEqual(t, chainSelector, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
	require.Empty(t, oaddrs)
}
