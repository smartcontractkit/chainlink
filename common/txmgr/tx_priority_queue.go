package txmgr

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
	txs       []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	idToIndex map[int64]int
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
		txs:       make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0, maxUnstarted),
		idToIndex: make(map[int64]int),
	}

	return &pq
}

func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Cap() int {
	return cap(pq.txs)
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Len() int {
	return len(pq.txs)
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Less(i, j int) bool {
	// We want Pop to give us the oldest, not newest, transaction based on creation time
	return pq.txs[i].CreatedAt.Before(pq.txs[j].CreatedAt)
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Swap(i, j int) {
	pq.txs[i], pq.txs[j] = pq.txs[j], pq.txs[i]
	pq.idToIndex[pq.txs[i].ID] = j
	pq.idToIndex[pq.txs[j].ID] = i
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Push(tx any) {
	pq.txs = append(pq.txs, tx.(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]))
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Pop() any {
	old := pq.txs
	n := len(old)
	tx := old[n-1]
	old[n-1] = nil // avoid memory leak
	pq.txs = old[0 : n-1]
	return tx
}

func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AddTx(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	pq.Lock()
	defer pq.Unlock()

	heap.Push(pq, tx)
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	return heap.Pop(pq).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) RemoveTxByID(id int64) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	if i, ok := pq.idToIndex[id]; ok {
		return heap.Remove(pq, i).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
	}

	return nil
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	pq.Lock()
	defer pq.Unlock()

	if len(pq.txs) == 0 {
		return nil
	}
	return pq.txs[0]
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	pq.Lock()
	defer pq.Unlock()

	clear(pq.txs)
}
