package workflows

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type stepRequest struct {
	stepRef string
	state   store.WorkflowExecution
}

// StepDefinition is the parsed representation of a step in a workflow.
//
// Within the workflow spec, they are called "Capability Properties".
type StepDefinition struct {
	ID     string         `json:"id" jsonschema:"required"`
	Ref    string         `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_]+$"`
	Inputs map[string]any `json:"inputs,omitempty"`
	Config map[string]any `json:"config" jsonschema:"required"`

	CapabilityType capabilities.CapabilityType `json:"-"`
}

// WorkflowSpec is the parsed representation of a workflow.
type WorkflowSpec struct {
	Triggers  []StepDefinition `json:"triggers" jsonschema:"required"`
	Actions   []StepDefinition `json:"actions,omitempty"`
	Consensus []StepDefinition `json:"consensus" jsonschema:"required"`
	Targets   []StepDefinition `json:"targets" jsonschema:"required"`
}

func (w *WorkflowSpec) Steps() []StepDefinition {
	s := []StepDefinition{}
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

	spec *WorkflowSpec
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

// step wraps a Vertex with additional context for execution that is mutated by the engine
type step struct {
	Vertex
	capability        capabilities.CallbackCapability
	config            *values.Map
	executionStrategy executionStrategy
}

type Vertex struct {
	StepDefinition
	dependencies []string
}

// DependencyGraph is an intermediate representation of a workflow wherein all the graph
// vertices are represented and validated. It is a static representation of the workflow dependencies.
type DependencyGraph struct {
	ID string
	graph.Graph[string, *Vertex]

	Triggers []*StepDefinition

	Spec *WorkflowSpec
}

// VID is an identifier for a Vertex that can be used to uniquely identify it in a graph.
// it represents the notion `hash` in the graph package AddVertex method.
// we refrain from naming it `hash` to avoid confusion with the hash function.
func (v *Vertex) VID() string {
	return v.Ref
}

type triggerCapability struct {
	StepDefinition
	trigger capabilities.TriggerCapability
	config  *values.Map
}

const (
	keywordTrigger = "trigger"
)

func Parse(yamlWorkflow string) (*workflow, error) {
	wf2, err := ParseDepedencyGraph(yamlWorkflow)
	if err != nil {
		return nil, err
	}
	return createWorkflow(wf2)
}

func ParseDepedencyGraph(yamlWorkflow string) (*DependencyGraph, error) {
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
	stepHash := func(s *Vertex) string {
		return s.VID()
	}
	g := graph.New(
		stepHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	err = g.AddVertex(&Vertex{
		StepDefinition: StepDefinition{Ref: keywordTrigger},
	})
	if err != nil {
		return nil, err
	}

	// Next, let's populate the other entries in the graph.
	for _, s := range spec.Steps() {
		// TODO: The workflow format spec doesn't always require a `Ref`
		// to be provided (triggers and targets don't have a `Ref` for example).
		// To handle this, we default the `Ref` to the type, but ideally we
		// should find a better long-term way to handle this.
		if s.Ref == "" {
			s.Ref = s.ID
		}

		innerErr := g.AddVertex(&Vertex{StepDefinition: s})
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

	triggerSteps := []*StepDefinition{}
	for _, t := range spec.Triggers {
		tt := t
		triggerSteps = append(triggerSteps, &tt)
	}
	wf := &DependencyGraph{
		Spec:     &spec,
		Graph:    g,
		Triggers: triggerSteps,
	}
	return wf, err
}

// createWorkflow converts a StaticWorkflow to an executable workflow
// by adding metadata to the vertices that is owned by the workflow runtime.
func createWorkflow(wf2 *DependencyGraph) (*workflow, error) {
	out := &workflow{
		id:       wf2.ID,
		triggers: []*triggerCapability{},
		spec:     wf2.Spec,
	}

	for _, t := range wf2.Triggers {
		out.triggers = append(out.triggers, &triggerCapability{
			StepDefinition: *t,
		})
	}

	stepHash := func(s *step) string {
		// must use the same hash function as the DependencyGraph.
		// this ensures that the intermediate representation (DependencyGraph) and the workflow
		// representation label vertices with the same identifier, which in turn allows us to
		// to copy the edges from the intermediate representation to the executable representation.
		return s.Vertex.VID()
	}
	g := graph.New(
		stepHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	adjMap, err := wf2.Graph.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert intermediate representation to adjacency map: %w", err)
	}

	// copy the all the vertices from the intermediate graph to the executable workflow graph
	for vertexRef := range adjMap {
		v, innerErr := wf2.Graph.Vertex(vertexRef)
		if innerErr != nil {
			return nil, fmt.Errorf("failed to retrieve vertex for %s: %w", vertexRef, innerErr)
		}
		innerErr = g.AddVertex(&step{Vertex: *v})
		if innerErr != nil {
			return nil, fmt.Errorf("failed to add vertex to executable workflow %s: %w", vertexRef, innerErr)
		}
	}
	// now we can add all the edges. this works because we are using vertex hash function is the same in both graphs.
	// see comment on `stepHash` function.
	for vertexRef, edgeRefs := range adjMap {
		for edgeRef := range edgeRefs {
			innerErr := g.AddEdge(vertexRef, edgeRef)
			if innerErr != nil {
				return nil, fmt.Errorf("failed to add edge from '%s' to '%s': %w", vertexRef, edgeRef, innerErr)
			}
		}
	}
	out.Graph = g
	return out, nil
}
