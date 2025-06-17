package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func TestWorkflowExecution_DeepCopy(t *testing.T) {
	now := time.Now()
	inputs, err := values.NewMap(map[string]int{})
	require.NoError(t, err)
	step := &WorkflowExecutionStep{
		ExecutionID: "exec1",
		Ref:         "step1",
		Status:      StatusStarted,
		Inputs:      inputs,
		Outputs:     StepOutput{Value: values.NewString("output")},
		UpdatedAt:   &now,
	}

	updatedNow := now.Add(1 * time.Minute)
	finishedNow := now.Add(2 * time.Minute)
	original := WorkflowExecution{
		ExecutionID: "exec1",
		WorkflowID:  "workflow1",
		Status:      StatusStarted,
		CreatedAt:   &now,
		UpdatedAt:   &updatedNow,
		FinishedAt:  &finishedNow,
		Steps:       map[string]*WorkflowExecutionStep{"step1": step},
	}

	deepCopy := original.DeepCopy()

	// Check that the copied execution is equal to the original
	assert.Equal(t, original, deepCopy)

	// Check that the copied execution is a different instance
	assert.NotSame(t, &original, &deepCopy)

	// Check that the steps map are the same before modification of the original
	assert.Equal(t, original.Steps, deepCopy.Steps)

	// Add a new step to the original steps and check that the deep copy is not affected
	original.Steps["step2"] = &WorkflowExecutionStep{}

	// Check that the steps map is a different instance
	assert.NotEqual(t, original.Steps, deepCopy.Steps)

	// Check that the step inside the steps map is a different instance
	assert.NotSame(t, original.Steps["step1"], deepCopy.Steps["step1"])

	// Check that the inputs map is a different instance
	assert.NotSame(t, original.Steps["step1"].Inputs, deepCopy.Steps["step1"].Inputs)

	// Check that the outputs value is a different instance
	assert.NotSame(t, original.Steps["step1"].Outputs.Value, deepCopy.Steps["step1"].Outputs.Value)
}
