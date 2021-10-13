package log

import (
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore --structname ORM --filename orm.go

type ORM interface {
	FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error)
	WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID int32, qopts ...postgres.QOpt) (bool, error)
	MarkBroadcastConsumed(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...postgres.QOpt) error
}

type orm struct {
	db         *sqlx.DB
	evmChainID utils.Big
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, evmChainID big.Int) *orm {
	return &orm{db, *utils.NewBig(&evmChainID)}
}

func (o *orm) WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID int32, qopts ...postgres.QOpt) (consumed bool, err error) {
	query := `
		SELECT consumed FROM log_broadcasts
		WHERE block_hash = $1
		AND log_index = $2
		AND job_id = $3
		AND evm_chain_id = $4
    `
	args := []interface{}{
		blockHash,
		logIndex,
		jobID,
		o.evmChainID,
	}
	q := postgres.NewQ(o.db, qopts...)
	err = q.QueryRowx(query, args...).Scan(&consumed)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return consumed, err
}

func (o *orm) FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error) {
	var broadcasts []LogBroadcast
	query := `
		SELECT block_hash, log_index, job_id FROM log_broadcasts
		WHERE block_number >= $1
		AND block_number <= $2
		AND evm_chain_id = $3
		AND consumed = true
	`
	err := o.db.Select(&broadcasts, query, fromBlockNum, toBlockNum, o.evmChainID)
	if err != nil {
		return make([]LogBroadcast, 0), err
	}
	return broadcasts, err
}

func (o *orm) MarkBroadcastConsumed(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...postgres.QOpt) error {
	q := postgres.NewQ(o.db, qopts...)
	res, err := q.Exec(`INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, consumed, evm_chain_id) VALUES ($1, $2, $3, $4, NOW(), true, $5)`, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	if err != nil {
		return errors.Wrap(err, "while marking log broadcast as consumed")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "MarkBroadcastConsumed failed")
	}
	if rowsAffected == 0 {
		return errors.Errorf("cannot mark log broadcast as consumed: does not exist")
	}
	return nil
}

// LogBroadcast - gorm-compatible receive data from log_broadcasts table columns
type LogBroadcast struct {
	BlockHash  common.Hash
	LogIndex   uint
	JobID      int32
	EVMChainID utils.Big
}

func (b LogBroadcast) AsKey() LogBroadcastAsKey {
	return LogBroadcastAsKey{
		b.BlockHash,
		b.LogIndex,
		b.JobID,
	}
}

// LogBroadcastAsKey - used as key in a map to filter out already consumed logs
type LogBroadcastAsKey struct {
	BlockHash common.Hash
	LogIndex  uint
	JobId     int32
}

func NewLogBroadcastAsKey(log types.Log, listener Listener) LogBroadcastAsKey {
	return LogBroadcastAsKey{
		log.BlockHash,
		log.Index,
		listener.JobID(),
	}
}
