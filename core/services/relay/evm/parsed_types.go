package evm

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type parsedTypes struct {
	encoderDefs map[string]types.CodecEntry
	decoderDefs map[string]types.CodecEntry
}

func (parsed *parsedTypes) toCodec() (commontypes.RemoteCodec, error) {
	modByTypeName := map[string]codec.Modifier{}
	if err := addEntries(parsed.encoderDefs, modByTypeName); err != nil {
		return nil, err
	}
	if err := addEntries(parsed.decoderDefs, modByTypeName); err != nil {
		return nil, err
	}

	mod, err := codec.NewByItemTypeModifier(modByTypeName)
	if err != nil {
		return nil, err
	}
	underlying := &evmCodec{
		encoder:     &encoder{Definitions: parsed.encoderDefs},
		decoder:     &decoder{Definitions: parsed.decoderDefs},
		parsedTypes: parsed,
	}
	return codec.NewModifierCodec(underlying, mod, evmDecoderHooks...)
}

// addEntries extracts the mods from codecEntry and adds them to modByTypeName use with codec.NewByItemTypeModifier
// Since each input/output can have its own modifications, we need to keep track of them by type name
func addEntries(defs map[string]types.CodecEntry, modByTypeName map[string]codec.Modifier) error {
	for k, def := range defs {
		modByTypeName[k] = def.Modifier()
		_, err := def.Modifier().RetypeToOffChain(reflect.PointerTo(def.CheckedType()), k)
		if err != nil {
			return fmt.Errorf("%w: cannot retype %v: %w", commontypes.ErrInvalidConfig, k, err)
		}
	}
	return nil
}
