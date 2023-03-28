package txmgr

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/common/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

var ErrKeyNotUpdated = errors.New("evmTxStorageService: Key not updated")
var ErrInvalidQOpt = errors.New("evmTxStorageService: Invalid QOpt")

type evmTxStorageService struct {
	txmgrtypes.TxStorageService[*evmtypes.Address, big.Int, *evmtypes.TxHash, *evmtypes.BlockHash, NewTx[*evmtypes.Address], *evmtypes.Receipt, EthTx[*evmtypes.Address, *evmtypes.TxHash], EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], int64, int64]
	q         pg.Q
	logger    logger.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc
}

var _ txmgrtypes.TxStorageService[*evmtypes.Address, big.Int, *evmtypes.TxHash, *evmtypes.BlockHash, NewTx[*evmtypes.Address], *evmtypes.Receipt, EthTx[*evmtypes.Address, *evmtypes.TxHash], EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], int64, int64] = &evmTxStorageService{}

// Directly maps to database schema.
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

func dbReceiptFromEvmReceipt(evmReceipt *EvmReceipt) dbReceipt {
	txHash, _ := evmReceipt.TxHash.MarshalText()
	blockHash, _ := evmReceipt.BlockHash.MarshalText()

	return dbReceipt{
		ID:               evmReceipt.ID,
		TxHash:           common.BytesToHash(txHash),
		BlockHash:        common.BytesToHash(blockHash),
		BlockNumber:      evmReceipt.BlockNumber,
		TransactionIndex: evmReceipt.TransactionIndex,
		Receipt:          *evmReceipt.Receipt,
		CreatedAt:        evmReceipt.CreatedAt,
	}
}

func dbReceiptToEvmReceipt(receipt *dbReceipt) EvmReceipt {
	txHash := evmtypes.NewTxHash(receipt.TxHash)
	blockHash := evmtypes.NewBlockHash(receipt.TxHash)

	return EvmReceipt{
		ID:               receipt.ID,
		TxHash:           txHash,
		BlockHash:        blockHash,
		BlockNumber:      receipt.BlockNumber,
		TransactionIndex: receipt.TransactionIndex,
		Receipt:          &receipt.Receipt,
		CreatedAt:        receipt.CreatedAt,
	}
}

// Directly maps to database schema.
// Do not modify type unless you
// intend to modify the database schema
type dbRawReceipt = evmtypes.Receipt

// Directly maps to database schema.
// Do not modify type unless you
// intend to modify the database schema
//
// Does not map to a single database table.
// It's comprised of fields from different tables.
type dbReceiptPlus struct {
	ID           uuid.UUID        `db:"id"`
	Receipt      evmtypes.Receipt `db:"receipt"`
	FailOnRevert bool             `db:"FailOnRevert"`
}

func fromDBReceipts(rs []dbReceipt) []EvmReceipt {
	receipts := make([]EvmReceipt, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = dbReceiptToEvmReceipt(&rs[i])
	}
	return receipts
}

func fromDBReceiptsPlus(rs []dbReceiptPlus) []EvmReceiptPlus {
	receipts := make([]EvmReceiptPlus, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = EvmReceiptPlus{
			ID:           rs[i].ID,
			Receipt:      &rs[i].Receipt,
			FailOnRevert: rs[i].FailOnRevert,
		}
	}
	return receipts
}

func toDBRawReceipt(rs []*evmtypes.Receipt) []dbRawReceipt {
	receipts := make([]dbRawReceipt, len(rs))
	for i := 0; i < len(rs); i++ {
		receipts[i] = dbRawReceipt(*rs[i])
	}
	return receipts
}

func NewTxStorageService(
	db *sqlx.DB,
	lggr logger.Logger,
	cfg pg.QConfig,
) txmgrtypes.TxStorageService[*evmtypes.Address, big.Int, *evmtypes.TxHash, *evmtypes.BlockHash, NewTx[*evmtypes.Address], *evmtypes.Receipt, EthTx[*evmtypes.Address, *evmtypes.TxHash], EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], int64, int64] {
	namedLogger := lggr.Named("TxmORM")
	ctx, cancel := context.WithCancel(context.Background())
	q := pg.NewQ(db, namedLogger, cfg, pg.WithParentCtx(ctx))
	return &evmTxStorageService{
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

// TODO: create method to pass in new context to evmTxStorageService (which will also create a new pg.Q)

func (o *evmTxStorageService) Close() {
	o.ctxCancel()
}

func (o *evmTxStorageService) preloadTxAttempts(txs []EthTx[*evmtypes.Address, *evmtypes.TxHash]) error {
	// Preload TxAttempts
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	var attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]
	sql := `SELECT * FROM eth_tx_attempts WHERE eth_tx_id IN (?) ORDER BY id desc;`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = o.q.Rebind(query)
	if err = o.q.Select(&attempts, query, args...); err != nil {
		return err
	}
	// fill in attempts
	for _, attempt := range attempts {
		for i, tx := range txs {
			if tx.ID == attempt.EthTxID {
				txs[i].EthTxAttempts = append(txs[i].EthTxAttempts, attempt)
			}
		}
	}
	return nil
}

func (o *evmTxStorageService) PreloadEthTxes(attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	ethTxM := make(map[int64]EthTx[*evmtypes.Address, *evmtypes.TxHash])
	for _, attempt := range attempts {
		ethTxM[attempt.EthTxID] = EthTx[*evmtypes.Address, *evmtypes.TxHash]{}
	}
	ethTxIDs := make([]int64, len(ethTxM))
	var i int
	for id := range ethTxM {
		ethTxIDs[i] = id
		i++
	}
	ethTxs := make([]EthTx[*evmtypes.Address, *evmtypes.TxHash], len(ethTxIDs))
	if err := o.q.Select(&ethTxs, `SELECT * FROM eth_txes WHERE id = ANY($1)`, pq.Array(ethTxIDs)); err != nil {
		return errors.Wrap(err, "loadEthTxes failed")
	}
	for _, etx := range ethTxs {
		ethTxM[etx.ID] = etx
	}
	for i, attempt := range attempts {
		attempts[i].EthTx = ethTxM[attempt.EthTxID]
	}
	return nil
}

// EthTransactions returns all eth transactions without loaded relations
// limited by passed parameters.
func (o *evmTxStorageService) EthTransactions(offset, limit int) (txs []EthTx[*evmtypes.Address, *evmtypes.TxHash], count int, err error) {
	sql := `SELECT count(*) FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts)`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	if err = o.q.Select(&txs, sql, limit, offset); err != nil {
		return
	}
	return
}

// EthTransactionsWithAttempts returns all eth transactions with at least one attempt
// limited by passed parameters. Attempts are sorted by id.
func (o *evmTxStorageService) EthTransactionsWithAttempts(offset, limit int) (txs []EthTx[*evmtypes.Address, *evmtypes.TxHash], count int, err error) {
	sql := `SELECT count(*) FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts)`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	if err = o.q.Select(&txs, sql, limit, offset); err != nil {
		return
	}

	err = o.preloadTxAttempts(txs)
	return
}

// EthTxAttempts returns the last tx attempts sorted by created_at descending.
func (o *evmTxStorageService) EthTxAttempts(offset, limit int) (txs []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], count int, err error) {
	sql := `SELECT count(*) FROM eth_tx_attempts`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_tx_attempts ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	if err = o.q.Select(&txs, sql, limit, offset); err != nil {
		return
	}
	err = o.PreloadEthTxes(txs)
	return
}

