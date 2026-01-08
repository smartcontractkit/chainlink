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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	// dbPersistLoopInterval is the interval at which we check for failed persistence and attempt to save again
	dbPersistLoopInterval = 1 * time.Second
	// defaultFetchTimeout is the default timeout for fetching channel definitions.
	defaultFetchTimeout = 15 * time.Second
	// fetchRetryTimeout is the timeout for retrying to fetch channel definitions.
	fetchRetryTimeout = 4 * defaultFetchTimeout

	// MaxChannelsPerAdder is the maximum number of channels allowed per adder source. This limit
	// is enforced based on existing channels from the same source in currentDefinitions plus new
	// channels being added incrementally. The limit check occurs during processing, not on the
	// total file size.
	MaxChannelsPerAdder = 100

	// newChannelDefinitionEventName is the ABI event name for NewChannelDefinition events.
	newChannelDefinitionEventName = "NewChannelDefinition"
	// channelDefinitionAddedEventName is the ABI event name for ChannelDefinitionAdded events.
	channelDefinitionAddedEventName = "ChannelDefinitionAdded"

	// SourceUndefined represents an undefined channel definition source.
	SourceUndefined uint32 = 0
	// SourceOwner represents the owner source for channel definitions, which has full authority.
	SourceOwner uint32 = 1

	// SingleChannelDefinitionsFormat is the format of the channel definitions for a single source.
	SingleChannelDefinitionsFormat uint32 = 0

	// MultiChannelDefinitionsFormat is the format of the channel definitions for multiple sources.
	MultiChannelDefinitionsFormat uint32 = 1
)

var (
	// channelConfigStoreABI is the parsed ABI for the ChannelConfigStore contract.
	channelConfigStoreABI abi.ABI
	// NewChannelDefinition is the topic hash for the NewChannelDefinition event.
	NewChannelDefinition = (channel_config_store.ChannelConfigStoreNewChannelDefinition{}).Topic()
	// ChannelDefinitionAdded is the topic hash for the ChannelDefinitionAdded event.
	ChannelDefinitionAdded = (channel_config_store.ChannelConfigStoreChannelDefinitionAdded{}).Topic()
	// NoLimitSortAsc is a query configuration that sorts results by sequence in ascending order with no limit.
	NoLimitSortAsc = query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))

	channelDefinitionCacheCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "llo",
		Subsystem: "channeldefinitions",
		Name:      "channel_definition_cache_count",
		Help:      "Current count of channel definitions in the cache",
	},
		[]string{"source"},
	)
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
	StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns json.RawMessage, blockNum int64, format uint32) (err error)
	CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) error
}

var _ llotypes.ChannelDefinitionCache = &channelDefinitionCache{}

// LogPoller is an interface for querying blockchain logs. It provides methods to get the latest block,
// filter logs by expressions, and manage log filters.
type LogPoller interface {
	LatestBlock(ctx context.Context) (logpoller.Block, error)
	FilteredLogs(ctx context.Context, filter []query.Expression, limitAndSort query.LimitAndSort, queryName string) ([]logpoller.Log, error)
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
	UnregisterFilter(ctx context.Context, filterName string) error
}

// Option is a function type for configuring channelDefinitionCache options.
type Option func(*channelDefinitionCache)

// WithLogPollInterval returns an Option that sets the log polling interval for the cache.
func WithLogPollInterval(d time.Duration) Option {
	return func(c *channelDefinitionCache) {
		c.logPollInterval = d
	}
}

// Definitions holds the in-memory state of channel definitions for a channel definition cache.
// It tracks the latest block number processed, the version (for owner sources), and
// source definitions keyed by source ID.
type Definitions struct {
	LastBlockNum int64                             // The latest block number from which channel definitions were processed
	Version      uint32                            // The version number from the owner source (only updated for SourceOwner)
	Sources      map[uint32]types.SourceDefinition // Channel definitions grouped by source ID
}

