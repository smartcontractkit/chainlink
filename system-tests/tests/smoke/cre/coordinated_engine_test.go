package cre

import (
	"context"
	"testing"
	"time"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// coordinatedEngineDON is the workflow DON of the default topology
// (configs/workflow-gateway-capabilities-don.toml).
const coordinatedEngineDON = "workflow"

// Test_CRE_V2_CoordinatedEngine runs a cron workflow with CoordinatedEngineEnabled
// on, so its triggers are registered, delivered and ACKed by the node's trigger
// coordinator instead of by the engine itself.
//
// The flag is read once at engine creation, so it is applied before deploying.
// It is a global override on the shared environment, hence serial.
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_V2_CoordinatedEngine$' -timeout 20m -v
//
//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_V2_CoordinatedEngine(t *testing.T) {
	testLogger := framework.L
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`CoordinatedEngineEnabled = 'true'`))

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// t.Context() is already cancelled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowName := t_helpers.UniqueWorkflowName(testEnv, "coordinatedcron")
	workflowConfig := crontypes.WorkflowConfig{Schedule: "*/30 * * * * *"}
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 2*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))

	// The execution alone would also pass on the legacy engine; these prove it ran on the coordinated path.
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, coordinatedEngineDON,
		"Routing workflow to the coordinated engine", time.Minute, 5*time.Second)
	t_helpers.RequireContainerLogsForNodesetEventually(t, testEnv, coordinatedEngineDON,
		"Registered triggers via coordinator", time.Minute, 5*time.Second)
}
