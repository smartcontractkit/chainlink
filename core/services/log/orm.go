package log

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore --structname ORM --filename orm.go

type ORM interface {
	UpsertLog(log types.Log) error
	UpsertBroadcastForListener(log types.Log, jobID interface{}) error
	UpsertBroadcastsForListenerSinceBlock(blockNumber uint64, address common.Address, jobID interface{}) error
	WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID interface{}) (bool, error)
	MarkBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID interface{}) error
	UnconsumedLogsPriorToBlock(blockNumber uint64) ([]types.Log, error)
	DeleteLogAndBroadcasts(blockHash common.Hash, logIndex uint) error
	DeleteUnconsumedBroadcastsForListener(jobID interface{}) error
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
	data := []byte{}
	data = append(data, log.Data...)
	err := o.db.Exec(`
INSERT INTO eth_logs (block_hash, block_number, index, address, topics, data, created_at) VALUES (?,?,?,?,?,?,NOW())
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
    `, log.BlockHash, log.BlockNumber, log.Index, log.Address, pq.ByteaArray(topics), data).Error
	return err
}

func (o *orm) UpsertBroadcastForListener(log types.Log, jobID interface{}) error {
	var jobIDName string
	switch v := jobID.(type) {
	case models.JobID:
		jobIDName = "job_id"
	case int32:
		jobIDName = "job_id_v2"
	default:
		panic(fmt.Sprintf("unrecognised type for jobID: %T", v))
	}

	q := `
INSERT INTO log_broadcasts (eth_log_id, block_hash, block_number, log_index, %[1]s, consumed, created_at)
SELECT eth_logs.id, ?, ?, ?, ?, false, NOW() FROM eth_logs
	WHERE eth_logs.block_hash = ? AND eth_logs.index = ?
ON CONFLICT (%[1]s, block_hash, log_index) WHERE %[1]s IS NOT NULL DO UPDATE SET (
	eth_log_id,
	block_hash,
	block_number,
	log_index,
	%[1]s
) = (
	EXCLUDED.eth_log_id,
	EXCLUDED.block_hash,
	EXCLUDED.block_number,
	EXCLUDED.log_index,
	EXCLUDED.%[1]s
)
`
	args := []interface{}{
		log.BlockHash,
		log.BlockNumber,
		log.Index,
		jobID,
		log.BlockHash,
		log.Index,
	}

	stmt := fmt.Sprintf(q, jobIDName)
	query := o.db.Exec(stmt, args...)

	if query.Error != nil {
		switch v := query.Error.(type) {
		case *pgconn.PgError:
			if v.Code == "23503" && v.ConstraintName == "log_broadcasts_eth_log_id_fkey" {
				logger.Debugw(
					"log.ORM.UpsertBroadcastForListener: tried to upsert log_broadcast but the referenced eth_log was already deleted",
					"blockHash", log.BlockHash, "blockNum", log.BlockNumber, "index", log.Index, "address", log.Address, "jobID", jobID,
					"log", log,
				)
				// NOTE: This rare case indicates that the eth_log was deleted
				// simultaneously with this query. Return as if we inserted
				// successfully and the delete operation came immediately
				// afterwards, cascading to the newly inserted log_broadcast.
				return nil
			}
		}
		return errors.Wrap(query.Error, "while upserting broadcast for listener")
	} else if query.RowsAffected == 0 {
		return errors.Errorf("no eth_log was found with block_hash %s and index %v", log.BlockHash, log.Index)
	}

	return nil
}

func (o *orm) UpsertBroadcastsForListenerSinceBlock(blockNumber uint64, address common.Address, jobID interface{}) (err error) {
	var jobIDName string
	switch v := jobID.(type) {
	case models.JobID:
		jobIDName = "job_id"
	case int32:
		jobIDName = "job_id_v2"
	default:
		panic(fmt.Sprintf("unrecognised type for jobID: %T", v))
	}

	q := `
INSERT INTO log_broadcasts (eth_log_id, block_hash, block_number, log_index, %[1]s, consumed, created_at)
SELECT id, block_hash, block_number, index, ?, false, NOW() FROM eth_logs
	WHERE eth_logs.block_number >= ? AND address = ?
ON CONFLICT (%[1]s, block_hash, log_index) WHERE %[1]s IS NOT NULL DO UPDATE SET (
	eth_log_id,
	block_hash,
	block_number,
	log_index,
	%[1]s
) = (
	EXCLUDED.eth_log_id,
	EXCLUDED.block_hash,
	EXCLUDED.block_number,
	EXCLUDED.log_index,
	EXCLUDED.%[1]s
)`

	args := []interface{}{
		jobID,
		blockNumber,
		address,
	}

	stmt := fmt.Sprintf(q, jobIDName)
	err = o.db.Exec(stmt, args...).Error
	switch v := err.(type) {
	case *pgconn.PgError:
		if v.Code == "23503" && v.ConstraintName == "log_broadcasts_eth_log_id_fkey" {
			logger.Debugw(
				"log.ORM.UpsertBroadcastsForListenerSinceBlock: tried to upsert log_broadcasts but the referenced eth_log was already deleted",
				"blockNum", blockNumber, "address", address, "jobID", jobID,
			)
			// NOTE: This rare case indicates that the eth_log was deleted
			// simultaneously with this query. Return as if we inserted
			// successfully and the delete operation came immediately
			// afterwards, cascading to the newly inserted log_broadcast.
			return nil
		}
	}
	return errors.Wrapf(err, "while upserting broadcast for listener since block %v", blockNumber)
}

