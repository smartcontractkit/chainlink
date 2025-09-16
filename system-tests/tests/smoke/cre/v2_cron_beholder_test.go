package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
)

// smoke
func ExecuteCronBeholderTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	expectedBeholderLog := "Amazing workflow user log"
	timeout := 2 * time.Minute
	err := assertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "Cron (Beholder) test failed")
	testLogger.Info().Msg("Cron (Beholder) test completed")
}

// regression
var cronInvalidSchedulesTests = []struct {
	name            string
	invalidSchedule string
}{
	{"below default limit (30s)", "*/29 * * * * *"},
	{"negative", "*/-1 * * * * *"},
	{"inappropriately formatted", "*MON/1 * * * * *"},
}

func CronBeholderFailWithInvalidScheduleTest(t *testing.T, testEnv *TestEnvironment, invalidSchedule string) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	listenerCtx, messageChan, kafkaErrChan := startBeholder(t, testLogger, testEnv)

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: invalidSchedule,
	}
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	testLogger.Warn().Msgf("Expecting Cron workflow to fail with invalid schedule: %s", invalidSchedule)
	expectedBeholderLog := "expecting the error from a Beholder messages test validator" // no empty string, because it may match a message
	timeout := 1 * time.Minute
	expectedError := assertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.Error(t, expectedError, "Cron (Beholder) test failed. This test expects to fail with an error, but did not.")

	testLogger.Info().Msg("Cron (Beholder) fail test completed")
}
