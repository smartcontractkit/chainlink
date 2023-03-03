package gen

import (
	"fmt"
	"reflect"

	"github.com/leanovate/gopter"
)

type sliceShrinkOne struct {
	original      reflect.Value
	index         int
	elementShrink gopter.Shrink
}

func (s *sliceShrinkOne) Next() (interface{}, bool) {
	value, ok := s.elementShrink()
	if !ok {
		return nil, false
	}
	result := reflect.MakeSlice(s.original.Type(), s.original.Len(), s.original.Len())
	reflect.Copy(result, s.original)
	if value == nil {
		result.Index(s.index).Set(reflect.Zero(s.original.Type().Elem()))
	} else {
		result.Index(s.index).Set(reflect.ValueOf(value))
	}

	return result.Interface(), true
}

// SliceShrinkerOne creates a slice shrinker from a shrinker for the elements of the slice.
// The length of the slice will remains unchanged, instead each element is shrunk after the
// other.
func SliceShrinkerOne(elementShrinker gopter.Shrinker) gopter.Shrinker {
	return func(v interface{}) gopter.Shrink {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			panic(fmt.Sprintf("%#v is not a slice", v))
		}

		shrinks := make([]gopter.Shrink, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			sliceShrinkOne := &sliceShrinkOne{
				original:      rv,
				index:         i,
				elementShrink: elementShrinker(rv.Index(i).Interface()),
			}
			shrinks = append(shrinks, sliceShrinkOne.Next)
		}
		return gopter.ConcatShrinks(shrinks...)
	}
}

type sliceShrink struct {
	original    reflect.Value
	length      int
	offset      int
	chunkLength int
}

func (s *sliceShrink) Next() (interface{}, bool) {
	if s.chunkLength == 0 {
		return nil, false
	}
	value := reflect.AppendSlice(reflect.MakeSlice(s.original.Type(), 0, s.length-s.chunkLength), s.original.Slice(0, s.offset))
	s.offset += s.chunkLength
	if s.offset < s.length {
		value = reflect.AppendSlice(value, s.original.Slice(s.offset, s.length))
	} else {
		s.offset = 0
		s.chunkLength >>= 1
	}

	return value.Interface(), true
}

// SliceShrinker creates a slice shrinker from a shrinker for the elements of the slice.
// The length of the slice will be shrunk as well
func SliceShrinker(elementShrinker gopter.Shrinker) gopter.Shrinker {
	return func(v interface{}) gopter.Shrink {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			panic(fmt.Sprintf("%#v is not a slice", v))
		}
		sliceShrink := &sliceShrink{
			original:    rv,
			offset:      0,
			length:      rv.Len(),
			chunkLength: rv.Len() >> 1,
		}

		shrinks := make([]gopter.Shrink, 0, rv.Len()+1)
		shrinks = append(shrinks, sliceShrink.Next)
		for i := 0; i < rv.Len(); i++ {
			sliceShrinkOne := &sliceShrinkOne{
				original:      rv,
				index:         i,
				elementShrink: elementShrinker(rv.Index(i).Interface()),
			}
			shrinks = append(shrinks, sliceShrinkOne.Next)
		}
		return gopter.ConcatShrinks(shrinks...)
	}
}