// FindEthTxAttempt returns an individual EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]
func (o *evmTxStorageService) FindEthTxAttempt(hash *evmtypes.TxHash) (*EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], error) {
	ethTxAttempt := EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]{}
	sql := `SELECT * FROM eth_tx_attempts WHERE hash = $1`
	if err := o.q.Get(&ethTxAttempt, sql, hash); err != nil {
		return nil, err
	}
	// reuse the preload
	attempts := []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]{ethTxAttempt}
	err := o.PreloadEthTxes(attempts)
	return &attempts[0], err
}

// FindEthTxAttemptsByEthTxIDs returns a list of attempts by ETH Tx IDs
func (o *evmTxStorageService) FindEthTxAttemptsByEthTxIDs(ids []int64) ([]EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], error) {
	var attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]

	sql := `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1)`
	if err := o.q.Select(&attempts, sql, ids); err != nil {
		return nil, err
	}

	return attempts, nil
}

func (o *evmTxStorageService) FindEthTxByHash(hash *evmtypes.TxHash) (*EthTx[*evmtypes.Address, *evmtypes.TxHash], error) {
	var etx EthTx[*evmtypes.Address, *evmtypes.TxHash]

	err := o.q.Transaction(func(tx pg.Queryer) error {
		sql := `SELECT eth_txes.* FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts WHERE hash = $1)`
		if err := tx.Get(&etx, sql, hash); err != nil {
			return errors.Wrapf(err, "failed to find eth_tx with hash %d", hash)
		}

		return nil
	}, pg.OptReadOnlyTx())

	return &etx, errors.Wrap(err, "FindEthTxByHash failed")
}

// InsertEthTxAttempt inserts a new txAttempt into the database
func (o *evmTxStorageService) InsertEthTx(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash]) error {
	if etx.CreatedAt == (time.Time{}) {
		etx.CreatedAt = time.Now()
	}
	const insertEthTxSQL = `INSERT INTO eth_txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, initial_broadcast_at, created_at, state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, access_list, transmit_checker) VALUES (
:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :error, :broadcast_at, :initial_broadcast_at, :created_at, :state, :meta, :subject, :pipeline_task_run_id, :min_confirmations, :evm_chain_id, :access_list, :transmit_checker
) RETURNING *`
	err := o.q.GetNamed(insertEthTxSQL, etx, etx)
	return errors.Wrap(err, "InsertEthTx failed")
}

