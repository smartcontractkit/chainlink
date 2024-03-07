package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
// and is cached for a certain period of time based on defaultGlobalStatusDuration.
// This may lead to some false-positives, but in this case we want to be extra cautious and avoid executing any reorged messages.
//
// Additionally, to reduce the number of calls to the RPC, we cache RMN curse state for a certain period of
// time based on defaultRmnStatusDuration.
//
//go:generate mockery --quiet --name ChainHealthcheck --filename chain_health_mock.go --case=underscore
type ChainHealthcheck interface {
	// IsHealthy checks if the chain is healthy and returns true if it is, false otherwise
	// If forceRefresh is set to true, it will refresh the RMN curse state. Should be used in the Observation and ShouldTransmit phases of OCR2.
	// Otherwise, it will use the cached value of the RMN curse state.
	IsHealthy(ctx context.Context, forceRefresh bool) (bool, error)
}

const (
	// RMN curse state is refreshed every 20 seconds or when ForceIsHealthy is called
	defaultRmnStatusDuration    = 20 * time.Second
	defaultGlobalStatusDuration = 30 * time.Minute

	globalStatusKey = "globalStatus"
	rmnStatusKey    = "rmnCurseCheck"
)

type chainHealthcheck struct {
	cache                  *cache.Cache
	globalStatusKey        string
	rmnStatusKey           string
	globalStatusExpiration time.Duration
	rmnStatusExpiration    time.Duration

	lggr        logger.Logger
	onRamp      ccipdata.OnRampReader
	commitStore ccipdata.CommitStoreReader
}

func NewChainHealthcheck(
	lggr logger.Logger,
	onRamp ccipdata.OnRampReader,
	commitStore ccipdata.CommitStoreReader,
) *chainHealthcheck {
	return &chainHealthcheck{
		cache:                  cache.New(defaultRmnStatusDuration, 0),
		globalStatusKey:        globalStatusKey,
		rmnStatusKey:           rmnStatusKey,
		globalStatusExpiration: defaultGlobalStatusDuration,
		rmnStatusExpiration:    defaultRmnStatusDuration,

		lggr:        lggr,
		onRamp:      onRamp,
		commitStore: commitStore,
	}
}

func newChainHealthcheckWithCustomEviction(
	lggr logger.Logger,
	onRamp ccipdata.OnRampReader,
	commitStore ccipdata.CommitStoreReader,
	globalStatusDuration time.Duration,
	rmnStatusDuration time.Duration,
) *chainHealthcheck {
	return &chainHealthcheck{
		cache:                  cache.New(rmnStatusDuration, 0),
		rmnStatusKey:           rmnStatusKey,
		globalStatusKey:        globalStatusKey,
		globalStatusExpiration: globalStatusDuration,
		rmnStatusExpiration:    rmnStatusDuration,

		lggr:        lggr,
		onRamp:      onRamp,
		commitStore: commitStore,
	}
}

func (c *chainHealthcheck) IsHealthy(ctx context.Context, forceRefresh bool) (bool, error) {
	// Verify if flag is raised to indicate that the chain is not healthy
	// If set to false then immediately return false without checking the chain
	if healthy, found := c.cache.Get(c.globalStatusKey); found && !healthy.(bool) {
		return false, nil
	}

	if healthy, err := c.checkIfReadersAreHealthy(ctx); err != nil {
		return false, err
	} else if !healthy {
		c.cache.Set(c.globalStatusKey, false, c.globalStatusExpiration)
		return healthy, nil
	}

	if healthy, err := c.checkIfRMNsAreHealthy(ctx, forceRefresh); err != nil {
		return false, err
	} else if !healthy {
		c.cache.Set(c.globalStatusKey, false, c.globalStatusExpiration)
		return healthy, nil
	}
	return true, nil
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

func (c *chainHealthcheck) checkIfRMNsAreHealthy(ctx context.Context, forceFetch bool) (bool, error) {
	if !forceFetch {
		if healthy, found := c.cache.Get(c.rmnStatusKey); found {
			return healthy.(bool), nil
		}
	}

	healthy, err := c.fetchRMNCurseState(ctx)
	if err != nil {
		return false, err
	}

	c.cache.Set(c.rmnStatusKey, healthy, c.rmnStatusExpiration)
	return healthy, nil
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
