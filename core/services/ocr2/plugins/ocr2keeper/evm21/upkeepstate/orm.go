package upkeepstate

import (
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type ORM struct {
	chainID *big.Int
	q       pg.Q
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *ORM {
	return &ORM{
		chainID: chainID,
		q:       pg.NewQ(db, lggr.Named("ORM"), cfg),
	}
}

// InsertUpkeepState is idempotent and sets upkeep state values in db
func (o *ORM) InsertUpkeepState(state upkeepStateRecord, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)

	query := `INSERT INTO evm_upkeep_state (evm_chain_id, work_id, completion_state, block_number, added_at)
	  VALUES ($1::NUMERIC, $2, $3, $4::NUMERIC, $5)
	    ON CONFLICT (evm_chain_id, work_id)
	    DO UPDATE SET completion_state=$3, block_number=$4`

	return q.ExecQ(query, utils.NewBig(o.chainID), state.WorkID, state.CompletionState, state.BlockNumber, state.AddedAt)
}

// SelectStatesByWorkIDs searches the data store for stored states for the
// provided work ids and configured chain id
func (o *ORM) SelectStatesByWorkIDs(workIDs []string, qopts ...pg.QOpt) (states []upkeepStateRecord, err error) {
	q := o.q.WithOpts(qopts...)

	namedArgs := map[string]any{
		"chainID": utils.NewBig(o.chainID),
		"workIDs": workIDs,
	}

	query, args, err := sqlx.Named(`SELECT work_id, completion_state, block_number, added_at
	  FROM evm_upkeep_state
	  WHERE evm_chain_id = :chainID AND work_id IN (:workIDs)`, namedArgs)

	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	query = q.Rebind(query)

	err = q.Select(&states, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return states, err
}

// DeleteBeforeTime prunes stored states older than to the provided time
func (o *ORM) DeleteBeforeTime(end time.Time, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM evm_upkeep_state WHERE added_at BEFORE $1 AND evm_chain_id::NUMERIC = $2`, end, utils.NewBig(o.chainID))

	return err
}
