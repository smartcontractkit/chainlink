package storage

import (
	"errors"
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
	// maxQueuedTransactions is the max limit of UnstartedTransactions and ConfirmedTransactions structures.
	maxQueuedTransactions = 250
	// pruneSubset controls the subset of confirmed transactions to prune when the structure reaches its max limit.
	// i.e. if the value is 3 and the limit is 90, 30 transactions will be pruned.
	pruneSubset = 3
)

type InMemoryStore struct {
	sync.RWMutex
	lggr      logger.Logger
	txIDCount uint64
	address   common.Address
	chainID   *big.Int

	UnstartedTransactions   []*types.Transaction
	UnconfirmedTransactions map[uint64]*types.Transaction
	ConfirmedTransactions   map[uint64]*types.Transaction
	FatalTransactions       []*types.Transaction

	Transactions map[uint64]*types.Transaction
}

func NewInMemoryStore(lggr logger.Logger, address common.Address, chainID *big.Int) *InMemoryStore {
	return &InMemoryStore{
		lggr:                    logger.Named(lggr, "InMemoryStore"),
		address:                 address,
		chainID:                 chainID,
		UnconfirmedTransactions: make(map[uint64]*types.Transaction),
		ConfirmedTransactions:   make(map[uint64]*types.Transaction),
		Transactions:            make(map[uint64]*types.Transaction),
	}
}

func (m *InMemoryStore) AbandonPendingTransactions() {
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
}

func (m *InMemoryStore) AppendAttemptToTransaction(txNonce uint64, attempt *types.Attempt) error {
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

func (m *InMemoryStore) CountUnstartedTransactions() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.UnstartedTransactions)
}

func (m *InMemoryStore) CreateEmptyUnconfirmedTransaction(nonce uint64, gasLimit uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	emptyTx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           m.chainID,
		Nonce:             &nonce,
		FromAddress:       m.address,
		ToAddress:         common.Address{},
		Value:             big.NewInt(0),
		SpecifiedGasLimit: gasLimit,
		CreatedAt:         time.Now(),
		State:             types.TxUnconfirmed,
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	if _, exists := m.ConfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("a confirmed tx with the same nonce already exists: %v", m.ConfirmedTransactions[nonce])
	}

	m.UnconfirmedTransactions[nonce] = emptyTx
	m.Transactions[emptyTx.ID] = emptyTx

	return emptyTx.DeepCopy(), nil
}

func (m *InMemoryStore) CreateTransaction(txRequest *types.TxRequest) *types.Transaction {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++

	tx := &types.Transaction{
		ID:                m.txIDCount,
		IdempotencyKey:    txRequest.IdempotencyKey,
		ChainID:           m.chainID,
		FromAddress:       m.address,
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
		m.lggr.Warnf("Unstarted transactions queue for address: %v reached max limit of: %d. Dropping oldest transaction: %v.",
			m.address, maxQueuedTransactions, m.UnstartedTransactions[0])
		delete(m.Transactions, m.UnstartedTransactions[0].ID)
		m.UnstartedTransactions = m.UnstartedTransactions[1:maxQueuedTransactions]
	}

	txCopy := tx.DeepCopy()
	m.Transactions[txCopy.ID] = txCopy
	m.UnstartedTransactions = append(m.UnstartedTransactions, txCopy)
	return tx
}

func (m *InMemoryStore) FetchUnconfirmedTransactionAtNonceWithCount(latestNonce uint64) (txCopy *types.Transaction, unconfirmedCount int) {
	m.RLock()
	defer m.RUnlock()

	tx := m.UnconfirmedTransactions[latestNonce]
	if tx != nil {
		txCopy = tx.DeepCopy()
	}
	unconfirmedCount = len(m.UnconfirmedTransactions)
	return
}

