package ccip

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	soltestutils "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestCCIPSolCRIB(t *testing.T) {

	simChainTestKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	solTestKey := "57qbvFjTChfNwQxqkFZwjHp7xYoPZa7f9ow6GA59msfCH1g6onSjKUTrrLp4w1nAwbwQuit8YgJJ2AwT9BSwownC"

	// comment out when executing the test
	// t.Skip("Skipping test as this test should not be auto triggered")
	lggr := logger.Test(t)
	ctx, cancel := context.WithCancel(tests.Context(t))
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(*userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	cribEnvironment, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, cribEnvironment)
	userOverrides.Validate(t, cribEnvironment)

	allChainSelectors := cribEnvironment.AllChainSelectors()
	allSolChainSelectors := cribEnvironment.AllChainSelectorsSolana()

	e := testhelpers.DeployedEnv{
		Env:          *cribEnvironment,
		HomeChainSel: allSolChainSelectors[0],
		FeedChainSel: allChainSelectors[0],
	}
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := allSolChainSelectors[0]
	destChain := allChainSelectors[0]
	t.Log("All chain selectors:", allChainSelectors,
		", sol chain selectors:", allSolChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	t.Logf("Deployer key %v", *e.Env.SolChains[sourceChain].DeployerKey)
	t.Logf("Deployer pub key %v", e.Env.SolChains[sourceChain].DeployerKey.PublicKey())

	var (
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.SolChains[sourceChain].DeployerKey.PublicKey().Bytes(), 32)
		out      mt.TestCaseOutput
		setup    = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
			true,  // validateResp
		)
	)

	if false {
		// TODO: handle in setup
		deployer := *e.Env.SolChains[sourceChain].DeployerKey
		rpcClient := e.Env.SolChains[sourceChain].Client

		// create ATA for user
		tokenProgram := solana.TokenProgramID
		wSOL := solana.SolMint
		ixAtaUser, deployerWSOL, uerr := soltokens.CreateAssociatedTokenAccount(tokenProgram, wSOL, deployer.PublicKey(), deployer.PublicKey())
		require.NoError(t, uerr)

		billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(state.SolChains[sourceChain].Router)
		require.NoError(t, err)

		// Approve CCIP to transfer the user's token for billing
		ixApprove, err := soltokens.TokenApproveChecked(1e9, 9, tokenProgram, deployerWSOL, wSOL, billingSignerPDA, deployer.PublicKey(), []solana.PublicKey{})
		require.NoError(t, err)

		soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixAtaUser, ixApprove}, deployer, solconfig.DefaultCommitment)

		// fund user WSOL (transfer SOL + syncNative)
		transferAmount := 1.0 * solana.LAMPORTS_PER_SOL
		ixTransfer, err := soltokens.NativeTransfer(tokenProgram, transferAmount, deployer.PublicKey(), deployerWSOL)
		require.NoError(t, err)
		ixSync, err := soltokens.SyncNative(tokenProgram, deployerWSOL)
		require.NoError(t, err)
		soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixTransfer, ixSync}, deployer, solconfig.DefaultCommitment)
		// END: handle in setup
	}

	emptyEVMExtraArgsV2 := []byte{}
	ccip_router.SetProgramID(state.SolChains[sourceChain].Router)
	fee_quoter.SetProgramID(state.SolChains[sourceChain].FeeQuoter)
	ccip_offramp.SetProgramID(state.SolChains[sourceChain].OffRamp)

	t.Logf("Dest chain %v", state.Chains[destChain].Receiver.Address())
	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		extraArgs := emptyEVMExtraArgsV2
		latestHead, err := e.Env.Chains[destChain].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               true,
				Nonce:                  nonce,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              extraArgs, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
							Context: ctx,
							Start:   latestHead.Number.Uint64(),
						})
						require.NoError(t, err)
						require.True(t, iter.Next())
						// MessageReceived doesn't emit the data unfortunately, so can't check that.
					},
				},
			},
		)

		_ = out // avoid unused error
	})
}
