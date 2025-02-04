package messagingtest

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

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
	Receiver               []byte
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

func sleepAndReplay(t *testing.T, e testhelpers.DeployedEnv, chainSelectors ...uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	for _, selector := range chainSelectors {
		family, err := chain_selectors.GetSelectorFamily(selector)
		require.NoError(t, err)
		// log replay is only available on EVM
		if family == chain_selectors.FamilyEVM {
			replayBlocks[selector] = 1
		}
	}
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
		ctx := context.Background()
		client := tc.Env.SolChains[tc.DestChain].Client
		noncePDA, err := solstate.FindNoncePDA(tc.SourceChain, solana.PublicKeyFromBytes(tc.Sender), tc.OnchainState.SolChains[tc.DestChain].Router)
		require.NoError(tc.T, err)
		var nonceCounterAccount ccip_router.Nonce
		err = solcommon.GetAccountDataBorshInto(ctx, client, noncePDA, solconfig.DefaultCommitment, &nonceCounterAccount)
		require.NoError(tc.T, err, "failed to get nonce account info")
		latestNonce = nonceCounterAccount.Counter
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
			Receiver:     common.LeftPadBytes(tc.Receiver, 32),
			Data:         tc.MsgData,
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    tc.ExtraArgs,
		})
	sourceDest := testhelpers.SourceDestPair{
		SourceChainSelector: tc.SourceChain,
		DestChainSelector:   tc.DestChain,
	}
	expectedSeqNum := map[testhelpers.SourceDestPair]uint64{
		sourceDest: msgSentEvent.SequenceNumber,
	}
	expectedSeqNumExec := map[testhelpers.SourceDestPair][]uint64{
		sourceDest: {msgSentEvent.SequenceNumber},
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
			execStates[sourceDest][msgSentEvent.SequenceNumber],
			"wrong execution state for seq nr %d, expected %d, got %d",
			msgSentEvent.SequenceNumber,
			tc.ExpectedExecutionState,
			execStates[sourceDest][msgSentEvent.SequenceNumber],
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
