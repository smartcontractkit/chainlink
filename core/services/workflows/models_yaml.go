package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/shopspring/decimal"
	"sigs.k8s.io/yaml"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

func GenerateJsonSchema() ([]byte, error) {
	schema := jsonschema.Reflect(&workflowSpecYaml{})

	return json.MarshalIndent(schema, "", "  ")
}

func ParseWorkflowSpecYaml(data string) (WorkflowSpec, error) {
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
func (w workflowSpecYaml) toWorkflowSpec() WorkflowSpec {
	triggers := make([]StepDefinition, 0, len(w.Triggers))
	for _, t := range w.Triggers {
		sd := t.toStepDefinition()
		sd.CapabilityType = capabilities.CapabilityTypeTrigger
		triggers = append(triggers, sd)
	}

	actions := make([]StepDefinition, 0, len(w.Actions))
	for _, a := range w.Actions {
		sd := a.toStepDefinition()
		sd.CapabilityType = capabilities.CapabilityTypeAction
		actions = append(actions, sd)
	}

	consensus := make([]StepDefinition, 0, len(w.Consensus))
	for _, c := range w.Consensus {
		sd := c.toStepDefinition()
		sd.CapabilityType = capabilities.CapabilityTypeConsensus
		consensus = append(consensus, sd)
	}

	targets := make([]StepDefinition, 0, len(w.Targets))
	for _, t := range w.Targets {
		sd := t.toStepDefinition()
		sd.CapabilityType = capabilities.CapabilityTypeTarget
		targets = append(targets, sd)
	}

	return WorkflowSpec{
		Triggers:  triggers,
		Actions:   actions,
		Consensus: consensus,
		Targets:   targets,
	}
}

type mapping map[string]any

func (m *mapping) UnmarshalJSON(b []byte) error {
	mp := map[string]any{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()

	err := d.Decode(&mp)
	if err != nil {
		return err
	}

	nm, err := convertNumbers(mp)
	if err != nil {
		return err
	}

	*m = (mapping)(nm)
	return err
}

func convertNumber(el any) (any, error) {
	switch elv := el.(type) {
	case json.Number:
		if strings.Contains(elv.String(), ".") {
			f, err := elv.Float64()
			if err == nil {
				return decimal.NewFromFloat(f), nil
			}
		}

		return elv.Int64()
	default:
		return el, nil
	}
}

func convertNumbers(m map[string]any) (map[string]any, error) {
	nm := map[string]any{}
	for k, v := range m {
		switch tv := v.(type) {
		case map[string]any:
			cm, err := convertNumbers(tv)
			if err != nil {
				return nil, err
			}

			nm[k] = cm
		case []any:
			na := make([]any, len(tv))
			for i, v := range tv {
				cv, err := convertNumber(v)
				if err != nil {
					return nil, err
				}

				na[i] = cv
			}

			nm[k] = na
		default:
			cv, err := convertNumber(v)
			if err != nil {
				return nil, err
			}

			nm[k] = cv
		}
	}

	return nm, nil
}

func (m mapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(m))
}

// stepDefinitionYaml is the YAML representation of a step in a workflow.
//
// It allows for multiple ways of defining a step, which we later
// convert to a single representation, `stepDefinition`.
type stepDefinitionYaml struct {
	// A universally unique name for a capability will be defined under the “id” property. The uniqueness will, eventually, be enforced in the Capability Registry.
	//
	// Semver must be used to specify the version of the Capability at the end of the id field. Capability versions must be immutable.
	//
	// Initially, we will require major versions. This will ease upgrades early on while we develop the infrastructure.
	//
	// Eventually, we might support minor version and specific version pins. This will allow workflow authors to have flexibility when selecting the version, and node operators will be able to determine when they should update their capabilities.
	//
	// There are two ways to specify an id - using a string as a fully qualified ID or a structured table. When using a table, labels are ordered alphanumerically and joined into a string following a
	//  {name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}@{version}
	// pattern.
	//
	// The “id” supports [a-z0-9_-:] characters followed by an @ and [semver regex] at the end.
	//
	// Validation must throw an error if:
	//
	// Unsupported characters are used.
	// (For Keystone only.) More specific than a major version is specified.
	//
	// Example (string)
	//  id: read_chain:chain_ethereum:network_mainnet@1
	//
	// Example (table)
	//
	//  id:
	//    name: read_chain
	//    version: 1
	//    labels:
	//      chain: ethereum
	//      network: mainnet
	//
	// [semver regex]: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	ID stepDefinitionID `json:"id" jsonschema:"required"`

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
	Ref string `json:"ref,omitempty" jsonschema:"pattern=^[a-z0-9_-]+$"`

	// Capabilities can specify an additional optional ”inputs” property. It allows specifying a dependency on the result of one or more other capabilities. These are always runtime values that cannot be provided upfront. It takes a map of the argument name internal to the capability and an explicit reference to the values.
	//
	// References are specified using the [id].[ref].[path_to_value] pattern.
	//
	// The interpolation of “inputs” is allowed
	//
	// Validation must throw an error if:
	//  - Input reference cannot be resolved.
	//  - Input is defined on triggers
	// NOTE: Should introduce a custom validator to cover trigger case
	Inputs mapping `json:"inputs,omitempty"`

	// The configuration of a Capability will be done using the “config” property. Each capability is responsible for defining an external interface used during setup. This interface may be unique or identical, meaning multiple Capabilities might use the same configuration properties.
	//
	// The interpolation of “inputs”
	//
	// Interpolation of self inputs is allowed from within the “config” property.
	//
	// Example
	//  targets:
	//    - id: write_polygon_mainnet@1
	//      inputs:
	//        report:
	//          - consensus.evm_median.outputs.report
	//      config:
	//        address: "0xaabbcc"
	//        method: "updateFeedValues(report bytes, role uint8)"
	//        params: [$(inputs.report), 1]
	Config mapping `json:"config" jsonschema:"required"`
}

// toStepDefinition converts a stepDefinitionYaml to a stepDefinition.
//
// `stepDefinition` is the converged representation of a step in a workflow.
func (s stepDefinitionYaml) toStepDefinition() StepDefinition {
	return StepDefinition{
		Ref:    s.Ref,
		ID:     s.ID.String(),
		Inputs: s.Inputs,
		Config: s.Config,
	}
}

// stepDefinitionID represents both the string and table representations of the "id" field in a stepDefinition.
type stepDefinitionID struct {
	idStr   string
	idTable *stepDefinitionTableID
}

func (s stepDefinitionID) String() string {
	if s.idStr != "" {
		return s.idStr
	}

	return s.idTable.String()
}

func (s *stepDefinitionID) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into a map to determine if it's a string or a table
	var m string
	err := json.Unmarshal(data, &m)
	if err == nil {
		s.idStr = m
		return nil
	}

	// If the JSON data is a table, unmarshal it into a stepDefinitionTableID
	var table stepDefinitionTableID
	err = json.Unmarshal(data, &table)
	if err != nil {
		return err
	}
	s.idTable = &table
	return nil
}

