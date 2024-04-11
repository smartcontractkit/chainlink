package values

import (
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Map struct {
	Underlying map[string]Value
}

func EmptyMap() *Map {
	return &Map{
		Underlying: map[string]Value{},
	}
}

func NewMap(m map[string]any) (*Map, error) {
	mv := map[string]Value{}
	for k, v := range m {
		val, err := Wrap(v)
		if err != nil {
			return nil, err
		}

		mv[k] = val
	}

	return &Map{
		Underlying: mv,
	}, nil
}

func (m *Map) proto() *pb.Value {
	pm := map[string]*pb.Value{}
	for k, v := range m.Underlying {
		pm[k] = Proto(v)
	}

	return pb.NewMapValue(pm)
}

func (m *Map) Unwrap() (any, error) {
	nm := map[string]any{}
	return nm, m.UnwrapTo(&nm)
}

func (m *Map) UnwrapTo(to any) error {
	c := &mapstructure.DecoderConfig{
		Result: to,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapValueToMap,
			unwrapsValues,
		),
	}

	d, err := mapstructure.NewDecoder(c)
	if err != nil {
		return err
	}

	return d.Decode(m.Underlying)
}

func mapValueToMap(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f != reflect.TypeOf(map[string]Value{}) {
		return data, nil
	}

	switch t {
	// If the destination type is `map[string]any` or `any`,
	// fully unwrap the values.Map.
	// We have to handle the `any` case here as otherwise UnwrapTo won't work on
	// maps recursively
	case reflect.TypeOf(map[string]any{}), reflect.TypeOf((*any)(nil)).Elem():
		dv := data.(map[string]Value)
		d := map[string]any{}
		for k, v := range dv {
			unw, err := Unwrap(v)
			if err != nil {
				return nil, err
			}

			d[k] = unw
		}

		return d, nil
	}
	return data, nil
}

func baseTypesEqual(f reflect.Type, t reflect.Type) bool {
	if f.Kind() == reflect.Pointer {
		f = f.Elem()
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return f == t
}

// unwrapsValues takes a value of type `f` and tries to convert it to a value of type `t`
func unwrapsValues(f reflect.Type, t reflect.Type, data any) (any, error) {
	// First, check if t and f have the same base type. If they do,
	// we don't need to do anything further. `mapstructure` will
	// automatically convert between values and pointers to get the right result.
	if baseTypesEqual(f, t) {
		return data, nil
	}

	valueType := reflect.TypeOf((*Value)(nil)).Elem()

	// Next, if f is a `Value`, we'll try to transform it to `t`,
	// but only if `t` is not itself a `Value`.
	// This avoids the following cases which we handle differently:
	// - f and t are the same concrete value type -- handled above.
	// - data is a concrete value and t represents the `Value` interface type.
	//   This is compatible and we'll handle it on line 137 by returning `data`.
	// - f and t are different concrete value types -- we can't handle that
	//   here, so we'll just return data untransformed.
	// In all other cases, we want to rely on data's UnwrapTo implementation
	// to try to get the right result.
	if f.Implements(valueType) && !t.Implements(valueType) {
		dv := data.(Value)

		n := reflect.New(t).Interface()
		err := dv.UnwrapTo(n)
		if err != nil {
			return nil, err
		}

		if reflect.TypeOf(n).Elem() == t {
			return n, nil
		}
	}

	return data, nil
}
