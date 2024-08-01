package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

type nilShrink struct {
	done bool
}

func (s *nilShrink) Next() (interface{}, bool) {
	if !s.done {
		s.done = true
		return nil, true
	}
	return nil, false
}

// PtrShrinker convert a value shrinker to a pointer to value shrinker
func PtrShrinker(elementShrinker gopter.Shrinker) gopter.Shrinker {
	return func(v interface{}) gopter.Shrink {
		if v == nil {
			return gopter.NoShrink
		}
		elem := reflect.ValueOf(v).Elem()
		if !elem.IsValid() || !elem.CanInterface() {
			return gopter.NoShrink
		}
		rt := reflect.TypeOf(v)
		elementShink := elementShrinker(reflect.ValueOf(v).Elem().Interface())

		nilShrink := &nilShrink{}
		return gopter.ConcatShrinks(
			nilShrink.Next,
			elementShink.Map(func(elem interface{}) interface{} {
				slice := reflect.MakeSlice(reflect.SliceOf(rt.Elem()), 0, 1)
				slice = reflect.Append(slice, reflect.ValueOf(elem))

				return slice.Index(0).Addr().Interface()
			}),
		)
	}
}
