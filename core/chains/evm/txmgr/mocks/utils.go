package mocks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type MockEvmTxManager = txmgrmocks.TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return txmgrmocks.NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee](t)
}
