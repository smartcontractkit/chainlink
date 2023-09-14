package llo

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/stream_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// // TODO: needs to be populated asynchronously from onchain ConfigurationStore
// type ChannelDefinitionCache interface {
//     // TODO: Would this necessarily need to be scoped by contract address?
//     Definitions() ChannelDefinitions
//     services.Service
// }

type ChannelDefinitionCacheORM interface {
	// TODO: What about delete/cleanup?
	LoadChannelDefinitions(ctx context.Context, addr common.Address) (cd commontypes.ChannelDefinitions, blockNum int64, err error)
	StoreChannelDefinitions(ctx context.Context, cd commontypes.ChannelDefinitions) (err error)
}

var streamConfigStoreABI abi.ABI

func init() {
	var err error
	streamConfigStoreABI, err = abi.JSON(strings.NewReader(stream_config_store.StreamConfigStoreABI))
	if err != nil {
		panic(err)
	}
}

var _ commontypes.ChannelDefinitionCache = &channelDefinitionCache{}

type channelDefinitionCache struct {
	services.StateMachine

	orm ChannelDefinitionCacheORM

	filterName string
	lp         logpoller.LogPoller
	fromBlock  int64
	addr       common.Address
	lggr       logger.Logger

	definitionsMu       sync.RWMutex
	definitions         commontypes.ChannelDefinitions
	definitionsBlockNum int64

	wg     sync.WaitGroup
	chStop chan struct{}
}

var (
	topicNewChannelDefinition     = (stream_config_store.StreamConfigStoreNewChannelDefinition{}).Topic()
	topicChannelDefinitionRemoved = (stream_config_store.StreamConfigStoreChannelDefinitionRemoved{}).Topic()
	topicNewProductionConfig      = (stream_config_store.StreamConfigStoreNewProductionConfig{}).Topic()
	topicNewStagingConfig         = (stream_config_store.StreamConfigStoreNewStagingConfig{}).Topic()
	topicPromoteStagingConfig     = (stream_config_store.StreamConfigStorePromoteStagingConfig{}).Topic()

	allTopics = []common.Hash{topicNewChannelDefinition, topicChannelDefinitionRemoved, topicNewProductionConfig, topicNewStagingConfig, topicPromoteStagingConfig}
)

func NewChannelDefinitionCache(lggr logger.Logger, orm ChannelDefinitionCacheORM, lp logpoller.LogPoller, addr common.Address, fromBlock int64) commontypes.ChannelDefinitionCache {
	filterName := logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String())
	return &channelDefinitionCache{
		services.StateMachine{},
		orm,
		filterName,
		lp,
		0, // TODO: fromblock needs to be loaded from DB cache somehow because we don't want to scan all logs every time we start this job
		addr,
		lggr.Named("ChannelDefinitionCache").With("addr", addr),
		sync.RWMutex{},
		make(commontypes.ChannelDefinitions),
		fromBlock,
		sync.WaitGroup{},
		make(chan struct{}),
	}
}

// TODO: Needs a way to subscribe/unsubscribe to contracts

func (c *channelDefinitionCache) Start(ctx context.Context) error {
	// TODO: Initial load, then poll
	// TODO: needs to be populated asynchronously from onchain ConfigurationStore
	return c.StartOnce("ChannelDefinitionCache", func() (err error) {
		err = c.lp.RegisterFilter(logpoller.Filter{Name: c.filterName, EventSigs: allTopics, Addresses: []common.Address{c.addr}}, pg.WithParentCtx(ctx))
		if err != nil {
			return err
		}
		c.definitions, c.definitionsBlockNum, err = c.orm.LoadChannelDefinitions(ctx, c.addr)
		if err != nil {
			return err
		}
		c.wg.Add(1)
		go c.poll()
		return nil
	})
}

// TODO: make this configurable?
const pollInterval = 5 * time.Second

func (c *channelDefinitionCache) poll() {
	defer c.wg.Done()

	pollT := time.NewTicker(utils.WithJitter(pollInterval))

	for {
		select {
		case <-c.chStop:
			return
		case <-pollT.C:
			latest, err := c.lp.LatestBlock()
			if err != nil {
				panic("TODO")
			}
			toBlock := latest.BlockNumber
			// TODO: Pass context

			fromBlock := c.definitionsBlockNum

			if toBlock <= fromBlock {
				continue
			}

			// NOTE: We assume that log poller returns logs in ascending order chronologically
			logs, err := c.lp.LogsWithSigs(fromBlock, toBlock, []common.Hash{}, c.addr)
			if err != nil {
				// TODO: retry?
				panic(err)
			}
			for _, log := range logs {
				if err := c.applyLog(log); err != nil {
					// TODO: handle errors
					panic(err)
				}
			}

			c.definitionsBlockNum = toBlock
		}
	}
}

func (c *channelDefinitionCache) applyLog(log logpoller.Log) error {
	switch log.EventSig {
	case topicNewChannelDefinition:
		unpacked := new(stream_config_store.StreamConfigStoreNewChannelDefinition)

		err := streamConfigStoreABI.UnpackIntoInterface(unpacked, "NewChannelDefinition", log.Data)
		if err != nil {
			return fmt.Errorf("failed to unpack log data: %w", err)
		}

		c.applyNewChannelDefinition(unpacked)
	case topicChannelDefinitionRemoved:
		unpacked := new(stream_config_store.StreamConfigStoreChannelDefinitionRemoved)

		err := streamConfigStoreABI.UnpackIntoInterface(unpacked, "ChannelDefinitionRemoved", log.Data)
		if err != nil {
			return fmt.Errorf("failed to unpack log data: %w", err)
		}

		c.applyChannelDefinitionRemoved(unpacked)
	default:
		panic("TODO")
	}
	return nil
}

func (c *channelDefinitionCache) applyNewChannelDefinition(log *stream_config_store.StreamConfigStoreNewChannelDefinition) {
	rf := string(log.ChannelDefinition.ReportFormat[:])
	streamIDs := make([]commontypes.StreamID, len(log.ChannelDefinition.StreamIDs))
	for i, streamID := range log.ChannelDefinition.StreamIDs {
		streamIDs[i] = commontypes.StreamID(string(streamID[:]))
	}
	c.definitionsMu.Lock()
	defer c.definitionsMu.Unlock()
	c.definitions[log.ChannelId] = commontypes.ChannelDefinition{
		ReportFormat:  commontypes.LLOReportFormat(rf),
		ChainSelector: log.ChannelDefinition.ChainSelector,
		StreamIDs:     streamIDs,
	}
}

func (c *channelDefinitionCache) applyChannelDefinitionRemoved(log *stream_config_store.StreamConfigStoreChannelDefinitionRemoved) {
	c.definitionsMu.Lock()
	defer c.definitionsMu.Unlock()
	delete(c.definitions, log.ChannelId)
}

func (c *channelDefinitionCache) Close() error {
	// TODO
	// TODO: unregister filter (on job delete)?
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

func (c *channelDefinitionCache) Definitions() commontypes.ChannelDefinitions {
	c.definitionsMu.RLock()
	defer c.definitionsMu.RUnlock()
	return maps.Clone(c.definitions)
}
