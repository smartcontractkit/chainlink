package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

// KeeperBenchmarkTestReporter enables reporting on the keeper benchmark test
type KeeperBenchmarkTestReporter struct {
	Reports                        []KeeperBenchmarkTestReport `json:"reports"`
	ReportMutex                    sync.Mutex
	AttemptedChainlinkTransactions []*nodeclient.TransactionsData `json:"attemptedChainlinkTransactions"`
	NumRevertedUpkeeps             int64
	NumStaleUpkeepReports          int64
	Summary                        KeeperBenchmarkTestSummary `json:"summary"`

	namespace                 string
	keeperReportFile          string
	attemptedTransactionsFile string
	keeperSummaryFile         string
}

type KeeperBenchmarkTestSummary struct {
	Load       KeeperBenchmarkTestLoad    `json:"load"`
	Config     KeeperBenchmarkTestConfig  `json:"config"`
	Metrics    KeeperBenchmarkTestMetrics `json:"metrics"`
	TestInputs map[string]interface{}     `json:"testInputs"`
	StartTime  int64                      `json:"startTime"`
	EndTime    int64                      `json:"endTime"`
}

type KeeperBenchmarkTestLoad struct {
	TotalCheckGasPerBlock           int64   `json:"totalCheckGasPerBlock"`
	TotalPerformGasPerBlock         int64   `json:"totalPerformGasPerBlock"`
	AverageExpectedPerformsPerBlock float64 `json:"averageExpectedPerformsPerBlock"`
}

type KeeperBenchmarkTestConfig struct {
	Chainlink map[string]map[string]string `json:"chainlink"`
	Geth      map[string]map[string]string `json:"geth"`
}

type KeeperBenchmarkTestMetrics struct {
	Delay                         map[string]interface{} `json:"delay"`
	PercentWithinSLA              float64                `json:"percentWithinSLA"`
	PercentRevert                 float64                `json:"percentRevert"`
	PercentStale                  float64                `json:"percentStale"`
	TotalTimesEligible            int64                  `json:"totalTimesEligible"`
	TotalTimesPerformed           int64                  `json:"totalTimesPerformed"`
	TotalStaleReports             int64                  `json:"totalStaleReports"`
	AverageActualPerformsPerBlock float64                `json:"averageActualPerformsPerBlock"`
}

// KeeperBenchmarkTestReport holds a report information for a single Upkeep Consumer contract
type KeeperBenchmarkTestReport struct {
	RegistryAddress       string  `json:"registryAddress"`
	ContractAddress       string  `json:"contractAddress"`
	TotalEligibleCount    int64   `json:"totalEligibleCount"`
	TotalSLAMissedUpkeeps int64   `json:"totalSLAMissedUpkeeps"`
	TotalPerformedUpkeeps int64   `json:"totalPerformedUpkeeps"`
	AllCheckDelays        []int64 `json:"allCheckDelays"` // List of the delays since checkUpkeep for all performs
}

func (k *KeeperBenchmarkTestReporter) SetNamespace(namespace string) {
	k.namespace = namespace
}

