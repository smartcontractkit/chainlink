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
	)

	t.Run("data message to eoa", func(t *testing.T) {
		// check latest nonce
		latestNonce, err := state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		require.Equal(t, nonce, latestNonce)

		startBlocks := make(map[uint64]*uint64)
		seqNum := ccipdeployment.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(common.HexToAddress("0x1234").Bytes(), 32),
			Data:         []byte("hello 1234"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		expectedSeqNum := make(map[uint64]uint64)
		expectedSeqNum[destChain] = seqNum

		// hack
		sleepAndReplay(t, e, sourceChain, destChain)
		replayed = true

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)

		// check the sender latestNonce on the dest, should be incremented
		latestNonce, err = state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		require.Equal(t, nonce+1, latestNonce)
		nonce = latestNonce
		t.Logf("confirmed nonce bump for sender %x, latestNonce %d", sender, latestNonce)
	})

	t.Run("message to contract not implementing CCIPReceiver", func(t *testing.T) {
		// check latest nonce
		latestNonce, err := state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		// nonce will be zero if previous test case isn't run.
		// but should be 1 if previous test case is run.
		require.Equal(t, nonce, latestNonce)

		startBlocks := make(map[uint64]*uint64)
		seqNum := ccipdeployment.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
			// FeeQuoter doesn't implement CCIPReceiver
			Receiver:     common.LeftPadBytes(state.Chains[destChain].FeeQuoter.Address().Bytes(), 32),
			Data:         []byte("hello FeeQuoter"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		expectedSeqNum := make(map[uint64]uint64)
		expectedSeqNum[destChain] = seqNum

		if !replayed {
			sleepAndReplay(t, e, sourceChain, destChain)
			replayed = true
		}

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)

		// check the sender latestNonce on the dest, should be incremented
		latestNonce, err = state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		require.Equal(t, nonce+1, latestNonce)
		nonce = latestNonce
		t.Logf("confirmed nonce bump for sender %x, latestNonce %d", sender, latestNonce)
	})

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		// check latest nonce
		latestNonce, err := state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		// nonce will be zero if previous test case isn't run.
		// but should be 1 if previous test case is run.
		require.Equal(t, nonce, latestNonce)

		startBlocks := make(map[uint64]*uint64)
		ccipReceiver := state.Chains[destChain].Receiver
		seqNum := ccipdeployment.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(ccipReceiver.Address().Bytes(), 32),
			Data:         []byte("hello real CCIPReceiver"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		expectedSeqNum := make(map[uint64]uint64)
		expectedSeqNum[destChain] = seqNum

		if !replayed {
			sleepAndReplay(t, e, sourceChain, destChain)
			replayed = true
		}

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)

		iter, err := ccipReceiver.FilterMessageReceived(nil)
		require.NoError(t, err)
		require.True(t, iter.Next())
		// MessageReceived doesn't emit the data unfortunately, so can't check that.

		// check the sender latestNonce on the dest, should be incremented
		latestNonce, err = state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		require.Equal(t, nonce+1, latestNonce)
		nonce = latestNonce
		t.Logf("confirmed nonce bump for sender %x, latestNonce %d", sender, latestNonce)
	})

	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {
		// check latest nonce
		latestNonce, err := state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		// nonce will be zero if previous test case isn't run.
		// but should be 1 if previous test case is run.
		require.Equal(t, nonce, latestNonce)

		// headerBefore, err := e.Env.Chains[destChain].Client.HeaderByNumber(tests.Context(t), nil)
		// require.NoError(t, err)

		startBlocks := make(map[uint64]*uint64)
		ccipReceiver := state.Chains[destChain].Receiver
		seqNum := ccipdeployment.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(ccipReceiver.Address().Bytes(), 32),
			Data:         []byte("hello real CCIPReceiver"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    ccipdeployment.MakeExtraArgsV2(1, false),
		})
		expectedSeqNum := make(map[uint64]uint64)
		expectedSeqNum[destChain] = seqNum

		if !replayed {
			sleepAndReplay(t, e, sourceChain, destChain)
			replayed = true
		}

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)

		// Message would not be emitted, not enough gas to emit log.
		// TODO: this is still returning a log, probably the older one since FAILURE is the execution state.
		// Not enough ctx in the message received log to confirm that it's from another test.
		// iter, err := ccipReceiver.FilterMessageReceived(&bind.FilterOpts{
		// 	Start: headerBefore.Number.Uint64(),
		// })
		// require.NoError(t, err)
		// require.False(t, iter.Next(), "MessageReceived should not be emitted in this test case since gas is too low")

		// check the sender latestNonce on the dest, should be incremented
		latestNonce, err = state.Chains[destChain].NonceManager.GetInboundNonce(&bind.CallOpts{
			Context: tests.Context(t),
		}, sourceChain, sender)
		require.NoError(t, err)
		require.Equal(t, nonce+1, latestNonce)
		nonce = latestNonce
		t.Logf("confirmed nonce bump for sender %x, latestNonce %d", sender, latestNonce)

		// TODO: manually exec?
	})
}

func sleepAndReplay(t *testing.T, e ccipdeployment.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}
