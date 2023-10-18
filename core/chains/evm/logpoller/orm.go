package logpoller

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ORM represents the persistent data access layer used by the log poller. At this moment, it's a bit leaky abstraction, because
// it exposes some of the database implementation details (e.g. pg.Q). Ideally it should be agnostic and could be applied to any persistence layer.
// What is more, LogPoller should not be aware of the underlying database implementation and delegate all the queries to the ORM.
type ORM interface {
	Q() pg.Q
	InsertLogs(logs []Log, qopts ...pg.QOpt) error
	InsertBlock(blockHash common.Hash, blockNumber int64, blockTimestamp time.Time, lastFinalizedBlockNumber int64, qopts ...pg.QOpt) error
	InsertFilter(filter Filter, qopts ...pg.QOpt) error

	LoadFilters(qopts ...pg.QOpt) (map[string]Filter, error)
	DeleteFilter(name string, qopts ...pg.QOpt) error

	DeleteBlocksAfter(start int64, qopts ...pg.QOpt) error
	DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error
	DeleteLogsAfter(start int64, qopts ...pg.QOpt) error
	DeleteExpiredLogs(qopts ...pg.QOpt) error

	GetBlocksRange(start int64, end int64, qopts ...pg.QOpt) ([]LogPollerBlock, error)
	SelectBlockByNumber(blockNumber int64, qopts ...pg.QOpt) (*LogPollerBlock, error)
	SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error)

	SelectLogs(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsWithSigs(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsCreatedAfter(address common.Address, eventSig common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectLatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs Confirmations, qopts ...pg.QOpt) (*Log, error)
	SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) (int64, error)

	SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsByBlockRange(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsCreatedAfter(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectIndexedLogsByTxHash(eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	SelectLogsDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
}

type DbORM struct {
	chainID *big.Int
	q       pg.Q
}

// NewORM creates a DbORM scoped to chainID.
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
func (o *DbORM) InsertBlock(blockHash common.Hash, blockNumber int64, blockTimestamp time.Time, finalizedBlock int64, qopts ...pg.QOpt) error {
	args, err := newQueryArgs(o.chainID).
		withCustomHashArg("block_hash", blockHash).
		withCustomArg("block_number", blockNumber).
		withCustomArg("block_timestamp", blockTimestamp).
		withCustomArg("finalized_block_number", finalizedBlock).
		toArgs()
	if err != nil {
		return err
	}
	return o.q.WithOpts(qopts...).ExecQNamed(`
			INSERT INTO evm.log_poller_blocks 
				(evm_chain_id, block_hash, block_number, block_timestamp, finalized_block_number, created_at) 
      		VALUES (:evm_chain_id, :block_hash, :block_number, :block_timestamp, :finalized_block_number, NOW()) 
			ON CONFLICT DO NOTHING`, args)
}

// InsertFilter is idempotent.
//
// Each address/event pair must have a unique job id, so it may be removed when the job is deleted.
// If a second job tries to overwrite the same pair, this should fail.
func (o *DbORM) InsertFilter(filter Filter, qopts ...pg.QOpt) (err error) {
	args, err := newQueryArgs(o.chainID).
		withCustomArg("name", filter.Name).
		withCustomArg("retention", filter.Retention).
		withAddressArray(filter.Addresses).
		withEventSigArray(filter.EventSigs).
		toArgs()
	if err != nil {
		return err
	}
	// '::' has to be escaped in the query string
	// https://github.com/jmoiron/sqlx/issues/91, https://github.com/jmoiron/sqlx/issues/428
	return o.q.WithOpts(qopts...).ExecQNamed(`
		INSERT INTO evm.log_poller_filters
	  		(name, evm_chain_id, retention, created_at, address, event)
		SELECT * FROM
			(SELECT :name, :evm_chain_id ::::NUMERIC, :retention ::::BIGINT, NOW()) x,
			(SELECT unnest(:address_array ::::BYTEA[]) addr) a,
			(SELECT unnest(:event_sig_array ::::BYTEA[]) ev) e
		ON CONFLICT (name, evm_chain_id, address, event) 
		DO UPDATE SET retention=:retention ::::BIGINT`, args)
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

func (o *DbORM) SelectBlockByHash(hash common.Hash, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm.log_poller_blocks WHERE block_hash = $1 AND evm_chain_id = $2`, hash, utils.NewBig(o.chainID)); err != nil {
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

func (o *DbORM) SelectLatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs Confirmations, qopts ...pg.QOpt) (*Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT * FROM evm.logs
			WHERE evm_chain_id = :evm_chain_id
			AND event_sig = :event_sig
			AND address = :address
			AND block_number <= %s
			ORDER BY (block_number, log_index) DESC LIMIT 1`, nestedBlockNumberQuery(confs))
	var l Log
	if err := o.q.WithOpts(qopts...).GetNamed(query, &l, args); err != nil {
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

		err := q.ExecQNamed(`
			INSERT INTO evm.logs 
				(evm_chain_id, log_index, block_hash, block_number, block_timestamp, address, event_sig, topics, tx_hash, data, created_at) 
				VALUES 
				(:evm_chain_id, :log_index, :block_hash, :block_number, :block_timestamp, :address, :event_sig, :topics, :tx_hash, :data, NOW()) 
				ON CONFLICT DO NOTHING`, logs[start:end])

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
	args, err := newQueryArgs(o.chainID).
		withStartBlock(start).
		withEndBlock(end).
		toArgs()
	if err != nil {
		return nil, err
	}

	var logs []Log
	err = o.q.SelectNamed(&logs, `
        SELECT * FROM evm.logs 
        	WHERE evm_chain_id = :evm_chain_id
        	AND block_number >= :start_block 
        	AND block_number <= :end_block 
        	ORDER BY (block_number, log_index, created_at)`, args)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsByBlockRangeFilter finds the logs in a given block range.
func (o *DbORM) SelectLogs(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withStartBlock(start).
		withEndBlock(end).
		toArgs()
	if err != nil {
		return nil, err
	}
	var logs []Log
	err = o.q.WithOpts(qopts...).SelectNamed(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id 
			AND address = :address
			AND event_sig = :event_sig  
			AND block_number >= :start_block 
			AND block_number <= :end_block
			ORDER BY (block_number, log_index)`, args)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsCreatedAfter finds logs created after some timestamp.
func (o *DbORM) SelectLogsCreatedAfter(address common.Address, eventSig common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withBlockTimestampAfter(after).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT * FROM evm.logs 
				WHERE evm_chain_id = :evm_chain_id
				AND address = :address
				AND event_sig = :event_sig
				AND block_timestamp > :block_timestamp_after
				AND block_number <= %s
				ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))

	var logs []Log
	if err = o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsWithSigsByBlockRangeFilter finds the logs in the given block range with the given event signatures
// emitted from the given address.
func (o *DbORM) SelectLogsWithSigs(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) (logs []Log, err error) {
	args, err := newQueryArgs(o.chainID).
		withAddress(address).
		withEventSigArray(eventSigs).
		withStartBlock(start).
		withEndBlock(end).
		toArgs()
	if err != nil {
		return nil, err
	}

	q := o.q.WithOpts(qopts...)
	err = q.SelectNamed(&logs, `
			SELECT * FROM evm.logs
				WHERE evm_chain_id = :evm_chain_id
				AND address = :address
				AND event_sig = ANY(:event_sig_array)
				AND block_number BETWEEN :start_block AND :end_block
				ORDER BY (block_number, log_index)`, args)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return logs, err
}

func (o *DbORM) GetBlocksRange(start int64, end int64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	args, err := newQueryArgs(o.chainID).
		withStartBlock(start).
		withEndBlock(end).
		toArgs()
	if err != nil {
		return nil, err
	}
	var blocks []LogPollerBlock
	err = o.q.WithOpts(qopts...).SelectNamed(&blocks, `
        SELECT * FROM evm.log_poller_blocks 
			WHERE block_number >= :start_block 
			AND block_number <= :end_block
			AND evm_chain_id = :evm_chain_id
			ORDER BY block_number ASC`, args)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SelectLatestLogEventSigsAddrsWithConfs finds the latest log by (address, event) combination that matches a list of Addresses and list of events
func (o *DbORM) SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgs(o.chainID).
		withAddressArray(addresses).
		withEventSigArray(eventSigs).
		withStartBlock(fromBlock).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT * FROM evm.logs WHERE (block_number, address, event_sig) IN (
			SELECT MAX(block_number), address, event_sig FROM evm.logs 
				WHERE evm_chain_id = :evm_chain_id 
				AND event_sig = ANY(:event_sig_array) 
				AND address = ANY(:address_array) 
				AND block_number > :start_block 
				AND block_number <= %s
			GROUP BY event_sig, address
		)
		ORDER BY block_number ASC`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err := o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return logs, nil
}

// SelectLatestBlockNumberEventSigsAddrsWithConfs finds the latest block number that matches a list of Addresses and list of events. It returns 0 if there is no matching block
func (o *DbORM) SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) (int64, error) {
	args, err := newQueryArgs(o.chainID).
		withEventSigArray(eventSigs).
		withAddressArray(addresses).
		withStartBlock(fromBlock).
		withConfs(confs).
		toArgs()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(block_number), 0) FROM evm.logs
			WHERE evm_chain_id = :evm_chain_id 
			AND event_sig = ANY(:event_sig_array) 
			AND address = ANY(:address_array) 
			AND block_number > :start_block 
			AND block_number <= %s`, nestedBlockNumberQuery(confs))
	var blockNumber int64
	if err := o.q.WithOpts(qopts...).GetNamed(query, &blockNumber, args); err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (o *DbORM) SelectLogsDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withWordIndex(wordIndex).
		withWordValueMin(wordValueMin).
		withWordValueMax(wordValueMax).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id
			AND address = :address 
			AND event_sig = :event_sig
			AND substring(data from 32*:word_index+1 for 32) >= :word_value_min
			AND substring(data from 32*:word_index+1 for 32) <= :word_value_max
			AND block_number <= %s
			ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err := o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectLogsDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withWordIndex(wordIndex).
		withWordValueMin(wordValueMin).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id
			AND address = :address
			AND event_sig = :event_sig
			AND substring(data from 32*:word_index+1 for 32) >= :word_value_min
			AND block_number <= %s
			ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err = o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withTopicIndex(topicIndex).
		withTopicValueMin(topicValueMin).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT * FROM evm.logs
			WHERE evm_chain_id = :evm_chain_id
			AND address = :address 
			AND event_sig = :event_sig
			AND topics[:topic_index] >= :topic_value_min
			AND block_number <= %s
			ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err = o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withTopicIndex(topicIndex).
		withTopicValueMin(topicValueMin).
		withTopicValueMax(topicValueMax).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
			SELECT * FROM evm.logs 
				WHERE evm_chain_id = :evm_chain_id
				AND address = :address
				AND event_sig = :event_sig
				AND topics[:topic_index] >= :topic_value_min
				AND topics[:topic_index] <= :topic_value_max
				AND block_number <= %s
			ORDER BY (evm.logs.block_number, evm.logs.log_index)`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err := o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withTopicIndex(topicIndex).
		withTopicValues(topicValues).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id
			AND address = :address
			AND event_sig = :event_sig
			AND topics[:topic_index] = ANY(:topic_values)
			AND block_number <= %s
			ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))
	var logs []Log
	if err := o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectIndexedLogsByBlockRangeFilter finds the indexed logs in a given block range.
func (o *DbORM) SelectIndexedLogsByBlockRange(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withTopicIndex(topicIndex).
		withTopicValues(topicValues).
		withStartBlock(start).
		withEndBlock(end).
		toArgs()
	if err != nil {
		return nil, err
	}
	var logs []Log
	err = o.q.WithOpts(qopts...).SelectNamed(&logs, `
		SELECT * FROM evm.logs 
				WHERE evm_chain_id = :evm_chain_id 
				AND address = :address
				AND event_sig = :event_sig
				AND topics[:topic_index] = ANY(:topic_values)
				AND block_number >= :start_block
				AND block_number <= :end_block
				ORDER BY (block_number, log_index)`, args)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsCreatedAfter(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgsForEvent(o.chainID, address, eventSig).
		withBlockTimestampAfter(after).
		withConfs(confs).
		withTopicIndex(topicIndex).
		withTopicValues(topicValues).
		toArgs()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id
			AND address = :address
			AND event_sig = :event_sig
			AND topics[:topic_index] = ANY(:topic_values)
			AND block_timestamp > :block_timestamp_after
			AND block_number <= %s
			ORDER BY (block_number, log_index)`, nestedBlockNumberQuery(confs))

	var logs []Log
	if err = o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *DbORM) SelectIndexedLogsByTxHash(eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgs(o.chainID).
		withTxHash(txHash).
		withEventSig(eventSig).
		toArgs()
	if err != nil {
		return nil, err
	}
	var logs []Log
	err = o.q.WithOpts(qopts...).SelectNamed(&logs, `
		SELECT * FROM evm.logs 
			WHERE evm_chain_id = :evm_chain_id
			AND tx_hash = :tx_hash
			AND event_sig = :event_sig
			ORDER BY (block_number, log_index)`, args)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectIndexedLogsWithSigsExcluding query's for logs that have signature A and exclude logs that have a corresponding signature B, matching is done based on the topic index both logs should be inside the block range and have the minimum number of confirmations
func (o *DbORM) SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	args, err := newQueryArgs(o.chainID).
		withAddress(address).
		withTopicIndex(topicIndex).
		withStartBlock(startBlock).
		withEndBlock(endBlock).
		withCustomHashArg("sigA", sigA).
		withCustomHashArg("sigB", sigB).
		withConfs(confs).
		toArgs()
	if err != nil {
		return nil, err
	}

	nestedQuery := nestedBlockNumberQuery(confs)
	query := fmt.Sprintf(`
		SELECT * FROM   evm.logs
		WHERE   evm_chain_id = :evm_chain_id
		AND     address = :address
		AND     event_sig = :sigA
		AND 	block_number BETWEEN :start_block AND :end_block
		AND 	block_number <= %s		
		EXCEPT
		SELECT     a.* FROM       evm.logs AS a
		INNER JOIN evm.logs B
		ON         a.evm_chain_id = b.evm_chain_id
		AND        a.address = b.address
		AND        a.topics[:topic_index] = b.topics[:topic_index]
		AND        a.event_sig = :sigA
		AND        b.event_sig = :sigB
	    AND 	   b.block_number BETWEEN :start_block AND :end_block
		AND		   b.block_number <= %s
		ORDER BY block_number,log_index ASC`, nestedQuery, nestedQuery)
	var logs []Log
	if err := o.q.WithOpts(qopts...).SelectNamed(&logs, query, args); err != nil {
		return nil, err
	}
	return logs, nil
}

func nestedBlockNumberQuery(confs Confirmations) string {
	if confs == Finalized {
		return `
				(SELECT finalized_block_number 
				FROM evm.log_poller_blocks 
				WHERE evm_chain_id = :evm_chain_id 
				ORDER BY block_number DESC LIMIT 1) `
	}
	// Intentionally wrap with greatest() function and don't return negative block numbers when :confs > :block_number
	// It doesn't impact logic of the outer query, because block numbers are never less or equal to 0 (guarded by log_poller_blocks_block_number_check)
	return `
			(SELECT greatest(block_number - :confs, 0) 
			FROM evm.log_poller_blocks 	
			WHERE evm_chain_id = :evm_chain_id 
			ORDER BY block_number DESC LIMIT 1) `

}
