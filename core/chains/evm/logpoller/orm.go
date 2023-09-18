package logpoller

import (
	"context"
	"database/sql"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ORM represents the persistent data access layer used by the log poller. At this moment, it's a bit leaky, because
// it exposes some of the database implementation details (e.g. pg.Q). Ideally it should be agnostic and could be applied to any persistence layer
type ORM interface {
	Q() pg.Q
	InsertLogs(logs []Log, qopts ...pg.QOpt) error
	InsertBlock(h common.Hash, n int64, t time.Time, qopts ...pg.QOpt) error
	InsertFilter(filter Filter, qopts ...pg.QOpt) error

	LoadFilters(qopts ...pg.QOpt) (map[string]Filter, error)
	DeleteFilter(name string, qopts ...pg.QOpt) error

	DeleteBlocksAfter(start int64, qopts ...pg.QOpt) error
	DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error
	DeleteLogsAfter(start int64, qopts ...pg.QOpt) error
	DeleteExpiredLogs(qopts ...pg.QOpt) error

	GetBlocksRange(start uint64, end uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error)
	SelectBlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error)
	SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error)

	SelectLogs(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsWithSigs(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsCreatedAfter(eventSig common.Hash, address common.Address, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectLatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error)
	SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) (int64, error)

	SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsByBlockRange(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsCreatedAfter(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsUntilBlockHashDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, untilBlockHash common.Hash, qopts ...pg.QOpt) ([]Log, error)
}

type DbORM struct {
	chainID *big.Int
	q       pg.Q
}

// NewORM creates an DbORM scoped to chainID.
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *DbORM {
	namedLogger := lggr.Named("Configs")
	q := pg.NewQ(db, namedLogger, cfg)
	return &DbORM{
		chainID: chainID,
		q:       q,
	}
}

func (o *DbORM) Q() pg.Q {
	return o.q
}

// InsertBlock is idempotent to support replays.
func (o *DbORM) InsertBlock(h common.Hash, n int64, t time.Time, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`INSERT INTO evm.log_poller_blocks (evm_chain_id, block_hash, block_number, block_timestamp, created_at) 
      VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT DO NOTHING`, utils.NewBig(o.chainID), h[:], n, t)
	return err
}

// InsertFilter is idempotent.
//
// Each address/event pair must have a unique job id, so it may be removed when the job is deleted.
// If a second job tries to overwrite the same pair, this should fail.
func (o *DbORM) InsertFilter(filter Filter, qopts ...pg.QOpt) (err error) {
	q := o.q.WithOpts(qopts...)
	addresses := make([][]byte, 0)
	events := make([][]byte, 0)

	for _, addr := range filter.Addresses {
		addresses = append(addresses, addr.Bytes())
	}
	for _, ev := range filter.EventSigs {
		events = append(events, ev.Bytes())
	}
	return q.ExecQ(`INSERT INTO evm.log_poller_filters
	  (name, evm_chain_id, retention, created_at, address, event)
		SELECT * FROM
			(SELECT $1, $2::NUMERIC, $3::BIGINT, NOW()) x,
			(SELECT unnest($4::BYTEA[]) addr) a,
			(SELECT unnest($5::BYTEA[]) ev) e
		ON CONFLICT (name, evm_chain_id, address, event) DO UPDATE SET retention=$3::BIGINT;`,
		filter.Name, utils.NewBig(o.chainID), filter.Retention, addresses, events)
}

