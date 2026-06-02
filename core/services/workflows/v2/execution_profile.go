package v2

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

// To avoid excessive memory usage, we limit the number of steps in the execution profile.
const maxExecutionProfileSteps = 100

type executionProfileStep struct {
	stepID       string
	capabilityID string
	startTime    time.Time
	endTime      time.Time
	hasError     bool
}

type executionProfileCollector struct {
	mu    sync.Mutex
	steps map[string]executionProfileStep
	order []string
}

func newExecutionProfileCollector() *executionProfileCollector {
	return &executionProfileCollector{
		steps: make(map[string]executionProfileStep),
	}
}

func (c *executionProfileCollector) recordStepStart(stepID, capabilityID string, start time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.steps[stepID]; !ok {
		if len(c.order) >= maxExecutionProfileSteps {
			return
		}
		c.order = append(c.order, stepID)
	}
	c.steps[stepID] = executionProfileStep{
		stepID:       stepID,
		capabilityID: capabilityID,
		startTime:    start,
	}
}

func (c *executionProfileCollector) recordStepEnd(stepID string, end time.Time, hasError bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	step, ok := c.steps[stepID]
	if !ok {
		return
	}
	step.endTime = end
	step.hasError = hasError
	c.steps[stepID] = step
}

func (c *executionProfileCollector) stepInputs() []events.ExecutionProfileStepInput {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]events.ExecutionProfileStepInput, 0, len(c.order))
	for _, stepID := range c.order {
		if step, ok := c.steps[stepID]; ok {
			out = append(out, events.ExecutionProfileStepInput{
				StepID:       step.stepID,
				StartTime:    step.startTime,
				EndTime:      step.endTime,
				CapabilityID: step.capabilityID,
				HasError:     step.hasError,
			})
		}
	}
	return out
}
