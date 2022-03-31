package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"math/big"
	"sort"
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
	pollPeriod time.Duration // poll period set by block production rate
	finalityDepth int64 // finality depth is taken to mean that block (head - finality) is finalized
	backfillBatchSize int64 // batch size to use when backfilling finalized logs

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
func (lp *LogPoller) MergeFilter(topics []common.Hash, address common.Address) {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	lp.addresses[address] = struct{}{}
	// [[A, B], [C]] + [[D], [], [E]] = [[A, B, D], [C], [E]]
	for i := 0; i < len(topics); i++ {
		if lp.topics[i] == nil {
			lp.topics[i] = make(map[common.Hash]struct{})
		}
		lp.topics[i][topics[i]] = struct{}{}
	}
}

func (lp *LogPoller) filterAddresses() []common.Address {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	var addresses []common.Address
	for addr := range lp.addresses {
		addresses = append(addresses, addr)
	}
	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i][:], addresses[j][:]) < 0
	})
	return addresses
}

func (lp *LogPoller) filterTopics() [][]common.Hash {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	var topics [][]common.Hash
	for idx := 0; idx < len(lp.topics); idx++ {
		var topicPosition []common.Hash
		for topic := range lp.topics[idx] {
			topicPosition = append(topicPosition, topic)
		}
		sort.Slice(topicPosition, func(i, j int) bool {
			return bytes.Compare(topicPosition[i][:], topicPosition[j][:]) < 0
		})
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
			lp.lggr.Warnw("Replay requested", "from", fromBlock)
			start = fromBlock
		case <-tick:
			tick = time.After(lp.pollPeriod)
			if start == 0 {
				lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(lp.ctx))
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						lp.lggr.Warnw("unable to get starting block", "err", err)
						continue
					}
					// Otherwise this is the first poll _ever_ on a new chain.
					// Only safe thing to do is to start at finality depth behind tip.
					latest, err := lp.ec.BlockByNumber(context.Background(), nil)
					if err != nil {
						lp.lggr.Warnw("unable to get latest for first poll", "err", err)
						continue
					}
					// Do not support polling chains with don't even have finality depth worth of blocks.
					// Could conceivably support this but not worth the effort.
					if int64(latest.NumberU64()) < lp.finalityDepth {
						lp.lggr.Criticalw("insufficient number of blocks on chain, log poller exiting", "err", err, "latest", latest.NumberU64())
						return
					}
					start = int64(latest.NumberU64()) - lp.finalityDepth + 1
					continue
				}
				start = lastProcessed.BlockNumber + 1
				continue
			}
			start = lp.PollAndSaveLogs(lp.ctx, start)
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
			EventSig:    l.Topics[0].Bytes(), // First topic is always event signature.
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
		from := i
		to := min(i+lp.backfillBatchSize-1, end)
		utils.RetryWithBackoff(ctx, func() bool {
			logs, err = lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
				FromBlock: big.NewInt(from),
				ToBlock:   big.NewInt(to),
				Addresses: lp.filterAddresses(),
				Topics:    lp.filterTopics(),
			})
			if err != nil {
				lp.lggr.Warnw("Unable query for logs, retrying", "err", err, "from", from, "to", to)
				return true
			}
			return false
		})
		if len(logs) == 0 {
			continue
		}
		lp.lggr.Infow("Backfill found logs", "from", from, "to", to, "logs", len(logs))
		utils.RetryWithBackoff(ctx, func() bool {
			if err := lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs)); err != nil {
				lp.lggr.Warnw("Unable to insert logs logs, retrying", "err", err, "from", from, "to", to)
				return true
			}
			return false
		})
	}
	return end + 1
}

