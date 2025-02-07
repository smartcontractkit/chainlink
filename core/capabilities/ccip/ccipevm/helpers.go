package ccipevm

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	svmV1DecodeName    = "decodeSVMExtraArgsV1"
	evmV1DecodeName    = "decodeEVMExtraArgsV1"
	evmV2DecodeName    = "decodeEVMExtraArgsV2"
	evmDestExecDataKey = "destGasAmount"
)

var (
	abiUint32               = ABITypeOrPanic("uint32")
	TokenDestGasOverheadABI = abi.Arguments{
		{
			Type: abiUint32,
		},
	}
)

func decodeExtraArgsV1V2(extraArgs []byte) (gasLimit *big.Int, err error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	var method string
	if bytes.Equal(extraArgs[:4], evmExtraArgsV1Tag) {
		method = evmV1DecodeName
	} else if bytes.Equal(extraArgs[:4], evmExtraArgsV2Tag) {
		method = evmV2DecodeName
	} else {
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}
	ifaces, err := messageHasherABI.Methods[method].Inputs.UnpackValues(extraArgs[4:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args v1: %w", err)
	}
	// gas limit is always the first argument, and allow OOO isn't set explicitly
	// on the message.
	_, ok := ifaces[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("expected *big.Int, got %T", ifaces[0])
	}
	return ifaces[0].(*big.Int), nil
}

// abiEncodeMethodInputs encodes the inputs for a method call.
// example abi: `[{ "name" : "method", "type": "function", "inputs": [{"name": "a", "type": "uint256"}]}]`
func abiEncodeMethodInputs(abiDef abi.ABI, inputs ...interface{}) ([]byte, error) {
	packed, err := abiDef.Pack("method", inputs...)
	if err != nil {
		return nil, err
	}
	return packed[4:], nil // remove the method selector
}

func ABITypeOrPanic(t string) abi.Type {
	abiType, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(err)
	}
	return abiType
}

// Decodes the given bytes into a uint32, based on the encoding of destGasAmount in FeeQuoter.sol
func decodeTokenDestGasOverhead(destExecData []byte) (uint32, error) {
	ifaces, err := TokenDestGasOverheadABI.UnpackValues(destExecData)
	if err != nil {
		return 0, fmt.Errorf("abi decode TokenDestGasOverheadABI: %w", err)
	}
	_, ok := ifaces[0].(uint32)
	if !ok {
		return 0, fmt.Errorf("expected uint32, got %T", ifaces[0])
	}
	return ifaces[0].(uint32), nil
}
