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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

//go:generate mockery --quiet --name LogPoller --output ./mocks/ --case=underscore --structname LogPoller --filename log_poller.go
type LogPoller interface {
	services.ServiceCtx
	Replay(ctx context.Context, fromBlock int64) error
	ReplayAsync(fromBlock int64)
	RegisterFilter(filter Filter, qopts ...pg.QOpt) error
	UnregisterFilter(name string, qopts ...pg.QOpt) error
	HasFilter(name string) bool
	LatestBlock(qopts ...pg.QOpt) (LogPollerBlock, error)
	GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error)

	// General querying
	Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error)
	LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error)
	LogsCreatedAfter(eventSig common.Hash, address common.Address, time time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs Confirmations, qopts ...pg.QOpt) (*Log, error)
	LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	LatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) (int64, error)

	// Content based querying
	IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsByBlockRange(start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsCreatedAfter(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsByTxHash(eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	IndexedLogsWithSigsExcluding(address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
	LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error)
}

type Confirmations int

const (
	Finalized   = Confirmations(-1)
	Unconfirmed = Confirmations(0)
)

type LogPollerTest interface {
	LogPoller
	PollAndSaveLogs(ctx context.Context, currentBlockNumber int64)
	BackupPollAndSaveLogs(ctx context.Context, backupPollerBlockDelay int64)
	Filter(from, to *big.Int, bh *common.Hash) ethereum.FilterQuery
	GetReplayFromBlock(ctx context.Context, requested int64) (int64, error)
	PruneOldBlocks(ctx context.Context) error
}

type Client interface {
	HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error)
	HeadByHash(ctx context.Context, n common.Hash) (*evmtypes.Head, error)
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	ConfiguredChainID() *big.Int
}

var (
	_                       LogPollerTest = &logPoller{}
	ErrReplayRequestAborted               = errors.New("aborted, replay request cancelled")
	ErrReplayInProgress                   = errors.New("replay request cancelled, but replay is already in progress")
	ErrLogPollerShutdown                  = errors.New("replay aborted due to log poller shutdown")
)

type logPoller struct {
	utils.StartStopOnce
	ec                       Client
	orm                      ORM
	lggr                     logger.Logger
	pollPeriod               time.Duration // poll period set by block production rate
	useFinalityTag           bool          // indicates whether logPoller should use chain's finality or pick a fixed depth for finality
	finalityDepth            int64         // finality depth is taken to mean that block (head - finality) is finalized. If `useFinalityTag` is set to true, this value is ignored, because finalityDepth is fetched from chain
	keepFinalizedBlocksDepth int64         // the number of blocks behind the last finalized block we keep in database
	backfillBatchSize        int64         // batch size to use when backfilling finalized logs
	rpcBatchSize             int64         // batch size to use for fallback RPC calls made in GetBlocks
	backupPollerNextBlock    int64

	filterMu        sync.RWMutex
	filters         map[string]Filter
	filterDirty     bool
	cachedAddresses []common.Address
	cachedEventSigs []common.Hash

	replayStart    chan int64
	replayComplete chan error
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewLogPoller creates a log poller. Note there is an assumption
// that blocks can be processed faster than they are produced for the given chain, or the poller will fall behind.
// Block processing involves the following calls in steady state (without reorgs):
//   - eth_getBlockByNumber - headers only (transaction hashes, not full transaction objects),
//   - eth_getLogs - get the logs for the block
//   - 1 db read latest block - for checking reorgs
//   - 1 db tx including block write and logs write to logs.
//
// How fast that can be done depends largely on network speed and DB, but even for the fastest
// support chain, polygon, which has 2s block times, we need RPCs roughly with <= 500ms latency
func NewLogPoller(orm ORM, ec Client, lggr logger.Logger, pollPeriod time.Duration,
	useFinalityTag bool, finalityDepth int64, backfillBatchSize int64, rpcBatchSize int64, keepFinalizedBlocksDepth int64) *logPoller {

	return &logPoller{
		ec:                       ec,
		orm:                      orm,
		lggr:                     lggr.Named("LogPoller"),
		replayStart:              make(chan int64),
		replayComplete:           make(chan error),
		pollPeriod:               pollPeriod,
		finalityDepth:            finalityDepth,
		useFinalityTag:           useFinalityTag,
		backfillBatchSize:        backfillBatchSize,
		rpcBatchSize:             rpcBatchSize,
		keepFinalizedBlocksDepth: keepFinalizedBlocksDepth,
		filters:                  make(map[string]Filter),
		filterDirty:              true, // Always build Filter on first call to cache an empty filter if nothing registered yet.
	}
}

type Filter struct {
	Name      string // see FilterName(id, args) below
	EventSigs evmtypes.HashArray
	Addresses evmtypes.AddressArray
	Retention time.Duration
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

// Contains returns true if this filter already fully Contains a
// filter passed to it.
func (filter *Filter) Contains(other *Filter) bool {
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
//
//	RegisterFilter(event1, addr1)
//	RegisterFilter(event2, addr2)
//
// will result in the poller saving (event1, addr2) or (event2, addr1) as well, should it exist.
// Generally speaking this is harmless. We enforce that EventSigs and Addresses are non-empty,
// which means that anonymous events are not supported and log.Topics >= 1 always (log.Topics[0] is the event signature).
// The filter may be unregistered later by Filter.Name
func (lp *logPoller) RegisterFilter(filter Filter, qopts ...pg.QOpt) error {
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
		if existingFilter.Contains(&filter) {
			// Nothing new in this Filter
			return nil
		}
		lp.lggr.Warnw("Updating existing filter with more events or addresses", "filter", filter)
	} else {
		lp.lggr.Debugw("Creating new filter", "filter", filter)
	}

	if err := lp.orm.InsertFilter(filter, qopts...); err != nil {
		return errors.Wrap(err, "RegisterFilter failed to save filter to db")
	}
	lp.filters[filter.Name] = filter
	lp.filterDirty = true
	return nil
}

func (lp *logPoller) UnregisterFilter(name string, qopts ...pg.QOpt) error {
	lp.filterMu.Lock()
	defer lp.filterMu.Unlock()

	_, ok := lp.filters[name]
	if !ok {
		lp.lggr.Errorf("Filter %s not found", name)
		return nil
	}

	if err := lp.orm.DeleteFilter(name, qopts...); err != nil {
		return errors.Wrapf(err, "Failed to delete filter %s", name)
	}
	delete(lp.filters, name)
	lp.filterDirty = true
	return nil
}

// HasFilter returns true if the log poller has an active filter with the given name.
func (lp *logPoller) HasFilter(name string) bool {
	lp.filterMu.RLock()
	defer lp.filterMu.RUnlock()

	_, ok := lp.filters[name]
	return ok
}

func (lp *logPoller) Filter(from, to *big.Int, bh *common.Hash) ethereum.FilterQuery {
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
// If ctx is cancelled before the replay request has been initiated, ErrReplayRequestAborted is returned.  If the replay
// is already in progress, the replay will continue and ErrReplayInProgress will be returned.  If the client needs a
// guarantee that the replay is complete before proceeding, it should either avoid cancelling or retry until nil is returned
func (lp *logPoller) Replay(ctx context.Context, fromBlock int64) error {
	lp.lggr.Debugf("Replaying from block %d", fromBlock)
	latest, err := lp.ec.HeadByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if fromBlock < 1 || fromBlock > latest.Number {
		return errors.Errorf("Invalid replay block number %v, acceptable range [1, %v]", fromBlock, latest.Number)
	}
	// Block until replay notification accepted or cancelled.
	select {
	case lp.replayStart <- fromBlock:
	case <-ctx.Done():
		return errors.Wrap(ErrReplayRequestAborted, ctx.Err().Error())
	}
	// Block until replay complete or cancelled.
	select {
	case err = <-lp.replayComplete:
		return err
	case <-ctx.Done():
		// Note: this will not abort the actual replay, it just means the client gave up on waiting for it to complete
		lp.wg.Add(1)
		go lp.recvReplayComplete()
		return ErrReplayInProgress
	}
}

func (lp *logPoller) recvReplayComplete() {
	err := <-lp.replayComplete
	if err != nil {
		lp.lggr.Error(err)
	}
	lp.wg.Done()
}

// Asynchronous wrapper for Replay()
func (lp *logPoller) ReplayAsync(fromBlock int64) {
	lp.wg.Add(1)
	go func() {
		if err := lp.Replay(context.Background(), fromBlock); err != nil {
			lp.lggr.Error(err)
		}
		lp.wg.Done()
	}()
}

func (lp *logPoller) Start(parentCtx context.Context) error {
	return lp.StartOnce("LogPoller", func() error {
		ctx, cancel := context.WithCancel(parentCtx)
		lp.ctx = ctx
		lp.cancel = cancel
		lp.wg.Add(1)
		go lp.run()
		return nil
	})
}

func (lp *logPoller) Close() error {
	return lp.StopOnce("LogPoller", func() error {
		select {
		case lp.replayComplete <- ErrLogPollerShutdown:
		default:
		}
		lp.cancel()
		lp.wg.Wait()
		return nil
	})
}

func (lp *logPoller) Name() string {
	return lp.lggr.Name()
}

func (lp *logPoller) HealthReport() map[string]error {
	return map[string]error{lp.Name(): lp.Healthy()}
}

func (lp *logPoller) GetReplayFromBlock(ctx context.Context, requested int64) (int64, error) {
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
	defer lp.wg.Done()
	logPollTick := time.After(0)
	// stagger these somewhat, so they don't all run back-to-back
	backupLogPollTick := time.After(100 * time.Millisecond)
	blockPruneTick := time.After(3 * time.Second)
	logPruneTick := time.After(5 * time.Second)
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
		case fromBlockReq := <-lp.replayStart:
			fromBlock, err := lp.GetReplayFromBlock(lp.ctx, fromBlockReq)
			if err == nil {
				if !filtersLoaded {
					lp.lggr.Warnw("Received replayReq before filters loaded", "fromBlock", fromBlock, "requested", fromBlockReq)
					if err = loadFilters(); err != nil {
						lp.lggr.Errorw("Failed loading filters during Replay", "err", err, "fromBlock", fromBlock)
					}
				}
				if err == nil {
					// Serially process replay requests.
					lp.lggr.Infow("Executing replay", "fromBlock", fromBlock, "requested", fromBlockReq)
					lp.PollAndSaveLogs(lp.ctx, fromBlock)
				}
			} else {
				lp.lggr.Errorw("Error executing replay, could not get fromBlock", "err", err)
			}
			select {
			case <-lp.ctx.Done():
				// We're shutting down, notify client and exit
				select {
				case lp.replayComplete <- ErrReplayRequestAborted:
				default:
				}
				return
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
				latestBlock, latestFinalizedBlockNumber, err := lp.latestBlocks(lp.ctx)
				if err != nil {
					lp.lggr.Warnw("Unable to get latest for first poll", "err", err)
					continue
				}
				// Do not support polling chains which don't even have finality depth worth of blocks.
				// Could conceivably support this but not worth the effort.
				// Need last finalized block number to be higher than 0
				if latestFinalizedBlockNumber <= 0 {
					lp.lggr.Warnw("Insufficient number of blocks on chain, waiting for finality depth", "err", err, "latest", latestBlock.Number)
					continue
				}
				// Starting at the first finalized block. We do not backfill the first finalized block.
				start = latestFinalizedBlockNumber
			} else {
				start = lastProcessed.BlockNumber + 1
			}
			lp.PollAndSaveLogs(lp.ctx, start)
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
				lp.lggr.Warnw("Backup log poller ran before filters loaded, skipping")
				continue
			}
			lp.BackupPollAndSaveLogs(lp.ctx, backupPollerBlockDelay)
		case <-blockPruneTick:
			blockPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 1000))
			if err := lp.PruneOldBlocks(lp.ctx); err != nil {
				lp.lggr.Errorw("Unable to prune old blocks", "err", err)
			}
		case <-logPruneTick:
			logPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 2401)) // = 7^5 avoids common factors with 1000
			if err := lp.orm.DeleteExpiredLogs(pg.WithParentCtx(lp.ctx)); err != nil {
				lp.lggr.Error(err)
			}
		}
	}
}

