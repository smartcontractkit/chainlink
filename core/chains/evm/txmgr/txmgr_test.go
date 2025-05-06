package txmgr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	evmtxm "github.com/smartcontractkit/chainlink-evm/pkg/txm"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func makeTestEvmTxm(
	t testing.TB, db *sqlx.DB, ethClient evmclient.Client, estimator gas.EvmFeeEstimator, ccfg txmgr.ChainConfig, fcfg txmgr.FeeConfig, txConfig evmconfig.Transactions, dbConfig txmgr.DatabaseConfig, listenerConfig txmgr.ListenerConfig, keyStore keys.ChainStore) (txmgr.TxManager, error) {
	lggr := logger.Test(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}

	ht := headstest.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), ethClient, lggr, ht, lpOpts)

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
		nil,
		dbConfig,
		listenerConfig,
		ethClient,
		lggr,
		lp,
		keyStore,
		estimator,
		ht,
		nil)
}

func TestTxm_SendNativeToken_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := testutils.NewSqlxDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1).ToInt()

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	estimator, err := gas.NewEstimator(logger.Test(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), &keystest.FakeChainStore{})
	require.NoError(t, err)

	_, err = txm.SendNativeToken(tests.Context(t), big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send native token to zero address")
}

func TestTxm_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))

	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	ethClient := clienttest.NewClientWithDefaultChainID(t)

	estimator, err := gas.NewEstimator(logger.Test(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), ethKeyStore)
	require.NoError(t, err)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)
		evmConfig.MaxQueued = uint64(1)
		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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

		cltest.AssertCount(t, db, "evm.txes", 1)

		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, txmgrcommon.TxUnstarted)
		assert.Equal(t, gasLimit, etx.FeeLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, big.Int(assets.NewEthValue(0)), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)

	t.Run("with queue at capacity does not insert eth_tx", func(t *testing.T) {
		evmConfig.MaxQueued = uint64(1)
		_, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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
		evmConfig.MaxQueued = uint64(3)
		id := uuid.New()
		tx1, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)

		tx2, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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

	t.Run("returns error if eth key is not enabled", func(t *testing.T) {
		rndAddr := testutils.NewAddress()
		_, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    rndAddr,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, evmtxm.NotEnabledError{})
		assert.ErrorIs(t, err, evmtxm.NotEnabledError{FromAddress: rndAddr})
		var as evmtxm.NotEnabledError
		if assert.ErrorAs(t, err, &as) {
			assert.Equal(t, rndAddr.String(), as.FromAddress.String())
		}
	})

	t.Run("simulate transmit checker", func(t *testing.T) {
		testutils.MustExec(t, db, `DELETE FROM evm.txes`)

		checker := txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		}
		evmConfig.MaxQueued = uint64(1)
		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "evm.txes", 1)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		var c txmgr.TransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)
	})

	t.Run("meta and vrf checker", func(t *testing.T) {
		testutils.MustExec(t, db, `DELETE FROM evm.txes`)
		testDefaultSubID := uint64(2)
		testDefaultMaxLink := "1000000000000000000"
		testDefaultMaxEth := "2000000000000000000"
		// max uint256 is 1.1579209e+77
		testDefaultGlobalSubID := crypto.Keccak256Hash([]byte("sub id")).String()
		jobID := int32(25)
		requestID := common.HexToHash("abcd")
		requestTxHash := common.HexToHash("dcba")
		meta := &txmgr.TxMeta{
			JobID:         &jobID,
			RequestID:     &requestID,
			RequestTxHash: &requestTxHash,
			MaxLink:       &testDefaultMaxLink, // 1e18
			MaxEth:        &testDefaultMaxEth,  // 2e18
			SubID:         &testDefaultSubID,
			GlobalSubID:   &testDefaultGlobalSubID,
		}
		evmConfig.MaxQueued = uint64(1)
		checker := txmgr.TransmitCheckerSpec{
			CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
			VRFCoordinatorAddress: testutils.NewAddressPtr(),
		}
		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           meta,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
			Checker:        checker,
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "evm.txes", 1)
		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)
		require.Equal(t, meta, m)

		var c txmgr.TransmitCheckerSpec
		require.NotNil(t, etx.TransmitChecker)
		require.NoError(t, json.Unmarshal(*etx.TransmitChecker, &c))
		require.Equal(t, checker, c)
	})

	t.Run("forwards tx when a proper forwarder is set up", func(t *testing.T) {
		testutils.MustExec(t, db, `DELETE FROM evm.txes`)
		testutils.MustExec(t, db, `DELETE FROM evm.forwarders`)
		evmConfig.MaxQueued = uint64(1)

		// Create mock forwarder, mock authorizedsenders call.
		form := forwarders.NewORM(db)
		fwdrAddr := testutils.NewAddress()
		fwdr, err := form.CreateForwarder(tests.Context(t), fwdrAddr, ubig.Big(cltest.FixtureChainID))
		require.NoError(t, err)
		require.Equal(t, fwdr.Address, fwdrAddr)

		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:      fromAddress,
			ToAddress:        toAddress,
			EncodedPayload:   payload,
			FeeLimit:         gasLimit,
			ForwarderAddress: fwdr.Address,
			Strategy:         txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)
		cltest.AssertCount(t, db, "evm.txes", 1)

		var dbEtx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEtx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		m, err := etx.GetMeta()
		require.NoError(t, err)

		require.NotEqual(t, etx.ToAddress, *m.FwdrDestAddress)
		require.Equal(t, toAddress, *m.FwdrDestAddress)
		require.NotNil(t, m.FwdrDestAddress)
		require.Equal(t, etx.ToAddress.String(), fwdrAddr.String())
	})

	t.Run("insert Tx successfully with a IdempotencyKey", func(t *testing.T) {
		evmConfig.MaxQueued = uint64(3)
		id := uuid.New()
		idempotencyKey := "1"
		_, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			IdempotencyKey:    &idempotencyKey,
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that IdempotencyKey", func(t *testing.T) {
		evmConfig.MaxQueued = uint64(3)
		id := uuid.New()
		idempotencyKey := "2"
		tx1, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			IdempotencyKey:    &idempotencyKey,
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		})
		assert.NoError(t, err)

		tx2, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			IdempotencyKey:    &idempotencyKey,
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
}

