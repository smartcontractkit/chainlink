package txmgr

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/null"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

var (
	ErrKeyNotUpdated = errors.New("evmTxStore: Key not updated")
	// ErrCouldNotGetReceipt is the error string we save if we reach our finality depth for a confirmed transaction without ever getting a receipt
	// This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = "could not get receipt"
)

// EvmTxStore combines the txmgr tx store interface and the interface needed for the API to read from the tx DB
//
//go:generate mockery --quiet --name EvmTxStore --output ./mocks/ --case=underscore
type EvmTxStore interface {
	// redeclare TxStore for mockery
	txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	TxStoreWebApi
}

// TxStoreWebApi encapsulates the methods that are not used by the txmgr and only used by the various web controllers and readers
type TxStoreWebApi interface {
	FindTxAttemptConfirmedByTxIDs(ctx context.Context, ids []int64) ([]TxAttempt, error)
	FindTxByHash(ctx context.Context, hash common.Hash) (*Tx, error)
	Transactions(ctx context.Context, offset, limit int) ([]Tx, int, error)
	TxAttempts(ctx context.Context, offset, limit int) ([]TxAttempt, int, error)
	TransactionsWithAttempts(ctx context.Context, offset, limit int) ([]Tx, int, error)
	FindTxAttempt(ctx context.Context, hash common.Hash) (*TxAttempt, error)
	FindTxWithAttempts(ctx context.Context, etxID int64) (etx Tx, err error)
}

type TestEvmTxStore interface {
	EvmTxStore

	// methods only used for testing purposes
	InsertReceipt(ctx context.Context, receipt *evmtypes.Receipt) (int64, error)
	InsertTx(ctx context.Context, etx *Tx) error
	FindTxAttemptsByTxIDs(ctx context.Context, ids []int64) ([]TxAttempt, error)
	InsertTxAttempt(ctx context.Context, attempt *TxAttempt) error
	LoadTxesAttempts(ctx context.Context, etxs []*Tx) error
	GetFatalTransactions(ctx context.Context) (txes []*Tx, err error)
	GetAllTxes(ctx context.Context) (txes []*Tx, err error)
	GetAllTxAttempts(ctx context.Context) (attempts []TxAttempt, err error)
	CountTxesByStateAndSubject(ctx context.Context, state txmgrtypes.TxState, subject uuid.UUID) (count int, err error)
	FindTxesByFromAddressAndState(ctx context.Context, fromAddress common.Address, state string) (txes []*Tx, err error)
	UpdateTxAttemptBroadcastBeforeBlockNum(ctx context.Context, id int64, blockNum uint) error
}

type evmTxStore struct {
	q         sqlutil.DataSource
	logger    logger.SugaredLogger
	ctx       context.Context
	ctxCancel context.CancelFunc
}

var _ EvmTxStore = (*evmTxStore)(nil)
var _ TestEvmTxStore = (*evmTxStore)(nil)

// Directly maps to columns of database table "evm.receipts".
// Do not modify type unless you
// intend to modify the database schema
type dbReceipt struct {
	ID               int64
	TxHash           common.Hash
	BlockHash        common.Hash
	BlockNumber      int64
	TransactionIndex uint
	Receipt          evmtypes.Receipt
	CreatedAt        time.Time
}

func DbReceiptFromEvmReceipt(evmReceipt *evmtypes.Receipt) dbReceipt {
	return dbReceipt{
		TxHash:           evmReceipt.TxHash,
		BlockHash:        evmReceipt.BlockHash,
		BlockNumber:      evmReceipt.BlockNumber.Int64(),
		TransactionIndex: evmReceipt.TransactionIndex,
		Receipt:          *evmReceipt,
	}
}

func DbReceiptToEvmReceipt(receipt *dbReceipt) *evmtypes.Receipt {
	return &receipt.Receipt
}

// Directly maps to onchain receipt schema.
type rawOnchainReceipt = evmtypes.Receipt

func (o *evmTxStore) Transact(ctx context.Context, readOnly bool, fn func(*evmTxStore) error) (err error) {
	opts := &sqlutil.TxOptions{TxOptions: sql.TxOptions{ReadOnly: readOnly}}
	return sqlutil.Transact(ctx, o.new, o.q, opts, fn)
}

// new returns a NewORM like o, but backed by q.
func (o *evmTxStore) new(q sqlutil.DataSource) *evmTxStore { return NewTxStore(q, o.logger) }

// Directly maps to some columns of few database tables.
// Does not map to a single database table.
// It's comprised of fields from different tables.
type dbReceiptPlus struct {
	ID           uuid.UUID        `db:"pipeline_task_run_id"`
	Receipt      evmtypes.Receipt `db:"receipt"`
	FailOnRevert bool             `db:"FailOnRevert"`
}

func fromDBReceipts(rs []dbReceipt) []*evmtypes.Receipt {
	receipts := make([]*evmtypes.Receipt, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = DbReceiptToEvmReceipt(&rs[i])
	}
	return receipts
}

func fromDBReceiptsPlus(rs []dbReceiptPlus) []ReceiptPlus {
	receipts := make([]ReceiptPlus, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = ReceiptPlus{
			ID:           rs[i].ID,
			Receipt:      &rs[i].Receipt,
			FailOnRevert: rs[i].FailOnRevert,
		}
	}
	return receipts
}

func toOnchainReceipt(rs []*evmtypes.Receipt) []rawOnchainReceipt {
	receipts := make([]rawOnchainReceipt, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = *rs[i]
	}
	return receipts
}

// Directly maps to columns of database table "evm.txes".
// This is exported, as tests and other external code still directly reads DB using this schema.
type DbEthTx struct {
	ID             int64
	IdempotencyKey *string
	Nonce          *int64
	FromAddress    common.Address
	ToAddress      common.Address
	EncodedPayload []byte
	Value          assets.Eth
	// GasLimit on the EthTx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	GasLimit uint64
	Error    nullv4.String
	// BroadcastAt is updated every time an attempt for this eth_tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt *time.Time
	// InitialBroadcastAt is recorded once, the first ever time this eth_tx is sent
	CreatedAt time.Time
	State     txmgrtypes.TxState
	// Marshalled EvmTxMeta
	// Used for additional context around transactions which you want to log
	// at send time.
	Meta              *sqlutil.JSON
	Subject           uuid.NullUUID
	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  null.Uint32
	EVMChainID        ubig.Big
	// TransmitChecker defines the check that should be performed before a transaction is submitted on
	// chain.
	TransmitChecker    *sqlutil.JSON
	InitialBroadcastAt *time.Time
	// Marks tx requiring callback
	SignalCallback bool
	// Marks tx callback as signaled
	CallbackCompleted bool
}

func (db *DbEthTx) FromTx(tx *Tx) {
	db.ID = tx.ID
	db.IdempotencyKey = tx.IdempotencyKey
	db.FromAddress = tx.FromAddress
	db.ToAddress = tx.ToAddress
	db.EncodedPayload = tx.EncodedPayload
	db.Value = assets.Eth(tx.Value)
	db.GasLimit = tx.FeeLimit
	db.Error = tx.Error
	db.BroadcastAt = tx.BroadcastAt
	db.CreatedAt = tx.CreatedAt
	db.State = tx.State
	db.Meta = tx.Meta
	db.Subject = tx.Subject
	db.PipelineTaskRunID = tx.PipelineTaskRunID
	db.MinConfirmations = tx.MinConfirmations
	db.TransmitChecker = tx.TransmitChecker
	db.InitialBroadcastAt = tx.InitialBroadcastAt
	db.SignalCallback = tx.SignalCallback
	db.CallbackCompleted = tx.CallbackCompleted

	if tx.ChainID != nil {
		db.EVMChainID = *ubig.New(tx.ChainID)
	}
	if tx.Sequence != nil {
		n := tx.Sequence.Int64()
		db.Nonce = &n
	}
}

func (db DbEthTx) ToTx(tx *Tx) {
	tx.ID = db.ID
	if db.Nonce != nil {
		n := evmtypes.Nonce(*db.Nonce)
		tx.Sequence = &n
	}
	tx.IdempotencyKey = db.IdempotencyKey
	tx.FromAddress = db.FromAddress
	tx.ToAddress = db.ToAddress
	tx.EncodedPayload = db.EncodedPayload
	tx.Value = *db.Value.ToInt()
	tx.FeeLimit = db.GasLimit
	tx.Error = db.Error
	tx.BroadcastAt = db.BroadcastAt
	tx.CreatedAt = db.CreatedAt
	tx.State = db.State
	tx.Meta = db.Meta
	tx.Subject = db.Subject
	tx.PipelineTaskRunID = db.PipelineTaskRunID
	tx.MinConfirmations = db.MinConfirmations
	tx.ChainID = db.EVMChainID.ToInt()
	tx.TransmitChecker = db.TransmitChecker
	tx.InitialBroadcastAt = db.InitialBroadcastAt
	tx.SignalCallback = db.SignalCallback
	tx.CallbackCompleted = db.CallbackCompleted
}

func dbEthTxsToEvmEthTxs(dbEthTxs []DbEthTx) []Tx {
	evmEthTxs := make([]Tx, len(dbEthTxs))
	for i, dbTx := range dbEthTxs {
		dbTx.ToTx(&evmEthTxs[i])
	}
	return evmEthTxs
}

func dbEthTxsToEvmEthTxPtrs(dbEthTxs []DbEthTx, evmEthTxs []*Tx) {
	for i, dbTx := range dbEthTxs {
		evmEthTxs[i] = &Tx{}
		dbTx.ToTx(evmEthTxs[i])
	}
}

// Directly maps to columns of database table "evm.tx_attempts".
// This is exported, as tests and other external code still directly reads DB using this schema.
type DbEthTxAttempt struct {
	ID                      int64
	EthTxID                 int64
	GasPrice                *assets.Wei
	SignedRawTx             []byte
	Hash                    common.Hash
	BroadcastBeforeBlockNum *int64
	State                   string
	CreatedAt               time.Time
	ChainSpecificGasLimit   uint64
	TxType                  int
	GasTipCap               *assets.Wei
	GasFeeCap               *assets.Wei
}

