package ccipevm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pkg/logutil"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/message_hasher"
)

var (
	// bytes32 internal constant LEAF_DOMAIN_SEPARATOR = 0x0000000000000000000000000000000000000000000000000000000000000000;
	LEAF_DOMAIN_SEPARATOR = [32]byte{}

	// bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
	ANY_2_EVM_MESSAGE_HASH = utils.Keccak256Fixed([]byte("Any2EVMMessageHashV1"))

	messageHasherABI = types.MustGetABI(message_hasher.MessageHasherABI)

	// bytes4 public constant EVM_EXTRA_ARGS_V1_TAG = 0x97a657c9;
	evmExtraArgsV1Tag = hexutil.MustDecode("0x97a657c9")

	// bytes4 public constant EVM_EXTRA_ARGS_V2_TAG = 0x181dcf10;
	evmExtraArgsV2Tag = hexutil.MustDecode("0x181dcf10")

	// bytes4 public constant SVM_EXTRA_EXTRA_ARGS_V1_TAG = 0x1f3b3aba
	svmExtraArgsV1Tag = hexutil.MustDecode("0x1f3b3aba")
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "OnRamp 1.6.0"
type MessageHasherV1 struct {
	lggr           logger.Logger
	extraDataCodec ccipcommon.ExtraDataCodec
}

func NewMessageHasherV1(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr:           lggr,
		extraDataCodec: extraDataCodec,
	}
}

// Hash implements the MessageHasher interface.
// It constructs all of the inputs to the final keccak256 hash in Internal._hash(Any2EVMRampMessage).
// The main structure of the hash is as follows:
/*
	// Fixed-size message fields are included in nested hash to reduce stack pressure.
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        metadataHash,
        keccak256(
          abi.encode(
            original.header.messageId,
            original.receiver,
            original.header.sequenceNumber,
            original.gasLimit,
            original.header.nonce
          )
        ),
        keccak256(original.sender),
        keccak256(original.data),
        keccak256(abi.encode(original.tokenAmounts))
      )
    );
*/
func (h *MessageHasherV1) Hash(ctx context.Context, msg cciptypes.Message) (cciptypes.Bytes32, error) {
	lggr := logutil.WithContextValues(ctx, h.lggr)
	lggr = logger.With(
		lggr,
		"msgID", msg.Header.MessageID.String(),
		"ANY_2_EVM_MESSAGE_HASH", hexutil.Encode(ANY_2_EVM_MESSAGE_HASH[:]),
		"onrampAddress", msg.Header.OnRamp,
	)
	lggr.Debugw("hashing message", "msg", msg)

	any2EVM, onRamp, err := messageToAny2EVM(msg, h.extraDataCodec)
	if err != nil {
		return [32]byte{}, fmt.Errorf("convert message to Any2EVM: %w", err)
	}

	return h.hashAny2EVM(any2EVM, onRamp, lggr)
}

func (h *MessageHasherV1) hashAny2EVM(
	any2EVM message_hasher.InternalAny2EVMRampMessage,
	onRamp []byte,
	lggr logger.Logger) ([32]byte, error) {
	encodedRampTokenAmounts, err := h.abiEncode(
		"encodeAny2EVMTokenAmountsHashPreimage",
		any2EVM.TokenAmounts,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode token amounts: %w", err)
	}

	lggr.Debugw("token amounts preimage",
		"encodedRampTokenAmounts", hexutil.Encode(encodedRampTokenAmounts))

	metaDataHashInput, err := h.abiEncode(
		"encodeMetadataHashPreimage",
		ANY_2_EVM_MESSAGE_HASH,
		any2EVM.Header.SourceChainSelector,
		any2EVM.Header.DestChainSelector,
		utils.Keccak256Fixed(onRamp),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode metadata hash input: %w", err)
	}

	lggr.Debugw("metadata hash preimage",
		"metaDataHashInput", hexutil.Encode(metaDataHashInput))

	fixedSizeFieldsEncoded, err := h.abiEncode(
		"encodeFixedSizeFieldsHashPreimage",
		any2EVM.Header.MessageId,
		any2EVM.Receiver,
		any2EVM.Header.SequenceNumber,
		any2EVM.GasLimit,
		any2EVM.Header.Nonce,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode fixed size values: %w", err)
	}

	lggr.Debugw("fixed size fields has preimage",
		"fixedSizeFieldsEncoded", hexutil.Encode(fixedSizeFieldsEncoded))

	hashPreimage, err := h.abiEncode(
		"encodeFinalHashPreimage",
		LEAF_DOMAIN_SEPARATOR,
		utils.Keccak256Fixed(metaDataHashInput), // metaDataHash
		utils.Keccak256Fixed(fixedSizeFieldsEncoded),
		utils.Keccak256Fixed(any2EVM.Sender), // todo: this is not chain-agnostic
		utils.Keccak256Fixed(any2EVM.Data),
		utils.Keccak256Fixed(encodedRampTokenAmounts),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode packed values: %w", err)
	}

	msgHash := utils.Keccak256Fixed(hashPreimage)

	lggr.Debugw("final hash preimage and message hash result",
		"hashPreimage", hexutil.Encode(hashPreimage),
		"msgHash", hexutil.Encode(msgHash[:]),
	)

	return msgHash, nil
}

func (h *MessageHasherV1) abiEncode(method string, values ...interface{}) ([]byte, error) {
	res, err := messageHasherABI.Pack(method, values...)
	if err != nil {
		return nil, err
	}
	// trim the method selector.
	return res[4:], nil
}

