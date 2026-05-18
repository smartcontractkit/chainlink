package codec

import (
	"context"
	"fmt"
	"reflect"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type encoder struct {
	Definitions map[string]types.CodecEntry
}

var _ commontypes.Encoder = &encoder{}

func (e *encoder) Encode(_ context.Context, item any, itemType string) (res []byte, err error) {
	// nil values can cause abi.Arguments.Pack to panic.
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = fmt.Errorf("%w: cannot encode type", commontypes.ErrInvalidType)
		}
	}()
	info, ok := e.Definitions[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find definition for %s", commontypes.ErrInvalidType, itemType)
	}

	if len(info.Args()) == 0 {
		return info.EncodingPrefix(), nil
	} else if item == nil {
		return nil, fmt.Errorf("%w: cannot encode nil value for %s", commontypes.ErrInvalidType, itemType)
	}

	return encode(reflect.ValueOf(item), info)
}

func (e *encoder) GetMaxEncodingSize(_ context.Context, n int, itemType string) (int, error) {
	entry, ok := e.Definitions[itemType]
	if !ok {
		return 0, fmt.Errorf("%w: nil entry", commontypes.ErrInvalidType)
	}
	return entry.GetMaxSize(n)
}

func encode(item reflect.Value, info types.CodecEntry) ([]byte, error) {
	for item.Kind() == reflect.Pointer {
		item = reflect.Indirect(item)
	}
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := RepresentArray(item, info)
		if err != nil {
			return nil, err
		}
		return pack(info, native)
	case reflect.Struct, reflect.Map:
		values, err := UnrollItem(item, info)
		if err != nil {
			return nil, err
		}
		return pack(info, values...)
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}
}

func RepresentArray(item reflect.Value, info types.CodecEntry) (any, error) {
	length := item.Len()
	checkedType := info.CheckedType()
	checked := reflect.New(checkedType)
	iChecked := reflect.Indirect(checked)
	switch checkedType.Kind() {
	case reflect.Array:
		if checkedType.Len() != length {
			return nil, commontypes.ErrSliceWrongLen
		}
	case reflect.Slice:
		iChecked.Set(reflect.MakeSlice(checkedType, length, length))
	default:
		return nil, fmt.Errorf("%w: cannot encode %v as array", commontypes.ErrInvalidType, checkedType.Kind())
	}

	checkedElm := checkedType.Elem()
	for i := 0; i < length; i++ {
		tmp := reflect.New(checkedElm)
		if err := MapstructureDecode(item.Index(i).Interface(), tmp.Interface()); err != nil {
			return nil, err
		}
		iChecked.Index(i).Set(tmp.Elem())
	}
	native, err := info.ToNative(checked)
	if err != nil {
		return nil, err
	}

	return native.Elem().Interface(), nil
}

func UnrollItem(item reflect.Value, info types.CodecEntry) ([]any, error) {
	checkedType := info.CheckedType()
	if item.CanAddr() {
		item = item.Addr()
	}

	if item.Type() == reflect.PointerTo(checkedType) {
		var err error
		if item, err = info.ToNative(item); err != nil {
			return nil, err
		}
	} else if !info.IsNativePointer(item.Type()) {
		var err error
		checked := reflect.New(checkedType)
		if err = MapstructureDecode(item.Interface(), checked.Interface()); err != nil {
			return nil, err
		}
		if item, err = info.ToNative(checked); err != nil {
			return nil, err
		}
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
	return values, nil
}

func pack(info types.CodecEntry, values ...any) ([]byte, error) {
	bytes, err := info.Args().Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	withPrefix := info.EncodingPrefix()
	return append(withPrefix, bytes...), nil
}
