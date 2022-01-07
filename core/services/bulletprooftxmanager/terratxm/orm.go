package terratxm

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

type ORM struct {
	q pg.Q
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *ORM {
	namedLogger := lggr.Named("TerraTxmORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		q: q,
	}
}

func (o *ORM) InsertMsg(contractID string, msg []byte) (int64, error) {
	var tm TerraMsg
	err := o.q.Get(&tm, `INSERT INTO terra_msgs (contract_id, msg, state, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING *`, contractID, msg, Unstarted)
	if err != nil {
		return 0, err
	}
	return tm.ID, nil
}

func (o *ORM) SelectMsgsWithState(state State) ([]TerraMsg, error) {
	var msgs []TerraMsg
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE state = $1`, state); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (o *ORM) SelectMsgsWithIDs(ids []int64) ([]TerraMsg, error) {
	var msgs []TerraMsg
	if err := o.q.Select(&msgs, `SELECT * FROM terra_msgs WHERE id = ANY($1)`, ids); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (o *ORM) UpdateMsgsWithState(ids []int64, state State) error {
	res, err := o.q.Exec(`UPDATE terra_msgs SET state = $1, updated_at = NOW() WHERE id = ANY($2)`, state, ids)
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
