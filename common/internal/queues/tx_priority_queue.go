package queues

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
	ph *priorityHeap[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
}

// NewTxPriorityQueue returns a new txPrioirityQueue instance
func NewTxPriorityQueue[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](capacity int) *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	pq := TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		ph: NewPriorityHeap[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](capacity),
	}

	return &pq
}

// AddTx adds a transaction to the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AddTx(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if pq.ph.Len() == pq.ph.Cap() {
		heap.Pop(pq.ph)
	}

	heap.Push(pq.ph, tx)
}

// RemoveNextTx removes the next transaction to be processed from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if pq.ph.Len() == 0 {
		return nil
	}

	return heap.Pop(pq.ph).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
}

// RemoveTxByID removes the transaction with the given ID from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveTxByID(id int64) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if pq.ph.Len() == 0 {
		return nil
	}

	if i := pq.ph.FindIndexByID(id); i != -1 {
		return heap.Remove(pq.ph, i-1).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
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
	return pq.ph.Peek()
}

// Close clears the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	pq.ph.Close()
}

// Cap returns the capacity of the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Cap() int {
	return pq.ph.Cap()
}

// Len returns the length of the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Len() int {
	return pq.ph.Len()
}