func (db *DbEthTxAttempt) FromTxAttempt(attempt *TxAttempt) {
	db.ID = attempt.ID
	db.EthTxID = attempt.TxID
	db.GasPrice = attempt.TxFee.Legacy
	db.SignedRawTx = attempt.SignedRawTx
	db.Hash = attempt.Hash
	db.BroadcastBeforeBlockNum = attempt.BroadcastBeforeBlockNum
	db.CreatedAt = attempt.CreatedAt
	db.ChainSpecificGasLimit = attempt.ChainSpecificFeeLimit
	db.TxType = attempt.TxType
	db.GasTipCap = attempt.TxFee.DynamicTipCap
	db.GasFeeCap = attempt.TxFee.DynamicFeeCap

	// handle state naming difference between generic + EVM
	if attempt.State == txmgrtypes.TxAttemptInsufficientFunds {
		db.State = "insufficient_eth"
	} else {
		db.State = attempt.State.String()
	}
}

func DbEthTxAttemptStateToTxAttemptState(state string) txmgrtypes.TxAttemptState {
	if state == "insufficient_eth" {
		return txmgrtypes.TxAttemptInsufficientFunds
	}
	return txmgrtypes.NewTxAttemptState(state)
}

func (db DbEthTxAttempt) ToTxAttempt(attempt *TxAttempt) {
	attempt.ID = db.ID
	attempt.TxID = db.EthTxID
	attempt.SignedRawTx = db.SignedRawTx
	attempt.Hash = db.Hash
	attempt.BroadcastBeforeBlockNum = db.BroadcastBeforeBlockNum
	attempt.State = DbEthTxAttemptStateToTxAttemptState(db.State)
	attempt.CreatedAt = db.CreatedAt
	attempt.ChainSpecificFeeLimit = db.ChainSpecificGasLimit
	attempt.TxType = db.TxType
	attempt.TxFee = gas.EvmFee{
		Legacy:        db.GasPrice,
		DynamicTipCap: db.GasTipCap,
		DynamicFeeCap: db.GasFeeCap,
	}
}

func dbEthTxAttemptsToEthTxAttempts(dbEthTxAttempt []DbEthTxAttempt) []TxAttempt {
	evmEthTxAttempt := make([]TxAttempt, len(dbEthTxAttempt))
	for i, dbTxAttempt := range dbEthTxAttempt {
		dbTxAttempt.ToTxAttempt(&evmEthTxAttempt[i])
	}
	return evmEthTxAttempt
}

