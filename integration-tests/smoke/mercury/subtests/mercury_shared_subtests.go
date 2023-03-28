package subtests

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
	"github.com/rs/zerolog/log"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
	"github.com/test-go/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func RunTestGetReportNotFound(t *testing.T, te *mercury.TestEnv, feedId string) {
	t.Run(fmt.Sprintf("get report by feed id string and block number which does not exist, feedId: %s", feedId),
		func(t *testing.T) {
			t.Parallel()

			lastBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			queryBlockNum := lastBlockNum + 500

			reportStr, resp, err := te.MSClient.GetReportsByFeedId(feedId, queryBlockNum, client.StringFeedId)
			require.NoError(t, err, "Error getting report from Mercury Server")
			require.Equal(t, 404, resp.StatusCode)
			require.Empty(t, reportStr.ChainlinkBlob, "Report response should not contain chainlinkBlob")
		})
}

func RunTestsGetBulkReportsForRecentBlockNum(t *testing.T, te *mercury.TestEnv, feedIdStr string, feedIdType client.FeedIdType) {
	t.Run(fmt.Sprintf("bulk get reports by feed id %s, feedId: %s", feedIdType, feedIdStr),
		func(t *testing.T) {
			t.Parallel()

			lastBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			var queryBlockNum uint64
			switch te.EvmNetwork.ChainID {
			// Arbitrum Goerli is fast so query for older reports
			case 421613:
				queryBlockNum = lastBlockNum - 15
			default:
				queryBlockNum = lastBlockNum - 10
			}

			var limit uint64 = 100

			var feedId string
			if feedIdType == client.HexFeedId {
				feedId = fmt.Sprintf("0x%x", mercury.StringToByte32(feedIdStr))
			} else {
				feedId = feedIdStr
			}

			result, _, err := te.MSClient.BulkGetReportsByFeedId(feedId, queryBlockNum, limit, feedIdType)
			require.NoError(t, err, "Error getting reports from Mercury Server")
			for _, blob := range result.ChainlinkBlob {
				reportBytes, err := hex.DecodeString(blob[2:])
				require.NoError(t, err)
				reportCtx, err := mercuryactions.DecodeReport(reportBytes)
				require.NoError(t, err)
				log.Debug().Msgf("received report: %+v", reportCtx)
			}
		})
}

func RunTestGetReportByFeedIdForRecentBlockNum(t *testing.T, te *mercury.TestEnv, feedIdStr string, feedIdType client.FeedIdType) {
	t.Run(fmt.Sprintf("get report by feed id %s for the recent block number, feedId: %s", feedIdType, feedIdStr),
		func(t *testing.T) {
			t.Parallel()

			lastBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			var queryBlockNum uint64
			switch te.EvmNetwork.ChainID {
			// Arbitrum Goerli is fast so query for older reports
			case 421613:
				queryBlockNum = lastBlockNum - 15
			default:
				queryBlockNum = lastBlockNum - 10
			}

			var feedId string
			if feedIdType == client.HexFeedId {
				feedId = fmt.Sprintf("0x%x", mercury.StringToByte32(feedIdStr))
			} else {
				feedId = feedIdStr
			}

			reportStr, _, err := te.MSClient.GetReportsByFeedId(feedId, queryBlockNum, client.StringFeedId)
			require.NoError(t, err, "Error getting report from Mercury Server")
			require.NotEmpty(t, reportStr.ChainlinkBlob, "Report response does not contain chainlinkBlob")
			reportBytes, err := hex.DecodeString(reportStr.ChainlinkBlob[2:])
			require.NoError(t, err)
			reportCtx, err := mercuryactions.DecodeReport(reportBytes)
			require.NoError(t, err)
			log.Info().Msgf("received report: %+v", reportCtx)
		})
}

func RunTestGetReportByFeedIdStrFromWS(t *testing.T, te *mercury.TestEnv, feedId string) {
	t.Run("get report by feed id from /ws websocket", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		c, _, err := te.MSClient.DialWS(ctx, fmt.Sprintf("?feedIds=%s,abc1,def-", feedId))
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "")

		m := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c, &m)
		require.NoError(t, err, "failed read ws msg from instance")

		r, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		log.Info().Msgf("received report: %+v", r)
	})
}

