package logpoller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

//go:generate mockery --quiet --name LogPoller --output ./mocks/ --case=underscore --structname LogPoller --filename log_poller.go
type LogPoller interface {
	services.Service
	Replay(ctx context.Context, fromBlock int64) error
	ReplayAsync(fromBlock int64)
	RegisterFilter(ctx context.Context, filter Filter) error
	UnregisterFilter(ctx context.Context, name string) error
	HasFilter(name string) bool
	LatestBlock(ctx context.Context) (LogPollerBlock, error)
	GetBlocksRange(ctx context.Context, numbers []uint64) ([]LogPollerBlock, error)

	// General querying
	Logs(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address) ([]Log, error)
	LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]Log, error)
	LogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, time time.Time, confs Confirmations) ([]Log, error)
	LatestLogByEventSigWithConfs(ctx context.Context, eventSig common.Hash, address common.Address, confs Confirmations) (*Log, error)
	LatestLogEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations) ([]Log, error)
	LatestBlockByEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations) (int64, error)

	// Content based querying
	IndexedLogs(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs Confirmations) ([]Log, error)
	IndexedLogsByBlockRange(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash) ([]Log, error)
	IndexedLogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations) ([]Log, error)
	IndexedLogsByTxHash(ctx context.Context, eventSig common.Hash, address common.Address, txHash common.Hash) ([]Log, error)
	IndexedLogsTopicGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs Confirmations) ([]Log, error)
	IndexedLogsTopicRange(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs Confirmations) ([]Log, error)
	IndexedLogsWithSigsExcluding(ctx context.Context, address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs Confirmations) ([]Log, error)
	LogsDataWordRange(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations) ([]Log, error)
	LogsDataWordGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs Confirmations) ([]Log, error)
	LogsDataWordBetween(ctx context.Context, eventSig common.Hash, address common.Address, wordIndexMin, wordIndexMax int, wordValue common.Hash, confs Confirmations) ([]Log, error)
}

// GetLogsBatchElem hides away all the interface casting, so the fields can be accessed more easily, and with type safety
type GetLogsBatchElem rpc.BatchElem

func NewGetLogsReq(filter Filter) *GetLogsBatchElem {
	topics := make2DTopics(filter.EventSigs, filter.Topic2, filter.Topic3, filter.Topic4)

	params := map[string]interface{}{
		"address": []common.Address(filter.Addresses),
		"topics":  topics,
	}

	return &GetLogsBatchElem{
		Method: "eth_getLogs",
		Args:   []interface{}{params},
		Result: new([]types.Log),
	}
}

func (e GetLogsBatchElem) params() map[string]interface{} {
	return e.Args[0].(map[string]interface{})
}

func (e GetLogsBatchElem) Addresses() []common.Address {
	return e.params()["address"].([]common.Address)
}

func (e GetLogsBatchElem) SetAddresses(addresses []common.Address) {
	e.params()["address"] = addresses
}

func (e GetLogsBatchElem) Topics() [][]common.Hash {
	return e.params()["topics"].([][]common.Hash)
}

func (e GetLogsBatchElem) SetTopics(topics [][]common.Hash) {
	e.params()["topics"] = topics
}

func (e GetLogsBatchElem) FromBlock() *big.Int {
	fromBlock, ok := e.params()["fromBlock"].(*big.Int)
	if !ok {
		return nil
	}
	return fromBlock
}

func (e GetLogsBatchElem) ToBlock() *big.Int {
	toBlock, ok := e.params()["fromBlock"].(*big.Int)
	if !ok {
		return nil
	}
	return toBlock
}

func (e GetLogsBatchElem) BlockHash() *common.Hash {
	blockHash, ok := e.params()["blockHash"].(*common.Hash)
	if !ok {
		return nil
	}
	return blockHash
}

func (e GetLogsBatchElem) SetFromBlock(fromBlock *big.Int) {
	e.params()["fromBlock"] = fromBlock
}

type Confirmations int

const (
	Finalized   = Confirmations(-1)
	Unconfirmed = Confirmations(0)
)

