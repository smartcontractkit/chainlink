package wrokflowtesting

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func NewRunner(registry *Registry) workflow.Runner {
	return &runner{registry: registry}
}

type runner struct {
	registry *Registry
}

func (r *runner) Run(spec *workflow.Spec) error {
	var allSteps []workflow.StepDefinition
	allSteps = append(allSteps, spec.Actions...)
	allSteps = append(allSteps, spec.Consensus...)
	allSteps = append(allSteps, spec.Targets...)

	for trigger := r.registry.next(); trigger != nil; trigger = r.registry.next() {
		results := map[string]values.Value{}
		results["trigger"] = trigger.value
		// this can be more efficient, it's just to prove out making testing easier than running the whole engine
		if err := r.runStepsOnce(results, allSteps, spec.LocalExecutions); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) runStepsOnce(results map[string]values.Value, steps []workflow.StepDefinition, executions map[string]workflow.LocalCapability) error {
	hasMoreData := false
	for _, step := range steps {
		ran, err := r.runStepOnce(results, step, executions)
		hasMoreData = hasMoreData || ran
		if err != nil {
			return err
		}
	}

	if hasMoreData {
		return r.runStepsOnce(results, steps, executions)
	}
	return nil
}

func (r *runner) runStepOnce(results map[string]values.Value, step workflow.StepDefinition, executions map[string]workflow.LocalCapability) (bool, error) {
	if results[step.Ref] != nil {
		return false, nil
	}

	inputs := map[string]any{}
	for k, v := range step.Inputs {
		stepResult, ok := results[v]
		if !ok {
			return false, nil
		}
		inputs[k] = stepResult
	}

	wrapped, err := values.Wrap(inputs)
	if err != nil {
		return false, err
	}

	if action, ok := r.registry.remoteActionsAndConsensus[step.Ref]; ok {
		return true, runMock(results, step, action, wrapped)
	} else if target, ok := r.registry.remoteTargets[step.Ref]; ok {
		wrapped = wrapped.(*values.Map).Underlying["report"]
		results[step.Ref] = wrapped
		return false, target(wrapped)
	}

	return execute(results, step, executions, err, wrapped)
}

func runMock(results map[string]values.Value, step workflow.StepDefinition, action func(value values.Value) (values.Value, error), wrapped values.Value) error {
	result, err := action(wrapped)
	if err != nil {
		return err
	}
	results[step.Ref] = result
	return nil
}

func execute(results map[string]values.Value, step workflow.StepDefinition, executions map[string]workflow.LocalCapability, err error, wrapped values.Value) (bool, error) {
	result, cont, err := executions[step.Ref].Run(step.Ref, wrapped)
	if err != nil {
		return false, err
	}
	if !cont {
		return false, nil
	}

	results[step.Ref] = result
	return true, nil
}
