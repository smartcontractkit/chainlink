package ccip

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	soltestutils "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	soltesthelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/manualexechelpers"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func Test_CCIPMessaging_EVM2EVM(t *testing.T) {
	// fix the chain ids for the test so we can appropriately set finality depth numbers on the destination chain.
	chains := []chainsel.Chain{
		chainsel.GETH_TESTNET,  // source
		chainsel.TEST_90000001, // dest
	}
	var chainIDs = []uint64{
		chains[0].EvmChainID,
		chains[1].EvmChainID,
	}
	// Setup 2 chains and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithChainIDs(chainIDs),
		testhelpers.WithCLNodeConfigOpts(memory.WithFinalityDepths(map[uint64]uint32{
			chains[1].EvmChainID: 30, // make dest chain finality depth 30 so we can observe exec behavior
		})),
	)

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	require.Len(t, allChainSelectors, 2)
	sourceChain := chains[0].Selector
	destChain := chains[1].Selector
	require.Contains(t, allChainSelectors, sourceChain)
	require.Contains(t, allChainSelectors, destChain)
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)
	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.Chains[sourceChain].DeployerKey.From.Bytes(), 32)
		out      mt.TestCaseOutput
		setup    = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	monitorCtx, monitorCancel := context.WithCancel(ctx)
	ms := &monitorState{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorReExecutions(monitorCtx, t, state, destChain, ms)
	}()

	t.Run("data message to eoa", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				Receiver:               common.HexToAddress("0xdead").Bytes(),
				MsgData:                []byte("hello eoa"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {

					},
				},
			},
		)
	})

	t.Run("message to contract not implementing CCIPReceiver", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  &out.Nonce,
				Receiver:               state.Chains[destChain].FeeQuoter.Address().Bytes(),
				MsgData:                []byte("hello FeeQuoter"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call a contract not implementing CCIPReceiver
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  &out.Nonce,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              nil, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
							Context: ctx,
							Start:   latestHead,
						})
						require.NoError(t, err)
						require.True(t, iter.Next())
						// MessageReceived doesn't emit the data unfortunately, so can't check that.
					},
				},
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  &out.Nonce,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver with low exec gas"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(1, false), // 1 gas is too low.
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_FAILURE,      // state would be failed onchain due to low gas
			},
		)

		err := manualexechelpers.ManuallyExecuteAll(
			ctx,
			e.Env.Logger,
			state,
			e.Env,
			sourceChain,
			destChain,
			[]int64{
				int64(out.MsgSentEvent.Message.Header.SequenceNumber), //nolint:gosec // seqNr fits in int64
			},
			24*time.Hour, // lookbackDurationMsgs
			24*time.Hour, // lookbackDurationCommitReport
			24*time.Hour, // stepDuration
			true,         // reExecuteIfFailed
		)
		require.NoError(t, err)

		t.Logf("successfully manually executed message %x",
			out.MsgSentEvent.Message.Header.MessageId)
	})

	monitorCancel()
	wg.Wait()
	// there should be no re-executions.
	require.Equal(t, int32(0), ms.reExecutionsObserved.Load())
}

