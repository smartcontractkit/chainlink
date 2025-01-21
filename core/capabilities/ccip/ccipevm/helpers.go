package ccipevm

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
)

const (
	svmV1DecodeName = "decodeSVMExtraArgsV1"
	evmV1DecodeName = "decodeEVMExtraArgsV1"
	evmV2DecodeName = "decodeEVMExtraArgsV2"
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

func decodeExtraArgsSVMV1(extraArgs []byte) (*message_hasher.ClientSVMExtraArgsV1, error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 (i.e the extraArgs tag)", len(extraArgs))
	}

	if !bytes.Equal(extraArgs[:4], svmExtraArgsV1Tag) {
		return nil, fmt.Errorf("unknown extra args tag: %x", extraArgs)
	}

	ifaces, err := messageHasherABI.Methods[svmV1DecodeName].Inputs.UnpackValues(extraArgs[4:])
	if err != nil {
		return nil, fmt.Errorf("abi decode extra args v1: %w", err)
	}

	if len(ifaces) != 5 {
		return nil, fmt.Errorf("expected 5 inputs, got %d", len(ifaces))
	}

	_, ok := ifaces[0].(uint32)
	if !ok {
		return nil, fmt.Errorf("expected uint32, got %T", ifaces[0])
	}

	_, ok = ifaces[1].(uint64)
	if !ok {
		return nil, fmt.Errorf("expected uint64, got %T", ifaces[1])
	}

	_, ok = ifaces[2].(bool)
	if !ok {
		return nil, fmt.Errorf("expected bool, got %T", ifaces[2])
	}

	tokenReceiver, ok := ifaces[3].([]byte)
	if !ok || len(tokenReceiver) != 32 {
		return nil, fmt.Errorf("expected [32]byte, got %T or incorrect length %d", ifaces[3], len(tokenReceiver))
	}

	var tokenReceiverArray [32]byte
	copy(tokenReceiverArray[:], tokenReceiver)

	accounts, ok := ifaces[4].([][32]byte)
	if !ok {
		return nil, fmt.Errorf("expected [][32]byte, got %T", ifaces[4])
	}

	return &message_hasher.ClientSVMExtraArgsV1{
		ComputeUnits:             ifaces[0].(uint32),
		AccountIsWritableBitmap:  ifaces[1].(uint64),
		AllowOutOfOrderExecution: ifaces[2].(bool),
		TokenReceiver:            tokenReceiverArray,
		Accounts:                 accounts,
	}, nil
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
