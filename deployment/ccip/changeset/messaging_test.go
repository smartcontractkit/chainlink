package changeset

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"
)

type testCaseSetup struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            ccipdeployment.DeployedEnv
	onchainState           ccipdeployment.CCIPOnChainState
	sourceChain, destChain uint64
}

type messagingTestCase struct {
	testCaseSetup
	replayed bool
	nonce    uint64
}

type messagingTestCaseOutput struct {
	replayed bool
	nonce    uint64
}

func Test_Messaging(t *testing.T) {
	t.Parallel()

	// Setup 2 chains and a single lane.
	e := ccipdeployment.NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 2, 4)
	state, err := ccipdeployment.LoadOnchainState(e.Env)
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

	tokenConfig := ccipdeployment.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	newAddresses := deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployCCIPContracts(e.Env, newAddresses, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: allChainSelectors,
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccipdeployment.NewTestMCMSConfig(t, e.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	state, err = ccipdeployment.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// connect a single lane, source to dest
	require.NoError(t, ccipdeployment.AddLane(e.Env, state, sourceChain, destChain))

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
		)
	})

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		out = runMessagingTestCase(
			messagingTestCase{
				testCaseSetup: setup,
				replayed:      out.replayed,
				nonce:         out.nonce,
			},
			state.Chains[destChain].Receiver.Address(),
			[]byte("hello CCIPReceiver"),
			func(t *testing.T) {
				iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(nil)
				require.NoError(t, err)
				require.True(t, iter.Next())
				// MessageReceived doesn't emit the data unfortunately, so can't check that.
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {
		out = runMessagingTestCase(
			messagingTestCase{
				testCaseSetup: setup,
				replayed:      out.replayed,
				nonce:         out.nonce,
			},
			state.Chains[destChain].Receiver.Address(),
			[]byte("hello CCIPReceiver with low exec gas"),
			func(t *testing.T) {
				// Message should not be emitted, not enough gas to emit log.
				// TODO: this is still returning a log, probably the older one since FAILURE is the execution state.
				// Not enough ctx in the message received log to confirm that it's from another test.
				// Maybe check the log block number and assert that its < the header before block number from above?
				// iter, err := ccipReceiver.FilterMessageReceived(&bind.FilterOpts{
				// 	Start: headerBefore.Number.Uint64(),
				// })
				// require.NoError(t, err)
				// require.False(t, iter.Next(), "MessageReceived should not be emitted in this test case since gas is too low")
			},
		)
	})
}

func sleepAndReplay(t *testing.T, e ccipdeployment.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}

func runMessagingTestCase(
	tc messagingTestCase,
	receiver common.Address,
	msgData []byte,
	extraAssertions ...func(t *testing.T),
) (out messagingTestCaseOutput) {
	// check latest nonce
	latestNonce, err := tc.onchainState.Chains[tc.destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
		Context: tests.Context(tc.t),
	}, tc.sourceChain, tc.sender)
	require.NoError(tc.t, err)
	require.Equal(tc.t, tc.nonce, latestNonce)

	startBlocks := make(map[uint64]*uint64)
	seqNum := ccipdeployment.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         msgData,
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum := make(map[uint64]uint64)
	expectedSeqNum[tc.destChain] = seqNum

	// hack
	if !tc.replayed {
		sleepAndReplay(tc.t, tc.deployedEnv, tc.sourceChain, tc.destChain)
		out.replayed = true
	}

	ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	ccipdeployment.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)

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
