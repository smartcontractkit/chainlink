package storage

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type InMemoryStore struct {
	sync.RWMutex
	lggr      logger.Logger
	txIDCount uint64

	UnstartedTransactions   []*types.Transaction
	UnconfirmedTransactions map[uint64]*types.Transaction
	ConfirmedTransactions   map[uint64]*types.Transaction
	FatalTransactions       []*types.Transaction
}

func NewInMemoryStore(lggr logger.Logger) *InMemoryStore {
	return &InMemoryStore{
		lggr:                    logger.Named(lggr, "InMemoryStore"),
		UnconfirmedTransactions: make(map[uint64]*types.Transaction),
		ConfirmedTransactions:   make(map[uint64]*types.Transaction),
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

	return emptyTx.DeepCopy(), nil
}

func (m *InMemoryStore) CreateTransaction(_ context.Context, tx *types.Transaction) (uint64, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++

	tx.ID = m.txIDCount
	tx.CreatedAt = time.Now()
	tx.State = types.TxUnstarted

	m.UnstartedTransactions = append(m.UnstartedTransactions, tx.DeepCopy())
	return tx.ID, nil
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
			unconfirmedTransactionIDs = append(unconfirmedTransactionIDs, tx.ID)
			m.UnconfirmedTransactions[tx.Nonce] = tx
			delete(m.ConfirmedTransactions, tx.Nonce)
		}
	}
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
