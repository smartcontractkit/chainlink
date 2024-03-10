package testreporters

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

//TODO: This whole process can definitely be simplified and improved, but for some reason I'm getting brain block at the moment

// OCRSoakTestReporter collates all OCRAnswerUpdated events into a single report
type OCRSoakTestReporter struct {
	StartTime         time.Time
	AnomaliesDetected bool
	OCRVersion        string

	anomalies   [][]string
	timeLine    [][]string
	namespace   string
	csvLocation string
}

// TimeLineEvent represents a single event in the timeline
type TimeLineEvent interface {
	Time() time.Time
	CSV() [][]string
}

// TestIssue is a single RPC issue, either a disconnect or reconnect
type TestIssue struct {
	StartTime time.Time `toml:"startTime"`
	Message   string    `toml:"message"`
}

func (r *TestIssue) Time() time.Time {
	return r.StartTime
}

func (r *TestIssue) CSV() [][]string {
	return [][]string{{r.StartTime.Format("2006-01-02 15:04:05.00 MST"), "Test Issue!", r.Message}}
}

// OCRRoundState indicates that a round per contract should complete within this time with this answer
type OCRRoundState struct {
	StartTime      time.Time                `toml:"startTime"`
	EndTime        time.Time                `toml:"endTime"` // Time when the round should end, only used for analysis
	Answer         int64                    `toml:"answer"`
	Anomalous      bool                     `toml:"anomalous"`   // Whether the round was anomalous
	FoundEvents    map[string][]*FoundEvent `toml:"foundEvents"` // Address -> FoundEvents, possible to have multiple found events per round, and need to call it out
	TimeLineEvents []TimeLineEvent          `toml:"timeLineEvents"`

	anomalies [][]string
}

func (e *OCRRoundState) Time() time.Time {
	return e.StartTime
}

// CSV returns a CSV representation of the test state and all events
func (e *OCRRoundState) CSV() [][]string {
	rows := [][]string{{e.StartTime.Format("2006-01-02 15:04:05.00 MST"), fmt.Sprintf("Expecting new Answer: %d", e.Answer)}}
	rows = append(rows, e.anomalies...)
	return rows
}

// Validate checks that
// 1. There is a FoundEvent for every address
// 2. There is only one FoundEvent for every address
// 3. The answer is correct
func (e *OCRRoundState) Validate() bool {
	anomalies := [][]string{}
	for address, eventList := range e.FoundEvents {
		if len(eventList) == 0 {
			e.Anomalous = true
			anomalies = append(anomalies, []string{
				e.StartTime.Format("2006-01-02 15:04:05.00 MST"), "Anomaly Found!", fmt.Sprintf("No AnswerUpdated for address '%s'", address),
			})
		} else if len(eventList) > 1 {
			e.Anomalous = true
			anomalies = append(anomalies, []string{e.StartTime.Format("2006-01-02 15:04:05.00 MST"), "Anomaly Found!",
				fmt.Sprintf("Multiple AnswerUpdated for address '%s', possible double-transmission", address)},
			)
		} else {
			event := eventList[0]
			if event.Answer != e.Answer {
				e.Anomalous = true
				anomalies = append(e.anomalies, []string{e.StartTime.Format("2006-01-02 15:04:05.00 MST"), "Anomaly Found!",
					fmt.Sprintf("FoundEvent for address '%s' has wrong answer '%d'", address, event.Answer)},
				)
			}
		}
	}
	e.anomalies = anomalies
	return e.Anomalous
}

// FoundEvent is a single round update event
type FoundEvent struct {
	StartTime   time.Time
	BlockNumber uint64
	Address     string
	Answer      int64
	RoundID     uint64
}

func (a *FoundEvent) Time() time.Time {
	return a.StartTime
}

// CSV returns a CSV representation of the event
func (a *FoundEvent) CSV() [][]string {
	return [][]string{{
		a.StartTime.Format("2006-01-02 15:04:05.00 MST"),
		fmt.Sprintf("Address: %s", a.Address),
		fmt.Sprintf("Round: %d", a.RoundID),
		fmt.Sprintf("Answer: %d", a.Answer),
		fmt.Sprintf("Block: %d", a.BlockNumber),
	}}
}

