package logprovider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

var (
	// LogRetention is the amount of time to retain logs for.
	LogRetention = 24 * time.Hour
	// LogBackfillBuffer is the number of blocks from the latest block for which backfill is done when adding a filter in log poller
	LogBackfillBuffer = 100
)

func (p *logEventProvider) RefreshActiveUpkeeps(ctx context.Context, ids ...*big.Int) ([]*big.Int, error) {
	// Exploratory: investigate how we can batch the refresh
	if len(ids) == 0 {
		return nil, nil
	}
	p.lggr.Debugw("Refreshing active upkeeps", "upkeeps", len(ids))
	visited := make(map[string]bool, len(ids))
	for _, id := range ids {
		visited[id.String()] = false
	}
	inactiveIDs := p.filterStore.GetIDs(func(f upkeepFilter) bool {
		uid := f.upkeepID.String()
		_, ok := visited[uid]
		visited[uid] = true
		return !ok
	})
	var merr error
	if len(inactiveIDs) > 0 {
		p.lggr.Debugw("Removing inactive upkeeps", "upkeeps", len(inactiveIDs))
		for _, id := range inactiveIDs {
			if err := p.UnregisterFilter(ctx, id); err != nil {
				merr = errors.Join(merr, fmt.Errorf("failed to unregister filter: %s", id.String()))
			}
		}
	}
	var newIDs []*big.Int
	for id, ok := range visited {
		if !ok {
			uid, _ := new(big.Int).SetString(id, 10)
			newIDs = append(newIDs, uid)
		}
	}

	return newIDs, merr
}

func (p *logEventProvider) RegisterFilter(ctx context.Context, opts FilterOptions) error {
	upkeepID, cfg := opts.UpkeepID, opts.TriggerConfig
	if err := p.validateLogTriggerConfig(cfg); err != nil {
		return fmt.Errorf("invalid log trigger config: %w", err)
	}
	lpFilter := p.newLogFilter(upkeepID, cfg)

	// using lock to facilitate multiple events causing filter registration
	// at the same time.
	// Exploratory: consider using a q to handle registration requests
	p.registerLock.Lock()
	defer p.registerLock.Unlock()

	var filter upkeepFilter
	currentFilter := p.filterStore.Get(upkeepID)
	if currentFilter != nil {
		if currentFilter.configUpdateBlock > opts.UpdateBlock {
			// already registered with a config from a higher block number
			return fmt.Errorf("filter for upkeep with id %s already registered with newer config", upkeepID.String())
		} else if currentFilter.configUpdateBlock == opts.UpdateBlock {
			// already registered with the same config
			p.lggr.Debugf("filter for upkeep with id %s already registered with the same config", upkeepID.String())
			return nil
		}
		filter = *currentFilter
	} else { // new filter
		filter = upkeepFilter{
			upkeepID: upkeepID,
		}
	}
	filter.lastPollBlock = 0
	filter.lastRePollBlock = 0
	filter.configUpdateBlock = opts.UpdateBlock
	filter.selector = cfg.FilterSelector
	filter.addr = cfg.ContractAddress.Bytes()
	filter.topics = []common.Hash{cfg.Topic0, cfg.Topic1, cfg.Topic2, cfg.Topic3}

	if err := p.register(ctx, lpFilter, filter); err != nil {
		return fmt.Errorf("failed to register upkeep filter %s: %w", filter.upkeepID.String(), err)
	}

	return nil
}

// register registers the upkeep filter with the log poller and adds it to the filter store.
func (p *logEventProvider) register(ctx context.Context, lpFilter logpoller.Filter, ufilter upkeepFilter) error {
	latest, err := p.poller.LatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block while registering filter: %w", err)
	}
	lggr := logger.With(p.lggr, "upkeepID", ufilter.upkeepID.String())
	logPollerHasFilter := p.poller.HasFilter(lpFilter.Name)
	filterStoreHasFilter := p.filterStore.Has(ufilter.upkeepID)
	if filterStoreHasFilter {
		// removing filter in case of an update so we can recreate it with updated values
		lggr.Debugw("Upserting upkeep filter")
		err := p.poller.UnregisterFilter(ctx, lpFilter.Name)
		if err != nil {
			return fmt.Errorf("failed to upsert (unregister) upkeep filter %s: %w", ufilter.upkeepID.String(), err)
		}
	}
	if err := p.poller.RegisterFilter(ctx, lpFilter); err != nil {
		return err
	}
	p.filterStore.AddActiveUpkeeps(ufilter)
	if logPollerHasFilter {
		// already registered in DB before, no need to backfill
		return nil
	}
	backfillBlock := latest.BlockNumber - int64(LogBackfillBuffer)
	if backfillBlock < 1 {
		// New chain, backfill from start
		backfillBlock = 1
	}
	if int64(ufilter.configUpdateBlock) > backfillBlock {
		// backfill from config update block in case it is not too old
		backfillBlock = int64(ufilter.configUpdateBlock)
	}
	// NOTE: replys are planned to be done as part of RegisterFilter within logpoller
	lggr.Debugw("Backfilling logs for new upkeep filter", "backfillBlock", backfillBlock)
	p.poller.ReplayAsync(backfillBlock)

	return nil
}

func (p *logEventProvider) UnregisterFilter(ctx context.Context, upkeepID *big.Int) error {
	// Filter might have been unregistered already, only try to unregister if it exists
	if p.poller.HasFilter(p.filterName(upkeepID)) {
		if err := p.poller.UnregisterFilter(ctx, p.filterName(upkeepID)); err != nil {
			return fmt.Errorf("failed to unregister upkeep filter %s: %w", upkeepID.String(), err)
		}
	}
	p.filterStore.RemoveActiveUpkeeps(upkeepFilter{
		upkeepID: upkeepID,
	})
	return nil
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (p *logEventProvider) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	return logpoller.Filter{
		Name: p.filterName(upkeepID),
		// log poller filter treats this event sigs slice as an array of topic0
		// since we don't support multiple events right now, only put one topic0 here
		EventSigs: []common.Hash{common.BytesToHash(cfg.Topic0[:])},
		Addresses: []common.Address{cfg.ContractAddress},
		Retention: LogRetention,
	}
}

func (p *logEventProvider) validateLogTriggerConfig(cfg LogTriggerConfig) error {
	var zeroAddr common.Address
	var zeroBytes [32]byte
	if bytes.Equal(cfg.ContractAddress[:], zeroAddr[:]) {
		return errors.New("invalid contract address: zeroed")
	}
	if bytes.Equal(cfg.Topic0[:], zeroBytes[:]) {
		return errors.New("invalid topic0: zeroed")
	}
	s := cfg.FilterSelector
	if s >= 8 {
		p.lggr.Error("filter selector %d is invalid", s)
		return errors.New("invalid filter selector: larger or equal to 8")
	}
	return nil
}

func (p *logEventProvider) filterName(upkeepID *big.Int) string {
	return logpoller.FilterName("KeepersRegistry LogUpkeep", upkeepID.String())
}
