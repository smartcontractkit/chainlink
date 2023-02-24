package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/mathutil"
)

//go:generate mockery --quiet --name LogPoller --output ./mocks/ --case=underscore --structname LogPoller --filename log_poller.go
type LogPoller interface {
	services.ServiceCtx
	Replay(ctx context.Context, fromBlock int64) error
	RegisterFilter(filter Filter) error
	UnregisterFilter(name string) error
	LatestBlock(qopts ...pg.QOpt) (int64, error)
	GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error)
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

type Client interface {
	HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error)
	HeadByHash(ctx context.Context, n common.Hash) (*evmtypes.Head, error)
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	ChainID() *big.Int
}

var (
	_                          LogPoller = &logPoller{}
	ErrReplayAbortedByClient             = errors.New("replay aborted by client")
	ErrReplayAbortedOnShutdown           = errors.New("replay aborted, log poller shutdown")
)

type logPoller struct {
	utils.StartStopOnce
	ec                    Client
	orm                   *ORM
	lggr                  logger.Logger
	pollPeriod            time.Duration // poll period set by block production rate
	finalityDepth         int64         // finality depth is taken to mean that block (head - finality) is finalized
	keepBlocksDepth       int64         // the number of blocks behind the head for which we keep the blocks. Must be greater than finality depth + 1.
	backfillBatchSize     int64         // batch size to use when backfilling finalized logs
	rpcBatchSize          int64         // batch size to use for fallback RPC calls made in GetBlocks
	backupPollerNextBlock int64

	filterMu        sync.RWMutex
	filters         map[string]Filter
	filterDirty     bool
	cachedAddresses []common.Address
	cachedEventSigs []common.Hash

	replayStart    chan ReplayRequest
	replayComplete chan error
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
// Block processing involves the following calls in steady state (without reorgs):
// - eth_getBlockByNumber - headers only (transaction hashes, not full transaction objects),
// - eth_getLogs - get the logs for the block
// - 1 db read latest block - for checking reorgs
// - 1 db tx including block write and logs write to logs.
// How fast that can be done depends largely on network speed and DB, but even for the fastest
// support chain, polygon, which has 2s block times, we need RPCs roughly with <= 500ms latency
func NewLogPoller(orm *ORM, ec Client, lggr logger.Logger, pollPeriod time.Duration,
	finalityDepth int64, backfillBatchSize int64, rpcBatchSize int64, keepBlocksDepth int64) *logPoller {

	return &logPoller{
		ec:                ec,
		orm:               orm,
		lggr:              lggr,
		replayStart:       make(chan ReplayRequest),
		replayComplete:    make(chan error),
		done:              make(chan struct{}),
		pollPeriod:        pollPeriod,
		finalityDepth:     finalityDepth,
		backfillBatchSize: backfillBatchSize,
		rpcBatchSize:      rpcBatchSize,
		keepBlocksDepth:   keepBlocksDepth,
		filters:           make(map[string]Filter),
		filterDirty:       true, // Always build filter on first call to cache an empty filter if nothing registered yet.
	}
}

type Filter struct {
	Name      string // see FilterName(id, args) below
	EventSigs evmtypes.HashArray
	Addresses evmtypes.AddressArray
}

// FilterName is a suggested convenience function for clients to construct unique filter names
// to populate Name field of struct Filter
func FilterName(id string, args ...any) string {
	if len(args) == 0 {
		return id
	}
	s := &strings.Builder{}
	s.WriteString(id)
	s.WriteString(" - ")
	fmt.Fprintf(s, "%s", args[0])
	for _, a := range args[1:] {
		fmt.Fprintf(s, ":%s", a)
	}
	return s.String()
}

// contains returns true if this filter already fully contains a
// filter passed to it.
func (filter *Filter) contains(other *Filter) bool {
	if other == nil {
		return true
	}
	addresses := make(map[common.Address]interface{})
	for _, addr := range filter.Addresses {
		addresses[addr] = struct{}{}
	}
	events := make(map[common.Hash]interface{})
	for _, ev := range filter.EventSigs {
		events[ev] = struct{}{}
	}

	for _, addr := range other.Addresses {
		if _, ok := addresses[addr]; !ok {
			return false
		}
	}
	for _, ev := range other.EventSigs {
		if _, ok := events[ev]; !ok {
			return false
		}
	}
	return true
}

// RegisterFilter adds the provided EventSigs and Addresses to the log poller's log filter query.
// If any eventSig is emitted from any address, it will be captured by the log poller.
// If an event matching any of the given event signatures is emitted from any of the provided Addresses,
// the log poller will pick those up and save them. For topic specific queries see content based querying.
// Clients may choose to MergeFilter and then Replay in order to ensure desired logs are present.
// NOTE: due to constraints of the eth filter, there is "leakage" between successive MergeFilter calls, for example
// RegisterFilter(event1, addr1)
// RegisterFilter(event2, addr2)
// will result in the poller saving (event1, addr2) or (event2, addr1) as well, should it exist.
// Generally speaking this is harmless. We enforce that EventSigs and Addresses are non-empty,
// which means that anonymous events are not supported and log.Topics >= 1 always (log.Topics[0] is the event signature).
// The filter may be unregistered later by Filter.Name
func (lp *logPoller) RegisterFilter(filter Filter) error {
	if len(filter.Addresses) == 0 {
		return errors.Errorf("at least one address must be specified")
	}
	if len(filter.EventSigs) == 0 {
		return errors.Errorf("at least one event must be specified")
	}

	for _, eventSig := range filter.EventSigs {
		if eventSig == [common.HashLength]byte{} {
			return errors.Errorf("empty event sig")
		}
	}
	for _, addr := range filter.Addresses {
		if addr == [common.AddressLength]byte{} {
			return errors.Errorf("empty address")
		}
	}

	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()

	if existingFilter, ok := lp.filters[filter.Name]; ok {
		if existingFilter.contains(&filter) {
			// Nothing new in this filter
			return nil
		}
		lp.lggr.Warnw("Updating existing filter %s with more events or addresses", "filter.Name", filter.Name)
	} else {
		lp.lggr.Debugf("Creating new filter %s", filter.Name)
	}

	if err := lp.orm.InsertFilter(filter); err != nil {
		return errors.Wrap(err, "RegisterFilter failed to save filter to db")
	}
	lp.filters[filter.Name] = filter
	lp.filterDirty = true
	return nil
}

func (lp *logPoller) UnregisterFilter(name string) error {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()

	_, ok := lp.filters[name]
	if !ok {
		return errors.Errorf("filter %s not found", name)
	}
	if err := lp.orm.DeleteFilter(name); err != nil {
		return errors.Wrapf(err, "Failed to delete filter %s", name)
	}
	delete(lp.filters, name)
	lp.filterDirty = true
	return nil
}

func (lp *logPoller) filter(from, to *big.Int, bh *common.Hash) ethereum.FilterQuery {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()
	if !lp.filterDirty {
		return ethereum.FilterQuery{FromBlock: from, ToBlock: to, BlockHash: bh, Topics: [][]common.Hash{lp.cachedEventSigs}, Addresses: lp.cachedAddresses}
	}
	var (
		addresses  []common.Address
		eventSigs  []common.Hash
		addressMp  = make(map[common.Address]struct{})
		eventSigMp = make(map[common.Hash]struct{})
	)
	// Merge filters.
	for _, filter := range lp.filters {
		for _, addr := range filter.Addresses {
			addressMp[addr] = struct{}{}
		}
		for _, eventSig := range filter.EventSigs {
			eventSigMp[eventSig] = struct{}{}
		}
	}
	for addr := range addressMp {
		addresses = append(addresses, addr)
	}
	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i][:], addresses[j][:]) < 0
	})
	for eventSig := range eventSigMp {
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
		eventSigs = []common.Hash{}
	}
	lp.cachedAddresses = addresses
	lp.cachedEventSigs = eventSigs
	lp.filterDirty = false
	return ethereum.FilterQuery{FromBlock: from, ToBlock: to, BlockHash: bh, Topics: [][]common.Hash{eventSigs}, Addresses: addresses}
}

