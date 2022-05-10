package txmgr_test

import (
	"context"
	"encoding/json"
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
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	pgmocks "github.com/smartcontractkit/chainlink/core/services/pg/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestTxm_SendEther_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1)

	config := newMockConfig(t)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestLogger(t)
	checkerFactory := &testCheckerFactory{}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewPGCfg(true)),
		ethClient, lggr, 100*time.Millisecond, 2, 3)
	txm := txmgr.NewTxm(db, ethClient, config, nil, nil, lggr, checkerFactory, lp)

	_, err := txm.SendEther(big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestTxm_CheckEthTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertFatalErrorEthTx(t, borm, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	var n int64 = 0
	cltest.MustInsertInProgressEthTxWithAttempt(t, borm, n, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, 1, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS", maxUnconfirmedTransactions))

		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
		err = txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS", maxUnconfirmedTransactions))
	})

	t.Run("with different chain ID ignores txes", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions, *big.NewInt(42))
		require.NoError(t, err)
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := txmgr.CheckEthTxQueueCapacity(db, fromAddress, 0, cltest.FixtureChainID)
		require.NoError(t, err)
	})
}

func TestTxm_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	q := pg.NewQ(db, logger.TestLogger(t), cfg)
	count, err := txmgr.CountUnconfirmedTransactions(q, fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestTxm_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	q := pg.NewQ(db, logger.TestLogger(t), cfg)
	count, err := txmgr.CountUnstartedTransactions(q, fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}
func TestTxm_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	keyStore := cltest.NewKeyStore(t, db, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)
	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config := newMockConfig(t)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	lggr := logger.TestLogger(t)
	checkerFactory := &testCheckerFactory{}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewPGCfg(true)),
		ethClient, lggr, 100*time.Millisecond, 2, 3)
	txm := txmgr.NewTxm(db, ethClient, config, nil, nil, lggr, checkerFactory, lp)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, db.Get(&etx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
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
		_, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Meta:           nil,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Txm#CreateEthTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting ETH_MAX_QUEUED_TRANSACTIONS")
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		id := uuid.NewV4()
		tx1, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			GasLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgr.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		tx2, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			GasLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgr.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		assert.Equal(t, tx1.ID, tx2.ID)
	})

	t.Run("returns error if eth key state is missing or doesn't match chain ID", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Twice()
		rndAddr := testutils.NewAddress()
		_, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    rndAddr,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", rndAddr.Hex()))

		_, otherAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), *utils.NewBigI(1337), 0)

		_, err = txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    otherAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			GasLimit:       21000,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot send transaction on chain ID 0; eth key with address %s is pegged to chain ID 1337", otherAddress.Hex()))
	})

	t.Run("simulate transmit checker", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)

		checker := txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		}
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Strategy:       txmgr.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, db.Get(&etx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		var c txmgr.TransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)
	})

	t.Run("meta and vrf checker", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"
		jobID := int32(25)
		requestID := gethcommon.HexToHash("abcd")
		requestTxHash := gethcommon.HexToHash("dcba")
		meta := &txmgr.EthTxMeta{
			JobID:         jobID,
			RequestID:     requestID,
			RequestTxHash: requestTxHash,
			MaxLink:       &testDefaultMaxLink, // 1e18
			SubID:         &testDefaultSubID,
		}
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		checker := txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddress(),
		}
		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           meta,
			Strategy:       txmgr.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, db.Get(&etx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.Equal(t, meta, m)

		var c txmgr.TransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)
	})
}

func newMockTxStrategy(t *testing.T) *txmmocks.TxStrategy {
	strategy := new(txmmocks.TxStrategy)
	strategy.Test(t)
	return strategy
}

