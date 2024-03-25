package queues

import (
	"container/heap"
	"sync"

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
	sync.RWMutex
	ph *priorityHeap[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
}

// NewTxPriorityQueue returns a new txPriorityQueue instance
func NewTxPriorityQueue[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](capacity int) *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	return &TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		ph: NewPriorityHeap[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](capacity),
	}
}

// AddTx adds a transaction to the queue.
// If the queue is full, the oldest transaction is removed.
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AddTx(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	pq.Lock()

	if pq.ph.Len() == pq.ph.Cap() {
		heap.Pop(pq.ph)
	}
	heap.Push(pq.ph, tx)

	pq.Unlock()
}

// RemoveNextTx removes the next transaction to be processed from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	if pq.ph.Len() == 0 {
		return nil
	}

	return heap.Pop(pq.ph).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
}

// RemoveTxByID removes the transaction with the given ID from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveTxByID(id int64) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	if pq.ph.Len() == 0 {
		return nil
	}

	return pq._removeTxByID(id)
}

// PruneByTxIDs removes the transactions with the given IDs from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneByTxIDs(ids []int64) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	if pq.ph.Len() == 0 {
		return nil
	}

	removed := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, id := range ids {
		if tx := pq._removeTxByID(id); tx != nil {
			removed = append(removed, *tx)
		}
	}

	return removed
}

// PeekNextTx returns the next transaction to be processed without removing it from the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.RLock()
	defer pq.RUnlock()

	if pq.ph.Len() == 0 {
		return nil
	}
	return pq.ph.Peek()
}

// Close clears the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	pq.Lock()
	defer pq.Unlock()

	pq.ph.Close()
}

// Cap returns the capacity of the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Cap() int {
	pq.RLock()
	defer pq.RUnlock()

	return pq.ph.Cap()
}

// Len returns the length of the queue
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Len() int {
	pq.RLock()
	defer pq.RUnlock()

	return pq.ph.Len()
}

// _removeTxByID removes the transaction with the given ID from the queue.
// This method assumes that the caller holds the lock.
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _removeTxByID(id int64) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if i, ok := pq.ph.FindIndexByID(id); ok {
		return heap.Remove(pq.ph, i).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
	}

	return nil
}
