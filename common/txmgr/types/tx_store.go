package types

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

//go:generate mockery --quiet --name TxStore --output ./mocks/ --case=underscore
type TxStore[
	// Represents an account address, in native chain format.
	ADDR types.Hashable,
	// Represents a chain id to be used for the chain.
	CHAIN_ID ID,
	// Represents a unique Tx Hash for a chain
	TX_HASH types.Hashable,
	// Represents a unique Block Hash for a chain
	BLOCK_HASH types.Hashable,
	// Represents a onchain receipt object that a chain's RPC returns
	R ChainReceipt[TX_HASH, BLOCK_HASH],
	// Represents the sequence type for a chain. For example, nonce for EVM.
	SEQ Sequence,
	// Represents the chain specific fee type
	FEE feetypes.Fee,
] interface {
	UnstartedTxQueuePruner
	TxHistoryReaper[CHAIN_ID]
	CheckTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID, qopts ...pg.QOpt) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CreateTransaction(txRequest TxRequest[ADDR, TX_HASH], chainID CHAIN_ID, qopts ...pg.QOpt) (tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error
	Transactions(offset, limit int) ([]Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], int, error)
	TransactionsWithAttempts(offset, limit int) ([]Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], int, error)
	TxAttempts(offset, limit int) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], int, error)
	FindReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAIN_ID) (receiptsPlus []ReceiptPlus[R], err error)
	FindTxAttempt(hash TX_HASH) (*TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error)
	FindTxAttemptConfirmedByTxIDs(ids []int64) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error)
	FindTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxsRequiringResubmissionDueToInsufficientFunds(address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxAttemptsConfirmedMissingReceipt(chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxAttemptsByTxIDs(ids []int64) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error)
	FindTxAttemptsRequiringReceiptFetch(chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxByHash(hash TX_HASH) (*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], error)
	FindTxWithAttempts(etxID int64) (etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindTxWithNonce(fromAddress ADDR, seq SEQ) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	FindNextUnstartedTransactionFromAddress(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	GetTxInProgress(fromAddress ADDR, qopts ...pg.QOpt) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	GetInProgressTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], err error)
	HasInProgressTransaction(account ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (exists bool, err error)
	InsertTx(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error
	InsertTxAttempt(attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error
	LoadTxAttempts(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	LoadTxesAttempts(etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	MarkAllConfirmedMissingReceipt(chainID CHAIN_ID) (err error)
	MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAIN_ID, qopts ...pg.QOpt) error
	PreloadTxes(attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], broadcastAt time.Time) error
	SaveFetchedReceipts(receipts []R, chainID CHAIN_ID) (err error)
	SaveInProgressAttempt(attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error
	SaveInsufficientFundsAttempt(timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(oldAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], replacementAttempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	SaveSentAttempt(timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(blockNum int64, chainID CHAIN_ID) error
	UpdateBroadcastAts(now time.Time, etxIDs []int64) error
	UpdateKeyNextSequence(newNextNonce, currentNextNonce SEQ, address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	UpdateTxAttemptInProgressToBroadcast(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], NewAttemptState TxAttemptState, incrNextNonceCallback QueryerFunc, qopts ...pg.QOpt) error
	UpdateTxsUnconfirmed(ids []int64) error
	UpdateTxUnstartedToInProgress(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	UpdateTxFatalError(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], qopts ...pg.QOpt) error
	UpdateTxForRebroadcast(etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE], etxAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE]) error
	Close()
	Abandon(id CHAIN_ID, addr ADDR) error
}

type TxHistoryReaper[CHAIN_ID ID] interface {
	ReapTxHistory(minBlockNumberToKeep int64, timeThreshold time.Time, chainID CHAIN_ID) error
}

type UnstartedTxQueuePruner interface {
	PruneUnstartedTxQueue(queueSize uint32, subject uuid.UUID, qopts ...pg.QOpt) (n int64, err error)
}

// R is the raw unparsed transaction receipt
type ReceiptPlus[R any] struct {
	ID           uuid.UUID `db:"pipeline_run_id"`
	Receipt      R         `db:"receipt"`
	FailOnRevert bool      `db:"fail_on_revert"`
}

type QueryerFunc = func(tx pg.Queryer) error

type ChainReceipt[TX_HASH, BLOCK_HASH types.Hashable] interface {
	GetStatus() uint64
	GetTxHash() TX_HASH
	GetBlockNumber() *big.Int
	IsZero() bool
	IsUnmined() bool
	GetFeeUsed() uint64
	GetTransactionIndex() uint
	GetBlockHash() BLOCK_HASH
}