// channelDefinitionCache maintains an in-memory cache of channel definitions fetched from on-chain
// events and external URLs. It polls the blockchain for new channel definition events, fetches
// definitions from URLs, verifies SHA hashes, merges definitions from multiple sources according
// to authority rules, and persists source definitions (map[uint32]types.SourceDefinition) to the database.
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

	fetchTriggerCh chan types.Trigger

	definitionsMu sync.RWMutex
	definitions   Definitions

	persistMu         sync.RWMutex
	persistedBlockNum int64

	wg     sync.WaitGroup
	chStop services.StopChan
}

// HTTPClient is an interface for making HTTP requests. It matches the standard library's
// http.Client interface.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewChannelDefinitionCache creates a new channel definition cache that monitors on-chain events
// for channel definition updates. It configures log polling filters for both owner and adder events,
// sets up the initial state, and applies any provided options. The cache must be started via Start()
// before it begins polling and fetching definitions.
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
		fetchTriggerCh:  make(chan types.Trigger, 1),
		initialBlockNum: fromBlock,
		chStop:          make(chan struct{}),
		definitions: Definitions{
			Sources: make(map[uint32]types.SourceDefinition),
		},
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
// 3. persistLoop: Periodically persists the in-memory source definitions to the database
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

		var pd *types.PersistedDefinitions
		if pd, err = c.orm.LoadChannelDefinitions(ctx, c.addr, c.donID); err != nil {
			return err
		}

		c.definitions.Sources = make(map[uint32]types.SourceDefinition)
		if pd != nil {
			if pd.Format == MultiChannelDefinitionsFormat {
				var sources map[uint32]types.SourceDefinition
				if err := json.Unmarshal(pd.Definitions, &sources); err != nil {
					return fmt.Errorf("failed to unmarshal definitions: %w", err)
				}
				c.definitions.Sources = sources
			}
			c.definitions.Version = pd.Version
			c.definitions.LastBlockNum = pd.BlockNum
			c.persistedBlockNum = pd.BlockNum
			if pd.BlockNum+1 > c.initialBlockNum {
				c.initialBlockNum = pd.BlockNum
			}
		}

		c.lggr.Infow("started channel definition cache", "definitions", c.definitions, "initialBlockNum", c.initialBlockNum, "persistedBlockNum", c.persistedBlockNum, "definitionsVersion", c.definitions.Version)

		c.wg.Add(3)
		// We have three concurrent loops
		// 1. Poll chain for new logs
		// 2. Fetch latest definitions from URL and verify SHA, according to latest log
		// 3. Persist definitions to database
		go c.pollChainLoop()
		go c.fetchLatestLoop()
		go c.persistLoop()
		return nil
	})
}

// blockNumFromUint64 converts a uint64 block number to int64.
// This is safe as block numbers are well within int64 range.
func blockNumFromUint64(blockNum uint64) int64 {
	//nolint:gosec // disable G115
	return int64(blockNum)
}

// unpackOwnerLog unpacks and validates an owner log from logpoller.
// Returns the unpacked log and an error if unpacking or validation fails.
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

// unpackAdderLog unpacks and validates an adder log from logpoller.
// Returns the unpacked log and an error if unpacking or validation fails.
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
	//nolint:gosec // disable G115
	unpacked.ChannelAdderId = uint32(new(big.Int).SetBytes(log.Topics[2]).Uint64())
	//nolint:gosec // disable G115
	unpacked.Raw.BlockNumber = uint64(log.BlockNumber)

	// Validate donID matches
	if unpacked.DonId.Cmp(big.NewInt(int64(c.donID))) != 0 {
		return nil, fmt.Errorf("donID mismatch: expected %d, got %s", c.donID, unpacked.DonId.String())
	}

	return unpacked, nil
}

