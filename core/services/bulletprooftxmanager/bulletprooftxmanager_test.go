package bulletprooftxmanager_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	pgmocks "github.com/smartcontractkit/chainlink/core/services/pg/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestBulletproofTxManager_SendEther_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1)

	q := pg.NewQ(db, logger.TestLogger(t), cltest.NewTestGeneralConfig(t))
	_, err := bulletprooftxmanager.SendEther(q, big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestBulletproofTxManager_CheckEthTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewBulletproofTxManagerORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertFatalErrorEthTx(t, borm, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	var n int64 = 0
	cltest.MustInsertInProgressEthTxWithAttempt(t, borm, n, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, 1, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse", maxUnconfirmedTransactions))

		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
		err = bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)

		require.EqualError(t, err, fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse", maxUnconfirmedTransactions))
	})

	t.Run("with different chain ID ignores txes", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, *big.NewInt(42))
		require.NoError(t, err)
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, 0, cltest.FixtureChainID)
		require.NoError(t, err)
	})
}

func TestBulletproofTxManager_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewBulletproofTxManagerORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	q := pg.NewQ(db, logger.TestLogger(t), cfg)
	count, err := bulletprooftxmanager.CountUnconfirmedTransactions(q, fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestBulletproofTxManager_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewBulletproofTxManagerORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	q := pg.NewQ(db, logger.TestLogger(t), cfg)
	count, err := bulletprooftxmanager.CountUnstartedTransactions(q, fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}
func TestBulletproofTxManager_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewBulletproofTxManagerORM(t, db, cfg)

	keyStore := cltest.NewKeyStore(t, db, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)
	toAddress := cltest.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	lggr := logger.TestLogger(t)
	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, ethClient, config, nil, nil, lggr)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		etx, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, bulletprooftxmanager.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, db.Get(&etx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, bulletprooftxmanager.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		_, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      cltest.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Meta:           nil,
			Strategy:       bulletprooftxmanager.SendEveryStrategy{},
		})
		assert.EqualError(t, err, "BulletproofTxManager#CreateEthTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS is a sanity limit and should never happen under normal operation. This error is very unlikely to be a problem with Chainlink, and instead more likely to be caused by a problem with your eth node's connectivity. Check your eth node: it may not be broadcasting transactions to the network, or it might be overloaded and evicting Chainlink's transactions from its mempool. Increasing ETH_MAX_QUEUED_TRANSACTIONS is almost certainly not the correct action to take here unless you ABSOLUTELY know what you are doing, and will probably make things worse")
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		id := uuid.NewV4()
		tx1, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:       fromAddress,
			ToAddress:         cltest.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			GasLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          bulletprooftxmanager.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		tx2, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:       fromAddress,
			ToAddress:         cltest.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			GasLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          bulletprooftxmanager.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		assert.Equal(t, tx1.ID, tx2.ID)
	})

	t.Run("returns error if eth key state is missing or doesn't match chain ID", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Twice()
		rndAddr := cltest.NewAddress()
		_, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    rndAddr,
			ToAddress:      cltest.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Strategy:       bulletprooftxmanager.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", rndAddr.Hex()))

		_, otherAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), *utils.NewBigI(1337), 0)

		_, err = bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    otherAddress,
			ToAddress:      cltest.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Strategy:       bulletprooftxmanager.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot send transaction on chain ID 0; eth key with address %s is pegged to chain ID 1337", otherAddress.Hex()))
	})
}

func newMockTxStrategy(t *testing.T) *bptxmmocks.TxStrategy {
	strategy := new(bptxmmocks.TxStrategy)
	strategy.Test(t)
	strategy.On("Simulate").Return(true)
	return strategy
}

