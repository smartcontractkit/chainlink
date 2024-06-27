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
	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	ErrKeyNotUpdated = errors.New("evmTxStore: Key not updated")
	ErrInvalidQOpt   = errors.New("evmTxStore: Invalid QOpt")

	// ErrCouldNotGetReceipt is the error string we save if we reach our finality depth for a confirmed transaction without ever getting a receipt
	// This most likely happened because an external wallet used the account for this nonce
	ErrCouldNotGetReceipt = "could not get receipt"
)

// EvmTxStore combines the txmgr tx store interface and the interface needed for the the API to read from the tx DB
//
//go:generate mockery --quiet --name EvmTxStore --output ./mocks/ --case=underscore
type EvmTxStore interface {
	// redeclare TxStore for mockery
	txmgrtypes.TxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee]
	TxStoreWebApi
}

// TxStoreWebApi encapsulates the methods that are not used by the txmgr and only used by the various web controllers and readers
type TxStoreWebApi interface {
	FindTxAttemptConfirmedByTxIDs(ids []int64) ([]TxAttempt, error)
	FindTxByHash(hash common.Hash) (*Tx, error)
	Transactions(offset, limit int) ([]Tx, int, error)
	TxAttempts(offset, limit int) ([]TxAttempt, int, error)
	TransactionsWithAttempts(offset, limit int) ([]Tx, int, error)
	FindTxAttempt(hash common.Hash) (*TxAttempt, error)
}

type TestEvmTxStore interface {
	EvmTxStore

	// methods only used for testing purposes
	InsertReceipt(receipt *evmtypes.Receipt) (int64, error)
	InsertTx(etx *Tx) error
	FindTxAttemptsByTxIDs(ids []int64) ([]TxAttempt, error)
	FindTxWithAttempts(etxID int64) (etx Tx, err error)
	InsertTxAttempt(attempt *TxAttempt) error
	LoadTxesAttempts(etxs []*Tx, qopts ...pg.QOpt) error
}

type evmTxStore struct {
	q         pg.Q
	logger    logger.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc
}

var _ EvmTxStore = (*evmTxStore)(nil)
var _ TestEvmTxStore = (*evmTxStore)(nil)

// Directly maps to columns of database table "eth_receipts".
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

// Directly maps to some columns of few database tables.
// Does not map to a single database table.
// It's comprised of fields from different tables.
type dbReceiptPlus struct {
	ID           uuid.UUID        `db:"id"`
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
		receipts[i] = rawOnchainReceipt(*rs[i])
	}
	return receipts
}

// Directly maps to columns of database table "eth_txes".
// This is exported, as tests and other external code still directly reads DB using this schema.
type DbEthTx struct {
	ID             int64
	Nonce          *int64
	FromAddress    common.Address
	ToAddress      common.Address
	EncodedPayload []byte
	Value          assets.Eth
	// GasLimit on the EthTx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	GasLimit uint32
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
	Meta              *datatypes.JSON
	Subject           uuid.NullUUID
	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  null.Uint32
	EVMChainID        utils.Big
	// TransmitChecker defines the check that should be performed before a transaction is submitted on
	// chain.
	TransmitChecker    *datatypes.JSON
	InitialBroadcastAt *time.Time
}

func DbEthTxFromEthTx(ethTx *Tx) DbEthTx {
	tx := DbEthTx{
		ID:                 ethTx.ID,
		FromAddress:        ethTx.FromAddress,
		ToAddress:          ethTx.ToAddress,
		EncodedPayload:     ethTx.EncodedPayload,
		Value:              assets.Eth(ethTx.Value),
		GasLimit:           ethTx.FeeLimit,
		Error:              ethTx.Error,
		BroadcastAt:        ethTx.BroadcastAt,
		CreatedAt:          ethTx.CreatedAt,
		State:              ethTx.State,
		Meta:               ethTx.Meta,
		Subject:            ethTx.Subject,
		PipelineTaskRunID:  ethTx.PipelineTaskRunID,
		MinConfirmations:   ethTx.MinConfirmations,
		TransmitChecker:    ethTx.TransmitChecker,
		InitialBroadcastAt: ethTx.InitialBroadcastAt,
	}

	if ethTx.ChainID != nil {
		tx.EVMChainID = *utils.NewBig(ethTx.ChainID)
	}
	if ethTx.Sequence != nil {
		n := ethTx.Sequence.Int64()
		tx.Nonce = &n
	}

	return tx
}

func DbEthTxToEthTx(dbEthTx DbEthTx, evmEthTx *Tx) {
	evmEthTx.ID = dbEthTx.ID
	if dbEthTx.Nonce != nil {
		n := evmtypes.Nonce(*dbEthTx.Nonce)
		evmEthTx.Sequence = &n
	}
	evmEthTx.FromAddress = dbEthTx.FromAddress
	evmEthTx.ToAddress = dbEthTx.ToAddress
	evmEthTx.EncodedPayload = dbEthTx.EncodedPayload
	evmEthTx.Value = *dbEthTx.Value.ToInt()
	evmEthTx.FeeLimit = dbEthTx.GasLimit
	evmEthTx.Error = dbEthTx.Error
	evmEthTx.BroadcastAt = dbEthTx.BroadcastAt
	evmEthTx.CreatedAt = dbEthTx.CreatedAt
	evmEthTx.State = dbEthTx.State
	evmEthTx.Meta = dbEthTx.Meta
	evmEthTx.Subject = dbEthTx.Subject
	evmEthTx.PipelineTaskRunID = dbEthTx.PipelineTaskRunID
	evmEthTx.MinConfirmations = dbEthTx.MinConfirmations
	evmEthTx.ChainID = dbEthTx.EVMChainID.ToInt()
	evmEthTx.TransmitChecker = dbEthTx.TransmitChecker
	evmEthTx.InitialBroadcastAt = dbEthTx.InitialBroadcastAt
}

func dbEthTxsToEvmEthTxs(dbEthTxs []DbEthTx) []Tx {
	evmEthTxs := make([]Tx, len(dbEthTxs))
	for i, dbTx := range dbEthTxs {
		DbEthTxToEthTx(dbTx, &evmEthTxs[i])
	}
	return evmEthTxs
}

func dbEthTxsToEvmEthTxPtrs(dbEthTxs []DbEthTx, evmEthTxs []*Tx) {
	for i, dbTx := range dbEthTxs {
		evmEthTxs[i] = &Tx{}
		DbEthTxToEthTx(dbTx, evmEthTxs[i])
	}
}

// Directly maps to columns of database table "eth_tx_attempts".
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
	ChainSpecificGasLimit   uint32
	TxType                  int
	GasTipCap               *assets.Wei
	GasFeeCap               *assets.Wei
}

