package txm

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/mocks"
)

func TestLifecycle(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	storage := mocks.NewStorage(t)
	config := Config{}
	address := testutils.NewAddress()

	t.Run("fails to start if pending nonce call fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(0), errors.New("error")).Once()
		assert.Error(t, txm.Start(tests.Context(t)))
	})

	t.Run("tests lifecycle successfully without any transactions", func(t *testing.T) {
		lggr, _ := logger.TestObserved(t, zap.DebugLevel)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, storage, config, address)
		var nonce uint64 = 0
		// Start
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil).Once()
		// broadcast loop (may or may not be executed multiple times)
		client.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil)
		storage.On("UpdateUnstartedTransactionWithNonce", mock.Anything, address, mock.Anything).Return(nil, nil)
		// backfill loop (may or may not be executed multiple times)
		client.On("NonceAt", mock.Anything, address, nil).Return(nonce, nil)
		storage.On("MarkTransactionsConfirmed", mock.Anything, nonce, address).Return([]uint64{}, []uint64{}, nil)
		storage.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, nonce, address).Return(nil, 0)
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil)
		storage.On("CountUnstartedTransactions", mock.Anything, address).Return(0)

		assert.NoError(t, txm.Start(tests.Context(t)))
		assert.NoError(t, txm.Close())
	})

}

func TestTrigger(t *testing.T) {
	t.Parallel()

	t.Run("Trigger fails if Txm is unstarted", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), nil, nil, nil, nil, Config{}, common.Address{})
		txm.Trigger()
		assert.Error(t, txm.Trigger(), "Txm unstarted")
	})
}


func TestBroadcastTransaction(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	storage := mocks.NewStorage(t)
	config := Config{}
	address := testutils.NewAddress()

	t.Run("fails if batch call for pending and latest nonce fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		client.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("batch call error")).Once()
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "batch call error")
	})

	t.Run("fails if batch call for pending and latest nonce fails for one of them", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		//pending nonce
		client.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Error = errors.New("pending nonce failed")
		}).Return(nil).Once()
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pending nonce failed")

		// latest nonce
		client.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[1].Error = errors.New("latest nonce failed")
		}).Return(nil).Once()
		err = txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "latest nonce failed")
	})

	t.Run("throws a warning if maxInFlightTransactions are reached", func(t *testing.T) {
		pending := "0x100"
		latest := "0x0"
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		client.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &pending  // pending
			elems[1].Result =  &latest  // latest
		}).Return(nil).Once()
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Reached transaction limit")

	})
	t.Run("fails if UpdateUnstartedTransactionWithNonce fails", func(t *testing.T) {
		pending := "0x8"
		latest := "0x0"
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, config, address)
		txm.nonce.Store(0)
		client.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &pending  // pending
			elems[1].Result =  &latest  // latest
		}).Return(nil).Once()
		storage.On("UpdateUnstartedTransactionWithNonce", mock.Anything, address, mock.Anything).Return(nil, errors.New("update failed"))
		err := txm.broadcastTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
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
