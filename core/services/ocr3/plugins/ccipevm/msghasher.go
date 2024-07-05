package ccipevm

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
)

var (
	// bytes32 internal constant LEAF_DOMAIN_SEPARATOR = 0x0000000000000000000000000000000000000000000000000000000000000000;
	leafDomainSeparator = [32]byte{}

	// bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
	ANY_2_EVM_MESSAGE_HASH = utils.Keccak256Fixed([]byte("Any2EVMMessageHashV1"))

	messageHasherABI = types.MustGetABI(message_hasher.MessageHasherABI)
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "EVM2EVMMultiOnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	// TODO: move these to CCIPMsg instead?
	destChainSelector cciptypes.ChainSelector
	onrampAddress     []byte
}

func NewMessageHasherV1(
	onrampAddress []byte,
	destChainSelector cciptypes.ChainSelector,
) *MessageHasherV1 {
	return &MessageHasherV1{
		destChainSelector: destChainSelector,
		onrampAddress:     onrampAddress,
	}
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
		keccak256(encodedTokenAmounts),
		keccak256(encodedSourceTokenData),
	)
*/
func (h *MessageHasherV1) Hash(_ context.Context, msg cciptypes.CCIPMsg) (cciptypes.Bytes32, error) {
	tokenAmounts := make([]evm_2_evm_multi_onramp.ClientEVMTokenAmount, len(msg.TokenAmounts))
	for i, ta := range msg.TokenAmounts {
		tokenAmounts[i] = evm_2_evm_multi_onramp.ClientEVMTokenAmount{
			Token:  common.HexToAddress(string(ta.Token)),
			Amount: ta.Amount,
		}
	}
	encodedTokens, err := h.abiEncode("encodeTokenAmountsHashPreimage", tokenAmounts)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode token amounts: %w", err)
	}

	encodedSourceTokenData, err := h.abiEncode("encodeSourceTokenDataHashPreimage", msg.SourceTokenData)
	if err != nil {
		return [32]byte{}, fmt.Errorf("pack source token data: %w", err)
	}

	metaDataHashInput, err := h.abiEncode(
		"encodeMetadataHashPreimage",
		ANY_2_EVM_MESSAGE_HASH,
		uint64(msg.SourceChain),
		uint64(h.destChainSelector),
		h.onrampAddress,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode metadata hash input: %w", err)
	}

	var msgID [32]byte
	decoded, err := hex.DecodeString(msg.ID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode message ID: %w", err)
	}
	if len(decoded) != 32 {
		return [32]byte{}, fmt.Errorf("message ID must be 32 bytes")
	}
	copy(msgID[:], decoded)

	// NOTE: msg.Sender is not necessarily an EVM address since this is Any2EVM.
	// Accordingly, sender is defined as "bytes" in the onchain message definition
	// rather than "address".
	// However, its not clear how best to translate from Sender being a string representation
	// to bytes. For now, we assume that the string is hex encoded, but ideally Sender would
	// just be a byte array in the CCIPMsg struct that represents a sender encoded in the
	// source chain family encoding scheme.
	decodedSender, err := hex.DecodeString(
		strings.TrimPrefix(string(msg.Sender), "0x"),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode sender '%s': %w", msg.Sender, err)
	}
	fixedSizeFieldsEncoded, err := h.abiEncode(
		"encodeFixedSizeFieldsHashPreimage",
		msgID,
		decodedSender,
		common.HexToAddress(string(msg.Receiver)),
		uint64(msg.SeqNum),
		msg.ChainFeeLimit.Int,
		msg.Nonce,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode fixed size values: %w", err)
	}

	packedValues, err := h.abiEncode(
		"encodeFinalHashPreimage",
		leafDomainSeparator,
		utils.Keccak256Fixed(metaDataHashInput),
		utils.Keccak256Fixed(fixedSizeFieldsEncoded),
		utils.Keccak256Fixed(msg.Data),
		utils.Keccak256Fixed(encodedTokens),
		utils.Keccak256Fixed(encodedSourceTokenData),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode packed values: %w", err)
	}

	return utils.Keccak256Fixed(packedValues), nil
}

func (h *MessageHasherV1) abiEncode(method string, values ...interface{}) ([]byte, error) {
	res, err := messageHasherABI.Pack(method, values...)
	if err != nil {
		return nil, err
	}
	// trim the method selector.
	return res[4:], nil
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
