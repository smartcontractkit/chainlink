package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
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

	cacheAddress, err := deployment.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache")
	require.NoError(t, err)

	dataID := "0x01bb0467f50003040000000000000000"

	// without MCMS
	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
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
			changeset.UpdateDataIDProxyChangeset,
			types.UpdateDataIDProxyConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   common.HexToAddress(cacheAddress),
				ProxyAddresses: []common.Address{common.HexToAddress("0x11")},
				DataIDs:        []string{dataID},
			},
		),
	)
	require.NoError(t, err)

	// with MCMS
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

	newEnv, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			changeset.UpdateDataIDProxyChangeset,
			types.UpdateDataIDProxyConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   common.HexToAddress(cacheAddress),
				ProxyAddresses: []common.Address{common.HexToAddress("0x11")},
				DataIDs:        []string{dataID},
				McmsConfig: &types.MCMSConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)
}
