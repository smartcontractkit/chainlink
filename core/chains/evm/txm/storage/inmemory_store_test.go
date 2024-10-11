package storage

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

func TestAbandonPendingTransactions(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))
	t.Run("abandons unstarted and unconfirmed transactions", func(t *testing.T) {
		// Unstarted
		tx1 := insertUnstartedTransaction(m, fromAddress)
		tx2 := insertUnstartedTransaction(m, fromAddress)

		// Unconfirmed
		tx3, err := insertUnconfirmedTransaction(m, fromAddress, 3)
		assert.NoError(t, err)
		tx4, err := insertUnconfirmedTransaction(m, fromAddress, 4)
		assert.NoError(t, err)

		assert.NoError(t, m.AbandonPendingTransactions(tests.Context(t), fromAddress))

		assert.Equal(t, types.TxFatalError, tx1.State)
		assert.Equal(t, types.TxFatalError, tx2.State)
		assert.Equal(t, types.TxFatalError, tx3.State)
		assert.Equal(t, types.TxFatalError, tx4.State)
	})

	t.Run("skips all types apart from unstarted and unconfirmed transactions", func(t *testing.T) {
		// Fatal
		tx1 := insertFataTransaction(m, fromAddress)
		tx2 := insertFataTransaction(m, fromAddress)

		// Confirmed
		tx3, err := insertConfirmedTransaction(m, fromAddress, 3)
		assert.NoError(t, err)
		tx4, err := insertConfirmedTransaction(m, fromAddress, 4)
		assert.NoError(t, err)

		assert.NoError(t, m.AbandonPendingTransactions(tests.Context(t), fromAddress))

		assert.Equal(t, types.TxFatalError, tx1.State)
		assert.Equal(t, types.TxFatalError, tx2.State)
		assert.Equal(t, types.TxConfirmed, tx3.State)
		assert.Equal(t, types.TxConfirmed, tx4.State)

	})

}

func TestAppendAttemptToTransaction(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))

	_, err := insertUnconfirmedTransaction(m, fromAddress, 0) // txID = 1
	assert.NoError(t, err)
	_, err = insertConfirmedTransaction(m, fromAddress, 2) // txID = 1
	assert.NoError(t, err)

	t.Run("fails if corresponding unconfirmed transaction for attempt was not found", func(t *testing.T) {
		var nonce uint64 = 1
		newAttempt := &types.Attempt{
			TxID: 1,
		}
		assert.Error(t, m.AppendAttemptToTransaction(tests.Context(t), nonce, newAttempt))
	})

	t.Run("fails if unconfirmed transaction was found but has doesn't match the txID", func(t *testing.T) {
		var nonce uint64 = 0
		newAttempt := &types.Attempt{
			TxID: 2,
		}
		assert.Error(t, m.AppendAttemptToTransaction(tests.Context(t), nonce, newAttempt))
	})

	t.Run("appends attempt to transaction", func(t *testing.T) {
		var nonce uint64 = 0
		newAttempt := &types.Attempt{
			TxID: 1,
		}
		assert.NoError(t, m.AppendAttemptToTransaction(tests.Context(t), nonce, newAttempt))

	})
}

func TestCountUnstartedTransactions(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))
	n, _ := m.CountUnstartedTransactions(tests.Context(t), fromAddress)
	assert.Equal(t, 0, n)

	insertUnstartedTransaction(m, fromAddress)
	n, _ = m.CountUnstartedTransactions(tests.Context(t), fromAddress)
	assert.Equal(t, 1, n)

}

func TestCreateEmptyUnconfirmedTransaction(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))
	insertUnconfirmedTransaction(m, fromAddress, 0)

	t.Run("fails if unconfirmed transaction with the same nonce exists", func(t *testing.T) {
		_, err := m.CreateEmptyUnconfirmedTransaction(tests.Context(t), fromAddress, testutils.FixtureChainID, 0, 0)
		assert.Error(t, err)
	})

	t.Run("creates a new empty unconfirmed transaction", func(t *testing.T) {
		tx, err := m.CreateEmptyUnconfirmedTransaction(tests.Context(t), fromAddress, testutils.FixtureChainID, 1, 0)
		assert.NoError(t, err)
		assert.Equal(t, types.TxUnconfirmed, tx.State)
	})

}

func TestCreateTransaction(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))

	tx1 := &types.Transaction{}
	tx2 := &types.Transaction{}
	id1, err := m.CreateTransaction(tests.Context(t), tx1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), id1)

	id2, err := m.CreateTransaction(tests.Context(t), tx2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), id2)

	count, _ := m.CountUnstartedTransactions(tests.Context(t), fromAddress)
	assert.Equal(t, count, 2)

}

