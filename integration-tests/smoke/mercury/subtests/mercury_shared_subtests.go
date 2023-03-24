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
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum/mercury/exchanger"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
	"github.com/test-go/testify/require"
)

func RunTestMercuryServerHasReportForRecentBlockNum(t *testing.T, te *mercury.TestEnv, feedId string) {
	t.Run(fmt.Sprintf("test mercury server has report for the recent block number, feedId: %s", feedId),
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

			report, _, err := te.MSClient.GetReportsByFeedIdStr(feedId, queryBlockNum)
			require.NoError(t, err, "Error getting report from Mercury Server")
			require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
		})
}

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

			d := 3 * time.Second
			log.Info().Msgf("Wait for %s report to be generated and available on the mercury server..", d)
			time.Sleep(d)

			// Get report from mercury server
			msClient := client.NewMercuryServerClient(
				te.MSInfo.LocalUrl, te.MSInfo.UserId, te.MSInfo.UserKey)
			report, resp, err := msClient.CallGet(fmt.Sprintf("/client%s", fixedMerucyrUrlPath2))
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