// DeleteFilter removes all events,address pairs associated with the Filter
func (o *DbORM) DeleteFilter(name string, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm.log_poller_filters WHERE name = $1 AND evm_chain_id = $2`, name, utils.NewBig(o.chainID))
}

// LoadFiltersForChain returns all filters for this chain
func (o *DbORM) LoadFilters(qopts ...pg.QOpt) (map[string]Filter, error) {
	q := o.q.WithOpts(qopts...)
	rows := make([]Filter, 0)
	err := q.Select(&rows, `SELECT name,
			ARRAY_AGG(DISTINCT address)::BYTEA[] AS addresses, 
			ARRAY_AGG(DISTINCT event)::BYTEA[] AS event_sigs,
			MAX(retention) AS retention
		FROM evm.log_poller_filters WHERE evm_chain_id = $1
		GROUP BY name`, utils.NewBig(o.chainID))
	filters := make(map[string]Filter)
	for _, filter := range rows {
		filters[filter.Name] = filter
	}

	return filters, err
}

func (o *DbORM) SelectBlockByHash(h common.Hash, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm.log_poller_blocks WHERE block_hash = $1 AND evm_chain_id = $2`, h, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *DbORM) SelectBlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm.log_poller_blocks WHERE block_number = $1 AND evm_chain_id = $2`, n, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *DbORM) SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1`, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *DbORM) SelectLatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	q := o.q.WithOpts(qopts...)
	var l Log
	if err := q.Get(&l, `SELECT * FROM evm.logs 
         WHERE evm_chain_id = $1 
            AND event_sig = $2 
            AND address = $3 
            AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $4
        ORDER BY (block_number, log_index) DESC LIMIT 1`, utils.NewBig(o.chainID), eventSig, address, confs); err != nil {
		return nil, err
	}
	return &l, nil
}

// DeleteBlocksAfter delete all blocks after and including start.
func (o *DbORM) DeleteBlocksAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm.log_poller_blocks WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
}

// DeleteBlocksBefore delete all blocks before and including end.
func (o *DbORM) DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM evm.log_poller_blocks WHERE block_number <= $1 AND evm_chain_id = $2`, end, utils.NewBig(o.chainID))
	return err
}

func (o *DbORM) DeleteLogsAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm.logs WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
}

type Exp struct {
	Address      common.Address
	EventSig     common.Hash
	Expiration   time.Time
	TimeNow      time.Time
	ShouldDelete bool
}

func (o *DbORM) DeleteExpiredLogs(qopts ...pg.QOpt) error {
	qopts = append(qopts, pg.WithLongQueryTimeout())
	q := o.q.WithOpts(qopts...)

	return q.ExecQ(`WITH r AS
		( SELECT address, event, MAX(retention) AS retention
			FROM evm.log_poller_filters WHERE evm_chain_id=$1 
			GROUP BY evm_chain_id,address, event HAVING NOT 0 = ANY(ARRAY_AGG(retention))
		) DELETE FROM evm.logs l USING r
			WHERE l.evm_chain_id = $1 AND l.address=r.address AND l.event_sig=r.event
			AND l.created_at <= STATEMENT_TIMESTAMP() - (r.retention / 10^9 * interval '1 second')`, // retention is in nanoseconds (time.Duration aka BIGINT)
		utils.NewBig(o.chainID))
}

// InsertLogs is idempotent to support replays.
func (o *DbORM) InsertLogs(logs []Log, qopts ...pg.QOpt) error {
	for _, log := range logs {
		if o.chainID.Cmp(log.EvmChainId.ToInt()) != 0 {
			return errors.Errorf("invalid chainID in log got %v want %v", log.EvmChainId.ToInt(), o.chainID)
		}
	}
	q := o.q.WithOpts(qopts...)

	batchInsertSize := 4000
	for i := 0; i < len(logs); i += batchInsertSize {
		start, end := i, i+batchInsertSize
		if end > len(logs) {
			end = len(logs)
		}

		err := q.ExecQNamed(`INSERT INTO evm.logs 
(evm_chain_id, log_index, block_hash, block_number, block_timestamp, address, event_sig, topics, tx_hash, data, created_at) VALUES 
(:evm_chain_id, :log_index, :block_hash, :block_number, :block_timestamp, :address, :event_sig, :topics, :tx_hash, :data, NOW()) ON CONFLICT DO NOTHING`, logs[start:end])

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) && batchInsertSize > 500 {
				// In case of DB timeouts, try to insert again with a smaller batch upto a limit
				batchInsertSize /= 2
				i -= batchInsertSize // counteract +=batchInsertSize on next loop iteration
				continue
			}
			return err
		}
	}

	return nil
}

func (o *DbORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
	var logs []Log
	err := o.q.Select(&logs, `
        SELECT * FROM evm.logs 
        WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3
        ORDER BY (block_number, log_index, created_at)`, start, end, utils.NewBig(o.chainID))
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsByBlockRangeFilter finds the logs in a given block range.
func (o *DbORM) SelectLogs(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm.logs.block_number >= $1 AND evm.logs.block_number <= $2 AND evm.logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5 
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig.Bytes())
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsCreatedAfter finds logs created after some timestamp.
func (o *DbORM) SelectLogsCreatedAfter(address common.Address, eventSig common.Hash, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error) {
	minBlock, maxBlock, err := o.blocksRangeAfterTimestamp(after, confs, qopts...)
	if err != nil {
		return nil, err
	}

	var logs []Log
	q := o.q.WithOpts(qopts...)
	err = q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = $1 
			AND address = $2 
			AND event_sig = $3 	
			AND block_number > $4
			AND block_number <= $5
			ORDER BY (block_number, log_index)`, utils.NewBig(o.chainID), address, eventSig, minBlock, maxBlock)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsWithSigsByBlockRangeFilter finds the logs in the given block range with the given event signatures