type LogPollerTest interface {
	LogPoller
	PollAndSaveLogs(ctx context.Context, currentBlockNumber int64)
	BackupPollAndSaveLogs(ctx context.Context)
	EthGetLogsReqs(fromBlock, toBlock *big.Int, blockHash *common.Hash) []GetLogsBatchElem
	GetReplayFromBlock(ctx context.Context, requested int64) (int64, error)
	PruneOldBlocks(ctx context.Context) (bool, error)
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
	services.StateMachine
	ec                          Client
	orm                         ORM
	lggr                        logger.SugaredLogger
	pollPeriod                  time.Duration // poll period set by block production rate
	useFinalityTag              bool          // indicates whether logPoller should use chain's finality or pick a fixed depth for finality
	finalityDepth               int64         // finality depth is taken to mean that block (head - finality) is finalized. If `useFinalityTag` is set to true, this value is ignored, because finalityDepth is fetched from chain
	keepFinalizedBlocksDepth    int64         // the number of blocks behind the last finalized block we keep in database
	backfillBatchSize           int64         // batch size to use when backfilling finalized logs
	rpcBatchSize                int64         // batch size to use for fallback RPC calls made in GetBlocks
	logPrunePageSize            int64
	backupPollerNextBlock       int64 // next block to be processed by Backup LogPoller
	backupPollerBlockDelay      int64 // how far behind regular LogPoller should BackupLogPoller run. 0 = disabled
	filtersMu                   sync.RWMutex
	filters                     map[string]Filter
	newFilters                  map[string]struct{}                    // Set of filter names which have been added since cached reqs indices were last rebuilt
	removedFilters              []Filter                               // Slice of filters which have been removed or replaced since cached reqs indices were last rebuilt
	cachedReqsByAddress         map[common.Address][]*GetLogsBatchElem // Index of cached GetLogs requests, by contract address
	cachedReqsByEventsTopicsKey map[string]*GetLogsBatchElem           // Index of cached GetLogs requests, by eventTopicsKey

	replayStart    chan int64
	replayComplete chan error
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

type Opts struct {
	PollPeriod               time.Duration
	UseFinalityTag           bool
	FinalityDepth            int64
	BackfillBatchSize        int64
	RpcBatchSize             int64
	KeepFinalizedBlocksDepth int64
	BackupPollerBlockDelay   int64
	LogPrunePageSize         int64
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
func NewLogPoller(orm ORM, ec Client, lggr logger.Logger, opts Opts) *logPoller {
	ctx, cancel := context.WithCancel(context.Background())
	return &logPoller{
		ctx:                         ctx,
		cancel:                      cancel,
		ec:                          ec,
		orm:                         orm,
		lggr:                        logger.Sugared(logger.Named(lggr, "LogPoller")),
		replayStart:                 make(chan int64),
		replayComplete:              make(chan error),
		pollPeriod:                  opts.PollPeriod,
		backupPollerBlockDelay:      opts.BackupPollerBlockDelay,
		finalityDepth:               opts.FinalityDepth,
		useFinalityTag:              opts.UseFinalityTag,
		backfillBatchSize:           opts.BackfillBatchSize,
		rpcBatchSize:                opts.RpcBatchSize,
		keepFinalizedBlocksDepth:    opts.KeepFinalizedBlocksDepth,
		logPrunePageSize:            opts.LogPrunePageSize,
		cachedReqsByAddress:         make(map[common.Address][]*GetLogsBatchElem),
		cachedReqsByEventsTopicsKey: make(map[string]*GetLogsBatchElem),
		filters:                     make(map[string]Filter),
		newFilters:                  make(map[string]struct{}),
	}
}

type Filter struct {
	Name         string // see FilterName(id, args) below
	Addresses    evmtypes.AddressArray
	EventSigs    evmtypes.HashArray // list of possible values for eventsig (aka topic1)
	Topic2       evmtypes.HashArray // list of possible values for topic2
	Topic3       evmtypes.HashArray // list of possible values for topic3
	Topic4       evmtypes.HashArray // list of possible values for topic4
	Retention    time.Duration      // maximum amount of time to retain logs
	MaxLogsKept  uint64             // maximum number of logs to retain ( 0 = unlimited )
	LogsPerBlock uint64             // rate limit ( maximum # of logs per block, 0 = unlimited )
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
	for _, addr := range other.Addresses {
		if _, ok := addresses[addr]; !ok {
			return false
		}
	}

	return isTopicsSubset(
		make2DTopics(other.EventSigs, other.Topic2, other.Topic3, other.Topic4),
		make2DTopics(filter.EventSigs, filter.Topic2, filter.Topic3, filter.Topic4),
	)
}

type BytesRepresentable interface {
	Bytes() []byte
}

// sortByteArrays can sort a slice of byte arrays (eg common.Address or common.Hash)
// by comparing bytes.  It will also remove any duplicate entries found, and
// ensure that what's returned is a copy rather than the original
func sortDeDupByteArrays[T BytesRepresentable](vals []T) (sorted []T) {
	if len(vals) <= 1 {
		copy(sorted, vals)
		return vals
	}

	slices.SortStableFunc(vals, func(b1, b2 T) int {
		return bytes.Compare(
			b1.Bytes(),
			b2.Bytes())
	})

	res := []T{vals[0]}
	for _, val := range vals { // de-dupe
		if !bytes.Equal(val.Bytes(), res[len(res)-1].Bytes()) {
			res = append(res, val)
		}
	}
	return res
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
// Warnings/debug information is keyed by filter name.
func (lp *logPoller) RegisterFilter(ctx context.Context, filter Filter) error {
	if len(filter.Addresses) == 0 {
		return fmt.Errorf("at least one address must be specified")
	}
	if len(filter.EventSigs) == 0 {
		return fmt.Errorf("at least one event must be specified")
	}

	for _, eventSig := range filter.EventSigs {
		if eventSig == [common.HashLength]byte{} {
			return fmt.Errorf("empty event sig")
		}
	}
	for _, addr := range filter.Addresses {
		if addr == [common.AddressLength]byte{} {
			return fmt.Errorf("empty address")
		}
	}

	// Sort all of these, to speed up comparisons between topics & addresses of different filters
	filter.Addresses = sortDeDupByteArrays(filter.Addresses)
	filter.EventSigs = sortDeDupByteArrays(filter.EventSigs)
	filter.Topic2 = sortDeDupByteArrays(filter.Topic2)
	filter.Topic3 = sortDeDupByteArrays(filter.Topic3)
	filter.Topic4 = sortDeDupByteArrays(filter.Topic4)

	lp.filtersMu.Lock()
	defer lp.filtersMu.Unlock()

	if existingFilter, ok := lp.filters[filter.Name]; ok {
		if existingFilter.Contains(&filter) {
			// Nothing new in this Filter
			lp.lggr.Warnw("Filter already present, no-op", "name", filter.Name, "filter", filter)
			return nil
		}
		lp.lggr.Warnw("Updating existing filter with more events or addresses", "name", filter.Name, "filter", filter)
		lp.removedFilters = append(lp.removedFilters, existingFilter)
	}

	if err := lp.orm.InsertFilter(ctx, filter); err != nil {
		return fmt.Errorf("error inserting filter: %w", err)
	}
	lp.filters[filter.Name] = filter
	lp.newFilters[filter.Name] = struct{}{}

	lp.lggr.Debugw("RegisterFilter: registered new filter", "name", filter.Name, "addresses", filter.Addresses, "eventSigs", filter.EventSigs)
	return nil
}

// UnregisterFilter will remove the filter with the given name.
// If the name does not exist, it will log an error but not return an error.
// Warnings/debug information is keyed by filter name.
func (lp *logPoller) UnregisterFilter(ctx context.Context, name string) error {
	lp.filtersMu.Lock()
	defer lp.filtersMu.Unlock()

	_, ok := lp.filters[name]
	if !ok {
		lp.lggr.Warnw("Filter not found", "name", name)
		return nil
	}

	if err := lp.orm.DeleteFilter(ctx, name); err != nil {
		return fmt.Errorf("error deleting filter: %w", err)
	}

	lp.removedFilters = append(lp.removedFilters, lp.filters[name])
	delete(lp.filters, name)

	return nil
}

// HasFilter returns true if the log poller has an active filter with the given name.
func (lp *logPoller) HasFilter(name string) bool {
	lp.filtersMu.RLock()
	defer lp.filtersMu.RUnlock()

	_, ok := lp.filters[name]
	return ok
}

// Replay signals that the poller should resume from a new block.
// Blocks until the replay is complete.
// Replay can be used to ensure that filter modification has been applied for all blocks from "fromBlock" up to latest.
// If ctx is cancelled before the replay request has been initiated, ErrReplayRequestAborted is returned.  If the replay
// is already in progress, the replay will continue and ErrReplayInProgress will be returned.  If the client needs a
// guarantee that the replay is complete before proceeding, it should either avoid cancelling or retry until nil is returned
func (lp *logPoller) Replay(ctx context.Context, fromBlock int64) (err error) {
	defer func() {
		if errors.Is(err, context.Canceled) {
			err = ErrReplayRequestAborted
		}
	}()

	lp.lggr.Debugf("Replaying from block", fromBlock)
	latest, err := lp.ec.HeadByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if fromBlock < 1 || fromBlock > latest.Number {
		return fmt.Errorf("Invalid replay block number %v, acceptable range [1, %v]", fromBlock, latest.Number)
	}

	// Backfill all logs up to the latest saved finalized block outside the LogPoller's main loop.
	// This is safe, because chain cannot be rewinded deeper than that, so there must not be any race conditions.
	savedFinalizedBlockNumber, err := lp.savedFinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}
	if fromBlock <= savedFinalizedBlockNumber {
		err = lp.backfill(ctx, fromBlock, savedFinalizedBlockNumber)
		if err != nil {
			return err
		}
	}

	// Poll everything after latest finalized block in main loop to avoid concurrent writes during reorg
	// We assume that number of logs between saved finalized block and current head is small enough to be processed in main loop
	fromBlock = mathutil.Max(fromBlock, savedFinalizedBlockNumber+1)
	// Don't continue if latest block number is the same as saved finalized block number
	if fromBlock > latest.Number {
		return nil
	}
	// Block until replay notification accepted or cancelled.
	select {
	case lp.replayStart <- fromBlock:
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrReplayRequestAborted, ctx.Err())
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

// savedFinalizedBlockNumber returns the FinalizedBlockNumber saved with the last processed block in the db
// (latestFinalizedBlock at the time the last processed block was saved)
// If this is the first poll and no blocks are in the db, it returns 0
func (lp *logPoller) savedFinalizedBlockNumber(ctx context.Context) (int64, error) {
	latestProcessed, err := lp.LatestBlock(ctx)
	if err == nil {
		return latestProcessed.FinalizedBlockNumber, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return 0, err
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
		if err := lp.Replay(lp.ctx, fromBlock); err != nil {
			lp.lggr.Error(err)
		}
		lp.wg.Done()
	}()
}

func (lp *logPoller) Start(context.Context) error {
	return lp.StartOnce("LogPoller", func() error {
		lp.wg.Add(2)
		go lp.run()
		go lp.backgroundWorkerRun()
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
	lastProcessed, err := lp.orm.SelectLatestBlock(ctx)
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

func (lp *logPoller) loadFilters() error {
	lp.filtersMu.Lock()
	defer lp.filtersMu.Unlock()
	filters, err := lp.orm.LoadFilters(lp.ctx)

	if err != nil {
		return fmt.Errorf("Failed to load initial filters from db, retrying: %w", err)
	}

	for name, filter := range filters {
		lp.filters[name] = filter
		lp.newFilters[name] = struct{}{}
	}

	return nil
}

func (lp *logPoller) run() {
	defer lp.wg.Done()
	logPollTick := time.After(0)
	// stagger these somewhat, so they don't all run back-to-back
	backupLogPollTick := time.After(100 * time.Millisecond)
	filtersLoaded := false

	for {
		select {
		case <-lp.ctx.Done():
			return
		case fromBlockReq := <-lp.replayStart:
			lp.handleReplayRequest(fromBlockReq, filtersLoaded)
		case <-logPollTick:
			logPollTick = time.After(utils.WithJitter(lp.pollPeriod))
			if !filtersLoaded {
				if err := lp.loadFilters(); err != nil {
					lp.lggr.Errorw("Failed loading filters in main logpoller loop, retrying later", "err", err)
					continue
				}
				filtersLoaded = true
			}

			// Always start from the latest block in the db.
			var start int64
			lastProcessed, err := lp.orm.SelectLatestBlock(lp.ctx)
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
			if lp.backupPollerBlockDelay == 0 {
				continue // backup poller is disabled
			}
			// Backup log poller:  this serves as an emergency backup to protect against eventual-consistency behavior
			// of an rpc node (seen occasionally on optimism, but possibly could happen on other chains?).  If the first
			// time we request a block, no logs or incomplete logs come back, this ensures that every log is eventually
			// re-requested after it is finalized. This doesn't add much overhead, because we can request all of them
			// in one shot, since we don't need to worry about re-orgs after finality depth, and it runs far less
			// frequently than the primary log poller (instead of roughly once per block it runs once roughly once every
			// lp.backupPollerDelay blocks--with default settings about 100x less frequently).

			backupLogPollTick = time.After(utils.WithJitter(time.Duration(lp.backupPollerBlockDelay) * lp.pollPeriod))
			if !filtersLoaded {
				lp.lggr.Warnw("Backup log poller ran before filters loaded, skipping")
				continue
			}
			lp.BackupPollAndSaveLogs(lp.ctx)
		}
	}
}

func (lp *logPoller) backgroundWorkerRun() {
	defer lp.wg.Done()

	// Avoid putting too much pressure on the database by staggering the pruning of old blocks and logs.
	// Usually, node after restart will have some work to boot the plugins and other services.
	// Deferring first prune by minutes reduces risk of putting too much pressure on the database.
	blockPruneTick := time.After(5 * time.Minute)
	logPruneTick := time.After(10 * time.Minute)

	for {
		select {
		case <-lp.ctx.Done():
			return
		case <-blockPruneTick:
			blockPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 1000))
			if allRemoved, err := lp.PruneOldBlocks(lp.ctx); err != nil {
				lp.lggr.Errorw("Unable to prune old blocks", "err", err)
			} else if !allRemoved {
				// Tick faster when cleanup can't keep up with the pace of new blocks
				blockPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 100))
			}
		case <-logPruneTick:
			logPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 2401)) // = 7^5 avoids common factors with 1000
			if allRemoved, err := lp.PruneExpiredLogs(lp.ctx); err != nil {
				lp.lggr.Errorw("Unable to prune expired logs", "err", err)
			} else if !allRemoved {
				// Tick faster when cleanup can't keep up with the pace of new logs
				logPruneTick = time.After(utils.WithJitter(lp.pollPeriod * 241))
			}
		}
	}
}

