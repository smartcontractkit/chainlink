package txm

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/storage"
)

func TestLifecycle(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	config := Config{BlockTime: 10 * time.Millisecond}
	address := testutils.NewAddress()

	t.Run("fails to start if initial pending nonce call fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, nil, config, address)
		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(0), errors.New("error")).Once()
		assert.Error(t, txm.Start(tests.Context(t)))
	})

	t.Run("tests lifecycle successfully without any transactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		txStore := storage.NewInMemoryStore(lggr)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, config, address)
		var nonce uint64 = 0
		// Start
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil).Once()
		// backfill loop (may or may not be executed multiple times)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(nonce, nil)

		servicetest.Run(t, txm)
		tests.AssertLogEventually(t, observedLogs, "Backfill time elapsed")
	})

}

func TestTrigger(t *testing.T) {
	t.Parallel()

	t.Run("Trigger fails if Txm is unstarted", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), nil, nil, nil, nil, Config{}, common.Address{})
		txm.Trigger()
		assert.Error(t, txm.Trigger(), "Txm unstarted")
	})

	t.Run("executes Trigger", func(t *testing.T) {
		lggr := logger.Test(t)
		address := testutils.NewAddress()
		txStore := storage.NewInMemoryStore(lggr)
		client := mocks.NewClient(t)
		ab := mocks.NewAttemptBuilder(t)
		config := Config{BlockTime: 10 * time.Second}
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, config, address)
		var nonce uint64 = 0
		// Start
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil).Once()
		servicetest.Run(t, txm)
		assert.NoError(t, txm.Trigger())
	})
}

func TestBroadcastTransaction(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	config := Config{}
	address := testutils.NewAddress()

	t.Run("fails if FetchUnconfirmedTransactionAtNonceWithCount for unconfirmed transactions fails", func(t *testing.T) {
		mTxStore := mocks.NewStorage(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("call failed")).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, config, address)
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "call failed")
	})

	t.Run("throws a warning and returns if unconfirmed transactions exceed maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewStorage(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, int(maxInFlightTransactions+1), nil).Once()
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, config, address)
		txm.broadcastTransaction()
		tests.AssertLogEventually(t, observedLogs, "Reached transaction limit")
	})

	t.Run("checks pending nonce if unconfirmed transactions are more than 1/3 of maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewStorage(t)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, config, address)
		txm.nonce.Store(1)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, int(maxInFlightTransactions/3), nil).Twice()

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(0), nil).Once() // LocalNonce: 1, PendingNonce: 0
		txm.broadcastTransaction()

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(1), nil).Once() // LocalNonce: 1, PendingNonce: 1
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		txm.broadcastTransaction()
		tests.AssertLogCountEventually(t, observedLogs, "Reached transaction limit.", 1)

	})

	t.Run("fails if UpdateUnstartedTransactionWithNonce fails", func(t *testing.T) {
		mTxStore := mocks.NewStorage(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, nil).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, config, address)
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("call failed")).Once()
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "call failed")
	})

	t.Run("returns if there are no unstarted transactions", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStore(lggr)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, config, address)
		err := txm.broadcastTransaction()
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), txm.nonce.Load())
	})
}

func TestBackfillTransactions(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	storage := mocks.NewStorage(t)
	config := Config{}
	address := testutils.NewAddress()

	t.Run("fails if latest nonce fetching fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), errors.New("latest nonce fail")).Once()
		err := txm.backfillTransactions()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "latest nonce fail")
	})

	t.Run("fails if MarkTransactionsConfirmed fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), nil)
		storage.On("MarkTransactionsConfirmed", mock.Anything, mock.Anything, address).Return([]uint64{}, []uint64{}, errors.New("marking transactions confirmed failed"))
		err := txm.backfillTransactions()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "marking transactions confirmed failed")
	})
}
