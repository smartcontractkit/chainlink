package cre

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteConsensusTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")

	expectedBeholderLog := "Successfully passed all consensus tests"
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedBeholderLog, 4*time.Minute)
	testLogger.Info().Msg("Consensus capability test completed")
}
