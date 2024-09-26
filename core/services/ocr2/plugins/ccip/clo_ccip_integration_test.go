package ccip_test

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

func Test_CLOSpecApprovalFlow_pipeline(t *testing.T) {
	ccipTH := integrationtesthelpers.SetupCCIPIntegrationTH(
		t,
		testhelpers.SourceChainID,
		testhelpers.SourceChainSelector,
		testhelpers.DestChainID,
		testhelpers.DestChainSelector,
		ccip.DefaultSourceFinalityDepth,
		ccip.DefaultDestFinalityDepth,
	)

	tokenPricesUSDPipeline, linkUSD, ethUSD := ccipTH.CreatePricesPipeline(t)
	defer linkUSD.Close()
	defer ethUSD.Close()

	test_CLOSpecApprovalFlow(t, ccipTH, tokenPricesUSDPipeline, "")
}

func Test_CLOSpecApprovalFlow_dynamicPriceGetter(t *testing.T) {
	ccipTH := integrationtesthelpers.SetupCCIPIntegrationTH(
		t,
		testhelpers.SourceChainID,
		testhelpers.SourceChainSelector,
		testhelpers.DestChainID,
		testhelpers.DestChainSelector,
		ccip.DefaultSourceFinalityDepth,
		ccip.DefaultDestFinalityDepth,
	)

	//Set up the aggregators here to avoid modifying ccipTH.
	dstLinkAddr := ccipTH.Dest.LinkToken.Address()
	srcNativeAddr, err := ccipTH.Source.Router.GetWrappedNative(nil)
	require.NoError(t, err)
	aggDstNativeAddr := ccipTH.Dest.WrappedNative.Address()

	aggSrcNatAddr, _, aggSrcNat, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(ccipTH.Source.User, ccipTH.Source.Chain, 18, big.NewInt(2e18))
	require.NoError(t, err)
	_, err = aggSrcNat.UpdateRoundData(ccipTH.Source.User, big.NewInt(50), big.NewInt(17000000), big.NewInt(1000), big.NewInt(1000))
	require.NoError(t, err)
	ccipTH.Source.Chain.Commit()

	aggDstLnkAddr, _, aggDstLnk, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(ccipTH.Dest.User, ccipTH.Dest.Chain, 18, big.NewInt(3e18))
	require.NoError(t, err)
	ccipTH.Dest.Chain.Commit()
	_, err = aggDstLnk.UpdateRoundData(ccipTH.Dest.User, big.NewInt(50), big.NewInt(8000000), big.NewInt(1000), big.NewInt(1000))
	require.NoError(t, err)
	ccipTH.Dest.Chain.Commit()

	// Check content is ok on aggregator.
	tmp, err := aggDstLnk.LatestRoundData(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(t, big.NewInt(50), tmp.RoundId)
	require.Equal(t, big.NewInt(8000000), tmp.Answer)

	// deploy dest wrapped native aggregator
	aggDstNativeAggrAddr, _, aggDstNativeAggr, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(ccipTH.Dest.User, ccipTH.Dest.Chain, 18, big.NewInt(3e18))
	require.NoError(t, err)
	ccipTH.Dest.Chain.Commit()
	_, err = aggDstNativeAggr.UpdateRoundData(ccipTH.Dest.User, big.NewInt(50), big.NewInt(500000), big.NewInt(1000), big.NewInt(1000))
	require.NoError(t, err)
	ccipTH.Dest.Chain.Commit()

	priceGetterConfig := config.DynamicPriceGetterConfig{
		AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
			srcNativeAddr: {
				ChainID:                   ccipTH.Source.ChainID,
				AggregatorContractAddress: aggSrcNatAddr,
			},
			dstLinkAddr: {
				ChainID:                   ccipTH.Dest.ChainID,
				AggregatorContractAddress: aggDstLnkAddr,
			},
			aggDstNativeAddr: {
				ChainID:                   ccipTH.Dest.ChainID,
				AggregatorContractAddress: aggDstNativeAggrAddr,
			},
		},
		StaticPrices: map[common.Address]config.StaticPriceConfig{},
	}
	priceGetterConfigBytes, err := json.MarshalIndent(priceGetterConfig, "", " ")
	require.NoError(t, err)
	priceGetterConfigJson := string(priceGetterConfigBytes)

	test_CLOSpecApprovalFlow(t, ccipTH, "", priceGetterConfigJson)
}

func test_CLOSpecApprovalFlow(t *testing.T, ccipTH integrationtesthelpers.CCIPIntegrationTestHarness, tokenPricesUSDPipeline string, priceGetterConfiguration string) {
	jobParams := ccipTH.SetUpNodesAndJobs(t, tokenPricesUSDPipeline, priceGetterConfiguration, "http://blah.com")
	ccipTH.SetupFeedsManager(t)

	// Propose and approve new specs
	ccipTH.ApproveJobSpecs(t, jobParams)

	// Sanity check that CCIP works after CLO flow
	currentSeqNum := 1

	extraArgs, err := testhelpers.GetEVMExtraArgsV1(big.NewInt(200_003), false)
	require.NoError(t, err)

	msg := router.ClientEVM2AnyMessage{
		Receiver:     testhelpers.MustEncodeAddress(t, ccipTH.Dest.Receivers[0].Receiver.Address()),
		Data:         utils.RandomAddress().Bytes(),
		TokenAmounts: []router.ClientEVMTokenAmount{},
		FeeToken:     ccipTH.Source.LinkToken.Address(),
		ExtraArgs:    extraArgs,
	}
	fee, err := ccipTH.Source.Router.GetFee(nil, testhelpers.DestChainSelector, msg)
	require.NoError(t, err)

	_, err = ccipTH.Source.LinkToken.Approve(ccipTH.Source.User, ccipTH.Source.Router.Address(), new(big.Int).Set(fee))
	require.NoError(t, err)
	blockHash := ccipTH.Dest.Chain.Commit()
	// get the block number
	block, err := ccipTH.Dest.Chain.BlockByHash(context.Background(), blockHash)
	require.NoError(t, err)
	blockNumber := block.Number().Uint64() + 1 // +1 as a block will be mined for the request from EventuallyReportCommitted

	ccipTH.SendRequest(t, msg)
	ccipTH.AllNodesHaveReqSeqNum(t, currentSeqNum)
	ccipTH.EventuallyReportCommitted(t, currentSeqNum)
	ccipTH.EventuallyPriceRegistryUpdated(
		t,
		blockNumber,
		ccipTH.Source.ChainSelector,
		[]common.Address{ccipTH.Dest.LinkToken.Address(), ccipTH.Dest.WrappedNative.Address()},
		ccipTH.Source.WrappedNative.Address(),
	)

	executionLogs := ccipTH.AllNodesHaveExecutedSeqNums(t, currentSeqNum, currentSeqNum)
	assert.Len(t, executionLogs, 1)
	ccipTH.AssertExecState(t, executionLogs[0], testhelpers.ExecutionStateSuccess)
}
