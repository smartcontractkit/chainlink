package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

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

func CronChipIngressStackFailsWithInvalidScheduleTest(t *testing.T, testEnv *ttypes.TestEnvironment, invalidSchedule string) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "cronchipingressstack")

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartChipIngressStack(t, testLogger, testEnv)
	t.Cleanup(func() {
		// stop ChIP Ingress after the test to free the port, on which other tests will start the ChiP Test Sink
		err := t_helpers.StopChipIngressStack(testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
		require.NoError(t, err, "Failed to stop Chip Ingress stack")
	})

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: invalidSchedule,
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	testLogger.Warn().Msgf("Expecting Cron workflow to fail with invalid schedule: %s", invalidSchedule)
	// Not matched via UserLogs; engine init failure path ends the assertion with an error.
	unusedExpectedUserLog := "__unused_expected_user_log_for_negative_test__"
	timeout := 75 * time.Second
	expectedError := t_helpers.AssertChipIngressStackMessage(listenerCtx, t, unusedExpectedUserLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.Error(t, expectedError, "Cron (Chip Ingress stack) test failed. This test expects to fail with an error, but did not.")
	testLogger.Info().Msg("Cron (Chip Ingress stack) fail test completed")
}
