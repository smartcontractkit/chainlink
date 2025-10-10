package ccipsui

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-ccip/pkg/logutil"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var (
	// const LEAF_DOMAIN_SEPARATOR: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
	leafDomainSeparator = [32]byte{}

	// see sui_hash::keccak256(b"Any2SuiMessageHashV1") in calculate_metadata_hash
	any2SuiMessageHash = utils.Keccak256Fixed([]byte("Any2SuiMessageHashV1"))
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with ccip::offramp version 1.6.0
type MessageHasherV1 struct {
	lggr           logger.Logger
	extraDataCodec ccipocr3common.ExtraDataCodecBundle
}

type any2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

func NewMessageHasherV1(lggr logger.Logger, extraDataCodec ccipocr3common.ExtraDataCodecBundle) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr:           lggr,
		extraDataCodec: extraDataCodec,
	}
}

// Hash implements the MessageHasher interface.
// It constructs all of the inputs to the final keccak256 hash in Internal._hash(Any2EVMRampMessage).
// The main structure of the hash is as follows:
// Fixed-size message fields are included in nested hash to reduce stack pressure.
// This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
func (h *MessageHasherV1) Hash(ctx context.Context, msg ccipocr3common.Message) (ccipocr3common.Bytes32, error) {
	lggr := logutil.WithContextValues(ctx, h.lggr)
	lggr = logger.With(
		lggr,
		"msgID", msg.Header.MessageID.String(),
		"ANY_2_SUI_MESSAGE_HASH", hexutil.Encode(any2SuiMessageHash[:]),
		"onrampAddress", msg.Header.OnRamp,
	)
	lggr.Debugw("hashing message", "msg", msg)

	rampTokenAmounts := make([]any2SuiTokenTransfer, len(msg.TokenAmounts))
	for i, rta := range msg.TokenAmounts {
		destExecDataDecodedMap, err := h.extraDataCodec.DecodeTokenAmountDestExecData(rta.DestExecData, msg.Header.SourceChainSelector)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to decode dest exec data: %w", err)
		}

		destGasAmountValue, ok := destExecDataDecodedMap["destGasAmount"]
		if !ok {
			return [32]byte{}, errors.New("destGasAmount not found in destExecDataDecodedMap")
		}

		destGasAmount, ok := destGasAmountValue.(uint32)
		if !ok {
			return [32]byte{}, fmt.Errorf("invalid type for destGasAmount, expected uint32, got %T", destGasAmount)
		}

		lggr.Debugw("decoded dest gas amount",
			"destGasAmount", destGasAmount)

		destTokenAddress, err := addressBytesToBytes32(rta.DestTokenAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf("decode dest token address: %w", err)
		}

		lggr.Debugw("abi decoded dest token address",
			"destTokenAddress", destTokenAddress)

		rampTokenAmounts[i] = any2SuiTokenTransfer{
			SourcePoolAddress: rta.SourcePoolAddress,
			DestTokenAddress:  destTokenAddress,
			DestGasAmount:     destGasAmount,
			ExtraData:         rta.ExtraData,
			Amount:            rta.Amount.Int,
		}
	}

	// one difference from EVM is that we don't left pad the OnRamp to 32 bytes here, we use the source chain's canonical bytes encoding directly.
	metaDataHashInput, err := computeMetadataHash(uint64(msg.Header.SourceChainSelector), uint64(msg.Header.DestChainSelector), msg.Header.OnRamp)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode metadata hash input: %w", err)
	}

	lggr.Debugw("metadata hash preimage",
		"metaDataHashInput", hexutil.Encode(metaDataHashInput[:]))

	// Need to decode the extra args to get the gas limit.
	// TODO: we assume that extra args is always abi-encoded for now, but we need
	// to decode according to source chain selector family. We should add a family
	// lookup API to the chain-selectors library.

	decodedExtraArgsMap, err := h.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return [32]byte{}, err
	}

	gasLimit, tokenReceiver, err := parseExtraDataMap(decodedExtraArgsMap)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode extra args to get gas limit: %w", err)
	}

	lggr.Debugw("decoded msg gas limit", "gasLimit", gasLimit)

	receiverAddress, err := addressBytesToBytes32(msg.Receiver)
	if err != nil {
		return [32]byte{}, err
	}

	msgHash, err := computeMessageDataHash(metaDataHashInput, msg.Header.MessageID, receiverAddress, uint64(msg.Header.SequenceNumber), gasLimit, tokenReceiver, msg.Header.Nonce, msg.Sender, msg.Data, rampTokenAmounts)
	if err != nil {
		return [32]byte{}, err
	}

	lggr.Debugw("final message hash result",
		"msgHash", hexutil.Encode(msgHash[:]),
	)

	return msgHash, nil
}

