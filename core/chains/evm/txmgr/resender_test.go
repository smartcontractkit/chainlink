package txmgr_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
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
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func Test_EthResender_resendUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})
	ccfg := evmtest.NewChainScopedConfig(t, cfg)

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress3 := cltest.MustInsertRandomKey(t, ethKeyStore)

	txStore := cltest.NewTestTxStore(t, db, logCfg)

	originalBroadcastAt := time.Unix(1616509100, 0)

	txConfig := ccfg.EVM().Transactions()
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

	estimator := gasmocks.NewEvmEstimator(t)
	newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
	ge := ccfg.EVM().GasEstimator()
	feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ethKeyStore, 100*time.Millisecond, ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ge), ccfg.EVM().Transactions(), ccfg.Database())
	servicetest.Run(t, er)

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

	err := er.XXXTestResendUnconfirmed(ctx)
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
		err1 := er.XXXTestResendUnconfirmed(ctx)
		require.NoError(t, err1)

		err2 := er.XXXTestResendUnconfirmed(ctx)
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
	lggr := logger.Test(t)

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		estimator := gasmocks.NewEvmEstimator(t)
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
		er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ethKeyStore, 100*time.Millisecond, ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ge), ccfg.EVM().Transactions(), ccfg.Database())

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

		func() {
			servicetest.Run(t, er)

			cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, time.Second)
		}()

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
		ec := newEthResender(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == uint64(overrideGasLimit) &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{1}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ec := newEthResender(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == uint64(etx1.FeeLimit) &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(1)}, gasPriceWei, fromAddress, 0))
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ec := newEthResender(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == uint64(overrideGasLimit)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == uint64(overrideGasLimit)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(1), (2)}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ec := newEthResender(t, txStore, ethClient, config, ethKeyStore)

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

		require.NoError(t, ec.ForceRebroadcast(testutils.Context(t), nonces, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ec := newEthResender(t, txStore, ethClient, config, ethKeyStore)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && uint32(tx.Gas()) == config.EVM().GasEstimator().LimitDefault()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(testutils.Context(t), []evmtypes.Nonce{(0)}, gasPriceWei, fromAddress, 0))
	})
}

func newEthResender(t testing.TB, txStore txmgr.EvmTxStore, ethClient client.Client, config evmconfig.ChainScopedConfig, ks keystore.Eth) *txmgr.Resender {
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	estimator := gas.NewWrappedEvmEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(ge, ge.BlockHistory(), lggr)
	}, ge.EIP1559DynamicFees(), nil)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ks, estimator)
	er := txmgr.NewEvmResender(lggr, txStore, txmgr.NewEvmTxmClient(ethClient), txBuilder, ks, txmgrcommon.DefaultResenderPollInterval,
		txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), config.Database())

	servicetest.Run(t, er)
	return er
}