func BenchmarkCreateTransaction(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	//txStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(b, db)

	_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(b)

	ethClient := clienttest.NewClient(b)
	ethClient.On("ConfiguredChainID").Return(big.NewInt(0)).Maybe()

	estimator, err := gas.NewEstimator(logger.Test(b), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(b, err)
	txm, err := makeTestEvmTxm(b, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), keys.NewChainStore(keystore.NewEthSigner(kst.Eth(), new(big.Int)), new(big.Int)))
	require.NoError(b, err)

	subject := uuid.New()
	strategy := newMockTxStrategy(b)
	strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
	strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)
	for n := 0; n < b.N; n++ {
		txm.CreateTransaction(tests.Context(b), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
	}
}

func newMockTxStrategy(t testing.TB) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func TestTxm_CreateTransaction_OutOfEth(t *testing.T) {
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	otherAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))

	gasLimit := uint64(1000)
	toAddress := testutils.NewAddress()

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	estimator, err := gas.NewEstimator(logger.Test(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), ethKeyStore)
	require.NoError(t, err)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)

		evmConfig.MaxQueued = uint64(1)
		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherAddress)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		assert.NoError(t, err)

		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, commonutils.JustError(db.Exec(`DELETE FROM evm.txes WHERE from_address = $1`, fromAddress)))

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		evmConfig.MaxQueued = uint64(1)

		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		require.NoError(t, err)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, commonutils.JustError(db.Exec(`DELETE FROM evm.txes WHERE from_address = $1`, fromAddress)))

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 42, fromAddress)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		evmConfig.MaxQueued = uint64(1)
		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		require.NoError(t, err)
		require.Equal(t, payload, etx.EncodedPayload)
	})
}

