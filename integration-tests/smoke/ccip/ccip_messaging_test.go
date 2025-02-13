package ccip

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
)

func Test_CCIPMessaging(t *testing.T) {
	// Setup 2 chains and a single lane.
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t)

	state, err := changeset.LoadOnchainState(e.Env)
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
			true,  // validateResp
		)
	)

	t.Run("data message to eoa", func(t *testing.T) {
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  nonce,
				Receiver:               common.HexToAddress("0xdead"),
				MsgData:                []byte("hello eoa"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
			},
		)
	})

	t.Run("message to contract not implementing CCIPReceiver", func(t *testing.T) {
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  out.Nonce,
				Receiver:               state.Chains[destChain].FeeQuoter.Address(),
				MsgData:                []byte("hello FeeQuoter"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call a contract not implementing CCIPReceiver
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		latestHead, err := e.Env.Chains[destChain].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  out.Nonce,
				Receiver:               state.Chains[destChain].Receiver.Address(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              nil, // default extraArgs
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
	})

	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {
		latestHead, err := e.Env.Chains[destChain].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               out.Replayed,
				Nonce:                  out.Nonce,
				Receiver:               state.Chains[destChain].Receiver.Address(),
				MsgData:                []byte("hello CCIPReceiver with low exec gas"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(1, false), // 1 gas is too low.
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_FAILURE,      // state would be failed onchain due to low gas
			},
		)

		manuallyExecute(ctx, t, latestHead.Number.Uint64(), state, destChain, out, sourceChain, e, sender)

		t.Logf("successfully manually executed message %x",
			out.MsgSentEvent.Message.Header.MessageId)
	})
}

func manuallyExecute(
	ctx context.Context,
	t *testing.T,
	startBlock uint64,
	state changeset.CCIPOnChainState,
	destChain uint64,
	out mt.TestCaseOutput,
	sourceChain uint64,
	e testhelpers.DeployedEnv,
	sender []byte,
) {
	merkleRoot := getMerkleRoot(
		ctx,
		t,
		state.Chains[destChain].OffRamp,
		out.MsgSentEvent.SequenceNumber,
		startBlock,
	)
	messageHash := getMessageHash(
		ctx,
		t,
		state.Chains[destChain].OffRamp,
		sourceChain,
		out.MsgSentEvent.SequenceNumber,
		out.MsgSentEvent.Message.Header.MessageId,
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
							MessageId:           out.MsgSentEvent.Message.Header.MessageId,
							SourceChainSelector: sourceChain,
							DestChainSelector:   destChain,
							SequenceNumber:      out.MsgSentEvent.SequenceNumber,
							Nonce:               out.MsgSentEvent.Message.Header.Nonce,
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

	newExecutionState, err := state.Chains[destChain].OffRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, sourceChain, out.MsgSentEvent.SequenceNumber)
	require.NoError(t, err)
	require.Equal(t, uint8(testhelpers.EXECUTION_STATE_SUCCESS), newExecutionState)
}

func getMerkleRoot(
	ctx context.Context,
	t *testing.T,
	offRamp offramp.OffRampInterface,
	seqNr,
	startBlock uint64,
) (merkleRoot [32]byte) {
	iter, err := offRamp.FilterCommitReportAccepted(&bind.FilterOpts{
		Context: ctx,
		Start:   startBlock,
	})
	require.NoError(t, err)
	for iter.Next() {
		for _, mr := range iter.Event.BlessedMerkleRoots {
			if mr.MinSeqNr >= seqNr || mr.MaxSeqNr <= seqNr {
				return mr.MerkleRoot
			}
		}
		// todo: dedup
		// ------------------------------
		for _, mr := range iter.Event.UnblessedMerkleRoots {
			if mr.MinSeqNr >= seqNr || mr.MaxSeqNr <= seqNr {
				return mr.MerkleRoot
			}
		}
		// ------------------------------
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
	offRamp offramp.OffRampInterface,
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
