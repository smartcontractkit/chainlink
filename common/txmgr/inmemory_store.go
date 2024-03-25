package txmgr

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/chains/label"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

var (
	// ErrInvalidChainID is returned when the chain ID is invalid
	ErrInvalidChainID = errors.New("invalid chain ID")
	// ErrTxnNotFound is returned when a transaction is not found
	ErrTxnNotFound = errors.New("transaction not found")
	// ErrExistingIdempotencyKey is returned when a transaction with the same idempotency key already exists
	ErrExistingIdempotencyKey = errors.New("transaction with idempotency key already exists")
	// ErrAddressNotFound is returned when an address is not found
	ErrAddressNotFound = errors.New("address not found")
	// ErrSequenceNotFound is returned when a sequence is not found
	ErrSequenceNotFound = errors.New("sequence not found")
	// ErrCouldNotGetReceipt is the error string we save if we reach our finality depth for a confirmed transaction without ever getting a receipt
	// This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = errors.New("could not get receipt")
)

// inMemoryStore is a simple wrapper around a persistent TxStore and KeyStore. It handles all the transaction state in memory.
type inMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	lggr    logger.SugaredLogger
	chainID CHAIN_ID

	maxUnstarted      uint64
	keyStore          txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	persistentTxStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE]

	addressStatesLock sync.RWMutex
	addressStates     map[ADDR]*addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]
}

// NewInMemoryStore returns a new inMemoryStore
func NewInMemoryStore[
	CHAIN_ID types.ID,
	ADDR, TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	ctx context.Context,
	lggr logger.SugaredLogger,
	chainID CHAIN_ID,
	keyStore txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	persistentTxStore txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE],
	config txmgrtypes.InMemoryStoreConfig,
) (*inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error) {
	ms := inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{
		lggr:              lggr,
		chainID:           chainID,
		keyStore:          keyStore,
		persistentTxStore: persistentTxStore,

		addressStates: map[ADDR]*addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]{},
	}

	ms.maxUnstarted = config.MaxQueued()
	if ms.maxUnstarted <= 0 {
		ms.maxUnstarted = 10000
	}

	txs, err := persistentTxStore.GetAllTransactions(ctx, chainID)
	if err != nil {
		return nil, fmt.Errorf("address_state: initialization: %w", err)
	}
	addressesToTxs := map[ADDR][]txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, tx := range txs {
		at, exists := addressesToTxs[tx.FromAddress]
		if !exists {
			at = []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
		}
		at = append(at, tx)
		addressesToTxs[tx.FromAddress] = at
	}
	for fromAddr, txs := range addressesToTxs {
		ms.addressStates[fromAddr] = newAddressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE](lggr, chainID, fromAddr, ms.maxUnstarted, txs)
	}

	return &ms, nil
}

// CreateTransaction creates a new transaction for a given txRequest.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CreateTransaction(
	ctx context.Context,
	txRequest txmgrtypes.TxRequest[ADDR, TX_HASH],
	chainID CHAIN_ID,
) (
	txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	error,
) {
	return txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}, nil
}

// FindTxWithIdempotencyKey returns a transaction with the given idempotency key
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}

	// Check if the transaction is in the pending queue of all address states
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		if tx := as.findTxWithIdempotencyKey(idempotencyKey); tx != nil {
			return ms.deepCopyTx(*tx), nil
		}
	}

	return nil, nil
}

// CheckTxQueueCapacity checks if the queue capacity has been reached for a given address
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID) error {
	if maxQueuedTransactions == 0 {
		return nil
	}
	if ms.chainID.String() != chainID.String() {
		return nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return nil
	}

	count := uint64(as.countTransactionsByState(TxUnstarted))
	if count >= maxQueuedTransactions {
		return fmt.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}

	return nil
}

