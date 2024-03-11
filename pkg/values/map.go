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

func unwrapsValues(f reflect.Type, t reflect.Type, data any) (any, error) {
	valueType := reflect.TypeOf((*Value)(nil)).Elem()
	if f.Implements(valueType) {
		dv := data.(Value)
		unw, err := Unwrap(dv)
		if err != nil {
			return data, nil
		}

		switch t {
		case reflect.TypeOf(unw):
			return unw, nil

		// Handle integer types exceptionally;
		// This is because ints are handled as int64s
		// in the values library.
		// TODO: refactor this so that we inspect the destination type
		// and just call UnwrapTo using an instantiated pointer of that type.
		case reflect.TypeOf(int(0)):
			var i int
			err := dv.UnwrapTo(&i)
			if err != nil {
				return nil, err
			}

			return i, nil
		case reflect.TypeOf(uint(0)):
			var i uint
			err := dv.UnwrapTo(&i)
			if err != nil {
				return nil, err
			}

			return i, nil
		case reflect.TypeOf(uint64(0)):
			var i uint
			err := dv.UnwrapTo(&i)
			if err != nil {
				return nil, err
			}

			return i, nil
		}
	}

	return data, nil
}
