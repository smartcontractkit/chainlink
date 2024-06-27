package ccipevm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "EVM2EVMMultiOnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	metaDataHash        [32]byte
	leafDomainSeparator [32]byte
}

func NewMessageHasherV1(metaDataHash [32]byte) *MessageHasherV1 {
	return &MessageHasherV1{
		metaDataHash:        metaDataHash,
		leafDomainSeparator: [32]byte{},
	}
}

func (h *MessageHasherV1) Hash(_ context.Context, msg cciptypes.CCIPMsg) (cciptypes.Bytes32, error) {
	type tokenAmount struct {
		Token  common.Address
		Amount *big.Int
	}
	tokenAmounts := make([]tokenAmount, len(msg.TokenAmounts))
	for i, ta := range msg.TokenAmounts {
		tokenAmounts[i] = tokenAmount{
			Token:  common.HexToAddress(string(ta.Token)),
			Amount: ta.Amount,
		}
	}
	encodedTokens, err := h.abiEncode(`[{"components": [{"name":"token","type":"address"},
		{"name":"amount","type":"uint256"}], "type":"tuple[]"}]`, tokenAmounts)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode token amounts: %w", err)
	}

	bytesArray, err := abi.NewType("bytes[]", "bytes[]", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("new abi type bytes[]: %w", err)
	}

	encodedSourceTokenData, err := abi.Arguments{abi.Argument{Type: bytesArray}}.PackValues([]interface{}{msg.SourceTokenData})
	if err != nil {
		return [32]byte{}, fmt.Errorf("pack source token data: %w", err)
	}

	packedFixedSizeValues, err := h.abiEncode(
		`[{"name": "sender", "type":"address"},
			{"name": "receiver", "type":"address"},
			{"name": "sequenceNumber", "type":"uint64"},
			{"name": "gasLimit", "type":"uint256"},
			{"name": "strict", "type":"bool"},
			{"name": "nonce", "type":"uint64"},
			{"name": "feeToken","type": "address"},
			{"name": "feeTokenAmount","type": "uint256"}]`,
		common.HexToAddress(string(msg.Sender)),
		common.HexToAddress(string(msg.Receiver)),
		uint64(msg.SeqNum),
		msg.ChainFeeLimit.Int,
		msg.Strict,
		msg.Nonce,
		common.HexToAddress(string(msg.FeeToken)),
		msg.FeeTokenAmount.Int,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode fixed size values: %w", err)
	}
	fixedSizeValuesHash := utils.Keccak256Fixed(packedFixedSizeValues)

	packedValues, err := h.abiEncode(
		`[{"name": "leafDomainSeparator","type":"bytes32"},
			{"name": "metadataHash", "type":"bytes32"},
			{"name": "fixedSizeValuesHash", "type":"bytes32"},
			{"name": "dataHash", "type":"bytes32"},
			{"name": "tokenAmountsHash", "type":"bytes32"},
			{"name": "sourceTokenDataHash", "type":"bytes32"}]`,
		h.leafDomainSeparator,
		h.metaDataHash,
		fixedSizeValuesHash,
		utils.Keccak256Fixed(msg.Data),
		utils.Keccak256Fixed(encodedTokens),
		utils.Keccak256Fixed(encodedSourceTokenData),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("abi encode packed values: %w", err)
	}

	return utils.Keccak256Fixed(packedValues), nil
}

func (h *MessageHasherV1) abiEncode(abiStr string, values ...interface{}) ([]byte, error) {
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