func NewTxStore(
	db sqlutil.DataSource,
	lggr logger.Logger,
) *evmTxStore {
	namedLogger := logger.Named(lggr, "TxmStore")
	ctx, cancel := context.WithCancel(context.Background())
	return &evmTxStore{
		q:         db,
		logger:    logger.Sugared(namedLogger),
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

const insertIntoEthTxAttemptsQuery = `
INSERT INTO evm.tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap)
VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap)
RETURNING *;
`

func (o *evmTxStore) Close() {
	o.ctxCancel()
}

func (o *evmTxStore) preloadTxAttempts(ctx context.Context, txs []Tx) error {
	// Preload TxAttempts
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	var dbAttempts []DbEthTxAttempt
	sql := `SELECT * FROM evm.tx_attempts WHERE eth_tx_id IN (?) ORDER BY id desc;`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = o.q.Rebind(query)
	if err = o.q.SelectContext(ctx, &dbAttempts, query, args...); err != nil {
		return err
	}
	// fill in attempts
	for _, dbAttempt := range dbAttempts {
		for i, tx := range txs {
			if tx.ID == dbAttempt.EthTxID {
				var attempt TxAttempt
				dbAttempt.ToTxAttempt(&attempt)
				txs[i].TxAttempts = append(txs[i].TxAttempts, attempt)
			}
		}
	}
	return nil
}

func (o *evmTxStore) PreloadTxes(ctx context.Context, attempts []TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	return o.preloadTxesAtomic(ctx, attempts)
}

// Only to be used for atomic transactions internal to the tx store
func (o *evmTxStore) preloadTxesAtomic(ctx context.Context, attempts []TxAttempt) error {
	ethTxM := make(map[int64]Tx)
	for _, attempt := range attempts {
		ethTxM[attempt.TxID] = Tx{}
	}
	ethTxIDs := make([]int64, len(ethTxM))
	var i int
	for id := range ethTxM {
		ethTxIDs[i] = id
		i++
	}
	dbEthTxs := make([]DbEthTx, len(ethTxIDs))
	if err := o.q.SelectContext(ctx, &dbEthTxs, `SELECT * FROM evm.txes WHERE id = ANY($1)`, pq.Array(ethTxIDs)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxes failed")
	}
	for _, dbEtx := range dbEthTxs {
		etx := ethTxM[dbEtx.ID]
		dbEtx.ToTx(&etx)
		ethTxM[etx.ID] = etx
	}
	for i, attempt := range attempts {
		attempts[i].Tx = ethTxM[attempt.TxID]
	}
	return nil
}

// Transactions returns all eth transactions without loaded relations
// limited by passed parameters.
func (o *evmTxStore) Transactions(ctx context.Context, offset, limit int) (txs []Tx, count int, err error) {
	sql := `SELECT count(*) FROM evm.txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM evm.tx_attempts)`
	if err = o.q.GetContext(ctx, &count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM evm.txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM evm.tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	var dbEthTxs []DbEthTx
	if err = o.q.SelectContext(ctx, &dbEthTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxsToEvmEthTxs(dbEthTxs)
	return
}

// TransactionsWithAttempts returns all eth transactions with at least one attempt
// limited by passed parameters. Attempts are sorted by id.
func (o *evmTxStore) TransactionsWithAttempts(ctx context.Context, offset, limit int) (txs []Tx, count int, err error) {
	sql := `SELECT count(*) FROM evm.txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM evm.tx_attempts)`
	if err = o.q.GetContext(ctx, &count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM evm.txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM evm.tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	var dbTxs []DbEthTx
	if err = o.q.SelectContext(ctx, &dbTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxsToEvmEthTxs(dbTxs)
	err = o.preloadTxAttempts(ctx, txs)
	return
}

// TxAttempts returns the last tx attempts sorted by created_at descending.
func (o *evmTxStore) TxAttempts(ctx context.Context, offset, limit int) (txs []TxAttempt, count int, err error) {
	sql := `SELECT count(*) FROM evm.tx_attempts`
	if err = o.q.GetContext(ctx, &count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM evm.tx_attempts ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	var dbTxs []DbEthTxAttempt
	if err = o.q.SelectContext(ctx, &dbTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxAttemptsToEthTxAttempts(dbTxs)
	err = o.preloadTxesAtomic(ctx, txs)
	return
}

// FindTxAttempt returns an individual TxAttempt
func (o *evmTxStore) FindTxAttempt(ctx context.Context, hash common.Hash) (*TxAttempt, error) {
	dbTxAttempt := DbEthTxAttempt{}
	sql := `SELECT * FROM evm.tx_attempts WHERE hash = $1`
	if err := o.q.GetContext(ctx, &dbTxAttempt, sql, hash); err != nil {
		return nil, err
	}
	// reuse the preload
	var attempt TxAttempt
	dbTxAttempt.ToTxAttempt(&attempt)
	attempts := []TxAttempt{attempt}
	err := o.preloadTxesAtomic(ctx, attempts)
	return &attempts[0], err
}

// FindTxAttemptsByTxIDs returns a list of attempts by ETH Tx IDs
func (o *evmTxStore) FindTxAttemptsByTxIDs(ctx context.Context, ids []int64) ([]TxAttempt, error) {
	sql := `SELECT * FROM evm.tx_attempts WHERE eth_tx_id = ANY($1)`
	var dbTxAttempts []DbEthTxAttempt
	if err := o.q.SelectContext(ctx, &dbTxAttempts, sql, ids); err != nil {
		return nil, err
	}
	return dbEthTxAttemptsToEthTxAttempts(dbTxAttempts), nil
}

func (o *evmTxStore) FindTxByHash(ctx context.Context, hash common.Hash) (*Tx, error) {
	var dbEtx DbEthTx
	err := o.Transact(ctx, true, func(orm *evmTxStore) error {
		sql := `SELECT evm.txes.* FROM evm.txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM evm.tx_attempts WHERE hash = $1)`
		if err := orm.q.GetContext(ctx, &dbEtx, sql, hash); err != nil {
			return pkgerrors.Wrapf(err, "failed to find eth_tx with hash %d", hash)
		}
		return nil
	})

	var etx Tx
	dbEtx.ToTx(&etx)
	return &etx, pkgerrors.Wrap(err, "FindEthTxByHash failed")
}

// InsertTx inserts a new evm tx into the database
func (o *evmTxStore) InsertTx(ctx context.Context, etx *Tx) error {
	if etx.CreatedAt == (time.Time{}) {
		etx.CreatedAt = time.Now()
	}
	const insertEthTxSQL = `INSERT INTO evm.txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, initial_broadcast_at, created_at, state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, transmit_checker, idempotency_key, signal_callback, callback_completed) VALUES (
:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :error, :broadcast_at, :initial_broadcast_at, :created_at, :state, :meta, :subject, :pipeline_task_run_id, :min_confirmations, :evm_chain_id, :transmit_checker, :idempotency_key, :signal_callback, :callback_completed
) RETURNING *`
	var dbTx DbEthTx
	dbTx.FromTx(etx)

	query, args, err := o.q.BindNamed(insertEthTxSQL, &dbTx)
	if err != nil {
		return pkgerrors.Wrap(err, "InsertTx failed to bind named")
	}
	err = o.q.GetContext(ctx, &dbTx, query, args...)
	dbTx.ToTx(etx)
	return pkgerrors.Wrap(err, "InsertTx failed")
}

// InsertTxAttempt inserts a new txAttempt into the database
func (o *evmTxStore) InsertTxAttempt(ctx context.Context, attempt *TxAttempt) error {
	var dbTxAttempt DbEthTxAttempt
	dbTxAttempt.FromTxAttempt(attempt)
	query, args, err := o.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbTxAttempt)
	if err != nil {
		return pkgerrors.Wrap(err, "InsertTxAttempt failed to bind named")
	}
	err = o.q.GetContext(ctx, &dbTxAttempt, query, args...)
	dbTxAttempt.ToTxAttempt(attempt)
	return pkgerrors.Wrap(err, "InsertTxAttempt failed")
}

// InsertReceipt only used in tests. Use SaveFetchedReceipts instead
func (o *evmTxStore) InsertReceipt(ctx context.Context, receipt *evmtypes.Receipt) (int64, error) {
	// convert to database representation
	r := DbReceiptFromEvmReceipt(receipt)

	const insertEthReceiptSQL = `INSERT INTO evm.receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at) VALUES (
:tx_hash, :block_hash, :block_number, :transaction_index, :receipt, NOW()
) RETURNING *`
	query, args, err := o.q.BindNamed(insertEthReceiptSQL, &r)
	if err != nil {
		return 0, pkgerrors.Wrap(err, "InsertReceipt failed to bind named")
	}
	err = o.q.GetContext(ctx, &r, query, args...)
	return r.ID, pkgerrors.Wrap(err, "InsertReceipt failed")
}

func (o *evmTxStore) GetFatalTransactions(ctx context.Context) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		stmt := `SELECT * FROM evm.txes WHERE state = 'fatal_error'`
		var dbEtxs []DbEthTx
		if err = orm.q.SelectContext(ctx, &dbEtxs, stmt); err != nil {
			return fmt.Errorf("failed to load evm.txes: %w", err)
		}
		txes = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
		err = orm.LoadTxesAttempts(ctx, txes)
		if err != nil {
			return fmt.Errorf("failed to load evm.tx_attempts: %w", err)
		}
		return nil
	})

	return txes, nil
}

// FindTxWithAttempts finds the Tx with its attempts and receipts preloaded
func (o *evmTxStore) FindTxWithAttempts(ctx context.Context, etxID int64) (etx Tx, err error) {
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtx DbEthTx
		if err = orm.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE id = $1 ORDER BY created_at ASC, id ASC`, etxID); err != nil {
			return pkgerrors.Wrapf(err, "failed to find evm.tx with id %d", etxID)
		}
		dbEtx.ToTx(&etx)
		if err = orm.loadTxAttemptsAtomic(ctx, &etx); err != nil {
			return pkgerrors.Wrapf(err, "failed to load evm.tx_attempts for evm.tx with id %d", etxID)
		}
		if err = orm.loadEthTxAttemptsReceipts(ctx, &etx); err != nil {
			return pkgerrors.Wrapf(err, "failed to load evm.receipts for evm.tx with id %d", etxID)
		}
		return nil
	})
	return etx, pkgerrors.Wrap(err, "FindTxWithAttempts failed")
}

func (o *evmTxStore) FindTxAttemptConfirmedByTxIDs(ctx context.Context, ids []int64) ([]TxAttempt, error) {
	var txAttempts []TxAttempt
	err := o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbAttempts []DbEthTxAttempt
		if err := orm.q.SelectContext(ctx, &dbAttempts, `SELECT eta.*
		FROM evm.tx_attempts eta
			join evm.receipts er on eta.hash = er.tx_hash where eta.eth_tx_id = ANY($1) ORDER BY eta.gas_price DESC, eta.gas_tip_cap DESC`, ids); err != nil {
			return err
		}
		txAttempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		return loadConfirmedAttemptsReceipts(ctx, orm.q, txAttempts)
	})
	return txAttempts, pkgerrors.Wrap(err, "FindTxAttemptConfirmedByTxIDs failed")
}

// Only used internally for atomic transactions
func (o *evmTxStore) LoadTxesAttempts(ctx context.Context, etxs []*Tx) error {
	ethTxIDs := make([]int64, len(etxs))
	ethTxesM := make(map[int64]*Tx, len(etxs))
	for i, etx := range etxs {
		etx.TxAttempts = nil // this will overwrite any previous preload
		ethTxIDs[i] = etx.ID
		ethTxesM[etx.ID] = etxs[i]
	}
	var dbTxAttempts []DbEthTxAttempt
	if err := o.q.SelectContext(ctx, &dbTxAttempts, `SELECT * FROM evm.tx_attempts WHERE eth_tx_id = ANY($1) ORDER BY evm.tx_attempts.gas_price DESC, evm.tx_attempts.gas_tip_cap DESC`, pq.Array(ethTxIDs)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxesAttempts failed to load evm.tx_attempts")
	}
	for _, dbAttempt := range dbTxAttempts {
		etx := ethTxesM[dbAttempt.EthTxID]
		var attempt TxAttempt
		dbAttempt.ToTxAttempt(&attempt)
		etx.TxAttempts = append(etx.TxAttempts, attempt)
	}
	return nil
}

func (o *evmTxStore) LoadTxAttempts(ctx context.Context, etx *Tx) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	return o.loadTxAttemptsAtomic(ctx, etx)
}

// Only to be used for atomic transactions internal to the tx store
func (o *evmTxStore) loadTxAttemptsAtomic(ctx context.Context, etx *Tx) error {
	return o.LoadTxesAttempts(ctx, []*Tx{etx})
}

func (o *evmTxStore) loadEthTxAttemptsReceipts(ctx context.Context, etx *Tx) (err error) {
	return o.loadEthTxesAttemptsReceipts(ctx, []*Tx{etx})
}

func (o *evmTxStore) loadEthTxesAttemptsReceipts(ctx context.Context, etxs []*Tx) (err error) {
	if len(etxs) == 0 {
		return nil
	}
	attemptHashM := make(map[common.Hash]*TxAttempt, len(etxs)) // len here is lower bound
	attemptHashes := make([][]byte, len(etxs))                  // len here is lower bound
	for _, etx := range etxs {
		for i, attempt := range etx.TxAttempts {
			attemptHashM[attempt.Hash] = &etx.TxAttempts[i]
			attemptHashes = append(attemptHashes, attempt.Hash.Bytes())
		}
	}
	var rs []dbReceipt
	if err = o.q.SelectContext(ctx, &rs, `SELECT * FROM evm.receipts WHERE tx_hash = ANY($1)`, pq.Array(attemptHashes)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxesAttemptsReceipts failed to load evm.receipts")
	}

	var receipts []*evmtypes.Receipt = fromDBReceipts(rs)

	for _, receipt := range receipts {
		attempt := attemptHashM[receipt.TxHash]
		// Although the attempts struct supports multiple receipts, the expectation for EVM is that there is only one receipt
		// per tx and therefore attempt too.
		attempt.Receipts = append(attempt.Receipts, receipt)
	}
	return nil
}

func loadConfirmedAttemptsReceipts(ctx context.Context, q sqlutil.DataSource, attempts []TxAttempt) error {
	byHash := make(map[string]*TxAttempt, len(attempts))
	hashes := make([][]byte, len(attempts))
	for i, attempt := range attempts {
		byHash[attempt.Hash.String()] = &attempts[i]
		hashes = append(hashes, attempt.Hash.Bytes())
	}
	var rs []dbReceipt
	if err := q.SelectContext(ctx, &rs, `SELECT * FROM evm.receipts WHERE tx_hash = ANY($1)`, pq.Array(hashes)); err != nil {
		return pkgerrors.Wrap(err, "loadConfirmedAttemptsReceipts failed to load evm.receipts")
	}
	var receipts []*evmtypes.Receipt = fromDBReceipts(rs)
	for _, receipt := range receipts {
		attempt := byHash[receipt.TxHash.String()]
		attempt.Receipts = append(attempt.Receipts, receipt)
	}
	return nil
}

// FindTxAttemptsRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time (up to limit)
func (o *evmTxStore) FindTxAttemptsRequiringResend(ctx context.Context, olderThan time.Time, maxInFlightTransactions uint32, chainID *big.Int, address common.Address) (attempts []TxAttempt, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	var dbAttempts []DbEthTxAttempt
	// this select distinct works because of unique index on evm.txes
	// (evm_chain_id, from_address, nonce)
	err = o.q.SelectContext(ctx, &dbAttempts, `
SELECT DISTINCT ON (evm.txes.nonce) evm.tx_attempts.*
FROM evm.tx_attempts
JOIN evm.txes ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE evm.tx_attempts.state <> 'in_progress' AND evm.txes.broadcast_at <= $1 AND evm_chain_id = $2 AND from_address = $3
ORDER BY evm.txes.nonce ASC, evm.tx_attempts.gas_price DESC, evm.tx_attempts.gas_tip_cap DESC
LIMIT $4
`, olderThan, chainID.String(), address, limit)

	attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
	return attempts, pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringResend failed to load evm.tx_attempts")
}

func (o *evmTxStore) UpdateBroadcastAts(ctx context.Context, now time.Time, etxIDs []int64) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	// Deliberately do nothing on NULL broadcast_at because that indicates the
	// tx has been moved into a state where broadcast_at is not relevant, e.g.
	// fatally errored.
	//
	// Since EthConfirmer/EthResender can race (totally OK since highest
	// priced transaction always wins) we only want to update broadcast_at if
	// our version is later.
	_, err := o.q.ExecContext(ctx, `UPDATE evm.txes SET broadcast_at = $1 WHERE id = ANY($2) AND broadcast_at < $1`, now, pq.Array(etxIDs))
	return pkgerrors.Wrap(err, "updateBroadcastAts failed to update evm.txes")
}

// SetBroadcastBeforeBlockNum updates already broadcast attempts with the
// current block number. This is safe no matter how old the head is because if
// the attempt is already broadcast it _must_ have been before this head.
func (o *evmTxStore) SetBroadcastBeforeBlockNum(ctx context.Context, blockNum int64, chainID *big.Int) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	_, err := o.q.ExecContext(ctx,
		`UPDATE evm.tx_attempts
SET broadcast_before_block_num = $1 
FROM evm.txes
WHERE evm.tx_attempts.broadcast_before_block_num IS NULL AND evm.tx_attempts.state = 'broadcast'
AND evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.txes.evm_chain_id = $2`,
		blockNum, chainID.String(),
	)
	return pkgerrors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
}

func (o *evmTxStore) FindTxAttemptsConfirmedMissingReceipt(ctx context.Context, chainID *big.Int) (attempts []TxAttempt, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbAttempts []DbEthTxAttempt
	err = o.q.SelectContext(ctx, &dbAttempts,
		`SELECT DISTINCT ON (evm.tx_attempts.eth_tx_id) evm.tx_attempts.*
		FROM evm.tx_attempts
		JOIN evm.txes ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.txes.state = 'confirmed_missing_receipt'
		WHERE evm_chain_id = $1
		ORDER BY evm.tx_attempts.eth_tx_id ASC, evm.tx_attempts.gas_price DESC, evm.tx_attempts.gas_tip_cap DESC`,
		chainID.String())
	if err != nil {
		err = pkgerrors.Wrap(err, "FindEtxAttemptsConfirmedMissingReceipt failed to query")
	}
	attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
	return
}

func (o *evmTxStore) UpdateTxsUnconfirmed(ctx context.Context, ids []int64) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	_, err := o.q.ExecContext(ctx, `UPDATE evm.txes SET state='unconfirmed' WHERE id = ANY($1)`, pq.Array(ids))

	if err != nil {
		return pkgerrors.Wrap(err, "UpdateEthTxsUnconfirmed failed to execute")
	}
	return nil
}

func (o *evmTxStore) FindTxAttemptsRequiringReceiptFetch(ctx context.Context, chainID *big.Int) (attempts []TxAttempt, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbAttempts []DbEthTxAttempt
		err = orm.q.SelectContext(ctx, &dbAttempts, `
SELECT evm.tx_attempts.* FROM evm.tx_attempts
JOIN evm.txes ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.txes.state IN ('unconfirmed', 'confirmed_missing_receipt') AND evm.txes.evm_chain_id = $1
WHERE evm.tx_attempts.state != 'insufficient_eth'
ORDER BY evm.txes.nonce ASC, evm.tx_attempts.gas_price DESC, evm.tx_attempts.gas_tip_cap DESC
`, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load evm.tx_attempts")
		}
		attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		err = orm.preloadTxesAtomic(ctx, attempts)
		return pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load evm.txes")
	})
	return
}

func (o *evmTxStore) SaveFetchedReceipts(ctx context.Context, r []*evmtypes.Receipt, chainID *big.Int) (err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	receipts := toOnchainReceipt(r)
	if len(receipts) == 0 {
		return nil
	}

	// Notes on this query:
	//
	// # Receipts insert
	// Conflict on (tx_hash, block_hash) shouldn't be possible because there
	// should only ever be one receipt for an eth_tx.
	//
	// ASIDE: This is because we mark confirmed atomically with receipt insert
	// in this query, and delete receipts upon marking unconfirmed - see
	// markForRebroadcast.
	//
	// If a receipt with the same (tx_hash, block_hash) exists then the
	// transaction is marked confirmed which means we _should_ never get here.
	// However, even so, it still shouldn't be an error to upsert a receipt we
	// already have.
	//
	// # EthTxAttempts update
	// It should always be safe to mark the attempt as broadcast here because
	// if it were not successfully broadcast how could it possibly have a
	// receipt?
	//
	// This state is reachable for example if the eth node errors so the
	// attempt was left in_progress but the transaction was actually accepted
	// and mined.
	//
	// # EthTxes update
	// Should be self-explanatory. If we got a receipt, the eth_tx is confirmed.
	//
	var valueStrs []string
	var valueArgs []interface{}
	for _, r := range receipts {
		var receiptJSON []byte
		receiptJSON, err = json.Marshal(r)
		if err != nil {
			return pkgerrors.Wrap(err, "saveFetchedReceipts failed to marshal JSON")
		}
		valueStrs = append(valueStrs, "(?,?,?,?,?,NOW())")
		valueArgs = append(valueArgs, r.TxHash, r.BlockHash, r.BlockNumber.Int64(), r.TransactionIndex, receiptJSON)
	}
	valueArgs = append(valueArgs, chainID.String())

	/* #nosec G201 */
	sql := `
	WITH inserted_receipts AS (
		INSERT INTO evm.receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at)
		VALUES %s
		ON CONFLICT (tx_hash, block_hash) DO UPDATE SET
			block_number = EXCLUDED.block_number,
			transaction_index = EXCLUDED.transaction_index,
			receipt = EXCLUDED.receipt
		RETURNING evm.receipts.tx_hash, evm.receipts.block_number
	),
	updated_eth_tx_attempts AS (
		UPDATE evm.tx_attempts
		SET
			state = 'broadcast',
			broadcast_before_block_num = COALESCE(evm.tx_attempts.broadcast_before_block_num, inserted_receipts.block_number)
		FROM inserted_receipts
		WHERE inserted_receipts.tx_hash = evm.tx_attempts.hash
		RETURNING evm.tx_attempts.eth_tx_id
	)
	UPDATE evm.txes
	SET state = 'confirmed'
	FROM updated_eth_tx_attempts
	WHERE updated_eth_tx_attempts.eth_tx_id = evm.txes.id
	AND evm_chain_id = ?
	`

	stmt := fmt.Sprintf(sql, strings.Join(valueStrs, ","))

	stmt = sqlx.Rebind(sqlx.DOLLAR, stmt)

	_, err = o.q.ExecContext(ctx, stmt, valueArgs...)
	return pkgerrors.Wrap(err, "SaveFetchedReceipts failed to save receipts")
}

// MarkAllConfirmedMissingReceipt
// It is possible that we can fail to get a receipt for all evm.tx_attempts
// even though a transaction with this nonce has long since been confirmed (we
// know this because transactions with higher nonces HAVE returned a receipt).
//
// This can probably only happen if an external wallet used the account (or
// conceivably because of some bug in the remote eth node that prevents it
// from returning a receipt for a valid transaction).
//
// In this case we mark these transactions as 'confirmed_missing_receipt' to
// prevent gas bumping.
//
// NOTE: We continue to attempt to resend evm.txes in this state on
// every head to guard against the extremely rare scenario of nonce gap due to
// reorg that excludes the transaction (from another wallet) that had this
// nonce (until finality depth is reached, after which we make the explicit
// decision to give up). This is done in the EthResender.
//
// We will continue to try to fetch a receipt for these attempts until all
// attempts are below the finality depth from current head.
func (o *evmTxStore) MarkAllConfirmedMissingReceipt(ctx context.Context, chainID *big.Int) (err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	res, err := o.q.ExecContext(ctx, `
UPDATE evm.txes
SET state = 'confirmed_missing_receipt'
FROM (
	SELECT from_address, MAX(nonce) as max_nonce 
	FROM evm.txes
	WHERE state = 'confirmed' AND evm_chain_id = $1
	GROUP BY from_address
) AS max_table
WHERE state = 'unconfirmed'
	AND evm_chain_id = $1
	AND nonce < max_table.max_nonce
	AND evm.txes.from_address = max_table.from_address
	`, chainID.String())
	if err != nil {
		return pkgerrors.Wrap(err, "markAllConfirmedMissingReceipt failed")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return pkgerrors.Wrap(err, "markAllConfirmedMissingReceipt RowsAffected failed")
	}
	if rowsAffected > 0 {
		o.logger.Infow(fmt.Sprintf("%d transactions missing receipt", rowsAffected), "n", rowsAffected)
	}
	return
}

func (o *evmTxStore) GetInProgressTxAttempts(ctx context.Context, address common.Address, chainID *big.Int) (attempts []TxAttempt, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbAttempts []DbEthTxAttempt
		err = orm.q.SelectContext(ctx, &dbAttempts, `
SELECT evm.tx_attempts.* FROM evm.tx_attempts
INNER JOIN evm.txes ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
WHERE evm.tx_attempts.state = 'in_progress' AND evm.txes.from_address = $1 AND evm.txes.evm_chain_id = $2
`, address, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed to load evm.tx_attempts")
		}
		attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		err = orm.preloadTxesAtomic(ctx, attempts)
		return pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed to load evm.txes")
	})
	return attempts, pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed")
}

// Find confirmed txes requiring callback but have not yet been signaled
func (o *evmTxStore) FindTxesPendingCallback(ctx context.Context, blockNum int64, chainID *big.Int) (receiptsPlus []ReceiptPlus, err error) {
	var rs []dbReceiptPlus

	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.q.SelectContext(ctx, &rs, `
	SELECT evm.txes.pipeline_task_run_id, evm.receipts.receipt, COALESCE((evm.txes.meta->>'FailOnRevert')::boolean, false) "FailOnRevert" FROM evm.txes
	INNER JOIN evm.tx_attempts ON evm.txes.id = evm.tx_attempts.eth_tx_id
	INNER JOIN evm.receipts ON evm.tx_attempts.hash = evm.receipts.tx_hash
	WHERE evm.txes.pipeline_task_run_id IS NOT NULL AND evm.txes.signal_callback = TRUE AND evm.txes.callback_completed = FALSE
	AND evm.receipts.block_number <= ($1 - evm.txes.min_confirmations) AND evm.txes.evm_chain_id = $2
	`, blockNum, chainID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions pending pipeline resume callback: %w", err)
	}
	receiptsPlus = fromDBReceiptsPlus(rs)
	return
}

// Update tx to mark that its callback has been signaled
func (o *evmTxStore) UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunId uuid.UUID, chainId *big.Int) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	_, err := o.q.ExecContext(ctx, `UPDATE evm.txes SET callback_completed = TRUE WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, pipelineTaskRunId, chainId.String())
	if err != nil {
		return fmt.Errorf("failed to mark callback completed for transaction: %w", err)
	}
	return nil
}

func (o *evmTxStore) FindLatestSequence(ctx context.Context, fromAddress common.Address, chainId *big.Int) (nonce evmtypes.Nonce, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	sql := `SELECT nonce FROM evm.txes WHERE from_address = $1 AND evm_chain_id = $2 AND nonce IS NOT NULL ORDER BY nonce DESC LIMIT 1`
	err = o.q.GetContext(ctx, &nonce, sql, fromAddress, chainId.String())
	return
}

// FindTxWithIdempotencyKey returns any broadcast ethtx with the given idempotencyKey and chainID
func (o *evmTxStore) FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID *big.Int) (etx *Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtx DbEthTx
	err = o.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE idempotency_key = $1 and evm_chain_id = $2`, idempotencyKey, chainID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, pkgerrors.Wrap(err, "FindTxWithIdempotencyKey failed to load evm.txes")
	}
	etx = new(Tx)
	dbEtx.ToTx(etx)
	return
}

// FindTxWithSequence returns any broadcast ethtx with the given nonce
func (o *evmTxStore) FindTxWithSequence(ctx context.Context, fromAddress common.Address, nonce evmtypes.Nonce) (etx *Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	etx = new(Tx)
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtx DbEthTx
		err = orm.q.GetContext(ctx, &dbEtx, `
SELECT * FROM evm.txes WHERE from_address = $1 AND nonce = $2 AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
`, fromAddress, nonce.Int64())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxWithNonce failed to load evm.txes")
		}
		dbEtx.ToTx(etx)
		err = orm.loadTxAttemptsAtomic(ctx, etx)
		return pkgerrors.Wrap(err, "FindEthTxWithNonce failed to load evm.tx_attempts")
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

func updateEthTxAttemptUnbroadcast(ctx context.Context, orm *evmTxStore, attempt TxAttempt) error {
	if attempt.State != txmgrtypes.TxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	_, err := orm.q.ExecContext(ctx, `UPDATE evm.tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = $1`, attempt.ID)
	return pkgerrors.Wrap(err, "updateEthTxAttemptUnbroadcast failed")
}

func updateEthTxUnconfirm(ctx context.Context, orm *evmTxStore, etx Tx) error {
	if etx.State != txmgr.TxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	_, err := orm.q.ExecContext(ctx, `UPDATE evm.txes SET state = 'unconfirmed' WHERE id = $1`, etx.ID)
	return pkgerrors.Wrap(err, "updateEthTxUnconfirm failed")
}

func deleteEthReceipts(ctx context.Context, orm *evmTxStore, etxID int64) (err error) {
	_, err = orm.q.ExecContext(ctx, `
DELETE FROM evm.receipts
USING evm.tx_attempts
WHERE evm.receipts.tx_hash = evm.tx_attempts.hash
AND evm.tx_attempts.eth_tx_id = $1
	`, etxID)
	return pkgerrors.Wrap(err, "deleteEthReceipts failed")
}

func (o *evmTxStore) UpdateTxForRebroadcast(ctx context.Context, etx Tx, etxAttempt TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		if err := deleteEthReceipts(ctx, orm, etx.ID); err != nil {
			return pkgerrors.Wrapf(err, "deleteEthReceipts failed for etx %v", etx.ID)
		}
		if err := updateEthTxUnconfirm(ctx, orm, etx); err != nil {
			return pkgerrors.Wrapf(err, "updateEthTxUnconfirm failed for etx %v", etx.ID)
		}
		return updateEthTxAttemptUnbroadcast(ctx, orm, etxAttempt)
	})
}

func (o *evmTxStore) FindTransactionsConfirmedInBlockRange(ctx context.Context, highBlockNumber, lowBlockNumber int64, chainID *big.Int) (etxs []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtxs []DbEthTx
		err = orm.q.SelectContext(ctx, &dbEtxs, `
SELECT DISTINCT evm.txes.* FROM evm.txes
INNER JOIN evm.tx_attempts ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.tx_attempts.state = 'broadcast'
INNER JOIN evm.receipts ON evm.receipts.tx_hash = evm.tx_attempts.hash
WHERE evm.txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number BETWEEN $1 AND $2 AND evm_chain_id = $3
ORDER BY nonce ASC
`, lowBlockNumber, highBlockNumber, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load evm.txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		if err = orm.LoadTxesAttempts(ctx, etxs); err != nil {
			return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load evm.tx_attempts")
		}
		err = orm.loadEthTxesAttemptsReceipts(ctx, etxs)
		return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load evm.receipts")
	})
	return etxs, pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed")
}

func (o *evmTxStore) FindEarliestUnconfirmedBroadcastTime(ctx context.Context, chainID *big.Int) (broadcastAt nullv4.Time, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		if err = orm.q.QueryRowxContext(ctx, `SELECT min(initial_broadcast_at) FROM evm.txes WHERE state = 'unconfirmed' AND evm_chain_id = $1`, chainID.String()).Scan(&broadcastAt); err != nil {
			return fmt.Errorf("failed to query for unconfirmed eth_tx count: %w", err)
		}
		return nil
	})
	return broadcastAt, err
}

func (o *evmTxStore) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context, chainID *big.Int) (earliestUnconfirmedTxBlock nullv4.Int, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		err = orm.q.QueryRowxContext(ctx, `
SELECT MIN(broadcast_before_block_num) FROM evm.tx_attempts
JOIN evm.txes ON evm.txes.id = evm.tx_attempts.eth_tx_id
WHERE evm.txes.state = 'unconfirmed'
AND evm_chain_id = $1`, chainID.String()).Scan(&earliestUnconfirmedTxBlock)
		if err != nil {
			return fmt.Errorf("failed to query for earliest unconfirmed tx block: %w", err)
		}
		return nil
	})
	return earliestUnconfirmedTxBlock, err
}

func (o *evmTxStore) IsTxFinalized(ctx context.Context, blockHeight int64, txID int64, chainID *big.Int) (finalized bool, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()

	var count int32
	err = o.q.GetContext(ctx, &count, `
    SELECT COUNT(evm.receipts.receipt) FROM evm.txes
    INNER JOIN evm.tx_attempts ON evm.txes.id = evm.tx_attempts.eth_tx_id
    INNER JOIN evm.receipts ON evm.tx_attempts.hash = evm.receipts.tx_hash
    WHERE evm.receipts.block_number <= ($1 - evm.txes.min_confirmations)
    AND evm.txes.id = $2 AND evm.txes.evm_chain_id = $3`, blockHeight, txID, chainID.String())
	if err != nil {
		return false, fmt.Errorf("failed to retrieve transaction reciepts: %w", err)
	}
	return count > 0, nil
}

func (o *evmTxStore) saveAttemptWithNewState(ctx context.Context, attempt TxAttempt, broadcastAt time.Time) error {
	var dbAttempt DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt)
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if _, err := orm.q.ExecContext(ctx, `UPDATE evm.txes SET broadcast_at = $1 WHERE id = $2 AND broadcast_at < $1`, broadcastAt, dbAttempt.EthTxID); err != nil {
			return pkgerrors.Wrap(err, "saveAttemptWithNewState failed to update evm.txes")
		}
		_, err := orm.q.ExecContext(ctx, `UPDATE evm.tx_attempts SET state=$1 WHERE id=$2`, dbAttempt.State, dbAttempt.ID)
		return pkgerrors.Wrap(err, "saveAttemptWithNewState failed to update evm.tx_attempts")
	})
}

func (o *evmTxStore) SaveInsufficientFundsAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if !(attempt.State == txmgrtypes.TxAttemptInProgress || attempt.State == txmgrtypes.TxAttemptInsufficientFunds) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = txmgrtypes.TxAttemptInsufficientFunds
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()
	return pkgerrors.Wrap(o.saveAttemptWithNewState(ctx, *attempt, broadcastAt), "saveInsufficientEthAttempt failed")
}

func (o *evmTxStore) saveSentAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = txmgrtypes.TxAttemptBroadcast
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return pkgerrors.Wrap(o.saveAttemptWithNewState(ctx, *attempt, broadcastAt), "saveSentAttempt failed")
}

func (o *evmTxStore) SaveSentAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	return o.saveSentAttempt(ctx, timeout, attempt, broadcastAt)
}

func (o *evmTxStore) SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err := o.Transact(ctx, false, func(orm *evmTxStore) error {
		if err := orm.saveSentAttempt(ctx, timeout, attempt, broadcastAt); err != nil {
			return err
		}
		if _, err := orm.q.ExecContext(ctx, `UPDATE evm.txes SET state = 'confirmed_missing_receipt' WHERE id = $1`, attempt.TxID); err != nil {
			return pkgerrors.Wrap(err, "failed to update evm.txes")
		}
		return nil
	})
	return pkgerrors.Wrap(err, "SaveConfirmedMissingReceiptAttempt failed")
}

func (o *evmTxStore) DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("DeleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("DeleteInProgressAttempt: expected attempt to have an id")
	}
	_, err := o.q.ExecContext(ctx, `DELETE FROM evm.tx_attempts WHERE id = $1`, attempt.ID)
	return pkgerrors.Wrap(err, "DeleteInProgressAttempt failed")
}

// SaveInProgressAttempt inserts or updates an attempt
func (o *evmTxStore) SaveInProgressAttempt(ctx context.Context, attempt *TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("SaveInProgressAttempt failed: attempt state must be in_progress")
	}
	var dbAttempt DbEthTxAttempt
	dbAttempt.FromTxAttempt(attempt)
	// Insert is the usual mode because the attempt is new
	if attempt.ID == 0 {
		query, args, e := o.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if e != nil {
			return pkgerrors.Wrap(e, "SaveInProgressAttempt failed to BindNamed")
		}
		e = o.q.GetContext(ctx, &dbAttempt, query, args...)
		dbAttempt.ToTxAttempt(attempt)
		return pkgerrors.Wrap(e, "SaveInProgressAttempt failed to insert into evm.tx_attempts")
	}
	// Update only applies to case of insufficient eth and simply changes the state to in_progress
	res, err := o.q.ExecContext(ctx, `UPDATE evm.tx_attempts SET state=$1, broadcast_before_block_num=$2 WHERE id=$3`, dbAttempt.State, dbAttempt.BroadcastBeforeBlockNum, dbAttempt.ID)
	if err != nil {
		return pkgerrors.Wrap(err, "SaveInProgressAttempt failed to update evm.tx_attempts")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return pkgerrors.Wrap(err, "SaveInProgressAttempt failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return pkgerrors.Wrapf(sql.ErrNoRows, "SaveInProgressAttempt tried to update evm.tx_attempts but no rows matched id %d", attempt.ID)
	}
	return nil
}

func (o *evmTxStore) GetAbandonedTransactionsByBatch(ctx context.Context, chainID *big.Int, enabledAddrs []common.Address, offset, limit uint) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()

	var enabledAddrsBytea [][]byte
	for _, addr := range enabledAddrs {
		enabledAddrsBytea = append(enabledAddrsBytea, addr[:])
	}

	// TODO: include confirmed txes https://smartcontract-it.atlassian.net/browse/BCI-2920
	query := `SELECT * FROM evm.txes WHERE state <> 'fatal_error' AND state <> 'confirmed' AND evm_chain_id = $1 
                       AND from_address <> ALL($2) ORDER BY nonce ASC OFFSET $3 LIMIT $4`

	var dbEtxs []DbEthTx
	if err = o.q.SelectContext(ctx, &dbEtxs, query, chainID.String(), enabledAddrsBytea, offset, limit); err != nil {
		return nil, fmt.Errorf("failed to load evm.txes: %w", err)
	}
	txes = make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)

	return txes, err
}

func (o *evmTxStore) GetTxByID(ctx context.Context, id int64) (txe *Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()

	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		stmt := `SELECT * FROM evm.txes WHERE id = $1`
		var dbEtxs []DbEthTx
		if err = orm.q.SelectContext(ctx, &dbEtxs, stmt, id); err != nil {
			return fmt.Errorf("failed to load evm.txes: %w", err)
		}
		txes := make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
		if len(txes) != 1 {
			return fmt.Errorf("failed to get tx with id %v", id)
		}
		txe = txes[0]
		err = o.LoadTxesAttempts(ctx, txes)
		if err != nil {
			return fmt.Errorf("failed to load evm.tx_attempts: %w", err)
		}
		return nil
	})

	return txe, nil
}

// FindTxsRequiringGasBump returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
//
// It also returns evm.txes that are unconfirmed with no evm.tx_attempts
func (o *evmTxStore) FindTxsRequiringGasBump(ctx context.Context, address common.Address, blockNum, gasBumpThreshold, depth int64, chainID *big.Int) (etxs []*Tx, err error) {
	if gasBumpThreshold == 0 {
		return
	}
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		stmt := `
SELECT evm.txes.* FROM evm.txes
LEFT JOIN evm.tx_attempts ON evm.txes.id = evm.tx_attempts.eth_tx_id AND (broadcast_before_block_num > $4 OR broadcast_before_block_num IS NULL OR evm.tx_attempts.state != 'broadcast')
WHERE evm.txes.state = 'unconfirmed' AND evm.tx_attempts.id IS NULL AND evm.txes.from_address = $1 AND evm.txes.evm_chain_id = $2
	AND (($3 = 0) OR (evm.txes.id IN (SELECT id FROM evm.txes WHERE state = 'unconfirmed' AND from_address = $1 ORDER BY nonce ASC LIMIT $3)))
ORDER BY nonce ASC
`
		var dbEtxs []DbEthTx
		if err = orm.q.SelectContext(ctx, &dbEtxs, stmt, address, chainID.String(), depth, blockNum-gasBumpThreshold); err != nil {
			return pkgerrors.Wrap(err, "FindEthTxsRequiringGasBump failed to load evm.txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		err = orm.LoadTxesAttempts(ctx, etxs)
		return pkgerrors.Wrap(err, "FindEthTxsRequiringGasBump failed to load evm.tx_attempts")
	})
	return
}

// FindTxsRequiringResubmissionDueToInsufficientFunds returns transactions
// that need to be re-sent because they hit an out-of-eth error on a previous
// block
func (o *evmTxStore) FindTxsRequiringResubmissionDueToInsufficientFunds(ctx context.Context, address common.Address, chainID *big.Int) (etxs []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtxs []DbEthTx
		err = orm.q.SelectContext(ctx, &dbEtxs, `
SELECT DISTINCT evm.txes.* FROM evm.txes
INNER JOIN evm.tx_attempts ON evm.txes.id = evm.tx_attempts.eth_tx_id AND evm.tx_attempts.state = 'insufficient_eth'
WHERE evm.txes.from_address = $1 AND evm.txes.state = 'unconfirmed' AND evm.txes.evm_chain_id = $2
ORDER BY nonce ASC
`, address, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load evm.txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		err = orm.LoadTxesAttempts(ctx, etxs)
		return pkgerrors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load evm.tx_attempts")
	})
	return
}

// markOldTxesMissingReceiptAsErrored
//
// Once eth_tx has all of its attempts broadcast before some cutoff threshold
// without receiving any receipts, we mark it as fatally errored (never sent).
//
// The job run will also be marked as errored in this case since we never got a
// receipt and thus cannot pass on any transaction hash
func (o *evmTxStore) MarkOldTxesMissingReceiptAsErrored(ctx context.Context, blockNum int64, finalityDepth uint32, chainID *big.Int) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	// cutoffBlockNum is a block height
	// Any 'confirmed_missing_receipt' eth_tx with all attempts older than this block height will be marked as errored
	// We will not try to query for receipts for this transaction any more
	cutoff := blockNum - int64(finalityDepth)
	if cutoff <= 0 {
		return nil
	}
	if cutoff <= 0 {
		return nil
	}
	// note: if QOpt passes in a sql.Tx this will reuse it
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		type etx struct {
			ID    int64
			Nonce int64
		}
		var data []etx
		err := orm.q.SelectContext(ctx, &data, `
UPDATE evm.txes
SET state='fatal_error', nonce=NULL, error=$1, broadcast_at=NULL, initial_broadcast_at=NULL
FROM (
	SELECT e1.id, e1.nonce, e1.from_address FROM evm.txes AS e1 WHERE id IN (
		SELECT e2.id FROM evm.txes AS e2
		INNER JOIN evm.tx_attempts ON e2.id = evm.tx_attempts.eth_tx_id
		WHERE e2.state = 'confirmed_missing_receipt'
		AND e2.evm_chain_id = $3
		GROUP BY e2.id
		HAVING max(evm.tx_attempts.broadcast_before_block_num) < $2
	)
	FOR UPDATE OF e1
) e0
WHERE e0.id = evm.txes.id
RETURNING e0.id, e0.nonce`, ErrCouldNotGetReceipt, cutoff, chainID.String())

		if err != nil {
			return pkgerrors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to query")
		}

		// We need this little lookup table because we have to have the nonce
		// from the first query, BEFORE it was updated/nullified
		lookup := make(map[int64]etx)
		for _, d := range data {
			lookup[d.ID] = d
		}
		etxIDs := make([]int64, len(data))
		for i := 0; i < len(data); i++ {
			etxIDs[i] = data[i].ID
		}

		type result struct {
			ID                         int64
			FromAddress                common.Address
			MaxBroadcastBeforeBlockNum int64
			TxHashes                   pq.ByteaArray
		}

		var results []result
		err = orm.q.SelectContext(ctx, &results, `
SELECT e.id, e.from_address, max(a.broadcast_before_block_num) AS max_broadcast_before_block_num, array_agg(a.hash) AS tx_hashes
FROM evm.txes e
INNER JOIN evm.tx_attempts a ON e.id = a.eth_tx_id
WHERE e.id = ANY($1)
GROUP BY e.id
`, etxIDs)

		if err != nil {
			return pkgerrors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to load additional data")
		}

		for _, r := range results {
			nonce := lookup[r.ID].Nonce
			txHashesHex := make([]common.Address, len(r.TxHashes))
			for i := 0; i < len(r.TxHashes); i++ {
				txHashesHex[i] = common.BytesToAddress(r.TxHashes[i])
			}

			orm.logger.Criticalw(fmt.Sprintf("eth_tx with ID %v expired without ever getting a receipt for any of our attempts. "+
				"Current block height is %v, transaction was broadcast before block height %v. This transaction may not have not been sent and will be marked as fatally errored. "+
				"This can happen if there is another instance of chainlink running that is using the same private key, or if "+
				"an external wallet has been used to send a transaction from account %s with nonce %v."+
				" Please note that Chainlink requires exclusive ownership of it's private keys and sharing keys across multiple"+
				" chainlink instances, or using the chainlink keys with an external wallet is NOT SUPPORTED and WILL lead to missed transactions",
				r.ID, blockNum, r.MaxBroadcastBeforeBlockNum, r.FromAddress, nonce), "ethTxID", r.ID, "nonce", nonce, "fromAddress", r.FromAddress, "txHashes", txHashesHex)
		}

		return nil
	})
}

func (o *evmTxStore) SaveReplacementInProgressAttempt(ctx context.Context, oldAttempt TxAttempt, replacementAttempt *TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if oldAttempt.State != txmgrtypes.TxAttemptInProgress || replacementAttempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		if _, err := orm.q.ExecContext(ctx, `DELETE FROM evm.tx_attempts WHERE id=$1`, oldAttempt.ID); err != nil {
			return pkgerrors.Wrap(err, "saveReplacementInProgressAttempt failed to delete from evm.tx_attempts")
		}
		var dbAttempt DbEthTxAttempt
		dbAttempt.FromTxAttempt(replacementAttempt)
		query, args, e := orm.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if e != nil {
			return pkgerrors.Wrap(e, "saveReplacementInProgressAttempt failed to BindNamed")
		}
		e = orm.q.GetContext(ctx, &dbAttempt, query, args...)
		dbAttempt.ToTxAttempt(replacementAttempt)
		return pkgerrors.Wrap(e, "saveReplacementInProgressAttempt failed to insert replacement attempt")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func (o *evmTxStore) FindNextUnstartedTransactionFromAddress(ctx context.Context, fromAddress common.Address, chainID *big.Int) (*Tx, error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtx DbEthTx
	err := o.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2 ORDER BY value ASC, created_at ASC, id ASC`, fromAddress, chainID.String())
	etx := new(Tx)
	dbEtx.ToTx(etx)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to FindNextUnstartedTransactionFromAddress")
	}

	return etx, nil
}

