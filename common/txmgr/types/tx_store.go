package types

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// TxStore is a superset of all the needed persistence layer methods
type TxStore[
	// Represents an account address, in native chain format.
	ADDR types.Hashable,
	// Represents a chain id to be used for the chain.
	CHAIN_ID types.ID,
	// Represents a unique Tx Hash for a chain
	TX_HASH types.Hashable,
	// Represents a unique Block Hash for a chain
	BLOCK_HASH types.Hashable,
	// Represents a onchain receipt object that a chain's RPC returns
	R ChainReceipt[TX_HASH, BLOCK_HASH],
	// Represents the sequence type for a chain. For example, nonce for EVM.
	SEQ types.Sequence,
	// Represents the chain specific fee type
	FEE feetypes.Fee,
] interface {
	UnstartedTxQueuePruner
	TxHistoryReaper[CHAIN_ID]
	TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE]

	// Find confirmed txes beyond the minConfirmations param that require callback but have not yet been signaled
	FindTxesPendingCallback(ctx context.Context, blockNum int64, chainID CHAIN_ID) (receiptsPlus []ReceiptPlus[R], err error)
	// Update tx to mark that its callback has been signaled
	UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainId CHAIN_ID) error
	SaveFetchedReceipts(ctx context.Context, r []R, state TxState, errorMsg *string, chainID CHAIN_ID) error

	// additional methods for tx store management
	CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID) (err error)
	Close()
	Abandon(ctx context.Context, id CHAIN_ID, addr ADDR) error
	// Find transactions by a field in the TxMeta blob and transaction states
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []TxState, chainID *big.Int) (tx []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided by transaction states
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []TxState, chainID *big.Int) (tx []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (tx []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []TxState, chainID *big.Int) (tx []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (tx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
}

// TransactionStore contains the persistence layer methods needed to manage Txs and TxAttempts
type TransactionStore[
	ADDR types.Hashable,
	CHAIN_ID types.ID,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	CountUnconfirmedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (count uint32, err error)
	CountTransactionsByState(ctx context.Context, state TxState, chainID CHAIN_ID) (count uint32, err error)
	CountUnstartedTransactions(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (count uint32, err error)
	CreateTransaction(ctx context.Context, txRequest TxRequest[ADDR, TX_HASH], chainID CHAIN_ID) (tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	FindLatestSequence(ctx context.Context, fromAddress ADDR, chainId CHAIN_ID) (SEQ, error)
	FindTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindTxsRequiringResubmissionDueToInsufficientFunds(ctx context.Context, address ADDR, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindTxAttemptsConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindTxAttemptsRequiringReceiptFetch(ctx context.Context, chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindTxAttemptsRequiringResend(ctx context.Context, olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Search for Tx using the idempotencyKey and chainID
	FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CHAIN_ID) (tx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Search for Tx using the fromAddress and sequence
	FindTxWithSequence(ctx context.Context, fromAddress ADDR, seq SEQ) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindNextUnstartedTransactionFromAddress(ctx context.Context, fromAddress ADDR, chainID CHAIN_ID) (*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error)

	// FindTransactionsConfirmedInBlockRange retrieves tx with attempts and partial receipt values for optimization purpose
	FindTransactionsConfirmedInBlockRange(ctx context.Context, highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindEarliestUnconfirmedBroadcastTime(ctx context.Context, chainID CHAIN_ID) (null.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context, chainID CHAIN_ID) (null.Int, error)
	GetTxInProgress(ctx context.Context, fromAddress ADDR) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetInProgressTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetAbandonedTransactionsByBatch(ctx context.Context, chainID CHAIN_ID, enabledAddrs []ADDR, offset, limit uint) (txs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetTxByID(ctx context.Context, id int64) (tx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	HasInProgressTransaction(ctx context.Context, account ADDR, chainID CHAIN_ID) (exists bool, err error)
	LoadTxAttempts(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	MarkAllConfirmedMissingReceipt(ctx context.Context, chainID CHAIN_ID) (err error)
	MarkOldTxesMissingReceiptAsErrored(ctx context.Context, blockNum int64, latestFinalizedBlockNum int64, chainID CHAIN_ID) error
	PreloadTxes(ctx context.Context, attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error
	SaveInProgressAttempt(ctx context.Context, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	SaveInsufficientFundsAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(ctx context.Context, oldAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], replacementAttempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	SaveSentAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(ctx context.Context, blockNum int64, chainID CHAIN_ID) error
	UpdateBroadcastAts(ctx context.Context, now time.Time, etxIDs []int64) error
	UpdateTxAttemptInProgressToBroadcast(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], NewAttemptState TxAttemptState) error
	// Update tx to mark that its callback has been signaled
	UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainId CHAIN_ID) error
	UpdateTxsUnconfirmed(ctx context.Context, ids []int64) error
	UpdateTxUnstartedToInProgress(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	UpdateTxFatalError(ctx context.Context, etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
	UpdateTxForRebroadcast(ctx context.Context, etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], etxAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) error
}

type TxHistoryReaper[CHAIN_ID types.ID] interface {
	ReapTxHistory(ctx context.Context, timeThreshold time.Time, chainID CHAIN_ID) error
}

type UnstartedTxQueuePruner interface {
	PruneUnstartedTxQueue(ctx context.Context, queueSize uint32, subject uuid.UUID) (ids []int64, err error)
}

// R is the raw unparsed transaction receipt
type ReceiptPlus[R any] struct {
	ID           uuid.UUID `db:"pipeline_run_id"`
	Receipt      R         `db:"receipt"`
	FailOnRevert bool      `db:"fail_on_revert"`
}

type ChainReceipt[TX_HASH, BLOCK_HASH types.Hashable] interface {
	GetStatus() uint64
	GetTxHash() TX_HASH
	GetBlockNumber() *big.Int
	IsZero() bool
	IsUnmined() bool
	GetFeeUsed() uint64
	GetTransactionIndex() uint
	GetBlockHash() BLOCK_HASH
	GetRevertReason() *string
}
