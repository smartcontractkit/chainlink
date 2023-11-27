package cache

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// CacheSet holds a set of mercury caches keyed by server URL
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
	return cs.StopOnce("CacheSet", func() error {
		cs.Lock()
		defer cs.Unlock()
		caches := make([]io.Closer, len(cs.caches))
		var i int
		for _, c := range cs.caches {
			caches[i] = c
			i++
		}
		if err := services.MultiCloser(caches).Close(); err != nil {
			return err
		}
		cs.caches = nil
		return nil
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
	sURL := client.ServerURL()
	// HOT PATH
	cs.RLock()
	c, exists := cs.caches[sURL]
	cs.RUnlock()
	if exists {
		return c, nil
	}

	// COLD PATH
	cs.Lock()
	defer cs.Unlock()
	c, exists = cs.caches[sURL]
	if exists {
		return c, nil
	}
	cfg := Config{
		Logger:          cs.lggr.With("serverURL", sURL),
		LatestReportTTL: cs.latestPriceTTL,
		MaxStaleAge:     cs.maxStaleAge,
	}
	c = newMemCache(client, cfg)
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	cs.caches[sURL] = c
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
