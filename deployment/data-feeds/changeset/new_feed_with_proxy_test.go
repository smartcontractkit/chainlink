package changeset_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestNewFeedWithProxy(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

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

	cacheAddress, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache")
	require.NoError(t, err)

	timeLockAddress, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "RBACTimelock")
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

	dataid := "0x01bb0467f50003040000000000000000"
	dataid2 := "0x01475851f90003320000000000000000"

	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.NewFeedWithProxyChangeset,
			types.NewFeedWithProxyConfig{
				ChainSelector:    chainSelector,
				AccessController: common.HexToAddress("0x00"),
				DataIDs:          []string{dataid, dataid2},
				Descriptions:     []string{"feed1", "feed2"},
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

	addrs, err := newEnv.ExistingAddresses.AddressesForChain(chainSelector)
	fmt.Println(addrs)
	require.NoError(t, err)
	// 2 AggregatorProxy, DataFeedsCache, CallProxy, RBACTimelock, ProposerManyChainMultiSig, BypasserManyChainMultiSig, CancellerManyChainMultiSig
	require.Len(t, addrs, 8)
}
