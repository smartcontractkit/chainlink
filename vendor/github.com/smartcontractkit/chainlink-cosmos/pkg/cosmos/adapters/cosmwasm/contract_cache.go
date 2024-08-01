package cosmwasm

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
)

var _ median.MedianContract = (*ContractCache)(nil)

type ContractCache struct {
	cfg    config.Config
	reader *OCR2Reader
	lggr   logger.Logger

	stop, done chan struct{}

	configMu    sync.RWMutex
	configTS    time.Time
	configBlock uint64
	config      types.ContractConfig

	transMu         sync.RWMutex
	transTS         time.Time
	digest          types.ConfigDigest
	epoch           uint32
	round           uint8
	latestAnswer    *big.Int
	latestTimestamp time.Time
}

func NewContractCache(cfg config.Config, reader *OCR2Reader, lggr logger.Logger) *ContractCache {
	return &ContractCache{
		cfg:    cfg,
		reader: reader,
		lggr:   lggr,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func (cc *ContractCache) Start() error {
	// We synchronously update the config on start so that
	// when OCR starts there is config available (if possible).
	// Avoids confusing "contract has not been configured" OCR errors.
	ctx, cancel := utils.ContextFromChan(cc.stop)
	defer cancel()
	if err := cc.updateConfig(ctx); err != nil {
		cc.lggr.Warnf("failed to populate initial config: %v", err)
	}
	go cc.poll()
	return nil
}

func (cc *ContractCache) Close() error {
	close(cc.stop)
	select {
	case <-time.After(time.Second):
		// can't rely on clients to cancel
	case <-cc.done:
	}
	return nil
}

func (cc *ContractCache) poll() {
	defer close(cc.done)
	tick := time.After(0)
	for {
		select {
		case <-cc.stop:
			return
		case <-tick:
			ctx, cancel := utils.ContextFromChan(cc.stop)
			if err := cc.updateConfig(ctx); err != nil {
				cc.lggr.Errorf("Failed to update config: %v", err)
			}
			if ctx.Err() != nil { // b/c client doesn't use ctx
				return
			}
			if err := cc.updateTransmission(ctx); err != nil {
				cc.lggr.Errorf("Failed to update transmission: %v", err)
			}
			cancel()
			tick = time.After(utils.WithJitter(cc.cfg.OCR2CachePollPeriod()))
		}
	}
}

func (cc *ContractCache) updateConfig(ctx context.Context) error {
	changedInBlock, configDigest, err := cc.reader.LatestConfigDetails(ctx)
	if err != nil {
		return fmt.Errorf("fetch latest config details: %w", err)
	}
	if err = ctx.Err(); err != nil { // b/c client doesn't use ctx
		return err
	}
	now := time.Now()
	cc.configMu.Lock()
	same := cc.configBlock == changedInBlock && cc.config.ConfigDigest == configDigest
	if same {
		cc.configTS = now // refresh TTL
	}
	cc.configMu.Unlock()
	if same {
		return nil
	}
	contractConfig, err := cc.reader.LatestConfig(ctx, changedInBlock)
	if err != nil {
		return fmt.Errorf("fetch latest config, block %d %w", changedInBlock, err)
	}
	now = time.Now()
	cc.configMu.Lock()
	cc.configTS = now
	cc.configBlock = changedInBlock
	cc.config = contractConfig
	cc.configMu.Unlock()
	cc.lggr.Infof("updated config. [config %v, config block %v]",
		contractConfig, changedInBlock)
	return nil
}

func (cc *ContractCache) updateTransmission(ctx context.Context) error {
	digest, epoch, round, latestAnswer, latestTimestamp, err := cc.reader.LatestTransmissionDetails(ctx)
	if err != nil {
		return fmt.Errorf("fetch latest transmission: %w", err)
	}
	now := time.Now()
	cc.transMu.Lock()
	cc.transTS = now
	cc.digest = digest
	cc.epoch = epoch
	cc.round = round
	cc.latestAnswer = latestAnswer
	cc.latestTimestamp = latestTimestamp
	cc.transMu.Unlock()
	cc.lggr.Infof("updated transmission details. [epoch %v, round %v, answer %v, ts %v]",
		epoch, round, latestAnswer, latestTimestamp)
	return nil
}

func (cc *ContractCache) checkTS(ts time.Time) error {
	if ts.IsZero() {
		return fmt.Errorf("contract cache not yet initialized")
	} else if since := time.Since(ts); since > cc.cfg.OCR2CacheTTL() {
		return fmt.Errorf("contract cache expired: value cached %s ago", since)
	}
	return nil
}

func (cc *ContractCache) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	cc.configMu.RLock()
	ts := cc.configTS
	changedInBlock = cc.configBlock
	configDigest = cc.config.ConfigDigest
	cc.configMu.RUnlock()
	err = cc.checkTS(ts)
	return
}

func (cc *ContractCache) LatestConfig(ctx context.Context, changedInBlock uint64) (contractConfig types.ContractConfig, err error) {
	cc.configMu.RLock()
	ts := cc.configTS
	contractConfig = cc.config
	cachedBlock := cc.configBlock
	cc.configMu.RUnlock()
	err = cc.checkTS(ts)
	if err == nil && cachedBlock != changedInBlock {
		err = fmt.Errorf("failed to get config from %d: latest config in cache is from %d", changedInBlock, cachedBlock)
	}
	return
}

func (cc *ContractCache) LatestTransmissionDetails(ctx context.Context) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	cc.transMu.RLock()
	ts := cc.transTS
	configDigest = cc.digest
	epoch = cc.epoch
	round = cc.round
	latestAnswer = cc.latestAnswer
	latestTimestamp = cc.latestTimestamp
	cc.transMu.RUnlock()
	err = cc.checkTS(ts)
	return
}

// LatestRoundRequested returns the configDigest, epoch, and round from the latest
// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
// return a result if the latest such event was emitted in a block b such that
// b.timestamp < tip.timestamp - lookback.
//
// If no event is found, LatestRoundRequested should return zero values, not an error.
// An error should only be returned if an actual error occurred during execution,
// e.g. because there was an error querying the blockchain or the database.
//
// As an optimization, this function may also return zero values, if no
// RoundRequested event has been emitted after the latest NewTransmission event.
func (cc *ContractCache) LatestRoundRequested(ctx context.Context, lookback time.Duration) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	// Not supporting this feature initially, rounds are frequent enough.
	return types.ConfigDigest{}, 0, 0, nil
}
