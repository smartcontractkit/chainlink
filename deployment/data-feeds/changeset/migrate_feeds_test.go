package changeset_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestMigrateFeeds(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(changeset.DeployCacheChangeset, types.DeployConfig{
			ChainsToDeploy: []uint64{selector},
			Labels:         []string{"data-feeds"},
			Qualifier:      "data-feeds",
		}),
	)
	require.NoError(t, err)

	records := rt.State().DataStore.Addresses().Filter(datastore.AddressRefByType("DataFeedsCache"))
	require.Len(t, records, 1)
	cacheAddress := records[0].Address

	err = rt.Exec(
		runtime.ChangesetTask(changeset.SetFeedAdminChangeset, types.SetFeedAdminConfig{
			ChainSelector: selector,
			CacheAddress:  common.HexToAddress(cacheAddress),
			AdminAddress:  common.HexToAddress(rt.Environment().BlockChains.EVMChains()[selector].DeployerKey.From.Hex()),
			IsAdmin:       true,
		}),
		runtime.ChangesetTask(changeset.MigrateFeedsChangeset, types.MigrationConfig{
			ChainSelector: selector,
			CacheAddress:  common.HexToAddress(cacheAddress),
			Proxies: []*types.MigrationSchema{
				{
					Address:     "0x33442400910b7B03316fe47eF8fC7bEd54Bca407",
					FeedID:      "0x01bb0467f50003040000000000000000",
					Description: "TEST / USD",
					TypeAndVersion: cldf.TypeAndVersion{
						Type:    "AggregatorProxy",
						Version: *semver.MustParse("1.0.0"),
					},
				},
				{
					Address:     "0x43442400910b7B03316fe47eF8fC7bEd54Bca407",
					FeedID:      "0x01b40467f50003040000000000000000",
					Description: "LINK / USD",
					TypeAndVersion: cldf.TypeAndVersion{
						Type:    "AggregatorProxy",
						Version: *semver.MustParse("1.0.0"),
					},
				},
			},
			WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
				{
					AllowedSender:        common.HexToAddress("0x22"),
					AllowedWorkflowOwner: common.HexToAddress("0x33"),
					AllowedWorkflowName:  changeset.HashedWorkflowName("test"),
				},
			},
		}),
	)
	require.NoError(t, err)

	addresses, err := rt.State().DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addresses, 3) // DataFeedsCache and two migrated proxies
}
