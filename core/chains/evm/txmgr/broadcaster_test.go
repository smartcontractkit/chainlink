package txmgr_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// NewEthBroadcaster creates a new txmgr.EthBroadcaster for use in testing.
func NewTestEthBroadcaster(
	t testing.TB,
	txStore txmgr.TestEvmTxStore,
	ethClient client.Client,
	keyStore keystore.Eth,
	gconfig chainlink.GeneralConfig,
	config evmconfig.ChainScopedConfig,
	checkerFactory txmgr.TransmitCheckerFactory,
	nonceAutoSync bool,
	nonceTracker txmgr.NonceTracker,
) *txmgr.Broadcaster {
	t.Helper()

	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()

	estimator := gas.NewEvmFeeEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(config.EVM().GasEstimator(), nil, ge.BlockHistory(), lggr, nil)
	}, ge.EIP1559DynamicFees(), ge)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, keyStore, estimator)
	ethBroadcaster := txmgrcommon.NewBroadcaster(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(config.EVM().GasEstimator()), config.EVM().Transactions(), gconfig.Database().Listener(), keyStore, txBuilder, nonceTracker, lggr, checkerFactory, nonceAutoSync)

	// Mark instance as test
	ethBroadcaster.XXXTestDisableUnstartedTxAutoProcessing()
	servicetest.Run(t, ethBroadcaster)
	return ethBroadcaster
}

func TestEthBroadcaster_Lifecycle(t *testing.T) {
	cfg, db := heavyweight.FullTestDBV2(t, nil)
	txStore := cltest.NewTestTxStore(t, db)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), evmcfg.EVM().GasEstimator(), ethKeyStore, estimator)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil)
	eb := txmgr.NewEvmBroadcaster(
		txStore,
		txmClient,
		txmgr.NewEvmTxmConfig(evmcfg.EVM()),
		txmgr.NewEvmTxmFeeConfig(evmcfg.EVM().GasEstimator()),
		evmcfg.EVM().Transactions(),
		cfg.Database().Listener(),
		ethKeyStore,
		txBuilder,
		logger.Test(t),
		&testCheckerFactory{},
		false,
	)

	// Can't close an unstarted instance
	err := eb.Close()
	require.Error(t, err)
	ctx := testutils.Context(t)

	// Can start a new instance
	err = eb.Start(ctx)
	require.NoError(t, err)

	// Can successfully close once
	err = eb.Close()
	require.NoError(t, err)

	// Can't start more than once (Broadcaster uses services.StateMachine)
	err = eb.Start(ctx)
	require.Error(t, err)
	// Can't close more than once (Broadcaster uses services.StateMachine)
	err = eb.Close()
	require.Error(t, err)

	// Can't closeInternal unstarted instance
	require.Error(t, eb.XXXTestCloseInternal())

	// Can successfully startInternal a previously closed instance
	require.NoError(t, eb.XXXTestStartInternal(ctx))
	// Can't startInternal already started instance
	require.Error(t, eb.XXXTestStartInternal(ctx))
	// Can successfully closeInternal again
	require.NoError(t, eb.XXXTestCloseInternal())
}

