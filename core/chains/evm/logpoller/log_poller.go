package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
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
	utils.StartStopOnce
	ec                client.Client
	orm               *ORM
	lggr              logger.Logger
	pollPeriod        time.Duration // poll period set by block production rate
	finalityDepth     int64         // finality depth is taken to mean that block (head - finality) is finalized
	backfillBatchSize int64         // batch size to use when backfilling finalized logs

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
// Blocks until the replay starts.
func (lp *LogPoller) Replay(ctx context.Context, fromBlock int64) error {
	latest, err := lp.ec.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if fromBlock < 1 || uint64(fromBlock) > latest.NumberU64() {
		return errors.Errorf("Invalid replay block number %v, acceptable range [1, %v]", fromBlock, latest.NumberU64())
	}
	lp.replay <- fromBlock
	return nil
}

func (lp *LogPoller) Start(parentCtx context.Context) error {
	return lp.StartOnce("LogPoller", func() error {
		ctx, cancel := context.WithCancel(parentCtx)
		lp.ctx = ctx
		lp.cancel = cancel
		go lp.run()
		return nil
	})
}

func (lp *LogPoller) Close() error {
	return lp.StopOnce("LogPoller", func() error {
		lp.cancel()
		<-lp.done
		return nil
	})
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
			tick = time.After(utils.WithJitter(lp.pollPeriod))
			if start != 0 {
				start = lp.pollAndSaveLogs(lp.ctx, start)
				continue
			}
			// Otherwise, still need initial start
			lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(lp.ctx))
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					lp.lggr.Errorw("unable to get starting block", "err", err)
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
					lp.lggr.Warnw("insufficient number of blocks on chain, waiting for finality depth", "err", err, "latest", latest.NumberU64())
					continue
				}
				start = int64(latest.NumberU64()) - lp.finalityDepth + 1
			} else {
				start = lastProcessed.BlockNumber + 1
			}
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
	for from := start; from <= end; from += lp.backfillBatchSize {
		var (
			logs []types.Log
			err  error
		)
		to := min(from+lp.backfillBatchSize-1, end)
		// Retry forever to query for logs,
		// unblocked by resolving node connectivity issues.
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
		// Retry forever to save logs,
		// unblocked by resolving db connectivity issues.
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

func (lp *LogPoller) maybeHandleReorg(ctx context.Context, currentBlockNumber, latestBlockNumber int64) (*types.Block, bool, int64, error) {
	currentBlock, err1 := lp.ec.BlockByNumber(ctx, big.NewInt(currentBlockNumber))
	if err1 != nil {
		lp.lggr.Warnw("Unable to get currentBlock", "err", err1, "currentBlockNumber", currentBlockNumber)
		return nil, false, currentBlockNumber, err1
	}
	// Does this currentBlock point to the same parent that we have saved?
	// If not, there was a reorg, so we need to rewind.
	expectedParent, err1 := lp.orm.SelectBlockByNumber(currentBlockNumber - 1)
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		// If err is not a no rows error, assume transient db issue and retry
		lp.lggr.Warnw("Unable to read latestBlockNumber currentBlock saved", "err", err1, "currentBlockNumber", currentBlockNumber)
		return nil, false, 0, errors.New("Unable to read latestBlockNumber currentBlock saved")
	}
	// We will not have the previous currentBlock on initial poll or after a backfill.
	havePreviousBlock := !errors.Is(err1, sql.ErrNoRows)
	if havePreviousBlock && (currentBlock.ParentHash() != expectedParent.BlockHash) {
		// There can be another reorg while we're finding the LCA.
		// That is ok, since we'll detect it on the next iteration.
		// Since we go currentBlock by currentBlock for unfinalized logs, the mismatch starts at currentBlockNumber currentBlock - 1.
		lca, err2 := lp.findLCA(currentBlock.ParentHash())
		if err2 != nil {
			lp.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err2)
			return nil, false, 0, errors.New("Unable to find LCA after reorg, retrying")
		}

		lp.lggr.Infow("Re-org detected", "lca", lca, "currentBlockNumber", currentBlockNumber, "latestBlockNumber", latestBlockNumber)
		// We truncate all the blocks and logs after the LCA.
		// We could preserve the logs for forensics, since its possible
		// that applications see them and take action upon it, however that
		// results in significantly slower reads since we must then compute
		// the canonical set per read. Typically if an application took action on a log
		// it would be saved elsewhere e.g. eth_txes, so it seems better to just support the fast reads.
		// Its also nicely analogous to reading from the chain itself.
		err2 = lp.orm.q.Transaction(func(tx pg.Queryer) error {
			// These deletes are bounded by reorg depth, so they are
			// fast and should not slow down the log readers.
			err2 = lp.orm.DeleteRangeBlocks(lca+1, latestBlockNumber, pg.WithQueryer(tx))
			if err2 != nil {
				lp.lggr.Warnw("Unable to clear reorged blocks, retrying", "err", err2)
				return err2
			}
			err2 = lp.orm.DeleteLogs(lca+1, latestBlockNumber, pg.WithQueryer(tx))
			if err2 != nil {
				lp.lggr.Warnw("Unable to clear reorged logs, retrying", "err", err2)
				return err2
			}
			return nil
		})
		if err2 != nil {
			// If we crash or fail to update state we simply do not increment currentBlockNumber so we'll detect the same
			// reorg (if still present) and retry.
			return nil, false, 0, err2
		}
		return currentBlock, true, lca + 1, nil
	}
	return currentBlock, false, 0, nil
}

