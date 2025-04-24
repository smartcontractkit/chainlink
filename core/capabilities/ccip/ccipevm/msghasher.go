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
	ag_binary "github.com/gagliardetto/binary"
	chainsel "github.com/smartcontractkit/chain-selectors"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink-ccip/pkg/logutil"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
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
	sourceChainFamily, err := chainsel.GetSelectorFamily(uint64(msg.Header.SourceChainSelector))
	if err != nil {
		return [32]byte{}, fmt.Errorf("get source chain family: %w", err)
	}

	lggr := logutil.WithContextValues(ctx, h.lggr)
	lggr = logger.With(
		lggr,
		"msgID", msg.Header.MessageID.String(),
		"ANY_2_EVM_MESSAGE_HASH", hexutil.Encode(ANY_2_EVM_MESSAGE_HASH[:]),
		"onrampAddress", msg.Header.OnRamp,
		"sourceChainFamily", sourceChainFamily,
	)
	lggr.Debugw("hashing message", "msg", msg)

	var rampTokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
	for _, rta := range msg.TokenAmounts {
		destExecDataDecodedMap, err := h.extraDataCodec.DecodeTokenAmountDestExecData(rta.DestExecData, msg.Header.SourceChainSelector)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to decode dest exec data: %w", err)
		}

		destGasAmount, err := extractDestGasAmountFromMap(destExecDataDecodedMap)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed extract dest gas amount from decoded extradata map: %w", err)
		}

		lggr.Debugw("decoded dest gas amount",
			"destGasAmount", destGasAmount)

		destTokenAddress, err := abiDecodeAddress(rta.DestTokenAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf("decode dest token address: %w", err)
		}

		lggr.Debugw("abi decoded dest token address",
			"destTokenAddress", destTokenAddress)

		var sourcePoolAddr []byte
		if sourceChainFamily == chainsel.FamilyEVM {
			// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
			// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
			sourcePoolAddr, err = abiEncodeAddress(common.BytesToAddress(rta.SourcePoolAddress))
			if err != nil {
				return [32]byte{}, fmt.Errorf("abi encode source pool address: %w", err)
			}
		} else {
			sourcePoolAddr = rta.SourcePoolAddress
		}

		lggr.Debugw("resolved token amount fields", "sourcePoolAddress", sourcePoolAddr, "destTokenAddress", destTokenAddress, "destGasAmount", destGasAmount)
		rampTokenAmounts = append(rampTokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: sourcePoolAddr,
			DestTokenAddress:  destTokenAddress,
			DestGasAmount:     destGasAmount,
			ExtraData:         rta.ExtraData,
			Amount:            rta.Amount.Int,
		})
	}

	encodedRampTokenAmounts, err := h.abiEncode(
		"encodeAny2EVMTokenAmountsHashPreimage",
		rampTokenAmounts,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode token amounts: %w", err)
	}

	lggr.Debugw("token amounts preimage",
		"encodedRampTokenAmounts", hexutil.Encode(encodedRampTokenAmounts))

	metaDataHashInput, err := h.abiEncode(
		"encodeMetadataHashPreimage",
		ANY_2_EVM_MESSAGE_HASH,
		uint64(msg.Header.SourceChainSelector),
		uint64(msg.Header.DestChainSelector),
		// TODO: this is evm-specific padding, fix.
		// no-op if the onramp is already 32 bytes.
		utils.Keccak256Fixed(common.LeftPadBytes(msg.Header.OnRamp, 32)),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode metadata hash input: %w", err)
	}

	lggr.Debugw("metadata hash preimage",
		"metaDataHashInput", hexutil.Encode(metaDataHashInput))

	// Need to decode the extra args to get the gas limit.
	// TODO: we assume that extra args is always abi-encoded for now, but we need
	// to decode according to source chain selector family. We should add a family
	// lookup API to the chain-selectors library.

	decodedExtraArgsMap, err := h.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return [32]byte{}, err
	}

	gasLimit, err := parseExtraArgsMap(decodedExtraArgsMap)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode extra args to get gas limit: %w", err)
	}

	lggr.Debugw("decoded msg gas limit", "gasLimit", gasLimit)

	fixedSizeFieldsEncoded, err := h.abiEncode(
		"encodeFixedSizeFieldsHashPreimage",
		msg.Header.MessageID,
		common.BytesToAddress(msg.Receiver),
		uint64(msg.Header.SequenceNumber),
		gasLimit,
		msg.Header.Nonce,
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
		utils.Keccak256Fixed(common.LeftPadBytes(msg.Sender, 32)), // todo: this is not chain-agnostic
		utils.Keccak256Fixed(msg.Data),
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

func (h *MessageHasherV1) abiEncode(method string, values ...any) ([]byte, error) {
	res, err := messageHasherABI.Pack(method, values...)
	if err != nil {
		return nil, err
	}
	// trim the method selector.
	return res[4:], nil
}

func abiDecodeType[T any](argsABI abi.Arguments, data []byte) (T, error) {
	raw, err := argsABI.UnpackValues(data)
	if err != nil {
		val := *new(T)
		return val, fmt.Errorf("abi decode %T: %w", val, err)
	}

	val := *abi.ConvertType(raw[0], new(T)).(*T)
	return val, nil
}

var uint32ABI abi.Arguments = abi.Arguments{{Type: utils.MustAbiType("uint32", nil)}}

func abiEncodeUint32(data uint32) ([]byte, error) {
	return uint32ABI.Pack(data)
}

func abiDecodeUint32(data []byte) (uint32, error) {
	return abiDecodeType[uint32](uint32ABI, data)
}

var addressABI abi.Arguments = abi.Arguments{{Type: utils.MustAbiType("address", nil)}}

// abiEncodeAddress encodes the given address as a solidity address.
// TODO: this is potentially incorrect for nonEVM sources.
// we need to revisit.
// e.g on Solana, we would be abi.encode()ing bytes or bytes32.
// encoding 20 bytes as a solidity bytes is not the same as encoding a 20 byte address
// or a bytes32.
func abiEncodeAddress(data common.Address) ([]byte, error) {
	return addressABI.Pack(data)
}

func abiDecodeAddress(data []byte) (common.Address, error) {
	return abiDecodeType[common.Address](addressABI, data)
}

func parseExtraArgsMap(input map[string]any) (*big.Int, error) {
	var outputGas *big.Int
	for fieldName, fieldValue := range input {
		lowercase := strings.ToLower(fieldName)
		switch lowercase {
		case "gaslimit":
			if val, ok := fieldValue.(*big.Int); ok {
				outputGas = val
				return outputGas, nil
			} else {
				// when source chain is svm, the gas limit is an ag_binary.Uint128 struct instead of *big.Int
				if val, ok := fieldValue.(ag_binary.Uint128); ok {
					outputGas = val.BigInt()
					return outputGas, nil
				}
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

func SerializeExtraArgs(tag []byte, method string, inputs ...any) ([]byte, error) {
	v, err := messageHasherABI.Methods[method].Inputs.Pack(inputs...)
	return append(tag, v...), err
}

func SerializeEVMExtraArgsV1(data message_hasher.ClientEVMExtraArgsV1) ([]byte, error) {
	return SerializeExtraArgs(evmExtraArgsV1Tag, "encodeEVMExtraArgsV1", data)
}

func SerializeClientGenericExtraArgsV2(data message_hasher.ClientGenericExtraArgsV2) ([]byte, error) {
	return SerializeExtraArgs(evmExtraArgsV2Tag, "encodeGenericExtraArgsV2", data)
}

func SerializeClientSVMExtraArgsV1(data message_hasher.ClientSVMExtraArgsV1) ([]byte, error) {
	return SerializeExtraArgs(svmExtraArgsV1Tag, "encodeSVMExtraArgsV1", data)
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
