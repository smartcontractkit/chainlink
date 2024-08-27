package llo

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"maps"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

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
	channelConfigStoreABI     abi.ABI
	topicNewChannelDefinition = (channel_config_store.ChannelConfigStoreNewChannelDefinition{}).Topic()

	allTopics = []common.Hash{topicNewChannelDefinition}
)

func init() {
	var err error
	channelConfigStoreABI, err = abi.JSON(strings.NewReader(channel_config_store.ChannelConfigStoreABI))
	if err != nil {
		panic(err)
	}
}

type ChannelDefinitionCacheORM interface {
	// TODO: What about delete/cleanup?
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
	LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *PersistedDefinitions, err error)
	StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) (err error)
}

var _ llotypes.ChannelDefinitionCache = &channelDefinitionCache{}

type LogPoller interface {
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
	LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error)
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
	chStop chan struct{}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewChannelDefinitionCache(lggr logger.Logger, orm ChannelDefinitionCacheORM, client HTTPClient, lp logpoller.LogPoller, addr common.Address, donID uint32, fromBlock int64, options ...Option) llotypes.ChannelDefinitionCache {
	filterName := logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String(), donID)
	cdc := &channelDefinitionCache{
		orm:             orm,
		client:          client,
		httpLimit:       MaxChannelDefinitionsFileSize,
		filterName:      filterName,
		lp:              lp,
		logPollInterval: defaultLogPollInterval,
		addr:            addr,
		donID:           donID,
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
		donIDTopic := common.BigToHash(big.NewInt(int64(c.donID)))
		err = c.lp.RegisterFilter(ctx, logpoller.Filter{Name: c.filterName, EventSigs: allTopics, Topic2: []common.Hash{donIDTopic}, Addresses: []common.Address{c.addr}})
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

	pollT := services.NewTicker(c.logPollInterval)
	defer pollT.Stop()

	for {
		select {
		case <-c.chStop:
			return
		case <-pollT.C:
			// failures will be tried again on the next tick
			if err := c.readLogs(); err != nil {
				c.lggr.Errorw("Failed to fetch channel definitions from chain", "err", err)
				continue
			}
		}
	}
}

func (c *channelDefinitionCache) readLogs() (err error) {
	ctx, cancel := services.StopChan(c.chStop).NewCtx()
	defer cancel()
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

	// NOTE: We assume that log poller returns logs in order of block_num, log_index ASC
	logs, err := c.lp.LogsWithSigs(ctx, fromBlock, toBlock, allTopics, c.addr)
	if err != nil {
		return err
	}

	for _, log := range logs {
		switch log.EventSig {
		case topicNewChannelDefinition:
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
				c.lggr.Warnw("Got log for unexpected donID", "donID", unpacked.DonId.String(), "expectedDonID", c.donID)
				// ignore logs for other donIDs
				// NOTE: shouldn't happen anyway since log poller filters on
				// donID
				continue
			}

			c.newLogMu.Lock()
			if c.newLog == nil || unpacked.Version > c.newLog.Version {
				// assume that donID is correct due to log poller filtering
				c.lggr.Infow("Got new channel definitions from chain", "version", unpacked.Version, "blockNumber", log.BlockNumber, "sha", fmt.Sprintf("%x", unpacked.Sha), "url", unpacked.Url)
				c.newLog = unpacked
				c.newLogCh <- unpacked
			}
			c.newLogMu.Unlock()

		default:
			// ignore unrecognized logs
			continue
		}
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

	var fetchCh chan struct{}

	for {
		select {
		case latest := <-c.newLogCh:
			// kill the old retry loop if any
			if fetchCh != nil {
				close(fetchCh)
			}

			fetchCh = make(chan struct{})

			c.wg.Add(1)
			go c.fetchLoop(fetchCh, latest)

		case <-c.chStop:
			return
		}
	}
}

func (c *channelDefinitionCache) fetchLoop(closeCh chan struct{}, log *channel_config_store.ChannelConfigStoreNewChannelDefinition) {
	defer c.wg.Done()
	b := utils.NewHTTPFetchBackoff()
	var attemptCnt int

	ctx, cancel := services.StopChan(c.chStop).NewCtx()
	defer cancel()

	err := c.fetchAndSetChannelDefinitions(ctx, log)
	if err == nil {
		c.lggr.Debugw("Set new channel definitions", "donID", c.donID, "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha))
		return
	}
	c.lggr.Warnw("Error while fetching channel definitions", "donID", c.donID, "version", log.Version, "url", log.Url, "sha", fmt.Sprintf("%x", log.Sha), "err", err, "attempt", attemptCnt)

	for {
		select {
		case <-closeCh:
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

	if memoryVersion, persistedVersion, err := c.persist(context.Background()); err != nil {
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
		bodyBytes, err := ioutil.ReadAll(body)
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

	// TODO: we could delete the old logs from logpoller here actually
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
	return
}

// Checks persisted version and tries to save if necessary on a periodic timer
// Simple backup in case database persistence fails
func (c *channelDefinitionCache) failedPersistLoop() {
	defer c.wg.Done()

	ctx, cancel := services.StopChan(c.chStop).NewCtx()
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
	// TODO: unregister filter (on job delete)?
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
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
