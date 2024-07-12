package ccipevm

import (
	"context"
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
)

var (
	// bytes32 internal constant LEAF_DOMAIN_SEPARATOR = 0x0000000000000000000000000000000000000000000000000000000000000000;
	leafDomainSeparator = [32]byte{}

	// bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
	ANY_2_EVM_MESSAGE_HASH = utils.Keccak256Fixed([]byte("Any2EVMMessageHashV1"))

	messageHasherABI = types.MustGetABI(message_hasher.MessageHasherABI)

	// bytes4 public constant EVM_EXTRA_ARGS_V1_TAG = 0x97a657c9;
	evmExtraArgsV1Tag = hexutil.MustDecode("0x97a657c9")

	// bytes4 public constant EVM_EXTRA_ARGS_V2_TAG = 0x181dcf10;
	evmExtraArgsV2Tag = hexutil.MustDecode("0x181dcf10")
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "EVM2EVMMultiOnRamp 1.6.0-dev"
type MessageHasherV1 struct{}

func NewMessageHasherV1() *MessageHasherV1 {
	return &MessageHasherV1{}
}

// Hash implements the MessageHasher interface.
// It constructs all of the inputs to the final keccak256 hash in Internal._hash(Any2EVMRampMessage).
// The main structure of the hash is as follows:
/*
	keccak256(
		leafDomainSeparator,
		keccak256(any_2_evm_message_hash, header.sourceChainSelector, header.destinationChainSelector, onRamp),
		keccak256(fixedSizeMessageFields),
		keccak256(messageData),
		keccak256(encodedRampTokenAmounts),
	)
*/
func (h *MessageHasherV1) Hash(_ context.Context, msg cciptypes.Message) (cciptypes.Bytes32, error) {
	var rampTokenAmounts []message_hasher.InternalRampTokenAmount
	for _, rta := range msg.TokenAmounts {
		rampTokenAmounts = append(rampTokenAmounts, message_hasher.InternalRampTokenAmount{
			SourcePoolAddress: rta.SourcePoolAddress,
			DestTokenAddress:  rta.DestTokenAddress,
			ExtraData:         rta.ExtraData,
			Amount:            rta.Amount.Int,
		})
	}
	encodedRampTokenAmounts, err := abiEncode("encodeTokenAmountsHashPreimage", rampTokenAmounts)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode token amounts: %w", err)
	}

	metaDataHashInput, err := abiEncode(
		"encodeMetadataHashPreimage",
		ANY_2_EVM_MESSAGE_HASH,
		uint64(msg.Header.SourceChainSelector),
		uint64(msg.Header.DestChainSelector),
		[]byte(msg.Header.OnRamp),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode metadata hash input: %w", err)
	}

	// Need to decode the extra args to get the gas limit.
	// TODO: we assume that extra args is always abi-encoded for now, but we need
	// to decode according to source chain selector family. We should add a family
	// lookup API to the chain-selectors library.
	gasLimit, err := decodeExtraArgsV1V2(msg.ExtraArgs)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode extra args: %w", err)
	}

	fixedSizeFieldsEncoded, err := abiEncode(
		"encodeFixedSizeFieldsHashPreimage",
		msg.Header.MessageID,
		[]byte(msg.Sender),
		common.BytesToAddress(msg.Receiver),
		uint64(msg.Header.SequenceNumber),
		gasLimit,
		msg.Header.Nonce,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode fixed size values: %w", err)
	}

	packedValues, err := abiEncode(
		"encodeFinalHashPreimage",
		leafDomainSeparator,
		utils.Keccak256Fixed(metaDataHashInput),
		utils.Keccak256Fixed(fixedSizeFieldsEncoded),
		utils.Keccak256Fixed(msg.Data),
		utils.Keccak256Fixed(encodedRampTokenAmounts),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode packed values: %w", err)
	}

	return utils.Keccak256Fixed(packedValues), nil
}

func abiEncode(method string, values ...interface{}) ([]byte, error) {
	res, err := messageHasherABI.Pack(method, values...)
	if err != nil {
		return nil, err
	}
	// trim the method selector.
	return res[4:], nil
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