// pollAndSaveLogs On startup/crash current is the first block after the last processed block.
func (lp *LogPoller) pollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
	lp.lggr.Infow("Polling for logs", "currentBlockNumber", currentBlockNumber)
	// Get latestBlockNumber block on chain
	latestBlock, err1 := lp.ec.BlockByNumber(ctx, nil)
	if err1 != nil {
		lp.lggr.Warnw("Unable to get latestBlockNumber block", "err", err1, "currentBlockNumber", currentBlockNumber)
		return currentBlockNumber
	}
	latestBlockNumber := latestBlock.Number().Int64()
	if currentBlockNumber > latestBlockNumber {
		lp.lggr.Debugw("No new blocks since last poll", "currentBlockNumber", currentBlockNumber, "latestBlockNumber", currentBlockNumber)
		return currentBlockNumber
	}
	// Possibly handle a reorg
	_, reorgDetected, newPollBlockNumber, err1 := lp.maybeHandleReorg(ctx, currentBlockNumber, latestBlockNumber)
	if err1 != nil {
		// Continuously retry from same block on any error in reorg handling.
		return currentBlockNumber
	}
	// If we did detect a reorg, we'll have a new block number to start from (LCA+1)
	// so let's resume polling from there.
	if reorgDetected {
		currentBlockNumber = newPollBlockNumber
	}

	// Backfill finalized blocks if we can for performance.
	// E.g. 1<-2<-3(currentBlockNumber)<-4<-5<-6<-7(latestBlockNumber), finality is 2. So 3,4,5 can be batched.
	// start = currentBlockNumber = 3, end = latestBlockNumber - finality = 7-2 = 5 (inclusive range).
	if (latestBlockNumber - currentBlockNumber) >= lp.finalityDepth {
		lp.lggr.Infow("Backfilling logs", "start", currentBlockNumber, "end", latestBlockNumber-lp.finalityDepth)
		currentBlockNumber = lp.backfill(ctx, currentBlockNumber, latestBlockNumber-lp.finalityDepth)
	}

	for currentBlockNumber <= latestBlockNumber {
		// Same reorg detection on unfinalized blocks.
		// Get currentBlockNumber block
		currentBlock, reorgDetected, newPollBlock, err2 := lp.maybeHandleReorg(ctx, currentBlockNumber, latestBlockNumber)
		if err2 != nil {
			return currentBlockNumber
		}
		if reorgDetected {
			currentBlockNumber = newPollBlock
			continue
		}

		h := currentBlock.Hash()
		logs, err2 := lp.ec.FilterLogs(ctx, ethereum.FilterQuery{
			BlockHash: &h,
			Addresses: lp.filterAddresses(),
			Topics:    lp.filterTopics(),
		})
		if err2 != nil {
			lp.lggr.Warnw("Unable query for logs, retrying", "err", err2, "block", currentBlock.Number())
			return currentBlockNumber
		}
		lp.lggr.Infow("Unfinalized log query", "logs", len(logs), "currentBlockNumber", currentBlockNumber)
		err2 = lp.orm.q.Transaction(func(q pg.Queryer) error {
			if err3 := lp.orm.InsertBlock(currentBlock.Hash(), currentBlock.Number().Int64()); err3 != nil {
				return err3
			}
			if len(logs) == 0 {
				return nil
			}
			return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs))
		})
		if err2 != nil {
			// If we're unable to insert, don't increment currentBlockNumber and just retry
			lp.lggr.Warnw("Unable to save logs, retrying", "err", err2, "block", currentBlock.Number())
			return currentBlockNumber
		}
		currentBlockNumber++
	}
	return currentBlockNumber
}

