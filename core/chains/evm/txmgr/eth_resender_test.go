package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_EthResender_resendUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	lggr := logger.TestLogger(t)
	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})
	evmcfg := txmgr.NewEvmTxmConfig(evmtest.NewChainScopedConfig(t, cfg))

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress3 := cltest.MustInsertRandomKey(t, ethKeyStore)

	txStore := cltest.NewTxStore(t, db, logCfg)

	originalBroadcastAt := time.Unix(1616509100, 0)

	var addr1TxesRawHex, addr2TxesRawHex, addr3TxesRawHex []string
	// fewer than EvmMaxInFlightTransactions
	for i := uint32(0); i < evmcfg.EvmMaxInFlightTransactions()/2; i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress, originalBroadcastAt)
		addr1TxesRawHex = append(addr1TxesRawHex, hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx))
	}

	// exactly EvmMaxInFlightTransactions
	for i := uint32(0); i < evmcfg.EvmMaxInFlightTransactions(); i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress2, originalBroadcastAt)
		addr2TxesRawHex = append(addr2TxesRawHex, hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx))
	}

	// more than EvmMaxInFlightTransactions
	for i := uint32(0); i < evmcfg.EvmMaxInFlightTransactions()*2; i++ {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(i), fromAddress3, originalBroadcastAt)
		addr3TxesRawHex = append(addr3TxesRawHex, hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx))
	}

	er := txmgr.NewEthResender(lggr, txStore, ethClient, ethKeyStore, 100*time.Millisecond, evmcfg)

	t.Run("sends up to EvmMaxInFlightTransactions per key", func(t *testing.T) {
		ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(elems []rpc.BatchElem) bool {
			resentHex := make([]string, len(elems))
			for i, elem := range elems {
				resentHex[i] = elem.Args[0].(string)
			}
			assert.Len(t, elems, len(addr1TxesRawHex)+len(addr2TxesRawHex)+int(evmcfg.EvmMaxInFlightTransactions()))
			// All addr1TxesRawHex should be included
			for _, addr := range addr1TxesRawHex {
				assert.Contains(t, resentHex, addr)
			}
			// All addr2TxesRawHex should be included
			for _, addr := range addr1TxesRawHex {
				assert.Contains(t, resentHex, addr)
			}
			// Up to limit EvmMaxInFlightTransactions addr3TxesRawHex should be included
			for i, addr := range addr1TxesRawHex {
				if i > int(evmcfg.EvmMaxInFlightTransactions()) {
					// Above limit EvmMaxInFlightTransactions addr3TxesRawHex should NOT be included
					assert.NotContains(t, resentHex, addr)
				} else {
					assert.Contains(t, resentHex, addr)
				}
			}
			return true
		})).Run(func(args mock.Arguments) {}).Return(nil)
		err := er.ResendUnconfirmed()
		require.NoError(t, err)

	})
}

func Test_EthResender_alertUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	lggr, o := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	// Set this to the smallest non-zero value possible for the attempt to be eligible for resend
	delay := models.MustNewDuration(1 * time.Nanosecond)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0] = &v2.EVMConfig{
			Chain: v2.Defaults(utils.NewBig(big.NewInt(0)), &v2.Chain{
				Transactions: v2.Transactions{ResendAfterThreshold: delay},
			}),
		}
	})
	evmcfg := txmgr.NewEvmTxmConfig(evmtest.NewChainScopedConfig(t, cfg))

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	txStore := cltest.NewTxStore(t, db, logCfg)

	originalBroadcastAt := time.Unix(1616509100, 0)
	er := txmgr.NewEthResender(lggr, txStore, ethClient, ethKeyStore, 100*time.Millisecond, evmcfg)

	t.Run("alerts only once for unconfirmed transaction attempt within the unconfirmedTxAlertDelay duration", func(t *testing.T) {
		_ = cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(1), fromAddress, originalBroadcastAt)

		ethClient.On("BatchCallContextAll", mock.Anything, mock.Anything).Return(nil)

		// Try to resend the same unconfirmed attempt twice within the unconfirmedTxAlertDelay to only receive one alert
		err1 := er.ResendUnconfirmed()
		require.NoError(t, err1)

		err2 := er.ResendUnconfirmed()
		require.NoError(t, err2)
		testutils.WaitForLogMessageCount(t, o, "TxAttempt has been unconfirmed for more than: ", 1)
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
	txStore := cltest.NewTxStore(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	evmcfg := txmgr.NewEvmTxmConfig(evmtest.NewChainScopedConfig(t, cfg))
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	lggr := logger.TestLogger(t)

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		er := txmgr.NewEthResender(lggr, txStore, ethClient, ethKeyStore, 100*time.Millisecond, evmcfg)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress, originalBroadcastAt)
		etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress, originalBroadcastAt)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress, time.Now().Add(1*time.Hour))

		// First batch of 1
		ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx)
		})).Return(nil)
		// Second batch of 1
		ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx2.EthTxAttempts[0].SignedRawTx)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// It should update BroadcastAt even if there is an error here
			elems[0].Error = errors.New("kaboom")
		})

		func() {
			er.Start()
			defer er.Stop()

			cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, time.Second)
		}()

		var dbEtx txmgr.DbEthTx
		err := db.Get(&dbEtx, `SELECT * FROM eth_txes WHERE id = $1`, etx.ID)
		require.NoError(t, err)
		var dbEtx2 txmgr.DbEthTx
		err = db.Get(&dbEtx2, `SELECT * FROM eth_txes WHERE id = $1`, etx2.ID)
		require.NoError(t, err)

		assert.Greater(t, dbEtx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		assert.Greater(t, dbEtx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}
