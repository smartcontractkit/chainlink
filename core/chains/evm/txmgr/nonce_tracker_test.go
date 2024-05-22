package txmgr_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	clientmock "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txstoremock "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestNonceTracker_LoadSequenceMap(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr1 := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	addr2 := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781140")
	enabledAddresses := []common.Address{addr1, addr2}

	t.Run("set next nonce using entries from tx table", func(t *testing.T) {
		randNonce1 := testutils.NewRandomPositiveInt64()
		randNonce2 := testutils.NewRandomPositiveInt64()
		txStore.On("FindLatestSequence", mock.Anything, addr1, chainID).Return(types.Nonce(randNonce1), nil).Once()
		txStore.On("FindLatestSequence", mock.Anything, addr2, chainID).Return(types.Nonce(randNonce2), nil).Once()

		nonceTracker.LoadNextSequences(ctx, enabledAddresses)
		seq, err := nonceTracker.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(randNonce1+1), seq)
		seq, err = nonceTracker.GetNextSequence(ctx, addr2)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(randNonce2+1), seq)
	})

	t.Run("set next nonce using client when not found in tx table", func(t *testing.T) {
		var emptyNonce types.Nonce
		txStore.On("FindLatestSequence", mock.Anything, addr1, chainID).Return(emptyNonce, errors.New("no rows")).Once()
		txStore.On("FindLatestSequence", mock.Anything, addr2, chainID).Return(emptyNonce, errors.New("no rows")).Once()

		randNonce1 := testutils.NewRandomPositiveInt64()
		randNonce2 := testutils.NewRandomPositiveInt64()
		client.On("PendingNonceAt", mock.Anything, addr1).Return(uint64(randNonce1), nil).Once()
		client.On("PendingNonceAt", mock.Anything, addr2).Return(uint64(randNonce2), nil).Once()

		nonceTracker.LoadNextSequences(ctx, enabledAddresses)
		seq, err := nonceTracker.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(randNonce1), seq)
		seq, err = nonceTracker.GetNextSequence(ctx, addr2)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(randNonce2), seq)
	})
}

func TestNonceTracker_syncOnChain(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")

	t.Run("throws error if RPC call fails", func(t *testing.T) {
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()

		err := nonceTracker.SyncOnChain(ctx, addr, types.Nonce(2))
		require.Error(t, err)
	})

	t.Run("uses local nonce instead of on-chain nonce if on-chain nonce is lower", func(t *testing.T) {
		nonce := 2
		newNonce := 5
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(nonce), nil).Once()

		enabledAddresses := []common.Address{}
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		// syncOnChain will set the next sequence even if the address is not present in the map
		err := nonceTracker.SyncOnChain(ctx, addr, types.Nonce(newNonce))
		require.NoError(t, err)

		seq, err := nonceTracker.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(newNonce), seq)
	})

	t.Run("fast forwards nonce if on-chain nonce is higher than local nonce", func(t *testing.T) {
		nonce := 10
		onChainNonce := 5
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(nonce), nil).Once()

		enabledAddresses := []common.Address{}
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		// syncOnChain will set the next sequence even if the address is not present in the map
		err := nonceTracker.SyncOnChain(ctx, addr, types.Nonce(onChainNonce))
		require.NoError(t, err)

		seq, err := nonceTracker.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(nonce), seq)
	})
}

func TestNonceTracker_SyncSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	enabledAddresses := []common.Address{addr}

	t.Run("syncs sequence successfully", func(t *testing.T) {
		txStoreNonce := 2
		onChainNonce := 3
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(txStoreNonce), nil).Once()
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		var chStop services.StopChan
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(onChainNonce), nil).Once()
		nonceTracker.SyncSequence(ctx, addr, chStop)

		seq, err := nonceTracker.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(onChainNonce), seq)
	})

	t.Run("retries if on-chain syncing fails", func(t *testing.T) {
		txStoreNonce := 2
		onChainNonce := 3
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(txStoreNonce), nil).Once()
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		var chStop services.StopChan
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(onChainNonce), nil).Once()
		nonceTracker.SyncSequence(ctx, addr, chStop)

		seq, err := nonceTracker.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(onChainNonce), seq)
	})
}

func TestNonceTracker_GetNextSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")

	t.Run("fails to get sequence if address doesn't exist in map", func(t *testing.T) {
		_, err := nonceTracker.GetNextSequence(ctx, addr)
		require.Error(t, err)
	})

	t.Run("fails to get sequence if address doesn't exist in map and is disabled", func(t *testing.T) {
		_, err := nonceTracker.GetNextSequence(ctx, addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("address disabled: %s", addr.Hex()))
	})

	t.Run("fails to get sequence if address is enabled, doesn't exist in map, and getSequenceForAddr fails", func(t *testing.T) {
		enabledAddresses := []common.Address{addr}
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(0), errors.New("no rows")).Twice()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Twice()
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		_, err := nonceTracker.GetNextSequence(ctx, addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("failed to find next sequence for address: %s", addr.Hex()))
	})

	t.Run("gets next sequence successfully if there is no entry in map but address is enabled and getSequenceForAddr is successful", func(t *testing.T) {
		txStoreNonce := 4
		enabledAddresses := []common.Address{addr}
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(0), errors.New("no rows")).Once()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()
		nonceTracker.LoadNextSequences(ctx, enabledAddresses)

		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(txStoreNonce), nil).Once()
		seq, err := nonceTracker.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(txStoreNonce+1), seq)
	})
}

func TestNonceTracker_GenerateNextSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	enabledAddresses := []common.Address{addr}

	randNonce := testutils.NewRandomPositiveInt64()
	txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(randNonce), nil).Once()
	nonceTracker.LoadNextSequences(ctx, enabledAddresses)
	seq, err := nonceTracker.GetNextSequence(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, types.Nonce(randNonce+1), seq) // Local nonce should be highest nonce in DB + 1

	nonceTracker.GenerateNextSequence(addr, types.Nonce(randNonce+1))

	seq, err = nonceTracker.GetNextSequence(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, types.Nonce(randNonce+2), seq) // GenerateNextSequence increases local nonce by 1
}

func Test_SetNonceAfterInit(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	nonceTracker := txmgr.NewNonceTracker(logger.Test(t), txStore, txmgr.NewEvmTxmClient(client, nil))

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	enabledAddresses := []common.Address{addr}
	randNonce := testutils.NewRandomPositiveInt64()
	client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("failed to retrieve nonce at startup")).Once()
	client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(randNonce), nil).Once()
	nonceTracker.LoadNextSequences(ctx, enabledAddresses)

	nonce, err := nonceTracker.GetNextSequence(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, randNonce, int64(nonce))

	// Test that the new nonce is set in the map and does not need a client call to retrieve on subsequent calls
	nonce, err = nonceTracker.GetNextSequence(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, randNonce, int64(nonce))
}
