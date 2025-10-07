package evm

import (
	"context"
	"encoding/json"
	"fmt"

	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/abi"
	"github.com/smartcontractkit/chainlink-evm/pkg/codec"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

const (
	abiConfigFieldName    = "abi"
	subabiConfigFieldName = "subabi"
	encoderName           = "user"
)

type capEncoder struct {
	codec commontypes.RemoteCodec
}

var _ consensustypes.Encoder = (*capEncoder)(nil)

func NewEVMEncoder(m *values.Map) (consensustypes.Encoder, error) {
	// parse the "inner" encoder config - user-defined fields
	abiConfig, ok := m.Underlying[abiConfigFieldName]
	if !ok {
		return nil, fmt.Errorf("required field %s is missing", abiConfigFieldName)
	}
	wrappedSelector, err := abiConfig.Unwrap()
	if err != nil {
		return nil, err
	}
	selectorStr, ok := wrappedSelector.(string)
	if !ok {
		return nil, fmt.Errorf("expected %s to be a string", abiConfigFieldName)
	}
	selector, err := abi.ParseSelector("inner(" + selectorStr + ")")
	if err != nil {
		return nil, err
	}
	jsonSelector, err := json.Marshal(selector.Inputs)
	if err != nil {
		return nil, err
	}

	chainCodecConfig := config.ChainCodecConfig{
		TypeABI: string(jsonSelector),
	}

	var subabi map[string]string
	subabiConfig, ok := m.Underlying[subabiConfigFieldName]
	if ok {
		err2 := subabiConfig.UnwrapTo(&subabi)
		if err2 != nil {
			return nil, err2
		}
		codecs, err2 := makePreCodecModifierCodecs(subabi)
		if err2 != nil {
			return nil, err2
		}
		chainCodecConfig.ModifierConfigs = commoncodec.ModifiersConfig{
			&commoncodec.PreCodecModifierConfig{
				Fields: subabi,
				Codecs: codecs,
			},
		}
	}

	codecConfig := config.CodecConfig{Configs: map[string]config.ChainCodecConfig{
		encoderName: chainCodecConfig,
	}}

	c, err := codec.NewCodec(codecConfig)
	if err != nil {
		return nil, err
	}

	return &capEncoder{codec: c}, nil
}

func makePreCodecModifierCodecs(subabi map[string]string) (map[string]commontypes.RemoteCodec, error) {
	codecs := map[string]commontypes.RemoteCodec{}
	for _, abiStr := range subabi {
		selector, err := abi.ParseSelector("inner(" + abiStr + ")")
		if err != nil {
			return nil, err
		}
		jsonSelector, err := json.Marshal(selector.Inputs)
		if err != nil {
			return nil, err
		}
		emptyName := ""
		codecConfig := config.CodecConfig{Configs: map[string]config.ChainCodecConfig{
			emptyName: {
				TypeABI: string(jsonSelector),
			},
		}}
		codec, err := codec.NewCodec(codecConfig)
		if err != nil {
			return nil, err
		}
		codecs[abiStr] = codec
	}
	return codecs, nil
}

func (c *capEncoder) Encode(ctx context.Context, input values.Map) ([]byte, error) {
	unwrappedInput, err := input.Unwrap()
	if err != nil {
		return nil, err
	}
	unwrappedMap, ok := unwrappedInput.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected unwrapped input to be a map")
	}
	userPayload, err := c.codec.Encode(ctx, unwrappedMap, encoderName)
	if err != nil {
		return nil, err
	}

	metaMap, ok := input.Underlying[consensustypes.MetadataFieldName]
	if !ok {
		return nil, fmt.Errorf("expected metadata field to be present: %s", consensustypes.MetadataFieldName)
	}

	var meta consensustypes.Metadata
	err = metaMap.UnwrapTo(&meta)
	if err != nil {
		return nil, err
	}

	encodedMeta, err := meta.Encode()
	if err != nil {
		return nil, err
	}
	return append(encodedMeta, userPayload...), nil
}
