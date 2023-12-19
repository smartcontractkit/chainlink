package evm

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type decoder struct {
	Definitions map[string]*codecEntry
	lggr        logger.Logger
}

var _ commontypes.Decoder = &decoder{}

func (m *decoder) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	m.lggr.Infof("!!!!!!!!!!\nDecode\n%s\n!!!!!!!!!!\n", itemType)
	info, ok := m.Definitions[itemType]
	if !ok {
		m.lggr.Errorf("!!!!!!!!!!\nDecode err not found type\n%s\n!!!!!!!!!!\n", itemType)
		return commontypes.ErrInvalidType
	}

	decode, err := extractDecoding(info, raw)
	if err != nil {
		m.lggr.Errorf("!!!!!!!!!!\nDecode err: %v\n%s\n!!!!!!!!!!\n", err, itemType)
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
		return setElements(length, rDecode, iInto, m.lggr)
	case reflect.Slice:
		iInto := reflect.Indirect(reflect.ValueOf(into))
		length := rDecode.Len()
		iInto.Set(reflect.MakeSlice(iInto.Type(), length, length))
		return setElements(length, rDecode, iInto, m.lggr)
	default:
		return mapstructureDecode(decode, into, m.lggr)
	}
}

func (m *decoder) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	m.lggr.Infof("!!!!!!!!!!\nGetMaxDecodingSize\n%s\n\n%v\n!!!!!!!!!!\n", itemType, m.Definitions)

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

func setElements(length int, rDecode reflect.Value, iInto reflect.Value, lggr logger.Logger) error {
	for i := 0; i < length; i++ {
		if err := mapstructureDecode(rDecode.Index(i).Interface(), iInto.Index(i).Addr().Interface(), lggr); err != nil {
			return err
		}
	}

	return nil
}

func mapstructureDecode(src, dest any, lggr logger.Logger) error {
	mDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(evmDecoderHooks...),
		Result:     dest,
		Squash:     true,
	})
	if err != nil || mDecoder.Decode(src) != nil {
		lggr.Errorf("!!!!!!!!!!\nDecode item error: %v\n%v\n!!!!!!!!!!\n", err, mDecoder.Decode(src))
		return commontypes.ErrInvalidType
	}
	lggr.Infof("!!!!!!!!!!\nDecode item success\n%#v\n%#v\n!!!!!!!!!!\n", dest, src)
	return nil
}
