package bulletprooftxmanager_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	pgmocks "github.com/smartcontractkit/chainlink/core/services/postgres/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBulletproofTxManager_SendEther_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := pgtest.NewGormDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1)

	_, err := bulletprooftxmanager.SendEther(db, from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestBulletproofTxManager_CheckEthTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnstartedEthTx(t, db, otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertFatalErrorEthTx(t, db, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	var n int64 = 0
	cltest.MustInsertInProgressEthTxWithAttempt(t, db, n, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, 1)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithAttempt(t, db, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnstartedEthTx(t, db, fromAddress)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	cltest.MustInsertUnstartedEthTx(t, db, fromAddress)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse", maxUnconfirmedTransactions))

		cltest.MustInsertUnstartedEthTx(t, db, fromAddress)
		err = bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.Error(t, err)

		require.EqualError(t, err, fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse", maxUnconfirmedTransactions))
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, 0)
		require.NoError(t, err)
	})
}

func TestBulletproofTxManager_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress)

	count, err := bulletprooftxmanager.CountUnconfirmedTransactions(db, fromAddress)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestBulletproofTxManager_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	cltest.MustInsertUnstartedEthTx(t, db, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, db, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, db, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress)

	count, err := bulletprooftxmanager.CountUnstartedTransactions(db, fromAddress)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}
func TestBulletproofTxManager_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)

	key := cltest.MustInsertRandomKey(t, db, 0)
	fromAddress := key.Address.Address()
	toAddress := cltest.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")

	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, nil, config, nil, nil, nil)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		etx, err := bptxm.CreateEthTransaction(db, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, bulletprooftxmanager.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, bulletprooftxmanager.EthTx{}, 1)

		require.NoError(t, db.First(&etx).Error)

		assert.Equal(t, etx.State, bulletprooftxmanager.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		_, err := bptxm.CreateEthTransaction(db, fromAddress, cltest.NewAddress(), []byte{1, 2, 3}, 21000, nil, bulletprooftxmanager.SendEveryStrategy{})
		assert.EqualError(t, err, "BulletproofTxManager#CreateEthTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse")
	})
}

func TestBulletproofTxManager_CreateEthTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewGormDB(t)

	thisKey := cltest.MustInsertRandomKey(t, db, 1)
	otherKey := cltest.MustInsertRandomKey(t, db, 1)

	fromAddress := thisKey.Address.Address()
	gasLimit := uint64(1000)
	toAddress := cltest.NewAddress()

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, nil, config, nil, nil, nil)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 0, otherKey.Address.Address())
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)

		etx, err := bptxm.CreateEthTransaction(db, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})

	require.NoError(t, db.Exec(`DELETE FROM eth_txes WHERE from_address = ?`, thisKey.Address.Address()).Error)

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 0, thisKey.Address.Address())
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)

		etx, err := bptxm.CreateEthTransaction(db, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})

	require.NoError(t, db.Exec(`DELETE FROM eth_txes WHERE from_address = ?`, thisKey.Address.Address()).Error)

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 0, 42, thisKey.Address.Address())
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)

		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		etx, err := bptxm.CreateEthTransaction(db, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})
}

func TestBulletproofTxManager_Lifecycle(t *testing.T) {
	db := pgtest.NewGormDB(t)

	ethClient := new(mocks.Client)
	config := new(bptxmmocks.Config)
	kst := new(ksmocks.EthKeyStoreInterface)
	advisoryLocker := &postgres.NullAdvisoryLocker{}
	eventBroadcaster := new(pgmocks.EventBroadcaster)

	config.On("EthTxResendAfterThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperInterval").Return(1 * time.Hour)
	config.On("EthMaxInFlightTransactions").Return(uint32(42))
	config.On("EthFinalityDepth").Maybe().Return(uint(42))
	config.On("GasEstimatorMode").Return("FixedPrice")
	kst.On("AllKeys").Return([]ethkey.Key{}, nil).Once()

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened)

	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, ethClient, config, kst, advisoryLocker, eventBroadcaster)

	head := cltest.Head(42)
	// It should not hang or panic
	bptxm.OnNewLongestChain(context.Background(), *head)

	sub := new(pgmocks.Subscription)
	sub.On("Events").Return(make(<-chan postgres.Event))
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)
	config.On("EthNonceAutoSync").Return(true)
	config.On("EthGasBumpThreshold").Return(uint64(1))

	require.NoError(t, bptxm.Start())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	bptxm.OnNewLongestChain(ctx, *head)
	require.NoError(t, ctx.Err())

	key := cltest.MustGenerateRandomKey(t)

	kst.On("AllKeys").Return([]ethkey.Key{key}, nil).Once()
	sub.On("Close").Return()
	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.timerCtx"), key.Address.Address()).Return(uint64(0), nil)
	config.On("TriggerFallbackDBPollInterval").Return(1 * time.Hour)
	keyChangeCh <- struct{}{}

	require.NoError(t, bptxm.Close())

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
	kst.AssertExpectations(t)
	eventBroadcaster.AssertExpectations(t)
	unsub.AwaitOrFail(t, 1*time.Second)
}
