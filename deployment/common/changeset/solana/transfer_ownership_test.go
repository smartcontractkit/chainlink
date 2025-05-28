package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanachangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTransferToMCMSToTimelockSolana(t *testing.T) {
	t.Parallel()
	// --- arrange ---
	log := logger.TestLogger(t)
	envConfig := memory.MemoryEnvironmentConfig{Chains: 0, SolChains: 1}
	env := memory.NewMemoryEnvironment(t, log, zapcore.InfoLevel, envConfig)
	solanaSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]

	commonchangeset.SetPreloadedSolanaAddresses(t, env, solanaSelector)
	chainState := deployMCMS(t, env, solanaSelector)
	fundSignerPDAs(t, env, solanaSelector, chainState)

	configuredChangeset := commonchangeset.Configure(
		&solanachangesets.TransferMCMSToTimelockSolana{},
		solanachangesets.TransferMCMSToTimelockSolanaConfig{
			Chains:  []uint64{solanaSelector},
			MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
		},
	)
	// validate initial owner
	deployer := env.BlockChains.SolanaChains()[solanaSelector].DeployerKey.PublicKey()
	assertOwner(t, env, solanaSelector, chainState, deployer)

	// --- act ---
	_, _, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
	require.NoError(t, err)

	// --- assert ---
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	assertOwner(t, env, solanaSelector, chainState, timelockSignerPDA)
}

func deployMCMS(t *testing.T, env cldf.Environment, selector uint64) *state.MCMSWithTimelockStateSolana {
	t.Helper()

	solanaChain := env.BlockChains.SolanaChains()[selector]
	addressBook := cldf.NewMemoryAddressBook()
	mcmsConfig := commontypes.MCMSWithTimelockConfigV2{
		Canceller:        proposalutils.SingleGroupMCMSV2(t),
		Bypasser:         proposalutils.SingleGroupMCMSV2(t),
		Proposer:         proposalutils.SingleGroupMCMSV2(t),
		TimelockMinDelay: big.NewInt(1),
	}

	chainState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, solanaChain, addressBook, mcmsConfig)
	require.NoError(t, err)
	err = env.ExistingAddresses.Merge(addressBook)
	require.NoError(t, err)

	return chainState
}

func assertOwner(
	t *testing.T, env cldf.Environment, selector uint64, chainState *state.MCMSWithTimelockStateSolana, owner solana.PublicKey,
) {
	assertMCMOwner(t, owner, state.GetMCMConfigPDA(chainState.McmProgram, chainState.ProposerMcmSeed), env, selector)
	assertMCMOwner(t, owner, state.GetMCMConfigPDA(chainState.McmProgram, chainState.CancellerMcmSeed), env, selector)
	assertMCMOwner(t, owner, state.GetMCMConfigPDA(chainState.McmProgram, chainState.BypasserMcmSeed), env, selector)
	assertTimelockOwner(t, owner, state.GetTimelockConfigPDA(chainState.TimelockProgram, chainState.TimelockSeed), env, selector)
	assertAccessControllerOwner(t, owner, chainState.ProposerAccessControllerAccount, env, selector)
	assertAccessControllerOwner(t, owner, chainState.ExecutorAccessControllerAccount, env, selector)
	assertAccessControllerOwner(t, owner, chainState.CancellerAccessControllerAccount, env, selector)
	assertAccessControllerOwner(t, owner, chainState.BypasserAccessControllerAccount, env, selector)
}

func assertMCMOwner(
	t *testing.T, want solana.PublicKey, configPDA solana.PublicKey, env cldf.Environment, selector uint64,
) {
	t.Helper()
	var config mcmBindings.MultisigConfig
	err := env.BlockChains.SolanaChains()[selector].GetAccountDataBorshInto(env.GetContext(), configPDA, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func assertTimelockOwner(
	t *testing.T, want solana.PublicKey, configPDA solana.PublicKey, env cldf.Environment, selector uint64,
) {
	t.Helper()
	var config timelockBindings.Config
	err := env.BlockChains.SolanaChains()[selector].GetAccountDataBorshInto(env.GetContext(), configPDA, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func assertAccessControllerOwner(
	t *testing.T, want solana.PublicKey, account solana.PublicKey, env cldf.Environment, selector uint64,
) {
	t.Helper()
	var config accessControllerBindings.AccessController
	err := env.BlockChains.SolanaChains()[selector].GetAccountDataBorshInto(env.GetContext(), account, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func fundSignerPDAs(
	t *testing.T, env cldf.Environment, chainSelector uint64, chainState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()
	solChain := env.BlockChains.SolanaChains()[chainSelector]
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(chainState.McmProgram, chainState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	err := memory.FundSolanaAccounts(env.GetContext(), signerPDAs, 1, solChain.Client)
	require.NoError(t, err)
}
