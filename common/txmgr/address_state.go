package txmgr

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/internal/queues"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// addressState is the state of all transactions for a given address.
// It holds information about all transactions for a given address, including unstarted, in-progress, unconfirmed, confirmed, and fatal errored transactions.
// It is designed to help transition transactions between states and to provide information about the current state of transactions for a given address.
type addressState[
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
	idempotencyKeyToTx     map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	attemptHashToTxAttempt map[TX_HASH]*txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	unstartedTxs           *queues.TxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
	inprogressTx           *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	// NOTE: below each map's key is the transaction ID that is assigned via the persistent datastore
	unconfirmedTxs             map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmedMissingReceiptTxs map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	confirmedTxs               map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	allTxs                     map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	fatalErroredTxs            map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
}

// newAddressState returns a new addressState instance with initialized transaction state
func newAddressState[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	lggr logger.SugaredLogger,
	chainID CHAIN_ID,
	fromAddress ADDR,
	maxUnstarted uint64,
	txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE] {
	if maxUnstarted == 0 {
		panic("new_address_state: MaxUnstarted queue size must be greater than 0")
	}

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
		if len(tx.TxAttempts) > 0 {
			txAttemptCount += len(tx.TxAttempts)
		}
	}

	as := addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		lggr:        lggr,
		chainID:     chainID,
		fromAddress: fromAddress,

		idempotencyKeyToTx:         make(map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], idempotencyKeysCount),
		unstartedTxs:               queues.NewTxPriorityQueue[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](int(maxUnstarted)),
		inprogressTx:               nil,
		unconfirmedTxs:             make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxUnconfirmed]),
		confirmedMissingReceiptTxs: make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxConfirmedMissingReceipt]),
		confirmedTxs:               make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxConfirmed]),
		allTxs:                     make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs)),
		fatalErroredTxs:            make(map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], counts[TxFatalError]),
		attemptHashToTxAttempt:     make(map[TX_HASH]*txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], txAttemptCount),
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
		default:
			panic(fmt.Sprintf("unknown transaction state: %q", tx.State))
		}
		as.allTxs[tx.ID] = &tx
		if tx.IdempotencyKey != nil {
			as.idempotencyKeyToTx[*tx.IdempotencyKey] = &tx
		}
		for i := 0; i < len(tx.TxAttempts); i++ {
			txAttempt := tx.TxAttempts[i]
			as.attemptHashToTxAttempt[txAttempt.Hash] = &txAttempt
		}
	}

	return &as
}

// countTransactionsByState returns the number of transactions that are in the given state
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) countTransactionsByState(txState txmgrtypes.TxState) int {
	return 0
}

// findTxWithIdempotencyKey returns the transaction with the given idempotency key. If no transaction is found, nil is returned.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxWithIdempotencyKey(key string) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return nil
}

// applyToTxsByState calls the given function for each transaction in the given states.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions in the given states are considered.
// If no txStates are provided, all transactions are considered.
// Any transaction states that are unknown will cause a panic.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) applyToTxsByState(
	txStates []txmgrtypes.TxState,
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
	as.Lock()
	defer as.Unlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if len(txStates) == 0 {
		as._applyToTxs(as.allTxs, fn, txIDs...)
		return
	}

	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogressTx != nil {
				fn(as.inprogressTx)
			}
		case TxUnconfirmed:
			as._applyToTxs(as.unconfirmedTxs, fn, txIDs...)
		case TxConfirmedMissingReceipt:
			as._applyToTxs(as.confirmedMissingReceiptTxs, fn, txIDs...)
		case TxConfirmed:
			as._applyToTxs(as.confirmedTxs, fn, txIDs...)
		case TxFatalError:
			as._applyToTxs(as.fatalErroredTxs, fn, txIDs...)
		case TxUnstarted:
			nfn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
				if tx.State == TxUnstarted {
					fn(tx)
				}
			}
			as._applyToTxs(as.allTxs, nfn, txIDs...)
		default:
			panic(fmt.Sprintf("unknown transaction state: %q", txState))
		}
	}
}

// findTxAttempts returns all attempts for the given transactions that match the given filters.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions are considered.
// If no txStates are provided, all transactions are considered.
// The txFilter is applied to the transactions and the txAttemptFilter is applied to the attempts.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxAttempts(
	txStates []txmgrtypes.TxState,
	txFilter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txAttemptFilter func(*txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return nil
}