func TestTxm_Lifecycle(t *testing.T) {
	db := testutils.NewSqlxDB(t)

	ethClient := clienttest.NewClientWithDefaultChainID(t)

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	config.SetFinalityDepth(uint32(42))

	evmConfig.RpcDefaultBatchSize = uint32(4)
	evmConfig.ResendAfterThreshold = 1 * time.Hour
	evmConfig.ReaperThreshold = 1 * time.Hour
	evmConfig.ReaperInterval = 1 * time.Hour

	kst := &keystest.FakeChainStore{}

	head := cltest.Head(42)
	finalizedHead := cltest.Head(0)

	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(finalizedHead, nil).Maybe()

	estimator, err := gas.NewEstimator(logger.Test(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst)
	require.NoError(t, err)

	// It should not hang or panic
	txm.OnNewLongestChain(tests.Context(t), head)

	evmConfig.BumpThreshold = uint64(1)

	require.NoError(t, txm.Start(tests.Context(t)))

	ctx, cancel := context.WithTimeout(tests.Context(t), 5*time.Second)
	t.Cleanup(cancel)
	txm.OnNewLongestChain(ctx, head)
	require.NoError(t, ctx.Err())

	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.cancelCtx"), common.Address{}).Return(uint64(0), nil).Maybe()

	require.NoError(t, txm.Close())
}

func TestTxm_Reset(t *testing.T) {
	t.Parallel()

	// Lots of boilerplate setup since we actually want to test start/stop of EthBroadcaster/EthConfirmer
	db := testutils.NewSqlxDB(t)

	_, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	memKS := keystest.NewMemoryChainStore()
	addr := memKS.MustCreate(t)
	addr2 := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))

	txStore := cltest.NewTestTxStore(t, db)
	// 4 confirmed tx from addr1
	for i := int64(0); i < 4; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr)
	}
	// 2 confirmed from addr2
	for i := int64(0); i < 2; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr2)
	}

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, nil).Maybe()
	ethClient.On("BatchCallContextAll", mock.Anything, mock.Anything).Return(nil).Maybe()
	ethClient.On("NonceAt", mock.Anything, addr, mock.Anything).Return(uint64(128), nil).Maybe()
	ethClient.On("NonceAt", mock.Anything, addr2, mock.Anything).Return(uint64(44), nil).Maybe()

	estimator, err := gas.NewEstimator(logger.Test(t), ethClient, evmConfig.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), ethKeyStore)
	require.NoError(t, err)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, addr2)
	for i := 0; i < 1000; i++ {
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 4+int64(i), addr)
	}

	t.Run("returns error if not started", func(t *testing.T) {
		err := txm.Reset(addr, false)
		require.Error(t, err)
		assert.EqualError(t, err, "not started")
	})

	servicetest.Run(t, txm)

	t.Run("returns no error if started", func(t *testing.T) {
		err := txm.Reset(addr, false)
		require.NoError(t, err)
	})

	t.Run("deletes relevant evm.txes if abandon=true", func(t *testing.T) {
		err := txm.Reset(addr, true)
		require.NoError(t, err)

		var s string
		err = db.Get(&s, `SELECT error FROM evm.txes WHERE from_address = $1 AND state = 'fatal_error'`, addr)
		require.NoError(t, err)
		assert.Equal(t, "abandoned", s)

		// the other address didn't get touched
		var count int
		err = db.Get(&count, `SELECT count(*) FROM evm.txes WHERE from_address = $1 AND state = 'fatal_error'`, addr2)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestTxm_GetTransactionFee(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	memKS := keystest.NewMemoryChainStore()
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))
	feeLimit := uint64(10_000)

	_, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	h99 := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	h99.IsFinalized.Store(true)
	head := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(h99)

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head.Parent.Load(), nil).Once()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	feeEstimator.On("Start", mock.Anything).Return(nil).Once()
	feeEstimator.On("Close", mock.Anything).Return(nil).Once()
	feeEstimator.On("OnNewLongestChain", mock.Anything, mock.Anything).Once()
	txm, err := makeTestEvmTxm(t, db, ethClient, feeEstimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), ethKeyStore)
	require.NoError(t, err)
	servicetest.Run(t, txm)

	txm.OnNewLongestChain(ctx, head)

	t.Run("returns error if receipt not found", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		_, err := txm.GetTransactionFee(ctx, idempotencyKey)
		require.Error(t, err, fmt.Sprintf("failed to find receipt with IdempotencyKey: %s", idempotencyKey))
	})

	t.Run("returns error for unstarted state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		tx := &txmgr.Tx{
			IdempotencyKey: &idempotencyKey,
			FromAddress:    fromAddress,
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       feeLimit,
			State:          txmgrcommon.TxUnstarted,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)

		attemptD := cltest.NewDynamicFeeEthTxAttempt(t, tx.ID)
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attemptD))

		// insert receipt
		var r txmgr.Receipt
		r = newEthReceipt(42, utils.NewHash(), attemptD.Hash, 0x1)
		_, err = txStore.InsertReceipt(ctx, &r.Receipt)
		require.NoError(t, err)

		_, err = txm.GetTransactionFee(ctx, idempotencyKey)
		require.Error(t, err, "tx status is not finalized")
	})

	t.Run("returns correct fee for finalized state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)

		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxFinalized,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)

		attemptD := cltest.NewDynamicFeeEthTxAttempt(t, tx.ID)
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attemptD))

		// insert receipt
		var r txmgr.Receipt
		r = newEthReceipt(42, utils.NewHash(), attemptD.Hash, 0x1)
		expFee := r.Receipt.EffectiveGasPrice.Int64() * int64(r.Receipt.GasUsed)
		_, err = txStore.InsertReceipt(ctx, &r.Receipt)
		require.NoError(t, err)

		fee, err := txm.GetTransactionFee(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(expFee), fee.TransactionFee)
	})
}