// buildFilterExprs builds filter expressions by appending block range filters to base expressions.
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
// the last processed block to the latest available block. It fetches adder events
// (ChannelDefinitionAdded) and owner events (NewChannelDefinition) separately, each sorted
// individually by block number (ascending), and processes them separately by passing each
// batch to processLogs for unpacking and trigger generation.
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

	exprs := buildFilterExprs(c.adderFilterExprs, fromBlock, toBlock)
	logs, err := c.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ChannelDefinitionCachePoller - NewAdderChannelDefinition")
	if err != nil {
		return err
	}
	c.lggr.Debugw("read adder logs", "fromBlock", fromBlock, "toBlock", toBlock, "logsCount", len(logs))
	c.processLogs(logs)

	exprs = buildFilterExprs(c.ownerFilterExprs, fromBlock, toBlock)
	logs, err = c.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ChannelDefinitionCachePoller - NewOwnerChannelDefinition")
	if err != nil {
		return err
	}
	c.lggr.Debugw("read owner logs", "fromBlock", fromBlock, "toBlock", toBlock, "logsCount", len(logs))
	c.processLogs(logs)

	return nil
}

// scanFromBlockNum returns the next block number to scan from, ensuring no gaps between
// persisted and in-memory state.
// It returns the max between the in-memory definitions block number and the initial block number.
func (c *channelDefinitionCache) scanFromBlockNum() int64 {
	c.definitionsMu.RLock()
	defer c.definitionsMu.RUnlock()
	return max(c.definitions.LastBlockNum, c.initialBlockNum)
}

// processLogs unpacks channel definition logs into fetch triggers by extracting URL, SHA hash,
// block number, and source information. It validates logs and handles unpacking errors gracefully,
// continuing to process remaining logs even if individual logs fail. Valid triggers are sent to
// the fetch channel for asynchronous processing by fetchLatestLoop.
func (c *channelDefinitionCache) processLogs(logs []logpoller.Log) {
	for _, log := range logs {
		var trigger types.Trigger
		switch log.EventSig {
		case NewChannelDefinition:
			unpacked, err := c.unpackOwnerLog(log)
			if err != nil {
				// Log warning but continue processing other logs
				c.lggr.Warnw("Failed to unpack owner log", "err", err, "blockNumber", log.BlockNumber)
				continue
			}
			trigger = types.Trigger{
				Source:   SourceOwner,
				URL:      unpacked.Url,
				SHA:      unpacked.Sha,
				LogIndex: log.LogIndex,
				BlockNum: blockNumFromUint64(unpacked.Raw.BlockNumber),
				Version:  unpacked.Version,
				TxHash:   log.TxHash,
			}
		case ChannelDefinitionAdded:
			unpacked, err := c.unpackAdderLog(log)
			if err != nil {
				// Log warning but continue processing other logs
				c.lggr.Warnw("Failed to unpack adder log", "err", err, "blockNumber", log.BlockNumber)
				continue
			}
			trigger = types.Trigger{
				Source:   unpacked.ChannelAdderId,
				URL:      unpacked.Url,
				SHA:      unpacked.Sha,
				LogIndex: log.LogIndex,
				BlockNum: blockNumFromUint64(unpacked.Raw.BlockNumber),
				TxHash:   log.TxHash,
			}
		default:
			c.lggr.Warnw("Unknown log event signature",
				"blockNumber", log.BlockNumber, "eventSig", log.EventSig, "logHash", log.TxHash.Hex())
			continue
		}

		c.lggr.Debugw("Got new logs", "source", trigger.Source, "url", trigger.URL, "sha", hex.EncodeToString(trigger.SHA[:]), "blockNum", trigger.BlockNum)
		select {
		case c.fetchTriggerCh <- trigger:
		case <-c.chStop:
			return
		}
	}
}

type chOpts struct {
	FeedID common.Hash `json:"feedID"`
}

// extractFeedID attempts to extract the FeedID from channel options JSON.
// Returns the FeedID if found, or an empty hash if not found or if parsing fails.
func extractFeedID(opts llotypes.ChannelOpts) common.Hash {
	if len(opts) == 0 {
		return common.Hash{}
	}

	var optsJSON chOpts
	if err := json.Unmarshal(opts, &optsJSON); err != nil {
		// If unmarshaling fails, return empty hash (not all channel types have FeedID)
		return common.Hash{}
	}
	return optsJSON.FeedID
}