func (lp *logPoller) handleReplayRequest(fromBlockReq int64, filtersLoaded bool) {
	fromBlock, err := lp.GetReplayFromBlock(lp.ctx, fromBlockReq)
	if err == nil {
		if !filtersLoaded {
			lp.lggr.Warnw("Received replayReq before filters loaded", "fromBlock", fromBlock, "requested", fromBlockReq)
			if err = lp.loadFilters(); err != nil {
				lp.lggr.Errorw("Failed loading filters during Replay", "err", err, "fromBlock", fromBlock)
			}
		}
		if err == nil {
			// Serially process replay requests.
			lp.lggr.Infow("Executing replay", "fromBlock", fromBlock, "requested", fromBlockReq)
			lp.PollAndSaveLogs(lp.ctx, fromBlock)
			lp.lggr.Infow("Executing replay finished", "fromBlock", fromBlock, "requested", fromBlockReq)
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
}

func (lp *logPoller) BackupPollAndSaveLogs(ctx context.Context) {
	if lp.backupPollerNextBlock == 0 {
		lastProcessed, err := lp.orm.SelectLatestBlock(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lp.lggr.Warnw("Backup log poller ran before first successful log poller run, skipping")
			} else {
				lp.lggr.Errorw("Backup log poller unable to get starting block", "err", err)
			}
			return
		}
		// If this is our first run, start from block min(lastProcessed.FinalizedBlockNumber-1, lastProcessed.BlockNumber-backupPollerBlockDelay)
		backupStartBlock := mathutil.Min(lastProcessed.FinalizedBlockNumber-1, lastProcessed.BlockNumber-lp.backupPollerBlockDelay)
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
		lp.lggr.Infow("Backup poller started backfilling logs", "start", lp.backupPollerNextBlock, "end", lastSafeBackfillBlock)
		if err = lp.backfill(ctx, lp.backupPollerNextBlock, lastSafeBackfillBlock); err != nil {
			// If there's an error backfilling, we can just return and retry from the last block saved
			// since we don't save any blocks on backfilling. We may re-insert the same logs but thats ok.
			lp.lggr.Warnw("Backup poller failed", "err", err)
			return
		}
		lp.lggr.Infow("Backup poller finished backfilling", "start", lp.backupPollerNextBlock, "end", lastSafeBackfillBlock)
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
			EvmChainId: ubig.New(chainID),
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
		gethLogs, err := lp.batchFetchLogs(ctx, big.NewInt(from), big.NewInt(to), nil)
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
		err = lp.orm.InsertLogsWithBlock(ctx, convertLogs(gethLogs, blocks, lp.lggr, lp.ec.ConfiguredChainID()), blocks[len(blocks)-1])
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
			return nil, fmt.Errorf("Got nil block for %d", currentBlockNumber)
		}
		if currentBlock.Number != currentBlockNumber {
			lp.lggr.Warnw("Unable to get currentBlock, rpc returned incorrect block", "currentBlockNumber", currentBlockNumber, "got", currentBlock.Number)
			return nil, fmt.Errorf("Block mismatch have %d want %d", currentBlock.Number, currentBlockNumber)
		}
	}
	// Does this currentBlock point to the same parent that we have saved?
	// If not, there was a reorg, so we need to rewind.

	expectedParent, err1 := lp.orm.SelectBlockByNumber(ctx, currentBlockNumber-1)
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
		err2 = lp.orm.DeleteLogsAndBlocksAfter(ctx, blockAfterLCA.Number)
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
		logs, err = lp.batchFetchLogs(ctx, nil, nil, &h)
		if err != nil {
			lp.lggr.Warnw("Unable to query for logs, retrying", "err", err, "block", currentBlockNumber)
			return
		}
		lp.lggr.Debugw("Unfinalized log query", "logs", len(logs), "currentBlockNumber", currentBlockNumber, "blockHash", currentBlock.Hash, "timestamp", currentBlock.Timestamp.Unix())
		block := NewLogPollerBlock(h, currentBlockNumber, currentBlock.Timestamp, latestFinalizedBlockNumber)
		err = lp.orm.InsertLogsWithBlock(
			ctx,
			convertLogs(logs, []LogPollerBlock{block}, lp.lggr, lp.ec.ConfiguredChainID()),
			block,
		)
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
	lp.lggr.Debugw("Latest blocks read from chain", "latest", latest.Number, "finalized", finalized.Number)
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
		ourParentBlockHash, err := lp.orm.SelectBlockByNumber(ctx, parent.Number)
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
// Returns whether all blocks eligible for pruning were removed. If logPrunePageSize is set to 0, it will always return true.
func (lp *logPoller) PruneOldBlocks(ctx context.Context) (bool, error) {
	latestBlock, err := lp.orm.SelectLatestBlock(ctx)
	if err != nil {
		return false, err
	}
	if latestBlock == nil {
		// No blocks saved yet.
		return true, nil
	}
	if latestBlock.FinalizedBlockNumber <= lp.keepFinalizedBlocksDepth {
		// No-op, keep all blocks
		return true, nil
	}
	// 1-2-3-4-5(finalized)-6-7(latest), keepFinalizedBlocksDepth=3
	// Remove <= 2
	rowsRemoved, err := lp.orm.DeleteBlocksBefore(
		ctx,
		latestBlock.FinalizedBlockNumber-lp.keepFinalizedBlocksDepth,
		lp.logPrunePageSize,
	)
	return lp.logPrunePageSize == 0 || rowsRemoved < lp.logPrunePageSize, err
}

// PruneExpiredLogs logs that are older than their retention period defined in Filter.
// Returns whether all logs eligible for pruning were removed. If logPrunePageSize is set to 0, it will always return true.
func (lp *logPoller) PruneExpiredLogs(ctx context.Context) (bool, error) {
	rowsRemoved, err := lp.orm.DeleteExpiredLogs(ctx, lp.logPrunePageSize)
	return lp.logPrunePageSize == 0 || rowsRemoved < lp.logPrunePageSize, err
}

// Logs returns logs matching topics and address (exactly) in the given block range,
// which are canonical at time of query.
func (lp *logPoller) Logs(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address) ([]Log, error) {
	return lp.orm.SelectLogs(ctx, start, end, address, eventSig)
}

func (lp *logPoller) LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]Log, error) {
	return lp.orm.SelectLogsWithSigs(ctx, start, end, address, eventSigs)
}