func RunAllWSTests(t *testing.T, te *mercury.TestEnv, feedId string) {
	RunTestGetReportByFeedIdStrFromWS(t, te, feedId)

	t.Run("open 2 /ws connections sequentially and get reports", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		c, _, err := te.MSClient.DialWS(ctx, fmt.Sprintf("?feedIds=%s,abc1,def-", feedId))
		require.NoError(t, err)

		m := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c, &m)
		require.NoError(t, err, "failed read ws msg from instance")

		r, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		log.Info().Msgf("received report: %+v", r)

		err = c.Close(websocket.StatusNormalClosure, "")
		require.NoError(t, err)

		c2, _, err := te.MSClient.DialWS(ctx, fmt.Sprintf("?feedIds=%s,abc1,def-", feedId))
		require.NoError(t, err)
		defer c2.Close(websocket.StatusNormalClosure, "")

		m2 := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c2, &m2)
		require.NoError(t, err, "failed read ws msg from instance")

		r2, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		log.Info().Msgf("received report: %+v", r2)
	})

	t.Run("read 2 reports by feed id from /ws", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		c, _, err := te.MSClient.DialWS(ctx, fmt.Sprintf("?feedIds=%s,abc1,def-", feedId))
		require.NoError(t, err)
		defer c.Close(websocket.StatusNormalClosure, "")

		m := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c, &m)
		require.NoError(t, err, "failed read ws msg from instance")

		r, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		log.Info().Msgf("received report: %+v", r)

		m2 := client.NewReportWSMessage{}
		err = wsjson.Read(context.Background(), c, &m2)
		require.NoError(t, err, "failed read ws msg from instance")

		r2, err := mercuryactions.DecodeReport(m.FullReport)
		require.NoError(t, err)
		log.Info().Msgf("received report: %+v", r2)
	})
}

func RunTestReportVerificationWithVerifierContract(t *testing.T, te *mercury.TestEnv, verifierProxy contracts.VerifierProxy, feedId string) {
	t.Run("verify report using verifier contract",
		func(t *testing.T) {
			t.Parallel()

			lastBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			var queryBlockNum uint64
			switch te.EvmNetwork.ChainID {
			// Arbitrum Goerli is fast so query for older reports
			case 421613:
				queryBlockNum = lastBlockNum - 15
			default:
				queryBlockNum = lastBlockNum - 10
			}

			reportStr, _, err := te.MSClient.GetReportsByFeedId(feedId, queryBlockNum, client.StringFeedId)
			require.NoError(t, err, "Error getting report from Mercury Server")
			require.NotEmpty(t, reportStr.ChainlinkBlob, "Report response does not contain chainlinkBlob")
			reportBytes, err := hex.DecodeString(reportStr.ChainlinkBlob[2:])
			require.NoError(t, err)
			reportCtx, err := mercuryactions.DecodeReport(reportBytes)
			require.NoError(t, err)
			log.Info().Msgf("Decoded report: %+v", reportCtx)

			err = verifierProxy.Verify(reportBytes)
			require.NoError(t, err)
		})
}

// This will fail if https://smartcontract-it.atlassian.net/browse/MERC-337 not resolved
func RunTestReportVerificationWithExchangerContract(t *testing.T, te *mercury.TestEnv,
	exchangerContract contracts.Exchanger, feedId string) {
	feedIdBytes := mercury.StringToByte32(feedId)

	t.Run(fmt.Sprintf("test report verfification using Exchanger.ResolveTradeWithReport call, feedId: %s", feedId),
		func(t *testing.T) {
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
			err := exchangerContract.CommitTrade(commitmentHash)
			require.NoError(t, err)

			// Resove the trade and get mercry server url
			encodedCommitment, err := mercury.CreateEncodedCommitment(order)
			require.NoError(t, err)
			mercuryUrlPath, err := exchangerContract.ResolveTrade(encodedCommitment)
			require.NoError(t, err)
			// feedIdHex param is still not fixed in the Exchanger contract. Should be feedIDHex
			fixedMerucyrUrlPath := strings.Replace(mercuryUrlPath, "feedIdHex", "feedIDHex", -1)
			fixedMerucyrUrlPath2 := strings.Replace(fixedMerucyrUrlPath, "L2Blocknumber", "blockNumber", -1)

			d := 2 * time.Second
			log.Info().Msgf("Wait for %s report to be generated and available on the mercury server..", d)
			time.Sleep(d)

			// Get report from mercury server
			report, resp, err := te.MSClient.CallGet(fmt.Sprintf("/client%s", fixedMerucyrUrlPath2))
			log.Info().Msgf("Got response from Mercury server. Response: %v. Report: %s", resp, report)
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
			log.Info().Interface("TradeExecuted", tradeExecuted).Msg("ResolveTradeWithReport logs")
		})
}