func (k *KeeperBenchmarkTestReporter) WriteReport(folderLocation string) error {
	k.keeperReportFile = filepath.Join(folderLocation, "./benchmark_report.csv")
	k.keeperSummaryFile = filepath.Join(folderLocation, "./benchmark_summary.json")
	// k.keeperSummaryCsvFile = filepath.Join(folderLocation, "./benchmark_summary.csv")
	k.attemptedTransactionsFile = filepath.Join(folderLocation, "./attempted_transactions_report.json")
	keeperReportFile, err := os.Create(k.keeperReportFile)
	if err != nil {
		return err
	}
	defer keeperReportFile.Close()

	keeperReportWriter := csv.NewWriter(keeperReportFile)
	var totalEligibleCount, totalPerformed, totalMissedSLA, totalReverted, totalStaleReports int64
	var allDelays []int64
	for _, report := range k.Reports {
		totalEligibleCount += report.TotalEligibleCount
		totalPerformed += report.TotalPerformedUpkeeps
		totalMissedSLA += report.TotalSLAMissedUpkeeps

		allDelays = append(allDelays, report.AllCheckDelays...)
	}
	totalReverted = k.NumRevertedUpkeeps
	totalStaleReports = k.NumStaleUpkeepReports
	pctWithinSLA := (1.0 - float64(totalMissedSLA)/float64(totalEligibleCount)) * 100
	var pctReverted, pctStale float64
	if totalPerformed > 0 {
		pctReverted = (float64(totalReverted) / float64(totalPerformed)) * 100
		pctStale = (float64(totalStaleReports) / float64(totalPerformed)) * 100
	}

	err = keeperReportWriter.Write([]string{"Full Test Summary"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{
		"Total Times Eligible",
		"Total Performed",
		"Total Reverted",
		"Total Stale Reports",
		"Average Perform Delay",
		"Median Perform Delay",
		"90th pct Perform Delay",
		"99th pct Perform Delay",
		"Max Perform Delay",
		"Percent Within SLA",
		"Percent Revert",
		"Percent Stale",
	})
	if err != nil {
		return err
	}
	avg, median, ninetyPct, ninetyNinePct, maxVal := IntListStats(allDelays)
	err = keeperReportWriter.Write([]string{
		fmt.Sprint(totalEligibleCount),
		fmt.Sprint(totalPerformed),
		fmt.Sprint(totalReverted),
		fmt.Sprint(totalStaleReports),
		fmt.Sprintf("%.2f", avg),
		fmt.Sprint(median),
		fmt.Sprint(ninetyPct),
		fmt.Sprint(ninetyNinePct),
		fmt.Sprint(maxVal),
		fmt.Sprintf("%.2f%%", pctWithinSLA),
		fmt.Sprintf("%.2f%%", pctReverted),
		fmt.Sprintf("%.2f%%", pctStale),
	})
	if err != nil {
		return err
	}
	keeperReportWriter.Flush()
	log.Info().
		Int64("Total Times Eligible", totalEligibleCount).
		Int64("Total Performed", totalPerformed).
		Int64("Total Reverted", totalReverted).
		Float64("Average Perform Delay", avg).
		Int64("Median Perform Delay", median).
		Int64("90th pct Perform Delay", ninetyPct).
		Int64("99th pct Perform Delay", ninetyNinePct).
		Int64("Max Perform Delay", maxVal).
		Float64("Percent Within SLA", pctWithinSLA).
		Float64("Percent Reverted", pctReverted).
		Msg("Calculated Aggregate Results")

	err = keeperReportWriter.Write([]string{
		"Contract Index",
		"RegistryAddress",
		"Contract Address",
		"Total Times Eligible",
		"Total Performed Upkeeps",
		"Average Perform Delay",
		"Median Perform Delay",
		"90th pct Perform Delay",
		"99th pct Perform Delay",
		"Largest Perform Delay",
		"Percent Within SLA",
	})
	if err != nil {
		return err
	}

	for contractIndex, report := range k.Reports {
		avg, median, ninetyPct, ninetyNinePct, maxVal = IntListStats(report.AllCheckDelays)
		err = keeperReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.RegistryAddress,
			report.ContractAddress,
			fmt.Sprint(report.TotalEligibleCount),
			fmt.Sprint(report.TotalPerformedUpkeeps),
			fmt.Sprintf("%.2f", avg),
			fmt.Sprint(median),
			fmt.Sprint(ninetyPct),
			fmt.Sprint(ninetyNinePct),
			fmt.Sprint(maxVal),
			fmt.Sprintf("%.2f%%", (1.0-float64(report.TotalSLAMissedUpkeeps)/float64(report.TotalEligibleCount))*100),
		})
		if err != nil {
			return err
		}
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

	log.Info().Msg("Successfully wrote report on Keeper Benchmark")

	k.Summary.Metrics.Delay = map[string]interface{}{
		"mean":   avg,
		"median": median,
		"90p":    ninetyPct,
		"99p":    ninetyNinePct,
		"max":    maxVal,
	}
	k.Summary.Metrics.PercentWithinSLA = pctWithinSLA
	k.Summary.Metrics.PercentRevert = pctReverted
	k.Summary.Metrics.TotalTimesEligible = totalEligibleCount
	k.Summary.Metrics.TotalTimesPerformed = totalPerformed
	k.Summary.Metrics.TotalStaleReports = totalStaleReports
	k.Summary.Metrics.PercentStale = pctStale
	if k.Summary.TestInputs["BlockRange"] != nil {
		k.Summary.Metrics.AverageActualPerformsPerBlock = float64(totalPerformed) / float64(k.Summary.TestInputs["BlockRange"].(int64))
	}
	// TODO: Set test expectations
	/* Expect(int64(pctWithinSLA)).Should(BeNumerically(">=", int64(80)), "Expected PercentWithinSLA to be greater than or equal to 80, but got %f", pctWithinSLA)
	Expect(int64(pctReverted)).Should(BeNumerically("<=", int64(10)), "Expected PercentRevert to be less than or equal to 10, but got %f", pctReverted)
	Expect(k.Summary.Metrics.AverageActualPerformsPerBlock).Should(BeNumerically("~", k.Summary.Load.AverageExpectedPerformsPerBlock, 10), "Expected PercentRevert to be less than 10, but got %f", pctReverted) */

	res, err := json.MarshalIndent(k.Summary, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(k.keeperSummaryFile, res, 0600)
	if err != nil {
		return err
	}

	log.Info().Str("Summary", string(res)).Msg("Successfully wrote summary on Keeper Benchmark")

	return nil
}

// SendSlackNotification sends a slack notification on the results of the test
func (k *KeeperBenchmarkTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client, grafanaUrlProvider testreporters.GrafanaURLProvider) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := ":white_check_mark: Automation Benchmark Test FINISHED :white_check_mark:"
	if testFailed {
		headerText = ":x: Automation Benchmark Test FAILED :x:"
	}
	messageBlocks := testreporters.CommonSlackNotificationBlocks(
		headerText, k.namespace, k.keeperReportFile,
	)
	ts, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	grafanaUrl, err := grafanaUrlProvider.GetGrafanaBaseURL()
	if err != nil {
		return err
	}

	dashboardUrl, err := grafanaUrlProvider.GetGrafanaDashboardURL()
	if err != nil {
		return err
	}

	formattedDashboardURL := fmt.Sprintf("%s%s?from=%d&to=%d&var-namespace=%s&var-cl_node=chainlink-0-0", grafanaUrl, dashboardUrl, k.Summary.StartTime, k.Summary.EndTime, k.namespace)
	log.Info().Str("Dashboard", formattedDashboardURL).Msg("Dashboard URL")

	if err := testreporters.UploadSlackFile(slackClient, slack.UploadFileV2Parameters{
		Title:           fmt.Sprintf("Automation Benchmark Test Summary %s", k.namespace),
		Filename:        fmt.Sprintf("automation_benchmark_summary_%s.json", k.namespace),
		File:            k.keeperSummaryFile,
		InitialComment:  fmt.Sprintf("Automation Benchmark Test Summary %s.\n<%s|Test Dashboard> ", k.namespace, formattedDashboardURL),
		Channel:         testreporters.SlackChannel,
		ThreadTimestamp: ts,
	}); err != nil {
		return err
	}

	if err := testreporters.UploadSlackFile(slackClient, slack.UploadFileV2Parameters{
		Title:           fmt.Sprintf("Automation Benchmark Test Report %s", k.namespace),
		Filename:        fmt.Sprintf("automation_benchmark_report_%s.csv", k.namespace),
		File:            k.keeperReportFile,
		InitialComment:  fmt.Sprintf("Automation Benchmark Test Report %s", k.namespace),
		Channel:         testreporters.SlackChannel,
		ThreadTimestamp: ts,
	}); err != nil {
		return err
	}
	return testreporters.UploadSlackFile(slackClient, slack.UploadFileV2Parameters{
		Title:           fmt.Sprintf("Automation Benchmark Attempted Chainlink Txs %s", k.namespace),
		Filename:        fmt.Sprintf("attempted_cl_txs_%s.json", k.namespace),
		File:            k.attemptedTransactionsFile,
		InitialComment:  fmt.Sprintf("Automation Benchmark Attempted Txs %s", k.namespace),
		Channel:         testreporters.SlackChannel,
		ThreadTimestamp: ts,
	})
}

// intListStats helper calculates some statistics on an int list: avg, median, 90pct, 99pct, max
//
//nolint:revive
func IntListStats(in []int64) (float64, int64, int64, int64, int64) {
	length := len(in)
	if length == 0 {
		return 0, 0, 0, 0, 0
	}
	sort.Slice(in, func(i, j int) bool { return in[i] < in[j] })
	var sum int64
	for _, num := range in {
		sum += num
	}
	return float64(sum) / float64(length), in[int(math.Floor(float64(length)*0.5))], in[int(math.Floor(float64(length)*0.9))], in[int(math.Floor(float64(length)*0.99))], in[length-1]
}
