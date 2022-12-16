package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

// KeeperBlockTimeTestReporter enables reporting on the keeper block time test
type KeeperBlockTimeTestReporter struct {
	Reports                        []KeeperBlockTimeTestReport `json:"reports"`
	ReportMutex                    sync.Mutex
	AttemptedChainlinkTransactions []*client.TransactionsData `json:"attemptedChainlinkTransactions"`

	namespace                 string
	keeperReportFile          string
	attemptedTransactionsFile string
}

// KeeperBlockTimeTestReport holds a report information for a single Upkeep Consumer contract
type KeeperBlockTimeTestReport struct {
	ContractAddress        string  `json:"contractAddress"`
	TotalExpectedUpkeeps   int64   `json:"totalExpectedUpkeeps"`
	TotalSuccessfulUpkeeps int64   `json:"totalSuccessfulUpkeeps"`
	AllMissedUpkeeps       []int64 `json:"allMissedUpkeeps"` // List of each time an upkeep was missed, represented by how many blocks it was missed by
}

func (k *KeeperBlockTimeTestReporter) SetNamespace(namespace string) {
	k.namespace = namespace
}

func (k *KeeperBlockTimeTestReporter) WriteReport(folderLocation string) error {
	k.keeperReportFile = filepath.Join(folderLocation, "./block_time_report.csv")
	k.attemptedTransactionsFile = filepath.Join(folderLocation, "./attempted_transactions_report.json")
	keeperReportFile, err := os.Create(k.keeperReportFile)
	if err != nil {
		return err
	}
	defer keeperReportFile.Close()

	keeperReportWriter := csv.NewWriter(keeperReportFile)
	err = keeperReportWriter.Write([]string{
		"Contract Index",
		"Contract Address",
		"Total Expected Upkeeps",
		"Total Successful Upkeeps",
		"Total Missed Upkeeps",
		"Average Blocks Missed",
		"Largest Missed Upkeep",
		"Percent Successful",
	})
	if err != nil {
		return err
	}
	var totalExpected, totalSuccessful, totalMissed, worstMiss int64
	for contractIndex, report := range k.Reports {
		avg, max := int64AvgMax(report.AllMissedUpkeeps)
		err = keeperReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalExpectedUpkeeps),
			fmt.Sprint(report.TotalSuccessfulUpkeeps),
			fmt.Sprint(len(report.AllMissedUpkeeps)),
			fmt.Sprint(avg),
			fmt.Sprint(max),
			fmt.Sprintf("%.2f%%", (float64(report.TotalSuccessfulUpkeeps)/float64(report.TotalExpectedUpkeeps))*100),
		})
		totalExpected += report.TotalExpectedUpkeeps
		totalSuccessful += report.TotalSuccessfulUpkeeps
		totalMissed += int64(len(report.AllMissedUpkeeps))
		worstMiss = int64(math.Max(float64(max), float64(worstMiss)))
		if err != nil {
			return err
		}
	}
	keeperReportWriter.Flush()

	err = keeperReportWriter.Write([]string{"Full Test Summary"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{"Total Expected", "Total Successful", "Total Missed", "Worst Miss", "Total Percent"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{
		fmt.Sprint(totalExpected),
		fmt.Sprint(totalSuccessful),
		fmt.Sprint(totalMissed),
		fmt.Sprint(worstMiss),
		fmt.Sprintf("%.2f%%", (float64(totalSuccessful)/float64(totalExpected))*100)})
	if err != nil {
		return err
	}
	keeperReportWriter.Flush()

	txs, err := json.Marshal(k.AttemptedChainlinkTransactions)
	if err != nil {
		return err
	}
	err = os.WriteFile(k.attemptedTransactionsFile, txs, 0600)
	if err != nil {
		return err
	}

	log.Info().Msg("Successfully wrote report on Keeper Block Timing")
	return nil
}

// SendSlackNotification sends a slack notification on the results of the test
func (k *KeeperBlockTimeTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := ":white_check_mark: Keeper Block Time Test PASSED :white_check_mark:"
	if testFailed {
		headerText = ":x: Keeper Block Time Test FAILED :x:"
	}
	messageBlocks := testreporters.CommonSlackNotificationBlocks(
		t, slackClient, headerText, k.namespace, k.keeperReportFile, testreporters.SlackUserID, testFailed,
	)
	ts, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	if err := testreporters.UploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("Keeper Block Time Test Report %s", k.namespace),
		Filetype:        "csv",
		Filename:        fmt.Sprintf("keeper_block_time_%s.csv", k.namespace),
		File:            k.keeperReportFile,
		InitialComment:  fmt.Sprintf("Keeper Block Time Test Report %s", k.namespace),
		Channels:        []string{testreporters.SlackChannel},
		ThreadTimestamp: ts,
	}); err != nil {
		return err
	}
	return testreporters.UploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("Keeper Block Time Attempted Chainlink Txs %s", k.namespace),
		Filetype:        "json",
		Filename:        fmt.Sprintf("attempted_cl_txs_%s.json", k.namespace),
		File:            k.attemptedTransactionsFile,
		InitialComment:  fmt.Sprintf("Keeper Block Time Attempted Txs %s", k.namespace),
		Channels:        []string{testreporters.SlackChannel},
		ThreadTimestamp: ts,
	})
}

// int64AvgMax helper calculates the avg and the max values in a list
func int64AvgMax(in []int64) (float64, int64) {
	var sum int64
	var max int64
	if len(in) == 0 {
		return 0, 0
	}
	for _, num := range in {
		sum += num
		max = int64(math.Max(float64(max), float64(num)))
	}
	return float64(sum) / float64(len(in)), max
}
