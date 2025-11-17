package channeldefinitions

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	clhttp "github.com/smartcontractkit/chainlink-common/pkg/http"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// MaxChannelDefinitionsFileSize is a sanity limit to avoid OOM for a
	// maliciously large file. It should be much larger than any real expected
	// channel definitions file.
	MaxChannelDefinitionsFileSize = 25 * 1024 * 1024 // 25MB
	// How often we query logpoller for new logs
	defaultLogPollInterval = 1 * time.Second
	// How often we check for failed persistence and attempt to save again
	dbPersistLoopInterval = 1 * time.Second

	// MaxChannelsPerAdder is the maximum number of channels allowed in a single adder definition
	// file. This limit is enforced on the total number of channels in the definition file before
	// any processing occurs.
	MaxChannelsPerAdder = 100
	// MaxAdderAdditionsPerDefinition is the maximum number of new channels an adder can add in a
	// single definition update. Only channels that don't already exist in the current definitions
	// are counted toward this limit.
	MaxAdderAdditionsPerDefinition = 10

	newChannelDefinitionEventName   = "NewChannelDefinition"
	channelDefinitionAddedEventName = "ChannelDefinitionAdded"

	SourceUndefined uint32 = 0
	SourceOwner     uint32 = 1
)

var (
	channelConfigStoreABI  abi.ABI
	NewChannelDefinition   = (channel_config_store.ChannelConfigStoreNewChannelDefinition{}).Topic()
	ChannelDefinitionAdded = (channel_config_store.ChannelConfigStoreChannelDefinitionAdded{}).Topic()
	NoLimitSortAsc         = query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))

	errAdderAdditionsLimitExceeded   = errors.New("adder additions per definition limit exceeded")
	errChannelsPerAdderLimitExceeded = errors.New("channels per adder limit exceeded")
)

func init() {
	var err error
	channelConfigStoreABI, err = abi.JSON(strings.NewReader(channel_config_store.ChannelConfigStoreABI))
	if err != nil {
		panic(err)
	}
}

type ChannelDefinitionCacheORM interface {
	LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *types.PersistedDefinitions, err error)
	StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) (err error)
	CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) error
}

var _ llotypes.ChannelDefinitionCache = &channelDefinitionCache{}

type LogPoller interface {
	LatestBlock(ctx context.Context) (logpoller.Block, error)
	FilteredLogs(ctx context.Context, filter []query.Expression, limitAndSort query.LimitAndSort, queryName string) ([]logpoller.Log, error)
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
	UnregisterFilter(ctx context.Context, filterName string) error
}

type Option func(*channelDefinitionCache)

func WithLogPollInterval(d time.Duration) Option {
	return func(c *channelDefinitionCache) {
		c.logPollInterval = d
	}
}

type fetchTrigger struct {
	source   uint32
	url      string
	sha      [32]byte
	blockNum int64
	version  uint32
}