func DbEthTxAttemptFromEthTxAttempt(ethTxAttempt *TxAttempt) DbEthTxAttempt {
	dbTx := DbEthTxAttempt{
		ID:                      ethTxAttempt.ID,
		EthTxID:                 ethTxAttempt.TxID,
		GasPrice:                ethTxAttempt.TxFee.Legacy,
		SignedRawTx:             ethTxAttempt.SignedRawTx,
		Hash:                    ethTxAttempt.Hash,
		BroadcastBeforeBlockNum: ethTxAttempt.BroadcastBeforeBlockNum,
		CreatedAt:               ethTxAttempt.CreatedAt,
		ChainSpecificGasLimit:   ethTxAttempt.ChainSpecificFeeLimit,
		TxType:                  ethTxAttempt.TxType,
		GasTipCap:               ethTxAttempt.TxFee.DynamicTipCap,
		GasFeeCap:               ethTxAttempt.TxFee.DynamicFeeCap,
	}

	// handle state naming difference between generic + EVM
	if ethTxAttempt.State == txmgrtypes.TxAttemptInsufficientFunds {
		dbTx.State = "insufficient_eth"
	} else {
		dbTx.State = ethTxAttempt.State.String()
	}

	return dbTx
}

func DbEthTxAttemptStateToTxAttemptState(state string) txmgrtypes.TxAttemptState {
	if state == "insufficient_eth" {
		return txmgrtypes.TxAttemptInsufficientFunds
	}
	return txmgrtypes.NewTxAttemptState(state)
}

func DbEthTxAttemptToEthTxAttempt(dbEthTxAttempt DbEthTxAttempt, evmAttempt *TxAttempt) {
	evmAttempt.ID = dbEthTxAttempt.ID
	evmAttempt.TxID = dbEthTxAttempt.EthTxID
	evmAttempt.SignedRawTx = dbEthTxAttempt.SignedRawTx
	evmAttempt.Hash = dbEthTxAttempt.Hash
	evmAttempt.BroadcastBeforeBlockNum = dbEthTxAttempt.BroadcastBeforeBlockNum
	evmAttempt.State = DbEthTxAttemptStateToTxAttemptState(dbEthTxAttempt.State)
	evmAttempt.CreatedAt = dbEthTxAttempt.CreatedAt
	evmAttempt.ChainSpecificFeeLimit = dbEthTxAttempt.ChainSpecificGasLimit
	evmAttempt.TxType = dbEthTxAttempt.TxType
	evmAttempt.TxFee = gas.EvmFee{
		Legacy:        dbEthTxAttempt.GasPrice,
		DynamicTipCap: dbEthTxAttempt.GasTipCap,
		DynamicFeeCap: dbEthTxAttempt.GasFeeCap,
	}
}

func dbEthTxAttemptsToEthTxAttempts(dbEthTxAttempt []DbEthTxAttempt) []TxAttempt {
	evmEthTxAttempt := make([]TxAttempt, len(dbEthTxAttempt))
	for i, dbTxAttempt := range dbEthTxAttempt {
		DbEthTxAttemptToEthTxAttempt(dbTxAttempt, &evmEthTxAttempt[i])
	}
	return evmEthTxAttempt
}

