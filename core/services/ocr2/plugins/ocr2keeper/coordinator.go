package ocr2keeper

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	registry "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	plugintypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

type performLogProvider struct {
	subscriptions map[string]chan plugintypes.PerformLog
	mu            sync.RWMutex
}

func (p *performLogProvider) Subscribe() (string, chan plugintypes.PerformLog) {
	id := uuid.NewString()
	s := make(chan plugintypes.PerformLog, 100)

	p.mu.Lock()
	p.subscriptions[id] = s
	p.mu.Unlock()

	return id, s
}

func (p *performLogProvider) Unsubscribe(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch, ok := p.subscriptions[id]
	if ok {
		close(ch)
		delete(p.subscriptions, id)
	}
}

func (p *performLogProvider) Broadcast(log plugintypes.PerformLog) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, chPerformLogs := range p.subscriptions {
		go func(ch chan plugintypes.PerformLog) {
			t := time.NewTimer(1 * time.Second)
			select {
			case <-t.C:
				return
			case ch <- log:
				t.Stop()
				return
			}
		}(chPerformLogs)
	}
}

type LogCoordinator struct {
	utils.StartStopOnce
	logger          logger.Logger
	logPoller       logpoller.LogPoller
	registryAddress common.Address
	lookbackBlocks  int64
	provider        *performLogProvider
	registry        *registry.KeeperRegistry
	stop            chan struct{}
}

func NewLogCoordinator(
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
) (*LogCoordinator, error) {
	var err error

	contract, err := registry.NewKeeperRegistry(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.MergeFilter([]common.Hash{
		registry.KeeperRegistryUpkeepPerformed{}.Topic(),
	}, []common.Address{registryAddress})
	if err != nil {
		return nil, err
	}

	return &LogCoordinator{
		logger:          logger,
		logPoller:       logPoller,
		registryAddress: registryAddress,
		lookbackBlocks:  lookbackBlocks,
		provider:        &performLogProvider{subscriptions: make(map[string]chan plugintypes.PerformLog)},
		registry:        contract,
		stop:            make(chan struct{}),
	}, nil
}

func (c *LogCoordinator) PerformLogProvider() plugintypes.PerformLogProvider {
	return c.provider
}

func (c *LogCoordinator) run() {
	defer func() {
		// gracefully handle panics by printing the stack trace and restarting
		// the process after some wait time
		if err := recover(); err != nil {
			c.logger.Errorf("log coordinator run routine had a panic: %s", err)

			// print the stack trace
			debug.PrintStack()

			// restart the process after a wait time
			<-time.After(1 * time.Second)

			c.logger.Errorf("log coordinator run routine is restarting")
			go c.run()
		}
	}()

	cadence := time.Second
	t := time.NewTimer(cadence)
	for {
		select {
		case <-t.C:
			ctx, cancel := context.WithTimeout(context.Background(), cadence)

			end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
			if err != nil {
				c.logger.Errorf("failed to get latest block from log poller: %s", err)
				cancel()
				continue
			}

			logs, err := c.logPoller.LogsWithSigs(
				end-c.lookbackBlocks,
				end,
				[]common.Hash{
					registry.KeeperRegistryTransmitted{}.Topic(),
				},
				c.registryAddress,
				pg.WithParentCtx(ctx),
			)
			if err != nil {
				c.logger.Errorf("failed to collect logs from log poller: %s", err)
				cancel()
				continue
			}
			cancel()

			performed, err := c.unmarshalLogs(logs)
			if err != nil {
				c.logger.Errorf("failed to unmarshal logs: %s", err)
				cancel()
				continue
			}

			for _, p := range performed {
				// broadcast log to subscribers
				c.provider.Broadcast(plugintypes.PerformLog{
					Key: plugintypes.UpkeepKey(fmt.Sprintf("%d|%s", p.Raw.BlockNumber, p.Id.String())),
				})
			}

			t.Reset(cadence)
		case <-c.stop:
			t.Stop()
			return
		}
	}
}

func (c *LogCoordinator) Start(ctx context.Context) error {
	return c.StartOnce("OCR2KeepersCoordinator", func() error {
		go c.run()
		return nil
	})
}

func (c *LogCoordinator) Close() error {
	return c.StopOnce("OCR2KeepersCoordinator", func() error {
		close(c.stop)
		return nil
	})
}

func (c *LogCoordinator) unmarshalLogs(logs []logpoller.Log) ([]*registry.KeeperRegistryUpkeepPerformed, error) {
	var err error
	performed := []*registry.KeeperRegistryUpkeepPerformed{}

	for _, log := range logs {
		rawLog := toGethLog(log)
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return performed, err
		}

		switch l := abilog.(type) {
		case *registry.KeeperRegistryUpkeepPerformed:
			performed = append(performed, l)
		}
	}

	return performed, err
}

func toGethLog(lg logpoller.Log) types.Log {
	var topics []common.Hash
	for _, b := range lg.Topics {
		topics = append(topics, common.BytesToHash(b))
	}
	return types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		Topics:      topics,
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
	}
}