func TestFetchUnconfirmedTransactionAtNonceWithCount(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))

	tx, count, _ := m.FetchUnconfirmedTransactionAtNonceWithCount(tests.Context(t), 0, fromAddress)
	assert.Nil(t, tx)
	assert.Equal(t, 0, count)

	var nonce uint64 = 0
	insertUnconfirmedTransaction(m, fromAddress, nonce)
	tx, count, _ = m.FetchUnconfirmedTransactionAtNonceWithCount(tests.Context(t), nonce, fromAddress)
	assert.Equal(t, tx.Nonce, nonce)
	assert.Equal(t, 1, count)

}

func TestMarkTransactionsConfirmed(t *testing.T) {

	fromAddress := testutils.NewAddress()

	t.Run("returns 0 if there are no transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		un, cn, err := m.MarkTransactionsConfirmed(tests.Context(t), 100, fromAddress)
		assert.NoError(t, err)
		assert.Equal(t, len(un), 0)
		assert.Equal(t, len(cn), 0)
	})

	t.Run("confirms transaction with nonce lower than the latest", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		ctx1, err := insertUnconfirmedTransaction(m, fromAddress, 0)
		assert.NoError(t, err)

		ctx2, err := insertUnconfirmedTransaction(m, fromAddress, 1)
		assert.NoError(t, err)

		ctxs, utxs, err := m.MarkTransactionsConfirmed(tests.Context(t), 1, fromAddress)
		assert.NoError(t, err)
		assert.Equal(t, types.TxConfirmed, ctx1.State)
		assert.Equal(t, types.TxUnconfirmed, ctx2.State)
		assert.Equal(t, ctxs[0], ctx1.ID)
		assert.Equal(t, 0, len(utxs))
	})

	t.Run("unconfirms transaction with nonce equal to or higher than the latest", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		ctx1, err := insertConfirmedTransaction(m, fromAddress, 0)
		assert.NoError(t, err)

		ctx2, err := insertConfirmedTransaction(m, fromAddress, 1)
		assert.NoError(t, err)

		ctxs, utxs, err := m.MarkTransactionsConfirmed(tests.Context(t), 1, fromAddress)
		assert.NoError(t, err)
		assert.Equal(t, types.TxConfirmed, ctx1.State)
		assert.Equal(t, types.TxUnconfirmed, ctx2.State)
		assert.Equal(t, utxs[0], ctx2.ID)
		assert.Equal(t, 0, len(ctxs))
	})
}

func TestMarkUnconfirmedTransactionPurgeable(t *testing.T) {

	fromAddress := testutils.NewAddress()
	m := NewInMemoryStore(logger.Test(t))

	// fails if tx was not found
	err := m.MarkUnconfirmedTransactionPurgeable(tests.Context(t), 0)
	assert.Error(t, err)

	tx, err := insertUnconfirmedTransaction(m, fromAddress, 0)
	assert.NoError(t, err)
	err = m.MarkUnconfirmedTransactionPurgeable(tests.Context(t), 0)
	assert.NoError(t, err)
	assert.Equal(t, true, tx.IsPurgeable)
}

func TestUpdateTransactionBroadcast(t *testing.T) {

	fromAddress := testutils.NewAddress()
	hash := testutils.NewHash()
	t.Run("fails if unconfirmed transaction was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		var nonce uint64 = 0
		assert.Error(t, m.UpdateTransactionBroadcast(tests.Context(t), 0, nonce, hash))
	})

	t.Run("fails if attempt was not found for a given transaction", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		var nonce uint64 = 0
		tx, err := insertUnconfirmedTransaction(m, fromAddress, nonce)
		assert.NoError(t, err)
		assert.Error(t, m.UpdateTransactionBroadcast(tests.Context(t), 0, nonce, hash))

		// Attempt with different hash
		attempt := &types.Attempt{TxID: tx.ID, Hash: testutils.NewHash()}
		tx.Attempts = append(tx.Attempts, attempt)
		assert.Error(t, m.UpdateTransactionBroadcast(tests.Context(t), 0, nonce, hash))
	})

	t.Run("updates transaction's and attempt's broadcast times", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		var nonce uint64 = 0
		tx, err := insertUnconfirmedTransaction(m, fromAddress, nonce)
		assert.NoError(t, err)
		attempt := &types.Attempt{TxID: tx.ID, Hash: hash}
		tx.Attempts = append(tx.Attempts, attempt)
		assert.NoError(t, m.UpdateTransactionBroadcast(tests.Context(t), 0, nonce, hash))
		assert.False(t, tx.LastBroadcastAt.IsZero())
		assert.False(t, attempt.BroadcastAt.IsZero())
	})
}

