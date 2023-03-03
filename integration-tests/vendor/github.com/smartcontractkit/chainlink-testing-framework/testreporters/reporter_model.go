package testreporters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

// TestReporter is a general interface for all test reporters
type TestReporter interface {
	WriteReport(folderLocation string) error
	SendSlackNotification(t *testing.T, slackClient *slack.Client) error
	SetNamespace(namespace string)
}

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir string = "logs"
)

// WriteTeardownLogs attempts to download the logs of all ephemeral test deployments onto the test runner, also writing
// a test report if one is provided. A failing log level also enables you to fail a test based on what level logs the
// Chainlink nodes have thrown during their test.
func WriteTeardownLogs(
	t *testing.T,
	env *environment.Environment,
	optionalTestReporter TestReporter,
	failingLogLevel zapcore.Level, // Chainlink core uses zapcore for logging https://docs.chain.link/chainlink-nodes/v1/configuration#log_level
) error {
	logsPath := filepath.Join(DefaultArtifactsDir, fmt.Sprintf("%s-%s-%d", t.Name(), env.Cfg.Namespace, time.Now().Unix()))
	if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
		log.Warn().Err(err).Msg("Error trying to collect pod logs")
		return err
	}
	logFiles, err := findAllLogFilesToScan(logsPath)
	if err != nil {
		log.Warn().Err(err).Msg("Error looking for pod logs")
		return err
	}
	verifyLogsGroup := &errgroup.Group{}
	for _, f := range logFiles {
		file := f
		verifyLogsGroup.Go(func() error {
			return verifyLogFile(file, failingLogLevel)
		})
	}
	assert.NoError(t, verifyLogsGroup.Wait(), "Found a concerning log")

	if t.Failed() || optionalTestReporter != nil {
		if err := SendReport(t, env, logsPath, optionalTestReporter); err != nil {
			log.Warn().Err(err).Msg("Error writing test report")
		}
	}
	return nil
}

// SendReport writes a test report and sends a Slack notification if the test provides one
func SendReport(t *testing.T, env *environment.Environment, logsPath string, optionalTestReporter TestReporter) error {
	if optionalTestReporter != nil {
		log.Info().Msg("Writing Test Report")
		optionalTestReporter.SetNamespace(env.Cfg.Namespace)
		err := optionalTestReporter.WriteReport(logsPath)
		if err != nil {
			return err
		}
		err = optionalTestReporter.SendSlackNotification(t, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// findAllLogFilesToScan walks through log files pulled from all pods, and gets all chainlink node logs
func findAllLogFilesToScan(directoryPath string) (logFilesToScan []*os.File, err error) {
	logFilePaths := []string{}
	err = filepath.Walk(directoryPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			logFilePaths = append(logFilePaths, path)
		}
		return nil
	})

	for _, filePath := range logFilePaths {
		if strings.Contains(filePath, "node.log") {
			logFileToScan, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}
			logFilesToScan = append(logFilesToScan, logFileToScan)
		}
	}
	return logFilesToScan, err
}

// allowedLogMessage is a log message that might be thrown by a Chainlink node during a test, but is not a concern
type allowedLogMessage struct {
	message string
	reason  string
	level   zapcore.Level
}

var allowedLogMessages = []allowedLogMessage{
	{
		message: "No EVM primary nodes available: 0/1 nodes are alive",
		reason:  "Sometimes geth gets unlucky in the start up process and the Chainlink node starts before geth is ready",
		level:   zapcore.DPanicLevel,
	},
}

// verifyLogFile verifies that a log file
func verifyLogFile(file *os.File, failingLogLevel zapcore.Level) error {
	// nolint
	defer file.Close()

	var (
		zapLevel zapcore.Level
		err      error
	)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		jsonLogLine := scanner.Text()
		if !strings.HasPrefix(jsonLogLine, "{") { // don't bother with non-json lines
			if strings.HasPrefix(jsonLogLine, "panic") { // unless it's a panic
				return fmt.Errorf("found panic: %s", jsonLogLine)
			}
			continue
		}
		jsonMapping := map[string]any{}

		if err = json.Unmarshal([]byte(jsonLogLine), &jsonMapping); err != nil {
			return err
		}
		logLevel, ok := jsonMapping["level"].(string)
		if !ok {
			return fmt.Errorf("found no log level in chainlink log line: %s", jsonLogLine)
		}

		if logLevel == "crit" { // "crit" is a custom core type they map to DPanic
			zapLevel = zapcore.DPanicLevel
		} else {
			zapLevel, err = zapcore.ParseLevel(logLevel)
			if err != nil {
				return fmt.Errorf("'%s' not a valid zapcore level", logLevel)
			}
		}

		if zapLevel > failingLogLevel {
			logErr := fmt.Errorf("found log at level '%s', failing any log level higher than %s: %s", logLevel, zapLevel.String(), jsonLogLine)
			logMessage, hasMessage := jsonMapping["msg"]
			if !hasMessage {
				return logErr
			}
			for _, allowedLog := range allowedLogMessages {
				if strings.Contains(logMessage.(string), allowedLog.message) {
					log.Warn().
						Str("Reason", allowedLog.reason).
						Str("Level", allowedLog.level.CapitalString()).
						Str("Msg", logMessage.(string)).
						Msg("Found allowed log message, ignoring")
				}
			}
		}
	}
	return nil
}