// Replay signals that the poller should resume from a new block.
// Blocks until the replay is complete.
// Replay can be used to ensure that filter modification has been applied for all blocks from "fromBlock" up to latest.
func (lp *logPoller) Replay(ctx context.Context, fromBlock int64) error {
	latest, err := lp.ec.HeadByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if fromBlock < 1 || fromBlock > latest.Number {
		return errors.Errorf("Invalid replay block number %v, acceptable range [1, %v]", fromBlock, latest.Number)
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
	case err := <-lp.replayComplete:
		return err
	case <-lp.ctx.Done():
		return ErrReplayAbortedOnShutdown
	case <-ctx.Done():
		return ErrReplayAbortedByClient
	}
}

func (lp *logPoller) Start(parentCtx context.Context) error {
	if lp.keepBlocksDepth < (lp.finalityDepth + 1) {
		// We add 1 since for reorg detection on the first unfinalized block
		// we need to keep 1 finalized block.
		return errors.Errorf("keepBlocksDepth %d must be greater than finality %d + 1", lp.keepBlocksDepth, lp.finalityDepth)
	}
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

func (lp *logPoller) Name() string {
	return lp.lggr.Name()
}

func (lp *logPoller) HealthReport() map[string]error {
	return map[string]error{lp.Name(): lp.Healthy()}
}

func (lp *logPoller) getReplayFromBlock(ctx context.Context, requested int64) (int64, error) {
	lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			// Real DB error
			return 0, err
		}
		// Nothing in db, use requested
		return requested, nil
	}
	// We have lastProcessed, take min(requested, lastProcessed).
	// This is to avoid replaying from a block later than what we have in the DB
	// and skipping blocks.
	return mathutil.Min(requested, lastProcessed.BlockNumber), nil
}

