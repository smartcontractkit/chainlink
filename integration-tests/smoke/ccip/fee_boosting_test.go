package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type feeboostTestCase struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            changeset.DeployedEnv
	onchainState           changeset.CCIPOnChainState
	initialPrices          changeset.InitialPrices
	priceFeedPrices        priceFeedPrices
	sourceChain, destChain uint64
}

type priceFeedPrices struct {
	linkPrice *big.Int
	wethPrice *big.Int
}

// TODO: find a way to reuse the same test setup for all tests
func Test_CCIPFeeBoosting(t *testing.T) {
	setupTestEnv := func(t *testing.T, numChains int) (changeset.DeployedEnv, changeset.CCIPOnChainState, []uint64) {
		e, _, _ := testsetups.NewLocalDevEnvironment(
			t, logger.TestLogger(t),
			deployment.E18Mult(5),
			big.NewInt(9e8))

		state, err := changeset.LoadOnchainState(e.Env)
		require.NoError(t, err)

		allChainSelectors := maps.Keys(e.Env.Chains)
		require.Len(t, allChainSelectors, numChains)
		return e, state, allChainSelectors
	}

	t.Run("boost needed due to WETH price increase (also covering gas price inscrease)", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: changeset.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  changeset.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: deployment.E18Mult(5),
				wethPrice: big.NewInt(9.9e8), // increase from 9e8 to 9.9e8
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})

	t.Run("boost needed due to LINK price decrease", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: changeset.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  changeset.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: big.NewInt(4.5e18), // decrease from 5e18 to 4.5e18
				wethPrice: big.NewInt(9e8),
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})
}

func runFeeboostTestCase(tc feeboostTestCase) {
	require.NoError(tc.t, changeset.AddLane(tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, tc.initialPrices))

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	msgSentEvent := changeset.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(tc.onchainState.Chains[tc.destChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("message that needs fee boosting"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[changeset.SourceDestPair{
		SourceChainSelector: tc.sourceChain,
		DestChainSelector:   tc.destChain,
	}] = msgSentEvent.SequenceNumber

	// hack
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[tc.sourceChain] = 1
	replayBlocks[tc.destChain] = 1
	changeset.ReplayLogs(tc.t, tc.deployedEnv.Env.Offchain, replayBlocks)

	changeset.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	changeset.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
}
