package log

import (
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore --structname ORM --filename orm.go

type ORM interface {
	FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error)
	WasBroadcastConsumed(q sqlx.QueryerContext, blockHash common.Hash, logIndex uint, jobID JobIdSelect) (bool, error)
	MarkBroadcastConsumed(e sqlx.ExecerContext, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID JobIdSelect) error
}

type orm struct {
	db *sqlx.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB) *orm {
	return &orm{db}
}

func (o *orm) WasBroadcastConsumed(q sqlx.QueryerContext, blockHash common.Hash, logIndex uint, jobID JobIdSelect) (consumed bool, err error) {
	var jobIDValue interface{}
	var jobIDName = "job_id"
	if jobID.IsV2 {
		jobIDName = "job_id_v2"
		jobIDValue = jobID.JobIDV2
	} else {
		jobIDValue = jobID.JobIDV1
	}
	query := `
        SELECT consumed FROM log_broadcasts
        WHERE block_hash = $1
        AND log_index = $2
        AND %s = $3
    `
	args := []interface{}{
		blockHash,
		logIndex,
		jobIDValue,
	}

	stmt := fmt.Sprintf(query, jobIDName)
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err = sqlx.GetContext(ctx, q, &consumed, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return consumed, err
}

func (o *orm) FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error) {
	var broadcasts []LogBroadcast
	query := `
		SELECT block_hash, log_index, job_id, job_id_v2 FROM log_broadcasts
		WHERE block_number >= $1
		AND block_number <= $2
		AND consumed = true
	`
	err := o.db.Select(&broadcasts, query, fromBlockNum, toBlockNum)
	if err != nil {
		return make([]LogBroadcast, 0), err
	}
	return broadcasts, err
}

func (o *orm) MarkBroadcastConsumed(e sqlx.ExecerContext, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID JobIdSelect) error {
	var jobID1Value interface{}
	var jobID2Value interface{}

	if jobID.IsV2 {
		jobID2Value = jobID.JobIDV2
	} else {
		jobID1Value = jobID.JobIDV1
	}

	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	query, err := e.ExecContext(ctx, `
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, job_id_v2, created_at, consumed) VALUES ($1, $2, $3, $4, $5, NOW(), true)
    `, blockHash, blockNumber, logIndex, jobID1Value, jobID2Value)
	if err != nil {
		return errors.Wrap(err, "while marking log broadcast as consumed")
	}
	if rows, err := query.RowsAffected(); rows == 0 || err != nil {
		return errors.Errorf("cannot mark log broadcast as consumed: does not exist")
	}
	return nil
}

// LogBroadcast - gorm-compatible receive data from log_broadcasts table columns
type LogBroadcast struct {
	BlockHash common.Hash
	LogIndex  uint

	JobId   uuid.NullUUID
	JobIdV2 null.Int64
}

func (b LogBroadcast) JobID() JobIdSelect {
	if b.JobIdV2.Valid {
		return NewJobIdV2(int32(b.JobIdV2.Int64))
	}
	return NewJobIdV1(models.JobID(b.JobId.UUID))
}

func (b LogBroadcast) AsKey() LogBroadcastAsKey {
	return LogBroadcastAsKey{
		b.BlockHash,
		b.LogIndex,
		b.JobID().String(),
	}
}

// LogBroadcastAsKey - used as key in a map to filter out already consumed logs
type LogBroadcastAsKey struct {
	BlockHash common.Hash
	LogIndex  uint
	JobId     string
}

func NewLogBroadcastAsKey(log types.Log, listener Listener) LogBroadcastAsKey {
	return LogBroadcastAsKey{
		log.BlockHash,
		log.Index,
		NewJobIdFromListener(listener).String(),
	}
}
