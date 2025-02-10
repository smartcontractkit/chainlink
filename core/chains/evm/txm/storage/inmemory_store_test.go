package storage

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

func TestAbandonPendingTransactions(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	t.Run("abandons unstarted and unconfirmed transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		// Unstarted
		tx1 := insertUnstartedTransaction(m)
		tx2 := insertUnstartedTransaction(m)

		// Unconfirmed
		tx3, err := insertUnconfirmedTransaction(m, 3)
		require.NoError(t, err)
		tx4, err := insertUnconfirmedTransaction(m, 4)
		require.NoError(t, err)

		m.AbandonPendingTransactions()

		assert.Equal(t, txmgr.TxFatalError, tx1.State)
		assert.Equal(t, txmgr.TxFatalError, tx2.State)
		assert.Equal(t, txmgr.TxFatalError, tx3.State)
		assert.Equal(t, txmgr.TxFatalError, tx4.State)
	})

	t.Run("skips all types apart from unstarted and unconfirmed transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		// Fatal
		tx1 := insertFataTransaction(m)
		tx2 := insertFataTransaction(m)

		// Confirmed
		tx3, err := insertConfirmedTransaction(m, 3)
		require.NoError(t, err)
		tx4, err := insertConfirmedTransaction(m, 4)
		require.NoError(t, err)

		m.AbandonPendingTransactions()

		assert.Equal(t, txmgr.TxFatalError, tx1.State)
		assert.Equal(t, txmgr.TxFatalError, tx2.State)
		assert.Equal(t, txmgr.TxConfirmed, tx3.State)
		assert.Equal(t, txmgr.TxConfirmed, tx4.State)
		assert.Len(t, m.Transactions, 2) // tx1, tx2 were dropped
	})
}

func TestAppendAttemptToTransaction(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)

	_, err := insertUnconfirmedTransaction(m, 10) // txID = 1, nonce = 10
	require.NoError(t, err)
	_, err = insertConfirmedTransaction(m, 2) // txID = 2, nonce = 2
	require.NoError(t, err)

	t.Run("fails if corresponding unconfirmed transaction for attempt was not found", func(t *testing.T) {
		var nonce uint64 = 1
		newAttempt := &types.Attempt{}
		err := m.AppendAttemptToTransaction(nonce, newAttempt)
		require.Error(t, err)
		require.ErrorContains(t, err, "unconfirmed tx was not found")
	})

	t.Run("fails if unconfirmed transaction was found but doesn't match the txID", func(t *testing.T) {
		var nonce uint64 = 10
		newAttempt := &types.Attempt{
			TxID: 2,
		}
		err := m.AppendAttemptToTransaction(nonce, newAttempt)
		require.Error(t, err)
		require.ErrorContains(t, err, "attempt points to a different txID")
	})

	t.Run("appends attempt to transaction", func(t *testing.T) {
		var nonce uint64 = 10
		newAttempt := &types.Attempt{
			TxID: 1,
		}
		require.NoError(t, m.AppendAttemptToTransaction(nonce, newAttempt))
		tx, _ := m.FetchUnconfirmedTransactionAtNonceWithCount(10)
		assert.Len(t, tx.Attempts, 1)
		assert.Equal(t, uint16(1), tx.AttemptCount)
		assert.False(t, tx.Attempts[0].CreatedAt.IsZero())
	})
}

func TestCountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)

	assert.Equal(t, 0, m.CountUnstartedTransactions())

	insertUnstartedTransaction(m)
	assert.Equal(t, 1, m.CountUnstartedTransactions())

	_, err := insertConfirmedTransaction(m, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, m.CountUnstartedTransactions())
}

