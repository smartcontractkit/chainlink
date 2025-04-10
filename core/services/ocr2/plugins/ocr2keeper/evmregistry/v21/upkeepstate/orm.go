package upkeepstate

import (
	"context"
	"math/big"
	"time"

	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
)

type orm struct {
	chainID *ubig.Big
	ds      sqlutil.DataSource
}

type persistedStateRecord struct {
	UpkeepID            *ubig.Big
	WorkID              string
	CompletionState     uint8
	BlockNumber         int64
	IneligibilityReason uint8
	InsertedAt          time.Time
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID *big.Int, ds sqlutil.DataSource) *orm {
	return &orm{
		chainID: ubig.New(chainID),
		ds:      ds,
	}
}

// BatchInsertRecords is idempotent and sets upkeep state values in db
func (o *orm) BatchInsertRecords(ctx context.Context, state []persistedStateRecord) error {
	if len(state) == 0 {
		return nil
	}

	type row struct {
		EvmChainId          *ubig.Big
		WorkId              string
		CompletionState     uint8
		BlockNumber         int64
		InsertedAt          time.Time
		UpkeepId            *ubig.Big
		IneligibilityReason uint8
	}

	var rows []row
	for _, record := range state {
		rows = append(rows, row{
			EvmChainId:          o.chainID,
			WorkId:              record.WorkID,
			CompletionState:     record.CompletionState,
			BlockNumber:         record.BlockNumber,
			InsertedAt:          record.InsertedAt,
			UpkeepId:            record.UpkeepID,
			IneligibilityReason: record.IneligibilityReason,
		})
	}

	_, err := o.ds.NamedExecContext(ctx, `INSERT INTO evm.upkeep_states
(evm_chain_id, work_id, completion_state, block_number, inserted_at, upkeep_id, ineligibility_reason) VALUES
(:evm_chain_id, :work_id, :completion_state, :block_number, :inserted_at, :upkeep_id, :ineligibility_reason) ON CONFLICT (evm_chain_id, work_id) DO NOTHING`, rows)
	return err
}

// SelectStatesByWorkIDs searches the data store for stored states for the
// provided work ids and configured chain id
func (o *orm) SelectStatesByWorkIDs(ctx context.Context, workIDs []string) (states []persistedStateRecord, err error) {
	err = o.ds.SelectContext(ctx, &states, `SELECT upkeep_id, work_id, completion_state, block_number, ineligibility_reason, inserted_at
	  FROM evm.upkeep_states
	  WHERE work_id = ANY($1) AND evm_chain_id = $2::NUMERIC`, pq.Array(workIDs), o.chainID)

	if err != nil {
		return nil, err
	}

	return states, err
}

// DeleteExpired prunes stored states older than to the provided time
func (o *orm) DeleteExpired(ctx context.Context, expired time.Time) error {
	_, err := o.ds.ExecContext(ctx, `DELETE FROM evm.upkeep_states WHERE inserted_at <= $1 AND evm_chain_id::NUMERIC = $2`, expired, o.chainID)

	return err
}
