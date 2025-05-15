package mocks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

type MockEvmTxManager = TxManager[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return NewTxManager[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee](t)
}
