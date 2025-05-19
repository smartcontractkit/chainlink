package cache

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
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

	lggr   logger.SugaredLogger
	caches map[string]Cache

	cfg Config
}

func NewCacheSet(lggr logger.Logger, cfg Config) CacheSet {
	return newCacheSet(lggr, cfg)
}

func newCacheSet(lggr logger.Logger, cfg Config) *cacheSet {
	return &cacheSet{
		sync.RWMutex{},
		services.StateMachine{},
		logger.Sugared(lggr).Named("CacheSet"),
		make(map[string]Cache),
		cfg,
	}
}

func (cs *cacheSet) Start(context.Context) error {
	return cs.StartOnce("CacheSet", func() error {
		cs.lggr.Debugw("CacheSet starting", "config", cs.cfg, "cachingEnabled", cs.cfg.LatestReportTTL > 0)
		return nil
	})
}

func (cs *cacheSet) Close() error {
	return cs.StopOnce("CacheSet", func() error {
		cs.Lock()
		defer cs.Unlock()
		caches := maps.Values(cs.caches)
		if err := services.MultiCloser(caches).Close(); err != nil {
			return err
		}
		cs.caches = nil
		return nil
	})
}

func (cs *cacheSet) Get(ctx context.Context, client Client) (f Fetcher, err error) {
	if cs.cfg.LatestReportTTL == 0 {
		// caching disabled
		return nil, nil
	}
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
	c = newMemCache(cs.lggr, client, cs.cfg)
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
	caches := maps.Values(cs.caches)
	cs.RUnlock()
	for _, c := range caches {
		services.CopyHealth(report, c.HealthReport())
	}
	return report
}
func (cs *cacheSet) Name() string { return cs.lggr.Name() }
