package types

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// TX and TXATTEMPT will be converted from a generic type to a struct at a future date to enforce design and type checks
type ORM[ADDR any, CHAINID any, HASH any, R any, TX any, TXATTEMPT any, TXID any, TXMETA any] interface {
	CheckEthTxQueueCapacity(fromAddress ADDR, maxQueuedTransactions uint64, chainID CHAINID, opts ...any) (err error)
	CountUnconfirmedTransactions(fromAddress ADDR, chainID CHAINID, opts ...any) (count uint32, err error)
	CountUnstartedTransactions(fromAddress ADDR, chainID CHAINID, opts ...any) (count uint32, err error)
	// CreateEthTransaction(newTx any, chainID CHAINID, opts ...any) (tx TX, err error)
	// DeleteInProgressAttempt(ctx context.Context, attempt TXATTEMPT, opts ...any) error
	// EthTransactions(offset, limit int, opts ...any) ([]TX, int, error)
	// EthTransactionsWithAttempts(offset, limit int, opts ...any) ([]TX, int, error)
	// EthTxAttempts(offset, limit int, opts ...any) ([]TXATTEMPT, int, error)
	// FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID CHAINID, opts ...any) (receipts []ReceiptPlus[R], err error)
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
	// InsertEthReceipt(receipt *Receipt[R, HASH], opts ...any) error
	// InsertEthTx(etx *TX, opts ...any) error
	// InsertEthTxAttempt(attempt *TXATTEMPT, opts ...any) error
	// LoadEthTxAttempts(etx *TX, opts ...any) error
	// LoadEthTxesAttempts(etxs []*TX, opts ...any) error
	// MarkAllConfirmedMissingReceipt(chainID CHAINID, opts ...any) (err error)
	// MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID CHAINID, opts ...any) error
	// PreloadEthTxes(attempts []TXATTEMPT, opts ...any) error
	// PruneUnstartedEthTxQueue(queueSize uint32, subject uuid.UUID, opts ...any) (n int64, err error)
	// SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TXATTEMPT, broadcastAt time.Time, opts ...any) error
	// SaveFetchedReceipts(receipts []R, chainID CHAINID, opts ...any) (err error)
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

// type TxStrategy interface {
// 	// Subject will be saved to eth_txes.subject if not null
// 	Subject() uuid.NullUUID
// 	// PruneQueue is called after eth_tx insertion
// 	PruneQueue(orm ORM, q pg.Queryer) (n int64, err error)
// }

// type NewTx[ADDR any, TX any] struct {
// 	Tx TX
// 	ForwarderAddress ADDR

// 	Strategy txmgr.TxStrategy

// 	// Checker defines the check that should be run before a transaction is submitted on chain.
// 	Checker TransmitCheckerSpec
// }

// //go:generate mockery --quiet --name TxStrategy --output ../mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
// type TxStrategy interface {
// 	// Subject will be saved to eth_txes.subject if not null
// 	Subject() uuid.NullUUID
// 	// PruneQueue is called after eth_tx insertion
// 	PruneQueue(orm ORM, opt any) (n int64, err error)
// }

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

// const (
// 	TxAttemptInProgress = TxAttemptState("in_progress")
// 	// EthTxAttemptInsufficientEth = EthTxAttemptState("insufficient_eth")
// 	TxAttemptInsufficientFunds = TxAttemptState("insufficient_funds")
// 	TxAttemptBroadcast         = TxAttemptState("broadcast")
// )

// TXM type for tracking tx in TXM
// type EthTx struct {
// type Tx[ADD any, ADDR any, GASPRICE any, HASH any, ID any, META any, R any, TOKEN any] struct {
// 	ID int64
// 	//	Nonce          *int64
// 	Metadata       *META // this should be a number type for sorting
// 	FromAddress    ADDR
// 	ToAddress      ADDR
// 	EncodedPayload []byte
// 	// Value          assets.Eth
// 	Value TOKEN
// 	// GasLimit on the EthTx is always the conceptual gas limit, which is not
// 	// necessarily the same as the on-chain encoded value (i.e. Optimism)
// 	GasLimit uint32
// 	Error    null.String
// 	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
// 	// In almost all cases it will be within a second or so of the actual send time.
// 	BroadcastAt *time.Time
// 	// InitialBroadcastAt is recorded once, the first ever time this eth_tx is sent
// 	InitialBroadcastAt *time.Time
// 	CreatedAt          time.Time
// 	// State              EthTxState
// 	State TxState
// 	// EthTxAttempts      []EthTxAttempt `json:"-"`
// 	TxAttempts []TxAttempt[ADD, ADDR, GASPRICE, HASH, ID, META, R, TOKEN] `json:"-"`
// 	// Marshalled EthTxMeta
// 	// Used for additional context around transactions which you want to log
// 	// at send time.
// 	Meta    *datatypes.JSON
// 	Subject uuid.NullUUID
// 	// EVMChainID utils.Big
// 	ChainID ID

// 	PipelineTaskRunID uuid.NullUUID
// 	MinConfirmations  cnull.Uint32

// 	// AccessList is optional and only has an effect on DynamicFee transactions
// 	// on chains that support it (e.g. Ethereum Mainnet after London hard fork)
// 	// AccessList NullableEIP2930AccessList
// 	// flexible parameter for passing information like the EIP2930AccessList (or even a struct wrapping multiple parameters together)
// 	AdditionalTxData ADD // then passed to a chain specific tx builder

// 	// TransmitChecker defines the check that should be performed before a transaction is submitted on
// 	// chain.
// 	TransmitChecker *datatypes.JSON
// }

// type TxAttempt[ADD any, ADDR any, GASPRICE any, HASH any, ID any, META any, R any, TOKEN any] struct {
// 	ID int64
// 	// EthTxID int64
// 	TxID int64
// 	// EthTx   EthTx
// 	Tx Tx[ADD, ADDR, GASPRICE, HASH, ID, META, R, TOKEN]

// 	// GasPrice applies to LegacyTx
// 	// GasPrice *assets.Wei
// 	// GasTipCap and GasFeeCap are used instead for DynamicFeeTx
// 	// GasTipCap *assets.Wei // combine into a struct with FeeCap
// 	// GasFeeCap *assets.Wei // combine into a struct with TipCap
// 	GasPrice *GASPRICE

// 	// ChainSpecificGasLimit on the EthTxAttempt is always the same as the on-chain encoded value for gas limit
// 	ChainSpecificGasLimit   uint32
// 	SignedRawTx             []byte
// 	Hash                    HASH
// 	CreatedAt               time.Time
// 	BroadcastBeforeBlockNum *int64
// 	State                   TxAttemptState
// 	EthReceipts             []Receipt[R, HASH] `json:"-"`
// 	TxType                  int
// }

// type EthTxState string

// const (
// 	TxUnstarted               = TxState("unstarted")
// 	TxInProgress              = TxState("in_progress")
// 	TxFatalError              = TxState("fatal_error")
// 	TxUnconfirmed             = TxState("unconfirmed")
// 	TxConfirmed               = TxState("confirmed")
// 	TxConfirmedMissingReceipt = TxState("confirmed_missing_receipt")
// )