func (lp *logPoller) LogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, after time.Time, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectLogsCreatedAfter(ctx, address, eventSig, after, confs)
}

// IndexedLogs finds all the logs that have a topic value in topicValues at index topicIndex.
func (lp *logPoller) IndexedLogs(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectIndexedLogs(ctx, address, eventSig, topicIndex, topicValues, confs)
}

// IndexedLogsByBlockRange finds all the logs that have a topic value in topicValues at index topicIndex within the block range
func (lp *logPoller) IndexedLogsByBlockRange(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash) ([]Log, error) {
	return lp.orm.SelectIndexedLogsByBlockRange(ctx, start, end, address, eventSig, topicIndex, topicValues)
}

func (lp *logPoller) IndexedLogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectIndexedLogsCreatedAfter(ctx, address, eventSig, topicIndex, topicValues, after, confs)
}

func (lp *logPoller) IndexedLogsByTxHash(ctx context.Context, eventSig common.Hash, address common.Address, txHash common.Hash) ([]Log, error) {
	return lp.orm.SelectIndexedLogsByTxHash(ctx, address, eventSig, txHash)
}

// LogsDataWordGreaterThan note index is 0 based.
func (lp *logPoller) LogsDataWordGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectLogsDataWordGreaterThan(ctx, address, eventSig, wordIndex, wordValueMin, confs)
}