func (m *InMemoryStore) MarkConfirmedAndReorgedTransactions(latestNonce uint64) ([]*types.Transaction, []uint64, error) {
	m.Lock()
	defer m.Unlock()

	var confirmedTransactions []*types.Transaction
	for _, tx := range m.UnconfirmedTransactions {
		if tx.Nonce == nil {
			return nil, nil, fmt.Errorf("nonce for txID: %v is empty", tx.ID)
		}
		if *tx.Nonce < latestNonce {
			tx.State = types.TxConfirmed
			confirmedTransactions = append(confirmedTransactions, tx.DeepCopy())
			m.ConfirmedTransactions[*tx.Nonce] = tx
			delete(m.UnconfirmedTransactions, *tx.Nonce)
		}
	}

	var unconfirmedTransactionIDs []uint64
	for _, tx := range m.ConfirmedTransactions {
		if tx.Nonce == nil {
			return nil, nil, fmt.Errorf("nonce for txID: %v is empty", tx.ID)
		}
		if *tx.Nonce >= latestNonce {
			tx.State = types.TxUnconfirmed
			tx.LastBroadcastAt = time.Time{} // Mark reorged transaction as if it wasn't broadcasted before
			unconfirmedTransactionIDs = append(unconfirmedTransactionIDs, tx.ID)
			m.UnconfirmedTransactions[*tx.Nonce] = tx
			delete(m.ConfirmedTransactions, *tx.Nonce)
		}
	}

	if len(m.ConfirmedTransactions) >= maxQueuedTransactions {
		prunedTxIDs := m.pruneConfirmedTransactions()
		m.lggr.Debugf("Confirmed transactions map for address: %v reached max limit of: %d. Pruned 1/3 of the oldest confirmed transactions. TxIDs: %v",
			m.address, maxQueuedTransactions, prunedTxIDs)
	}
	sort.Slice(confirmedTransactions, func(i, j int) bool { return confirmedTransactions[i].ID < confirmedTransactions[j].ID })
	sort.Slice(unconfirmedTransactionIDs, func(i, j int) bool { return unconfirmedTransactionIDs[i] < unconfirmedTransactionIDs[j] })
	return confirmedTransactions, unconfirmedTransactionIDs, nil
}

func (m *InMemoryStore) MarkUnconfirmedTransactionPurgeable(nonce uint64) error {
	m.Lock()
	defer m.Unlock()

	tx, exists := m.UnconfirmedTransactions[nonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx with nonce: %d was not found", nonce)
	}

	tx.IsPurgeable = true

	return nil
}

func (m *InMemoryStore) UpdateTransactionBroadcast(txID uint64, txNonce uint64, attemptHash common.Hash) error {
	m.Lock()
	defer m.Unlock()

	unconfirmedTx, exists := m.UnconfirmedTransactions[txNonce]
	if !exists {
		return fmt.Errorf("unconfirmed tx was not found for nonce: %d - txID: %v", txNonce, txID)
	}

	// Set the same time for both the tx and its attempt
	now := time.Now()
	unconfirmedTx.LastBroadcastAt = now
	if unconfirmedTx.InitialBroadcastAt.IsZero() {
		unconfirmedTx.InitialBroadcastAt = now
	}
	a, err := unconfirmedTx.FindAttemptByHash(attemptHash)
	if err != nil {
		return err
	}
	a.BroadcastAt = now

	return nil
}

func (m *InMemoryStore) UpdateUnstartedTransactionWithNonce(nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.UnstartedTransactions) == 0 {
		m.lggr.Debugf("Unstarted transactions queue is empty for address: %v", m.address)
		return nil, nil
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	tx := m.UnstartedTransactions[0]
	tx.Nonce = &nonce
	tx.State = types.TxUnconfirmed

	m.UnstartedTransactions = m.UnstartedTransactions[1:]
	m.UnconfirmedTransactions[nonce] = tx

	return tx.DeepCopy(), nil
}

// Shouldn't call lock because it's being called by a method that already has the lock
func (m *InMemoryStore) pruneConfirmedTransactions() []uint64 {
	noncesToPrune := make([]uint64, 0, len(m.ConfirmedTransactions))
	for nonce := range m.ConfirmedTransactions {
		noncesToPrune = append(noncesToPrune, nonce)
	}
	if len(noncesToPrune) == 0 {
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
func (m *InMemoryStore) DeleteAttemptForUnconfirmedTx(transactionNonce uint64, attempt *types.Attempt) error {
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

func (m *InMemoryStore) MarkTxFatal(*types.Transaction) error {
	return errors.New("not implemented")
}

// Orchestrator
func (m *InMemoryStore) FindTxWithIdempotencyKey(idempotencyKey *string) *types.Transaction {
	m.RLock()
	defer m.RUnlock()

	if idempotencyKey != nil {
		for _, tx := range m.Transactions {
			if tx.IdempotencyKey != nil && tx.IdempotencyKey == idempotencyKey {
				return tx.DeepCopy()
			}
		}
	}

	return nil
}
