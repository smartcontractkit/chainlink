package ccipevm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

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

// CCIPMsgToAny2EVMMessage converts a ccipocr3.Message object to an offramp.InternalAny2EVMRampMessage object.
// These are typically used to create the execution report for EVM.
func CCIPMsgToAny2EVMMessage(msg ccipocr3.Message, codec ccipcommon.ExtraDataCodec) (offramp.InternalAny2EVMRampMessage, error) {
	var tokenAmounts []offramp.InternalAny2EVMTokenTransfer
	for _, rta := range msg.TokenAmounts {
		decodedMap, err := codec.DecodeTokenAmountDestExecData(rta.DestExecData, msg.Header.SourceChainSelector)
		if err != nil {
			return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to decode dest gas amount: %w", err)
		}
		destGasAmount, err := extractDestGasAmountFromMap(decodedMap)
		if err != nil {
			return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to extract dest gas amount from decodec map: %w", err)
		}

		tokenAmounts = append(tokenAmounts, offramp.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: common.LeftPadBytes(rta.SourcePoolAddress, 32),
			DestTokenAddress:  common.BytesToAddress(rta.DestTokenAddress),
			ExtraData:         rta.ExtraData[:],
			Amount:            rta.Amount.Int,
			DestGasAmount:     destGasAmount,
		})
	}

	decodeMap, err := codec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to decode extra args: %w", err)
	}

	gasLimit, err := parseExtraArgsMap(decodeMap)
	if err != nil {
		return offramp.InternalAny2EVMRampMessage{}, fmt.Errorf("failed to get gasLimit by parsing extra args map: %w", err)
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
