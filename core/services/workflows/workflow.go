package workflows

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Capability struct {
	Type   string         `yaml:"type"`
	Ref    string         `yaml:"ref"`
	Inputs map[string]any `yaml:"inputs"`
	Config map[string]any `yaml:"config"`
}

type workflowSpec struct {
	Triggers  []Capability `yaml:"triggers"`
	Actions   []Capability `yaml:"actions"`
	Consensus []Capability `yaml:"consensus"`
	Targets   []Capability `yaml:"targets"`
}

func (w *workflowSpec) steps() []Capability {
	s := []Capability{}
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
	Capability
	dependencies     []string
	cachedCapability capabilities.CallbackExecutable
	cachedConfig     *values.Map
}

type triggerCapability struct {
	Capability
	cachedTrigger capabilities.TriggerCapability
}

const (
	keywordTrigger = "trigger"
)

func Parse(yamlWorkflow string) (*workflow, error) {
	wfs := &workflowSpec{}
	err := yaml.Unmarshal([]byte(yamlWorkflow), wfs)
	if err != nil {
		return nil, err
	}

	// Construct and validate the graph. We instantiate an
	// empty graph with just one starting entry: `trigger`.
	// This provides the starting point for our graph and
	// points to all dependent nodes.
	nodeHash := func(n *node) string {
		return n.Ref
	}
	g := graph.New(
		nodeHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	err = g.AddVertex(&node{
		Capability: Capability{Ref: keywordTrigger},
	})
	if err != nil {
		return nil, err
	}

	for _, s := range wfs.steps() {
		if s.Ref == "" {
			s.Ref = s.Type
		}

		err := g.AddVertex(&node{Capability: s})
		if err != nil {
			return nil, fmt.Errorf("cannot add vertex %s: %w", s.Ref, err)
		}
	}

	nodeRefs, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	for nodeRef := range nodeRefs {
		node, err := g.Vertex(nodeRef)
		if err != nil {
			return nil, err
		}

		refs, innerErr := findRefs(node.Inputs)
		if innerErr != nil {
			return nil, innerErr
		}
		node.dependencies = refs

		for _, r := range refs {
			err = g.AddEdge(r, node.Ref)
			if err != nil {
				return nil, err
			}
		}
	}

	triggerNodes := []*triggerCapability{}
	for _, t := range wfs.Triggers {
		triggerNodes = append(triggerNodes, &triggerCapability{
			Capability: t,
		})
	}
	wf := &workflow{
		spec:     wfs,
		Graph:    g,
		triggers: triggerNodes,
	}
	return wf, err
}
