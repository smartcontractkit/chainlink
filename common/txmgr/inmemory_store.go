package txmgr

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"gopkg.in/guregu/null.v4"
)

var (
	// ErrInvalidChainID is returned when the chain ID is invalid
	ErrInvalidChainID = fmt.Errorf("invalid chain ID")
	// ErrTxnNotFound is returned when a transaction is not found
	ErrTxnNotFound = fmt.Errorf("transaction not found")
	// ErrExistingIdempotencyKey is returned when a transaction with the same idempotency key already exists
	ErrExistingIdempotencyKey = fmt.Errorf("transaction with idempotency key already exists")
	// ErrAddressNotFound is returned when an address is not found
	ErrAddressNotFound = fmt.Errorf("address not found")
	// ErrSequenceNotFound is returned when a sequence is not found
	ErrSequenceNotFound = fmt.Errorf("sequence not found")
)

type PersistentTxStore[
	ADDR types.Hashable,
	CHAIN_ID types.ID,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	UnstartedTransactions(limit, offset int, fromAddress ADDR, chainID CHAIN_ID) ([]txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], int, error)
	UnconfirmedTransactions(limit, offset int, fromAddress ADDR, chainID CHAIN_ID) ([]txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], int, error)
	ConfirmedTransactions(limit, offset int, fromAddress ADDR, chainID CHAIN_ID) ([]txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], int, error)
	ConfirmedMissingReceiptTransactions(limit, offset int, fromAddress ADDR, chainID CHAIN_ID) ([]txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], int, error)
}

type InMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	chainID CHAIN_ID

	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	txStore  PersistentTxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	addressStatesLock sync.RWMutex
	addressStates     map[ADDR]*AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
}

// NewInMemoryStore returns a new InMemoryStore
func NewInMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	chainID CHAIN_ID,
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	txStore PersistentTxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
) (*InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	ms := InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		chainID:  chainID,
		keyStore: keyStore,
		txStore:  txStore,

		addressStates: map[ADDR]*AddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{},
	}

	maxUnstarted := 50
	addresses, err := keyStore.EnabledAddressesForChain(chainID)
	if err != nil {
		return nil, fmt.Errorf("new_in_memory_store: %w", err)
	}
	for _, fromAddr := range addresses {
		as, err := NewAddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](chainID, fromAddr, maxUnstarted, txStore)
		if err != nil {
			return nil, fmt.Errorf("new_in_memory_store: %w", err)
		}

		ms.addressStates[fromAddr] = as
	}

	return &ms, nil
}

// CreateTransaction creates a new transaction for a given txRequest.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH], chainID CHAIN_ID) (tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if ms.chainID.String() != chainID.String() {
		return tx, fmt.Errorf("create_transaction: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[tx.FromAddress]
	if !ok {
		return tx, fmt.Errorf("create_transaction: %w", ErrAddressNotFound)
	}

	// Persist Transaction to persistent storage
	tx, err = ms.txStore.CreateTransaction(ctx, txRequest, chainID)
	if err != nil {
		return tx, fmt.Errorf("create_transaction: %w", err)
	}

	// TODO(jtw); HANDLE PRUNING STEP

	// Add the request to the Unstarted channel to be processed by the Broadcaster
	if err := as.addTxToUnstarted(&tx); err != nil {
		return tx, fmt.Errorf("create_transaction: %w", err)
	}

	return tx, nil
}

// FindTxWithIdempotencyKey returns a transaction with the given idempotency key
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_tx_with_idempotency_key: %w", ErrInvalidChainID)
	}

	// Check if the transaction is in the pending queue of all address states
	ms.addressStatesLock.Lock()
	defer ms.addressStatesLock.Unlock()
	for _, as := range ms.addressStates {
		if tx := as.findTxWithIdempotencyKey(idempotencyKey); tx != nil {
			return ms.deepCopyTx(tx), nil
		}
	}

	return nil, fmt.Errorf("find_tx_with_idempotency_key: %w", ErrTxnNotFound)

}

