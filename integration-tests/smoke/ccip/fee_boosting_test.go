package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type feeboostTestCase struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            changeset.DeployedEnv
	onchainState           changeset.CCIPOnChainState
	priceFeedPrices        priceFeedPrices
	sourceChain, destChain uint64
}

type priceFeedPrices struct {
	linkPrice *big.Int
	wethPrice *big.Int
}

var (
	initialPrices = changeset.InitialPrices{
		LinkPrice: deployment.E18Mult(5),
		WethPrice: deployment.E18Mult(9),
		GasPrice:  changeset.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
	}
)

func Test_CCIPFeeBoosting(t *testing.T) {
	setupTestEnv := func(t *testing.T, numChains int) (changeset.DeployedEnv, changeset.CCIPOnChainState, []uint64) {
		e, _, _ := testsetups.NewLocalDevEnvironment(t, logger.TestLogger(t), deployment.E18Mult(5), deployment.E18Mult(9), nil)

		state, err := changeset.LoadOnchainState(e.Env)
		require.NoError(t, err)

		allChainSelectors := maps.Keys(e.Env.Chains)
		require.Len(t, allChainSelectors, numChains)
		sourceChain := allChainSelectors[0]
		destChain := allChainSelectors[1]
		t.Log("All chain selectors:", allChainSelectors,
			", home chain selector:", e.HomeChainSel,
			", feed chain selector:", e.FeedChainSel,
			", source chain selector:", sourceChain,
			", dest chain selector:", destChain,
		)

		laneCfg := changeset.LaneConfig{
			SourceSelector:        sourceChain,
			DestSelector:          destChain,
			InitialPricesBySource: initialPrices,
			FeeQuoterDestChain:    changeset.DefaultFeeQuoterDestChainConfig(),
		}

		require.NoError(t, changeset.AddLane(e.Env, state, laneCfg, false))
		return e, state, allChainSelectors
	}

	e, state, chains := setupTestEnv(t, 2)

	t.Run("boost needed due to WETH price increase (also covering gas price increase)", func(t *testing.T) {
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			priceFeedPrices: priceFeedPrices{
				wethPrice: new(big.Int).Mul(
					big.NewInt(99),
					new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil),
				), // increase from 9e18 to 9.9e18
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})

	// t.Run("boost needed due to LINK price decrease", func(t *testing.T) {
	// 	runFeeboostTestCase(feeboostTestCase{
	// 		t:            t,
	// 		sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
	// 		deployedEnv:  e,
	// 		onchainState: state,
	// 		priceFeedPrices: priceFeedPrices{
	// 			linkPrice: big.NewInt(4.5e18), // decrease from 5e18 to 4.5e18
	// 		},
	// 		sourceChain: chains[0],
	// 		destChain:   chains[1],
	// 	})
	// })
}

func runFeeboostTestCase(tc feeboostTestCase) {
	// Set initial prices
	// setPrices(tc, initialPrices.LinkPrice, initialPrices.WethPrice)

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)
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
	expectedSeqNumExec[changeset.SourceDestPair{
		SourceChainSelector: tc.sourceChain,
		DestChainSelector:   tc.destChain,
	}] = []uint64{msgSentEvent.SequenceNumber}

	// Update prices
	setPrices(tc, tc.priceFeedPrices.linkPrice, tc.priceFeedPrices.wethPrice)

	require.True(tc.t, willTriggerFeeBoosting(msgSentEvent, tc))

	// hack
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[tc.sourceChain] = 1
	replayBlocks[tc.destChain] = 1
	changeset.ReplayLogs(tc.t, tc.deployedEnv.Env.Offchain, replayBlocks)

	changeset.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	changeset.ConfirmExecWithSeqNrsForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNumExec, startBlocks)
}

func setPrices(tc feeboostTestCase, linkPrice, wethPrice *big.Int) {
	feedSelector := tc.deployedEnv.FeedChainSel

	if linkPrice != nil {
		require.NoError(tc.t, tc.updatePrice(changeset.LinkSymbol, linkPrice))
		require.NoError(tc.t, changeset.ConfirmPriceUpdate(
			tc.t,
			tc.deployedEnv.Env.Chains[feedSelector],
			tc.onchainState,
			changeset.LinkSymbol,
			linkPrice,
		))
	}

	if wethPrice != nil {
		require.NoError(tc.t, tc.updatePrice(changeset.WethSymbol, wethPrice))
		require.NoError(tc.t, changeset.ConfirmPriceUpdate(
			tc.t,
			tc.deployedEnv.Env.Chains[feedSelector],
			tc.onchainState,
			changeset.WethSymbol,
			wethPrice,
		))
	}
}

