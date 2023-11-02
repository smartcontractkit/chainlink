package txmgr

import (
	"context"
	"fmt"
	"sync"

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
	// ErrExistingPilelineTaskRunId is returned when a transaction with the same pipeline task run id already exists
	ErrExistingPilelineTaskRunId = fmt.Errorf("transaction with pipeline task run id already exists")
)

// Store and update all transaction state as files
// Read from the files to restore state at startup
// Delete files when transactions are completed or reaped

// Life of a Transaction
// 1. Transaction Request is created
// 2. Transaction Request is submitted to the Transaction Manager
// 3. Transaction Manager creates and persists a new transaction (unstarted) from the transaction request (not persisted)
// 4. Transaction Manager sends the transaction (unstarted) to the Broadcaster Unstarted Queue
// 4. Transaction Manager prunes the Unstarted Queue based on the transaction prune strategy

// NOTE(jtw): Gets triggered by postgres Events
// NOTE(jtw): Only one transaction per address can be in_progress at a time
// NOTE(jtw): Only one broadcasted attempt exists per transaction the rest are errored or abandoned
// 1. Broadcaster assigns a sequence number to the transaction
// 2. Broadcaster creates and persists a new transaction attempt (in_progress) from the transaction (in_progress)
// 3. Broadcaster asks the Checker to check if the transaction should not be sent
// 4. Broadcaster asks the Attempt builder to figure out gas fee for the transaction
// 5. Broadcaster attempts to send the Transaction to TransactionClient to be published on-chain
// 6. Broadcaster updates the transaction attempt (broadcast) and transaction (unconfirmed)
// 7. Broadcaster increments global sequence number for address for next transaction attempt

// NOTE(jtw): Only one receipt should exist per confirmed transaction
// 1. Confirmer listens and reads new Head events from the Chain
// 2. Confirmer sets the last known block number for the transaction attempts that have been broadcast
// 3. Confirmer checks for missing receipts for transactions that have been broadcast
// 4. Confirmer sets transactions that have failed to (unconfirmed) which will be retried by the resender
// 5. Confirmer sets transactions that have been confirmed to (confirmed) and creates a new receipt which is persisted

type InMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	chainID CHAIN_ID

	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	txStore  txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	pendingLock sync.Mutex
	// NOTE(jtw): we might need to watch out for txns that finish and are removed from the pending map
	pendingIdempotencyKeys    map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	pendingPipelineTaskRunIds map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	unstarted  map[ADDR]chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	inprogress map[ADDR]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
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
	txStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
) (*InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	tm := InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		chainID:  chainID,
		keyStore: keyStore,
		txStore:  txStore,

		pendingIdempotencyKeys:    map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		pendingPipelineTaskRunIds: map[string]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},

		unstarted:  map[ADDR]chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
		inprogress: map[ADDR]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{},
	}

	addresses, err := keyStore.EnabledAddressesForChain(chainID)
	if err != nil {
		return nil, err
	}
	for _, fromAddr := range addresses {
		// Channel Buffer is set to something high to prevent blocking and allow the pruning to happen
		tm.unstarted[fromAddr] = make(chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 100)
	}

	return &tm, nil
}

// Close closes the InMemoryStore
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
	// Close the event recorder
	ms.txStore.Close()

	// Close all channels
	for _, ch := range ms.unstarted {
		close(ch)
	}

	// Clear all pending requests
	ms.pendingLock.Lock()
	clear(ms.pendingIdempotencyKeys)
	clear(ms.pendingPipelineTaskRunIds)
	ms.pendingLock.Unlock()
}

// Abandon removes all transactions for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Abandon(ctx context.Context, chainID CHAIN_ID, addr ADDR) error {
	if ms.chainID.String() != chainID.String() {
		return ErrInvalidChainID
	}

	// Mark all persisted transactions as abandoned
	if err := ms.txStore.Abandon(ctx, chainID, addr); err != nil {
		return err
	}

	// Mark all unstarted transactions as abandoned
	close(ms.unstarted[addr])
	for tx := range ms.unstarted[addr] {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}
	// reset the unstarted channel
	ms.unstarted[addr] = make(chan *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], 100)

	// Mark all inprogress transactions as abandoned
	if tx, ok := ms.inprogress[addr]; ok {
		tx.State = TxFatalError
		tx.Sequence = nil
		tx.Error = null.NewString("abandoned", true)
	}
	ms.inprogress[addr] = nil

	// TODO(jtw): Mark all unconfirmed transactions as abandoned

	// Mark all pending transactions as abandoned
	for _, tx := range ms.pendingIdempotencyKeys {
		if tx.FromAddress == addr {
			tx.State = TxFatalError
			tx.Sequence = nil
			tx.Error = null.NewString("abandoned", true)
		}
	}
	for _, tx := range ms.pendingPipelineTaskRunIds {
		if tx.FromAddress == addr {
			tx.State = TxFatalError
			tx.Sequence = nil
			tx.Error = null.NewString("abandoned", true)
		}
	}
	// TODO(jtw): SHOULD THE REAPER BE RESPONSIBLE FOR CLEARING THE PENDING MAPS?

	return nil
}

// CreateTransaction creates a new transaction for a given txRequest.
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH], chainID CHAIN_ID) (txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	// Persist Transaction to persistent storage
	tx, err := ms.txStore.CreateTransaction(ctx, txRequest, chainID)
	if err != nil {
		return tx, err
	}
	return tx, ms.sendTxToBroadcaster(tx)
}

