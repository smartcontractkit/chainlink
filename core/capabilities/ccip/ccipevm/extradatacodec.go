package ccipevm

import (
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

const (
	svmV1DecodeStructName = "decodeSVMExtraArgsStruct"
	evmV1DecodeName       = "decodeEVMExtraArgsV1"
	evmV2DecodeName       = "decodeEVMExtraArgsV2"
	evmDestExecDataKey    = "destGasAmount"
)

// ExtraDataDecoder is a concrete implementation of SourceChainExtraDataCodec
type ExtraDataDecoder struct{}

// DecodeDestExecDataToMap reformats bytes into a chain agnostic map[string]interface{} representation for dest exec data
func (d ExtraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]interface{}, error) {
	destGasAmount, err := abiDecodeUint32(destExecData)
	if err != nil {
		return nil, fmt.Errorf("decode dest gas amount: %w", err)
	}

	return map[string]interface{}{
		evmDestExecDataKey: destGasAmount,
	}, nil
}

// DecodeExtraArgsToMap reformats bytes into a chain agnostic map[string]any representation for extra args
func (d ExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var method string
	var extraByteOffset = 4
	switch string(extraArgs[:4]) {
	case string(evmExtraArgsV1Tag):
		// for EVMExtraArgs, the first four bytes is the method name
		method = evmV1DecodeName
	case string(evmExtraArgsV2Tag):
		method = evmV2DecodeName
	case string(svmExtraArgsV1Tag):
		method = svmV1DecodeStructName
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}

	output := make(map[string]any)
	args := make(map[string]interface{})
	err := messageHasherABI.Methods[method].Inputs.UnpackIntoMap(args, extraArgs[extraByteOffset:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args %v: %w", method, err)
	}

	switch method {
	case evmV1DecodeName, evmV2DecodeName:
		for k, val := range args {
			output[k] = val
		}
	case svmV1DecodeStructName:
		// NOTE: the cast only works with this particular struct definition, including the json tags
		extraArgsStruct, ok := args["extraArgs"].(struct {
			ComputeUnits             uint32      `json:"computeUnits"`
			AccountIsWritableBitmap  uint64      `json:"accountIsWritableBitmap"`
			AllowOutOfOrderExecution bool        `json:"allowOutOfOrderExecution"`
			TokenReceiver            [32]uint8   `json:"tokenReceiver"`
			Accounts                 [][32]uint8 `json:"accounts"`
		})
		if !ok {
			return nil, fmt.Errorf("solana extra args struct is not the equivalent of message_hasher.ClientSVMExtraArgsV1")
		}
		output["computeUnits"] = extraArgsStruct.ComputeUnits
		output["accountIsWritableBitmap"] = extraArgsStruct.AccountIsWritableBitmap
		output["allowOutOfOrderExecution"] = extraArgsStruct.AllowOutOfOrderExecution
		output["tokenReceiver"] = extraArgsStruct.TokenReceiver
		output["accounts"] = extraArgsStruct.Accounts
	default:
		return nil, fmt.Errorf("unknown extra args method: %s", method)
	}

	return output, nil
}

// Ensure ExtraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &ExtraDataDecoder{}
