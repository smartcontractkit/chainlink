package testreporters

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

type VRFV2SoakTestReporter struct {
	Reports     map[string]*VRFV2SoakTestReport // contractAddress: Report
	namespace   string
	csvLocation string
}

type VRFV2SoakTestReport struct {
	ContractAddress string
	TotalRounds     uint

	averageRoundTime  time.Duration
	LongestRoundTime  time.Duration
	ShortestRoundTime time.Duration
	totalRoundTimes   time.Duration

	averageRoundBlocks  uint
	LongestRoundBlocks  uint
	ShortestRoundBlocks uint
	totalBlockLengths   uint
}

// SetNamespace sets the namespace of the report for clean reports
func (o *VRFV2SoakTestReporter) SetNamespace(namespace string) {
	o.namespace = namespace
}

// WriteReport writes VRFV2 Soak test report to logs
func (o *VRFV2SoakTestReporter) WriteReport(folderLocation string) error {
	for _, report := range o.Reports {
		report.averageRoundBlocks = report.totalBlockLengths / report.TotalRounds
		report.averageRoundTime = time.Duration(report.totalRoundTimes.Nanoseconds() / int64(report.TotalRounds))
	}
	if err := o.writeCSV(folderLocation); err != nil {
		return err
	}

	log.Info().Msg("VRFV2 Soak Test Report")
	log.Info().Msg("--------------------")
	for contractAddress, report := range o.Reports {
		log.Info().
			Str("Contract Address", report.ContractAddress).
			Uint("Total Rounds Processed", report.TotalRounds).
			Str("Average Round Time", fmt.Sprint(report.averageRoundTime)).
			Str("Longest Round Time", fmt.Sprint(report.LongestRoundTime)).
			Str("Shortest Round Time", fmt.Sprint(report.ShortestRoundTime)).
			Uint("Average Round Blocks", report.averageRoundBlocks).
			Uint("Longest Round Blocks", report.LongestRoundBlocks).
			Uint("Shortest Round Blocks", report.ShortestRoundBlocks).
			Msg(contractAddress)
	}
	log.Info().Msg("--------------------")
	return nil
}

// SendNotification sends a slack message to a slack webhook and uploads test artifacts
func (o *VRFV2SoakTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := ":white_check_mark: VRFV2 Soak Test PASSED :white_check_mark:"
	if testFailed {
		headerText = ":x: VRFV2 Soak Test FAILED :x:"
	}
	messageBlocks := testreporters.CommonSlackNotificationBlocks(
		t, slackClient, headerText, o.namespace, o.csvLocation, testreporters.SlackUserID, testFailed,
	)
	ts, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	return testreporters.UploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("VRFV2 Soak Test Report %s", o.namespace),
		Filetype:        "csv",
		Filename:        fmt.Sprintf("vrfv2_soak_%s.csv", o.namespace),
		File:            o.csvLocation,
		InitialComment:  fmt.Sprintf("VRFV2 Soak Test Report %s.", o.namespace),
		Channels:        []string{testreporters.SlackChannel},
		ThreadTimestamp: ts,
	})
}

// UpdateReport updates the report based on the latest info
func (o *VRFV2SoakTestReport) UpdateReport(roundTime time.Duration, blockLength uint) {
	// Updates min values from default 0
	if o.ShortestRoundBlocks == 0 {
		o.ShortestRoundBlocks = blockLength
	}
	if o.ShortestRoundTime == 0 {
		o.ShortestRoundTime = roundTime
	}
	o.TotalRounds++
	o.totalRoundTimes += roundTime
	o.totalBlockLengths += blockLength
	if roundTime >= o.LongestRoundTime {
		o.LongestRoundTime = roundTime
	}
	if roundTime <= o.ShortestRoundTime {
		o.ShortestRoundTime = roundTime
	}
	if blockLength >= o.LongestRoundBlocks {
		o.LongestRoundBlocks = blockLength
	}
	if blockLength <= o.ShortestRoundBlocks {
		o.ShortestRoundBlocks = blockLength
	}
}

// writes a CSV report on the test runner
func (o *VRFV2SoakTestReporter) writeCSV(folderLocation string) error {
	reportLocation := filepath.Join(folderLocation, "./vrfv2_soak_report.csv")
	log.Debug().Str("Location", reportLocation).Msg("Writing VRFV2 report")
	o.csvLocation = reportLocation
	vrfv2ReportFile, err := os.Create(reportLocation)
	if err != nil {
		return err
	}
	defer vrfv2ReportFile.Close()

	vrfv2ReportWriter := csv.NewWriter(vrfv2ReportFile)
	err = vrfv2ReportWriter.Write([]string{
		"Contract Index",
		"Contract Address",
		"Total Rounds Processed",
		"Average Round Time",
		"Longest Round Time",
		"Shortest Round Time",
		"Average Round Blocks",
		"Longest Round Blocks",
		"Shortest Round Blocks",
	})
	if err != nil {
		return err
	}
	for contractIndex, report := range o.Reports {
		err = vrfv2ReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalRounds),
			fmt.Sprint(report.averageRoundTime),
			fmt.Sprint(report.LongestRoundTime),
			fmt.Sprint(report.ShortestRoundTime),
			fmt.Sprint(report.averageRoundBlocks),
			fmt.Sprint(report.LongestRoundBlocks),
			fmt.Sprint(report.ShortestRoundBlocks),
		})
		if err != nil {
			return err
		}
	}
	vrfv2ReportWriter.Flush()

	log.Info().Str("Location", reportLocation).Msg("Wrote CSV file")
	return nil
}