// FindLatestSequence returns the latest sequence number for a given address
// It is used to initialize the in-memory sequence map in the broadcaster
// TODO(jtw): this is until we have a abstracted Sequencer Component which can be used instead
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindLatestSequence(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (seq SEQ, err error) {
	// Query the persistent store
	return ms.persistentTxStore.FindLatestSequence(ctx, fromAddress, chainID)
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions for a given address.
// Unconfirmed transactions are transactions that have been broadcast but not confirmed on-chain.
// NOTE(jtw): used to calculate total inflight transactions
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnconfirmedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	if ms.chainID.String() != chainID.String() {
		return 0, nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return 0, nil
	}

	return uint32(as.countTransactionsByState(TxUnconfirmed)), nil
}

// CountUnstartedTransactions returns the number of unstarted transactions for a given address.
// Unstarted transactions are transactions that have not been broadcast yet.
// NOTE(jtw): used to calculate total inflight transactions
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountUnstartedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (uint32, error) {
	if ms.chainID.String() != chainID.String() {
		return 0, nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return 0, nil
	}

	return uint32(as.countTransactionsByState(TxUnstarted)), nil
}

// UpdateTxUnstartedToInProgress updates a transaction from unstarted to in_progress.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxUnstartedToInProgress(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// GetTxInProgress returns the in_progress transaction for a given address.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetTxInProgress(ctx context.Context, fromAddress ADDR) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	return nil, nil
}

// UpdateTxAttemptInProgressToBroadcast updates a transaction attempt from in_progress to broadcast.
// It also updates the transaction state to unconfirmed.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxAttemptInProgressToBroadcast(
	ctx context.Context,
	tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	newAttemptState txmgrtypes.TxAttemptState,
) error {
	return nil
}

// FindNextUnstartedTransactionFromAddress returns the next unstarted transaction for a given address.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindNextUnstartedTransactionFromAddress(_ context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR, chainID CHAIN_ID) error {
	return nil
}

// SaveReplacementInProgressAttempt saves a replacement attempt for a transaction that is in_progress.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveReplacementInProgressAttempt(
	ctx context.Context,
	oldAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	replacementAttempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) error {
	return nil
}

// UpdateTxFatalError updates a transaction to fatal_error.
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxFatalError(ctx context.Context, tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

// Close closes the inMemoryStore
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Close() {
}

// Abandon removes all transactions for a given address
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) Abandon(ctx context.Context, chainID CHAIN_ID, addr ADDR) error {
	return nil
}

// SetBroadcastBeforeBlockNum sets the broadcast_before_block_num for a given chain ID
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SetBroadcastBeforeBlockNum(ctx context.Context, blockNum int64, chainID CHAIN_ID) error {
	return nil
}

// FindTxAttemptsConfirmedMissingReceipt returns all transactions that are confirmed but missing a receipt
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) (
	[]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	error,
) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}

	txFilter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.TxAttempts != nil && len(tx.TxAttempts) > 0
	}
	txAttemptFilter := func(attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return true
	}
	states := []txmgrtypes.TxState{TxConfirmedMissingReceipt}
	attempts := []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		attempts = append(attempts, as.findTxAttempts(states, txFilter, txAttemptFilter)...)
	}
	// sort by tx_id ASC, gas_price DESC, gas_tip_cap DESC
	sort.SliceStable(attempts, func(i, j int) bool {
		if attempts[i].TxID == attempts[j].TxID {
			var iGasPrice, jGasPrice evmgas.EvmFee
			// TODO: FIGURE OUT HOW TO GET GAS PRICE AND GAS TIP CAP FROM TxFee

			if iGasPrice.Legacy.Cmp(jGasPrice.Legacy) == 0 {
				// sort by gas_tip_cap DESC
				return iGasPrice.DynamicFeeCap.Cmp(jGasPrice.DynamicFeeCap) > 0
			} else {
				// sort by gas_price DESC
				return iGasPrice.Legacy.Cmp(jGasPrice.Legacy) > 0
			}
		}

		return attempts[i].TxID < attempts[j].TxID
	})

	// deep copy the attempts
	var eAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, attempt := range attempts {
		eAttempts = append(eAttempts, ms.deepCopyTxAttempt(attempt.Tx, attempt))
	}

	return eAttempts, nil
}

// UpdateBroadcastAts updates the broadcast_at time for a given set of attempts
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateBroadcastAts(ctx context.Context, now time.Time, txIDs []int64) error {
	return nil
}

// UpdateTxsUnconfirmed updates the unconfirmed transactions for a given set of ids
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxsUnconfirmed(ctx context.Context, txIDs []int64) error {
	return nil
}