// buildFeedIDMap extracts FeedIDs from channel definitions and builds a map
// from FeedID to channel ID for collision detection.
func buildFeedIDMap(definitions llotypes.ChannelDefinitions) map[common.Hash]uint32 {
	feedIDToChannelID := make(map[common.Hash]uint32)
	for channelID, def := range definitions {
		feedID := extractFeedID(def.Opts)
		if feedID != (common.Hash{}) {
			feedIDToChannelID[feedID] = channelID
		}
	}
	return feedIDToChannelID
}

// mergeDefinitions reconciles new channel definitions with the current set according to source
// authority rules. Owner definitions (SourceOwner) have full authority: they can add, update, or
// tombstone (delete) channels. Missing channels in newDefinitions are not automatically removed;
// channels must be explicitly tombstoned to be removed. Adder definitions (non-owner sources) have
// limited authority: they can only add new channels and cannot overwrite or tombstone existing ones.
//
// Adder limits are enforced:
//   - MaxChannelsPerAdder: The limit is enforced based on existing channels from the same source
//     in currentDefinitions plus new channels being added incrementally. The check occurs before
//     each new channel addition. Existing channels that are already in currentDefinitions are
//     skipped and do not count toward new additions.
//
// FeedID uniqueness is enforced:
//   - All channels must have unique FeedIDs in their options. If a new channel has a FeedID that
//     collides with an existing channel, the new channel is logged and skipped (not added).
func (c *channelDefinitionCache) mergeDefinitions(source uint32, currentDefinitions llotypes.ChannelDefinitions, newDefinitions llotypes.ChannelDefinitions, feedIDToChannelID map[common.Hash]uint32) {
	// Count the number of channels for adder sources in the current definitions
	var numberOfChannels uint32
	if source > SourceOwner {
		for _, def := range currentDefinitions {
			if def.Source == source {
				numberOfChannels++
			}
		}
	}

	// process new definitions in a deterministic order
	channelIDs := make([]llotypes.ChannelID, 0, len(newDefinitions))
	for channelID := range newDefinitions {
		channelIDs = append(channelIDs, channelID)
	}
	sort.Slice(channelIDs, func(i, j int) bool {
		return channelIDs[i] < channelIDs[j]
	})

	for _, channelID := range channelIDs {
		def := newDefinitions[channelID]

		// Check for FeedID collision before adding the channel
		newFeedID := extractFeedID(def.Opts)
		if newFeedID != (common.Hash{}) {
			if existingChannelID, exists := feedIDToChannelID[newFeedID]; exists && existingChannelID != channelID {
				c.lggr.Warnw("feedID collision detected, skipping channel definition",
					"channelID", channelID, "feedID", newFeedID.Hex(), "existingChannelID", existingChannelID, "source", source)
				continue
			}
		}

		switch {
		case source == SourceOwner:
			currentDefinitions[channelID] = def

			// Update FeedID map after adding the channel
			if newFeedID != (common.Hash{}) {
				feedIDToChannelID[newFeedID] = channelID
			}

		case source > SourceOwner:
			if def.Tombstone {
				c.lggr.Debugw("invalid channel tombstone, cannot be added by source",
					"channelID", channelID, "source", source)
				continue
			}

			if existing, exists := currentDefinitions[channelID]; exists {
				if existing.Source != def.Source {
					c.lggr.Debugw("channel adder conflict, skipping definition",
						"channelID", channelID, "existingSourceID", existing.Source, "newSourceID", def.Source)
				}
				// Adders do not overwrite existing definitions, they can only add new ones
				continue
			}

			// stop processing new definitions if the adder limit is exceeded
			if numberOfChannels >= MaxChannelsPerAdder {
				c.lggr.Warnw("adder limit exceeded, skipping remaining definitions for source",
					"source", source, "numberOfChannels", numberOfChannels, "max", MaxChannelsPerAdder)
				return
			}

			currentDefinitions[channelID] = def
			numberOfChannels++
			// Update FeedID map after adding the channel
			if newFeedID != (common.Hash{}) {
				feedIDToChannelID[newFeedID] = channelID
			}

		default:
			c.lggr.Warnw("undefined source, skipping definition",
				"channelID", channelID, "source", source)
			continue
		}
	}
}

