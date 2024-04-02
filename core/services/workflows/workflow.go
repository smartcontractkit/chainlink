package workflows

import (
	"encoding/json"
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/invopop/jsonschema"
	"sigs.k8s.io/yaml"

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
	// A universally unique name for a capability will be defined under the “type” property. The uniqueness will, eventually, be enforced in the Capability Registry . Semver must be used to specify the version of the Capability at the end of the type field. Capability versions must be immutable.
	//
	// Initially, we will require major versions. This will ease upgrades early on while we develop the infrastructure.
	//
	// Eventually, we might support minor version and specific version pins. This will allow workflow authors to have flexibility when selecting the version, and node operators will be able to determine when they should update their capabilities.
	//
	// There are two ways to specify a type - using a string as a fully qualified ID or a structured table. When using a table, tags are ordered alphanumerically and joined into a string following a
	//  {type}:{tag1_key}_{tag1_value}:{tag2_key}_{tag2_value}@{version}
	// pattern.
	//
	// The “type” supports [a-z0-9_-:] characters followed by an @ and [semver regex] at the end.
	//
	// Validation must throw an error if:
	//
	// Unsupported characters are used.
	// (For Keystone only.) More specific than a major version is specified.
	//
	// Example (string)
	//  type: read_chain:chain_ethereum:network_mainnet@1
	//
	// Example (table)
	//
	//  type:
	//    name: read_chain
	//    version: 1
	//    tags:
	//      chain: ethereum
	//      network: mainnet
	//
	// [semver regex]: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	Type string `json:"type" jsonschema:"required,pattern=^[a-z0-9_\\-:]+@(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"`

	// Actions and Consensus capabilities have a required “ref” property that must be unique within a Workflow file (not universally) This property enables referencing outputs and is required because Actions and Consensus always need to be referenced in the following phases. Triggers can optionally specify  if they need to be referenced.
	//
	// The “ref” supports [a-z0-9_] characters.
	//
	// Validation must throw an error if:
	//  - Unsupported characters are used.
	//  - The same “ref” appears in the workflow multiple times.
	//  - “ref” is used on a Target capability.
	//  - “ref” has a circular reference.
	//
	// NOTE: Should introduce a custom validator to cover trigger case
	Ref string `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_]+$"`

	// Capabilities can specify an additional optional ”inputs” property. It allows specifying a dependency on the result of one or more other capabilities. These are always runtime values that cannot be provided upfront. It takes a map of the argument name internal to the capability and an explicit reference to the values.
	//
	// References are specified using the [type].[ref].[path_to_value] pattern.
	//
	// The interpolation of “inputs” is allowed
	//
	// Validation must throw an error if:
	//  - Input reference cannot be resolved.
	//  - Input is defined on triggers
	// NOTE: Should introduce a custom validator to cover trigger case
	Inputs map[string]any `json:"inputs,omitempty"`

	// The configuration of a Capability will be done using the “config” property. Each capability is responsible for defining an external interface used during setup. This interface may be unique or identical, meaning multiple Capabilities might use the same configuration properties.
	//
	// The interpolation of “inputs”
	//
	// Interpolation of self inputs is allowed from within the “config” property.
	//
	// Example
	//  targets:
	//    - type: write_polygon_mainnet@1
	//      inputs:
	//        report:
	//          - consensus.evm_median.outputs.report
	//      config:
	//        address: "0xaabbcc"
	//        method: "updateFeedValues(report bytes, role uint8)"
	//        params: [$(inputs.report), 1]
	Config map[string]any `json:"config" jsonschema:"required"`
}

// workflowSpec is the parsed representation of a workflow.
type workflowSpec struct {
	// Triggers define a starting condition for the workflow, based on specific events or conditions.
	Triggers []stepDefinition `json:"triggers" jsonschema:"required"`
	// Actions represent a discrete operation within the workflow, potentially transforming input data.
	Actions []stepDefinition `json:"actions,omitempty"`
	// Consensus encapsulates the logic for aggregating and validating the results from various nodes.
	Consensus []stepDefinition `json:"consensus" jsonschema:"required"`
	// Targets represents the final step of the workflow, delivering the processed data to a specified location.
	Targets []stepDefinition `json:"targets" jsonschema:"required"`
}

func GenerateJsonSchema() ([]byte, error) {
	schema := jsonschema.Reflect(&workflowSpec{})

	return json.MarshalIndent(schema, "", "  ")
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
	dependencies []string
	capability   capabilities.CallbackExecutable
	config       *values.Map
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
	spec := &workflowSpec{}
	err := yaml.Unmarshal([]byte(yamlWorkflow), spec)
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
		spec:     spec,
		Graph:    g,
		triggers: triggerSteps,
	}
	return wf, err
}
