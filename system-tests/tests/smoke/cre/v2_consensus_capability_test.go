package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteConsensusTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	// Start subscriber BEFORE registering workflow to capture the current sequence
	// and avoid processing old events from previous test runs
	channel, listenerCtx, cancelFn, startErr := t_helpers.StartChipEventsSubscriber(t.Context(), t_helpers.ChipEventsSubscriberConfig{
		ChipServerURL: "http://localhost:8081", // HTTP API port, not gRPC port
		Logger:        testLogger,
	})
	require.NoError(t, startErr, "Failed to start workflow events subscriber")
	t.Cleanup(cancelFn)

	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")

	channels := t_helpers.FanOutChipEvents(listenerCtx, channel, 2)
	go func() {
		t_helpers.LogChipEvents(listenerCtx, channels[0], testLogger)
	}()

	expectedBeholderLog := "Successfully passed all consensus tests"
	// err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, 4*time.Minute)
	err := t_helpers.AssertChipEventMatchedByNodes(listenerCtx, workflowID, 2, channels[1], t_helpers.GetUserLogMatcherFn(expectedBeholderLog), 4*time.Minute, testLogger)
	require.NoError(t, err, "Consensus capability test failed, Beholder should not return an error")
	testLogger.Info().Msg("Consensus capability test completed")
}
