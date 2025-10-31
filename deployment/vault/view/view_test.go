package view

import (
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestVault_NoChains(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{Chains: 0})

	viewMarshaler, err := Vault(env, nil)
	require.NoError(t, err)
	require.NotNil(t, viewMarshaler)

	view := viewMarshaler.(*VaultView)
	require.Empty(t, view.TimelockBalances)
	require.Empty(t, view.WhitelistedAddresses)
	require.Empty(t, view.MCMSWithTimelock)
}

func TestGenerateVaultView_WithoutTimelock(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{Chains: 1})

	chainSelectors := getChainSelectors(env)
	require.Len(t, chainSelectors, 1)

	env.DataStore = datastore.NewMemoryDataStore().Seal()

	view, err := GenerateVaultView(env, chainSelectors)
	require.NoError(t, err)
	require.NotNil(t, view)

	require.Empty(t, view.TimelockBalances)

	require.Len(t, view.WhitelistedAddresses, len(chainSelectors))
	for _, sel := range chainSelectors {
		require.Empty(t, view.WhitelistedAddresses[sel])
	}

	require.Empty(t, view.MCMSWithTimelock)
}

func TestGenerateVaultView_WithMCMSAndWhitelist(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{Chains: 2})

	chainSelectors := getChainSelectors(env)
	require.Len(t, chainSelectors, 2)

	env = setupMCMS(t, env, chainSelectors)

	whitelistByChain := map[uint64][]types.WhitelistAddress{}
	for i, sel := range chainSelectors {
		addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
		if i == 1 {
			addr = common.HexToAddress("0x2222222222222222222222222222222222222222")
		}
		whitelistByChain[sel] = []types.WhitelistAddress{{
			Address:     addr.Hex(),
			Description: "recipient",
			Labels:      []string{"test"},
		}}
	}
	output, err := changeset.SetWhitelistChangeset.Apply(env, types.SetWhitelistConfig{WhitelistByChain: whitelistByChain})
	require.NoError(t, err)
	require.NotNil(t, output.DataStore)
	env.DataStore = output.DataStore.Seal()

	view, err := GenerateVaultView(env, chainSelectors)
	require.NoError(t, err)
	require.NotNil(t, view)

	require.Len(t, view.TimelockBalances, len(chainSelectors))
	require.Len(t, view.WhitelistedAddresses, len(chainSelectors))
	for _, sel := range chainSelectors {
		entries := view.WhitelistedAddresses[sel]
		require.Len(t, entries, 1)
		require.Equal(t, whitelistByChain[sel][0].Address, entries[0].Address)
	}
	require.Len(t, view.MCMSWithTimelock, len(chainSelectors))
}

func setupMCMS(t *testing.T, env cldf.Environment, chainSelectors []uint64) cldf.Environment {
	t.Helper()
	timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, sel := range chainSelectors {
		timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
	}

	updatedEnv, err := commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			timelockCfgs,
		),
	)
	require.NoError(t, err)
	return updatedEnv
}

func getChainSelectors(env cldf.Environment) []uint64 {
	selectors := make([]uint64, 0)
	for sel := range env.BlockChains.EVMChains() {
		selectors = append(selectors, sel)
	}
	sort.Slice(selectors, func(i, j int) bool {
		return selectors[i] < selectors[j]
	})
	return selectors
}
