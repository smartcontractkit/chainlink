package queues_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/common/internal/queues"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func TestTxPriorityQueue(t *testing.T) {
	capacity := 5
	t.Run("transactions can be added to queue", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(100, 0)},
			{ID: 1, CreatedAt: time.Unix(200, 0)},
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		require.Equal(t, len(txs), pq.Len())

		assert.Equal(t, txs[0].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[0].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("transactions get ordered by createdAt", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(500, 0)}, // 4
			{ID: 1, CreatedAt: time.Unix(300, 0)}, // 2
			{ID: 2, CreatedAt: time.Unix(100, 0)}, // 0
			{ID: 3, CreatedAt: time.Unix(200, 0)}, // 1
			{ID: 4, CreatedAt: time.Unix(400, 0)}, // 3
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		require.Equal(t, len(txs), pq.Len())

		assert.Equal(t, txs[2].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[2].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[4].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[4].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[0].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[0].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("transactions can be added to full queue and keep capacity limits", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(100, 0)}, // dropped
			{ID: 1, CreatedAt: time.Unix(200, 0)}, // dropped
			{ID: 2, CreatedAt: time.Unix(300, 0)},
			{ID: 3, CreatedAt: time.Unix(400, 0)},
			{ID: 4, CreatedAt: time.Unix(500, 0)},
			{ID: 5, CreatedAt: time.Unix(600, 0)},
			{ID: 6, CreatedAt: time.Unix(700, 0)},
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		require.Equal(t, capacity, pq.Len())

		assert.Equal(t, txs[2].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[2].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[4].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[4].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[5].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[5].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[6].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[6].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("remove oldest transactions first if over capacity", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 3, CreatedAt: time.Unix(400, 0)},
			{ID: 2, CreatedAt: time.Unix(300, 0)}, // oldest
			{ID: 4, CreatedAt: time.Unix(500, 0)},
			{ID: 6, CreatedAt: time.Unix(700, 0)},
			{ID: 0, CreatedAt: time.Unix(100, 0)}, // Dropped
			{ID: 1, CreatedAt: time.Unix(200, 0)}, // Dropped
			{ID: 5, CreatedAt: time.Unix(600, 0)},
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		assert.Equal(t, capacity, pq.Len())

		assert.Equal(t, txs[1].ID, pq.PeekNextTx().ID)
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("access oldest transactions when using peek and RemoveNextTx", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(400, 0)},
			{ID: 1, CreatedAt: time.Unix(300, 0)},
			{ID: 2, CreatedAt: time.Unix(100, 0)}, // oldest
			{ID: 3, CreatedAt: time.Unix(200, 0)},
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		assert.Equal(t, len(txs), pq.Len())
		assert.Equal(t, txs[2].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[2].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[0].ID, pq.PeekNextTx().ID)
		assert.Equal(t, txs[0].ID, pq.RemoveNextTx().ID)
		assert.Nil(t, pq.PeekNextTx())
		assert.Nil(t, pq.RemoveNextTx())
		assert.Equal(t, 0, pq.Len())
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("transactions can be removed by using RemoveTxByID", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(100, 0)},
			{ID: 1, CreatedAt: time.Unix(200, 0)},
			{ID: 2, CreatedAt: time.Unix(300, 0)}, // should be removed
			{ID: 3, CreatedAt: time.Unix(400, 0)},
			{ID: 4, CreatedAt: time.Unix(500, 0)},
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		require.Equal(t, capacity, pq.Len())

		txIDToRemove := int64(2)
		removedTx := pq.RemoveTxByID(txIDToRemove)
		require.NotNil(t, removedTx)
		require.Equal(t, txIDToRemove, removedTx.ID)

		require.Equal(t, len(txs)-1, pq.Len())
		assert.Equal(t, txs[0].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[1].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[4].ID, pq.RemoveNextTx().ID)
		assert.Nil(t, pq.PeekNextTx())
		assert.Equal(t, capacity, pq.Cap())
	})
	t.Run("transactions can be removed by using PruneByTxIDs", func(t *testing.T) {
		pq := queues.NewTxPriorityQueue[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, evmgas.EvmFee,
		](capacity)
		require.Equal(t, capacity, pq.Cap())
		defer func(t *testing.T) {
			pq.Close()
			assert.Equal(t, 0, pq.Len())
			assert.Equal(t, 0, pq.Cap())
		}(t)

		txs := []txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, evmgas.EvmFee]{
			{ID: 0, CreatedAt: time.Unix(100, 0)}, // should be pruned
			{ID: 1, CreatedAt: time.Unix(200, 0)},
			{ID: 2, CreatedAt: time.Unix(300, 0)}, // should be pruned
			{ID: 3, CreatedAt: time.Unix(400, 0)},
			{ID: 4, CreatedAt: time.Unix(500, 0)}, // should be pruned
		}
		for i := 0; i < len(txs); i++ {
			pq.AddTx(&txs[i])
		}
		require.Equal(t, capacity, pq.Len())

		txIDsToBePruned := []int64{0, 2, 4}
		removed := pq.PruneByTxIDs(txIDsToBePruned)
		require.Equal(t, len(txIDsToBePruned), len(removed))
		assert.Equal(t, txs[0].ID, removed[0].ID)
		assert.Equal(t, txs[2].ID, removed[1].ID)
		assert.Equal(t, txs[4].ID, removed[2].ID)

		assert.Equal(t, txs[1].ID, pq.RemoveNextTx().ID)
		assert.Equal(t, txs[3].ID, pq.RemoveNextTx().ID)
		assert.Nil(t, pq.PeekNextTx())
		assert.Equal(t, capacity, pq.Cap())
	})
}
