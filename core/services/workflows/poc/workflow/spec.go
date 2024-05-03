package workflow

import (
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type Spec struct {
	Triggers        []StepDefinition `json:"triggers" jsonschema:"required"`
	Actions         []StepDefinition `json:"actions,omitempty"`
	Consensus       []StepDefinition `json:"consensus" jsonschema:"required"`
	Targets         []StepDefinition `json:"targets" jsonschema:"required"`
	LocalExecutions map[string]LocalCapability
}

func (w *Spec) ToWorkflow() Workflow {
	refs := map[string]*StepDefinition{}
	addAllRefs(refs, w.Triggers)
	addAllRefs(refs, w.Actions)
	addAllRefs(refs, w.Consensus)
	return Workflow{refs: refs}
}

func addAllRefs(refs map[string]*StepDefinition, steps []StepDefinition) {
	for _, s := range steps {
		refs[s.Ref] = &s
	}
}

type StepDefinition struct {
	TypeRef        string
	Ref            string                      `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_]+$"`
	Inputs         map[string]any              `json:"inputs,omitempty"`
	CapabilityType capabilities.CapabilityType `json:"-"`
}

type Workflow struct {
	refs map[string]*StepDefinition
}
