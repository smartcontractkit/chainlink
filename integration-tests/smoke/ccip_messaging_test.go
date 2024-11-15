package smoke

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

func Test_CCIPMessaging(t *testing.T) {
	// Setup 2 chains and a single lane.
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	e, _, _ := testsetups.NewLocalDevEnvironment(t, lggr)

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

	tokenConfig := ccipdeployment.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
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
			nil,                                    // default extraArgs
			ccipdeployment.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
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
			nil,                                    // default extraArgs
			ccipdeployment.EXECUTION_STATE_SUCCESS, // success because offRamp won't call a contract not implementing CCIPReceiver
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
			nil, // default extraArgs
			ccipdeployment.EXECUTION_STATE_SUCCESS,
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
			ccipdeployment.MakeEVMExtraArgsV2(1, false), // 1 gas is too low.
			ccipdeployment.EXECUTION_STATE_FAILURE,      // state would be failed onchain due to low gas
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
	seqNum := ccipdeployment.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         msgData,
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    extraArgs,
	})
	expectedSeqNum := make(map[uint64]uint64)
	expectedSeqNum[tc.destChain] = seqNum

	// hack
	if !tc.replayed {
		sleepAndReplay(tc.t, tc.deployedEnv, tc.sourceChain, tc.destChain)
		out.replayed = true
	}

	ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	execStates := ccipdeployment.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)

	require.Equalf(
		tc.t,
		expectedExecutionState,
		execStates[seqNum],
		"wrong execution state for seq nr %d, expected %d, got %d",
		seqNum,
		expectedExecutionState,
		execStates[seqNum],
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
