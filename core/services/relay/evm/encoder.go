package evm

import (
	"context"
	"fmt"
	"reflect"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type encoder struct {
	Definitions map[string]*codecEntry
}

var _ commontypes.Encoder = &encoder{}

func (e *encoder) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	info, ok := e.Definitions[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find definition for %s", commontypes.ErrInvalidType, itemType)
	}

	if item == nil {
		cpy := make([]byte, len(info.encodingPrefix))
		copy(cpy, info.encodingPrefix)
		return cpy, nil
	}

	return encode(reflect.ValueOf(item), info)
}

func (e *encoder) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return e.Definitions[itemType].GetMaxSize(n)
}

func encode(item reflect.Value, info *codecEntry) ([]byte, error) {
	for item.Kind() == reflect.Pointer {
		item = reflect.Indirect(item)
	}
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		return encodeArray(item, info)
	case reflect.Struct, reflect.Map:
		return encodeItem(item, info)
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}
}

func encodeArray(item reflect.Value, info *codecEntry) ([]byte, error) {
	length := item.Len()
	var native reflect.Value
	switch info.checkedType.Kind() {
	case reflect.Array:
		if info.checkedType.Len() != length {
			return nil, commontypes.ErrSliceWrongLen
		}
		native = reflect.New(info.nativeType).Elem()
	case reflect.Slice:
		native = reflect.MakeSlice(info.nativeType, length, length)
	default:
		return nil, fmt.Errorf("%w: cannot encode %v as array", commontypes.ErrInvalidType, info.checkedType.Kind())
	}

	checkedElm := info.checkedType.Elem()
	nativeElm := info.nativeType.Elem()
	for i := 0; i < length; i++ {
		tmp := reflect.New(checkedElm)
		if err := mapstructureDecode(item.Index(i).Interface(), tmp.Interface()); err != nil {
			return nil, err
		}
		native.Index(i).Set(reflect.NewAt(nativeElm, tmp.UnsafePointer()).Elem())
	}

	return pack(info, native.Interface())
}

func encodeItem(item reflect.Value, info *codecEntry) ([]byte, error) {
	if item.Type() == reflect.PointerTo(info.checkedType) {
		item = reflect.NewAt(info.nativeType, item.UnsafePointer())
	} else if item.Type() != reflect.PointerTo(info.nativeType) {
		checked := reflect.New(info.checkedType)
		if err := mapstructureDecode(item.Interface(), checked.Interface()); err != nil {
			return nil, err
		}
		item = reflect.NewAt(info.nativeType, checked.UnsafePointer())
	}

	item = reflect.Indirect(item)
	length := item.NumField()
	values := make([]any, length)
	iType := item.Type()
	for i := 0; i < length; i++ {
		if iType.Field(i).IsExported() {
			values[i] = item.Field(i).Interface()
		}
	}

	return pack(info, values...)
}

func pack(info *codecEntry, values ...any) ([]byte, error) {
	bytes, err := info.Args.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	withPrefix := make([]byte, 0, len(info.encodingPrefix)+len(bytes))
	withPrefix = append(withPrefix, info.encodingPrefix...)
	return append(withPrefix, bytes...), nil
}