// findTxs returns all transactions that match the given filters.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions are considered.
// If no txStates are provided, all transactions are considered.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxs(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	as.RLock()
	defer as.RUnlock()

	// if txStates is empty then apply the filter to only the as.allTransactions map
	if len(txStates) == 0 {
		return as._findTxs(as.allTxs, filter, txIDs...)
	}

	var txs []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, txState := range txStates {
		switch txState {
		case TxInProgress:
			if as.inprogressTx != nil && filter(as.inprogressTx) {
				txs = append(txs, *as.inprogressTx)
			}
		case TxUnconfirmed:
			txs = append(txs, as._findTxs(as.unconfirmedTxs, filter, txIDs...)...)
		case TxConfirmedMissingReceipt:
			txs = append(txs, as._findTxs(as.confirmedMissingReceiptTxs, filter, txIDs...)...)
		case TxConfirmed:
			txs = append(txs, as._findTxs(as.confirmedTxs, filter, txIDs...)...)
		case TxFatalError:
			txs = append(txs, as._findTxs(as.fatalErroredTxs, filter, txIDs...)...)
		case TxUnstarted:
			fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
				return tx.State == TxUnstarted && filter(tx)
			}
			txs = append(txs, as._findTxs(as.allTxs, fn, txIDs...)...)
		default:
			panic(fmt.Sprintf("unknown transaction state: %q", txState))
		}
	}

	return txs
}

// pruneUnstartedTxQueue removes the transactions with the given IDs from the unstarted transaction queue.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) pruneUnstartedTxQueue(ids []int64) {
}

// reapConfirmedTxs removes confirmed transactions that are older than the given time threshold and have receipts older than the given block number threshold.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) reapConfirmedTxs(minBlockNumberToKeep int64, timeThreshold time.Time) {
	as.Lock()
	defer as.Unlock()

	for _, tx := range as.confirmedTxs {
		if len(tx.TxAttempts) == 0 {
			continue
		}
		if tx.CreatedAt.After(timeThreshold) {
			continue
		}

		for i := 0; i < len(tx.TxAttempts); i++ {
			if len(tx.TxAttempts[i].Receipts) == 0 {
				continue
			}
			if tx.TxAttempts[i].Receipts[0].GetBlockNumber() == nil || tx.TxAttempts[i].Receipts[0].GetBlockNumber().Int64() >= minBlockNumberToKeep {
				continue
			}
			as._deleteTx(tx.ID)
		}
	}
}

// reapFatalErroredTxs removes fatal errored transactions that are older than the given time threshold.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) reapFatalErroredTxs(timeThreshold time.Time) {
	as.Lock()
	defer as.Unlock()

	for _, tx := range as.fatalErroredTxs {
		if tx.CreatedAt.After(timeThreshold) {
			continue
		}
		as._deleteTx(tx.ID)
	}
}

// deleteTxs removes the transactions with the given IDs from the address state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deleteTxs(txIDs ...int64) {
	as.Lock()
	defer as.Unlock()

	as._deleteTxs(txIDs...)
}

// deleteTxAttempts removes the attempts with the given IDs from the address state.
// It removes the attempts from the hash lookup map and from the transaction.
// If an attempt is not found in the hash lookup map, it is ignored.
// If a transaction is not found in the allTxs map, it is ignored.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deleteTxAttempts(txAttempts ...txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
	as.Lock()
	defer as.Unlock()

	for _, txAttempt := range txAttempts {
		// remove the attempt from the hash lookup map
		delete(as.attemptHashToTxAttempt, txAttempt.Hash)
		// remove the attempt from the transaction
		if tx := as.allTxs[txAttempt.TxID]; tx != nil {
			var removeIndex int
			for i := 0; i < len(tx.TxAttempts); i++ {
				if tx.TxAttempts[i].ID == txAttempt.ID {
					removeIndex = i
					break
				}
			}
			tx.TxAttempts = append(tx.TxAttempts[:removeIndex], tx.TxAttempts[removeIndex+1:]...)
		}
	}
}

// peekNextUnstartedTx returns the next unstarted transaction in the queue without removing it from the unstarted queue.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekNextUnstartedTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// peekInProgressTx returns the in-progress transaction without removing it from the in-progress state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekInProgressTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// addTxToUnstarted adds the given transaction to the unstarted queue.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) addTxToUnstarted(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

// moveUnstartedToInProgress moves the next unstarted transaction to the in-progress state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnstartedToInProgress(
	etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// moveConfirmedMissingReceiptToUnconfirmed moves the confirmed missing receipt transaction to the unconfirmed state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedMissingReceiptToUnconfirmed(
	txID int64,
) error {
	return nil
}

// moveInProgressToUnconfirmed moves the in-progress transaction to the unconfirmed state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToUnconfirmed(
	etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// moveUnconfirmedToConfirmed moves the unconfirmed transaction to the confirmed state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmed(
	receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
) error {
	return nil
}

// moveTxToFatalError moves a transaction to the fatal error state.
// It returns an error if there is no transaction with the given ID.
// It returns an error if the transaction is not in an expected state.
// Unknown transaction states will cause a panic this includes Unconfirmed and Confirmed transactions.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveTxToFatalError(
	txID int64, txError null.String,
) error {
	as.Lock()
	defer as.Unlock()

	tx := as.allTxs[txID]
	if tx == nil {
		return fmt.Errorf("move_tx_to_fatal_error: no transaction with ID %d", txID)
	}
	originalState := tx.State

	// Move the transaction to the fatal error state
	as._moveTxToFatalError(tx, txError)

	// Remove the transaction from its original state
	switch originalState {
	case TxUnstarted:
		_ = as.unstartedTxs.RemoveTxByID(txID)
	case TxInProgress:
		as.inprogressTx = nil
	case TxConfirmedMissingReceipt:
		delete(as.confirmedMissingReceiptTxs, tx.ID)
	case TxFatalError:
		// Already in fatal error state
		return nil
	default:
		panic(fmt.Sprintf("unknown transaction state: %q", tx.State))
	}

	return nil
}

