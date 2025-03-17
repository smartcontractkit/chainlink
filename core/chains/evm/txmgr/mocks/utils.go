package mocks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
)

type MockEvmTxManager = txmgrmocks.TxManager[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return txmgrmocks.NewTxManager[*big.Int, *types.Head, common.Address, common.Hash, common.Hash, types.Nonce, gas.EvmFee](t)
}
