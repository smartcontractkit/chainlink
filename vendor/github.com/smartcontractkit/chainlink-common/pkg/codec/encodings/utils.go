package encodings

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func SafeDecode[T interface{}](raw []byte, size int, call func([]byte) T) (T, []byte, error) {
	if len(raw) < size {
		var t T
		return t, nil, fmt.Errorf("%w: not enough bytes to decode type", types.ErrInvalidEncoding)
	}
	return call(raw[:size]), raw[size:], nil
}

func EncodeEach(value reflect.Value, into []byte, tc TypeCodec) ([]byte, error) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, fmt.Errorf("%w: value must be a slice or an array", types.ErrInvalidType)
	}

	numElements := value.Len()
	for i := 0; i < numElements; i++ {
		var err error
		into, err = tc.Encode(value.Index(i).Interface(), into)
		if err != nil {
			return nil, err
		}
	}
	return into, nil
}

func DecodeEach(encoded []byte, into reflect.Value, numElements int, tc TypeCodec) (any, []byte, error) {
	switch into.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, nil, fmt.Errorf("%w: value must be a slice or an array", types.ErrInvalidType)
	}

	if into.Len() < numElements {
		return nil, nil, fmt.Errorf("%w: not enough elements in slice or array", types.ErrSliceWrongLen)
	}

	remaining := encoded
	for i := 0; i < numElements; i++ {
		element, bytes, err := tc.Decode(remaining)
		if err != nil {
			return nil, nil, err
		}
		remaining = bytes

		elm := into.Index(i)
		if !elm.CanSet() {
			return nil, nil, fmt.Errorf("%w: cannot set element %d", types.ErrInternal, i)
		}
		elm.Set(reflect.ValueOf(element))
	}

	return into.Interface(), remaining, nil
}

func IndirectIfPointer(value reflect.Value) (reflect.Value, error) {
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("%w: nil pointer", types.ErrInvalidType)
		}

		value = reflect.Indirect(value)
	}

	return value, nil
}