func (o *orm) WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID interface{}) (consumed bool, err error) {
	var jobIDName string
	switch v := jobID.(type) {
	case models.JobID:
		jobIDName = "job_id"
	case int32:
		jobIDName = "job_id_v2"
	default:
		panic(fmt.Sprintf("unrecognised type for jobID: %T", v))
	}

	q := `
SELECT consumed FROM log_broadcasts
WHERE block_hash = ?
AND log_index = ?
AND %s = ?
`

	args := []interface{}{
		blockHash,
		logIndex,
		jobID,
	}

	stmt := fmt.Sprintf(q, jobIDName)
	err = o.db.Raw(stmt, args...).Row().Scan(&consumed)

	return consumed, err
}

func (o *orm) MarkBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID interface{}) error {
	var jobIDName string
	switch v := jobID.(type) {
	case models.JobID:
		jobIDName = "job_id"
	case int32:
		jobIDName = "job_id_v2"
	default:
		panic(fmt.Sprintf("unrecognised type for jobID: %T", v))
	}

	q := `
UPDATE log_broadcasts SET consumed = true
WHERE block_hash = ?
AND log_index = ?
AND %s = ?
`
	args := []interface{}{
		blockHash,
		logIndex,
		jobID,
	}

	stmt := fmt.Sprintf(q, jobIDName)
	query := o.db.Exec(stmt, args...)

	if query.Error != nil {
		return errors.Wrap(query.Error, "while marking log broadcast as consumed")
	} else if query.RowsAffected == 0 {
		return errors.Errorf("cannot mark log broadcast as consumed: does not exist")
	}
	return nil
}

func (o *orm) UnconsumedLogsPriorToBlock(blockNumber uint64) ([]types.Log, error) {
	logs, err := FetchLogs(o.db, `
        SELECT d.block_hash, d.block_number, d.index, d.address, d.topics, d.data FROM
		(
			SELECT DISTINCT ON (eth_logs.id) eth_logs.* FROM eth_logs
			INNER JOIN log_broadcasts ON eth_logs.id = log_broadcasts.eth_log_id
			WHERE eth_logs.block_number < $1 AND log_broadcasts.consumed = false
			ORDER BY eth_logs.id
		) d
        ORDER BY d.order_received, d.block_number, d.index ASC;
    `, blockNumber)
	if err != nil {
		logger.Errorw("could not fetch logs to broadcast", "error", err)
		return nil, err
	}
	return logs, nil
}

func (o *orm) DeleteLogAndBroadcasts(blockHash common.Hash, logIndex uint) error {
	return o.db.Exec(`DELETE FROM eth_logs WHERE block_hash = ? AND index = ?`, blockHash, logIndex).Error
}

func (o *orm) DeleteUnconsumedBroadcastsForListener(jobID interface{}) error {
	var jobIDName string
	switch v := jobID.(type) {
	case models.JobID:
		jobIDName = "job_id"
	case int32:
		jobIDName = "job_id_v2"
	default:
		panic(fmt.Sprintf("unrecognised type for jobID: %T", v))
	}

	q := `DELETE FROM log_broadcasts WHERE %s = ? AND consumed = false`

	stmt := fmt.Sprintf(q, jobIDName)
	return o.db.Exec(stmt, jobID).Error
}

type logRow struct {
	BlockHash   common.Hash
	BlockNumber uint64
	Index       uint
	Address     common.Address
	Topics      pq.ByteaArray
	Data        []byte
}

func FetchLogs(db *gorm.DB, query string, args ...interface{}) (logs []types.Log, err error) {
	d, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "FetchLogs failed")
	}
	rows, err := d.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "FetchLogs query failed")
	}
	defer logger.ErrorIfCalling(rows.Close)
	for rows.Next() {
		var lr logRow
		err := rows.Scan(&lr.BlockHash, &lr.BlockNumber, &lr.Index, &lr.Address, &lr.Topics, &lr.Data)
		if err != nil {
			return nil, errors.Wrap(err, "FetchLogs scan failed")
		}

		topics := make([]common.Hash, len(lr.Topics))
		bytesTopics := [][]byte(lr.Topics)
		for j, topic := range bytesTopics {
			topics[j] = common.BytesToHash(topic)
		}
		log := types.Log{
			Address:     lr.Address,
			Topics:      topics,
			Data:        lr.Data,
			BlockNumber: lr.BlockNumber,
			BlockHash:   lr.BlockHash,
			Index:       lr.Index,
		}
		logs = append(logs, log)
	}
	return logs, nil
}