// Failure to load next sequnce map should not fail Broadcaster startup
func TestEthBroadcaster_LoadNextSequenceMapFailure_StartupSuccess(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), evmcfg.EVM().GasEstimator(), ethKeyStore, estimator)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), errors.New("Getting on-chain nonce failed"))
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	eb := txmgr.NewEvmBroadcaster(
		txStore,
		txmClient,
		txmgr.NewEvmTxmConfig(evmcfg.EVM()),
		txmgr.NewEvmTxmFeeConfig(evmcfg.EVM().GasEstimator()),
		evmcfg.EVM().Transactions(),
		cfg.Database().Listener(),
		ethKeyStore,
		txBuilder,
		logger.Test(t),
		&testCheckerFactory{},
		false,
	)

	// Instance starts without error even if loading next sequence map fails
	err := eb.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, eb.Close()) })
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ctx := testutils.Context(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	checkerFactory := &txmgr.CheckerFactory{Client: ethClient}

	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
	ethClient.On("PendingNonceAt", mock.Anything, otherAddress).Return(uint64(0), nil).Once()
	lggr := logger.Test(t)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, checkerFactory, false, nonceTracker)
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	timeNow := time.Now()

	encodedPayload := []byte{1, 2, 3}
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint64(242)
	checker := txmgr.TransmitCheckerSpec{
		CheckerType: txmgr.TransmitCheckerTypeSimulate,
	}

	t.Run("no eth_txes at all", func(t *testing.T) {
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)
	})

	t.Run("eth_txes exist for a different from address", func(t *testing.T) {
		mustCreateUnstartedTx(t, txStore, otherAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)
	})

	t.Run("existing eth_txes with broadcast_at or error", func(t *testing.T) {
		nonce := evmtypes.Nonce(342)
		errStr := "some error"

		etxUnconfirmed := txmgr.Tx{
			Sequence:           &nonce,
			FromAddress:        fromAddress,
			ToAddress:          toAddress,
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              txmgrcommon.TxUnconfirmed,
		}
		etxWithError := txmgr.Tx{
			Sequence:       nil,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			Error:          null.StringFrom(errStr),
			State:          txmgrcommon.TxFatalError,
		}

		require.NoError(t, txStore.InsertTx(ctx, &etxUnconfirmed))
		require.NoError(t, txStore.InsertTx(ctx, &etxWithError))

		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)
	})

	t.Run("sends 3 EthTxs in order with higher value last, and lower values starting from the earliest", func(t *testing.T) {
		// Higher value
		expensiveEthTx := txmgr.Tx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          big.Int(assets.NewEthValue(242)),
			FeeLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          txmgrcommon.TxUnstarted,
		}
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(2) && tx.Value().Cmp(big.NewInt(242)) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Earlier
		tr := int32(99)
		b, err := json.Marshal(txmgr.TxMeta{JobID: &tr})
		require.NoError(t, err)
		meta := sqlutil.JSON(b)
		earlierEthTx := txmgr.Tx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			FeeLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          txmgrcommon.TxUnstarted,
			Meta:           &meta,
		}
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(0) {
				return false
			}
			require.Equal(t, evmcfg.EVM().ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, evmcfg.EVM().GasEstimator().PriceDefault().ToInt(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.String(), tx.Value().String())
			require.Equal(t, earlierEthTx.EncodedPayload, tx.Data())
			return true
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Later
		laterEthTx := txmgr.Tx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 1},
			Value:          value,
			FeeLimit:       gasLimit,
			CreatedAt:      time.Unix(1, 0),
			State:          txmgrcommon.TxUnstarted,
		}
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(1) {
				return false
			}
			require.Equal(t, evmcfg.EVM().ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, evmcfg.EVM().GasEstimator().PriceDefault().ToInt(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.String(), tx.Value().String())
			require.Equal(t, laterEthTx.EncodedPayload, tx.Data())
			return true
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Insertion order deliberately reversed to test ordering
		require.NoError(t, txStore.InsertTx(ctx, &expensiveEthTx))
		require.NoError(t, txStore.InsertTx(ctx, &laterEthTx))
		require.NoError(t, txStore.InsertTx(ctx, &earlierEthTx))

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check earlierEthTx and it's attempt
		// This was the earlier one sent so it has the lower nonce
		earlierTransaction, err := txStore.FindTxWithAttempts(ctx, earlierEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, earlierTransaction.FromAddress)
		assert.Equal(t, fromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Sequence)
		assert.Equal(t, evmtypes.Nonce(0), *earlierTransaction.Sequence)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.NotNil(t, earlierTransaction.InitialBroadcastAt)
		assert.Len(t, earlierTransaction.TxAttempts, 1)
		var m txmgr.TxMeta
		err = json.Unmarshal(*earlierEthTx.Meta, &m)
		require.NoError(t, err)
		assert.NotNil(t, m.JobID)
		assert.Equal(t, tr, *m.JobID)

		attempt := earlierTransaction.TxAttempts[0]

		assert.Equal(t, earlierTransaction.ID, attempt.TxID)
		assert.NotNil(t, attempt.TxFee.Legacy)
		assert.Nil(t, attempt.TxFee.DynamicTipCap)
		assert.Nil(t, attempt.TxFee.DynamicFeeCap)
		assert.Equal(t, evmcfg.EVM().GasEstimator().PriceDefault(), attempt.TxFee.Legacy)

		_, err = txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.Receipts, 0)

		// Check laterEthTx and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := txStore.FindTxWithAttempts(ctx, laterEthTx.ID)
		require.NoError(t, err)
		assert.False(t, earlierTransaction.Error.Valid)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, fromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Sequence)
		assert.Equal(t, evmtypes.Nonce(1), *laterTransaction.Sequence)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.NotNil(t, earlierTransaction.InitialBroadcastAt)
		assert.Len(t, laterTransaction.TxAttempts, 1)

		attempt = laterTransaction.TxAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.TxID)
		assert.Equal(t, evmcfg.EVM().GasEstimator().PriceDefault(), attempt.TxFee.Legacy)

		_, err = txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.Receipts, 0)
	})

	rnd := int64(1000000000 + rand.Intn(5000))
	cfg = configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
		c.EVM[0].GasEstimator.TipCapDefault = assets.NewWeiI(rnd)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.NewWeiI(rnd + 1)
		c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(rnd + 2)
	})
	evmcfg = evmtest.NewChainScopedConfig(t, cfg)
	ethClient.On("PendingNonceAt", mock.Anything, otherAddress).Return(uint64(1), nil).Once()
	nonceTracker = txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb = NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, checkerFactory, false, nonceTracker)

	t.Run("sends transactions with type 0x2 in EIP-1559 mode", func(t *testing.T) {
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(343) && tx.Value().Cmp(big.NewInt(242)) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, []byte{42, 42, 0}, gasLimit, big.Int(assets.NewEthValue(242)), &cltest.FixtureChainID)
		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check eipTxWithAl and it's attempt
		// This was the earlier one sent so it has the lower nonce
		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.False(t, etx.Error.Valid)
		require.NotNil(t, etx.FromAddress)
		assert.Equal(t, fromAddress, etx.FromAddress)
		require.NotNil(t, etx.Sequence)
		assert.Equal(t, evmtypes.Nonce(343), *etx.Sequence)
		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.Len(t, etx.TxAttempts, 1)

		attempt := etx.TxAttempts[0]

		assert.Equal(t, etx.ID, attempt.TxID)
		assert.Nil(t, attempt.TxFee.Legacy)
		assert.Equal(t, rnd, attempt.TxFee.DynamicTipCap.ToInt().Int64())
		assert.Equal(t, rnd+1, attempt.TxFee.DynamicFeeCap.ToInt().Int64())

		_, err = txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.Receipts, 0)
	})

	t.Run("transaction simulation", func(t *testing.T) {
		t.Run("when simulation succeeds, sends tx as normal", func(t *testing.T) {
			txRequest := txmgr.TxRequest{
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				EncodedPayload: []byte{42, 0, 0},
				Value:          big.Int(assets.NewEthValue(442)),
				FeeLimit:       gasLimit,
				Strategy:       txmgrcommon.NewSendEveryStrategy(),
				Checker: txmgr.TransmitCheckerSpec{
					CheckerType: txmgr.TransmitCheckerTypeSimulate,
				},
			}
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(344) && tx.Value().Cmp(big.NewInt(442)) == 0
			}), fromAddress).Return(commonclient.Successful, nil).Once()
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				if fmt.Sprintf("%s", callarg["value"]) == "0x1ba" { // 442
					assert.Equal(t, txRequest.FromAddress, callarg["from"])
					assert.Equal(t, &txRequest.ToAddress, callarg["to"])
					assert.Equal(t, hexutil.Uint64(txRequest.FeeLimit), callarg["gas"])
					assert.Nil(t, callarg["gasPrice"])
					assert.Nil(t, callarg["maxFeePerGas"])
					assert.Nil(t, callarg["maxPriorityFeePerGas"])
					assert.Equal(t, (*hexutil.Big)(&txRequest.Value), callarg["value"])
					assert.Equal(t, hexutil.Bytes(txRequest.EncodedPayload), callarg["data"])
					return true
				}
				return false
			}), "latest").Return(nil).Once()

			ethTx := mustCreateUnstartedTxFromEvmTxRequest(t, txStore, txRequest, &cltest.FixtureChainID)

			{
				retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
				assert.NoError(t, err)
				assert.False(t, retryable)
			}

			// Check ethtx was sent
			ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgrcommon.TxUnconfirmed, ethTx.State)
		})

		t.Run("with unknown error, sends tx as normal", func(t *testing.T) {
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(345) && tx.Value().Cmp(big.NewInt(542)) == 0
			}), fromAddress).Return(commonclient.Successful, nil).Once()
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x21e" // 542
			}), "latest").Return(errors.New("this is not a revert, something unexpected went wrong")).Once()

			ethTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID,
				txRequestWithChecker(checker),
				txRequestWithValue(big.Int(assets.NewEthValue(542))))

			{
				retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
				assert.NoError(t, err)
				assert.False(t, retryable)
			}

			ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgrcommon.TxUnconfirmed, ethTx.State)
		})

		t.Run("on revert, marks tx as fatally errored and does not send", func(t *testing.T) {
			jerr := client.JsonError{
				Code:    42,
				Message: "oh no, it reverted",
				Data:    []byte{42, 166, 34},
			}
			ethClient.On("CallContext", mock.Anything, mock.AnythingOfType("*hexutil.Bytes"), "eth_call", mock.MatchedBy(func(callarg map[string]interface{}) bool {
				return fmt.Sprintf("%s", callarg["value"]) == "0x282" // 642
			}), "latest").Return(&jerr).Once()

			ethTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID,
				txRequestWithChecker(checker),
				txRequestWithValue(big.Int(assets.NewEthValue(642))))
			{
				retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
				assert.NoError(t, err)
				assert.False(t, retryable)
			}

			ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
			require.NoError(t, err)
			assert.Equal(t, txmgrcommon.TxFatalError, ethTx.State)
			assert.True(t, ethTx.Error.Valid)
			assert.Equal(t, "transaction reverted during simulation: json-rpc error { Code = 42, Message = 'oh no, it reverted', Data = 'KqYi' }", ethTx.Error.String)
		})
	})
}

