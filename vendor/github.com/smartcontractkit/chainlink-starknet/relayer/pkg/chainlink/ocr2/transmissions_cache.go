package ocr2

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ Tracker = (*transmissionsCache)(nil)
var _ median.MedianContract = (*transmissionsCache)(nil)

type transmissionsCache struct {
	transmissionDetails TransmissionDetails
	tdLock              sync.RWMutex
	tdLastCheckedAt     time.Time

	stop, done chan struct{}

	reader Reader
	cfg    Config
	lggr   logger.Logger
}

func NewTransmissionsCache(cfg Config, reader Reader, lggr logger.Logger) *transmissionsCache {
	return &transmissionsCache{
		cfg:    cfg,
		reader: reader,
		lggr:   lggr,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
		transmissionDetails: TransmissionDetails{
			LatestAnswer: big.NewInt(0), // should always return at least 0 and not nil
		},
	}
}

func (c *transmissionsCache) updateTransmission(ctx context.Context) error {
	digest, epoch, round, answer, timestamp, err := c.reader.LatestTransmissionDetails(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't fetch latest transmission details")
	}

	c.tdLock.Lock()
	defer c.tdLock.Unlock()
	c.tdLastCheckedAt = time.Now()
	c.transmissionDetails = TransmissionDetails{
		Digest:          digest,
		Epoch:           epoch,
		Round:           round,
		LatestAnswer:    answer,
		LatestTimestamp: timestamp,
	}

	c.lggr.Debugw("transmission cache update", "details", c.transmissionDetails)

	return nil
}

func (c *transmissionsCache) Start() error {
	ctx, cancel := utils.ContextFromChan(c.stop)
	defer cancel()
	if err := c.updateTransmission(ctx); err != nil {
		c.lggr.Warnf("failed to populate initial transmission details: %v", err)
	}
	go c.poll()
	return nil
}

func (c *transmissionsCache) Close() error {
	close(c.stop)
	return nil
}

func (c *transmissionsCache) poll() {
	defer close(c.done)
	tick := time.After(0)
	for {
		select {
		case <-c.stop:
			return
		case <-tick:
			ctx, cancel := utils.ContextFromChan(c.stop)

			if err := c.updateTransmission(ctx); err != nil {
				c.lggr.Errorf("Failed to update transmission: %v", err)
			}
			cancel()

			tick = time.After(utils.WithJitter(c.cfg.OCR2CachePollPeriod()))
		}
	}
}

func (c *transmissionsCache) LatestTransmissionDetails(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	c.tdLock.RLock()
	defer c.tdLock.RUnlock()
	configDigest = c.transmissionDetails.Digest
	epoch = c.transmissionDetails.Epoch
	round = c.transmissionDetails.Round
	latestAnswer = c.transmissionDetails.LatestAnswer
	latestTimestamp = c.transmissionDetails.LatestTimestamp
	err = c.assertTransmissionsNotStale()
	return
}

func (c *transmissionsCache) LatestRoundRequested(
	ctx context.Context,
	lookback time.Duration,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	c.tdLock.RLock()
	defer c.tdLock.RUnlock()
	configDigest = c.transmissionDetails.Digest
	epoch = c.transmissionDetails.Epoch
	round = c.transmissionDetails.Round
	err = c.assertTransmissionsNotStale()
	return
}

func (c *transmissionsCache) assertTransmissionsNotStale() error {
	if c.tdLastCheckedAt.IsZero() {
		return errors.New("transmissions cache not yet initialized")
	}

	if since := time.Since(c.tdLastCheckedAt); since > c.cfg.OCR2CacheTTL() {
		return fmt.Errorf("transmissions cache expired: checked last %s ago", since)
	}

	return nil
}
