package llo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ChannelDefinitionCacheORM interface {
	// TODO: What about delete/cleanup?
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
	LoadChannelDefinitions(ctx context.Context, addr common.Address) (dfns llotypes.ChannelDefinitions, blockNum int64, err error)
	StoreChannelDefinitions(ctx context.Context, addr common.Address, dfns llotypes.ChannelDefinitions, blockNum int64) (err error)
}

var channelConfigStoreABI abi.ABI

func init() {
	var err error
	channelConfigStoreABI, err = abi.JSON(strings.NewReader(channel_config_store.ChannelConfigStoreABI))
	if err != nil {
		panic(err)
	}
}

var _ llotypes.ChannelDefinitionCache = &channelDefinitionCache{}

type channelDefinitionCache struct {
	services.StateMachine

	orm ChannelDefinitionCacheORM

	filterName string
	lp         logpoller.LogPoller
	fromBlock  int64
	addr       common.Address
	lggr       logger.Logger

	definitionsMu       sync.RWMutex
	definitions         llotypes.ChannelDefinitions
	definitionsBlockNum int64

	wg     sync.WaitGroup
	chStop chan struct{}
}

var (
	topicNewChannelDefinition     = (channel_config_store.ChannelConfigStoreNewChannelDefinition{}).Topic()
	topicChannelDefinitionRemoved = (channel_config_store.ChannelConfigStoreChannelDefinitionRemoved{}).Topic()

	allTopics = []common.Hash{topicNewChannelDefinition, topicChannelDefinitionRemoved}
)

func NewChannelDefinitionCache(lggr logger.Logger, orm ChannelDefinitionCacheORM, lp logpoller.LogPoller, addr common.Address, fromBlock int64) llotypes.ChannelDefinitionCache {
	filterName := logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String())
	return &channelDefinitionCache{
		services.StateMachine{},
		orm,
		filterName,
		lp,
		0,
		addr,
		lggr.Named("ChannelDefinitionCache").With("addr", addr, "fromBlock", fromBlock),
		sync.RWMutex{},
		nil,
		fromBlock,
		sync.WaitGroup{},
		make(chan struct{}),
	}
}

func (c *channelDefinitionCache) Start(ctx context.Context) error {
	// Initial load from DB, then async poll from chain thereafter
	return c.StartOnce("ChannelDefinitionCache", func() (err error) {
		err = c.lp.RegisterFilter(ctx, logpoller.Filter{Name: c.filterName, EventSigs: allTopics, Addresses: []common.Address{c.addr}})
		if err != nil {
			return err
		}
		if definitions, definitionsBlockNum, err := c.orm.LoadChannelDefinitions(ctx, c.addr); err != nil {
			return err
		} else if definitions != nil {
			c.definitions = definitions
			c.definitionsBlockNum = definitionsBlockNum
		} else {
			// ensure non-nil map ready for assignment later
			c.definitions = make(llotypes.ChannelDefinitions)
			// leave c.definitionsBlockNum as provided fromBlock argument
		}
		c.wg.Add(1)
		go c.poll()
		return nil
	})
}

// TODO: make this configurable?
const pollInterval = 1 * time.Second

func (c *channelDefinitionCache) poll() {
	defer c.wg.Done()

	pollT := time.NewTicker(utils.WithJitter(pollInterval))

	for {
		select {
		case <-c.chStop:
			return
		case <-pollT.C:
			if n, err := c.fetchFromChain(); err != nil {
				// TODO: retry with backoff?
				// https://smartcontract-it.atlassian.net/browse/MERC-3653
				c.lggr.Errorw("Failed to fetch channel definitions from chain", "err", err)
				continue
			} else {
				if n > 0 {
					c.lggr.Infow("Updated channel definitions", "nLogs", n, "definitionsBlockNum", c.definitionsBlockNum)
				} else {
					c.lggr.Debugw("No new channel definitions", "nLogs", 0, "definitionsBlockNum", c.definitionsBlockNum)
				}
			}
		}
	}
}

func (c *channelDefinitionCache) fetchFromChain() (nLogs int, err error) {
	// TODO: Pass context
	ctx, cancel := services.StopChan(c.chStop).NewCtx()
	defer cancel()
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
	latest, err := c.lp.LatestBlock(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		c.lggr.Debug("Logpoller has no logs yet, skipping poll")
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	toBlock := latest.BlockNumber

	fromBlock := c.definitionsBlockNum

	if toBlock <= fromBlock {
		return 0, nil
	}

	// NOTE: We assume that log poller returns logs in ascending order chronologically
	logs, err := c.lp.LogsWithSigs(ctx, fromBlock, toBlock, allTopics, c.addr)
	if err != nil {
		// TODO: retry?
		// https://smartcontract-it.atlassian.net/browse/MERC-3653
		return 0, err
	}
	for _, log := range logs {
		if err = c.applyLog(log); err != nil {
			return 0, err
		}
	}

	// Use context.Background() here because we want to try to save even if we
	// are closing
	if err = c.orm.StoreChannelDefinitions(context.Background(), c.addr, c.Definitions(), toBlock); err != nil {
		return 0, err
	}

	c.definitionsBlockNum = toBlock

	return len(logs), nil
}

func (c *channelDefinitionCache) applyLog(log logpoller.Log) error {
	switch log.EventSig {
	case topicNewChannelDefinition:
		unpacked := new(channel_config_store.ChannelConfigStoreNewChannelDefinition)

		err := channelConfigStoreABI.UnpackIntoInterface(unpacked, "NewChannelDefinition", log.Data)
		if err != nil {
			return fmt.Errorf("failed to unpack log data: %w", err)
		}

		c.applyNewChannelDefinition(unpacked)
	case topicChannelDefinitionRemoved:
		unpacked := new(channel_config_store.ChannelConfigStoreChannelDefinitionRemoved)

		err := channelConfigStoreABI.UnpackIntoInterface(unpacked, "ChannelDefinitionRemoved", log.Data)
		if err != nil {
			return fmt.Errorf("failed to unpack log data: %w", err)
		}

		c.applyChannelDefinitionRemoved(unpacked)
	default:
		// don't return error here, we want to ignore unrecognized logs and
		// continue rather than interrupting the loop
		c.lggr.Errorw("Unexpected log topic", "topic", log.EventSig.Hex())
	}
	return nil
}

func (c *channelDefinitionCache) applyNewChannelDefinition(log *channel_config_store.ChannelConfigStoreNewChannelDefinition) {
	streamIDs := make([]llotypes.StreamID, len(log.ChannelDefinition.StreamIDs))
	copy(streamIDs, log.ChannelDefinition.StreamIDs)
	c.definitionsMu.Lock()
	defer c.definitionsMu.Unlock()
	c.definitions[log.ChannelId] = llotypes.ChannelDefinition{
		ReportFormat:  llotypes.ReportFormat(log.ChannelDefinition.ReportFormat),
		ChainSelector: log.ChannelDefinition.ChainSelector,
		StreamIDs:     streamIDs,
	}
}

func (c *channelDefinitionCache) applyChannelDefinitionRemoved(log *channel_config_store.ChannelConfigStoreChannelDefinitionRemoved) {
	c.definitionsMu.Lock()
	defer c.definitionsMu.Unlock()
	delete(c.definitions, log.ChannelId)
}

func (c *channelDefinitionCache) Close() error {
	// TODO: unregister filter (on job delete)?
	// https://smartcontract-it.atlassian.net/browse/MERC-3653
	return c.StopOnce("ChannelDefinitionCache", func() error {
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
