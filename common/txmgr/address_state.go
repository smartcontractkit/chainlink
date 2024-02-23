package txmgr

import (
	"sync"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// AddressState is the state of all transactions for a given address
type AddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	lggr        logger.SugaredLogger
	chainID     CHAIN_ID
	fromAddress ADDR

	sync.RWMutex
	idempotencyKeyToTx map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	unstartedTxs       *TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	inprogressTx       *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// NOTE: currently the unconfirmed map's key is the transaction ID that is assigned via the postgres DB
	unconfirmedTxs             map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmedMissingReceiptTxs map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmedTxs               map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	allTxs                     map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	fatalErroredTxs            map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	attemptHashToTxAttempt     map[TX_HASH]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// NewAddressState returns a new AddressState instance with initialized transaction state
func NewAddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	lggr logger.SugaredLogger,
	chainID CHAIN_ID,
	fromAddress ADDR,
	maxUnstarted int,
	txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) (*AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	// Count the number of transactions in each state to reduce the number of map resizes
	counts := map[txmgrtypes.TxState]int{
		TxUnstarted:               0,
		TxInProgress:              0,
		TxUnconfirmed:             0,
		TxConfirmedMissingReceipt: 0,
		TxConfirmed:               0,
		TxFatalError:              0,
	}
	var idempotencyKeysCount int
	var txAttemptCount int
	for _, tx := range txs {
		counts[tx.State]++
		if tx.IdempotencyKey != nil {
			idempotencyKeysCount++
		}
		if tx.State == TxUnconfirmed {
			txAttemptCount += len(tx.TxAttempts)
		}
	}

	as := AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		lggr:        lggr,
		chainID:     chainID,
		fromAddress: fromAddress,

		idempotencyKeyToTx:         make(map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], idempotencyKeysCount),
		unstartedTxs:               NewTxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](maxUnstarted),
		inprogressTx:               nil,
		unconfirmedTxs:             make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxUnconfirmed]),
		confirmedMissingReceiptTxs: make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxConfirmedMissingReceipt]),
		confirmedTxs:               make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxConfirmed]),
		allTxs:                     make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs)),
		fatalErroredTxs:            make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxFatalError]),
		attemptHashToTxAttempt:     make(map[TX_HASH]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], txAttemptCount),
	}

	// Load all transactions supplied
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		switch tx.State {
		case TxUnstarted:
			as.unstartedTxs.AddTx(&tx)
		case TxInProgress:
			as.inprogressTx = &tx
		case TxUnconfirmed:
			as.unconfirmedTxs[tx.ID] = &tx
		case TxConfirmedMissingReceipt:
			as.confirmedMissingReceiptTxs[tx.ID] = &tx
		case TxConfirmed:
			as.confirmedTxs[tx.ID] = &tx
		case TxFatalError:
			as.fatalErroredTxs[tx.ID] = &tx
		}
		as.allTxs[tx.ID] = &tx
		if tx.IdempotencyKey != nil {
			as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
		}
		for _, txAttempt := range tx.TxAttempts {
			as.attemptHashToTxAttempt[txAttempt.Hash] = txAttempt
		}
	}

	return &as, nil
}

// CountTransactionsByState returns the number of transactions that are in the given state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountTransactionsByState(txState txmgrtypes.TxState) int {
	as.RLock()
	defer as.RUnlock()

	switch txState {
	case TxUnstarted:
		return as.unstarted.Len()
	case TxInProgress:
		if as.inprogress != nil {
			return 1
		}
	case TxUnconfirmed:
		return len(as.unconfirmed)
	case TxConfirmedMissingReceipt:
		return len(as.confirmedMissingReceipt)
	case TxConfirmed:
		return len(as.confirmed)
	case TxFatalError:
		return len(as.fatalErrored)
	}

	return 0
}

// FindTxWithIdempotencyKey returns the transaction with the given idempotency key.
// If no transaction is found, nil is returned.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithIdempotencyKey(key string) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	return as.idempotencyKeyToTx[key]
}

// ApplyToTxsByState calls the given function for each transaction in the given states.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions in the given states are considered.
// If no txStates are provided, all transactions are considered.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ApplyToTxsByState(
	txStates []txmgrtypes.TxState,
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
	as.Lock()
	defer as.Unlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if len(txStates) == 0 {
		as.applyToTxs(as.allTransactions, fn, txIDs...)
		return
	}

	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil {
				fn(as.inprogress)
			}
		case TxUnconfirmed:
			as.applyToTxs(as.unconfirmed, fn, txIDs...)
		case TxConfirmedMissingReceipt:
			as.applyToTxs(as.confirmedMissingReceipt, fn, txIDs...)
		case TxConfirmed:
			as.applyToTxs(as.confirmed, fn, txIDs...)
		case TxFatalError:
			as.applyToTxs(as.fatalErrored, fn, txIDs...)
		}
	}
}

