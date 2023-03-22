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

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

func TestMercuryManyFeedsManyVerifiers(t *testing.T) {
	l := utils.GetTestLogger(t)

	var (
		feedIds = [][32]byte{mercury.StringToByte32("feed-1"), mercury.StringToByte32("feed-2")}
	)

	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	testEnv.AddEvmNetwork()

	err = testEnv.AddDON()
	require.NoError(t, err)

	ocrConfig, err := testEnv.BuildOCRConfig()
	require.NoError(t, err)

	_, _, err = testEnv.AddMercuryServer(nil)
	require.NoError(t, err)

	verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy1")
	require.NoError(t, err)
	exchangerContract, err := testEnv.AddExchangerContract("exchanger1", verifierProxyContract.Address(),
		"", 255)
	require.NoError(t, err)

	for i, feedId := range feedIds {
		verifierContractId := fmt.Sprintf("verifier_%d", i)
		verifierContract, err := testEnv.AddVerifierContract(verifierContractId, verifierProxyContract.Address())
		require.NoError(t, err)

		blockNumber, err := testEnv.SetConfigAndInitializeVerifierContract(
			fmt.Sprintf("setAndInitializeVerifier%d", i),
			verifierContractId,
			"verifierProxy1",
			feedId,
			*ocrConfig,
		)
		require.NoError(t, err)

		err = testEnv.AddBootstrapJob(fmt.Sprintf("createBoostrap%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)

		err = testEnv.AddOCRJobs(fmt.Sprintf("createOcrJobs%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)
	}

	err = testEnv.WaitForReportsInMercuryDb(feedIds)
	require.NoError(t, err)

	for _, feedId := range feedIds {
		feedIdStr := mercury.Byte32ToString(feedId)

		t.Run(fmt.Sprintf("test mercury server has report for the latest block number, feedId: %s", feedIdStr),
			func(t *testing.T) {
				latestBlockNum, err := testEnv.EvmClient.LatestBlockNumber(context.Background())
				require.NoError(t, err, "Err getting latest block number")

				report, _, err := testEnv.MSClient.GetReportsByFeedIdStr(feedIdStr, latestBlockNum-5)
				require.NoError(t, err, "Error getting report from Mercury Server")
				require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
			})

		t.Run(fmt.Sprintf("test report verfification using Exchanger.ResolveTradeWithReport call, feedId: %s", feedIdStr),
			func(t *testing.T) {
				order := mercury.Order{
					FeedID:       feedId,
					CurrencySrc:  mercury.StringToByte32("1"),
					CurrencyDst:  mercury.StringToByte32("2"),
					AmountSrc:    big.NewInt(1),
					MinAmountDst: big.NewInt(2),
					Sender:       common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40b222"),
					Receiver:     common.HexToAddress("c7ca5f083dce8c0034e9a6033032ec576d40bf45"),
				}

				// Commit to a trade
				commitmentHash := mercury.CreateCommitmentHash(order)
				err := exchangerContract.CommitTrade(commitmentHash)
				require.NoError(t, err)

				// Resove the trade and get mercry server url
				encodedCommitment, err := mercury.CreateEncodedCommitment(order)
				require.NoError(t, err)
				mercuryUrlPath, err := exchangerContract.ResolveTrade(encodedCommitment)
				require.NoError(t, err)
				// feedIdHex param is still not fixed in the Exchanger contract. Should be feedIDHex
				fixedMerucyrUrlPath := strings.Replace(mercuryUrlPath, "feedIdHex", "feedIDHex", -1)

				// Get report from mercury server
				msClient := client.NewMercuryServerClient(
					testEnv.MSInfo.LocalUrl, testEnv.MSInfo.UserId, testEnv.MSInfo.UserKey)
				r, resp, err := msClient.CallGet(fmt.Sprintf("/client%s", fixedMerucyrUrlPath))
				report := r.(map[string]interface{})
				l.Info().Msgf("Got response from Mercury server. Response: %v. Report: %s", resp, report)
				require.NoError(t, err, "Error getting report from Mercury Server")
				require.NotEmpty(t, report["chainlinkBlob"], "Report response does not contain chainlinkBlob")
				reportBlob := report["chainlinkBlob"].(string)

				// Resolve the trade with report
				reportBytes, err := hex.DecodeString(reportBlob[2:])
				require.NoError(t, err)
				receipt, err := exchangerContract.ResolveTradeWithReport(reportBytes, encodedCommitment)
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
}
