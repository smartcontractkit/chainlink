package mocks

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type MockEvmTxManager = TxManager[*big.Int, *evmtypes.Head, evmtypes.Address, evmtypes.TxHash, evmtypes.BlockHash]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return NewTxManager[*big.Int, *evmtypes.Head, evmtypes.Address, evmtypes.TxHash, evmtypes.BlockHash](t)
}

type MockEvmTxStore = mocks.TxStore[evmtypes.Address, *big.Int, evmtypes.TxHash, evmtypes.BlockHash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce]

func NewMockEvmTxStore(t *testing.T) *MockEvmTxStore {
	return mocks.NewTxStore[evmtypes.Address, *big.Int, evmtypes.TxHash, evmtypes.BlockHash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce](t)
}