func TestEthBroadcaster_TransmitChecking(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ctx := testutils.Context(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	checkerFactory := &testCheckerFactory{}
	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, checkerFactory, false, nonceTracker)

	checker := txmgr.TransmitCheckerSpec{
		CheckerType: txmgr.TransmitCheckerTypeSimulate,
	}
	t.Run("when transmit checking times out, sends tx as normal", func(t *testing.T) {
		// Checker will return a canceled error
		checkerFactory.err = context.Canceled

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == 0 && tx.Value().Cmp(big.NewInt(442)) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		ethTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID,
			txRequestWithValue(big.Int(assets.NewEthValue(442))),
			txRequestWithChecker(checker))
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check ethtx was sent
		ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, ethTx.State)
	})

	t.Run("when transmit checking succeeds, sends tx as normal", func(t *testing.T) {
		// Checker will return no error
		checkerFactory.err = nil

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == 1 && tx.Value().Cmp(big.NewInt(442)) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		ethTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID,
			txRequestWithValue(big.Int(assets.NewEthValue(442))),
			txRequestWithChecker(checker))
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check ethtx was sent
		ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, ethTx.State)
	})

	t.Run("when transmit errors, fatally error transaction", func(t *testing.T) {
		// Checker will return a fatal error
		checkerFactory.err = errors.New("fatal checker error")

		ethTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID, txRequestWithChecker(checker))
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check ethtx was sent
		ethTx, err := txStore.FindTxWithAttempts(ctx, ethTx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxFatalError, ethTx.State)
		assert.True(t, ethTx.Error.Valid)
		assert.Equal(t, "fatal checker error", ethTx.Error.String)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_OptimisticLockingOnEthTx(t *testing.T) {
	// non-transactional DB needed because we deliberately test for FK violation
	cfg, db := heavyweight.FullTestDBV2(t, nil)
	txStore := cltest.NewTestTxStore(t, db)
	ccfg := evmtest.NewChainScopedConfig(t, cfg)
	evmcfg := txmgr.NewEvmTxmConfig(ccfg.EVM())
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ccfg.EVM().GasEstimator(), ethKeyStore, estimator)

	chStartEstimate := make(chan struct{})
	chBlock := make(chan struct{})

	estimator.On("GetFee", mock.Anything, mock.Anything, mock.Anything, ccfg.EVM().GasEstimator().PriceMaxKey(fromAddress)).Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint64(500), nil).Run(func(_ mock.Arguments) {
		close(chStartEstimate)
		<-chBlock
	}).Once()
	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	eb := txmgr.NewEvmBroadcaster(
		txStore,
		txmClient,
		evmcfg,
		txmgr.NewEvmTxmFeeConfig(ccfg.EVM().GasEstimator()),
		ccfg.EVM().Transactions(),
		cfg.Database().Listener(),
		ethKeyStore,
		txBuilder,
		logger.Test(t),
		&testCheckerFactory{},
		false,
	)
	eb.XXXTestDisableUnstartedTxAutoProcessing()

	// Start instance of broadcaster
	servicetest.Run(t, eb)

	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)

	go func() {
		select {
		case <-chStartEstimate:
		case <-time.After(5 * time.Second):
			t.Log("timed out waiting for estimator to be called")
			return
		}

		// Simulate a "PruneQueue" call
		assert.NoError(t, commonutils.JustError(db.Exec(`DELETE FROM evm.txes WHERE state = 'unstarted'`)))
		close(chBlock)
	}()

	{
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)
	}
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success_WithMultiplier(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// Configured gas price changed
		lm := decimal.RequireFromString("1.3")
		c.EVM[0].GasEstimator.LimitMultiplier = &lm
	})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

	ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		assert.Equal(t, int(1600), int(tx.Gas()))
		return true
	}), fromAddress).Return(commonclient.Successful, nil).Once()

	txRequest := txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411"),
		EncodedPayload: []byte{42, 42, 0},
		Value:          big.Int(assets.NewEthValue(242)),
		FeeLimit:       1231,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
	}
	mustCreateUnstartedTxFromEvmTxRequest(t, txStore, txRequest, &cltest.FixtureChainID)

	// Do the thing
	{
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)
	}
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_ResumingFromCrash(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	nextNonce := evmtypes.Nonce(916714082576372851)
	firstNonce := nextNonce
	secondNonce := nextNonce + 1
	cfg := configtest.NewGeneralConfig(t, nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ctx := testutils.Context(t)

	t.Run("cannot be more than one transaction per address in an unfinished state", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		firstInProgress := txmgr.Tx{
			FromAddress:    fromAddress,
			Sequence:       &firstNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			Error:          null.String{},
			State:          txmgrcommon.TxInProgress,
		}

		secondInProgress := txmgr.Tx{
			FromAddress:    fromAddress,
			Sequence:       &secondNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			FeeLimit:       gasLimit,
			Error:          null.String{},
			State:          txmgrcommon.TxInProgress,
		}

		require.NoError(t, txStore.InsertTx(ctx, &firstInProgress))
		err := txStore.InsertTx(ctx, &secondInProgress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id\" (SQLSTATE 23505)")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so evm.key_states.next_nonce has not been
		// incremented yet
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
	})

	t.Run("previous run assigned nonce and broadcast but it fatally errored before we could save", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved the nonce to the eth_tx
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		}), fromAddress).Return(commonclient.Fatal, errors.New("exceeds block gas limit")).Once()

		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.True(t, etx.Error.Valid)
		assert.Equal(t, "exceeds block gas limit", etx.Error.String)
		assert.Len(t, etx.TxAttempts, 0)
	})

	t.Run("previous run assigned nonce and broadcast and is now in mempool", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved the nonce to the eth_tx
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		}), fromAddress).Return(commonclient.Successful, errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("previous run assigned nonce and broadcast and now the transaction has been confirmed", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved the nonce to the eth_tx
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		}), fromAddress).Return(commonclient.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		require.NotNil(t, etx.BroadcastAt)
		assert.Equal(t, *etx.BroadcastAt, etx.CreatedAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("previous run assigned nonce and then failed to reach node for some reason and node is still down", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved the nonce to the eth_tx
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		}), fromAddress).Return(commonclient.Retryable, failedToReachNodeError).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), failedToReachNodeError.Error())
		assert.True(t, retryable)

		// Check it was left in the unfinished state
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.Equal(t, nextNonce, *etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("previous run assigned nonce and broadcast transaction then crashed and rebooted with a different configured gas price", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)

		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.RandomKey{Nonce: nextNonce.Int64()}.MustInsertWithState(t, ethKeyStore)

		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			// Configured gas price changed
			c.EVM[0].GasEstimator.PriceDefault = assets.NewWeiI(500000000000)
		})
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

		// Crashed right after we commit the database transaction that saved the nonce to the eth_tx
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, txStore, firstNonce, fromAddress)
		require.Len(t, inProgressEthTx.TxAttempts, 1)
		attempt := inProgressEthTx.TxAttempts[0]

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			// Ensure that the gas price is the same as the original attempt
			s, e := txmgr.GetGethSignedTx(attempt.SignedRawTx)
			require.NoError(t, e)
			return tx.Nonce() == uint64(firstNonce) && tx.GasPrice().Int64() == s.GasPrice().Int64()
		}), fromAddress).Return(commonclient.Successful, errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		{
			retryable, err := eb.ProcessUnstartedTxs(testutils.Context(t), fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)
		}

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		s, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)
		assert.Equal(t, int64(342), s.GasPrice().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})
}

