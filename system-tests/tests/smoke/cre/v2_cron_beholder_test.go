package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// smoke
func ExecuteCronBeholderTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	expectedBeholderLog := "Amazing workflow user log"

	err := t_helpers.AssertBeholderMessage(listenerCtx, t, expectedBeholderLog, testLogger, messageChan, kafkaErrChan, 4*time.Minute)
	require.NoError(t, err, "Cron (Beholder) test failed")
	testLogger.Info().Msg("Cron (Beholder) test completed")

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorklfowEngineInitErrorLog, expectedBeholderLog, 2*time.Minute)
}
