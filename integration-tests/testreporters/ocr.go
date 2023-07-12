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

// OCRSoakTestReporter collates all OCRAnswerUpdated events into a single report
type OCRSoakTestReporter struct {
	ExpectedRoundDuration time.Duration
	AnomaliesDetected     bool

	timeLine    []*OCRTestState
	namespace   string
	csvLocation string
}

type OCRTestState struct {
	Time    time.Time
	Message string
}

// SetNamespace sets the namespace of the report for clean reports
func (o *OCRSoakTestReporter) SetNamespace(namespace string) {
	o.namespace = namespace
}

func (o *OCRSoakTestReporter) RecordEvents(expectedEvents, actualEvents []*OCRTestState) {
	expectedEventIndex, actualEventIndex := 0, 0
	for expectedEventIndex < len(expectedEvents) || actualEventIndex < len(actualEvents) {
		if expectedEventIndex >= len(expectedEvents) {
			o.timeLine = append(o.timeLine, actualEvents[actualEventIndex])
			actualEventIndex++
		} else if actualEventIndex >= len(actualEvents) {
			o.timeLine = append(o.timeLine, expectedEvents[expectedEventIndex])
			expectedEventIndex++
		} else if expectedEvents[expectedEventIndex].Time.Before(actualEvents[actualEventIndex].Time) {
			o.timeLine = append(o.timeLine, expectedEvents[expectedEventIndex])
			expectedEventIndex++
		} else {
			o.timeLine = append(o.timeLine, actualEvents[actualEventIndex])
			actualEventIndex++
		}
	}
}

// WriteReport writes OCR Soak test report to logs
func (o *OCRSoakTestReporter) WriteReport(folderLocation string) error {
	log.Debug().Msg("Writing OCR Soak Test Report")
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
		headerText, o.namespace, o.csvLocation,
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
	o.csvLocation = reportLocation
	ocrReportFile, err := os.Create(reportLocation)
	if err != nil {
		return err
	}
	defer ocrReportFile.Close()

	ocrReportWriter := csv.NewWriter(ocrReportFile)

	err = ocrReportWriter.Write([]string{
		"Time",
		"Message",
	})
	if err != nil {
		return err
	}

	for _, event := range o.timeLine {
		err = ocrReportWriter.Write([]string{
			event.Time.String(),
			event.Message,
		})
		if err != nil {
			return err
		}
	}

	ocrReportWriter.Flush()

	log.Info().Str("Location", reportLocation).Msg("Wrote CSV file")
	return nil
}