// emitted from the given address.
func (o *DbORM) SelectLogsWithSigs(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) (logs []Log, err error) {
	q := o.q.WithOpts(qopts...)
	sigs := make([][]byte, 0, len(eventSigs))
	for _, sig := range eventSigs {
		sigs = append(sigs, sig.Bytes())
	}
	a := map[string]any{
		"start":     start,
		"end":       end,
		"chainid":   utils.NewBig(o.chainID),
		"address":   address,
		"EventSigs": sigs,
	}
	query, args, err := sqlx.Named(
		`
SELECT
	*
FROM evm.logs
WHERE evm.logs.block_number BETWEEN :start AND :end
	AND evm.logs.evm_chain_id = :chainid
	AND evm.logs.address = :address
	AND evm.logs.event_sig IN (:EventSigs)
ORDER BY (evm.logs.block_number, evm.logs.log_index)`, a)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx Named")
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx In")
	}
	query = q.Rebind(query)
	err = q.Select(&logs, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return logs, err
}

func (o *DbORM) GetBlocksRange(start uint64, end uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	var blocks []LogPollerBlock
	q := o.q.WithOpts(qopts...)
	err := q.Select(&blocks, `
        SELECT * FROM evm.log_poller_blocks 
        WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3
        ORDER BY block_number ASC`, start, end, utils.NewBig(o.chainID))
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SelectLatestLogEventSigsAddrsWithConfs finds the latest log by (address, event) combination that matches a list of Addresses and list of events
func (o *DbORM) SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	sigs := concatBytes(eventSigs)
	addrs := concatBytes(addresses)

	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm.logs WHERE (block_number, address, event_sig) IN (
			SELECT MAX(block_number), address, event_sig FROM evm.logs 
				WHERE evm_chain_id = $1 AND
				    event_sig = ANY($2) AND
					address = ANY($3) AND
		   			block_number > $4 AND
					block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $5
			GROUP BY event_sig, address
		)
		ORDER BY block_number ASC
	`, o.chainID.Int64(), sigs, addrs, fromBlock, confs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return logs, nil
}

// SelectLatestBlockNumberEventSigsAddrsWithConfs finds the latest block number that matches a list of Addresses and list of events. It returns 0 if there is no matching block
func (o *DbORM) SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) (int64, error) {
	var blockNumber int64
	sigs := concatBytes(eventSigs)
	addrs := concatBytes(addresses)

	q := o.q.WithOpts(qopts...)
	err := q.Get(&blockNumber, `
			SELECT COALESCE(MAX(block_number), 0) FROM evm.logs 
				WHERE evm_chain_id = $1 AND
				    event_sig = ANY($2) AND
					address = ANY($3) AND
					block_number > $4 AND
					block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $5`,
		o.chainID.Int64(), sigs, addrs, fromBlock, confs)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (o *DbORM) SelectLogsDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND substring(data from 32*$4+1 for 32) <= $6
			AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $7
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), wordValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectLogsDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $6
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $6
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectLogsUntilBlockHashDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, untilBlockHash common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Transaction(func(tx pg.Queryer) error {
		// We want to mimic the behaviour of the ETH RPC which errors if blockhash not found.
		var block LogPollerBlock
		if err := tx.Get(&block,
			`SELECT * FROM evm.log_poller_blocks 
					WHERE evm_chain_id = $1 AND block_hash = $2`, utils.NewBig(o.chainID), untilBlockHash); err != nil {
			return err
		}
		return q.Select(&logs,
			`SELECT * FROM evm.logs 
			WHERE evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND block_number <= $6 
			ORDER BY (block_number, log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), block.BlockNumber)
	})
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND topics[$4] <= $6
			AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $7
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), topicValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	q := o.q.WithOpts(qopts...)
	var logs []Log
	topicValuesBytes := concatBytes(topicValues)
	// Add 1 since postgresql arrays are 1-indexed.
	err := q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] = ANY($5)
			AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $6
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValuesBytes, confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectIndexedLogsByBlockRangeFilter finds the indexed logs in a given block range.
func (o *DbORM) SelectIndexedLogsByBlockRange(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	topicValuesBytes := concatBytes(topicValues)
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm.logs.block_number >= $1 AND evm.logs.block_number <= $2 AND evm.logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5
			AND topics[$6] = ANY($7)
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValuesBytes)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func validateTopicIndex(index int) error {
	// Only topicIndex 1 through 3 is valid. 0 is the event sig and only 4 total topics are allowed
	if !(index == 1 || index == 2 || index == 3) {
		return errors.Errorf("invalid index for topic: %d", index)
	}
	return nil
}

