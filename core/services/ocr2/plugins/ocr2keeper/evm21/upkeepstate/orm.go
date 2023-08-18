package upkeepstate

import (
	"math/big"
	"time"

	"github.com/lib/pq"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type orm struct {
	chainID *utils.Big
	q       pg.Q
}

type persistedStateRecord struct {
	UpkeepID            *utils.Big
	WorkID              string
	CompletionState     uint8
	BlockNumber         int64
	IneligibilityReason uint8
	InsertedAt          time.Time
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *orm {
	return &orm{
		chainID: utils.NewBig(chainID),
		q:       pg.NewQ(db, lggr.Named("ORM"), cfg),
	}
}

// InsertUpkeepState is idempotent and sets upkeep state values in db
func (o *orm) InsertUpkeepState(state persistedStateRecord, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)

	query := `INSERT INTO evm_upkeep_states (evm_chain_id, work_id, completion_state, block_number, inserted_at, upkeep_id, ineligibility_reason)
	  VALUES ($1::NUMERIC, $2, $3, $4, $5, $6::NUMERIC, $7)
	    ON CONFLICT (evm_chain_id, work_id)
	    DO NOTHING`

	return q.ExecQ(query, o.chainID, state.WorkID, state.CompletionState, state.BlockNumber, state.InsertedAt, state.UpkeepID, state.IneligibilityReason)
}

// SelectStatesByWorkIDs searches the data store for stored states for the
// provided work ids and configured chain id
func (o *orm) SelectStatesByWorkIDs(workIDs []string, qopts ...pg.QOpt) (states []persistedStateRecord, err error) {
	q := o.q.WithOpts(qopts...)

	err = q.Select(&states, `SELECT upkeep_id, work_id, completion_state, block_number, ineligibility_reason, inserted_at
	  FROM evm_upkeep_states
	  WHERE work_id = ANY($1) AND evm_chain_id = $2::NUMERIC`, pq.Array(workIDs), o.chainID)

	if err != nil {
		return nil, err
	}

	return states, err
}

// DeleteExpired prunes stored states older than to the provided time
func (o *orm) DeleteExpired(expired time.Time, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM evm_upkeep_states WHERE inserted_at <= $1 AND evm_chain_id::NUMERIC = $2`, expired, o.chainID)

	return err
}