// FindTxAttemptsRequiringReceiptFetch returns all transactions that are missing a receipt
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsRequiringReceiptFetch(ctx context.Context, chainID CHAIN_ID) (
	attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	err error,
) {
	if ms.chainID.String() != chainID.String() {
		return attempts, nil
	}

	txFilterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return len(tx.TxAttempts) > 0
	}
	txAttemptFilterFn := func(attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return attempt.State != txmgrtypes.TxAttemptInsufficientFunds
	}
	states := []txmgrtypes.TxState{TxUnconfirmed, TxConfirmedMissingReceipt}
	attempts = []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		attempts = append(attempts, as.findTxAttempts(states, txFilterFn, txAttemptFilterFn)...)
	}
	// sort by sequence ASC, gas_price DESC, gas_tip_cap DESC
	// TODO: FIGURE OUT HOW TO GET GAS PRICE AND GAS TIP CAP FROM TxFee
	slices.SortFunc(attempts, func(a, b txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int {
		aSequence, bSequence := a.Tx.Sequence, b.Tx.Sequence
		if aSequence == nil || bSequence == nil {
			return 0
		}

		return cmp.Compare((*aSequence).Int64(), (*bSequence).Int64())
	})

	// deep copy the attempts
	var eAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, attempt := range attempts {
		eAttempts = append(eAttempts, ms.deepCopyTxAttempt(attempt.Tx, attempt))
	}

	return eAttempts, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesPendingCallback(ctx context.Context, blockNum int64, chainID CHAIN_ID) (
	[]txmgrtypes.ReceiptPlus[R],
	error,
) {
	if ms.chainID.String() != chainID.String() {
		panic("invalid chain ID")
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return false
		}

		for i := 0; i < len(tx.TxAttempts); i++ {
			if len(tx.TxAttempts[i].Receipts) == 0 {
				continue
			}

			if !tx.PipelineTaskRunID.Valid || !tx.SignalCallback || tx.CallbackCompleted {
				continue
			}
			receipt := tx.TxAttempts[i].Receipts[0]
			minConfirmations := int64(tx.MinConfirmations.Uint32)
			if receipt.GetBlockNumber() != nil &&
				receipt.GetBlockNumber().Int64() <= (blockNum-minConfirmations) {
				return true
			}
		}

		return false

	}
	states := []txmgrtypes.TxState{TxConfirmed}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		txs = append(txs, as.findTxs(states, filterFn)...)
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

		for j := 0; j < len(tx.TxAttempts); j++ {
			if len(tx.TxAttempts[j].Receipts) == 0 {
				continue
			}
			receiptsPlus[i] = txmgrtypes.ReceiptPlus[R]{
				ID:           tx.PipelineTaskRunID.UUID,
				Receipt:      tx.TxAttempts[j].Receipts[0].(R),
				FailOnRevert: failOnRevert,
			}
		}
		clear(meta)
	}

	return receiptsPlus, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainId CHAIN_ID) error {
	return nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveFetchedReceipts(ctx context.Context, receipts []R, chainID CHAIN_ID) error {
	return nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (
	[]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	error,
) {
	if ms.chainID.String() != chainID.String() {
		panic("invalid chain ID")
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.Meta == nil {
			return false
		}
		meta := map[string]interface{}{}
		if err := json.Unmarshal(json.RawMessage(*tx.Meta), &meta); err != nil {
			return false
		}
		return isMetaValueEqual(meta[metaField], metaValue)
	}
	txs := []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		for _, tx := range as.findTxs(states, filterFn) {
			etx := ms.deepCopyTx(tx)
			txs = append(txs, etx)
		}
	}

	return txs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		panic("invalid chain ID")
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
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		for _, tx := range as.findTxs(states, filterFn) {
			etx := ms.deepCopyTx(tx)
			txs = append(txs, etx)
		}
	}

	return txs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		panic("invalid chain ID")
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.Meta == nil {
			return false
		}
		meta := map[string]interface{}{}
		if err := json.Unmarshal(json.RawMessage(*tx.Meta), &meta); err != nil {
			return false
		}
		if _, ok := meta[metaField]; !ok {
			return false
		}
		if len(tx.TxAttempts) == 0 {
			return false
		}

		for _, attempt := range tx.TxAttempts {
			if len(attempt.Receipts) == 0 {
				continue
			}
			if attempt.Receipts[0].GetBlockNumber() == nil {
				continue
			}
			return attempt.Receipts[0].GetBlockNumber().Int64() >= blockNum
		}

		return false
	}

	txs := []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		for _, tx := range as.findTxs(nil, filterFn) {
			etx := ms.deepCopyTx(tx)
			txs = append(txs, etx)
		}
	}

	return txs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []big.Int, states []txmgrtypes.TxState, chainID *big.Int) (tx []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("find_txes_with_attempts_and_receipts_by_ids_and_state: %w", ErrInvalidChainID)
	}

	filterFn := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return true
	}

	txIDs := make([]int64, len(ids))
	for i, id := range ids {
		txIDs[i] = id.Int64()
	}

	txsLock := sync.Mutex{}
	txs := []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	wg := sync.WaitGroup{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		wg.Add(1)
		go func(as *addressState[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) {
			for _, tx := range as.findTxs(states, filterFn, txIDs...) {
				etx := ms.deepCopyTx(tx)
				txsLock.Lock()
				txs = append(txs, etx)
				txsLock.Unlock()
			}
			wg.Done()
		}(as)
	}
	wg.Wait()

	return txs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PruneUnstartedTxQueue(ctx context.Context, queueSize uint32, subject uuid.UUID) ([]int64, error) {
	return nil, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) ReapTxHistory(ctx context.Context, minBlockNumberToKeep int64, timeThreshold time.Time, chainID CHAIN_ID) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) CountTransactionsByState(_ context.Context, state txmgrtypes.TxState, chainID CHAIN_ID) (uint32, error) {
	if ms.chainID.String() != chainID.String() {
		return 0, nil
	}

	var total int
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		total += as.countTransactionsByState(state)
	}

	return uint32(total), nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) DeleteInProgressAttempt(ctx context.Context, attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringResubmissionDueToInsufficientFunds(_ context.Context, address ADDR, chainID CHAIN_ID) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[address]
	if !ok {
		return nil, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return false
		}
		for _, attempt := range tx.TxAttempts {
			if attempt.State == txmgrtypes.TxAttemptInsufficientFunds {
				return true
			}
		}
		return false
	}
	states := []txmgrtypes.TxState{TxUnconfirmed}
	txs := as.findTxs(states, filter)
	// sort by sequence ASC
	slices.SortFunc(txs, func(a, b txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int {
		aSequence, bSequence := a.Sequence, b.Sequence
		if aSequence == nil || bSequence == nil {
			return 0
		}

		return cmp.Compare((*aSequence).Int64(), (*bSequence).Int64())
	})

	etxs := make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs))
	for i, tx := range txs {
		etxs[i] = ms.deepCopyTx(tx)
	}

	return etxs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxAttemptsRequiringResend(_ context.Context, olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) ([]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		panic("invalid chain ID")
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[address]
	if !ok {
		return nil, nil
	}

	txFilter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return false
		}
		return tx.BroadcastAt.Before(olderThan) || tx.BroadcastAt.Equal(olderThan)
	}
	txAttemptFilter := func(attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return attempt.State != txmgrtypes.TxAttemptInProgress
	}
	states := []txmgrtypes.TxState{TxUnconfirmed, TxConfirmedMissingReceipt}
	attempts := as.findTxAttempts(states, txFilter, txAttemptFilter)
	// sort by sequence ASC, gas_price DESC, gas_tip_cap DESC
	slices.SortFunc(attempts, func(a, b txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int {
		aSequence, bSequence := a.Tx.Sequence, b.Tx.Sequence
		if aSequence == nil || bSequence == nil {
			return 0
		}
		// TODO: FIGURE OUT HOW TO GET GAS PRICE AND GAS TIP CAP FROM TxFee
		/*
			v, ok := a.TxFee.(*gas.EvmFee)
			if !ok {
				panic("invalid gas fee")
			}
			fmt.Println("hereh", v)
		*/

		return cmp.Compare((*aSequence).Int64(), (*bSequence).Int64())
	})
	// LIMIT by maxInFlightTransactions
	if maxInFlightTransactions > 0 && len(attempts) > int(maxInFlightTransactions) {
		attempts = attempts[:maxInFlightTransactions]
	}

	// deep copy the attempts
	var eAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, attempt := range attempts {
		eAttempts = append(eAttempts, ms.deepCopyTxAttempt(attempt.Tx, attempt))
	}

	return eAttempts, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxWithSequence(_ context.Context, fromAddress ADDR, seq SEQ) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[fromAddress]
	if !ok {
		return nil, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.Sequence == nil {
			return false
		}

		return (*tx.Sequence).String() == seq.String()
	}
	states := []txmgrtypes.TxState{TxConfirmed, TxConfirmedMissingReceipt, TxUnconfirmed}
	txs := as.findTxs(states, filter)
	if len(txs) == 0 {
		return nil, nil
	}

	return ms.deepCopyTx(txs[0]), nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTransactionsConfirmedInBlockRange(_ context.Context, highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return false
		}
		for _, attempt := range tx.TxAttempts {
			if attempt.State != txmgrtypes.TxAttemptBroadcast {
				continue
			}
			if len(attempt.Receipts) == 0 {
				continue
			}
			if attempt.Receipts[0].GetBlockNumber() == nil {
				continue
			}
			blockNum := attempt.Receipts[0].GetBlockNumber().Int64()
			if blockNum >= lowBlockNumber && blockNum <= highBlockNumber {
				return true
			}
		}

		return false
	}
	states := []txmgrtypes.TxState{TxConfirmed, TxConfirmedMissingReceipt}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		ts := as.findTxs(states, filter)
		txs = append(txs, ts...)
	}
	// sort by sequence ASC
	slices.SortFunc(txs, func(a, b txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int {
		aSequence, bSequence := a.Sequence, b.Sequence
		if aSequence == nil || bSequence == nil {
			return 0
		}

		return cmp.Compare((*aSequence).Int64(), (*bSequence).Int64())
	})

	etxs := make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs))
	for i, tx := range txs {
		etxs[i] = ms.deepCopyTx(tx)
	}

	return etxs, nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedBroadcastTime(ctx context.Context, chainID CHAIN_ID) (null.Time, error) {
	if ms.chainID.String() != chainID.String() {
		return null.Time{}, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.InitialBroadcastAt != nil
	}
	states := []txmgrtypes.TxState{TxUnconfirmed}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		etxs := as.findTxs(states, filter)
		txs = append(txs, etxs...)
	}

	var minInitialBroadcastAt time.Time
	for _, tx := range txs {
		if tx.InitialBroadcastAt == nil {
			continue
		}
		if minInitialBroadcastAt.IsZero() || tx.InitialBroadcastAt.Before(minInitialBroadcastAt) {
			minInitialBroadcastAt = *tx.InitialBroadcastAt
		}
	}
	if minInitialBroadcastAt.IsZero() {
		return null.Time{}, nil
	}

	return null.TimeFrom(minInitialBroadcastAt), nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context, chainID CHAIN_ID) (null.Int, error) {
	if ms.chainID.String() != chainID.String() {
		return null.Int{}, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return false
		}

		for _, attempt := range tx.TxAttempts {
			if attempt.BroadcastBeforeBlockNum != nil {
				return true
			}
		}

		return false
	}
	states := []txmgrtypes.TxState{TxUnconfirmed}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		etxs := as.findTxs(states, filter)
		txs = append(txs, etxs...)
	}

	var minBroadcastBeforeBlockNum int64
	for _, tx := range txs {
		if minBroadcastBeforeBlockNum == 0 || *tx.TxAttempts[0].BroadcastBeforeBlockNum < minBroadcastBeforeBlockNum {
			minBroadcastBeforeBlockNum = *tx.TxAttempts[0].BroadcastBeforeBlockNum
		}
	}
	if minBroadcastBeforeBlockNum == 0 {
		return null.Int{}, nil
	}

	return null.IntFrom(minBroadcastBeforeBlockNum), nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetInProgressTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) ([]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, fmt.Errorf("get_in_progress_tx_attempts: %w", ErrInvalidChainID)
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[address]
	if !ok {
		return nil, fmt.Errorf("get_in_progress_tx_attempts: %w", ErrAddressNotFound)
	}

	txFilter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.TxAttempts != nil && len(tx.TxAttempts) > 0
	}
	txAttemptFilter := func(attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return attempt.State == txmgrtypes.TxAttemptInProgress
	}
	states := []txmgrtypes.TxState{TxConfirmed, TxConfirmedMissingReceipt, TxUnconfirmed}
	attempts := as.findTxAttempts(states, txFilter, txAttemptFilter)

	// deep copy the attempts
	var eAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, attempt := range attempts {
		eAttempts = append(eAttempts, ms.deepCopyTxAttempt(attempt.Tx, attempt))
	}

	return eAttempts, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetNonFatalTransactions(ctx context.Context, chainID CHAIN_ID) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.State != TxFatalError
	}
	txs := []txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		etxs := as.findTxs(nil, filter)
		txs = append(txs, etxs...)
	}

	etxs := make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs))
	for i, tx := range txs {
		etxs[i] = ms.deepCopyTx(tx)
	}

	return etxs, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) GetTxByID(_ context.Context, id int64) (*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.ID == id
	}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		txs := as.findTxs(nil, filter, id)
		if len(txs) > 0 {
			return ms.deepCopyTx(txs[0]), nil
		}
	}

	return nil, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) HasInProgressTransaction(_ context.Context, account ADDR, chainID CHAIN_ID) (bool, error) {
	if ms.chainID.String() != chainID.String() {
		return false, fmt.Errorf("has_in_progress_transaction: %w", ErrInvalidChainID)
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[account]
	if !ok {
		return false, fmt.Errorf("has_in_progress_transaction: %w", ErrAddressNotFound)
	}

	n := as.countTransactionsByState(TxInProgress)

	return n > 0, nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) LoadTxAttempts(_ context.Context, etx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[etx.FromAddress]
	if !ok {
		return nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return tx.ID == etx.ID
	}
	txAttempts := []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{}
	for _, tx := range as.findTxs(nil, filter, etx.ID) {
		for _, txAttempt := range tx.TxAttempts {
			txAttempts = append(txAttempts, ms.deepCopyTxAttempt(*etx, txAttempt))
		}
	}
	etx.TxAttempts = txAttempts

	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) PreloadTxes(_ context.Context, attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	if len(attempts) == 0 {
		return nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[attempts[0].Tx.FromAddress]
	if !ok {
		return nil
	}

	txIDs := make([]int64, len(attempts))
	for i, attempt := range attempts {
		txIDs[i] = attempt.TxID
	}
	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return true
	}
	txs := as.findTxs(nil, filter, txIDs...)
	for i, attempt := range attempts {
		for _, tx := range txs {
			if tx.ID == attempt.TxID {
				attempts[i].Tx = *ms.deepCopyTx(tx)
			}
		}
	}

	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveInProgressAttempt(ctx context.Context, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveInsufficientFundsAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) SaveSentAttempt(ctx context.Context, timeout time.Duration, attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) UpdateTxForRebroadcast(ctx context.Context, etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], etxAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) IsTxFinalized(ctx context.Context, blockHeight int64, txID int64, chainID CHAIN_ID) (bool, error) {
	if ms.chainID.String() != chainID.String() {
		return false, nil
	}

	txFilter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if tx.ID != txID {
			return false
		}

		for _, attempt := range tx.TxAttempts {
			if len(attempt.Receipts) == 0 {
				continue
			}
			// there can only be one receipt per attempt
			if attempt.Receipts[0].GetBlockNumber() == nil {
				continue
			}
			return attempt.Receipts[0].GetBlockNumber().Int64() <= (blockHeight - int64(tx.MinConfirmations.Uint32))
		}

		return false
	}
	txAttemptFilter := func(attempt *txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		return len(attempt.Receipts) > 0
	}
	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	for _, as := range ms.addressStates {
		txas := as.findTxAttempts(nil, txFilter, txAttemptFilter, txID)
		if len(txas) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) FindTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) ([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error) {
	if ms.chainID.String() != chainID.String() {
		return nil, nil
	}
	if gasBumpThreshold == 0 {
		return nil, nil
	}

	ms.addressStatesLock.RLock()
	defer ms.addressStatesLock.RUnlock()
	as, ok := ms.addressStates[address]
	if !ok {
		return nil, nil
	}

	filter := func(tx *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) bool {
		if len(tx.TxAttempts) == 0 {
			return true
		}
		for _, attempt := range tx.TxAttempts {
			if attempt.BroadcastBeforeBlockNum == nil || *attempt.BroadcastBeforeBlockNum > blockNum-gasBumpThreshold || attempt.State != txmgrtypes.TxAttemptBroadcast {
				return false
			}
		}

		return true
	}
	states := []txmgrtypes.TxState{TxUnconfirmed}
	txs := as.findTxs(states, filter)

	// sort by sequence ASC
	slices.SortFunc(txs, func(a, b txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) int {
		aSequence, bSequence := a.Sequence, b.Sequence
		if aSequence == nil || bSequence == nil {
			return 0
		}

		return cmp.Compare((*aSequence).Int64(), (*bSequence).Int64())
	})

	if depth > 0 {
		// LIMIT by depth
		if len(txs) > int(depth) {
			txs = txs[:depth]
		}
	}

	etxs := make([]*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(txs))
	for i, tx := range txs {
		etxs[i] = ms.deepCopyTx(tx)
	}

	return etxs, nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkAllConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) error {
	return nil
}
func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) MarkOldTxesMissingReceiptAsErrored(ctx context.Context, blockNum int64, finalityDepth uint32, chainID CHAIN_ID) error {
	return nil
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deepCopyTx(
	tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) *txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	copyTx := txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		ID:                 tx.ID,
		IdempotencyKey:     tx.IdempotencyKey,
		Sequence:           tx.Sequence,
		FromAddress:        tx.FromAddress,
		ToAddress:          tx.ToAddress,
		EncodedPayload:     make([]byte, len(tx.EncodedPayload)),
		Value:              *new(big.Int).Set(&tx.Value),
		FeeLimit:           tx.FeeLimit,
		Error:              tx.Error,
		BroadcastAt:        tx.BroadcastAt,
		InitialBroadcastAt: tx.InitialBroadcastAt,
		CreatedAt:          tx.CreatedAt,
		State:              tx.State,
		TxAttempts:         make([]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], len(tx.TxAttempts)),
		Meta:               tx.Meta,
		Subject:            tx.Subject,
		ChainID:            tx.ChainID,
		PipelineTaskRunID:  tx.PipelineTaskRunID,
		MinConfirmations:   tx.MinConfirmations,
		TransmitChecker:    tx.TransmitChecker,
		SignalCallback:     tx.SignalCallback,
		CallbackCompleted:  tx.CallbackCompleted,
	}

	// Copy the EncodedPayload
	copy(copyTx.EncodedPayload, tx.EncodedPayload)

	// Copy the TxAttempts
	for i, attempt := range tx.TxAttempts {
		copyTx.TxAttempts[i] = ms.deepCopyTxAttempt(copyTx, attempt)
	}

	return &copyTx
}