// CheckTxQueueCapacity checks if the queue capacity has been reached for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID) error {
	if maxQueuedTransactions == 0 {
		return nil
	}
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("check_tx_queue_capacity: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return fmt.Errorf("check_tx_queue_capacity: %w", ErrAddressNotFound)
	}

	count := uint64(as.unstartedCount())
	if count >= maxQueuedTransactions {
		return fmt.Errorf("check_tx_queue_capacity: cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}

	return nil
}

/////
// BROADCASTER FUNCTIONS
/////

// FindLatestSequence returns the latest sequence number for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindLatestSequence(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (seq SEQ, err error) {
	// query the persistent storage since this method only gets called when the broadcaster is starting up.
	// It is used to initialize the in-memory sequence map in the broadcaster
	// NOTE(jtw): should the nextSequenceMap be moved to the in-memory store?

	if ms.chainID.String() != chainID.String() {
		return seq, fmt.Errorf("find_latest_sequence: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return seq, fmt.Errorf("find_latest_sequence: %w", ErrAddressNotFound)
	}

	seq = as.findLatestSequence()
	if seq.Int64() == 0 {
		return seq, fmt.Errorf("find_latest_sequence: %w", ErrSequenceNotFound)
	}

	return seq, nil
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions for a given address.
// Unconfirmed transactions are transactions that have been broadcast but not confirmed on-chain.
// NOTE(jtw): used to calculate total inflight transactions
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnconfirmedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	if ms.chainID.String() != chainID.String() {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrAddressNotFound)
	}

	return uint32(as.unconfirmedCount()), nil
}

// CountUnstartedTransactions returns the number of unstarted transactions for a given address.
// Unstarted transactions are transactions that have not been broadcast yet.
// NOTE(jtw): used to calculate total inflight transactions
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnstartedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	if ms.chainID.String() != chainID.String() {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return 0, fmt.Errorf("count_unstarted_transactions: %w", ErrAddressNotFound)
	}

	return uint32(as.unstartedCount()), nil
}

// UpdateTxUnstartedToInProgress updates a transaction from unstarted to in_progress.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxUnstartedToInProgress(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	if tx.Sequence == nil {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: in_progress transaction must have a sequence number")
	}
	if tx.State != TxUnstarted {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: can only transition to in_progress from unstarted, transaction is currently %s", tx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: attempt state must be in_progress")
	}
	as, ok := ms.addressStates[tx.FromAddress]
	if !ok {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: %w", ErrAddressNotFound)
	}

	// Persist to persistent storage
	if err := ms.txStore.UpdateTxUnstartedToInProgress(ctx, tx, attempt); err != nil {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: %w", err)
	}
	tx.TxAttempts = append(tx.TxAttempts, *attempt)

	// Update in address state in memory
	if err := as.moveUnstartedToInProgress(tx); err != nil {
		return fmt.Errorf("update_tx_unstarted_to_in_progress: %w", err)
	}

	return nil
}

// GetTxInProgress returns the in_progress transaction for a given address.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetTxInProgress(ctx context.Context, fromAddress ADDR) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return nil, fmt.Errorf("get_tx_in_progress: %w", ErrAddressNotFound)
	}

	tx, err := as.peekInProgressTx()
	if tx == nil {
		return nil, fmt.Errorf("get_tx_in_progress: %w", err)
	}

	// NOTE(jtw): should this exist in the in-memory store? or just the persistent store?
	// NOTE(jtw): where should this live?
	if len(tx.TxAttempts) != 1 || tx.TxAttempts[0].State != txmgrtypes.TxAttemptInProgress {
		return nil, fmt.Errorf("get_tx_in_progress: expected in_progress transaction %v to have exactly one unsent attempt. "+
			"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", tx.ID)
	}

	return ms.deepCopyTx(tx), nil
}

// UpdateTxAttemptInProgressToBroadcast updates a transaction attempt from in_progress to broadcast.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxAttemptInProgressToBroadcast(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	newAttemptState txmgrtypes.TxAttemptState,
) error {
	if tx.BroadcastAt == nil {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: unconfirmed transaction must have broadcast_at time")
	}
	if tx.InitialBroadcastAt == nil {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: unconfirmed transaction must have initial_broadcast_at time")
	}
	if tx.State != TxInProgress {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: can only transition to unconfirmed from in_progress, transaction is currently %s", tx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: attempt must be in in_progress state")
	}
	if newAttemptState != txmgrtypes.TxAttemptBroadcast {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: new attempt state must be broadcast, got: %s", newAttemptState)
	}

	// Persist to persistent storage
	if err := ms.txStore.UpdateTxAttemptInProgressToBroadcast(ctx, tx, attempt, newAttemptState); err != nil {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: %w", err)
	}
	// Ensure that the tx state is updated to unconfirmed since this is a chain agnostic operation
	attempt.State = newAttemptState

	as, ok := ms.addressStates[tx.FromAddress]
	if !ok {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: %w", ErrAddressNotFound)
	}
	if err := as.moveInProgressToUnconfirmed(attempt); err != nil {
		return fmt.Errorf("update_tx_attempt_in_progress_to_broadcast: %w", err)
	}

	return nil
}

