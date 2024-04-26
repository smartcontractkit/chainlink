package workflows

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type stepRequest struct {
	stepRef string
	state   executionState
}

// stepDefinition is the parsed representation of a step in a workflow.
//
// Within the workflow spec, they are called "Capability Properties".
type stepDefinition struct {
	// TODO: Rename this, type here refers to the capability ID, not its type.
	Type   string         `json:"type" jsonschema:"required"`
	Ref    string         `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_]+$"`
	Inputs map[string]any `json:"inputs,omitempty"`
	Config map[string]any `json:"config" jsonschema:"required"`

	CapabilityType capabilities.CapabilityType `json:"-"`
}

// workflowSpec is the parsed representation of a workflow.
type workflowSpec struct {
	Triggers  []stepDefinition `json:"triggers" jsonschema:"required"`
	Actions   []stepDefinition `json:"actions,omitempty"`
	Consensus []stepDefinition `json:"consensus" jsonschema:"required"`
	Targets   []stepDefinition `json:"targets" jsonschema:"required"`
}

func (w *workflowSpec) steps() []stepDefinition {
	s := []stepDefinition{}
	s = append(s, w.Actions...)
	s = append(s, w.Consensus...)
	s = append(s, w.Targets...)
	return s
}

// workflow is a directed graph of nodes, where each node is a step.
//
// triggers are special steps that are stored separately, they're
// treated differently due to their nature of being the starting
// point of a workflow.
type workflow struct {
	id string
	graph.Graph[string, *step]

	triggers []*triggerCapability

	spec *workflowSpec
}

func (w *workflow) walkDo(start string, do func(s *step) error) error {
	var outerErr error
	err := graph.BFS(w.Graph, start, func(ref string) bool {
		n, err := w.Graph.Vertex(ref)
		if err != nil {
			outerErr = err
			return true
		}

		err = do(n)
		if err != nil {
			outerErr = err
			return true
		}

		return false
	})
	if err != nil {
		return err
	}

	return outerErr
}

func (w *workflow) dependents(start string) ([]*step, error) {
	steps := []*step{}
	m, err := w.Graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	adj, ok := m[start]
	if !ok {
		return nil, fmt.Errorf("could not find step with ref %s", start)
	}

	for adjacentRef := range adj {
		n, err := w.Graph.Vertex(adjacentRef)
		if err != nil {
			return nil, err
		}

		steps = append(steps, n)
	}

	return steps, nil
}

// step wraps a stepDefinition with additional context for dependencies and execution
type step struct {
	stepDefinition
	dependencies      []string
	capability        capabilities.CallbackCapability
	config            *values.Map
	executionStrategy executionStrategy
}

type triggerCapability struct {
	stepDefinition
	trigger capabilities.TriggerCapability
	config  *values.Map
}

const (
	keywordTrigger = "trigger"
)

func Parse(yamlWorkflow string) (*workflow, error) {
	spec, err := ParseWorkflowSpecYaml(yamlWorkflow)
	if err != nil {
		return nil, err
	}

	// Construct and validate the graph. We instantiate an
	// empty graph with just one starting entry: `trigger`.
	// This provides the starting point for our graph and
	// points to all dependent steps.
	// Note: all triggers are represented by a single step called
	// `trigger`. This is because for workflows with multiple triggers
	// only one trigger will have started the workflow.
	stepHash := func(s *step) string {
		return s.Ref
	}
	g := graph.New(
		stepHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	err = g.AddVertex(&step{
		stepDefinition: stepDefinition{Ref: keywordTrigger},
	})
	if err != nil {
		return nil, err
	}

	// Next, let's populate the other entries in the graph.
	for _, s := range spec.steps() {
		// TODO: The workflow format spec doesn't always require a `Ref`
		// to be provided (triggers and targets don't have a `Ref` for example).
		// To handle this, we default the `Ref` to the type, but ideally we
		// should find a better long-term way to handle this.
		if s.Ref == "" {
			s.Ref = s.Type
		}

		innerErr := g.AddVertex(&step{stepDefinition: s})
		if innerErr != nil {
			return nil, fmt.Errorf("cannot add vertex %s: %w", s.Ref, innerErr)
		}
	}

	stepRefs, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	// Next, let's iterate over the steps and populate
	// any edges.
	for stepRef := range stepRefs {
		step, innerErr := g.Vertex(stepRef)
		if innerErr != nil {
			return nil, innerErr
		}

		refs, innerErr := findRefs(step.Inputs)
		if innerErr != nil {
			return nil, innerErr
		}
		step.dependencies = refs

		if stepRef != keywordTrigger && len(refs) == 0 {
			return nil, errors.New("all non-trigger steps must have a dependent ref")
		}

		for _, r := range refs {
			innerErr = g.AddEdge(r, step.Ref)
			if innerErr != nil {
				return nil, innerErr
			}
		}
	}

	triggerSteps := []*triggerCapability{}
	for _, t := range spec.Triggers {
		triggerSteps = append(triggerSteps, &triggerCapability{
			stepDefinition: t,
		})
	}
	wf := &workflow{
		spec:     &spec,
		Graph:    g,
		triggers: triggerSteps,
	}
	return wf, err
}
