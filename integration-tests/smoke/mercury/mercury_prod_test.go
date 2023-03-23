package mercury

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

// To run this test, provide path to env config file in MERCURY_ENV_CONFIG_PATH
// Example:
//
// {
//     "id": "TestSmokeMercuryProd",
//     "chainId": 420,
//     "feedId": "feed-1",
//     "contracts": {
//         "verifierProxy": "0x42973a598f94Dd6A14a1F2E9CB336Fe88672Fa79"
//     },
//     "mercuryServer": {
//         "remoteUrl": "http://10.14.90.115:3000",
//         "userId": "02185d5a-f1ee-40d1-a52a-bf39871b614c",
//         "userKey": "admintestkey"
//     }
// }

func TestSmokeMercuryProd(t *testing.T) {
	l := utils.GetTestLogger(t)

	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)
	require.NoError(t, err)
	if testEnv.C == nil {
		t.Skip("Test is skipped because env config file was not provided")
	}
	feedId := testEnv.C.FeedId
	require.NotEmpty(t, feedId, "'feedId' needs to be provided in the env config file")

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})

	err = testEnv.AddEvmNetwork()
	require.NoError(t, err)

	msClient := client.NewMercuryServerClient(
		testEnv.MSInfo.RemoteUrl, testEnv.MSInfo.UserId, testEnv.MSInfo.UserKey)

	verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy")
	require.NoError(t, err)
	exchangerContract, err := testEnv.AddExchangerContract("exchanger", verifierProxyContract.Address(),
		"", 255)
	require.NoError(t, err)

	t.Run("get report by feed id str for the latest block number-2", func(t *testing.T) {
		// latestBlockNum, err := testEnv.EvmClient.LatestBlockNumber(context.Background())
		// require.NoError(t, err, "Err getting latest block number")

		reportData, _, err := msClient.GetReportsByFeedIdStr(feedId, 12905278)
		require.NoError(t, err)
		require.NotEmpty(t, reportData.ChainlinkBlob, "received empty ChainlinkBlob")
		reportBytes, err := hex.DecodeString(reportData.ChainlinkBlob[2:])
		require.NoError(t, err)
		reportCtx, err := mercuryactions.DecodeReport(reportBytes)
		require.NoError(t, err)
		l.Info().Msgf("received report: %+v", reportCtx)
	})

	t.Run("get report by feed id hex for the latest block number-2", func(t *testing.T) {
		// latestBlockNum, err := testEnv.EvmClient.LatestBlockNumber(context.Background())
		// require.NoError(t, err, "Err getting latest block number")

		feedIdHex := fmt.Sprintf("0x%x", mercury.StringToByte32(feedId))
		reportData, _, err := msClient.GetReportsByFeedIdHex(feedIdHex, 12905278)
		require.NoError(t, err)
		require.NotEmpty(t, reportData.ChainlinkBlob, "received empty ChainlinkBlob")
		reportBytes, err := hex.DecodeString(reportData.ChainlinkBlob[2:])
		require.NoError(t, err)
		reportCtx, err := mercuryactions.DecodeReport(reportBytes)
		require.NoError(t, err)
		l.Info().Msgf("received report: %+v", reportCtx)
	})

	t.Run("get report by feed id from /ws websocket", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		c, _, err := msClient.DialWS(ctx)
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "")

		m := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c, &m)
		require.NoError(t, err, "failed read ws msg from instance")

		r, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		l.Info().Msgf("received report: %+v", r)
	})

	t.Run("get report and verify it on chain using Exchanger.ResolveTradeWithReport call",
		func(t *testing.T) {
			order := mercury.Order{
				FeedID:       mercury.StringToByte32(feedId),
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
			report, resp, err := msClient.CallGet(fmt.Sprintf("/client%s", fixedMerucyrUrlPath))
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
