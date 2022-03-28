package logpoller

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type LogPoller struct {
	ec client.Client
	orm *ORM
	lggr logger.Logger
	addresses []common.Address
	topics [][]common.Hash
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
	// load last block processed and its hash from the db
	start := int64(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			_ = lp.pollAndSaveLogs(ctx, start)
			tick = time.After(15*time.Second) // poll period set by block production rate
		}
	}
}

// On startup/crash start is the most recent saveBlockHash
func (lp *LogPoller) pollAndSaveLogs(ctx context.Context, current int64) error {
	// Get latest block on chain
	latestBlock, err := lp.ec.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}
	latest := latestBlock.Number().Int64();
	// TODO: batch fill up until unfinalized logs.
	for current <= latest {
		block, _ := lp.ec.BlockByNumber(ctx, big.NewInt(current))
		// Does this block point to the same parent that we have saved?
		// If not, there was a reorg, so we need to rewind.
		expectedParent, _ := lp.orm.SelectBlockByNumber(current - 1)
		if !bytes.Equal(block.ParentHash().Bytes(), expectedParent.BlockHash.Bytes()) {
			// There can be another reorg while we're finding the LCA.
			// That is ok, since we'll detect it on the next iteration of current.
			current = lp.findLCA(block.ParentHash())
			// We truncate all the blocks we had up until
			// TODO: Handle cleaning up on restart
			_ = lp.orm.DeleteRangeBlocks(block.Number().Int64(), current)
			continue
		}
		h := block.Hash()
		logs, err := lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
			BlockHash: &h,
			Addresses: lp.addresses,
			Topics:    lp.topics,
		})
		if err != nil {
			// TODO log warn
			continue
		}
		err = lp.orm.q.Transaction(func(q pg.Queryer) error {
			if err := lp.orm.InsertBlock(block.Hash(), block.Number().Int64()); err != nil {
				return err
			}
			if len(logs) > 0 {
				var lgs []Log
				for _, l := range logs {
					lgs = append(lgs, Log{
						EvmChainId:     utils.NewBig(lp.ec.ChainID()),
						LogIndex:       int64(l.Index),
						BlockHash:      l.BlockHash,
						// We assume block numbers fit in int64
						// in many places.
						BlockNumber:    int64(l.BlockNumber),
						// TODO: Event signature
						//EventSignature: l.,
						Address:        l.Address,
						TxHash:         l.TxHash,
						Data:           l.Data,
					})
				}
				return lp.orm.InsertLogs(lgs)
			}
			return nil
		})
		if err != nil {
			// If we're unable to insert just retry
			// TODO log warn
			continue
		}

		// Continue from the end boundary block.
		// Ok to poll the same block again, saveLogs is idempotent.
		current++
	}
	return nil
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
