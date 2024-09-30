package workflows

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_KeystoneContextLabels(t *testing.T) {
	ctx := testutils.Context(t)

	expValues := KeystoneWorkflowLabels{
		WorkflowID:          fmt.Sprintf("Test%v", WorkflowID),
		WorkflowExecutionID: fmt.Sprintf("Test%v", WorkflowExecutionID),
	}
	hydratedContext1 := NewKeystoneContext(ctx, expValues)

	actValues, err := GetKeystoneLabelsFromContext(hydratedContext1)
	require.NoError(t, err)
	require.Equal(t, expValues.WorkflowID, actValues.WorkflowID)
	require.Equal(t, expValues.WorkflowExecutionID, actValues.WorkflowExecutionID)

	hydratedContext2 := NewKeystoneContext(ctx, KeystoneWorkflowLabels{})
	hydratedContext2a, err := KeystoneContextWithLabel(hydratedContext2, WorkflowID, fmt.Sprintf("Test%v", WorkflowID))
	require.NoError(t, err)
	actValues2, err := GetKeystoneLabelsFromContext(hydratedContext2a)
	require.NoError(t, err)
	require.Equal(t, expValues.WorkflowID, actValues2.WorkflowID)

}

func Test_ComposeLabeledMsg(t *testing.T) {
	unlabeledMsgFormat := "test message: %v"
	unlabeledMsgVal := fmt.Errorf("test error: %d", 0)

	testWorkflowId := "TestWorkflowID"
	testWorkflowExecutionId := "TestWorkflowExecutionID"

	ctx := testutils.Context(t)
	labeledCtx := NewKeystoneContext(ctx, KeystoneWorkflowLabels{
		WorkflowID:          testWorkflowId,
		WorkflowExecutionID: testWorkflowExecutionId,
	})

	actualLabeledMsg, err := composeLabeledMsg(labeledCtx, unlabeledMsgFormat, unlabeledMsgVal)
	require.NoError(t, err)

	expectedLabeledMsg := fmt.Sprintf("%v.%v.%v", testWorkflowId, testWorkflowExecutionId, fmt.Sprintf(unlabeledMsgFormat, unlabeledMsgVal))
	require.Equal(t, expectedLabeledMsg, actualLabeledMsg)

}
