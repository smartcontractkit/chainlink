package smoke

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

// var configs mercury.TestEnvConfigs

// func TestMain(m *testing.M) {
// 	logging.Init()

// 	configs = mercury.LoadSavedEnvConfigs()

// 	exitVal := m.Run()

// 	configs.Save()

// 	os.Exit(exitVal)
// }

func TestMercurySmoke2(t *testing.T) {
	testEnv, err := mercury.NewEnv(t.Name(), "smoke")

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)
}

// TODO: tests that use different env should use different env configs. Save list of tests into json with [{name: testName, testEnvConfig: {}}]
func TestMercurySmoke(t *testing.T) {
	l := actions.GetTestLogger(t)

	testEnv, err := mercury.NewEnv(t.Name(), "smoke")

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	var (
		feedId      = testEnv.FeedIds[0]
		feedIdBytes = mercury.StringToByte32(feedId)
	)

	testEnv.AddEvmNetwork()

	// // each call updates test env struct
	// // when reconstructing from config, each function should do it individually
	// err = testEnv.AddDON()
	// require.NoError(t, err)

	// ocrConfig, err := testEnv.BuildOCRConfig()
	// require.NoError(t, err)

	// err = testEnv.AddMercuryServer(nil, nil)
	// require.NoError(t, err)

	// // saved in the test env state in list of verifier contracts
	// // verifier1 will be saved in env config so that it can be loaded when config reused
	// verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy")
	// require.NoError(t, err)
	// verifierContract1, err := testEnv.AddVerifierContract("verifier1", verifierProxyContract.Address())
	// require.NoError(t, err)
	// exchangerContract, err := testEnv.AddExchangerContract("exchanger1", verifierProxyContract.Address(),
	// 	"", 255)
	// _ = exchangerContract

	// require.NoError(t, err)

	// // Initialize first feed
	// err = verifierContract1.SetConfig(feedIdBytes, *ocrConfig)
	// require.NoError(t, err)
	// configDetails, err := verifierContract1.LatestConfigDetails(feedIdBytes)
	// require.NoError(t, err)
	// verifierProxyContract.InitializeVerifier(configDetails.ConfigDigest, verifierContract1.Address())

	// // verifierContract2 := testEnv.AddVerifierContract("verifier2")

	// err = testEnv.AddBootstrapJob(verifierContract1.Address(), uint64(configDetails.BlockNumber), feedId)
	// require.NoError(t, err)

	// // testEnv.AddBootstrapJob("bootstrap-2", verifierContract2)

	// err = testEnv.AddOCRJobs(verifierContract1.Address(), uint64(configDetails.BlockNumber), feedId)
	// require.NoError(t, err)

	// err = testEnv.WaitForReportsInMercuryDb([]string{feedId})
	// require.NoError(t, err)

	t.Run("multiple feeds using separate verifier contracts", func(t *testing.T) {
		// setup multiple job specs with different feed ids

	})

	t.Run("test mercury server has report for the latest block number", func(t *testing.T) {
		t.Skip()

		latestBlockNum, err := testEnv.EvmClient.LatestBlockNumber(context.Background())
		_ = latestBlockNum
		require.NoError(t, err, "Err getting latest block number")
		report, _, err := testEnv.MSClient.GetReports(feedId, latestBlockNum-5)
		require.NoError(t, err, "Error getting report from Mercury Server")
		require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
		err = mercuryactions.ValidateReport([]byte(report.ChainlinkBlob))
		require.NoError(t, err, "Error validating mercury report")
	})

	t.Run("test report verfification using Exchanger.ResolveTradeWithReport call", func(t *testing.T) {
		t.Skip()

		order := mercury.Order{
			FeedID:       feedIdBytes,
			CurrencySrc:  mercury.StringToByte32("1"),
			CurrencyDst:  mercury.StringToByte32("2"),
			AmountSrc:    big.NewInt(1),
			MinAmountDst: big.NewInt(2),
			Sender:       common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40b222"),
			Receiver:     common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40bf45"),
		}

		// Commit to a trade
		commitmentHash := mercury.CreateCommitmentHash(order)
		err := testEnv.ExchangerContract.CommitTrade(commitmentHash)
		require.NoError(t, err)

		// Resove the trade and get mercry server url
		encodedCommitment, err := mercury.CreateEncodedCommitment(order)
		require.NoError(t, err)
		mercuryUrlPath, err := testEnv.ExchangerContract.ResolveTrade(encodedCommitment)
		require.NoError(t, err)
		// feedIdHex param is still not fixed in the Exchanger contract. Should be feedIDHex
		fixedMerucyrUrlPath := strings.Replace(mercuryUrlPath, "feedIdHex", "feedIDHex", -1)

		// Get report from mercury server
		msClient := client.NewMercuryServerClient(
			testEnv.MSInfo.LocalUrl, testEnv.MSInfo.AdminId, testEnv.MSInfo.AdminKey)
		report, resp, err := msClient.CallGet(fmt.Sprintf("/client%s", fixedMerucyrUrlPath))
		l.Info().Msgf("Got response from Mercury server. Response: %v. Report: %s", resp, report)
		require.NoError(t, err, "Error getting report from Mercury Server")
		require.NotEmpty(t, report["chainlinkBlob"], "Report response does not contain chainlinkBlob")
		reportBlob := report["chainlinkBlob"].(string)

		// Resolve the trade with report
		reportBytes, err := hex.DecodeString(reportBlob[2:])
		require.NoError(t, err)
		receipt, err := testEnv.ExchangerContract.ResolveTradeWithReport(reportBytes, encodedCommitment)
		require.NoError(t, err)

		// Get transaction logs
		exchangerABI, err := abi.JSON(strings.NewReader(exchanger.ExchangerABI))
		require.NoError(t, err)
		tradeExecuted := map[string]interface{}{}
		err = exchangerABI.UnpackIntoMap(tradeExecuted, "TradeExecuted", receipt.Logs[1].Data)
		require.NoError(t, err)
		l.Info().Interface("TradeExecuted", tradeExecuted).Msg("ResolveTradeWithReport logs")
	})
}
