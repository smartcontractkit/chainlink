package mocks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type MockEvmTxManager = TxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash]

func NewMockEvmTxManager(t *testing.T) *MockEvmTxManager {
	return NewTxManager[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash](t)
}

type MockEvmTxStore = mocks.TxStore[common.Address, *big.Int, common.Hash, common.Hash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce]

func NewMockEvmTxStore(t *testing.T) *MockEvmTxStore {
	return mocks.NewTxStore[common.Address, *big.Int, common.Hash, common.Hash, txmgr.EvmNewTx, *evmtypes.Receipt, txmgr.EvmTx, txmgr.EvmTxAttempt, evmtypes.Nonce](t)
}
