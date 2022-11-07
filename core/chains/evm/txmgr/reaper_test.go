package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newReaperWithChainID(t *testing.T, db *sqlx.DB, cfg txmgr.ReaperConfig, cid big.Int) *txmgr.Reaper {
	return txmgr.NewReaper(logger.TestLogger(t), db, cfg, cid)
}

func newReaper(t *testing.T, db *sqlx.DB, cfg txmgr.ReaperConfig) *txmgr.Reaper {
	return newReaperWithChainID(t, db, cfg, cltest.FixtureChainID)
}

func TestReaper_ReapEthTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	var nonce int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EvmFinalityDepth").Return(uint32(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
	})

	// Confirmed in block number 5
	cltest.MustInsertConfirmedEthTxWithReceipt(t, borm, from, nonce, 5)

	t.Run("skips if threshold=0", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EvmFinalityDepth").Return(uint32(10))
		config.On("EthTxReaperThreshold").Return(0 * time.Second)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)

		cltest.AssertCount(t, db, "eth_txes", 1)
	})

	t.Run("doesn't touch ethtxes with different chain ID", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EvmFinalityDepth").Return(uint32(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := newReaperWithChainID(t, db, config, *big.NewInt(42))

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "eth_txes", 1)
	})

	t.Run("deletes confirmed eth_txes that exceed the age threshold with at least ETH_FINALITY_DEPTH blocks above their receipt", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EvmFinalityDepth").Return(uint32(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := newReaper(t, db, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "eth_txes", 1)

		pgtest.MustExec(t, db, `UPDATE eth_txes SET created_at=$1`, oneDayAgo)

		err = r.ReapEthTxes(12)
		assert.NoError(t, err)
		// Didn't delete because eth_tx although old enough, was still within ETH_FINALITY_DEPTH of the current head
		cltest.AssertCount(t, db, "eth_txes", 1)

		err = r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past ETH_FINALITY_DEPTH
		cltest.AssertCount(t, db, "eth_txes", 0)
	})

	cltest.MustInsertFatalErrorEthTx(t, borm, from)

	t.Run("deletes errored eth_txes that exceed the age threshold", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EvmFinalityDepth").Return(uint32(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

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