func TestTxm_GetTransactionStatus(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	memKS := keystest.NewMemoryChainStore()
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))
	feeLimit := uint64(10_000)

	_, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	h99 := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	h99.IsFinalized.Store(true)
	head := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(h99)

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head.Parent.Load(), nil).Once()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	feeEstimator.On("Start", mock.Anything).Return(nil).Once()
	feeEstimator.On("Close", mock.Anything).Return(nil).Once()
	feeEstimator.On("OnNewLongestChain", mock.Anything, mock.Anything).Once()
	txm, err := makeTestEvmTxm(t, db, ethClient, feeEstimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), ethKeyStore)
	require.NoError(t, err)
	servicetest.Run(t, txm)

	txm.OnNewLongestChain(ctx, head)

	t.Run("returns error if transaction not found", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.Error(t, err, fmt.Sprintf("failed to find transaction with IdempotencyKey: %s", idempotencyKey))
		require.Equal(t, commontypes.Unknown, state)
	})

	t.Run("returns unknown for unstarted state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		tx := &txmgr.Tx{
			IdempotencyKey: &idempotencyKey,
			FromAddress:    fromAddress,
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       feeLimit,
			State:          txmgrcommon.TxUnstarted,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Unknown, state)
	})

	t.Run("returns unknown for in-progress state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		nonce := evmtypes.Nonce(0)
		tx := &txmgr.Tx{
			Sequence:       &nonce,
			IdempotencyKey: &idempotencyKey,
			FromAddress:    fromAddress,
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       feeLimit,
			State:          txmgrcommon.TxInProgress,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Unknown, state)
	})

	t.Run("returns pending for unconfirmed state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxUnconfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Pending, state)
	})

	t.Run("returns unconfirmed for confirmed state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		attempt := cltest.NewLegacyEthTxAttempt(t, tx.ID)
		err = txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)
		// Insert receipt for unfinalized block num
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attempt.Hash)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Unconfirmed, state)
	})

	t.Run("returns finalized for finalized state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxFinalized,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		attempt := cltest.NewLegacyEthTxAttempt(t, tx.ID)
		err = txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)
		// Insert receipt for finalized block num
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, head.Parent.Load().Hash, attempt.Hash)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Finalized, state)
	})

	t.Run("returns pending for confirmed missing receipt state", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmedMissingReceipt,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.NoError(t, err)
		require.Equal(t, commontypes.Pending, state)
	})

	t.Run("returns fatal for fatal error state with terminally stuck error", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		// Test the internal terminally stuck error returns Fatal
		nonce := evmtypes.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxFatalError,
			Error:              null.NewString(evmclient.TerminallyStuckMsg, true),
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.Equal(t, commontypes.Fatal, state)
		require.Error(t, err)
		require.Equal(t, evmclient.TerminallyStuckMsg, err.Error())

		// Test a terminally stuck client error returns Fatal
		nonce = evmtypes.Nonce(1)
		idempotencyKey = uuid.New().String()
		terminallyStuckClientError := "failed to add tx to the pool: not enough step counters to continue the execution"
		tx = &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxFatalError,
			Error:              null.NewString(terminallyStuckClientError, true),
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err = txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err = txm.GetTransactionStatus(ctx, idempotencyKey)
		require.Equal(t, commontypes.Fatal, state)
		require.Error(t, err)
		require.Equal(t, terminallyStuckClientError, err.Error())
	})

	t.Run("returns failed for fatal error state with other error", func(t *testing.T) {
		idempotencyKey := uuid.New().String()
		fromAddress := memKS.MustCreate(t)
		errorMsg := "something went wrong"
		tx := &txmgr.Tx{
			IdempotencyKey: &idempotencyKey,
			FromAddress:    fromAddress,
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       feeLimit,
			State:          txmgrcommon.TxFatalError,
			Error:          null.NewString(errorMsg, true),
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		state, err := txm.GetTransactionStatus(ctx, idempotencyKey)
		require.Equal(t, commontypes.Failed, state)
		require.Error(t, err, errorMsg)
	})
}

