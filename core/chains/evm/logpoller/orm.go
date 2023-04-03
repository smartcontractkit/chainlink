package logpoller

import (
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

type ORM struct {
	chainID *big.Int
	q       pg.Q
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *ORM {
	namedLogger := lggr.Named("Configs")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		chainID: chainID,
		q:       q,
	}
}

// InsertBlock is idempotent to support replays.
func (o *ORM) InsertBlock(h common.Hash, n int64, t time.Time, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`INSERT INTO evm_log_poller_blocks (evm_chain_id, block_hash, block_number, block_timestamp, created_at) 
      VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT DO NOTHING`, utils.NewBig(o.chainID), h[:], n, t)
	return err
}

// InsertFilter is idempotent.
//
// Each address/event pair must have a unique job id, so it may be removed when the job is deleted.
// If a second job tries to overwrite the same pair, this should fail.
func (o *ORM) InsertFilter(filter Filter, qopts ...pg.QOpt) (err error) {
	q := o.q.WithOpts(qopts...)
	addresses := make([][]byte, 0)
	events := make([][]byte, 0)

	for _, addr := range filter.Addresses {
		addresses = append(addresses, addr.Bytes())
	}
	for _, ev := range filter.EventSigs {
		events = append(events, ev.Bytes())
	}
	return q.ExecQ(`INSERT INTO evm_log_poller_filters
								(name, evm_chain_id, created_at, address, event)
								SELECT * FROM
									(SELECT $1, $2::NUMERIC, NOW()) x,
									(SELECT unnest($3::BYTEA[]) addr) a,
									(SELECT unnest($4::BYTEA[]) ev) e
								ON CONFLICT (name, evm_chain_id, address, event) DO NOTHING`,
		filter.Name, utils.NewBig(o.chainID), addresses, events)
}

// DeleteFilter removes all events,address pairs associated with the Filter
func (o *ORM) DeleteFilter(name string, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm_log_poller_filters WHERE name = $1 AND evm_chain_id = $2`, name, utils.NewBig(o.chainID))
}

// LoadFiltersForChain returns all filters for this chain
func (o *ORM) LoadFilters(qopts ...pg.QOpt) (map[string]Filter, error) {
	q := o.q.WithOpts(qopts...)
	rows := make([]Filter, 0)
	err := q.Select(&rows, `SELECT name, ARRAY_AGG(DISTINCT address)::BYTEA[] AS addresses, ARRAY_AGG(DISTINCT event)::BYTEA[] AS event_sigs
									FROM evm_log_poller_filters WHERE evm_chain_id = $1 GROUP BY name`, utils.NewBig(o.chainID))
	filters := make(map[string]Filter)
	for _, filter := range rows {
		filters[filter.Name] = filter
	}

	return filters, err
}

func (o *ORM) SelectBlockByHash(h common.Hash, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm_log_poller_blocks WHERE block_hash = $1 AND evm_chain_id = $2`, h, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectBlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm_log_poller_blocks WHERE block_number = $1 AND evm_chain_id = $2`, n, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1`, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectLatestLogEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	q := o.q.WithOpts(qopts...)
	var l Log
	if err := q.Get(&l, `SELECT * FROM evm_logs 
         WHERE evm_chain_id = $1 
            AND event_sig = $2 
            AND address = $3 
            AND (block_number + $4) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
        ORDER BY (block_number, log_index) DESC LIMIT 1`, utils.NewBig(o.chainID), eventSig, address, confs); err != nil {
		return nil, err
	}
	return &l, nil
}

// DeleteBlocksAfter delete all blocks after and including start.
func (o *ORM) DeleteBlocksAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm_log_poller_blocks WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
}

// DeleteBlocksBefore delete all blocks before and including end.
func (o *ORM) DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM evm_log_poller_blocks WHERE block_number <= $1 AND evm_chain_id = $2`, end, utils.NewBig(o.chainID))
	return err
}

func (o *ORM) DeleteLogsAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM evm_logs WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
}

// InsertLogs is idempotent to support replays.
func (o *ORM) InsertLogs(logs []Log, qopts ...pg.QOpt) error {
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

		err := q.ExecQNamed(`INSERT INTO evm_logs 