type channelDefinitionCache struct {
	services.StateMachine

	orm       ChannelDefinitionCacheORM
	client    HTTPClient
	httpLimit int64

	filterName       string
	lp               LogPoller
	logPollInterval  time.Duration
	addr             common.Address
	donID            uint32
	donIDTopic       common.Hash
	ownerFilterExprs []query.Expression
	adderFilterExprs []query.Expression
	lggr             logger.SugaredLogger
	initialBlockNum  int64

	fetchMutex     sync.Mutex
	fetchTriggerCh chan fetchTrigger

	definitionsMu       sync.RWMutex
	definitions         llotypes.ChannelDefinitions
	definitionsBlockNum int64
	definitionsVersion  uint32

	persistMu         sync.RWMutex
	persistedBlockNum int64

	wg     sync.WaitGroup
	chStop services.StopChan
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewChannelDefinitionCache(lggr logger.Logger, orm ChannelDefinitionCacheORM, client HTTPClient, lp logpoller.LogPoller, addr common.Address, donID uint32, fromBlock int64, options ...Option) llotypes.ChannelDefinitionCache {

	cdc := &channelDefinitionCache{
		orm:             orm,
		client:          client,
		httpLimit:       MaxChannelDefinitionsFileSize,
		filterName:      types.ChannelDefinitionCacheFilterName(addr, donID),
		lp:              lp,
		logPollInterval: defaultLogPollInterval,
		addr:            addr,
		donID:           donID,
		donIDTopic:      common.BigToHash(big.NewInt(int64(donID))),
		lggr:            logger.Sugared(lggr).Named("ChannelDefinitionCache").With("addr", addr, "fromBlock", fromBlock),
		fetchTriggerCh:  make(chan fetchTrigger, 1),
		initialBlockNum: fromBlock,
		chStop:          make(chan struct{}),
	}

	cdc.ownerFilterExprs = []query.Expression{
		logpoller.NewAddressFilter(addr),
		logpoller.NewEventSigFilter(NewChannelDefinition),
		logpoller.NewEventByTopicFilter(1, []logpoller.HashedValueComparator{
			{Values: []common.Hash{cdc.donIDTopic}, Operator: primitives.Eq},
		}),
		// Optimize for fast pickup of new channel definitions.
		// On Arbitrum, finalization can take a long time.
		query.Confidence(primitives.Unconfirmed),
	}

	cdc.adderFilterExprs = []query.Expression{
		logpoller.NewAddressFilter(addr),
		logpoller.NewEventSigFilter(ChannelDefinitionAdded),
		logpoller.NewEventByTopicFilter(1, []logpoller.HashedValueComparator{
			{Values: []common.Hash{cdc.donIDTopic}, Operator: primitives.Eq},
		}),
		// Optimize for fast pickup of new channel definitions.
		// On Arbitrum, finalization can take a long time.
		query.Confidence(primitives.Unconfirmed),
	}

	for _, option := range options {
		option(cdc)
	}
	return cdc
}

// Start initializes the channel definition cache by loading persisted state from the database,
// registering logpoller filters, and launching three concurrent asynchronous loops:
// 1. pollChainLoop: Periodically queries logpoller for new channel definition events
// 2. fetchLatestLoop: Receives fetch triggers and coordinates fetching definitions from URLs
// 3. failedPersistLoop: Periodically retries failed database persistence operations
// All loops run until the cache is stopped via Close().
func (c *channelDefinitionCache) Start(ctx context.Context) error {
	return c.StartOnce("ChannelDefinitionCache", func() (err error) {
		err = c.lp.RegisterFilter(ctx, logpoller.Filter{
			Name:      c.filterName,
			EventSigs: []common.Hash{NewChannelDefinition, ChannelDefinitionAdded},
			Topic2:    []common.Hash{c.donIDTopic},
			Addresses: []common.Address{c.addr},
		})

		if err != nil {
			return err
		}

		if pd, err := c.orm.LoadChannelDefinitions(ctx, c.addr, c.donID); err != nil {
			return err
		} else if pd != nil {
			c.definitions = pd.Definitions
			c.definitionsVersion = pd.Version
			c.definitionsBlockNum = pd.BlockNum
			c.persistedBlockNum = pd.BlockNum
			if pd.BlockNum+1 > c.initialBlockNum {
				c.initialBlockNum = pd.BlockNum
			}
		} else {
			// ensure non-nil map ready for assignment later
			c.definitions = make(llotypes.ChannelDefinitions)
			// leave c.initialBlockNum as provided fromBlock argument
		}

		c.wg.Add(3)
		// We have three concurrent loops
		// 1. Poll chain for new logs
		// 2. Fetch latest definitions from URL and verify SHA, according to latest log
		// 3. Retry persisting records to DB, if it failed
		go c.pollChainLoop()
		go c.fetchLatestLoop()
		go c.failedPersistLoop()
		return nil
	})
}

// blockNumFromUint64 converts a uint64 block number to int64
// This is safe as block numbers are well within int64 range
func blockNumFromUint64(blockNum uint64) int64 {
	//nolint:gosec // disable G115
	return int64(blockNum)
}

// unpackOwnerLog unpacks and validates an owner log from logpoller
// Returns the unpacked log and an error if unpacking or validation fails
func (c *channelDefinitionCache) unpackOwnerLog(log logpoller.Log) (*channel_config_store.ChannelConfigStoreNewChannelDefinition, error) {
	if log.EventSig != NewChannelDefinition {
		return nil, fmt.Errorf("log event signature mismatch: expected %x, got %x", NewChannelDefinition, log.EventSig)
	}

	unpacked := new(channel_config_store.ChannelConfigStoreNewChannelDefinition)
	err := channelConfigStoreABI.UnpackIntoInterface(unpacked, newChannelDefinitionEventName, log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log data: %w", err)
	}

	if len(log.Topics) < 2 {
		return nil, fmt.Errorf("log missing expected topics: got %d, expected at least 2", len(log.Topics))
	}

	unpacked.DonId = new(big.Int).SetBytes(log.Topics[1])
	//nolint:gosec // disable G115
	unpacked.Raw.BlockNumber = uint64(log.BlockNumber)

	// Validate donID matches
	if unpacked.DonId.Cmp(big.NewInt(int64(c.donID))) != 0 {
		return nil, fmt.Errorf("donID mismatch: expected %d, got %s", c.donID, unpacked.DonId.String())
	}

	return unpacked, nil
}

// unpackAdderLog unpacks and validates an adder log from logpoller
// Returns the unpacked log and an error if unpacking or validation fails
func (c *channelDefinitionCache) unpackAdderLog(log logpoller.Log) (*channel_config_store.ChannelConfigStoreChannelDefinitionAdded, error) {
	if log.EventSig != ChannelDefinitionAdded {
		return nil, fmt.Errorf("log event signature mismatch: expected %x, got %x", ChannelDefinitionAdded, log.EventSig)
	}

	unpacked := new(channel_config_store.ChannelConfigStoreChannelDefinitionAdded)
	err := channelConfigStoreABI.UnpackIntoInterface(unpacked, channelDefinitionAddedEventName, log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack adder log data: %w", err)
	}

	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("adder log missing expected topics: got %d, expected at least 3", len(log.Topics))
	}

	unpacked.DonId = new(big.Int).SetBytes(log.Topics[1])
	unpacked.ChannelAdderId = uint32(new(big.Int).SetBytes(log.Topics[2]).Uint64())
	//nolint:gosec // disable G115
	unpacked.Raw.BlockNumber = uint64(log.BlockNumber)

	// Validate donID matches
	if unpacked.DonId.Cmp(big.NewInt(int64(c.donID))) != 0 {
		return nil, fmt.Errorf("donID mismatch: expected %d, got %s", c.donID, unpacked.DonId.String())
	}

	return unpacked, nil
}

