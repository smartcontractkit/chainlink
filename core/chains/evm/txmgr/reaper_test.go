package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func newReaperWithChainID(t *testing.T, db txmgrtypes.TxHistoryReaper[*big.Int], cfg txmgrtypes.ReaperChainConfig, txConfig txmgrtypes.ReaperTransactionsConfig, cid *big.Int) *txmgr.Reaper {
	return txmgr.NewEvmReaper(logger.Test(t), db, cfg, txConfig, cid)
}

func newReaper(t *testing.T, db txmgrtypes.TxHistoryReaper[*big.Int], cfg txmgrtypes.ReaperChainConfig, txConfig txmgrtypes.ReaperTransactionsConfig) *txmgr.Reaper {
	return newReaperWithChainID(t, db, cfg, txConfig, &cltest.FixtureChainID)
}

type reaperConfig struct {
	reaperInterval  time.Duration
	reaperThreshold time.Duration
}

func (r *reaperConfig) ReaperInterval() time.Duration {
	return r.reaperInterval
}

func (r *reaperConfig) ReaperThreshold() time.Duration {
	return r.reaperThreshold
}

func TestReaper_ReapTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)
	var nonce int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
	})

	// Confirmed in block number 5
	mustInsertConfirmedEthTxWithReceipt(t, txStore, from, nonce, 5)

	t.Run("skips if threshold=0", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 0 * time.Second}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)

		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("doesn't touch ethtxes with different chain ID", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaperWithChainID(t, txStore, config, tc, big.NewInt(42))

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("deletes confirmed evm.txes that exceed the age threshold with at least EVM.FinalityDepth blocks above their receipt", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET created_at=$1`, oneDayAgo)

		err = r.ReapTxes(12)
		assert.NoError(t, err)
		// Didn't delete because eth_tx although old enough, was still within EVM.FinalityDepth of the current head
		cltest.AssertCount(t, db, "evm.txes", 1)

		err = r.ReapTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past EVM.FinalityDepth
		cltest.AssertCount(t, db, "evm.txes", 0)
	})

	mustInsertFatalErrorEthTx(t, txStore, from)

	t.Run("deletes errored evm.txes that exceed the age threshold", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)
		config.On("FinalityDepth").Return(uint32(10))

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		require.NoError(t, utils.JustError(db.Exec(`UPDATE evm.txes SET created_at=$1`, oneDayAgo)))

		err = r.ReapTxes(42)
		assert.NoError(t, err)
		// Deleted because it is old enough now
		cltest.AssertCount(t, db, "evm.txes", 0)
	})
}