func NewTxStore(
	db *sqlx.DB,
	lggr logger.Logger,
	cfg pg.QConfig,
) *evmTxStore {
	namedLogger := lggr.Named("TxmStore")
	ctx, cancel := context.WithCancel(context.Background())
	q := pg.NewQ(db, namedLogger, cfg, pg.WithParentCtx(ctx))
	return &evmTxStore{
		q:         q,
		logger:    namedLogger,
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

const insertIntoEthTxAttemptsQuery = `
INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap)
VALUES (:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap)
RETURNING *;
`

// TODO: create method to pass in new context to evmTxStore (which will also create a new pg.Q)

func (o *evmTxStore) Close() {
	o.ctxCancel()
}

func (o *evmTxStore) preloadTxAttempts(txs []Tx) error {
	// Preload TxAttempts
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	var dbAttempts []DbEthTxAttempt
	sql := `SELECT * FROM eth_tx_attempts WHERE eth_tx_id IN (?) ORDER BY id desc;`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = o.q.Rebind(query)
	if err = o.q.Select(&dbAttempts, query, args...); err != nil {
		return err
	}
	// fill in attempts
	for _, dbAttempt := range dbAttempts {
		for i, tx := range txs {
			if tx.ID == dbAttempt.EthTxID {
				var attempt TxAttempt
				DbEthTxAttemptToEthTxAttempt(dbAttempt, &attempt)
				txs[i].TxAttempts = append(txs[i].TxAttempts, attempt)
			}
		}
	}
	return nil
}

func (o *evmTxStore) PreloadTxes(attempts []TxAttempt, qopts ...pg.QOpt) error {
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
	qq := o.q.WithOpts(qopts...)
	if err := qq.Select(&dbEthTxs, `SELECT * FROM eth_txes WHERE id = ANY($1)`, pq.Array(ethTxIDs)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxes failed")
	}
	for _, dbEtx := range dbEthTxs {
		etx := ethTxM[dbEtx.ID]
		DbEthTxToEthTx(dbEtx, &etx)
		ethTxM[etx.ID] = etx
	}
	for i, attempt := range attempts {
		attempts[i].Tx = ethTxM[attempt.TxID]
	}
	return nil
}

// Transactions returns all eth transactions without loaded relations
// limited by passed parameters.
func (o *evmTxStore) Transactions(offset, limit int) (txs []Tx, count int, err error) {
	sql := `SELECT count(*) FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts)`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	var dbEthTxs []DbEthTx
	if err = o.q.Select(&dbEthTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxsToEvmEthTxs(dbEthTxs)
	return
}

// TransactionsWithAttempts returns all eth transactions with at least one attempt
// limited by passed parameters. Attempts are sorted by id.
func (o *evmTxStore) TransactionsWithAttempts(offset, limit int) (txs []Tx, count int, err error) {
	sql := `SELECT count(*) FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts)`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	var dbTxs []DbEthTx
	if err = o.q.Select(&dbTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxsToEvmEthTxs(dbTxs)
	err = o.preloadTxAttempts(txs)
	return
}

// TxAttempts returns the last tx attempts sorted by created_at descending.
func (o *evmTxStore) TxAttempts(offset, limit int) (txs []TxAttempt, count int, err error) {
	sql := `SELECT count(*) FROM eth_tx_attempts`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_tx_attempts ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	var dbTxs []DbEthTxAttempt
	if err = o.q.Select(&dbTxs, sql, limit, offset); err != nil {
		return
	}
	txs = dbEthTxAttemptsToEthTxAttempts(dbTxs)
	err = o.PreloadTxes(txs)
	return
}

// FindTxAttempt returns an individual TxAttempt
func (o *evmTxStore) FindTxAttempt(hash common.Hash) (*TxAttempt, error) {
	dbTxAttempt := DbEthTxAttempt{}
	sql := `SELECT * FROM eth_tx_attempts WHERE hash = $1`
	if err := o.q.Get(&dbTxAttempt, sql, hash); err != nil {
		return nil, err
	}
	// reuse the preload
	var attempt TxAttempt
	DbEthTxAttemptToEthTxAttempt(dbTxAttempt, &attempt)
	attempts := []TxAttempt{attempt}
	err := o.PreloadTxes(attempts)
	return &attempts[0], err
}

// FindTxAttemptsByTxIDs returns a list of attempts by ETH Tx IDs
func (o *evmTxStore) FindTxAttemptsByTxIDs(ids []int64) ([]TxAttempt, error) {
	sql := `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1)`
	var dbTxAttempts []DbEthTxAttempt
	if err := o.q.Select(&dbTxAttempts, sql, ids); err != nil {
		return nil, err
	}
	return dbEthTxAttemptsToEthTxAttempts(dbTxAttempts), nil
}

func (o *evmTxStore) FindTxByHash(hash common.Hash) (*Tx, error) {
	var dbEtx DbEthTx
	err := o.q.Transaction(func(tx pg.Queryer) error {
		sql := `SELECT eth_txes.* FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts WHERE hash = $1)`
		if err := tx.Get(&dbEtx, sql, hash); err != nil {
			return pkgerrors.Wrapf(err, "failed to find eth_tx with hash %d", hash)
		}
		return nil
	}, pg.OptReadOnlyTx())

	var etx Tx
	DbEthTxToEthTx(dbEtx, &etx)
	return &etx, pkgerrors.Wrap(err, "FindEthTxByHash failed")
}

// InsertTx inserts a new evm tx into the database
func (o *evmTxStore) InsertTx(etx *Tx) error {
	if etx.CreatedAt == (time.Time{}) {
		etx.CreatedAt = time.Now()
	}
	const insertEthTxSQL = `INSERT INTO eth_txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, initial_broadcast_at, created_at, state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, transmit_checker) VALUES (
:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :error, :broadcast_at, :initial_broadcast_at, :created_at, :state, :meta, :subject, :pipeline_task_run_id, :min_confirmations, :evm_chain_id, :transmit_checker
) RETURNING *`
	dbTx := DbEthTxFromEthTx(etx)
	err := o.q.GetNamed(insertEthTxSQL, &dbTx, &dbTx)
	DbEthTxToEthTx(dbTx, etx)
	return pkgerrors.Wrap(err, "InsertTx failed")
}

// InsertTxAttempt inserts a new txAttempt into the database
func (o *evmTxStore) InsertTxAttempt(attempt *TxAttempt) error {
	dbTxAttempt := DbEthTxAttemptFromEthTxAttempt(attempt)
	err := o.q.GetNamed(insertIntoEthTxAttemptsQuery, &dbTxAttempt, &dbTxAttempt)
	DbEthTxAttemptToEthTxAttempt(dbTxAttempt, attempt)
	return pkgerrors.Wrap(err, "InsertTxAttempt failed")
}

// InsertReceipt only used in tests. Use SaveFetchedReceipts instead
func (o *evmTxStore) InsertReceipt(receipt *evmtypes.Receipt) (int64, error) {
	// convert to database representation
	r := DbReceiptFromEvmReceipt(receipt)

	const insertEthReceiptSQL = `INSERT INTO eth_receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at) VALUES (
:tx_hash, :block_hash, :block_number, :transaction_index, :receipt, NOW()
) RETURNING *`
	err := o.q.GetNamed(insertEthReceiptSQL, &r, &r)

	return r.ID, pkgerrors.Wrap(err, "InsertReceipt failed")
}

// FindTxWithAttempts finds the Tx with its attempts and receipts preloaded
func (o *evmTxStore) FindTxWithAttempts(etxID int64) (etx Tx, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var dbEtx DbEthTx
		if err = tx.Get(&dbEtx, `SELECT * FROM eth_txes WHERE id = $1 ORDER BY created_at ASC, id ASC`, etxID); err != nil {
			return pkgerrors.Wrapf(err, "failed to find eth_tx with id %d", etxID)
		}
		DbEthTxToEthTx(dbEtx, &etx)
		if err = o.LoadTxAttempts(&etx, pg.WithQueryer(tx)); err != nil {
			return pkgerrors.Wrapf(err, "failed to load eth_tx_attempts for eth_tx with id %d", etxID)
		}
		if err = loadEthTxAttemptsReceipts(tx, &etx); err != nil {
			return pkgerrors.Wrapf(err, "failed to load eth_receipts for eth_tx with id %d", etxID)
		}
		return nil
	}, pg.OptReadOnlyTx())
	return etx, pkgerrors.Wrap(err, "FindTxWithAttempts failed")
}

func (o *evmTxStore) FindTxAttemptConfirmedByTxIDs(ids []int64) ([]TxAttempt, error) {
	var txAttempts []TxAttempt
	err := o.q.Transaction(func(tx pg.Queryer) error {
		var dbAttempts []DbEthTxAttempt
		if err := tx.Select(&dbAttempts, `SELECT eta.*
		FROM eth_tx_attempts eta
			join eth_receipts er on eta.hash = er.tx_hash where eta.eth_tx_id = ANY($1) ORDER BY eta.gas_price DESC, eta.gas_tip_cap DESC`, ids); err != nil {
			return err
		}
		txAttempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		return loadConfirmedAttemptsReceipts(tx, txAttempts)
	}, pg.OptReadOnlyTx())
	return txAttempts, pkgerrors.Wrap(err, "FindTxAttemptConfirmedByTxIDs failed")
}

func (o *evmTxStore) LoadTxesAttempts(etxs []*Tx, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	ethTxIDs := make([]int64, len(etxs))
	ethTxesM := make(map[int64]*Tx, len(etxs))
	for i, etx := range etxs {
		etx.TxAttempts = nil // this will overwrite any previous preload
		ethTxIDs[i] = etx.ID
		ethTxesM[etx.ID] = etxs[i]
	}
	var dbTxAttempts []DbEthTxAttempt
	if err := qq.Select(&dbTxAttempts, `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1) ORDER BY eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`, pq.Array(ethTxIDs)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxesAttempts failed to load eth_tx_attempts")
	}
	for _, dbAttempt := range dbTxAttempts {
		etx := ethTxesM[dbAttempt.EthTxID]
		var attempt TxAttempt
		DbEthTxAttemptToEthTxAttempt(dbAttempt, &attempt)
		etx.TxAttempts = append(etx.TxAttempts, attempt)
	}
	return nil
}

func (o *evmTxStore) LoadTxAttempts(etx *Tx, qopts ...pg.QOpt) error {
	return o.LoadTxesAttempts([]*Tx{etx}, qopts...)
}

func loadEthTxAttemptsReceipts(q pg.Queryer, etx *Tx) (err error) {
	return loadEthTxesAttemptsReceipts(q, []*Tx{etx})
}

func loadEthTxesAttemptsReceipts(q pg.Queryer, etxs []*Tx) (err error) {
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
	if err = q.Select(&rs, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(attemptHashes)); err != nil {
		return pkgerrors.Wrap(err, "loadEthTxesAttemptsReceipts failed to load eth_receipts")
	}

	var receipts []*evmtypes.Receipt = fromDBReceipts(rs)

	for _, receipt := range receipts {
		attempt := attemptHashM[receipt.TxHash]
		attempt.Receipts = append(attempt.Receipts, receipt)
	}
	return nil
}

func loadConfirmedAttemptsReceipts(q pg.Queryer, attempts []TxAttempt) error {
	byHash := make(map[string]*TxAttempt, len(attempts))
	hashes := make([][]byte, len(attempts))
	for i, attempt := range attempts {
		byHash[attempt.Hash.String()] = &attempts[i]
		hashes = append(hashes, attempt.Hash.Bytes())
	}
	var rs []dbReceipt
	if err := q.Select(&rs, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(hashes)); err != nil {
		return pkgerrors.Wrap(err, "loadConfirmedAttemptsReceipts failed to load eth_receipts")
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
func (o *evmTxStore) FindTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID *big.Int, address common.Address) (attempts []TxAttempt, err error) {
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	var dbAttempts []DbEthTxAttempt
	// this select distinct works because of unique index on eth_txes
	// (evm_chain_id, from_address, nonce)
	err = o.q.Select(&dbAttempts, `
SELECT DISTINCT ON (eth_txes.nonce) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= $1 AND evm_chain_id = $2 AND from_address = $3
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
LIMIT $4
`, olderThan, chainID.String(), address, limit)

	attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
	return attempts, pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringResend failed to load eth_tx_attempts")
}

func (o *evmTxStore) UpdateBroadcastAts(now time.Time, etxIDs []int64) error {
	// Deliberately do nothing on NULL broadcast_at because that indicates the
	// tx has been moved into a state where broadcast_at is not relevant, e.g.
	// fatally errored.
	//
	// Since EthConfirmer/EthResender can race (totally OK since highest
	// priced transaction always wins) we only want to update broadcast_at if
	// our version is later.
	_, err := o.q.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = ANY($2) AND broadcast_at < $1`, now, pq.Array(etxIDs))
	return pkgerrors.Wrap(err, "updateBroadcastAts failed to update eth_txes")
}

// SetBroadcastBeforeBlockNum updates already broadcast attempts with the
// current block number. This is safe no matter how old the head is because if
// the attempt is already broadcast it _must_ have been before this head.
func (o *evmTxStore) SetBroadcastBeforeBlockNum(blockNum int64, chainID *big.Int) error {
	_, err := o.q.Exec(
		`UPDATE eth_tx_attempts
SET broadcast_before_block_num = $1 
FROM eth_txes
WHERE eth_tx_attempts.broadcast_before_block_num IS NULL AND eth_tx_attempts.state = 'broadcast'
AND eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.evm_chain_id = $2`,
		blockNum, chainID.String(),
	)
	return pkgerrors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
}

func (o *evmTxStore) FindTxAttemptsConfirmedMissingReceipt(chainID *big.Int) (attempts []TxAttempt, err error) {
	var dbAttempts []DbEthTxAttempt
	err = o.q.Select(&dbAttempts,
		`SELECT DISTINCT ON (eth_tx_attempts.eth_tx_id) eth_tx_attempts.*
		FROM eth_tx_attempts
		JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state = 'confirmed_missing_receipt'
		WHERE evm_chain_id = $1
		ORDER BY eth_tx_attempts.eth_tx_id ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`,
		chainID.String())
	if err != nil {
		err = pkgerrors.Wrap(err, "FindEtxAttemptsConfirmedMissingReceipt failed to query")
	}
	attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
	return
}

func (o *evmTxStore) UpdateTxsUnconfirmed(ids []int64) error {
	_, err := o.q.Exec(`UPDATE eth_txes SET state='unconfirmed' WHERE id = ANY($1)`, pq.Array(ids))

	if err != nil {
		return pkgerrors.Wrap(err, "UpdateEthTxsUnconfirmed failed to execute")
	}
	return nil
}

func (o *evmTxStore) FindTxAttemptsRequiringReceiptFetch(chainID *big.Int) (attempts []TxAttempt, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var dbAttempts []DbEthTxAttempt
		err = tx.Select(&dbAttempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt') AND eth_txes.evm_chain_id = $1
WHERE eth_tx_attempts.state != 'insufficient_eth'
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
`, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load eth_tx_attempts")
		}
		attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		err = o.PreloadTxes(attempts, pg.WithQueryer(tx))
		return pkgerrors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return
}

func (o *evmTxStore) SaveFetchedReceipts(r []*evmtypes.Receipt, chainID *big.Int) (err error) {
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
		INSERT INTO eth_receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at)
		VALUES %s
		ON CONFLICT (tx_hash, block_hash) DO UPDATE SET
			block_number = EXCLUDED.block_number,
			transaction_index = EXCLUDED.transaction_index,
			receipt = EXCLUDED.receipt
		RETURNING eth_receipts.tx_hash, eth_receipts.block_number
	),
	updated_eth_tx_attempts AS (
		UPDATE eth_tx_attempts
		SET
			state = 'broadcast',
			broadcast_before_block_num = COALESCE(eth_tx_attempts.broadcast_before_block_num, inserted_receipts.block_number)
		FROM inserted_receipts
		WHERE inserted_receipts.tx_hash = eth_tx_attempts.hash
		RETURNING eth_tx_attempts.eth_tx_id
	)
	UPDATE eth_txes
	SET state = 'confirmed'
	FROM updated_eth_tx_attempts
	WHERE updated_eth_tx_attempts.eth_tx_id = eth_txes.id
	AND evm_chain_id = ?
	`

	stmt := fmt.Sprintf(sql, strings.Join(valueStrs, ","))

	stmt = sqlx.Rebind(sqlx.DOLLAR, stmt)

	err = o.q.ExecQ(stmt, valueArgs...)
	return pkgerrors.Wrap(err, "SaveFetchedReceipts failed to save receipts")
}

// MarkAllConfirmedMissingReceipt
// It is possible that we can fail to get a receipt for all eth_tx_attempts
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
// NOTE: We continue to attempt to resend eth_txes in this state on
// every head to guard against the extremely rare scenario of nonce gap due to
// reorg that excludes the transaction (from another wallet) that had this
// nonce (until finality depth is reached, after which we make the explicit
// decision to give up). This is done in the EthResender.
//
// We will continue to try to fetch a receipt for these attempts until all
// attempts are below the finality depth from current head.
func (o *evmTxStore) MarkAllConfirmedMissingReceipt(chainID *big.Int) (err error) {
	res, err := o.q.Exec(`
UPDATE eth_txes
SET state = 'confirmed_missing_receipt'
FROM (
	SELECT from_address, MAX(nonce) as max_nonce 
	FROM eth_txes
	WHERE state = 'confirmed' AND evm_chain_id = $1
	GROUP BY from_address
) AS max_table
WHERE state = 'unconfirmed'
	AND evm_chain_id = $1
	AND nonce < max_table.max_nonce
	AND eth_txes.from_address = max_table.from_address
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
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		var dbAttempts []DbEthTxAttempt
		err = tx.Select(&dbAttempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
WHERE eth_tx_attempts.state = 'in_progress' AND eth_txes.from_address = $1 AND eth_txes.evm_chain_id = $2
`, address, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_tx_attempts")
		}
		attempts = dbEthTxAttemptsToEthTxAttempts(dbAttempts)
		err = o.PreloadTxes(attempts, pg.WithQueryer(tx))
		return pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return attempts, pkgerrors.Wrap(err, "getInProgressEthTxAttempts failed")
}

func (o *evmTxStore) FindReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID *big.Int) (receiptsPlus []ReceiptPlus, err error) {
	var rs []dbReceiptPlus

	err = o.q.SelectContext(ctx, &rs, `
	SELECT pipeline_task_runs.id, eth_receipts.receipt, COALESCE((eth_txes.meta->>'FailOnRevert')::boolean, false) "FailOnRevert" FROM pipeline_task_runs
	INNER JOIN pipeline_runs ON pipeline_runs.id = pipeline_task_runs.pipeline_run_id
	INNER JOIN eth_txes ON eth_txes.pipeline_task_run_id = pipeline_task_runs.id
	INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id
	INNER JOIN eth_receipts ON eth_tx_attempts.hash = eth_receipts.tx_hash
	WHERE pipeline_runs.state = 'suspended' AND eth_receipts.block_number <= ($1 - eth_txes.min_confirmations) AND eth_txes.evm_chain_id = $2
	`, blockNum, chainID.String())

	receiptsPlus = fromDBReceiptsPlus(rs)
	return
}

// FindTxWithSequence returns any broadcast ethtx with the given nonce
func (o *evmTxStore) FindTxWithSequence(fromAddress common.Address, nonce evmtypes.Nonce) (etx *Tx, err error) {
	etx = new(Tx)
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var dbEtx DbEthTx
		err = tx.Get(&dbEtx, `
SELECT * FROM eth_txes WHERE from_address = $1 AND nonce = $2 AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
`, fromAddress, nonce.Int64())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxWithNonce failed to load eth_txes")
		}
		DbEthTxToEthTx(dbEtx, etx)
		err = o.LoadTxAttempts(etx, pg.WithQueryer(tx))
		return pkgerrors.Wrap(err, "FindEthTxWithNonce failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

func updateEthTxAttemptUnbroadcast(q pg.Queryer, attempt TxAttempt) error {
	if attempt.State != txmgrtypes.TxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	_, err := q.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = $1`, attempt.ID)
	return pkgerrors.Wrap(err, "updateEthTxAttemptUnbroadcast failed")
}

