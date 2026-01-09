package messagingtest

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// Use this when testhelpers.DeployedEnv is available (usually in ephemeral test environments).
func NewTestSetupWithDeployedEnv(
	t *testing.T,
	depEnv testhelpers.DeployedEnv,
	onchainState stateview.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	sender []byte,
	testRouter bool,
) TestSetup {
	return TestSetup{
		T:            t,
		Sender:       sender,
		Env:          depEnv.Env,
		DeployedEnv:  depEnv,
		OnchainState: onchainState,
		SourceChain:  sourceChain,
		DestChain:    destChain,
		TestRouter:   testRouter,
	}
}

// Use this when testhelpers.DeployedEnv is not available (usually in long-running test environments like staging).
func NewTestSetup(
	t *testing.T,
	env cldf.Environment,
	onchainState stateview.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	sender []byte,
	testRouter bool,
) TestSetup {
	return TestSetup{
		T:      t,
		Sender: sender,
		Env:    env,
		// no DeployedEnv
		OnchainState: onchainState,
		SourceChain:  sourceChain,
		DestChain:    destChain,
		TestRouter:   testRouter,
	}
}

type TestSetup struct {
	T            *testing.T
	Sender       []byte
	Env          cldf.Environment
	DeployedEnv  testhelpers.DeployedEnv
	OnchainState stateview.CCIPOnChainState
	SourceChain  uint64
	DestChain    uint64
	TestRouter   bool
}

type TestCase struct {
	TestSetup
	ValidationType         ValidationType
	Replayed               bool
	Nonce                  *uint64
	Receiver               []byte
	MsgData                []byte
	ExtraArgs              []byte
	FeeToken               string
	ExpectedExecutionState int
	ExpRevertOnSource      bool
	ExtraAssertions        []func(t *testing.T)
	NumberOfMessages       int // number of messages to send, use same data and extraArgs
	UseMulticall3          bool
}

type ValidationType int

const (
	ValidationTypeNone ValidationType = iota
	ValidationTypeCommit
	ValidationTypeExec // will validate both commit and exec
)

type TestCaseOutput struct {
	Replayed         bool
	Nonce            uint64
	MsgSentEvent     *ccipclient.AnyMsgSentEvent
	AllMsgSentEvents []*ccipclient.AnyMsgSentEvent
}

func getLatestNonce(tc TestCase) uint64 {
	family, err := chain_selectors.GetSelectorFamily(tc.DestChain)
	require.NoError(tc.T, err)

	var latestNonce uint64
	switch family {
	case chain_selectors.FamilyTon:
		// No nonce management on Ton, just return the test value
		return *tc.Nonce
	}

	destAdapter := tc.DeployedEnv.Adapters[tc.DestChain]
	latestNonce, err = destAdapter.GetInboundNonce(tc.T.Context(), tc.Sender, tc.SourceChain)
	require.NoError(tc.T, err)
	return latestNonce
}