// FindTxAttempts returns all attempts for the given transactions that match the given filters.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions are considered.
// If no txStates are provided, all transactions are considered.
// The txFilter is applied to the transactions and the txAttemptFilter is applied to the attempts.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttempts(
	txStates []txmgrtypes.TxState,
	txFilter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txAttemptFilter func(*txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if len(txStates) == 0 {
		return as.findTxAttempts(as.allTransactions, txFilter, txAttemptFilter, txIDs...)
	}

	var txAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil && txFilter(as.inprogress) {
				for i := 0; i < len(as.inprogress.TxAttempts); i++ {
					txAttempt := as.inprogress.TxAttempts[i]
					if txAttemptFilter(&txAttempt) {
						txAttempts = append(txAttempts, txAttempt)
					}
				}
			}
		case TxUnconfirmed:
			txAttempts = append(txAttempts, as.findTxAttempts(as.unconfirmed, txFilter, txAttemptFilter, txIDs...)...)
		case TxConfirmedMissingReceipt:
			txAttempts = append(txAttempts, as.findTxAttempts(as.confirmedMissingReceipt, txFilter, txAttemptFilter, txIDs...)...)
		case TxConfirmed:
			txAttempts = append(txAttempts, as.findTxAttempts(as.confirmed, txFilter, txAttemptFilter, txIDs...)...)
		case TxFatalError:
			txAttempts = append(txAttempts, as.findTxAttempts(as.fatalErrored, txFilter, txAttemptFilter, txIDs...)...)
		}
	}

	return txAttempts
}

// FindTxs returns all transactions that match the given filters.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions are considered.
// If no txStates are provided, all transactions are considered.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxs(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if len(txStates) == 0 {
		return as.findTxs(as.allTransactions, filter, txIDs...)
	}

	var txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogress != nil && filter(as.inprogress) {
				txs = append(txs, *as.inprogress)
			}
		case TxUnconfirmed:
			txs = append(txs, as.findTxs(as.unconfirmed, filter, txIDs...)...)
		case TxConfirmedMissingReceipt:
			txs = append(txs, as.findTxs(as.confirmedMissingReceipt, filter, txIDs...)...)
		case TxConfirmed:
			txs = append(txs, as.findTxs(as.confirmed, filter, txIDs...)...)
		case TxFatalError:
			txs = append(txs, as.findTxs(as.fatalErrored, filter, txIDs...)...)
		}
	}

	return txs
}

// PruneUnstartedTxQueue removes the transactions with the given IDs from the unstarted transaction queue.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneUnstartedTxQueue(ids []int64) {
}

// DeleteTxs removes the transactions with the given IDs from the address state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) DeleteTxs(txs ...txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
}

// PeekNextUnstartedTx returns the next unstarted transaction in the queue without removing it from the unstarted queue.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekNextUnstartedTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// PeekInProgressTx returns the in-progress transaction without removing it from the in-progress state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PeekInProgressTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// AddTxToUnstarted adds the given transaction to the unstarted queue.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) AddTxToUnstarted(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

// MoveUnstartedToInProgress moves the next unstarted transaction to the in-progress state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveUnstartedToInProgress(
	etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// MoveConfirmedMissingReceiptToUnconfirmed moves the confirmed missing receipt transaction to the unconfirmed state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveConfirmedMissingReceiptToUnconfirmed(
	txID int64,
) error {
	return nil
}

// MoveInProgressToUnconfirmed moves the in-progress transaction to the unconfirmed state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveInProgressToUnconfirmed(
	etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// MoveUnconfirmedToConfirmed moves the unconfirmed transaction to the confirmed state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveUnconfirmedToConfirmed(
	receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
) error {
	return nil
}

// MoveUnstartedToFatalError moves the unstarted transaction to the fatal error state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveUnstartedToFatalError(
	etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txError null.String,
) error {
	return nil
}

// MoveInProgressToFatalError moves the in-progress transaction to the fatal error state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveInProgressToFatalError(txError null.String) error {
	return nil
}

// MoveConfirmedMissingReceiptToFatalError moves the confirmed missing receipt transaction to the fatal error state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveConfirmedMissingReceiptToFatalError(
	etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txError null.String,
) error {
	return nil
}

// MoveUnconfirmedToConfirmedMissingReceipt moves the unconfirmed transaction to the confirmed missing receipt state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveUnconfirmedToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// MoveInProgressToConfirmedMissingReceipt moves the in-progress transaction to the confirmed missing receipt state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveInProgressToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// MoveConfirmedToUnconfirmed moves the confirmed transaction to the unconfirmed state.
func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MoveConfirmedToUnconfirmed(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) applyToTxs(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
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

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxAttempts(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txFilter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txAttemptFilter func(*txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	var txAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// if txIDs is not empty then only apply the filter to those transactions
	if len(txIDs) > 0 {
		for _, txID := range txIDs {
			tx := txIDsToTx[txID]
			if tx != nil && txFilter(tx) {
				for i := 0; i < len(tx.TxAttempts); i++ {
					txAttempt := tx.TxAttempts[i]
					if txAttemptFilter(&txAttempt) {
						txAttempts = append(txAttempts, txAttempt)
					}
				}
			}
		}
		return txAttempts
	}

	// if txIDs is empty then apply the filter to all transactions
	for _, tx := range txIDsToTx {
		if txFilter(tx) {
			for i := 0; i < len(tx.TxAttempts); i++ {
				txAttempt := tx.TxAttempts[i]
				if txAttemptFilter(&txAttempt) {
					txAttempts = append(txAttempts, txAttempt)
				}
			}
		}
	}

	return txAttempts
}

func (as *AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxs(
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
