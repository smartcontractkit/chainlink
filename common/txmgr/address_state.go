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
			panic("unknown transaction state")
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
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) applyToTxsByState(
	txStates []txmgrtypes.TxState,
	fn func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]),
	txIDs ...int64,
) {
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
	return nil
}

// pruneUnstartedTxQueue removes the transactions with the given IDs from the unstarted transaction queue.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) pruneUnstartedTxQueue(ids []int64) {
}

// deleteTxs removes the transactions with the given IDs from the address state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deleteTxs(txs ...txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
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
// The supplied txAttempt is added to the transaction.
// It returns an error if there is already a transaction in progress.
// It returns an error if there is no unstarted transaction to move to in_progress.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnstartedToInProgress(
	txID int64, seq SEQ, broadcastAt time.Time, initialBroadcastAt time.Time,
	txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	as.Lock()
	defer as.Unlock()

	if as.inprogressTx != nil {
		return fmt.Errorf("move_unstarted_to_in_progress: address %s already has a transaction in progress", as.fromAddress)
	}

	tx := as.unstartedTxs.RemoveTxByID(txID)
	if tx == nil {
		return fmt.Errorf("move_unstarted_to_in_progress: no unstarted transaction to move to in_progress")
	}
	tx.State = TxInProgress
	tx.TxAttempts = []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{txAttempt}
	tx.Sequence = &seq
	tx.BroadcastAt = &broadcastAt
	tx.InitialBroadcastAt = &initialBroadcastAt

	as.attemptHashToTxAttempt[txAttempt.Hash] = &txAttempt
	as.inprogressTx = tx

	return nil
}

// moveConfirmedMissingReceiptToUnconfirmed moves the confirmed missing receipt transaction to the unconfirmed state.
// It returns an error if there is no confirmed missing receipt transaction with the given ID.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedMissingReceiptToUnconfirmed(
	txID int64,
) error {
	as.Lock()
	defer as.Unlock()

	tx, ok := as.confirmedMissingReceiptTxs[txID]
	if !ok || tx == nil {
		return fmt.Errorf("move_confirmed_missing_receipt_to_unconfirmed: no confirmed_missing_receipt transaction with ID %d", txID)
	}

	tx.State = TxUnconfirmed
	as.unconfirmedTxs[tx.ID] = tx
	delete(as.confirmedMissingReceiptTxs, tx.ID)

	return nil
}

// moveInProgressToUnconfirmed moves the in-progress transaction to the unconfirmed state.
// It returns an error if there is no transaction in progress.
// It returns an error if the transaction in progress does not have an attempt with the given ID.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToUnconfirmed(
	txError null.String, broadcastAt time.Time, initialBroadcastAt time.Time,
	txAttemptID int64,
) error {
	as.Lock()
	defer as.Unlock()

	tx := as.inprogressTx
	if tx == nil {
		return fmt.Errorf("move_in_progress_to_unconfirmed: no transaction in progress")
	}

	found := false
	for i := 0; i < len(tx.TxAttempts); i++ {
		txAttempt := tx.TxAttempts[i]
		if txAttempt.ID == txAttemptID {
			txAttempt.State = txmgrtypes.TxAttemptBroadcast
			tx.TxAttempts[i] = txAttempt
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("move_in_progress_to_unconfirmed: transaction in progress does not have an attempt with ID %d", txAttemptID)
	}

	tx.State = TxUnconfirmed
	tx.Error = txError
	tx.BroadcastAt = &broadcastAt
	tx.InitialBroadcastAt = &initialBroadcastAt

	as.unconfirmedTxs[tx.ID] = tx
	as.inprogressTx = nil

	return nil
}

// moveUnconfirmedToConfirmed moves the unconfirmed transaction to the confirmed state.
// It returns an error if there is no unconfirmed transaction with the given ID.
// It returns an error if there is no attempt with the given receipt.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmed(
	receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
) error {
	as.Lock()
	defer as.Unlock()

	txAttempt, ok := as.attemptHashToTxAttempt[receipt.GetTxHash()]
	if !ok {
		return fmt.Errorf("move_unconfirmed_to_confirmed: no attempt found with receipt %v", receipt)
	}
	// TODO(jtw): not sure how to set blocknumber, transactionindex, and receipt on conflict
	txAttempt.Receipts = []txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH]{receipt}
	tx, ok := as.unconfirmedTxs[txAttempt.TxID]
	if !ok {
		return fmt.Errorf("move_unconfirmed_to_confirmed: no unconfirmed transaction with ID %d", txAttempt.TxID)
	}
	txAttempt.State = txmgrtypes.TxAttemptBroadcast
	if txAttempt.BroadcastBeforeBlockNum == nil {
		blockNum := receipt.GetBlockNumber().Int64()
		txAttempt.BroadcastBeforeBlockNum = &blockNum
	}
	tx.State = TxConfirmed

	return nil
}

// moveTxToFatalError moves a transaction to the fatal error state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveTxToFatalError(
	txID int64, txError null.String,
) error {
	return nil
}

// moveUnconfirmedToConfirmedMissingReceipt moves the unconfirmed transaction to the confirmed missing receipt state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// moveInProgressToConfirmedMissingReceipt moves the in-progress transaction to the confirmed missing receipt state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// moveConfirmedToUnconfirmed moves the confirmed transaction to the unconfirmed state.
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedToUnconfirmed(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}