func (o *evmTxStorageService) InsertEthTxAttempt(attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	const insertEthTxAttemptSQL = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap) VALUES (
:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap
) RETURNING *`
	err := o.q.GetNamed(insertEthTxAttemptSQL, attempt, attempt)
	return errors.Wrap(err, "InsertEthTxAttempt failed")
}

func (o *evmTxStorageService) InsertEthReceipt(receipt *EvmReceipt) error {
	// convert to database representation
	r := dbReceiptFromEvmReceipt(receipt)

	const insertEthReceiptSQL = `INSERT INTO eth_receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at) VALUES (
:tx_hash, :block_hash, :block_number, :transaction_index, :receipt, NOW()
) RETURNING *`
	err := o.q.GetNamed(insertEthReceiptSQL, &r, &r)

	// method expects original (destination) receipt struct to be updated
	*receipt = dbReceiptToEvmReceipt(&r)

	return errors.Wrap(err, "InsertEthReceipt failed")
}

// FindEthTxWithAttempts finds the EthTx[*evmtypes.Address, *evmtypes.TxHash] with its attempts and receipts preloaded
func (o *evmTxStorageService) FindEthTxWithAttempts(etxID int64) (etx EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Get(&etx, `SELECT * FROM eth_txes WHERE id = $1 ORDER BY created_at ASC, id ASC`, etxID); err != nil {
			return errors.Wrapf(err, "failed to find eth_tx with id %d", etxID)
		}
		if err = o.LoadEthTxAttempts(&etx, pg.WithQueryer(tx)); err != nil {
			return errors.Wrapf(err, "failed to load eth_tx_attempts for eth_tx with id %d", etxID)
		}
		if err = loadEthTxAttemptsReceipts(tx, &etx); err != nil {
			return errors.Wrapf(err, "failed to load eth_receipts for eth_tx with id %d", etxID)
		}
		return nil
	}, pg.OptReadOnlyTx())
	return etx, errors.Wrap(err, "FindEthTxWithAttempts failed")
}

func (o *evmTxStorageService) FindEthTxAttemptConfirmedByEthTxIDs(ids []int64) ([]EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], error) {
	var attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]
	err := o.q.Transaction(func(tx pg.Queryer) error {
		if err := tx.Select(&attempts, `SELECT eta.*
		FROM eth_tx_attempts eta
			join eth_receipts er on eta.hash = er.tx_hash where eta.eth_tx_id = ANY($1) ORDER BY eta.gas_price DESC, eta.gas_tip_cap DESC`, ids); err != nil {
			return err
		}
		return loadConfirmedAttemptsReceipts(tx, attempts)
	}, pg.OptReadOnlyTx())
	return attempts, errors.Wrap(err, "FindEthTxAttemptConfirmedByEthTxIDs failed")
}

func (o *evmTxStorageService) LoadEthTxesAttempts(etxs []*EthTx[*evmtypes.Address, *evmtypes.TxHash], qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	ethTxIDs := make([]int64, len(etxs))
	ethTxesM := make(map[int64]*EthTx[*evmtypes.Address, *evmtypes.TxHash], len(etxs))
	for i, etx := range etxs {
		etx.EthTxAttempts = nil // this will overwrite any previous preload
		ethTxIDs[i] = etx.ID
		ethTxesM[etx.ID] = etxs[i]
	}
	var ethTxAttempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]
	if err := qq.Select(&ethTxAttempts, `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1) ORDER BY eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`, pq.Array(ethTxIDs)); err != nil {
		return errors.Wrap(err, "loadEthTxesAttempts failed to load eth_tx_attempts")
	}
	for _, attempt := range ethTxAttempts {
		etx := ethTxesM[attempt.EthTxID]
		etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	}
	return nil
}

func (o *evmTxStorageService) LoadEthTxAttempts(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], qopts ...pg.QOpt) error {
	return o.LoadEthTxesAttempts([]*EthTx[*evmtypes.Address, *evmtypes.TxHash]{etx}, qopts...)
}

func loadEthTxAttemptsReceipts[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, etx *EthTx[*evmtypes.Address, *evmtypes.TxHash]) (err error) {
	return loadEthTxesAttemptsReceipts(q, []*EthTx[*evmtypes.Address, *evmtypes.TxHash]{etx})
}

func loadEthTxesAttemptsReceipts[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, etxs []*EthTx[*evmtypes.Address, *evmtypes.TxHash]) (err error) {
	if len(etxs) == 0 {
		return nil
	}
	attemptHashM := make(map[string]*EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], len(etxs)) // len here is lower bound
	attemptHashes := make([][]byte, len(etxs))                                                     // len here is lower bound
	for _, etx := range etxs {
		for i, attempt := range etx.EthTxAttempts {
			attemptHashM[attempt.Hash.String()] = &etx.EthTxAttempts[i]
			bytes, err := attempt.Hash.MarshalText()
			if err != nil {
				return errors.Wrap(err, "loadEthTxesAttemptsReceipts failed to marshal attempt hash")
			}
			attemptHashes = append(attemptHashes, bytes)
		}
	}
	var rs []dbReceipt
	if err = q.Select(&rs, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(attemptHashes)); err != nil {
		return errors.Wrap(err, "loadEthTxesAttemptsReceipts failed to load eth_receipts")
	}

	var receipts []EvmReceipt = fromDBReceipts(rs)

	for _, receipt := range receipts {
		attempt := attemptHashM[receipt.TxHash.String()]
		attempt.EthReceipts = append(attempt.EthReceipts, receipt)
	}
	return nil
}

func loadConfirmedAttemptsReceipts[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	byHash := make(map[string]*EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], len(attempts))
	hashes := make([][]byte, len(attempts))
	for i, attempt := range attempts {
		byHash[attempt.Hash.String()] = &attempts[i]
		bytes, err := attempt.Hash.MarshalText()
		if err != nil {
			return errors.Wrap(err, "loadConfirmedAttemptsReceipts failed to load eth_receipts")
		}
		hashes = append(hashes, bytes)
	}
	var rs []dbReceipt
	if err := q.Select(&rs, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(hashes)); err != nil {
		return errors.Wrap(err, "loadConfirmedAttemptsReceipts failed to load eth_receipts")
	}
	var receipts []EvmReceipt = fromDBReceipts(rs)
	for _, receipt := range receipts {
		attempt := byHash[receipt.TxHash.String()]
		attempt.EthReceipts = append(attempt.EthReceipts, receipt)
	}
	return nil
}

// FindEthTxAttemptsRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time (up to limit)
func (o *evmTxStorageService) FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID big.Int, address *evmtypes.Address) (attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], err error) {
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	// this select distinct works because of unique index on eth_txes
	// (evm_chain_id, from_address, nonce)
	err = o.q.Select(&attempts, `