func TestCreateEmptyUnconfirmedTransaction(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
	_, err := insertUnconfirmedTransaction(m, 1)
	require.NoError(t, err)
	_, err = insertConfirmedTransaction(m, 0)
	require.NoError(t, err)

	t.Run("fails if unconfirmed transaction with the same nonce exists", func(t *testing.T) {
		_, err := m.CreateEmptyUnconfirmedTransaction(1, 0)
		require.Error(t, err)
	})

	t.Run("fails if confirmed transaction with the same nonce exists", func(t *testing.T) {
		_, err := m.CreateEmptyUnconfirmedTransaction(0, 0)
		require.Error(t, err)
	})

	t.Run("creates a new empty unconfirmed transaction", func(t *testing.T) {
		tx, err := m.CreateEmptyUnconfirmedTransaction(2, 0)
		require.NoError(t, err)
		assert.Equal(t, txmgr.TxUnconfirmed, tx.State)
	})
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()

	t.Run("creates new transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		now := time.Now()
		txR1 := &types.TxRequest{}
		txR2 := &types.TxRequest{}
		tx1 := m.CreateTransaction(txR1)
		assert.Equal(t, uint64(0), tx1.ID)
		assert.LessOrEqual(t, now, tx1.CreatedAt)

		tx2 := m.CreateTransaction(txR2)
		assert.Equal(t, uint64(1), tx2.ID)
		assert.LessOrEqual(t, now, tx2.CreatedAt)

		assert.Equal(t, 2, m.CountUnstartedTransactions())
	})

	t.Run("prunes oldest unstarted transactions if limit is reached", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		overshot := 5
		for i := 0; i < maxQueuedTransactions+overshot; i++ {
			r := &types.TxRequest{}
			tx := m.CreateTransaction(r)
			//nolint:gosec // this won't overflow
			assert.Equal(t, uint64(i), tx.ID)
		}
		// total shouldn't exceed maxQueuedTransactions
		assert.Equal(t, maxQueuedTransactions, m.CountUnstartedTransactions())
		// earliest tx ID should be the same amount of the number of transactions that we dropped
		tx, err := m.UpdateUnstartedTransactionWithNonce(0)
		require.NoError(t, err)
		//nolint:gosec // this won't overflow
		assert.Equal(t, uint64(overshot), tx.ID)
	})
}

func TestFetchUnconfirmedTransactionAtNonceWithCount(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)

	tx, count := m.FetchUnconfirmedTransactionAtNonceWithCount(0)
	assert.Nil(t, tx)
	assert.Equal(t, 0, count)

	var nonce uint64
	_, err := insertUnconfirmedTransaction(m, nonce)
	require.NoError(t, err)
	tx, count = m.FetchUnconfirmedTransactionAtNonceWithCount(0)
	assert.Equal(t, *tx.Nonce, nonce)
	assert.Equal(t, 1, count)
}

func TestMarkConfirmedAndReorgedTransactions(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()

	t.Run("returns 0 if there are no transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		un, cn, err := m.MarkConfirmedAndReorgedTransactions(100)
		require.NoError(t, err)
		assert.Empty(t, un)
		assert.Empty(t, cn)
	})

	t.Run("confirms transaction with nonce lower than the latest", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		ctx1, err := insertUnconfirmedTransaction(m, 0)
		require.NoError(t, err)

		ctx2, err := insertUnconfirmedTransaction(m, 1)
		require.NoError(t, err)

		ctxs, utxs, err := m.MarkConfirmedAndReorgedTransactions(1)
		require.NoError(t, err)
		assert.Equal(t, txmgr.TxConfirmed, ctx1.State)
		assert.Equal(t, txmgr.TxUnconfirmed, ctx2.State)
		assert.Equal(t, ctxs[0].ID, ctx1.ID) // Ensure order
		assert.Empty(t, utxs)
	})

	t.Run("state remains the same if nonce didn't change", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		ctx1, err := insertConfirmedTransaction(m, 0)
		require.NoError(t, err)

		ctx2, err := insertUnconfirmedTransaction(m, 1)
		require.NoError(t, err)

		ctxs, utxs, err := m.MarkConfirmedAndReorgedTransactions(1)
		require.NoError(t, err)
		assert.Equal(t, txmgr.TxConfirmed, ctx1.State)
		assert.Equal(t, txmgr.TxUnconfirmed, ctx2.State)
		assert.Empty(t, ctxs)
		assert.Empty(t, utxs)
	})

	t.Run("unconfirms transaction with nonce equal to or higher than the latest", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		ctx1, err := insertConfirmedTransaction(m, 0)
		require.NoError(t, err)

		ctx2, err := insertConfirmedTransaction(m, 1)
		require.NoError(t, err)

		ctxs, utxs, err := m.MarkConfirmedAndReorgedTransactions(1)
		require.NoError(t, err)
		assert.Equal(t, txmgr.TxConfirmed, ctx1.State)
		assert.Equal(t, txmgr.TxUnconfirmed, ctx2.State)
		assert.Equal(t, utxs[0], ctx2.ID)
		assert.Empty(t, ctxs)
	})

	t.Run("logs an error during confirmation if a transaction with the same nonce already exists", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		m := NewInMemoryStore(lggr, fromAddress, testutils.FixtureChainID)
		_, err := insertConfirmedTransaction(m, 0)
		require.NoError(t, err)
		_, err = insertUnconfirmedTransaction(m, 0)
		require.NoError(t, err)

		_, _, err = m.MarkConfirmedAndReorgedTransactions(1)
		require.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, "Another confirmed transaction with the same nonce exists")
	})

	t.Run("prunes confirmed transactions map if it reaches the limit", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		overshot := 5
		for i := 0; i < maxQueuedTransactions+overshot; i++ {
			//nolint:gosec // this won't overflow
			_, err := insertConfirmedTransaction(m, uint64(i))
			require.NoError(t, err)
		}
		assert.Len(t, m.ConfirmedTransactions, maxQueuedTransactions+overshot)
		//nolint:gosec // this won't overflow
		_, _, err := m.MarkConfirmedAndReorgedTransactions(uint64(maxQueuedTransactions + overshot))
		require.NoError(t, err)
		assert.Len(t, m.ConfirmedTransactions, 170)
	})
}

