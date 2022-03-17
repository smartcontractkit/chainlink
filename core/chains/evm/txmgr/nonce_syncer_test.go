package txmgr_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NonceSyncer_SyncAll(t *testing.T) {
	t.Parallel()

	t.Run("returns error if PendingNonceAt fails", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		cfg := cltest.NewTestGeneralConfig(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

		_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), errors.New("something exploded"))

		ns := txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		sendingKeys := cltest.MustSendingKeyStates(t, ethKeyStore)
		err := ns.SyncAll(context.Background(), sendingKeys)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")

		cltest.AssertCount(t, db, "eth_txes", 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 0)

		assertDatabaseNonce(t, db, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce reflects local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		cfg := cltest.NewTestGeneralConfig(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

		_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), nil)

		ns := txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		sendingKeys := cltest.MustSendingKeyStates(t, ethKeyStore)
		require.NoError(t, ns.SyncAll(context.Background(), sendingKeys))

		cltest.AssertCount(t, db, "eth_txes", 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 0)

		assertDatabaseNonce(t, db, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce is behind local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := cltest.NewTestGeneralConfig(t)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

		k1, _ := cltest.MustInsertRandomKey(t, ethKeyStore, int64(32))

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return k1.Address.Address() == addr
		})).Return(uint64(31), nil)

		ns := txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		sendingKeys := cltest.MustSendingKeyStates(t, ethKeyStore)
		require.NoError(t, ns.SyncAll(context.Background(), sendingKeys))

		cltest.AssertCount(t, db, "eth_txes", 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 0)

		assertDatabaseNonce(t, db, k1.Address.Address(), 32)

		ethClient.AssertExpectations(t)
	})

	t.Run("fast forwards if chain nonce is ahead of local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := cltest.NewTestGeneralConfig(t)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

		_, key1 := cltest.MustInsertRandomKey(t, ethKeyStore, int64(0))
		_, key2 := cltest.MustInsertRandomKey(t, ethKeyStore, int64(32))

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// Nothing to do for key2
			return key2 == addr
		})).Return(uint64(32), nil)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// key1 has chain nonce of 5 which is ahead of local nonce 0
			return key1 == addr
		})).Return(uint64(5), nil)

		ns := txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		sendingKeys := cltest.MustSendingKeyStates(t, ethKeyStore)
		require.NoError(t, ns.SyncAll(context.Background(), sendingKeys))

		assertDatabaseNonce(t, db, key1, 5)

		ethClient.AssertExpectations(t)
	})

	t.Run("counts 'in_progress' eth_tx as bumping the local next nonce by 1", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := cltest.NewTestGeneralConfig(t)
		borm := cltest.NewTxmORM(t, db, cfg)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

		_, key1 := cltest.MustInsertRandomKey(t, ethKeyStore, int64(0))

		cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 1, key1)

		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// key1 has chain nonce of 1 which is ahead of keys.next_nonce (0)
			// by 1, but does not need to change when taking into account the in_progress tx
			return key1 == addr
		})).Return(uint64(1), nil)
		ns := txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		sendingKeys := cltest.MustSendingKeyStates(t, ethKeyStore)
		require.NoError(t, ns.SyncAll(context.Background(), sendingKeys))
		assertDatabaseNonce(t, db, key1, 0)

		ethClient.AssertExpectations(t)

		ethClient = cltest.NewEthClientMockWithDefaultChain(t)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// key1 has chain nonce of 2 which is ahead of keys.next_nonce (0)
			// by 2, but only ahead by 1 if we count the in_progress tx as +1
			return key1 == addr
		})).Return(uint64(2), nil)
		ns = txmgr.NewNonceSyncer(db, logger.TestLogger(t), cfg, ethClient)

		require.NoError(t, ns.SyncAll(context.Background(), sendingKeys))
		assertDatabaseNonce(t, db, key1, 1)

		ethClient.AssertExpectations(t)
	})
}

func assertDatabaseNonce(t *testing.T, db *sqlx.DB, address common.Address, nonce int64) {
	t.Helper()

	var nextNonce int64
	err := db.Get(&nextNonce, `SELECT next_nonce FROM eth_key_states WHERE address = $1`, address)
	require.NoError(t, err)
	assert.Equal(t, nonce, nextNonce)
}
