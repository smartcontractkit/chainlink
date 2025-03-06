package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
	cache "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestNewFeedWithProxy(t *testing.T) {
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
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
			map[uint64]commonTypes.MCMSWithTimelockConfigV2{
				chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	cacheAddress, err := deployment.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache")
	require.NoError(t, err)

	timeLockAddress, err := deployment.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "RBACTimelock")
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
			deployment.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2),
			commonChangesets.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					chainSelector: {common.HexToAddress(cacheAddress)},
				},
				MinDelay: 0,
			},
		),
	)
	require.NoError(t, err)

	dataid, _ := shared.ConvertHexToBytes16("01bb0467f50003040000000000000000")

	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.NewFeedWithProxyChangeset,
			types.NewFeedWithProxyConfig{
				ChainSelector:    chainSelector,
				AccessController: common.HexToAddress("0x00"),
				DataID:           dataid,
				Description:      "test2",
				WorkflowMetadata: []cache.DataFeedsCacheWorkflowMetadata{
					cache.DataFeedsCacheWorkflowMetadata{
						AllowedSender:        common.HexToAddress("0x22"),
						AllowedWorkflowOwner: common.HexToAddress("0x33"),
						AllowedWorkflowName:  shared.HashedWorkflowName("test"),
					},
				},
				McmsConfig: &types.MCMSConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)

	addrs, err := newEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	// AggregatorProxy, DataFeedsCache, CallProxy, RBACTimelock, ProposerManyChainMultiSig, BypasserManyChainMultiSig, CancellerManyChainMultiSig
	require.Len(t, addrs, 7)
}
