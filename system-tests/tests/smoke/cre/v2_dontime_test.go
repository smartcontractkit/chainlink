package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteDonTimeTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/time_consensus/main.go"
	workflowName := "timebeholder"

	// Start subscriber BEFORE registering workflow to capture the current sequence
	// and avoid processing old events from previous test runs
	channel, listenerCtx, cancelFn, startErr := t_helpers.StartChipEventsSubscriber(t.Context(), t_helpers.ChipEventsSubscriberConfig{
		ChipServerURL: "http://localhost:8081",
		Logger:        testLogger,
	})
	require.NoError(t, startErr, "Failed to start chip events subscriber")
	t.Cleanup(cancelFn)

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	channels := t_helpers.FanOutChipEvents(listenerCtx, channel, 2)
	go func() {
		t_helpers.LogChipEvents(listenerCtx, channels[0], testLogger)
	}()

	expectedBeholderLog := "Verified consensus on DON Time"
	err := t_helpers.AssertChipEventMatchedByNodes(listenerCtx, workflowID, 2, channels[1], t_helpers.GetUserLogMatcherFn(expectedBeholderLog), timeout, testLogger)
	require.NoError(t, err, "DON Time test failed, Beholder should not return an error")
	testLogger.Info().Msg("DON Time test completed")
}
