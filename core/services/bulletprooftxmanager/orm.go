package bulletprooftxmanager

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
)

type ORM interface {
	EthTransactionsWithAttempts(offset, limit int) ([]EthTx, int, error)
	EthTxAttempts(offset, limit int) ([]EthTxAttempt, int, error)
	FindEthTxAttempt(hash common.Hash) (*EthTxAttempt, error)
	InsertEthTxAttempt(txAttempt *EthTxAttempt) error
}

type orm struct {
	db *sqlx.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB) ORM {
	return &orm{db}
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
	query = o.db.Rebind(query)
	if err = o.db.Select(&attempts, query, args...); err != nil {
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
	query = o.db.Rebind(query)
	if err = o.db.Select(&txs, query, args...); err != nil {
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

// EthTransactionsWithAttempts returns all eth transactions with at least one attempt
// limited by passed parameters. Attempts are sorted by id.
func (o *orm) EthTransactionsWithAttempts(offset, limit int) (txs []EthTx, count int, err error) {
	sql := `SELECT count(*) FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts)`
	if err = o.db.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_txes WHERE id IN (SELECT DISTINCT eth_tx_id FROM eth_tx_attempts) ORDER BY id desc LIMIT $1 OFFSET $2`
	if err = o.db.Select(&txs, sql, limit, offset); err != nil {
		return
	}

	err = o.preloadTxAttempts(txs)
	return
}

// EthTxAttempts returns the last tx attempts sorted by created_at descending.
func (o *orm) EthTxAttempts(offset, limit int) (txs []EthTxAttempt, count int, err error) {
	sql := `SELECT count(*) FROM eth_tx_attempts`
	if err = o.db.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM eth_tx_attempts ORDER BY created_at desc LIMIT $1 OFFSET $2`
	if err = o.db.Select(&txs, sql, limit, offset); err != nil {
		return
	}
	err = o.preloadTxes(txs)
	return
}

// FindEthTxAttempt returns an individual EthTxAttempt
func (o *orm) FindEthTxAttempt(hash common.Hash) (*EthTxAttempt, error) {
	ethTxAttempt := EthTxAttempt{}
	sql := `SELECT * FROM eth_tx_attempts WHERE hash = $1`
	if err := o.db.Get(&ethTxAttempt, sql, hash); err != nil {
		return nil, err
	}
	// reuse the preload
	attempts := []EthTxAttempt{ethTxAttempt}
	err := o.preloadTxes(attempts)
	return &attempts[0], err
}

// InsertEthTxAttempt inserts a new txAttempt into the database
func (o *orm) InsertEthTxAttempt(txAttempt *EthTxAttempt) error {
	query, args, err := o.db.BindNamed(insertIntoEthTxAttemptsQuery, txAttempt)
	if err != nil {
		return errors.Wrap(err, "InsertEthTxAttempt failed to BindNamed")
	}
	return errors.Wrap(o.db.Get(txAttempt, query, args...), "InsertEthTxAttempt failed to insert")
}
