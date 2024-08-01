package values

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Unwrappable interface {
	Unwrap() (any, error)
	UnwrapTo(any) error
}

type Value interface {
	proto() *pb.Value

	Copy() Value
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
	case uint64:
		return NewInt64(int64(tv)), nil
	case uint:
		return NewInt64(int64(tv)), nil
	case uint32:
		return NewInt64(int64(tv)), nil
	case *big.Int:
		return NewBigInt(tv), nil
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

	// Handle slices, structs, and pointers to structs
	val := reflect.ValueOf(v)
	// nolint
	switch val.Kind() {
	// Better complex type support for maps
	case reflect.Map:
		m := make(map[string]any, val.Len())
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			ks, ok := k.Interface().(string)
			if !ok {
				return nil, fmt.Errorf("could not wrap into value %+v", v)
			}
			v := iter.Value()
			m[ks] = v.Interface()
		}
		return NewMap(m)
	// Better complex type support for slices
	case reflect.Slice:
		s := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			s[i] = item
		}
		return NewList(s)
	case reflect.Struct:
		return createMapFromStruct(v)
	case reflect.Pointer:
		if reflect.Indirect(reflect.ValueOf(v)).Kind() == reflect.Struct {
			return createMapFromStruct(reflect.Indirect(reflect.ValueOf(v)).Interface())
		}
	}

	return nil, fmt.Errorf("could not wrap into value: %+v", v)
}

func WrapMap(a any) (*Map, error) {
	v, err := Wrap(a)
	if err != nil {
		return nil, err
	}

	vm, ok := v.(*Map)
	if !ok {
		return nil, fmt.Errorf("could not wrap %+v to map", a)
	}

	return vm, nil
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

func ProtoMap(v *Map) *pb.Map {
	return Proto(v).GetMapValue()
}

func FromProto(val *pb.Value) (Value, error) {
	if val == nil {
		return nil, nil
	}

	switch val.Value.(type) {
	case nil:
		return nil, nil
	case *pb.Value_StringValue:
		return NewString(val.GetStringValue()), nil
	case *pb.Value_BoolValue:
		return NewBool(val.GetBoolValue()), nil
	case *pb.Value_DecimalValue:
		return fromDecimalValueProto(val.GetDecimalValue()), nil
	case *pb.Value_Int64Value:
		return NewInt64(val.GetInt64Value()), nil
	case *pb.Value_BytesValue:
		return NewBytes(val.GetBytesValue()), nil
	case *pb.Value_ListValue:
		return FromListValueProto(val.GetListValue())
	case *pb.Value_MapValue:
		return FromMapValueProto(val.GetMapValue())
	case *pb.Value_BigintValue:
		return fromBigIntValueProto(val.GetBigintValue()), nil
	}

	return nil, fmt.Errorf("unsupported type %T: %+v", val, val)
}

func FromMapValueProto(mv *pb.Map) (*Map, error) {
	if mv == nil {
		return nil, nil
	}

	nm := map[string]Value{}
	for k, v := range mv.Fields {
		inner, err := FromProto(v)
		if err != nil {
			return nil, err
		}
		nm[k] = inner
	}
	return &Map{Underlying: nm}, nil
}

func FromListValueProto(lv *pb.List) (*List, error) {
	if lv == nil {
		return nil, nil
	}

	nl := []Value{}
	for _, el := range lv.Fields {
		inner, err := FromProto(el)
		if err != nil {
			return nil, err
		}

		nl = append(nl, inner)
	}
	return &List{Underlying: nl}, nil
}

func fromDecimalValueProto(dec *pb.Decimal) *Decimal {
	if dec == nil {
		return nil
	}

	dc := decimal.NewFromBigInt(protoToBigInt(dec.Coefficient), dec.Exponent)
	return NewDecimal(dc)
}

func protoToBigInt(biv *pb.BigInt) *big.Int {
	if biv == nil {
		return nil
	}

	av := &big.Int{}
	av = av.SetBytes(biv.AbsVal)

	if biv.Sign < 0 {
		av.Neg(av)
	}

	return av
}

func fromBigIntValueProto(biv *pb.BigInt) *BigInt {
	return NewBigInt(protoToBigInt(biv))
}

func createMapFromStruct(v any) (Value, error) {
	var resultMap map[string]interface{}
	err := mapstructure.Decode(v, &resultMap)
	if err != nil {
		return nil, err
	}
	return NewMap(resultMap)
}

func unwrapTo[T any](underlying T, to any) error {
	switch tb := to.(type) {
	case *T:
		if tb == nil {
			return fmt.Errorf("cannot unwrap to nil pointer")
		}
		*tb = underlying
	case *any:
		if tb == nil {
			return fmt.Errorf("cannot unwrap to nil pointer")
		}
		*tb = underlying
	default:
		return fmt.Errorf("cannot unwrap to value of type: %T", to)
	}

	return nil
}
