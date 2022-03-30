package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type LogPoller struct {
	ec   client.Client
	orm  *ORM
	lggr logger.Logger
	// poll period set by block production rate
	pollPeriod                       time.Duration
	finalityDepth, backfillBatchSize int64

	filterMu  sync.Mutex
	addresses map[common.Address]struct{}
	topics    map[int]map[common.Hash]struct{}

	replay chan int64
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func NewLogPoller(orm *ORM, ec client.Client, lggr logger.Logger, pollPeriod time.Duration, finalityDepth, backfillBatchSize int64) *LogPoller {
	return &LogPoller{
		ec:                ec,
		orm:               orm,
		lggr:              lggr,
		replay:            make(chan int64),
		done:              make(chan struct{}),
		pollPeriod:        pollPeriod,
		finalityDepth:     finalityDepth,
		backfillBatchSize: backfillBatchSize,
		addresses:         make(map[common.Address]struct{}),
		topics:            make(map[int]map[common.Hash]struct{}),
	}
}

// MergeFilter will update the filter with the new topics and addresses.
// Clients may chose to MergeFilter and then replay in order to ensure desired logs are present.
func (lp *LogPoller) MergeFilter(topics [][]common.Hash, addresses []common.Address) {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	for _, addr := range addresses {
		lp.addresses[addr] = struct{}{}
	}
	// [[A, B], [C]] + [[D], [], [E]] = [[A, B, D], [C], [E]]
	for i := 0; i < len(topics); i++ {
		if lp.topics[i] == nil {
			lp.topics[i] = make(map[common.Hash]struct{})
		}
		for j := 0; j < len(topics[i]); j++ {
			lp.topics[i][topics[i][j]] = struct{}{}
		}
	}
}

func (lp *LogPoller) FilterAddresses() []common.Address {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	var addresses []common.Address
	for addr := range lp.addresses {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (lp *LogPoller) FilterTopics() [][]common.Hash {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	var topics [][]common.Hash
	for i := 0; i < len(lp.topics); i++ {
		var topicPosition []common.Hash
		// Order not important within each position.
		for topic := range lp.topics[i] {
			topicPosition = append(topicPosition, topic)
		}
		topics = append(topics, topicPosition)
	}
	return topics
}

// Replay signals that the poller should resume from a new block.
func (lp *LogPoller) Replay(fromBlock int64) {
	lp.replay <- fromBlock
}

func (lp *LogPoller) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	lp.ctx = ctx
	lp.cancel = cancel
	go lp.run()
	return nil
}

func (lp *LogPoller) Close() error {
	lp.cancel()
	<-lp.done
	return nil
}

func (lp *LogPoller) Ready() error {
	return nil
}

func (lp *LogPoller) Healthy() error {
	return nil
}

func (lp *LogPoller) run() {
	defer close(lp.done)
	tick := time.After(0)
	var start int64
	for {
		select {
		case <-lp.ctx.Done():
			return
		case fromBlock := <-lp.replay:
			start = fromBlock
		case <-tick:
			if start == 0 {
				lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(lp.ctx))
				if err != nil {
					lp.lggr.Warnw("unable to get starting block", "err", err)
					continue
				}
				start = lastProcessed.BlockNumber + 1
				continue
			}
			start = lp.PollAndSaveLogs(lp.ctx, start)
			tick = time.After(lp.pollPeriod)
		}
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func convertLogs(chainID *big.Int, logs []types.Log) []Log {
	var lgs []Log
	for _, l := range logs {
		lgs = append(lgs, Log{
			EvmChainId: utils.NewBig(chainID),
			LogIndex:   int64(l.Index),
			BlockHash:  l.BlockHash,
			// We assume block numbers fit in int64
			// in many places.
			BlockNumber: int64(l.BlockNumber),
			Topics:      convertTopics(l.Topics),
			Address:     l.Address,
			TxHash:      l.TxHash,
			Data:        l.Data,
		})
	}
	return lgs
}

func convertTopics(topics []common.Hash) [][]byte {
	var topicsForDB [][]byte
	for _, t := range topics {
		topicsForDB = append(topicsForDB, t.Bytes())
	}
	return topicsForDB
}

func (lp *LogPoller) backfill(ctx context.Context, start, end int64) int64 {
	for i := start; i <= end; i += lp.backfillBatchSize {
		var (
			logs []types.Log
			err  error
		)
		utils.RetryWithBackoff(ctx, func() bool {
			logs, err = lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
				FromBlock: big.NewInt(i),
				ToBlock:   big.NewInt(min(i+lp.backfillBatchSize, end)),
				Addresses: lp.FilterAddresses(),
				Topics:    lp.FilterTopics(),
			})
			if err != nil {
				lp.lggr.Warnw("Unable query for logs, retrying", "err", err, "from", i, "to", min(i+lp.backfillBatchSize, end))
				return true
			}
			return false
		})
		if len(logs) == 0 {
			continue
		}
		utils.RetryWithBackoff(ctx, func() bool {
			if err := lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs)); err != nil {
				lp.lggr.Warnw("Unable to insert logs logs, retrying", "err", err, "from", i, "to", min(i+lp.backfillBatchSize, end))
				return true
			}
			return false
		})
	}
	return end + 1
}