func (lp *LogPoller) findLCA(h common.Hash) (int64, error) {
	// Find the first place where our chain and their chain have the same block,
	// that block number is the LCA.
	block, err := lp.ec.BlockByHash(context.Background(), h)
	if err != nil {
		return 0, err
	}
	blockNumber := block.Number().Int64()
	startBlockNumber := blockNumber
	for blockNumber >= (startBlockNumber - lp.finalityDepth) {
		ourBlockHash, err := lp.orm.SelectBlockByNumber(blockNumber)
		if err != nil {
			return 0, err
		}
		if block.Hash() == ourBlockHash.BlockHash {
			// If we do have the blockhash, that is the LCA
			return blockNumber, nil
		}
		blockNumber--
		block, err = lp.ec.BlockByHash(context.Background(), block.ParentHash())
		if err != nil {
			return 0, err
		}
	}
	lp.lggr.Criticalw("Reorg greater than finality depth detected", "finality", lp.finalityDepth)
	return 0, errors.New("reorg greater than finality depth")
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *LogPoller) Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsByBlockRangeFilter(start, end, address, eventSig[:], qopts...)
}

func (lp *LogPoller) LatestBlock(qopts ...pg.QOpt) (int64, error) {
	b, err := lp.orm.SelectLatestBlock(qopts...)
	if err != nil {
		return 0, err
	}
	return b.BlockNumber, nil
}

// LatestLogByEventSigWithConfs finds the latest log that has confs number of blocks on top of the log.
func (lp *LogPoller) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	log, err := lp.orm.SelectLatestLogEventSigWithConfs(eventSig, address, confs, qopts...)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (lp *LogPoller) LatestLogEventSigsAddrs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.LatestLogEventSigsAddrs(fromBlock, addresses, eventSigs, qopts...)
}

// IndexedLogs finds all the logs that have a topic value in topicValues at index topicIndex.
func (lp *LogPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogs(address, eventSig[:], topicIndex, topicValues, confs, qopts...)
}

// Index is 0 based.
func (lp *LogPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectDataWordGreaterThan(address, eventSig[:], wordIndex, wordValueMin, confs, qopts...)
}

// Index is 0 based.
func (lp *LogPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectDataWordRange(address, eventSig[:], wordIndex, wordValueMin, wordValueMax, confs, qopts...)
}

// IndexedLogs finds all the logs that have a topic value greater than topicValueMin at index topicIndex.
// Only works for integer topics.
func (lp *LogPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexLogsTopicGreaterThan(address, eventSig[:], topicIndex, topicValueMin, confs, qopts...)
}

func (lp *LogPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexLogsTopicRange(address, eventSig[:], topicIndex, topicValueMin, topicValueMax, confs, qopts...)
}

func EvmWord(i uint64) common.Hash {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BytesToHash(b)
}
