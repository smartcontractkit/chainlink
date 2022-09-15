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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/mathutil"
)

//go:generate mockery --name LogPoller --output ./mocks/ --case=underscore --structname LogPoller --filename log_poller.go
type LogPoller interface {
	services.ServiceCtx
	Replay(ctx context.Context, fromBlock int64) error
	MergeFilter(eventSigs []common.Hash, addresses []common.Address) error
	LatestBlock(qopts ...pg.QOpt) (int64, error)
	GetBlocks(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error)

	// General querying
	Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error)
	LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error)
	LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error)
	LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]Log, error)

	// Content based querying
	IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
	LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error)
}

var (
	_                          LogPoller = &logPoller{}
	ErrReplayAbortedByClient             = errors.New("replay aborted by client")
	ErrReplayAbortedOnShutdown           = errors.New("replay aborted, log poller shutdown")
)

type logPoller struct {
	utils.StartStopOnce
	ec                client.Client
	orm               *ORM
	lggr              logger.Logger
	pollPeriod        time.Duration // poll period set by block production rate
	finalityDepth     int64         // finality depth is taken to mean that block (head - finality) is finalized
	backfillBatchSize int64         // batch size to use when backfilling finalized logs
	rpcBatchSize      int64         // batch size to use for fallback RPC calls made in GetBlocks

	filterMu  sync.RWMutex
	addresses map[common.Address]struct{}
	eventSigs map[common.Hash]struct{}

	replayStart    chan ReplayRequest
	replayComplete chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
	done           chan struct{}
}

type ReplayRequest struct {
	fromBlock int64
	ctx       context.Context
}

// NewLogPoller creates a log poller. Note there is an assumption
// that blocks can be processed faster than they are produced for the given chain, or the poller will fall behind.
// Block processing involves in the steady state (non-reorg case):
// - eth_getBlockByNumber - headers only (transaction hashes, not full transaction objects),
// - 1 db read latest block
// - 1 db tx including block write and logs write to logs.
// How fast that can be done depends largely on network speed and DB, but even for the fastest
// support chain, polygon, which has 2s block times, we need RPCs roughly with <= 500ms latency
func NewLogPoller(orm *ORM, ec client.Client, lggr logger.Logger, pollPeriod time.Duration, finalityDepth, backfillBatchSize, rpcBatchSize int64) *logPoller {
	return &logPoller{
		ec:                ec,
		orm:               orm,
		lggr:              lggr,
		replayStart:       make(chan ReplayRequest),
		replayComplete:    make(chan struct{}),
		done:              make(chan struct{}),
		pollPeriod:        pollPeriod,
		finalityDepth:     finalityDepth,
		backfillBatchSize: backfillBatchSize,
		rpcBatchSize:      rpcBatchSize,
		addresses:         make(map[common.Address]struct{}),
		eventSigs:         make(map[common.Hash]struct{}),
	}
}

// MergeFilter adds the provided eventSigs and addresses to the log poller's log filter query.
// If an event matching any of the given event signatures is emitted from any of the provided addresses,
// the log poller will pick those up and save them. For topic specific queries see content based querying.
// Clients may choose to MergeFilter and then Replay in order to ensure desired logs are present.
// NOTE: due to constraints of the eth filter, there is "leakage" between successive MergeFilter calls, for example
// MergeFilter(event1, addr1)
// MergeFilter(event2, addr2)
// will result in the poller saving (event1, addr2) or (event2, addr1) as well, should it exist.
// Generally speaking this is harmless. We enforce that eventSigs and addresses are non-empty,
// which means that anonymous events are not supported and log.Topics >= 1 always (log.Topics[0] is the event signature).
func (lp *logPoller) MergeFilter(eventSigs []common.Hash, addresses []common.Address) error {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	// Add any valid unique eventSigs or addresses.
	for _, eventSig := range eventSigs {
		// Force specification of both address and eventSig to avoid
		// events from unknown addresses.
		if eventSig == [common.HashLength]byte{} {
			return errors.Errorf("empty event sig")
		}
		lp.eventSigs[eventSig] = struct{}{}
	}
	for _, addr := range addresses {
		if addr == [common.AddressLength]byte{} {
			return errors.Errorf("empty address")
		}
		lp.addresses[addr] = struct{}{}
	}
	return nil
}

