package txmgr

import (
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
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
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
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) countTransactionsByState(txState txmgrtypes.TxState) int {
	return 0
}

// findTxWithIdempotencyKey returns the transaction with the given idempotency key. If no transaction is found, nil is returned.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxWithIdempotencyKey(key string) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return nil
}

// applyToTxsByState calls the given function for each transaction in the given states.
// If txIDs are provided, only the transactions with those IDs are considered.
// If no txIDs are provided, all transactions in the given states are considered.
// If no txStates are provided, all transactions are considered.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
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
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
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
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) findTxs(
	txStates []txmgrtypes.TxState,
	filter func(*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool,
	txIDs ...int64,
) []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	return nil
}

// pruneUnstartedTxQueue removes the transactions with the given IDs from the unstarted transaction queue.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) pruneUnstartedTxQueue(ids []int64) {
}

// deleteTxs removes the transactions with the given IDs from the address state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deleteTxs(txs ...txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
}

// peekNextUnstartedTx returns the next unstarted transaction in the queue without removing it from the unstarted queue.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekNextUnstartedTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// peekInProgressTx returns the in-progress transaction without removing it from the in-progress state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) peekInProgressTx() (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// addTxToUnstarted adds the given transaction to the unstarted queue.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) addTxToUnstarted(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

// moveUnstartedToInProgress moves the next unstarted transaction to the in-progress state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnstartedToInProgress(
	etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// moveConfirmedMissingReceiptToUnconfirmed moves the confirmed missing receipt transaction to the unconfirmed state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedMissingReceiptToUnconfirmed(
	txID int64,
) error {
	return nil
}

// moveInProgressToUnconfirmed moves the in-progress transaction to the unconfirmed state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToUnconfirmed(
	etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	txAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// moveUnconfirmedToConfirmed moves the unconfirmed transaction to the confirmed state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmed(
	receipt txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
) error {
	return nil
}

// moveTxToFatalError moves a transaction to the fatal error state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveTxToFatalError(
	txID int64, txError null.String,
) error {
	return nil
}

// moveUnconfirmedToConfirmedMissingReceipt moves the unconfirmed transaction to the confirmed missing receipt state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveUnconfirmedToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// moveInProgressToConfirmedMissingReceipt moves the in-progress transaction to the confirmed missing receipt state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveInProgressToConfirmedMissingReceipt(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}

// moveConfirmedToUnconfirmed moves the confirmed transaction to the unconfirmed state.
//
//lint:ignore U1000 Ignore unused function temporarily while adding the framework
func (as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) moveConfirmedToUnconfirmed(attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}