// buildFilterExprs builds filter expressions by appending block range filters to base expressions
func buildFilterExprs(baseExprs []query.Expression, fromBlock, toBlock int64) []query.Expression {
	exprs := make([]query.Expression, 0, len(baseExprs)+2)
	exprs = append(exprs, baseExprs...)
	exprs = append(exprs,
		query.Block(strconv.FormatInt(fromBlock, 10), primitives.Gte),
		query.Block(strconv.FormatInt(toBlock, 10), primitives.Lte),
	)
	return exprs
}

// pollChainLoop is an asynchronous goroutine that periodically polls logpoller for new channel
// definition events (both owner and adder events). It processes logs sequentially by block number,
// unpacks them into fetch triggers, and sends triggers to the fetch channel for asynchronous
// processing. The loop runs until the cache is stopped, with failures logged and retried on
// the next polling interval.
func (c *channelDefinitionCache) pollChainLoop() {
	defer c.wg.Done()

	ctx, cancel := c.chStop.NewCtx()
	defer cancel()

	pollT := services.NewTicker(c.logPollInterval)
	defer pollT.Stop()

	for {
		select {
		case <-c.chStop:
			return
		case <-pollT.C:
			// failures will be tried again on the next tick
			if err := c.readLogs(ctx); err != nil {
				c.lggr.Errorw("Failed to fetch channel definitions from chain", "err", err)
				continue
			}
		}
	}
}

// readLogs queries logpoller for new channel definition events within the block range from
// the last processed block to the latest available block. It fetches both owner events
// (NewChannelDefinition) and adder events (ChannelDefinitionAdded), sorts all logs by block
// number to ensure sequential processing, and passes them to processLogs for unpacking and
// trigger generation.
func (c *channelDefinitionCache) readLogs(ctx context.Context) (err error) {
	latestBlock, err := c.lp.LatestBlock(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		c.lggr.Debug("Logpoller has no logs yet, skipping poll")
		return nil
	} else if err != nil {
		return err
	}

	toBlock := latestBlock.BlockNumber
	fromBlock := c.scanFromBlockNum()
	if toBlock <= fromBlock {
		return nil
	}

	logsToProcess := make([]logpoller.Log, 0)

	exprs := buildFilterExprs(c.adderFilterExprs, fromBlock, toBlock)
	logs, err := c.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ChannelDefinitionCachePoller - NewAdderChannelDefinition")
	if err != nil {
		return err
	}
	logsToProcess = append(logsToProcess, logs...)

	exprs = buildFilterExprs(c.ownerFilterExprs, fromBlock, toBlock)
	logs, err = c.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ChannelDefinitionCachePoller - NewOwnerChannelDefinition")
	if err != nil {
		return err
	}
	logsToProcess = append(logsToProcess, logs...)

	sort.Slice(logsToProcess, func(i, j int) bool {
		return logsToProcess[i].BlockNumber < logsToProcess[j].BlockNumber
	})

	c.processLogs(logsToProcess)

	return nil
}