func getLocalNextNonce(t *testing.T, nonceTracker txmgr.NonceTracker, fromAddress gethCommon.Address) uint64 {
	n, err := nonceTracker.GetNextSequence(testutils.Context(t), fromAddress)
	require.NoError(t, err)
	require.NotNil(t, n)
	return uint64(n)
}

// Note that all of these tests share the same database, and ordering matters.
// This in order to more deeply test ProcessUnstartedEthTxs over
// multiple runs with previous errors in the database.
func TestEthBroadcaster_ProcessUnstartedEthTxs_Errors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
	lggr := logger.Test(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmClient)
	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)
	ctx := testutils.Context(t)

	require.NoError(t, commonutils.JustError(db.Exec(`SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)))
	require.NoError(t, commonutils.JustError(db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)))

	t.Run("if external wallet sent a transaction from the account and now the nonce is one higher than it should be and we got replacement underpriced then we assume a previous transaction of ours was the one that succeeded, and hand off to EthConfirmer", func(t *testing.T) {
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		// First send, replacement underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(0)
		}), fromAddress).Return(commonclient.Successful, errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check that the transaction was saved correctly with its attempt
		// We assume success and hand off to eth confirmer to eventually mark it as failed
		var latestID int64
		var etx1 txmgr.Tx
		require.NoError(t, db.Get(&latestID, "SELECT max(id) FROM evm.txes"))
		etx1, err = txStore.FindTxWithAttempts(ctx, latestID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		assert.NotEqual(t, etx1.CreatedAt, *etx1.BroadcastAt)
		assert.NotNil(t, etx1.InitialBroadcastAt)
		require.NotNil(t, etx1.Sequence)
		assert.Equal(t, evmtypes.Nonce(0), *etx1.Sequence)
		assert.False(t, etx1.Error.Valid)
		assert.Len(t, etx1.TxAttempts, 1)

		// Check that the local nonce was incremented by one
		finalNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		require.NoError(t, err)
		require.NotNil(t, finalNextNonce)
		require.Equal(t, int64(1), int64(finalNextNonce))
	})

	t.Run("geth Client returns an error in the fatal errors category", func(t *testing.T) {
		fatalErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)

		t.Run("without callback", func(t *testing.T) {
			etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == localNextNonce
			}), fromAddress).Return(commonclient.Fatal, errors.New(fatalErrorExample)).Once()

			retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
			assert.NoError(t, err)
			assert.False(t, retryable)

			// Check it was saved correctly with its attempt
			etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
			require.NoError(t, err)

			assert.Nil(t, etx.BroadcastAt)
			assert.Nil(t, etx.InitialBroadcastAt)
			require.Nil(t, etx.Sequence)
			assert.True(t, etx.Error.Valid)
			assert.Contains(t, etx.Error.String, "exceeds block gas limit")
			assert.Len(t, etx.TxAttempts, 0)

			// Check that the key had its nonce reset
			var nonce evmtypes.Nonce
			nonce, err = nonceTracker.GetNextSequence(ctx, fromAddress)
			require.NoError(t, err)
			// Saved NextNonce must be the same as before because this transaction
			// was not accepted by the eth node and never can be
			require.Equal(t, int64(localNextNonce), int64(nonce))
		})

		t.Run("with callback", func(t *testing.T) {
			run := cltest.MustInsertPipelineRun(t, db)
			tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
			etx := txmgr.Tx{
				FromAddress:       fromAddress,
				ToAddress:         toAddress,
				EncodedPayload:    encodedPayload,
				Value:             value,
				FeeLimit:          gasLimit,
				State:             txmgrcommon.TxUnstarted,
				PipelineTaskRunID: uuid.NullUUID{UUID: tr.ID, Valid: true},
				SignalCallback:    true,
			}

			t.Run("with erroring callback bails out", func(t *testing.T) {
				require.NoError(t, txStore.InsertTx(testutils.Context(t), &etx))
				fn := func(ctx context.Context, id uuid.UUID, result interface{}, err error) error {
					return errors.New("something exploded in the callback")
				}

				eb.SetResumeCallback(fn)

				ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return tx.Nonce() == localNextNonce
				}), fromAddress).Return(commonclient.Fatal, errors.New(fatalErrorExample)).Once()

				retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
				require.Error(t, err)
				require.Contains(t, err.Error(), "something exploded in the callback")
				assert.True(t, retryable)
			})

			t.Run("calls resume with error", func(t *testing.T) {
				fn := func(ctx context.Context, id uuid.UUID, result interface{}, err error) error {
					require.Equal(t, id, tr.ID)
					require.Nil(t, result)
					require.Error(t, err)
					require.Contains(t, err.Error(), "fatal error while sending transaction: exceeds block gas limit")
					return nil
				}

				eb.SetResumeCallback(fn)

				ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return tx.Nonce() == localNextNonce
				}), fromAddress).Return(commonclient.Fatal, errors.New(fatalErrorExample)).Once()

				{
					retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
					assert.NoError(t, err)
					assert.False(t, retryable)
				}

				// same as the parent test, but callback is set by ctor
				t.Run("callback set by ctor", func(t *testing.T) {
					estimator := gas.NewEvmFeeEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
						return gas.NewFixedPriceEstimator(evmcfg.EVM().GasEstimator(), nil, evmcfg.EVM().GasEstimator().BlockHistory(), lggr, nil)
					}, evmcfg.EVM().GasEstimator().EIP1559DynamicFees(), evmcfg.EVM().GasEstimator())
					txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), evmcfg.EVM().GasEstimator(), ethKeyStore, estimator)
					localNextNonce = getLocalNextNonce(t, nonceTracker, fromAddress)
					eb2 := txmgr.NewEvmBroadcaster(txStore, txmClient, txmgr.NewEvmTxmConfig(evmcfg.EVM()), txmgr.NewEvmTxmFeeConfig(evmcfg.EVM().GasEstimator()), evmcfg.EVM().Transactions(), cfg.Database().Listener(), ethKeyStore, txBuilder, lggr, &testCheckerFactory{}, false)
					retryable, err := eb2.ProcessUnstartedTxs(ctx, fromAddress)
					assert.NoError(t, err)
					assert.False(t, retryable)
				})
			})
		})
	})

	eb.SetResumeCallback(nil)

	t.Run("geth Client fails with error indicating that the transaction was too expensive", func(t *testing.T) {
		TxFeeExceedsCapError := "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce, nil).Once()
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.ExceedsMaxFee, errors.New(TxFeeExceedsCapError)).Twice()
		// In the first case, the tx was NOT accepted into the mempool. In the case
		// of multiple RPC nodes, it is possible that it can be accepted by
		// another node even if the primary one returns "exceeds the configured
		// cap"

		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.Contains(t, err.Error(), "error while sending transaction")
		assert.True(t, retryable)

		// Check it was saved with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt) // Note that InitialBroadcastAt really means "InitialDefinitelySuccessfulBroadcastAt"
		assert.Equal(t, evmtypes.Nonce(localNextNonce), *etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		// Check that the key had its nonce reset
		var nonce evmtypes.Nonce
		nonce, err = nonceTracker.GetNextSequence(ctx, fromAddress)
		require.NoError(t, err)
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.Equal(t, int64(localNextNonce), int64(nonce))

		// On the second try, the tx has been accepted into the mempool
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce+1, nil).Once()

		retryable, err = eb.ProcessUnstartedTxs(ctx, fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt) // Note that InitialBroadcastAt really means "InitialDefinitelySuccessfulBroadcastAt"
		assert.Equal(t, evmtypes.Nonce(localNextNonce), *etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("eth Client call fails with an unexpected random error, and transaction was not accepted into mempool", func(t *testing.T) {
		retryableErrorExample := "some unknown error"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Unknown, errors.New(retryableErrorExample)).Once()
		// Nonce is the same as localNextNonce, implying that this sent transaction has not been accepted
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce, nil).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), retryableErrorExample)
		assert.True(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		// Now on the second run, it is successful
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		retryable, err = eb.ProcessUnstartedTxs(ctx, fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("eth client call fails with an unexpected random error, and the nonce check also subsequently fails", func(t *testing.T) {
		retryableErrorExample := "some unknown error"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Unknown, errors.New(retryableErrorExample)).Once()
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), errors.New("pending nonce fetch failed")).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), retryableErrorExample)
		require.Contains(t, err.Error(), "pending nonce fetch failed")
		assert.True(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		// Now on the second run, it is successful
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		retryable, err = eb.ProcessUnstartedTxs(ctx, fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("eth Client call fails with an unexpected random error, and transaction was accepted into mempool", func(t *testing.T) {
		retryableErrorExample := "some strange RPC returns an unexpected thing"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Unknown, errors.New(retryableErrorExample)).Once()
		// Nonce is one higher than localNextNonce, implying that despite the error, this sent transaction has been accepted into the mempool
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce+1, nil).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved correctly with its attempt, in a broadcast state
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("eth node returns underpriced transaction", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)

		// First was underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(evmcfg.EVM().GasEstimator().PriceDefault().ToInt()) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()

		// Second with gas bump was still underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(25000000000)) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()

		// Third succeeded
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(30000000000)) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, "30 gwei", attempt.TxFee.Legacy.String())
	})

	etxUnfinished := txmgr.Tx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          value,
		FeeLimit:       gasLimit,
		State:          txmgrcommon.TxUnstarted,
	}
	require.NoError(t, txStore.InsertTx(ctx, &etxUnfinished))

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Retryable, failedToReachNodeError).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
		assert.True(t, retryable)

		// Check it was left in the unfinished state
		etx, err := txStore.FindTxWithAttempts(ctx, etxUnfinished.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		assert.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, etx.TxAttempts[0].State)
	})

	t.Run("eth node returns temporarily underpriced transaction", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to ethConfirmer to bump gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)

		// Re-use the previously unfinished transaction, no need to insert new

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Successful, errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		assert.NoError(t, err)
		assert.False(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err := txStore.FindTxWithAttempts(ctx, etxUnfinished.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, "20 gwei", attempt.TxFee.Legacy.String())
	})

	t.Run("eth node returns underpriced transaction and bumping gas doesn't increase it", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		evmcfg2 := evmtest.NewChainScopedConfig(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.BumpMin = assets.NewWeiI(0)
			c.EVM[0].GasEstimator.BumpPercent = ptr[uint16](0)
		}))
		eb2 := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg2, &testCheckerFactory{}, false, nonceTracker)
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)

		// First was underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(evmcfg2.EVM().GasEstimator().PriceDefault().ToInt()) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()

		// Do the thing
		retryable, err := eb2.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped fee price of 20 gwei is equal to original fee price of 20 gwei. ACTION REQUIRED: This is a configuration error, you must increase either FeeEstimator.BumpPercent or FeeEstimator.BumpMin")
		assert.True(t, retryable)

		// TEARDOWN: Clear out the unsent tx before the next test
		pgtest.MustExec(t, db, `DELETE FROM evm.txes WHERE nonce = $1`, localNextNonce)
	})

	t.Run("eth tx is left in progress if eth node returns insufficient eth", func(t *testing.T) {
		insufficientEthError := "insufficient funds for transfer"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.InsufficientFunds, errors.New(insufficientEthError)).Once()

		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient funds for transfer")
		assert.True(t, retryable)

		// Check it was saved correctly with its attempt
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)
		assert.Nil(t, attempt.BroadcastBeforeBlockNum)
	})

	pgtest.MustExec(t, db, `DELETE FROM evm.txes`)

	t.Run("eth tx is left in progress if nonce is too high", func(t *testing.T) {
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		nonceGapError := "NonceGap, Future nonce. Expected nonce: " + strconv.FormatUint(localNextNonce, 10)
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		}), fromAddress).Return(commonclient.Retryable, errors.New(nonceGapError)).Once()

		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		assert.Contains(t, err.Error(), nonceGapError)
		assert.True(t, retryable)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Nil(t, etx.InitialBroadcastAt)
		require.NotNil(t, etx.Sequence)
		assert.False(t, etx.Error.Valid)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)
		assert.Nil(t, attempt.BroadcastBeforeBlockNum)

		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)
	})

	t.Run("eth node returns underpriced transaction and bumping gas doesn't increase it in EIP-1559 mode", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.

		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		evmcfg2 := evmtest.NewChainScopedConfig(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.EVM[0].GasEstimator.BumpMin = assets.NewWeiI(0)
			c.EVM[0].GasEstimator.BumpPercent = ptr[uint16](0)
		}))
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce, nil).Once()
		eb2 := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg2, &testCheckerFactory{}, false, nonceTracker)
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		underpricedError := "transaction underpriced"
		localNextNonce = getLocalNextNonce(t, nonceTracker, fromAddress)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(1)) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()

		// Check gas tip cap verification
		retryable, err := eb2.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped gas tip cap of 1 wei is less than or equal to original gas tip cap of 1 wei")
		assert.True(t, retryable)

		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)
	})

	t.Run("eth node returns underpriced transaction in EIP-1559 mode, bumps until inclusion", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)

		// Check gas tip cap verification
		evmcfg2 := evmtest.NewChainScopedConfig(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.EVM[0].GasEstimator.TipCapDefault = assets.NewWeiI(0)
		}))
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce, nil).Once()
		eb2 := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg2, &testCheckerFactory{}, false, nonceTracker)

		retryable, err := eb2.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), "specified gas tip cap of 0 is below min configured gas tip of 1 wei for key")
		assert.True(t, retryable)

		gasTipCapDefault := assets.NewWeiI(42)

		evmcfg2 = evmtest.NewChainScopedConfig(t, configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.EVM[0].GasEstimator.TipCapDefault = gasTipCapDefault
		}))
		localNextNonce = getLocalNextNonce(t, nonceTracker, fromAddress)
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(localNextNonce, nil).Once()
		eb2 = NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg2, &testCheckerFactory{}, false, nonceTracker)

		// Second was underpriced but above minimum
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(gasTipCapDefault.ToInt()) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()
		// Resend at the bumped price
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(0).Add(gasTipCapDefault.ToInt(), evmcfg2.EVM().GasEstimator().BumpMin().ToInt())) == 0
		}), fromAddress).Return(commonclient.Underpriced, errors.New(underpricedError)).Once()
		// Final bump succeeds
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasTipCap().Cmp(big.NewInt(0).Add(gasTipCapDefault.ToInt(), big.NewInt(0).Mul(evmcfg2.EVM().GasEstimator().BumpMin().ToInt(), big.NewInt(2)))) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		retryable, err = eb2.ProcessUnstartedTxs(ctx, fromAddress)
		require.NoError(t, err)
		assert.False(t, retryable)

		// TEARDOWN: Clear out the unsent tx before the next test
		pgtest.MustExec(t, db, `DELETE FROM evm.txes WHERE nonce = $1`, localNextNonce)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_KeystoreErrors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	localNonce := 0

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db)

	realKeystore := cltest.NewKeyStore(t, db)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, realKeystore.Eth())

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Once()
	ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
	lggr := logger.Test(t)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb := NewTestEthBroadcaster(t, txStore, ethClient, kst, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)
	ctx := testutils.Context(t)
	_, err := nonceTracker.GetNextSequence(ctx, fromAddress)
	require.NoError(t, err)

	t.Run("tx signing fails", func(t *testing.T) {
		etx := mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, gasLimit, value, &cltest.FixtureChainID)
		tx := *gethTypes.NewTx(&gethTypes.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.AnythingOfType("*types.Transaction"),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&tx, errors.New("could not sign transaction"))

		// Do the thing
		retryable, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not sign transaction")
		assert.True(t, retryable)

		// Check that the transaction is left in unstarted state
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnstarted, etx.State)
		assert.Len(t, etx.TxAttempts, 0)

		// Check that the key did not have its nonce incremented
		var nonce evmtypes.Nonce
		nonce, err = nonceTracker.GetNextSequence(ctx, fromAddress)
		require.NoError(t, err)
		require.Equal(t, int64(localNonce), int64(nonce))
	})
}

func TestEthBroadcaster_Trigger(t *testing.T) {
	t.Parallel()

	// Simple sanity check to make sure it doesn't block
	db := pgtest.NewSqlxDB(t)

	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.Test(t)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))
	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, &testCheckerFactory{}, false, nonceTracker)

	eb.Trigger(testutils.NewAddress())
	eb.Trigger(testutils.NewAddress())
	eb.Trigger(testutils.NewAddress())
}

func TestEthBroadcaster_SyncNonce(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)

	lggr, observed := logger.TestObserved(t, zapcore.DebugLevel)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].NonceAutoSync = ptr(true)
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	evmTxmCfg := txmgr.NewEvmTxmConfig(evmcfg.EVM())
	txStore := cltest.NewTestTxStore(t, db)

	kst := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.RandomKey{Disabled: false}.MustInsertWithState(t, kst)

	estimator := gas.NewEvmFeeEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(evmcfg.EVM().GasEstimator(), nil, evmcfg.EVM().GasEstimator().BlockHistory(), lggr, nil)
	}, evmcfg.EVM().GasEstimator().EIP1559DynamicFees(), evmcfg.EVM().GasEstimator())
	checkerFactory := &testCheckerFactory{}

	ge := evmcfg.EVM().GasEstimator()

	t.Run("does nothing if nonce sync is disabled", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, kst, estimator)

		kst := ksmocks.NewEth(t)
		addresses := []gethCommon.Address{fromAddress}
		kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(0), nil).Once()
		txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
		eb := txmgr.NewEvmBroadcaster(txStore, txmClient, evmTxmCfg, txmgr.NewEvmTxmFeeConfig(ge), evmcfg.EVM().Transactions(), cfg.Database().Listener(), kst, txBuilder, lggr, checkerFactory, false)
		err := eb.Start(ctx)
		assert.NoError(t, err)

		defer func() { assert.NoError(t, eb.Close()) }()

		testutils.WaitForLogMessage(t, observed, "Skipping sequence auto-sync")
	})
}

func TestEthBroadcaster_NonceTracker_InProgressTx(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	checkerFactory := &txmgr.CheckerFactory{Client: ethClient}
	lggr := logger.Test(t)
	ctx := testutils.Context(t)

	t.Run("maintains the proper nonce if there is an in-progress tx during startup", func(t *testing.T) {
		inProgressTxNonce := uint64(0)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == inProgressTxNonce
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Tx with nonce 0 in DB will set local nonce map to value to 1
		mustInsertInProgressEthTxWithAttempt(t, txStore, evmtypes.Nonce(inProgressTxNonce), fromAddress)
		nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))
		eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, cfg, evmcfg, checkerFactory, false, nonceTracker)

		// Check the local nonce map was set to 1 higher than in-progress tx nonce
		nonce := getLocalNextNonce(t, nonceTracker, fromAddress)
		require.Equal(t, inProgressTxNonce+1, nonce)

		_, err := eb.ProcessUnstartedTxs(ctx, fromAddress)
		require.NoError(t, err)

		// Check the local nonce map maintained nonce 1 higher than in-progress tx nonce
		nonce = getLocalNextNonce(t, nonceTracker, fromAddress)
		require.Equal(t, inProgressTxNonce+1, nonce)
	})
}

type testCheckerFactory struct {
	err error
}

func (t *testCheckerFactory) BuildChecker(spec txmgr.TransmitCheckerSpec) (txmgr.TransmitChecker, error) {
	return &testChecker{t.err}, nil
}

type testChecker struct {
	err error
}

func (t *testChecker) Check(
	_ context.Context,
	_ logger.SugaredLogger,
	_ txmgr.Tx,
	_ txmgr.TxAttempt,
) error {
	return t.err
}
