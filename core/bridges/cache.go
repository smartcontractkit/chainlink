package bridges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

const (
	CacheServiceName      = "BridgeCache"
	DefaultUpsertInterval = 5 * time.Second
)

type Cache struct {
	// dependencies and configurations
	ORM
	interval time.Duration

	// service state
	services.Service
	eng *services.Engine

	// data state
	bridgeTypesCache     sync.Map
	bridgeLastValueCache map[string]BridgeResponse
	mu                   sync.RWMutex
}

var _ ORM = (*Cache)(nil)
var _ services.Service = (*Cache)(nil)

func NewCache(base ORM, lggr logger.Logger, upsertInterval time.Duration) *Cache {
	c := &Cache{
		ORM:                  base,
		interval:             upsertInterval,
		bridgeLastValueCache: make(map[string]BridgeResponse),
	}
	c.Service, c.eng = services.Config{
		Name:  CacheServiceName,
		Start: c.start,
	}.NewServiceEngine(lggr)
	return c
}

func (c *Cache) WithDataSource(ds sqlutil.DataSource) ORM {
	return NewCache(NewORM(ds), c.eng, c.interval)
}

func (c *Cache) FindBridge(ctx context.Context, name BridgeName) (BridgeType, error) {
	if bridgeType, ok := c.bridgeTypesCache.Load(name); ok {
		return bridgeType.(BridgeType), nil
	}

	ormResult, err := c.ORM.FindBridge(ctx, name)
	if err == nil {
		c.bridgeTypesCache.Store(ormResult.Name, ormResult)
	}

	return ormResult, err
}

func (c *Cache) FindBridges(ctx context.Context, names []BridgeName) ([]BridgeType, error) {
	if len(names) == 0 {
		return nil, errors.New("at least one bridge name is required")
	}

	var (
		allFoundBts []BridgeType
		searchNames []BridgeName
	)

	for _, n := range names {
		if bridgeType, ok := c.bridgeTypesCache.Load(n); ok {
			allFoundBts = append(allFoundBts, bridgeType.(BridgeType))

			continue
		}

		searchNames = append(searchNames, n)
	}

	if len(allFoundBts) == len(names) {
		return allFoundBts, nil
	}

	bts, err := c.ORM.FindBridges(ctx, searchNames)
	if err != nil {
		return nil, err
	}

	for _, bt := range bts {
		c.bridgeTypesCache.Store(bt.Name, bt)
	}

	allFoundBts = append(allFoundBts, bts...)
	if len(allFoundBts) != len(names) {
		return nil, fmt.Errorf("not all bridges exist, asked for %v, exists %v", names, allFoundBts)
	}

	return allFoundBts, nil
}

func (c *Cache) DeleteBridgeType(ctx context.Context, bt *BridgeType) error {
	err := c.ORM.DeleteBridgeType(ctx, bt)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	// We delete regardless of the rows affected, in case it gets out of sync
	c.bridgeTypesCache.Delete(bt.Name)

	return err
}

func (c *Cache) BridgeTypes(ctx context.Context, offset int, limit int) ([]BridgeType, int, error) {
	return c.ORM.BridgeTypes(ctx, offset, limit)
}

func (c *Cache) CreateBridgeType(ctx context.Context, bt *BridgeType) error {
	err := c.ORM.CreateBridgeType(ctx, bt)
	if err != nil {
		return err
	}

	c.bridgeTypesCache.Store(bt.Name, *bt)

	return nil
}

func (c *Cache) UpdateBridgeType(ctx context.Context, bt *BridgeType, btr *BridgeTypeRequest) error {
	if err := c.ORM.UpdateBridgeType(ctx, bt, btr); err != nil {
		return err
	}

	c.bridgeTypesCache.Store(bt.Name, *bt)

	return nil
}

func (c *Cache) GetCachedResponse(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, error) {
	// prefer to get latest value from cache
	cached, inCache := c.latestValue(dotId, specId)
	if inCache && cached.FinishedAt.After(time.Now().Add(-maxElapsed)) {
		return cached.Value, nil
	}

	response, finishedAt, err := c.ORM.GetCachedResponseWithFinished(ctx, dotId, specId, maxElapsed)
	if err != nil {
		return nil, err
	}

	c.setValue(dotId, specId, BridgeResponse{
		DotID:      dotId,
		SpecID:     specId,
		Value:      response,
		FinishedAt: finishedAt,
	})

	return response, nil
}

func (c *Cache) UpsertBridgeResponse(ctx context.Context, dotId string, specId int32, response []byte) error {
	upsertTime := time.Now()

	// catch the rare case of a save race
	cached, inCache := c.latestValue(dotId, specId)
	if inCache && cached.FinishedAt.After(upsertTime) {
		return nil
	}

	c.setValue(dotId, specId, BridgeResponse{
		DotID:      dotId,
		SpecID:     specId,
		Value:      response,
		FinishedAt: upsertTime,
	})

	return nil
}

func (c *Cache) start(_ context.Context) error {
	ticker := services.TickerConfig{
		Initial:   c.interval,
		JitterPct: services.DefaultJitter,
	}.NewTicker(c.interval)
	c.eng.GoTick(ticker, c.doBulkUpsert)

	return nil
}

func (c *Cache) doBulkUpsert(ctx context.Context) {
	c.mu.RLock()
	values := maps.Values(c.bridgeLastValueCache)
	c.mu.RUnlock()

	if len(values) == 0 {
		return
	}

	if err := c.ORM.BulkUpsertBridgeResponse(ctx, values); err != nil {
		c.eng.Warnf("bulk upsert of bridge responses failed: %s", err.Error())
	}
}

func (c *Cache) latestValue(dotId string, specId int32) (BridgeResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, inCache := c.bridgeLastValueCache[responseKey(dotId, specId)]

	return cached, inCache
}

func (c *Cache) setValue(dotId string, specId int32, resp BridgeResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bridgeLastValueCache[responseKey(dotId, specId)] = resp
}

func responseKey(dotId string, specId int32) string {
	return fmt.Sprintf("%s||%d", dotId, specId)
}
