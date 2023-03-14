package types

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NEWTX, TX, TXATTEMPT will be converted from generic types to structs at a future date to enforce design and type checks
type TxStorageService[ADDR any, CHAINID any, HASH any, NEWTX any, R any, TX any, TXATTEMPT any, TXID any, TXMETA any] interface {
	CheckEthTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAINID, opts ...any) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAINID, opts ...any) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAINID, opts ...any) (count uint32, err error)
	CreateEthTransaction(newTx NEWTX, chainID CHAINID, opts ...any) (tx TX, err error)
	DeleteInProgressAttempt(attempt TXATTEMPT, opts ...any) error
	EthTransactions(offset, limit int, opts ...any) ([]TX, int, error)
	EthTransactionsWithAttempts(offset, limit int, opts ...any) ([]TX, int, error)
	// EthTxAttempts(offset, limit int, opts ...any) ([]TXATTEMPT, int, error)
	FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAINID, opts ...any) (receipts []ReceiptPlus[R], err error)
	// FindEthTxAttempt(hash HASH, opts ...any) (*TXATTEMPT, error)
	// FindEthTxAttemptConfirmedByEthTxIDs(ids []TXID, opts ...any) ([]TXATTEMPT, error)
	// FindEthTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAINID, opts ...any) (etxs []*TX, err error)
	// FindEthTxsRequiringResubmissionDueToInsufficientEth(address ADDR, chainID CHAINID, opts ...any) (etxs []TX, err error)
	// FindEtxAttemptsConfirmedMissingReceipt(chainID CHAINID, opts ...any) (attempts []TXATTEMPT, err error)
	// FindEthTxAttemptsByEthTxIDs(ids []TXID, opts ...any) ([]TXATTEMPT, error)
	// FindEthTxAttemptsRequiringReceiptFetch(chainID CHAINID, opts ...any) (attempts []TXATTEMPT, err error)
	// FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID CHAINID, address ADDR, opts ...any) (attempts []TXATTEMPT, err error)
	// FindEthTxByHash(hash HASH, opts ...any) (*TX, error)
	// FindEthTxWithAttempts(etxID TXID, opts ...any) (etx TX, err error)
	// FindEthTxWithNonce(fromAddress ADDR, nonce TXMETA, opts ...any) (etx *TX, err error)
	// FindNextUnstartedTransactionFromAddress(etx *TX, fromAddress ADDR, chainID CHAINID, opts ...any) error
	// FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID CHAINID, opts ...any) (etxs []*TX, err error)
	// GetEthTxInProgress(fromAddress ADDR, opts ...any) (etx *TX, err error)
	// GetInProgressEthTxAttempts(ctx context.Context, address ADDR, chainID CHAINID, opts ...any) (attempts []TXATTEMPT, err error)
	// HasInProgressTransaction(account ADDR, chainID CHAINID, opts ...any) (exists bool, err error)
	// // InsertEthReceipt only used in tests. Use SaveFetchedReceipts instead
	InsertEthReceipt(receipt *Receipt[R, HASH], opts ...any) error
	// InsertEthTx(etx *TX, opts ...any) error
	// InsertEthTxAttempt(attempt *TXATTEMPT, opts ...any) error
	// LoadEthTxAttempts(etx *TX, opts ...any) error
	// LoadEthTxesAttempts(etxs []*TX, opts ...any) error
	// MarkAllConfirmedMissingReceipt(chainID CHAINID, opts ...any) (err error)
	// MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAINID, opts ...any) error
	// PreloadEthTxes(attempts []TXATTEMPT, opts ...any) error
	PruneUnstartedEthTxQueue(queueSize uint32, subject uuid.UUID, opts ...any) (n int64, err error)
	// SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time, opts ...any) error
	SaveFetchedReceipts(receipts []R, chainID CHAINID, opts ...any) (err error)
	// SaveInProgressAttempt(attempt *TXATTEMPT, opts ...any) error
	// SaveInsufficientEthAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time, opts ...any) error
	// SaveReplacementInProgressAttempt(oldAttempt TXATTEMPT, replacementAttempt *TXATTEMPT, opts ...any) error
	// SaveSentAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time, opts ...any) error
	// SetBroadcastBeforeBlockNum(blockNum int64, chainID CHAINID, opts ...any) error
	// UpdateBroadcastAts(now time.Time, etxIDs []TXID, opts ...any) error
	// UpdateEthKeyNextNonce(newNextNonce, currentNextNonce TXMETA, address ADDR, chainID CHAINID, opts ...any) error
	// UpdateEthTxAttemptInProgressToBroadcast(etx *TX, attempt TXATTEMPT, NewAttemptState TxAttemptState, incrNextNonceCallback CallbackFunc, opts ...any) error
	// UpdateEthTxsUnconfirmed(ids []int64, opts ...any) error
	// UpdateEthTxUnstartedToInProgress(etx *TX, attempt *TXATTEMPT, opts ...any) error
	// UpdateEthTxFatalError(etx *TX, opts ...any) error
	// UpdateEthTxForRebroadcast(etx TX, etxAttempt TXATTEMPT, opts ...any) error
	// Close()
}

// TxStrategy controls how txes are queued and sent
//
//go:generate mockery --quiet --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved to eth_txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after eth_tx insertion
	PruneQueue(pruneService any, opt any) (n int64, err error)
}

// R is the raw transaction receipt
type ReceiptPlus[R any] struct {
	ID           uuid.UUID `db:"id"`
	Receipt      R         `db:"receipt"`
	FailOnRevert bool      `db:"FailOnRevert"`
}

type Receipt[R any, HASH any] struct {
	ID               int64
	TxHash           HASH
	BlockHash        HASH
	BlockNumber      int64
	TransactionIndex uint
	Receipt          R
	CreatedAt        time.Time
}

type TxAttemptState string

type TxState string

type CallbackFunc func(opts ...any) error