func (o *evmTxStore) UpdateTxFatalError(ctx context.Context, etx *Tx) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if etx.State != txmgr.TxInProgress && etx.State != txmgr.TxUnstarted {
		return pkgerrors.Errorf("can only transition to fatal_error from in_progress or unstarted, transaction is currently %s", etx.State)
	}
	if !etx.Error.Valid {
		return errors.New("expected error field to be set")
	}

	etx.Sequence = nil
	etx.State = txmgr.TxFatalError

	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		if _, err := orm.q.ExecContext(ctx, `DELETE FROM evm.tx_attempts WHERE eth_tx_id = $1`, etx.ID); err != nil {
			return pkgerrors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		var dbEtx DbEthTx
		dbEtx.FromTx(etx)
		err := pkgerrors.Wrap(orm.q.GetContext(ctx, &dbEtx, `UPDATE evm.txes SET state=$1, error=$2, broadcast_at=NULL, initial_broadcast_at=NULL, nonce=NULL WHERE id=$3 RETURNING *`, etx.State, etx.Error, etx.ID), "saveFatallyErroredTransaction failed to save eth_tx")
		dbEtx.ToTx(etx)
		return err
	})
}

// Updates eth attempt from in_progress to broadcast. Also updates the eth tx to unconfirmed.
func (o *evmTxStore) UpdateTxAttemptInProgressToBroadcast(ctx context.Context, etx *Tx, attempt TxAttempt, NewAttemptState txmgrtypes.TxAttemptState) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if etx.BroadcastAt == nil {
		return errors.New("unconfirmed transaction must have broadcast_at time")
	}
	if etx.InitialBroadcastAt == nil {
		return errors.New("unconfirmed transaction must have initial_broadcast_at time")
	}
	if etx.State != txmgr.TxInProgress {
		return pkgerrors.Errorf("can only transition to unconfirmed from in_progress, transaction is currently %s", etx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("attempt must be in in_progress state")
	}
	if NewAttemptState != txmgrtypes.TxAttemptBroadcast {
		return pkgerrors.Errorf("new attempt state must be broadcast, got: %s", NewAttemptState)
	}
	etx.State = txmgr.TxUnconfirmed
	attempt.State = NewAttemptState
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		var dbEtx DbEthTx
		dbEtx.FromTx(etx)
		if err := orm.q.GetContext(ctx, &dbEtx, `UPDATE evm.txes SET state=$1, error=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id = $5 RETURNING *`, dbEtx.State, dbEtx.Error, dbEtx.BroadcastAt, dbEtx.InitialBroadcastAt, dbEtx.ID); err != nil {
			return pkgerrors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx")
		}
		dbEtx.ToTx(etx)
		var dbAttempt DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt)
		if err := orm.q.GetContext(ctx, &dbAttempt, `UPDATE evm.tx_attempts SET state = $1 WHERE id = $2 RETURNING *`, dbAttempt.State, dbAttempt.ID); err != nil {
			return pkgerrors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx_attempt")
		}
		return nil
	})
}

