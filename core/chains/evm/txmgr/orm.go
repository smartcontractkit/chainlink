package txmgr

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	EthTransactions(offset, limit int) ([]EthTx, int, error)
	EthTransactionsWithAttempts(offset, limit int) ([]EthTx, int, error)
	EthTxAttempts(offset, limit int) ([]EthTxAttempt, int, error)
	FindEthTxAttempt(hash common.Hash) (*EthTxAttempt, error)
	FindEthTxAttemptConfirmedByEthTxIDs(ids []int64) ([]EthTxAttempt, error)
	FindEtxAttemptsConfirmedMissingReceipt(chainID big.Int) (attempts []EthTxAttempt, err error)
	FindEthTxAttemptsByEthTxIDs(ids []int64) ([]EthTxAttempt, error)
	FindEthTxAttemptsRequiringReceiptFetch(chainID big.Int) (attempts []EthTxAttempt, err error)
	FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID big.Int, address common.Address) (attempts []EthTxAttempt, err error)
	FindEthTxByHash(hash common.Hash) (*EthTx, error)
	FindEthTxWithAttempts(etxID int64) (etx EthTx, err error)
	// InsertEthReceipt only used in tests. Use SaveFetchedReceipts instead
	InsertEthReceipt(receipt *EthReceipt) error
	InsertEthTx(etx *EthTx) error
	InsertEthTxAttempt(attempt *EthTxAttempt) error
	MarkAllConfirmedMissingReceipt(chainID big.Int) (err error)
	SaveFetchedReceipts(receipts []evmtypes.Receipt, chainID big.Int) (err error)
	SetBroadcastBeforeBlockNum(blockNum int64, chainID big.Int) error
	UpdateBroadcastAts(now time.Time, etxIDs []int64) error
	UpdateEthTxsUnconfirmed(ids []int64) error
}

