package bulletprooftxmanager_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_EthResender_FindEthTxesRequiringResend(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	key := cltest.MustInsertRandomKey(t, store.DB)
	fromAddress := key.Address.Address()

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := bulletprooftxmanager.FindEthTxesRequiringResend(store.DB, olderThan, 10)
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	etxs := []bulletprooftxmanager.EthTx{
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, fromAddress, time.Unix(1616509100, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 1, fromAddress, time.Unix(1616509200, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress, time.Unix(1616509300, 0)),
	}
	attempt1_2 := newBroadcastEthTxAttempt(t, etxs[0].ID)
	attempt1_2.GasPrice = *utils.NewBig(big.NewInt(10))
	require.NoError(t, store.DB.Create(&attempt1_2).Error)

	attempt3_2 := newInProgressEthTxAttempt(t, etxs[2].ID)
	attempt3_2.GasPrice = *utils.NewBig(big.NewInt(10))
	require.NoError(t, store.DB.Create(&attempt3_2).Error)

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := bulletprooftxmanager.FindEthTxesRequiringResend(store.DB, olderThan, 0)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].EthTxAttempts[0].ID, attempts[1].ID)
	})

	t.Run("applies limit", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := bulletprooftxmanager.FindEthTxesRequiringResend(store.DB, olderThan, 1)
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
	})
}

func Test_EthResender_Start(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	db := store.DB
	// This can be anything as long as it isn't zero
	store.Config.Set("ETH_TX_RESEND_AFTER_THRESHOLD", "42h")
	// Set batch size low to test batching
	store.Config.Set("ETH_RPC_DEFAULT_BATCH_SIZE", "1")
	key := cltest.MustInsertRandomKey(t, store.DB)
	fromAddress := key.Address.Address()

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := new(mocks.Client)

		er := bulletprooftxmanager.NewEthResender(store.DB, ethClient, 100*time.Millisecond, store.Config)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, fromAddress, originalBroadcastAt)
		etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 1, fromAddress, originalBroadcastAt)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress, time.Now().Add(1*time.Hour))

		// First batch of 1
		ethClient.On("RoundRobinBatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx)
		})).Return(nil)
		// Second batch of 1
		ethClient.On("RoundRobinBatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx2.EthTxAttempts[0].SignedRawTx)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// It should update BroadcastAt even if there is an error here
			elems[0].Error = errors.New("kaboom")
		})

		func() {
			er.Start()
			defer er.Stop()

			cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
		}()

		err := store.DB.First(&etx).Error
		require.NoError(t, err)
		err = store.DB.First(&etx2).Error
		require.NoError(t, err)

		assert.Greater(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		assert.Greater(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}