// fetchLatestLoop is an asynchronous goroutine that receives fetch triggers from the poll chain
// loop via a channel. It coordinates fetching channel definitions from URLs, verifying SHA hashes,
// and storing them in c.definitions.Sources (the source definitions map).
// It spawns a separate goroutine (fetchLoop) for each trigger.
func (c *channelDefinitionCache) fetchLatestLoop() {
	defer c.wg.Done()

	var trigger types.Trigger
	for {
		select {
		case trigger = <-c.fetchTriggerCh:
			if trigger.Source == SourceUndefined {
				c.lggr.Warnw("Undefined source to fetch", "url", trigger.URL, "source", trigger.Source)
				continue
			}
			c.wg.Add(1)
			go c.fetchLoop(trigger)

		case <-c.chStop:
			return
		}
	}
}

// fetchLoop is a retry goroutine spawned when an initial fetch attempt fails in fetchLatestLoop.
// It uses exponential backoff to retry fetching channel definitions until either the fetch succeeds,
// fetchRetryTimeout is reached or the cache is stopped (context cache shutdown).
// This isolates retry logic from the main fetch loop, allowing it to continue processing new triggers
// while retries occur in the background.
func (c *channelDefinitionCache) fetchLoop(trigger types.Trigger) {
	defer c.wg.Done()
	var err error
	b := utils.NewHTTPFetchBackoff()

	ctx, cancel := c.chStop.CtxWithTimeout(fetchRetryTimeout)
	defer cancel()

	if err = c.fetchAndSetChannelDefinitions(ctx, trigger); err == nil {
		return
	}
	c.lggr.Warnw("Error while fetching channel definitions", "donID",
		c.donID, "err", err, "source", trigger.Source, "attempt", b.Attempt())

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(b.Duration()):
			if err := c.fetchAndSetChannelDefinitions(ctx, trigger); err != nil {
				c.lggr.Warnw("Error while fetching channel definitions", "donID",
					c.donID, "err", err, "source", trigger.Source, "attempt", b.Attempt())
				continue
			}
			return
		}
	}
}

// fetchAndSetChannelDefinitions orchestrates fetching and storing channel definitions from a trigger.
// It checks that the trigger block number is newer than the current state to avoid processing stale
// events, fetches definitions from the URL and verifies the SHA hash, then stores them in
// c.definitions.Sources keyed by source ID. It also updates c.definitions.LastBlockNum and, for owner
// sources, c.definitions.Version. The actual merging of source definitions happens later when
// Definitions() is called.
//
// Returns an error if fetching, SHA verification, or JSON decoding fails. Note that adder limit
// checks occur during merging in Definitions(), where violations are handled by logging warnings
// and stopping processing for that source, not by returning errors.
func (c *channelDefinitionCache) fetchAndSetChannelDefinitions(ctx context.Context, trigger types.Trigger) error {
	defs, err := c.fetchChannelDefinitions(ctx, trigger)
	if err != nil {
		return fmt.Errorf("failed to fetch channel definitions: %w", err)
	}

	c.definitionsMu.Lock()
	defer c.definitionsMu.Unlock()
	if sourceDef, exists := c.definitions.Sources[trigger.Source]; exists {
		switch {
		// don't process a trigger with an earlier block number
		case trigger.BlockNum < sourceDef.Trigger.BlockNum:
			return nil

		// don't process a trigger with the same block number and an earlier log index
		case trigger.BlockNum == sourceDef.Trigger.BlockNum && trigger.LogIndex <= sourceDef.Trigger.LogIndex:
			return nil
		}
	}

	c.definitions.Sources[trigger.Source] = types.SourceDefinition{
		Trigger:     trigger,
		Definitions: defs,
	}

	if trigger.Source == SourceOwner {
		c.definitions.Version = trigger.Version
	}

	if trigger.BlockNum > c.definitions.LastBlockNum {
		c.definitions.LastBlockNum = trigger.BlockNum
	}

	c.lggr.Infow("Set channel definitions for source",
		"source", trigger.Source, "blockNum", trigger.BlockNum, "url", trigger.URL, "sha", hex.EncodeToString(trigger.SHA[:]))

	return nil
}