func Test_CCIPMessaging_EVM2Solana(t *testing.T) {
	// Setup 2 chains (EVM and Solana) and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithSolChains(1),
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			if params.ExecuteOffChainConfig != nil {
				params.ExecuteOffChainConfig.InflightCacheExpiry = *config.MustNewDuration(1 * time.Hour)
				params.ExecuteOffChainConfig.MessageVisibilityInterval = *config.MustNewDuration(1 * time.Hour)
			}
			return params
		}),
	)

	// TODO: do this as part of setup
	testhelpers.DeploySolanaCcipReceiver(t, e.Env)

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	allSolChainSelectors := maps.Keys(e.Env.SolChains)
	sourceChain := allChainSelectors[0]
	destChain := allSolChainSelectors[0]
	t.Log("All chain selectors:", allChainSelectors,
		", sol chain selectors:", allSolChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)
	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
		// nonce    uint64 // Nonce not used as Solana check is skipped
		sender = common.LeftPadBytes(e.Env.Chains[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	receiverProgram := state.SolChains[destChain].Receiver
	receiver := receiverProgram.Bytes()
	receiverTargetAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverProgram)
	receiverExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverProgram)

	// message := ccip_router.SVM2AnyMessage{
	// 	Receiver:     validReceiverAddress[:],
	// 	FeeToken:     wsol.mint,
	// 	TokenAmounts: []ccip_router.SVMTokenAmount{{Token: token0.Mint.PublicKey(), Amount: 1}},
	// 	ExtraArgs:    emptyEVMExtraArgsV2,
	// }

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		accounts := [][32]byte{
			receiverExternalExecutionConfigPDA,
			receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap: solccip.GenerateBitMapForIndexes([]int{0, 1}),
			Accounts:                accounts,
			ComputeUnits:            80_000,
		})
		require.NoError(t, err)

		// check that counter is 0
		var receiverCounterAccount soltesthelpers.ReceiverCounter
		err = solcommon.GetAccountDataBorshInto(ctx, e.Env.SolChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
		require.NoError(t, err, "failed to get account info")
		require.Equal(t, uint8(0), receiverCounterAccount.Value)

		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  nil, // Solana nonce check is skipped
				Receiver:               receiver,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              extraArgs,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						var receiverCounterAccount soltesthelpers.ReceiverCounter
						err = solcommon.GetAccountDataBorshInto(ctx, e.Env.SolChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
						require.NoError(t, err, "failed to get account info")
						require.Equal(t, uint8(1), receiverCounterAccount.Value)
					},
				},
			},
		)
	})

	t.Run("message sequence: failure (too many accounts) -> success", func(t *testing.T) {
		// --- 1. First Message (Failure - Too Many Accounts) ---
		t.Log("Sending first message (expecting failure due to too many accounts)...")

		// Generate 60 dummy accounts
		numAccounts := 60
		accountsFailure := make([][32]byte, numAccounts)
		writableIndexes := []int{0, 1, 2} // Mark first 3 as writable
		for i := 0; i < numAccounts; i++ {
			accountsFailure[i] = common.HexToHash(fmt.Sprintf("0x%064d", i+1))
		}
		// Set required accounts
		accountsFailure[0] = receiverExternalExecutionConfigPDA
		accountsFailure[1] = receiverTargetAccountPDA
		accountsFailure[2] = solana.SystemProgramID

		extraArgsFailure, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap: solccip.GenerateBitMapForIndexes(writableIndexes),
			Accounts:                accountsFailure,
			ComputeUnits:            80_000,
		})
		require.NoError(t, err, "failed to serialize extra args for failing message")

		// Run the test case expecting failure
		// Use initial replayed=false and nonce=0
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType: mt.ValidationTypeCommit,
				TestSetup:      setup,
				Replayed:       out.Replayed,
				Nonce:          nil, // Nonce check skipped for Commit validation and Solana
				Receiver:       receiver,
				MsgData:        []byte("hello with too many accounts"),
				ExtraArgs:      extraArgsFailure,
			},
		)

		// --- 2. Second Message (Success) ---
		t.Log("Sending second message (expecting success)...")
		accountsSuccess := [][32]byte{ // Use valid accounts
			receiverExternalExecutionConfigPDA,
			receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		extraArgsSuccess, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap: solccip.GenerateBitMapForIndexes([]int{0, 1}), // Mark relevant accounts as writable
			Accounts:                accountsSuccess,
			ComputeUnits:            80_000,
		})
		require.NoError(t, err, "failed to serialize extra args for successful message")

		// Run the test case expecting success
		// Use Replayed and Nonce from the previous (failed) run's output stored in 'out'
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  nil, // Solana nonce check is skipped
				Receiver:               receiver,
				MsgData:                []byte("hello CCIPReceiver that should succeed"),
				ExtraArgs:              extraArgsSuccess,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						// Check counter is now 1
						var receiverCounterAccountAfterSuccess soltesthelpers.ReceiverCounter
						err = solcommon.GetAccountDataBorshInto(ctx, e.Env.SolChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccountAfterSuccess)
						require.NoError(t, err, "failed to get account info after second message")
						require.Equal(t, uint8(2), receiverCounterAccountAfterSuccess.Value, "Counter should have incremented to 2")
						t.Logf("Confirmed counter incremented to 2 after second (successful) message")
					},
				},
			},
		)
	})

	_ = out
}

func Test_CCIPMessaging_Solana2EVM(t *testing.T) {
	// Setup 2 chains (EVM and Solana) and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithSolChains(1))

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	allSolChainSelectors := maps.Keys(e.Env.SolChains)
	sourceChain := allSolChainSelectors[0]
	destChain := allChainSelectors[0]
	t.Log("All chain selectors:", allChainSelectors,
		", sol chain selectors:", allSolChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)
	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
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
		)
	)

	// TODO: handle in setup
	deployer := e.Env.SolChains[sourceChain].DeployerKey
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

	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixAtaUser, ixApprove}, *deployer, solconfig.DefaultCommitment)

	// fund user WSOL (transfer SOL + syncNative)
	transferAmount := 1.0 * solana.LAMPORTS_PER_SOL
	ixTransfer, err := soltokens.NativeTransfer(tokenProgram, transferAmount, deployer.PublicKey(), deployerWSOL)
	require.NoError(t, err)
	ixSync, err := soltokens.SyncNative(tokenProgram, deployerWSOL)
	require.NoError(t, err)
	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixTransfer, ixSync}, *deployer, solconfig.DefaultCommitment)
	// END: handle in setup

	emptyEVMExtraArgsV2 := []byte{}

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		extraArgs := emptyEVMExtraArgsV2
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              extraArgs, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
							Context: ctx,
							Start:   latestHead,
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

type monitorState struct {
	reExecutionsObserved atomic.Int32
}

func (s *monitorState) incReExecutions() {
	s.reExecutionsObserved.Add(1)
}

func monitorReExecutions(
	ctx context.Context,
	t *testing.T,
	state changeset.CCIPOnChainState,
	destChain uint64,
	ss *monitorState,
) {
	sink := make(chan *offramp.OffRampSkippedAlreadyExecutedMessage)
	sub, err := state.Chains[destChain].OffRamp.WatchSkippedAlreadyExecutedMessage(&bind.WatchOpts{
		Start: nil,
	}, sink)
	if err != nil {
		t.Fatalf("failed to subscribe to already executed msg stream: %s", err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			return
		case subErr := <-sub.Err():
			t.Fatalf("subscription error: %s", subErr.Error())
		case ev := <-sink:
			t.Logf("received an already executed event for seq nr %d and source chain %d",
				ev.SequenceNumber, ev.SourceChainSelector)
			ss.incReExecutions()
		}
	}
}
