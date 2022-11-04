package log

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ORM is the interface for log broadcasts.
//   - Unconsumed broadcasts are created just before notifying subscribers, who are responsible for marking them consumed.
//   - Pending broadcast block numbers are synced to the min from the pool (or deleted when empty)
//   - On reboot, backfill considers the min block number from unconsumed and pending broadcasts. Additionally, unconsumed
//     entries are removed and the pending broadcasts number updated.
type ORM interface {
	// FindBroadcasts returns broadcasts for a range of block numbers, both consumed and unconsumed.
	FindBroadcasts(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error)
	// CreateBroadcast inserts an unconsumed log broadcast for jobID.
	CreateBroadcast(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...pg.QOpt) error
	// WasBroadcastConsumed returns true if jobID consumed the log broadcast.
	WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID int32, qopts ...pg.QOpt) (bool, error)
	// MarkBroadcastConsumed marks the log broadcast as consumed by jobID.
	MarkBroadcastConsumed(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...pg.QOpt) error
	// MarkBroadcastsConsumed marks the log broadcasts as consumed by jobID.
	MarkBroadcastsConsumed(blockHashes []common.Hash, blockNumbers []uint64, logIndexes []uint, jobIDs []int32, qopts ...pg.QOpt) error
	// MarkBroadcastsUnconsumed marks all log broadcasts from all jobs on or after fromBlock as
	// unconsumed.
	MarkBroadcastsUnconsumed(fromBlock int64, qopts ...pg.QOpt) error

	// SetPendingMinBlock sets the minimum block number for which there are pending broadcasts in the pool, or nil if empty.
	SetPendingMinBlock(blockNum *int64, qopts ...pg.QOpt) error
	// GetPendingMinBlock returns the minimum block number for which there were pending broadcasts in the pool, or nil if it was empty.
	GetPendingMinBlock(qopts ...pg.QOpt) (blockNumber *int64, err error)

	// Reinitialize cleans up the database by removing any unconsumed broadcasts, then updating (if necessary) and
	// returning the pending minimum block number.
	Reinitialize(qopts ...pg.QOpt) (blockNumber *int64, err error)
}

type orm struct {
	q          pg.Q
	evmChainID utils.Big
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, evmChainID big.Int) *orm {
	return &orm{pg.NewQ(db, lggr, cfg), *utils.NewBig(&evmChainID)}
}

func (o *orm) WasBroadcastConsumed(blockHash common.Hash, logIndex uint, jobID int32, qopts ...pg.QOpt) (consumed bool, err error) {
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
	q := o.q.WithOpts(qopts...)
	err = q.Get(&consumed, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return consumed, err
}

func (o *orm) FindBroadcasts(fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error) {
	var broadcasts []LogBroadcast
	query := `
		SELECT block_hash, consumed, log_index, job_id FROM log_broadcasts
		WHERE block_number >= $1
		AND block_number <= $2
		AND evm_chain_id = $3
	`
	err := o.q.Select(&broadcasts, query, fromBlockNum, toBlockNum, o.evmChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find log broadcasts")
	}
	return broadcasts, err
}

func (o *orm) CreateBroadcast(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, updated_at, consumed, evm_chain_id)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), false, $5)
    `, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	return errors.Wrap(err, "failed to create log broadcast")
}

func (o *orm) MarkBroadcastConsumed(blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, updated_at, consumed, evm_chain_id)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), true, $5)
		ON CONFLICT (job_id, block_hash, log_index, evm_chain_id) DO UPDATE
		SET consumed = true, updated_at = NOW()
    `, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	return errors.Wrap(err, "failed to mark log broadcast as consumed")
}

