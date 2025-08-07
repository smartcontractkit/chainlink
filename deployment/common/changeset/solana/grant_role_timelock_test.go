package solana_test

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanachangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGrantRoleTimelockSolana(t *testing.T) {
	t.Skip("fails with Program is not deployed (DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA) in CI")
	t.Parallel()
	// --- arrange ---
	log := logger.TestLogger(t)
	envConfig := memory.MemoryEnvironmentConfig{Chains: 0, SolChains: 1}
	env := memory.NewMemoryEnvironment(t, log, zapcore.InfoLevel, envConfig)
	solanaSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	deployer := env.BlockChains.SolanaChains()[solanaSelector].DeployerKey
	rpcClient := env.BlockChains.SolanaChains()[solanaSelector].Client
	executors1 := randomSolanaAccounts(t, 2)
	executors2 := randomSolanaAccounts(t, 2)

	commonchangeset.SetPreloadedSolanaAddresses(t, env, solanaSelector)
	chainState := deployMCMS(t, env, solanaSelector)
	fundSignerPDAs(t, env, solanaSelector, chainState)

	// validate initial executors
	inspector := mcmssolanasdk.NewTimelockInspector(rpcClient)
	onChainExecutors, err := inspector.GetExecutors(t.Context(), timelockAddress(chainState))
	require.NoError(t, err)
	require.ElementsMatch(t, onChainExecutors, []string{deployer.PublicKey().String()})

	t.Run("without MCMS", func(t *testing.T) {
		nomcmsChangeset := commonchangeset.Configure(
			&solanachangesets.GrantRoleTimelockSolana{},
			solanachangesets.GrantRoleTimelockSolanaConfig{
				Role:     timelockbindings.Executor_Role,
				Accounts: map[uint64][]solana.PublicKey{solanaSelector: executors1},
			},
		)

		// --- act ---
		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{nomcmsChangeset})
		require.NoError(t, err)

		// --- assert ---
		onChainExecutors, err = inspector.GetExecutors(t.Context(), timelockAddress(chainState))
		require.NoError(t, err)
		require.ElementsMatch(t, onChainExecutors, []string{
			deployer.PublicKey().String(), executors1[0].String(), executors1[1].String(),
		})
	})

	t.Run("with MCMS", func(t *testing.T) {
		transferMCMSToTimelock(t, env, solanaSelector)

		mcmsChangeset := commonchangeset.Configure(
			&solanachangesets.GrantRoleTimelockSolana{},
			solanachangesets.GrantRoleTimelockSolanaConfig{
				Role:     timelockbindings.Executor_Role,
				Accounts: map[uint64][]solana.PublicKey{solanaSelector: executors2},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay:   1 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
			},
		)

		// --- act ---
		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{mcmsChangeset})
		require.NoError(t, err)

		// --- assert ---
		onChainExecutors, err = inspector.GetExecutors(t.Context(), timelockAddress(chainState))
		require.NoError(t, err)
		require.ElementsMatch(t, onChainExecutors, []string{
			deployer.PublicKey().String(),
			executors1[0].String(), executors1[1].String(),
			executors2[0].String(), executors2[1].String(),
		})
	})
}

func randomSolanaAccounts(t *testing.T, n int) []solana.PublicKey {
	t.Helper()
	accounts := make([]solana.PublicKey, n)
	for i := range n {
		privateKey, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		accounts[i] = privateKey.PublicKey()
	}

	return accounts
}

func transferMCMSToTimelock(t *testing.T, env cldf.Environment, selector uint64) {
	t.Helper()
	configuredChangeset := commonchangeset.Configure(
		&solanachangesets.TransferMCMSToTimelockSolana{},
		solanachangesets.TransferMCMSToTimelockSolanaConfig{
			Chains:  []uint64{selector},
			MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
		},
	)

	_, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
	require.NoError(t, err)
	t.Logf("transferred MCMS contracts to timelock")
}

func timelockAddress(chainState *state.MCMSWithTimelockStateSolana) string {
	return state.EncodeAddressWithSeed(chainState.TimelockProgram, chainState.TimelockSeed)
}
