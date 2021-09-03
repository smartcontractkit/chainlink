package log

import (
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore --structname ORM --filename orm.go

type ORM interface {
	FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error)
	WasBroadcastConsumed(tx *gorm.DB, blockHash common.Hash, logIndex uint, jobID int32) (bool, error)
	MarkBroadcastConsumed(tx *gorm.DB, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error
}

type orm struct {
	db         *gorm.DB
	evmChainID utils.Big
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, evmChainID big.Int) *orm {
	return &orm{db, *utils.NewBig(&evmChainID)}
}

func (o *orm) WasBroadcastConsumed(tx *gorm.DB, blockHash common.Hash, logIndex uint, jobID int32) (consumed bool, err error) {
	q := `
        SELECT consumed FROM log_broadcasts
        WHERE block_hash = ?
        AND log_index = ?
        AND job_id = ?
		AND evm_chain_id = ?
    `
	args := []interface{}{
		blockHash,
		logIndex,
		jobID,
		o.evmChainID,
	}

	err = tx.Raw(q, args...).Row().Scan(&consumed)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return consumed, err
}

func (o *orm) FindConsumedLogs(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error) {
	var broadcasts []LogBroadcast
	query := `
		SELECT block_hash, log_index, job_id FROM log_broadcasts
		WHERE block_number >= ?
		AND block_number <= ?
		AND evm_chain_id = ?
		AND consumed = true
	`
	err := o.db.Raw(query, fromBlockNum, toBlockNum, o.evmChainID).Find(&broadcasts).Error
	if err != nil {
		return make([]LogBroadcast, 0), err
	}
	return broadcasts, err
}

func (o *orm) MarkBroadcastConsumed(tx *gorm.DB, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error {
	query := tx.Exec(`
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, consumed, evm_chain_id) VALUES (?, ?, ?, ?, NOW(), true, ?)
    `, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	if query.Error != nil {
		return errors.Wrap(query.Error, "while marking log broadcast as consumed")
	} else if query.RowsAffected == 0 {
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
