package ocr2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

type Tracker interface {
	Start() error
	Close() error
	poll()
}

var _ Tracker = (*contractCache)(nil)
var _ types.ContractConfigTracker = (*contractCache)(nil)

type contractCache struct {
	contractConfig  ContractConfig
	blockHeight     uint64
	ccLock          sync.RWMutex
	ccLastCheckedAt time.Time

	stop, done chan struct{}

	reader Reader
	cfg    Config
	lggr   logger.Logger
}

func NewContractCache(cfg Config, reader Reader, lggr logger.Logger) *contractCache {
	return &contractCache{
		cfg:    cfg,
		reader: reader,
		lggr:   lggr,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func (c *contractCache) updateConfig(ctx context.Context) error {
	configBlock, configDigest, err := c.reader.LatestConfigDetails(ctx)
	if err != nil {
		return fmt.Errorf("couldn't fetch latest config details: %w", err)
	}

	c.ccLock.RLock()
	isSame := c.contractConfig.ConfigBlock == configBlock && c.contractConfig.Config.ConfigDigest == configDigest
	c.ccLock.RUnlock()

	var newConfig types.ContractConfig
	if !isSame {
		newConfig, err = c.reader.LatestConfig(ctx, configBlock)
		if err != nil {
			return fmt.Errorf("couldn't fetch latest config: %w", err)
		}
	}

	blockHeight, err := c.reader.LatestBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block height: %w", err)
	}

	c.lggr.Debugw("contract cache update", "blockHeight", blockHeight, "configBlock", configBlock, "configDigest", configDigest)

	c.ccLock.Lock()
	defer c.ccLock.Unlock()
	c.ccLastCheckedAt = time.Now()
	c.blockHeight = blockHeight
	if !isSame {
		c.contractConfig = ContractConfig{
			Config:      newConfig,
			ConfigBlock: configBlock,
		}
	}

	return nil
}

func (c *contractCache) Start() error {
	ctx, cancel := utils.ContextFromChan(c.stop)
	defer cancel()
	if err := c.updateConfig(ctx); err != nil {
		c.lggr.Warnf("Failed to populate initial config: %v", err)
	}
	go c.poll()
	return nil
}

func (c *contractCache) Close() error {
	close(c.stop)
	return nil
}

func (c *contractCache) poll() {
	defer close(c.done)
	tick := time.After(0)
	for {
		select {
		case <-c.stop:
			return
		case <-tick:
			ctx, cancel := utils.ContextFromChan(c.stop)

			if err := c.updateConfig(ctx); err != nil {
				c.lggr.Errorf("Failed to update config: %v", err)
			}
			cancel()

			tick = time.After(utils.WithJitter(c.cfg.OCR2CachePollPeriod()))
		}
	}
}

func (c *contractCache) Notify() <-chan struct{} {
	return nil
}

func (c *contractCache) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	c.ccLock.RLock()
	defer c.ccLock.RUnlock()
	changedInBlock = c.contractConfig.ConfigBlock
	configDigest = c.contractConfig.Config.ConfigDigest
	err = c.assertConfigNotStale()
	return
}

func (c *contractCache) LatestConfig(ctx context.Context, changedInBlock uint64) (config types.ContractConfig, err error) {
	c.ccLock.RLock()
	defer c.ccLock.RUnlock()
	config = c.contractConfig.Config
	err = c.assertConfigNotStale()
	return
}

func (c *contractCache) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	c.ccLock.RLock()
	defer c.ccLock.RUnlock()
	return c.blockHeight, c.assertConfigNotStale()
}

func (c *contractCache) assertConfigNotStale() error {
	if c.ccLastCheckedAt.IsZero() {
		return errors.New("contract config cache not yet initialized")
	}

	if since := time.Since(c.ccLastCheckedAt); since > c.cfg.OCR2CacheTTL() {
		return fmt.Errorf("contract config cache expired: checked last %s ago", since)
	}

	return nil
}
