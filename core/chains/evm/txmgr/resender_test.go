package txmgr_test

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/common/fee"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func ptr[T any](t T) *T { return &t }

func newTestChainScopedConfig(t *testing.T) evmconfig.ChainScopedConfig {
	cfg := configtest.NewTestGeneralConfig(t)
	return evmtest.NewChainScopedConfig(t, cfg)
}

func newEthResenderWithDefaultInterval(t testing.TB,
	txStore txmgr.EvmTxStore,
	ethClient client.Client,
	config evmconfig.ChainScopedConfig,
	ks keystore.Eth,
) *txmgr.Resender {
	return newEthResender(t, txStore, ethClient, config, ks, txmgrcommon.DefaultResenderPollInterval)
}

// newEthResender creates, starts, and eventually closes a Resender
func newEthResender(t testing.TB,
	txStore txmgr.EvmTxStore,
	ethClient client.Client,
	config evmconfig.ChainScopedConfig,
	ks keystore.Eth,
	pollInterval time.Duration,
) *txmgr.Resender {
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	estimator := gas.NewWrappedEvmEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(ge, ge.BlockHistory(), lggr)
	}, ge.EIP1559DynamicFees(), nil)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ks, estimator)
	er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ks, pollInterval,
		txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), config.Database())

	servicetest.Run(t, er)
	return er
}

func Test_EthResender_resendUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	ctx := testutils.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := newTestChainScopedConfig(t)

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress3 := cltest.MustInsertRandomKey(t, ethKeyStore)

	txStore := cltest.NewTestTxStore(t, db, logCfg)

	originalBroadcastAt := time.Unix(1616509100, 0)

	txConfig := cfg.EVM().Transactions()
	var addr1TxesRawHex, addr2TxesRawHex, addr3TxesRawHex []string
	// fewer than EvmMaxInFlightTransactions
	for i := uint32(0); i < txConfig.MaxInFlight()/2; i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress, originalBroadcastAt)
		addr1TxesRawHex = append(addr1TxesRawHex, hexutil.Encode(etx.TxAttempts[0].SignedRawTx))
	}

	// exactly EvmMaxInFlightTransactions
	for i := uint32(0); i < txConfig.MaxInFlight(); i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress2, originalBroadcastAt)
		addr2TxesRawHex = append(addr2TxesRawHex, hexutil.Encode(etx.TxAttempts[0].SignedRawTx))
	}

	// more than EvmMaxInFlightTransactions
	for i := uint32(0); i < txConfig.MaxInFlight()*2; i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress3, originalBroadcastAt)
		addr3TxesRawHex = append(addr3TxesRawHex, hexutil.Encode(etx.TxAttempts[0].SignedRawTx))
	}

	er := newEthResenderWithDefaultInterval(t, txStore, ethClient, cfg, ethKeyStore)

	var resentHex = make(map[string]struct{})
	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(elems []rpc.BatchElem) bool {
		for _, elem := range elems {
			resentHex[elem.Args[0].(string)] = struct{}{}
		}
		assert.Len(t, elems, len(addr1TxesRawHex)+len(addr2TxesRawHex)+int(txConfig.MaxInFlight()))
		// All addr1TxesRawHex should be included
		for _, addr := range addr1TxesRawHex {
			assert.Contains(t, resentHex, addr)
		}
		// All addr2TxesRawHex should be included
		for _, addr := range addr2TxesRawHex {
			assert.Contains(t, resentHex, addr)
		}
		// Up to limit EvmMaxInFlightTransactions addr3TxesRawHex should be included
		for i, addr := range addr3TxesRawHex {
			if i >= int(txConfig.MaxInFlight()) {
				// Above limit EvmMaxInFlightTransactions addr3TxesRawHex should NOT be included
				assert.NotContains(t, resentHex, addr)
			} else {
				assert.Contains(t, resentHex, addr)
			}
		}
		return true
	})).Run(func(args mock.Arguments) {}).Return(nil)

	err := er.ResendUnconfirmed(ctx)
	require.NoError(t, err)
}

