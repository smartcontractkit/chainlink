package terratxm

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages the data model for terra tx management.
type ORM struct {
	chainID string
	q       pg.Q
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID string, db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *ORM {
	namedLogger := lggr.Named("ORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		chainID: chainID,
		q:       q,
	}
}

// InsertMsg inserts a terra msg, assumed to be a serialized terra ExecuteContractMsg.
func (o *ORM) InsertMsg(contractID string, msg []byte) (int64, error) {
	var tm terra.Msg
	err := o.q.Get(&tm, `INSERT INTO terra_msgs (contract_id, raw, state, terra_chain_id, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING *`, contractID, msg, db.Unstarted, o.chainID)
	if err != nil {
		return 0, err
	}
	return tm.ID, nil
}

// SelectMsgsWithState selects the oldest messages with a given state up to limit.
func (o *ORM) SelectMsgsWithState(state db.State, limit int64) (terra.Msgs, error) {
	if limit < 1 {
		return terra.Msgs{}, errors.New("limit must be greater than 0")
	}
	var msgs terra.Msgs
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE state = $1 AND terra_chain_id = $2 ORDER BY created_at LIMIT $3`, state, o.chainID, limit); err != nil {
		return nil, err
	}
	return msgs, nil
}

// SelectMsgsWithIDs selects messages the given ids
func (o *ORM) SelectMsgsWithIDs(ids []int64) (terra.Msgs, error) {
	var msgs terra.Msgs
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE id = ANY($1)`, ids); err != nil {
		return nil, err
	}
	return msgs, nil
}

// UpdateMsgsWithState update the msgs with the given ids to the given state
// Note state transitions are validated at the db level.
func (o *ORM) UpdateMsgsWithState(ids []int64, state db.State, txHash *string, qopts ...pg.QOpt) error {
	if state == db.Broadcasted && txHash == nil {
		return errors.New("txHash is required when updating to broadcasted")
	}
	q := o.q.WithOpts(qopts...)
	var res sql.Result
	var err error
	if state == db.Broadcasted {
		res, err = q.Exec(`UPDATE terra_msgs SET state = $1, updated_at = NOW(), tx_hash = $2 WHERE id = ANY($3)`, state, *txHash, ids)
	} else {
		res, err = q.Exec(`UPDATE terra_msgs SET state = $1, updated_at = NOW() WHERE id = ANY($2)`, state, ids)
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
