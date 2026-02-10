package products

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	dfilter "github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// code mostly copied from CTFv1's 'lib/testreporters' package to avoid depending on it
type WarnAboutAllowedMsgs = bool

const (
	WarnAboutAllowedMsgs_Yes WarnAboutAllowedMsgs = true  //nolint: revive //we feel like using underscores
	WarnAboutAllowedMsgs_No  WarnAboutAllowedMsgs = false //nolint: revive //we feel like using underscores
)

// AllowedLogMessage is a log message that might be thrown by a Chainlink node during a test, but is not a concern
type AllowedLogMessage struct {
	message      string
	reason       string
	level        zapcore.Level
	logWhenFound WarnAboutAllowedMsgs
}

// NewAllowedLogMessage creates a new AllowedLogMessage. If logWhenFound is true, the log message will be printed to the
// console when found in the log file with Warn level (this can get noisy).
func NewAllowedLogMessage(message string, reason string, level zapcore.Level, logWhenFound WarnAboutAllowedMsgs) AllowedLogMessage {
	return AllowedLogMessage{
		message:      message,
		reason:       reason,
		level:        level,
		logWhenFound: logWhenFound,
	}
}

var defaultAllowedLogMessages = []AllowedLogMessage{
	{
		message: "No EVM primary nodes available: 0/1 nodes are alive",
		reason:  "Sometimes geth gets unlucky in the start up process and the Chainlink node starts before geth is ready",
		level:   zapcore.DPanicLevel,
	},
}

type ChainlinkNodeLogScannerSettings struct {
	FailingLogLevel zapcore.Level
	Threshold       uint
	AllowedMessages []AllowedLogMessage
}

func ScanLogs(l zerolog.Logger, settings ChainlinkNodeLogScannerSettings) error {
	logStream, lErr := framework.StreamContainerLogs(container.ListOptions{
		All: true,
		Filters: dfilter.NewArgs(dfilter.KeyValuePair{
			Key:   "label",
			Value: "framework=ctf",
		}),
	}, container.LogsOptions{ShowStdout: true, ShowStderr: true})

	if lErr != nil {
		return lErr
	}

	verifyLogsGroup := &errgroup.Group{}
	for _, stream := range logStream {
		verifyLogsGroup.Go(func() error {
			verifyErr := verifyLogStream(
				stream,
				settings.FailingLogLevel,
				settings.Threshold,
				settings.AllowedMessages...,
			)
			// ignore processing errors
			if verifyErr != nil && !strings.Contains(verifyErr.Error(), MultipleLogsAtLogLevelErr) &&
				!strings.Contains(verifyErr.Error(), OneLogAtLogLevelErr) {
				l.Error().Err(verifyErr).Msg("Error processing CL node logs")

				return nil

				// if it's not a processing error, we want to fail the test; we also can stop processing logs all together at this point
			} else if verifyErr != nil &&
				(strings.Contains(verifyErr.Error(), MultipleLogsAtLogLevelErr) ||
					strings.Contains(verifyErr.Error(), OneLogAtLogLevelErr)) {
				return verifyErr
			}
			return nil
		})
	}

	if logVerificationErr := verifyLogsGroup.Wait(); logVerificationErr != nil {
		return fmt.Errorf("found a concerning log in Chainlink Node logs: %w", logVerificationErr)
	}

	l.Info().Msg("Found no concerning entries in Chainlink Node logs")

	return nil
}

var (
	OneLogAtLogLevelErr       = "found log at level"
	MultipleLogsAtLogLevelErr = "found too many logs at level"
)

// verifyLogStream verifies that a log stream does not contain any logs at a level higher than the failingLogLevel. If it does,
// it will return an error. It also allows for a list of AllowedLogMessages to be passed in, which will be ignored if found
// in the log file. The failureThreshold is the number of logs at the failingLogLevel or higher that can be found before
// the function returns an error.
func verifyLogStream(stream io.ReadCloser, failingLogLevel zapcore.Level, failureThreshold uint, allowedMessages ...AllowedLogMessage) error {
	defer func(stream io.ReadCloser) {
		_ = stream.Close()
	}(stream)

	var err error
	scanner := bufio.NewScanner(stream)
	scanner.Split(bufio.ScanLines)

	allAllowedMessages := append([]AllowedLogMessage{}, defaultAllowedLogMessages...)
	allAllowedMessages = append(allAllowedMessages, allowedMessages...)

	var logsFound uint

	for scanner.Scan() {
		jsonLogLine := scanner.Text()
		// Docker API streams logs with multiplex header like this: \x02\x00\x00\x00\x00\x00\x00\x92
		// we need to strip it before processing logs
		idx := strings.IndexByte(jsonLogLine, '{')
		if idx >= 0 {
			jsonLogLine = jsonLogLine[idx:]
		}
		logsFound, err = scanLogLine(log.Logger, jsonLogLine, failingLogLevel, logsFound, failureThreshold, allAllowedMessages)
		if err != nil {
			return err
		}
	}
	return nil
}

