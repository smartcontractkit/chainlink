package cache

import (
	"context"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

// ChainHealthcheck checks the health of the both source and destination chain.
// Based on the values returned, CCIP can make a decision to stop or continue processing messages.
// There are four things verified here:
// 1. Source chain is healthy (this is verified by checking if source LogPoller saw finality violation)
// 2. Dest chain is healthy (this is verified by checking if destination LogPoller saw finality violation)
// 3. CommitStore is down (this is verified by checking if CommitStore is down and destination RMN is not cursed)
// 4. Source chain is cursed (this is verified by checking if source RMN is not cursed)
//
// Whenever any of the above checks fail, the chain is considered unhealthy and the CCIP should stop
// processing messages. Additionally, when the chain is unhealthy, this information is considered "sticky"
// and is cached for a certain period of time based on defaultGlobalStatusExpirationDuration.
// This may lead to some false-positives, but in this case we want to be extra cautious and avoid executing any reorged messages.
//
// Additionally, to reduce the number of calls to the RPC, we refresh RMN state in the background based on defaultRMNStateRefreshInterval
type ChainHealthcheck interface {
	job.ServiceCtx
	IsHealthy(ctx context.Context) (bool, error)
}

const (
	// RMN curse state is refreshed every 10 seconds
	defaultRMNStateRefreshInterval = 10 * time.Second
	// Whenever we mark the chain as unhealthy, we cache this information for 30 minutes
	defaultGlobalStatusExpirationDuration = 30 * time.Minute

	globalStatusKey = "globalStatus"
	rmnStatusKey    = "rmnCurseCheck"
)

type chainHealthcheck struct {
	cache                    *cache.Cache
	globalStatusKey          string
	rmnStatusKey             string
	globalStatusExpiration   time.Duration
	rmnStatusRefreshInterval time.Duration

	lggr        logger.Logger
	onRamp      ccipdata.OnRampReader
	commitStore ccipdata.CommitStoreReader

	services.StateMachine
	wg               *sync.WaitGroup
	backgroundCtx    context.Context //nolint:containedctx
	backgroundCancel context.CancelFunc
}

func NewChainHealthcheck(lggr logger.Logger, onRamp ccipdata.OnRampReader, commitStore ccipdata.CommitStoreReader) *chainHealthcheck {
	ctx, cancel := context.WithCancel(context.Background())

	ch := &chainHealthcheck{
		// Different keys use different expiration times, so we don't need to worry about the default value
		cache:                    cache.New(cache.NoExpiration, 0),
		rmnStatusKey:             rmnStatusKey,
		globalStatusKey:          globalStatusKey,
		globalStatusExpiration:   defaultGlobalStatusExpirationDuration,
		rmnStatusRefreshInterval: defaultRMNStateRefreshInterval,

		lggr:        lggr,
		onRamp:      onRamp,
		commitStore: commitStore,

		wg:               new(sync.WaitGroup),
		backgroundCtx:    ctx,
		backgroundCancel: cancel,
	}
	return ch
}

// newChainHealthcheckWithCustomEviction is used for testing purposes only. It doesn't start background worker
func newChainHealthcheckWithCustomEviction(lggr logger.Logger, onRamp ccipdata.OnRampReader, commitStore ccipdata.CommitStoreReader, globalStatusDuration time.Duration, rmnStatusRefreshInterval time.Duration) *chainHealthcheck {
	ctx, cancel := context.WithCancel(context.Background())

	return &chainHealthcheck{
		cache:                    cache.New(rmnStatusRefreshInterval, 0),
		rmnStatusKey:             rmnStatusKey,
		globalStatusKey:          globalStatusKey,
		globalStatusExpiration:   globalStatusDuration,
		rmnStatusRefreshInterval: rmnStatusRefreshInterval,

		lggr:        lggr,
		onRamp:      onRamp,
		commitStore: commitStore,

		wg:               new(sync.WaitGroup),
		backgroundCtx:    ctx,
		backgroundCancel: cancel,
	}
}

type rmnResponse struct {
	healthy bool
	err     error
}

func (c *chainHealthcheck) IsHealthy(ctx context.Context) (bool, error) {
	// Verify if flag is raised to indicate that the chain is not healthy
	// If set to false then immediately return false without checking the chain
	if cachedValue, found := c.cache.Get(c.globalStatusKey); found {
		healthy, ok := cachedValue.(bool)
		// If cached value is properly casted to bool and not healthy it means the sticky flag is raised
		// and should be returned immediately
		if !ok {
			c.lggr.Criticalw("Failed to cast cached value to sticky healthcheck", "value", cachedValue)
		} else if ok && !healthy {
			return false, nil
		}
	}

	// These checks are cheap and don't require any communication with the database or RPC
	if healthy, err := c.checkIfReadersAreHealthy(ctx); err != nil {
		return false, err
	} else if !healthy {
		c.markStickyStatusUnhealthy()
		return healthy, nil
	}

	// First call might initialize cache if it's not initialized yet. Otherwise, it will use the cached value
	if healthy, err := c.checkIfRMNsAreHealthy(ctx); err != nil {
		return false, err
	} else if !healthy {
		c.markStickyStatusUnhealthy()
		return healthy, nil
	}
	return true, nil
}

func (c *chainHealthcheck) Start(context.Context) error {
	return c.StateMachine.StartOnce("ChainHealthcheck", func() error {
		c.lggr.Info("Starting ChainHealthcheck")
		c.wg.Add(1)
		c.run()
		return nil
	})
}

func (c *chainHealthcheck) Close() error {
	return c.StateMachine.StopOnce("ChainHealthcheck", func() error {
		c.lggr.Info("Closing ChainHealthcheck")
		c.backgroundCancel()
		c.wg.Wait()
		return nil
	})
}

func (c *chainHealthcheck) run() {
	ticker := time.NewTicker(c.rmnStatusRefreshInterval)
	go func() {
		defer c.wg.Done()
		// Refresh the RMN state immediately after starting the background refresher
		_, _ = c.refresh(c.backgroundCtx)

		for {
			select {
			case <-c.backgroundCtx.Done():
				return
			case <-ticker.C:
				_, err := c.refresh(c.backgroundCtx)
				if err != nil {
					c.lggr.Errorw("Failed to refresh RMN state in the background", "err", err)
				}
			}
		}
	}()
}

func (c *chainHealthcheck) refresh(ctx context.Context) (bool, error) {
	healthy, err := c.fetchRMNCurseState(ctx)
	c.cache.Set(
		c.rmnStatusKey,
		rmnResponse{healthy, err},
		// Cache the value for 3 refresh intervals, this is just a defensive approach
		// that will enforce the RMN state to be refreshed in case of bg worker hiccup (it should never happen)
		3*c.rmnStatusRefreshInterval,
	)
	return healthy, err
}

// checkIfReadersAreHealthy checks if the source and destination chains are healthy by calling underlying LogPoller
// These calls are cheap because they don't require any communication with the database or RPC, so we don't have
// to cache the result of these calls.
func (c *chainHealthcheck) checkIfReadersAreHealthy(ctx context.Context) (bool, error) {
	sourceChainHealthy, err := c.onRamp.IsSourceChainHealthy(ctx)
	if err != nil {
		return false, errors.Wrap(err, "onRamp IsSourceChainHealthy errored")
	}

	destChainHealthy, err := c.commitStore.IsDestChainHealthy(ctx)
	if err != nil {
		return false, errors.Wrap(err, "commitStore IsDestChainHealthy errored")
	}

	if !sourceChainHealthy || !destChainHealthy {
		c.lggr.Criticalw(
			"Lane processing is stopped because source or destination chain is reported unhealthy",
			"sourceChainHealthy", sourceChainHealthy,
			"destChainHealthy", destChainHealthy,
		)
	}
	return sourceChainHealthy && destChainHealthy, nil
}

func (c *chainHealthcheck) checkIfRMNsAreHealthy(ctx context.Context) (bool, error) {
	if cachedValue, found := c.cache.Get(c.rmnStatusKey); found {
		rmn := cachedValue.(rmnResponse)
		return rmn.healthy, rmn.err
	}

	// If the value is not found in the cache, fetch the RMN curse state in a sync manner for the first time
	c.lggr.Info("Refreshing RMN state from the plugin routine, this should happen only once per lane during boot")
	return c.refresh(ctx)
}

func (c *chainHealthcheck) markStickyStatusUnhealthy() {
	c.cache.Set(c.globalStatusKey, false, c.globalStatusExpiration)
}

func (c *chainHealthcheck) fetchRMNCurseState(ctx context.Context) (bool, error) {
	var (
		eg                = new(errgroup.Group)
		isCommitStoreDown bool
		isSourceCursed    bool
	)

	eg.Go(func() error {
		var err error
		isCommitStoreDown, err = c.commitStore.IsDown(ctx)
		if err != nil {
			return errors.Wrap(err, "commitStore isDown check errored")
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		isSourceCursed, err = c.onRamp.IsSourceCursed(ctx)
		if err != nil {
			return errors.Wrap(err, "onRamp isSourceCursed errored")
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return false, err
	}

	if isCommitStoreDown || isSourceCursed {
		c.lggr.Criticalw(
			"Lane processing is stopped because source chain is cursed or CommitStore is down",
			"isCommitStoreDown", isCommitStoreDown,
			"isSourceCursed", isSourceCursed,
		)
		return false, nil
	}
	return true, nil
}
