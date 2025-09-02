package cre

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/proto"

	common_events "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflow_events "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteBeholderTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	testLogger.Info().Msg("Starting Beholder...")
	beholderConfigPath := testEnv.TestConfig.BeholderConfigPath
	bErr := startBeholderStackIfIsNotRunning(beholderConfigPath, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, bErr, "failed to start Beholder")

	chipConfig, chipErr := loadBeholderStackCache(beholderConfigPath)
	require.NoError(t, chipErr, "failed to load chip ingress cache")
	require.NotNil(t, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, "kafka external url is not set in the cache")
	require.NotEmpty(t, chipConfig.Kafka.Topics, "kafka topics are not set in the cache")

	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &None{}, workflowFileLocation)

	listenerCtx, cancelListener := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(func() {
		cancelListener()
	})

	kafkaErrChan := make(chan error, 1)
	messageChan := make(chan proto.Message, 10)

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	messageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflow_events.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &common_events.BaseMessage{}
		},
	}

	// Start listening for messages in the background
	go func() {
		listenForKafkaMessages(listenerCtx, testLogger, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, chipConfig.Kafka.Topics[0], messageTypes, messageChan, kafkaErrChan)
	}()

	expectedUserLog := "Amazing workflow user log"

	foundExpectedLog := make(chan bool, 1) // Channel to signal when expected log is found
	foundErrorLog := make(chan bool, 1)    // Channel to signal when engine initialization failure is detected
	receivedUserLogs := 0
	// Start message processor goroutine
	go func() {
		for {
			select {
			case <-listenerCtx.Done():
				return
			case msg := <-messageChan:
				// Process received messages
				switch typedMsg := msg.(type) {
				case *common_events.BaseMessage:
					if strings.Contains(typedMsg.Msg, "Workflow Engine initialization failed") {
						foundErrorLog <- true
					}
				case *workflow_events.UserLogs:
					testLogger.Info().
						Msg("🎉 Received UserLogs message in test")
					receivedUserLogs++

					for _, logLine := range typedMsg.LogLines {
						if strings.Contains(logLine.Message, expectedUserLog) {
							testLogger.Info().
								Str("expected_log", expectedUserLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
								Msg("🎯 Found expected user log message!")

							select {
							case foundExpectedLog <- true:
							default: // Channel might already have a value
							}
							return // Exit the processor goroutine
						}
						testLogger.Warn().
							Str("expected_log", expectedUserLog).
							Str("found_message", strings.TrimSpace(logLine.Message)).
							Msg("Received UserLogs message, but it does not match expected log")
					}
				default:
					// ignore other message types
				}
			}
		}
	}()

	testLogger.Info().
		Str("expected_log", expectedUserLog).
		Dur("timeout", timeout).
		Msg("Waiting for expected user log message or timeout")

	// Wait for either the expected log to be found, or engine initialization failure to be detected, or timeout (2 minutes)
	select {
	case <-foundExpectedLog:
		testLogger.Info().
			Str("expected_log", expectedUserLog).
			Msg("✅ Test completed successfully - found expected user log message!")
		return
	case <-foundErrorLog:
		require.Fail(t, "Test completed with error - found engine initialization failure message!")
	case <-time.After(timeout):
		testLogger.Error().Msg("Timed out waiting for expected user log message")
		if receivedUserLogs > 0 {
			testLogger.Warn().Int("received_user_logs", receivedUserLogs).Msg("Received some UserLogs messages, but none matched expected log")
		} else {
			testLogger.Warn().Msg("Did not receive any UserLogs messages")
		}
		require.Failf(t, "Timed out waiting for expected user log message", "Expected user log message: '%s' not found after %s", expectedUserLog, timeout.String())
	case err := <-kafkaErrChan:
		testLogger.Error().Err(err).Msg("Kafka listener encountered an error during execution. Ensure Beholder is running and accessible.")
		require.Fail(t, "Kafka listener failed", err.Error())
	}

	testLogger.Info().Msg("Beholder test completed")
}