func (lp *logPoller) filter(from, to *big.Int, bh *common.Hash) ethereum.FilterQuery {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	var (
		addresses []common.Address
		eventSigs []common.Hash
	)
	for addr := range lp.addresses {
		addresses = append(addresses, addr)
	}
	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i][:], addresses[j][:]) < 0
	})
	for eventSig := range lp.eventSigs {
		eventSigs = append(eventSigs, eventSig)
	}
	sort.Slice(eventSigs, func(i, j int) bool {
		return bytes.Compare(eventSigs[i][:], eventSigs[j][:]) < 0
	})
	if len(eventSigs) == 0 && len(addresses) == 0 {
		// If no filter specified, ignore everything.
		// This allows us to keep the log poller up and running with no filters present (e.g. no jobs on the node),
		// then as jobs are added dynamically start using their filters.
		addresses = []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000000")}
	}
	return ethereum.FilterQuery{FromBlock: from, ToBlock: to, BlockHash: bh, Topics: [][]common.Hash{eventSigs}, Addresses: addresses}
}

// Replay signals that the poller should resume from a new block.
// Blocks until the replay is complete.
// Replay can be used to ensure that filter modification has been applied for all blocks from "fromBlock" up to latest.
func (lp *logPoller) Replay(ctx context.Context, fromBlock int64) error {
	latest, err := lp.ec.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if fromBlock < 1 || fromBlock > latest.Number.Int64() {
		return errors.Errorf("Invalid replay block number %v, acceptable range [1, %v]", fromBlock, latest)
	}
	// Block until replay notification accepted or cancelled.
	select {
	case lp.replayStart <- ReplayRequest{fromBlock, ctx}:
	case <-lp.ctx.Done():
		return ErrReplayAbortedOnShutdown
	case <-ctx.Done():
		return ErrReplayAbortedByClient
	}
	// Block until replay complete or cancelled.
	select {
	case <-lp.replayComplete:
	case <-lp.ctx.Done():
		return ErrReplayAbortedOnShutdown
	case <-ctx.Done():
		return ErrReplayAbortedByClient
	}
	return nil
}

func (lp *logPoller) Start(parentCtx context.Context) error {
	return lp.StartOnce("LogPoller", func() error {
		ctx, cancel := context.WithCancel(parentCtx)
		lp.ctx = ctx
		lp.cancel = cancel
		go lp.run()
		return nil
	})
}

func (lp *logPoller) Close() error {
	return lp.StopOnce("LogPoller", func() error {
		lp.cancel()
		<-lp.done
		return nil
	})
}

func (lp *logPoller) run() {
	defer close(lp.done)
	tick := time.After(0)
	for {
		select {
		case <-lp.ctx.Done():
			return
		case replayReq := <-lp.replayStart:
			lp.lggr.Warnw("Executing replay", "fromBlock", replayReq.fromBlock)
			// Serially process replay requests.
			lp.pollAndSaveLogs(replayReq.ctx, replayReq.fromBlock)
			select {
			case <-lp.ctx.Done():
				// We're shutting down, lets return.
				return
			case <-replayReq.ctx.Done():
				// Client gave up, lets continue.
				continue
			case lp.replayComplete <- struct{}{}:
			}
		case <-tick:
			tick = time.After(utils.WithJitter(lp.pollPeriod))
			// Always start from the latest block in the db.
			var start int64
			lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(lp.ctx))
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					// Assume transient db reading issue, retry forever.
					lp.lggr.Errorw("unable to get starting block", "err", err)
					continue
				}
				// Otherwise this is the first poll _ever_ on a new chain.
				// Only safe thing to do is to start at the first finalized block.
				latest, err := lp.ec.HeaderByNumber(lp.ctx, nil)
				if err != nil {
					lp.lggr.Warnw("unable to get latest for first poll", "err", err)
					continue
				}
				latestNum := latest.Number.Int64()
				// Do not support polling chains with don't even have finality depth worth of blocks.
				// Could conceivably support this but not worth the effort.
				// Need finality depth + 1, no block 0.
				if latestNum <= lp.finalityDepth {
					lp.lggr.Warnw("insufficient number of blocks on chain, waiting for finality depth", "err", err, "latest", latestNum, "finality", lp.finalityDepth)
					continue
				}
				// Starting at the first finalized block. We do not backfill the first finalized block.
				start = latestNum - lp.finalityDepth
			} else {
				start = lastProcessed.BlockNumber + 1
			}
			lp.pollAndSaveLogs(lp.ctx, start)
		}
	}
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