// PollAndSaveLogs On startup/crash current is the first block after the last processed block.
func (lp *LogPoller) PollAndSaveLogs(ctx context.Context, current int64) int64 {
	// Get latest block on chain
	latestBlock, err1 := lp.ec.BlockByNumber(ctx, nil)
	if err1 != nil {
		lp.lggr.Warnw("Unable to get latest block", "err", err1, "current", current)
		return current
	}
	latest := latestBlock.Number().Int64()
	// 1<-2<-3(current)<-4<-5<-6<-7(latest). Finality is 2, so 3,4,5 can be batched.
	// start = current = 3, end = latest-current+1 = 7-3+1 = 5 (inclusive range).
	if (latest - current) > lp.finalityDepth {
		current = lp.backfill(ctx, current, latest-current+1)
	}

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
			// If err is not a no rows error, assume transient db issue and retry
			lp.lggr.Warnw("Unable to read latest block saved", "err", err, "current", current)
			return current
		}
		// We will not have the previous block on initial poll or after a backfill.
		havePreviousBlock := !errors.Is(err, sql.ErrNoRows)
		if havePreviousBlock && !bytes.Equal(block.ParentHash().Bytes(), expectedParent.BlockHash.Bytes()) {
			// There can be another reorg while we're finding the LCA.
			// That is ok, since we'll detect it on the next iteration of current.
			lca, err := lp.findLCA(block.ParentHash())
			if err != nil {
				lp.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err)
				return current
			}

			// We truncate all the blocks after the LCA.
			// TODO: We could mark all the logs after this reorg to be excluded
			// from canonical queries?
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
			Addresses: lp.FilterAddresses(),
			Topics:    lp.FilterTopics(),
		})
		if err != nil {
			lp.lggr.Warnw("Unable query for logs, retrying", "err", err, "block", block.Number())
			return current
		}
		err = lp.orm.q.Transaction(func(q pg.Queryer) error {
			if err := lp.orm.InsertBlock(block.Hash(), block.Number().Int64()); err != nil {
				return err
			}
			if len(logs) == 0 {
				return nil
			}
			return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs))
		})
		if err != nil {
			// If we're unable to insert, don't increment current and just retry
			lp.lggr.Warnw("Unable to save logs, retrying", "err", err, "block", block.Number())
			return current
		}
		current++
	}
	return current
}

func (lp *LogPoller) findLCA(h common.Hash) (int64, error) {
	// Find the first place where our chain and their chain have the same block,
	// that block number is the LCA.
	block, err := lp.ec.BlockByHash(context.Background(), h)
	if err != nil {
		return 0, err
	}
	ourBlockHash, err := lp.orm.SelectBlockByNumber(block.Number().Int64())
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(block.Hash().Bytes(), ourBlockHash.BlockHash.Bytes()) {
		return lp.findLCA(block.ParentHash())
	}
	// If we do have the blockhash, that is the LCA
	return block.Number().Int64(), nil
}

func (lp *LogPoller) CanonicalLogs(start, end int64, topics []common.Hash, address common.Address) ([]Log, error) {
	return lp.orm.SelectCanonicalLogsByBlockRangeTopicAddress(start, end, address, convertTopics(topics))
}
