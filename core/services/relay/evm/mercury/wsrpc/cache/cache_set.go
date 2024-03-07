package cache

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// WsrpcCacheSet holds a set of mercury caches keyed by server URL
type WsrpcCacheSet interface {
	services.Service
	Get(ctx context.Context, client WsrpcClient) (WsrpcFetcher, error)
}

type GrpcCacheSet interface {
	services.Service
	Get(ctx context.Context, client GrpcClient) (GrpcFetcher, error)
}

type cacheSet struct {
	sync.RWMutex
	services.StateMachine

	lggr logger.Logger

	cfg Config
}

type wsrpcCacheSet struct {
	*cacheSet
	caches map[string]WsrpcCache
}

type grpcCacheSet struct {
	*cacheSet
	caches map[string]GrpcCache
}

func NewWsrpcCacheSet(lggr logger.Logger, cfg Config) WsrpcCacheSet {
	cs := newCacheSet(lggr, cfg)
	
	return &wsrpcCacheSet{
		cacheSet: cs,
		caches: make(map[string]WsrpcCache),
	}
}

func NewGrpcCacheSet(lggr logger.Logger, cfg Config) GrpcCacheSet {
	cs := newCacheSet(lggr, cfg)
	
	return &grpcCacheSet{
		cacheSet: cs,
		caches: make(map[string]GrpcCache),
	}
}

func newCacheSet(lggr logger.Logger, cfg Config) *cacheSet {
	return &cacheSet{
		sync.RWMutex{},
		services.StateMachine{},
		lggr.Named("CacheSet"),
		cfg,
	}
}

func (cs *cacheSet) Start(context.Context) error {
	return cs.StartOnce("CacheSet", func() error {
		cs.lggr.Debugw("CacheSet starting", "config", cs.cfg, "cachingEnabled", cs.cfg.LatestReportTTL > 0)
		return nil
	})
}

// func (cs *cacheSet) Start(context.Context) error {
// 	return cs.StartOnce("CacheSet", func() error {
// 		cs.lggr.Debugw("CacheSet starting", "config", cs.cfg, "cachingEnabled", cs.cfg.LatestReportTTL > 0)
// 		return nil
// 	})
// }

func (cs *grpcCacheSet) Close() error {
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

func (cs *wsrpcCacheSet) Close() error {
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

func (cs *grpcCacheSet) Get(ctx context.Context, client GrpcClient) (f GrpcFetcher, err error) {
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

func (cs *wsrpcCacheSet) Get(ctx context.Context, client WsrpcClient) (f WsrpcFetcher, err error) {
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

func (cs *grpcCacheSet) get(ctx context.Context, client GrpcClient) (GrpcFetcher, error) {
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
	c = newMemCacheGrpc(cs.lggr, client, cs.cfg)
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	cs.caches[sURL] = c
	return c, nil
}

func (cs *wsrpcCacheSet) get(ctx context.Context, client WsrpcClient) (WsrpcFetcher, error) {
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
	c = newMemCacheWsrpc(cs.lggr, client, cs.cfg)
	if err := c.Start(ctx); err != nil {
		return nil, err
	}
	cs.caches[sURL] = c
	return c, nil
}

func (cs *grpcCacheSet) HealthReport() map[string]error {
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

func (cs *wsrpcCacheSet) HealthReport() map[string]error {
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
