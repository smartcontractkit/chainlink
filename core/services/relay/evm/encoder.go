package evm

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type encoder struct {
	Definitions map[string]*codecEntry
}

var _ commontypes.Encoder = &encoder{}

func (e *encoder) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	fmt.Printf("!!!!!!!!!!\nEncode: %#v\n%s\n!!!!!!!!!!\n", item, itemType)
	info, ok := e.Definitions[itemType]
	if !ok {
		fmt.Printf("!!!!!!!!!!\nEncode error not found\n%s\n!!!!!!!!!!\n", itemType)
		return nil, commontypes.ErrInvalidType
	}

	if item == nil {
		cpy := make([]byte, len(info.encodingPrefix))
		copy(cpy, info.encodingPrefix)
		return cpy, nil
	}

	b, err := encode(reflect.ValueOf(item), info)
	if err == nil {
		fmt.Printf("!!!!!!!!!!\nEncode success\n%s\n!!!!!!!!!!\n", itemType)
	} else {
		fmt.Printf("!!!!!!!!!!\nEncode error\n%v\n%s\n!!!!!!!!!!\n", err, itemType)
	}
	return b, err
}

func (e *encoder) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	b := make([]byte, 2048) // adjust buffer size to be larger than expected stack
	nb := runtime.Stack(b, false)
	s := string(b[:nb])
	fmt.Printf("!!!!!!!!!!\nGetMaxEncodingSize\n%s\n\n%v\n%s!!!!!!!!!!\n", itemType, e.Definitions, s)
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
		return nil, commontypes.ErrInvalidEncoding
	}
}

func encodeArray(item reflect.Value, info *codecEntry) ([]byte, error) {
	length := item.Len()
	var native reflect.Value
	switch info.checkedType.Kind() {
	case reflect.Array:
		if info.checkedType.Len() != length {
			return nil, commontypes.ErrWrongNumberOfElements
		}
		native = reflect.New(info.nativeType).Elem()
	case reflect.Slice:
		native = reflect.MakeSlice(info.nativeType, length, length)
	default:
		return nil, commontypes.ErrInvalidType
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
	if bytes, err := info.Args.Pack(values...); err == nil {
		withPrefix := make([]byte, 0, len(info.encodingPrefix)+len(bytes))
		withPrefix = append(withPrefix, info.encodingPrefix...)
		return append(withPrefix, bytes...), nil
	}

	return nil, commontypes.ErrInvalidType
}
