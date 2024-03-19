package workflows

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type stepRequest struct {
	executionID string
	stepRef     string
	state       executionState
}

type capabilityDefinition struct {
	Type   string         `yaml:"type"`
	Ref    string         `yaml:"ref"`
	Inputs map[string]any `yaml:"inputs"`
	Config map[string]any `yaml:"config"`
}

type workflowSpec struct {
	Triggers  []capabilityDefinition `yaml:"triggers"`
	Actions   []capabilityDefinition `yaml:"actions"`
	Consensus []capabilityDefinition `yaml:"consensus"`
	Targets   []capabilityDefinition `yaml:"targets"`
}

func (w *workflowSpec) steps() []capabilityDefinition {
	s := []capabilityDefinition{}
	s = append(s, w.Actions...)
	s = append(s, w.Consensus...)
	s = append(s, w.Targets...)
	return s
}

type workflow struct {
	graph.Graph[string, *node]

	triggers []*triggerCapability

	spec *workflowSpec
}

func (w *workflow) walkDo(start string, do func(n *node) error) error {
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

func (w *workflow) adjacentNodes(start string) ([]*node, error) {
	nodes := []*node{}
	m, err := w.Graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	adj, ok := m[start]
	if !ok {
		return nil, fmt.Errorf("could not find node with ref %s", start)
	}

	for adjacentRef := range adj {
		n, err := w.Graph.Vertex(adjacentRef)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

type node struct {
	capabilityDefinition
	dependencies []string
	capability   capabilities.CallbackExecutable
	config       *values.Map
}

type triggerCapability struct {
	capabilityDefinition
	trigger capabilities.TriggerCapability
}

const (
	keywordTrigger = "trigger"
)

func Parse(yamlWorkflow string) (*workflow, error) {
	spec := &workflowSpec{}
	err := yaml.Unmarshal([]byte(yamlWorkflow), spec)
	if err != nil {
		return nil, err
	}

	// Construct and validate the graph. We instantiate an
	// empty graph with just one starting entry: `trigger`.
	// This provides the starting point for our graph and
	// points to all dependent nodes.
	// Note: all triggers are represented by a single node called
	// `trigger`. This is because for workflows with multiple triggers
	// only one trigger will have started the workflow.
	nodeHash := func(n *node) string {
		return n.Ref
	}
	g := graph.New(
		nodeHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	err = g.AddVertex(&node{
		capabilityDefinition: capabilityDefinition{Ref: keywordTrigger},
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

		innerErr := g.AddVertex(&node{capabilityDefinition: s})
		if innerErr != nil {
			return nil, fmt.Errorf("cannot add vertex %s: %w", s.Ref, innerErr)
		}
	}

	nodeRefs, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	// Next, let's iterate over the nodes and populate
	// any edges.
	for nodeRef := range nodeRefs {
		node, innerErr := g.Vertex(nodeRef)
		if innerErr != nil {
			return nil, innerErr
		}

		refs, innerErr := findRefs(node.Inputs)
		if innerErr != nil {
			return nil, innerErr
		}
		node.dependencies = refs

		for _, r := range refs {
			innerErr = g.AddEdge(r, node.Ref)
			if innerErr != nil {
				return nil, innerErr
			}
		}
	}

	triggerNodes := []*triggerCapability{}
	for _, t := range spec.Triggers {
		triggerNodes = append(triggerNodes, &triggerCapability{
			capabilityDefinition: t,
		})
	}
	wf := &workflow{
		spec:     spec,
		Graph:    g,
		triggers: triggerNodes,
	}
	return wf, err
}