// moveUnconfirmedToConfirmedMissingReceipt moves the unconfirmed transaction to the confirmed missing receipt state.
// If there is no unconfirmed transaction with the given ID, an error is returned.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmedMissingReceipt(txID int64) error {
	as.Lock()
	defer as.Unlock()

	tx, ok := as.unconfirmedTxs[txID]
	if !ok || tx == nil {
		return fmt.Errorf("move_unconfirmed_to_confirmed_missing_receipt: no unconfirmed transaction with ID %d", txID)
	}
	tx.State = TxConfirmedMissingReceipt

	as.confirmedMissingReceiptTxs[tx.ID] = tx
	delete(as.unconfirmedTxs, tx.ID)

	return nil
}

// moveInProgressToConfirmedMissingReceipt moves the in-progress transaction to the confirmed missing receipt state.
// If there is no in-progress transaction, an error is returned.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToConfirmedMissingReceipt(txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	as.Lock()
	defer as.Unlock()

	tx := as.inprogressTx
	if tx == nil {
		return fmt.Errorf("move_in_progress_to_confirmed_missing_receipt: no transaction in progress")
	}
	if len(tx.TxAttempts) == 0 {
		return fmt.Errorf("move_in_progress_to_confirmed_missing_receipt: no attempts for transaction with ID %d", tx.ID)
	}
	if tx.BroadcastAt.Before(broadcastAt) {
		tx.BroadcastAt = &broadcastAt
	}
	tx.State = TxConfirmedMissingReceipt
	txAttempt.State = txmgrtypes.TxAttemptBroadcast
	tx.TxAttempts = append(tx.TxAttempts, txAttempt)

	as.confirmedMissingReceiptTxs[tx.ID] = tx
	as.inprogressTx = nil

	return nil
}

// moveConfirmedToUnconfirmed moves the confirmed transaction to the unconfirmed state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedToUnconfirmed(txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	as.Lock()
	defer as.Unlock()

	if txAttempt.State != txmgrtypes.TxAttemptBroadcast {
		return fmt.Errorf("attempt must be in broadcast state")
	}

	tx, ok := as.confirmedTxs[txAttempt.TxID]
	if !ok || tx == nil {
		return fmt.Errorf("no confirmed transaction with ID %d", txAttempt.TxID)
	}
	if len(tx.TxAttempts) == 0 {
		return fmt.Errorf("no attempts for transaction with ID %d", txAttempt.TxID)
	}
	tx.State = TxUnconfirmed

	// Delete the receipt from the attempt
	for i := 0; i < len(tx.TxAttempts); i++ {
		if tx.TxAttempts[i].ID == txAttempt.ID {
			tx.TxAttempts[i].Receipts = nil
			tx.TxAttempts[i].State = txmgrtypes.TxAttemptInProgress
			tx.TxAttempts[i].BroadcastBeforeBlockNum = nil
			break
		}
	}

	as.unconfirmedTxs[tx.ID] = tx
	delete(as.confirmedTxs, tx.ID)

	return nil
}

func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _applyToTxs(
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

func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _findTxs(
	txIDsToTx map[int64]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
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

func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _deleteTxs(txIDs ...int64) {
	for _, txID := range txIDs {
		as._deleteTx(txID)
	}
}

func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _deleteTx(txID int64) {
	tx, ok := as.allTxs[txID]
	if !ok {
		return
	}

	for i := 0; i < len(tx.TxAttempts); i++ {
		txAttemptHash := tx.TxAttempts[i].Hash
		delete(as.attemptHashToTxAttempt, txAttemptHash)
	}
	if tx.IdempotencyKey != nil {
		delete(as.idempotencyKeyToTx, *tx.IdempotencyKey)
	}
	if as.inprogressTx != nil && as.inprogressTx.ID == txID {
		as.inprogressTx = nil
	}
	as.unstartedTxs.RemoveTxByID(txID)
	delete(as.unconfirmedTxs, txID)
	delete(as.confirmedMissingReceiptTxs, txID)
	delete(as.confirmedTxs, txID)
	delete(as.fatalErroredTxs, txID)
	delete(as.allTxs, txID)
}

func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) _moveTxToFatalError(
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txError null.String,
) {
	tx.State = TxFatalError
	tx.Sequence = nil
	tx.BroadcastAt = nil
	tx.InitialBroadcastAt = nil
	tx.Error = txError
	as.fatalErroredTxs[tx.ID] = tx
}
