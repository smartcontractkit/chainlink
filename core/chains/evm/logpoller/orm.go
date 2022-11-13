package logpoller

import (
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ORM struct {
	chainID *big.Int
	q       pg.Q
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *ORM {
	namedLogger := lggr.Named("ORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &ORM{
		chainID: chainID,
		q:       q,
	}
}

// InsertBlock is idempotent to support replays.
func (o *ORM) InsertBlock(h common.Hash, n int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`INSERT INTO log_poller_blocks (evm_chain_id, block_hash, block_number, created_at) 
      VALUES ($1, $2, $3, NOW()) ON CONFLICT DO NOTHING`, utils.NewBig(o.chainID), h[:], n)
	return err
}

func (o *ORM) SelectBlockByHash(h common.Hash, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM log_poller_blocks WHERE block_hash = $1 AND evm_chain_id = $2`, h, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectBlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM log_poller_blocks WHERE block_number = $1 AND evm_chain_id = $2`, n, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error) {
	q := o.q.WithOpts(qopts...)
	var b LogPollerBlock
	if err := q.Get(&b, `SELECT * FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1`, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) SelectLatestLogEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	q := o.q.WithOpts(qopts...)
	var l Log
	if err := q.Get(&l, `SELECT * FROM logs 
         WHERE evm_chain_id = $1 
            AND event_sig = $2 
            AND address = $3 
            AND (block_number + $4) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
        ORDER BY (block_number, log_index) DESC LIMIT 1`, utils.NewBig(o.chainID), eventSig, address, confs); err != nil {
		return nil, err
	}
	return &l, nil
}

// DeleteBlocksAfter delete all blocks after and including start.
func (o *ORM) DeleteBlocksAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM log_poller_blocks WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
}

// DeleteBlocksBefore delete all blocks before and including end.
func (o *ORM) DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM log_poller_blocks WHERE block_number <= $1 AND evm_chain_id = $2`, end, utils.NewBig(o.chainID))
	return err
}

func (o *ORM) DeleteLogsAfter(start int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	return q.ExecQ(`DELETE FROM logs WHERE block_number >= $1 AND evm_chain_id = $2`, start, utils.NewBig(o.chainID))
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

		err := q.ExecQNamed(`INSERT INTO logs 
(evm_chain_id, log_index, block_hash, block_number, address, event_sig, topics, tx_hash, data, created_at) VALUES 
(:evm_chain_id, :log_index, :block_hash, :block_number, :address, :event_sig, :topics, :tx_hash, :data, NOW()) ON CONFLICT DO NOTHING`, logs[start:end])

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *ORM) selectLogsByBlockRange(start, end int64) ([]Log, error) {
	var logs []Log
	err := o.q.Select(&logs, `
        SELECT * FROM logs 
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
		SELECT * FROM logs 
			WHERE logs.block_number >= $1 AND logs.block_number <= $2 AND logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5 
			ORDER BY (logs.block_number, logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig.Bytes())
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
FROM logs
WHERE logs.block_number BETWEEN :start AND :end
	AND logs.evm_chain_id = :chainid
	AND logs.address = :address
	AND logs.event_sig IN (:EventSigs)
ORDER BY (logs.block_number, logs.log_index)`, a)
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

func (o *ORM) GetBlocks(blockNumbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	if len(blockNumbers) == 0 {
		return nil, nil
	}

	var blocks []LogPollerBlock
	q := o.q.WithOpts(qopts...)
	a := map[string]any{
		"blockNumbers": blockNumbers,
		"chainid":      utils.NewBig(o.chainID),
	}
	query, args, err := sqlx.Named(
		`
SELECT
	*
FROM log_poller_blocks 
WHERE evm_chain_id = :chainid
	AND block_number IN (:blockNumbers)
`, a)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx Named")
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx In")
	}
	query = q.Rebind(query)
	err = q.Select(&blocks, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return blocks, err
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
		SELECT * FROM logs WHERE (block_number, address, event_sig) IN (
			SELECT MAX(block_number), address, event_sig FROM logs 
				WHERE evm_chain_id = $1 AND
				    event_sig = ANY($2) AND
					address = ANY($3) AND
		   			block_number > $4 AND
					(block_number + $5) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
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
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND substring(data from 32*$4+1 for 32) <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), wordValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND substring(data from 32*$4+1 for 32) >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), wordIndex, wordValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] >= $5
			AND topics[$4] <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, topicValueMin.Bytes(), topicValueMax.Bytes(), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	q := o.q.WithOpts(qopts...)
	var logs []Log
	var topicValuesBytes [][]byte
	for _, topicValue := range topicValues {
		topicValuesBytes = append(topicValuesBytes, topicValue.Bytes())
	}
	// Add 1 since postgresql arrays are 1-indexed.
	err := q.Select(&logs, `
		SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND topics[$4] = ANY($5)
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig.Bytes(), topicIndex+1, pq.ByteaArray(topicValuesBytes), confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