func (lp *logPoller) run() {
	defer close(lp.done)
	logPollTick := time.After(0)
	// trigger first backup poller run shortly after first log poller run
	backupLogPollTick := time.After(100 * time.Millisecond)
	blockPruneTick := time.After(0)
	filtersLoaded := false

	loadFilters := func() error {
		lp.filterMu.Lock()
		defer lp.filterMu.Unlock()
		filters, err := lp.orm.LoadFilters(pg.WithParentCtx(lp.ctx))

		if err != nil {
			return errors.Wrapf(err, "Failed to load initial filters from db, retrying")
		}

		lp.filters = filters
		lp.filterDirty = true
		filtersLoaded = true
		return nil
	}

	for {
		select {
		case <-lp.ctx.Done():
			return
		case replayReq := <-lp.replayStart:
			fromBlock, err := lp.getReplayFromBlock(replayReq.ctx, replayReq.fromBlock)
			if err == nil {
				if !filtersLoaded {
					lp.lggr.Warnw("Received replayReq before filters loaded", "fromBlock", fromBlock, "requested", replayReq.fromBlock)
					if err = loadFilters(); err != nil {
						lp.lggr.Errorw("Failed loading filters during Replay", "err", err, "fromBlock", fromBlock)
					}
				} else {
					// Serially process replay requests.
					lp.lggr.Warnw("Executing replay", "fromBlock", fromBlock, "requested", replayReq.fromBlock)
					lp.pollAndSaveLogs(replayReq.ctx, fromBlock)
				}
			} else {
				lp.lggr.Errorw("Error executing replay, could not get fromBlock", "err", err)
			}
			select {
			case <-lp.ctx.Done():
				// We're shutting down, lets return.
				return
			case <-replayReq.ctx.Done():
				// Client gave up, lets continue.
				continue
			case lp.replayComplete <- err:
			}
		case <-logPollTick:
			logPollTick = time.After(utils.WithJitter(lp.pollPeriod))
			if !filtersLoaded {
				if err := loadFilters(); err != nil {
					lp.lggr.Errorw("Failed loading filters in main logpoller loop, retrying later", "err", err)
					continue
				}
			}

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
				latest, err := lp.ec.HeadByNumber(lp.ctx, nil)
				if err != nil {
					lp.lggr.Warnw("unable to get latest for first poll", "err", err)
					continue
				}
				latestNum := latest.Number
				// Do not support polling chains which don't even have finality depth worth of blocks.
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
		case <-backupLogPollTick:
			// Backup log poller:  this serves as an emergency backup to protect against eventual-consistency behavior
			// of an rpc node (seen occasionally on optimism, but possibly could happen on other chains?).  If the first
			// time we request a block, no logs or incomplete logs come back, this ensures that every log is eventually
			// re-requested after it is finalized.  This doesn't add much overhead, because we can request all of them
			// in one shot, since we don't need to worry about re-orgs after finality depth, and it runs 100x less
			// frequently than the primary log poller.

			// If pollPeriod is set to 1 block time, backup log poller will run once every 100 blocks
			const backupPollerBlockDelay = 100

			backupLogPollTick = time.After(utils.WithJitter(backupPollerBlockDelay * lp.pollPeriod))
			if !filtersLoaded {
				lp.lggr.Warnw("backup log poller ran before filters loaded, skipping")
				continue
			}

			if lp.backupPollerNextBlock == 0 {
				lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(lp.ctx))
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						lp.lggr.Warnw("backup log poller ran before first successful log poller run, skipping")
					} else {
						lp.lggr.Errorw("unable to get starting block", "err", err)
					}
					continue
				}

				// If this is our first run, start max(finalityDepth+1, backupPollerBlockDelay) blocks behind the last processed
				// (or at block 0 if whole blockchain is too short)
				lp.backupPollerNextBlock = lastProcessed.BlockNumber - mathutil.Max(lp.finalityDepth+1, backupPollerBlockDelay)
				if lastProcessed.BlockNumber > backupPollerBlockDelay {
					lp.backupPollerNextBlock = lastProcessed.BlockNumber - backupPollerBlockDelay
				}
			}

			latestBlock, err := lp.ec.HeadByNumber(lp.ctx, nil)
			if err != nil {
				lp.lggr.Warnw("backup logpoller failed to get latest block", "err", err)
				continue
			}

			lastSafeBackfillBlock := latestBlock.Number - lp.finalityDepth - 1
			if lastSafeBackfillBlock >= lp.backupPollerNextBlock {
				lp.lggr.Infow("Backup poller backfilling logs", "start", lp.backupPollerNextBlock, "end", lastSafeBackfillBlock)
				if err = lp.backfill(lp.ctx, lp.backupPollerNextBlock, lastSafeBackfillBlock); err != nil {
					// If there's an error backfilling, we can just return and retry from the last block saved
					// since we don't save any blocks on backfilling. We may re-insert the same logs but thats ok.
					lp.lggr.Warnw("Backup poller failed", "err", err)
					continue
				}
				lp.backupPollerNextBlock = lastSafeBackfillBlock + 1
			}
		case <-blockPruneTick:
			blockPruneTick = time.After(lp.pollPeriod * 1000)
			if err := lp.pruneOldBlocks(lp.ctx); err != nil {
				lp.lggr.Errorw("unable to prune old blocks", "err", err)
			}
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
			EventSig:    l.Topics[0], // First topic is always event signature.
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

