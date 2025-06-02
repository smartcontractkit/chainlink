package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestSetBundleFeedConfig(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(
		cldf_chain.WithFamily(chain_selectors.FamilyEVM),
	)[0]

	newEnv, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			changeset.DeployCacheChangeset,
			types.DeployConfig{
				ChainsToDeploy: []uint64{chainSelector},
				Labels:         []string{"data-feeds"},
			},
		),
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
			map[uint64]commonTypes.MCMSWithTimelockConfigV2{
				chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	cacheAddress, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache") //nolint:staticcheck // TODO: replace with DataStore when ready
	require.NoError(t, err)

	dataid := "0x01bb0467f50003040000000000000000"

	// without MCMS
	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.SetFeedAdminChangeset,
			types.SetFeedAdminConfig{
				ChainSelector: chainSelector,
				CacheAddress:  common.HexToAddress(cacheAddress),
				AdminAddress:  common.HexToAddress(env.BlockChains.EVMChains()[chainSelector].DeployerKey.From.Hex()),
				IsAdmin:       true,
			},
		),
		commonChangesets.Configure(
			changeset.SetBundleFeedConfigChangeset,
			types.SetFeedBundleConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   common.HexToAddress(cacheAddress),
				DataIDs:        []string{dataid},
				Descriptions:   []string{"test"},
				DecimalsMatrix: [][]uint8{{18, 8, 6}},
				WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
					{
						AllowedSender:        common.HexToAddress("0x22"),
						AllowedWorkflowOwner: common.HexToAddress("0x33"),
						AllowedWorkflowName:  changeset.HashedWorkflowName("test"),
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// with MCMS
	timeLockAddress, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "RBACTimelock") //nolint:staticcheck // TODO: replace with DataStore when ready
	require.NoError(t, err)

	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		// Set the admin to the timelock
		commonChangesets.Configure(
			changeset.SetFeedAdminChangeset,
			types.SetFeedAdminConfig{
				ChainSelector: chainSelector,
				CacheAddress:  common.HexToAddress(cacheAddress),
				AdminAddress:  common.HexToAddress(timeLockAddress),
				IsAdmin:       true,
			},
		),
		// Transfer cache ownership to MCMS
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2),
			commonChangesets.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chainSelector: {common.HexToAddress(cacheAddress)},
				},
				MCMSConfig: proposalutils.TimelockConfig{MinDelay: 0},
			},
		),
	)
	require.NoError(t, err)

	// Set the feed config with MCMS
	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.SetBundleFeedConfigChangeset,
			types.SetFeedBundleConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   common.HexToAddress(cacheAddress),
				DataIDs:        []string{dataid},
				Descriptions:   []string{"test2"},
				DecimalsMatrix: [][]uint8{{18, 8, 6}},
				WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
					{
						AllowedSender:        common.HexToAddress("0x22"),
						AllowedWorkflowOwner: common.HexToAddress("0x33"),
						AllowedWorkflowName:  changeset.HashedWorkflowName("test"),
					},
				},
				McmsConfig: &types.MCMSConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)
}