// MarkBroadcastsConsumed marks many broadcasts as consumed.
// The lengths of all the provided slices must be equal, otherwise an error is returned.
func (o *orm) MarkBroadcastsConsumed(blockHashes []common.Hash, blockNumbers []uint64, logIndexes []uint, jobIDs []int32, qopts ...pg.QOpt) error {
	if !utils.AllEqual(len(blockHashes), len(blockNumbers), len(logIndexes), len(jobIDs)) {
		return fmt.Errorf("all arg slice lengths must be equal, got: %d %d %d %d",
			len(blockHashes), len(blockNumbers), len(logIndexes), len(jobIDs),
		)
	}

	type input struct {
		BlockHash   common.Hash `db:"blockHash"`
		BlockNumber uint64      `db:"blockNumber"`
		LogIndex    uint        `db:"logIndex"`
		JobID       int32       `db:"jobID"`
		ChainID     utils.Big   `db:"chainID"`
	}
	inputs := make([]input, len(blockHashes))
	query := `
INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, updated_at, consumed, evm_chain_id)
VALUES (:blockHash, :blockNumber, :logIndex, :jobID, NOW(), NOW(), true, :chainID)
ON CONFLICT (job_id, block_hash, log_index, evm_chain_id) DO UPDATE
SET consumed = true, updated_at = NOW();
	`
	for i := range blockHashes {
		inputs[i] = input{
			BlockHash:   blockHashes[i],
			BlockNumber: blockNumbers[i],
			LogIndex:    logIndexes[i],
			JobID:       jobIDs[i],
			ChainID:     o.evmChainID,
		}
	}
	q := o.q.WithOpts(qopts...)
	_, err := q.NamedExec(query, inputs)
	return errors.Wrap(err, "mark broadcasts consumed")
}

// MarkBroadcastsUnconsumed implements the ORM interface.
func (o *orm) MarkBroadcastsUnconsumed(fromBlock int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        UPDATE log_broadcasts
        SET consumed = false
        WHERE block_number >= $1
		AND evm_chain_id = $2
        `, fromBlock, o.evmChainID)
	return errors.Wrap(err, "failed to mark broadcasts unconsumed")
}

func (o *orm) Reinitialize(qopts ...pg.QOpt) (*int64, error) {
	// Minimum block number from the set of unconsumed logs, which we'll remove later.
	minUnconsumed, err := o.getUnconsumedMinBlock(qopts...)
	if err != nil {
		return nil, err
	}
	// Minimum block number from the set of pending logs in the pool.
	minPending, err := o.GetPendingMinBlock(qopts...)
	if err != nil {
		return nil, err
	}
	if minUnconsumed == nil {
		// Nothing unconsumed to consider or cleanup, and pending minimum block number still stands.
		return minPending, nil
	}
	if minPending == nil || *minUnconsumed < *minPending {
		// Use the lesser minUnconsumed.
		minPending = minUnconsumed
		// Update the db so that we can safely delete the unconsumed entries.
		if err := o.SetPendingMinBlock(minPending, qopts...); err != nil {
			return nil, err
		}
	}
	// Safe to delete old unconsumed entries since the pending minimum block covers this range.
	if err := o.removeUnconsumed(qopts...); err != nil {
		return nil, err
	}
	return minPending, nil
}

func (o *orm) SetPendingMinBlock(blockNumber *int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        INSERT INTO log_broadcasts_pending (evm_chain_id, block_number, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (evm_chain_id) DO UPDATE SET block_number = $3, updated_at = NOW() 
    `, o.evmChainID, blockNumber, blockNumber)
	return errors.Wrap(err, "failed to set pending broadcast block number")
}

func (o *orm) GetPendingMinBlock(qopts ...pg.QOpt) (*int64, error) {
	q := o.q.WithOpts(qopts...)
	var blockNumber *int64
	err := q.Get(&blockNumber, `
        SELECT block_number FROM log_broadcasts_pending WHERE evm_chain_id = $1
    `, o.evmChainID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to get broadcasts pending number")
	}
	return blockNumber, nil
}

func (o *orm) getUnconsumedMinBlock(qopts ...pg.QOpt) (*int64, error) {
	q := o.q.WithOpts(qopts...)
	var blockNumber *int64
	err := q.Get(&blockNumber, `
        SELECT min(block_number) FROM log_broadcasts
			WHERE evm_chain_id = $1
			AND consumed = false
			AND block_number IS NOT NULL
    `, o.evmChainID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to get unconsumed broadcasts min block number")
	}
	return blockNumber, nil
}

func (o *orm) removeUnconsumed(qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
        DELETE FROM log_broadcasts
			WHERE evm_chain_id = $1
			AND consumed = false
			AND block_number IS NOT NULL
    `, o.evmChainID)
	return errors.Wrap(err, "failed to delete unconsumed broadcasts")
}

// LogBroadcast - data from log_broadcasts table columns
type LogBroadcast struct {
	BlockHash common.Hash
	Consumed  bool
	LogIndex  uint
	JobID     int32
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