// backfill will query FilterLogs in batches for logs in the
// block range [start, end] and save them to the db.
// Retries until ctx cancelled. Will return an error if cancelled
// or if there is an error backfilling.
func (lp *logPoller) backfill(ctx context.Context, start, end int64) error {
	for from := start; from <= end; from += lp.backfillBatchSize {
		to := mathutil.Min(from+lp.backfillBatchSize-1, end)
		logs, err := lp.ec.FilterLogs(ctx, lp.filter(big.NewInt(from), big.NewInt(to), nil))
		if err != nil {
			lp.lggr.Warnw("Unable query for logs, retrying", "err", err, "from", from, "to", to)
			return err
		}
		if len(logs) == 0 {
			continue
		}
		lp.lggr.Infow("Backfill found logs", "from", from, "to", to, "logs", len(logs))
		err = lp.orm.q.WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			return lp.orm.InsertLogs(convertLogs(lp.ec.ChainID(), logs), pg.WithQueryer(tx))
		})
		if err != nil {
			lp.lggr.Warnw("Unable to insert logs, retrying", "err", err, "from", from, "to", to)
			return err
		}
	}
	return nil
}

// getCurrentBlockMaybeHandleReorg accepts a block number
// and will return that block if its parent points to our last saved block.
// One can optionally pass the block header if it has already been queried to avoid an extra RPC call.
// If its parent does not point to our last saved block we know a reorg has occurred,
// so we:
// 1. Find the LCA by following parent hashes.
// 2. Delete all logs and blocks after the LCA
// 3. Return the LCA+1, i.e. our new current (unprocessed) block.
func (lp *logPoller) getCurrentBlockMaybeHandleReorg(ctx context.Context, currentBlockNumber int64, currentBlock *evmtypes.Head) (*evmtypes.Head, error) {
	var err1 error
	if currentBlock == nil {
		// If we don't have the current block already, lets get it.
		currentBlock, err1 = lp.ec.HeadByNumber(ctx, big.NewInt(currentBlockNumber))
		if err1 != nil {
			lp.lggr.Warnw("Unable to get currentBlock", "err", err1, "currentBlockNumber", currentBlockNumber)
			return nil, err1
		}
		// Additional sanity checks, don't necessarily trust the RPC.
		if currentBlock == nil {
			lp.lggr.Errorf("Unexpected nil block from RPC", "currentBlockNumber", currentBlockNumber)
			return nil, errors.Errorf("Got nil block for %d", currentBlockNumber)
		}
		if currentBlock.Number != currentBlockNumber {
			lp.lggr.Warnw("Unable to get currentBlock, rpc returned incorrect block", "currentBlockNumber", currentBlockNumber, "got", currentBlock.Number)
			return nil, errors.Errorf("Block mismatch have %d want %d", currentBlock.Number, currentBlockNumber)
		}
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
		lp.lggr.Infow("Do not have previous block, first poll ever on new chain or after backfill", "currentBlockNumber", currentBlockNumber)
		return currentBlock, nil
	}
	// Check for reorg.
	if currentBlock.ParentHash != expectedParent.BlockHash {
		// There can be another reorg while we're finding the LCA.
		// That is ok, since we'll detect it on the next iteration.
		// Since we go currentBlock by currentBlock for unfinalized logs, the mismatch starts at currentBlockNumber - 1.
		blockAfterLCA, err2 := lp.findBlockAfterLCA(ctx, currentBlock)
		if err2 != nil {
			lp.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err2)
			return nil, errors.New("Unable to find LCA after reorg, retrying")
		}

		lp.lggr.Infow("Reorg detected", "blockAfterLCA", blockAfterLCA.Number, "currentBlockNumber", currentBlockNumber)
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
			err3 := lp.orm.DeleteBlocksAfter(blockAfterLCA.Number, pg.WithQueryer(tx))
			if err3 != nil {
				lp.lggr.Warnw("Unable to clear reorged blocks, retrying", "err", err3)
				return err3
			}
			err3 = lp.orm.DeleteLogsAfter(blockAfterLCA.Number, pg.WithQueryer(tx))
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
	lp.lggr.Debugw("Polling for logs", "currentBlockNumber", currentBlockNumber)
	latestBlock, err := lp.ec.HeadByNumber(ctx, nil)
	if err != nil {
		lp.lggr.Warnw("Unable to get latestBlockNumber block", "err", err, "currentBlockNumber", currentBlockNumber)
		return
	}
	latestBlockNumber := latestBlock.Number
	if currentBlockNumber > latestBlockNumber {
		// Note there can also be a reorg "shortening" i.e. chain height decreases but TDD increases. In that case
		// we also just wait until the new tip is longer and then detect the reorg.
		lp.lggr.Debugw("No new blocks since last poll", "currentBlockNumber", currentBlockNumber, "latestBlockNumber", latestBlockNumber)
		return
	}
	var currentBlock *evmtypes.Head
	if currentBlockNumber == latestBlockNumber {
		// Can re-use our currentBlock and avoid an extra RPC call.
		currentBlock = latestBlock
	}
	// Possibly handle a reorg. For example if we crash, we'll be in the middle of processing unfinalized blocks.
	// Returns (currentBlock || LCA+1 if reorg detected, error)
	currentBlock, err = lp.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber, currentBlock)
	if err != nil {
		// If there's an error handling the reorg, we can't be sure what state the db was left in.
		// Resume from the latest block saved and retry.
		lp.lggr.Errorw("Unable to get current block, retrying", "err", err)
		return
	}
	currentBlockNumber = currentBlock.Number

	// backfill finalized blocks if we can for performance. If we crash during backfill, we
	// may reprocess logs.  Log insertion is idempotent so this is ok.
	// E.g. 1<-2<-3(currentBlockNumber)<-4<-5<-6<-7(latestBlockNumber), finality is 2. So 3,4 can be batched.
	// Although 5 is finalized, we still need to save it to the db for reorg detection if 6 is a reorg.
	// start = currentBlockNumber = 3, end = latestBlockNumber - finality - 1 = 7-2-1 = 4 (inclusive range).
	lastSafeBackfillBlock := latestBlockNumber - lp.finalityDepth - 1
	if lastSafeBackfillBlock >= currentBlockNumber {
		lp.lggr.Infow("Backfilling logs", "start", currentBlockNumber, "end", lastSafeBackfillBlock)
		if err = lp.backfill(ctx, currentBlockNumber, lastSafeBackfillBlock); err != nil {
			// If there's an error backfilling, we can just return and retry from the last block saved
			// since we don't save any blocks on backfilling. We may re-insert the same logs but thats ok.
			lp.lggr.Warnw("Unable to backfill finalized logs, retrying later", "err", err)
			return
		}
		currentBlockNumber = lastSafeBackfillBlock + 1
	}

	if currentBlockNumber > currentBlock.Number {
		// If we successfully backfilled we have logs up to and including lastSafeBackfillBlock,
		// now load the first unfinalized block.
		currentBlock, err = lp.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber, nil)
		if err != nil {
			// If there's an error handling the reorg, we can't be sure what state the db was left in.
			// Resume from the latest block saved.
			lp.lggr.Errorw("Unable to get current block", "err", err)
			return
		}
	}

	for {
		h := currentBlock.Hash
		var logs []types.Log
		logs, err = lp.ec.FilterLogs(ctx, lp.filter(nil, nil, &h))
		if err != nil {
			lp.lggr.Warnw("Unable to query for logs, retrying", "err", err, "block", currentBlockNumber)
			return
		}
		lp.lggr.Debugw("Unfinalized log query", "logs", len(logs), "currentBlockNumber", currentBlockNumber, "blockHash", currentBlock.Hash)
		err = lp.orm.q.WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			if err2 := lp.orm.InsertBlock(h, currentBlockNumber, pg.WithQueryer(tx)); err2 != nil {
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
		// Update current block.
		// Same reorg detection on unfinalized blocks.
		currentBlockNumber++
		if currentBlockNumber > latestBlockNumber {
			break
		}
		currentBlock, err = lp.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber, nil)
		if err != nil {
			// If there's an error handling the reorg, we can't be sure what state the db was left in.
			// Resume from the latest block saved.
			lp.lggr.Errorw("Unable to get current block", "err", err)
			return
		}
		currentBlockNumber = currentBlock.Number
	}
}