func updateEthTxUnconfirm(q pg.Queryer, etx Tx) error {
	if etx.State != txmgr.TxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	_, err := q.Exec(`UPDATE eth_txes SET state = 'unconfirmed' WHERE id = $1`, etx.ID)
	return pkgerrors.Wrap(err, "updateEthTxUnconfirm failed")
}

func deleteEthReceipts(q pg.Queryer, etxID int64) (err error) {
	_, err = q.Exec(`
DELETE FROM eth_receipts
USING eth_tx_attempts
WHERE eth_receipts.tx_hash = eth_tx_attempts.hash
AND eth_tx_attempts.eth_tx_id = $1
	`, etxID)
	return pkgerrors.Wrap(err, "deleteEthReceipts failed")
}

func (o *evmTxStore) UpdateTxForRebroadcast(etx Tx, etxAttempt TxAttempt) error {
	return o.q.Transaction(func(tx pg.Queryer) error {
		if err := deleteEthReceipts(tx, etx.ID); err != nil {
			return pkgerrors.Wrapf(err, "deleteEthReceipts failed for etx %v", etx.ID)
		}
		if err := updateEthTxUnconfirm(tx, etx); err != nil {
			return pkgerrors.Wrapf(err, "updateEthTxUnconfirm failed for etx %v", etx.ID)
		}
		return updateEthTxAttemptUnbroadcast(tx, etxAttempt)
	})
}

