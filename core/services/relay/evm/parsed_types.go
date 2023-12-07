package evm

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type parsedTypes struct {
	encoderDefs map[string]*CodecEntry
	decoderDefs map[string]*CodecEntry
}

func (parsed *parsedTypes) toCodec() (commontypes.RemoteCodec, error) {
	modByType := map[string]codec.Modifier{}
	if err := addEntries(parsed.encoderDefs, modByType); err != nil {
		return nil, err
	}
	if err := addEntries(parsed.decoderDefs, modByType); err != nil {
		return nil, err
	}

	mod, err := codec.NewByItemTypeModifier(modByType)
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

func addEntries(defs map[string]*CodecEntry, modByType map[string]codec.Modifier) error {
	for k, def := range defs {
		modByType[k] = def.mod
		_, err := def.mod.RetypeForOffChain(reflect.PointerTo(def.checkedType), k)
		if err != nil {
			return fmt.Errorf("%w: cannot retype %v: %w", commontypes.ErrInvalidConfig, k, err)
		}
	}
	return nil
}
