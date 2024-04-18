package capabilities

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
	jsonvalidate "github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Validator[Config any, Inputs any, Outputs any] struct {
	ValidatorArgs

	configType  Config
	inputsType  Inputs
	outputsType Outputs

	schemas map[string]string
}

var uriPrefix = "https://github.com/smartcontractkit/chainlink/capabilities/"

var _ Validatable = (*Validator[any, any, any])(nil)

type ValidatorArgs struct {
	Info CapabilityInfo

	// You can customize each one of the reflectors
	// or leave them nil to use the default reflector.
	//
	// You can also override the default reflector by setting
	// the DefaultReflector field.
	DefaultReflector *jsonschema.Reflector

	SchemaReflector  *jsonschema.Reflector
	ConfigReflector  *jsonschema.Reflector
	InputsReflector  *jsonschema.Reflector
	OutputsReflector *jsonschema.Reflector
}

func NewValidator[Config any, Inputs any, Outputs any](args ValidatorArgs) Validator[Config, Inputs, Outputs] {
	baseID := jsonschema.ID(uriPrefix)
	defaultReflector := &jsonschema.Reflector{ExpandedStruct: true, DoNotReference: true,
		BaseSchemaID: baseID,
	}
	if args.DefaultReflector != nil {
		defaultReflector = args.DefaultReflector
	}

	if args.SchemaReflector == nil {
		args.SchemaReflector = defaultReflector
	}

	if args.ConfigReflector == nil {
		args.ConfigReflector = defaultReflector
	}

	if args.InputsReflector == nil {
		args.InputsReflector = defaultReflector
	}

	if args.OutputsReflector == nil {
		args.OutputsReflector = defaultReflector
	}

	return Validator[Config, Inputs, Outputs]{
		ValidatorArgs: args,
		schemas:       make(map[string]string),
	}
}

func (v *Validator[Config, Inputs, Outputs]) Schema() (string, error) {
	type combined struct {
		Config  Config  `json:"config"`
		Inputs  Inputs  `json:"inputs"`
		Outputs Outputs `json:"outputs"`
	}
	c := combined{}
	type combinedWithInputs struct {
		Config  Config  `json:"config"`
		Outputs Outputs `json:"outputs"`
	}
	ci := combinedWithInputs{}

	var config interface{} = c.Config
	var inputs interface{} = c.Inputs
	var outputs interface{} = c.Outputs
	if config == nil {
		return "", errors.New("config is nil, please provide a config type")
	}
	if outputs == nil {
		return "", errors.New("outputs is nil, please provide an outputs type")
	}

	// we allow inputs to be nil, since triggers do not have inputss
	if inputs == nil {
		return schemaWith(*v.SchemaReflector, ci, v.schemas, "root", v.Info)
	}

	// print values of combined
	return schemaWith(*v.SchemaReflector, c, v.schemas, "root", v.Info)
}
func (v *Validator[Config, Inputs, Outputs]) ConfigSchema() (string, error) {
	return schemaWith(*v.ConfigReflector, v.configType, v.schemas, "config", v.Info)
}
func (v *Validator[Config, Inputs, Outputs]) InputsSchema() (string, error) {
	return schemaWith(*v.InputsReflector, v.inputsType, v.schemas, "inputs", v.Info)
}
func (v *Validator[Config, Inputs, Outputs]) OutputsSchema() (string, error) {
	return schemaWith(*v.OutputsReflector, v.outputsType, v.schemas, "outputs", v.Info)
}

func (v *Validator[Config, Inputs, Outputs]) ValidateInputs(inputs *values.Map) (*Inputs, error) {
	inputsSchema, err := v.InputsSchema()
	if err != nil {
		return nil, errors.Join(errors.New("validation error while validating inputs"), err)
	}

	return validateAgainstSchema[Inputs](inputs, inputsSchema)
}
func (v *Validator[Config, Inputs, Outputs]) ValidateConfig(config *values.Map) (*Config, error) {
	configSchema, err := v.ConfigSchema()
	if err != nil {
		return nil, errors.Join(errors.New("validation error while validating config"), err)
	}

	return validateAgainstSchema[Config](config, configSchema)
}
func (v *Validator[Config, Inputs, Outputs]) ValidateOutputs(outputs *values.Map) (*Outputs, error) {
	outputsSchema, err := v.OutputsSchema()
	if err != nil {
		return nil, errors.Join(errors.New("validation error while validating outputs"), err)
	}

	return validateAgainstSchema[Outputs](outputs, outputsSchema)
}

func validateAgainstSchema[DecodedValue any](value *values.Map, schema string) (*DecodedValue, error) {
	jsonSchema, err := jsonvalidate.CompileString(uriPrefix, schema)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, fmt.Errorf("cannot validate nil value against schema: %s", jsonSchema.Location)
	}
	// parse
	decodedValue := new(DecodedValue)
	err = value.UnwrapTo(decodedValue)
	if err != nil {
		return nil, err
	}

	// validate
	jsonValue, err := json.Marshal(decodedValue)
	if err != nil {
		return nil, err
	}
	var jsonRaw any
	err = json.Unmarshal(jsonValue, &jsonRaw)
	if err != nil {
		return nil, err
	}

	err = jsonSchema.Validate(jsonRaw)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error validating value %v", jsonRaw), err)
	}

	return decodedValue, err
}

func schemaWith(reflector jsonschema.Reflector, schemaType any, schemaCache map[string]string, key string, info CapabilityInfo) (string, error) {
	if schema, ok := schemaCache[key]; ok {
		return schema, nil
	}

	idPath := uriPrefix + info.ID + "/" + key
	if info.ID == "" {
		idPath = uriPrefix + key
	}

	id := jsonschema.ID(idPath)
	err := id.Validate()
	if err != nil {
		return "", errors.Join(errors.New("invalid schema ID"), err)
	}

	schema := reflector.Reflect(schemaType)
	schema.ID = id
	schema.Description = info.Description
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")

	if err != nil {
		return "", errors.Join(errors.New("unable to marshal schema"), err)
	}

	schemaCache[key] = string(schemaBytes)
	return string(schemaBytes), nil
}