func (lp *logPoller) BackupPollAndSaveLogs(ctx context.Context, backupPollerBlockDelay int64) {
	if lp.backupPollerNextBlock == 0 {
		lastProcessed, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(ctx))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lp.lggr.Warnw("Backup log poller ran before first successful log poller run, skipping")
			} else {
				lp.lggr.Errorw("Backup log poller unable to get starting block", "err", err)
			}
			return
		}
		// If this is our first run, start from block min(lastProcessed.FinalizedBlockNumber-1, lastProcessed.BlockNumber-backupPollerBlockDelay)
		backupStartBlock := mathutil.Min(lastProcessed.FinalizedBlockNumber-1, lastProcessed.BlockNumber-backupPollerBlockDelay)
		// (or at block 0 if whole blockchain is too short)
		lp.backupPollerNextBlock = mathutil.Max(backupStartBlock, 0)
	}

	_, latestFinalizedBlockNumber, err := lp.latestBlocks(ctx)
	if err != nil {
		lp.lggr.Warnw("Backup logpoller failed to get latest block", "err", err)
		return
	}

	lastSafeBackfillBlock := latestFinalizedBlockNumber - 1
	if lastSafeBackfillBlock >= lp.backupPollerNextBlock {
		lp.lggr.Infow("Backup poller backfilling logs", "start", lp.backupPollerNextBlock, "end", lastSafeBackfillBlock)
		if err = lp.backfill(ctx, lp.backupPollerNextBlock, lastSafeBackfillBlock); err != nil {
			// If there's an error backfilling, we can just return and retry from the last block saved
			// since we don't save any blocks on backfilling. We may re-insert the same logs but thats ok.
			lp.lggr.Warnw("Backup poller failed", "err", err)
			return
		}
		lp.backupPollerNextBlock = lastSafeBackfillBlock + 1
	}
}