// Find the first place where our chain and their chain have the same block,
// that block number is the LCA. Return the block after that, where we want to resume polling.
func (lp *logPoller) findBlockAfterLCA(ctx context.Context, current *evmtypes.Head) (*evmtypes.Head, error) {
	// Current is where the mismatch starts.
	// Check its parent to see if its the same as ours saved.
	parent, err := lp.ec.HeadByHash(ctx, current.ParentHash)
	if err != nil {
		return nil, err
	}
	blockAfterLCA := *current
	reorgStart := parent.Number
	// We expect reorgs up to the block after (current - finalityDepth),
	// since the block at (current - finalityDepth) is finalized.
	// We loop via parent instead of current so current always holds the LCA+1.
	// If the parent block number becomes < the first finalized block our reorg is too deep.
	for parent.Number >= (reorgStart - lp.finalityDepth) {
		ourParentBlockHash, err := lp.orm.SelectBlockByNumber(parent.Number, pg.WithParentCtx(ctx))
		if err != nil {
			return nil, err
		}
		if parent.Hash == ourParentBlockHash.BlockHash {
			// If we do have the blockhash, return blockAfterLCA
			return &blockAfterLCA, nil
		}
		// Otherwise get a new parent and update blockAfterLCA.
		blockAfterLCA = *parent
		parent, err = lp.ec.HeadByHash(ctx, parent.ParentHash)
		if err != nil {
			return nil, err
		}
	}
	lp.lggr.Criticalw("Reorg greater than finality depth detected", "max reorg depth", lp.finalityDepth-1)
	return nil, errors.New("Reorg greater than finality depth")
}