// scanFromBlockNum returns the next block number to scan from, ensuring no gaps between
// persisted and in-memory state. It uses the maximum of the in-memory definitions block number
// and the initial block number to prevent re-scanning blocks that have already been processed.
func (c *channelDefinitionCache) scanFromBlockNum() int64 {
	c.definitionsMu.RLock()
	blockNum := c.definitionsBlockNum
	c.definitionsMu.RUnlock()

	if blockNum > c.initialBlockNum {
		return blockNum
	}

	return c.initialBlockNum
}

// processLogs unpacks channel definition logs into fetch triggers by extracting URL, SHA hash,
// block number, and source information. It validates logs and handles unpacking errors gracefully,
// continuing to process remaining logs even if individual logs fail. Valid triggers are sent to
// the fetch channel for asynchronous processing by fetchLatestLoop.
func (c *channelDefinitionCache) processLogs(logs []logpoller.Log) {
	for _, log := range logs {
		var trigger fetchTrigger
		switch log.EventSig {
		case NewChannelDefinition:
			unpacked, err := c.unpackOwnerLog(log)
			if err != nil {
				// Log warning but continue processing other logs
				c.lggr.Warnw("Failed to unpack owner log", "err", err, "blockNumber", log.BlockNumber)
				continue
			}
			trigger = fetchTrigger{
				source:   SourceOwner,
				url:      unpacked.Url,
				sha:      unpacked.Sha,
				blockNum: blockNumFromUint64(unpacked.Raw.BlockNumber),
				version:  unpacked.Version,
			}
		case ChannelDefinitionAdded:
			unpacked, err := c.unpackAdderLog(log)
			if err != nil {
				// Log warning but continue processing other logs
				c.lggr.Warnw("Failed to unpack adder log", "err", err, "blockNumber", log.BlockNumber)
				continue
			}
			trigger = fetchTrigger{
				source:   unpacked.AdderId,
				url:      unpacked.Url,
				sha:      unpacked.Sha,
				blockNum: blockNumFromUint64(unpacked.Raw.BlockNumber),
			}
		default:
			c.lggr.Warnw("Unknown log event signature",
				"blockNumber", log.BlockNumber, "eventSig", log.EventSig, "logHash", log.TxHash.Hex())
			continue
		}
		c.lggr.Infow("Got new logs", "source", trigger.source, "url", trigger.url, "sha", hex.EncodeToString(trigger.sha[:]), "blockNum", trigger.blockNum)
		c.fetchTriggerCh <- trigger
	}
}

