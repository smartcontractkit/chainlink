package pipeline

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

type ObjectType int

const (
	NilType ObjectType = iota
	BoolType
	DecimalType
	StringType
	SliceType
	MapType
)

// ObjectParam represents a kind of any type that could be used by the
// memo task
type ObjectParam struct {
	Type         ObjectType
	BoolValue    BoolParam
	DecimalValue DecimalParam
	StringValue  StringParam
	SliceValue   SliceParam
	MapValue     MapParam
}

func (o ObjectParam) MarshalJSON() ([]byte, error) {
	switch o.Type {
	case NilType:
		return json.Marshal(nil)
	case BoolType:
		return json.Marshal(o.BoolValue)
	case DecimalType:
		return json.Marshal(o.DecimalValue.Decimal())
	case StringType:
		return json.Marshal(o.StringValue)
	case MapType:
		return json.Marshal(o.MapValue)
	case SliceType:
		return json.Marshal(o.SliceValue)
	}
	panic(fmt.Sprintf("Invalid type for ObjectParam %v", o.Type))
}

func (o ObjectParam) Marshal() (string, error) {
	b, err := o.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (o ObjectParam) String() string {
	value, err := o.Marshal()
	if err != nil {
		return fmt.Sprintf("<error Stringifying: %v>", err)
	}
	return value
}

func (o *ObjectParam) UnmarshalPipelineParam(val interface{}) error {
	switch v := val.(type) {
	case nil:
		o.Type = NilType
		return nil

	case bool:
		o.Type = BoolType
		o.BoolValue = BoolParam(v)
		return nil

	case uint8, uint16, uint32, uint64, uint, int8, int16, int32, int64, int, float32, float64, decimal.Decimal, *decimal.Decimal, big.Int, *big.Int:
		o.Type = DecimalType
		return o.DecimalValue.UnmarshalPipelineParam(v)

	case string:
		o.Type = StringType
		return o.StringValue.UnmarshalPipelineParam(v)

		// Maps
	case MapParam:
		o.Type = MapType
		o.MapValue = v
		return nil

	case map[string]interface{}:
		o.Type = MapType
		return o.MapValue.UnmarshalPipelineParam(v)

		// Slices
	case SliceParam:
		o.Type = SliceType
		o.SliceValue = v
		return nil

	case []interface{}:
		o.Type = SliceType
		return o.SliceValue.UnmarshalPipelineParam(v)

	case []int:
		o.Type = SliceType
		for _, value := range v {
			o.SliceValue = append(o.SliceValue, value)
		}
		return nil

	case []string:
		o.Type = SliceType
		for _, value := range v {
			o.SliceValue = append(o.SliceValue, value)
		}
		return nil

	case ObjectParam:
		o.Type = v.Type
		o.BoolValue = v.BoolValue
		o.MapValue = v.MapValue
		o.StringValue = v.StringValue
		o.DecimalValue = v.DecimalValue
		return nil
	}

	return fmt.Errorf("bad input for task: %T", val)
}