func (o *evmTxStore) FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID *big.Int) (etxs []*Tx, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var dbEtxs []DbEthTx
		err = tx.Select(&dbEtxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'broadcast'
INNER JOIN eth_receipts ON eth_receipts.tx_hash = eth_tx_attempts.hash
WHERE eth_txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number BETWEEN $1 AND $2 AND evm_chain_id = $3
ORDER BY nonce ASC
`, lowBlockNumber, highBlockNumber, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		if err = o.LoadTxesAttempts(etxs, pg.WithQueryer(tx)); err != nil {
			return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_tx_attempts")
		}
		err = loadEthTxesAttemptsReceipts(tx, etxs)
		return pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_receipts")
	}, pg.OptReadOnlyTx())
	return etxs, pkgerrors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed")
}

func saveAttemptWithNewState(q pg.Queryer, timeout time.Duration, logger logger.Logger, attempt TxAttempt, broadcastAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	dbAttempt := DbEthTxAttemptFromEthTxAttempt(&attempt)
	defer cancel()
	return pg.SqlxTransaction(ctx, q, logger, func(tx pg.Queryer) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if _, err := tx.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = $2 AND broadcast_at < $1`, broadcastAt, dbAttempt.EthTxID); err != nil {
			return pkgerrors.Wrap(err, "saveAttemptWithNewState failed to update eth_txes")
		}
		_, err := tx.Exec(`UPDATE eth_tx_attempts SET state=$1 WHERE id=$2`, dbAttempt.State, dbAttempt.ID)
		return pkgerrors.Wrap(err, "saveAttemptWithNewState failed to update eth_tx_attempts")
	})
}

func (o *evmTxStore) SaveInsufficientFundsAttempt(timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	if !(attempt.State == txmgrtypes.TxAttemptInProgress || attempt.State == txmgrtypes.TxAttemptInsufficientFunds) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = txmgrtypes.TxAttemptInsufficientFunds
	return pkgerrors.Wrap(saveAttemptWithNewState(o.q, timeout, o.logger, *attempt, broadcastAt), "saveInsufficientEthAttempt failed")
}

func saveSentAttempt(q pg.Queryer, timeout time.Duration, logger logger.Logger, attempt *TxAttempt, broadcastAt time.Time) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = txmgrtypes.TxAttemptBroadcast
	return pkgerrors.Wrap(saveAttemptWithNewState(q, timeout, logger, *attempt, broadcastAt), "saveSentAttempt failed")
}

func (o *evmTxStore) SaveSentAttempt(timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	return saveSentAttempt(o.q, timeout, o.logger, attempt, broadcastAt)
}