func TestMarkUnconfirmedTransactionPurgeable(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)

	// fails if tx was not found
	err := m.MarkUnconfirmedTransactionPurgeable(0)
	require.Error(t, err)

	tx, err := insertUnconfirmedTransaction(m, 0)
	require.NoError(t, err)
	err = m.MarkUnconfirmedTransactionPurgeable(0)
	require.NoError(t, err)
	assert.True(t, tx.IsPurgeable)
}

func TestUpdateTransactionBroadcast(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	hash := testutils.NewHash()
	t.Run("fails if unconfirmed transaction was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		var nonce uint64
		require.Error(t, m.UpdateTransactionBroadcast(0, nonce, hash))
	})

	t.Run("fails if attempt was not found for a given transaction", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		var nonce uint64
		tx, err := insertUnconfirmedTransaction(m, nonce)
		require.NoError(t, err)
		require.Error(t, m.UpdateTransactionBroadcast(0, nonce, hash))

		// Attempt with different hash
		attempt := &types.Attempt{TxID: tx.ID, Hash: testutils.NewHash()}
		tx.Attempts = append(tx.Attempts, attempt)
		require.Error(t, m.UpdateTransactionBroadcast(0, nonce, hash))
	})

	t.Run("updates transaction's and attempt's broadcast times", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		var nonce uint64
		tx, err := insertUnconfirmedTransaction(m, nonce)
		require.NoError(t, err)
		attempt := &types.Attempt{TxID: tx.ID, Hash: hash}
		tx.Attempts = append(tx.Attempts, attempt)
		require.NoError(t, m.UpdateTransactionBroadcast(0, nonce, hash))
		assert.False(t, tx.LastBroadcastAt.IsZero())
		assert.False(t, attempt.BroadcastAt.IsZero())
		assert.False(t, tx.InitialBroadcastAt.IsZero())
	})
}

func TestUpdateUnstartedTransactionWithNonce(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	t.Run("returns nil if there are no unstarted transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		tx, err := m.UpdateUnstartedTransactionWithNonce(0)
		require.NoError(t, err)
		assert.Nil(t, tx)
	})

	t.Run("fails if there is already another unconfirmed transaction with the same nonce", func(t *testing.T) {
		var nonce uint64
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		insertUnstartedTransaction(m)
		_, err := insertUnconfirmedTransaction(m, nonce)
		require.NoError(t, err)

		_, err = m.UpdateUnstartedTransactionWithNonce(nonce)
		require.Error(t, err)
	})

	t.Run("updates unstarted transaction to unconfirmed and assigns a nonce", func(t *testing.T) {
		var nonce uint64
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		insertUnstartedTransaction(m)

		tx, err := m.UpdateUnstartedTransactionWithNonce(nonce)
		require.NoError(t, err)
		assert.Equal(t, nonce, *tx.Nonce)
		assert.Equal(t, txmgr.TxUnconfirmed, tx.State)
		assert.Empty(t, m.UnstartedTransactions)
	})
}