SELECT DISTINCT ON (eth_txes.nonce) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= $1 AND evm_chain_id = $2 AND from_address = $3
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
LIMIT $4
`, olderThan, chainID.String(), address, limit)

	return attempts, errors.Wrap(err, "FindEthTxAttemptsRequiringResend failed to load eth_tx_attempts")
}

func (o *evmTxStorageService) UpdateBroadcastAts(now time.Time, etxIDs []int64) error {
	// Deliberately do nothing on NULL broadcast_at because that indicates the
	// tx has been moved into a state where broadcast_at is not relevant, e.g.
	// fatally errored.
	//
	// Since EthConfirmer/EthResender can race (totally OK since highest
	// priced transaction always wins) we only want to update broadcast_at if
	// our version is later.
	_, err := o.q.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = ANY($2) AND broadcast_at < $1`, now, pq.Array(etxIDs))
	return errors.Wrap(err, "updateBroadcastAts failed to update eth_txes")
}

// SetBroadcastBeforeBlockNum updates already broadcast attempts with the
// current block number. This is safe no matter how old the head is because if
// the attempt is already broadcast it _must_ have been before this head.
func (o *evmTxStorageService) SetBroadcastBeforeBlockNum(blockNum int64, chainID big.Int) error {
	_, err := o.q.Exec(
		`UPDATE eth_tx_attempts
SET broadcast_before_block_num = $1 
FROM eth_txes
WHERE eth_tx_attempts.broadcast_before_block_num IS NULL AND eth_tx_attempts.state = 'broadcast'
AND eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.evm_chain_id = $2`,
		blockNum, chainID.String(),
	)
	return errors.Wrap(err, "SetBroadcastBeforeBlockNum failed")
}

func (o *evmTxStorageService) FindEtxAttemptsConfirmedMissingReceipt(chainID big.Int) (attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], err error) {
	err = o.q.Select(&attempts,
		`SELECT DISTINCT ON (eth_tx_attempts.eth_tx_id) eth_tx_attempts.*
		FROM eth_tx_attempts
		JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state = 'confirmed_missing_receipt'
		WHERE evm_chain_id = $1
		ORDER BY eth_tx_attempts.eth_tx_id ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`,
		chainID.String())
	if err != nil {
		err = errors.Wrap(err, "FindEtxAttemptsConfirmedMissingReceipt failed to query")
	}
	return
}

func (o *evmTxStorageService) UpdateEthTxsUnconfirmed(ids []int64) error {
	_, err := o.q.Exec(`UPDATE eth_txes SET state='unconfirmed' WHERE id = ANY($1)`, pq.Array(ids))

	if err != nil {
		return errors.Wrap(err, "UpdateEthTxsUnconfirmed failed to execute")
	}
	return nil
}

func (o *evmTxStorageService) FindEthTxAttemptsRequiringReceiptFetch(chainID big.Int) (attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&attempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt') AND eth_txes.evm_chain_id = $1
WHERE eth_tx_attempts.state != 'insufficient_eth'
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
`, chainID.String())
		if err != nil {
			return errors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load eth_tx_attempts")
		}
		err = o.PreloadEthTxes(attempts)
		return errors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return
}

func (o *evmTxStorageService) SaveFetchedReceipts(r []*evmtypes.Receipt, chainID big.Int) (err error) {
	receipts := toDBRawReceipt(r)
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
			return errors.Wrap(err, "saveFetchedReceipts failed to marshal JSON")
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
	return errors.Wrap(err, "SaveFetchedReceipts failed to save receipts")
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
func (o *evmTxStorageService) MarkAllConfirmedMissingReceipt(chainID big.Int) (err error) {
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
		return errors.Wrap(err, "markAllConfirmedMissingReceipt failed")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "markAllConfirmedMissingReceipt RowsAffected failed")
	}
	if rowsAffected > 0 {
		o.logger.Infow(fmt.Sprintf("%d transactions missing receipt", rowsAffected), "n", rowsAffected)
	}
	return
}

func (o *evmTxStorageService) GetInProgressEthTxAttempts(ctx context.Context, address *evmtypes.Address, chainID big.Int) (attempts []EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], err error) {
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))
	err = qq.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&attempts, `
SELECT eth_tx_attempts.* FROM eth_tx_attempts
INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state in ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
WHERE eth_tx_attempts.state = 'in_progress' AND eth_txes.from_address = $1 AND eth_txes.evm_chain_id = $2
`, address, chainID.String())
		if err != nil {
			return errors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_tx_attempts")
		}
		err = o.PreloadEthTxes(attempts)
		return errors.Wrap(err, "getInProgressEthTxAttempts failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return attempts, errors.Wrap(err, "getInProgressEthTxAttempts failed")
}

func (o *evmTxStorageService) FindEthReceiptsPendingConfirmation(ctx context.Context, blockNum int64, chainID big.Int) (receiptsPlus []EvmReceiptPlus, err error) {
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

// FindEthTxWithNonce returns any broadcast ethtx with the given nonce
func (o *evmTxStorageService) FindEthTxWithNonce(fromAddress *evmtypes.Address, nonce int64) (etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	etx = new(EthTx[*evmtypes.Address, *evmtypes.TxHash])
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Get(etx, `
SELECT * FROM eth_txes WHERE from_address = $1 AND nonce = $2 AND state IN ('confirmed', 'confirmed_missing_receipt', 'unconfirmed')
`, fromAddress, nonce)
		if err != nil {
			return errors.Wrap(err, "FindEthTxWithNonce failed to load eth_txes")
		}
		err = o.LoadEthTxAttempts(etx, pg.WithQueryer(tx))
		return errors.Wrap(err, "FindEthTxWithNonce failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}

func updateEthTxAttemptUnbroadcast[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, attempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	if attempt.State != txmgrtypes.TxAttemptBroadcast {
		return errors.New("expected eth_tx_attempt to be broadcast")
	}
	_, err := q.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = NULL, state = 'in_progress' WHERE id = $1`, attempt.ID)
	return errors.Wrap(err, "updateEthTxAttemptUnbroadcast failed")
}