func Test_EthResender_alertUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	lggr, o := logger.TestObserved(t, zapcore.DebugLevel)
	ctx := testutils.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	// Set this to the smallest non-zero value possible for the attempt to be eligible for resend
	delay := models.MustNewDuration(1 * time.Nanosecond)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0] = &toml.EVMConfig{
			Chain: toml.Defaults(ubig.New(big.NewInt(0)), &toml.Chain{
				Transactions: toml.Transactions{ResendAfterThreshold: delay},
			}),
		}
	})
	ccfg := evmtest.NewChainScopedConfig(t, cfg)

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	txStore := cltest.NewTestTxStore(t, db, logCfg)

	originalBroadcastAt := time.Unix(1616509100, 0)
	estimator := gasmocks.NewEvmEstimator(t)
	newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
	ge := ccfg.EVM().GasEstimator()
	feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ethKeyStore, 100*time.Millisecond, ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ge), ccfg.EVM().Transactions(), ccfg.Database())
	servicetest.Run(t, er)

	t.Run("alerts only once for unconfirmed transaction attempt within the unconfirmedTxAlertDelay duration", func(t *testing.T) {
		_ = cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(1), fromAddress, originalBroadcastAt)

		ethClient.On("BatchCallContextAll", mock.Anything, mock.Anything).Return(nil)

		// Try to resend the same unconfirmed attempt twice within the unconfirmedTxAlertDelay to only receive one alert
		err1 := er.ResendUnconfirmed(ctx)
		require.NoError(t, err1)

		err2 := er.ResendUnconfirmed(ctx)
		require.NoError(t, err2)
		testutils.WaitForLogMessageCount(t, o, "TxAttempt has been unconfirmed for more than max duration", 1)
	})
}

func Test_EthResender_Start(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		// This can be anything as long as it isn't zero
		c.EVM[0].Transactions.ResendAfterThreshold = models.MustNewDuration(42 * time.Hour)
		// Set batch size low to test batching
		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](1)
	})
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ccfg := evmtest.NewChainScopedConfig(t, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress, originalBroadcastAt)
		etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress, originalBroadcastAt)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress, time.Now().Add(1*time.Hour))

		// First batch of 1
		ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx.TxAttempts[0].SignedRawTx)
		})).Return(nil)
		// Second batch of 1
		ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx2.TxAttempts[0].SignedRawTx)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// It should update BroadcastAt even if there is an error here
			elems[0].Error = errors.New("kaboom")
		})

		newEthResender(t, txStore, ethClient, ccfg, ethKeyStore, 100*time.Millisecond)

		cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, time.Second)

		var dbEtx txmgr.DbEthTx
		err := db.Get(&dbEtx, `SELECT * FROM evm.txes WHERE id = $1`, etx.ID)
		require.NoError(t, err)
		var dbEtx2 txmgr.DbEthTx
		err = db.Get(&dbEtx2, `SELECT * FROM evm.txes WHERE id = $1`, etx2.ID)
		require.NoError(t, err)

		assert.Greater(t, dbEtx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		assert.Greater(t, dbEtx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}

func TestEthResender_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())

	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	config := newTestChainScopedConfig(t)
	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())
	mustInsertInProgressEthTx(t, txStore, 0, fromAddress)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	gasPriceWei := gas.EvmFee{Legacy: assets.GWei(52)}
	overrideGasLimit := uint32(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == uint64(overrideGasLimit) &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{1}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == uint64(etx1.FeeLimit) &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(1)}, gasPriceWei, fromAddress, 0))
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == uint64(overrideGasLimit)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == uint64(overrideGasLimit)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(1), (2)}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(1)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(2)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		for i := 3; i <= 5; i++ {
			nonce := i
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(nonce) &&
					tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
					tx.Gas() == uint64(overrideGasLimit) &&
					*tx.To() == fromAddress &&
					tx.Value().Cmp(big.NewInt(0)) == 0 &&
					len(tx.Data()) == 0
			}), mock.Anything).Return(commonclient.Successful, nil).Once()
		}
		nonces := []evmtypes.Nonce{(1), (2), (3), (4), (5)}

		require.NoError(t, er.ForceRebroadcast(testutils.Context(t), nonces, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && uint32(tx.Gas()) == config.EVM().GasEstimator().LimitDefault()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(0)}, gasPriceWei, fromAddress, 0))
	})
}

