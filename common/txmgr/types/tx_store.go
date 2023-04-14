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
	ADDR any,
	// Represents a chain id to be used for the chain.
	CHAINID any,
	// Represents a unique Tx Hash for a chain
	TX_HASH types.Hashable[TX_HASH],
	// Represents a unique Tx Hash for a chain
	BLOCK_HASH types.Hashable[BLOCK_HASH],
	NEWTX any,
	// Represents a onchain receipt object that a chain's RPC returns
	R any,
	// Represents a tx object that is used by the TXM
	// TODO: Remove https://smartcontract-it.atlassian.net/browse/BCI-865
	TX any,
	// Represents a tx attempt object that is used by the TXM
	// TODO: Remove https://smartcontract-it.atlassian.net/browse/BCI-865
	TXATTEMPT any,
	// Represents a unique tx id for the Chainlink Node.
	// TODO: Remove https://smartcontract-it.atlassian.net/browse/BCI-865
	TXID any,
	// Represents the sequence type for a chain. For example, nonce for EVM.
	TXMETA any,
] interface {
	UnstartedTxQueuePruner
	CheckEthTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAINID, qopts ...pg.QOpt) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAINID, qopts ...pg.QOpt) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAINID, qopts ...pg.QOpt) (count uint32, err error)
	CreateEthTransaction(newTx NEWTX, chainID CHAINID, qopts ...pg.QOpt) (tx Transaction, err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TXATTEMPT) error
	EthTransactions(offset, limit int) ([]TX, int, error)
	EthTransactionsWithAttempts(offset, limit int) ([]TX, int, error)
	EthTxAttempts(offset, limit int) ([]TXATTEMPT, int, error)
	FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAINID) (receiptsPlus []ReceiptPlus[R], err error)
	FindEthTxAttempt(hash TX_HASH) (*TXATTEMPT, error)
	FindEthTxAttemptConfirmedByEthTxIDs(ids []TXID) ([]TXATTEMPT, error)
	FindEthTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CHAINID) (etxs []*TX, err error)
	FindEthTxsRequiringResubmissionDueToInsufficientEth(address ADDR, chainID CHAINID, qopts ...pg.QOpt) (etxs []*TX, err error)
	FindEtxAttemptsConfirmedMissingReceipt(chainID CHAINID) (attempts []TXATTEMPT, err error)
	FindEthTxAttemptsByEthTxIDs(ids []TXID) ([]TXATTEMPT, error)
	FindEthTxAttemptsRequiringReceiptFetch(chainID CHAINID) (attempts []TXATTEMPT, err error)
	FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID CHAINID, address ADDR) (attempts []TXATTEMPT, err error)
	FindEthTxByHash(hash TX_HASH) (*TX, error)
	FindEthTxWithAttempts(etxID TXID) (etx TX, err error)
	FindEthTxWithNonce(fromAddress ADDR, nonce TXMETA) (etx *TX, err error)
	FindNextUnstartedTransactionFromAddress(etx *TX, fromAddress ADDR, chainID CHAINID, qopts ...pg.QOpt) error
	FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID CHAINID) (etxs []*TX, err error)
	GetEthTxInProgress(fromAddress ADDR, qopts ...pg.QOpt) (etx *TX, err error)
	GetInProgressEthTxAttempts(ctx context.Context, address ADDR, chainID CHAINID) (attempts []TXATTEMPT, err error)
	HasInProgressTransaction(account ADDR, chainID CHAINID, qopts ...pg.QOpt) (exists bool, err error)
	// InsertEthReceipt only used in tests. Use SaveFetchedReceipts instead
	InsertEthReceipt(receipt *Receipt[R, TX_HASH, BLOCK_HASH]) error
	InsertEthTx(etx *TX) error
	InsertEthTxAttempt(attempt *TXATTEMPT) error
	LoadEthTxAttempts(etx *TX, qopts ...pg.QOpt) error
	LoadEthTxesAttempts(etxs []*TX, qopts ...pg.QOpt) error
	MarkAllConfirmedMissingReceipt(chainID CHAINID) (err error)
	MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAINID, qopts ...pg.QOpt) error
	PreloadEthTxes(attempts []TXATTEMPT, qopts ...pg.QOpt) error
	SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SaveFetchedReceipts(receipts []R, chainID CHAINID) (err error)
	SaveInProgressAttempt(attempt *TXATTEMPT) error
	SaveInsufficientEthAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(oldAttempt TXATTEMPT, replacementAttempt *TXATTEMPT, qopts ...pg.QOpt) error
	SaveSentAttempt(timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(blockNum int64, chainID CHAINID) error
	UpdateBroadcastAts(now time.Time, etxIDs []TXID) error
	UpdateEthKeyNextNonce(newNextNonce, currentNextNonce TXMETA, address ADDR, chainID CHAINID, qopts ...pg.QOpt) error
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
type Receipt[R any, TX_HASH types.Hashable[TX_HASH], BLOCK_HASH types.Hashable[BLOCK_HASH]] struct {
	ID               int64
	TxHash           TX_HASH
	BlockHash        BLOCK_HASH
	BlockNumber      int64
	TransactionIndex uint
	Receipt          R
	CreatedAt        time.Time
}

type QueryerFunc = func(tx pg.Queryer) error