// fetchChannelDefinitions fetches channel definitions from the URL specified in the trigger,
// verifies the response SHA3 hash matches the expected hash from the on-chain event, decodes
// the JSON response, and annotates each definition with its source identifier. Returns an
// error if the URL is invalid, the HTTP request fails, the hash verification fails, or the
// JSON cannot be decoded.
func (c *channelDefinitionCache) fetchChannelDefinitions(ctx context.Context, trigger types.Trigger) (llotypes.ChannelDefinitions, error) {
	u, err := url.ParseRequestURI(trigger.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %s: %w", trigger.URL, err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultFetchTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request for channel definitions URL %s: %w", trigger.URL, err)
	}
	request.Header.Set("Content-Type", "application/json")

	httpRequest := clhttp.Request{
		Client:  c.client,
		Request: request,
		Config:  clhttp.RequestConfig{SizeLimit: c.httpLimit},
		Logger:  c.lggr.Named("HTTPRequest").With("url", trigger.URL, "expectedSHA", hex.EncodeToString(trigger.SHA[:])),
	}

	reader, statusCode, _, err := httpRequest.SendRequestReader()
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to channel definitions URL %s: %w", trigger.URL, err)
	}

	if statusCode >= 400 {
		// NOTE: Truncate the returned body here as we don't want to spam the
		// logs with potentially huge messages
		body := http.MaxBytesReader(nil, reader, 1024)
		defer body.Close()
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("HTTP error from channel definitions URL %s (status %d): failed to read response body: %w (partial body: %s)", trigger.URL, statusCode, err, bodyBytes)
		}
		return nil, fmt.Errorf("HTTP error from channel definitions URL %s (status %d): %s", trigger.URL, statusCode, string(bodyBytes))
	}
	defer reader.Close()

	var buf bytes.Buffer
	// Use a teeReader to avoid excessive copying
	teeReader := io.TeeReader(reader, &buf)

	hash := sha3.New256()
	// Stream the data directly into the hash and copy to buf as we go
	if _, err := io.Copy(hash, teeReader); err != nil {
		return nil, fmt.Errorf("failed to read channel definitions response body from %s: %w", trigger.URL, err)
	}

	actualSha := hash.Sum(nil)
	if !bytes.Equal(trigger.SHA[:], actualSha) {
		return nil, fmt.Errorf("SHA3 mismatch for channel definitions from %s: expected %s, got %x", trigger.URL, hex.EncodeToString(trigger.SHA[:]), actualSha)
	}

	var cd llotypes.ChannelDefinitions
	decoder := json.NewDecoder(&buf)
	if err := decoder.Decode(&cd); err != nil {
		return nil, fmt.Errorf("failed to decode channel definitions JSON from %s: %w", trigger.URL, err)
	}

	// Annotate each definition with its source identifier.
	for channelID, def := range cd {
		def.Source = trigger.Source
		cd[channelID] = def
	}

	return cd, nil
}

// persist atomically writes the in-memory source definitions (c.definitions.Sources) to the database.
// Returns the memory and persisted block numbers along with any error that occurred during persistence.
func (c *channelDefinitionCache) persist(ctx context.Context) (int64, int64, error) {
	c.persistMu.Lock()
	defer c.persistMu.Unlock()

	c.definitionsMu.RLock()
	definitions := maps.Clone(c.definitions.Sources)
	definitionsBlockNum := c.definitions.LastBlockNum
	definitionsVersion := c.definitions.Version
	c.definitionsMu.RUnlock()

	if c.persistedBlockNum >= definitionsBlockNum {
		return definitionsBlockNum, c.persistedBlockNum, nil
	}

	definitionsJSON, err := json.Marshal(definitions)
	if err != nil {
		return definitionsBlockNum, c.persistedBlockNum, fmt.Errorf("failed to marshal definitions: %w", err)
	}

	err = c.orm.StoreChannelDefinitions(ctx, c.addr, c.donID, definitionsVersion,
		definitionsJSON, definitionsBlockNum, MultiChannelDefinitionsFormat)
	if err != nil {
		return definitionsBlockNum, c.persistedBlockNum, fmt.Errorf("failed to store definitions: %w", err)
	}

	c.persistedBlockNum = definitionsBlockNum
	return definitionsBlockNum, c.persistedBlockNum, nil
}

