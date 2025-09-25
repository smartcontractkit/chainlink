package cre

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/proto"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteConsensusTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	beholder, err := t_helpers.NewBeholder(testLogger, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, err, "failed to create beholder instance")

	ctxWithTimeout, cancelCtx := context.WithTimeout(t.Context(), 4*time.Minute)
	defer cancelCtx()

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	beholderMessageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflowevents.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &commonevents.BaseMessage{}
		},
	}

	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")
	beholderMsgChan, beholderErrChan := beholder.SubscribeToBeholderMessages(ctxWithTimeout, beholderMessageTypes)

	expectedUserLog := "Successfully fetched"
	var receivedResults []string
	nodeCount := 4

	// Check the beholder logs for the expected messages
	for {
		select {
		case <-ctxWithTimeout.Done():
			require.Fail(t, "Test timed out before completion")
		case err := <-beholderErrChan:
			require.FailNowf(t, "Kafka error received from Kafka %s", err.Error())
		case msg := <-beholderMsgChan:
			switch typedMsg := msg.(type) {
			case *commonevents.BaseMessage:
				// Log this as it can be useful for debugging
				testLogger.Debug().Msgf("Received BaseMessage from Beholder: %s", typedMsg.Msg)
			case *workflowevents.UserLogs:
				testLogger.Info().Msg("🎉 Received UserLogs message in test")

				for _, logLine := range typedMsg.LogLines {
					if strings.Contains(logLine.Message, "Consensus error") {
						testLogger.Warn().
							Str("message", strings.TrimSpace(logLine.Message)).
							Msgf("⚠️ Received consensus error from workflow: %s", logLine.Message)
					}

					if strings.Contains(logLine.Message, expectedUserLog) {
						testLogger.Info().
							Str("expected_log", expectedUserLog).
							Str("found_message", strings.TrimSpace(logLine.Message)).
							Msg("🎯 Found expected user log message!")

						// Extract the result value from the message
						pattern := `result=([0-9]+)`
						re := regexp.MustCompile(pattern)

						matches := re.FindStringSubmatch(logLine.Message)
						if len(matches) > 1 {
							receivedResults = append(receivedResults, matches[1])
						}

						// Each node will log the consensus result value so we expect at least node count instances of the value
						// in the beholder logs.
						valueCounts := make(map[string]int)
						for _, result := range receivedResults {
							valueCounts[result]++
							if valueCounts[result] >= nodeCount {
								testLogger.Info().Msgf("Found %d identical results for value %s, test has passed", nodeCount, result)
								return
							}
						}
					} else {
						testLogger.Info().Msgf("Received user message from Beholder: %s", typedMsg.LogLines)
					}
				}
			default:
				// ignore other message types
			}
		}
	}
}
