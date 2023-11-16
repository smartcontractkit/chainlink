package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NOTE: CacheSet holds a set of mercury caches keyed by server URL

type CacheSet interface {
	services.Service
	Get(ctx context.Context, client Client) (Fetcher, error)
}

var _ CacheSet = (*cacheSet)(nil)

type cacheSet struct {
	sync.RWMutex
	services.StateMachine

	lggr   logger.Logger
	caches map[string]Cache

	latestPriceTTL time.Duration
	maxStaleAge    time.Duration
}

func NewCacheSet(cfg Config) CacheSet {
	return newCacheSet(cfg)
}

func newCacheSet(cfg Config) *cacheSet {
	return &cacheSet{
		sync.RWMutex{},
		services.StateMachine{},
		cfg.Logger.Named("CacheSet"),
		make(map[string]Cache),
		cfg.LatestReportTTL,
		cfg.MaxStaleAge,
	}
}

func (cs *cacheSet) Start(context.Context) error {
	return cs.StartOnce("CacheSet", func() error {
		return nil
	})
}

func (cs *cacheSet) Close() error {
	return cs.StopOnce("CacheSet", func() (merr error) {
		cs.Lock()
		defer cs.Unlock()
		for _, c := range cs.caches {
			merr = errors.Join(merr, c.Close())
		}
		cs.caches = nil
		return
	})
}

func (cs *cacheSet) Get(ctx context.Context, client Client) (f Fetcher, err error) {
	ok := cs.IfStarted(func() {
		f, err = cs.get(ctx, client)
	})
	if !ok {
		return nil, fmt.Errorf("cacheSet must be started, but is: %v", cs.State())
	}
	return
}

func (cs *cacheSet) get(ctx context.Context, client Client) (Fetcher, error) {
	// HOT PATH
	cs.RLock()
	c, exists := cs.caches[client.ServerURL()]
	if exists {
		cs.RUnlock()
		return c, nil
	}
	cs.RUnlock()

	// COLD PATH
	cs.Lock()
	defer cs.Unlock()
	c, exists = cs.caches[client.ServerURL()]
	if exists {
		return c, nil
	}
	cfg := Config{
		Logger:          cs.lggr.With("serverURL", client.ServerURL()),
		LatestReportTTL: cs.latestPriceTTL,
		MaxStaleAge:     cs.maxStaleAge,
	}
	c = newMemCache(client, cfg)
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	cs.caches[client.ServerURL()] = c
	return c, nil
}

func (cs *cacheSet) HealthReport() map[string]error {
	report := map[string]error{
		cs.Name(): cs.Ready(),
	}
	cs.RLock()
	defer cs.RUnlock()
	for _, c := range cs.caches {
		services.CopyHealth(report, c.HealthReport())
	}
	return report
}
func (cs *cacheSet) Name() string { return cs.lggr.Name() }