// TODO(jtw): change naming to something more appropriate
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) sendTxToBroadcaster(tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	// TODO(jtw); HANDLE PRUNING STEP

	select {
	// Add the request to the Unstarted channel to be processed by the Broadcaster
	case ms.unstarted[tx.FromAddress] <- &tx:
		// Persist to persistent storage

		ms.pendingLock.Lock()
		if tx.IdempotencyKey != nil {
			ms.pendingIdempotencyKeys[*tx.IdempotencyKey] = &tx
		}
		if tx.PipelineTaskRunID.UUID.String() != "" {
			ms.pendingPipelineTaskRunIds[tx.PipelineTaskRunID.UUID.String()] = &tx
		}
		ms.pendingLock.Unlock()

		return nil
	default:
		// Return an error if the Manager Queue Capacity has been reached
		return fmt.Errorf("transaction manager queue capacity has been reached")
	}
}

// FindTxWithIdempotencyKey returns a transaction with the given idempotency key
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	// TODO(jtw): this is a change from current functionality... it returns nil, nil if nothing found in other implementations
	if ms.chainID.String() != chainID.String() {
		return nil, ErrInvalidChainID
	}
	if idempotencyKey == "" {
		return nil, fmt.Errorf("FindTxWithIdempotencyKey: idempotency key cannot be empty")
	}

	ms.pendingLock.Lock()
	defer ms.pendingLock.Unlock()

	tx, ok := ms.pendingIdempotencyKeys[idempotencyKey]
	if !ok {
		return nil, fmt.Errorf("FindTxWithIdempotencyKey: transaction not found")
	}

	return tx, nil
}

// CheckTxQueueCapacity checks if the queue capacity has been reached for a given address
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID) (err error) {
	if maxQueuedTransactions == 0 {
		return nil
	}
	if ms.chainID.String() != chainID.String() {
		return ErrInvalidChainID
	}
	if _, ok := ms.unstarted[fromAddress]; !ok {
		return fmt.Errorf("CheckTxQueueCapacity: address not found")
	}

	count := uint64(len(ms.unstarted[fromAddress]))
	if count >= maxQueuedTransactions {
		return fmt.Errorf("CheckTxQueueCapacity: cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}

	return nil
}

/*
// BROADCASTER FUNCTIONS
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) CountUnconfirmedTransactions(_ context.Context, fromAddress ADDR, chainID CHAIN_ID) (count uint32, err error) {
	if ms.chainID != chainID {
		return 0, ErrInvalidChainID
	}
	// TODO(jtw): NEED TO COMPLETE
	return 0, nil
}
func (ms *InMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) CountUnstartedTransactions(_ context.Context, fromAddress ADDR, chainID CHAIN_ID) (count uint32, err error) {
	if ms.chainID != chainID {
		return 0, ErrInvalidChainID
	}

	return uint32(len(ms.unstarted[fromAddress])), nil
}
func (ms *InMemoryStore) FindNextUnstartedTransactionFromAddress(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR, chainID CHAIN_ID) error {

}
func (ms *InMemoryStore) GetTxInProgress(ctx context.Context, fromAddress ADDR) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
}

func (ms *InMemoryStore) UpdateTxAttemptInProgressToBroadcast(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], NewAttemptState TxAttemptState, incrNextSequenceCallback QueryerFunc, qopts ...pg.QOpt) error {
}
func (ms *InMemoryStore) UpdateTxUnstartedToInProgress(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
}
func (ms *InMemoryStore) UpdateTxFatalError(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
}
func (ms *InMemoryStore) SaveReplacementInProgressAttempt(ctx context.Context, oldAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], replacementAttempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
}

// ProcessUnstartedTxs processes unstarted transactions
// TODO(jtw): SHOULD THIS BE CALLED THE BROADCASTER?
func (tm *TransactionManager) ProcessUnstartedTxs(ctx context.Context, fromAddress string) {
	// if there are more in flight transactions than the max then throttle using the InFlightTransactionRecheckInterval
	for {
		select {
		// NOTE: There can be at most one in_progress transaction per address.
		case txReq := <-tm.Unstarted[fromAddress]:
			// check number of in flight transactions to see if we can process more... MaxInFlight is for total inflight
			tm.inFlightWG.Wait()
			// Reserve a spot in the in flight transactions
			tm.inFlightWG.Done()

			// TODO(jtw): THERE ARE SOME CHANGES AROUND ERROR FUNCTIONALITY
			// EXAMPLE NO LONGER WILL ERROR IF THE NUMBER OF IN FLIGHT TRANSACTIONS IS EXCEEDED
			if err := tm.PublishToChain(txReq); err != nil {
				// TODO(jtw): Handle error properly
				fmt.Println(err)
			}

			// Free up a spot in the in flight transactions
			tm.inFlightWG.Add(1)
		case <-ctx.Done():
			return
		}
	}

}

// PublishToChain attempts to publish a transaction to the chain
// TODO(jtw): NO LONGER RETURNS AN ERROR IF FULL OF IN PROGRESS TRANSACTIONS... not sure if okay
func (tm *TransactionManager) PublishToChain(txReq TxRequest) error {
	// Handle an unstarted request
	// Get next sequence number from the KeyStore
	// ks.NextSequence(fromAddress, tm.ChainID)
	// Create a new transaction attempt to be put on chain
	// eb.NewTxAttempt(ctx, txReq, logger)

	// IT BLOCKS UNTIL THERE IS A SPOT IN THE IN PROGRESS TRANSACTIONS
	tm.Inprogress[txReq.FromAddress] <- txReq

	return nil
}
*/

/*
// Close closes the InMemoryStore
func (ms *InMemoryStore) Close() {
	// Close all channels
	for _, ch := range ms.Unstarted {
		close(ch)
	}
	for _, ch := range ms.Inprogress {
		close(ch)
	}

	// Clear all pending requests
	ms.pendingLock.Lock()
	clear(ms.pendingIdempotencyKeys)
	clear(ms.pendingPipelineTaskRunIds)
	ms.pendingLock.Unlock()
}
*/
