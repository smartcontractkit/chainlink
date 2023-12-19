package evm

import (
	"context"
	"reflect"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type encoder struct {
	Definitions map[string]*codecEntry
	lggr        logger.Logger
}

var _ commontypes.Encoder = &encoder{}

func (e *encoder) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	e.lggr.Infof("!!!!!!!!!!\nEncode: %#v\n%s\n!!!!!!!!!!\n", item, itemType)
	info, ok := e.Definitions[itemType]
	if !ok {
		e.lggr.Errorf("!!!!!!!!!!\nEncode error not found\n%s\n!!!!!!!!!!\n", itemType)
		return nil, commontypes.ErrInvalidType
	}

	if item == nil {
		cpy := make([]byte, len(info.encodingPrefix))
		copy(cpy, info.encodingPrefix)
		return cpy, nil
	}

	b, err := encode(reflect.ValueOf(item), info, e.lggr)
	if err == nil {
		e.lggr.Infof("!!!!!!!!!!\nEncode success\n%s\n!!!!!!!!!!\n", itemType)
	} else {
		e.lggr.Errorf("!!!!!!!!!!\nEncode error\n%v\n%s\n!!!!!!!!!!\n", err, itemType)
	}
	return b, err
}

func (e *encoder) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	e.lggr.Infof("!!!!!!!!!!\nGetMaxEncodingSize\n%s\n\n%v\n!!!!!!!!!!\n", itemType, e.Definitions)
	return e.Definitions[itemType].GetMaxSize(n)
}

func encode(item reflect.Value, info *codecEntry, lggr logger.Logger) ([]byte, error) {
	for item.Kind() == reflect.Pointer {
		item = reflect.Indirect(item)
	}
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		return encodeArray(item, info, lggr)
	case reflect.Struct, reflect.Map:
		return encodeItem(item, info, lggr)
	default:
		return nil, commontypes.ErrInvalidEncoding
	}
}

func encodeArray(item reflect.Value, info *codecEntry, lggr logger.Logger) ([]byte, error) {
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
		if err := mapstructureDecode(item.Index(i).Interface(), tmp.Interface(), lggr); err != nil {
			return nil, err
		}
		native.Index(i).Set(reflect.NewAt(nativeElm, tmp.UnsafePointer()).Elem())
	}

	return pack(info, native.Interface())
}

func encodeItem(item reflect.Value, info *codecEntry, lggr logger.Logger) ([]byte, error) {
	if item.Type() == reflect.PointerTo(info.checkedType) {
		item = reflect.NewAt(info.nativeType, item.UnsafePointer())
	} else if item.Type() != reflect.PointerTo(info.nativeType) {
		checked := reflect.New(info.checkedType)
		if err := mapstructureDecode(item.Interface(), checked.Interface(), lggr); err != nil {
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
