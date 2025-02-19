package messagingtest

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
)

// Use this when testhelpers.DeployedEnv is available (usually in ephemeral test environments).
func NewTestSetupWithDeployedEnv(
	t *testing.T,
	depEnv testhelpers.DeployedEnv,
	onchainState changeset.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	sender []byte,
	testRouter,
	validateResp bool,
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
		ValidateResp: validateResp,
	}
}

// Use this when testhelpers.DeployedEnv is not available (usually in long-running test environments like staging).
func NewTestSetup(
	t *testing.T,
	env deployment.Environment,
	onchainState changeset.CCIPOnChainState,
	sourceChain,
	destChain uint64,
	sender []byte,
	testRouter,
	validateResp bool,
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
		ValidateResp: validateResp,
	}
}

type TestSetup struct {
	T            *testing.T
	Sender       []byte
	Env          deployment.Environment
	DeployedEnv  testhelpers.DeployedEnv
	OnchainState changeset.CCIPOnChainState
	SourceChain  uint64
	DestChain    uint64
	TestRouter   bool
	ValidateResp bool
}

type TestCase struct {
	TestSetup
	Replayed               bool
	Nonce                  uint64
	Receiver               common.Address
	MsgData                []byte
	ExtraArgs              []byte
	ExpectedExecutionState int
	ExtraAssertions        []func(t *testing.T)
}

type TestCaseOutput struct {
	Replayed     bool
	Nonce        uint64
	MsgSentEvent *onramp.OnRampCCIPMessageSent
}

func sleepAndReplay(t *testing.T, e testhelpers.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1

	testhelpers.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}

func getLatestNonce(tc TestCase) uint64 {
	family, err := chain_selectors.GetSelectorFamily(tc.DestChain)
	require.NoError(tc.T, err)

	var latestNonce uint64
	switch family {
	case chain_selectors.FamilyEVM:
		latestNonce, err = tc.OnchainState.Chains[tc.DestChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(tc.T),
		}, tc.SourceChain, tc.Sender)
		require.NoError(tc.T, err)
	case chain_selectors.FamilySolana:
		// var nonceCounterAccount ccip_router.Nonce
		// err = common.GetAccountDataBorshInto(ctx, solanaGoClient, nonceEvmPDA, config.DefaultCommitment, &nonceCounterAccount)
		// require.NoError(t, err, "failed to get account info")
		// require.Equal(t, uint64(1), nonceCounterAccount.Counter)
	}
	return latestNonce
}

// Run runs a messaging test case.
func Run(tc TestCase) (out TestCaseOutput) {
	if tc.ValidateResp {
		// check latest nonce
		latestNonce := getLatestNonce(tc)
		require.Equal(tc.T, tc.Nonce, latestNonce)
	}

	startBlocks := make(map[uint64]*uint64)
	msgSentEvent := testhelpers.TestSendRequest(
		tc.T,
		tc.Env,
		tc.OnchainState,
		tc.SourceChain,
		tc.DestChain,
		tc.TestRouter,
		router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(tc.Receiver.Bytes(), 32),
			Data:         tc.MsgData,
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    tc.ExtraArgs,
		})
	expectedSeqNum := map[testhelpers.SourceDestPair]uint64{
		{
			SourceChainSelector: tc.SourceChain,
			DestChainSelector:   tc.DestChain,
		}: msgSentEvent.SequenceNumber,
	}
	expectedSeqNumExec := map[testhelpers.SourceDestPair][]uint64{
		{
			SourceChainSelector: tc.SourceChain,
			DestChainSelector:   tc.DestChain,
		}: {msgSentEvent.SequenceNumber},
	}
	out.MsgSentEvent = msgSentEvent

	// hack
	if !tc.Replayed {
		require.NotNil(tc.T, tc.DeployedEnv)
		sleepAndReplay(tc.T, tc.DeployedEnv, tc.SourceChain, tc.DestChain)
		out.Replayed = true
	}

	if tc.ValidateResp {
		commitStart := time.Now()
		testhelpers.ConfirmCommitForAllWithExpectedSeqNums(tc.T, tc.Env, tc.OnchainState, expectedSeqNum, startBlocks)
		tc.T.Logf("confirmed commit of seq nums %+v in %s", expectedSeqNum, time.Since(commitStart).String())
		execStart := time.Now()
		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(tc.T, tc.Env, tc.OnchainState, expectedSeqNumExec, startBlocks)
		tc.T.Logf("confirmed exec of seq nums %+v in %s", expectedSeqNumExec, time.Since(execStart).String())

		require.Equalf(
			tc.T,
			tc.ExpectedExecutionState,
			execStates[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SourceChain,
				DestChainSelector:   tc.DestChain,
			}][msgSentEvent.SequenceNumber],
			"wrong execution state for seq nr %d, expected %d, got %d",
			msgSentEvent.SequenceNumber,
			tc.ExpectedExecutionState,
			execStates[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SourceChain,
				DestChainSelector:   tc.DestChain,
			}][msgSentEvent.SequenceNumber],
		)

		// check the sender latestNonce on the dest, should be incremented
		latestNonce := getLatestNonce(tc)
		require.Equal(tc.T, tc.Nonce+1, latestNonce)
		out.Nonce = latestNonce
		tc.T.Logf("confirmed nonce bump for sender %x, latestNonce %d", tc.Sender, latestNonce)

		for _, assertion := range tc.ExtraAssertions {
			assertion(tc.T)
		}
	} else {
		tc.T.Logf("skipping validation of sent message")
	}

	return
}
