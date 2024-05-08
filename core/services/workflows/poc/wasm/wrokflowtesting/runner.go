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

// TODO do we need the .output and deepmap stuff here? I don't think so for code...
func (r *runner) Run(spec *workflow.Spec) error {
	for trigger := r.registry.next(); trigger != nil; trigger = r.registry.next() {
		results := map[string]values.Value{}
		results["trigger"] = trigger.value
		if err := runNext("trigger", results, spec); err != nil {
			return err
		}
	}
	return nil
}

func runNext(from string, results map[string]values.Value, spec *workflow.Spec) error {
	// could be more efficient with maps, but this is just for testing in a POC and I'm lazy
	for _, action := range spec.Actions {
		if _, ok := results[action.Ref]; ok {
			continue
		}

		for _, input := range action.StepDependencies {

		}
	}
}
