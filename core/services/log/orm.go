package log

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type ORM interface {
	UpsertLog(log types.Log) error
	UpsertUnconsumedLogBroadcastForListener(log types.Log, listener Listener) error
	WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID *models.ID, jobIDV2 int32) (bool, error)
	MarkBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID *models.ID, jobIDV2 int32) error
	UnconsumedLogsPriorToBlock(blockNumber uint64) ([]types.Log, error)
	DeleteLogAndBroadcasts(blockHash common.Hash, blockNumber uint64, logIndex uint) error
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

func (o *orm) UpsertLog(log types.Log) error {
	topics := make([][]byte, len(log.Topics))
	for i, topic := range log.Topics {
		x := make([]byte, len(topic))
		copy(x, topic[:])
		topics[i] = x
	}
	err := o.db.Exec(`
        INSERT INTO logs (block_hash, block_number, index, address, topics, data, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW())
        ON CONFLICT (block_hash, index) DO UPDATE SET (
            block_hash,
            block_number,
            index,
            address,
            topics,
            data
        ) = (
            EXCLUDED.block_hash,
            EXCLUDED.block_number,
            EXCLUDED.index,
            EXCLUDED.address,
            EXCLUDED.topics,
            EXCLUDED.data
        )
    `, log.BlockHash, log.BlockNumber, log.Index, log.Address, pq.ByteaArray(topics), log.Data).Error
	return err
}

func (o *orm) UpsertUnconsumedLogBroadcastForListener(log types.Log, listener Listener) error {
	if listener.IsV2Job() {
		return o.db.Exec(`
            INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id_v2, consumed, created_at)
            VALUES (?, ?, ?, ?, false, NOW())
            ON CONFLICT (job_id_v2, block_hash, log_index) DO UPDATE SET (
                block_hash,
                block_number,
                log_index,
                job_id_v2,
                consumed
            ) = (
                EXCLUDED.block_hash,
                EXCLUDED.block_number,
                EXCLUDED.log_index,
                EXCLUDED.job_id_v2,
                false
            )
        `, log.BlockHash, log.BlockNumber, log.Index, listener.JobIDV2()).Error
	} else {
		return o.db.Exec(`
            INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, consumed, created_at)
            VALUES (?, ?, ?, ?, false, NOW())
            ON CONFLICT (job_id, block_hash, log_index) DO UPDATE SET (
                block_hash,
                block_number,
                log_index,
                job_id,
                consumed
            ) = (
                EXCLUDED.block_hash,
                EXCLUDED.block_number,
                EXCLUDED.log_index,
                EXCLUDED.job_id,
                false
            )
        `, log.BlockHash, log.BlockNumber, log.Index, listener.JobID()).Error
	}
}

func (o *orm) WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID *models.ID, jobIDV2 int32) (bool, error) {
	var consumed struct{ Consumed bool }
	var err error
	if jobID == nil {
		err = o.db.Raw(`
            SELECT consumed FROM log_broadcasts
            WHERE block_hash = ?
            AND log_index = ?
            AND job_id IS NULL
            AND job_id_v2 = ?
        `, blockHash, logIndex, jobIDV2).Scan(&consumed).Error
	} else {
		err = o.db.Raw(`
            SELECT consumed FROM log_broadcasts
            WHERE block_hash = ?
            AND log_index = ?
            AND job_id = ?
            AND job_id_v2 IS NULL
        `, blockHash, logIndex, jobID).Scan(&consumed).Error
	}

	return consumed.Consumed, err
}

func (o *orm) MarkBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID *models.ID, jobIDV2 int32) error {
	var query *gorm.DB
	if jobID == nil {
		query = o.db.Exec(`
            UPDATE log_broadcasts SET consumed = true
            WHERE block_hash = ?
            AND log_index = ?
            AND job_id IS NULL
            AND job_id_v2 = ?
        `, blockHash, logIndex, jobIDV2)
	} else {
		query = o.db.Exec(`
            UPDATE log_broadcasts SET consumed = true
            WHERE block_hash = ?
            AND log_index = ?
            AND job_id = ?
            AND job_id_v2 IS NULL
        `, blockHash, logIndex, jobID)
	}
	if query.Error != nil {
		return query.Error
	} else if query.RowsAffected == 0 {
		return errors.Errorf("cannot mark log broadcast as consumed: does not exist")
	}
	return nil
}

func (o *orm) UnconsumedLogsPriorToBlock(blockNumber uint64) ([]types.Log, error) {
	type logRow struct {
		Address     common.Address
		Topics      pq.ByteaArray
		Data        []byte
		BlockNumber uint64
		BlockHash   common.Hash
		Index       uint
		Consumed    bool
	}

	var logRows []logRow
	err := o.db.Raw(`
        SELECT logs.*, bool_and(log_broadcasts.consumed) as consumed FROM logs
        LEFT JOIN log_broadcasts ON logs.block_hash = log_broadcasts.block_hash AND logs.index = log_broadcasts.log_index
        WHERE logs.block_number < ?
        GROUP BY logs.block_hash, logs.index, log_broadcasts.consumed
        HAVING consumed = false
        ORDER BY logs.block_number, logs.index ASC
    `, blockNumber).
		Scan(&logRows).Error
	if err != nil {
		logger.Errorw("could not fetch logs to broadcast", "error", err)
		return nil, err
	}
	logs := make([]types.Log, len(logRows))
	for i, log := range logRows {
		topics := make([]common.Hash, len(log.Topics))
		bytesTopics := [][]byte(log.Topics)
		for j, topic := range bytesTopics {
			topics[j] = common.BytesToHash(topic)
		}
		logs[i] = types.Log{
			Address:     log.Address,
			Topics:      topics,
			Data:        log.Data,
			BlockNumber: log.BlockNumber,
			BlockHash:   log.BlockHash,
			Index:       log.Index,
		}
	}
	return logs, nil
}

func (o *orm) DeleteLogAndBroadcasts(blockHash common.Hash, blockNumber uint64, logIndex uint) error {
	return o.db.Exec(`
        DELETE FROM logs WHERE block_hash = ? AND block_number = ? AND index = ?
    `, blockHash, blockNumber, logIndex).Error
}
