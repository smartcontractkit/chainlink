package txmgr

import (
	"container/heap"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// TxPriorityQueue is a priority queue of transactions prioritized by creation time. The oldest transaction is at the front of the queue.
type TxPriorityQueue[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	*txmgrtypes.PriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
}

// NewTxPriorityQueue returns a new TxPriorityQueue instance
func NewTxPriorityQueue[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](maxUnstarted int) *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	pq := TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		PriorityQueue: txmgrtypes.NewPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](maxUnstarted),
	}

	return &pq
}

// AddTx adds a transaction to the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AddTx(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	heap.Push(pq, tx)
}

// RemoveNextTx removes the next transaction to be processed from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return heap.Pop(pq).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
}

// RemoveTxByID removes the transaction with the given ID from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveTxByID(id int64) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if i := pq.FindIndexByID(id); i != -1 {
		return heap.Remove(pq, i).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
	}

	return nil
}

// PruneByTxIDs removes the transactions with the given IDs from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneByTxIDs(ids []int64) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	removed := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, id := range ids {
		if tx := pq.RemoveTxByID(id); tx != nil {
			removed = append(removed, *tx)
		}
	}

	return removed
}

// PeekNextTx returns the next transaction to be processed without removing it from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return pq.Peek()
}

// Close clears the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	pq.PriorityQueue.Close()
}

// Cap returns the capacity of the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Cap() int {
	return pq.PriorityQueue.Cap()
}
