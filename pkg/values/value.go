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
	case *pb.Value_BigintValue:
		return fromBigIntValueProto(val.GetBigintValue())
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

func fromBigIntValueProto(b []byte) *BigInt {
	i := big.Int{}
	bi := i.SetBytes(b)
	return NewBigInt(bi)
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
