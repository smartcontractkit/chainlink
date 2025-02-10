package ccipevm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-ccip/pkg/logutil"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
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
// - "OnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	lggr logger.Logger
}

func NewMessageHasherV1(lggr logger.Logger) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr: lggr,
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

	var rampTokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
	for _, rta := range msg.TokenAmounts {
		destGasAmount, err := abiDecodeUint32(rta.DestExecData)
		if err != nil {
			return [32]byte{}, fmt.Errorf("decode dest gas amount: %w", err)
		}

		lggr.Debugw("decoded dest gas amount",
			"destGasAmount", destGasAmount)

		// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
		// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
		// its unclear as of writing how we will handle non-EVM chains and their addresses.
		// e.g, will we encode them as bytes or bytes32?
		sourcePoolAddressABIEncodedAsAddress, err := abiEncodeAddress(common.BytesToAddress(rta.SourcePoolAddress))
		if err != nil {
			return [32]byte{}, fmt.Errorf("abi encode source pool address: %w", err)
		}

		lggr.Debugw("abi encoded source pool address as solidity address",
			"sourcePoolAddressABIEncodedAsAddress", hexutil.Encode(sourcePoolAddressABIEncodedAsAddress))

		destTokenAddress, err := abiDecodeAddress(rta.DestTokenAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf("decode dest token address: %w", err)
		}

		lggr.Debugw("abi decoded dest token address",
			"destTokenAddress", destTokenAddress)

		rampTokenAmounts = append(rampTokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: sourcePoolAddressABIEncodedAsAddress,
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
	gasLimit, err := decodeExtraArgsV1V2(msg.ExtraArgs)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode extra args: %w", err)
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

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
