package smoke

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	mercury_server "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/stretchr/testify/require"
)

func createCommitmentHash(order Order) common.Hash {
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	addressTy, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
	}

	bytes, _ := arguments.Pack(
		order.FeedID,
		order.CurrencySrc,
		order.CurrencyDst,
		order.AmountSrc,
		order.MinAmountDst,
		order.Sender,
		order.Receiver,
	)

	return crypto.Keccak256Hash(bytes)
}

type Order struct {
	FeedID       [32]byte
	CurrencySrc  [32]byte
	CurrencyDst  [32]byte
	AmountSrc    *big.Int
	MinAmountDst *big.Int
	Sender       common.Address
	Receiver     common.Address
}

func createEncodedCommitment(order Order) ([]byte, error) {
	// bytes32 feedID, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver
	orderType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "feedID", Type: "bytes32"},
		{Name: "currencySrc", Type: "bytes32"},
		{Name: "currencyDst", Type: "bytes32"},
		{Name: "amountSrc", Type: "uint256"},
		{Name: "minAmountDst", Type: "uint256"},
		{Name: "sender", Type: "address"},
		{Name: "receiver", Type: "address"},
	})
	var args abi.Arguments = []abi.Argument{{Type: orderType}}
	return args.Pack(order)
}

var feedId = testsetups.StringToByte32("ETH-USD-1")

// func TestContracts(t *testing.T) {
// 	testNetwork := networks.SelectedNetwork
// 	// evmConfig := eth.New(nil)
// 	// if !testNetwork.Simulated {
// 	// 	evmConfig = eth.New(&eth.Props{
// 	// 		NetworkName: testNetwork.Name,
// 	// 		Simulated:   testNetwork.Simulated,
// 	// 		WsURLs:      testNetwork.URLs,
// 	// 	})
// 	// }
// 	evmClient, err := blockchain.NewEVMClient(testNetwork, nil)
// 	require.NoError(t, err, "Error connecting to blockchain")

// 	contractDeployer, err := contracts.NewContractDeployer(evmClient)
// 	require.NoError(t, err, "Deploying contracts shouldn't fail")

// 	// accessController, err := contractDeployer.DeployReadAccessController()
// 	// require.NoError(t, err, "Error deploying ReadAccessController contract")

// 	// verifierProxy, err := contractDeployer.DeployVerifierProxy(accessController.Address())
// 	// Use zero address for access controller disables access control
// 	verifierProxy, err := contractDeployer.DeployVerifierProxy("0x0")
// 	require.NoError(t, err, "Error deploying VerifierProxy contract")

// 	verifier, err := contractDeployer.DeployVerifier(verifierProxy.Address())
// 	require.NoError(t, err, "Error deploying Verifier contract")
// 	_ = verifier
// }

func TestMercurySmoke(t *testing.T) {
	l := zerolog.New(zerolog.NewTestWriter(t))

	testEnv, isExistingTestEnv, testNetwork, chainlinkNodes,
		mercuryServerRemoteUrl,
		evmClient, mockServerClient, mercuryServerClient, msRpcPubKey := testsetups.SetupMercuryEnv(t, nil, nil)
	_ = isExistingTestEnv

	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig := testsetups.BuildMercuryOCRConfig(t, nodesWithoutBootstrap)
	verifier, verifierProxy, accessController, _ := testsetups.SetupMercuryContracts(t, evmClient,
		mercuryServerRemoteUrl, feedId, ocrConfig)
	_ = verifierProxy
	_ = accessController

	latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
	require.NoError(t, err)

	mercuryServerLocalUrl := testEnv.URLs[mercury_server.URLsKey][0]
	testsetups.SetupMercuryNodeJobs(t, chainlinkNodes, mockServerClient, verifier.Address(),
		feedId, latestBlockNum, mercuryServerLocalUrl, msRpcPubKey, testNetwork.ChainID, 0)

	verifier.SetConfig(feedId, ocrConfig)

	// Wait for the DON to start generating reports
	d := 160 * time.Second
	l.Info().Msgf("Sleeping for %s to wait for Mercury env to be ready..", d)
	time.Sleep(d)

	// mercuryLookupUrl := fmt.Sprintf("%s/client", mercuryServerRemoteUrl)
	// contractDeployer, err := contracts.NewContractDeployer(evmClient)
	// require.NoError(t, err, "Error in contract deployer")

	// exchangerContract, err := contractDeployer.DeployExchanger(verifierProxy.Address(), mercuryLookupUrl, 255)
	// require.NoError(t, err, "Error deploying Exchanger contract")
	// err = accessController.AddAccess(exchangerContract.Address())
	// require.NoError(t, err, "Error in AddAccess(exchanger.Address())")

	t.Run("test mercury server has report for the latest block number", func(t *testing.T) {
		latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
		require.NoError(t, err, "Err getting latest block number")
		report, _, err := mercuryServerClient.GetReports(string(feedId[:]), latestBlockNum-5)
		require.NoError(t, err, "Error getting report from Mercury Server")
		require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
	})

	// t.Run("test report verfification using Exchanger.ResolveTradeWithReport call", func(t *testing.T) {
	// 	order := Order{
	// 		FeedID:       feedId,
	// 		CurrencySrc:  StringToByte32("1"),
	// 		CurrencyDst:  StringToByte32("2"),
	// 		AmountSrc:    big.NewInt(1),
	// 		MinAmountDst: big.NewInt(2),
	// 		Sender:       common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40b222"),
	// 		Receiver:     common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40bf45"),
	// 	}

	// 	// Commit to a trade
	// 	commitmentHash := createCommitmentHash(order)
	// 	err = exchangerContract.CommitTrade(commitmentHash)
	// 	require.NoError(t, err)

	// 	// Resove the trade and get mercry server url
	// 	encodedCommitment, err := createEncodedCommitment(order)
	// 	require.NoError(t, err)
	// 	mercuryUrl, err := exchangerContract.ResolveTrade(encodedCommitment)
	// 	require.NoError(t, err)

	// 	// Get report from Mercury server
	// 	report := &client.GetReportsResult{}
	// 	resp, err := resty.New().R().SetResult(&report).Get(mercuryUrl)
	// 	l.Info().Msgf("Got response from Mercury server: %s", resp)
	// 	require.NoError(t, err, "Error getting report from Mercury Server")
	// 	require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")

	// 	// Resolve the trade with report
	// 	reportBytes, err := hex.DecodeString(report.ChainlinkBlob[2:])
	// 	require.NoError(t, err)
	// 	receipt, err := exchangerContract.ResolveTradeWithReport(reportBytes, encodedCommitment)
	// 	require.NoError(t, err)

	// 	// Get transaction logs
	// 	exchangerABI, err := abi.JSON(strings.NewReader(exchanger.ExchangerABI))
	// 	require.NoError(t, err)
	// 	tradeExecuted := map[string]interface{}{}
	// 	err = exchangerABI.UnpackIntoMap(tradeExecuted, "TradeExecuted", receipt.Logs[1].Data)
	// 	require.NoError(t, err)
	// 	l.Info().Interface("TradeExecuted", tradeExecuted).Msg("ResolveTradeWithReport logs")
	// })
}