func updateEthTxUnconfirm[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, etx EthTx[*evmtypes.Address, *evmtypes.TxHash]) error {
	if etx.State != EthTxConfirmed {
		return errors.New("expected eth_tx state to be confirmed")
	}
	_, err := q.Exec(`UPDATE eth_txes SET state = 'unconfirmed' WHERE id = $1`, etx.ID)
	return errors.Wrap(err, "updateEthTxUnconfirm failed")
}

func deleteEthReceipts(q pg.Queryer, etxID int64) (err error) {
	_, err = q.Exec(`
DELETE FROM eth_receipts
USING eth_tx_attempts
WHERE eth_receipts.tx_hash = eth_tx_attempts.hash
AND eth_tx_attempts.eth_tx_id = $1
	`, etxID)
	return errors.Wrap(err, "deleteEthReceipts failed")
}

func (o *evmTxStorageService) UpdateEthTxForRebroadcast(etx EthTx[*evmtypes.Address, *evmtypes.TxHash], etxAttempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	return o.q.Transaction(func(tx pg.Queryer) error {
		if err := deleteEthReceipts(tx, etx.ID); err != nil {
			return errors.Wrapf(err, "deleteEthReceipts failed for etx %v", etx.ID)
		}
		if err := updateEthTxUnconfirm(tx, etx); err != nil {
			return errors.Wrapf(err, "updateEthTxUnconfirm failed for etx %v", etx.ID)
		}
		return updateEthTxAttemptUnbroadcast(tx, etxAttempt)
	})
}

func (o *evmTxStorageService) FindTransactionsConfirmedInBlockRange(highBlockNumber, lowBlockNumber int64, chainID big.Int) (etxs []*EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&etxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'broadcast'
INNER JOIN eth_receipts ON eth_receipts.tx_hash = eth_tx_attempts.hash
WHERE eth_txes.state IN ('confirmed', 'confirmed_missing_receipt') AND block_number BETWEEN $1 AND $2 AND evm_chain_id = $3
ORDER BY nonce ASC
`, lowBlockNumber, highBlockNumber, chainID.String())
		if err != nil {
			return errors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_txes")
		}
		if err = o.LoadEthTxesAttempts(etxs, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_tx_attempts")
		}
		err = loadEthTxesAttemptsReceipts(tx, etxs)
		return errors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed to load eth_receipts")
	}, pg.OptReadOnlyTx())
	return etxs, errors.Wrap(err, "FindTransactionsConfirmedInBlockRange failed")
}

func saveAttemptWithNewState[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, timeout time.Duration, logger logger.Logger, attempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], broadcastAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return pg.SqlxTransaction(ctx, q, logger, func(tx pg.Queryer) error {
		// In case of null broadcast_at (shouldn't happen) we don't want to
		// update anyway because it indicates a state where broadcast_at makes
		// no sense e.g. fatal_error
		if _, err := tx.Exec(`UPDATE eth_txes SET broadcast_at = $1 WHERE id = $2 AND broadcast_at < $1`, broadcastAt, attempt.EthTxID); err != nil {
			return errors.Wrap(err, "saveAttemptWithNewState failed to update eth_txes")
		}
		_, err := tx.Exec(`UPDATE eth_tx_attempts SET state=$1 WHERE id=$2`, attempt.State, attempt.ID)
		return errors.Wrap(err, "saveAttemptWithNewState failed to update eth_tx_attempts")
	})
}

func (o *evmTxStorageService) SaveInsufficientEthAttempt(timeout time.Duration, attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], broadcastAt time.Time) error {
	if !(attempt.State == txmgrtypes.TxAttemptInProgress || attempt.State == txmgrtypes.TxAttemptInsufficientEth) {
		return errors.New("expected state to be either in_progress or insufficient_eth")
	}
	attempt.State = txmgrtypes.TxAttemptInsufficientEth
	return errors.Wrap(saveAttemptWithNewState(o.q, timeout, o.logger, *attempt, broadcastAt), "saveInsufficientEthAttempt failed")
}

func saveSentAttempt[ADDR types.Hashable, TX_HASH types.Hashable](q pg.Queryer, timeout time.Duration, logger logger.Logger, attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], broadcastAt time.Time) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected state to be in_progress")
	}
	attempt.State = txmgrtypes.TxAttemptBroadcast
	return errors.Wrap(saveAttemptWithNewState(q, timeout, logger, *attempt, broadcastAt), "saveSentAttempt failed")
}

func (o *evmTxStorageService) SaveSentAttempt(timeout time.Duration, attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], broadcastAt time.Time) error {
	return saveSentAttempt(o.q, timeout, o.logger, attempt, broadcastAt)
}

func (o *evmTxStorageService) SaveConfirmedMissingReceiptAttempt(ctx context.Context, timeout time.Duration, attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], broadcastAt time.Time) error {
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))
	err := qq.Transaction(func(tx pg.Queryer) error {
		if err := saveSentAttempt(tx, timeout, o.logger, attempt, broadcastAt); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE eth_txes SET state = 'confirmed_missing_receipt' WHERE id = $1`, attempt.EthTxID); err != nil {
			return errors.Wrap(err, "failed to update eth_txes")
		}
		return nil
	})
	return errors.Wrap(err, "SaveConfirmedMissingReceiptAttempt failed")
}

func (o *evmTxStorageService) DeleteInProgressAttempt(ctx context.Context, attempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	qq := o.q.WithOpts(pg.WithParentCtx(ctx))

	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("DeleteInProgressAttempt: expected attempt state to be in_progress")
	}
	if attempt.ID == 0 {
		return errors.New("DeleteInProgressAttempt: expected attempt to have an id")
	}
	_, err := qq.Exec(`DELETE FROM eth_tx_attempts WHERE id = $1`, attempt.ID)
	return errors.Wrap(err, "DeleteInProgressAttempt failed")
}

