package llo

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/big"
	"net/http"
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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	clhttp "github.com/smartcontractkit/chainlink/v2/core/utils/http"
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

	newChannelDefinitionEventName = "NewChannelDefinition"
)

var (
	channelConfigStoreABI abi.ABI
	NewChannelDefinition  = (channel_config_store.ChannelConfigStoreNewChannelDefinition{}).Topic()

	NoLimitSortAsc = query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))
)

func init() {
	var err error
	channelConfigStoreABI, err = abi.JSON(strings.NewReader(channel_config_store.ChannelConfigStoreABI))
	if err != nil {
		panic(err)
	}
}

type ChannelDefinitionCacheORM interface {
	LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *PersistedDefinitions, err error)
	StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) (err error)
	CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) error
}

var _ llotypes.ChannelDefinitionCache = &channelDefinitionCache{}

type LogPoller interface {
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
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

type channelDefinitionCache struct {
	services.StateMachine

	orm       ChannelDefinitionCacheORM
	client    HTTPClient
	httpLimit int64

	filterName      string
	lp              LogPoller
	logPollInterval time.Duration
	addr            common.Address
	donID           uint32
	donIDTopic      common.Hash
	filterExprs     []query.Expression
	lggr            logger.SugaredLogger
	initialBlockNum int64

	newLogMu sync.RWMutex
	newLog   *channel_config_store.ChannelConfigStoreNewChannelDefinition
	newLogCh chan *channel_config_store.ChannelConfigStoreNewChannelDefinition

	definitionsMu       sync.RWMutex
	definitions         llotypes.ChannelDefinitions
	definitionsVersion  uint32
	definitionsBlockNum int64

	persistMu        sync.RWMutex
	persistedVersion uint32

	wg     sync.WaitGroup
	chStop services.StopChan
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func filterName(addr common.Address, donID uint32) string {
	return logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String(), fmt.Sprintf("%d", donID))
}

func NewChannelDefinitionCache(lggr logger.Logger, orm ChannelDefinitionCacheORM, client HTTPClient, lp logpoller.LogPoller, addr common.Address, donID uint32, fromBlock int64, options ...Option) llotypes.ChannelDefinitionCache {
	filterName := logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String(), donID)
	donIDTopic := common.BigToHash(big.NewInt(int64(donID)))

	exprs := []query.Expression{
		logpoller.NewAddressFilter(addr),
		logpoller.NewEventSigFilter(NewChannelDefinition),
		logpoller.NewEventByTopicFilter(1, []logpoller.HashedValueComparator{
			{Values: []common.Hash{donIDTopic}, Operator: primitives.Eq},
		}),
		// NOTE: Optimize for fast pickup of new channel definitions. On
		// Arbitrum, finalization can take tens of minutes
		// (https://grafana.ops.prod.cldev.sh/d/e0453cc9-4b4a-41e1-9f01-7c21de805b39/blockchain-finality-and-gas?orgId=1&var-env=All&var-network_name=ethereum-testnet-sepolia-arbitrum-1&var-network_name=ethereum-mainnet-arbitrum-1&from=1732460992641&to=1732547392641)
		query.Confidence(primitives.Unconfirmed),
	}

	cdc := &channelDefinitionCache{
		orm:             orm,
		client:          client,
		httpLimit:       MaxChannelDefinitionsFileSize,
		filterName:      filterName,
		lp:              lp,
		logPollInterval: defaultLogPollInterval,
		addr:            addr,
		donID:           donID,
		donIDTopic:      donIDTopic,
		filterExprs:     exprs,
		lggr:            logger.Sugared(lggr).Named("ChannelDefinitionCache").With("addr", addr, "fromBlock", fromBlock),
		newLogCh:        make(chan *channel_config_store.ChannelConfigStoreNewChannelDefinition, 1),
		initialBlockNum: fromBlock,
		chStop:          make(chan struct{}),
	}
	for _, option := range options {
		option(cdc)
	}
	return cdc
}

