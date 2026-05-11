package cre

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	cronWorkflowLocation = "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
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
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "cronbeholderinvalid")

	// Beholder/Kafka listener
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)
	t.Cleanup(func() {
		err := t_helpers.StopBeholder(testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
		require.NoError(t, err, "Failed to stop Beholder")
	})

	// CHiP test sink
	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(
		t,
		t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh),
	)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})
	t_helpers.IgnoreUserLogs(t.Context(), userLogsCh)

	// Deploy workflow with invalid schedule
	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: invalidSchedule,
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, cronWorkflowLocation)

	timeout := 75 * time.Second

	// Assert Beholder/Kafka path
	testLogger.Warn().Msgf("Expecting Cron workflow to fail with invalid schedule: %s", invalidSchedule)
	expectedBeholderLog := "beholder found engine initialization failure message!"
	expectedError := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, timeout)
	require.Error(t, expectedError, "Cron (Beholder) test failed. This test expects to fail with an error, but did not.")
	testLogger.Info().Msg("Beholder/Kafka assertion passed")

	// Assert CHiP sink path
	msg := t_helpers.WatchBaseMessages(
		t,
		testLogger,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		timeout,
		t_helpers.WithBaseMessageWorkflowID(workflowID),
	)
	require.Contains(t, msg.Msg, t_helpers.WorkflowEngineInitErrorLog)
	testLogger.Info().Msg("Cron invalid schedule test completed (Beholder + CHiP sink)")
}
