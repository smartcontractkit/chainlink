package ccipevm

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
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

// decodeExtraArgsV1V2 decodes the given EVM Extra Args and extracts the gas limit
// that was specified.
func decodeExtraArgsV1V2(extraArgs []byte) (gasLimit *big.Int, err error) {
	if len(extraArgs) < 4 {
		return nil, fmt.Errorf("extra args too short: %d, should be at least 4 bytes long (i.e the extraArgs tag)", len(extraArgs))
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

// CCIPMsgToAny2EVMMessage converts a ccipocr3.Message object to an offramp.InternalAny2EVMRampMessage object.
// These are typically used to create the execution report for EVM.
func CCIPMsgToAny2EVMMessage(msg ccipocr3.Message) (offramp.InternalAny2EVMRampMessage, error) {
	var tokenAmounts []offramp.InternalAny2EVMTokenTransfer
	for _, rta := range msg.TokenAmounts {
		destGasAmount, err := abiDecodeUint32(rta.DestExecData)
		if err != nil {
			return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to decode dest gas amount: %w", err)
		}

		tokenAmounts = append(tokenAmounts, offramp.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: common.LeftPadBytes(rta.SourcePoolAddress, 32),
			DestTokenAddress:  common.BytesToAddress(rta.DestTokenAddress),
			ExtraData:         rta.ExtraData[:],
			Amount:            rta.Amount.Int,
			DestGasAmount:     destGasAmount,
		})
	}

	gasLimit, err := decodeExtraArgsV1V2(msg.ExtraArgs)
	if err != nil {
		return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to decode extra args: %w", err)
	}

	return offramp.InternalAny2EVMRampMessage{
		Header: offramp.InternalRampMessageHeader{
			MessageId:           msg.Header.MessageID,
			SourceChainSelector: uint64(msg.Header.SourceChainSelector),
			DestChainSelector:   uint64(msg.Header.DestChainSelector),
			SequenceNumber:      uint64(msg.Header.SequenceNumber),
			Nonce:               msg.Header.Nonce,
		},
		Sender:       common.LeftPadBytes(msg.Sender, 32),
		Data:         msg.Data,
		Receiver:     common.BytesToAddress(msg.Receiver),
		GasLimit:     gasLimit,
		TokenAmounts: tokenAmounts,
	}, nil
}

// EVM2AnyToCCIPMsg converts an offramp.InternalEVM2AnyRampMessage object to a ccipocr3.Message object.
// These are typically used to calculate the message hash.
func EVM2AnyToCCIPMsg(
	onrampAddress common.Address,
	any2EVM onramp.InternalEVM2AnyRampMessage,
) ccipocr3.Message {
	var tokenAmounts []ccipocr3.RampTokenAmount
	for _, ta := range any2EVM.TokenAmounts {
		tokenAmounts = append(tokenAmounts, ccipocr3.RampTokenAmount{
			SourcePoolAddress: ta.SourcePoolAddress.Bytes(),
			DestTokenAddress:  ta.DestTokenAddress,
			DestExecData:      ta.DestExecData,
			ExtraData:         ta.ExtraData,
			Amount:            ccipocr3.NewBigInt(ta.Amount),
		})
	}
	return ccipocr3.Message{
		Header: ccipocr3.RampMessageHeader{
			MessageID:           ccipocr3.Bytes32(any2EVM.Header.MessageId),
			SourceChainSelector: ccipocr3.ChainSelector(any2EVM.Header.SourceChainSelector),
			DestChainSelector:   ccipocr3.ChainSelector(any2EVM.Header.DestChainSelector),
			SequenceNumber:      ccipocr3.SeqNum(any2EVM.Header.SequenceNumber),
			Nonce:               any2EVM.Header.Nonce,
			OnRamp:              onrampAddress.Bytes(),
		},
		Sender:         any2EVM.Sender.Bytes(),
		Data:           any2EVM.Data,
		Receiver:       ccipocr3.UnknownAddress(any2EVM.Receiver),
		ExtraArgs:      any2EVM.ExtraArgs,
		FeeToken:       any2EVM.FeeToken.Bytes(),
		FeeTokenAmount: ccipocr3.NewBigInt(any2EVM.FeeTokenAmount),
		FeeValueJuels:  ccipocr3.NewBigInt(any2EVM.FeeValueJuels),
		TokenAmounts:   tokenAmounts,
	}
}

// BoolsToBitFlags transforms a list of boolean flags to a *big.Int encoded number.
func BoolsToBitFlags(bools []bool) *big.Int {
	encodedFlags := big.NewInt(0)
	for i := 0; i < len(bools); i++ {
		if bools[i] {
			encodedFlags.SetBit(encodedFlags, i, 1)
		}
	}
	return encodedFlags
}
