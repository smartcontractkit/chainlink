package cre

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workfloweventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

//nolint:paralleltest // parallel execution is controlled by CRE_TEST_PARALLEL_ENABLED
func Test_CRE_V2_WorkflowExecutionMetadata(t *testing.T) {
	// TODO: (CRE-5839) Remove this skip once baseline images supply execution metadata.
	t_helpers.SkipIfMixedEnv(t, "baseline nodes do not generate execution metadata, hence failing consensus; enable in version 2.75.0+")
	if parallelEnabled {
		t.Parallel()
	}
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))
	ExecuteWorkflowExecutionMetadataTest(t, testEnv)
}

func ExecuteWorkflowExecutionMetadataTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	// Start telemetry capture.
	userLogs := make(chan *workfloweventsv2.WorkflowUserLog, 1000)
	baseMessages := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(framework.L, userLogs, baseMessages))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogs, baseMessages)
	})

	// Deploy the workflow.
	workflowID := t_helpers.CompileAndDeployWorkflow[t_helpers.None](t, testEnv, framework.L,
		t_helpers.UniqueWorkflowName(testEnv, "execution-metadata"), nil, "./executiometadata/main.go",
		t_helpers.WithArtifactCopyDONTypes(crelib.WorkflowDON, crelib.ShardDON))
	workflowID = strings.ToLower(strings.TrimPrefix(workflowID, "0x"))
	require.NotNil(t, testEnv.Execution)
	owner := strings.ToLower(strings.TrimPrefix(testEnv.Execution.OwnerAddress.Hex(), "0x"))
	workflowDONs := make(map[string]*crelib.Don)
	for _, don := range testEnv.Dons.List() {
		if don.HasFlag(crelib.WorkflowDON) || don.HasFlag(crelib.ShardDON) {
			workflowDONs[strconv.FormatUint(don.ID, 10)] = don
		}
	}

	// Wait for successful metadata validation.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()
	watchCtx, cancelCause := context.WithCancelCause(ctx)
	defer cancelCause(nil)
	go t_helpers.FailOnBaseMessage(watchCtx, cancelCause, t, framework.L, baseMessages,
		t_helpers.WorkflowEngineInitErrorLog, t_helpers.WithBaseMessageWorkflowID(workflowID))
	observed, err := t_helpers.WaitForUserLog(watchCtx, framework.L, userLogs,
		"event=WORKFLOW_EXECUTION_METADATA", t_helpers.WithUserLogWorkflowID(workflowID))
	require.NoError(t, err)
	require.NotNil(t, observed)
	require.Contains(t, observed.GetMsg(), "Workflow execution metadata fetched successfully")

	// Compare live metadata with deployment and host values.
	require.NotEmpty(t, observed.GetWorkflowExecutionID())
	donID := strconv.FormatInt(int64(observed.GetCreInfo().GetDonID()), 10)
	require.Contains(t, workflowDONs, donID)
	don := workflowDONs[donID]
	fields := strings.Fields(observed.GetMsg())
	for _, expected := range []string{
		"event=WORKFLOW_EXECUTION_METADATA",
		"workflowID=" + workflowID,
		"executionID=" + observed.GetWorkflowExecutionID(),
		"workflowOwner=" + owner,
		"donID=" + donID,
		fmt.Sprintf("donF=%d", don.F),
	} {
		require.Contains(t, fields, expected)
	}
}
