package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_EthResender_FindEthTxAttemptsRequiringResend(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := txmgr.FindEthTxAttemptsRequiringResend(db, olderThan, 10, cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	etxs := []txmgr.EthTx{
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress, time.Unix(1616509100, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress, time.Unix(1616509200, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress, time.Unix(1616509300, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, borm, 3, fromAddress, time.Unix(1616509400, 0)),
	}
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etxs[0].ID)
	attempt1_2.GasPrice = utils.NewBig(big.NewInt(10))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))

	attempt3_2 := newInProgressLegacyEthTxAttempt(t, etxs[2].ID)
	attempt3_2.GasPrice = utils.NewBig(big.NewInt(10))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_2))

	attempt4_2 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_2.GasTipCap = utils.NewBig(big.NewInt(10))
	attempt4_2.GasFeeCap = utils.NewBig(big.NewInt(20))
	attempt4_2.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_2))
	attempt4_4 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_4.GasTipCap = utils.NewBig(big.NewInt(30))
	attempt4_4.GasFeeCap = utils.NewBig(big.NewInt(40))
	attempt4_4.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_4))
	attempt4_3 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_3.GasTipCap = utils.NewBig(big.NewInt(20))
	attempt4_3.GasFeeCap = utils.NewBig(big.NewInt(30))
	attempt4_3.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_3))

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txmgr.FindEthTxAttemptsRequiringResend(db, olderThan, 0, cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].EthTxAttempts[0].ID, attempts[1].ID)
	})

	t.Run("returns the highest price attempt for EIP-1559 transactions", func(t *testing.T) {
		olderThan := time.Unix(1616509400, 0)
		attempts, err := txmgr.FindEthTxAttemptsRequiringResend(db, olderThan, 0, cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, attempts, 4)
		assert.Equal(t, attempt4_4.ID, attempts[3].ID)
	})

	t.Run("applies limit", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txmgr.FindEthTxAttemptsRequiringResend(db, olderThan, 1, cltest.FixtureChainID)
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
	})
}

func Test_EthResender_Start(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	// This can be anything as long as it isn't zero
	d := 42 * time.Hour
	cfg.Overrides.GlobalEthTxResendAfterThreshold = &d
	// Set batch size low to test batching
	cfg.Overrides.GlobalEvmRPCDefaultBatchSize = null.IntFrom(1)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	lggr := logger.TestLogger(t)

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		er := txmgr.NewEthResender(lggr, db, ethClient, 100*time.Millisecond, evmcfg)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress, originalBroadcastAt)
		etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress, originalBroadcastAt)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress, time.Now().Add(1*time.Hour))

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

			cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
		}()

		err := db.Get(&etx, `SELECT * FROM eth_txes WHERE id = $1`, etx.ID)
		require.NoError(t, err)
		err = db.Get(&etx2, `SELECT * FROM eth_txes WHERE id = $1`, etx2.ID)
		require.NoError(t, err)

		assert.Greater(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		assert.Greater(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}
