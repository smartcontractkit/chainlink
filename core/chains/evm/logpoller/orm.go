package logpoller

import (
	"fmt"
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
func NewORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *ORM {
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
	_, err := q.Exec(`INSERT INTO log_poller_blocks (evm_chain_id, block_hash, block_number, created_at) 
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

func (o *ORM) DeleteRangeBlocks(start, end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM log_poller_blocks WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3`, start, end, utils.NewBig(o.chainID))
	return err
}

func (o *ORM) DeleteLogs(start, end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM logs WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3`, start, end, utils.NewBig(o.chainID))
	return err
}

// InsertLogs is idempotent to support replays.
func (o *ORM) InsertLogs(logs []Log, qopts ...pg.QOpt) error {
	for _, log := range logs {
		if o.chainID.Cmp(log.EvmChainId.ToInt()) != 0 {
			return errors.Errorf("invalid chainID in log got %v want %v", log.EvmChainId.ToInt(), o.chainID)
		}
	}
	q := o.q.WithOpts(qopts...)
	_, err := q.NamedExec(`INSERT INTO logs 
(evm_chain_id, log_index, block_hash, block_number, address, event_sig, topics, tx_hash, data, created_at) VALUES 
(:evm_chain_id, :log_index, :block_hash, :block_number, :address, :event_sig, :topics, :tx_hash, :data, NOW()) ON CONFLICT DO NOTHING`, logs)
	return err
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

// SelectLogsByBlockRangeFilter finds the latest logs by block.
// Assumes that logs inserted later for a given block are "more" canonical.
func (o *ORM) SelectLogsByBlockRangeFilter(start, end int64, address common.Address, eventSig []byte, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs, `
		SELECT * FROM logs 
			WHERE logs.block_number >= $1 AND logs.block_number <= $2 AND logs.evm_chain_id = $3 
			AND address = $4 AND event_sig = $5 
			ORDER BY (logs.block_number, logs.log_index)`, start, end, utils.NewBig(o.chainID), address, eventSig)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// LatestLogEventSigsAddrs finds the latest log by (address, event) combination that matches a list of addresses and list of events
func (o *ORM) LatestLogEventSigsAddrs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log

	sigs := [][]byte{}
	for _, sig := range eventSigs {
		sigs = append(sigs, sig.Bytes())
	}
	addrs := [][]byte{}
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
		   			block_number > $4
			GROUP BY event_sig, address
		)
		ORDER BY block_number ASC
	`, o.chainID.Int64(), pq.Array(sigs), pq.Array(addrs), fromBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return logs, nil
}

func (o *ORM) SelectDataWordRange(address common.Address, eventSig []byte, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND encode(substring(data from 32*$4+1 for 32), 'hex') >= $5
			AND encode(substring(data from 32*$4+1 for 32), 'hex') <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig, wordIndex, wordValueMin.String()[2:], wordValueMax.String()[2:], confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectDataWordGreaterThan(address common.Address, eventSig []byte, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND encode(substring(data from 32*$4+1 for 32), 'hex') >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig, wordIndex, wordValueMin.String()[2:], confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicGreaterThan(address common.Address, eventSig []byte, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND encode(topics[$4], 'hex') >= $5
			AND (block_number + $6) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig, topicIndex+1, topicValueMin.String()[2:], confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexLogsTopicRange(address common.Address, eventSig []byte, topicIndex int, topicValueMin, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	err := q.Select(&logs,
		`SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND encode(topics[$4], 'hex') >= $5
			AND encode(topics[$4], 'hex') <= $6
			AND (block_number + $7) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, utils.NewBig(o.chainID), address, eventSig, topicIndex+1, topicValueMin.String()[2:], topicValueMax.String()[2:], confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectIndexedLogs(address common.Address, eventSig []byte, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	var logs []Log
	q := o.q.WithOpts(qopts...)
	topicValuesList := fmt.Sprintf("'%s'", topicValues[0].String()[2:])
	for _, topicValue := range topicValues[1:] {
		topicValuesList += fmt.Sprintf(",'%s'", topicValue.String()[2:])
	}
	// Add 1 since arrays are 1-indexed.
	err := q.Select(&logs, fmt.Sprintf(`
		SELECT * FROM logs 
			WHERE logs.evm_chain_id = $1
			AND address = $2 AND event_sig = $3
			AND encode(topics[$4], 'hex') IN (%s) 
			AND (block_number + $5) <= (SELECT COALESCE(block_number, 0) FROM log_poller_blocks WHERE evm_chain_id = $1 ORDER BY block_number DESC LIMIT 1)
			ORDER BY (logs.block_number, logs.log_index)`, topicValuesList), utils.NewBig(o.chainID), address, eventSig, topicIndex+1, confs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
