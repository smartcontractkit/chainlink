package evm

import (
	"context"
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
	selector, err := abiutil.ParseSignature("inner(" + selectorStr + ")")
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
	// prepend workflowID and workflowExecutionID to the encoded user data
	workflowIDbytes, executionIDBytes, err := extractIDs(unwrappedMap)
	if err != nil {
		return nil, err
	}
	return append(append(workflowIDbytes, executionIDBytes...), userPayload...), nil
}

// extract workflowID and executionID from the input map, validate and align to 32 bytes
// NOTE: consider requiring them to be exactly 32 bytes to avoid issues with padding
func extractIDs(input map[string]any) ([]byte, []byte, error) {
	workflowID, ok := input[consensustypes.WorkflowIDFieldName].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected %s to be a string", consensustypes.WorkflowIDFieldName)
	}
	executionID, ok := input[consensustypes.ExecutionIDFieldName].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected %s to be a string", consensustypes.ExecutionIDFieldName)
	}
	if len(workflowID) > 32 || len(executionID) > 32 {
		return nil, nil, fmt.Errorf("IDs too long: %d, %d", len(workflowID), len(executionID))
	}
	alignedWorkflowID := make([]byte, idLen)
	copy(alignedWorkflowID, workflowID)
	alignedExecutionID := make([]byte, idLen)
	copy(alignedExecutionID, executionID)
	return alignedWorkflowID, alignedExecutionID, nil
}