type orm struct {
	q      pg.Q
	logger logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ORM {
	namedLogger := lggr.Named("TxmORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &orm{q, namedLogger}
}

func (o *orm) preloadTxAttempts(txs []EthTx) error {
	// Preload TxAttempts
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ID)
	}
	if len(ids) == 0 {
		return nil
	}
	var attempts []EthTxAttempt
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

func (o *orm) preloadTxes(attempts []EthTxAttempt) error {
	var ids []int64
	for _, attempt := range attempts {
		ids = append(ids, attempt.EthTxID)
	}
	if len(ids) == 0 {
		return nil
	}
	var txs []EthTx
	sql := `SELECT * FROM eth_txes WHERE id IN (?)`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	query = o.q.Rebind(query)
	if err = o.q.Select(&txs, query, args...); err != nil {
		return err
	}
	// fill in txs
	for _, tx := range txs {
		for i, attempt := range attempts {
			if tx.ID == attempt.EthTxID {
				attempts[i].EthTx = tx
			}
		}
	}
	return nil
}

// EthTransactions returns all eth transactions without loaded relations
// limited by passed parameters.
func (o *orm) EthTransactions(offset, limit int) (txs []EthTx, count int, err error) {
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
func (o *orm) EthTransactionsWithAttempts(offset, limit int) (txs []EthTx, count int, err error) {
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
func (o *orm) EthTxAttempts(offset, limit int) (txs []EthTxAttempt, count int, err error) {
	sql := `SELECT count(*) FROM eth_tx_attempts`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_tx_attempts ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	if err = o.q.Select(&txs, sql, limit, offset); err != nil {
		return
	}
	err = o.preloadTxes(txs)
	return
}

// FindEthTxAttempt returns an individual EthTxAttempt
func (o *orm) FindEthTxAttempt(hash common.Hash) (*EthTxAttempt, error) {
	ethTxAttempt := EthTxAttempt{}
	sql := `SELECT * FROM eth_tx_attempts WHERE hash = $1`
	if err := o.q.Get(&ethTxAttempt, sql, hash); err != nil {
		return nil, err
	}
	// reuse the preload
	attempts := []EthTxAttempt{ethTxAttempt}
	err := o.preloadTxes(attempts)
	return &attempts[0], err
}

// FindEthTxAttemptsByEthTxIDs returns a list of attempts by ETH Tx IDs
func (o *orm) FindEthTxAttemptsByEthTxIDs(ids []int64) ([]EthTxAttempt, error) {
	var attempts []EthTxAttempt

	sql := `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = ANY($1)`
	if err := o.q.Select(&attempts, sql, ids); err != nil {
		return nil, err
	}

	return attempts, nil
}

func (o *orm) FindEthTxByHash(hash common.Hash) (*EthTx, error) {
	var etx EthTx

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
func (o *orm) InsertEthTx(etx *EthTx) error {
	if etx.CreatedAt == (time.Time{}) {
		etx.CreatedAt = time.Now()
	}
	const insertEthTxSQL = `INSERT INTO eth_txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, initial_broadcast_at, created_at, state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, access_list, transmit_checker) VALUES (
:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :error, :broadcast_at, :initial_broadcast_at, :created_at, :state, :meta, :subject, :pipeline_task_run_id, :min_confirmations, :evm_chain_id, :access_list, :transmit_checker
) RETURNING *`
	err := o.q.GetNamed(insertEthTxSQL, etx, etx)
	return errors.Wrap(err, "InsertEthTx failed")
}

func (o *orm) InsertEthTxAttempt(attempt *EthTxAttempt) error {
	const insertEthTxAttemptSQL = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap) VALUES (
:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap
) RETURNING *`
	err := o.q.GetNamed(insertEthTxAttemptSQL, attempt, attempt)
	return errors.Wrap(err, "InsertEthTxAttempt failed")
}

func (o *orm) InsertEthReceipt(receipt *EthReceipt) error {
	const insertEthReceiptSQL = `INSERT INTO eth_receipts (tx_hash, block_hash, block_number, transaction_index, receipt, created_at) VALUES (
:tx_hash, :block_hash, :block_number, :transaction_index, :receipt, NOW()
) RETURNING *`
	err := o.q.GetNamed(insertEthReceiptSQL, receipt, receipt)
	return errors.Wrap(err, "InsertEthReceipt failed")
}

// FindEthTxWithAttempts finds the EthTx with its attempts and receipts preloaded
func (o *orm) FindEthTxWithAttempts(etxID int64) (etx EthTx, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Get(&etx, `SELECT * FROM eth_txes WHERE id = $1 ORDER BY created_at ASC, id ASC`, etxID); err != nil {
			return errors.Wrapf(err, "failed to find eth_tx with id %d", etxID)
		}
		if err = loadEthTxAttempts(tx, &etx); err != nil {
			return errors.Wrapf(err, "failed to load eth_tx_attempts for eth_tx with id %d", etxID)
		}
		if err = loadEthTxAttemptsReceipts(tx, &etx); err != nil {
			return errors.Wrapf(err, "failed to load eth_receipts for eth_tx with id %d", etxID)
		}
		return nil
	}, pg.OptReadOnlyTx())
	return etx, errors.Wrap(err, "FindEthTxWithAttempts failed")
}

func (o *orm) FindEthTxAttemptConfirmedByEthTxIDs(ids []int64) ([]EthTxAttempt, error) {
	var attempts []EthTxAttempt
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

func loadEthTxAttempts(q pg.Queryer, etx *EthTx) error {
	err := q.Select(&etx.EthTxAttempts, `SELECT * FROM eth_tx_attempts WHERE eth_tx_id = $1 ORDER BY eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC`, etx.ID)
	return errors.Wrapf(err, "failed to load ethtxattempts for eth tx %d", etx.ID)
}

func loadEthTxAttemptsReceipts(q pg.Queryer, etx *EthTx) (err error) {
	return loadEthTxesAttemptsReceipts(q, []*EthTx{etx})
}

func loadEthTxesAttemptsReceipts(q pg.Queryer, etxs []*EthTx) (err error) {
	if len(etxs) == 0 {
		return nil
	}
	attemptHashM := make(map[common.Hash]*EthTxAttempt, len(etxs)) // len here is lower bound
	attemptHashes := make([][]byte, len(etxs))                     // len here is lower bound
	for _, etx := range etxs {
		for i, attempt := range etx.EthTxAttempts {
			attemptHashM[attempt.Hash] = &etx.EthTxAttempts[i]
			attemptHashes = append(attemptHashes, attempt.Hash.Bytes())
		}
	}
	var receipts []EthReceipt
	if err = q.Select(&receipts, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(attemptHashes)); err != nil {
		return errors.Wrap(err, "loadEthTxesAttemptsReceipts failed to load eth_receipts")
	}
	for _, receipt := range receipts {
		attempt := attemptHashM[receipt.TxHash]
		attempt.EthReceipts = append(attempt.EthReceipts, receipt)
	}
	return nil
}

func loadConfirmedAttemptsReceipts(q pg.Queryer, attempts []EthTxAttempt) error {
	byHash := make(map[common.Hash]*EthTxAttempt, len(attempts))
	hashes := make([][]byte, len(attempts))
	for i, attempt := range attempts {
		byHash[attempt.Hash] = &attempts[i]
		hashes = append(hashes, attempt.Hash.Bytes())
	}
	var receipts []EthReceipt
	if err := q.Select(&receipts, `SELECT * FROM eth_receipts WHERE tx_hash = ANY($1)`, pq.Array(hashes)); err != nil {
		return errors.Wrap(err, "loadConfirmedAttemptsReceipts failed to load eth_receipts")
	}
	for _, receipt := range receipts {
		attempt := byHash[receipt.TxHash]
		attempt.EthReceipts = append(attempt.EthReceipts, receipt)
	}
	return nil
}

// FindEthTxAttemptsRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time (up to limit)
func (o *orm) FindEthTxAttemptsRequiringResend(olderThan time.Time, maxInFlightTransactions uint32, chainID big.Int, address common.Address) (attempts []EthTxAttempt, err error) {
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	// this select distinct works because of unique index on eth_txes
	// (evm_chain_id, from_address, nonce)
	err = o.q.Select(&attempts, `
SELECT DISTINCT ON (nonce) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= $1 AND evm_chain_id = $2 AND from_address = $3
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
LIMIT $4
`, olderThan, chainID.String(), address, limit)

	return attempts, errors.Wrap(err, "FindEthTxAttemptsRequiringResend failed to load eth_tx_attempts")
}

func (o *orm) UpdateBroadcastAts(now time.Time, etxIDs []int64) error {
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
func (o *orm) SetBroadcastBeforeBlockNum(blockNum int64, chainID big.Int) error {
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

func (o *orm) FindEtxAttemptsConfirmedMissingReceipt(chainID big.Int) (attempts []EthTxAttempt, err error) {
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

func (o *orm) UpdateEthTxsUnconfirmed(ids []int64) error {
	_, err := o.q.Exec(`UPDATE eth_txes SET state='unconfirmed' WHERE id = ANY($1)`, pq.Array(ids))

	if err != nil {
		return errors.Wrap(err, "UpdateEthTxsUnconfirmed failed to execute")
	}
	return nil
}

func (o *orm) FindEthTxAttemptsRequiringReceiptFetch(chainID big.Int) (attempts []EthTxAttempt, err error) {
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
		err = loadEthTxes(tx, attempts)
		return errors.Wrap(err, "FindEthTxAttemptsRequiringReceiptFetch failed to load eth_txes")
	}, pg.OptReadOnlyTx())
	return
}

func (o *orm) SaveFetchedReceipts(receipts []evmtypes.Receipt, chainID big.Int) (err error) {
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
func (o *orm) MarkAllConfirmedMissingReceipt(chainID big.Int) (err error) {
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