// FindNextUnstartedTransactionFromAddress returns the next unstarted transaction for a given address.
// NOTE(jtw): method signature is different from most other signatures where the tx is passed in and updated
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindNextUnstartedTransactionFromAddress(_ context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR, chainID CHAIN_ID) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: %w", ErrInvalidChainID)
	}
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: %w", ErrAddressNotFound)
	}

	// ensure that the address is not already busy with a transaction in progress
	if as.inprogress != nil {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: address %s is already busy with a transaction in progress", fromAddress)
	}

	var err error
	tx, err = as.peekNextUnstartedTx()
	if tx == nil {
		return fmt.Errorf("find_next_unstarted_transaction_from_address: %w", err)
	}

	return nil
}

// SaveReplacementInProgressAttempt saves a replacement attempt for a transaction that is in_progress.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveReplacementInProgressAttempt(
	ctx context.Context,
	oldAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	replacementAttempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	if oldAttempt.State != txmgrtypes.TxAttemptInProgress || replacementAttempt.State != txmgrtypes.TxAttemptInProgress {
		return fmt.Errorf("save_replacement_in_progress_attempt: expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return fmt.Errorf("save_replacement_in_progress_attempt: expected oldattempt to have an ID")
	}

	// Persist to persistent storage
	if err := ms.txStore.SaveReplacementInProgressAttempt(ctx, oldAttempt, replacementAttempt); err != nil {
		return fmt.Errorf("save_replacement_in_progress_attempt: %w", err)
	}

	// Update in memory store
	as, ok := ms.addressStates[oldAttempt.Tx.FromAddress]
	if !ok {
		return fmt.Errorf("save_replacement_in_progress_attempt: %w", ErrAddressNotFound)
	}
	tx, err := as.peekInProgressTx()
	if tx == nil {
		return fmt.Errorf("save_replacement_in_progress_attempt: %w", err)
	}
	var found bool
	for i := 0; i < len(tx.TxAttempts); i++ {
		if tx.TxAttempts[i].ID == oldAttempt.ID {
			tx.TxAttempts[i] = *replacementAttempt
			found = true
		}
	}
	if !found {
		tx.TxAttempts = append(tx.TxAttempts, *replacementAttempt)
		// NOTE(jtw): should this log a warning?
	}

	return nil
}

// UpdateTxFatalError updates a transaction to fatal_error.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxFatalError(ctx context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	if tx.State != TxInProgress {
		return fmt.Errorf("update_tx_fatal_error: can only transition to fatal_error from in_progress, transaction is currently %s", tx.State)
	}
	if !tx.Error.Valid {
		return fmt.Errorf("update_tx_fatal_error: expected error field to be set")
	}

	// Persist to persistent storage
	if err := ms.txStore.UpdateTxFatalError(ctx, tx); err != nil {
		return fmt.Errorf("update_tx_fatal_error: %w", err)
	}

	// Ensure that the tx state is updated to fatal_error since this is a chain agnostic operation
	tx.Sequence = nil
	tx.State = TxFatalError

	return fmt.Errorf("update_tx_fatal_error: not implemented")
}

// Close closes the InMemoryStore
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	// Close the event recorder
	ms.txStore.Close()

	// Clear all address states
	ms.addressStatesLock.Lock()
	for _, as := range ms.addressStates {
		as.close()
	}
	clear(ms.addressStates)
	ms.addressStatesLock.Unlock()
}

// Abandon removes all transactions for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Abandon(ctx context.Context, chainID CHAIN_ID, addr ADDR) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("abandon: %w", ErrInvalidChainID)
	}

	// Mark all persisted transactions as abandoned
	if err := ms.txStore.Abandon(ctx, chainID, addr); err != nil {
		return err
	}

	// check that the address exists in the unstarted transactions
	as, ok := ms.addressStates[addr]
	if !ok {
		return fmt.Errorf("abandon: %w", ErrAddressNotFound)
	}
	as.abandon()

	return nil
}

