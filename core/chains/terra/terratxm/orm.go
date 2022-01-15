package terratxm

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

// ORM manages the data model for terra tx management.
type ORM struct {
	q pg.Q
}

// NewORM creates an ORM
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *ORM {
	namedLogger := lggr.Named("ORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		q: q,
	}
}

// InsertMsg inserts a terra msg, assumed to be a serialized terra ExecuteContractMsg.
func (o *ORM) InsertMsg(contractID string, msg []byte) (int64, error) {
	var tm TerraMsg
	err := o.q.Get(&tm, `INSERT INTO terra_msgs (contract_id, msg, state, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING *`, contractID, msg, Unstarted)
	if err != nil {
		return 0, err
	}
	return tm.ID, nil
}

// SelectMsgsWithState selects all messages with a given state
func (o *ORM) SelectMsgsWithState(state State) ([]TerraMsg, error) {
	var msgs []TerraMsg
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE state = $1`, state); err != nil {
		return nil, err
	}
	return msgs, nil
}

// SelectMsgsWithIDs selects messages the given ids
func (o *ORM) SelectMsgsWithIDs(ids []int64) ([]TerraMsg, error) {
	var msgs []TerraMsg
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE id = ANY($1)`, ids); err != nil {
		return nil, err
	}
	return msgs, nil
}

// UpdateMsgsWithState update the msgs with the given ids to the given state
// TODO: could enforce state transitions here too
func (o *ORM) UpdateMsgsWithState(ids []int64, state State, txHash *string, qopts ...pg.QOpt) error {
	if state == Broadcasted && txHash == nil {
		return errors.New("txHash is required when updating to broadcasted")
	}
	q := o.q.WithOpts(qopts...)
	var res sql.Result
	var err error
	if state == Broadcasted {
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