// convertLogs converts an array of geth logs ([]type.Log) to an array of logpoller logs ([]Log)
//
//	Block timestamps are extracted from blocks param.  If len(blocks) == 1, the same timestamp from this block
//	will be used for all logs.  If len(blocks) == len(logs) then the block number of each block is used for the
//	corresponding log.  Any other length for blocks is invalid.
func convertLogs(logs []types.Log, blocks []LogPollerBlock, lggr logger.Logger, chainID *big.Int) []Log {
	var lgs []Log
	blockTimestamp := time.Now()
	if len(logs) == 0 {
		return lgs
	}
	if len(blocks) != 1 && len(blocks) != len(logs) {
		lggr.Errorf("AssumptionViolation:  invalid params passed to convertLogs, length of blocks must either be 1 or match length of logs")
		return lgs
	}

	for i, l := range logs {
		if i == 0 || len(blocks) == len(logs) {
			blockTimestamp = blocks[i].BlockTimestamp
		}
		lgs = append(lgs, Log{
			EvmChainId: utils.NewBig(chainID),
			LogIndex:   int64(l.Index),
			BlockHash:  l.BlockHash,
			// We assume block numbers fit in int64
			// in many places.
			BlockNumber:    int64(l.BlockNumber),
			BlockTimestamp: blockTimestamp,
			EventSig:       l.Topics[0], // First topic is always event signature.
			Topics:         convertTopics(l.Topics),
			Address:        l.Address,
			TxHash:         l.TxHash,
			Data:           l.Data,
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

func (lp *logPoller) blocksFromLogs(ctx context.Context, logs []types.Log) (blocks []LogPollerBlock, err error) {
	var numbers []uint64
	for _, log := range logs {
		numbers = append(numbers, log.BlockNumber)
	}

	return lp.GetBlocksRange(ctx, numbers)
}

const jsonRpcLimitExceeded = -32005 // See https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md

// backfill will query FilterLogs in batches for logs in the
// block range [start, end] and save them to the db.
// Retries until ctx cancelled. Will return an error if cancelled
// or if there is an error backfilling.
func (lp *logPoller) backfill(ctx context.Context, start, end int64) error {
	batchSize := lp.backfillBatchSize
	for from := start; from <= end; from += batchSize {
		to := mathutil.Min(from+batchSize-1, end)
		gethLogs, err := lp.ec.FilterLogs(ctx, lp.Filter(big.NewInt(from), big.NewInt(to), nil))
		if err != nil {
			var rpcErr client.JsonError
			if errors.As(err, &rpcErr) {
				if rpcErr.Code != jsonRpcLimitExceeded {
					lp.lggr.Errorw("Unable to query for logs", "err", err, "from", from, "to", to)
					return err
				}
			}
			if batchSize == 1 {
				lp.lggr.Criticalw("Too many log results in a single block, failed to retrieve logs! Node may be running in a degraded state.", "err", err, "from", from, "to", to, "LogBackfillBatchSize", lp.backfillBatchSize)
				return err
			}
			batchSize /= 2
			lp.lggr.Warnw("Too many log results, halving block range batch size.  Consider increasing LogBackfillBatchSize if this happens frequently", "err", err, "from", from, "to", to, "newBatchSize", batchSize, "LogBackfillBatchSize", lp.backfillBatchSize)
			from -= batchSize // counteract +=batchSize on next loop iteration, so starting block does not change
			continue
		}
		if len(gethLogs) == 0 {
			continue
		}
		blocks, err := lp.blocksFromLogs(ctx, gethLogs)
		if err != nil {
			return err
		}

		lp.lggr.Debugw("Backfill found logs", "from", from, "to", to, "logs", len(gethLogs), "blocks", blocks)
		err = lp.orm.Q().WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			return lp.orm.InsertLogs(convertLogs(gethLogs, blocks, lp.lggr, lp.ec.ConfiguredChainID()), pg.WithQueryer(tx))
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
		blockAfterLCA, err2 := lp.findBlockAfterLCA(ctx, currentBlock, expectedParent.FinalizedBlockNumber)
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
		// it would be saved elsewhere e.g. evm.txes, so it seems better to just support the fast reads.
		// Its also nicely analogous to reading from the chain itself.
		err2 = lp.orm.Q().WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
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

// PollAndSaveLogs On startup/crash current is the first block after the last processed block.
// currentBlockNumber is the block from where new logs are to be polled & saved. Under normal
// conditions this would be equal to lastProcessed.BlockNumber + 1.
func (lp *logPoller) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) {
	lp.lggr.Debugw("Polling for logs", "currentBlockNumber", currentBlockNumber)
	// Intentionally not using logPoller.finalityDepth directly but the latestFinalizedBlockNumber returned from lp.latestBlocks()
	// latestBlocks knows how to pick a proper latestFinalizedBlockNumber based on the logPoller's configuration
	latestBlock, latestFinalizedBlockNumber, err := lp.latestBlocks(ctx)
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
	lastSafeBackfillBlock := latestFinalizedBlockNumber - 1
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
		logs, err = lp.ec.FilterLogs(ctx, lp.Filter(nil, nil, &h))
		if err != nil {
			lp.lggr.Warnw("Unable to query for logs, retrying", "err", err, "block", currentBlockNumber)
			return
		}
		lp.lggr.Debugw("Unfinalized log query", "logs", len(logs), "currentBlockNumber", currentBlockNumber, "blockHash", currentBlock.Hash, "timestamp", currentBlock.Timestamp.Unix())
		err = lp.orm.Q().WithOpts(pg.WithParentCtx(ctx)).Transaction(func(tx pg.Queryer) error {
			if err2 := lp.orm.InsertBlock(h, currentBlockNumber, currentBlock.Timestamp, latestFinalizedBlockNumber, pg.WithQueryer(tx)); err2 != nil {
				return err2
			}
			if len(logs) == 0 {
				return nil
			}
			return lp.orm.InsertLogs(convertLogs(logs,
				[]LogPollerBlock{{BlockNumber: currentBlockNumber,
					BlockTimestamp: currentBlock.Timestamp}},
				lp.lggr,
				lp.ec.ConfiguredChainID(),
			), pg.WithQueryer(tx))
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

// Returns information about latestBlock, latestFinalizedBlockNumber
// If finality tag is not enabled, latestFinalizedBlockNumber is calculated as latestBlockNumber - lp.finalityDepth (configured param)
// Otherwise, we return last finalized block number returned from chain
func (lp *logPoller) latestBlocks(ctx context.Context) (*evmtypes.Head, int64, error) {
	// If finality is not enabled, we can only fetch the latest block
	if !lp.useFinalityTag {
		// Example:
		// finalityDepth = 2
		// Blocks: 1->2->3->4->5(latestBlock)
		// latestFinalizedBlockNumber would be 3
		latestBlock, err := lp.ec.HeadByNumber(ctx, nil)
		if err != nil {
			return nil, 0, err
		}
		// If chain has fewer blocks than finalityDepth, return 0
		return latestBlock, mathutil.Max(latestBlock.Number-lp.finalityDepth, 0), nil
	}

	// If finality is enabled, we need to get the latest and finalized blocks.
	blocks, err := lp.batchFetchBlocks(ctx, []string{rpc.LatestBlockNumber.String(), rpc.FinalizedBlockNumber.String()}, 2)
	if err != nil {
		return nil, 0, err
	}
	latest := blocks[0]
	finalized := blocks[1]
	return latest, finalized.Number, nil
}

// Find the first place where our chain and their chain have the same block,
// that block number is the LCA. Return the block after that, where we want to resume polling.
func (lp *logPoller) findBlockAfterLCA(ctx context.Context, current *evmtypes.Head, latestFinalizedBlockNumber int64) (*evmtypes.Head, error) {
	// Current is where the mismatch starts.
	// Check its parent to see if its the same as ours saved.
	parent, err := lp.ec.HeadByHash(ctx, current.ParentHash)
	if err != nil {
		return nil, err
	}
	blockAfterLCA := *current
	// We expect reorgs up to the block after latestFinalizedBlock
	// We loop via parent instead of current so current always holds the LCA+1.
	// If the parent block number becomes < the first finalized block our reorg is too deep.
	// This can happen only if finalityTag is not enabled and fixed finalityDepth is provided via config.
	for parent.Number >= latestFinalizedBlockNumber {
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
	lp.lggr.Criticalw("Reorg greater than finality depth detected", "finalityTag", lp.useFinalityTag, "current", current.Number, "latestFinalized", latestFinalizedBlockNumber)
	rerr := errors.New("Reorg greater than finality depth")
	lp.SvcErrBuffer.Append(rerr)
	return nil, rerr
}

// PruneOldBlocks removes blocks that are > lp.keepFinalizedBlocksDepth behind the latest finalized block.
func (lp *logPoller) PruneOldBlocks(ctx context.Context) error {
	latestBlock, err := lp.orm.SelectLatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return err
	}
	if latestBlock == nil {
		// No blocks saved yet.
		return nil
	}
	if latestBlock.FinalizedBlockNumber <= lp.keepFinalizedBlocksDepth {
		// No-op, keep all blocks
		return nil
	}
	// 1-2-3-4-5(finalized)-6-7(latest), keepFinalizedBlocksDepth=3
	// Remove <= 2
	return lp.orm.DeleteBlocksBefore(latestBlock.FinalizedBlockNumber-lp.keepFinalizedBlocksDepth, pg.WithParentCtx(ctx))
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *logPoller) Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogs(start, end, address, eventSig, qopts...)
}

func (lp *logPoller) LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsWithSigs(start, end, address, eventSigs, qopts...)
}

func (lp *logPoller) LogsCreatedAfter(eventSig common.Hash, address common.Address, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsCreatedAfter(address, eventSig, after, confs, qopts...)
}

// IndexedLogs finds all the logs that have a topic value in topicValues at index topicIndex.
func (lp *logPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogs(address, eventSig, topicIndex, topicValues, confs, qopts...)
}

// IndexedLogsByBlockRange finds all the logs that have a topic value in topicValues at index topicIndex within the block range
func (lp *logPoller) IndexedLogsByBlockRange(start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsByBlockRange(start, end, address, eventSig, topicIndex, topicValues, qopts...)
}

func (lp *logPoller) IndexedLogsCreatedAfter(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsCreatedAfter(address, eventSig, topicIndex, topicValues, after, confs, qopts...)
}

func (lp *logPoller) IndexedLogsByTxHash(eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsByTxHash(eventSig, txHash, qopts...)
}

// LogsDataWordGreaterThan note index is 0 based.
func (lp *logPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsDataWordGreaterThan(address, eventSig, wordIndex, wordValueMin, confs, qopts...)
}

// LogsDataWordRange note index is 0 based.
func (lp *logPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLogsDataWordRange(address, eventSig, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
}

// IndexedLogsTopicGreaterThan finds all the logs that have a topic value greater than topicValueMin at index topicIndex.
// Only works for integer topics.
func (lp *logPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsTopicGreaterThan(address, eventSig, topicIndex, topicValueMin, confs, qopts...)
}

func (lp *logPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsTopicRange(address, eventSig, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
}

// LatestBlock returns the latest block the log poller is on. It tracks blocks to be able
// to detect reorgs.
func (lp *logPoller) LatestBlock(qopts ...pg.QOpt) (LogPollerBlock, error) {
	b, err := lp.orm.SelectLatestBlock(qopts...)
	if err != nil {
		return LogPollerBlock{}, err
	}

	return *b, nil
}

func (lp *logPoller) BlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	return lp.orm.SelectBlockByNumber(n, qopts...)
}

// LatestLogByEventSigWithConfs finds the latest log that has confs number of blocks on top of the log.
func (lp *logPoller) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs Confirmations, qopts ...pg.QOpt) (*Log, error) {
	return lp.orm.SelectLatestLogByEventSigWithConfs(eventSig, address, confs, qopts...)
}

func (lp *logPoller) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectLatestLogEventSigsAddrsWithConfs(fromBlock, addresses, eventSigs, confs, qopts...)
}

func (lp *logPoller) LatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) (int64, error) {
	return lp.orm.SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock, eventSigs, addresses, confs, qopts...)
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
	minRequestedBlock := int64(mathutil.Min(numbers[0], numbers[1:]...))
	maxRequestedBlock := int64(mathutil.Max(numbers[0], numbers[1:]...))
	lpBlocks, err := lp.orm.GetBlocksRange(int64(minRequestedBlock), int64(maxRequestedBlock), qopts...)
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
	blocksFoundFromRPC, err := lp.fillRemainingBlocksFromRPC(ctx, blocksRequested, blocksFound)
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
	blocksRequested map[uint64]struct{},
	blocksFound map[uint64]LogPollerBlock,
) (map[uint64]LogPollerBlock, error) {
	var remainingBlocks []string
	for num := range blocksRequested {
		if _, ok := blocksFound[num]; !ok {
			remainingBlocks = append(remainingBlocks, hexutil.EncodeBig(new(big.Int).SetUint64(num)))
		}
	}

	if len(remainingBlocks) > 0 {
		lp.lggr.Debugw("Falling back to RPC for blocks not found in log poller blocks table",
			"remainingBlocks", remainingBlocks)
	}

	evmBlocks, err := lp.batchFetchBlocks(ctx, remainingBlocks, lp.rpcBatchSize)
	if err != nil {
		return nil, err
	}

	logPollerBlocks := make(map[uint64]LogPollerBlock)
	for _, head := range evmBlocks {
		logPollerBlocks[uint64(head.Number)] = LogPollerBlock{
			EvmChainId:     head.EVMChainID,
			BlockHash:      head.Hash,
			BlockNumber:    head.Number,
			BlockTimestamp: head.Timestamp,
			CreatedAt:      head.Timestamp,
		}
	}
	return logPollerBlocks, nil
}

func (lp *logPoller) batchFetchBlocks(ctx context.Context, blocksRequested []string, batchSize int64) ([]*evmtypes.Head, error) {
	reqs := make([]rpc.BatchElem, 0, len(blocksRequested))
	for _, num := range blocksRequested {
		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{num, false},
			Result: &evmtypes.Head{},
		}
		reqs = append(reqs, req)
	}

	for i := 0; i < len(reqs); i += int(batchSize) {
		j := i + int(batchSize)
		if j > len(reqs) {
			j = len(reqs)
		}

		err := lp.ec.BatchCallContext(ctx, reqs[i:j])
		if err != nil {
			return nil, err
		}
	}

	var blocks = make([]*evmtypes.Head, 0, len(reqs))
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
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// IndexedLogsWithSigsExcluding returns the set difference(A-B) of logs with signature sigA and sigB, matching is done on the topics index
//
// For example, query to retrieve unfulfilled requests by querying request log events without matching fulfillment log events.
// The order of events is not significant. Both logs must be inside the block range and have the minimum number of confirmations
func (lp *logPoller) IndexedLogsWithSigsExcluding(address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return lp.orm.SelectIndexedLogsWithSigsExcluding(eventSigA, eventSigB, topicIndex, address, fromBlock, toBlock, confs, qopts...)
}

func EvmWord(i uint64) common.Hash {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BytesToHash(b)
}
