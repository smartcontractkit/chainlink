package changeset_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	cache "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/data_feeds_cache"
)

func TestSetFeedConfig(t *testing.T) {
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
	require.NoError(t, err)
	// End of pre-requisite contracts

	dataid, _ := shared.ConvertHexToBytes16("01bb0467f50003040000000000000000")

	resp, err := changeset.SetFeedConfigChangeset(env, types.SetFeedDecimalConfig{
		ChainSelector: chainSelector,
		CacheAddress:  common.HexToAddress(cacheAddress),
		DataIDs:       [][16]byte{dataid},
		Descriptions:  []string{"test"},
		WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
			cache.DataFeedsCacheWorkflowMetadata{
				AllowedSender:        common.HexToAddress("0x22"),
				AllowedWorkflowOwner: common.HexToAddress("0x33"),
				AllowedWorkflowName:  shared.HashedWorkflowName("test"),
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
}
