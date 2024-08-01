package ccipevm

import (
	"bytes"
	"fmt"
	"math/big"
)

func decodeExtraArgsV1V2(extraArgs []byte) (gasLimit *big.Int, err error) {
	var method string
	if bytes.Equal(extraArgs[:4], evmExtraArgsV1Tag) {
		method = "decodeEVMExtraArgsV1"
	} else if bytes.Equal(extraArgs[:4], evmExtraArgsV2Tag) {
		method = "decodeEVMExtraArgsV2"
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
