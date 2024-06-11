package evm

import (
	"context"
	"encoding/binary"
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
	var err error
	var result []byte

	// 1. Version (1 byte)
	if meta.Version > 255 {
		return nil, fmt.Errorf("version must be between 0 and 255")
	}
	result = append(result, byte(meta.Version))

	// 2. Execution ID (32 bytes)
	if result, err = decodeAndAppend(meta.ExecutionID, 32, result, "ExecutionID"); err != nil {
		return nil, err
	}

	// 3. Timestamp (4 bytes)
	tsBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tsBytes, meta.Timestamp)
	result = append(result, tsBytes...)

	// 4. DON ID (4 bytes)
	if result, err = decodeAndAppend(meta.DONID, 4, result, "DONID"); err != nil {
		return nil, err
	}

	// 5. DON config version (4 bytes)
	cfgVersionBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(cfgVersionBytes, meta.DONConfigVersion)
	result = append(result, cfgVersionBytes...)

	// 6. Workflow ID / spec hash (32 bytes)
	if result, err = decodeAndAppend(meta.WorkflowID, 32, result, "WorkflowID"); err != nil {
		return nil, err
	}

	// 7. Workflow Name (10 bytes)
	if result, err = decodeAndAppend(meta.WorkflowName, 10, result, "WorkflowName"); err != nil {
		return nil, err
	}

	// 8. Workflow Owner (20 bytes)
	if result, err = decodeAndAppend(meta.WorkflowOwner, 20, result, "WorkflowOwner"); err != nil {
		return nil, err
	}

	// 9. Report ID (2 bytes)
	if result, err = decodeAndAppend(meta.ReportID, 2, result, "ReportID"); err != nil {
		return nil, err
	}

	return append(result, userPayload...), nil
}

func decodeAndAppend(id string, expectedLen int, prevResult []byte, logName string) ([]byte, error) {
	b, err := hex.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("failed to hex-decode %s (%s): %w", logName, id, err)
	}
	if len(b) != expectedLen {
		return nil, fmt.Errorf("incorrect length for id %s (%s), expected %d bytes, got %d", logName, id, expectedLen, len(b))
	}
	return append(prevResult, b...), nil
}