// scanLogLine scans a log line for a failing log level, returning the number of failing logs found so far. It returns an error if the failure threshold is reached or if any panic is found
// or if there's no log level found. It also takes a list of allowed messages that are ignored if found.
func scanLogLine(log zerolog.Logger, jsonLogLine string, failingLogLevel zapcore.Level, foundSoFar, failureThreshold uint, allowedMessages []AllowedLogMessage) (uint, error) {
	var zapLevel zapcore.Level
	var err error

	if !strings.HasPrefix(jsonLogLine, "{") { // don't bother with non-json lines
		if strings.HasPrefix(jsonLogLine, "panic") { // unless it's a panic
			return 0, fmt.Errorf("found panic: %s", jsonLogLine)
		}
		return foundSoFar, nil
	}
	jsonMapping := map[string]any{}

	if err = json.Unmarshal([]byte(jsonLogLine), &jsonMapping); err != nil {
		// This error can occur anytime someone uses %+v in a log message, ignoring
		return foundSoFar, nil
	}
	logLevel, ok := jsonMapping["level"].(string)
	if !ok {
		return 0, fmt.Errorf("found no log level in chainlink log line: %s", jsonLogLine)
	}

	if logLevel == "crit" { // "crit" is a custom core type they map to DPanic
		zapLevel = zapcore.DPanicLevel
	} else {
		zapLevel, err = zapcore.ParseLevel(logLevel)
		if err != nil {
			return 0, fmt.Errorf("'%s' not a valid zapcore level", logLevel)
		}
	}

	if zapLevel >= failingLogLevel {
		logErr := fmt.Errorf("%s '%s', failing any log level higher than %s: %s", OneLogAtLogLevelErr, logLevel, zapLevel.String(), jsonLogLine)
		if failureThreshold > 1 {
			logErr = fmt.Errorf("%s '%s' or above; failure threshold of %d reached; last error found: %s", MultipleLogsAtLogLevelErr, logLevel, failureThreshold, jsonLogLine)
		}
		logMessage, hasMessage := jsonMapping["msg"]
		if !hasMessage {
			foundSoFar++
			if foundSoFar >= failureThreshold {
				return foundSoFar, logErr
			}
			return foundSoFar, nil
		}

		for _, allowedLog := range allowedMessages {
			if strings.Contains(logMessage.(string), allowedLog.message) {
				if allowedLog.logWhenFound {
					log.Warn().
						Str("Reason", allowedLog.reason).
						Str("Level", allowedLog.level.CapitalString()).
						Str("Msg", logMessage.(string)).
						Msg("Found allowed log message, ignoring")
				}

				return foundSoFar, nil
			}
		}

		foundSoFar++
		if foundSoFar >= failureThreshold {
			return foundSoFar, logErr
		}
	}

	return foundSoFar, nil
}

var defaultAllowedMessages = []AllowedLogMessage{
	NewAllowedLogMessage("Failed to get LINK balance", "Happens only when we deploy LINK token for test purposes. Harmless.", zapcore.ErrorLevel, WarnAboutAllowedMsgs_No),
	NewAllowedLogMessage("Error stopping job service", "It's a known issue with lifecycle. There's ongoing work that will fix it.", zapcore.DPanicLevel, WarnAboutAllowedMsgs_No),
	NewAllowedLogMessage(
		"No live RPC nodes available",
		"Networking or infra issues can cause brief disconnections from the node to RPC nodes, especially at startup. This isn't a concern as long as the test passes otherwise",
		zapcore.DPanicLevel,
		WarnAboutAllowedMsgs_Yes,
	),
}

var defaultSettings = ChainlinkNodeLogScannerSettings{
	FailingLogLevel: zapcore.DPanicLevel,
	Threshold:       1, // we want to fail on the first concerning log
	AllowedMessages: defaultAllowedMessages,
}

func DefaultSettings(extraAllowedMessages ...AllowedLogMessage) ChainlinkNodeLogScannerSettings {
	allowedMessages := append([]AllowedLogMessage{}, defaultAllowedMessages...)
	allowedMessages = append(allowedMessages, extraAllowedMessages...)
	return ChainlinkNodeLogScannerSettings{
		FailingLogLevel: defaultSettings.FailingLogLevel,
		Threshold:       defaultSettings.Threshold,
		AllowedMessages: allowedMessages,
	}
}