func (lp *logPoller) backfill(ctx context.Context, start, end int64) int64 {
	for from := start; from <= end; from += lp.backfillBatchSize {
		var (
			logs []types.Log
			err  error
		)
		to := mathutil.Min(from+lp.backfillBatchSize-1, end)
		// Retry forever to query for logs,
		// unblocked by resolving node connectivity issues.
		utils.RetryWithBackoff(ctx, func() bool {
			logs, err = lp.ec.FilterLogs(ctx, lp.filter(big.NewInt(from), big.NewInt(to), nil))
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
			err = lp.orm.q.WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
				return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs), pg.WithQueryer(tx))
			})
			if err != nil {
				lp.lggr.Warnw("Unable to insert logs logs, retrying", "err", err, "from", from, "to", to)
				return true
			}
			return false
		})
	}
	return end + 1
}

// getCurrentBlockMaybeHandleReorg accepts a block number
// and will return that block if its parent points to our last saved block.
// If its parent does not point to our last saved block we know a reorg has occurred.
// In that case return the LCA+1, i.e. our new current (unprocessed) block.
func (lp *logPoller) getCurrentBlockMaybeHandleReorg(ctx context.Context, currentBlockNumber int64) (*types.Header, error) {
	currentBlock, err1 := lp.ec.HeaderByNumber(ctx, big.NewInt(currentBlockNumber))
	if err1 != nil {
		lp.lggr.Warnw("Unable to get currentBlock", "err", err1, "currentBlockNumber", currentBlockNumber)
		return nil, err1
	}
	// Does this currentBlock point to the same parent that we have saved?
	// If not, there was a reorg, so we need to rewind.
	expectedParent, err1 := lp.orm.SelectBlockByNumber(currentBlockNumber-1, pg.WithParentCtx(ctx))
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		// If err is not a 'no rows' error, assume transient db issue and retry
		lp.lggr.Warnw("Unable to read latestBlockNumber currentBlock saved", "err", err1, "currentBlockNumber", currentBlockNumber)
		return nil, errors.New("Unable to read latestBlockNumber currentBlock saved")
	}
	// We will not have the previous currentBlock on initial poll.
	havePreviousBlock := err1 == nil
	if !havePreviousBlock {
		lp.lggr.Infow("Do not have previous block, first poll ever on new chain or after backfill")
		return currentBlock, nil
	}
	// Check for reorg.
	if currentBlock.ParentHash != expectedParent.BlockHash {
		// There can be another reorg while we're finding the LCA.
		// That is ok, since we'll detect it on the next iteration.
		// Since we go currentBlock by currentBlock for unfinalized logs, the mismatch starts at currentBlockNumber currentBlock - 1.
		blockAfterLCA, err2 := lp.findBlockAfterLCA(ctx, currentBlock)
		if err2 != nil {
			lp.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err2)
			return nil, errors.New("Unable to find LCA after reorg, retrying")
		}

		lp.lggr.Infow("Reorg detected", "blockAfterLCA", blockAfterLCA.Number.Int64(), "currentBlockNumber", currentBlockNumber)
		// We truncate all the blocks and logs after the LCA.
		// We could preserve the logs for forensics, since its possible
		// that applications see them and take action upon it, however that
		// results in significantly slower reads since we must then compute
		// the canonical set per read. Typically, if an application took action on a log
		// it would be saved elsewhere e.g. eth_txes, so it seems better to just support the fast reads.
		// Its also nicely analogous to reading from the chain itself.
		err2 = lp.orm.q.WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			// These deletes are bounded by reorg depth, so they are
			// fast and should not slow down the log readers.
			err3 := lp.orm.DeleteRangeBlocks(blockAfterLCA.Number.Int64(), currentBlockNumber, pg.WithQueryer(tx))
			if err3 != nil {
				lp.lggr.Warnw("Unable to clear reorged blocks, retrying", "err", err3)
				return err3
			}
			err3 = lp.orm.DeleteLogs(blockAfterLCA.Number.Int64(), currentBlockNumber, pg.WithQueryer(tx))
			if err3 != nil {
				lp.lggr.Warnw("Unable to clear reorged logs, retrying", "err", err3)
				return err3
			}
			return nil
		})
		if err2 != nil {
			// If we error on db commit, we can't know if the tx went through or not.
			// We return an error here which will cause us to restart polling from lastBlockSaved + 1
			return nil, err2
		}
		return blockAfterLCA, nil
	}
	// No reorg, return current block.
	return currentBlock, nil
}

