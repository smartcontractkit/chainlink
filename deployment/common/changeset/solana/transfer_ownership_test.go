package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
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
	solanaSelector := env.AllChainSelectorsSolana()[0]

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
	deployer := env.SolChains[solanaSelector].DeployerKey.PublicKey()
	assertOwner(t, env, solanaSelector, chainState, deployer)

	// --- act ---
	_, _, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
	require.NoError(t, err)

	// --- assert ---
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	assertOwner(t, env, solanaSelector, chainState, timelockSignerPDA)
}

func deployMCMS(t *testing.T, env deployment.Environment, selector uint64) *state.MCMSWithTimelockStateSolana {
	t.Helper()

	solanaChain := env.SolChains[selector]
	addressBook := deployment.NewMemoryAddressBook()
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
	t *testing.T, env deployment.Environment, selector uint64, chainState *state.MCMSWithTimelockStateSolana, owner solana.PublicKey,
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
	t *testing.T, want solana.PublicKey, configPDA solana.PublicKey, env deployment.Environment, selector uint64,
) {
	t.Helper()
	var config mcmBindings.MultisigConfig
	err := env.SolChains[selector].GetAccountDataBorshInto(env.GetContext(), configPDA, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func assertTimelockOwner(
	t *testing.T, want solana.PublicKey, configPDA solana.PublicKey, env deployment.Environment, selector uint64,
) {
	t.Helper()
	var config timelockBindings.Config
	err := env.SolChains[selector].GetAccountDataBorshInto(env.GetContext(), configPDA, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func assertAccessControllerOwner(
	t *testing.T, want solana.PublicKey, account solana.PublicKey, env deployment.Environment, selector uint64,
) {
	t.Helper()
	var config accessControllerBindings.AccessController
	err := env.SolChains[selector].GetAccountDataBorshInto(env.GetContext(), account, &config)
	require.NoError(t, err)
	require.Equal(t, want, config.Owner)
}

func fundSignerPDAs(
	t *testing.T, env deployment.Environment, chainSelector uint64, chainState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()
	solChain := env.SolChains[chainSelector]
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(chainState.McmProgram, chainState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	memory.FundSolanaAccounts(env.GetContext(), t, signerPDAs, 1, solChain.Client)
}
