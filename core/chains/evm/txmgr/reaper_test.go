package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func newReaperWithChainID(t *testing.T, db txmgrtypes.TxHistoryReaper[*big.Int], txConfig txmgrtypes.ReaperTransactionsConfig, cid *big.Int) *txmgr.Reaper {
	return txmgr.NewEvmReaper(logger.Test(t), db, txConfig, cid)
}

func newReaper(t *testing.T, db txmgrtypes.TxHistoryReaper[*big.Int], txConfig txmgrtypes.ReaperTransactionsConfig) *txmgr.Reaper {
	return newReaperWithChainID(t, db, txConfig, &cltest.FixtureChainID)
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

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)
	var nonce int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
	})

	// Confirmed in block number 5
	mustInsertConfirmedEthTxWithReceipt(t, txStore, from, nonce, 5)

	t.Run("skips if threshold=0", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 0 * time.Second}

		r := newReaper(t, txStore, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)

		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("doesn't touch ethtxes with different chain ID", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaperWithChainID(t, txStore, tc, big.NewInt(42))

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("deletes finalized evm.txes that exceed the age threshold", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		testutils.MustExec(t, db, `UPDATE evm.txes SET created_at=$1, state='finalized'`, oneDayAgo)

		err = r.ReapTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past the age threshold
		cltest.AssertCount(t, db, "evm.txes", 0)
	})

	mustInsertFatalErrorEthTx(t, txStore, from)

	t.Run("deletes errored evm.txes that exceed the age threshold", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, tc)

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

	mustInsertConfirmedEthTxWithReceipt(t, txStore, from, 0, 42)

	t.Run("deletes confirmed evm.txes that exceed the age threshold", func(t *testing.T) {
		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, tc)

		err := r.ReapTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		testutils.MustExec(t, db, `UPDATE evm.txes SET created_at=$1`, oneDayAgo)

		err = r.ReapTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past the age threshold
		cltest.AssertCount(t, db, "evm.txes", 0)
	})
}
