package txmgr_test

import (
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func newTestEvmTrackerSetup(t *testing.T) (*txmgr.Tracker, txmgr.TestEvmTxStore) {
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	chainID := big.NewInt(0)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	var enabledAddresses []common.Address
	_, addr1 := cltest.MustInsertRandomKey(t, ethKeyStore, *ubig.NewI(chainID.Int64()))
	_, addr2 := cltest.MustInsertRandomKey(t, ethKeyStore, *ubig.NewI(chainID.Int64()))
	enabledAddresses = append(enabledAddresses, addr1, addr2)
	lggr := logger.Test(t)
	return txmgr.NewEvmTracker(txStore, keys.NewStore(keystore.NewEthSigner(ethKeyStore, chainID)), chainID, lggr), txStore
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

	tracker, _ := newTestEvmTrackerSetup(t)
	ctx := tests.Context(t)

	require.NoError(t, tracker.Start(ctx))
	require.True(t, tracker.IsStarted())

	t.Run("stop tracker", func(t *testing.T) {
		require.NoError(t, tracker.Close())
		require.False(t, tracker.IsStarted())
	})
}

func TestEvmTracker_AddressTracking(t *testing.T) {
	t.Parallel()

	t.Run("track abandoned addresses", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		tracker, txStore := newTestEvmTrackerSetup(t)
		inProgressAddr := cltest.MustGenerateRandomKey(t).Address
		unstartedAddr := cltest.MustGenerateRandomKey(t).Address
		unconfirmedAddr := cltest.MustGenerateRandomKey(t).Address
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		_ = mustInsertInProgressEthTxWithAttempt(t, txStore, 123, inProgressAddr)
		_ = cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, unconfirmedAddr)
		_ = mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)
		_ = mustCreateUnstartedTx(t, txStore, unstartedAddr, cltest.MustGenerateRandomKey(t).Address, []byte{}, 0, big.Int{}, ethClient.ConfiguredChainID())

		servicetest.Run(t, tracker)

		require.Eventually(t, func() bool {
			addrs := tracker.GetAbandonedAddresses()
			t.Logf("Addresses: %v", addrs)
			return !slices.Contains(addrs, inProgressAddr) &&
				!slices.Contains(addrs, unstartedAddr) &&
				slices.Contains(addrs, unconfirmedAddr)
		}, tests.WaitTimeout(t), time.Second)
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
	ctx := tests.Context(t)

	tracker, txStore := newTestEvmTrackerSetup(t)
	addr1 := cltest.MustGenerateRandomKey(t).Address
	addr2 := cltest.MustGenerateRandomKey(t).Address
	tx1 := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, addr1)
	tx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, addr2)

	tracker.XXXTestSetTTL(time.Nanosecond)
	servicetest.Run(t, tracker)

	require.Eventually(t, func() bool {
		addrs := tracker.GetAbandonedAddresses()
		return !slices.Contains(addrs, addr1) && !slices.Contains(addrs, addr2)
	}, tests.WaitTimeout(t), time.Second)

	require.Eventually(t, func() bool {
		fatalTxes, err := txStore.GetFatalTransactions(ctx)
		require.NoError(t, err)
		return containsID(fatalTxes, tx1.ID) && containsID(fatalTxes, tx2.ID)
	}, tests.WaitTimeout(t), time.Second)
}
