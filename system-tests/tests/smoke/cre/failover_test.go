package cre

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	stvault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteFailoverManualSwapTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	shardDONs := testEnv.Dons.DonsWithFlag(cre.ShardDON)
	require.GreaterOrEqual(t, len(shardDONs), 2, "Expected at least 2 shard DONs for failover test")

	shardLeaderDON := getShardZeroDon(t, testEnv)

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}
	expectedUserLog := "Amazing workflow user log"

	defaultOwner := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"

	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err, "failed to start linking service")
	linkingService.SetOwnerOrg(defaultOwner, "org_test_failover")

	primaryAssignmentTOML := `
static_default_assignment = [0,1]
hashed_default_assignment = false

[per_org_assignment]
  org_test_failover = [0,1]
`

	proposeAndApproveShardAssignmentJob(t, testEnv, shardLeaderDON, primaryAssignmentTOML, testLogger)

	const numWorkflows = 3
	workflowIDs := make([]string, 0, numWorkflows)
	for i := range numWorkflows {
		workflowName := fmt.Sprintf("failover-run%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs = append(workflowIDs, workflowID)
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed workflows for failover test")

	workflowToShardIndex := make(map[string]uint32, len(workflowIDs))
	for _, wfID := range workflowIDs {
		workflowToShardIndex[wfID] = 0
	}

	nodeP2PIDToShardIndex := buildNodeP2PIDToShardIndex(t, testEnv)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	execTimeout := 3 * time.Minute
	timeoutCtx, cancelTimeout := context.WithTimeout(t.Context(), execTimeout)
	defer cancelTimeout()
	execCtx, cancelCause := context.WithCancelCause(timeoutCtx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(execCtx, cancelCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	testLogger.Info().Msg("Phase 1: Verify workflows execute on primary shard (shard 0)")
	executedWorkflows := waitForAllWorkflowsExecuted(execCtx, t, testLogger, userLogsCh, workflowIDs, workflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, execTimeout)
	require.Len(t, executedWorkflows, len(workflowIDs), "Not all workflows executed on primary shard")
	testLogger.Info().Int("executedCount", len(executedWorkflows)).Msg("All workflows executed on primary shard (shard 0)")

	for i := range numWorkflows {
		workflowName := fmt.Sprintf("failover_swap%d", i)
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)
		workflowIDs[i] = workflowID
		workflowToShardIndex[workflowID] = 1
	}
	testLogger.Info().Strs("workflowIDs", workflowIDs).Msg("Deployed fresh workflows for swap phase")

	for _, don := range shardDONs {
		proposeAndApproveShardAssignmentJob(t, testEnv, don, secondaryAssignmentTOML, testLogger)
	}

	testLogger.Info().Msg("Phase 2: After swap [0,1]->[1,0], verify workflows execute on new primary (shard 1)")

	swappedWorkflowIDs := workflowIDs
	swappedWorkflowToShardIndex := make(map[string]uint32, len(swappedWorkflowIDs))
	for _, wfID := range swappedWorkflowIDs {
		swappedWorkflowToShardIndex[wfID] = 1
	}

	swapTimeout := 5 * time.Minute
	swapTimeoutCtx, cancelSwapTimeout := context.WithTimeout(t.Context(), swapTimeout)
	defer cancelSwapTimeout()
	swapExecCtx, cancelSwapCause := context.WithCancelCause(swapTimeoutCtx)
	defer cancelSwapCause(nil)
	go t_helpers.FailOnBaseMessage(swapExecCtx, cancelSwapCause, t, testLogger, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)

	swapUserLogsCh := make(chan *workflowevents.UserLogs, 1000)
	swapBaseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	swapServer := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, swapUserLogsCh, swapBaseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, swapServer, swapUserLogsCh, swapBaseMessageCh)
	})

	swappedExecuted := waitForAllWorkflowsExecuted(swapExecCtx, t, testLogger, swapUserLogsCh, swappedWorkflowIDs, swappedWorkflowToShardIndex, nodeP2PIDToShardIndex, expectedUserLog, swapTimeout)
	require.Len(t, swappedExecuted, len(swappedWorkflowIDs), "Not all workflows executed on new primary shard after swap")
	testLogger.Info().Int("executedCount", len(swappedExecuted)).Msg("All workflows executed on new primary shard (shard 1) after manual failover swap")
}

const secondaryAssignmentTOML = `
static_default_assignment = [1,0]
hashed_default_assignment = false

[per_org_assignment]
  org_test_failover = [1,0]
`
