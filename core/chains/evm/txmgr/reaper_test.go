package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
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
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)
	var nonce int64

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(time.Now())
		assert.NoError(t, err)
	})

	// Finalized with confirmed receipt in block number 5
	mustInsertFinalizedEthTxWithReceipt(t, txStore, from, nonce, 5)
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	t.Run("skips if threshold=0", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 0 * time.Second}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(tomorrow)
		assert.NoError(t, err)

		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("doesn't touch ethtxes with different chain ID", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaperWithChainID(t, txStore, config, tc, big.NewInt(42))

		// use time in the future to ensure that just created txs match the filter
		err := r.ReapTxes(tomorrow)
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "evm.txes", 1)
	})

	t.Run("deletes finalized evm.txes that exceed the age threshold", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(now)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		// use time in the future to ensure that just created txs match the filter
		err = r.ReapTxes(tomorrow)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was old enough
		cltest.AssertCount(t, db, "evm.txes", 0)
	})

	mustInsertFatalErrorEthTx(t, txStore, from)

	t.Run("deletes errored evm.txes that exceed the age threshold", func(t *testing.T) {
		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(now)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		err = r.ReapTxes(tomorrow)
		assert.NoError(t, err)
		// Deleted because it is old enough now
		cltest.AssertCount(t, db, "evm.txes", 0)
	})
	t.Run("does not delete old confirmed txs", func(t *testing.T) {
		mustInsertConfirmedEthTxWithReceipt(t, txStore, from, 2, 2)

		config := txmgrmocks.NewReaperConfig(t)

		tc := &reaperConfig{reaperThreshold: 1 * time.Hour}

		r := newReaper(t, txStore, config, tc)

		err := r.ReapTxes(now)
		assert.NoError(t, err)
		// Didn't delete because eth_tx not finalized
		cltest.AssertCount(t, db, "evm.txes", 1)

		err = r.ReapTxes(tomorrow)
		assert.NoError(t, err)
		// Didn't delete old enough but not finalized
		cltest.AssertCount(t, db, "evm.txes", 1)
	})
}
