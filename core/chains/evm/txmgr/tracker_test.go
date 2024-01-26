package txmgr_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const waitTime = 5 * time.Millisecond

func newTestEvmTrackerSetup(t *testing.T) (*txmgr.Tracker, txmgr.TestEvmTxStore, keystore.Eth, []common.Address) {
	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	chainID := big.NewInt(0)
	enabledAddresses := generateEnabledAddresses(t, ethKeyStore, chainID)
	lggr := logger.TestLogger(t)
	return txmgr.NewEvmTracker(txStore, ethKeyStore, chainID, lggr), txStore, ethKeyStore, enabledAddresses
}

func generateEnabledAddresses(t *testing.T, keyStore keystore.Eth, chainID *big.Int) []common.Address {
	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, keyStore, *ubig.NewI(chainID.Int64()))
	_, addr2 := cltest.MustInsertRandomKey(t, keyStore, *ubig.NewI(chainID.Int64()))
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

func TestEvmTracker_Initialization(t *testing.T) {
	t.Skip("BCI-2638 tracker disabled")
	t.Parallel()

	tracker, _, _, _ := newTestEvmTrackerSetup(t)

	err := tracker.Start(context.Background())
	require.NoError(t, err)
	require.True(t, tracker.IsStarted())

	t.Run("stop tracker", func(t *testing.T) {
		err := tracker.Close()
		require.NoError(t, err)
		require.False(t, tracker.IsStarted())
	})
}

func TestEvmTracker_AddressTracking(t *testing.T) {
	t.Skip("BCI-2638 tracker disabled")
	t.Parallel()

	t.Run("track abandoned addresses", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		inProgressAddr := cltest.MustGenerateRandomKey(t).Address
		unstartedAddr := cltest.MustGenerateRandomKey(t).Address
		unconfirmedAddr := cltest.MustGenerateRandomKey(t).Address
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = mustInsertInProgressEthTxWithAttempt(t, txStore, 123, inProgressAddr)
		_ = cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, unconfirmedAddr)
		_ = mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)
		_ = mustCreateUnstartedTx(t, txStore, unstartedAddr, cltest.MustGenerateRandomKey(t).Address, []byte{}, 0, big.Int{}, ethClient.ConfiguredChainID())

		err := tracker.Start(context.Background())
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		addrs := tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, inProgressAddr)
		require.NotContains(t, addrs, unstartedAddr)
		require.Contains(t, addrs, confirmedAddr)
		require.Contains(t, addrs, unconfirmedAddr)
	})

	t.Run("stop tracking finalized tx", func(t *testing.T) {
		t.Skip("BCI-2638 tracker disabled")
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)

		err := tracker.Start(context.Background())
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		addrs := tracker.GetAbandonedAddresses()
		require.Contains(t, addrs, confirmedAddr)

		// deliver block past minConfirmations to finalize tx
		tracker.XXXDeliverBlock(10)
		time.Sleep(waitTime)

		addrs = tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, confirmedAddr)
	})
}

func TestEvmTracker_ExceedingTTL(t *testing.T) {
	t.Skip("BCI-2638 tracker disabled")
	t.Parallel()

	t.Run("confirmed but unfinalized transaction still tracked", func(t *testing.T) {
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		addr1 := cltest.MustGenerateRandomKey(t).Address
		_ = mustInsertConfirmedEthTxWithReceipt(t, txStore, addr1, 123, 1)

		err := tracker.Start(context.Background())
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		require.Contains(t, tracker.GetAbandonedAddresses(), addr1)
	})

	t.Run("exceeding ttl", func(t *testing.T) {
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		addr1 := cltest.MustGenerateRandomKey(t).Address
		addr2 := cltest.MustGenerateRandomKey(t).Address
		tx1 := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
		tx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

		tracker.XXXTestSetTTL(time.Nanosecond)
		err := tracker.Start(context.Background())
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		time.Sleep(waitTime)
		require.NotContains(t, tracker.GetAbandonedAddresses(), addr1, addr2)

		fatalTxes, err := txStore.GetFatalTransactions(context.Background())
		require.NoError(t, err)
		require.True(t, containsID(fatalTxes, tx1.ID))
		require.True(t, containsID(fatalTxes, tx2.ID))
	})
}