func (o *evmTxStore) SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt, broadcastAt time.Time) error {
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))
	err := qq.Transaction(func(tx pg.Queryer) error {
		if err := saveSentAttempt(tx, timeout, o.logger, attempt, broadcastAt); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE eth_txes SET state = 'confirmed_missing_receipt' WHERE id = $1`, attempt.TxID); err != nil {
			return pkgerrors.Wrap(err, "failed to update eth_txes")

		}
		return nil
	})
	return pkgerrors.Wrap(err, "SaveConfirmedMissingReceiptAttempt failed")
}

func (o *evmTxStore) DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt) error {
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))

	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("DeleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("DeleteInProgressAttempt: expected attempt to have an id")
	}
	_, err := qq.Exec(`DELETE FROM eth_tx_attempts WHERE id = $1`, attempt.ID)
	return pkgerrors.Wrap(err, "DeleteInProgressAttempt failed")
}

// SaveInProgressAttempt inserts or updates an attempt
func (o *evmTxStore) SaveInProgressAttempt(attempt *TxAttempt) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("SaveInProgressAttempt failed: attempt state must be in_progress")
	}
	dbAttempt := DbEthTxAttemptFromEthTxAttempt(attempt)
	// Insert is the usual mode because the attempt is new
	if attempt.ID == 0 {
		query, args, e := o.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if e != nil {
			return pkgerrors.Wrap(e, "SaveInProgressAttempt failed to BindNamed")
		}
		e = o.q.Get(&dbAttempt, query, args...)
		DbEthTxAttemptToEthTxAttempt(dbAttempt, attempt)
		return pkgerrors.Wrap(e, "SaveInProgressAttempt failed to insert into eth_tx_attempts")
	}
	// Update only applies to case of insufficient eth and simply changes the state to in_progress
	res, err := o.q.Exec(`UPDATE eth_tx_attempts SET state=$1, broadcast_before_block_num=$2 WHERE id=$3`, dbAttempt.State, dbAttempt.BroadcastBeforeBlockNum, dbAttempt.ID)
	if err != nil {
		return pkgerrors.Wrap(err, "SaveInProgressAttempt failed to update eth_tx_attempts")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return pkgerrors.Wrap(err, "SaveInProgressAttempt failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return pkgerrors.Wrapf(sql.ErrNoRows, "SaveInProgressAttempt tried to update eth_tx_attempts but no rows matched id %d", attempt.ID)
	}
	return nil
}

// FindTxsRequiringGasBump returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
//
// It also returns eth_txes that are unconfirmed with no eth_tx_attempts
func (o *evmTxStore) FindTxsRequiringGasBump(ctx context.Context, address common.Address, blockNum, gasBumpThreshold, depth int64, chainID *big.Int) (etxs []*Tx, err error) {
	if gasBumpThreshold == 0 {
		return
	}
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		stmt := `
SELECT eth_txes.* FROM eth_txes
LEFT JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND (broadcast_before_block_num > $4 OR broadcast_before_block_num IS NULL OR eth_tx_attempts.state != 'broadcast')
WHERE eth_txes.state = 'unconfirmed' AND eth_tx_attempts.id IS NULL AND eth_txes.from_address = $1 AND eth_txes.evm_chain_id = $2
	AND (($3 = 0) OR (eth_txes.id IN (SELECT id FROM eth_txes WHERE state = 'unconfirmed' AND from_address = $1 ORDER BY nonce ASC LIMIT $3)))
ORDER BY nonce ASC
`
		var dbEtxs []DbEthTx
		if err = tx.Select(&dbEtxs, stmt, address, chainID.String(), depth, blockNum-gasBumpThreshold); err != nil {
			return pkgerrors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		err = o.LoadTxesAttempts(etxs, pg.WithQueryer(tx))
		return pkgerrors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	return
}

// FindTxsRequiringResubmissionDueToInsufficientFunds returns transactions
// that need to be re-sent because they hit an out-of-eth error on a previous
// block
func (o *evmTxStore) FindTxsRequiringResubmissionDueToInsufficientFunds(address common.Address, chainID *big.Int, qopts ...pg.QOpt) (etxs []*Tx, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Transaction(func(tx pg.Queryer) error {
		var dbEtxs []DbEthTx
		err = tx.Select(&dbEtxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'insufficient_eth'
WHERE eth_txes.from_address = $1 AND eth_txes.state = 'unconfirmed' AND eth_txes.evm_chain_id = $2
ORDER BY nonce ASC
`, address, chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_txes")
		}
		etxs = make([]*Tx, len(dbEtxs))
		dbEthTxsToEvmEthTxPtrs(dbEtxs, etxs)
		err = o.LoadTxesAttempts(etxs, pg.WithQueryer(tx))
		return pkgerrors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	return
}

