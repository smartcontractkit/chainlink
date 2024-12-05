package txm

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/storage"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

func TestLifecycle(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	address1 := testutils.NewAddress()
	address2 := testutils.NewAddress()
	assert.NotEqual(t, address1, address2)
	addresses := []common.Address{address1, address2}
	keystore := mocks.NewKeystore(t)

	t.Run("retries if initial pending nonce call fails", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		config := Config{BlockTime: 1 * time.Minute}
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		keystore.On("EnabledAddressesForChain", mock.Anything, mock.Anything).Return([]common.Address{address1}, nil).Once()
		txm := NewTxm(lggr, testutils.FixtureChainID, client, nil, txStore, nil, config, keystore)
		client.On("PendingNonceAt", mock.Anything, address1).Return(uint64(0), errors.New("error")).Once()
		client.On("PendingNonceAt", mock.Anything, address1).Return(uint64(0), nil).Once()
		require.NoError(t, txm.Start(tests.Context(t)))
		tests.AssertLogEventually(t, observedLogs, "Error when fetching initial pending nonce")
	})

	t.Run("tests lifecycle successfully without any transactions", func(t *testing.T) {
		config := Config{BlockTime: 200 * time.Millisecond}
		keystore.On("EnabledAddressesForChain", mock.Anything, mock.Anything).Return(addresses, nil).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		require.NoError(t, txStore.Add(addresses...))
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		var nonce uint64
		// Start
		client.On("PendingNonceAt", mock.Anything, address1).Return(nonce, nil).Once()
		client.On("PendingNonceAt", mock.Anything, address2).Return(nonce, nil).Once()
		// backfill loop (may or may not be executed multiple times)
		client.On("NonceAt", mock.Anything, address1, mock.Anything).Return(nonce, nil).Maybe()
		client.On("NonceAt", mock.Anything, address2, mock.Anything).Return(nonce, nil).Maybe()

		servicetest.Run(t, txm)
		tests.AssertLogEventually(t, observedLogs, "Backfill time elapsed")
	})
}

func TestTrigger(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)
	t.Run("Trigger fails if Txm is unstarted", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		txm := NewTxm(lggr, nil, nil, nil, nil, nil, Config{}, keystore)
		txm.Trigger(address)
		tests.AssertLogEventually(t, observedLogs, "Txm unstarted")
	})

	t.Run("executes Trigger", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		require.NoError(t, txStore.Add(address))
		client := mocks.NewClient(t)
		ab := mocks.NewAttemptBuilder(t)
		config := Config{BlockTime: 1 * time.Minute, RetryBlockThreshold: 10}
		keystore.On("EnabledAddressesForChain", mock.Anything, mock.Anything).Return([]common.Address{address}, nil)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		var nonce uint64
		// Start
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil).Once()
		servicetest.Run(t, txm)
		txm.Trigger(address)
	})
}

func TestBroadcastTransaction(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	config := Config{}
	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)

	t.Run("fails if FetchUnconfirmedTransactionAtNonceWithCount for unconfirmed transactions fails", func(t *testing.T) {
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("call failed")).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		bo, err := txm.broadcastTransaction(ctx, address)
		require.Error(t, err)
		assert.False(t, bo)
		require.ErrorContains(t, err, "call failed")
	})

	t.Run("throws a warning and returns if unconfirmed transactions exceed maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, maxInFlightTransactions+1, nil).Once()
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.True(t, bo)
		require.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, "Reached transaction limit")
	})

	t.Run("checks pending nonce if unconfirmed transactions are more than 1/3 of maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewTxStore(t)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		txm.setNonce(address, 1)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, maxInFlightTransactions/3, nil).Twice()

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(0), nil).Once() // LocalNonce: 1, PendingNonce: 0
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.True(t, bo)
		require.NoError(t, err)

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(1), nil).Once() // LocalNonce: 1, PendingNonce: 1
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		bo, err = txm.broadcastTransaction(ctx, address)
		assert.False(t, bo)
		require.NoError(t, err)
		tests.AssertLogCountEventually(t, observedLogs, "Reached transaction limit.", 1)
	})

	t.Run("fails if UpdateUnstartedTransactionWithNonce fails", func(t *testing.T) {
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, nil).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("call failed")).Once()
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.False(t, bo)
		require.Error(t, err)
		require.ErrorContains(t, err, "call failed")
	})

	t.Run("returns if there are no unstarted transactions", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		require.NoError(t, txStore.Add(address))
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		bo, err := txm.broadcastTransaction(ctx, address)
		require.NoError(t, err)
		assert.False(t, bo)
		assert.Equal(t, uint64(0), txm.getNonce(address))
	})

	t.Run("picks a new tx and creates a new attempt then sends it and updates the broadcast time", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		require.NoError(t, txStore.Add(address))
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		txm.setNonce(address, 8)
		IDK := "IDK"
		txRequest := &types.TxRequest{
			Data:              []byte{100},
			IdempotencyKey:    &IDK,
			ChainID:           testutils.FixtureChainID,
			FromAddress:       address,
			ToAddress:         testutils.NewAddress(),
			SpecifiedGasLimit: 22000,
		}
		tx, err := txm.CreateTransaction(tests.Context(t), txRequest)
		require.NoError(t, err)
		attempt := &types.Attempt{
			TxID:     tx.ID,
			Fee:      gas.EvmFee{GasPrice: assets.NewWeiI(1)},
			GasLimit: 22000,
		}
		ab.On("NewAttempt", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(attempt, nil).Once()
		client.On("SendTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		bo, err := txm.broadcastTransaction(ctx, address)
		require.NoError(t, err)
		assert.False(t, bo)
		assert.Equal(t, uint64(9), txm.getNonce(address))
		tx, err = txStore.FindTxWithIdempotencyKey(tests.Context(t), &IDK)
		require.NoError(t, err)
		assert.Len(t, tx.Attempts, 1)
		var zeroTime time.Time
		assert.Greater(t, tx.LastBroadcastAt, zeroTime)
		assert.Greater(t, tx.Attempts[0].BroadcastAt, zeroTime)
		assert.Greater(t, tx.InitialBroadcastAt, zeroTime)
	})
}

func TestBackfillTransactions(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	storage := mocks.NewTxStore(t)
	config := Config{}
	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)

	t.Run("fails if latest nonce fetching fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, nil, config, keystore)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), errors.New("latest nonce fail")).Once()
		bo, err := txm.backfillTransactions(ctx, address)
		require.Error(t, err)
		assert.False(t, bo)
		require.ErrorContains(t, err, "latest nonce fail")
	})

	t.Run("fails if MarkConfirmedAndReorgedTransactions fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, nil, config, keystore)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), nil).Once()
		storage.On("MarkConfirmedAndReorgedTransactions", mock.Anything, mock.Anything, address).
			Return([]*types.Transaction{}, []uint64{}, errors.New("marking transactions confirmed failed")).Once()
		bo, err := txm.backfillTransactions(ctx, address)
		require.Error(t, err)
		assert.False(t, bo)
		require.ErrorContains(t, err, "marking transactions confirmed failed")
	})
}
