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

// smoke
func ExecuteCronBeholderTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	channel, listenerCtx, cancelFn, startErr := t_helpers.StartWorkflowEventsSubscriber(t.Context(), t_helpers.GetStandardWorkflowEventsSubscriberConfig(testEnv, workflowID))
	require.NoError(t, startErr, "Failed to start workflow events subscriber")
	t.Cleanup(cancelFn)

	channels := t_helpers.FanOutWorkflowEvents(listenerCtx, channel, 2)
	go func() {
		t_helpers.LogWorkflowEvent(listenerCtx, channels[0])
	}()

	err := t_helpers.AssertWorkflowEventMatched(listenerCtx, 2, channels[1], t_helpers.GetUserLogMatcherFn("Amazing workflow user log"), 2*time.Minute, testLogger)
	require.NoError(t, err, "Cron (Beholder) test failed")
	testLogger.Info().Msg("Cron (Beholder) test completed")
}
