package v2

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutionProfileCollector(t *testing.T) {
	t.Parallel()

	collector := newExecutionProfileCollector()
	start := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	stepEnd := start.Add(50 * time.Millisecond)
	collector.recordStepStart("1", "cap@1.0.0", start)
	collector.recordStepEnd("1", stepEnd, false)

	steps := collector.stepInputs()
	require.Len(t, steps, 1)
	assert.Equal(t, "1", steps[0].StepID)
	assert.Equal(t, "cap@1.0.0", steps[0].CapabilityID)
	assert.Equal(t, start, steps[0].StartTime)
	assert.Equal(t, stepEnd, steps[0].EndTime)
	assert.False(t, steps[0].HasError)
}

func TestExecutionProfileCollectorInsertionOrder(t *testing.T) {
	t.Parallel()

	collector := newExecutionProfileCollector()
	start := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	collector.recordStepStart("1", "cap@1.0.0", start.Add(10*time.Millisecond))
	collector.recordStepStart("2", "cap@2.0.0", start.Add(20*time.Millisecond))
	collector.recordStepEnd("1", start.Add(30*time.Millisecond), false)
	collector.recordStepEnd("2", start.Add(40*time.Millisecond), true)

	steps := collector.stepInputs()
	require.Len(t, steps, 2)
	assert.Equal(t, "1", steps[0].StepID)
	assert.Equal(t, "2", steps[1].StepID)
	assert.False(t, steps[0].HasError)
	assert.True(t, steps[1].HasError)
}

func TestExecutionProfileCollectorMaxSteps(t *testing.T) {
	t.Parallel()

	collector := newExecutionProfileCollector()
	start := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	for i := range maxExecutionProfileSteps + 5 {
		stepID := strconv.Itoa(i)
		collector.recordStepStart(stepID, "cap@1.0.0", start.Add(time.Duration(i)*time.Millisecond))
		collector.recordStepEnd(stepID, start.Add(time.Duration(i+1)*time.Millisecond), false)
	}

	steps := collector.stepInputs()
	require.Len(t, steps, maxExecutionProfileSteps)
	assert.Equal(t, "0", steps[0].StepID)
	assert.Equal(t, strconv.Itoa(maxExecutionProfileSteps-1), steps[len(steps)-1].StepID)
}