func TestDeleteAttemptForUnconfirmedTx(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	t.Run("fails if corresponding unconfirmed transaction for attempt was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		var nonce uint64
		tx := &types.Transaction{Nonce: &nonce}
		attempt := &types.Attempt{TxID: 0}
		err := m.DeleteAttemptForUnconfirmedTx(*tx.Nonce, attempt)
		require.Error(t, err)
	})

	t.Run("fails if corresponding unconfirmed attempt for txID was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		_, err := insertUnconfirmedTransaction(m, 0)
		require.NoError(t, err)

		attempt := &types.Attempt{TxID: 2, Hash: testutils.NewHash()}
		err = m.DeleteAttemptForUnconfirmedTx(0, attempt)

		require.Error(t, err)
	})

	t.Run("deletes attempt of unconfirmed transaction", func(t *testing.T) {
		hash := testutils.NewHash()
		var nonce uint64
		m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
		tx, err := insertUnconfirmedTransaction(m, nonce)
		require.NoError(t, err)

		attempt := &types.Attempt{TxID: 0, Hash: hash}
		tx.Attempts = append(tx.Attempts, attempt)
		err = m.DeleteAttemptForUnconfirmedTx(nonce, attempt)
		require.NoError(t, err)

		assert.Empty(t, tx.Attempts)
	})
}

func TestFindTxWithIdempotencyKey(t *testing.T) {
	t.Parallel()
	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
	tx, err := insertConfirmedTransaction(m, 0)
	require.NoError(t, err)

	ik := "IK"
	tx.IdempotencyKey = &ik
	itx := m.FindTxWithIdempotencyKey(ik)
	assert.Equal(t, ik, *itx.IdempotencyKey)

	uik := "Unknown"
	itx = m.FindTxWithIdempotencyKey(uik)
	assert.Nil(t, itx)
}

func TestPruneConfirmedTransactions(t *testing.T) {
	t.Parallel()
	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t), fromAddress, testutils.FixtureChainID)
	total := 5
	for i := 0; i < total; i++ {
		//nolint:gosec // this won't overflow
		_, err := insertConfirmedTransaction(m, uint64(i))
		require.NoError(t, err)
	}
	prunedTxIDs := m.pruneConfirmedTransactions()
	left := total - total/pruneSubset
	assert.Len(t, m.ConfirmedTransactions, left)
	assert.Len(t, prunedTxIDs, total/pruneSubset)
}

func insertUnstartedTransaction(m *InMemoryStore) *types.Transaction {
	m.Lock()
	defer m.Unlock()

	var nonce uint64
	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             &nonce,
		FromAddress:       m.address,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             txmgr.TxUnstarted,
	}

	m.UnstartedTransactions = append(m.UnstartedTransactions, tx)
	m.Transactions[tx.ID] = tx
	return tx
}

func insertUnconfirmedTransaction(m *InMemoryStore, nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             &nonce,
		FromAddress:       m.address,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             txmgr.TxUnconfirmed,
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	m.UnconfirmedTransactions[nonce] = tx
	m.Transactions[tx.ID] = tx
	return tx, nil
}

func insertConfirmedTransaction(m *InMemoryStore, nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             &nonce,
		FromAddress:       m.address,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             txmgr.TxConfirmed,
	}

	if _, exists := m.ConfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("a confirmed tx with the same nonce already exists: %v", m.ConfirmedTransactions[nonce])
	}

	m.ConfirmedTransactions[nonce] = tx
	m.Transactions[tx.ID] = tx
	return tx, nil
}

func insertFataTransaction(m *InMemoryStore) *types.Transaction {
	m.Lock()
	defer m.Unlock()

	var nonce uint64
	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             &nonce,
		FromAddress:       m.address,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             txmgr.TxFatalError,
	}

	m.FatalTransactions = append(m.FatalTransactions, tx)
	m.Transactions[tx.ID] = tx
	return tx
}
