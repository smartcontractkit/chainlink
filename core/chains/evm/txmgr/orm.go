package txmgr

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	EthTransactions(offset, limit int) ([]EthTx, int, error)
	EthTransactionsWithAttempts(offset, limit int) ([]EthTx, int, error)
	EthTxAttempts(offset, limit int) ([]EthTxAttempt, int, error)
	FindEthTxAttempt(hash common.Hash) (*EthTxAttempt, error)
	FindEthTxAttemptsByEthTxIDs(ids []int64) ([]EthTxAttempt, error)
	FindEthTxByHash(hash common.Hash) (*EthTx, error)
	InsertEthTxAttempt(attempt *EthTxAttempt) error
	InsertEthTx(etx *EthTx) error
	InsertEthReceipt(receipt *EthReceipt) error
	FindEthTxWithAttempts(etxID int64) (etx EthTx, err error)
}

type orm struct {
	q      pg.Q
	logger logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) ORM {
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
