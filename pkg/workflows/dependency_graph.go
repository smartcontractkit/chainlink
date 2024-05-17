package workflows

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

const (
	KeywordTrigger = "trigger"
)

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

type Vertex struct {
	StepDefinition
	Dependencies []string
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

func ParseDependencyGraph(yamlWorkflow string) (*DependencyGraph, error) {
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
		StepDefinition: StepDefinition{Ref: KeywordTrigger},
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
		step.Dependencies = refs

		if stepRef != KeywordTrigger && len(refs) == 0 {
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

var (
	InterpolationTokenRe = regexp.MustCompile(`^\$\((\S+)\)$`)
)

// findRefs takes an `inputs` map and returns a list of all the step references
// contained within it.
func findRefs(inputs map[string]any) ([]string, error) {
	refs := []string{}
	_, err := DeepMap(
		inputs,
		// This function is called for each string in the map
		// for each string, we iterate over each match of the interpolation token
		// - if there are no matches, return no reference
		// - if there is one match, return the reference
		// - if there are multiple matches (in the case of a multi-part state reference), return just the step ref
		func(el string) (any, error) {
			matches := InterpolationTokenRe.FindStringSubmatch(el)
			if len(matches) < 2 {
				return el, nil
			}

			m := matches[1]
			parts := strings.Split(m, ".")
			if len(parts) < 1 {
				return nil, fmt.Errorf("invalid ref %s", m)
			}

			refs = append(refs, parts[0])
			return el, nil
		},
	)
	return refs, err
}

// DeepMap recursively applies a transformation function
// over each string within:
//
//   - a map[string]any
//   - a []any
//   - a string
func DeepMap(input any, transform func(el string) (any, error)) (any, error) {
	// in the case of a string, simply apply the transformation
	// in the case of a map, recurse and apply the transformation to each value
	// in the case of a list, recurse and apply the transformation to each element
	switch tv := input.(type) {
	case string:
		nv, err := transform(tv)
		if err != nil {
			return nil, err
		}

		return nv, nil
	case mapping:
		// coerce mapping to map[string]any
		mp := map[string]any(tv)

		nm := map[string]any{}
		for k, v := range mp {
			nv, err := DeepMap(v, transform)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case map[string]any:
		nm := map[string]any{}
		for k, v := range tv {
			nv, err := DeepMap(v, transform)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case []any:
		a := []any{}
		for _, el := range tv {
			ne, err := DeepMap(el, transform)
			if err != nil {
				return nil, err
			}

			a = append(a, ne)
		}
		return a, nil
	}

	return nil, fmt.Errorf("cannot traverse item %+v of type %T", input, input)
}