// SetBroadcastBeforeBlockNum sets the broadcast_before_block_num for a given chain ID
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetBroadcastBeforeBlockNum(ctx context.Context, blockNum int64, chainID CHAIN_ID) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("set_broadcast_before_block_num: %w", ErrInvalidChainID)
	}

	// Persist to persistent storage
	if err := ms.txStore.SetBroadcastBeforeBlockNum(ctx, blockNum, chainID); err != nil {
		return fmt.Errorf("set_broadcast_before_block_num: %w", err)
	}

	fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
		if tx.TxAttempts == nil || len(tx.TxAttempts) == 0 {
			return
		}
		// TODO(jtw): how many tx_attempts are actually stored in the db for each tx? It looks like its only 1
		attempt := tx.TxAttempts[0]
		if attempt.State == txmgrtypes.TxAttemptBroadcast && attempt.BroadcastBeforeBlockNum == nil &&
			tx.ChainID.String() == chainID.String() {
			tx.TxAttempts[0].BroadcastBeforeBlockNum = &blockNum
		}
	}
	for _, as := range ms.addressStates {
		as.ApplyToTxs(nil, fn)
	}

	return nil
}

// FindTxAttemptsConfirmedMissingReceipt returns all transactions that are confirmed but missing a receipt
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) ([]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_next_unstarted_transaction_from_address: %w", ErrInvalidChainID)
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.TxAttempts != nil && len(tx.TxAttempts) > 0 {
			return tx.ChainID.String() == chainID.String()
		}

		return false
	}
	states := []txmgrtypes.TxState{TxConfirmedMissingReceipt}
	attempts := []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range ms.addressStates {
		attempts = append(attempts, as.FetchTxAttempts(states, filter)...)
	}
	// sort by tx_id ASC, gas_price DESC, gas_tip_cap DESC
	// TODO

	return attempts, nil
}

// UpdateBroadcastAts updates the broadcast_at time for a given set of attempts
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateBroadcastAts(ctx context.Context, now time.Time, txIDs []int64) error {
	// Persist to persistent storage
	if err := ms.txStore.UpdateBroadcastAts(ctx, now, txIDs); err != nil {
		return err
	}

	// Update in memory store
	fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
		if tx.BroadcastAt != nil {
			tx.BroadcastAt = &now
		}
	}

	for _, as := range ms.addressStates {
		as.ApplyToTxs(nil, fn, txIDs...)
	}

	return nil
}

// UpdateTxsUnconfirmed updates the unconfirmed transactions for a given set of ids
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxsUnconfirmed(ctx context.Context, txIDs []int64) error {
	// Persist to persistent storage
	if err := ms.txStore.UpdateTxsUnconfirmed(ctx, txIDs); err != nil {
		return err
	}

	// Update in memory store
	fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
		tx.State = TxUnconfirmed
	}

	for _, as := range ms.addressStates {
		as.ApplyToTxs(nil, fn, txIDs...)
	}

	return nil
}

