package values

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Unwrappable interface {
	Unwrap() (any, error)
}

type Value interface {
	proto() *pb.Value

	Unwrappable
}

func Wrap(v any) (Value, error) {
	switch tv := v.(type) {
	case map[string]any:
		return NewMap(tv)
	case string:
		return NewString(tv), nil
	case bool:
		return NewBool(tv), nil
	case []byte:
		return NewBytes(tv), nil
	case []any:
		return NewList(tv)
	case decimal.Decimal:
		return NewDecimal(tv), nil
	case int64:
		return NewInt64(tv), nil
	case int:
		return NewInt64(int64(tv)), nil
	case nil:
		return nil, nil

	// Transparently wrap values.
	// This is helpful for recursive wrapping of values.
	case *Map:
		return tv, nil
	case *List:
		return tv, nil
	case *String:
		return tv, nil
	case *Bytes:
		return tv, nil
	case *Decimal:
		return tv, nil
	case *Int64:
		return tv, nil
	}

	return nil, fmt.Errorf("could not wrap into value: %+v", v)
}

func Unwrap(v Value) (any, error) {
	if v == nil {
		return nil, nil
	}

	return v.Unwrap()
}

func Proto(v Value) *pb.Value {
	if v == nil {
		return &pb.Value{}
	}

	return v.proto()
}

func FromProto(val *pb.Value) Value {
	if val == nil {
		return nil
	}

	switch val.Value.(type) {
	case nil:
		return nil
	case *pb.Value_StringValue:
		return NewString(val.GetStringValue())
	case *pb.Value_BoolValue:
		return NewBool(val.GetBoolValue())
	case *pb.Value_DecimalValue:
		return fromDecimalValueProto(val.GetDecimalValue())
	case *pb.Value_Int64Value:
		return NewInt64(val.GetInt64Value())
	case *pb.Value_BytesValue:
		return NewBytes(val.GetBytesValue())
	case *pb.Value_ListValue:
		return FromListValueProto(val.GetListValue())
	case *pb.Value_MapValue:
		return FromMapValueProto(val.GetMapValue())
	}

	panic(fmt.Errorf("unsupported type %T: %+v", val, val))
}

func FromMapValueProto(mv *pb.Map) *Map {
	nm := map[string]Value{}
	for k, v := range mv.Fields {
		nm[k] = FromProto(v)
	}
	return &Map{Underlying: nm}
}

func FromListValueProto(lv *pb.List) *List {
	nl := []Value{}
	for _, el := range lv.Fields {
		nl = append(nl, FromProto(el))
	}
	return &List{Underlying: nl}
}

func fromDecimalValueProto(decStr string) *Decimal {
	dec, err := decimal.NewFromString(decStr)
	if err != nil {
		panic(err)
	}

	return NewDecimal(dec)
}
