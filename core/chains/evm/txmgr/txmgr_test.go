package txmgr_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/config"
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
	t *testing.T, db *sqlx.DB, ethClient evmclient.Client, estimator gas.EvmFeeEstimator, ccfg txmgr.ChainConfig, fcfg txmgr.FeeConfig, txConfig evmconfig.Transactions, dbConfig txmgr.DatabaseConfig, listenerConfig txmgr.ListenerConfig, keyStore keystore.Eth, eventBroadcaster pg.EventBroadcaster) (txmgr.TxManager, error) {
	lggr := logger.TestLogger(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewQConfig(true)), ethClient, lggr, 100*time.Millisecond, 2, 3, 2, 1000)

	// logic for building components (from evm/evm_txm.go) -------
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", fcfg.BumpTxDepth(),
		"maxInFlightTransactions", txConfig.MaxInFlight(),
		"maxQueuedTransactions", txConfig.MaxQueued(),
		"nonceAutoSync", ccfg.NonceAutoSync(),
		"limitDefault", fcfg.LimitDefault(),
	)

	return txmgr.NewTxm(
		db,
		ccfg,
		fcfg,
		txConfig,
		dbConfig,
		listenerConfig,
		ethClient,
		lggr,
		lp,
		keyStore,
		eventBroadcaster,
		estimator)
}

func TestTxm_SendNativeToken_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1).ToInt()

	config, dbConfig, evmConfig := makeConfigs(t)

	keyStore := cltest.NewKeyStore(t, db, dbConfig).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), keyStore, nil)
	require.NoError(t, err)

	_, err = txm.SendNativeToken(big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send native token to zero address")
}

func TestTxm_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	kst := cltest.NewKeyStore(t, db, cfg.Database())

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth(), 0)
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	config, dbConfig, evmConfig := makeConfigs(t)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst.Eth(), nil)
	require.NoError(t, err)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)
		evmConfig.maxQueued = uint64(1)
		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, txmgrcommon.TxUnstarted)
		assert.Equal(t, gasLimit, etx.FeeLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, big.Int(assets.NewEthValue(0)), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "eth_txes", 1)

		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, txmgrcommon.TxUnstarted)
		assert.Equal(t, gasLimit, etx.FeeLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, big.Int(assets.NewEthValue(0)), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		evmConfig.maxQueued = uint64(1)
		_, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Meta:           nil,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Txm#CreateTransaction: cannot create transaction; too many unstarted transactions in the queue (1/1). WARNING: Hitting EVM.Transactions.MaxQueued")
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		evmConfig.maxQueued = uint64(3)
		id := uuid.New()
		tx1, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)

		tx2, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)

		assert.Equal(t, tx1.GetID(), tx2.GetID())
	})

	t.Run("returns error if eth key state is missing or doesn't match chain ID", func(t *testing.T) {
		rndAddr := testutils.NewAddress()
		_, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    rndAddr,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", rndAddr.String()))

		_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth(), *utils.NewBigI(1337))

		_, err = txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    otherAddress,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot send transaction from %s on chain ID 0: eth key with address %s exists but is has not been enabled for chain 0 (enabled only for chain IDs: 1337)", otherAddress.Hex(), otherAddress.Hex()))
	})

	t.Run("simulate transmit checker", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)

		checker := txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		}
		evmConfig.maxQueued = uint64(1)
		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

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
		meta := &txmgr.TxMeta{
			JobID:         &jobID,
			RequestID:     &requestID,
			RequestTxHash: &requestTxHash,
			MaxLink:       &testDefaultMaxLink, // 1e18
			SubID:         &testDefaultSubID,
		}
		evmConfig.maxQueued = uint64(1)
		checker := txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
		}
		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           meta,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.Equal(t, meta, m)

		var c txmgr.TransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)
	})

	t.Run("forwards tx when a proper forwarder is set up", func(t *testing.T) {
		pgtest.MustExec(t, db, `DELETE FROM eth_txes`)
		pgtest.MustExec(t, db, `DELETE FROM evm_forwarders`)
		evmConfig.maxQueued = uint64(1)

		// Create mock forwarder, mock authorizedsenders call.
		form := forwarders.NewORM(db, logger.TestLogger(t), cfg.Database())
		fwdrAddr := testutils.NewAddress()
		fwdr, err := form.CreateForwarder(fwdrAddr, utils.Big(cltest.FixtureChainID))
		require.NoError(t, err)
		require.Equal(t, fwdr.Address, fwdrAddr)

		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:      fromAddress,
			ToAddress:        toAddress,
			EncodedPayload:   payload,
			FeeLimit:         gasLimit,
			ForwarderAddress: fwdr.Address,
			Strategy:         txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "eth_txes", 1)

		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.NotNil(t, m.FwdrDestAddress)
		require.Equal(t, etx.ToAddress.String(), fwdrAddr.String())
	})
}

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