func abiDecodeUint32(data []byte) (uint32, error) {
	raw, err := utils.ABIDecode(`[{ "type": "uint32" }]`, data)
	if err != nil {
		return 0, fmt.Errorf("abi decode uint32: %w", err)
	}

	val := *abi.ConvertType(raw[0], new(uint32)).(*uint32)
	return val, nil
}

func abiEncodeUint32(data uint32) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "uint32" }]`, data)
}

// abiEncodeAddress encodes the given address as a solidity address.
// TODO: this is potentially incorrect for nonEVM sources.
// we need to revisit.
// e.g on Solana, we would be abi.encode()ing bytes or bytes32.
// encoding 20 bytes as a solidity bytes is not the same as encoding a 20 byte address
// or a bytes32.
func abiEncodeAddress(data common.Address) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "address" }]`, data)
}

func abiDecodeAddress(data []byte) (common.Address, error) {
	raw, err := utils.ABIDecode(`[{ "type": "address" }]`, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("abi decode address: %w", err)
	}

	val := *abi.ConvertType(raw[0], new(common.Address)).(*common.Address)
	return val, nil
}

func parseExtraDataMap(input map[string]any) (*big.Int, error) {
	var outputGas *big.Int
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "gaslimit":
			// Expect [][32]byte
			if val, ok := fieldValue.(*big.Int); ok {
				outputGas = val
				return outputGas, nil
			} else {
				return nil, fmt.Errorf("unexpected type for gas limit: %T", fieldValue)
			}
		default:
			// no error here, as we only need the keys to gasLimit, other keys can be skipped without like AllowOutOfOrderExecution	etc.
		}
	}
	return outputGas, errors.New("gas limit not found in extra data map")
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
			} else {
				return 0, errors.New("invalid type for destgasamount, expected uint32")
			}
		default:
		}
	}

	return 0, errors.New("invalid token message, dest gas amount not found in the DestExecDataDecoded map")
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)

// messageToAny2EVM converts a ccip message to an InternalAny2EVMRampMessage.
// It handles source-chain specific quirks like padding addresses to 32 bytes where applicable.
// It also returns the onRamp address as a 32-byte padded address where applicable.
func messageToAny2EVM(msg cciptypes.Message, extraDataCodec ccipcommon.ExtraDataCodec) (ret message_hasher.InternalAny2EVMRampMessage, onRamp []byte, err error) {
	sourceFamily, err := chainsel.GetSelectorFamily(uint64(msg.Header.SourceChainSelector))
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, nil, err
	}

	decodedExtraArgsMap, err := extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, nil, err
	}

	gasLimit, err := parseExtraDataMap(decodedExtraArgsMap)
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, nil, fmt.Errorf("decode extra args to get gas limit: %w", err)
	}

	// we can fill out some fields here optimistically, mostly in the header.
	any2EVM := message_hasher.InternalAny2EVMRampMessage{
		Header: message_hasher.InternalRampMessageHeader{
			MessageId:           msg.Header.MessageID,
			SourceChainSelector: uint64(msg.Header.SourceChainSelector),
			DestChainSelector:   uint64(msg.Header.DestChainSelector),
			SequenceNumber:      uint64(msg.Header.SequenceNumber),
			Nonce:               msg.Header.Nonce,
		},
		// Data is always just passed through.
		Data:     msg.Data,
		GasLimit: gasLimit,
		// Since the receiver is on EVM, we can assume its at least 20 bytes long.
		// This kind of thing would be checked onchain at the source.
		Receiver: common.BytesToAddress(msg.Receiver),
		// Sender: , // this needs to be handled on a chain-specific basis.
		// TokenAmounts: , // this needs to be handled on a chain-specific basis.
	}
	switch sourceFamily {
	case chainsel.FamilyEVM:
		// if sourceFamily is EVM, then the sender was parsed as an `address` of 20 bytes long, so we have to left-pad it.
		any2EVM.Sender = common.LeftPadBytes(msg.Sender, 32)

		var rampTokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
		for _, rta := range msg.TokenAmounts {
			destGasAmount, err := abiDecodeUint32(rta.DestExecData)
			if err != nil {
				return message_hasher.InternalAny2EVMRampMessage{}, nil, fmt.Errorf("decode dest gas amount: %w", err)
			}

			// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
			// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
			// its unclear as of writing how we will handle non-EVM chains and their addresses.
			// e.g, will we encode them as bytes or bytes32?
			sourcePoolAddressABIEncodedAsAddress, err := abiEncodeAddress(common.BytesToAddress(rta.SourcePoolAddress))
			if err != nil {
				return message_hasher.InternalAny2EVMRampMessage{}, nil, fmt.Errorf("abi encode source pool address: %w", err)
			}

			destTokenAddress, err := abiDecodeAddress(rta.DestTokenAddress)
			if err != nil {
				return message_hasher.InternalAny2EVMRampMessage{}, nil, fmt.Errorf("decode dest token address: %w", err)
			}

			rampTokenAmounts = append(rampTokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
				SourcePoolAddress: sourcePoolAddressABIEncodedAsAddress,
				DestTokenAddress:  destTokenAddress,
				DestGasAmount:     destGasAmount,
				ExtraData:         rta.ExtraData,
				Amount:            rta.Amount.Int,
			})
		}

		any2EVM.TokenAmounts = rampTokenAmounts
		onRamp = common.LeftPadBytes(msg.Header.OnRamp, 32)
	// case chainsel.FamilySolana: // TODO: implement
	// case chainsel.FamilyAptos:  // TODO: implement
	default:
		return message_hasher.InternalAny2EVMRampMessage{}, nil, fmt.Errorf("source chain selector family %s not supported", sourceFamily)
	}

	return any2EVM, onRamp, nil
}
