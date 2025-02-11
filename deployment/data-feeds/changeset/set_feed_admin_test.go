package changeset_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestSetCacheAdmin(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	ab, _ := changeset.DeployCacheChangeset(env, types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"},
	})
	addresses, _ := ab.AddressBook.Addresses()

	chainAddresses, _ := ab.AddressBook.AddressesForChain(chainSelector)
	var cacheAddress string
	for address, tv := range chainAddresses {
		if strings.Contains(tv.String(), "DataFeedsCache") {
			cacheAddress = address
		}
		break
	}
	env.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(addresses)

	resp, err := changeset.SetFeedAdminChangeset(env, types.SetFeedAdminConfig{
		ChainSelector: chainSelector,
		CacheAddress:  common.HexToAddress(cacheAddress),
		AdminAddress:  common.HexToAddress("0x123"),
		IsAdmin:       true,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

}
