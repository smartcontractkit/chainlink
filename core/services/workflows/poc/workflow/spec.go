package workflow

type Spec struct {
	SpecBase
	LocalExecutions map[string]LocalCapability
}

type SpecBase struct {
	Triggers  []StepDefinition `json:"triggers" jsonschema:"required"`
	Actions   []StepDefinition `json:"actions,omitempty"`
	Consensus []StepDefinition `json:"consensus" jsonschema:"required"`
	Targets   []StepDefinition `json:"targets" jsonschema:"required"`
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
	TypeRef string
	Ref     string            `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_]+$"`
	Inputs  map[string]string `json:"inputs,omitempty"`
	// Ideally, values.Value should be able to serialize anything that's kind is int, but hack around that for now.
	CapabilityType int `json:"-"`
}

type Workflow struct {
	refs map[string]*StepDefinition
}