// Updates eth tx from unstarted to in_progress and inserts in_progress eth attempt
func (o *evmTxStore) UpdateTxUnstartedToInProgress(ctx context.Context, etx *Tx, attempt *TxAttempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if etx.Sequence == nil {
		return errors.New("in_progress transaction must have nonce")
	}
	if etx.State != txmgr.TxUnstarted {
		return pkgerrors.Errorf("can only transition to in_progress from unstarted, transaction is currently %s", etx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("attempt state must be in_progress")
	}
	etx.State = txmgr.TxInProgress
	return o.Transact(ctx, false, func(orm *evmTxStore) error {
		// If a replay was triggered while unconfirmed transactions were pending, they will be marked as fatal_error => abandoned.
		// In this case, we must remove the abandoned attempt from evm.tx_attempts before replacing it with a new one.  In any other
		// case, we uphold the constraint, leaving the original tx attempt as-is and returning the constraint violation error.
		//
		// Note:  the record of the original abandoned transaction will remain in evm.txes, only the attempt is replaced.  (Any receipt
		// associated with the abandoned attempt would also be lost, although this shouldn't happen since only unconfirmed transactions
		// can be abandoned.)
		res, err2 := orm.q.ExecContext(ctx, `DELETE FROM evm.tx_attempts a USING evm.txes t
			WHERE t.id = a.eth_tx_id AND a.hash = $1 AND t.state = $2 AND t.error = 'abandoned'`,
			attempt.Hash, txmgr.TxFatalError,
		)

		if err2 != nil {
			// If the DELETE fails, we don't want to abort before at least attempting the INSERT. tx hash conflicts with
			// abandoned transactions can only happen after a nonce reset. If the node is operating normally but there is
			// some unexpected issue with the DELETE query, blocking the txmgr from sending transactions would be risky
			// and could potentially get the node stuck. If the INSERT is going to succeed then we definitely want to continue.
			// And even if the INSERT fails, an error message showing the txmgr is having trouble inserting tx's in the db may be
			// easier to understand quickly if there is a problem with the node.
			o.logger.Errorw("Ignoring unexpected db error while checking for txhash conflict", "err", err2)
		} else if rows, err := res.RowsAffected(); err != nil {
			o.logger.Errorw("Ignoring unexpected db error reading rows affected while checking for txhash conflict", "err", err)
		} else if rows > 0 {
			o.logger.Debugf("Replacing abandoned tx with tx hash %s with tx_id=%d with identical tx hash", attempt.Hash, attempt.TxID)
		}

		var dbAttempt DbEthTxAttempt
		dbAttempt.FromTxAttempt(attempt)
		query, args, err := orm.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if err != nil {
			return pkgerrors.Wrap(err, "UpdateTxUnstartedToInProgress failed to BindNamed")
		}
		err = orm.q.GetContext(ctx, &dbAttempt, query, args...)
		if err != nil {
			var pqErr *pgconn.PgError
			if isPqErr := errors.As(err, &pqErr); isPqErr &&
				pqErr.SchemaName == "evm" &&
				pqErr.ConstraintName == "eth_tx_attempts_eth_tx_id_fkey" {
				return txmgr.ErrTxRemoved
			}
			if err != nil {
				return pkgerrors.Wrap(err, "UpdateTxUnstartedToInProgress failed to create eth_tx_attempt")
			}
		}
		dbAttempt.ToTxAttempt(attempt)
		var dbEtx DbEthTx
		dbEtx.FromTx(etx)
		err = orm.q.GetContext(ctx, &dbEtx, `UPDATE evm.txes SET nonce=$1, state=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id=$5 RETURNING *`, etx.Sequence, etx.State, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID)
		dbEtx.ToTx(etx)
		return pkgerrors.Wrap(err, "UpdateTxUnstartedToInProgress failed to update eth_tx")
	})
}

// GetTxInProgress returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnstartedEthTxs loop.
// It may or may not have been broadcast to an eth node.
func (o *evmTxStore) GetTxInProgress(ctx context.Context, fromAddress common.Address) (etx *Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	etx = new(Tx)
	if err != nil {
		return etx, pkgerrors.Wrap(err, "getInProgressEthTx failed")
	}
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtx DbEthTx
		err = orm.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE from_address = $1 and state = 'in_progress'`, fromAddress)
		if errors.Is(err, sql.ErrNoRows) {
			etx = nil
			return nil
		} else if err != nil {
			return pkgerrors.Wrap(err, "GetTxInProgress failed while loading eth tx")
		}
		dbEtx.ToTx(etx)
		if err = o.loadTxAttemptsAtomic(ctx, etx); err != nil {
			return pkgerrors.Wrap(err, "GetTxInProgress failed while loading EthTxAttempts")
		}
		if len(etx.TxAttempts) != 1 || etx.TxAttempts[0].State != txmgrtypes.TxAttemptInProgress {
			return pkgerrors.Errorf("invariant violation: expected in_progress transaction %v to have exactly one unsent attempt. "+
				"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", etx.ID)
		}
		return nil
	})

	return etx, pkgerrors.Wrap(err, "getInProgressEthTx failed")
}

func (o *evmTxStore) HasInProgressTransaction(ctx context.Context, account common.Address, chainID *big.Int) (exists bool, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.q.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM evm.txes WHERE state = 'in_progress' AND from_address = $1 AND evm_chain_id = $2)`, account, chainID.String())
	return exists, pkgerrors.Wrap(err, "hasInProgressTransaction failed")
}