// persistLoop is an asynchronous goroutine that periodically persists the in-memory source definitions to the database.
func (c *channelDefinitionCache) persistLoop() {
	defer c.wg.Done()
	ctx, cancel := c.chStop.NewCtx()
	defer cancel()

	for {
		select {
		case <-time.After(dbPersistLoopInterval):
			if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
				c.lggr.Warnw("Failed to persist channel definitions", "err", err, "memoryVersion", memoryVersion,
					"persistedVersion", persistedVersion)
			}
		case <-c.chStop:
			// Try one final persist with a short-ish timeout, then return
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
				c.lggr.Errorw("Failed to persist channel definitions on shutdown",
					"err", err, "memoryVersion", memoryVersion, "persistedVersion", persistedVersion)
			}
			return
		}
	}
}

// Close stops the channel definition cache by canceling all contexts, closing the stop channel,
// and waiting for all goroutines to finish. It implements the services.Service interface.
func (c *channelDefinitionCache) Close() error {
	return c.StopOnce("ChannelDefinitionCache", func() error {
		// Cancel all contexts by closing the stop channel and wait for all goroutines to finish
		close(c.chStop)
		c.wg.Wait()
		return nil
	})
}

// HealthReport returns a health report map containing the cache's health status.
// It implements the services.Service interface.
func (c *channelDefinitionCache) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.Healthy()}
	return report
}

// Name returns the name of the channel definition cache service.
// It implements the services.Service interface.
func (c *channelDefinitionCache) Name() string { return c.lggr.Name() }

// Definitions merges all source definitions stored in c.definitions.Sources with the provided previous
// outcome definitions and returns the merged result. It starts with a clone of the prev parameter,
// applying source authority rules and adder limits. If adder limit violations occur, warnings are
// logged and processing stops for that source, but processing continues with other sources. After merging all
// sources, it does not update any in-memory fields (merging is read-only). Persistence of source definitions
// happens separately via the persistLoop goroutine, not directly triggered by this method.
// This is the main method that performs the actual reconciliation of channel definitions from
// multiple sources with the previous outcome definitions.
func (c *channelDefinitionCache) Definitions(prev llotypes.ChannelDefinitions) llotypes.ChannelDefinitions {
	c.definitionsMu.RLock()
	defer c.definitionsMu.RUnlock()

	channelDefinitionCacheCount.
		WithLabelValues("previous_outcome").Set(float64(len(prev)))

	// nothing to merge
	if len(c.definitions.Sources) == 0 {
		return prev
	}

	merged := maps.Clone(prev)
	if merged == nil {
		merged = make(llotypes.ChannelDefinitions)
	}

	src := make([]types.SourceDefinition, 0, len(c.definitions.Sources))
	for _, sourceDefinition := range c.definitions.Sources {
		src = append(src, sourceDefinition)
	}

	// process definitions deterministically
	sort.Slice(src, func(i, j int) bool {
		if src[i].Trigger.BlockNum == src[j].Trigger.BlockNum {
			return src[i].Trigger.LogIndex < src[j].Trigger.LogIndex
		}
		return src[i].Trigger.BlockNum < src[j].Trigger.BlockNum
	})

	feedIDToChannelID := buildFeedIDMap(merged)
	for _, sourceDefinition := range src {
		channelDefinitionCacheCount.
			WithLabelValues(strconv.Itoa(int(sourceDefinition.Trigger.Source))).Set(float64(len(sourceDefinition.Definitions)))
		c.lggr.Debugw("merging definitions", "source", sourceDefinition.Trigger.Source)
		c.mergeDefinitions(sourceDefinition.Trigger.Source, merged, sourceDefinition.Definitions, feedIDToChannelID)
	}

	c.lggr.Debugw("returning merged definitions", "definitions", merged)
	return merged
}