func (c *channelDefinitionCache) Start(ctx context.Context) error {
	// Initial load from DB, then async poll from chain thereafter
	return c.StartOnce("ChannelDefinitionCache", func() (err error) {
		err = c.lp.RegisterFilter(ctx, logpoller.Filter{Name: c.filterName, EventSigs: []common.Hash{NewChannelDefinition}, Topic2: []common.Hash{c.donIDTopic}, Addresses: []common.Address{c.addr}})
		if err != nil {
			return err
		}
		if pd, err := c.orm.LoadChannelDefinitions(ctx, c.addr, c.donID); err != nil {
			return err
		} else if pd != nil {
			c.definitions = pd.Definitions
			c.initialBlockNum = pd.BlockNum + 1
			c.definitionsVersion = uint32(pd.Version)
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

////////////////////////////////////////////////////////////////////
// Log Polling
////////////////////////////////////////////////////////////////////

// pollChainLoop periodically checks logpoller for new logs
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

	exprs := make([]query.Expression, 0, len(c.filterExprs)+2)
	exprs = append(exprs, c.filterExprs...)
	exprs = append(exprs,
		query.Block(strconv.FormatInt(fromBlock, 10), primitives.Gte),
		query.Block(strconv.FormatInt(toBlock, 10), primitives.Lte),
	)

	logs, err := c.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ChannelDefinitionCachePoller - NewChannelDefinition")
	if err != nil {
		return err
	}

	for _, log := range logs {
		if log.EventSig != NewChannelDefinition {
			// ignore unrecognized logs
			continue
		}
		unpacked := new(channel_config_store.ChannelConfigStoreNewChannelDefinition)

		err := channelConfigStoreABI.UnpackIntoInterface(unpacked, newChannelDefinitionEventName, log.Data)
		if err != nil {
			return fmt.Errorf("failed to unpack log data: %w", err)
		}
		if len(log.Topics) < 2 {
			// should never happen but must guard against unexpected panics
			c.lggr.Warnw("Log missing expected topics", "log", log)
			continue
		}
		unpacked.DonId = new(big.Int).SetBytes(log.Topics[1])

		if unpacked.DonId.Cmp(big.NewInt(int64(c.donID))) != 0 {
			// skip logs for other donIDs, shouldn't happen given the
			// FilterLogs call, but belts and braces
			continue
		}

		c.newLogMu.Lock()
		if c.newLog == nil || unpacked.Version > c.newLog.Version {
			c.lggr.Infow("Got new channel definitions from chain", "version", unpacked.Version, "blockNumber", log.BlockNumber, "sha", fmt.Sprintf("%x", unpacked.Sha), "url", unpacked.Url)
			c.newLog = unpacked
			c.newLogCh <- unpacked
		}
		c.newLogMu.Unlock()
	}

	return nil
}

func (c *channelDefinitionCache) scanFromBlockNum() int64 {
	c.newLogMu.RLock()
	defer c.newLogMu.RUnlock()
	if c.newLog != nil {
		return int64(c.newLog.Raw.BlockNumber) + 1
	}
	return c.initialBlockNum
}

////////////////////////////////////////////////////////////////////
// Fetch channel definitions from URL based on latest log
////////////////////////////////////////////////////////////////////

// fetchLatestLoop waits for new logs and tries on a loop to fetch the channel definitions from the specified url
func (c *channelDefinitionCache) fetchLatestLoop() {
	defer c.wg.Done()

	var cancel context.CancelFunc = func() {}

	for {
		select {
		case latest := <-c.newLogCh:
			// kill the old retry loop if any
			cancel()

			var ctx context.Context
			ctx, cancel = context.WithCancel(context.Background())

			c.wg.Add(1)
			go c.fetchLoop(ctx, latest)

		case <-c.chStop:
			// kill the old retry loop if any
			cancel()
			return
		}
	}
}

func (c *channelDefinitionCache) fetchLoop(ctx context.Context, log *channel_config_store.ChannelConfigStoreNewChannelDefinition) {
	defer c.wg.Done()
	b := utils.NewHTTPFetchBackoff()
	var attemptCnt int

	err := c.fetchAndSetChannelDefinitions(ctx, log)
	if err == nil {
		c.lggr.Debugw("Set new channel definitions", "donID", c.donID, "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha))
		return
	}
	c.lggr.Warnw("Error while fetching channel definitions", "donID", c.donID, "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha), "err", err, "attempt", attemptCnt)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(b.Duration()):
			attemptCnt++
			err := c.fetchAndSetChannelDefinitions(ctx, log)
			if err != nil {
				c.lggr.Warnw("Error while fetching channel definitions", "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha), "err", err, "attempt", attemptCnt)
				continue
			}
			c.lggr.Debugw("Set new channel definitions", "donID", c.donID, "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha))
			return
		}
	}
}

func (c *channelDefinitionCache) fetchAndSetChannelDefinitions(ctx context.Context, log *channel_config_store.ChannelConfigStoreNewChannelDefinition) error {
	c.definitionsMu.RLock()
	if log.Version <= c.definitionsVersion {
		c.definitionsMu.RUnlock()
		return nil
	}
	c.definitionsMu.RUnlock()

	cd, err := c.fetchChannelDefinitions(ctx, log.Url, log.Sha)
	if err != nil {
		return err
	}
	c.definitionsMu.Lock()
	if log.Version <= c.definitionsVersion {
		c.definitionsMu.Unlock()
		return nil
	}
	c.definitions = cd
	c.definitionsBlockNum = int64(log.Raw.BlockNumber)
	c.definitionsVersion = log.Version
	c.definitionsMu.Unlock()

	if memoryVersion, persistedVersion, err := c.persist(ctx); err != nil {
		// If this fails, the failedPersistLoop will try again
		c.lggr.Warnw("Failed to persist channel definitions", "err", err, "memoryVersion", memoryVersion, "persistedVersion", persistedVersion)
	}

	return nil
}

func (c *channelDefinitionCache) fetchChannelDefinitions(ctx context.Context, url string, expectedSha [32]byte) (llotypes.ChannelDefinitions, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http.Request; %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	httpRequest := clhttp.HTTPRequest{
		Client:  c.client,
		Request: request,
		Config:  clhttp.HTTPRequestConfig{SizeLimit: c.httpLimit},
		Logger:  c.lggr.Named("HTTPRequest").With("url", url, "expectedSHA", fmt.Sprintf("%x", expectedSha)),
	}

	reader, statusCode, _, err := httpRequest.SendRequestReader()
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	defer reader.Close()

	if statusCode >= 400 {
		// NOTE: Truncate the returned body here as we don't want to spam the
		// logs with potentially huge messages
		body := http.MaxBytesReader(nil, reader, 1024)
		defer body.Close()
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("got error from %s: (status code: %d, error reading response body: %w, response body: %s)", url, statusCode, err, bodyBytes)
		}
		return nil, fmt.Errorf("got error from %s: (status code: %d, response body: %s)", url, statusCode, string(bodyBytes))
	}

	var buf bytes.Buffer
	// Use a teeReader to avoid excessive copying
	teeReader := io.TeeReader(reader, &buf)

	hash := sha3.New256()
	// Stream the data directly into the hash and copy to buf as we go
	if _, err := io.Copy(hash, teeReader); err != nil {
		return nil, fmt.Errorf("failed to read from body: %w", err)
	}

	actualSha := hash.Sum(nil)
	if !bytes.Equal(expectedSha[:], actualSha) {
		return nil, fmt.Errorf("SHA3 mismatch: expected %x, got %x", expectedSha, actualSha)
	}

	var cd llotypes.ChannelDefinitions
	decoder := json.NewDecoder(&buf)
	if err := decoder.Decode(&cd); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return cd, nil
}

////////////////////////////////////////////////////////////////////
// Persistence
////////////////////////////////////////////////////////////////////

func (c *channelDefinitionCache) persist(ctx context.Context) (memoryVersion, persistedVersion uint32, err error) {
	c.persistMu.RLock()
	persistedVersion = c.persistedVersion
	c.persistMu.RUnlock()

	c.definitionsMu.RLock()
	memoryVersion = c.definitionsVersion
	dfns := c.definitions
	blockNum := c.definitionsBlockNum
	c.definitionsMu.RUnlock()

	if memoryVersion <= persistedVersion {
		return
	}

	if err = c.orm.StoreChannelDefinitions(ctx, c.addr, c.donID, memoryVersion, dfns, blockNum); err != nil {
		return
	}

	c.persistMu.Lock()
	defer c.persistMu.Unlock()
	if memoryVersion > c.persistedVersion {
		persistedVersion = memoryVersion
		c.persistedVersion = persistedVersion
	}

	// NOTE: We could, in theory, delete the old logs from logpoller here since
	// they are no longer needed. But logpoller does not currently support
	// that, and in any case, the number is likely to be small so not worth
	// worrying about.
	return
}

// Checks persisted version and tries to save if necessary on a periodic timer
// Simple backup in case database persistence fails
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
			ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
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
