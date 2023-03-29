package cosmostxm

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// ORM manages the data model for cosmos tx management.
type ORM struct {
	chainID string
	q       pg.Q
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID string, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *ORM {
	namedLogger := lggr.Named("Configs")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		chainID: chainID,
		q:       q,
	}
}

// InsertMsg inserts a cosmos msg, assumed to be a serialized cosmos ExecuteContractMsg.
func (o *ORM) InsertMsg(contractID, typeURL string, msg []byte, qopts ...pg.QOpt) (int64, error) {
	var tm adapters.Msg
	q := o.q.WithOpts(qopts...)
	err := q.Get(&tm, `INSERT INTO cosmos_msgs (contract_id, type, raw, state, cosmos_chain_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING *`, contractID, typeURL, msg, db.Unstarted, o.chainID)
	if err != nil {
		return 0, err
	}
	return tm.ID, nil
}

// UpdateMsgsContract updates messages for the given contract.
func (o *ORM) UpdateMsgsContract(contractID string, from, to db.State, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`UPDATE cosmos_msgs SET state = $1, updated_at = NOW()
	WHERE cosmos_chain_id = $2 AND contract_id = $3 AND state = $4`, to, o.chainID, contractID, from)
	if err != nil {
		return err
	}
	return nil
}

// GetMsgsState returns the oldest messages with a given state up to limit.
func (o *ORM) GetMsgsState(state db.State, limit int64, qopts ...pg.QOpt) (adapters.Msgs, error) {
	if limit < 1 {
		return adapters.Msgs{}, errors.New("limit must be greater than 0")
	}
	q := o.q.WithOpts(qopts...)
	var msgs adapters.Msgs
	if err := q.Select(&msgs, `SELECT * FROM cosmos_msgs WHERE state = $1 AND cosmos_chain_id = $2 ORDER BY id ASC LIMIT $3`, state, o.chainID, limit); err != nil {
		return nil, err
	}
	return msgs, nil
}

// GetMsgs returns any messages matching ids.
func (o *ORM) GetMsgs(ids ...int64) (adapters.Msgs, error) {
	var msgs adapters.Msgs
	if err := o.q.Select(&msgs, `SELECT * FROM cosmos_msgs WHERE id = ANY($1)`, ids); err != nil {
		return nil, err
	}
	return msgs, nil
}

// UpdateMsgs updates msgs with the given ids.
// Note state transitions are validated at the db level.
func (o *ORM) UpdateMsgs(ids []int64, state db.State, txHash *string, qopts ...pg.QOpt) error {
	if state == db.Broadcasted && txHash == nil {
		return errors.New("txHash is required when updating to broadcasted")
	}
	q := o.q.WithOpts(qopts...)
	var res sql.Result
	var err error
	if state == db.Broadcasted {
		res, err = q.Exec(`UPDATE cosmos_msgs SET state = $1, updated_at = NOW(), tx_hash = $2 WHERE id = ANY($3)`, state, *txHash, ids)
	} else {
		res, err = q.Exec(`UPDATE cosmos_msgs SET state = $1, updated_at = NOW() WHERE id = ANY($2)`, state, ids)
	}
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if int(count) != len(ids) {
		return errors.Errorf("expected %d records updated, got %d", len(ids), count)
	}
	return nil
}
