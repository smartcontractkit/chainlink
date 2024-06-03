package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	abiutil "github.com/smartcontractkit/chainlink/v2/core/chains/evm/abi"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	abiConfigFieldName = "abi"
	encoderName        = "user"
	idLen              = 32
)

type capEncoder struct {
	codec commontypes.RemoteCodec
}

var _ consensustypes.Encoder = (*capEncoder)(nil)

func NewEVMEncoder(config *values.Map) (consensustypes.Encoder, error) {
	// parse the "inner" encoder config - user-defined fields
	wrappedSelector, err := config.Underlying[abiConfigFieldName].Unwrap()
	if err != nil {
		return nil, err
	}
	selectorStr, ok := wrappedSelector.(string)
	if !ok {
		return nil, fmt.Errorf("expected %s to be a string", abiConfigFieldName)
	}
	selector, err := abiutil.ParseSelector("inner(" + selectorStr + ")")
	if err != nil {
		return nil, err
	}
	jsonSelector, err := json.Marshal(selector.Inputs)
	if err != nil {
		return nil, err
	}

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		encoderName: {TypeABI: string(jsonSelector)},
	}}
	c, err := NewCodec(codecConfig)
	if err != nil {
		return nil, err
	}

	return &capEncoder{codec: c}, nil
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

	return prependMetadataFields(meta, userPayload)
}

func prependMetadataFields(meta consensustypes.Metadata, userPayload []byte) ([]byte, error) {
	// TODO: use all 7 fields from Metadata struct
	result := []byte{}
	workflowID, err := decodeID(meta.WorkflowID, idLen)
	if err != nil {
		return nil, err
	}
	result = append(result, workflowID...)

	donID, err := decodeID(meta.DONID, 4)
	if err != nil {
		return nil, err
	}
	result = append(result, donID...)

	executionID, err := decodeID(meta.ExecutionID, idLen)
	if err != nil {
		return nil, err
	}
	result = append(result, executionID...)

	workflowOwner, err := decodeID(meta.WorkflowOwner, 20)
	if err != nil {
		return nil, err
	}
	result = append(result, workflowOwner...)
	return append(result, userPayload...), nil
}

func decodeID(id string, expectedLen int) ([]byte, error) {
	b, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	if len(b) != expectedLen {
		return nil, fmt.Errorf("incorrect length for id %s, expected %d bytes, got %d", id, expectedLen, len(b))
	}
	return b, nil
}
