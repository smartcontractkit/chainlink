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

func generateEnabledAddresses(t *testing.T, keyStore keystore.Eth) []common.Address {
	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, keyStore)
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore)
	enabledAddresses = append(enabledAddresses, addr1, addr2)
	return enabledAddresses
}

func containsID(txes []*txmgr.Tx, id int64) bool {
	for _, tx := range txes {
		if tx.ID == id {
			return true
		}
	}
	return false
}

func TestEthTracker_Initialization(t *testing.T) {
	t.Parallel()

	tracker, _, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	// error before setting enabled addresses
	err := tracker.TrackAbandonedTxes(ctx)
	assert.Error(t, err)
	assert.False(t, tracker.IsTracking())

	err = tracker.SetEnabledAddresses(generateEnabledAddresses(t, keyStore))
	assert.NoError(t, err)
	err = tracker.TrackAbandonedTxes(ctx)
	assert.NoError(t, err)
	assert.True(t, tracker.IsTracking())

	// error tracking already enabled
	err = tracker.TrackAbandonedTxes(ctx)
	assert.Error(t, err)

	t.Run("reset tracker", func(t *testing.T) {
		tracker.Reset()
		assert.False(t, tracker.IsTracking())
		err = tracker.SetEnabledAddresses(generateEnabledAddresses(t, keyStore))
		assert.NoError(t, err)
		err = tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)
	})
}

func TestEthTracker_AddressTracking(t *testing.T) {
	t.Parallel()

	tracker, txStore, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	t.Run("track abandoned addresses", func(t *testing.T) {
		tracker.Reset()
		err := tracker.SetEnabledAddresses(generateEnabledAddresses(t, keyStore))
		assert.NoError(t, err)
		addr1 := cltest.MustGenerateRandomKey(t).Address
		addr2 := cltest.MustGenerateRandomKey(t).Address

		// Insert abandoned transactions
		_ = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
		_ = cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

		err = tracker.TrackAbandonedTxes(ctx)
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
	enabledAddresses := generateEnabledAddresses(t, keyStore)

	t.Run("in progress transaction still valid", func(t *testing.T) {
		err := tracker.SetEnabledAddresses(enabledAddresses)
		assert.NoError(t, err)

		addr1 := cltest.MustGenerateRandomKey(t).Address
		_ = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)

		err = tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		err = tracker.HandleAbandonedTxes(ctx)
		assert.NoError(t, err)
		assert.Contains(t, tracker.GetAbandonedAddresses(), addr1)
	})

	t.Run("in progress transaction exceeding ttl", func(t *testing.T) {
		tracker.Reset()
		tracker.XXXTestSetTTL(time.Nanosecond)

		err := tracker.SetEnabledAddresses(enabledAddresses)
		assert.NoError(t, err)

		addr1 := cltest.MustGenerateRandomKey(t).Address
		tx1 := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)

		err = tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		// Ensure tx1 is finalized as fatal for exceeding ttl
		err = tracker.HandleAbandonedTxes(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, tracker.GetAbandonedAddresses(), addr1)

		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		assert.NoError(t, err)
		assert.True(t, containsID(fatalTxes, tx1.ID))
	})

	t.Run("unconfirmed transaction exceeding ttl", func(t *testing.T) {
		tracker.Reset()
		tracker.XXXTestSetTTL(time.Nanosecond)

		err := tracker.SetEnabledAddresses(enabledAddresses)
		assert.NoError(t, err)

		addr1 := cltest.MustGenerateRandomKey(t).Address
		tx1 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr1)

		err = tracker.TrackAbandonedTxes(ctx)
		assert.NoError(t, err)

		// Ensure tx1 is finalized as fatal for exceeding ttl
		err = tracker.HandleAbandonedTxes(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, tracker.GetAbandonedAddresses(), addr1)

		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		assert.NoError(t, err)
		assert.True(t, containsID(fatalTxes, tx1.ID))
	})
}