func newTxStore(t testing.TB, db *sqlx.DB) txmgr.EvmTxStore {
	return txmgr.NewTxStore(db, logger.Test(t))
}

func newEthReceipt(blockNumber int64, blockHash common.Hash, txHash common.Hash, status uint64) txmgr.Receipt {
	transactionIndex := uint(cltest.NewRandomPositiveInt64())

	receipt := evmtypes.Receipt{
		BlockNumber:       big.NewInt(blockNumber),
		BlockHash:         blockHash,
		TxHash:            txHash,
		TransactionIndex:  transactionIndex,
		GasUsed:           123,
		EffectiveGasPrice: big.NewInt(55),
		Status:            status,
	}

	r := txmgr.Receipt{
		BlockNumber:      blockNumber,
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
		Receipt:          receipt,
	}
	return r
}

func mustInsertEthReceipt(t testing.TB, txStore txmgr.TestEvmTxStore, blockNumber int64, blockHash common.Hash, txHash common.Hash) txmgr.Receipt {
	r := newEthReceipt(blockNumber, blockHash, txHash, 0x1)
	id, err := txStore.InsertReceipt(tests.Context(t), &r.Receipt)
	require.NoError(t, err)
	r.ID = id
	return r
}

func mustInsertRevertedEthReceipt(t *testing.T, txStore txmgr.TestEvmTxStore, blockNumber int64, blockHash common.Hash, txHash common.Hash) txmgr.Receipt {
	r := newEthReceipt(blockNumber, blockHash, txHash, 0x0)
	id, err := txStore.InsertReceipt(tests.Context(t), &r.Receipt)
	require.NoError(t, err)
	r.ID = id
	return r
}

// Inserts into evm.receipts but does not update evm.txes or evm.tx_attempts
func mustInsertConfirmedEthTxWithReceipt(t testing.TB, txStore txmgr.TestEvmTxStore, fromAddress common.Address, nonce, blockNum int64) (etx txmgr.Tx) {
	etx = cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, nonce, blockNum, fromAddress)
	mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx.TxAttempts[0].Hash)
	return etx
}

func mustInsertFatalErrorEthTx(t testing.TB, txStore txmgr.TestEvmTxStore, fromAddress common.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.Error = null.StringFrom("something exploded")
	etx.State = txmgrcommon.TxFatalError

	require.NoError(t, txStore.InsertTx(tests.Context(t), &etx))
	return etx
}

func mustInsertUnconfirmedEthTxWithAttemptState(t testing.TB, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, txAttemptState txmgrtypes.TxAttemptState, opts ...interface{}) txmgr.Tx {
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress, opts...)
	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
	ctx := tests.Context(t)

	tx := cltest.NewLegacyTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txAttemptState
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	var err error
	etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.Tx {
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress, opts...)
	attempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
	ctx := tests.Context(t)

	addr := testutils.NewAddress()
	dtx := types.DynamicFeeTx{
		ChainID:   big.NewInt(0),
		Nonce:     uint64(nonce),
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1),
		Gas:       242,
		To:        &addr,
		Value:     big.NewInt(342),
		Data:      []byte{2, 3, 4},
	}
	tx := types.NewTx(&dtx)
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	var err error
	etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address) txmgr.Tx {
	timeNow := time.Now()
	etx := cltest.NewEthTx(fromAddress)
	ctx := tests.Context(t)

	etx.BroadcastAt = &timeNow
	etx.InitialBroadcastAt = &timeNow
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxUnconfirmed
	require.NoError(t, txStore.InsertTx(ctx, &etx))
	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

	tx := cltest.NewLegacyTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgrtypes.TxAttemptInsufficientFunds
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	var err error
	etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
	t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, broadcastBeforeBlockNum int64,
	broadcastAt time.Time, fromAddress common.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	ctx := tests.Context(t)

	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxConfirmedMissingReceipt
	require.NoError(t, txStore.InsertTx(ctx, &etx))
	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx.TxAttempts = append(etx.TxAttempts, attempt)
	return etx
}

