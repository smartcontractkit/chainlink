package evm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type decoder struct {
	Definitions map[string]types.CodecEntry
}

var _ commontypes.Decoder = &decoder{}

func (m *decoder) Decode(_ context.Context, raw []byte, into any, itemType string) error {
	info, ok := m.Definitions[itemType]
	if !ok {
		return fmt.Errorf("%w: cannot find definition for %s", commontypes.ErrInvalidType, itemType)
	}

	decode, err := extractDecoding(info, raw)
	if err != nil {
		return err
	}

	rDecode := reflect.ValueOf(decode)
	switch rDecode.Kind() {
	case reflect.Array:
		return m.decodeArray(into, rDecode)
	case reflect.Slice:
		iInto := reflect.Indirect(reflect.ValueOf(into))
		length := rDecode.Len()
		iInto.Set(reflect.MakeSlice(iInto.Type(), length, length))
		return setElements(length, rDecode, iInto)
	default:
		return mapstructureDecode(decode, into)
	}
}

func (m *decoder) decodeArray(into any, rDecode reflect.Value) error {
	iInto := reflect.Indirect(reflect.ValueOf(into))
	length := rDecode.Len()
	if length != iInto.Len() {
		return commontypes.ErrSliceWrongLen
	}
	iInto.Set(reflect.New(iInto.Type()).Elem())
	return setElements(length, rDecode, iInto)
}

func (m *decoder) GetMaxDecodingSize(_ context.Context, n int, itemType string) (int, error) {
	entry, ok := m.Definitions[itemType]
	if !ok {
		return 0, fmt.Errorf("%w: nil entry", commontypes.ErrInvalidType)
	}
	return entry.GetMaxSize(n)
}

func extractDecoding(info types.CodecEntry, raw []byte) (any, error) {
	unpacked := map[string]any{}
	args := info.Args()
	if err := args.UnpackIntoMap(unpacked, raw); err != nil {
		return nil, fmt.Errorf("%w: %w: for args %#v", commontypes.ErrInvalidEncoding, err, args)
	}
	var decode any = unpacked

	if noName, ok := unpacked[""]; ok {
		decode = noName
	}
	return decode, nil
}

func setElements(length int, rDecode reflect.Value, iInto reflect.Value) error {
	for i := 0; i < length; i++ {
		if err := mapstructureDecode(rDecode.Index(i).Interface(), iInto.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}

func mapstructureDecode(src, dest any) error {
	mDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(evmDecoderHooks...),
		Result:     dest,
		Squash:     true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	if err = mDecoder.Decode(src); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}
	return nil
}