func (ms *inMemoryStore[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) deepCopyTxAttempt(
	tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
) txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	copyAttempt := txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		ID:                      attempt.ID,
		TxID:                    attempt.TxID,
		Tx:                      tx,
		TxFee:                   attempt.TxFee,
		ChainSpecificFeeLimit:   attempt.ChainSpecificFeeLimit,
		SignedRawTx:             make([]byte, len(attempt.SignedRawTx)),
		Hash:                    attempt.Hash,
		CreatedAt:               attempt.CreatedAt,
		BroadcastBeforeBlockNum: attempt.BroadcastBeforeBlockNum,
		State:                   attempt.State,
		Receipts:                make([]txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH], len(attempt.Receipts)),
		TxType:                  attempt.TxType,
	}

	copy(copyAttempt.SignedRawTx, attempt.SignedRawTx)
	copy(copyAttempt.Receipts, attempt.Receipts)

	return copyAttempt
}

func isMetaValueEqual(v interface{}, metaValue string) bool {
	switch v := v.(type) {
	case string:
		return v == metaValue
	case int:
		o, err := strconv.ParseInt(metaValue, 10, 64)
		if err != nil {
			return false
		}
		return v == int(o)
	case uint32:
		o, err := strconv.ParseUint(metaValue, 10, 32)
		if err != nil {
			return false
		}
		return v == uint32(o)
	case uint64:
		o, err := strconv.ParseUint(metaValue, 10, 64)
		if err != nil {
			return false
		}
		return v == o
	case int32:
		o, err := strconv.ParseInt(metaValue, 10, 32)
		if err != nil {
			return false
		}
		return v == int32(o)
	case int64:
		o, err := strconv.ParseInt(metaValue, 10, 64)
		if err != nil {
			return false
		}
		return v == o
	case float32:
		o, err := strconv.ParseFloat(metaValue, 32)
		if err != nil {
			return false
		}
		return v == float32(o)
	case float64:
		o, err := strconv.ParseFloat(metaValue, 64)
		if err != nil {
			return false
		}
		return v == o
	case bool:
		o, err := strconv.ParseBool(metaValue)
		if err != nil {
			return false
		}
		return v == o
	}

	return false
}
