package logpoller

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
	if err := q.Get(&b, `SELECT * FROM log_poller_blocks WHERE block_number = (SELECT MAX(block_number) FROM log_poller_blocks) AND evm_chain_id = $1`, utils.NewBig(o.chainID)); err != nil {
		return nil, err
	}
	return &b, nil
}

func (o *ORM) DeleteRangeBlocks(start, end int64, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.Exec(`DELETE FROM log_poller_blocks WHERE block_number >= $1 AND block_number <= $2 AND evm_chain_id = $3`, start, end, utils.NewBig(o.chainID))
	return err
}

func (o *ORM) InsertLogs(logs []Log, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	_, err := q.NamedExec(`INSERT INTO logs 
(evm_chain_id, log_index, block_hash, block_number, address, topics, tx_hash, data, created_at) VALUES 
(:evm_chain_id, :log_index, :block_hash, :block_number, :address, :topics, :tx_hash, :data, NOW()) ON CONFLICT DO NOTHING`, logs)
	return err
}

func (o *ORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
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

func (o *ORM) SelectCanonicalLogsByBlockRange(start, end int64) ([]Log, error) {
	var logs []Log
	err := o.q.Select(&logs, `
		SELECT * FROM logs 
            JOIN (SELECT block_number, max(created_at) AS created_at FROM logs GROUP BY block_number) AS latest_blocks 
			ON latest_blocks.block_number = logs.block_number AND latest_blocks.created_at = logs.created_at
		 	WHERE logs.block_number >= $1 AND logs.block_number <= $2 AND logs.evm_chain_id = $3 
            ORDER BY (logs.block_number, logs.log_index)`, start, end, utils.NewBig(o.chainID))
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (o *ORM) SelectCanonicalLogsByBlockRangeTopicAddress(start, end int64, address common.Address, topics [][]byte) ([]Log, error) {
	var logs []Log
	err := o.q.Select(&logs, `
		SELECT * FROM logs 
            JOIN (SELECT block_number, max(created_at) AS created_at FROM logs GROUP BY block_number) AS latest_blocks 
			ON latest_blocks.block_number = logs.block_number AND latest_blocks.created_at = logs.created_at
		 	WHERE logs.block_number >= $1 AND logs.block_number <= $2 AND logs.evm_chain_id = $3 
				AND address = $4 AND topics = $5 
            ORDER BY (logs.block_number, logs.log_index)`, start, end, utils.NewBig(o.chainID), address, topics)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
