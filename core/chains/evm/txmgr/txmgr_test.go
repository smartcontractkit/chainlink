package txmgr_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/builder"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	pgmocks "github.com/smartcontractkit/chainlink/v2/core/services/pg/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func makeTestEvmTxm(
	t *testing.T, db *sqlx.DB, ethClient evmclient.Client, cfg txmgr.Config, keyStore keystore.Eth, eventBroadcaster pg.EventBroadcaster) (txmgr.EvmTxManager, error) {
	lggr := logger.TestLogger(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewQConfig(true)), ethClient, lggr, 100*time.Millisecond, 2, 3, 2, 1000)

	// logic for building components (from evm/evm_txm.go) -------
	lggr.Infow("Initializing EVM transaction manager",
		"gasBumpTxDepth", cfg.EvmGasBumpTxDepth(),
		"maxInFlightTransactions", cfg.EvmMaxInFlightTransactions(),
		"maxQueuedTransactions", cfg.EvmMaxQueuedTransactions(),
		"nonceAutoSync", cfg.EvmNonceAutoSync(),
		"gasLimitDefault", cfg.EvmGasLimitDefault(),
	)

	// build estimator from factory
	estimator := gas.NewEstimator(lggr, ethClient, cfg)

	return builder.NewTxm(
		db,
		cfg,
		ethClient,
		lggr,
		lp,
		keyStore,
		eventBroadcaster,
		estimator)
}

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

	keyStore := cltest.NewKeyStore(t, db, config).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	txm, err := makeTestEvmTxm(t, db, ethClient, config, keyStore, nil)
	require.NoError(t, err)

	_, err = txm.SendEther(big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestTxm_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTxStore(t, db, cfg)
	kst := cltest.NewKeyStore(t, db, cfg)

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth(), 0)
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	config := newMockConfig(t)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	txm, err := makeTestEvmTxm(t, db, ethClient, config, kst.Eth(), nil)
	require.NoError(t, err)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		etx := tx.(txmgr.EvmTx)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "eth_txes", 1)

		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		config.AssertExpectations(t)
	})

	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		_, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    fromAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Meta:           nil,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Txm#CreateEthTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting EVM.Transactions.MaxQueued")

		config.AssertExpectations(t)
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		id := uuid.NewV4()
		tx1, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgr.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		config.On("EvmMaxQueuedTransactions").Return(uint64(3)).Once()
		tx2, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgr.SendEveryStrategy{},
		})
		assert.NoError(t, err)

		assert.Equal(t, tx1.GetID(), tx2.GetID())

		config.AssertExpectations(t)
	})

	t.Run("returns error if eth key state is missing or doesn't match chain ID", func(t *testing.T) {
		rndAddr := testutils.NewAddress()
		_, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    rndAddr,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", rndAddr.String()))

		_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth(), *utils.NewBigI(1337))

		_, err = txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    otherAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgr.SendEveryStrategy{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot send transaction from %s on chain ID 0: eth key with address %s exists but is has not been enabled for chain 0 (enabled only for chain IDs: 1337)", otherAddress.Hex(), otherAddress.Hex()))

		config.AssertExpectations(t)
	})

	t.Run("simulate transmit checker", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)

		checker := txmgr.EvmTransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		}
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Strategy:       txmgr.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)
		etx := tx.(txmgr.EvmTx)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		var c txmgr.EvmTransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)

		config.AssertExpectations(t)
	})

	t.Run("meta and vrf checker", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"
		jobID := int32(25)
		requestID := gethcommon.HexToHash("abcd")
		requestTxHash := gethcommon.HexToHash("dcba")
		meta := &txmgr.EthTxMeta{
			JobID:         &jobID,
			RequestID:     &requestID,
			RequestTxHash: &requestTxHash,
			MaxLink:       &testDefaultMaxLink, // 1e18
			SubID:         &testDefaultSubID,
		}
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()
		checker := txmgr.EvmTransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
		}
		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           meta,
			Strategy:       txmgr.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)
		etx := tx.(txmgr.EvmTx)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.Equal(t, meta, m)

		var c txmgr.EvmTransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)

		config.AssertExpectations(t)
	})

	t.Run("forwards tx when a proper forwarder is set up", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
		pgtest.MustExec(t, db, `DELETE FROM evm_forwarders`)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1)).Once()

		// Create mock forwarder, mock authorizedsenders call.
		form := forwarders.NewORM(db, logger.TestLogger(t), cfg)
		fwdrAddr := testutils.NewAddress()
		fwdr, err := form.CreateForwarder(fwdrAddr, utils.Big(cltest.FixtureChainID))
		require.NoError(t, err)
		require.Equal(t, fwdr.Address, fwdrAddr)

		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:      fromAddress,
			ToAddress:        toAddress,
			EncodedPayload:   payload,
			FeeLimit:         gasLimit,
			ForwarderAddress: fwdr.Address,
			Strategy:         txmgr.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)

		etx := tx.(txmgr.EvmTx)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.NotNil(t, m.FwdrDestAddress)
		require.Equal(t, etx.ToAddress.String(), fwdrAddr.String())

		config.AssertExpectations(t)
	})
}

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func newMockConfig(t *testing.T) *txmmocks.Config {
	// These are only used for logging, the exact value doesn't matter
	// It can be overridden in the test that uses it
	cfg := txmmocks.NewConfig(t)
	cfg.On("EvmGasBumpTxDepth").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmMaxInFlightTransactions").Return(uint32(42)).Maybe()
	cfg.On("EvmMaxQueuedTransactions").Return(uint64(42)).Maybe().Once()
	cfg.On("EvmNonceAutoSync").Return(true).Maybe()
	cfg.On("EvmGasLimitDefault").Return(uint32(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBatchSize").Return(uint32(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBlockDelay").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorBlockHistorySize").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").Return(uint16(42)).Maybe().Once()
	cfg.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmEIP1559DynamicFees").Return(false).Maybe().Twice()
	cfg.On("EvmGasBumpPercent").Return(uint16(42)).Maybe().Once()
	cfg.On("EvmGasBumpThreshold").Return(uint64(42)).Maybe()
	cfg.On("EvmGasBumpWei").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmGasFeeCapDefault").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmGasLimitMultiplier").Return(float32(42)).Maybe().Once()
	cfg.On("EvmGasPriceDefault").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmGasTipCapDefault").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmGasTipCapMinimum").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmMaxGasPriceWei").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmMinGasPriceWei").Return(assets.NewWeiI(42)).Maybe().Once()
	cfg.On("EvmUseForwarders").Return(true).Maybe()
	cfg.On("LogSQL").Maybe().Return(false)
	cfg.On("DatabaseDefaultQueryTimeout").Return(pg.DefaultQueryTimeout).Maybe()

	return cfg
}

