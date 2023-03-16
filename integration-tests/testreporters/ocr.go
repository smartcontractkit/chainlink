package testreporters

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

// OCRSoakTestReporter collates all OCRAnswerUpdated events into a single report
type OCRSoakTestReporter struct {
	ContractReports       map[string]*OCRSoakTestReport // contractAddress: Answers
	ExpectedRoundDuration time.Duration
	AnomaliesDetected     bool

	namespace   string
	csvLocation string
}

// SetNamespace sets the namespace of the report for clean reports
func (o *OCRSoakTestReporter) SetNamespace(namespace string) {
	o.namespace = namespace
}

// WriteReport writes OCR Soak test report to logs
func (o *OCRSoakTestReporter) WriteReport(folderLocation string) error {
	log.Debug().Msg("Writing OCR Soak Test Report")
	var reportGroup sync.WaitGroup
	for _, report := range o.ContractReports {
		reportGroup.Add(1)
		go func(report *OCRSoakTestReport) {
			defer reportGroup.Done()
			if report.ProcessOCRReport() {
				o.AnomaliesDetected = true
			}
		}(report)
	}
	reportGroup.Wait()
	log.Debug().Int("Count", len(o.ContractReports)).Msg("Processed OCR Soak Test Reports")
	return o.writeCSV(folderLocation)
}

// SendNotification sends a slack message to a slack webhook and uploads test artifacts
func (o *OCRSoakTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := ":white_check_mark: OCR Soak Test PASSED :white_check_mark:"
	if testFailed {
		headerText = ":x: OCR Soak Test FAILED :x:"
	} else if o.AnomaliesDetected {
		headerText = ":warning: OCR Soak Test Found Anomalies :warning:"
	}
	messageBlocks := testreporters.CommonSlackNotificationBlocks(
		t, slackClient, headerText, o.namespace, o.csvLocation, testreporters.SlackUserID, testFailed,
	)
	ts, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	return testreporters.UploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("OCR Soak Test Report %s", o.namespace),
		Filetype:        "csv",
		Filename:        fmt.Sprintf("ocr_soak_%s.csv", o.namespace),
		File:            o.csvLocation,
		InitialComment:  fmt.Sprintf("OCR Soak Test Report %s.", o.namespace),
		Channels:        []string{testreporters.SlackChannel},
		ThreadTimestamp: ts,
	})
}