// LogsDataWordRange note index is 0 based.
func (lp *logPoller) LogsDataWordRange(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectLogsDataWordRange(ctx, address, eventSig, wordIndex, wordValueMin, wordValueMax, confs)
}

// IndexedLogsTopicGreaterThan finds all the logs that have a topic value greater than topicValueMin at index topicIndex.
// Only works for integer topics.
func (lp *logPoller) IndexedLogsTopicGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectIndexedLogsTopicGreaterThan(ctx, address, eventSig, topicIndex, topicValueMin, confs)
}

func (lp *logPoller) IndexedLogsTopicRange(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectIndexedLogsTopicRange(ctx, address, eventSig, topicIndex, topicValueMin, topicValueMax, confs)
}

// LatestBlock returns the latest block the log poller is on. It tracks blocks to be able
// to detect reorgs.
func (lp *logPoller) LatestBlock(ctx context.Context) (LogPollerBlock, error) {
	b, err := lp.orm.SelectLatestBlock(ctx)
	if err != nil {
		return LogPollerBlock{}, err
	}

	return *b, nil
}

func (lp *logPoller) BlockByNumber(ctx context.Context, n int64) (*LogPollerBlock, error) {
	return lp.orm.SelectBlockByNumber(ctx, n)
}

// LatestLogByEventSigWithConfs finds the latest log that has confs number of blocks on top of the log.
func (lp *logPoller) LatestLogByEventSigWithConfs(ctx context.Context, eventSig common.Hash, address common.Address, confs Confirmations) (*Log, error) {
	return lp.orm.SelectLatestLogByEventSigWithConfs(ctx, eventSig, address, confs)
}

