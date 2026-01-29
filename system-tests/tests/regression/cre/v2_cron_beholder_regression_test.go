package cre

import (
	"testing"
	"time"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

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
		Schedule: invalidSchedule,
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	testLogger.Warn().Msgf("Expecting Cron workflow to fail with invalid schedule: %s", invalidSchedule)
	expectedBeholderLog := "beholder found engine initialization failure message!"

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorklfowEngineInitErrorLog, expectedBeholderLog, 75*time.Second)
	testLogger.Info().Msg("Cron (Beholder) fail test completed")
}