// pollAndSaveLogs On startup/crash current is the first block after the last processed block.
// currentBlockNumber is the block from where new logs are to be polled & saved. Under normal
// conditions this would be equal to lastProcessed.BlockNumber + 1.
func (lp *logPoller) pollAndSaveLogs(ctx context.Context, currentBlockNumber int64) {
	lp.lggr.Infow("Polling for logs", "currentBlockNumber", currentBlockNumber)
	latestBlock, err := lp.ec.HeaderByNumber(ctx, nil)
	if err != nil {
		lp.lggr.Warnw("Unable to get latestBlockNumber block", "err", err, "currentBlockNumber", currentBlockNumber)
		return
	}
	latestBlockNumber := latestBlock.Number.Int64()
	if currentBlockNumber > latestBlockNumber {
		// Note there can also be a reorg "shortening" i.e. chain height decreases but TDD increases. In that case
		// we also just wait until the new tip is longer and then detect the reorg.
		lp.lggr.Debugw("No new blocks since last poll", "currentBlockNumber", currentBlockNumber, "latestBlockNumber", latestBlockNumber)
		return
	}
	// Possibly handle a reorg. For example if we crash, we'll be in the middle of processing unfinalized blocks.
	// Returns (currentBlock || LCA+1 if reorg detected, error)
	currentBlock, err := lp.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber)
	if err != nil {
		// If there's an error handling the reorg, we can't be sure what state the db was left in.
		// Resume from the latest block saved.
		lp.lggr.Errorw("Unable to get current block", "err", err)
		return
	}
	currentBlockNumber = currentBlock.Number.Int64()
	// Backfill finalized blocks if we can for performance. If we crash during backfill, we may reprocess logs.
	// Log insertion is idempotent so this is ok.
	// E.g. 1<-2<-3(currentBlockNumber)<-4<-5<-6<-7(latestBlockNumber), finality is 2. So 3,4 can be batched.
	// Although 5 is finalized, we still need to save it to the db for reorg detection if 6 is a reorg.
	// start = currentBlockNumber = 3, end = latestBlockNumber - finality - 1 = 7-2-1 = 4 (inclusive range).
	lastSafeBackfillBlock := latestBlockNumber - lp.finalityDepth - 1
	if lastSafeBackfillBlock >= currentBlockNumber {
		lp.lggr.Infow("Backfilling logs", "start", currentBlockNumber, "end", lastSafeBackfillBlock)
		currentBlockNumber = lp.backfill(ctx, currentBlockNumber, lastSafeBackfillBlock)
	}

	for currentBlockNumber <= latestBlockNumber {
		// Same reorg detection on unfinalized blocks.
		currentBlock, err = lp.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber)
		if err != nil {
			// If there's an error handling the reorg, we can't be sure what state the db was left in.
			// Resume from the latest block saved.
			lp.lggr.Errorw("Unable to get current block", "err", err)
			return
		}
		currentBlockNumber = currentBlock.Number.Int64()
		h := currentBlock.Hash()
		var logs []types.Log
		logs, err = lp.ec.FilterLogs(ctx, lp.filter(nil, nil, &h))
		if err != nil {
			lp.lggr.Warnw("Unable to query for logs, retrying", "err", err, "block", currentBlockNumber)
			return
		}
		lp.lggr.Infow("Unfinalized log query", "logs", len(logs), "currentBlockNumber", currentBlockNumber, "blockHash", currentBlock.Hash())
		err = lp.orm.q.WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			if err2 := lp.orm.InsertBlock(currentBlock.Hash(), currentBlockNumber, pg.WithQueryer(tx)); err2 != nil {
				return err2
			}
			if len(logs) == 0 {
				return nil
			}
			return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs), pg.WithQueryer(tx))
		})
		if err != nil {
			lp.lggr.Warnw("Unable to save logs resuming from last saved block + 1", "err", err, "block", currentBlockNumber)
			return
		}
		currentBlockNumber++
	}
}

