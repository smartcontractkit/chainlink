package codec

import (
	"fmt"
	"reflect"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ParsedTypes struct {
	EncoderDefs map[string]types.CodecEntry
	DecoderDefs map[string]types.CodecEntry
}

func (parsed *ParsedTypes) ToCodec() (commontypes.RemoteCodec, error) {
	modByTypeName := map[string]commoncodec.Modifier{}
	if err := AddEntries(parsed.EncoderDefs, modByTypeName); err != nil {
		return nil, err
	}
	if err := AddEntries(parsed.DecoderDefs, modByTypeName); err != nil {
		return nil, err
	}

	mod, err := commoncodec.NewByItemTypeModifier(modByTypeName)
	if err != nil {
		return nil, err
	}
	underlying := &evmCodec{
		encoder:     &encoder{Definitions: parsed.EncoderDefs},
		decoder:     &decoder{Definitions: parsed.DecoderDefs},
		ParsedTypes: parsed,
	}
	return commoncodec.NewModifierCodec(underlying, mod, DecoderHooks...)
}

// AddEntries extracts the mods from codecEntry and adds them to modByTypeName use with codec.NewByItemTypeModifier
// Since each input/output can have its own modifications, we need to keep track of them by type name
func AddEntries(defs map[string]types.CodecEntry, modByTypeName map[string]commoncodec.Modifier) error {
	for k, def := range defs {
		modByTypeName[k] = def.Modifier()
		_, err := def.Modifier().RetypeToOffChain(reflect.PointerTo(def.CheckedType()), k)
		if err != nil {
			return fmt.Errorf("%w: cannot retype %v: %w", commontypes.ErrInvalidConfig, k, err)
		}
	}
	return nil
}