func (s *stepDefinitionID) MarshalJSON() ([]byte, error) {
	if s.idStr != "" {
		return json.Marshal(s.idStr)
	}

	return json.Marshal(s.idTable)
}

// JSONSchema returns the JSON schema for a stepDefinitionID.
//
// The schema is a oneOf schema that allows either a string or a table.
func (stepDefinitionID) JSONSchema() *jsonschema.Schema {
	reflector := jsonschema.Reflector{DoNotReference: true, ExpandedStruct: true}
	tableSchema := reflector.Reflect(&stepDefinitionTableID{})
	stringSchema := &jsonschema.Schema{
		ID:      "string",
		Pattern: "^[a-z0-9_\\-:]+@(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$",
	}

	return &jsonschema.Schema{
		Title: "id",
		OneOf: []*jsonschema.Schema{
			stringSchema,
			tableSchema,
		},
	}
}

// stepDefinitionTableID is the structured representation of a stepDefinitionID.
type stepDefinitionTableID struct {
	Name    string            `json:"name"`
	Version string            `json:"version" jsonschema:"pattern=(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"`
	Labels  map[string]string `json:"labels"`
}

// String returns the string representation of a stepDefinitionTableID.
//
// It follows the format:
//
//	{name}:{label1_key}_{label1_value}:{label2_key}_{label2_value}@{version}
//
// where labels are ordered alphanumerically.
func (s stepDefinitionTableID) String() string {
	labels := make([]string, 0, len(s.Labels))
	for k, v := range s.Labels {
		labels = append(labels, fmt.Sprintf("%s_%s", k, v))
	}
	slices.Sort(labels)

	return fmt.Sprintf("%s:%s@%s", s.Name, strings.Join(labels, ":"), s.Version)
}
