package storage

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

const (
	maxQueuedTransactions = 250
	pruneSubset           = 3
)

type InMemoryStore struct {
	sync.RWMutex
	lggr      logger.Logger
	txIDCount uint64

	UnstartedTransactions   []*types.Transaction
	UnconfirmedTransactions map[uint64]*types.Transaction
	ConfirmedTransactions   map[uint64]*types.Transaction
	FatalTransactions       []*types.Transaction

	Transactions map[uint64]*types.Transaction
}

func NewInMemoryStore(lggr logger.Logger) *InMemoryStore {
	return &InMemoryStore{
		lggr:                    logger.Named(lggr, "InMemoryStore"),
		UnconfirmedTransactions: make(map[uint64]*types.Transaction),
		ConfirmedTransactions:   make(map[uint64]*types.Transaction),
		Transactions:            make(map[uint64]*types.Transaction),
	}
}

func (m *InMemoryStore) AbandonPendingTransactions(context.Context, common.Address) error {
	m.Lock()
	defer m.Unlock()

	for _, tx := range m.UnstartedTransactions {
		tx.State = types.TxFatalError
	}
	m.FatalTransactions = m.UnstartedTransactions
	m.UnstartedTransactions = []*types.Transaction{}

	for _, tx := range m.UnconfirmedTransactions {
		tx.State = types.TxFatalError
		m.FatalTransactions = append(m.FatalTransactions, tx)
	}
	m.UnconfirmedTransactions = make(map[uint64]*types.Transaction)

	return nil
}

func (m *InMemoryStore) AppendAttemptToTransaction(_ context.Context, txNonce uint64, attempt *types.Attempt) error {
	m.Lock()
	defer m.Unlock()

	tx, exists := m.UnconfirmedTransactions[txNonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx was not found for nonce: %d - txID: %v", txNonce, attempt.TxID)
	}

	if tx.ID != attempt.TxID {
		return fmt.Errorf("unconfirmed tx with nonce exists but attempt points to a different txID. Found Tx: %v - txID: %v", m.UnconfirmedTransactions[txNonce], attempt.TxID)
	}

	attempt.CreatedAt = time.Now()
	attempt.ID = uint64(len(tx.Attempts)) // Attempts are not collectively tracked by the in-memory store so attemptIDs are not unique between transactions and can be reused.
	tx.AttemptCount++
	m.UnconfirmedTransactions[txNonce].Attempts = append(m.UnconfirmedTransactions[txNonce].Attempts, attempt.DeepCopy())

	return nil
}

func (m *InMemoryStore) CountUnstartedTransactions(context.Context, common.Address) (int, error) {
	m.RLock()
	defer m.RUnlock()

	return len(m.UnstartedTransactions), nil
}

func (m *InMemoryStore) CreateEmptyUnconfirmedTransaction(ctx context.Context, fromAddress common.Address, chainID *big.Int, nonce uint64, limit uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	emptyTx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           chainID,
		Nonce:             nonce,
		FromAddress:       fromAddress,
		ToAddress:         common.Address{},
		Value:             big.NewInt(0),
		SpecifiedGasLimit: limit,
		CreatedAt:         time.Now(),
		State:             types.TxUnconfirmed,
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	m.UnconfirmedTransactions[nonce] = emptyTx
	m.Transactions[emptyTx.ID] = emptyTx

	return emptyTx.DeepCopy(), nil
}

func (m *InMemoryStore) CreateTransaction(_ context.Context, txRequest *types.TxRequest) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++

	tx := &types.Transaction{
		ID:                m.txIDCount,
		IdempotencyKey:    txRequest.IdempotencyKey,
		ChainID:           txRequest.ChainID,
		FromAddress:       txRequest.FromAddress,
		ToAddress:         txRequest.ToAddress,
		Value:             txRequest.Value,
		Data:              txRequest.Data,
		SpecifiedGasLimit: txRequest.SpecifiedGasLimit,
		CreatedAt:         time.Now(),
		State:             types.TxUnstarted,
		Meta:              txRequest.Meta,
		MinConfirmations:  txRequest.MinConfirmations,
		PipelineTaskRunID: txRequest.PipelineTaskRunID,
		SignalCallback:    txRequest.SignalCallback,
	}

	if len(m.UnstartedTransactions) == maxQueuedTransactions {
		m.lggr.Warnf("Unstarted transactions queue reached max limit of: %d. Dropping oldest transaction: %v.",
			maxQueuedTransactions, m.UnstartedTransactions[0])
		delete(m.Transactions, m.UnstartedTransactions[0].ID)
		m.UnstartedTransactions = m.UnstartedTransactions[1:maxQueuedTransactions]
	}

	copy := tx.DeepCopy()
	m.Transactions[copy.ID] = copy
	m.UnstartedTransactions = append(m.UnstartedTransactions, copy)
	return tx, nil
}

func (m *InMemoryStore) FetchUnconfirmedTransactionAtNonceWithCount(_ context.Context, latestNonce uint64, _ common.Address) (txCopy *types.Transaction, unconfirmedCount int, err error) {
	m.RLock()
	defer m.RUnlock()

	tx := m.UnconfirmedTransactions[latestNonce]
	if tx != nil {
		txCopy = tx.DeepCopy()
	}
	unconfirmedCount = len(m.UnconfirmedTransactions)
	return
}

