package evm

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/mitchellh/mapstructure"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type decoder struct {
	Definitions map[string]*codecEntry
}

var _ commontypes.Decoder = &decoder{}

func (m *decoder) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	fmt.Printf("!!!!!!!!!!\\nDecode\\n%s\\n!!!!!!!!!!\\n", itemType)
	info, ok := m.Definitions[itemType]
	if !ok {
		fmt.Printf("!!!!!!!!!!\\nDecode err not found type\\n%s\\n!!!!!!!!!!\\n", itemType)
		return commontypes.ErrInvalidType
	}

	decode, err := extractDecoding(info, raw)
	if err != nil {
		fmt.Printf("!!!!!!!!!!\\nDecode err: %v\\n%s\\n!!!!!!!!!!\\n", err, itemType)
		return err
	}

	rDecode := reflect.ValueOf(decode)
	switch rDecode.Kind() {
	case reflect.Array:
		iInto := reflect.Indirect(reflect.ValueOf(into))
		length := rDecode.Len()
		if length != iInto.Len() {
			return commontypes.ErrWrongNumberOfElements
		}
		iInto.Set(reflect.New(iInto.Type()).Elem())
		return setElements(length, rDecode, iInto)
	case reflect.Slice:
		iInto := reflect.Indirect(reflect.ValueOf(into))
		length := rDecode.Len()
		iInto.Set(reflect.MakeSlice(iInto.Type(), length, length))
		return setElements(length, rDecode, iInto)
	default:
		return mapstructureDecode(decode, into)
	}
}

func (m *decoder) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	b := make([]byte, 2048) // adjust buffer size to be larger than expected stack
	nb := runtime.Stack(b, false)
	s := string(b[:nb])
	fmt.Printf("!!!!!!!!!!\\nGetMaxDecodingSize\\n%s\\n\\n%v\\n%s!!!!!!!!!!\\n", itemType, m.Definitions, s)

	return m.Definitions[itemType].GetMaxSize(n)
}

func extractDecoding(info *codecEntry, raw []byte) (any, error) {
	unpacked := map[string]any{}
	if err := info.Args.UnpackIntoMap(unpacked, raw); err != nil {
		return nil, commontypes.ErrInvalidEncoding
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
	if err != nil || mDecoder.Decode(src) != nil {
		fmt.Printf("!!!!!!!!!!\\nDecode item error: %v\\n%v\\n!!!!!!!!!!\\n", err, mDecoder.Decode(src))
		return commontypes.ErrInvalidType
	}
	return nil
}
