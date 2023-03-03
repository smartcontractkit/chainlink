package gen

import (
	"fmt"
	"reflect"

	"github.com/leanovate/gopter"
)

type mapShrinkOne struct {
	original         reflect.Value
	key              reflect.Value
	keyShrink        gopter.Shrink
	elementShrink    gopter.Shrink
	state            bool
	keyExhausted     bool
	lastKey          interface{}
	elementExhausted bool
	lastElement      interface{}
}

func (s *mapShrinkOne) nextKeyValue() (interface{}, interface{}, bool) {
	for !s.keyExhausted && !s.elementExhausted {
		s.state = !s.state
		if s.state && !s.keyExhausted {
			value, ok := s.keyShrink()
			if ok {
				s.lastKey = value
				return s.lastKey, s.lastElement, true
			}
			s.keyExhausted = true
		} else if !s.state && !s.elementExhausted {
			value, ok := s.elementShrink()
			if ok {
				s.lastElement = value
				return s.lastKey, s.lastElement, true
			}
			s.elementExhausted = true
		}
	}
	return nil, nil, false
}

func (s *mapShrinkOne) Next() (interface{}, bool) {
	nextKey, nextValue, ok := s.nextKeyValue()
	if !ok {
		return nil, false
	}
	result := reflect.MakeMapWithSize(s.original.Type(), s.original.Len())
	for _, key := range s.original.MapKeys() {
		if !reflect.DeepEqual(key.Interface(), s.key.Interface()) {
			result.SetMapIndex(key, s.original.MapIndex(key))
		}
	}
	result.SetMapIndex(reflect.ValueOf(nextKey), reflect.ValueOf(nextValue))

	return result.Interface(), true
}

// MapShrinkerOne creates a map shrinker from a shrinker for the key values of a map.
// The length of the map will remain (mostly) unchanged, instead each key value pair is
// shrunk after the other.
func MapShrinkerOne(keyShrinker, elementShrinker gopter.Shrinker) gopter.Shrinker {
	return func(v interface{}) gopter.Shrink {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Map {
			panic(fmt.Sprintf("%#v is not a map", v))
		}

		keys := rv.MapKeys()
		shrinks := make([]gopter.Shrink, 0, len(keys))
		for _, key := range keys {
			mapShrinkOne := &mapShrinkOne{
				original:      rv,
				key:           key,
				keyShrink:     keyShrinker(key.Interface()),
				lastKey:       key.Interface(),
				elementShrink: elementShrinker(rv.MapIndex(key).Interface()),
				lastElement:   rv.MapIndex(key).Interface(),
			}
			shrinks = append(shrinks, mapShrinkOne.Next)
		}
		return gopter.ConcatShrinks(shrinks...)
	}
}

type mapShrink struct {
	original     reflect.Value
	originalKeys []reflect.Value
	length       int
	offset       int
	chunkLength  int
}

func (s *mapShrink) Next() (interface{}, bool) {
	if s.chunkLength == 0 {
		return nil, false
	}
	keys := make([]reflect.Value, 0, s.length-s.chunkLength)
	keys = append(keys, s.originalKeys[0:s.offset]...)
	s.offset += s.chunkLength
	if s.offset < s.length {
		keys = append(keys, s.originalKeys[s.offset:s.length]...)
	} else {
		s.offset = 0
		s.chunkLength >>= 1
	}

	result := reflect.MakeMapWithSize(s.original.Type(), len(keys))
	for _, key := range keys {
		result.SetMapIndex(key, s.original.MapIndex(key))
	}

	return result.Interface(), true
}

// MapShrinker creates a map shrinker from shrinker for the key values.
// The length of the map will be shrunk as well
func MapShrinker(keyShrinker, elementShrinker gopter.Shrinker) gopter.Shrinker {
	return func(v interface{}) gopter.Shrink {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Map {
			panic(fmt.Sprintf("%#v is not a Map", v))
		}
		keys := rv.MapKeys()
		mapShrink := &mapShrink{
			original:     rv,
			originalKeys: keys,
			offset:       0,
			length:       rv.Len(),
			chunkLength:  rv.Len() >> 1,
		}

		shrinks := make([]gopter.Shrink, 0, rv.Len()+1)
		shrinks = append(shrinks, mapShrink.Next)
		for _, key := range keys {
			mapShrinkOne := &mapShrinkOne{
				original:      rv,
				key:           key,
				keyShrink:     keyShrinker(key.Interface()),
				lastKey:       key.Interface(),
				elementShrink: elementShrinker(rv.MapIndex(key).Interface()),
				lastElement:   rv.MapIndex(key).Interface(),
			}
			shrinks = append(shrinks, mapShrinkOne.Next)
		}
		return gopter.ConcatShrinks(shrinks...)
	}
}
