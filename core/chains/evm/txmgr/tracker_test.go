package txmgr_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const waitTime = 5 * time.Millisecond

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

	err := tracker.Start(ctx, generateEnabledAddresses(t, keyStore))
	require.NoError(t, err)
	require.True(t, tracker.IsStarted())

	t.Run("reset tracker", func(t *testing.T) {
		tracker.Stop()
		require.False(t, tracker.IsStarted())

		err := tracker.Start(ctx, generateEnabledAddresses(t, keyStore))
		require.NoError(t, err)
		require.True(t, tracker.IsStarted())
	})
}

func TestEthTracker_AddressTracking(t *testing.T) {
	t.Parallel()

	tracker, txStore, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	t.Run("track abandoned addresses", func(t *testing.T) {
		inProgressAddr := cltest.MustGenerateRandomKey(t).Address
		unconfirmedAddr := cltest.MustGenerateRandomKey(t).Address
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, inProgressAddr)
		_ = cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, unconfirmedAddr)
		_ = cltest.MustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)

		err := tracker.Start(ctx, generateEnabledAddresses(t, keyStore))
		defer tracker.Stop()
		require.NoError(t, err)

		addrs := tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, inProgressAddr)
		require.Contains(t, addrs, confirmedAddr)
		require.Contains(t, addrs, unconfirmedAddr)
	})

	t.Run("stop tracking finalized tx", func(t *testing.T) {
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = cltest.MustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)

		err := tracker.Start(ctx, generateEnabledAddresses(t, keyStore))
		defer tracker.Stop()
		require.NoError(t, err)

		// deliver block past minConfirmations to finalize tx
		tracker.XXXDeliverBlock(10)
		time.Sleep(waitTime)

		addrs := tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, confirmedAddr)
	})
}

func TestEthTracker_ExceedingTTL(t *testing.T) {
	t.Parallel()

	tracker, txStore, keyStore := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)
	enabledAddresses := generateEnabledAddresses(t, keyStore)

	t.Run("confirmed but unfinalized transaction still tracked", func(t *testing.T) {
		addr1 := cltest.MustGenerateRandomKey(t).Address
		_ = cltest.MustInsertConfirmedEthTxWithReceipt(t, txStore, addr1, 123, 1)

		err := tracker.Start(ctx, enabledAddresses)
		defer tracker.Stop()
		require.NoError(t, err)
		require.Contains(t, tracker.GetAbandonedAddresses(), addr1)
	})

	t.Run("exceeding ttl", func(t *testing.T) {
		addr1 := cltest.MustGenerateRandomKey(t).Address
		addr2 := cltest.MustGenerateRandomKey(t).Address
		tx1 := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
		tx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

		tracker.XXXTestSetTTL(time.Nanosecond)
		err := tracker.Start(ctx, enabledAddresses)
		defer tracker.Stop()
		require.NoError(t, err)

		time.Sleep(waitTime)
		require.NotContains(t, tracker.GetAbandonedAddresses(), addr1, addr2)

		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		require.NoError(t, err)
		require.True(t, containsID(fatalTxes, tx1.ID))
		require.True(t, containsID(fatalTxes, tx2.ID))
	})
}
