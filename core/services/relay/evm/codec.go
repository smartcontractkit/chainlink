package evm

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewCodec(conf types.CodecConfig) (commontypes.Codec, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]*CodecEntry{},
		decoderDefs: map[string]*CodecEntry{},
	}

	for k, v := range conf.ChainCodecConfigs {
		args := abi.Arguments{}
		if err := json.Unmarshal(([]byte)(v.TypeAbi), &args); err != nil {
			return nil, err
		}

		item := &CodecEntry{Args: args}
		if err := item.Init(); err != nil {
			return nil, err
		}

		parsed.encoderDefs[k] = item
		parsed.decoderDefs[k] = item
	}

	return codecFromTypes(parsed), nil
}

func codecFromTypes(parsed *parsedTypes) *evmCodec {
	return &evmCodec{
		encoder: &encoder{Definitions: parsed.encoderDefs},
		decoder: &decoder{Definitions: parsed.decoderDefs},
		types:   parsed,
	}
}

var _ commontypes.TypeProvider = &evmCodec{}

type evmCodec struct {
	*encoder
	*decoder
	types *parsedTypes
}

type parsedTypes struct {
	encoderDefs map[string]*CodecEntry
	decoderDefs map[string]*CodecEntry
}

func (c *evmCodec) CreateType(itemType string, forEncoding bool) (any, error) {
	var itemTypes map[string]*CodecEntry
	if forEncoding {
		itemTypes = c.types.encoderDefs
	} else {
		itemTypes = c.types.decoderDefs
	}

	def, ok := itemTypes[itemType]
	if !ok {
		return nil, commontypes.ErrInvalidType
	}

	return def.checkedType, nil
}
