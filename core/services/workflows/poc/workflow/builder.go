package workflow

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

func NewWorkflowBuilder[O any](
	trigger capabilities.Trigger[O], additionalTriggers ...capabilities.Trigger[O]) (*Root, *Builder[O]) {
	tr := &triggerRunner[O]{Trigger: trigger}
	runner := &multiTriggerRunner[O]{
		triggers: map[string]capability{tr.Ref(): tr},
	}

	root := &Root{
		capabilities: []capability{runner},
		names:        map[string]bool{runner.Ref(): true},
		open:         map[string]bool{},
		spec: &Spec{
			Triggers:        []StepDefinition{capabilityToStepDef(tr)},
			LocalExecutions: map[string]LocalCapability{},
		},
	}
	root.names[tr.Ref()] = true

	for _, t := range additionalTriggers {
		tr = &triggerRunner[O]{Trigger: t}
		root.names[t.Ref()] = true
		runner.triggers[t.Ref()] = tr
		root.spec.Triggers = append(root.spec.Triggers, capabilityToStepDef(tr))
		root.capabilities = append(root.capabilities, tr)
	}

	root.open[runner.Ref()] = true
	return root, &Builder[O]{root: root, current: runner}
}

func capabilityToStepDef(tr capability) StepDefinition {
	return StepDefinition{
		TypeRef:        tr.Type(),
		Ref:            tr.Ref(),
		Inputs:         tr.Inputs(),
		CapabilityType: tr.capabilityType(),
	}
}

// Root is NOT thread safe, do not call methods with it concurrently
type Root struct {
	capabilities []capability
	built        bool
	lock         sync.Mutex
	names        map[string]bool
	open         map[string]bool
	spec         *Spec
}

func (wr *Root) Build() (*Spec, error) {
	wr.lock.Lock()
	defer wr.lock.Unlock()
	if wr.built {
		return nil, errors.New("cannot build workflow twice")
	}

	var open []string
	for key, c := range wr.open {
		if c {
			open = append(open, key)
		}
	}

	if len(open) > 0 {
		return nil, fmt.Errorf("workflow has unused steps: %v", open)
	}

	wr.built = true
	return wr.spec, nil
}

type Builder[O any] struct {
	root    *Root
	current capability
}

func AddStep[I, O any](wb *Builder[I], a capabilities.Action[I, O]) (*Builder[O], error) {
	wb.root.lock.Lock()
	defer wb.root.lock.Unlock()
	if wb.root.built {
		return nil, errors.New("cannot add steps after workflow has been built")
	}

	if wb.root.names[a.Ref()] {
		return nil, fmt.Errorf("name %s already exists as a step", a.Ref())
	}
	wb.root.names[a.Ref()] = true
	wb.root.open[a.Ref()] = true
	wb.root.open[wb.current.Ref()] = false

	ar := &actionRunner[I, O]{
		nonTriggerCapability: nonTriggerCapability{
			inputs: map[string]any{"action": wb.current.Outputs()},
			ref:    a.Ref(),
		},
		Action: a,
	}

	wb.root.spec.Actions = append(wb.root.spec.Actions, capabilityToStepDef(ar))
	wb.root.capabilities = append(wb.root.capabilities, ar)

	wb.root.spec.LocalExecutions[a.Ref()] = ar
	return &Builder[O]{
		root:    wb.root,
		current: ar,
	}, nil
}

func AddConsensus[I, O any](wb *Builder[I], c capabilities.Consensus[I, O]) (*Builder[capabilities.ConsensusResult[O]], error) {
	wb.root.lock.Lock()
	defer wb.root.lock.Unlock()
	if wb.root.built {
		return nil, errors.New("cannot add steps after workflow has been built")
	}

	if wb.root.names[c.Ref()] {
		return nil, fmt.Errorf("name %s already exists as a step", c.Ref())
	}
	wb.root.names[c.Ref()] = true
	wb.root.open[c.Ref()] = true
	wb.root.open[wb.current.Ref()] = false

	cr := &consensusRunner[I, O]{
		nonTriggerCapability: nonTriggerCapability{
			inputs: map[string]any{"report": wb.current.Outputs()},
			ref:    c.Ref(),
		},
		Consensus: c,
	}

	wb.root.spec.Consensus = append(wb.root.spec.Consensus, capabilityToStepDef(cr))
	wb.root.capabilities = append(wb.root.capabilities, cr)

	wb.root.spec.LocalExecutions[c.Ref()] = cr

	return &Builder[capabilities.ConsensusResult[O]]{
		root:    wb.root,
		current: cr,
	}, nil
}

func AddTarget[O any](wb *Builder[capabilities.ConsensusResult[O]], t capabilities.Target[O]) error {
	wb.root.lock.Lock()
	defer wb.root.lock.Unlock()
	if wb.root.built {
		return errors.New("cannot add steps after workflow has been built")
	}

	wb.root.open[wb.current.Ref()] = false
	if wb.root.names[t.Ref()] {
		return fmt.Errorf("name %s already exists as a step", t.Ref())
	}
	wb.root.names[t.Ref()] = true

	tr := &targetRunner[O]{
		inputs: map[string]any{"report": wb.current.Outputs()},
		Target: t,
	}
	wb.root.spec.Targets = append(wb.root.spec.Targets, capabilityToStepDef(tr))
	wb.root.capabilities = append(wb.root.capabilities, tr)

	return nil
}
