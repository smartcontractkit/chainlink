package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func newReaperWithChainID(t *testing.T, db *sqlx.DB, cfg txmgrtypes.ReaperConfig, cid big.Int) *txmgr.Reaper {
	return txmgr.NewReaper(logger.TestLogger(t), db, cfg, cid)
}

func newReaper(t *testing.T, db *sqlx.DB, cfg txmgrtypes.ReaperConfig) *txmgr.Reaper {
	return newReaperWithChainID(t, db, cfg, cltest.FixtureChainID)
}

func TestReaper_ReapEthTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTxStore(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	var nonce int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))
		config.On("TxReaperThreshold").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
	})

	// Confirmed in block number 5
	cltest.MustInsertConfirmedEthTxWithReceipt(t, txStore, from, nonce, 5)

	t.Run("skips if threshold=0", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("TxReaperThreshold").Return(0 * time.Second)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)

		cltest.AssertCount(t, db, "eth_txes", 1)
	})

	t.Run("doesn't touch ethtxes with different chain ID", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))
		config.On("TxReaperThreshold").Return(1 * time.Hour)

		r := newReaperWithChainID(t, db, config, *big.NewInt(42))

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "eth_txes", 1)
	})

	t.Run("deletes confirmed eth_txes that exceed the age threshold with at least EVM.FinalityDepth blocks above their receipt", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))
		config.On("TxReaperThreshold").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "eth_txes", 1)

		pgtest.MustExec(t, db, `UPDATE eth_txes SET created_at=$1`, oneDayAgo)

		err = r.ReapEthTxes(12)
		assert.NoError(t, err)
		// Didn't delete because eth_tx although old enough, was still within EVM.FinalityDepth of the current head
		cltest.AssertCount(t, db, "eth_txes", 1)

		err = r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past EVM.FinalityDepth
		cltest.AssertCount(t, db, "eth_txes", 0)
	})

	cltest.MustInsertFatalErrorEthTx(t, txStore, from)

	t.Run("deletes errored eth_txes that exceed the age threshold", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))
		config.On("TxReaperThreshold").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, utils.JustError(db.Exec(`UPDATE eth_txes SET created_at=$1`, oneDayAgo)))

		err = r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Deleted because it is old enough now
		cltest.AssertCount(t, db, "eth_txes", 0)
	})
}