(evm_chain_id, log_index, block_hash, block_number, block_timestamp, address, event_sig, topics, tx_hash, data, created_at) VALUES 
(:evm_chain_id, :log_index, :block_hash, :block_number, :block_timestamp, :address, :event_sig, :topics, :tx_hash, :data, NOW()) ON CONFLICT DO NOTHING`, logs[start:end])

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *ORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
	var logs []Log
	err := o.q.Select(&logs, `
        SELECT * FROM evm_logs 
        WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3
        ORDER BY (block_number, log_index, created_at)`, start, end, utils.NewBig(o.chainID))
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsByBlockRangeFilter finds the logs in a given block range.
func (o *ORM) SelectLogsByBlockRangeFilter(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm_logs 
			WHERE evm_logs.block_number >= $1 AND evm_logs.block_number <= $2 AND evm_logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5 
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig.Bytes())
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectLogsWithSigsByBlockRangeFilter finds the logs in the given block range with the given event signatures
// emitted from the given address.
func (o *ORM) SelectLogsWithSigsByBlockRangeFilter(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) (logs []Log, err error) {
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
FROM evm_logs
WHERE evm_logs.block_number BETWEEN :start AND :end
	AND evm_logs.evm_chain_id = :chainid
	AND evm_logs.address = :address
	AND evm_logs.event_sig IN (:EventSigs)
ORDER BY (evm_logs.block_number, evm_logs.log_index)`, a)
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

func (o *ORM) GetBlocksRange(start uint64, end uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	var blocks []LogPollerBlock
	q := o.q.WithOpts(qopts...)
	err := q.Select(&blocks, `
        SELECT * FROM evm_log_poller_blocks 
        WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3
        ORDER BY block_number ASC`, start, end, utils.NewBig(o.chainID))
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SelectLatestLogEventSigsAddrsWithConfs finds the latest log by (address, event) combination that matches a list of Addresses and list of events
func (o *ORM) SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log

	var sigs [][]byte
	for _, sig := range eventSigs {
		sigs = append(sigs, sig.Bytes())
	}
	var addrs [][]byte
	for _, addr := range addresses {
		addrs = append(addrs, addr.Bytes())
	}

	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm_logs WHERE (block_number, address, event_sig) IN (
			SELECT MAX(block_number), address, event_sig FROM evm_logs 
				WHERE evm_chain_id = $1 AND
				    event_sig = ANY($2) AND
					address = ANY($3) AND
		   			block_number > $4 AND
					(block_number + $5) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			GROUP BY event_sig, address
		)
		ORDER BY block_number ASC
	`, o.chainID.Int64(), pq.ByteaArray(sigs), pq.ByteaArray(addrs), fromBlock, confs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return logs, nil
}

func (o *ORM) SelectDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm_logs 
			WHERE evm_logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND substring(data from 32*$4+1 for 32) <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), wordValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm_logs 
			WHERE evm_logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm_logs 
			WHERE evm_logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM evm_logs 
			WHERE evm_logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND topics[$4] <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), topicValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	q := o.q.WithOpts(qopts...)
	var logs []Log
	var topicValuesBytes [][]byte
	for _, topicValue := range topicValues {
		topicValuesBytes = append(topicValuesBytes, topicValue.Bytes())
	}
	// Add 1 since postgresql arrays are 1-indexed.
	err := q.Select(&logs, `
		SELECT * FROM evm_logs 
			WHERE evm_logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] = ANY($5)
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, pq.ByteaArray(topicValuesBytes), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SelectIndexedLogsByBlockRangeFilter finds the indexed logs in a given block range.
func (o *ORM) SelectIndexedLogsByBlockRangeFilter(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	var logs []Log
	var topicValuesBytes [][]byte
	for _, topicValue := range topicValues {
		topicValuesBytes = append(topicValuesBytes, topicValue.Bytes())
	}
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM evm_logs 
			WHERE evm_logs.block_number >= $1 AND evm_logs.block_number <= $2 AND evm_logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5
			AND topics[$6] = ANY($7)
			ORDER BY (evm_logs.block_number, evm_logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, pq.ByteaArray(topicValuesBytes))
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

// SelectIndexedLogsWithSigsExcluding query's for logs that have signature A and exclude logs that have a corresponding signature B, matching is done based on the topic index both logs should be inside the block range and have the minimum number of confirmations
func (o *ORM) SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := validateTopicIndex(topicIndex); err != nil {
		return nil, err
	}

	q := o.q.WithOpts(qopts...)
	var logs []Log

	err := q.Select(&logs, `
		SELECT *
		FROM   evm_logs
		WHERE  evm_chain_id = $1
		AND    address = $2
		AND    event_sig = $3
		AND block_number BETWEEN $6 AND $7
		AND (block_number + $8) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
		
		EXCEPT
		
		SELECT     a.*
		FROM       evm_logs AS a
		INNER JOIN evm_logs B
		ON         a.evm_chain_id = b.evm_chain_id
		AND        a.address = b.address
		AND        a.topics[$5] = b.topics[$5]
		AND        a.event_sig = $3
		AND        b.event_sig = $4
	    AND 	   b.block_number BETWEEN $6 AND $7
		AND (b.block_number + $8) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)

		ORDER BY block_number,log_index ASC
			`, utils.NewBig(o.chainID), address, sigA.Bytes(), sigB.Bytes(), topicIndex+1, startBlock, endBlock, confs)
	if err != nil {
		return nil, err
	}
	return logs, nil

}
