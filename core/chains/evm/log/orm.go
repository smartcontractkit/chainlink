package log

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// ORM is the interface for log broadcasts.
//   - Unconsumed broadcasts are created just before notifying subscribers, who are responsible for marking them consumed.
//   - Pending broadcast block numbers are synced to the min from the pool (or deleted when empty)
//   - On reboot, backfill considers the min block number from unconsumed and pending broadcasts. Additionally, unconsumed
//     entries are removed and the pending broadcasts number updated.
type ORM interface {
	// FindBroadcasts returns broadcasts for a range of block numbers, both consumed and unconsumed.
	FindBroadcasts(ctx context.Context, fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error)
	// CreateBroadcast inserts an unconsumed log broadcast for jobID.
	CreateBroadcast(ctx context.Context, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error
	// WasBroadcastConsumed returns true if jobID consumed the log broadcast.
	WasBroadcastConsumed(ctx context.Context, blockHash common.Hash, logIndex uint, jobID int32) (bool, error)
	// MarkBroadcastConsumed marks the log broadcast as consumed by jobID.
	MarkBroadcastConsumed(ctx context.Context, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error
	// MarkBroadcastsUnconsumed marks all log broadcasts from all jobs on or after fromBlock as
	// unconsumed.
	MarkBroadcastsUnconsumed(ctx context.Context, fromBlock int64) error

	// SetPendingMinBlock sets the minimum block number for which there are pending broadcasts in the pool, or nil if empty.
	SetPendingMinBlock(ctx context.Context, blockNum *int64) error
	// GetPendingMinBlock returns the minimum block number for which there were pending broadcasts in the pool, or nil if it was empty.
	GetPendingMinBlock(ctx context.Context) (blockNumber *int64, err error)

	// Reinitialize cleans up the database by removing any unconsumed broadcasts, then updating (if necessary) and
	// returning the pending minimum block number.
	Reinitialize(ctx context.Context) (blockNumber *int64, err error)

	WithDataSource(sqlutil.DataSource) ORM
}

type orm struct {
	ds         sqlutil.DataSource
	evmChainID ubig.Big
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, evmChainID big.Int) *orm {
	return &orm{ds, *ubig.New(&evmChainID)}
}

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM {
	return &orm{ds, o.evmChainID}
}

func (o *orm) WasBroadcastConsumed(ctx context.Context, blockHash common.Hash, logIndex uint, jobID int32) (consumed bool, err error) {
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
	err = o.ds.GetContext(ctx, &consumed, query, args...)
	if pkgerrors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return consumed, err
}

func (o *orm) FindBroadcasts(ctx context.Context, fromBlockNum int64, toBlockNum int64) ([]LogBroadcast, error) {
	var broadcasts []LogBroadcast
	query := `
		SELECT block_hash, consumed, log_index, job_id FROM log_broadcasts
		WHERE block_number >= $1
		AND block_number <= $2
		AND evm_chain_id = $3
	`
	err := o.ds.SelectContext(ctx, &broadcasts, query, fromBlockNum, toBlockNum, o.evmChainID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to find log broadcasts")
	}
	return broadcasts, err
}

func (o *orm) CreateBroadcast(ctx context.Context, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error {
	_, err := o.ds.ExecContext(ctx, `
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, updated_at, consumed, evm_chain_id)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), false, $5)
    `, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	return pkgerrors.Wrap(err, "failed to create log broadcast")
}

func (o *orm) MarkBroadcastConsumed(ctx context.Context, blockHash common.Hash, blockNumber uint64, logIndex uint, jobID int32) error {
	_, err := o.ds.ExecContext(ctx, `
        INSERT INTO log_broadcasts (block_hash, block_number, log_index, job_id, created_at, updated_at, consumed, evm_chain_id)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), true, $5)
		ON CONFLICT (job_id, block_hash, log_index, evm_chain_id) DO UPDATE
		SET consumed = true, updated_at = NOW()
    `, blockHash, blockNumber, logIndex, jobID, o.evmChainID)
	return pkgerrors.Wrap(err, "failed to mark log broadcast as consumed")
}

// MarkBroadcastsUnconsumed implements the ORM interface.
func (o *orm) MarkBroadcastsUnconsumed(ctx context.Context, fromBlock int64) error {
	_, err := o.ds.ExecContext(ctx, `
        UPDATE log_broadcasts
        SET consumed = false
        WHERE block_number >= $1
		AND evm_chain_id = $2
        `, fromBlock, o.evmChainID)
	return pkgerrors.Wrap(err, "failed to mark broadcasts unconsumed")
}

func (o *orm) Reinitialize(ctx context.Context) (*int64, error) {
	// Minimum block number from the set of unconsumed logs, which we'll remove later.
	minUnconsumed, err := o.getUnconsumedMinBlock(ctx)
	if err != nil {
		return nil, err
	}
	// Minimum block number from the set of pending logs in the pool.
	minPending, err := o.GetPendingMinBlock(ctx)
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
		if err := o.SetPendingMinBlock(ctx, minPending); err != nil {
			return nil, err
		}
	}
	// Safe to delete old unconsumed entries since the pending minimum block covers this range.
	if err := o.removeUnconsumed(ctx); err != nil {
		return nil, err
	}
	return minPending, nil
}

func (o *orm) SetPendingMinBlock(ctx context.Context, blockNumber *int64) error {
	_, err := o.ds.ExecContext(ctx, `
        INSERT INTO log_broadcasts_pending (evm_chain_id, block_number, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (evm_chain_id) DO UPDATE SET block_number = $3, updated_at = NOW() 
    `, o.evmChainID, blockNumber, blockNumber)
	return pkgerrors.Wrap(err, "failed to set pending broadcast block number")
}

func (o *orm) GetPendingMinBlock(ctx context.Context) (*int64, error) {
	var blockNumber *int64
	err := o.ds.GetContext(ctx, &blockNumber, `
        SELECT block_number FROM log_broadcasts_pending WHERE evm_chain_id = $1
    `, o.evmChainID)
	if pkgerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get broadcasts pending number")
	}
	return blockNumber, nil
}

func (o *orm) getUnconsumedMinBlock(ctx context.Context) (*int64, error) {
	var blockNumber *int64
	err := o.ds.GetContext(ctx, &blockNumber, `
        SELECT min(block_number) FROM log_broadcasts
			WHERE evm_chain_id = $1
			AND consumed = false
			AND block_number IS NOT NULL
    `, o.evmChainID)
	if pkgerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get unconsumed broadcasts min block number")
	}
	return blockNumber, nil
}

func (o *orm) removeUnconsumed(ctx context.Context) error {
	_, err := o.ds.ExecContext(ctx, `
        DELETE FROM log_broadcasts
			WHERE evm_chain_id = $1
			AND consumed = false
			AND block_number IS NOT NULL
    `, o.evmChainID)
	return pkgerrors.Wrap(err, "failed to delete unconsumed broadcasts")
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