// mergeDefinitions reconciles new channel definitions with the current set according to source
// authority rules. Owner definitions (SourceOwner) have full authority: they can add, update, or
// tombstone (delete) channels. Adder definitions (non-owner sources) have limited authority: they
// can only add new channels and cannot overwrite or tombstone existing ones.
//
// Adder limits are enforced:
//   - MaxChannelsPerAdder: The total number of channels in a single adder definition file cannot
//     exceed this limit. This is checked before processing any channels.
//   - MaxAdderAdditionsPerDefinition: The number of new channels an adder can add in a single
//     definition is limited. Only channels that don't already exist are counted toward this limit.
//     Existing channels from the same adder are skipped and not counted.
//
// Returns an error if adder limits are exceeded. The error will be one of:
//   - errChannelsPerAdderLimitExceeded: When the definition file contains too many channels
//   - errAdderAdditionsLimitExceeded: When the adder tries to add too many new channels
func (c *channelDefinitionCache) mergeDefinitions(source uint32, currentDefinitions llotypes.ChannelDefinitions, newDefinitions llotypes.ChannelDefinitions) (llotypes.ChannelDefinitions, error) {
	if source > SourceOwner {
		if len(newDefinitions) > MaxChannelsPerAdder {
			return nil, fmt.Errorf("%w: %d, max %d",
				errChannelsPerAdderLimitExceeded, len(newDefinitions), MaxChannelsPerAdder)
		}
	}

	var adderAdditions int
	for channelID, def := range newDefinitions {
		switch def.Source {
		case SourceOwner:
			if def.Tombstone {
				delete(currentDefinitions, channelID)
				continue
			}
			currentDefinitions[channelID] = def

		case source:
			if def.Tombstone {
				c.lggr.Warnw("invalid channel tombstone, cannot be added by adders",
					"channelID", channelID, "adderID", def.Source)
				continue
			}
			if existing, exists := currentDefinitions[channelID]; exists {
				if existing.Source != def.Source {
					c.lggr.Warnw("channel adder conflict, skipping definition",
						"channelID", channelID, "existingSourceID", existing.Source, "newSourceID", def.Source)
				}
				// Adders do not overwrite existing definitions, they can only add new ones
				continue
			}
			adderAdditions++
			currentDefinitions[channelID] = def

		default:
			c.lggr.Warnw("undefined source, skipping definition",
				"channelID", channelID, "source", def.Source, "triggerSource", source)
			continue
		}
	}

	if source > SourceOwner && adderAdditions > MaxAdderAdditionsPerDefinition {
		return nil, fmt.Errorf("%w: %d, max %d",
			errAdderAdditionsLimitExceeded, adderAdditions, MaxAdderAdditionsPerDefinition)
	}
	return currentDefinitions, nil
}

// fetchLatestLoop is an asynchronous goroutine that receives fetch triggers from the poll chain
// loop via a channel. It coordinates fetching channel definitions from URLs, verifying SHA hashes,
// and updating the in-memory cache. If an initial fetch fails, it spawns a separate retry goroutine
// (fetchLoop) with exponential backoff. The loop manages context cancellation to ensure proper
// cleanup when the cache is stopped, canceling any in-flight fetch operations.
func (c *channelDefinitionCache) fetchLatestLoop() {
	defer c.wg.Done()

	var cancel context.CancelFunc = func() {}
	var trigger fetchTrigger
	for {
		select {
		case trigger = <-c.fetchTriggerCh:
			if trigger.source == SourceUndefined {
				c.lggr.Warnw("Undefined source to fetch", "url", trigger.url, "source", trigger.source)
				continue
			}
			cancel()

			var ctx context.Context
			ctx, cancel = c.chStop.NewCtx()

			if err := c.fetchAndSetChannelDefinitions(ctx, trigger); err != nil {
				c.lggr.Warnw("Error while fetching channel definitions", "donID", c.donID, "err", err, "source", trigger.source)
				c.wg.Add(1)
				go c.fetchLoop(ctx, trigger)
			}

		case <-c.chStop:
			cancel()
			return
		}
	}
}

// fetchLoop is a retry goroutine spawned when an initial fetch attempt fails in fetchLatestLoop.
// It uses exponential backoff to retry fetching channel definitions until either the fetch succeeds
// or the context is canceled (e.g., during cache shutdown). This isolates retry logic from the
// main fetch loop, allowing it to continue processing new triggers while retries occur in the
// background.
//
// Special handling for adder limit errors: If the error is due to adder limits being exceeded
// (errAdderAdditionsLimitExceeded or errChannelsPerAdderLimitExceeded), retries are stopped
// immediately and the block number is updated to prevent reprocessing the same trigger. These
// errors indicate a permanent configuration issue that won't be resolved by retrying but
// by submitting a new definition file.
func (c *channelDefinitionCache) fetchLoop(ctx context.Context, trigger fetchTrigger) {
	defer c.wg.Done()
	var err error
	b := utils.NewHTTPFetchBackoff()

	if err = c.fetchAndSetChannelDefinitions(ctx, trigger); err == nil {
		return
	}

	if errors.Is(err, errAdderAdditionsLimitExceeded) || errors.Is(err, errChannelsPerAdderLimitExceeded) {
		c.lggr.Errorw("Error while fetching channel definitions due to limits exceeded, stopping retries",
			"donID", c.donID, "err", err, "source", trigger.source)

		c.definitionsMu.Lock()
		c.definitionsBlockNum = trigger.blockNum
		c.definitionsMu.Unlock()
		return
	}

	c.lggr.Warnw("Error while fetching channel definitions", "donID",
		c.donID, "err", err, "source", trigger.source, "attempt", b.Attempt())

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(b.Duration()):
			if err := c.fetchAndSetChannelDefinitions(ctx, trigger); err != nil {
				if errors.Is(err, errAdderAdditionsLimitExceeded) || errors.Is(err, errChannelsPerAdderLimitExceeded) {
					c.lggr.Errorw("Error while fetching channel definitions due to limits exceeded, stopping retries",
						"donID", c.donID, "err", err, "source", trigger.source)
					return
				}
				c.lggr.Warnw("Error while fetching channel definitions", "donID",
					c.donID, "err", err, "source", trigger.source, "attempt", b.Attempt())
				continue
			}
			return
		}
	}
}

