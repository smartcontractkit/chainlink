package bulletprooftxmanager_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_EthResender_FindEthTxesRequiringResend(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	key := cltest.MustInsertRandomKey(t, store.DB)
	fromAddress := key.Address.Address()

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := bulletprooftxmanager.FindEthTxesRequiringResend(store.DB, olderThan)
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	etxs := []models.EthTx{
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0, fromAddress, time.Unix(1616509100, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, fromAddress, time.Unix(1616509200, 0)),
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2, fromAddress, time.Unix(1616509300, 0)),
	}
	attempt1_2 := newBroadcastEthTxAttempt(t, etxs[0].ID, store)
	attempt1_2.GasPrice = *utils.NewBig(big.NewInt(10))
	require.NoError(t, store.DB.Create(&attempt1_2).Error)

	attempt3_2 := newInProgressEthTxAttempt(t, etxs[2].ID, store)
	attempt3_2.GasPrice = *utils.NewBig(big.NewInt(10))
	require.NoError(t, store.DB.Create(&attempt3_2).Error)

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := bulletprooftxmanager.FindEthTxesRequiringResend(store.DB, olderThan)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].EthTxAttempts[0].ID, attempts[1].ID)
	})
}

func Test_EthResender_Start(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	key := cltest.MustInsertRandomKey(t, store.DB)
	fromAddress := key.Address.Address()

	t.Run("resends transactions that have been languishing unconfirmed for too long", func(t *testing.T) {
		ethClient := new(mocks.Client)

		er := bulletprooftxmanager.NewEthResender(store.DB, ethClient, 100*time.Millisecond, 1*time.Hour)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0, fromAddress, originalBroadcastAt)
		etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, fromAddress, originalBroadcastAt)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2, fromAddress, time.Now().Add(1*time.Hour))

		ethClient.On("RoundRobinBatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_sendRawTransaction" && b[0].Args[0] == hexutil.Encode(etx.EthTxAttempts[0].SignedRawTx) &&
				b[1].Method == "eth_sendRawTransaction" && b[1].Args[0] == hexutil.Encode(etx2.EthTxAttempts[0].SignedRawTx)
		})).Return(nil)

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
