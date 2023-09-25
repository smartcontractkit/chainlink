package hashlib

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type LeafHasherInterface[H Hash] interface {
	HashLeaf(log types.Log) (H, error)
}

var LeafDomainSeparator = [1]byte{0x00}

func GetMetaDataHash[H Hash](ctx Ctx[H], prefix [32]byte, sourceChainSelector uint64, onRampId common.Address, destChainSelector uint64) H {
	paddedOnRamp := onRampId.Hash()
	return ctx.Hash(utils.ConcatBytes(prefix[:], math.U256Bytes(big.NewInt(0).SetUint64(sourceChainSelector)), math.U256Bytes(big.NewInt(0).SetUint64(destChainSelector)), paddedOnRamp[:]))
}

type LeafHasher struct {
	metaDataHash [32]byte
	ctx          Ctx[[32]byte]
}

func NewLeafHasher(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx Ctx[[32]byte]) *LeafHasher {
	return &LeafHasher{
		metaDataHash: GetMetaDataHash(ctx, ctx.Hash([]byte("EVM2EVMMessageHashV2")), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
	}
}

var _ LeafHasherInterface[[32]byte] = &LeafHasher{}

func (t *LeafHasher) HashLeaf(log types.Log) ([32]byte, error) {
	message, err := abihelpers.DecodeOffRampMessage(log.Data)
	if err != nil {
		return [32]byte{}, err
	}

	encodedTokens, err := abihelpers.TokenAmountsArgs.PackValues([]interface{}{message.TokenAmounts})
	if err != nil {
		return [32]byte{}, err
	}

	bytesArray, err := abi.NewType("bytes[]", "bytes[]", nil)
	if err != nil {
		return [32]byte{}, err
	}

	encodedSourceTokenData, err := abi.Arguments{abi.Argument{Type: bytesArray}}.PackValues([]interface{}{message.SourceTokenData})
	if err != nil {
		return [32]byte{}, err
	}

	packedFixedSizeValues, err := utils.ABIEncode(
		`[
{"name": "sender", "type":"address"},
{"name": "receiver", "type":"address"},
{"name": "sequenceNumber", "type":"uint64"},
{"name": "gasLimit", "type":"uint256"},
{"name": "strict", "type":"bool"},
{"name": "nonce", "type":"uint64"},
{"name": "feeToken","type": "address"},
{"name": "feeTokenAmount","type": "uint256"}
]`,
		message.Sender,
		message.Receiver,
		message.SequenceNumber,
		message.GasLimit,
		message.Strict,
		message.Nonce,
		message.FeeToken,
		message.FeeTokenAmount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	fixedSizeValuesHash := t.ctx.Hash(packedFixedSizeValues)

	packedValues, err := utils.ABIEncode(
		`[
{"name": "leafDomainSeparator","type":"bytes1"},
{"name": "metadataHash", "type":"bytes32"},
{"name": "fixedSizeValuesHash", "type":"bytes32"},
{"name": "dataHash", "type":"bytes32"},
{"name": "tokenAmountsHash", "type":"bytes32"},
{"name": "sourceTokenDataHash", "type":"bytes32"}
]`,
		LeafDomainSeparator,
		t.metaDataHash,
		fixedSizeValuesHash,
		t.ctx.Hash(message.Data),
		t.ctx.Hash(encodedTokens),
		t.ctx.Hash(encodedSourceTokenData),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return t.ctx.Hash(packedValues), nil
}
