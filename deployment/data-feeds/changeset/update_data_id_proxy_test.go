package changeset_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestUpdateDataIDProxyMap(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	// Deploy and configure pre-requisite contracts
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

	_, err := changeset.SetFeedAdminChangeset(env, types.SetFeedAdminConfig{
		ChainSelector: chainSelector,
		CacheAddress:  common.HexToAddress(cacheAddress),
		AdminAddress:  common.HexToAddress(env.Chains[chainSelector].DeployerKey.From.Hex()),
		IsAdmin:       true,
	})
	// End of pre-requisite contracts

	dataID, _ := shared.ConvertHexToBytes16("01bb0467f50003040000000000000000")

	resp, err := changeset.UpdateDataIDProxyChangeset(env, types.UpdateDataIDProxyConfig{
		ChainSelector: chainSelector,
		CacheAddress:  common.HexToAddress(cacheAddress),
		Proxies:       []common.Address{common.HexToAddress("0x11")},
		DataIDs:       [][16]byte{dataID},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

}