func (o *evmTxStore) countTransactionsWithState(ctx context.Context, fromAddress common.Address, state txmgrtypes.TxState, chainID *big.Int) (count uint32, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.q.GetContext(ctx, &count, `SELECT count(*) FROM evm.txes WHERE from_address = $1 AND state = $2 AND evm_chain_id = $3`,
		fromAddress, state, chainID.String())
	return count, pkgerrors.Wrap(err, "failed to countTransactionsWithState")
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func (o *evmTxStore) CountUnconfirmedTransactions(ctx context.Context, fromAddress common.Address, chainID *big.Int) (count uint32, err error) {
	return o.countTransactionsWithState(ctx, fromAddress, txmgr.TxUnconfirmed, chainID)
}

// CountTransactionsByState returns the number of transactions with any fromAddress in the given state
func (o *evmTxStore) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState, chainID *big.Int) (count uint32, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.q.GetContext(ctx, &count, `SELECT count(*) FROM evm.txes WHERE state = $1 AND evm_chain_id = $2`,
		state, chainID.String())
	if err != nil {
		return 0, fmt.Errorf("failed to CountTransactionsByState: %w", err)
	}
	return count, nil
}

// CountUnstartedTransactions returns the number of unconfirmed transactions
func (o *evmTxStore) CountUnstartedTransactions(ctx context.Context, fromAddress common.Address, chainID *big.Int) (count uint32, err error) {
	return o.countTransactionsWithState(ctx, fromAddress, txmgr.TxUnstarted, chainID)
}

