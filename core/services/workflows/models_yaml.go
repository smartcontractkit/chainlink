package workflows

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/invopop/jsonschema"
	"sigs.k8s.io/yaml"
)

func GenerateJsonSchema() ([]byte, error) {
	schema := jsonschema.Reflect(&workflowSpecYaml{})

	return json.MarshalIndent(schema, "", "  ")
}

func ParseWorkflowSpecYaml(data string) (workflowSpec, error) {
	w := workflowSpecYaml{}
	err := yaml.Unmarshal([]byte(data), &w)

	return w.toWorkflowSpec(), err
}

// workflowSpecYaml is the YAML representation of a workflow spec.
//
// It allows for multiple ways of defining a workflow spec, which we later
// convert to a single representation, `workflowSpec`.
type workflowSpecYaml struct {
	// Triggers define a starting condition for the workflow, based on specific events or conditions.
	Triggers []stepDefinitionYaml `json:"triggers" jsonschema:"required"`
	// Actions represent a discrete operation within the workflow, potentially transforming input data.
	Actions []stepDefinitionYaml `json:"actions,omitempty"`
	// Consensus encapsulates the logic for aggregating and validating the results from various nodes.
	Consensus []stepDefinitionYaml `json:"consensus" jsonschema:"required"`
	// Targets represents the final step of the workflow, delivering the processed data to a specified location.
	Targets []stepDefinitionYaml `json:"targets" jsonschema:"required"`
}

// toWorkflowSpec converts a workflowSpecYaml to a workflowSpec.
//
// We support multiple ways of defining a workflow spec yaml,
// but internally we want to work with a single representation.
func (w workflowSpecYaml) toWorkflowSpec() workflowSpec {
	triggers := make([]stepDefinition, 0, len(w.Triggers))
	for _, t := range w.Triggers {
		triggers = append(triggers, t.toStepDefinition())
	}

	actions := make([]stepDefinition, 0, len(w.Actions))
	for _, a := range w.Actions {
		actions = append(actions, a.toStepDefinition())
	}

	consensus := make([]stepDefinition, 0, len(w.Consensus))
	for _, c := range w.Consensus {
		consensus = append(consensus, c.toStepDefinition())
	}

	targets := make([]stepDefinition, 0, len(w.Targets))
	for _, t := range w.Targets {
		targets = append(targets, t.toStepDefinition())
	}

	return workflowSpec{
		Triggers:  triggers,
		Actions:   actions,
		Consensus: consensus,
		Targets:   targets,
	}
}

// stepDefinitionYaml is the YAML representation of a step in a workflow.
//
// It allows for multiple ways of defining a step, which we later
// convert to a single representation, `stepDefinition`.
type stepDefinitionYaml struct {
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
	Type stepDefinitionType `json:"type" jsonschema:"required"`

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

// toStepDefinition converts a stepDefinitionYaml to a stepDefinition.
//
// `stepDefinition` is the converged representation of a step in a workflow.
func (s stepDefinitionYaml) toStepDefinition() stepDefinition {
	return stepDefinition{
		Ref:    s.Ref,
		Type:   s.Type.String(),
		Inputs: s.Inputs,
		Config: s.Config,
	}
}

// stepDefinitionType represents both the string and table representations of the "type" field in a stepDefinition.
type stepDefinitionType struct {
	typeStr   string
	typeTable *stepDefinitionTableType
}

func (s stepDefinitionType) String() string {
	if s.typeStr != "" {
		return s.typeStr
	}

	return s.typeTable.String()
}

func (s *stepDefinitionType) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into a map to determine if it's a string or a table
	var m string
	err := json.Unmarshal(data, &m)
	if err == nil {
		s.typeStr = m
		return nil
	}

	// If the JSON data is a table, unmarshal it into a stepDefinitionTableType
	var table stepDefinitionTableType
	err = json.Unmarshal(data, &table)
	if err != nil {
		return err
	}
	s.typeTable = &table
	return nil
}

func (s *stepDefinitionType) MarshalJSON() ([]byte, error) {
	if s.typeStr != "" {
		return json.Marshal(s.typeStr)
	}

	return json.Marshal(s.typeTable)
}

// JSONSchema returns the JSON schema for a stepDefinitionType.
//
// The schema is a oneOf schema that allows either a string or a table.
func (stepDefinitionType) JSONSchema() *jsonschema.Schema {
	reflector := jsonschema.Reflector{DoNotReference: true, ExpandedStruct: true}
	tableSchema := reflector.Reflect(&stepDefinitionTableType{})
	stringSchema := &jsonschema.Schema{
		Type:    "string",
		Pattern: "^[a-z0-9_\\-:]+@(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$",
	}

	return &jsonschema.Schema{
		Title: "type",
		OneOf: []*jsonschema.Schema{
			stringSchema,
			tableSchema,
		},
	}
}

// stepDefinitionTableType is the structured representation of a stepDefinitionType.
type stepDefinitionTableType struct {
	Name    string            `json:"name"`
	Version string            `json:"version" jsonschema:"pattern=(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"`
	Tags    map[string]string `json:"tags"`
}

// String returns the string representation of a stepDefinitionTableType.
//
// It follows the format:
//
//	{name}:{tag1_key}_{tag1_value}:{tag2_key}_{tag2_value}@{version}
//
// where tags are ordered alphanumerically.
func (s stepDefinitionTableType) String() string {
	tags := make([]string, 0, len(s.Tags))
	for k, v := range s.Tags {
		tags = append(tags, fmt.Sprintf("%s_%s", k, v))
	}
	slices.Sort(tags)

	return fmt.Sprintf("%s:%s@%s", s.Name, strings.Join(tags, ":"), s.Version)
}
