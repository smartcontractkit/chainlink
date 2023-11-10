package evm

import (
	"context"

	"github.com/smartcontractkit/chainlink-relay/pkg/codec"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type mapDecoder struct {
	Definitions map[string]*CodecEntry
}

var _ relaytypes.MapDecoder = &mapDecoder{}

func (m *mapDecoder) DecodeSingle(ctx context.Context, raw []byte, itemType string) (map[string]any, error) {
	info, ok := m.Definitions[itemType]
	if !ok {
		return nil, relaytypes.InvalidTypeError{}
	}

	values := map[string]any{}
	if err := info.Args.UnpackIntoMap(values, raw); err != nil {
		return nil, relaytypes.InvalidEncodingError{}
	}

	if noName, ok := values[""]; ok {
		values = map[string]any{}
		if err := mapstructureDecode(noName, &values); err != nil {
			return nil, relaytypes.InvalidEncodingError{}
		}
	}

	args := info.unwrappedArgs
	fields := make([]string, len(args))
	for i, arg := range args {
		fields[i] = arg.Name
	}
	return values, codec.VerifyFieldMaps(fields, values)
}

func (m *mapDecoder) DecodeMany(ctx context.Context, raw []byte, itemType string) ([]map[string]any, error) {
	decoded, err := m.DecodeSingle(ctx, raw, itemType)
	if err != nil {
		return nil, err
	}
	return codec.SplitValueFields(decoded)
}

func (m *mapDecoder) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return GetMaxSizeFormEntry(n, m.Definitions[itemType])
}