// SaveInProgressAttempt inserts or updates an attempt
func (o *evmTxStorageService) SaveInProgressAttempt(attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash]) error {
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("SaveInProgressAttempt failed: attempt state must be in_progress")
	}
	// Insert is the usual mode because the attempt is new
	if attempt.ID == 0 {
		query, args, e := o.q.BindNamed(insertIntoEthTxAttemptsQuery, attempt)
		if e != nil {
			return errors.Wrap(e, "SaveInProgressAttempt failed to BindNamed")
		}
		return errors.Wrap(o.q.Get(attempt, query, args...), "SaveInProgressAttempt failed to insert into eth_tx_attempts")
	}
	// Update only applies to case of insufficient eth and simply changes the state to in_progress
	res, err := o.q.Exec(`UPDATE eth_tx_attempts SET state=$1, broadcast_before_block_num=$2 WHERE id=$3`, attempt.State, attempt.BroadcastBeforeBlockNum, attempt.ID)
	if err != nil {
		return errors.Wrap(err, "SaveInProgressAttempt failed to update eth_tx_attempts")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "SaveInProgressAttempt failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrapf(sql.ErrNoRows, "SaveInProgressAttempt tried to update eth_tx_attempts but no rows matched id %d", attempt.ID)
	}
	return nil
}

// FindEthTxsRequiringGasBump returns transactions that have all
// attempts which are unconfirmed for at least gasBumpThreshold blocks,
// limited by limit pending transactions
//
// It also returns eth_txes that are unconfirmed with no eth_tx_attempts
func (o *evmTxStorageService) FindEthTxsRequiringGasBump(ctx context.Context, address *evmtypes.Address, blockNum, gasBumpThreshold, depth int64, chainID big.Int) (etxs []*EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
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
		if err = tx.Select(&etxs, stmt, address, chainID.String(), depth, blockNum-gasBumpThreshold); err != nil {
			return errors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_txes")
		}
		err = o.LoadEthTxesAttempts(etxs, pg.WithQueryer(tx))
		return errors.Wrap(err, "FindEthTxsRequiringGasBump failed to load eth_tx_attempts")
	}, pg.OptReadOnlyTx())
	return
}

// FindEthTxsRequiringResubmissionDueToInsufficientEth returns transactions
// that need to be re-sent because they hit an out-of-eth error on a previous
// block
func (o *evmTxStorageService) FindEthTxsRequiringResubmissionDueToInsufficientEth(address *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) (etxs []*EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Transaction(func(tx pg.Queryer) error {
		err = tx.Select(&etxs, `
SELECT DISTINCT eth_txes.* FROM eth_txes
INNER JOIN eth_tx_attempts ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_tx_attempts.state = 'insufficient_eth'
WHERE eth_txes.from_address = $1 AND eth_txes.state = 'unconfirmed' AND eth_txes.evm_chain_id = $2
ORDER BY nonce ASC
`, address, chainID.String())
		if err != nil {
			return errors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_txes")
		}

		err = o.LoadEthTxesAttempts(etxs, pg.WithQueryer(tx))
		return errors.Wrap(err, "FindEthTxsRequiringResubmissionDueToInsufficientEth failed to load eth_tx_attempts")
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
func (o *evmTxStorageService) MarkOldTxesMissingReceiptAsErrored(blockNum int64, finalityDepth uint32, chainID big.Int, qopts ...pg.QOpt) error {
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
			return errors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to query")
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
			return errors.Wrap(err, "markOldTxesMissingReceiptAsErrored failed to load additional data")
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
				r.ID, blockNum, r.MaxBroadcastBeforeBlockNum, r.FromAddress.String(), nonce), "ethTxID", r.ID, "nonce", nonce, "fromAddress", r.FromAddress, "txHashes", txHashesHex)
		}

		return nil
	})
}

func (o *evmTxStorageService) SaveReplacementInProgressAttempt(oldAttempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], replacementAttempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	if oldAttempt.State != txmgrtypes.TxAttemptInProgress || replacementAttempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("expected attempts to be in_progress")
	}
	if oldAttempt.ID == 0 {
		return errors.New("expected oldAttempt to have an ID")
	}
	return qq.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE id=$1`, oldAttempt.ID); err != nil {
			return errors.Wrap(err, "saveReplacementInProgressAttempt failed to delete from eth_tx_attempts")
		}
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, replacementAttempt)
		if e != nil {
			return errors.Wrap(e, "saveReplacementInProgressAttempt failed to BindNamed")
		}
		return errors.Wrap(tx.Get(replacementAttempt, query, args...), "saveReplacementInProgressAttempt failed to insert replacement attempt")
	})
}

// Finds earliest saved transaction that has yet to be broadcast from the given address
func (o *evmTxStorageService) FindNextUnstartedTransactionFromAddress(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], fromAddress *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	err := qq.Get(etx, `SELECT * FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2 ORDER BY value ASC, created_at ASC, id ASC`, fromAddress, chainID.String())
	return errors.Wrap(err, "failed to FindNextUnstartedTransactionFromAddress")
}

func (o *evmTxStorageService) UpdateEthTxFatalError(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)

	if etx.State != EthTxInProgress {
		return errors.Errorf("can only transition to fatal_error from in_progress, transaction is currently %s", etx.State)
	}
	if !etx.Error.Valid {
		return errors.New("expected error field to be set")
	}

	etx.Nonce = nil
	etx.State = EthTxFatalError

	return qq.Transaction(func(tx pg.Queryer) error {
		if _, err := tx.Exec(`DELETE FROM eth_tx_attempts WHERE eth_tx_id = $1`, etx.ID); err != nil {
			return errors.Wrapf(err, "saveFatallyErroredTransaction failed to delete eth_tx_attempt with eth_tx.ID %v", etx.ID)
		}
		return errors.Wrap(tx.Get(etx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=NULL, initial_broadcast_at=NULL, nonce=NULL WHERE id=$3 RETURNING *`, etx.State, etx.Error, etx.ID), "saveFatallyErroredTransaction failed to save eth_tx")
	})
}

