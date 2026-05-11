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
	cronWorkflowLocation    = "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	cronExpectedBeholderLog = "Amazing workflow user log"
)

// smoke
func ExecuteCronBeholderTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "cronbeholder")

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	t.Cleanup(func() {
		// stop ChIP Ingress after the test to free the port, on which other tests will start the ChiP Test Sink
		err := t_helpers.StopBeholder(testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
		require.NoError(t, err, "Failed to stop Beholder")
	})

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, cronWorkflowLocation)

	err := t_helpers.AssertBeholderMessage(listenerCtx, t, cronExpectedBeholderLog, testLogger, messageChan, kafkaErrChan, 4*time.Minute)
	require.NoError(t, err, "Cron (Beholder) test failed")
	testLogger.Info().Msg("Cron (Beholder) test completed")
}

func ExecuteChipIngressBatchingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "chipingressbatch")

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

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, cronWorkflowLocation)

	t_helpers.WatchWorkflowLogs(
		t,
		testLogger,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		cronExpectedBeholderLog,
		4*time.Minute,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)

	testLogger.Info().Msg("ChipIngress batching test completed")
}
