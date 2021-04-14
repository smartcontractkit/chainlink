package bulletprooftxmanager_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NonceSyncer_SyncAll(t *testing.T) {
	t.Parallel()

	t.Run("returns error if PendingNonceAt fails", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), errors.New("something exploded"))

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		err := ns.SyncAll(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce reflects local nonce", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), nil)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce is behind local nonce", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store, int64(32))

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(31), nil)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 32)

		ethClient.AssertExpectations(t)
	})

	t.Run("fast forwards if chain nonce is ahead of local nonce", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, acct1 := cltest.MustAddRandomKeyToKeystore(t, store, int64(0))
		_, acct2 := cltest.MustAddRandomKeyToKeystore(t, store, int64(32))

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// Nothing to do for acct2
			return acct2 == addr
		})).Return(uint64(32), nil)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// acct1 has chain nonce of 5 which is ahead of local nonce 0
			return acct1 == addr
		})).Return(uint64(5), nil)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		assertDatabaseNonce(t, store, acct1, 5)

		ethClient.AssertExpectations(t)
	})

	t.Run("counts 'in_progress' eth_tx as bumping the local next nonce by 1", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		_, acct1 := cltest.MustAddRandomKeyToKeystore(t, store, int64(0))

		cltest.MustInsertInProgressEthTxWithAttempt(t, store, 1, acct1)

		ethClient := new(mocks.Client)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// acct1 has chain nonce of 1 which is ahead of keys.next_nonce (0)
			// by 1, but does not need to change when taking into account the in_progress tx
			return acct1 == addr
		})).Return(uint64(1), nil)
		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))
		assertDatabaseNonce(t, store, acct1, 0)

		ethClient.AssertExpectations(t)

		ethClient = new(mocks.Client)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// acct1 has chain nonce of 2 which is ahead of keys.next_nonce (0)
			// by 2, but only ahead by 1 if we count the in_progress tx as +1
			return acct1 == addr
		})).Return(uint64(2), nil)
		ns = bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))
		assertDatabaseNonce(t, store, acct1, 1)

		ethClient.AssertExpectations(t)
	})
}

func assertDatabaseNonce(t *testing.T, store *store.Store, from common.Address, nonce int64) {
	t.Helper()

	k, err := store.KeyByAddress(from)
	require.NoError(t, err)
	assert.Equal(t, nonce, k.NextNonce)
}