func TestTxm_CreateEthTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTxStore(t, db, cfg)
	etKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	thisKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)
	otherKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)

	fromAddress := thisKey.Address
	evmFromAddress := fromAddress
	gasLimit := uint32(1000)
	toAddress := testutils.NewAddress()

	config := newMockConfig(t)
	config.On("EthTxResendAfterThreshold").Return(time.Duration(0))
	config.On("EthTxReaperThreshold").Return(time.Duration(0))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	kst := cltest.NewKeyStore(t, db, cfg)
	txm, err := makeTestEvmTxm(t, db, ethClient, config, kst.Eth(), nil)
	require.NoError(t, err)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, pg.WithParentCtx(context.Background()))
		assert.NoError(t, err)

		etx := tx.(txmgr.EvmTx)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		etx := tx.(txmgr.EvmTx)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 42, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		config.On("EvmMaxQueuedTransactions").Return(uint64(1))
		tx, err := txm.CreateEthTransaction(txmgr.EvmNewTx{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		etx := tx.(txmgr.EvmTx)
		require.Equal(t, payload, etx.EncodedPayload)
	})
}

func TestTxm_Lifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	config := newMockConfig(t)
	kst := ksmocks.NewEth(t)
	eventBroadcaster := pgmocks.NewEventBroadcaster(t)

	config.On("EthTxResendAfterThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperThreshold").Return(1 * time.Hour)
	config.On("EthTxReaperInterval").Return(1 * time.Hour)
	config.On("EvmMaxInFlightTransactions").Return(uint32(42))
	config.On("EvmFinalityDepth").Maybe().Return(uint32(42))
	config.On("GasEstimatorMode").Return("FixedPrice")
	config.On("LogSQL").Return(false).Maybe()
	config.On("EvmRPCDefaultBatchSize").Return(uint32(4)).Maybe()
	kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return([]gethcommon.Address{}, nil)

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened)
	txm, err := makeTestEvmTxm(t, db, ethClient, config, kst, eventBroadcaster)
	require.NoError(t, err)

	head := cltest.Head(42)
	// It should not hang or panic
	txm.OnNewLongestChain(testutils.Context(t), head)

	sub := pgmocks.NewSubscription(t)
	sub.On("Events").Return(make(<-chan pg.Event))
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)
	config.On("EvmGasBumpThreshold").Return(uint64(1))

	require.NoError(t, txm.Start(testutils.Context(t)))

	ctx, cancel := context.WithTimeout(testutils.Context(t), 5*time.Second)
	t.Cleanup(cancel)
	txm.OnNewLongestChain(ctx, head)
	require.NoError(t, ctx.Err())

	keyState := cltest.MustGenerateRandomKeyState(t)

	addr := []gethcommon.Address{keyState.Address.Address()}
	kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return(addr, nil)
	sub.On("Close").Return()
	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.cancelCtx"), gethcommon.Address{}).Return(uint64(0), nil).Maybe()
	config.On("TriggerFallbackDBPollInterval").Return(1 * time.Hour).Maybe()
	keyChangeCh <- struct{}{}

	require.NoError(t, txm.Close())
	unsub.AwaitOrFail(t, 1*time.Second)
}

