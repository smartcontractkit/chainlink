package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestMigrateFeeds(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]

	newEnv, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			changeset.DeployCacheChangeset,
			types.DeployConfig{
				ChainsToDeploy: []uint64{chainSelector},
				Labels:         []string{"data-feeds"},
			},
		),
	)
	require.NoError(t, err)

	cacheAddress, err := deployment.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache")
	require.NoError(t, err)

	resp, err := commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.SetFeedAdminChangeset,
			types.SetFeedAdminConfig{
				ChainSelector: chainSelector,
				CacheAddress:  common.HexToAddress(cacheAddress),
				AdminAddress:  common.HexToAddress(env.Chains[chainSelector].DeployerKey.From.Hex()),
				IsAdmin:       true,
			},
		),
		commonChangesets.Configure(
			changeset.MigrateFeedsChangeset,
			types.MigrationConfig{
				ChainSelector: chainSelector,
				CacheAddress:  common.HexToAddress(cacheAddress),
				InputFileName: "testdata/migrate_feeds.json",
				InputFS:       testFS,
				WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
					cache.DataFeedsCacheWorkflowMetadata{
						AllowedSender:        common.HexToAddress("0x22"),
						AllowedWorkflowOwner: common.HexToAddress("0x33"),
						AllowedWorkflowName:  changeset.HashedWorkflowName("test"),
					},
				},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	addresses, err := resp.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addresses, 3) // DataFeedsCache and two migrated proxies
}
