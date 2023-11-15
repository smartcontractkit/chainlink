package txmgr

import (
	"container/heap"
	"context"
	"fmt"
	"sync"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"gopkg.in/guregu/null.v4"
)

// AddressState is the state of a given from address
type AddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	chainID     CHAIN_ID
	fromAddress ADDR

	lock               sync.RWMutex
	idempotencyKeyToTx map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	unstarted          *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	inprogress         *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// TODO(jtw): using the TX ID as the key for the map might not make sense since the ID is set by the
	// postgres DB which creates a dependency on the postgres DB. We should consider creating a UUID or ULID
	// TX ID -> TX
	unconfirmed map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// NewAddressState returns a new AddressState instance
func NewAddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](chainID CHAIN_ID, fromAddress ADDR, maxUnstarted int) *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	as := AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		chainID:            chainID,
		fromAddress:        fromAddress,
		idempotencyKeyToTx: map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		unstarted:          NewTxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](maxUnstarted),
		inprogress:         nil,
		unconfirmed:        map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
	}

	return &as
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Initialize(txStore PersistentTxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error {
	as.lock.Lock()
	defer as.lock.Unlock()

	// Load all unstarted transactions from persistent storage
	offset := 0
	limit := 50
	for {
		txs, count, err := txStore.UnstartedTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.unstarted.AddTx(&tx)
			if tx.IdempotencyKey != nil {
				as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
			}
		}
		if count <= offset+limit {
			break
		}
		offset += limit
	}

	// Load all in progress transactions from persistent storage
	ctx := context.Background()
	tx, err := txStore.GetTxInProgress(ctx, as.fromAddress)
	if err != nil {
		return fmt.Errorf("address_state: initialization: %w", err)
	}
	as.inprogress = tx
	if tx.IdempotencyKey != nil {
		as.idempotencyKeyToTx[*tx.IdempotencyKey] = tx
	}

	// Load all unconfirmed transactions from persistent storage
	offset = 0
	limit = 50
	for {
		txs, count, err := txStore.UnconfirmedTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.unconfirmed[tx.ID] = &tx
			if tx.IdempotencyKey != nil {
				as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
			}
		}
		if count <= offset+limit {
			break
		}
		offset += limit
	}

	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) close() {
	as.lock.Lock()
	defer as.lock.Unlock()

	as.unstarted.Close()
	as.unstarted = nil
	as.inprogress = nil
	clear(as.unconfirmed)
	clear(as.idempotencyKeyToTx)
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) unstartedCount() int {
	as.lock.RLock()
	defer as.lock.RUnlock()

	return as.unstarted.Len()
}
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) unconfirmedCount() int {
	as.lock.RLock()
	defer as.lock.RUnlock()

	return len(as.unconfirmed)
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxWithIdempotencyKey(key string) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.lock.RLock()
	defer as.lock.RUnlock()

	return as.idempotencyKeyToTx[key]
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findLatestSequence() SEQ {
	as.lock.RLock()
	defer as.lock.RUnlock()

	var maxSeq SEQ
	if as.inprogress != nil && as.inprogress.Sequence != nil {
		if (*as.inprogress.Sequence).Int64() > maxSeq.Int64() {
			maxSeq = *as.inprogress.Sequence
		}
	}
	for _, tx := range as.unconfirmed {
		if tx.Sequence == nil {
			continue
		}
		if (*tx.Sequence).Int64() > maxSeq.Int64() {
			maxSeq = *tx.Sequence
		}
	}

	return maxSeq
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekNextUnstartedTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	as.lock.RLock()
	defer as.lock.RUnlock()

	tx := as.unstarted.PeekNextTx()
	if tx == nil {
		return nil, fmt.Errorf("peek_next_unstarted_tx: %w (address: %s)", ErrTxnNotFound, as.fromAddress)
	}

	return tx, nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekInProgressTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	as.lock.RLock()
	defer as.lock.RUnlock()

	tx := as.inprogress
	if tx == nil {
		return nil, fmt.Errorf("peek_in_progress_tx: %w (address: %s)", ErrTxnNotFound, as.fromAddress)
	}

	return tx, nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) addTxToUnstarted(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as.lock.Lock()
	defer as.lock.Unlock()

	if as.unstarted.Len() >= as.unstarted.Cap() {
		return fmt.Errorf("move_tx_to_unstarted: address %s unstarted queue capactiry has been reached", as.fromAddress)
	}

	as.unstarted.AddTx(tx)
	if tx.IdempotencyKey != nil {
		as.idempotencyKeyToTx[*tx.IdempotencyKey] = tx
	}

	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnstartedToInProgress(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as.lock.Lock()
	defer as.lock.Unlock()

	if as.inprogress != nil {
		return fmt.Errorf("move_unstarted_to_in_progress: address %s already has a transaction in progress", as.fromAddress)
	}

	if tx != nil {
		// if tx is not nil then remove the tx from the unstarted queue
		// TODO(jtw): what should be the unique idenitifier for each transaction? ID is being set by the postgres DB
		tx = as.unstarted.RemoveTxByID(tx.ID)
	} else {
		// if tx is nil then pop the next unstarted transaction
		tx = as.unstarted.RemoveNextTx()
	}
	if tx == nil {
		return fmt.Errorf("move_unstarted_to_in_progress: no unstarted transaction to move to in_progress")
	}
	tx.State = TxInProgress
	as.inprogress = tx

	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToUnconfirmed(
	txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	as.lock.Lock()
	defer as.lock.Unlock()

	tx := as.inprogress
	if tx == nil {
		return fmt.Errorf("move_in_progress_to_unconfirmed: no transaction in progress")
	}
	tx.State = TxUnconfirmed

	var found bool
	for i := 0; i < len(tx.TxAttempts); i++ {
		if tx.TxAttempts[i].ID == txAttempt.ID {
			tx.TxAttempts[i] = txAttempt
			found = true
		}
	}
	if !found {
		// NOTE(jtw): this would mean that the TxAttempt did not exist for the Tx
		// NOTE(jtw): should this log a warning?
		// NOTE(jtw): can this happen?
		tx.TxAttempts = append(tx.TxAttempts, txAttempt)
	}

	as.unconfirmed[tx.ID] = tx
	as.inprogress = nil

	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) abandon() {
	as.lock.Lock()
	defer as.lock.Unlock()

	for as.unstarted.Len() > 0 {
		tx := as.unstarted.RemoveNextTx()
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}

	if as.inprogress != nil {
		as.inprogress.State = TxFatalError
		as.inprogress.Sequence = nil
		as.inprogress.Error = null.NewString("abandoned", true)
		as.inprogress = nil
	}
	for _, tx := range as.unconfirmed {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}
	for _, tx := range as.idempotencyKeyToTx {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}

	clear(as.unconfirmed)
}

// TxPriorityQueue is a priority queue of transactions prioritized by creation time. The oldest transaction is at the front of the queue.
type TxPriorityQueue[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	sync.Mutex
	txs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
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
		txs: make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 0, maxUnstarted),
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
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Push(tx any) {
	pq.txs = append(pq.txs, tx.(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]))
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Pop() any {
	pq.Lock()
	defer pq.Unlock()

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

	for i, tx := range pq.txs {
		if tx.ID == id {
			return heap.Remove(pq, i).(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE])
		}
	}

	return nil
}
func (pq *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekNextTx() *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
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