// Find the first place where our chain and their chain have the same block,
// that block number is the LCA. Return the block after that, where we want to resume polling.
func (lp *logPoller) findBlockAfterLCA(ctx context.Context, current *types.Header) (*types.Header, error) {
	// Current is where the mismatch starts.
	// Check its parent to see if its the same as ours saved.
	parent, err := lp.ec.HeaderByHash(ctx, current.ParentHash)
	if err != nil {
		return nil, err
	}
	blockAfterLCA := *current
	reorgStart := parent.Number.Int64()
	// We expected reorgs up to the block after (current - finalityDepth),
	// since the block at (current - finalityDepth) is finalized.
	// We loop via parent instead of current so current always holds the LCA+1.
	// If the parent block number becomes < the first finalized block our reorg is too deep.
	for parent.Number.Int64() >= (reorgStart - lp.finalityDepth) {
		ourParentBlockHash, err := lp.orm.SelectBlockByNumber(parent.Number.Int64(), pg.WithParentCtx(ctx))
		if err != nil {
			return nil, err
		}
		if parent.Hash() == ourParentBlockHash.BlockHash {
			// If we do have the blockhash, return blockAfterLCA
			return &blockAfterLCA, nil
		}
		// Otherwise get a new parent and update blockAfterLCA.
		blockAfterLCA = *parent
		parent, err = lp.ec.HeaderByHash(ctx, parent.ParentHash)
		if err != nil {
			return nil, err
		}
	}
	lp.lggr.Criticalw("Reorg greater than finality depth detected", "max reorg depth", lp.finalityDepth-1)
	return nil, errors.New("Reorg greater than finality depth")
}

func (lp *logPoller) assertInFilter(eventSigs []common.Hash, addresses []common.Address) error {
	lp.filterMu.RLock()
	defer lp.filterMu.RUnlock()
	for _, eventSig := range eventSigs {
		if _, ok := lp.eventSigs[eventSig]; !ok {
			return errors.Errorf("eventSig %x not registered", eventSig)
		}
	}
	for _, addr := range addresses {
		if _, ok := lp.addresses[addr]; !ok {
			return errors.Errorf("address %x not registered", addr)
		}
	}
	return nil
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *logPoller) Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectLogsByBlockRangeFilter(start, end, address, eventSig[:], qopts...)
}

func (lp *logPoller) LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter(eventSigs, []common.Address{address}); err != nil {
		return nil, err
	}
	sigs := make([][]byte, 0, len(eventSigs))
	for _, sig := range eventSigs {
		sigs = append(sigs, sig.Bytes())
	}
	return lp.orm.SelectLogsWithSigsByBlockRangeFilter(start, end, address, sigs, qopts...)
}

// IndexedLogs finds all the logs that have a topic value in topicValues at index topicIndex.
func (lp *logPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectIndexedLogs(address, eventSig[:], topicIndex, topicValues, confs, qopts...)
}

// LogsDataWordGreaterThan note index is 0 based.
func (lp *logPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectDataWordGreaterThan(address, eventSig[:], wordIndex, wordValueMin, confs, qopts...)
}

// LogsDataWordRange note index is 0 based.
func (lp *logPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectDataWordRange(address, eventSig[:], wordIndex, wordValueMin, wordValueMax, confs, qopts...)
}

// IndexedLogsTopicGreaterThan finds all the logs that have a topic value greater than topicValueMin at index topicIndex.
// Only works for integer topics.
func (lp *logPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectIndexLogsTopicGreaterThan(address, eventSig[:], topicIndex, topicValueMin, confs, qopts...)
}

