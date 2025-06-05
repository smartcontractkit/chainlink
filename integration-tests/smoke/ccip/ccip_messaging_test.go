package ccip

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	soltesthelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/manualexechelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
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
		sender   = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
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
				Receiver:               state.MustGetEVMChainState(destChain).FeeQuoter.Address().Bytes(),
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
				Receiver:               state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              nil, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.MustGetEVMChainState(destChain).Receiver.FilterMessageReceived(&bind.FilterOpts{
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
				Receiver:               state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(),
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

func Test_CCIPMessaging_EVM2SolanaMultiExecReports(t *testing.T) {
	// Setup 2 chains (EVM and Solana) and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithSolChains(1),
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.InflightCacheExpiry = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MessageVisibilityInterval = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			params.ExecuteOffChainConfig.MaxReportMessages = 1
			params.ExecuteOffChainConfig.MaxSingleChainReports = 1
			return params
		}),
	)

	testhelpers.DeploySolanaCcipReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilyEVM))
	allSolChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySolana))
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
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
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

	solChains := e.Env.BlockChains.SolanaChains()

	accounts := [][32]byte{
		receiverExternalExecutionConfigPDA,
		receiverTargetAccountPDA,
		solana.SystemProgramID,
	}

	extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
		AccountIsWritableBitmap:  solccip.GenerateBitMapForIndexes([]int{0, 1}),
		Accounts:                 accounts,
		ComputeUnits:             80_000,
		AllowOutOfOrderExecution: true,
	})
	require.NoError(t, err)

	// check that counter is 0
	var receiverCounterAccount soltesthelpers.ReceiverCounter
	err = solcommon.GetAccountDataBorshInto(ctx, solChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
	require.NoError(t, err, "failed to get account info")
	// Already sent a message earlier
	require.Equal(t, uint8(0), receiverCounterAccount.Value)

	numMessages := 5
	_ = mt.Run(
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
			NumberOfMessages:       numMessages,
		},
	)

	t.Logf("Checking TransmittedEvents")
	done := make(chan any)
	defer close(done)
	offrampAddress := state.SolChains[destChain].OffRamp
	sink, errCh := testhelpers.SolEventEmitter[solccip.EventTransmitted](
		t,
		solChains[destChain].Client,
		offrampAddress, "Transmitted", 0, done)
	timeout := time.NewTimer(tests.WaitTimeout(t) - 5*time.Second) // 5 seconds buffer for cleanup
	defer timeout.Stop()

	// create a map/set to keep track of transmittedEvent.SequenceNumber
	sequenceNumbers := make(map[uint64]int)

	for {
		select {
		case transmittedEvent := <-sink:
			if transmittedEvent.OcrPluginType == uint8(cctypes.PluginTypeCCIPExec) {
				ocrSeqNr := transmittedEvent.SequenceNumber
				_, exists := sequenceNumbers[ocrSeqNr]
				if !exists {
					sequenceNumbers[ocrSeqNr] = 0
				}
				sequenceNumbers[ocrSeqNr]++
				t.Logf("Current sequence numbers: %+v", sequenceNumbers)
				// All exec reports should have the same sequence number
				if len(sequenceNumbers) > 1 {
					t.Errorf("More than one sequence number: %+v", sequenceNumbers)
					break
				}

				// All messages were reported with different reports but same sequence number
				if sequenceNumbers[ocrSeqNr] >= numMessages {
					return
				}
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			require.Fail(t, "Timed out waiting for all messages to complete")
		}
	}
}

func Test_CCIPMessaging_EVM2Solana(t *testing.T) {
	// Setup 2 chains (EVM and Solana) and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithSolChains(1),
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.InflightCacheExpiry = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MessageVisibilityInterval = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			params.ExecuteOffChainConfig.MaxReportMessages = 1
			params.ExecuteOffChainConfig.MaxSingleChainReports = 1
			return params
		}),
	)

	// TODO: do this as part of setup
	testhelpers.DeploySolanaCcipReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilyEVM))
	allSolChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySolana))
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
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
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

	solChains := e.Env.BlockChains.SolanaChains()

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		accounts := [][32]byte{
			receiverExternalExecutionConfigPDA,
			receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap:  solccip.GenerateBitMapForIndexes([]int{0, 1}),
			Accounts:                 accounts,
			ComputeUnits:             80_000,
			AllowOutOfOrderExecution: true,
		})
		require.NoError(t, err)

		// check that counter is 0
		var receiverCounterAccount soltesthelpers.ReceiverCounter
		err = solcommon.GetAccountDataBorshInto(ctx, solChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
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
						err = solcommon.GetAccountDataBorshInto(ctx, solChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
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
						err = solcommon.GetAccountDataBorshInto(ctx, solChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccountAfterSuccess)
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

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	allSolChainSelectors := maps.Keys(e.Env.BlockChains.SolanaChains())
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
		sender   = common.LeftPadBytes(e.Env.BlockChains.SolanaChains()[sourceChain].DeployerKey.PublicKey().Bytes(), 32)
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
				Receiver:               state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				FeeToken:               "",        // use native SOL - internally this will be converted to wSOL via Sync Native
				ExtraArgs:              extraArgs, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.MustGetEVMChainState(destChain).Receiver.FilterMessageReceived(&bind.FilterOpts{
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
	state stateview.CCIPOnChainState,
	destChain uint64,
	ss *monitorState,
) {
	sink := make(chan *offramp.OffRampSkippedAlreadyExecutedMessage)
	sub, err := state.MustGetEVMChainState(destChain).OffRamp.WatchSkippedAlreadyExecutedMessage(&bind.WatchOpts{
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
