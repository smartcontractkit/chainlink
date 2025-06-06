package messagingtest

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
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
	ExtraAssertions        []func(t *testing.T)
	NumberOfMessages       int // number of messages to send, use same data and extraArgs
}

type ValidationType int

const (
	ValidationTypeNone ValidationType = iota
	ValidationTypeCommit
	ValidationTypeExec // will validate both commit and exec
)

type TestCaseOutput struct {
	Replayed     bool
	Nonce        uint64
	MsgSentEvent *onramp.OnRampCCIPMessageSent
}

func getLatestNonce(tc TestCase) uint64 {
	family, err := chain_selectors.GetSelectorFamily(tc.DestChain)
	require.NoError(tc.T, err)

	var latestNonce uint64
	switch family {
	case chain_selectors.FamilyEVM:
		latestNonce, err = tc.OnchainState.Chains[tc.DestChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tc.T.Context(),
		}, tc.SourceChain, tc.Sender)
		require.NoError(tc.T, err)
	case chain_selectors.FamilySolana:
		ctx := tc.T.Context()
		client := tc.Env.BlockChains.SolanaChains()[tc.DestChain].Client
		// TODO: solcommon.FindNoncePDA expected the sender to be a solana pubkey
		chainSelectorLE := solcommon.Uint64ToLE(tc.DestChain)
		noncePDA, _, err := solana.FindProgramAddress([][]byte{[]byte("nonce"), chainSelectorLE, tc.Sender}, tc.OnchainState.SolChains[tc.DestChain].Router)
		require.NoError(tc.T, err)
		var nonceCounterAccount ccip_router.Nonce
		// we ignore the error because the account might not exist yet
		_ = solcommon.GetAccountDataBorshInto(ctx, client, noncePDA, solconfig.DefaultCommitment, &nonceCounterAccount)
		latestNonce = nonceCounterAccount.Counter
	}
	return latestNonce
}

// Run runs a messaging test case.
func Run(t *testing.T, tc TestCase) (out TestCaseOutput) {
	// we need to reference the inner testing.T inside a t.Run
	tc.T = t

	startBlocks := make(map[uint64]*uint64)

	family, err := chain_selectors.GetSelectorFamily(tc.SourceChain)
	require.NoError(tc.T, err)

	var msg any
	switch family {
	case chain_selectors.FamilyEVM:
		feeToken := common.HexToAddress("0x0")
		if len(tc.FeeToken) > 0 {
			feeToken = common.HexToAddress(tc.FeeToken)
		}

		msg = router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(tc.Receiver, 32),
			Data:         tc.MsgData,
			TokenAmounts: nil,
			FeeToken:     feeToken,
			ExtraArgs:    tc.ExtraArgs,
		}
	case chain_selectors.FamilySolana:
		feeToken := solana.PublicKey{}
		if len(tc.FeeToken) > 0 {
			feeToken, err = solana.PublicKeyFromBase58(tc.FeeToken)
			require.NoError(t, err)
		}

		msg = ccip_router.SVM2AnyMessage{
			Receiver:     common.LeftPadBytes(tc.Receiver, 32),
			TokenAmounts: nil,
			Data:         tc.MsgData,
			FeeToken:     feeToken,
			ExtraArgs:    tc.ExtraArgs,
		}

	default:
		tc.T.Errorf("unsupported source chain: %v", family)
	}

	if tc.NumberOfMessages == 0 {
		tc.NumberOfMessages = 1 // default to sending one message if not specified
	}

	expectedSeqNumRange := map[testhelpers.SourceDestPair]ccipocr3.SeqNumRange{}
	expectedSeqNumExec := map[testhelpers.SourceDestPair][]uint64{}
	msgSentEvents := make([]*onramp.OnRampCCIPMessageSent, tc.NumberOfMessages)
	sourceDest := testhelpers.SourceDestPair{
		SourceChainSelector: tc.SourceChain,
		DestChainSelector:   tc.DestChain,
	}

	// send all messages first, then validate them
	for i := 0; i < tc.NumberOfMessages; i++ {
		msgSentEventLocal := testhelpers.TestSendRequest(
			tc.T,
			tc.Env,
			tc.OnchainState,
			tc.SourceChain,
			tc.DestChain,
			tc.TestRouter,
			msg)

		_, ok := expectedSeqNumRange[sourceDest]
		if !ok {
			expectedSeqNumRange[sourceDest] = ccipocr3.SeqNumRange{ccipocr3.SeqNum(msgSentEventLocal.SequenceNumber)}
		}
		expectedSeqNumRange[sourceDest] = ccipocr3.SeqNumRange{expectedSeqNumRange[sourceDest].Start(),
			ccipocr3.SeqNum(msgSentEventLocal.SequenceNumber)}

		expectedSeqNumExec[sourceDest] = append(expectedSeqNumExec[sourceDest], msgSentEventLocal.SequenceNumber)
		// TODO: If this feature is needed more we can refactor the function to return a slice of events
		// return only last msg event
		out.MsgSentEvent = msgSentEventLocal
		msgSentEvents[i] = msgSentEventLocal
	}

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
		testhelpers.ConfirmCommitForAllWithExpectedSeqNums(tc.T, tc.Env, tc.OnchainState, expectedSeqNumRange, startBlocks)
		tc.T.Logf("confirmed commit of seq nums %+v in %s", expectedSeqNumRange, time.Since(commitStart).String())
		// Explicitly log that only commit was validated if only Commit was requested
		tc.T.Logf("only commit validation was performed")

	case ValidationTypeExec: // will validate both commit and exec
		// First, validate commit
		commitStart := time.Now()
		testhelpers.ConfirmCommitForAllWithExpectedSeqNums(tc.T, tc.Env, tc.OnchainState, expectedSeqNumRange, startBlocks)
		tc.T.Logf("confirmed commit of seq nums %+v in %s", expectedSeqNumRange, time.Since(commitStart).String())

		// Then, validate execution
		execStart := time.Now()
		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(tc.T, tc.Env, tc.OnchainState, expectedSeqNumExec, startBlocks)
		tc.T.Logf("confirmed exec of seq nums %+v in %s", expectedSeqNumExec, time.Since(execStart).String())

		for _, msgSentEvent := range msgSentEvents {
			require.Equalf(
				tc.T,
				tc.ExpectedExecutionState,
				execStates[sourceDest][msgSentEvent.SequenceNumber],
				"wrong execution state for seq nr %d, expected %d, got %d",
				msgSentEvent.SequenceNumber,
				tc.ExpectedExecutionState,
				execStates[sourceDest][msgSentEvent.SequenceNumber],
			)
		}

		family, err := chain_selectors.GetSelectorFamily(tc.DestChain)
		require.NoError(tc.T, err)

		// Solana doesn't support catching CPI errors, so nonces can't be ordered
		unorderedExec := family == chain_selectors.FamilySolana

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
