package storage

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

const StoreNotFoundForAddress string = "InMemoryStore for address: %v not found"

type InMemoryStoreManager struct {
	lggr             logger.Logger
	chainID          *big.Int
	InMemoryStoreMap map[common.Address]*InMemoryStore
}

func NewInMemoryStoreManager(lggr logger.Logger, chainID *big.Int) *InMemoryStoreManager {
	inMemoryStoreMap := make(map[common.Address]*InMemoryStore)
	return &InMemoryStoreManager{
		lggr:             lggr,
		chainID:          chainID,
		InMemoryStoreMap: inMemoryStoreMap}
}

func (m *InMemoryStoreManager) AbandonPendingTransactions(_ context.Context, fromAddress common.Address) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		store.AbandonPendingTransactions()
		return nil
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) Add(addresses ...common.Address) (err error) {
	for _, address := range addresses {
		if _, exists := m.InMemoryStoreMap[address]; exists {
			err = errors.Join(err, fmt.Errorf("address %v already exists in store manager", address))
		}
		m.InMemoryStoreMap[address] = NewInMemoryStore(m.lggr, address, m.chainID)
	}
	return
}

func (m *InMemoryStoreManager) AppendAttemptToTransaction(_ context.Context, txNonce uint64, fromAddress common.Address, attempt *types.Attempt) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.AppendAttemptToTransaction(txNonce, attempt)
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) CountUnstartedTransactions(fromAddress common.Address) (int, error) {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.CountUnstartedTransactions(), nil
	}
	return 0, fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) CreateEmptyUnconfirmedTransaction(_ context.Context, fromAddress common.Address, nonce uint64, gasLimit uint64) (*types.Transaction, error) {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.CreateEmptyUnconfirmedTransaction(nonce, gasLimit)
	}
	return nil, fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) CreateTransaction(_ context.Context, txRequest *types.TxRequest) (*types.Transaction, error) {
	if store, exists := m.InMemoryStoreMap[txRequest.FromAddress]; exists {
		return store.CreateTransaction(txRequest), nil
	}
	return nil, fmt.Errorf(StoreNotFoundForAddress, txRequest.FromAddress)
}

func (m *InMemoryStoreManager) FetchUnconfirmedTransactionAtNonceWithCount(_ context.Context, nonce uint64, fromAddress common.Address) (tx *types.Transaction, count int, err error) {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		tx, count = store.FetchUnconfirmedTransactionAtNonceWithCount(nonce)
		return
	}
	return nil, 0, fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) MarkConfirmedAndReorgedTransactions(_ context.Context, nonce uint64, fromAddress common.Address) (confirmedTxs []*types.Transaction, unconfirmedTxIDs []uint64, err error) {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		confirmedTxs, unconfirmedTxIDs, err = store.MarkConfirmedAndReorgedTransactions(nonce)
		return
	}
	return nil, nil, fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) MarkUnconfirmedTransactionPurgeable(_ context.Context, nonce uint64, fromAddress common.Address) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.MarkUnconfirmedTransactionPurgeable(nonce)
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) UpdateTransactionBroadcast(_ context.Context, txID uint64, nonce uint64, attemptHash common.Hash, fromAddress common.Address) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.UpdateTransactionBroadcast(txID, nonce, attemptHash)
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) UpdateUnstartedTransactionWithNonce(_ context.Context, fromAddress common.Address, nonce uint64) (*types.Transaction, error) {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.UpdateUnstartedTransactionWithNonce(nonce)
	}
	return nil, fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) DeleteAttemptForUnconfirmedTx(_ context.Context, nonce uint64, attempt *types.Attempt, fromAddress common.Address) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.DeleteAttemptForUnconfirmedTx(nonce, attempt)
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) MarkTxFatal(_ context.Context, tx *types.Transaction, fromAddress common.Address) error {
	if store, exists := m.InMemoryStoreMap[fromAddress]; exists {
		return store.MarkTxFatal(tx)
	}
	return fmt.Errorf(StoreNotFoundForAddress, fromAddress)
}

func (m *InMemoryStoreManager) FindTxWithIdempotencyKey(_ context.Context, idempotencyKey string) (*types.Transaction, error) {
	for _, store := range m.InMemoryStoreMap {
		tx := store.FindTxWithIdempotencyKey(idempotencyKey)
		if tx != nil {
			return tx, nil
		}
	}
	return nil, nil
}
