package txm

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
	txstoremock "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func NewTestSequenceSyncer(t testing.TB, txStore SequenceSyncerTxStore, client SequenceSyncerClient) *sequenceSyncer {
	t.Helper()

	lggr := logger.Test(t)
	return NewSequenceSyncer(lggr, txStore, client)
}

func TestSequenceSyncer_LoadSequenceMap(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	ss := NewTestSequenceSyncer(t, txStore, client)

	addr1 := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	addr2 := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781140")
	enabledAddresses := []common.Address{addr1, addr2}

	t.Run("set next nonce using entries from tx table", func(t *testing.T) {
		txStore.On("FindLatestSequence", mock.Anything, addr1, chainID).Return(types.Nonce(2), nil).Once()
		txStore.On("FindLatestSequence", mock.Anything, addr2, chainID).Return(types.Nonce(5), nil).Once()

		ss.LoadNextSequenceMap(ctx, enabledAddresses)
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(3), seq)

	})

	t.Run("set next nonce using client when not found in tx table", func(t *testing.T) {
		var emptyNonce types.Nonce
		txStore.On("FindLatestSequence", mock.Anything, addr1, chainID).Return(emptyNonce, errors.New("no rows")).Once()
		txStore.On("FindLatestSequence", mock.Anything, addr2, chainID).Return(emptyNonce, errors.New("no rows")).Once()

		client.On("PendingNonceAt", mock.Anything, addr1).Return(uint64(2), nil).Once()
		client.On("PendingNonceAt", mock.Anything, addr2).Return(uint64(5), nil).Once()

		ss.LoadNextSequenceMap(ctx, enabledAddresses)
		seq, err := ss.GetNextSequence(ctx, addr1)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(2), seq)

	})

}

func TestSequenceSyncer_syncOnChain(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	ss := NewTestSequenceSyncer(t, txStore, client)

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")

	t.Run("throws error if RPC call fails", func(t *testing.T) {
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()

		err := ss.syncOnChain(ctx, addr, types.Nonce(2))
		require.Error(t, err)
	})

	t.Run("uses local nonce instead of on-chain nonce if on-chain nonce is lower", func(t *testing.T) {
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(2), nil).Once()

		enabledAddresses := []common.Address{}
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		// syncOnChain will set the next sequence even if the address is not present in the map
		err := ss.syncOnChain(ctx, addr, types.Nonce(5))
		require.NoError(t, err)

		seq, err := ss.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(5), seq)
	})

	t.Run("fast forwards nonce if on-chain nonce is higher than local nonce", func(t *testing.T) {
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(10), nil).Once()

		enabledAddresses := []common.Address{}
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		// syncOnChain will set the next sequence even if the address is not present in the map
		err := ss.syncOnChain(ctx, addr, types.Nonce(5))
		require.NoError(t, err)

		seq, err := ss.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(10), seq)
	})

}

func TestSequenceSyncer_IncreamentNextSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	ss := NewTestSequenceSyncer(t, txStore, client)

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	enabledAddresses := []common.Address{addr}

	txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(2), nil).Once()
	ss.LoadNextSequenceMap(ctx, enabledAddresses)
	ss.IncrementNextSequence(addr)

	seq, err := ss.GetNextSequence(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, types.Nonce(4), seq) // LoadNextSequenceMap increases local nonce by 1
}

func TestSequenceSyncer_SyncSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	ss := NewTestSequenceSyncer(t, txStore, client)

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")
	enabledAddresses := []common.Address{addr}

	t.Run("syncs sequence successfully", func(t *testing.T) {
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(2), nil).Once()
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		var chStop services.StopChan
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(3), nil).Once()
		ss.SyncSequence(ctx, addr, chStop)

		seq, err := ss.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(3), seq)
	})

	t.Run("retries if on-chain syncing fails", func(t *testing.T) {
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(2), nil).Once()
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		var chStop services.StopChan
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(3), nil).Once()
		ss.SyncSequence(ctx, addr, chStop)

		seq, err := ss.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(3), seq)
	})
}

func TestSequenceSyncer_GetNextSequence(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	chainID := big.NewInt(0)
	txStore := txstoremock.NewEvmTxStore(t)

	client := clientmock.NewClient(t)
	client.On("ConfiguredChainID").Return(chainID)

	ss := NewTestSequenceSyncer(t, txStore, client)

	addr := common.HexToAddress("0xd5e099c71b797516c10ed0f0d895f429c2781142")

	t.Run("fails to get sequence if address doesn't exist in map", func(t *testing.T) {
		_, err := ss.GetNextSequence(ctx, addr)
		require.Error(t, err)

	})

	t.Run("fails to get sequence if address doesn't exist in map and is disabled", func(t *testing.T) {
		_, err := ss.GetNextSequence(ctx, addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("address disabled: %s", addr.Hex()))
	})

	t.Run("fails to get sequence if address is enabled, doesn't exist in map, and getSequenceForAddr fails", func(t *testing.T) {
		enabledAddresses := []common.Address{addr}
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(0), errors.New("no rows")).Twice()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Twice()
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		_, err := ss.GetNextSequence(ctx, addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("failed to find next sequence for address: %s", addr.Hex()))
	})

	t.Run("gets next sequence successfully if there is no entry in map but address is enabled and getSequenceForAddr is successful", func(t *testing.T) {
		enabledAddresses := []common.Address{addr}
		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(0), errors.New("no rows")).Once()
		client.On("PendingNonceAt", mock.Anything, addr).Return(uint64(0), errors.New("RPC unavailable")).Once()
		ss.LoadNextSequenceMap(ctx, enabledAddresses)

		txStore.On("FindLatestSequence", mock.Anything, addr, chainID).Return(types.Nonce(4), nil).Once()
		seq, err := ss.GetNextSequence(ctx, addr)
		require.NoError(t, err)
		require.Equal(t, types.Nonce(5), seq)

	})
}