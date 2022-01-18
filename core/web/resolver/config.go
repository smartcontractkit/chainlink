package resolver

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
)

type ConfigItemResolver struct {
	key string
	cfg ConfigItemValue
}

func NewConfigItem(key string, value ConfigItemValue) *ConfigItemResolver {
	return &ConfigItemResolver{key: key, cfg: value}
}

func (r *ConfigItemResolver) Key() string {
	return r.key
}

func (r *ConfigItemResolver) Value() ConfigItemValue {
	return r.cfg
}

type ConfigItemValue struct {
	Value interface{} `json:"value"`
}

func (ConfigItemValue) ImplementsGraphQLType(name string) bool {
	return name == "ConfigItemValue"
}

func (c *ConfigItemValue) UnmarshalGraphQL(input interface{}) error {
	switch t := input.(type) {
	case int32:
		c.Value = t
		return nil
	case string:
		c.Value = t
		return nil
	case []string:
		c.Value = t
		return nil
	default:
		return errors.New("wrong type")
	}
}

func (c ConfigItemValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

type ConfigResolver struct {
	cfg config.EnvPrinter
}

func NewConfig(cfg config.EnvPrinter) *ConfigResolver {
	return &ConfigResolver{cfg: cfg}
}

func (r *ConfigResolver) Items() []*ConfigItemResolver {
	var cfgs []*ConfigItemResolver

	schemaT := reflect.TypeOf(envvar.ConfigSchema{})
	t := reflect.TypeOf(r.cfg)
	v := reflect.ValueOf(r.cfg)

	for i := 0; i < t.NumField(); i++ {
		item := t.Field(i)

		// Using the same logic as the renderer does to render the configuration values
		schemaItem, ok := schemaT.FieldByName(item.Name)
		if !ok {
			continue
		}
		envName, ok := schemaItem.Tag.Lookup("env")
		if !ok {
			continue
		}
		field := v.FieldByIndex(item.Index)

		var cfg *ConfigItemResolver

		if stringer, ok := field.Interface().(fmt.Stringer); ok {
			if stringer != reflect.Zero(reflect.TypeOf(stringer)).Interface() {
				cfg = NewConfigItem(envName, ConfigItemValue{Value: stringer.String()})
			}
		} else {
			cfg = NewConfigItem(envName, ConfigItemValue{Value: fmt.Sprintf("%v", field)})
		}

		if cfg != nil {
			cfgs = append(cfgs, cfg)
		}
	}

	return cfgs
}

type ConfigPayloadResolver struct {
	cfg config.EnvPrinter
}

func NewConfigPayload(cfg config.EnvPrinter) *ConfigPayloadResolver {
	return &ConfigPayloadResolver{cfg: cfg}
}

func (r *ConfigPayloadResolver) Items() []*ConfigItemResolver {
	return NewConfig(r.cfg).Items()
}