func (o *evmTxStore) CheckTxQueueCapacity(ctx context.Context, fromAddress common.Address, maxQueuedTransactions uint64, chainID *big.Int) (err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	if maxQueuedTransactions == 0 {
		return nil
	}
	var count uint64
	err = o.q.GetContext(ctx, &count, `SELECT count(*) FROM evm.txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2`, fromAddress, chainID.String())
	if err != nil {
		err = pkgerrors.Wrap(err, "CheckTxQueueCapacity query failed")
		return
	}

	if count >= maxQueuedTransactions {
		err = pkgerrors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}
	return
}

func (o *evmTxStore) CreateTransaction(ctx context.Context, txRequest TxRequest, chainID *big.Int) (tx Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtx DbEthTx
	err = o.Transact(ctx, false, func(orm *evmTxStore) error {
		if txRequest.PipelineTaskRunID != nil {
			err = orm.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, txRequest.PipelineTaskRunID, chainID.String())
			// If no eth_tx matches (the common case) then continue
			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return pkgerrors.Wrap(err, "CreateEthTransaction")
				}
				// if a previous transaction for this task run exists, immediately return it
				return nil
			}
		}
		err = orm.q.GetContext(ctx, &dbEtx, `
INSERT INTO evm.txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, transmit_checker, idempotency_key, signal_callback)
VALUES (
$1,$2,$3,$4,$5,'unstarted',NOW(),$6,$7,$8,$9,$10,$11,$12,$13
)
RETURNING "txes".*
`, txRequest.FromAddress, txRequest.ToAddress, txRequest.EncodedPayload, assets.Eth(txRequest.Value), txRequest.FeeLimit, txRequest.Meta, txRequest.Strategy.Subject(), chainID.String(), txRequest.MinConfirmations, txRequest.PipelineTaskRunID, txRequest.Checker, txRequest.IdempotencyKey, txRequest.SignalCallback)
		if err != nil {
			return pkgerrors.Wrap(err, "CreateEthTransaction failed to insert evm tx")
		}
		return nil
	})
	var etx Tx
	dbEtx.ToTx(&etx)
	return etx, err
}

