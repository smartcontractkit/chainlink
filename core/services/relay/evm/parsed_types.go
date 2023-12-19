package evm

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type parsedTypes struct {
	encoderDefs map[string]*codecEntry
	decoderDefs map[string]*codecEntry
}

func (parsed *parsedTypes) toCodec(lggr logger.Logger) (commontypes.RemoteCodec, error) {
	modByTypeName := map[string]codec.Modifier{}
	if err := addEntries(parsed.encoderDefs, modByTypeName); err != nil {
		lggr.Errorf("!!!!!!!!!!\nto codec add encoder entries err\n%#v\n!!!!!!!!!!\n%", err)
		return nil, err
	}
	if err := addEntries(parsed.decoderDefs, modByTypeName); err != nil {
		lggr.Errorf("!!!!!!!!!!\nto codec add decoder entries err\n%#v\n!!!!!!!!!!\n%", err)
		return nil, err
	}

	mod, err := codec.NewByItemTypeModifier(modByTypeName)
	if err != nil {
		lggr.Errorf("!!!!!!!!!!\nto codec mod by type err\n%#v\n!!!!!!!!!!\n%", err)
		return nil, err
	}
	underlying := &evmCodec{
		encoder:     &encoder{Definitions: parsed.encoderDefs, lggr: lggr},
		decoder:     &decoder{Definitions: parsed.decoderDefs, lggr: lggr},
		parsedTypes: parsed,
	}
	mc, err := codec.NewModifierCodec(underlying, mod, evmDecoderHooks...)
	lggr.Errorf("!!!!!!!!!!\nnow modifier codec: has error?\n%v\n%#v\n!!!!!!!!!!\n%", err != nil, err)
	return mc, err
}

// addEntries extracts the mods from codecEntry and adds them to modByTypeName use with codec.NewByItemTypeModifier
// Since each input/output can have its own modifications, we need to keep track of them by type name
func addEntries(defs map[string]*codecEntry, modByTypeName map[string]codec.Modifier) error {
	for k, def := range defs {
		modByTypeName[k] = def.mod
		_, err := def.mod.RetypeForOffChain(reflect.PointerTo(def.checkedType), k)
		if err != nil {
			return fmt.Errorf("%w: cannot retype %v: %w", commontypes.ErrInvalidConfig, k, err)
		}
	}
	return nil
}
