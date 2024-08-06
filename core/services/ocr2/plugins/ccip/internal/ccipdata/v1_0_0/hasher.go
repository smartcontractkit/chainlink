package v1_0_0

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

const (
	MetaDataHashPrefix = "EVM2EVMMessageEvent"
)

var LeafDomainSeparator = [1]byte{0x00}

type LeafHasher struct {
	metaDataHash [32]byte
	ctx          hashutil.Hasher[[32]byte]
	onRamp       *evm_2_evm_onramp_1_0_0.EVM2EVMOnRamp
}

func GetMetaDataHash[H hashutil.Hash](ctx hashutil.Hasher[H], prefix [32]byte, sourceChainSelector uint64, onRampId common.Address, destChainSelector uint64) H {
	paddedOnRamp := common.BytesToHash(onRampId[:])
	return ctx.Hash(utils.ConcatBytes(prefix[:], math.U256Bytes(big.NewInt(0).SetUint64(sourceChainSelector)), math.U256Bytes(big.NewInt(0).SetUint64(destChainSelector)), paddedOnRamp[:]))
}

func NewLeafHasher(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashutil.Hasher[[32]byte], onRamp *evm_2_evm_onramp_1_0_0.EVM2EVMOnRamp) *LeafHasher {
	return &LeafHasher{
		metaDataHash: GetMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefix)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasher) HashLeaf(log types.Log) ([32]byte, error) {
	message, err := t.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return [32]byte{}, err
	}
	encodedTokens, err := abihelpers.ABIEncode(
		`[
{"components": [{"name":"token","type":"address"},{"name":"amount","type":"uint256"}], "type":"tuple[]"}]`, message.Message.TokenAmounts)
	if err != nil {
		return [32]byte{}, err
	}

	packedValues, err := abihelpers.ABIEncode(
		`[
{"name": "leafDomainSeparator","type":"bytes1"},
{"name": "metadataHash", "type":"bytes32"},
{"name": "sequenceNumber", "type":"uint64"},
{"name": "nonce", "type":"uint64"},
{"name": "sender", "type":"address"},
{"name": "receiver", "type":"address"},
{"name": "dataHash", "type":"bytes32"},
{"name": "tokenAmountsHash", "type":"bytes32"},
{"name": "gasLimit", "type":"uint256"},
{"name": "strict", "type":"bool"},
{"name": "feeToken","type": "address"},
{"name": "feeTokenAmount","type": "uint256"}
]`,
		LeafDomainSeparator,
		t.metaDataHash,
		message.Message.SequenceNumber,
		message.Message.Nonce,
		message.Message.Sender,
		message.Message.Receiver,
		t.ctx.Hash(message.Message.Data),
		t.ctx.Hash(encodedTokens),
		message.Message.GasLimit,
		message.Message.Strict,
		message.Message.FeeToken,
		message.Message.FeeTokenAmount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return t.ctx.Hash(packedValues), nil
}