// RecordEvents takes in a list of test states and RPC issues, orders them, and records them in the timeline
func (o *OCRSoakTestReporter) RecordEvents(testStates []*OCRRoundState, testIssues []*TestIssue) {
	events := []TimeLineEvent{}
	for _, expectedEvent := range testStates {
		if expectedEvent.Validate() {
			o.AnomaliesDetected = true
			o.anomalies = append(o.anomalies, expectedEvent.anomalies...)
		}
		events = append(events, expectedEvent)
		events = append(events, expectedEvent.TimeLineEvents...)
	}
	if len(testIssues) > 0 {
		o.AnomaliesDetected = true
	}
	for _, testIssue := range testIssues {
		events = append(events, testIssue)
		o.anomalies = append(o.anomalies, testIssue.CSV()...)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time().Before(events[j].Time())
	})
	for _, event := range events {
		o.timeLine = append(o.timeLine, event.CSV()...)
	}
}

// SetNamespace sets the namespace of the report for clean reports
func (o *OCRSoakTestReporter) SetNamespace(namespace string) {
	o.namespace = namespace
}

// WriteReport writes OCR Soak test report to a CSV file and final report
func (o *OCRSoakTestReporter) WriteReport(folderLocation string) error {
	log.Debug().Msgf("Writing OCRv%s Soak Test Report", o.OCRVersion)
	reportLocation := filepath.Join(folderLocation, "./ocr_soak_report.csv")
	o.csvLocation = reportLocation
	ocrReportFile, err := os.Create(reportLocation)
	if err != nil {
		return err
	}
	defer ocrReportFile.Close()

	ocrReportWriter := csv.NewWriter(ocrReportFile)

	err = ocrReportWriter.Write([]string{fmt.Sprintf("OCRv%s Soak Test Report", o.OCRVersion)})
	if err != nil {
		return err
	}

	err = ocrReportWriter.Write([]string{
		"Namespace",
		o.namespace,
		"Started At",
		o.StartTime.Format("2006-01-02 15:04:05.00 MST"),
		"Test Duration",
		time.Since(o.StartTime).String(),
	})
	if err != nil {
		return err
	}

	err = ocrReportWriter.Write([]string{})
	if err != nil {
		return err
	}

	if len(o.anomalies) > 0 {
		err = ocrReportWriter.Write([]string{"Anomalies Found"})
		if err != nil {
			return err
		}

		err = ocrReportWriter.WriteAll(o.anomalies)
		if err != nil {
			return err
		}

		err = ocrReportWriter.Write([]string{})
		if err != nil {
			return err
		}
	}

	err = ocrReportWriter.Write([]string{"Timeline"})
	if err != nil {
		return err
	}

	err = ocrReportWriter.Write([]string{
		"Time",
		"Message",
	})
	if err != nil {
		return err
	}

	err = ocrReportWriter.WriteAll(o.timeLine)
	if err != nil {
		return err
	}

	ocrReportWriter.Flush()

	log.Info().Str("Location", reportLocation).Msg("Wrote CSV file")
	return nil
}

// SendNotification sends a slack message to a slack webhook and uploads test artifacts
func (o *OCRSoakTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client, _ testreporters.GrafanaURLProvider) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := fmt.Sprintf(":white_check_mark: OCRv%s Soak Test PASSED :white_check_mark:", o.OCRVersion)
	if testFailed {
		headerText = ":x: OCR Soak Test FAILED :x:"
	} else if o.AnomaliesDetected {
		headerText = ":warning: OCR Soak Test Found Anomalies :warning:"
	}
	messageBlocks := testreporters.CommonSlackNotificationBlocks(
		headerText, fmt.Sprintf("%s | Test took: %s", o.namespace, time.Since(o.StartTime).Truncate(time.Second).String()), o.csvLocation,
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
