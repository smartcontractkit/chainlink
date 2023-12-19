package txmgr

import (
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
	txStore     PersistentTxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	sync.RWMutex
	idempotencyKeyToTx map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	unstarted          *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	inprogress         *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// NOTE: currently the unconfirmed map's key is the transaction ID that is assigned via the postgres DB
	unconfirmed             map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmedMissingReceipt map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmed               map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	allTransactions         map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// NewAddressState returns a new AddressState instance
func NewAddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	chainID CHAIN_ID,
	fromAddress ADDR,
	maxUnstarted int,
	txStore PersistentTxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
) (*AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	as := AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		chainID:     chainID,
		fromAddress: fromAddress,
		txStore:     txStore,

		idempotencyKeyToTx:      map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		unstarted:               NewTxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](maxUnstarted),
		inprogress:              nil,
		unconfirmed:             map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		confirmedMissingReceipt: map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		confirmed:               map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		allTransactions:         map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
	}

	as.Lock()
	defer as.Unlock()

	// Load all unstarted transactions from persistent storage
	offset := 0
	limit := 50
	for {
		txs, count, err := txStore.UnstartedTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return nil, fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.unstarted.AddTx(&tx)
			as.allTransactions[tx.ID] = &tx
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
		return nil, fmt.Errorf("address_state: initialization: %w", err)
	}
	as.inprogress = tx
	if tx != nil {
		if tx.IdempotencyKey != nil {
			as.idempotencyKeyToTx[*tx.IdempotencyKey] = tx
		}
		as.allTransactions[tx.ID] = tx
	}

	// Load all unconfirmed transactions from persistent storage
	offset = 0
	limit = 50
	for {
		txs, count, err := txStore.UnconfirmedTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return nil, fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.unconfirmed[tx.ID] = &tx
			as.allTransactions[tx.ID] = &tx
			if tx.IdempotencyKey != nil {
				as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
			}
		}
		if count <= offset+limit {
			break
		}
		offset += limit
	}

	// Load all confirmed transactions from persistent storage
	offset = 0
	limit = 50
	for {
		txs, count, err := txStore.ConfirmedTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return nil, fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.confirmed[tx.ID] = &tx
			as.allTransactions[tx.ID] = &tx
			if tx.IdempotencyKey != nil {
				as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
			}
		}
		if count <= offset+limit {
			break
		}
		offset += limit
	}

	// Load all unconfirmed transactions from persistent storage
	offset = 0
	limit = 50
	for {
		txs, count, err := txStore.ConfirmedMissingReceiptTransactions(offset, limit, as.fromAddress, as.chainID)
		if err != nil {
			return nil, fmt.Errorf("address_state: initialization: %w", err)
		}
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			as.confirmedMissingReceipt[tx.ID] = &tx
			as.allTransactions[tx.ID] = &tx
			if tx.IdempotencyKey != nil {
				as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
			}
		}
		if count <= offset+limit {
			break
		}
		offset += limit
	}

	return &as, nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) close() {
	as.Lock()
	defer as.Unlock()

	as.unstarted.Close()
	as.unstarted = nil
	as.inprogress = nil
	clear(as.unconfirmed)
	clear(as.idempotencyKeyToTx)
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) unstartedCount() int {
	as.RLock()
	defer as.RUnlock()

	return as.unstarted.Len()
}
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) unconfirmedCount() int {
	as.RLock()
	defer as.RUnlock()

	return len(as.unconfirmed)
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxWithIdempotencyKey(key string) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	return as.idempotencyKeyToTx[key]
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) confirmedMissingReceiptTxs() []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	var txAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, tx := range as.confirmedMissingReceipt {
		txAttempts = append(txAttempts, tx.TxAttempts...)
	}

	return txAttempts
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findLatestSequence() SEQ {
	as.RLock()
	defer as.RUnlock()

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

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ApplyToTxs(
	txStates []txmgrtypes.TxState,
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
	as.Lock()
	defer as.Unlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if txStates == nil || len(txStates) == 0 {
		as.applyToStorage(as.allTransactions, fn, txIDs...)
		return
	}

	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil {
				fn(as.inprogress)
			}
		case TxUnconfirmed:
			as.applyToStorage(as.unconfirmed, fn, txIDs...)
		case TxConfirmedMissingReceipt:
			as.applyToStorage(as.confirmedMissingReceipt, fn, txIDs...)
		case TxConfirmed:
			as.applyToStorage(as.confirmed, fn, txIDs...)
		}
	}
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) applyToStorage(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
	as.Lock()
	defer as.Unlock()

	// if txIDs is not empty then only apply the filter to those transactions
	if len(txIDs) > 0 {
		for _, txID := range txIDs {
			tx := txIDsToTx[txID]
			if tx != nil {
				fn(tx)
			}
		}
		return
	}

	// if txIDs is empty then apply the filter to all transactions
	for _, tx := range txIDsToTx {
		fn(tx)
	}
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FetchTxAttempts(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if txStates == nil || len(txStates) == 0 {
		return as.fetchTxAttemptsFromStorage(as.allTransactions, filter, txIDs...)
	}

	var txAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil && filter(as.inprogress) {
				txAttempts = append(txAttempts, as.inprogress.TxAttempts...)
			}
		case TxUnconfirmed:
			txAttempts = append(txAttempts, as.fetchTxAttemptsFromStorage(as.unconfirmed, filter, txIDs...)...)
		case TxConfirmedMissingReceipt:
			txAttempts = append(txAttempts, as.fetchTxAttemptsFromStorage(as.confirmedMissingReceipt, filter, txIDs...)...)
		case TxConfirmed:
			txAttempts = append(txAttempts, as.fetchTxAttemptsFromStorage(as.confirmed, filter, txIDs...)...)
		}
	}

	return txAttempts
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) fetchTxAttemptsFromStorage(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	var txAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// if txIDs is not empty then only apply the filter to those transactions
	if len(txIDs) > 0 {
		for _, txID := range txIDs {
			tx := txIDsToTx[txID]
			if tx != nil && filter(tx) {
				txAttempts = append(txAttempts, tx.TxAttempts...)
			}
		}
		return txAttempts
	}

	// if txIDs is empty then apply the filter to all transactions
	for _, tx := range txIDsToTx {
		if filter(tx) {
			txAttempts = append(txAttempts, tx.TxAttempts...)
		}
	}

	return txAttempts
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FetchTxs(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if txStates == nil || len(txStates) == 0 {
		return as.fetchTxsFromStorage(as.allTransactions, filter, txIDs...)
	}

	var txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil && filter(as.inprogress) {
				txs = append(txs, *as.inprogress)
			}
		case TxUnconfirmed:
			txs = append(txs, as.fetchTxsFromStorage(as.unconfirmed, filter, txIDs...)...)
		case TxConfirmedMissingReceipt:
			txs = append(txs, as.fetchTxsFromStorage(as.confirmedMissingReceipt, filter, txIDs...)...)
		case TxConfirmed:
			txs = append(txs, as.fetchTxsFromStorage(as.confirmed, filter, txIDs...)...)
		}
	}

	return txs
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) fetchTxsFromStorage(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	var txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// if txIDs is not empty then only apply the filter to those transactions
	if len(txIDs) > 0 {
		for _, txID := range txIDs {
			tx := txIDsToTx[txID]
			if tx != nil && filter(tx) {
				txs = append(txs, *tx)
			}
		}
		return txs
	}

	// if txIDs is empty then apply the filter to all transactions
	for _, tx := range txIDsToTx {
		if filter(tx) {
			txs = append(txs, *tx)
		}
	}

	return txs
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneUnstartedTxQueue(queueSize uint32, filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool) int {
	as.RLock()
	defer as.RUnlock()

	return len(as.unstarted.Prune(int(queueSize), filter))
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekNextUnstartedTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	as.RLock()
	defer as.RUnlock()

	tx := as.unstarted.PeekNextTx()
	if tx == nil {
		return nil, fmt.Errorf("peek_next_unstarted_tx: %w (address: %s)", ErrTxnNotFound, as.fromAddress)
	}

	return tx, nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekInProgressTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	as.RLock()
	defer as.RUnlock()

	tx := as.inprogress
	if tx == nil {
		return nil, fmt.Errorf("peek_in_progress_tx: %w (address: %s)", ErrTxnNotFound, as.fromAddress)
	}

	return tx, nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) addTxToUnstarted(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as.Lock()
	defer as.Unlock()

	if as.unstarted.Len() >= as.unstarted.Cap() {
		return fmt.Errorf("move_tx_to_unstarted: address %s unstarted queue capactiry has been reached", as.fromAddress)
	}

	as.unstarted.AddTx(tx)
	as.allTransactions[tx.ID] = tx
	if tx.IdempotencyKey != nil {
		as.idempotencyKeyToTx[*tx.IdempotencyKey] = tx
	}

	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnstartedToInProgress(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as.Lock()
	defer as.Unlock()

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
	as.Lock()
	defer as.Unlock()

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
			break
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
	as.Lock()
	defer as.Unlock()

	for as.unstarted.Len() > 0 {
		tx := as.unstarted.RemoveNextTx()
		abandon[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](tx)
	}

	if as.inprogress != nil {
		tx := as.inprogress
		abandon[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](tx)
		as.inprogress = nil
	}
	for _, tx := range as.unconfirmed {
		abandon[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](tx)
	}
	for _, tx := range as.idempotencyKeyToTx {
		abandon[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](tx)
	}

	clear(as.unconfirmed)
}

func abandon[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	if tx == nil {
		return
	}

	tx.State = TxFatalError
	tx.Sequence = nil
	tx.Error = null.NewString("abandoned", true)
}