// Updates eth attempt from in_progress to broadcast. Also updates the eth tx to unconfirmed.
// Before it updates both tables though it increments the next nonce from the keystore
// One of the more complicated signatures. We have to accept variable pg.QOpt and QueryerFunc arguments
func (o *evmTxStorageService) UpdateEthTxAttemptInProgressToBroadcast(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], attempt EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], NewAttemptState txmgrtypes.TxAttemptState, incrNextNonceCallback txmgrtypes.QueryerFunc, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)

	if etx.BroadcastAt == nil {
		return errors.New("unconfirmed transaction must have broadcast_at time")
	}
	if etx.InitialBroadcastAt == nil {
		return errors.New("unconfirmed transaction must have initial_broadcast_at time")
	}
	if etx.State != EthTxInProgress {
		return errors.Errorf("can only transition to unconfirmed from in_progress, transaction is currently %s", etx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("attempt must be in in_progress state")
	}
	if NewAttemptState != txmgrtypes.TxAttemptBroadcast {
		return errors.Errorf("new attempt state must be broadcast, got: %s", NewAttemptState)
	}
	etx.State = EthTxUnconfirmed
	attempt.State = NewAttemptState
	return qq.Transaction(func(tx pg.Queryer) error {
		if err := incrNextNonceCallback(tx); err != nil {
			return errors.Wrap(err, "SaveEthTxAttempt failed on incrNextNonceCallback")
		}
		if err := tx.Get(etx, `UPDATE eth_txes SET state=$1, error=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id = $5 RETURNING *`, etx.State, etx.Error, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID); err != nil {
			return errors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx")
		}
		if err := tx.Get(&attempt, `UPDATE eth_tx_attempts SET state = $1 WHERE id = $2 RETURNING *`, attempt.State, attempt.ID); err != nil {
			return errors.Wrap(err, "SaveEthTxAttempt failed to save eth_tx_attempt")
		}
		return nil
	})
}

// Updates eth tx from unstarted to in_progress and inserts in_progress eth attempt
func (o *evmTxStorageService) UpdateEthTxUnstartedToInProgress(etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], attempt *EthTxAttempt[*evmtypes.Address, *evmtypes.TxHash], qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	if etx.Nonce == nil {
		return errors.New("in_progress transaction must have nonce")
	}
	if etx.State != EthTxUnstarted {
		return errors.Errorf("can only transition to in_progress from unstarted, transaction is currently %s", etx.State)
	}
	if attempt.State != txmgrtypes.TxAttemptInProgress {
		return errors.New("attempt state must be in_progress")
	}
	etx.State = EthTxInProgress
	return qq.Transaction(func(tx pg.Queryer) error {
		query, args, e := tx.BindNamed(insertIntoEthTxAttemptsQuery, attempt)
		if e != nil {
			return errors.Wrap(e, "failed to BindNamed")
		}
		err := tx.Get(attempt, query, args...)
		if err != nil {
			var pqErr *pgconn.PgError
			isPqErr := errors.As(err, &pqErr)
			if isPqErr && pqErr.ConstraintName == "eth_tx_attempts_eth_tx_id_fkey" {
				return errEthTxRemoved
			}
			return errors.Wrap(err, "UpdateEthTxUnstartedToInProgress failed to create eth_tx_attempt")
		}
		err = tx.Get(etx, `UPDATE eth_txes SET nonce=$1, state=$2, broadcast_at=$3, initial_broadcast_at=$4 WHERE id=$5 RETURNING *`, etx.Nonce, etx.State, etx.BroadcastAt, etx.InitialBroadcastAt, etx.ID)
		return errors.Wrap(err, "UpdateEthTxUnstartedToInProgress failed to update eth_tx")
	})
}

// GetEthTxInProgress returns either 0 or 1 transaction that was left in
// an unfinished state because something went screwy the last time. Most likely
// the node crashed in the middle of the ProcessUnstartedEthTxs loop.
// It may or may not have been broadcast to an eth node.
func (o *evmTxStorageService) GetEthTxInProgress(fromAddress *evmtypes.Address, qopts ...pg.QOpt) (etx *EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	qq := o.q.WithOpts(qopts...)
	etx = new(EthTx[*evmtypes.Address, *evmtypes.TxHash])
	bytes, err := fromAddress.MarshalText()
	if err != nil {
		return etx, errors.Wrap(err, "getInProgressEthTx failed")
	}
	err = qq.Transaction(func(tx pg.Queryer) error {
		err = qq.Get(etx, `SELECT * FROM eth_txes WHERE from_address = $1 and state = 'in_progress'`, bytes)
		if errors.Is(err, sql.ErrNoRows) {
			etx = nil
			return nil
		} else if err != nil {
			return errors.Wrap(err, "GetEthTxInProgress failed while loading eth tx")
		}
		if err = o.LoadEthTxAttempts(etx, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "GetEthTxInProgress failed while loading EthTxAttempts")
		}
		if len(etx.EthTxAttempts) != 1 || etx.EthTxAttempts[0].State != txmgrtypes.TxAttemptInProgress {
			return errors.Errorf("invariant violation: expected in_progress transaction %v to have exactly one unsent attempt. "+
				"Your database is in an inconsistent state and this node will not function correctly until the problem is resolved", etx.ID)
		}
		return nil
	})

	return etx, errors.Wrap(err, "getInProgressEthTx failed")
}

