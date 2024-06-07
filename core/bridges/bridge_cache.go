package bridges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	BridgeCacheServiceName   = "BridgeCache"
	DefaultUpsertInterval    = 500 * time.Millisecond
	DefaultBulkInsertTimeout = 250 * time.Millisecond
)

type BridgeCache struct {
	// dependencies and configurations
	ORM
	lggr     logger.Logger
	interval time.Duration

	// service state
	services.StateMachine
	wg   sync.WaitGroup
	stop chan struct{}

	// data state
	bridgeTypesCache     sync.Map
	bridgeLastValueCache map[string]BridgeResponse
	mu                   sync.RWMutex
}

var _ ORM = (*BridgeCache)(nil)
var _ services.Service = (*BridgeCache)(nil)

func NewBridgeCache(base ORM, lggr logger.Logger, upsertInterval time.Duration) *BridgeCache {
	return &BridgeCache{
		ORM:                  base,
		lggr:                 lggr,
		interval:             upsertInterval,
		stop:                 make(chan struct{}, 1),
		bridgeLastValueCache: make(map[string]BridgeResponse),
	}
}

func (c *BridgeCache) WithDataSource(ds sqlutil.DataSource) ORM {
	return NewBridgeCache(NewORM(ds), c.lggr, c.interval)
}

func (c *BridgeCache) FindBridge(ctx context.Context, name BridgeName) (BridgeType, error) {
	if bridgeType, ok := c.bridgeTypesCache.Load(name); ok {
		return bridgeType.(BridgeType), nil
	}

	ormResult, err := c.ORM.FindBridge(ctx, name)
	if err == nil {
		c.bridgeTypesCache.Store(ormResult.Name, ormResult)
	}

	return ormResult, err
}

func (c *BridgeCache) FindBridges(ctx context.Context, names []BridgeName) ([]BridgeType, error) {
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

func (c *BridgeCache) DeleteBridgeType(ctx context.Context, bt *BridgeType) error {
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

func (c *BridgeCache) BridgeTypes(ctx context.Context, offset int, limit int) ([]BridgeType, int, error) {
	return c.ORM.BridgeTypes(ctx, offset, limit)
}

func (c *BridgeCache) CreateBridgeType(ctx context.Context, bt *BridgeType) error {
	err := c.ORM.CreateBridgeType(ctx, bt)
	if err != nil {
		return err
	}

	c.bridgeTypesCache.Store(bt.Name, *bt)

	return nil
}

func (c *BridgeCache) UpdateBridgeType(ctx context.Context, bt *BridgeType, btr *BridgeTypeRequest) error {
	if err := c.ORM.UpdateBridgeType(ctx, bt, btr); err != nil {
		return err
	}

	c.bridgeTypesCache.Store(bt.Name, *bt)

	return nil
}

func (c *BridgeCache) GetCachedResponse(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, error) {
	// prefer to get latest value from cache
	c.mu.RLock()
	cached, inCache := c.bridgeLastValueCache[responseKey(dotId, specId)]
	c.mu.RUnlock()

	if inCache && cached.FinishedAt.After(time.Now().Add(-maxElapsed)) {
		return cached.Value, nil
	}

	response, finishedAt, err := c.ORM.GetCachedResponseWithFinished(ctx, dotId, specId, maxElapsed)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.bridgeLastValueCache[responseKey(dotId, specId)] = BridgeResponse{
		DotID:      dotId,
		SpecID:     specId,
		Value:      response,
		FinishedAt: finishedAt,
	}
	c.mu.Unlock()

	return response, nil
}

func (c *BridgeCache) UpsertBridgeResponse(ctx context.Context, dotId string, specId int32, response []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bridgeLastValueCache[responseKey(dotId, specId)] = BridgeResponse{
		DotID:      dotId,
		SpecID:     specId,
		Value:      response,
		FinishedAt: time.Now(),
	}

	return nil
}

func (c *BridgeCache) Start(context.Context) error {
	return c.StartOnce(BridgeCacheServiceName, func() error {
		c.wg.Add(1)

		go c.run()

		return nil
	})
}

func (c *BridgeCache) Close() error {
	return c.StopOnce(BridgeCacheServiceName, func() error {
		c.stop <- struct{}{}
		c.wg.Wait()

		return nil
	})
}

func (c *BridgeCache) Ready() error {
	return c.StateMachine.Ready()
}

func (c *BridgeCache) HealthReport() map[string]error {
	return map[string]error{BridgeCacheServiceName: c.Healthy()}
}

func (c *BridgeCache) Name() string {
	return BridgeCacheServiceName
}

func (c *BridgeCache) run() {
	defer c.wg.Done()

	for {
		timer := time.NewTimer(utils.WithJitter(c.interval))

		select {
		case <-timer.C:
			c.doBulkUpsert()
		case <-c.stop:
			timer.Stop()

			return
		}
	}
}

func (c *BridgeCache) doBulkUpsert() {
	values := make([]BridgeResponse, 0, len(c.bridgeLastValueCache))

	c.mu.RLock()

	for _, value := range c.bridgeLastValueCache {
		values = append(values, value)
	}

	c.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), DefaultBulkInsertTimeout)
	defer cancel()

	if err := c.ORM.BulkUpsertBridgeResponse(ctx, values); err != nil {
		c.lggr.Warnf("bulk upsert of bridge responses failed: %s", err.Error())
	}
}

func responseKey(dotId string, specId int32) string {
	return fmt.Sprintf("%s||%d", dotId, specId)
}