type databaseConfig struct {
	config.Database
	defaultQueryTimeout time.Duration
}

func (d *databaseConfig) DefaultQueryTimeout() time.Duration {
	return d.defaultQueryTimeout
}

func (d *databaseConfig) LogSQL() bool {
	return false
}

type listenerConfig struct {
	config.Listener
}

func (l *listenerConfig) FallbackPollInterval() time.Duration {
	return 1 * time.Minute
}

func (d *databaseConfig) Listener() config.Listener {
	return &listenerConfig{}
}

type evmConfig struct {
	evmconfig.EVM
	maxInFlight          uint32
	reaperInterval       time.Duration
	reaperThreshold      time.Duration
	resendAfterThreshold time.Duration
	bumpThreshold        uint64
	maxQueued            uint64
}

func (e *evmConfig) Transactions() evmconfig.Transactions {
	return &transactionsConfig{e: e}
}

func (e *evmConfig) GasEstimator() evmconfig.GasEstimator {
	return &gasEstimatorConfig{bumpThreshold: e.bumpThreshold}
}

func (e *evmConfig) NonceAutoSync() bool { return true }

func (e *evmConfig) FinalityDepth() uint32 { return 42 }

type gasEstimatorConfig struct {
	bumpThreshold uint64
}

func (g *gasEstimatorConfig) BlockHistory() evmconfig.BlockHistory {
	return &blockHistoryConfig{}
}

func (g *gasEstimatorConfig) EIP1559DynamicFees() bool             { return false }
func (g *gasEstimatorConfig) LimitDefault() uint32                 { return 42 }
func (g *gasEstimatorConfig) BumpPercent() uint16                  { return 42 }
func (g *gasEstimatorConfig) BumpThreshold() uint64                { return g.bumpThreshold }
func (g *gasEstimatorConfig) BumpMin() *assets.Wei                 { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) FeeCapDefault() *assets.Wei           { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) PriceDefault() *assets.Wei            { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) TipCapDefault() *assets.Wei           { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) TipCapMin() *assets.Wei               { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) LimitMax() uint32                     { return 0 }
func (g *gasEstimatorConfig) LimitMultiplier() float32             { return 0 }
func (g *gasEstimatorConfig) BumpTxDepth() uint32                  { return 42 }
func (g *gasEstimatorConfig) LimitTransfer() uint32                { return 42 }
func (g *gasEstimatorConfig) PriceMax() *assets.Wei                { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) PriceMin() *assets.Wei                { return assets.NewWeiI(42) }
func (g *gasEstimatorConfig) Mode() string                         { return "FixedPrice" }
func (g *gasEstimatorConfig) LimitJobType() evmconfig.LimitJobType { return &limitJobTypeConfig{} }
func (e *gasEstimatorConfig) PriceMaxKey(addr common.Address) *assets.Wei {
	return assets.NewWeiI(42)
}

type limitJobTypeConfig struct {
}

func (l *limitJobTypeConfig) OCR() *uint32    { return ptr(uint32(0)) }
func (l *limitJobTypeConfig) OCR2() *uint32   { return ptr(uint32(0)) }
func (l *limitJobTypeConfig) DR() *uint32     { return ptr(uint32(0)) }
func (l *limitJobTypeConfig) FM() *uint32     { return ptr(uint32(0)) }
func (l *limitJobTypeConfig) Keeper() *uint32 { return ptr(uint32(0)) }
func (l *limitJobTypeConfig) VRF() *uint32    { return ptr(uint32(0)) }

type blockHistoryConfig struct {
	evmconfig.BlockHistory
}

func (b *blockHistoryConfig) BatchSize() uint32                 { return 42 }
func (b *blockHistoryConfig) BlockDelay() uint16                { return 42 }
func (b *blockHistoryConfig) BlockHistorySize() uint16          { return 42 }
func (b *blockHistoryConfig) EIP1559FeeCapBufferBlocks() uint16 { return 42 }
func (b *blockHistoryConfig) TransactionPercentile() uint16     { return 42 }

