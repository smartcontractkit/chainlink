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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func makeTestEvmTxm(
	t *testing.T, db *sqlx.DB, ethClient evmclient.Client, estimator gas.EvmFeeEstimator, ccfg txmgr.ChainConfig, fcfg txmgr.FeeConfig, txConfig evmconfig.Transactions, dbConfig txmgr.DatabaseConfig, listenerConfig txmgr.ListenerConfig, keyStore keystore.Eth) (txmgr.TxManager, error) {
	lggr := logger.Test(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), ethClient, lggr, lpOpts)

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
		estimator)
}

func TestTxm_SendNativeToken_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1).ToInt()

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	keyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	estimator := gas.NewEstimator(logger.Test(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), keyStore)
	require.NoError(t, err)

	_, err = txm.SendNativeToken(tests.Context(t), big.NewInt(0), from, to, *value, 21000)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send native token to zero address")
}

func TestTxm_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db)

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	estimator := gas.NewEstimator(logger.Test(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst.Eth())
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

	t.Run("returns error if eth key state is missing or doesn't match chain ID", func(t *testing.T) {
		rndAddr := testutils.NewAddress()
		_, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    rndAddr,
			ToAddress:      testutils.NewAddress(),
			EncodedPayload: []byte{1, 2, 3},
			FeeLimit:       21000,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("no eth key exists with address %s", rndAddr.String()))

		_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth(), *ubig.NewI(1337))

		_, err = txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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
		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)

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
		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)
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
		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)
		pgtest.MustExec(t, db, `DELETE FROM evm.forwarders`)
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

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func TestTxm_CreateTransaction_OutOfEth(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	etKeyStore := cltest.NewKeyStore(t, db).Eth()

	thisKey, _ := cltest.RandomKey{Nonce: 1}.MustInsert(t, etKeyStore)
	otherKey, _ := cltest.RandomKey{Nonce: 1}.MustInsert(t, etKeyStore)

	fromAddress := thisKey.Address
	evmFromAddress := fromAddress
	gasLimit := uint64(1000)
	toAddress := testutils.NewAddress()

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	estimator := gas.NewEstimator(logger.Test(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), etKeyStore)
	require.NoError(t, err)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)

		evmConfig.MaxQueued = uint64(1)
		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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

	require.NoError(t, commonutils.JustError(db.Exec(`DELETE FROM evm.txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has any transactions with insufficient eth errors, inserts it anyway", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		evmConfig.MaxQueued = uint64(1)

		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    evmFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		})
		require.NoError(t, err)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, commonutils.JustError(db.Exec(`DELETE FROM evm.txes WHERE from_address = $1`, thisKey.Address)))

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 42, thisKey.Address)
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{})
		strategy.On("PruneQueue", mock.Anything, mock.Anything).Return(nil, nil)

		evmConfig.MaxQueued = uint64(1)
		etx, err := txm.CreateTransaction(tests.Context(t), txmgr.TxRequest{
			FromAddress:    evmFromAddress,
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
	db := pgtest.NewSqlxDB(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	kst := ksmocks.NewEth(t)

	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	config.SetFinalityDepth(uint32(42))
	config.RpcDefaultBatchSize = uint32(4)

	evmConfig.ResendAfterThreshold = 1 * time.Hour
	evmConfig.ReaperThreshold = 1 * time.Hour
	evmConfig.ReaperInterval = 1 * time.Hour

	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return([]common.Address{}, nil)

	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges", mock.Anything).Return(keyChangeCh, unsub.ItHappened)
	estimator := gas.NewEstimator(logger.Test(t), ethClient, config, evmConfig.GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), dbConfig, dbConfig.Listener(), kst)
	require.NoError(t, err)

	head := cltest.Head(42)
	// It should not hang or panic
	txm.OnNewLongestChain(tests.Context(t), head)

	evmConfig.BumpThreshold = uint64(1)

	require.NoError(t, txm.Start(tests.Context(t)))

	ctx, cancel := context.WithTimeout(tests.Context(t), 5*time.Second)
	t.Cleanup(cancel)
	txm.OnNewLongestChain(ctx, head)
	require.NoError(t, ctx.Err())

	keyState := cltest.MustGenerateRandomKeyState(t)

	addr := []common.Address{keyState.Address.Address()}
	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addr, nil)
	ethClient.On("PendingNonceAt", mock.AnythingOfType("*context.cancelCtx"), common.Address{}).Return(uint64(0), nil).Maybe()
	keyChangeCh <- struct{}{}

	require.NoError(t, txm.Close())
	unsub.AwaitOrFail(t, 1*time.Second)
}

func TestTxm_Reset(t *testing.T) {
	t.Parallel()

	// Lots of boilerplate setup since we actually want to test start/stop of EthBroadcaster/EthConfirmer
	db := pgtest.NewSqlxDB(t)
	gcfg := configtest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := cltest.NewKeyStore(t, db)

	_, addr := cltest.RandomKey{}.MustInsert(t, kst.Eth())
	_, addr2 := cltest.RandomKey{}.MustInsert(t, kst.Eth())
	txStore := cltest.NewTestTxStore(t, db)
	// 4 confirmed tx from addr1
	for i := int64(0); i < 4; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr)
	}
	// 2 confirmed from addr2
	for i := int64(0); i < 2; i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, i, i*42+1, addr2)
	}

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, nil)
	ethClient.On("BatchCallContextAll", mock.Anything, mock.Anything).Return(nil).Maybe()
	ethClient.On("PendingNonceAt", mock.Anything, addr).Return(uint64(128), nil).Maybe()
	ethClient.On("PendingNonceAt", mock.Anything, addr2).Return(uint64(44), nil).Maybe()

	estimator := gas.NewEstimator(logger.Test(t), ethClient, cfg.EVM(), cfg.EVM().GasEstimator())
	txm, err := makeTestEvmTxm(t, db, ethClient, estimator, cfg.EVM(), cfg.EVM().GasEstimator(), cfg.EVM().Transactions(), gcfg.Database(), gcfg.Database().Listener(), kst.Eth())
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

func newTxStore(t *testing.T, db *sqlx.DB) txmgr.EvmTxStore {
	return txmgr.NewTxStore(db, logger.Test(t))
}

func newEthReceipt(blockNumber int64, blockHash common.Hash, txHash common.Hash, status uint64) txmgr.Receipt {
	transactionIndex := uint(cltest.NewRandomPositiveInt64())

	receipt := evmtypes.Receipt{
		BlockNumber:      big.NewInt(blockNumber),
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
		Status:           status,
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

func mustInsertEthReceipt(t *testing.T, txStore txmgr.TestEvmTxStore, blockNumber int64, blockHash common.Hash, txHash common.Hash) txmgr.Receipt {
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
func mustInsertConfirmedEthTxWithReceipt(t *testing.T, txStore txmgr.TestEvmTxStore, fromAddress common.Address, nonce, blockNum int64) (etx txmgr.Tx) {
	etx = cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, nonce, blockNum, fromAddress)
	mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx.TxAttempts[0].Hash)
	return etx
}

func mustInsertConfirmedEthTxBySaveFetchedReceipts(t *testing.T, txStore txmgr.TestEvmTxStore, fromAddress common.Address, nonce int64, blockNum int64, chainID big.Int) (etx txmgr.Tx) {
	etx = cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, nonce, blockNum, fromAddress)
	receipt := evmtypes.Receipt{
		TxHash:           etx.TxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(nonce),
		TransactionIndex: uint(1),
	}
	err := txStore.SaveFetchedReceipts(tests.Context(t), []*evmtypes.Receipt{&receipt}, txmgrcommon.TxConfirmed, nil, &chainID)
	require.NoError(t, err)
	return etx
}

func mustInsertFatalErrorEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, fromAddress common.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.Error = null.StringFrom("something exploded")
	etx.State = txmgrcommon.TxFatalError

	require.NoError(t, txStore.InsertTx(tests.Context(t), &etx))
	return etx
}

func mustInsertUnconfirmedEthTxWithAttemptState(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, txAttemptState txmgrtypes.TxAttemptState, opts ...interface{}) txmgr.Tx {
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

func mustInsertInProgressEthTxWithAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce evmtypes.Nonce, fromAddress common.Address) txmgr.Tx {
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