func (o *evmTxStorageService) HasInProgressTransaction(account *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) (exists bool, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Get(&exists, `SELECT EXISTS(SELECT 1 FROM eth_txes WHERE state = 'in_progress' AND from_address = $1 AND evm_chain_id = $2)`, account, chainID.String())
	return exists, errors.Wrap(err, "hasInProgressTransaction failed")
}

func (o *evmTxStorageService) UpdateEthKeyNextNonce(newNextNonce, currentNextNonce int64, address *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) error {
	qq := o.q.WithOpts(qopts...)
	return qq.Transaction(func(tx pg.Queryer) error {
		//  We filter by next_nonce here as an optimistic lock to make sure it
		//  didn't get changed out from under us. Shouldn't happen but can't hurt.
		res, err := tx.Exec(`UPDATE evm_key_states SET next_nonce = $1, updated_at = $2 WHERE address = $3 AND next_nonce = $4 AND evm_chain_id = $5`, newNextNonce, time.Now(), address, currentNextNonce, chainID.String())
		if err != nil {
			return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to update keys.next_nonce")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "NonceSyncer#fastForwardNonceIfNecessary failed to get RowsAffected")
		}
		if rowsAffected == 0 {
			return ErrKeyNotUpdated
		}
		return nil
	})
}

func (o *evmTxStorageService) countTransactionsWithState(fromAddress *evmtypes.Address, state EthTxState, chainID big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	qq := o.q.WithOpts(qopts...)
	err = qq.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = $2 AND evm_chain_id = $3`,
		fromAddress, state, chainID.String())
	return count, errors.Wrap(err, "failed to countTransactionsWithState")
}

// CountUnconfirmedTransactions returns the number of unconfirmed transactions
func (o *evmTxStorageService) CountUnconfirmedTransactions(fromAddress *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	return o.countTransactionsWithState(fromAddress, EthTxUnconfirmed, chainID, qopts...)
}

// CountUnstartedTransactions returns the number of unconfirmed transactions
func (o *evmTxStorageService) CountUnstartedTransactions(fromAddress *evmtypes.Address, chainID big.Int, qopts ...pg.QOpt) (count uint32, err error) {
	return o.countTransactionsWithState(fromAddress, EthTxUnstarted, chainID, qopts...)
}

func (o *evmTxStorageService) CheckEthTxQueueCapacity(fromAddress *evmtypes.Address, maxQueuedTransactions uint64, chainID big.Int, qopts ...pg.QOpt) (err error) {
	qq := o.q.WithOpts(qopts...)
	if maxQueuedTransactions == 0 {
		return nil
	}
	var count uint64
	err = qq.Get(&count, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state = 'unstarted' AND evm_chain_id = $2`, fromAddress, chainID.String())
	if err != nil {
		err = errors.Wrap(err, "CheckEthTxQueueCapacity query failed")
		return
	}

	if count >= maxQueuedTransactions {
		err = errors.Errorf("cannot create transaction; too many unstarted transactions in the queue (%v/%v). %s", count, maxQueuedTransactions, label.MaxQueuedTransactionsWarning)
	}
	return
}

func (o *evmTxStorageService) CreateEthTransaction(newTx NewTx[*evmtypes.Address], chainID big.Int, qopts ...pg.QOpt) (etx EthTx[*evmtypes.Address, *evmtypes.TxHash], err error) {
	qq := o.q.WithOpts(qopts...)
	value := 0
	err = qq.Transaction(func(tx pg.Queryer) error {
		if newTx.PipelineTaskRunID != nil {
			err = tx.Get(&etx, `SELECT * FROM eth_txes WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, newTx.PipelineTaskRunID, chainID.String())
			// If no eth_tx matches (the common case) then continue
			if !errors.Is(err, sql.ErrNoRows) {
				if err != nil {
					return errors.Wrap(err, "CreateEthTransaction")
				}
				// if a previous transaction for this task run exists, immediately return it
				return nil
			}
		}
		err := tx.Get(&etx, `
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, subject, evm_chain_id, min_confirmations, pipeline_task_run_id, transmit_checker)
VALUES (
$1,$2,$3,$4,$5,'unstarted',NOW(),$6,$7,$8,$9,$10,$11
)
RETURNING "eth_txes".*
`, newTx.FromAddress, newTx.ToAddress, newTx.EncodedPayload, value, newTx.GasLimit, newTx.Meta, newTx.Strategy.Subject(), chainID.String(), newTx.MinConfirmations, newTx.PipelineTaskRunID, newTx.Checker)
		if err != nil {
			return errors.Wrap(err, "CreateEthTransaction failed to insert eth_tx")
		}

		pruned, err := newTx.Strategy.PruneQueue(o, pg.WithQueryer(tx))
		if err != nil {
			return errors.Wrap(err, "CreateEthTransaction failed to prune eth_txes")
		}
		if pruned > 0 {
			o.logger.Warnw(fmt.Sprintf("Dropped %d old transactions from transaction queue", pruned), "fromAddress", newTx.FromAddress, "toAddress", newTx.ToAddress, "meta", newTx.Meta, "subject", newTx.Strategy.Subject(), "replacementID", etx.ID)
		}
		return nil
	})
	return
}

func (o *evmTxStorageService) PruneUnstartedTxQueue(queueSize uint32, subject uuid.UUID, qopts ...pg.QOpt) (n int64, err error) {
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
			return errors.Wrap(err, "DeleteUnstartedEthTx failed")
		}
		n, err = res.RowsAffected()
		return err
	})
	return
}
