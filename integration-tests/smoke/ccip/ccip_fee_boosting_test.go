package ccip

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var (
	linkPrice = deployment.E18Mult(100)
	wethPrice = deployment.E18Mult(4000)
)

func Test_CCIPFeeBoosting(t *testing.T) {
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithOCRConfigOverride(func(params *changeset.CCIPOCRParams) {
			params.ExecuteOffChainConfig.RelativeBoostPerWaitHour = 1
			// Disable token price updates
			params.CommitOffChainConfig.TokenPriceBatchWriteFrequency = *config.MustNewDuration(1_000_000 * time.Hour)
			// Disable gas price updates
			params.CommitOffChainConfig.RemoteGasPriceBatchWriteFrequency = *config.MustNewDuration(1_000_000 * time.Hour)
			// Disable token price updates
			params.CommitOffChainConfig.TokenInfo = nil
		}),
	)

	state, err := changeset.LoadOnchainState(e.Env)
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

	// TODO: discrepancy between client and the gas estimator gas price to be fixed - hardcoded for now
	// fetchedGasPriceDest, err := e.Env.Chains[destChain].Client.SuggestGasPrice(tests.Context(t))
	fetchedGasPriceDest := big.NewInt(20e9) // 20 Gwei = default gas price
	require.NoError(t, err)
	originalGasPriceDestUSD := new(big.Int).Div(
		new(big.Int).Mul(fetchedGasPriceDest, wethPrice),
		big.NewInt(1e18),
	)
	t.Log("Gas price on dest chain (USD):", originalGasPriceDestUSD)

	// Adjust destination gas price on source fee quoter to 95% of the current value
	adjustedGasPriceDest :=
		new(big.Int).Div(
			new(big.Int).Mul(originalGasPriceDestUSD, big.NewInt(99)),
			big.NewInt(100),
		)
	t.Log("Adjusted gas price on dest chain:", adjustedGasPriceDest)

	feeQuoterCfg := changeset.DefaultFeeQuoterDestChainConfig()
	// the default adds 10% to the gas price, we want to increase it
	// to make sure the fee boosting will be finished in proper time for testing
	feeQuoterCfg.GasMultiplierWeiPerEth = 120e16

	testhelpers.AddLane(t, &e, sourceChain, destChain, false,
		map[uint64]*big.Int{
			destChain: testhelpers.ToPackedFee(adjustedGasPriceDest, big.NewInt(0)),
		},
		map[common.Address]*big.Int{
			state.Chains[sourceChain].LinkToken.Address(): linkPrice,
			state.Chains[sourceChain].Weth9.Address():     wethPrice,
		},
		feeQuoterCfg)

	// Update token prices in destination chain FeeQuoter
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(changeset.UpdateFeeQuoterPricesChangeset),
			Config: changeset.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]changeset.FeeQuoterPriceUpdatePerSource{
					destChain: {
						TokenPrices: map[common.Address]*big.Int{
							state.Chains[destChain].LinkToken.Address(): linkPrice,
							state.Chains[destChain].Weth9.Address():     wethPrice,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)

	latesthdr, err := e.Env.Chains[sourceChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("message that needs fee boosting"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	startBlocks[sourceChain] = &block
	expectedSeqNum[testhelpers.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}] = msgSentEvent.SequenceNumber
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}] = []uint64{msgSentEvent.SequenceNumber}

	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(changeset.UpdateFeeQuoterPricesChangeset),
			Config: changeset.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]changeset.FeeQuoterPriceUpdatePerSource{
					sourceChain: {
						GasPrices: map[uint64]*big.Int{
							destChain: originalGasPriceDestUSD,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	// Confirm gas prices are updated
	srcFeeQuoter := state.Chains[sourceChain].FeeQuoter
	err = testhelpers.ConfirmGasPriceUpdated(t, e.Env.Chains[destChain], srcFeeQuoter, 0, originalGasPriceDestUSD)
	require.NoError(t, err)

	// Confirm that fee boosting will be triggered
	require.True(t, willTriggerFeeBoosting(t, msgSentEvent, state, sourceChain, destChain))

	// hack
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[sourceChain] = 1
	replayBlocks[destChain] = 1
	testhelpers.ReplayLogs(t, e.Env.Offchain, replayBlocks)

	// Confirm that the message is committed and executed
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
}

// TODO: Find a more accurate way to determine if fee boosting will be triggered
func willTriggerFeeBoosting(
	t *testing.T,
	msgSentEvent *onramp.OnRampCCIPMessageSent,
	state changeset.CCIPOnChainState,
	srcChain, destChain uint64) bool {
	msg := convertToMessage(msgSentEvent.Message)
	t.Log("\n=== Fee Boosting Analysis ===")
	t.Logf("Src Chain: %d", msg.Header.SourceChainSelector)
	t.Logf("Dest Chain: %d", msg.Header.DestChainSelector)

	ep := ccipevm.NewGasEstimateProvider()
	chainState, exists := state.Chains[srcChain]
	require.True(t, exists)
	feeQuoter := chainState.FeeQuoter

	premium, err := feeQuoter.GetPremiumMultiplierWeiPerEth(&bind.CallOpts{Context: context.Background()}, chainState.Weth9.Address())
	require.NoError(t, err)
	t.Logf("Premium: %d", premium)

	// Get LINK price
	linkPrice, err := feeQuoter.GetTokenPrice(&bind.CallOpts{Context: context.Background()}, chainState.LinkToken.Address())
	require.NoError(t, err)
	t.Logf("LINK Price: %s", linkPrice.Value.String())
	t.Logf("Juels in message: %s", msg.FeeValueJuels.String())

	// Calculate fee in USD token
	fee := new(big.Int).Div(
		new(big.Int).Mul(linkPrice.Value, msg.FeeValueJuels.Int),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
	)
	t.Logf("Fee paid (in USD token): %s", fee.String())

	// Calculate message gas
	messageGas := new(big.Int).SetUint64(ep.CalculateMessageMaxGas(msg))
	t.Logf("Estimated message gas: %s", messageGas.String())

	// Get token and gas prices
	nativeTokenAddress := chainState.Weth9.Address()
	tokenAndGasPrice, err := feeQuoter.GetTokenAndGasPrices(&bind.CallOpts{Context: context.Background()}, nativeTokenAddress, destChain)
	require.NoError(t, err)
	t.Logf("Raw gas price (uint224): %s for chain: %d", tokenAndGasPrice.GasPriceValue.String(), destChain)

	// Extract uint112 gas price
	gasPrice, err := convertGasPriceToUint112(tokenAndGasPrice.GasPriceValue)
	require.NoError(t, err)
	t.Logf("Extracted gas price (uint112): %s", gasPrice.String())
	t.Logf("Native token price: %s", tokenAndGasPrice.TokenPrice.String())

	// Calculate total execution cost
	execCost := new(big.Int).Mul(messageGas, gasPrice)
	t.Logf("Total execution cost: %s", execCost.String())

	// Check if fee boosting will trigger
	willBoost := execCost.Cmp(fee) > 0
	t.Logf("\nWill fee boosting trigger? %v", willBoost)
	t.Logf("Execution cost / Fee ratio: %.2f",
		new(big.Float).Quo(
			new(big.Float).SetInt(execCost),
			new(big.Float).SetInt(fee),
		),
	)

	return execCost.Cmp(fee) > 0
}

func convertGasPriceToUint112(gasPrice *big.Int) (*big.Int, error) {
	// Create a mask for uint112 (112 bits of 1s)
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 112), big.NewInt(1))

	// Extract the lower 112 bits using AND operation
	result := new(big.Int).And(gasPrice, mask)

	return result, nil
}