// This is the equivalent of ccip_offramp::calculate_message_hash.
// This is similar to the EVM version, except for 32-byte addresses and no dynamic offsets.
func computeMessageDataHash(
	metadataHash [32]byte,
	messageID [32]byte,
	receiver [32]byte,
	sequenceNumber uint64,
	gasLimit *big.Int,
	tokenReceiver [32]byte,
	nonce uint64,
	sender []byte,
	data []byte,
	tokenAmounts []any2SuiTokenTransfer,
) ([32]byte, error) {
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint64 ABI type: %w", err)
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint256 ABI type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create bytes32 ABI type: %w", err)
	}

	headerArgs := abi.Arguments{
		{Type: bytes32Type}, // messageID
		{Type: bytes32Type}, // receiver as bytes32
		{Type: uint64Type},  // sequenceNumber
		{Type: uint256Type}, // gasLimit
		{Type: bytes32Type}, // tokenReceiver
		{Type: uint64Type},  // nonce
	}
	headerEncoded, err := headerArgs.Pack(
		messageID,
		receiver,
		sequenceNumber,
		gasLimit,
		tokenReceiver,
		nonce,
	)
	if err != nil {
		return [32]byte{}, err
	}
	headerHash := crypto.Keccak256Hash(headerEncoded)

	senderHash := crypto.Keccak256Hash(sender)

	dataHash := crypto.Keccak256Hash(data)

	// Manually encode tokens to match the Move implementation, because abi.Pack has different behavior
	// for dynamic types.
	var tokenHashData []byte
	tokenHashData = append(tokenHashData, encodeUint256(big.NewInt(int64(len(tokenAmounts))))...)
	for _, token := range tokenAmounts {
		tokenHashData = append(tokenHashData, encodeBytes(token.SourcePoolAddress)...)
		tokenHashData = append(tokenHashData, token.DestTokenAddress[:]...)
		tokenHashData = append(tokenHashData, encodeUint32(token.DestGasAmount)...)
		tokenHashData = append(tokenHashData, encodeBytes(token.ExtraData)...)
		tokenHashData = append(tokenHashData, encodeUint256(token.Amount)...)
	}
	tokenAmountsHash := crypto.Keccak256Hash(tokenHashData)

	finalArgs := abi.Arguments{
		{Type: bytes32Type}, // LEAF_DOMAIN_SEPARATOR
		{Type: bytes32Type}, // metadataHash
		{Type: bytes32Type}, // headerHash
		{Type: bytes32Type}, // senderHash
		{Type: bytes32Type}, // dataHash
		{Type: bytes32Type}, // tokenAmountsHash
	}

	finalEncoded, err := finalArgs.Pack(
		leafDomainSeparator,
		metadataHash,
		headerHash,
		senderHash,
		dataHash,
		tokenAmountsHash,
	)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(finalEncoded), nil
}

// This is the equivalent of ccip_offramp::calculate_metadata_hash.
// This is similar to the EVM version, except for the separator, 32-byte addresses, and no dynamic offsets.
func computeMetadataHash(
	sourceChainSelector uint64,
	destinationChainSelector uint64,
	onRamp []byte,
) ([32]byte, error) {
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint64 ABI type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create bytes32 ABI type: %w", err)
	}

	onRampHash := crypto.Keccak256Hash(onRamp)

	args := abi.Arguments{
		{Type: bytes32Type}, // ANY_2_SUI_MESSAGE_HASH
		{Type: uint64Type},  // sourceChainSelector
		{Type: uint64Type},  // destinationChainSelector (i_chainSelector)
		{Type: bytes32Type}, // onRamp
	}

	encoded, err := args.Pack(
		any2SuiMessageHash,
		sourceChainSelector,
		destinationChainSelector,
		onRampHash,
	)
	if err != nil {
		return [32]byte{}, err
	}

	metadataHash := crypto.Keccak256Hash(encoded)
	return metadataHash, nil
}

func encodeUint256(n *big.Int) []byte {
	return common.LeftPadBytes(n.Bytes(), 32)
}

func encodeUint32(n uint32) []byte {
	return common.LeftPadBytes(new(big.Int).SetUint64(uint64(n)).Bytes(), 32)
}

func encodeBytes(b []byte) []byte {
	encodedLength := common.LeftPadBytes(big.NewInt(int64(len(b))).Bytes(), 32)
	padLen := (32 - (len(b) % 32)) % 32
	result := make([]byte, 32+len(b)+padLen)
	copy(result[:32], encodedLength)
	copy(result[32:], b)
	return result
}

func parseExtraDataMap(input map[string]any) (*big.Int, [32]byte, error) {
	outputGas, ok := input["gasLimit"]
	if !ok {
		return nil, [32]byte{}, errors.New("gas limit not found in extra data map")
	}
	outputGasInt, ok := outputGas.(*big.Int)
	if !ok {
		return nil, [32]byte{}, errors.New("gas limit not a *big.Int")
	}

	tokenReceiver, ok := input["tokenReceiver"]
	if !ok {
		return nil, [32]byte{}, errors.New("token receiver not found in extra data map")
	}
	tokenReceiverBytes, ok := tokenReceiver.([32]byte)
	if !ok {
		return nil, [32]byte{}, errors.New("token receiver not a [32]byte")
	}
	return outputGasInt, tokenReceiverBytes, nil
}

func extractDestGasAmountFromMap(input map[string]any) (uint32, error) {
	// Iterate through the expected fields in the struct
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "destgasamount":
			// Expect uint32
			if val, ok := fieldValue.(uint32); ok {
				return val, nil
			}
			return 0, errors.New("invalid type for destgasamount, expected uint32")
		default:
		}
	}

	return 0, errors.New("invalid token message, dest gas amount not found in the DestExecDataDecoded map")
}

func addressBytesToBytes32(addr []byte) ([32]byte, error) {
	if len(addr) > 32 {
		return [32]byte{}, fmt.Errorf("invalid Sui address length, expected 32, got %d", len(addr))
	}
	var result [32]byte
	// Left pad by copying to the end of the 32 byte array
	copy(result[32-len(addr):], addr)
	return result, nil
}

// Interface compliance check
var _ ccipocr3common.MessageHasher = (*MessageHasherV1)(nil)