// pruneOldBlocks removes blocks that are > lp.ancientBlockDepth behind the head.
func (lp *logPoller) pruneOldBlocks(ctx context.Context) error {
	latest, err := lp.ec.HeadByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if latest == nil {
		return errors.Errorf("received nil block from RPC")
	}
	if latest.Number <= lp.keepBlocksDepth {
		// No-op, keep all blocks
		return nil
	}
	// 1-2-3-4-5(latest), keepBlocksDepth=3
	// Remove <= 2
	return lp.orm.DeleteBlocksBefore(latest.Number-lp.keepBlocksDepth, pg.WithParentCtx(ctx))
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *logPoller) Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsByBlockRangeFilter(start, end, address, eventSig, qopts...)
}

func (lp *logPoller) LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsWithSigsByBlockRangeFilter(start, end, address, eventSigs, qopts...)
}

// IndexedLogs finds all the logs that have a topic value in topicValues at index topicIndex.
func (lp *logPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogs(address, eventSig, topicIndex, topicValues, confs, qopts...)
}

// LogsDataWordGreaterThan note index is 0 based.
func (lp *logPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectDataWordGreaterThan(address, eventSig, wordIndex, wordValueMin, confs, qopts...)
}

// LogsDataWordRange note index is 0 based.
func (lp *logPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectDataWordRange(address, eventSig, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
}

// IndexedLogsTopicGreaterThan finds all the logs that have a topic value greater than topicValueMin at index topicIndex.
// Only works for integer topics.
func (lp *logPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexLogsTopicGreaterThan(address, eventSig, topicIndex, topicValueMin, confs, qopts...)
}

func (lp *logPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexLogsTopicRange(address, eventSig, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
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
	return lp.orm.SelectLatestLogEventSigWithConfs(eventSig, address, confs, qopts...)
}

func (lp *logPoller) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLatestLogEventSigsAddrsWithConfs(fromBlock, addresses, eventSigs, confs, qopts...)
}

// GetBlocksRange tries to get the specified block numbers from the log pollers
// blocks table. It falls back to the RPC for any unfulfilled requested blocks.
func (lp *logPoller) GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	var blocks []LogPollerBlock

	// Do nothing if no blocks are requested.
	if len(numbers) == 0 {
		return blocks, nil
	}

	// Assign the requested blocks to a mapping.
	blocksRequested := make(map[uint64]struct{})
	for _, b := range numbers {
		blocksRequested[b] = struct{}{}
	}

	// Retrieve all blocks within this range from the log poller.
	blocksFound := make(map[uint64]LogPollerBlock)
	qopts = append(qopts, pg.WithParentCtx(ctx))
	minRequestedBlock := mathutil.Min(numbers[0], numbers[1:]...)
	maxRequestedBlock := mathutil.Max(numbers[0], numbers[1:]...)
	lpBlocks, err := lp.orm.GetBlocksRange(minRequestedBlock, maxRequestedBlock, qopts...)
	if err != nil {
		lp.lggr.Warnw("Error while retrieving blocks from log pollers blocks table. Falling back to RPC...", "requestedBlocks", numbers, "err", err)
	} else {
		for _, b := range lpBlocks {
			if _, ok := blocksRequested[uint64(b.BlockNumber)]; ok {
				// Only fill requested blocks.
				blocksFound[uint64(b.BlockNumber)] = b
			}
		}
		lp.lggr.Debugw("Got blocks from log poller", "blockNumbers", maps.Keys(blocksFound))
	}

	// Fill any remaining blocks from the client.
	blocksFoundFromRPC, err := lp.fillRemainingBlocksFromRPC(ctx, numbers, blocksFound)
	if err != nil {
		return nil, err
	}
	for num, b := range blocksFoundFromRPC {
		blocksFound[num] = b
	}

	var blocksNotFound []uint64
	for _, num := range numbers {
		b, ok := blocksFound[num]
		if !ok {
			blocksNotFound = append(blocksNotFound, num)
		}
		blocks = append(blocks, b)
	}

	if len(blocksNotFound) > 0 {
		return nil, errors.Errorf("blocks were not found in db or RPC call: %v", blocksNotFound)
	}

	return blocks, nil
}