// fetchAndSetChannelDefinitions orchestrates the complete fetch-merge-update cycle for channel
// definitions. It checks that the trigger block number is newer than the current state to avoid
// processing stale events, fetches definitions from the URL and verifies the SHA hash, merges
// the new definitions with the current in-memory state using reconciliation rules (which may
// enforce adder limits), atomically updates the cache with the merged result and new block number,
// and triggers persistence to the database.
//
// Returns an error if fetching, SHA verification, JSON decoding, or merging fails. Merge errors
// include adder limit violations (errAdderAdditionsLimitExceeded, errChannelsPerAdderLimitExceeded)
// which should not be retried as they indicate permanent configuration issues.
func (c *channelDefinitionCache) fetchAndSetChannelDefinitions(ctx context.Context, trigger fetchTrigger) error {
	c.fetchMutex.Lock()
	defer c.fetchMutex.Unlock()

	c.definitionsMu.RLock()
	currentBlockNum := c.definitionsBlockNum
	currentDefinitions := maps.Clone(c.definitions)
	c.definitionsMu.RUnlock()

	if trigger.blockNum <= currentBlockNum {
		return nil
	}

	defs, err := c.fetchChannelDefinitions(ctx, trigger)
	if err != nil {
		return fmt.Errorf("failed to fetch channel definitions: %w", err)
	}

	mergedDefinitions, err := c.mergeDefinitions(trigger.source, currentDefinitions, defs)
	if err != nil {
		return err
	}

	// Update definitions with the merged result
	c.definitionsMu.Lock()

	c.definitions = mergedDefinitions
	c.definitionsBlockNum = trigger.blockNum

	// Use owner version if available, otherwise keep current version (adders don't increment version)
	if trigger.source == SourceOwner {
		c.definitionsVersion = trigger.version
	}
	c.definitionsMu.Unlock()

	c.lggr.Infow("Set channel definitions",
		"donID", c.donID, "version", trigger.version, "sha", hex.EncodeToString(trigger.sha[:]),
		"blockNum", trigger.blockNum, "url", trigger.url, "source", trigger.source)

	if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
		// If this fails, the failedPersistLoop will try again
		c.lggr.Warnw("Failed to persist channel definitions", "err", err, "memoryVersion", memoryVersion, "persistedVersion", persistedVersion)
	}

	return nil
}

// fetchChannelDefinitions fetches channel definitions from the URL specified in the trigger,
// verifies the response SHA3 hash matches the expected hash from the on-chain event, decodes
// the JSON response, and annotates each definition with its source identifier. Returns an
// error if the URL is invalid, the HTTP request fails, the hash verification fails, or the
// JSON cannot be decoded.
func (c *channelDefinitionCache) fetchChannelDefinitions(ctx context.Context, trigger fetchTrigger) (llotypes.ChannelDefinitions, error) {
	u, err := url.ParseRequestURI(trigger.url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %s: %w", trigger.url, err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request for channel definitions URL %s: %w", trigger.url, err)
	}
	request.Header.Set("Content-Type", "application/json")

	httpRequest := clhttp.Request{
		Client:  c.client,
		Request: request,
		Config:  clhttp.RequestConfig{SizeLimit: c.httpLimit},
		Logger:  c.lggr.Named("HTTPRequest").With("url", trigger.url, "expectedSHA", hex.EncodeToString(trigger.sha[:])),
	}

	reader, statusCode, _, err := httpRequest.SendRequestReader()
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to channel definitions URL %s: %w", trigger.url, err)
	}
	defer reader.Close()

	if statusCode >= 400 {
		// NOTE: Truncate the returned body here as we don't want to spam the
		// logs with potentially huge messages
		body := http.MaxBytesReader(nil, reader, 1024)
		defer body.Close()
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("HTTP error from channel definitions URL %s (status %d): failed to read response body: %w (partial body: %s)", trigger.url, statusCode, err, bodyBytes)
		}
		return nil, fmt.Errorf("HTTP error from channel definitions URL %s (status %d): %s", trigger.url, statusCode, string(bodyBytes))
	}

	var buf bytes.Buffer
	// Use a teeReader to avoid excessive copying
	teeReader := io.TeeReader(reader, &buf)

	hash := sha3.New256()
	// Stream the data directly into the hash and copy to buf as we go
	if _, err := io.Copy(hash, teeReader); err != nil {
		return nil, fmt.Errorf("failed to read channel definitions response body from %s: %w", trigger.url, err)
	}

	actualSha := hash.Sum(nil)
	if !bytes.Equal(trigger.sha[:], actualSha) {
		return nil, fmt.Errorf("SHA3 mismatch for channel definitions from %s: expected %x, got %x", trigger.url, trigger.sha, actualSha)
	}

	var cd llotypes.ChannelDefinitions
	decoder := json.NewDecoder(&buf)
	if err := decoder.Decode(&cd); err != nil {
		return nil, fmt.Errorf("failed to decode channel definitions JSON from %s: %w", trigger.url, err)
	}

	// apply source to the definitions
	for channelID, def := range cd {
		def.Source = trigger.source
		cd[channelID] = def
	}

	return cd, nil
}