func TestUpdateUnstartedTransactionWithNonce(t *testing.T) {

	fromAddress := testutils.NewAddress()
	t.Run("returns nil if there are no unstarted transactions", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		tx, err := m.UpdateUnstartedTransactionWithNonce(tests.Context(t), fromAddress, 0)
		assert.NoError(t, err)
		assert.Nil(t, tx)
	})

	t.Run("fails if there is already another unstarted transaction with the same nonce", func(t *testing.T) {
		var nonce uint64 = 0
		m := NewInMemoryStore(logger.Test(t))
		insertUnstartedTransaction(m, fromAddress)
		_, err := insertUnconfirmedTransaction(m, fromAddress, nonce)
		assert.NoError(t, err)

		_, err = m.UpdateUnstartedTransactionWithNonce(tests.Context(t), fromAddress, nonce)
		assert.Error(t, err)
	})

	t.Run("updates unstarted transaction to unconfirmed and assigns a nonce", func(t *testing.T) {
		var nonce uint64 = 0
		m := NewInMemoryStore(logger.Test(t))
		insertUnstartedTransaction(m, fromAddress)

		tx, err := m.UpdateUnstartedTransactionWithNonce(tests.Context(t), fromAddress, nonce)
		assert.NoError(t, err)
		assert.Equal(t, nonce, tx.Nonce)
		assert.Equal(t, types.TxUnconfirmed, tx.State)
	})
}

func TestDeleteAttemptForUnconfirmedTx(t *testing.T) {

	fromAddress := testutils.NewAddress()
	t.Run("fails if corresponding unconfirmed transaction for attempt was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		tx := &types.Transaction{Nonce: 0}
		attempt := &types.Attempt{TxID: 0}
		err := m.DeleteAttemptForUnconfirmedTx(tests.Context(t), tx.Nonce, attempt)
		assert.Error(t, err)
	})

	t.Run("fails if corresponding unconfirmed attempt for txID was not found", func(t *testing.T) {
		m := NewInMemoryStore(logger.Test(t))
		_, err := insertUnconfirmedTransaction(m, fromAddress, 0)
		assert.NoError(t, err)

		attempt := &types.Attempt{TxID: 2, Hash: testutils.NewHash()}
		err = m.DeleteAttemptForUnconfirmedTx(tests.Context(t), 0, attempt)

		assert.Error(t, err)
	})

	t.Run("deletes attempt of unconfirmed transaction", func(t *testing.T) {
		hash := testutils.NewHash()
		var nonce uint64 = 0
		m := NewInMemoryStore(logger.Test(t))
		tx, err := insertUnconfirmedTransaction(m, fromAddress, nonce)
		assert.NoError(t, err)

		attempt := &types.Attempt{TxID: 0, Hash: hash}
		tx.Attempts = append(tx.Attempts, attempt)
		err = m.DeleteAttemptForUnconfirmedTx(tests.Context(t), nonce, attempt)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(tx.Attempts))
	})
}

func insertUnstartedTransaction(m *InMemoryStore, fromAddress common.Address) *types.Transaction {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             0,
		FromAddress:       fromAddress,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             types.TxUnstarted,
	}

	m.UnstartedTransactions = append(m.UnstartedTransactions, tx)
	return tx
}

func insertUnconfirmedTransaction(m *InMemoryStore, fromAddress common.Address, nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             nonce,
		FromAddress:       fromAddress,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             types.TxUnconfirmed,
	}

	if _, exists := m.UnconfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("an unconfirmed tx with the same nonce already exists: %v", m.UnconfirmedTransactions[nonce])
	}

	m.UnconfirmedTransactions[nonce] = tx
	return tx, nil
}

func insertConfirmedTransaction(m *InMemoryStore, fromAddress common.Address, nonce uint64) (*types.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             nonce,
		FromAddress:       fromAddress,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             types.TxConfirmed,
	}

	if _, exists := m.ConfirmedTransactions[nonce]; exists {
		return nil, fmt.Errorf("a confirmed tx with the same nonce already exists: %v", m.ConfirmedTransactions[nonce])
	}

	m.ConfirmedTransactions[nonce] = tx
	return tx, nil
}

func insertFataTransaction(m *InMemoryStore, fromAddress common.Address) *types.Transaction {
	m.Lock()
	defer m.Unlock()

	m.txIDCount++
	tx := &types.Transaction{
		ID:                m.txIDCount,
		ChainID:           testutils.FixtureChainID,
		Nonce:             0,
		FromAddress:       fromAddress,
		ToAddress:         testutils.NewAddress(),
		Value:             big.NewInt(0),
		SpecifiedGasLimit: 0,
		CreatedAt:         time.Now(),
		State:             types.TxFatalError,
	}

	m.FatalTransactions = append(m.FatalTransactions, tx)
	return tx
}
