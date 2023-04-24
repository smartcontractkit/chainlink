package types

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// NEWTX, TX, TXATTEMPT will be converted from generic types to structs at a future date to enforce design and type checks
//
//go:generate mockery --quiet --name TxStore --output ./mocks/ --case=underscore
type TxStore[
	// Represents an account address, in native chain format. TODO: Should implement Hashable
	ADDR types.Hashable,
	// Represents a chain id to be used for the chain.
	CHAIN_ID ID,
	// Represents a unique Tx Hash for a chain
	TX_HASH types.Hashable,
	// Represents a unique Block Hash for a chain
	BLOCK_HASH types.Hashable,
	NEWTX any,
	// Represents a onchain receipt object that a chain's RPC returns
	R any,
	// Represents a tx object that is used by the TXM
	// TODO: Remove https://smartcontract-it.atlassian.net/browse/BCI-865
	TX any,
	// Represents a tx attempt object that is used by the TXM
	// TODO: Remove https://smartcontract-it.atlassian.net/browse/BCI-865
	TXATTEMPT any,
	// Represents the sequence type for a chain. For example, nonce for EVM.
	SEQ Sequence,
] interface {
	UnstartedTxQueuePruner
	CheckEthTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAIN_ID, qopts ...pg.QOpt) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (count uint32, err error)
	CreateEthTransaction(newTx NEWTX, chainID CHAIN_ID, qopts ...pg.QOpt) (tx Transaction, err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TXATTEMPT) error
	EthTransactions(offset, limit int) ([]TX, int, error)
	EthTransactionsWithAttempts(offset, limit int) ([]TX, int, error)
	EthTxAttempts(offset, limit int) ([]TXATTEMPT, int, error)
	FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAIN_ID) (receiptsPlus []ReceiptPlus[R], err error)
	FindEthTxAttempt(hash TX_HASH) (*TXATTEMPT, error)
	FindEthTxAttemptConfirmedByEthTxIDs(ids []int64) ([]TXATTEMPT, error)
	FindEthTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAIN_ID) (etxs []*TX, err error)
	FindEthTxsRequiringResubmissionDueToInsufficientEth(address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (etxs []*TX, err error)
	FindEtxAttemptsConfirmedMissingReceipt(chainID CHAIN_ID) (attempts []TXATTEMPT, err error)
	FindEthTxAttemptsByEthTxIDs(ids []int64) ([]TXATTEMPT, error)
	FindEthTxAttemptsRequiringReceiptFetch(chainID CHAIN_ID) (attempts []TXATTEMPT, err error)
	FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID CHAIN_ID, address ADDR) (attempts []TXATTEMPT, err error)
	FindEthTxByHash(hash TX_HASH) (*TX, error)
	FindEthTxWithAttempts(etxID int64) (etx TX, err error)
	FindEthTxWithNonce(fromAddress ADDR, seq SEQ) (etx *TX, err error)
	FindNextUnstartedTransactionFromAddress(etx *TX, fromAddress ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID CHAIN_ID) (etxs []*TX, err error)
	GetEthTxInProgress(fromAddress ADDR, qopts ...pg.QOpt) (etx *TX, err error)
	GetInProgressEthTxAttempts(ctx context.Context, address ADDR, chainID CHAIN_ID) (attempts []TXATTEMPT, err error)
	HasInProgressTransaction(account ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (exists bool, err error)
	// InsertEthReceipt only used in tests. Use SaveFetchedReceipts instead
	InsertEthReceipt(receipt *Receipt[R, TX_HASH, BLOCK_HASH]) error
	InsertEthTx(etx *TX) error
	InsertEthTxAttempt(attempt *TXATTEMPT) error
	LoadEthTxAttempts(etx *TX, qopts ...pg.QOpt) error
	LoadEthTxesAttempts(etxs []*TX, qopts ...pg.QOpt) error
	MarkAllConfirmedMissingReceipt(chainID CHAIN_ID) (err error)
	MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAIN_ID, qopts ...pg.QOpt) error
	PreloadEthTxes(attempts []TXATTEMPT, qopts ...pg.QOpt) error
	SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SaveFetchedReceipts(receipts []R, chainID CHAIN_ID) (err error)
	SaveInProgressAttempt(attempt *TXATTEMPT) error
	SaveInsufficientEthAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(oldAttempt TXATTEMPT, replacementAttempt *TXATTEMPT, qopts ...pg.QOpt) error
	SaveSentAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(blockNum int64, chainID CHAIN_ID) error
	UpdateBroadcastAts(now time.Time, etxIDs []int64) error
	UpdateEthKeyNextNonce(newNextNonce, currentNextNonce SEQ, address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) error
	UpdateEthTxAttemptInProgressToBroadcast(etx *TX, attempt TXATTEMPT, NewAttemptState TxAttemptState, incrNextNonceCallback QueryerFunc, qopts ...pg.QOpt) error
	UpdateEthTxsUnconfirmed(ids []int64) error
	UpdateEthTxUnstartedToInProgress(etx *TX, attempt *TXATTEMPT, qopts ...pg.QOpt) error
	UpdateEthTxFatalError(etx *TX, qopts ...pg.QOpt) error
	UpdateEthTxForRebroadcast(etx TX, etxAttempt TXATTEMPT) error
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