// FindTxAttemptsRequiringReceiptFetch returns all transactions that are missing a receipt
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsRequiringReceiptFetch(ctx context.Context, chainID CHAIN_ID) (attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if ms.chainID.String() != chainID.String() {
		return attempts, fmt.Errorf("find_tx_attempts_requiring_receipt_fetch: %w", ErrInvalidChainID)
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.TxAttempts != nil && len(tx.TxAttempts) > 0 {
			attempt := tx.TxAttempts[0]
			return attempt.State != txmgrtypes.TxAttemptInsufficientFunds
		}

		return false
	}
	states := []txmgrtypes.TxState{TxUnconfirmed, TxConfirmedMissingReceipt}
	attempts = []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range ms.addressStates {
		attempts = append(attempts, as.FetchTxAttempts(states, filterFn)...)
	}
	// sort by sequence ASC, gas_price DESC, gas_tip_cap DESC
	// TODO

	return attempts, nil
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesPendingCallback(ctx context.Context, blockNum int64, chainID CHAIN_ID) ([]txmgrtypes.ReceiptPlus[R], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_txes_pending_callback: %w", ErrInvalidChainID)
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.TxAttempts == nil || len(tx.TxAttempts) == 0 {
			return false
		}

		if tx.TxAttempts[0].Receipts == nil || len(tx.TxAttempts[0].Receipts) == 0 {
			return false
		}

		if tx.PipelineTaskRunID.Valid && tx.SignalCallback && !tx.CallbackCompleted &&
			tx.TxAttempts[0].Receipts[0].GetBlockNumber() != nil &&
			big.NewInt(blockNum-int64(tx.MinConfirmations.Uint32)).Cmp(tx.TxAttempts[0].Receipts[0].GetBlockNumber()) > 0 {
			return true
		}

		return false

	}
	states := []txmgrtypes.TxState{TxConfirmed}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range ms.addressStates {
		txs = append(txs, as.FetchTxs(states, filterFn)...)
	}

	receiptsPlus := make([]txmgrtypes.ReceiptPlus[R], len(txs))
	meta := map[string]interface{}{}
	for i, tx := range txs {
		if err := json.Unmarshal(json.RawMessage(*tx.Meta), &meta); err != nil {
			return nil, err
		}
		failOnRevert := false
		if v, ok := meta["FailOnRevert"].(bool); ok {
			failOnRevert = v
		}

		receiptsPlus[i] = txmgrtypes.ReceiptPlus[R]{
			ID:           tx.PipelineTaskRunID.UUID,
			Receipt:      (tx.TxAttempts[0].Receipts[0]).(R),
			FailOnRevert: failOnRevert,
		}
		clear(meta)
	}

	return receiptsPlus, nil
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainId CHAIN_ID) error {
	if ms.chainID.String() != chainId.String() {
		return fmt.Errorf("update_tx_callback_completed: %w", ErrInvalidChainID)
	}

	// Persist to persistent storage
	if err := ms.txStore.UpdateTxCallbackCompleted(ctx, pipelineTaskRunRid, chainId); err != nil {
		return err
	}

	// Update in memory store
	fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
		if tx.PipelineTaskRunID.UUID == pipelineTaskRunRid {
			tx.CallbackCompleted = true
		}
	}
	for _, as := range ms.addressStates {
		as.ApplyToTxs(nil, fn)
	}

	return nil
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveFetchedReceipts(ctx context.Context, receipts []R, chainID CHAIN_ID) error {
	if ms.chainID.String() != chainID.String() {
		return fmt.Errorf("save_fetched_receipts: %w", ErrInvalidChainID)
	}

	// Persist to persistent storage
	if err := ms.txStore.SaveFetchedReceipts(ctx, receipts, chainID); err != nil {
		return err
	}

	// convert receipts to map
	receiptsMap := map[TX_HASH]R{}
	for _, receipt := range receipts {
		receiptsMap[receipt.GetTxHash()] = receipt
	}

	// Update in memory store
	fn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) {
		if tx.TxAttempts == nil || len(tx.TxAttempts) == 0 {
			return
		}
		attempt := tx.TxAttempts[0]
		receipt, ok := receiptsMap[attempt.Hash]
		if !ok {
			return
		}

		if attempt.Receipts != nil && len(attempt.Receipts) > 0 &&
			attempt.Receipts[0].GetBlockNumber() != nil && receipt.GetBlockNumber() != nil &&
			attempt.Receipts[0].GetBlockNumber().Cmp(receipt.GetBlockNumber()) == 0 {
			return
		}
		// TODO(jtw): this needs to be finished

		attempt.State = txmgrtypes.TxAttemptBroadcast
		if attempt.BroadcastBeforeBlockNum == nil {
			blocknum := receipt.GetBlockNumber().Int64()
			attempt.BroadcastBeforeBlockNum = &blocknum
		}
		attempt.Receipts = []txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH]{receipt}

		tx.State = TxConfirmed
	}
	for _, as := range ms.addressStates {
		as.ApplyToTxs(nil, fn)
	}

	return nil
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_txes_by_meta_field_and_states: %w", ErrInvalidChainID)
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.Meta == nil {
			return false
		}
		meta := map[string]interface{}{}
		if err := json.Unmarshal(json.RawMessage(*tx.Meta), &meta); err != nil {
			return false
		}
		if v, ok := meta[metaField].(string); ok {
			return v == metaValue
		}

		return false
	}
	txs := []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range ms.addressStates {
		for _, tx := range as.FetchTxs(states, filterFn) {
			txs = append(txs, &tx)
		}
	}

	return txs, nil
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_txes_with_meta_field_by_states: %w", ErrInvalidChainID)
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.Meta == nil {
			return false
		}
		meta := map[string]interface{}{}
		if err := json.Unmarshal(json.RawMessage(*tx.Meta), &meta); err != nil {
			return false
		}
		if _, ok := meta[metaField]; ok {
			return true
		}

		return false
	}

	txs := []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, as := range ms.addressStates {
		for _, tx := range as.FetchTxs(states, filterFn) {
			txs = append(txs, &tx)
		}
	}

	return txs, nil
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return ms.txStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, metaField, blockNum, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []big.Int, states []txmgrtypes.TxState, chainID *big.Int) (tx []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx, ids, states, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneUnstartedTxQueue(ctx context.Context, queueSize uint32, subject uuid.UUID) (int64, error) {
	return ms.txStore.PruneUnstartedTxQueue(ctx, queueSize, subject)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ReapTxHistory(ctx context.Context, minBlockNumberToKeep int64, timeThreshold time.Time, chainID CHAIN_ID) error {
	return ms.txStore.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState, chainID CHAIN_ID) (uint32, error) {
	return ms.txStore.CountTransactionsByState(ctx, state, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) DeleteInProgressAttempt(ctx context.Context, attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return ms.txStore.DeleteInProgressAttempt(ctx, attempt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTxsRequiringGasBump(ctx, address, blockNum, gasBumpThreshold, depth, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringResubmissionDueToInsufficientFunds(ctx context.Context, address ADDR, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, address, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsRequiringResend(ctx context.Context, olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) (attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTxAttemptsRequiringResend(ctx, olderThan, maxInFlightTransactions, chainID, address)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithSequence(ctx context.Context, fromAddress ADDR, seq SEQ) (etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTxWithSequence(ctx, fromAddress, seq)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTransactionsConfirmedInBlockRange(ctx context.Context, highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) (etxs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.FindTransactionsConfirmedInBlockRange(ctx, highBlockNumber, lowBlockNumber, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedBroadcastTime(ctx context.Context, chainID CHAIN_ID) (null.Time, error) {
	return ms.txStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context, chainID CHAIN_ID) (null.Int, error) {
	return ms.txStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetInProgressTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) (attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.GetInProgressTxAttempts(ctx, address, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetNonFatalTransactions(ctx context.Context, chainID CHAIN_ID) (txs []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.GetNonFatalTransactions(ctx, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetTxByID(ctx context.Context, id int64) (tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	return ms.txStore.GetTxByID(ctx, id)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HasInProgressTransaction(ctx context.Context, account ADDR, chainID CHAIN_ID) (exists bool, err error) {
	return ms.txStore.HasInProgressTransaction(ctx, account, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) LoadTxAttempts(ctx context.Context, etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return ms.txStore.LoadTxAttempts(ctx, etx)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkAllConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) (err error) {
	return ms.txStore.MarkAllConfirmedMissingReceipt(ctx, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkOldTxesMissingReceiptAsErrored(ctx context.Context, blockNum int64, finalityDepth uint32, chainID CHAIN_ID) error {
	return ms.txStore.MarkOldTxesMissingReceiptAsErrored(ctx, blockNum, finalityDepth, chainID)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PreloadTxes(ctx context.Context, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return ms.txStore.PreloadTxes(ctx, attempts)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return ms.txStore.SaveConfirmedMissingReceiptAttempt(ctx, timeout, attempt, broadcastAt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveInProgressAttempt(ctx context.Context, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return ms.txStore.SaveInProgressAttempt(ctx, attempt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveInsufficientFundsAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return ms.txStore.SaveInsufficientFundsAttempt(ctx, timeout, attempt, broadcastAt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveSentAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return ms.txStore.SaveSentAttempt(ctx, timeout, attempt, broadcastAt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxForRebroadcast(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], etxAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return ms.txStore.UpdateTxForRebroadcast(ctx, etx, etxAttempt)
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) IsTxFinalized(ctx context.Context, blockHeight int64, txID int64, chainID CHAIN_ID) (finalized bool, err error) {
	return ms.txStore.IsTxFinalized(ctx, blockHeight, txID, chainID)
}

func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deepCopyTx(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	etx := *tx
	etx.TxAttempts = make([]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(tx.TxAttempts))
	copy(etx.TxAttempts, tx.TxAttempts)

	return &etx
}