// PollAndSaveLogs On startup/crash current is the first block after the last processed block.
func (lp *LogPoller) PollAndSaveLogs(ctx context.Context, current int64) int64 {
	lp.lggr.Infow("Polling for logs", "current", current)
	// Get latest block on chain
	latestBlock, err1 := lp.ec.BlockByNumber(ctx, nil)
	if err1 != nil {
		lp.lggr.Warnw("Unable to get latest block", "err", err1, "current", current)
		return current
	}
	latest := latestBlock.Number().Int64()
	// E.g. 1<-2<-3(current)<-4<-5<-6<-7(latest), finality is 2. So 3,4,5 can be batched.
	// start = current = 3, end = latest - finality = 7-2 = 5 (inclusive range).
	if (latest - current) >= lp.finalityDepth {
		lp.lggr.Infow("Backfilling logs", "start", current, "end", latest-lp.finalityDepth)
		current = lp.backfill(ctx, current, latest-lp.finalityDepth)
	}

	for current <= latest {
		block, err2 := lp.ec.BlockByNumber(ctx, big.NewInt(current))
		if err2 != nil {
			lp.lggr.Warnw("Unable to get block", "err", err2, "current", current)
			return current
		}
		// Does this block point to the same parent that we have saved?
		// If not, there was a reorg, so we need to rewind.
		expectedParent, err2 := lp.orm.SelectBlockByNumber(current - 1)
		if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
			// If err is not a no rows error, assume transient db issue and retry
			lp.lggr.Warnw("Unable to read latest block saved", "err", err2, "current", current)
			return current
		}
		// We will not have the previous block on initial poll or after a backfill.
		havePreviousBlock := !errors.Is(err2, sql.ErrNoRows)
		if havePreviousBlock && !bytes.Equal(block.ParentHash().Bytes(), expectedParent.BlockHash.Bytes()) {
			// There can be another reorg while we're finding the LCA.
			// That is ok, since we'll detect it on the next iteration.
			// Since we go block by block for unfinalized logs, the mismatch starts at current block - 1.
			lca, err3 := lp.findLCA(block.ParentHash(), block.Number().Int64()-1)
			if err3 != nil {
				lp.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err3)
				return current
			}

			lp.lggr.Infow("Re-org detected", "lca", lca, "current", current)
			// We truncate all the blocks and logs after the LCA.
			// We could preserve the logs for forensics, since its possible
			// that applications see them and take action upon it, however that
			// results in significantly slower reads since we must then compute
			// the canonical set per read. Typically if an application took action on a log
			// it would be saved elsewhere e.g. eth_txes, so it seems better to just support the fast reads.
			// Its also nicely analogous to reading from the chain itself.
			err3 = lp.orm.q.Transaction(func(tx pg.Queryer) error {
				// These deletes are bounded by reorg depth, so they are
				// fast and should not slow down the log readers.
				err3 = lp.orm.DeleteRangeBlocks(lca+1, latest, pg.WithQueryer(tx))
				if err3 != nil {
					lp.lggr.Warnw("Unable to clear reorged blocks, retrying", "err", err3)
					return err3
				}
				err3 = lp.orm.DeleteLogs(lca+1, latest, pg.WithQueryer(tx))
				if err3 != nil {
					lp.lggr.Warnw("Unable to clear reorged logs, retrying", "err", err3)
					return err3
				}
				return nil
			})
			if err3 != nil {
				// If we crash or fail to update state we simply do not increment current so we'll detect the same
				// reorg (if still present) and retry.
				return current
			}
			current = lca + 1
			continue
		}

		h := block.Hash()
		logs, err2 := lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
			BlockHash: &h,
			Addresses: lp.filterAddresses(),
			Topics:    lp.filterTopics(),
		})
		if err2 != nil {
			lp.lggr.Warnw("Unable query for logs, retrying", "err", err2, "block", block.Number())
			return current
		}
		lp.lggr.Infow("Unfinalized log query", "logs", len(logs), "current", current)
		err2 = lp.orm.q.Transaction(func(q pg.Queryer) error {
			if err3 := lp.orm.InsertBlock(block.Hash(), block.Number().Int64()); err3 != nil {
				return err3
			}
			if len(logs) == 0 {
				return nil
			}
			return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs))
		})
		if err2 != nil {
			// If we're unable to insert, don't increment current and just retry
			lp.lggr.Warnw("Unable to save logs, retrying", "err", err2, "block", block.Number())
			return current
		}
		current++
	}
	return current
}

func (lp *LogPoller) findLCA(h common.Hash, mismatchStart int64) (int64, error) {
	// Find the first place where our chain and their chain have the same block,
	// that block number is the LCA.
	block, err := lp.ec.BlockByHash(context.Background(), h)
	if err != nil {
		return 0, err
	}
	number := block.Number().Int64()
	if (mismatchStart - number) > lp.finalityDepth {
		lp.lggr.Criticalw("Reorg greater than finality depth detected", "depth", (mismatchStart - number), "finality", lp.finalityDepth)
		return 0, errors.New("reorg greater than finality depth")
	}
	ourBlockHash, err := lp.orm.SelectBlockByNumber(block.Number().Int64())
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(block.Hash().Bytes(), ourBlockHash.BlockHash.Bytes()) {
		return lp.findLCA(block.ParentHash(), mismatchStart)
	}
	// If we do have the blockhash, that is the LCA
	return block.Number().Int64(), nil
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *LogPoller) Logs(start, end int64, eventSig common.Hash, address common.Address) ([]Log, error) {
	return lp.orm.SelectLogsByBlockRangeFilter(start, end, address, eventSig[:])
}
