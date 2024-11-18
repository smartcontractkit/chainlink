package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type testCaseSetup struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            ccdeploy.DeployedEnv
	onchainState           ccdeploy.CCIPOnChainState
	sourceChain, destChain uint64
}

type messagingTestCase struct {
	testCaseSetup
	replayed bool
	nonce    uint64
}

type messagingTestCaseOutput struct {
	replayed     bool
	nonce        uint64
	msgSentEvent *onramp.OnRampCCIPMessageSent
}

func Test_CCIPMessaging(t *testing.T) {
	// Setup 2 chains and a single lane.
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	state, err := ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	require.Len(t, allChainSelectors, 2)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)
	output, err := changeset.DeployPrerequisites(e.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))

	tokenConfig := ccdeploy.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	// Apply migration
	output, err = changeset.InitialDeploy(e.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: allChainSelectors,
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// connect a single lane, source to dest
	require.NoError(t, ccdeploy.AddLaneWithDefaultPrices(e.Env, state, sourceChain, destChain))

	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.Chains[sourceChain].DeployerKey.From.Bytes(), 32)
		out      messagingTestCaseOutput
		setup    = testCaseSetup{
			t:            t,
			sender:       sender,
			deployedEnv:  e,
			onchainState: state,
			sourceChain:  sourceChain,
			destChain:    destChain,
		}
	)

	t.Run("data message to eoa", func(t *testing.T) {
		out = runMessagingTestCase(messagingTestCase{
			testCaseSetup: setup,
			replayed:      replayed,
			nonce:         nonce,
		},
			common.HexToAddress("0xdead"),
			[]byte("hello eoa"),
			nil,                              // default extraArgs
			ccdeploy.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
		)
	})

	t.Run("message to contract not implementing CCIPReceiver", func(t *testing.T) {
		out = runMessagingTestCase(
			messagingTestCase{
				testCaseSetup: setup,
				replayed:      out.replayed,
				nonce:         out.nonce,
			},
			state.Chains[destChain].FeeQuoter.Address(),
			[]byte("hello FeeQuoter"),
			nil,                              // default extraArgs
			ccdeploy.EXECUTION_STATE_SUCCESS, // success because offRamp won't call a contract not implementing CCIPReceiver
		)
	})

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		latestHead, err := e.Env.Chains[destChain].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		out = runMessagingTestCase(
			messagingTestCase{
				testCaseSetup: setup,
				replayed:      out.replayed,
				nonce:         out.nonce,
			},
			state.Chains[destChain].Receiver.Address(),
			[]byte("hello CCIPReceiver"),
			nil, // default extraArgs
			ccdeploy.EXECUTION_STATE_SUCCESS,
			func(t *testing.T) {
				iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
					Context: ctx,
					Start:   latestHead.Number.Uint64(),
				})
				require.NoError(t, err)
				require.True(t, iter.Next())
				// MessageReceived doesn't emit the data unfortunately, so can't check that.
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {
		latestHead, err := e.Env.Chains[destChain].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		out = runMessagingTestCase(
			messagingTestCase{
				testCaseSetup: setup,
				replayed:      out.replayed,
				nonce:         out.nonce,
			},
			state.Chains[destChain].Receiver.Address(),
			[]byte("hello CCIPReceiver with low exec gas"),
			ccdeploy.MakeEVMExtraArgsV2(1, false), // 1 gas is too low.
			ccdeploy.EXECUTION_STATE_FAILURE,      // state would be failed onchain due to low gas
		)

		manuallyExecute(ctx, t, latestHead.Number.Uint64(), state, destChain, out, sourceChain, e, sender)

		t.Logf("successfully manually executed message %x",
			out.msgSentEvent.Message.Header.MessageId)
	})
}

func manuallyExecute(
	ctx context.Context,
	t *testing.T,
	startBlock uint64,
	state ccdeploy.CCIPOnChainState,
	destChain uint64,
	out messagingTestCaseOutput,
	sourceChain uint64,
	e ccdeploy.DeployedEnv,
	sender []byte,
) {
	merkleRoot := getMerkleRoot(
		ctx,
		t,
		state.Chains[destChain].OffRamp,
		out.msgSentEvent.SequenceNumber,
		startBlock,
	)
	messageHash := getMessageHash(
		ctx,
		t,
		state.Chains[destChain].OffRamp,
		sourceChain,
		out.msgSentEvent.SequenceNumber,
		out.msgSentEvent.Message.Header.MessageId,
		startBlock,
	)
	tree, err := merklemulti.NewTree(hashutil.NewKeccak(), [][32]byte{messageHash})
	require.NoError(t, err)
	proof, err := tree.Prove([]int{0})
	require.NoError(t, err)
	require.Equal(t, merkleRoot, tree.Root())

	tx, err := state.Chains[destChain].OffRamp.ManuallyExecute(
		e.Env.Chains[destChain].DeployerKey,
		[]offramp.InternalExecutionReport{
			{
				SourceChainSelector: sourceChain,
				Messages: []offramp.InternalAny2EVMRampMessage{
					{
						Header: offramp.InternalRampMessageHeader{
							MessageId:           out.msgSentEvent.Message.Header.MessageId,
							SourceChainSelector: sourceChain,
							DestChainSelector:   destChain,
							SequenceNumber:      out.msgSentEvent.SequenceNumber,
							Nonce:               out.msgSentEvent.Message.Header.Nonce,
						},
						Sender:       sender,
						Data:         []byte("hello CCIPReceiver with low exec gas"),
						Receiver:     state.Chains[destChain].Receiver.Address(),
						GasLimit:     big.NewInt(1),
						TokenAmounts: []offramp.InternalAny2EVMTokenTransfer{},
					},
				},
				OffchainTokenData: [][][]byte{
					{},
				},
				Proofs:        proof.Hashes,
				ProofFlagBits: boolsToBitFlags(proof.SourceFlags),
			},
		},
		[][]offramp.OffRampGasLimitOverride{
			{
				{
					ReceiverExecutionGasLimit: big.NewInt(200_000),
					TokenGasOverrides:         nil,
				},
			},
		},
	)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[destChain], tx, err)
	require.NoError(t, err, "failed to send/confirm manuallyExecute tx")

	newExecutionState, err := state.Chains[destChain].OffRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, sourceChain, out.msgSentEvent.SequenceNumber)
	require.NoError(t, err)
	require.Equal(t, uint8(ccdeploy.EXECUTION_STATE_SUCCESS), newExecutionState)
}