func TestBulletproofTxManager_CreateEthTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewBulletproofTxManagerORM(t, db, cfg)
	etKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	thisKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)
	otherKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)

	fromAddress := thisKey.Address.Address()
	gasLimit := uint64(1000)
	toAddress := cltest.NewAddress()

	config := new(bptxmmocks.Config)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestLogger(t)
	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, ethClient, config, nil, nil, lggr)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, otherKey.Address.Address())
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)

		etx, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address.Address())))

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, thisKey.Address.Address())
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)

		etx, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address.Address())))

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 0, 42, thisKey.Address.Address())
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)

		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		etx, err := bptxm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
		strategy.AssertExpectations(t)
	})
}

func TestBulletproofTxManager_Lifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	config := new(bptxmmocks.Config)
	config.Test(t)
	kst := new(ksmocks.Eth)
	kst.Test(t)
	eventBroadcaster := new(pgmocks.EventBroadcaster)
	eventBroadcaster.Test(t)

	config.On("EthTxResendAfterThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperInterval").Return(1 * time.Hour)
	config.On("EvmMaxInFlightTransactions").Return(uint32(42))
	config.On("EvmFinalityDepth").Maybe().Return(uint32(42))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)
	kst.On("GetStatesForChain", &cltest.FixtureChainID).Return([]ethkey.State{}, nil).Once()

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened)
	lggr := logger.TestLogger(t)
	bptxm := bulletprooftxmanager.NewBulletproofTxManager(db, ethClient, config, kst, eventBroadcaster, lggr)

	head := cltest.Head(42)
	// It should not hang or panic
	bptxm.OnNewLongestChain(context.Background(), head)

	sub := new(pgmocks.Subscription)
	sub.On("Events").Return(make(<-chan pg.Event))
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)
	config.On("EvmNonceAutoSync").Return(true)
	config.On("EvmGasBumpThreshold").Return(uint64(1))

	require.NoError(t, bptxm.Start())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	bptxm.OnNewLongestChain(ctx, head)
	require.NoError(t, ctx.Err())

	keyState := cltest.MustGenerateRandomKeyState(t)

	kst.On("GetStatesForChain", &cltest.FixtureChainID).Return([]ethkey.State{keyState}, nil).Once()
	sub.On("Close").Return()
	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.timerCtx"), keyState.Address.Address()).Return(uint64(0), nil)
	config.On("TriggerFallbackDBPollInterval").Return(1 * time.Hour)
	keyChangeCh <- struct{}{}

	require.NoError(t, bptxm.Close())

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
	kst.AssertExpectations(t)
	eventBroadcaster.AssertExpectations(t)
	unsub.AwaitOrFail(t, 1*time.Second)
}

func TestBulletproofTxManager_SignTx(t *testing.T) {
	t.Parallel()

	addr := gethcommon.HexToAddress("0xb921F7763960b296B9cbAD586ff066A18D749724")
	to := gethcommon.HexToAddress("0xb921F7763960b296B9cbAD586ff066A18D749724")
	tx := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    42,
		To:       &to,
		Value:    big.NewInt(142),
		Gas:      242,
		GasPrice: big.NewInt(342),
		Data:     []byte{1, 2, 3},
	})

	t.Run("returns correct hash for non-okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		cfg := new(bptxmmocks.Config)
		cfg.Test(t)
		cfg.On("ChainType").Return(chains.ChainType(""))
		kst := new(ksmocks.Eth)
		kst.Test(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := bulletprooftxmanager.NewChainKeyStore(*chainID, cfg, kst)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.Hex())
	})
	t.Run("returns correct hash for okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		cfg := new(bptxmmocks.Config)
		cfg.Test(t)
		cfg.On("ChainType").Return(chains.ExChain)
		kst := new(ksmocks.Eth)
		kst.Test(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := bulletprooftxmanager.NewChainKeyStore(*chainID, cfg, kst)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.NotEqual(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.Hex(), "expected okex chain hash to be different from non-okex-chain hash")
		require.Equal(t, "0x1458742e3ba53316481eb18237ced517a536c1cdef61e7b7fb2a9569d84e41a6", hash.Hex())
	})
}