func convertToMessage(msg onramp.InternalEVM2AnyRampMessage) cciptypes.Message {
	// Convert header
	header := cciptypes.RampMessageHeader{
		MessageID:           cciptypes.Bytes32(msg.Header.MessageId),
		SourceChainSelector: cciptypes.ChainSelector(msg.Header.SourceChainSelector),
		DestChainSelector:   cciptypes.ChainSelector(msg.Header.DestChainSelector),
		SequenceNumber:      cciptypes.SeqNum(msg.Header.SequenceNumber),
		Nonce:               msg.Header.Nonce,
	}

	// Convert token amounts
	tokenAmounts := make([]cciptypes.RampTokenAmount, len(msg.TokenAmounts))
	for i, ta := range msg.TokenAmounts {
		tokenAmounts[i] = cciptypes.RampTokenAmount{
			SourcePoolAddress: cciptypes.UnknownAddress(ta.SourcePoolAddress.Bytes()),
			DestTokenAddress:  cciptypes.UnknownAddress(ta.DestTokenAddress),
			ExtraData:         cciptypes.Bytes(ta.ExtraData),
			Amount:            cciptypes.BigInt{Int: ta.Amount},
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
		FeeTokenAmount: cciptypes.BigInt{Int: msg.FeeTokenAmount},
		FeeValueJuels:  cciptypes.BigInt{Int: msg.FeeValueJuels},
		TokenAmounts:   tokenAmounts,
	}
}