func TestEthResender_FindTxsRequiringRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	evmFromAddress := fromAddress
	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	mustInsertConfirmedEthTx(t, txStore, nonce, fromAddress)
	nonce++

	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	evmOtherAddress := otherAddress

	lggr := logger.Test(t)

	er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, ethKeyStore)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	mustInsertInProgressEthTx(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx1.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt1_1.ID))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt2_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt2_1.ID))

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etxWithoutAttempts := cltest.NewEthTx(fromAddress)
	{
		n := evmtypes.Nonce(nonce)
		etxWithoutAttempts.Sequence = &n
	}
	now := time.Now()
	etxWithoutAttempts.BroadcastAt = &now
	etxWithoutAttempts.InitialBroadcastAt = &now
	etxWithoutAttempts.State = txmgrcommon.TxUnconfirmed
	require.NoError(t, txStore.InsertTx(&etxWithoutAttempts))
	nonce++

	t.Run("does nothing if the transaction is from a different address than the one given", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmOtherAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	t.Run("returns nothing for different chain id", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, big.NewInt(42))
		require.NoError(t, err)

		require.Len(t, etxs, 0)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt3_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_1.ID))

	// NOTE: It should ignore qualifying eth_txes from a different address
	etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, otherAddress)
	attemptOther1 := etxOther.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attemptOther1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attemptOther1.ID))

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	t.Run("returns nothing if threshold is zero", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, 0, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 0)
	})

	t.Run("does not return more transactions for gas bumping than gasBumpThreshold", func(t *testing.T) {
		// Unconfirmed txes in DB are:
		// (unnamed) (nonce 2)
		// etx1 (nonce 3)
		// etx2 (nonce 4)
		// etxWithoutAttempts (nonce 5)
		// etx3 (nonce 6) - ready for bump
		// etx4 (nonce 7) - ready for bump
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 4, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1) // returns etxWithoutAttempts only - eligible for gas bumping because it technically doesn't have any attempts within gasBumpThreshold blocks
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)

		etxs, err = er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 5, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)

		// Zero limit disables it
		etxs, err = er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 0, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt4_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt4_1.ID))

	t.Run("ignores pending transactions for another key", func(t *testing.T) {
		// Re-use etx3 nonce for another key, it should not affect the results for this key
		etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, (*etx3.Sequence).Int64(), otherAddress)
		aOther := etxOther.TxAttempts[0]
		dbAttempt = txmgr.DbEthTxAttempt{}
		dbAttempt.FromTxAttempt(&aOther)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, aOther.ID))

		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 6, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt3_2))

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_3 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt3_3))

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		require.Equal(t, evmtypes.Nonce(5), *etxWithoutAttempts.Sequence)
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		require.Equal(t, evmtypes.Nonce(7), *etx4.Sequence)
	})

	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etxWithoutAttempts.ID)
	attempt0_1.State = txmgrtypes.TxAttemptInsufficientFunds
	require.NoError(t, txStore.InsertTxAttempt(&attempt0_1))

	// This attempt has insufficient_eth, but there is also another attempt4_1
	// which is old enough, so this will be caught by both queries and should
	// not be duplicated
	attempt4_2 := cltest.NewLegacyEthTxAttempt(t, etx4.ID)
	attempt4_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt4_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt4_2))

	etx5 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	nonce++

	// This etx has one attempt that is too new, which would exclude it from
	// the gas bumping query, but it should still be caught by the insufficient
	// eth query
	etx6 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	attempt6_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt6_2.BroadcastBeforeBlockNum = &tooNew
	attempt6_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30001)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt6_2))

	t.Run("returns unique attempts requiring resubmission due to insufficient eth, ordered by nonce asc", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 4)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		assert.Equal(t, etx5.ID, etxs[2].ID)
		assert.Equal(t, *etx5.Sequence, *(etxs[2].Sequence))
		assert.Equal(t, etx6.ID, etxs[3].ID)
		assert.Equal(t, *etx6.Sequence, *(etxs[3].Sequence))
	})

	t.Run("applies limit", func(t *testing.T) {
		etxs, err := er.FindTxsRequiringRebroadcast(testutils.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 2, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
	})
}

