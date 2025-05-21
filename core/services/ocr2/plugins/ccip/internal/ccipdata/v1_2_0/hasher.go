package v1_2_0

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"

	evm_2_evm_onramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	MetaDataHashPrefix = "EVM2EVMMessageHashV2"
)

var LeafDomainSeparator = [1]byte{0x00}

type LeafHasher struct {
	metaDataHash [32]byte
	ctx          hashutil.Hasher[[32]byte]
	onRamp       *evm_2_evm_onramp_1_2_0.EVM2EVMOnRamp
}

func GetMetaDataHash[H hashutil.Hash](ctx hashutil.Hasher[H], prefix [32]byte, sourceChainSelector uint64, onRampID common.Address, destChainSelector uint64) H {
	paddedOnRamp := common.BytesToHash(onRampID[:])
	return ctx.Hash(utils.ConcatBytes(prefix[:], math.U256Bytes(big.NewInt(0).SetUint64(sourceChainSelector)), math.U256Bytes(big.NewInt(0).SetUint64(destChainSelector)), paddedOnRamp[:]))
}

func NewLeafHasher(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashutil.Hasher[[32]byte], onRamp *evm_2_evm_onramp_1_2_0.EVM2EVMOnRamp) *LeafHasher {
	return &LeafHasher{
		metaDataHash: GetMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefix)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasher) HashLeaf(log types.Log) ([32]byte, error) {
	msg, err := t.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return [32]byte{}, err
	}
	message := msg.Message
	encodedTokens, err := abihelpers.ABIEncode(
		`[
{"components": [{"name":"token","type":"address"},{"name":"amount","type":"uint256"}], "type":"tuple[]"}]`, message.TokenAmounts)
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

	packedFixedSizeValues, err := abihelpers.ABIEncode(
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

	packedValues, err := abihelpers.ABIEncode(
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
