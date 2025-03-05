package ccipevm

import (
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExtraDataDecoder is a concrete implementation of ExtraDataDecoder
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
	var extraByteOffset int
	switch string(extraArgs[:4]) {
	case string(evmExtraArgsV1Tag):
		// for EVMExtraArgs, the first four bytes is the method name
		method = evmV1DecodeName
		extraByteOffset = 4
	case string(evmExtraArgsV2Tag):
		method = evmV2DecodeName
		extraByteOffset = 4
	case string(svmExtraArgsV1Tag):
		// for SVMExtraArgs there's the four bytes plus another 32 bytes padding for the dynamic array
		// TODO this is a temporary solution, the evm on-chain side will fix it, so the offset should just be 4 instead of 36
		// https://smartcontract-it.atlassian.net/browse/BCI-4180
		method = svmV1DecodeName
		extraByteOffset = 36
	default:
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}

	output := make(map[string]any)
	args := make(map[string]interface{})
	err := messageHasherABI.Methods[method].Inputs.UnpackIntoMap(args, extraArgs[extraByteOffset:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args %v: %w", method, err)
	}

	for k, val := range args {
		output[k] = val
	}
	return output, nil
}