func (o *DbORM) SelectIndexedLogsCreatedAfter(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error) {
	minBlock, maxBlock, err := o.blocksRangeAfterTimestamp(after, confs, qopts...)
	if err != nil {
		return nil, err
	}
	var logs []Log
	q := o.q.WithOpts(qopts...)
	topicValuesBytes := concatBytes(topicValues)
	// Add 1 since postgresql arrays are 1-indexed.
	err = q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND address = $2 
			AND event_sig = $3
			AND topics[$4] = ANY($5)
			AND block_number > $6
			AND block_number <= $7
			ORDER BY (block_number, log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValuesBytes, minBlock, maxBlock)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsByTxHash(eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	q := o.q.WithOpts(qopts...)
	var logs []Log
	err := q.Select(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm.logs.evm_chain_id = $1
			AND tx_hash = $2
			AND event_sig = $3
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`,
		utils.NewBig(o.chainID), txHash.Bytes(), eventSig.Bytes())
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectIndexedLogsWithSigsExcluding query's for logs that have signature A and exclude logs that have a corresponding signature B, matching is done based on the topic index both logs should be inside the block range and have the minimum number of confirmations
func (o *DbORM) SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	q := o.q.WithOpts(qopts...)
	var logs []Log

	err := q.Select(&logs, `
		SELECT *
		FROM   evm.logs
		WHERE  evm_chain_id = $1
		AND    address = $2
		AND    event_sig = $3
		AND block_number BETWEEN $6 AND $7
		AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $8
		
		EXCEPT
		
		SELECT     a.*
		FROM       evm.logs AS a
		INNER JOIN evm.logs B
		ON         a.evm_chain_id = b.evm_chain_id
		AND        a.address = b.address
		AND        a.topics[$5] = b.topics[$5]
		AND        a.event_sig = $3
		AND        b.event_sig = $4
	    AND 	   b.block_number BETWEEN $6 AND $7
		AND		   b.block_number <= (SELECT COALESCE(block_number, 0) FROM evm.log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1) - $8

		ORDER BY block_number,log_index ASC
			`, utils.NewBig(o.chainID), address, sigA.Bytes(), sigB.Bytes(), topicIndex+1, startBlock, endBlock, confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) blocksRangeAfterTimestamp(after time.Time, confs int, qopts ...pg.QOpt) (int64, int64, error) {
	type blockRange struct {
		MinBlockNumber int64 `db:"min_block"`
		MaxBlockNumber int64 `db:"max_block"`
	}

	var br blockRange
	q := o.q.WithOpts(qopts...)
	err := q.Get(&br, `
		SELECT 
		    coalesce(min(block_number), 0) as min_block, 
		    coalesce(max(block_number), 0) as max_block
		FROM evm.log_poller_blocks 
		WHERE evm_chain_id = $1
		AND block_timestamp > $2`, utils.NewBig(o.chainID), after)
	if err != nil {
		return 0, 0, err
	}
	return br.MinBlockNumber, br.MaxBlockNumber - int64(confs), nil
}

type bytesProducer interface {
	Bytes() []byte
}

func concatBytes[T bytesProducer](byteSlice []T) pq.ByteaArray {
	var output [][]byte
	for _, b := range byteSlice {
		output = append(output, b.Bytes())
	}
	return output
}
