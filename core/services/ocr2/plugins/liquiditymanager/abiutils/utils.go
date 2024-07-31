package abiutils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

// UnpackUint256 ABI decodes a single uint256 from the given data.
func UnpackUint256(data []byte) (*big.Int, error) {
	decoded, err := utils.ABIDecode(`[{"type": "uint256"}]`, data)
	if err != nil {
		return nil, err
	}

	if len(decoded) != 1 {
		return nil, fmt.Errorf("expected 1 element, got %d", len(decoded))
	}

	num := *abi.ConvertType(decoded[0], new(*big.Int)).(**big.Int)
	return num, nil
}