func (lp *logPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectIndexLogsTopicRange(address, eventSig[:], topicIndex, topicValueMin, topicValueMax, confs, qopts...)
}

// LatestBlock returns the latest block the log poller is on. It tracks blocks to be able
// to detect reorgs.
func (lp *logPoller) LatestBlock(qopts ...pg.QOpt) (int64, error) {
	b, err := lp.orm.SelectLatestBlock(qopts...)
	if err != nil {
		return 0, err
	}

	return b.BlockNumber, nil
}

func (lp *logPoller) BlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	return lp.orm.SelectBlockByNumber(n, qopts...)
}

// LatestLogByEventSigWithConfs finds the latest log that has confs number of blocks on top of the log.
func (lp *logPoller) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	if err := lp.assertInFilter([]common.Hash{eventSig}, []common.Address{address}); err != nil {
		return nil, err
	}
	return lp.orm.SelectLatestLogEventSigWithConfs(eventSig, address, confs, qopts...)
}

func (lp *logPoller) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]Log, error) {
	if err := lp.assertInFilter(eventSigs, addresses); err != nil {
		return nil, err
	}
	return lp.orm.SelectLatestLogEventSigsAddrsWithConfs(fromBlock, addresses, eventSigs, confs, qopts...)
}

// GetBlocks tries to get the specified block numbers from the log pollers
// blocks table. Returns the blocks it was able to find, empty slice if none.
// When the log poller does not have requested blocks, it falls back
// to RPC to fetch the missing blocks.
// response contains blocks in the same order as "numbers" in request parameters
// the first context parameter takes precedence over contexts passed through qopts
func (lp *logPoller) GetBlocks(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	blocksFound := make(map[uint64]LogPollerBlock)
	qopts = append(qopts, pg.WithParentCtx(ctx))
	lpBlocks, err := lp.orm.GetBlocks(numbers, qopts...)
	if err != nil {
		lp.lggr.Warnw("Error while retrieving blocks from log pollers blocks table. Falling back to RPC...", "requestedBlocks", numbers, "err", err)
	}
	for _, b := range lpBlocks {
		blocksFound[uint64(b.BlockNumber)] = b
	}

	// Fallback to RPC for blocks not found in log poller blocks table
	var reqs []rpc.BatchElem
	for _, num := range numbers {
		if _, ok := blocksFound[num]; !ok {
			req := rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{hexutil.EncodeBig(big.NewInt(0).SetUint64(num)), false},
				Result: &evmtypes.Head{},
			}
			reqs = append(reqs, req)
		}
	}

	for i := 0; i < len(reqs); i += int(lp.rpcBatchSize) {
		j := i + int(lp.rpcBatchSize)
		if j > len(reqs) {
			j = len(reqs)
		}

		err := lp.ec.BatchCallContext(ctx, reqs[i:j])

		if err != nil {
			return nil, err
		}
	}

	for _, r := range reqs {
		if r.Error != nil {
			return nil, r.Error
		}
		block, is := r.Result.(*evmtypes.Head)

		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", &evmtypes.Head{}, r.Result)
		}
		if block == nil {
			return nil, errors.New("invariant violation: got nil block")
		}
		if block.Hash == (common.Hash{}) {
			return nil, errors.Errorf("missing block hash for block number: %d", block.Number)
		}
		if block.Number < 0 {
			return nil, errors.Errorf("expected block number to be >= to 0, got %d", block.Number)
		}
		blocksFound[uint64(block.Number)] = LogPollerBlock{
			EvmChainId:  block.EVMChainID,
			BlockHash:   block.Hash,
			BlockNumber: block.Number,
			CreatedAt:   block.Timestamp,
		}
	}

	var blocks []LogPollerBlock
	for _, num := range numbers {
		b, ok := blocksFound[num]
		if !ok {
			return nil, errors.Errorf("block: %d was not found in db or RPC call", num)
		}
		blocks = append(blocks, b)
	}

	return blocks, nil
}

func EvmWord(i uint64) common.Hash {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BytesToHash(b)
}
