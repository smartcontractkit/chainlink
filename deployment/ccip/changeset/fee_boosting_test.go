package changeset

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type feeboostTestCase struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            ccipdeployment.DeployedEnv
	onchainState           ccipdeployment.CCIPOnChainState
	initialPrices          ccipdeployment.InitialPrices
	priceFeedPrices        priceFeedPrices
	sourceChain, destChain uint64
}

type priceFeedPrices struct {
	linkPrice *big.Int
	wethPrice *big.Int
}

func Test_FeeBoosting(t *testing.T) {
	t.Parallel()

	setupTestEnv := func(t *testing.T, numChains int) (ccipdeployment.DeployedEnv, ccipdeployment.CCIPOnChainState, []uint64) {
		e := ccipdeployment.NewMemoryEnvironmentWithJobsAndPrices(
			t, logger.TestLogger(t),
			numChains, 4,
			deployment.E18Mult(5), big.NewInt(9e8))
		state, err := ccipdeployment.LoadOnchainState(e.Env)
		require.NoError(t, err)

		allChainSelectors := maps.Keys(e.Env.Chains)
		require.Len(t, allChainSelectors, numChains)

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

		return e, state, allChainSelectors
	}

	t.Run("boost needed due to WETH price increase", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: ccipdeployment.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  ccipdeployment.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: deployment.E18Mult(5),
				wethPrice: big.NewInt(9.9e8),
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
			initialPrices: ccipdeployment.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  ccipdeployment.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: big.NewInt(4.5e18),
				wethPrice: big.NewInt(9e8),
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})

	t.Run("boost needed due to gas price increase", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: ccipdeployment.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  ccipdeployment.ToPackedFee(big.NewInt(1.75e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: deployment.E18Mult(5),
				wethPrice: big.NewInt(9e8),
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})
}

func runFeeboostTestCase(tc feeboostTestCase) {
	require.NoError(tc.t, ccipdeployment.AddLane(tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, tc.initialPrices))

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[uint64]uint64)
	seqNum := ccipdeployment.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(tc.onchainState.Chains[tc.destChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("message that needs fee boosting"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[tc.destChain] = seqNum

	sleepAndReplay(tc.t, tc.deployedEnv, tc.sourceChain, tc.destChain)

	ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	ccipdeployment.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
}

func sleepAndReplay(t *testing.T, e ccipdeployment.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}
