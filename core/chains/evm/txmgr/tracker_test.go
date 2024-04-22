package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const waitTime = 5 * time.Millisecond

func newTestEvmTrackerSetup(t *testing.T) (*txmgr.Tracker, txmgr.TestEvmTxStore, keystore.Eth, []common.Address) {
	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	chainID := big.NewInt(0)
	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, ethKeyStore, *ubig.NewI(chainID.Int64()))
	_, addr2 := cltest.MustInsertRandomKey(t, ethKeyStore, *ubig.NewI(chainID.Int64()))
	enabledAddresses = append(enabledAddresses, addr1, addr2)
	lggr := logger.TestLogger(t)
	return txmgr.NewEvmTracker(txStore, ethKeyStore, chainID, lggr), txStore, ethKeyStore, enabledAddresses
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
	t.Parallel()

	tracker, _, _, _ := newTestEvmTrackerSetup(t)
	ctx := testutils.Context(t)

	require.NoError(t, tracker.Start(ctx))
	require.True(t, tracker.IsStarted())

	t.Run("stop tracker", func(t *testing.T) {
		require.NoError(t, tracker.Close())
		require.False(t, tracker.IsStarted())
	})
}

func TestEvmTracker_AddressTracking(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

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

		err := tracker.Start(ctx)
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		time.Sleep(waitTime)
		addrs := tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, inProgressAddr)
		require.NotContains(t, addrs, unstartedAddr)
		require.Contains(t, addrs, unconfirmedAddr)
	})

	/* TODO: finalized tx state https://smartcontract-it.atlassian.net/browse/BCI-2920
	t.Run("stop tracking finalized tx", func(t *testing.T) {
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)

		err := tracker.Start(ctx)
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		// deliver block before minConfirmations
		tracker.XXXDeliverBlock(1)
		time.Sleep(waitTime)

		addrs := tracker.GetAbandonedAddresses()
		require.Contains(t, addrs, confirmedAddr)

		// deliver block past minConfirmations to finalize tx
		tracker.XXXDeliverBlock(10)
		time.Sleep(waitTime)

		addrs = tracker.GetAbandonedAddresses()
		require.NotContains(t, addrs, confirmedAddr)
	})
	*/
}

func TestEvmTracker_ExceedingTTL(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	t.Run("exceeding ttl", func(t *testing.T) {
		tracker, txStore, _, _ := newTestEvmTrackerSetup(t)
		addr1 := cltest.MustGenerateRandomKey(t).Address
		addr2 := cltest.MustGenerateRandomKey(t).Address
		tx1 := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
		tx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

		tracker.XXXTestSetTTL(time.Nanosecond)
		err := tracker.Start(ctx)
		require.NoError(t, err)
		defer func(tracker *txmgr.Tracker) {
			err = tracker.Close()
			require.NoError(t, err)
		}(tracker)

		time.Sleep(100 * waitTime)
		require.NotContains(t, tracker.GetAbandonedAddresses(), addr1, addr2)

		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		require.NoError(t, err)
		require.True(t, containsID(fatalTxes, tx1.ID))
		require.True(t, containsID(fatalTxes, tx2.ID))
	})
}