// Run runs a messaging test case.
func Run(t *testing.T, tc TestCase) (out TestCaseOutput) {
	// we need to reference the inner testing.T inside a t.Run
	tc.T = t

	var startBlock *uint64

	sourceFamily, err := chain_selectors.GetSelectorFamily(tc.SourceChain)
	require.NoError(tc.T, err)

	sourceAdapter, ok := tc.DeployedEnv.Adapters[tc.SourceChain]
	if !ok {
		tc.T.Errorf("unsupported source chain: %v", tc.SourceChain)
	}
	destAdapter, ok := tc.DeployedEnv.Adapters[tc.DestChain]
	if !ok {
		tc.T.Errorf("unsupported dest chain: %v", tc.DestChain)
	}

	msg, err := sourceAdapter.BuildMessage(testhelpers.MessageComponents{
		DestChainSelector: tc.DestChain,
		Receiver:          tc.Receiver,
		Data:              tc.MsgData,
		FeeToken:          tc.FeeToken,
		ExtraArgs:         tc.ExtraArgs,
		TokenAmounts:      []testhelpers.TokenAmount{},
	})
	require.NoError(tc.T, err)

	if tc.NumberOfMessages == 0 {
		tc.NumberOfMessages = 1 // default to sending one message if not specified
	}

	expectedSeqNumRange := ccipocr3.SeqNumRange{}
	expectedSeqNumExec := []uint64{}
	msgSentEvents := make([]*ccipclient.AnyMsgSentEvent, tc.NumberOfMessages)

	// send all messages first, then validate them
	if tc.UseMulticall3 {
		require.Equal(t, chain_selectors.FamilyEVM, sourceFamily, "only EVM family supported for multicall3 usage")

		msgEVM2Any, ok := msg.(router.ClientEVM2AnyMessage)
		require.True(tc.T, ok, "expected EVM message type")

		onRamp := tc.OnchainState.MustGetEVMChainState(tc.SourceChain).OnRamp
		nextSeqNum, err := onRamp.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: tc.T.Context()}, tc.DestChain)
		require.NoError(tc.T, err)

		calls, totalValue, err := testhelpers.GenMessagesForMulticall3(
			tc.T.Context(),
			tc.OnchainState.MustGetEVMChainState(tc.SourceChain).Router,
			tc.DestChain,
			tc.NumberOfMessages,
			msgEVM2Any,
		)
		require.NoError(tc.T, err)

		sender := tc.Env.BlockChains.EVMChains()[tc.SourceChain].DeployerKey
		currBalance, err := tc.Env.BlockChains.EVMChains()[tc.SourceChain].Client.BalanceAt(tc.T.Context(), sender.From, nil)
		require.NoError(tc.T, err)
		//nolint:testifylint // incorrect lint, GreaterOrEqual can't be used with *big.Int.
		require.True(tc.T, currBalance.Cmp(totalValue) >= 0, "sender balance should be greater than or equal to total value")

		tx, err := tc.OnchainState.MustGetEVMChainState(tc.SourceChain).Multicall3.Aggregate3Value(
			&bind.TransactOpts{
				From:   sender.From,
				Signer: sender.Signer,
				Value:  totalValue,
			},
			calls,
		)
		require.NoError(tc.T, err)

		_, err = cldf.ConfirmIfNoError(tc.Env.BlockChains.EVMChains()[tc.SourceChain], tx, err)
		require.NoError(tc.T, err)

		// check that the message was emitted
		var expectedSeqNums []uint64
		for i := range tc.NumberOfMessages {
			//nolint:gosec // disable G115
			expectedSeqNums = append(expectedSeqNums, nextSeqNum+uint64(i))
		}
		iter, err := tc.OnchainState.MustGetEVMChainState(tc.SourceChain).OnRamp.FilterCCIPMessageSent(
			nil, []uint64{tc.DestChain}, expectedSeqNums,
		)
		require.NoError(tc.T, err)

		for i := 0; i < tc.NumberOfMessages; i++ {
			require.True(tc.T, iter.Next())
			msgSentEvents[i] = &ccipclient.AnyMsgSentEvent{
				SequenceNumber: iter.Event.SequenceNumber,
				RawEvent:       iter.Event,
			}
		}
		out.MsgSentEvent = msgSentEvents[len(msgSentEvents)-1]

		// Expect a single root with this sequence number range.
		expectedSeqNumRange = ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(msgSentEvents[0].SequenceNumber),
			ccipocr3.SeqNum(msgSentEvents[len(msgSentEvents)-1].SequenceNumber),
		}
		// Expect all messages to be executed.
		for i := range msgSentEvents {
			expectedSeqNumExec = append(expectedSeqNumExec, msgSentEvents[i].SequenceNumber)
		}
	} else {
		// Send sequentially
		for i := 0; i < tc.NumberOfMessages; i++ {
			msgSentEventLocal := testhelpers.TestSendRequest(
				tc.T,
				tc.Env,
				tc.OnchainState,
				tc.SourceChain,
				tc.DestChain,
				tc.TestRouter,
				msg)

			if i == 0 {
				expectedSeqNumRange = ccipocr3.SeqNumRange{ccipocr3.SeqNum(msgSentEventLocal.SequenceNumber)}
			}
			expectedSeqNumRange = ccipocr3.SeqNumRange{expectedSeqNumRange.Start(), ccipocr3.SeqNum(msgSentEventLocal.SequenceNumber)}

			expectedSeqNumExec = append(expectedSeqNumExec, msgSentEventLocal.SequenceNumber)
			// TODO: If this feature is needed more we can refactor the function to return a slice of events
			// return only last msg event
			out.MsgSentEvent = msgSentEventLocal
			msgSentEvents[i] = msgSentEventLocal
		}
	}
	// return all message sent events.
	out.AllMsgSentEvents = msgSentEvents

	// HACK: if the node booted or the logpoller filters got registered after ccipSend,
	// we need to replay missed logs
	if !tc.Replayed {
		require.NotNil(tc.T, tc.DeployedEnv)
		testhelpers.SleepAndReplay(tc.T, tc.DeployedEnv.Env, 30*time.Second, tc.SourceChain, tc.DestChain)
		out.Replayed = true
	}

	// Perform validation based on ValidationType
	switch tc.ValidationType {
	case ValidationTypeCommit:
		commitStart := time.Now()
		destAdapter.ValidateCommit(t, tc.SourceChain, startBlock, expectedSeqNumRange)
		tc.T.Logf("confirmed commit of seq nums %+v in %s", expectedSeqNumRange, time.Since(commitStart).String())
		// Explicitly log that only commit was validated if only Commit was requested
		tc.T.Logf("only commit validation was performed")
	case ValidationTypeExec: // will validate both commit and exec
		// First, validate commit
		commitStart := time.Now()
		destAdapter.ValidateCommit(t, tc.SourceChain, startBlock, expectedSeqNumRange)
		tc.T.Logf("confirmed commit of seq nums %+v in %s", expectedSeqNumRange, time.Since(commitStart).String())

		// Then, validate execution
		execStart := time.Now()
		execStates := destAdapter.ValidateExec(t, tc.SourceChain, startBlock, expectedSeqNumExec)
		tc.T.Logf("confirmed exec of seq nums %+v in %s", expectedSeqNumExec, time.Since(execStart).String())

		for _, msgSentEvent := range msgSentEvents {
			require.Equalf(
				tc.T,
				tc.ExpectedExecutionState,
				execStates[msgSentEvent.SequenceNumber],
				"wrong execution state for seq nr %d, expected %d, got %d",
				msgSentEvent.SequenceNumber,
				tc.ExpectedExecutionState,
				execStates[msgSentEvent.SequenceNumber],
			)
		}

		family, err := chain_selectors.GetSelectorFamily(tc.DestChain)
		require.NoError(tc.T, err)

		unorderedExec := true
		// Only EVM (and by extension Tron) support ordered execution
		switch family {
		case chain_selectors.FamilyEVM:
			unorderedExec = false
		case chain_selectors.FamilyTron:
			unorderedExec = false
		}

		if !unorderedExec {
			latestNonce := getLatestNonce(tc)
			// Check if Nonce is non-nil before comparing. Nonce check only makes sense if it was explicitly provided.
			if tc.Nonce != nil {
				require.Equal(tc.T, *tc.Nonce+1, latestNonce)
				out.Nonce = latestNonce
				tc.T.Logf("confirmed nonce bump for sender %x, expected %d, got latestNonce %d", tc.Sender, *tc.Nonce+1, latestNonce)
			} else {
				tc.T.Logf("skipping nonce bump check for sender %x as initial nonce was nil, latestNonce %d", tc.Sender, latestNonce)
			}
		}

		for _, assertion := range tc.ExtraAssertions {
			assertion(tc.T)
		}

	case ValidationTypeNone:
		tc.T.Logf("skipping validation of sent message")
	}
	return
}