func willTriggerFeeBoosting(msgSentEvent *onramp.OnRampCCIPMessageSent, tc feeboostTestCase) bool {
	msg := ConvertToMessage(msgSentEvent.Message)
	fmt.Println("\n=== Fee Boosting Analysis ===")
	fmt.Printf("Message ID: %x\n", msg.Header.MessageID)

	ep := ccipevm.NewGasEstimateProvider()
	chainState, exists := tc.onchainState.Chains[tc.sourceChain]
	require.True(tc.t, exists)
	feeQuoter := chainState.FeeQuoter

	// Get LINK price
	linkPrice, err := feeQuoter.GetTokenPrice(&bind.CallOpts{Context: context.Background()}, chainState.LinkToken.Address())
	require.NoError(tc.t, err)
	fmt.Printf("LINK Price: %s\n", linkPrice.Value.String())

	// Calculate fee in native token terms
	fee := new(big.Int).Div(
		new(big.Int).Mul(linkPrice.Value, msg.FeeValueJuels.Int),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
	)
	fmt.Printf("Fee paid (in native token): %s\n", fee.String())

	// Calculate message gas
	messageGas := new(big.Int).SetUint64(ep.CalculateMessageMaxGas(msg))
	fmt.Printf("Estimated message gas: %s\n", messageGas.String())

	// Get token and gas prices
	nativeTokenAddress := chainState.Weth9.Address()
	tokenAndGasPrice, err := feeQuoter.GetTokenAndGasPrices(&bind.CallOpts{Context: context.Background()}, nativeTokenAddress, tc.destChain)
	require.NoError(tc.t, err)
	fmt.Printf("Raw gas price (uint224): %s\n", tokenAndGasPrice.GasPriceValue.String())

	// Extract uint112 gas price
	gasPrice, err := ConvertGasPriceToUint112(tokenAndGasPrice.GasPriceValue)
	require.NoError(tc.t, err)
	fmt.Printf("Extracted gas price (uint112): %s\n", gasPrice.String())
	fmt.Printf("Native token price: %s\n", tokenAndGasPrice.TokenPrice.String())

	// Calculate execution fee
	tmp := new(big.Int).Mul(gasPrice, tokenAndGasPrice.TokenPrice)
	executionFee := tmp.Div(tmp, big.NewInt(1e18))
	fmt.Printf("Execution fee per gas: %s\n", executionFee.String())

	// Calculate total execution cost
	execCost := new(big.Int).Mul(messageGas, executionFee)
	fmt.Printf("Total execution cost: %s\n", execCost.String())

	// Check if fee boosting will trigger
	willBoost := execCost.Cmp(fee) > 0
	fmt.Printf("\nWill fee boosting trigger? %v\n", willBoost)
	fmt.Printf("Execution cost / Fee ratio: %.2f\n",
		new(big.Float).Quo(
			new(big.Float).SetInt(execCost),
			new(big.Float).SetInt(fee),
		),
	)

	return execCost.Cmp(fee) > 0
}

func ConvertGasPriceToUint112(gasPrice *big.Int) (*big.Int, error) {
	// Create a mask for uint112 (112 bits of 1s)
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 112), big.NewInt(1))

	// Extract the lower 112 bits using AND operation
	result := new(big.Int).And(gasPrice, mask)

	return result, nil
}

func ConvertToMessage(msg onramp.InternalEVM2AnyRampMessage) cciptypes.Message {
	// Convert header
	header := cciptypes.RampMessageHeader{
		MessageID:           cciptypes.Bytes32(msg.Header.MessageId),
		SourceChainSelector: cciptypes.ChainSelector(msg.Header.SourceChainSelector),
		DestChainSelector:   cciptypes.ChainSelector(msg.Header.DestChainSelector),
		SequenceNumber:      cciptypes.SeqNum(msg.Header.SequenceNumber),
		Nonce:               msg.Header.Nonce,
		MsgHash:             cciptypes.Bytes32{},        // This will be populated by the plugin
		OnRamp:              cciptypes.UnknownAddress{}, // This will be populated by the CCIP reader
	}

	// Convert token amounts
	tokenAmounts := make([]cciptypes.RampTokenAmount, len(msg.TokenAmounts))
	for i, ta := range msg.TokenAmounts {
		tokenAmounts[i] = cciptypes.RampTokenAmount{
			SourcePoolAddress: cciptypes.UnknownAddress(ta.SourcePoolAddress.Bytes()),
			DestTokenAddress:  cciptypes.UnknownAddress(ta.DestTokenAddress),
			ExtraData:         cciptypes.Bytes(ta.ExtraData),
			Amount:            cciptypes.BigInt{ta.Amount},
			DestExecData:      cciptypes.Bytes(ta.DestExecData),
		}
	}

	return cciptypes.Message{
		Header:         header,
		Sender:         cciptypes.UnknownAddress(msg.Sender.Bytes()),
		Data:           cciptypes.Bytes(msg.Data),
		Receiver:       cciptypes.UnknownAddress(msg.Receiver),
		ExtraArgs:      cciptypes.Bytes(msg.ExtraArgs),
		FeeToken:       cciptypes.UnknownAddress(msg.FeeToken.Bytes()),
		FeeTokenAmount: cciptypes.BigInt{msg.FeeTokenAmount},
		FeeValueJuels:  cciptypes.BigInt{msg.FeeValueJuels},
		TokenAmounts:   tokenAmounts,
	}
}

func (tc *feeboostTestCase) updatePrice(symbol changeset.TokenSymbol, price *big.Int) error {
	chainSelector := tc.deployedEnv.FeedChainSel
	chainState, exists := tc.onchainState.Chains[chainSelector]
	if !exists {
		return fmt.Errorf("chain state not found for selector: %d", chainSelector)
	}

	feed, exists := chainState.USDFeeds[symbol]
	if !exists {
		return fmt.Errorf("feed not found for token symbol %s on chain %d", symbol, chainSelector)
	}

	// Create mock aggregator instance
	aggr, err := mock_v3_aggregator_contract.NewMockV3AggregatorContract(feed.Address(), tc.deployedEnv.Env.Chains[chainSelector].Client)
	if err != nil {
		return errors.Wrapf(err, "creating aggregator instance for %s on chain %d", symbol, chainSelector)
	}

	// Update price
	auth := tc.deployedEnv.Env.Chains[chainSelector].DeployerKey
	tx, err := aggr.UpdateAnswer(auth, price)
	if err != nil {
		return errors.Wrapf(err, "updating %s price on chain %d", symbol, chainSelector)
	}
	if _, err := deployment.ConfirmIfNoError(tc.deployedEnv.Env.Chains[tc.sourceChain], tx, err); err != nil {
		return err
	}

	return nil
}