func newMockConfig(t *testing.T) *txmmocks.Config {
	// These are only used for logging, the exact value doesn't matter
	// It can be overridden in the test that uses it
	cfg := new(txmmocks.Config)
	cfg.Test(t)
	cfg.On("EvmGasBumpTxDepth").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmMaxInFlightTransactions").Return(uint32(42)).Maybe().Once()
	cfg.On("EvmMaxQueuedTransactions").Return(uint64(42)).Maybe().Once()
	cfg.On("EvmNonceAutoSync").Return(true).Maybe()
	cfg.On("EvmGasLimitDefault").Return(uint64(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBatchSize").Return(uint32(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBlockDelay").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBlockHistorySize").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmEIP1559DynamicFees").Return(false).Maybe().Once()
	cfg.On("EvmGasBumpPercent").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmGasBumpThreshold").Return(uint64(42)).Maybe().Once()
	cfg.On("EvmGasBumpWei").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmGasFeeCapDefault").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmGasLimitMultiplier").Return(float32(42)).Maybe().Once()
	cfg.On("EvmGasPriceDefault").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmGasTipCapDefault").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmGasTipCapMinimum").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmMaxGasPriceWei").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmMinGasPriceWei").Return(big.NewInt(42)).Maybe().Once()
	cfg.On("EvmUseForwarders").Return(false).Maybe()
	cfg.On("LogSQL").Maybe().Return(false)

	return cfg
}

func TestTxm_CreateEthTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	etKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	thisKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)
	otherKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)

	fromAddress := thisKey.Address.Address()
	gasLimit := uint64(1000)
	toAddress := testutils.NewAddress()

	config := newMockConfig(t)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestLogger(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewPGCfg(true)),
		ethClient, lggr, 100*time.Millisecond, 2, 3)
	txm := txmgr.NewTxm(db, ethClient, config, nil, nil, lggr, &testCheckerFactory{}, lp)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, otherKey.Address.Address())
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.AnythingOfType("*sqlx.Tx")).Return(int64(0), nil)

		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
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

		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
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
		etx, err := txm.CreateEthTransaction(txmgr.NewTx{
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

func TestTxm_Lifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	config := newMockConfig(t)
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
	config.On("LogSQL").Return(false).Maybe()
	config.On("EvmRPCDefaultBatchSize").Return(uint32(4)).Maybe()
	kst.On("GetStatesForChain", &cltest.FixtureChainID).Return([]ethkey.State{}, nil).Once()

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened)
	lggr := logger.TestLogger(t)
	checkerFactory := &testCheckerFactory{}

	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewPGCfg(true)),
		ethClient, lggr, 100*time.Millisecond, 2, 3)
	txm := txmgr.NewTxm(db, ethClient, config, kst, eventBroadcaster, lggr, checkerFactory, lp)

	head := cltest.Head(42)
	// It should not hang or panic
	txm.OnNewLongestChain(context.Background(), head)

	sub := new(pgmocks.Subscription)
	sub.On("Events").Return(make(<-chan pg.Event))
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)
	config.On("EvmGasBumpThreshold").Return(uint64(1))

	require.NoError(t, txm.Start(testutils.Context(t)))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	txm.OnNewLongestChain(ctx, head)
	require.NoError(t, ctx.Err())

	keyState := cltest.MustGenerateRandomKeyState(t)

	kst.On("GetStatesForChain", &cltest.FixtureChainID).Return([]ethkey.State{keyState}, nil).Once()
	sub.On("Close").Return()
	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.cancelCtx"), keyState.Address.Address()).Return(uint64(0), nil).Maybe()
	config.On("TriggerFallbackDBPollInterval").Return(1 * time.Hour).Maybe()
	keyChangeCh <- struct{}{}

	require.NoError(t, txm.Close())

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
	kst.AssertExpectations(t)
	eventBroadcaster.AssertExpectations(t)
	unsub.AwaitOrFail(t, 1*time.Second)
}

func TestTxm_SignTx(t *testing.T) {
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
		cfg := new(txmmocks.Config)
		cfg.Test(t)
		cfg.On("ChainType").Return(config.ChainType(""))
		kst := new(ksmocks.Eth)
		kst.Test(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewChainKeyStore(*chainID, cfg, kst)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.Hex())
	})
	// okex used to have a custom hash but now this just verifies that is it the same
	t.Run("returns correct hash for okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		cfg := new(txmmocks.Config)
		cfg.Test(t)
		kst := new(ksmocks.Eth)
		kst.Test(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewChainKeyStore(*chainID, cfg, kst)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.Hex())
	})
}