type fnMock struct{ called atomic.Bool }

func (fm *fnMock) Fn() {
	swapped := fm.called.CompareAndSwap(false, true)
	if !swapped {
		panic("func called more than once")
	}
}

func (fm *fnMock) AssertNotCalled(t *testing.T) {
	assert.False(t, fm.called.Load())
}

func (fm *fnMock) AssertCalled(t *testing.T) {
	assert.True(t, fm.called.Load())
}

func TestTxm_Reset(t *testing.T) {
	t.Parallel()

	// Lots of boilerplate setup since we actually want to test start/stop of EthBroadcaster/EthConfirmer
	db := pgtest.NewSqlxDB(t)
	gcfg := configtest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := cltest.NewKeyStore(t, db, cfg)

	_, addr := cltest.MustInsertRandomKey(t, kst.Eth(), 5)
	_, addr2 := cltest.MustInsertRandomKey(t, kst.Eth(), 3)
	txStore := cltest.NewTxStore(t, db, cfg)
	// 4 confirmed tx from addr1
	for i := int64(0); i < 4; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr)
	}
	// 2 confirmed from addr2
	for i := int64(0); i < 2; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr2)
	}

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), nil)
	ethClient.On("PendingNonceAt", mock.Anything, addr2).Return(uint64(0), nil)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, nil)
	ethClient.On("BatchCallContextAll", mock.Anything, mock.Anything).Return(nil).Maybe()
	eventBroadcaster := pgmocks.NewEventBroadcaster(t)
	sub := pgmocks.NewSubscription(t)
	sub.On("Events").Return(make(<-chan pg.Event))
	sub.On("Close")
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)

	txm, err := makeTestEvmTxm(t, db, ethClient, cfg, kst.Eth(), eventBroadcaster)
	require.NoError(t, err)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, addr2)
	for i := 0; i < 1000; i++ {
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 4+int64(i), addr)
	}

	t.Run("returns error if not started", func(t *testing.T) {
		f := new(fnMock)

		err := txm.Reset(f.Fn, addr, false)
		require.Error(t, err)
		assert.EqualError(t, err, "not started")

		f.AssertNotCalled(t)
	})

	require.NoError(t, txm.Start(testutils.Context(t)))
	defer func() { assert.NoError(t, txm.Close()) }()

	t.Run("calls function if started", func(t *testing.T) {
		f := new(fnMock)

		err := txm.Reset(f.Fn, addr, false)
		require.NoError(t, err)

		f.AssertCalled(t)
	})

	t.Run("calls function and deletes relevant eth_txes if abandon=true", func(t *testing.T) {
		f := new(fnMock)

		err := txm.Reset(f.Fn, addr, true)
		require.NoError(t, err)

		f.AssertCalled(t)

		var s string
		err = db.Get(&s, `SELECT error FROM eth_txes WHERE from_address = $1 AND state = 'fatal_error'`, addr)
		require.NoError(t, err)
		assert.Equal(t, "abandoned", s)

		// the other address didn't get touched
		var count int
		err = db.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = 'fatal_error'`, addr2)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