func (lp *logPoller) LatestLogEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectLatestLogEventSigsAddrsWithConfs(ctx, fromBlock, addresses, eventSigs, confs)
}

func (lp *logPoller) LatestBlockByEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations) (int64, error) {
	return lp.orm.SelectLatestBlockByEventSigsAddrsWithConfs(ctx, fromBlock, eventSigs, addresses, confs)
}

// LogsDataWordBetween retrieves a slice of Log records that match specific criteria.
// Besides generic filters like eventSig, address and confs, it also verifies data content against wordValue
// data[wordIndexMin] <= wordValue <= data[wordIndexMax].
//
// Passing the same value for wordIndexMin and wordIndexMax will check the equality of the wordValue at that index.
// Leading to returning logs matching: data[wordIndexMin] == wordValue.
//
// This function is particularly useful for filtering logs by data word values and their positions within the event data.
// It returns an empty slice if no logs match the provided criteria.
func (lp *logPoller) LogsDataWordBetween(ctx context.Context, eventSig common.Hash, address common.Address, wordIndexMin, wordIndexMax int, wordValue common.Hash, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectLogsDataWordBetween(ctx, address, eventSig, wordIndexMin, wordIndexMax, wordValue, confs)
}

// GetBlocksRange tries to get the specified block numbers from the log pollers
// blocks table. It falls back to the RPC for any unfulfilled requested blocks.
func (lp *logPoller) GetBlocksRange(ctx context.Context, numbers []uint64) ([]LogPollerBlock, error) {
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
	minRequestedBlock := int64(mathutil.Min(numbers[0], numbers[1:]...))
	maxRequestedBlock := int64(mathutil.Max(numbers[0], numbers[1:]...))
	lpBlocks, err := lp.orm.GetBlocksRange(ctx, minRequestedBlock, maxRequestedBlock)
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
		return nil, fmt.Errorf("blocks were not found in db or RPC call: %v", blocksNotFound)
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

// mergeAddressesIntoGetLogsReq merges a new list of addresses into a GetLogs req,
// while preserving sort order and removing duplicates
func mergeAddressesIntoGetLogsReq(req *GetLogsBatchElem, newAddresses []common.Address) {
	var merged []common.Address
	var i, j int
	addresses := req.Addresses()

	for i < len(addresses) && j < len(newAddresses) {
		cmp := bytes.Compare(newAddresses[j].Bytes(), addresses[i].Bytes())
		if cmp < 0 {
			merged = append(merged, newAddresses[j])
			j++
		} else if cmp > 0 {
			merged = append(merged, addresses[i])
			i++
		} else {
			merged = append(merged, addresses[i])
			i++
			j++ // only keep original, skip duplicate
			continue
		}
	}

	// Append remaining elements, if any
	merged = append(merged, newAddresses[j:]...)
	merged = append(merged, addresses[i:]...)

	req.SetAddresses(merged)
}

func make2DTopics(eventSigs, topics2, topics3, topics4 []common.Hash) (res [][]common.Hash) {
	topics := [][]common.Hash{eventSigs, topics2, topics3, topics4}
	lastTopic := len(topics) - 1
	for lastTopic >= 0 && topics[lastTopic] == nil {
		lastTopic--
	}

	res = make([][]common.Hash, lastTopic+1)
	copy(res, topics)

	return res
}

// isTopicsSubset returns true if all of the sets in the list topicsA are subsets of or equal to the sets in topicsB.
// topicsA and topicsB each contain 4 (or any equal number of) sets of slices of topic values.
//
// Interpreting A & B as filters on the same contract address, "true" means that anything matching A will match B
//
// Assumptions:
// - every element of topicsA & topicsB are sorted lists containing no duplicates
func isTopicsSubset(topicsA [][]common.Hash, topicsB [][]common.Hash) bool {
	if len(topicsB) > len(topicsA) {
		return false // If topicsB requires a larger number of topics to be emitted, then B is a narrower filter than A
	}
	for i := range topicsB { // doesn't matter what topics[j] for j > len(topicsB) is, as that can only narrows filter A further
		if len(topicsB[i]) == 0 {
			continue // nil/empty list of topics matches all values, so topicsA[n] automatically a subset
		}
		if len(topicsA[i]) == 0 {
			return false // topicsA[n] matches all values, but not topicsB[n], so topicsA is not a subset
		}
		topicsMapB := make(map[common.Hash]interface{})
		for _, b := range topicsB[i] {
			topicsMapB[b] = struct{}{}
		}
		for _, a := range topicsA[i] {
			if _, ok := topicsMapB[a]; !ok {
				return false
			}
		}
	}
	return true
}

func makeEventsTopicsKey(filter Filter) string {
	// eventsTopicsKey is constructed to uniquely identify the particular combination of
	// eventSigs and topics sets a filter has. Because we don't want the key to depend on
	// the order of eventSigs, or the order of topic values for a specific topic index, we
	// must make sure these 4 lists are sorted in the same way

	size := len(filter.EventSigs[0])*(len(filter.EventSigs)+len(filter.Topic2)+len(filter.Topic3)+len(filter.Topic4)) + 4
	var eventsTopicsKey = make([]byte, 0, size)

	appendHashes := func(hashes []common.Hash) {
		for _, h := range hashes {
			eventsTopicsKey = append(eventsTopicsKey, h[:]...)
		}
		eventsTopicsKey = append(eventsTopicsKey, 0xFF) // separator
	}
	appendHashes(filter.EventSigs)
	appendHashes(filter.Topic2)
	appendHashes(filter.Topic3)
	appendHashes(filter.Topic4)
	return hex.EncodeToString(eventsTopicsKey)
}

func compareBlockNumbers(n, m *big.Int) (cmp int) {
	if n != nil && m != nil {
		return int(m.Uint64() - n.Uint64())
	}
	if n == nil {
		cmp--
	}
	if m == nil {
		cmp++
	}
	return cmp
}

// Exposes ethGetLogsReqs to tests, casting the results to []GetLogBatchElem and sorting them to make
// the output more predictable and convenient for assertions
func (lp *logPoller) EthGetLogsReqs(fromBlock, toBlock *big.Int, blockHash *common.Hash) []GetLogsBatchElem {
	rawReqs := lp.ethGetLogsReqs(fromBlock, toBlock, blockHash)
	reqs := make([]GetLogsBatchElem, len(rawReqs))
	for i := range rawReqs {
		reqs[i] = GetLogsBatchElem(rawReqs[i])
	}

	slices.SortStableFunc(reqs, func(a, b GetLogsBatchElem) int {
		nilA, nilB := a.BlockHash() == nil, b.BlockHash() == nil
		if nilA && !nilB {
			return -1
		} else if !nilA && nilB {
			return 1
		}
		if !nilB && !nilA {
			if cmp := bytes.Compare(a.BlockHash()[:], b.BlockHash()[:]); cmp != 0 {
				return cmp
			}
		}

		if cmp := compareBlockNumbers(a.FromBlock(), b.FromBlock()); cmp != 0 {
			return cmp
		}

		if cmp := compareBlockNumbers(a.ToBlock(), b.ToBlock()); cmp != 0 {
			return cmp
		}

		addressesA, addressesB := a.Addresses(), b.Addresses()
		if len(addressesA) != len(addressesB) {
			return len(addressesA) - len(addressesB)
		}
		for i := range addressesA {
			if cmp := bytes.Compare(addressesA[i][:], addressesB[i][:]); cmp != 0 {
				return cmp
			}
		}

		topicsA, topicsB := a.Topics(), b.Topics()
		if len(topicsA) != len(topicsB) { // should both be 4, but may as well handle more general case
			return len(topicsA) - len(topicsB)
		}
		for i := range topicsA {
			if len(topicsA[i]) != len(topicsB[i]) {
				return len(topicsA[i]) - len(topicsB[i])
			}
			for j := range topicsA[i] {
				if cmp := bytes.Compare(topicsA[i][j][:], topicsB[i][j][:]); cmp != 0 {
					return cmp
				}
			}
		}
		return 0 // a and b are identical
	})
	return reqs
}

// The topics passed into RegisterFilter are 2D slices backed by arrays that the
// caller could change at any time. We must have our own deep copy that's
// immutable and thread safe. Similarly, we don't want to pass our mutex-protected
// copy down the stack while sending batch requests and waiting for responses
func copyTopics(topics [][]common.Hash) (clone [][]common.Hash) {
	clone = make([][]common.Hash, len(topics))
	for i, topic := range topics {
		clone[i] = make([]common.Hash, len(topic))
		copy(clone[i], topics[i])
	}
	return clone
}

// ethGetLogsReqs generates a batched rpc reqs for all logs matching registered filters,
// copying cached reqs and filling in block range/hash if none of the registered filters have changed
func (lp *logPoller) ethGetLogsReqs(fromBlock, toBlock *big.Int, blockHash *common.Hash) []rpc.BatchElem {
	lp.filtersMu.Lock()

	if len(lp.removedFilters) != 0 || len(lp.newFilters) != 0 {
		deletedAddresses := map[common.Address]struct{}{}
		deletedEventsTopicsKeys := map[string]struct{}{}

		lp.lggr.Debugw("ethGetLogsReqs: dirty cache, rebuilding reqs indices", "removedFilters", lp.removedFilters, "newFilters", lp.newFilters)

		// First, remove any reqs corresponding to removed filters
		// Some of them we may still need, they will be rebuilt on the next pass
		for _, filter := range lp.removedFilters {
			eventsTopicsKey := makeEventsTopicsKey(filter)
			deletedEventsTopicsKeys[eventsTopicsKey] = struct{}{}
			delete(lp.cachedReqsByEventsTopicsKey, eventsTopicsKey)
			for _, address := range filter.Addresses {
				deletedAddresses[address] = struct{}{}
				delete(lp.cachedReqsByAddress, address)
			}
		}
		lp.removedFilters = nil

		// Merge/add any new filters.
		for _, filter := range lp.filters {
			var newReq *GetLogsBatchElem

			eventsTopicsKey := makeEventsTopicsKey(filter)
			_, isNew := lp.newFilters[filter.Name]

			_, hasDeletedTopics := deletedEventsTopicsKeys[eventsTopicsKey]
			var hasDeletedAddress bool
			for _, addr := range filter.Addresses {
				if _, hasDeletedAddress = deletedAddresses[addr]; hasDeletedAddress {
					break
				}
			}

			if !(isNew || hasDeletedTopics || hasDeletedAddress) {
				continue // only rebuild reqs associated with new filters or those sharing topics or addresses with a removed filter
			}

			if req, ok2 := lp.cachedReqsByEventsTopicsKey[eventsTopicsKey]; ok2 {
				// merge this filter with other filters with the same events and topics lists
				mergeAddressesIntoGetLogsReq(req, filter.Addresses)
				continue
			}

			for _, addr := range filter.Addresses {
				if reqsForAddress, ok2 := lp.cachedReqsByAddress[addr]; !ok2 {
					newReq = NewGetLogsReq(filter)
					lp.cachedReqsByEventsTopicsKey[eventsTopicsKey] = newReq
					lp.cachedReqsByAddress[addr] = []*GetLogsBatchElem{newReq}
				} else {
					newTopics := make2DTopics(filter.EventSigs, filter.Topic2, filter.Topic3, filter.Topic4)

					for i, req := range reqsForAddress {
						topics := req.Topics()
						if isTopicsSubset(newTopics, topics) {
							// Already covered by existing req
							break
						} else if isTopicsSubset(topics, newTopics) {
							// Replace existing req by new req which includes it
							reqsForAddress[i] = NewGetLogsReq(filter)
							lp.cachedReqsByAddress[addr] = reqsForAddress
							break
						}
						// Nothing similar enough found for this address, add a new req
						lp.cachedReqsByEventsTopicsKey[eventsTopicsKey] = NewGetLogsReq(filter)
						lp.cachedReqsByAddress[addr] = append(reqsForAddress, lp.cachedReqsByEventsTopicsKey[eventsTopicsKey])
					}
				}
			}
		}
		lp.newFilters = make(map[string]struct{})
	}
	lp.filtersMu.Unlock()

	blockParams := map[string]interface{}{}
	if blockHash != nil {
		blockParams["blockHash"] = blockHash
	}
	if fromBlock != nil {
		blockParams["fromBlock"] = rpc.BlockNumber(fromBlock.Uint64()).String()
	}
	if toBlock != nil {
		blockParams["toBlock"] = rpc.BlockNumber(toBlock.Uint64()).String()
	}

	// Fill fromBlock, toBlock, & blockHash while deep-copying cached reqs into a result array
	reqs := make([]rpc.BatchElem, 0, len(lp.cachedReqsByEventsTopicsKey))
	for _, req := range lp.cachedReqsByEventsTopicsKey {
		addresses := make([]common.Address, len(req.Addresses()))
		copy(addresses, req.Addresses())
		topics := copyTopics(req.Topics())

		params := maps.Clone(blockParams)
		params["address"] = addresses
		params["topics"] = topics

		reqs = append(reqs, rpc.BatchElem{
			Method: req.Method,
			Args:   []interface{}{params},
			Result: new([]types.Log),
		})
	}

	return reqs
}

// batchFetchLogs fetches logs for either a single block by block hash, or by block range,
// rebuilding the cached reqs if necessary, sending them to the rpc server, and parsing the results
// Requests for different filters are sent in parallel batches. For block range requests, the
// block range is also broken up into serial batches
func (lp *logPoller) batchFetchLogs(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, blockHash *common.Hash) ([]types.Log, error) {
	reqs := lp.ethGetLogsReqs(fromBlock, toBlock, blockHash)

	lp.lggr.Debugw("batchFetchLogs: sending batched requests", "rpcBatchSize", lp.rpcBatchSize, "numReqs", len(reqs))
	if err := lp.sendBatchedRequests(ctx, lp.rpcBatchSize, reqs); err != nil {
		return nil, err
	}

	var logs []types.Log
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
		res, ok := req.Result.(*[]types.Log)
		if !ok {
			return nil, fmt.Errorf("expected result type %T from eth_getLogs request, got %T", res, req.Result)
		}
		logs = append(logs, *res...)
	}
	return logs, nil
}

