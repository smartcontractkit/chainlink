package types

import (
	"context"
	"time"

	"github.com/google/uuid"

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
	R ChainReceipt[TX_HASH],
	// Represents the sequence type for a chain. For example, nonce for EVM.
	SEQ Sequence,
	FEE Fee,
	ADD any,
] interface {
	UnstartedTxQueuePruner
	CheckEthTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID, qopts ...pg.QOpt) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CreateEthTransaction(newTx NewTx[ADDR, TX_HASH], chainID CHAIN_ID, qopts ...pg.QOpt) (tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) error
	EthTransactions(offset, limit int) ([]Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], int, error)
	EthTransactionsWithAttempts(offset, limit int) ([]Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], int, error)
	EthTxAttempts(offset, limit int) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], int, error)
	FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAIN_ID) (receiptsPlus []ReceiptPlus[R], err error)
	FindEthTxAttempt(hash TX_HASH) (*TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], error)
	FindEthTxAttemptConfirmedByEthTxIDs(ids []int64) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], error)
	FindEthTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEthTxsRequiringResubmissionDueToInsufficientEth(address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEtxAttemptsConfirmedMissingReceipt(chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEthTxAttemptsByEthTxIDs(ids []int64) ([]TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], error)
	FindEthTxAttemptsRequiringReceiptFetch(chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEthTxByHash(hash TX_HASH) (*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], error)
	FindEthTxWithAttempts(etxID int64) (etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindEthTxWithNonce(fromAddress ADDR, seq SEQ) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	FindNextUnstartedTransactionFromAddress(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) (etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	GetEthTxInProgress(fromAddress ADDR, qopts ...pg.QOpt) (etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	GetInProgressEthTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) (attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], err error)
	HasInProgressTransaction(account ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (exists bool, err error)
	// InsertEthReceipt only used in tests. Use SaveFetchedReceipts instead
	InsertEthReceipt(receipt *Receipt[R, TX_HASH, BLOCK_HASH]) error
	InsertEthTx(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) error
	InsertEthTxAttempt(attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) error
	LoadEthTxAttempts(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	LoadEthTxesAttempts(etxs []*Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	MarkAllConfirmedMissingReceipt(chainID CHAIN_ID) (err error)
	MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAIN_ID, qopts ...pg.QOpt) error
	PreloadEthTxes(attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], broadcastAt time.Time) error
	SaveFetchedReceipts(receipts []R, chainID CHAIN_ID) (err error)
	SaveInProgressAttempt(attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) error
	SaveInsufficientEthAttempt(timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(oldAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], replacementAttempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	SaveSentAttempt(timeout time.Duration, attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(blockNum int64, chainID CHAIN_ID) error
	UpdateBroadcastAts(now time.Time, etxIDs []int64) error
	UpdateEthKeyNextNonce(newNextNonce, currentNextNonce SEQ, address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	UpdateEthTxAttemptInProgressToBroadcast(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], NewAttemptState TxAttemptState, incrNextNonceCallback QueryerFunc, qopts ...pg.QOpt) error
	UpdateEthTxsUnconfirmed(ids []int64) error
	UpdateEthTxUnstartedToInProgress(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], attempt *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	UpdateEthTxFatalError(etx *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], qopts ...pg.QOpt) error
	UpdateEthTxForRebroadcast(etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], etxAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) error
	Close()
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

// R is the raw unparsed transaction receipt
type Receipt[R any, TX_HASH types.Hashable, BLOCK_HASH types.Hashable] struct {
	ID               int64
	TxHash           TX_HASH
	BlockHash        BLOCK_HASH
	BlockNumber      int64
	TransactionIndex uint
	Receipt          R
	CreatedAt        time.Time
}

type QueryerFunc = func(tx pg.Queryer) error

type ChainReceipt[TX_HASH types.Hashable] interface {
	GetStatus() uint64
	GetTxHash() TX_HASH
	GetBlockNumber() int64
}