func getMerkleRoot(
	ctx context.Context,
	t *testing.T,
	offRamp *offramp.OffRamp,
	seqNr,
	startBlock uint64,
) (merkleRoot [32]byte) {
	iter, err := offRamp.FilterCommitReportAccepted(&bind.FilterOpts{
		Context: ctx,
		Start:   startBlock,
	})
	require.NoError(t, err)
	for iter.Next() {
		for _, mr := range iter.Event.MerkleRoots {
			if mr.MinSeqNr >= seqNr || mr.MaxSeqNr <= seqNr {
				return mr.MerkleRoot
			}
		}
	}
	require.Fail(
		t,
		fmt.Sprintf("no merkle root found for seq nr %d", seqNr),
	)
	return merkleRoot
}

func getMessageHash(
	ctx context.Context,
	t *testing.T,
	offRamp *offramp.OffRamp,
	sourceChainSelector,
	seqNr uint64,
	msgID [32]byte,
	startBlock uint64,
) (messageHash [32]byte) {
	iter, err := offRamp.FilterExecutionStateChanged(
		&bind.FilterOpts{
			Context: ctx,
			Start:   startBlock,
		},
		[]uint64{sourceChainSelector},
		[]uint64{seqNr},
		[][32]byte{msgID},
	)
	require.NoError(t, err)
	require.True(t, iter.Next())
	require.Equal(t, sourceChainSelector, iter.Event.SourceChainSelector)
	require.Equal(t, seqNr, iter.Event.SequenceNumber)
	require.Equal(t, msgID, iter.Event.MessageId)

	return iter.Event.MessageHash
}

func sleepAndReplay(t *testing.T, e ccdeploy.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	ccdeploy.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}

func runMessagingTestCase(
	tc messagingTestCase,
	receiver common.Address,
	msgData []byte,
	extraArgs []byte,
	expectedExecutionState int,
	extraAssertions ...func(t *testing.T),
) (out messagingTestCaseOutput) {
	// check latest nonce
	latestNonce, err := tc.onchainState.Chains[tc.destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
		Context: tests.Context(tc.t),
	}, tc.sourceChain, tc.sender)
	require.NoError(tc.t, err)
	require.Equal(tc.t, tc.nonce, latestNonce)

	startBlocks := make(map[uint64]*uint64)
	msgSentEvent := ccdeploy.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         msgData,
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    extraArgs,
	})
	expectedSeqNum := make(map[uint64]uint64)
	expectedSeqNum[tc.destChain] = msgSentEvent.SequenceNumber
	out.msgSentEvent = msgSentEvent

	// hack
	if !tc.replayed {
		sleepAndReplay(tc.t, tc.deployedEnv, tc.sourceChain, tc.destChain)
		out.replayed = true
	}

	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	execStates := ccdeploy.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)

	require.Equalf(
		tc.t,
		expectedExecutionState,
		execStates[msgSentEvent.SequenceNumber],
		"wrong execution state for seq nr %d, expected %d, got %d",
		msgSentEvent.SequenceNumber,
		expectedExecutionState,
		execStates[msgSentEvent.SequenceNumber],
	)

	// check the sender latestNonce on the dest, should be incremented
	latestNonce, err = tc.onchainState.Chains[tc.destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
		Context: tests.Context(tc.t),
	}, tc.sourceChain, tc.sender)
	require.NoError(tc.t, err)
	require.Equal(tc.t, tc.nonce+1, latestNonce)
	out.nonce = latestNonce
	tc.t.Logf("confirmed nonce bump for sender %x, latestNonce %d", tc.sender, latestNonce)

	for _, assertion := range extraAssertions {
		assertion(tc.t)
	}

	return
}

// boolsToBitFlags transforms a list of boolean flags to a *big.Int encoded number.
func boolsToBitFlags(bools []bool) *big.Int {
	encodedFlags := big.NewInt(0)
	for i := 0; i < len(bools); i++ {
		if bools[i] {
			encodedFlags.SetBit(encodedFlags, i, 1)
		}
	}
	return encodedFlags
}