func (lp *logPoller) sendBatchedRequests(ctx context.Context, batchSize int64, reqs []rpc.BatchElem) error {
	for i := 0; i < len(reqs); i += int(batchSize) {
		j := i + int(batchSize)
		if j > len(reqs) {
			j = len(reqs)
		}

		err := lp.ec.BatchCallContext(ctx, reqs[i:j])
		if err != nil {
			return err
		}
	}
	return nil
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

	err := lp.sendBatchedRequests(ctx, batchSize, reqs)
	if err != nil {
		return nil, err
	}

	var blocks = make([]*evmtypes.Head, 0, len(reqs))
	for _, r := range reqs {
		if r.Error != nil {
			return nil, r.Error
		}
		block, is := r.Result.(*evmtypes.Head)

		if !is {
			return nil, fmt.Errorf("expected result to be a %T, got %T", &evmtypes.Head{}, r.Result)
		}
		if block == nil {
			return nil, errors.New("invariant violation: got nil block")
		}
		if block.Hash == (common.Hash{}) {
			return nil, fmt.Errorf("missing block hash for block number: %d", block.Number)
		}
		if block.Number < 0 {
			return nil, fmt.Errorf("expected block number to be >= to 0, got %d", block.Number)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// IndexedLogsWithSigsExcluding returns the set difference(A-B) of logs with signature sigA and sigB, matching is done on the topics index
//
// For example, query to retrieve unfulfilled requests by querying request log events without matching fulfillment log events.
// The order of events is not significant. Both logs must be inside the block range and have the minimum number of confirmations
func (lp *logPoller) IndexedLogsWithSigsExcluding(ctx context.Context, address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs Confirmations) ([]Log, error) {
	return lp.orm.SelectIndexedLogsWithSigsExcluding(ctx, eventSigA, eventSigB, topicIndex, address, fromBlock, toBlock, confs)
}

func EvmWord(i uint64) common.Hash {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BytesToHash(b)
}