type transactionsConfig struct {
	evmconfig.Transactions
	e *evmConfig
}

func (_ *transactionsConfig) ForwardersEnabled() bool             { return true }
func (t *transactionsConfig) MaxInFlight() uint32                 { return t.e.maxInFlight }
func (t *transactionsConfig) MaxQueued() uint64                   { return t.e.maxQueued }
func (t *transactionsConfig) ReaperInterval() time.Duration       { return t.e.reaperInterval }
func (t *transactionsConfig) ReaperThreshold() time.Duration      { return t.e.reaperThreshold }
func (t *transactionsConfig) ResendAfterThreshold() time.Duration { return t.e.resendAfterThreshold }

type mockConfig struct {
	evmConfig           *evmConfig
	rpcDefaultBatchSize uint32
	finalityDepth       uint32
	finalityTag         bool
}

func (c *mockConfig) EVM() evmconfig.EVM {
	return c.evmConfig
}

func (c *mockConfig) NonceAutoSync() bool         { return true }
func (c *mockConfig) ChainType() config.ChainType { return "" }
func (c *mockConfig) FinalityDepth() uint32       { return c.finalityDepth }
func (c *mockConfig) FinalityTag() bool           { return c.finalityTag }
func (c *mockConfig) RPCDefaultBatchSize() uint32 { return c.rpcDefaultBatchSize }

func makeConfigs(t *testing.T) (*mockConfig, *databaseConfig, *evmConfig) {
	db := &databaseConfig{defaultQueryTimeout: pg.DefaultQueryTimeout}
	ec := &evmConfig{bumpThreshold: 42, maxInFlight: uint32(42), maxQueued: uint64(0), reaperInterval: time.Duration(0), reaperThreshold: time.Duration(0)}
	config := &mockConfig{evmConfig: ec}
	return config, db, ec
}

func TestTxm_CreateTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	etKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	thisKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)
	otherKey, _ := cltest.MustInsertRandomKey(t, etKeyStore, 1)

	fromAddress := thisKey.Address
	evmFromAddress := fromAddress
	gasLimit := uint32(1000)
	toAddress := testutils.NewAddress()

	config, dbConfig, evmConfig := makeConfigs(t)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	kst := cltest.NewKeyStore(t, db, dbConfig)
	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst.Eth(), nil)
	require.NoError(t, err)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)

		evmConfig.maxQueued = uint64(1)
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, pg.WithParentCtx(context.Background()))
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		evmConfig.maxQueued = uint64(1)

		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, utils.JustError(db.Exec(`DELETE FROM eth_txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 42, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)

		evmConfig.maxQueued = uint64(1)
		etx, err := txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)
		require.Equal(t, payload, etx.EncodedPayload)
	})
}

func TestTxm_Lifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	kst := ksmocks.NewEth(t)
	eventBroadcaster := pgmocks.NewEventBroadcaster(t)

	config, dbConfig, evmConfig := makeConfigs(t)
	config.finalityDepth = uint32(42)
	config.rpcDefaultBatchSize = uint32(4)

	evmConfig.resendAfterThreshold = 1 * time.Hour
	evmConfig.reaperThreshold = 1 * time.Hour
	evmConfig.reaperInterval = 1 * time.Hour

	kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return([]gethcommon.Address{}, nil)

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened)
	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst, eventBroadcaster)
	require.NoError(t, err)

	head := cltest.Head(42)
	// It should not hang or panic
	txm.OnNewLongestChain(testutils.Context(t), head)

	sub := pgmocks.NewSubscription(t)
	sub.On("Events").Return(make(<-chan pg.Event))
	eventBroadcaster.On("Subscribe", "insert_on_eth_txes", "").Return(sub, nil)
	evmConfig.bumpThreshold = uint64(1)

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
	kst := cltest.NewKeyStore(t, db, cfg.Database())

	_, addr := cltest.MustInsertRandomKey(t, kst.Eth(), 5)
	_, addr2 := cltest.MustInsertRandomKey(t, kst.Eth(), 3)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
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

	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, cfg.EVM(), cfg.EVM().GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, cfg.EVM(), cfg.EVM().GasEstimator(), cfg.EVM().Transactions(), cfg.Database(), cfg.Database().Listener(), kst.Eth(), eventBroadcaster)
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
