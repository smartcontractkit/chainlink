package txmgr_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NonceSyncer_Sync(t *testing.T) {
	t.Parallel()

	t.Run("returns error if PendingNonceAt fails", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

		_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

		ns := txmgr.NewNonceSyncer(txStore, logger.TestLogger(t), ethClient)

		ethClient.On("PendingNonceAt", mock.Anything, from).Return(uint64(0), errors.New("something exploded"))
		_, err := ns.Sync(testutils.Context(t), from, types.Nonce(0))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")

		cltest.AssertCount(t, db, "evm.txes", 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 0)
	})

	t.Run("does nothing if chain nonce reflects local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

		_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

		ns := txmgr.NewNonceSyncer(txStore, logger.TestLogger(t), ethClient)

		ethClient.On("PendingNonceAt", mock.Anything, from).Return(uint64(0), nil)

		nonce, err := ns.Sync(testutils.Context(t), from, 0)
		require.Equal(t, nonce.Int64(), int64(0))
		require.NoError(t, err)

		cltest.AssertCount(t, db, "evm.txes", 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 0)
	})

	t.Run("does nothing if chain nonce is behind local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ks := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		_, fromAddress := cltest.RandomKey{Nonce: 32}.MustInsert(t, ks)

		ns := txmgr.NewNonceSyncer(txStore, logger.TestLogger(t), ethClient)

		// Used to mock the chain nonce
		ethClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(5), nil)
		nonce, err := ns.Sync(testutils.Context(t), fromAddress, types.Nonce(32))
		require.Equal(t, nonce.Int64(), int64(32))
		require.NoError(t, err)

		cltest.AssertCount(t, db, "evm.txes", 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 0)
	})

	t.Run("fast forwards if chain nonce is ahead of local nonce", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

		_, key1 := cltest.MustInsertRandomKey(t, ethKeyStore)
		_, key2 := cltest.RandomKey{Nonce: 32}.MustInsert(t, ethKeyStore)

		key1LocalNonce := types.Nonce(0)
		key2LocalNonce := types.Nonce(32)

		ns := txmgr.NewNonceSyncer(txStore, logger.TestLogger(t), ethClient)

		// Used to mock the chain nonce
		ethClient.On("PendingNonceAt", mock.Anything, key1).Return(uint64(5), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, key2).Return(uint64(32), nil).Once()

		syncerNonce, err := ns.Sync(testutils.Context(t), key1, key1LocalNonce)
		require.NoError(t, err)
		require.Greater(t, syncerNonce, key1LocalNonce)

		syncerNonce, err = ns.Sync(testutils.Context(t), key2, key2LocalNonce)
		require.NoError(t, err)
		require.Equal(t, syncerNonce, key2LocalNonce)
	})
}
