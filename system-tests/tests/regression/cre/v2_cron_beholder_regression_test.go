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

// regression
var cronInvalidSchedulesTests = []struct {
	name            string
	invalidSchedule string
}{
	{"negative", "*/-1 * * * * *"},
	{"below default limit", "*/29 * * * * *"},
	{"inappropriately formatted", "*MON/1 * * * * *"},
}

func CronBeholderFailsWithInvalidScheduleTest(t *testing.T, testEnv *ttypes.TestEnvironment, invalidSchedule string) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: invalidSchedule,
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	channel, listenerCtx, cancelFn, startErr := t_helpers.StartWorkflowEventsSubscriber(t.Context(), t_helpers.GetStandardWorkflowEventsSubscriberConfig(testEnv, workflowID))
	require.NoError(t, startErr, "Failed to start workflow events subscriber")

	testLogger.Warn().Msgf("Expecting Cron workflow to fail with invalid schedule: %s", invalidSchedule)
	expectedBeholderLog := "beholder found engine initialization failure message!"
	timeout := 75 * time.Second
	expectedError := t_helpers.AssertWorkflowEventMatched(listenerCtx, cancelFn, 2, channel, t_helpers.GetUserLogMatcherFn(expectedBeholderLog), timeout, testLogger)
	require.Error(t, expectedError, "Cron (Beholder) test failed. This test expects to fail with an error, but did not.")

	testLogger.Info().Msg("Cron (Beholder) fail test completed")
}