// markOldTxesMissingReceiptAsErrored
//
// Once eth_tx has all of its attempts broadcast before some cutoff threshold
// without receiving any receipts, we mark it as fatally errored (never sent).
//
// The job run will also be marked as errored in this case since we never got a
// receipt and thus cannot pass on any transaction hash
func (o *evmTxStore) MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID *big.Int, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
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
	return qq.Transaction(func(q pg.Queryer) error {
		type etx struct {
			ID    int64
			Nonce int64
		}
		var data []etx
		err := q.Select(&data, `
UPDATE eth_txes
SET state='fatal_error', nonce=NULL, error=$1, broadcast_at=NULL, initial_broadcast_at=NULL
FROM (
	SELECT e1.id, e1.nonce, e1.from_address FROM eth_txes AS e1 WHERE id IN (
		SELECT e2.id FROM eth_txes AS e2
		INNER JOIN eth_tx_attempts ON e2.id = eth_tx_attempts.eth_tx_id
		WHERE e2.state = 'confirmed_missing_receipt'
		AND e2.evm_chain_id = $3
		GROUP BY e2.id
		HAVING max(eth_tx_attempts.broadcast_before_block_num) < $2
	)
	FOR UPDATE OF e1
) e0
WHERE e0.id = eth_txes.id
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
		err = q.Select(&results, `
SELECT e.id, e.from_address, max(a.broadcast_before_block_num) AS max_broadcast_before_block_num, array_agg(a.hash) AS tx_hashes
FROM eth_txes e
INNER JOIN eth_tx_attempts a ON e.id = a.eth_tx_id
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

			o.logger.Criticalw(fmt.Sprintf("eth_tx with ID %v expired without ever getting a receipt for any of our attempts. "+
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

func (o *evmTxStore) SaveReplacementInProgressAttempt(oldAttempt TxAttempt, replacementAttempt *TxAttempt, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	if oldAttempt.State != txmgrtypes.TxAttemptInProgress || replacementAttempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return qq.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE id=$1`, oldAttempt.ID); err != nil {
			return pkgerrors.Wrap(err, "saveReplacementInProgressAttempt failed to delete from eth_tx_attempts")
		}
		dbAttempt := DbEthTxAttemptFromEthTxAttempt(replacementAttempt)
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if e != nil {
			return pkgerrors.Wrap(e, "saveReplacementInProgressAttempt failed to BindNamed")
		}
		e = tx.Get(&dbAttempt, query, args...)
		DbEthTxAttemptToEthTxAttempt(dbAttempt, replacementAttempt)
		return pkgerrors.Wrap(e, "saveReplacementInProgressAttempt failed to insert replacement attempt")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func (o *evmTxStore) FindNextUnstartedTransactionFromAddress(etx *Tx, fromAddress common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	var dbEtx DbEthTx
	err := qq.Get(&dbEtx, `SELECT * FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2 ORDER BY value ASC, created_at ASC, id ASC`, fromAddress, chainID.String())
	DbEthTxToEthTx(dbEtx, etx)
	return pkgerrors.Wrap(err, "failed to FindNextUnstartedTransactionFromAddress")
}

func (o *evmTxStore) UpdateTxFatalError(etx *Tx, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)

	if etx.State != txmgr.TxInProgress {
		return pkgerrors.Errorf("can only transition to fatal_error from in_progress, transaction is currently %s", etx.State)
	}
	if !etx.Error.Valid {
		return errors.New("expected error field to be set")
	}

	etx.Sequence = nil
	etx.State = txmgr.TxFatalError

	return qq.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = $1`, etx.ID); err != nil {
			return pkgerrors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		dbEtx := DbEthTxFromEthTx(etx)
		err := pkgerrors.Wrap(tx.Get(&dbEtx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=NULL, initial_broadcast_at=NULL, nonce=NULL WHERE id=$3 RETURNING *`, etx.State, etx.Error, etx.ID), "saveFatallyErroredTransaction failed to save eth_tx")
		DbEthTxToEthTx(dbEtx, etx)
		return err
	})
}

// Updates eth attempt from in_progress to broadcast. Also updates the eth tx to unconfirmed.
// Before it updates both tables though it increments the next nonce from the keystore
// One of the more complicated signatures. We have to accept variable pg.QOpt and QueryerFunc arguments
func (o *evmTxStore) UpdateTxAttemptInProgressToBroadcast(etx *Tx, attempt TxAttempt, NewAttemptState txmgrtypes.TxAttemptState, incrNextNonceCallback txmgrtypes.QueryerFunc, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)

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
	return qq.Transaction(func(tx pg.Queryer) error {
		if err := incrNextNonceCallback(tx); err != nil {
			return pkgerrors.Wrap(err, "SaveEthTxAttempt failed on incrNextNonceCallback")
		}
		dbEtx := DbEthTxFromEthTx(etx)
		if err := tx.Get(&dbEtx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id = $5 RETURNING *`, dbEtx.State, dbEtx.Error, dbEtx.BroadcastAt, dbEtx.InitialBroadcastAt, dbEtx.ID); err != nil {
			return pkgerrors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx")
		}
		DbEthTxToEthTx(dbEtx, etx)
		dbAttempt := DbEthTxAttemptFromEthTxAttempt(&attempt)
		if err := tx.Get(&dbAttempt, `UPDATE eth_tx_attempts SET state = $1 WHERE id = $2 RETURNING *`, dbAttempt.State, dbAttempt.ID); err != nil {
			return pkgerrors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx_attempt")
		}
		return nil
	})
}

// Updates eth tx from unstarted to in_progress and inserts in_progress eth attempt
func (o *evmTxStore) UpdateTxUnstartedToInProgress(etx *Tx, attempt *TxAttempt, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
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
	return qq.Transaction(func(tx pg.Queryer) error {
		// If a replay was triggered while unconfirmed transactions were pending, they will be marked as fatal_error => abandoned.
		// In this case, we must remove the abandoned attempt from eth_tx_attempts before replacing it with a new one.  In any other
		// case, we uphold the constraint, leaving the original tx attempt as-is and returning the constraint violation error.
		//
		// Note:  the record of the original abandoned transaction will remain in eth_txes, only the attempt is replaced.  (Any receipt
		// associated with the abandoned attempt would also be lost, although this shouldn't happen since only unconfirmed transactions
		// can be abandoned.)
		_, err := tx.Exec(`DELETE FROM eth_tx_attempts a USING eth_txes t
			WHERE t.id = a.eth_tx_id AND a.hash = $1 AND t.state = $2 AND t.error = 'abandoned'`,
			attempt.Hash, txmgr.TxFatalError,
		)
		if err == nil {
			o.logger.Debugf("Replacing abandoned tx with tx hash %s with tx_id=%d with identical tx hash", attempt.Hash, attempt.TxID)
		} else if errors.Is(err, sql.ErrNoRows) {
			return err
		}

		dbAttempt := DbEthTxAttemptFromEthTxAttempt(attempt)
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, &dbAttempt)
		if e != nil {
			return pkgerrors.Wrap(e, "failed to BindNamed")
		}
		err = tx.Get(&dbAttempt, query, args...)
		if err != nil {
			var pqErr *pgconn.PgError
			if isPqErr := errors.As(err, &pqErr); isPqErr {
				switch pqErr.ConstraintName {
				case "eth_tx_attempts_eth_tx_id_fkey":
					return txmgr.ErrTxRemoved
				default:
				}
			}
			if err != nil {
				return pkgerrors.Wrap(err, "UpdateTxUnstartedToInProgress failed to create eth_tx_attempt")
			}
		}
		DbEthTxAttemptToEthTxAttempt(dbAttempt, attempt)
		dbEtx := DbEthTxFromEthTx(etx)
		err = tx.Get(&dbEtx, `UPDATE eth_txes SET nonce=$1, state=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id=$5 RETURNING *`, etx.Sequence, etx.State, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID)
		DbEthTxToEthTx(dbEtx, etx)
		return pkgerrors.Wrap(err, "UpdateTxUnstartedToInProgress failed to update eth_tx")
	})
}

// GetTxInProgress returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnstartedEthTxs loop.
// It may or may not have been broadcast to an eth node.
func (o *evmTxStore) GetTxInProgress(fromAddress common.Address, qopts ...pg.QOpt) (etx *Tx, err error) {
	qq := o.q.WithOpts(qopts...)
	etx = new(Tx)
	if err != nil {
		return etx, pkgerrors.Wrap(err, "getInProgressEthTx failed")
	}
	err = qq.Transaction(func(tx pg.Queryer) error {
		var dbEtx DbEthTx
		err = tx.Get(&dbEtx, `SELECT * FROM eth_txes WHERE from_address = $1 and state = 'in_progress'`, fromAddress)
		if errors.Is(err, sql.ErrNoRows) {
			etx = nil
			return nil
		} else if err != nil {
			return pkgerrors.Wrap(err, "GetTxInProgress failed while loading eth tx")
		}
		DbEthTxToEthTx(dbEtx, etx)
		if err = o.LoadTxAttempts(etx, pg.WithQueryer(tx)); err != nil {
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

func (o *evmTxStore) HasInProgressTransaction(account common.Address, chainID *big.Int, qopts ...pg.QOpt) (exists bool, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Get(&exists, `SELECT EXISTS(SELECT 1 FROM eth_txes WHERE state = 'in_progress' AND from_address = $1 AND evm_chain_id = $2)`, account, chainID.String())
	return exists, pkgerrors.Wrap(err, "hasInProgressTransaction failed")
}

func (o *evmTxStore) UpdateKeyNextSequence(newNextNonce, currentNextNonce evmtypes.Nonce, address common.Address, chainID *big.Int, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	return qq.Transaction(func(tx pg.Queryer) error {
		//  We filter by next_nonce here as an optimistic lock to make sure it
		//  didn't get changed out from under us. Shouldn't happen but can't hurt.
		res, err := tx.Exec(`UPDATE evm_key_states SET next_nonce = $1, updated_at = $2 WHERE address = $3 AND next_nonce = $4 AND evm_chain_id = $5`, newNextNonce.Int64(), time.Now(), address, currentNextNonce.Int64(), chainID.String())
		if err != nil {
			return pkgerrors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to update keys.next_nonce")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return pkgerrors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to get RowsAffected")
		}
		if rowsAffected == 0 {
			return ErrKeyNotUpdated
		}
		return nil
	})
}

func (o *evmTxStore) countTransactionsWithState(fromAddress common.Address, state txmgrtypes.TxState, chainID *big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = $2 AND evm_chain_id = $3`,
		fromAddress, state, chainID.String())
	return count, pkgerrors.Wrap(err, "failed to countTransactionsWithState")
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func (o *evmTxStore) CountUnconfirmedTransactions(fromAddress common.Address, chainID *big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	return o.countTransactionsWithState(fromAddress, txmgr.TxUnconfirmed, chainID, qopts...)
}

// CountUnstartedTransactions returns the number of unconfirmed transactions
func (o *evmTxStore) CountUnstartedTransactions(fromAddress common.Address, chainID *big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	return o.countTransactionsWithState(fromAddress, txmgr.TxUnstarted, chainID, qopts...)
}

func (o *evmTxStore) CheckTxQueueCapacity(fromAddress common.Address, maxQueuedTransactions uint64, chainID *big.Int, qopts ...pg.QOpt) (err error) {
	qq := o.q.WithOpts(qopts...)
	if maxQueuedTransactions == 0 {
		return nil
	}
	var count uint64
	err = qq.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2`, fromAddress, chainID.String())
	if err != nil {
		err = pkgerrors.Wrap(err, "CheckTxQueueCapacity query failed")
		return
	}

	if count >= maxQueuedTransactions {
		err = pkgerrors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}
	return
}

func (o *evmTxStore) CreateTransaction(txRequest TxRequest, chainID *big.Int, qopts ...pg.QOpt) (tx Tx, err error) {
	var dbEtx DbEthTx
	qq := o.q.WithOpts(qopts...)
	err = qq.Transaction(func(tx pg.Queryer) error {
		if txRequest.PipelineTaskRunID != nil {

			err = tx.Get(&dbEtx, `SELECT * FROM eth_txes WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, txRequest.PipelineTaskRunID, chainID.String())
			// If no eth_tx matches (the common case) then continue
			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return pkgerrors.Wrap(err, "CreateEthTransaction")
				}
				// if a previous transaction for this task run exists, immediately return it
				return nil
			}
		}
		err = tx.Get(&dbEtx, `
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, transmit_checker)
VALUES (
$1,$2,$3,$4,$5,'unstarted',NOW(),$6,$7,$8,$9,$10,$11
)
RETURNING "eth_txes".*
`, txRequest.FromAddress, txRequest.ToAddress, txRequest.EncodedPayload, assets.Eth(txRequest.Value), txRequest.FeeLimit, txRequest.Meta, txRequest.Strategy.Subject(), chainID.String(), txRequest.MinConfirmations, txRequest.PipelineTaskRunID, txRequest.Checker)
		if err != nil {
			return pkgerrors.Wrap(err, "CreateEthTransaction failed to insert eth_tx")
		}
		var pruned int64
		pruned, err = txRequest.Strategy.PruneQueue(o, pg.WithQueryer(tx))
		if err != nil {
			return pkgerrors.Wrap(err, "CreateEthTransaction failed to prune eth_txes")
		}
		if pruned > 0 {
			o.logger.Warnw(fmt.Sprintf("Dropped %d old transactions from transaction queue", pruned), "fromAddress", txRequest.FromAddress, "toAddress", txRequest.ToAddress, "meta", txRequest.Meta, "subject", txRequest.Strategy.Subject(), "replacementID", dbEtx.ID)
		}
		return nil
	})
	var etx Tx
	DbEthTxToEthTx(dbEtx, &etx)
	return etx, err
}

func (o *evmTxStore) PruneUnstartedTxQueue(queueSize uint32, subject uuid.UUID, qopts ...pg.QOpt) (n int64, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Transaction(func(tx pg.Queryer) error {
		res, err := qq.Exec(`
DELETE FROM eth_txes
WHERE state = 'unstarted' AND subject = $1 AND
id < (
	SELECT min(id) FROM (
		SELECT id
		FROM eth_txes
		WHERE state = 'unstarted' AND subject = $2
		ORDER BY id DESC
		LIMIT $3
	) numbers
)`, subject, subject, queueSize)
		if err != nil {
			return pkgerrors.Wrap(err, "DeleteUnstartedEthTx failed")
		}
		n, err = res.RowsAffected()
		return err
	})
	return
}

func (o *evmTxStore) ReapTxHistory(minBlockNumberToKeep int64, timeThreshold time.Time, chainID *big.Int) error {
	// Delete old confirmed eth_txes
	// NOTE that this relies on foreign key triggers automatically removing
	// the eth_tx_attempts and eth_receipts linked to every eth_tx
	err := pg.Batch(func(_, limit uint) (count uint, err error) {
		res, err := o.q.Exec(`
WITH old_enough_receipts AS (
	SELECT tx_hash FROM eth_receipts
	WHERE block_number < $1
	ORDER BY block_number ASC, id ASC
	LIMIT $2
)
DELETE FROM eth_txes
USING old_enough_receipts, eth_tx_attempts
WHERE eth_tx_attempts.eth_tx_id = eth_txes.id
AND eth_tx_attempts.hash = old_enough_receipts.tx_hash
AND eth_txes.created_at < $3
AND eth_txes.state = 'confirmed'
AND evm_chain_id = $4`, minBlockNumberToKeep, limit, timeThreshold, chainID.String())
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to delete old confirmed eth_txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	})
	if err != nil {
		return pkgerrors.Wrap(err, "TxmReaper#reapEthTxes batch delete of confirmed eth_txes failed")
	}
	// Delete old 'fatal_error' eth_txes
	err = pg.Batch(func(_, limit uint) (count uint, err error) {
		res, err := o.q.Exec(`
DELETE FROM eth_txes
WHERE created_at < $1
AND state = 'fatal_error'
AND evm_chain_id = $2`, timeThreshold, chainID.String())
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to delete old fatally errored eth_txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, pkgerrors.Wrap(err, "ReapTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	})
	if err != nil {
		return pkgerrors.Wrap(err, "TxmReaper#reapEthTxes batch delete of fatally errored eth_txes failed")
	}

	return nil
}

func (o *evmTxStore) Abandon(chainID *big.Int, addr common.Address) error {
	_, err := o.q.Exec(`UPDATE eth_txes SET state='fatal_error', nonce = NULL, error = 'abandoned' WHERE state IN ('unconfirmed', 'in_progress', 'unstarted') AND evm_chain_id = $1 AND from_address = $2`, chainID.String(), addr)
	return err
}