// writes a CSV report on the test runner
func (o *OCRSoakTestReporter) writeCSV(folderLocation string) error {
	reportLocation := filepath.Join(folderLocation, "./ocr_soak_report.csv")
	log.Debug().Str("Location", reportLocation).Msg("Writing OCR report")
	o.csvLocation = reportLocation
	ocrReportFile, err := os.Create(reportLocation)
	if err != nil {
		return err
	}
	defer ocrReportFile.Close()

	ocrReportWriter := csv.NewWriter(ocrReportFile)

	err = ocrReportWriter.Write([]string{
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
	for contractAddress, report := range o.ContractReports {
		err = ocrReportWriter.Write([]string{
			contractAddress,
			fmt.Sprint(report.totalRounds),
			report.averageRoundTime.Truncate(time.Second).String(),
			report.longestRoundTime.Truncate(time.Second).String(),
			report.shortestRoundTime.Truncate(time.Second).String(),
			fmt.Sprint(report.averageRoundBlocks),
			fmt.Sprint(report.longestRoundBlocks),
			fmt.Sprint(report.shortestRoundBlocks),
		})
		if err != nil {
			return err
		}
	}

	err = ocrReportWriter.Write([]string{})
	if err != nil {
		return err
	}

	// Anomalous reports
	err = ocrReportWriter.Write([]string{"Updates With Anomalies"})
	if err != nil {
		return err
	}

	for _, report := range o.ContractReports {
		if len(report.AnomalousAnswerIndexes) > 0 {
			err = ocrReportWriter.Write(report.csvHeaders())
			if err != nil {
				return err
			}
		}
		for _, values := range report.anomalousCSVValues() {
			err = ocrReportWriter.Write(values)
			if err != nil {
				return err
			}
		}
	}

	err = ocrReportWriter.Write([]string{})
	if err != nil {
		return err
	}

	// All reports
	err = ocrReportWriter.Write([]string{"All Updated Answers"})
	if err != nil {
		return err
	}

	for _, report := range o.ContractReports {
		if len(report.UpdatedAnswers) > 0 {
			err = ocrReportWriter.Write(report.csvHeaders())
			if err != nil {
				return err
			}
		}
		for _, values := range report.allCSVValues() {
			err = ocrReportWriter.Write(values)
			if err != nil {
				return err
			}
		}
	}

	ocrReportWriter.Flush()

	log.Info().Str("Location", reportLocation).Msg("Wrote CSV file")
	return nil
}

// OCRSoakTestReport holds all answered rounds and summary data for an OCR contract
type OCRSoakTestReport struct {
	ContractAddress        string
	UpdatedAnswers         []*OCRAnswerUpdated
	AnomalousAnswerIndexes []int
	ExpectedRoundDuration  time.Duration

	newRoundExpected       bool
	newRoundExpectedId     uint64
	newRoundExpectedAnswer int
	newRoundStartTime      time.Time
	newRoundStartBlock     uint64

	totalRounds         uint64
	longestRoundTime    time.Duration
	shortestRoundTime   time.Duration
	averageRoundTime    time.Duration
	longestRoundBlocks  uint64
	shortestRoundBlocks uint64
	averageRoundBlocks  uint64
}

// NewOCRSoakTestReport initializes a new soak test report for a new tests
func NewOCRSoakTestReport(contractAddress string, startingAnswer int, expectedRoundDuration time.Duration) *OCRSoakTestReport {
	return &OCRSoakTestReport{
		ContractAddress:        contractAddress,
		UpdatedAnswers:         make([]*OCRAnswerUpdated, 0),
		AnomalousAnswerIndexes: make([]int, 0),
		ExpectedRoundDuration:  expectedRoundDuration,
		newRoundExpected:       true,
		newRoundExpectedId:     1,
		newRoundExpectedAnswer: startingAnswer,
	}
}

// ProcessOCRReport summarizes all data collected from OCR rounds, and returns if there are any anomalies detected
func (o *OCRSoakTestReport) ProcessOCRReport() bool {
	log.Debug().Str("OCR Address", o.ContractAddress).Msg("Processing OCR Soak Report")
	o.AnomalousAnswerIndexes = make([]int, 0)

	var (
		totalRoundBlocks uint64
		totalRoundTime   time.Duration
	)

	o.longestRoundTime = 0
	o.shortestRoundTime = math.MaxInt64
	o.longestRoundBlocks = 0
	o.shortestRoundBlocks = math.MaxUint64
	for index, updatedAnswer := range o.UpdatedAnswers {
		if updatedAnswer.ProcessAnomalies(o.ExpectedRoundDuration) {
			o.AnomalousAnswerIndexes = append(o.AnomalousAnswerIndexes, index)
		}
		if !updatedAnswer.Anomalous { // Anomalous answers can have outlier values that throw averages off
			o.totalRounds++
			updatedAnswer.RoundDuration = updatedAnswer.UpdatedTime.Sub(updatedAnswer.StartingTime)
			updatedAnswer.BlockDuration = updatedAnswer.UpdatedBlockNum - updatedAnswer.StartingBlockNum
			totalRoundTime += updatedAnswer.RoundDuration
			totalRoundBlocks += updatedAnswer.BlockDuration
			if o.longestRoundTime < updatedAnswer.RoundDuration {
				o.longestRoundTime = updatedAnswer.RoundDuration
			}
			if o.shortestRoundTime > updatedAnswer.RoundDuration {
				o.shortestRoundTime = updatedAnswer.RoundDuration
			}
			if o.longestRoundBlocks < updatedAnswer.BlockDuration {
				o.longestRoundBlocks = updatedAnswer.BlockDuration
			}
			if o.shortestRoundBlocks > updatedAnswer.BlockDuration {
				o.shortestRoundBlocks = updatedAnswer.BlockDuration
			}
		}
	}
	o.averageRoundBlocks = totalRoundBlocks / o.totalRounds
	o.averageRoundTime = totalRoundTime / time.Duration(o.totalRounds)
	return len(o.AnomalousAnswerIndexes) > 0
}

// NewAnswerUpdated records a new round, updating expectations for the next one if necessary. Returns true if the answer
// comes in as fully expected.
func (o *OCRSoakTestReport) NewAnswerUpdated(newAnswer *OCRAnswerUpdated) bool {
	fullyExpected := newAnswer.UpdatedRoundId == o.newRoundExpectedId &&
		o.newRoundExpected &&
		o.newRoundExpectedAnswer == newAnswer.UpdatedAnswer
	newAnswer.ContractAddress = o.ContractAddress
	newAnswer.ExpectedAnswer = o.newRoundExpectedAnswer
	newAnswer.ExpectedUpdate = o.newRoundExpected
	newAnswer.ExpectedRoundId = o.newRoundExpectedId
	if newAnswer.UpdatedRoundId >= o.newRoundExpectedId {
		o.newRoundExpectedId = newAnswer.UpdatedRoundId + 1
	}

	if fullyExpected { // Expected round came in correctly
		newAnswer.StartingBlockNum = o.newRoundStartBlock
		newAnswer.StartingTime = o.newRoundStartTime
		o.newRoundExpected = false
	}

	o.UpdatedAnswers = append(o.UpdatedAnswers, newAnswer)
	log.Info().
		Uint64("Updated Round ID", newAnswer.UpdatedRoundId).
		Uint64("Expected Round ID", newAnswer.ExpectedRoundId).
		Int("Updated Answer", newAnswer.UpdatedAnswer).
		Int("Expected Answer", newAnswer.ExpectedAnswer).
		Str("Address", o.ContractAddress).
		Uint64("Block Number", newAnswer.UpdatedBlockNum).
		Str("Block Hash", newAnswer.UpdatedBlockHash).
		Str("Event Tx Hash", newAnswer.RoundTxHash).
		Msg("Answer Updated")
	return fullyExpected
}

// NewAnswerExpected indicates that we're expecting a new answer on an OCR contract
func (o *OCRSoakTestReport) NewAnswerExpected(answer int, startingBlock uint64) {
	o.newRoundExpected = true
	o.newRoundExpectedAnswer = answer
	o.newRoundStartTime = time.Now()
	o.newRoundStartBlock = startingBlock
	log.Debug().
		Str("Address", o.ContractAddress).
		Uint64("Expected Round ID", o.newRoundExpectedId).
		Int("Expected Answer", o.newRoundExpectedAnswer).
		Msg("Expecting a New OCR Round")
}

func (o *OCRSoakTestReport) csvHeaders() []string {
	return []string{
		"Contract Address",
		"Update Expected?",
		"Expected Round ID",
		"Event Round ID",
		"On-Chain Round ID",
		"Round Start Time",
		"Round End Time",
		"Round Duration",
		"Round Triggered Block Number",
		"Round Updated Block Number",
		"Round Block Duration",
		"Round Updated Block Hash",
		"Round Tx Hash",
		"Expected Answer",
		"Event Answer",
		"On-Chain Answer",
		"Anomalous?",
		"Anomalies",
	}
}

// returns CSV formatted values for all updated answers
func (o *OCRSoakTestReport) allCSVValues() [][]string {
	csvValues := [][]string{}
	for _, updatedAnswer := range o.UpdatedAnswers {
		csvValues = append(csvValues, updatedAnswer.toCSV())
	}
	return csvValues
}

// returns all CSV formatted values of anomalous answers
func (o *OCRSoakTestReport) anomalousCSVValues() [][]string {
	csvValues := [][]string{}
	for _, anomalousIndex := range o.AnomalousAnswerIndexes {
		updatedAnswer := o.UpdatedAnswers[anomalousIndex]
		csvValues = append(csvValues, updatedAnswer.toCSV())
	}
	return csvValues
}

// OCRAnswerUpdated records details of an OCRAnswerUpdated event and compares them against expectations
type OCRAnswerUpdated struct {
	// metadata
	ContractAddress  string
	ExpectedUpdate   bool
	StartingBlockNum uint64
	UpdatedBlockHash string
	RoundTxHash      string
	BlockDuration    uint64
	StartingTime     time.Time
	RoundDuration    time.Duration

	// round data
	ExpectedRoundId uint64
	ExpectedAnswer  int

	UpdatedRoundId  uint64
	UpdatedBlockNum uint64
	UpdatedTime     time.Time
	UpdatedAnswer   int

	OnChainRoundId uint64
	OnChainAnswer  int

	Anomalous bool
	Anomalies []string
}

// ProcessAnomalies checks received data against expected data of the updated answer, returning if anything mismatches
func (o *OCRAnswerUpdated) ProcessAnomalies(expectedRoundDuration time.Duration) bool {
	if o.UpdatedRoundId == 0 && o.OnChainRoundId == 0 {
		o.Anomalous = true
		o.Anomalies = []string{fmt.Sprintf("Test likely ended before the round could be confirmed. Check for round %d on chain", o.ExpectedRoundId)}
		return o.Anomalous
	}

	var isAnomaly bool
	anomalies := []string{}
	if !o.ExpectedUpdate {
		isAnomaly = true
		anomalies = append(anomalies, "Unexpected new round, possible double transmission")
	}

	if o.ExpectedRoundId != o.UpdatedRoundId || o.ExpectedRoundId != o.OnChainRoundId {
		isAnomaly = true
		anomalies = append(anomalies, "RoundID mismatch, possible double transmission")
	}
	if o.ExpectedAnswer != o.UpdatedAnswer || o.ExpectedAnswer != o.OnChainAnswer {
		isAnomaly = true
		anomalies = append(anomalies, "! ANSWER MISMATCH !")
	}
	if o.RoundDuration > expectedRoundDuration {
		isAnomaly = true
		anomalies = append(anomalies, fmt.Sprintf(
			"Round took %s to complete, longer than expected time of %s", o.RoundDuration, expectedRoundDuration.String()),
		)
	}
	o.Anomalous, o.Anomalies = isAnomaly, anomalies
	return isAnomaly
}

func (o *OCRAnswerUpdated) toCSV() []string {
	var ( // Values that could be affected by anomalies
		startTime          string
		roundTimeDuration  string
		startingBlock      string
		roundBlockDuration string
	)

	if o.StartingTime.IsZero() {
		startTime = "Unknown"
		roundTimeDuration = "Unknown"
	} else {
		startTime = o.StartingTime.Truncate(time.Second).String()
		roundTimeDuration = o.UpdatedTime.Sub(o.StartingTime).Truncate(time.Second).String()
	}

	if o.StartingBlockNum == 0 {
		startingBlock = "Unknown"
		roundBlockDuration = "Unknown"
	} else {
		startingBlock = fmt.Sprint(o.StartingBlockNum)
		roundBlockDuration = fmt.Sprint(o.UpdatedBlockNum - o.StartingBlockNum)
	}

	return []string{
		o.ContractAddress,
		fmt.Sprint(o.ExpectedUpdate),
		fmt.Sprint(o.ExpectedRoundId),
		fmt.Sprint(o.UpdatedRoundId),
		fmt.Sprint(o.OnChainRoundId),
		startTime,
		o.UpdatedTime.Truncate(time.Second).String(),
		roundTimeDuration,
		startingBlock,
		fmt.Sprint(o.UpdatedBlockNum),
		roundBlockDuration,
		o.UpdatedBlockHash,
		o.RoundTxHash,
		fmt.Sprint(o.ExpectedAnswer),
		fmt.Sprint(o.UpdatedAnswer),
		fmt.Sprint(o.OnChainAnswer),
		fmt.Sprint(o.Anomalous),
		strings.Join(o.Anomalies, " | "),
	}
}
