package mocks

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type MockEvmTxStore = TxStore[evmtypes.Address, *big.Int, evmtypes.TxHash, evmtypes.BlockHash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce]

func NewMockEvmTxStore(t *testing.T) *MockEvmTxStore {
	return NewTxStore[evmtypes.Address, *big.Int, evmtypes.TxHash, evmtypes.BlockHash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce](t)
}
