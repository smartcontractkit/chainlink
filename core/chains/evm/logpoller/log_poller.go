package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type LogPoller struct {
	ec        client.Client
	orm       *ORM
	lggr      logger.Logger
	addresses []common.Address
	topics    [][]common.Hash
}

func NewLogPoller(orm *ORM, ec client.Client, lggr logger.Logger) *LogPoller {
	return &LogPoller{ec: ec, orm: orm, lggr: lggr}
}

func (lp *LogPoller) Start(ctx context.Context) error {
	go lp.run(ctx)
	return nil
}

func (lp *LogPoller) run(ctx context.Context) {
	tick := time.After(0)
	var start *int64
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			if start == nil {
				lastProcessed, err := lp.orm.SelectLatestBlock()
				if err != nil {
					lp.lggr.Warnw("unable to get starting block", "err", err)
					continue
				}
				start = &lastProcessed.BlockNumber
			}
			newStart := lp.pollAndSaveLogs(ctx, *start)
			start = &newStart
			tick = time.After(15 * time.Second) // poll period set by block production rate
		}
	}
}

// On startup/crash start is the most recent saveBlockHash
func (lp *LogPoller) pollAndSaveLogs(ctx context.Context, current int64) int64 {
	// Get latest block on chain
	latestBlock, err := lp.ec.BlockByNumber(ctx, nil)
	if err != nil {
		return current
	}
	latest := latestBlock.Number().Int64()
	// TODO: batch FilterLogs up until unfinalized logs.

	for current <= latest {
		block, err := lp.ec.BlockByNumber(ctx, big.NewInt(current))
		if err != nil {
			lp.lggr.Warnw("Unable to get block", "err", err, "current", current)
			return current
		}
		// Does this block point to the same parent that we have saved?
		// If not, there was a reorg, so we need to rewind.
		expectedParent, err := lp.orm.SelectBlockByNumber(current - 1)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			// If not a no rows error, assume transient db issue and retry
			return current
		}
		// The very first time we poll, we will not have the previous block.
		havePreviousBlock := !errors.Is(err, sql.ErrNoRows)
		if havePreviousBlock && !bytes.Equal(block.ParentHash().Bytes(), expectedParent.BlockHash.Bytes()) {
			// There can be another reorg while we're finding the LCA.
			// That is ok, since we'll detect it on the next iteration of current.
			lca := lp.findLCA(block.ParentHash())
			// We truncate all the blocks after the LCA.
			err = lp.orm.DeleteRangeBlocks(lca+1, latest)
			if err != nil {
				lp.lggr.Warnw("Unable to clear reorged blocks, retrying", "err", err)
				return current
			}
			current = lca + 1
			continue
		}
		h := block.Hash()
		logs, err := lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
			BlockHash: &h,
			Addresses: lp.addresses,
			Topics:    lp.topics,
		})
		if err != nil {
			lp.lggr.Warnw("Unable query for logs, retrying", "err", err, "block", block.Number())
			return current
		}
		err = lp.orm.q.Transaction(func(q pg.Queryer) error {
			if err := lp.orm.InsertBlock(block.Hash(), block.Number().Int64()); err != nil {
				return err
			}
			if len(logs) > 0 {
				var lgs []Log
				for _, l := range logs {
					var topics [][]byte
					for _, t := range l.Topics {
						topics = append(topics, t.Bytes())
					}
					lgs = append(lgs, Log{
						EvmChainId: utils.NewBig(lp.ec.ChainID()),
						LogIndex:   int64(l.Index),
						BlockHash:  l.BlockHash,
						// We assume block numbers fit in int64
						// in many places.
						BlockNumber: int64(l.BlockNumber),
						Topics:      topics,
						Address:     l.Address,
						TxHash:      l.TxHash,
						Data:        l.Data,
					})
				}
				return lp.orm.InsertLogs(lgs)
			}
			return nil
		})
		if err != nil {
			// If we're unable to insert just retry
			lp.lggr.Warnw("Unable to save logs, retrying", "err", err, "block", block.Number())
			return current
		}

		// Continue from the end boundary block.
		// Ok to poll the same block again, saveLogs is idempotent.
		current++
	}
	return current
}

func (lp *LogPoller) findLCA(h common.Hash) int64 {
	// Find the first place where our chain and their chain have the same block,
	// that block number is the LCA.
	block, _ := lp.ec.BlockByHash(context.Background(), h)
	ourBlockHash, _ := lp.orm.SelectBlockByNumber(block.Number().Int64())
	if !bytes.Equal(block.Hash().Bytes(), ourBlockHash.BlockHash.Bytes()) {
		return lp.findLCA(block.ParentHash())
	}
	// If we do have the blockhash, that is the LCA
	return block.Number().Int64()
}

func (lp *LogPoller) Stop() error {
	return nil
}

func (lp *LogPoller) CanonicalLogs(start, end int64, topic common.Hash) []Log {
	// Group by blockhash and take the latest
	return nil
}