func (o *evmTxStore) PruneUnstartedTxQueue(ctx context.Context, queueSize uint32, subject uuid.UUID) (ids []int64, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, false, func(orm *evmTxStore) error {
		err := orm.q.SelectContext(ctx, &ids, `
DELETE FROM evm.txes
WHERE state = 'unstarted' AND subject = $1 AND
id < (
	SELECT min(id) FROM (
		SELECT id
		FROM evm.txes
		WHERE state = 'unstarted' AND subject = $2
		ORDER BY id DESC
		LIMIT $3
	) numbers
) RETURNING id`, subject, subject, queueSize)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("PruneUnstartedTxQueue failed: %w", err)
		}
		return err
	})
	return
}

func (o *evmTxStore) ReapTxHistory(ctx context.Context, minBlockNumberToKeep int64, timeThreshold time.Time, chainID *big.Int) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()

	// Delete old confirmed evm.txes
	// NOTE that this relies on foreign key triggers automatically removing
	// the evm.tx_attempts and evm.receipts linked to every eth_tx
	const batchSize = 1000
	err := sqlutil.Batch(func(_, limit uint) (count uint, err error) {
		res, err := o.q.ExecContext(ctx, `
WITH old_enough_receipts AS (
	SELECT tx_hash FROM evm.receipts
	WHERE block_number < $1
	ORDER BY block_number ASC, id ASC
	LIMIT $2
)
DELETE FROM evm.txes
USING old_enough_receipts, evm.tx_attempts
WHERE evm.tx_attempts.eth_tx_id = evm.txes.id
AND evm.tx_attempts.hash = old_enough_receipts.tx_hash
AND evm.txes.created_at < $3
AND evm.txes.state = 'confirmed'
AND evm_chain_id = $4`, minBlockNumberToKeep, limit, timeThreshold, chainID.String())
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to delete old confirmed evm.txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	}, batchSize)
	if err != nil {
		return pkgerrors.Wrap(err, "TxmReaper#reapEthTxes batch delete of confirmed evm.txes failed")
	}
	// Delete old 'fatal_error' evm.txes
	err = sqlutil.Batch(func(_, limit uint) (count uint, err error) {
		res, err := o.q.ExecContext(ctx, `
DELETE FROM evm.txes
WHERE created_at < $1
AND state = 'fatal_error'
AND evm_chain_id = $2`, timeThreshold, chainID.String())
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to delete old fatally errored evm.txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	}, batchSize)
	if err != nil {
		return pkgerrors.Wrap(err, "TxmReaper#reapEthTxes batch delete of fatally errored evm.txes failed")
	}

	return nil
}

func (o *evmTxStore) Abandon(ctx context.Context, chainID *big.Int, addr common.Address) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	_, err := o.q.ExecContext(ctx, `UPDATE evm.txes SET state='fatal_error', nonce = NULL, error = 'abandoned' WHERE state IN ('unconfirmed', 'in_progress', 'unstarted') AND evm_chain_id = $1 AND from_address = $2`, chainID.String(), addr)
	return err
}

// Find transactions by a field in the TxMeta blob and transaction states
func (o *evmTxStore) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) ([]*Tx, error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtxs []DbEthTx
	sql := fmt.Sprintf("SELECT * FROM evm.txes WHERE evm_chain_id = $1 AND meta->>'%s' = $2 AND state = ANY($3)", metaField)
	err := o.q.SelectContext(ctx, &dbEtxs, sql, chainID.String(), metaValue, pq.Array(states))
	txes := make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
	return txes, pkgerrors.Wrap(err, "failed to FindTxesByMetaFieldAndStates")
}

// Find transactions with a non-null TxMeta field that was provided by transaction states
func (o *evmTxStore) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtxs []DbEthTx
	sql := fmt.Sprintf("SELECT * FROM evm.txes WHERE meta->'%s' IS NOT NULL AND state = ANY($1) AND evm_chain_id = $2", metaField)
	err = o.q.SelectContext(ctx, &dbEtxs, sql, pq.Array(states), chainID.String())
	txes = make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
	return txes, pkgerrors.Wrap(err, "failed to FindTxesWithMetaFieldByStates")
}

// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
func (o *evmTxStore) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtxs []DbEthTx
	sql := fmt.Sprintf("SELECT et.* FROM evm.txes et JOIN evm.tx_attempts eta on et.id = eta.eth_tx_id JOIN evm.receipts er on eta.hash = er.tx_hash WHERE et.meta->'%s' IS NOT NULL AND er.block_number >= $1 AND et.evm_chain_id = $2", metaField)
	err = o.q.SelectContext(ctx, &dbEtxs, sql, blockNum, chainID.String())
	txes = make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
	return txes, pkgerrors.Wrap(err, "failed to FindTxesWithMetaFieldByReceiptBlockNum")
}

// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
func (o *evmTxStore) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtxs []DbEthTx
		if err = orm.q.SelectContext(ctx, &dbEtxs, `SELECT * FROM evm.txes WHERE id = ANY($1) AND state = ANY($2) AND evm_chain_id = $3`, pq.Array(ids), pq.Array(states), chainID.String()); err != nil {
			return pkgerrors.Wrapf(err, "failed to find evm.txes")
		}
		txes = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
		if err = orm.LoadTxesAttempts(ctx, txes); err != nil {
			return pkgerrors.Wrapf(err, "failed to load evm.tx_attempts for evm.tx")
		}
		if err = orm.loadEthTxesAttemptsReceipts(ctx, txes); err != nil {
			return pkgerrors.Wrapf(err, "failed to load evm.receipts for evm.tx")
		}
		return nil
	})
	return txes, pkgerrors.Wrap(err, "FindTxesWithAttemptsAndReceiptsByIdsAndState failed")
}

// For testing only, get all txes in the DB
func (o *evmTxStore) GetAllTxes(ctx context.Context) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbEtxs []DbEthTx
	sql := "SELECT * FROM evm.txes"
	err = o.q.SelectContext(ctx, &dbEtxs, sql)
	txes = make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
	return txes, err
}

// For testing only, get all tx attempts in the DB
func (o *evmTxStore) GetAllTxAttempts(ctx context.Context) (attempts []TxAttempt, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	var dbAttempts []DbEthTxAttempt
	sql := "SELECT * FROM evm.tx_attempts"
	err = o.q.SelectContext(ctx, &dbAttempts, sql)
	attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
	return attempts, err
}

func (o *evmTxStore) CountTxesByStateAndSubject(ctx context.Context, state txmgrtypes.TxState, subject uuid.UUID) (count int, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	sql := "SELECT COUNT(*) FROM evm.txes WHERE state = $1 AND subject = $2"
	err = o.q.GetContext(ctx, &count, sql, state, subject)
	return count, err
}

func (o *evmTxStore) FindTxesByFromAddressAndState(ctx context.Context, fromAddress common.Address, state string) (txes []*Tx, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	sql := "SELECT * FROM evm.txes WHERE from_address = $1 AND state = $2"
	var dbEtxs []DbEthTx
	err = o.q.SelectContext(ctx, &dbEtxs, sql, fromAddress, state)
	txes = make([]*Tx, len(dbEtxs))
	dbEthTxsToEvmEthTxPtrs(dbEtxs, txes)
	return txes, err
}

func (o *evmTxStore) UpdateTxAttemptBroadcastBeforeBlockNum(ctx context.Context, id int64, blockNum uint) error {
	var cancel context.CancelFunc
	ctx, cancel = o.mergeContexts(ctx)
	defer cancel()
	sql := "UPDATE evm.tx_attempts SET broadcast_before_block_num = $1 WHERE eth_tx_id = $2"
	_, err := o.q.ExecContext(ctx, sql, blockNum, id)
	return err
}

// Returns a context that contains the values of the provided context,
// and which is canceled when either the provided context or TxStore parent context is canceled.
func (o *evmTxStore) mergeContexts(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelCauseFunc
	ctx, cancel = context.WithCancelCause(ctx)
	stop := context.AfterFunc(o.ctx, func() {
		cancel(context.Cause(o.ctx))
	})
	return ctx, func() {
		stop()
		cancel(context.Canceled)
	}
}