// persist atomically writes the in-memory channel definitions to the database if they are newer
// than the currently persisted state. It uses persistMu to prevent concurrent write operations
// and only performs the database write if definitionsBlockNum is greater than persistedBlockNum.
// Returns the memory and persisted block numbers along with any error that occurred during
// persistence.
func (c *channelDefinitionCache) persist(ctx context.Context) (memoryBlockNum, persistedBlockNum int64, err error) {
	c.persistMu.Lock()
	defer c.persistMu.Unlock()

	c.definitionsMu.RLock()
	definitionsVersion := c.definitionsVersion
	definitions := c.definitions
	definitionsBlockNum := c.definitionsBlockNum
	c.definitionsMu.RUnlock()

	if c.persistedBlockNum >= definitionsBlockNum {
		return definitionsBlockNum, c.persistedBlockNum, nil
	}

	if err = c.orm.StoreChannelDefinitions(ctx, c.addr, c.donID, definitionsVersion, definitions, definitionsBlockNum); err != nil {
		return definitionsBlockNum, c.persistedBlockNum, err
	}

	c.persistedBlockNum = definitionsBlockNum

	// NOTE: We could, in theory, delete the old logs from logpoller here since
	// they are no longer needed. But logpoller does not currently support
	// that, and in any case, the number is likely to be small so not worth
	// worrying about.
	return definitionsBlockNum, c.persistedBlockNum, nil
}

// failedPersistLoop is an asynchronous goroutine that periodically checks if in-memory channel
// definitions need to be persisted to the database and retries any failed persistence operations.
// It runs on a fixed interval and attempts to catch up any definitions that failed to persist
// during normal operation. On shutdown, it attempts one final persist operation with a timeout
// to ensure data is not lost when the cache stops.
func (c *channelDefinitionCache) failedPersistLoop() {
	defer c.wg.Done()

	ctx, cancel := c.chStop.NewCtx()
	defer cancel()

	for {
		select {
		case <-time.After(dbPersistLoopInterval):
			if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
				c.lggr.Warnw("Failed to persist channel definitions", "err", err, "memoryVersion", memoryVersion, "persistedVersion", persistedVersion)
			}
		case <-c.chStop:
			// Try one final persist with a short-ish timeout, then return
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
				c.lggr.Errorw("Failed to persist channel definitions on shutdown", "err", err, "memoryVersion", memoryVersion, "persistedVersion", persistedVersion)
			}
			return
		}
	}
}

func (c *channelDefinitionCache) Close() error {
	return c.StopOnce("ChannelDefinitionCache", func() error {
		// Cancel all contexts but try one final persist before closing
		close(c.chStop)
		c.wg.Wait()
		return nil
	})
}

func (c *channelDefinitionCache) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.Healthy()}
	return report
}

func (c *channelDefinitionCache) Name() string { return c.lggr.Name() }

func (c *channelDefinitionCache) Definitions() llotypes.ChannelDefinitions {
	c.definitionsMu.RLock()
	defer c.definitionsMu.RUnlock()
	return maps.Clone(c.definitions)
}
