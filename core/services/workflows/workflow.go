package workflows

import (
	"fmt"

	"gopkg.in/yaml.v3"

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
	*graph[*node]

	triggers []*triggerCapability

	spec *workflowSpec
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
	nodes := map[string]*node{
		keywordTrigger: {Capability: Capability{Ref: keywordTrigger}},
	}
	adjacencies := map[string]map[string]struct{}{
		keywordTrigger: {},
	}
	graph := &graph[*node]{
		adjacencies: adjacencies,
		nodes:       nodes,
	}
	for _, s := range wfs.steps() {
		// For steps that don't have a ref, use
		// the node's type as a default.
		if s.Ref == "" {
			s.Ref = s.Type
		}

		_, ok := nodes[s.Ref]
		if ok {
			return nil, fmt.Errorf("duplicate reference %s found in workflow spec", s.Ref)
		}

		nodes[s.Ref] = &node{Capability: s}
		adjacencies[s.Ref] = map[string]struct{}{}
	}

	for _, nd := range nodes {
		refs, innerErr := findRefs(nd.Inputs)
		if innerErr != nil {
			return nil, innerErr
		}
		nd.dependencies = refs

		for _, r := range refs {
			_, ok := nodes[r]
			if !ok && r != keywordTrigger {
				return nil, fmt.Errorf("invalid reference %s found in workflow spec", r)
			}

			adjacencies[r][nd.Ref] = struct{}{}

			var found bool
			innerErr := graph.walkDo(nd.Ref, func(n *node) error {
				if n.Ref == r {
					found = true
					return nil
				}

				return nil
			})
			if innerErr != nil {
				return nil, innerErr
			}

			if found {
				return nil, fmt.Errorf("found circular relationship between %s and %s", r, nd.Ref)
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
		graph:    graph,
		triggers: triggerNodes,
	}
	return wf, err
}