func (m *InMemoryStore) MarkTransactionsConfirmed(_ context.Context, latestNonce uint64, _ common.Address) ([]uint64, []uint64, error) {
	m.Lock()
	defer m.Unlock()

	var confirmedTransactionIDs []uint64
	for _, tx := range m.UnconfirmedTransactions {
		if tx.Nonce < latestNonce {
			tx.State = types.TxConfirmed
			confirmedTransactionIDs = append(confirmedTransactionIDs, tx.ID)
			m.ConfirmedTransactions[tx.Nonce] = tx
			delete(m.UnconfirmedTransactions, tx.Nonce)
		}
	}

	var unconfirmedTransactionIDs []uint64
	for _, tx := range m.ConfirmedTransactions {
		if tx.Nonce >= latestNonce {
			tx.State = types.TxUnconfirmed
			tx.LastBroadcastAt = time.Time{} // Mark reorged transaction as if it wasn't broadcasted before
			unconfirmedTransactionIDs = append(unconfirmedTransactionIDs, tx.ID)
			m.UnconfirmedTransactions[tx.Nonce] = tx
			delete(m.ConfirmedTransactions, tx.Nonce)
		}
	}

	if len(m.ConfirmedTransactions) >= maxQueuedTransactions {
		prunedTxIDs := m.pruneConfirmedTransactions()
		m.lggr.Debugf("Confirmed transactions map reached max limit of: %d. Pruned 1/3 of the oldest confirmed transactions. TxIDs: %v", maxQueuedTransactions, prunedTxIDs)
	}
	sort.Slice(confirmedTransactionIDs, func(i, j int) bool { return confirmedTransactionIDs[i] < confirmedTransactionIDs[j] })
	sort.Slice(unconfirmedTransactionIDs, func(i, j int) bool { return unconfirmedTransactionIDs[i] < unconfirmedTransactionIDs[j] })
	return confirmedTransactionIDs, unconfirmedTransactionIDs, nil
}

func (m *InMemoryStore) MarkUnconfirmedTransactionPurgeable(_ context.Context, nonce uint64) error {
	m.Lock()
	defer m.Unlock()

	tx, exists := m.UnconfirmedTransactions[nonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx with nonce: %d was not found", nonce)
	}

	tx.IsPurgeable = true

	return nil
}

func (m *InMemoryStore) UpdateTransactionBroadcast(_ context.Context, txID uint64, txNonce uint64, attemptHash common.Hash) error {
	m.Lock()
	defer m.Unlock()

	unconfirmedTx, exists := m.UnconfirmedTransactions[txNonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx was not found for nonce: %d - txID: %v", txNonce, txID)
	}

	// Set the same time for both the tx and its attempt
	now := time.Now()
	unconfirmedTx.LastBroadcastAt = now
	a, err := unconfirmedTx.FindAttemptByHash(attemptHash)
	if err != nil {
		return err
	}
	a.BroadcastAt = now

	return nil
}

func (m *InMemoryStore) UpdateUnstartedTransactionWithNonce(_ context.Context, _ common.Address, nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.UnstartedTransactions) == 0 {
		m.lggr.Debug("Unstarted transaction queue is empty")
		return nil, nil
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	tx := m.UnstartedTransactions[0]
	tx.Nonce = nonce
	tx.State = types.TxUnconfirmed

	m.UnstartedTransactions = m.UnstartedTransactions[1:]
	m.UnconfirmedTransactions[nonce] = tx

	return tx.DeepCopy(), nil
}

// Shouldn't call lock because it's being called by a method that already has the lock
func (m *InMemoryStore) pruneConfirmedTransactions() []uint64 {
	var noncesToPrune []uint64
	for nonce := range m.ConfirmedTransactions {
		noncesToPrune = append(noncesToPrune, nonce)
	}
	if len(noncesToPrune) <= 0 {
		return nil
	}
	sort.Slice(noncesToPrune, func(i, j int) bool { return noncesToPrune[i] < noncesToPrune[j] })
	minNonce := noncesToPrune[len(noncesToPrune)/pruneSubset]

	var txIDsToPrune []uint64
	for nonce, tx := range m.ConfirmedTransactions {
		if nonce < minNonce {
			txIDsToPrune = append(txIDsToPrune, tx.ID)
			delete(m.Transactions, tx.ID)
			delete(m.ConfirmedTransactions, nonce)
		}
	}

	sort.Slice(txIDsToPrune, func(i, j int) bool { return txIDsToPrune[i] < txIDsToPrune[j] })
	return txIDsToPrune
}

// Error Handler
func (m *InMemoryStore) DeleteAttemptForUnconfirmedTx(_ context.Context, transactionNonce uint64, attempt *types.Attempt) error {
	m.Lock()
	defer m.Unlock()

	tx, exists := m.UnconfirmedTransactions[transactionNonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx was not found for nonce: %d - txID: %v", transactionNonce, attempt.TxID)
	}

	for i, a := range tx.Attempts {
		if a.Hash == attempt.Hash {
			tx.Attempts = append(tx.Attempts[:i], tx.Attempts[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("attempt with hash: %v for txID: %v was not found", attempt.Hash, attempt.TxID)
}

func (m *InMemoryStore) MarkTxFatal(context.Context, *types.Transaction) error {
	return fmt.Errorf("not implemented")
}

// Orchestrator
func (m *InMemoryStore) FindTxWithIdempotencyKey(_ context.Context, idempotencyKey *string) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	if idempotencyKey != nil {
		for _, tx := range m.Transactions {
			if tx.IdempotencyKey != nil && tx.IdempotencyKey == idempotencyKey {
				return tx.DeepCopy(), nil
			}
		}
	}

	return nil, nil
}