func (lp *logPoller) fillRemainingBlocksFromRPC(
	ctx context.Context,
	blocksRequested []uint64,
	blocksFound map[uint64]LogPollerBlock,
) (map[uint64]LogPollerBlock, error) {
	var reqs []rpc.BatchElem
	var remainingBlocks []uint64
	for _, num := range blocksRequested {
		if _, ok := blocksFound[num]; !ok {
			req := rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{hexutil.EncodeBig(big.NewInt(0).SetUint64(num)), false},
				Result: &evmtypes.Head{},
			}
			reqs = append(reqs, req)
			remainingBlocks = append(remainingBlocks, num)
		}
	}

	if len(remainingBlocks) > 0 {
		lp.lggr.Debugw("falling back to RPC for blocks not found in log poller blocks table",
			"remainingBlocks", remainingBlocks)
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

	var blocksFoundFromRPC = make(map[uint64]LogPollerBlock)
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
		blocksFoundFromRPC[uint64(block.Number)] = LogPollerBlock{
			EvmChainId:  block.EVMChainID,
			BlockHash:   block.Hash,
			BlockNumber: block.Number,
			CreatedAt:   block.Timestamp,
		}
	}

	return blocksFoundFromRPC, nil
}

func EvmWord(i uint64) common.Hash {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BytesToHash(b)
}
