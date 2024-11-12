package changeset

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_FeeBoosting(t *testing.T) {
	t.Parallel()

	type priceFeedPrices struct {
		linkPrice *big.Int
		wethPrice *big.Int
	}

	type scenario struct {
		name            string
		initialPrices   ccipdeployment.InitialPrices
		priceFeedPrices priceFeedPrices
	}

	setupTest := func(t *testing.T, initialPrices ccipdeployment.InitialPrices, priceFeedPrices priceFeedPrices) (
		ccipdeployment.DeployedEnv,
		ccipdeployment.CCIPOnChainState,
		uint64,
		uint64,
	) {
		e := ccipdeployment.NewMemoryEnvironmentWithJobsAndPrices(t, logger.TestLogger(t), 2, 4,
			priceFeedPrices.linkPrice, priceFeedPrices.wethPrice)
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

		require.NoError(t, ccipdeployment.AddLane(e.Env, state, sourceChain, destChain, initialPrices))

		startGasPriceTicker(state, sourceChain, destChain)
		startGasPriceTicker(state, destChain, sourceChain)

		return e, state, sourceChain, destChain
	}

	runScenario := func(t *testing.T, s scenario) {
		e, state, sourceChain, destChain := setupTest(t, s.initialPrices, s.priceFeedPrices)

		// Send message that should initially be too costly
		seqNum := ccipdeployment.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			Data:         []byte("message that needs fee boosting"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})

		sleepAndReplay(t, e, sourceChain, destChain)

		expectedSeqNum := make(map[uint64]uint64)
		expectedSeqNum[destChain] = seqNum
		startBlocks := make(map[uint64]*uint64)

		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, e.Env, state, expectedSeqNum, startBlocks)
	}

	scenarios := []scenario{
		// {
		// 	name: "boost needed due to WETH price increase",
		// 	initialPrices: ccipdeployment.InitialPrices{
		// 		LinkPrice: deployment.E18Mult(5),
		// 		WethPrice: deployment.E18Mult(9),
		// 		GasPrice:  big.NewInt(1.8e11),
		// 	},
		// 	priceFeedPrices: priceFeedPrices{
		// 		linkPrice: deployment.E18Mult(5),
		// 		wethPrice: big.NewInt(20.1e8),
		// 	},
		// },
		// {
		// 	name: "boost needed due to LINK price decrease",
		// 	initialPrices: ccipdeployment.InitialPrices{
		// 		LinkPrice: deployment.E18Mult(5),
		// 		WethPrice: deployment.E18Mult(9),
		// 		GasPrice:  big.NewInt(1.8e11),
		// 	},
		// 	priceFeedPrices: priceFeedPrices{
		// 		linkPrice: big.NewInt(2.24e18),
		// 		wethPrice: big.NewInt(9e8),
		// 	},
		// },
		{
			name: "boost needed due to gas price increase",
			initialPrices: ccipdeployment.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  toPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
				// GasPrice: new(big.Int).Or(
				// 	new(big.Int).Lsh(
				// 		big.NewInt(1.8e11), // 1.6e11 vs 1.8e11 as the gas is 20gwei with 1ETH = 9$
				// 		112,
				// 	),
				// 	big.NewInt(0),
				// ),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: deployment.E18Mult(5),
				wethPrice: big.NewInt(9e8),
			},
		},
	}

	for _, s := range scenarios {
		s := s // capture range variable
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			runScenario(t, s)
		})
	}
}

func toPackedFee(execFee, daFee *big.Int) *big.Int {
	daShifted := new(big.Int).Lsh(daFee, 112)
	return new(big.Int).Or(daShifted, execFee)
}

func sleepAndReplay(t *testing.T, e ccipdeployment.DeployedEnv, sourceChain, destChain uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}

func startGasPriceTicker(state ccipdeployment.CCIPOnChainState, src uint64, dest uint64) {
	fq := state.Chains[src].FeeQuoter
	gasPrice, err := fq.GetDestinationChainGasPrice(&bind.CallOpts{Context: context.Background()}, dest)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get gas price"))
	}
	fmt.Println("Gas price for chain", dest, "is", gasPrice)
}