func mustInsertInProgressEthTxWithAttempt(t testing.TB, txStore txmgr.TestEvmTxStore, nonce evmtypes.Nonce, fromAddress common.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	ctx := tests.Context(t)

	etx.Sequence = &nonce
	etx.State = txmgrcommon.TxInProgress
	require.NoError(t, txStore.InsertTx(ctx, &etx))
	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
	tx := cltest.NewLegacyTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()
	attempt.State = txmgrtypes.TxAttemptInProgress
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	var err error
	etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustCreateUnstartedGeneratedTx(t testing.TB, txStore txmgr.EvmTxStore, fromAddress common.Address, chainID *big.Int, opts ...func(*txmgr.TxRequest)) (tx txmgr.Tx) {
	txRequest := txmgr.TxRequest{
		FromAddress: fromAddress,
	}

	// Apply the default options
	withDefaults()(&txRequest)
	// Apply the optional parameters
	for _, opt := range opts {
		opt(&txRequest)
	}
	return mustCreateUnstartedTxFromEvmTxRequest(t, txStore, txRequest, chainID)
}

func withDefaults() func(*txmgr.TxRequest) {
	return func(tx *txmgr.TxRequest) {
		tx.ToAddress = testutils.NewAddress()
		tx.EncodedPayload = []byte{1, 2, 3}
		tx.Value = big.Int(assets.NewEthValue(142))
		tx.FeeLimit = uint64(1000000000)
		tx.Strategy = txmgrcommon.NewSendEveryStrategy()
		// Set default values for other fields if needed
	}
}

func mustCreateUnstartedTx(t testing.TB, txStore txmgr.EvmTxStore, fromAddress common.Address, toAddress common.Address, encodedPayload []byte, gasLimit uint64, value big.Int, chainID *big.Int) (tx txmgr.Tx) {
	txRequest := txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          value,
		FeeLimit:       gasLimit,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	}

	return mustCreateUnstartedTxFromEvmTxRequest(t, txStore, txRequest, chainID)
}

func mustCreateUnstartedTxFromEvmTxRequest(t testing.TB, txStore txmgr.EvmTxStore, txRequest txmgr.TxRequest, chainID *big.Int) (tx txmgr.Tx) {
	tx, err := txStore.CreateTransaction(tests.Context(t), txRequest, chainID)
	require.NoError(t, err)

	_, err = txRequest.Strategy.PruneQueue(tests.Context(t), txStore)
	require.NoError(t, err)

	return tx
}

func mustInsertUnstartedTx(t testing.TB, txStore txmgr.TestEvmTxStore, fromAddress common.Address) {
	etx := cltest.NewEthTx(fromAddress)
	ctx := tests.Context(t)
	require.NoError(t, txStore.InsertTx(ctx, &etx))
}

func txRequestWithStrategy(strategy txmgrtypes.TxStrategy) func(*txmgr.TxRequest) {
	return func(tx *txmgr.TxRequest) {
		tx.Strategy = strategy
	}
}

func txRequestWithChecker(checker txmgr.TransmitCheckerSpec) func(*txmgr.TxRequest) {
	return func(tx *txmgr.TxRequest) {
		tx.Checker = checker
	}
}
func txRequestWithValue(value big.Int) func(*txmgr.TxRequest) {
	return func(tx *txmgr.TxRequest) {
		tx.Value = value
	}
}

func txRequestWithIdempotencyKey(idempotencyKey string) func(*txmgr.TxRequest) {
	return func(tx *txmgr.TxRequest) {
		tx.IdempotencyKey = &idempotencyKey
	}
}
