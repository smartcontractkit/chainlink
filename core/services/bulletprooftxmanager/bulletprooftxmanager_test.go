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
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	pgmocks "github.com/smartcontractkit/chainlink/core/services/postgres/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBulletproofTxManager_SendEther_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1)

	_, err := bulletprooftxmanager.SendEther(store.DB, from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestBulletproofTxManager_CheckEthTxQueueCapacity(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	db := store.DB
	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnstartedEthTx(t, store, otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertFatalErrorEthTx(t, store, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	var n int64 = 0
	cltest.MustInsertInProgressEthTxWithAttempt(t, store, n, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, 1)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithAttempt(t, store, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnstartedEthTx(t, store, fromAddress)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	cltest.MustInsertUnstartedEthTx(t, store, fromAddress)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse", maxUnconfirmedTransactions))

		cltest.MustInsertUnstartedEthTx(t, store, fromAddress)
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

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2, fromAddress)

	count, err := bulletprooftxmanager.CountUnconfirmedTransactions(store.DB, fromAddress)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestBulletproofTxManager_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	key := cltest.MustInsertRandomKey(t, store.DB, 0)
	fromAddress := key.Address.Address()
	toAddress := cltest.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")

	bptxm := bulletprooftxmanager.NewBulletproofTxManager(store.DB, nil, config, nil, nil, nil)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		etx, err := bptxm.CreateEthTransaction(store.DB, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, models.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, store, models.EthTx{}, 1)

		require.NoError(t, store.ORM.DB.First(&etx).Error)

		assert.Equal(t, etx.State, models.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, store, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		_, err := bptxm.CreateEthTransaction(store.DB, fromAddress, cltest.NewAddress(), []byte{1, 2, 3}, 21000, nil, bulletprooftxmanager.SendEveryStrategy{})
		assert.EqualError(t, err, "BulletproofTxManager#CreateEthTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse")
	})
}

func TestBulletproofTxManager_CreateEthTransaction_OutOfEth(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	thisKey := cltest.MustInsertRandomKey(t, store.DB, 1)
	otherKey := cltest.MustInsertRandomKey(t, store.DB, 1)

	fromAddress := thisKey.Address.Address()
	gasLimit := uint64(1000)
	toAddress := cltest.NewAddress()

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	bptxm := bulletprooftxmanager.NewBulletproofTxManager(store.DB, nil, config, nil, nil, nil)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, store, 0, otherKey.Address.Address())
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)

		etx, err := bptxm.CreateEthTransaction(store.DB, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})

	require.NoError(t, store.DB.Exec(`DELETE FROM eth_txes WHERE from_address = ?`, thisKey.Address.Address()).Error)

	t.Run("if this key has any transactions with insufficient eth errors, skips transmission entirely", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, store, 0, thisKey.Address.Address())
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})

		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		_, err := bptxm.CreateEthTransaction(store.DB, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		require.EqualError(t, err, fmt.Sprintf("wallet is out of eth: %s", thisKey.Address.Hex()))
		strategy.AssertExpectations(t)
	})

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		require.NoError(t, store.DB.Exec(`UPDATE eth_tx_attempts SET state = 'broadcast'`).Error)
		require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET nonce = 0, state = 'confirmed', broadcast_at = NOW()`).Error)
		strategy := new(bptxmmocks.TxStrategy)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*gorm.DB")).Return(int64(0), nil)

		config.On("EthMaxQueuedTransactions").Return(uint64(1))
		etx, err := bptxm.CreateEthTransaction(store.DB, fromAddress, toAddress, payload, gasLimit, nil, strategy)
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})
}

func TestBulletproofTxManager_Lifecycle(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db := store.DB
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
