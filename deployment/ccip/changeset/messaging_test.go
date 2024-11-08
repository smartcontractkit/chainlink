package changeset

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"
)

func Test_Messaging(t *testing.T) {
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

	t.Run("data message to eoa", func(t *testing.T) {
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
		time.Sleep(30 * time.Second)
		replayBlocks := make(map[uint64]uint64)
		replayBlocks[sourceChain] = 1
		replayBlocks[destChain] = 1
		ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)
	})

	// TODO
	t.Run("message to contract not implementing CCIPReceiver", func(t *testing.T) {})

	// TODO
	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {})

	// TODO
	t.Run("message to contract implementing CCIPReceiver with low exec gas", func(t *testing.T) {})
}
