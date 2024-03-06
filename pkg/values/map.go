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
	for k, v := range m.Underlying {
		uv, err := Unwrap(v)
		if err != nil {
			return nil, err
		}

		nm[k] = uv
	}

	return nm, nil
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
	case reflect.TypeOf(map[string]any{}):
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

		// Handle ints exceptionally;
		// This is because ints are handled as int64s
		// in the values library.
		case reflect.TypeOf(int(0)):
			var i int
			err := dv.UnwrapTo(&i)
			if err != nil {
				return nil, err
			}

			return i, nil
		}
	}

	return data, nil
}
