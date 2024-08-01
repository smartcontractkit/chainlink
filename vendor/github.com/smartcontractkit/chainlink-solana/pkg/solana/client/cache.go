package client

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/monitor"
)

type CacheGetter[R any] func(ctx context.Context) (res R, slot uint64, err error)

// Cache is a generic implementation for caching data from the chain
type Cache[R any] struct {
	services.StateMachine

	// identifier
	metricName string
	Account    solana.PublicKey
	ChainID    string

	// stored answer
	resLock sync.RWMutex
	res     R
	resTime time.Time

	// dependencies
	getter CacheGetter[R]
	cfg    config.Config
	lggr   logger.Logger

	// polling
	done   chan struct{}
	stopCh services.StopChan
}

func NewCache[R any](metricName string, account solana.PublicKey, chainID string, cfg config.Config, getFunc CacheGetter[R], lggr logger.Logger) *Cache[R] {
	return &Cache[R]{
		metricName: metricName,
		Account:    account,
		ChainID:    chainID,
		getter:     getFunc,
		lggr:       lggr,
		cfg:        cfg,
	}
}

func (c *Cache[R]) Name() string {
	return c.lggr.Name()
}

// Start polling
func (c *Cache[R]) Start(ctx context.Context) error {
	return c.StartOnce("cache_"+c.metricName, func() error {
		c.done = make(chan struct{})
		c.stopCh = make(chan struct{})
		// We synchronously update the config on start so that
		// when OCR starts there is config available (if possible).
		// Avoids confusing "contract has not been configured" OCR errors.
		err := c.Fetch(ctx)
		if err != nil {
			c.lggr.Warnf("error in initial fetch %s", err)
		}
		go c.Poll()
		return nil
	})
}

// Close stops the polling
func (c *Cache[R]) Close() error {
	return c.StopOnce("cache_"+c.metricName, func() error {
		close(c.stopCh)
		<-c.done
		return nil
	})
}

// Poll contains the polling implementation
func (c *Cache[R]) Poll() {
	defer close(c.done)
	ctx, cancel := c.stopCh.NewCtx()
	defer cancel()
	c.lggr.Debugf("Starting polling: %s", c.Account)
	tick := time.After(0)
	for {
		select {
		case <-ctx.Done():
			c.lggr.Debugf("Stopping polling: %s", c.Account)
			return
		case <-tick:
			start := time.Now()
			err := c.Fetch(ctx)
			if err != nil {
				c.lggr.Errorf("error in Poll.fetch %s", err)
			}
			// Note negative duration will be immediately ready
			tick = time.After(utils.WithJitter(c.cfg.OCR2CachePollPeriod()) - time.Since(start))
		}
	}
}

// Read reads the latest result from memory with mutex and errors if timeout is exceeded
func (c *Cache[R]) Read() (R, error) {
	c.resLock.RLock()
	defer c.resLock.RUnlock()

	// check if stale timeout
	var err error
	if time.Since(c.resTime) > c.cfg.OCR2CacheTTL() {
		err = errors.New("error in Read: stale data, polling is likely experiencing errors")
	}
	return c.res, err
}

func (c *Cache[R]) Timestamp() time.Time {
	return c.resTime
}

func (c *Cache[R]) Fetch(ctx context.Context) error {
	c.lggr.Debugf("fetch for account: %s", c.Account)
	res, _, err := c.getter(ctx)
	if err != nil {
		return err
	}
	c.lggr.Debugf("latest fetched for account: %s, result: %v", c.Account, res)

	timestamp := time.Now()
	monitor.SetCacheTimestamp(timestamp, c.metricName, c.ChainID, c.Account.String())
	// acquire lock and write to state
	c.resLock.Lock()
	defer c.resLock.Unlock()
	c.res = res
	c.resTime = timestamp
	return nil
}
