package mocks

import (
	"math/big"
	"testing"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type MockEvmTxManager = TxManager[*big.Int, *evmtypes.Head, evmtypes.Address, evmtypes.TxHash, evmtypes.BlockHash]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return NewTxManager[*big.Int, *evmtypes.Head, evmtypes.Address, evmtypes.TxHash, evmtypes.BlockHash](t)
}
