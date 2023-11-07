package txmgr_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func newTestEvmTrackerSetup(t *testing.T) (*txmgr.Tracker, txmgr.TestEvmTxStore, keystore.Eth) {
	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	lggr := logger.TestLogger(t)
	return txmgr.NewEvmTracker(txStore, lggr), txStore, ethKeyStore
}

func TestEthTracker_AddressTracking(t *testing.T) {
	t.Parallel()

	tracker, txStore, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, keyStore)
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore)
	enabledAddresses = append(enabledAddresses, addr1, addr2)

	t.Run("unset enabledAddresses", func(t *testing.T) {
		err := tracker.TrackAbandonedTxes(ctx)
		assert.Error(t, err)

		tracker.SetEnabledAddresses(enabledAddresses)
		err = tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)
	})

	t.Run("track abandoned addresses", func(t *testing.T) {
		tracker.SetEnabledAddresses(enabledAddresses)
		addr1 := cltest.MustGenerateRandomKey(t).Address
		addr2 := cltest.MustGenerateRandomKey(t).Address

		// Insert abandoned transactions
		_ = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
		_ = cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

		err := tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		addrs := tracker.GetAbandonedAddresses()
		assert.Contains(t, addrs, addr1)
		assert.Contains(t, addrs, addr2)
	})
}

func TestEthTracker_ExceedingTTL(t *testing.T) {
	t.Parallel()

	tracker, txStore, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, keyStore)
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore)
	enabledAddresses = append(enabledAddresses, addr1, addr2)
	tracker.SetEnabledAddresses(enabledAddresses)

	t.Run("in progress transaction still valid", func(t *testing.T) {
		addr1 := cltest.MustGenerateRandomKey(t).Address
		_ = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)

		err := tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		tracker.HandleAbandonedTxes(ctx)
		assert.Contains(t, tracker.GetAbandonedAddresses(), addr1)
	})

	t.Run("in progress transaction exceeding ttl", func(t *testing.T) {
		tracker.XXXTestSetTTL(time.Nanosecond)

		addr1 := cltest.MustGenerateRandomKey(t).Address
		tx1 := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)

		err := tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		// Ensure tx1 is finalized as fatal for exceeding ttl
		tracker.HandleAbandonedTxes(ctx)
		assert.Empty(t, tracker.GetAbandonedAddresses())

		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		assert.NoError(t, err)
		assert.Equal(t, fatalTxes[0].ID, tx1.ID)
	})
}