func TestEthResender_RebroadcastWhereNecessary(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = (*assets.Wei)(assets.GWei(500))
	})
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return(addresses, nil).Maybe()
	// Use a mock keystore for this test
	keyChangeCh := make(chan struct{})
	unsub := cltest.NewAwaiter()
	kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened).Once()
	er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, kst)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
	nonce++
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

	t.Run("re-sends previous transaction on keystore error", func(t *testing.T) {
		// simulate bumped transaction that is somehow impossible to sign
		kst.On("SignTx", fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(*etx.Sequence)
			}),
			mock.Anything).Return(nil, errors.New("signing error")).Once()

		// Do the thing
		err := er.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("does nothing and continues on fatal error", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Sequence) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence)
		}), fromAddress).Return(commonclient.Fatal, errors.New("exceeds block gas limit")).Once()

		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("does nothing and continues if bumped attempt transaction was too expensive", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Sequence) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()

		// Once for the bumped attempt which exceeds limit
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) && tx.GasPrice().Int64() == int64(20000000000)
		}), fromAddress).Return(commonclient.ExceedsMaxFee, errors.New("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// Did not create an additional attempt
		require.Len(t, etx.TxAttempts, 1)

		// broadcast_at did not change
		require.Equal(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		require.Equal(t, etx.InitialBroadcastAt.Unix(), originalBroadcastAt.Unix())
	})

	var attempt1_2 txmgr.TxAttempt

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
	})
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_2.ID))
	var attempt1_3 txmgr.TxAttempt

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[2].ID)
		require.Equal(t, attempt1_2.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_3 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_3.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_3.State)
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_3.ID))
	var attempt1_4 txmgr.TxAttempt

	t.Run("saves new attempt even for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		receipt := evmtypes.Receipt{BlockNumber: big.NewInt(40)}
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				receipt.TxHash = tx.Hash()
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx.State)

		// Got the new attempt
		attempt1_4 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_4.TxFee.Legacy.ToInt().Int64())

		require.Len(t, etx.TxAttempts, 4)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[3].ID)
		require.Equal(t, attempt1_2.ID, etx.TxAttempts[2].ID)
		require.Equal(t, attempt1_3.ID, etx.TxAttempts[1].ID)
		require.Equal(t, attempt1_4.ID, etx.TxAttempts[0].ID)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[1].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[2].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[3].State)
	})

	// Mark original tx as confirmed, so we won't pick it up anymore
	pgtest.MustExec(t, db, `UPDATE evm.txes SET state = 'confirmed'`)

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_1.ID))
	var attempt2_2 txmgr.TxAttempt

	t.Run("saves in_progress attempt on temporary error and returns error", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Sequence
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Unknown, errors.New("some network error")).Once()

		// Do the thing
		err := er.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)

		// Old attempt is untouched
		require.Len(t, etx2.TxAttempts, 2)
		require.Equal(t, attempt2_1.ID, etx2.TxAttempts[1].ID)
		attempt2_1 = etx2.TxAttempts[1]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2_1.State)
		assert.Equal(t, oldEnough, *attempt2_1.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		attempt2_2 = etx2.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		// Do it again and move the attempt into "broadcast"
		n = *etx2.Sequence
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		// Attempt marked "broadcast"
		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)

		// New in_progress attempt saved
		require.Len(t, etx2.TxAttempts, 2)
		require.Equal(t, attempt2_2.ID, etx2.TxAttempts[0].ID)
		attempt2_2 = etx2.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)
	})

	// Set BroadcastBeforeBlockNum again so the next test will pick it up
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_2.ID))

	t.Run("assumes that 'nonce too low' error means confirmed_missing_receipt", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Sequence
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		// Creates new attempt as normal if currentHead is not high enough
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx2.State)

		// One new attempt saved
		require.Len(t, etx2.TxAttempts, 3)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[0].State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[1].State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[2].State)
	})

	// Original tx is confirmed, so we won't pick it up anymore
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_price=$2 WHERE id=$3 RETURNING *`, oldEnough, assets.NewWeiI(35000000000), attempt3_1.ID))

	var attempt3_2 txmgr.TxAttempt

	t.Run("saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(42000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 2)
		require.Equal(t, attempt3_1.ID, etx3.TxAttempts[1].ID)
		attempt3_2 = etx3.TxAttempts[0]

		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_2.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_2.ID))
	var attempt3_3 txmgr.TxAttempt

	t.Run("handles case where transaction is already known somehow", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(50400000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 3)
		attempt3_3 = etx3.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_3.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_3.ID))
	var attempt3_4 txmgr.TxAttempt

	t.Run("pretends it was accepted and continues the cycle if rejected for being temporarily underpriced", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to the next cycle to continue bumping gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."

		expectedBumpedGasPrice := big.NewInt(60480000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if one of the bumped transactions would exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		// Set price such that the next bump will exceed EVM.GasEstimator.PriceMax
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.TxFee.Legacy.ToInt()
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(60500000000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		er2 := newEthResenderWithDefaultInterval(t, txStore, ethClient, newCfg, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && gasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("already known")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, er2.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if the current price is exactly EVM.GasEstimator.PriceMax", func(t *testing.T) {
		// Set price such that the current price is already at EVM.GasEstimator.PriceMax
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.TxFee.Legacy.ToInt()
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(60480000000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		er2 := newEthResenderWithDefaultInterval(t, txStore, ethClient, newCfg, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && gasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("already known")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, er2.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	// The EIP-1559 etx and attempt
	etx4 := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress)
	attempt4_1 := etx4.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(35), assets.GWei(100), attempt4_1.ID))
	var attempt4_2 txmgr.TxAttempt

	t.Run("EIP-1559: bumps using EIP-1559 rules when existing attempts are of type 0x2", func(t *testing.T) {
		ethTx := *types.NewTx(&types.DynamicFeeTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx4.Sequence {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		// This is the new, EIP-1559 attempt
		gasTipCap := assets.GWei(42)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && gasTipCap.ToInt().Cmp(tx.GasTipCap()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		// A new, bumped attempt
		require.Len(t, etx4.TxAttempts, 2)
		attempt4_2 = etx4.TxAttempts[0]
		assert.Nil(t, attempt4_2.TxFee.Legacy)
		assert.Equal(t, assets.GWei(42).String(), attempt4_2.TxFee.DynamicTipCap.String())
		assert.Equal(t, assets.GWei(120).String(), attempt4_2.TxFee.DynamicFeeCap.String())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(999), assets.GWei(1000), attempt4_2.ID))

	t.Run("EIP-1559: resubmits at the old price and does not create a new attempt if one of the bumped EIP-1559 transactions would have its tip cap exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.GWei(1000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		er2 := newEthResenderWithDefaultInterval(t, txStore, ethClient, newCfg, ethKeyStore)

		// Third attempt failed to bump, resubmits old one instead
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && attempt4_2.Hash.String() == tx.Hash().String()
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		require.NoError(t, er2.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		// No new tx attempts
		require.Len(t, etx4.TxAttempts, 2)
		assert.Equal(t, assets.GWei(999).Int64(), etx4.TxAttempts[0].TxFee.DynamicTipCap.ToInt().Int64())
		assert.Equal(t, assets.GWei(1000).Int64(), etx4.TxAttempts[0].TxFee.DynamicFeeCap.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(45), assets.GWei(100), attempt4_2.ID))

	t.Run("EIP-1559: saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		// NOTE: This test case was empirically impossible when I tried it on eth mainnet (any EIP1559 transaction with a higher tip cap is accepted even if it's only 1 wei more) but appears to be possible on Polygon/Matic, probably due to poor design that applies the 10% minimum to the overall value (base fee + tip cap)
		expectedBumpedTipCap := assets.GWei(54)
		require.Greater(t, expectedBumpedTipCap.Int64(), attempt4_2.TxFee.DynamicTipCap.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx4.Sequence || expectedBumpedTipCap.ToInt().Cmp(tx.GasTipCap()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && expectedBumpedTipCap.ToInt().Cmp(tx.GasTipCap()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("replacement transaction underpriced")).Once()

		// Do it
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		require.Len(t, etx4.TxAttempts, 3)
		require.Equal(t, attempt4_1.ID, etx4.TxAttempts[2].ID)
		require.Equal(t, attempt4_2.ID, etx4.TxAttempts[1].ID)
		attempt4_3 := etx4.TxAttempts[0]

		assert.Equal(t, expectedBumpedTipCap.Int64(), attempt4_3.TxFee.DynamicTipCap.ToInt().Int64())
	})
}
func TestEthResender_RebroadcastWhereNecessary_WithConnectivityCheck(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	db := pgtest.NewSqlxDB(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	t.Run("should retry previous attempt if connectivity check failed for legacy transactions", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(false)
			c.EVM[0].GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.EVM[0].GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})
		ccfg := evmtest.NewChainScopedConfig(t, cfg)

		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		kst := ksmocks.NewEth(t)

		estimator := gasmocks.NewEvmEstimator(t)
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		estimator.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, uint32(0), pkgerrors.Wrapf(fee.ErrConnectivity, "transaction..."))
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, kst, feeEstimator)
		addresses := []gethCommon.Address{fromAddress}
		kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return(addresses, nil).Maybe()
		// Create resender with necessary state
		er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ethKeyStore, 100*time.Millisecond, ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ge), ccfg.EVM().Transactions(), ccfg.Database())
		servicetest.Run(t, er)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(commonclient.Successful, nil).Once()

		err := er.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("should retry previous attempt if connectivity check failed for dynamic transactions", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.EVM[0].GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.EVM[0].GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})
		ccfg := evmtest.NewChainScopedConfig(t, cfg)

		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		kst := ksmocks.NewEth(t)

		estimator := gasmocks.NewEvmEstimator(t)
		estimator.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.DynamicFee{}, uint32(0), pkgerrors.Wrapf(fee.ErrConnectivity, "transaction..."))
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		// Create confirmer with necessary state
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, kst, feeEstimator)
		addresses := []gethCommon.Address{fromAddress}
		kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return(addresses, nil).Maybe()
		er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ethKeyStore, 100*time.Millisecond, ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ge), ccfg.EVM().Transactions(), ccfg.Database())
		servicetest.Run(t, er)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(commonclient.Successful, nil).Once()

		err := er.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})
}

func TestEthResender_RebroadcastWhereNecessary_TerminallyUnderpriced_ThenGoesThrough(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = (*assets.Wei)(assets.GWei(500))
	})
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Use a mock keystore for this test
	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", &cltest.FixtureChainID).Return(addresses, nil).Maybe()
	currentHead := int64(30)
	oldEnough := 5
	nonce := int64(0)

	t.Run("terminally underpriced transaction with in_progress attempt is retried with more gas", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		keyChangeCh := make(chan struct{})
		unsub := cltest.NewAwaiter()
		kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened).Once()
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, kst)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, nonce, fromAddress, txmgrtypes.TxAttemptInProgress, originalBroadcastAt)
		require.Equal(t, originalBroadcastAt, *etx.BroadcastAt)
		nonce++
		attempt := etx.TxAttempts[0]
		signedTx, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)

		// Fail the first time with terminally underpriced.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Once()
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		kst.On("SignTx", mock.Anything, mock.Anything, mock.Anything).Return(
			signedTx, nil,
		).Once()
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in legacy mode", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		keyChangeCh := make(chan struct{})
		unsub := cltest.NewAwaiter()
		kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened).Once()
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, kst)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
		nonce++
		legacyAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&legacyAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, legacyAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		signedLegacyTx := new(types.Transaction)
		kst.On("SignTx", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Type() == 0x0 && tx.Nonce() == uint64(*etx.Sequence)
		}), mock.Anything).Return(
			signedLegacyTx, nil,
		).Run(func(args mock.Arguments) {
			unsignedLegacyTx := args.Get(1).(*types.Transaction)
			// Use the real keystore to do the actual signing
			thisSignedLegacyTx, err := ethKeyStore.SignTx(fromAddress, unsignedLegacyTx, testutils.FixtureChainID)
			require.NoError(t, err)
			*signedLegacyTx = *thisSignedLegacyTx
		}).Times(4) // 3 failures 1 success
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in EIP-1559 mode", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		keyChangeCh := make(chan struct{})
		unsub := cltest.NewAwaiter()
		kst.On("SubscribeToKeyChanges").Return(keyChangeCh, unsub.ItHappened).Once()
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, kst)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress)
		nonce++
		dxFeeAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&dxFeeAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, dxFeeAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("transaction underpriced")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		signedDxFeeTx := new(types.Transaction)
		kst.On("SignTx", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Type() == 0x2 && tx.Nonce() == uint64(*etx.Sequence)
		}), mock.Anything).Return(
			signedDxFeeTx, nil,
		).Run(func(args mock.Arguments) {
			unsignedDxFeeTx := args.Get(1).(*types.Transaction)
			// Use the real keystore to do the actual signing
			thisSignedDxFeeTx, err := ethKeyStore.SignTx(fromAddress, unsignedDxFeeTx, testutils.FixtureChainID)
			require.NoError(t, err)
			*signedDxFeeTx = *thisSignedDxFeeTx
		}).Times(4) // 3 failures 1 success
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})
}

func TestEthResender_RebroadcastWhereNecessary_WhenOutOfEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	_, err := ethKeyStore.EnabledKeysForChain(testutils.FixtureChainID)
	require.NoError(t, err)

	config := newTestChainScopedConfig(t)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))
	var attempt1_2 txmgr.TxAttempt

	insufficientEthError := errors.New("insufficient funds for gas * price + value")

	t.Run("saves attempt with state 'insufficient_eth' if eth node returns this error", func(t *testing.T) {
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
		assert.Nil(t, attempt1_2.BroadcastBeforeBlockNum)
	})

	t.Run("does not bump gas when previous error was 'out of eth', instead resubmits existing transaction", func(t *testing.T) {
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// The attempt is still "out of eth"
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
	})

	t.Run("saves the attempt as broadcast after node wallet has been topped up with sufficient balance", func(t *testing.T) {
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, config, ethKeyStore)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// Attempt is now 'broadcast'
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	t.Run("resubmitting due to insufficient eth is not limited by EVM.GasEstimator.BumpTxDepth", func(t *testing.T) {
		depth := 2
		etxCount := 4

		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.BumpTxDepth = ptr(uint32(depth))
		})
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)
		er := newEthResenderWithDefaultInterval(t, txStore, ethClient, evmcfg, ethKeyStore)

		for i := 0; i < etxCount; i++ {
			n := nonce
			mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(n)
			}), fromAddress).Return(commonclient.Successful, nil).Once()

			nonce++
		}

		require.NoError(t, er.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		var dbAttempts []txmgr.DbEthTxAttempt

		require.NoError(t, db.Select(&dbAttempts, "SELECT * FROM evm.tx_attempts WHERE state = 'insufficient_eth'"))
		require.Len(t, dbAttempts, 0)
	})
}
