package logprovider

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

var (
	// LogRetention is the amount of time to retain logs for.
	LogRetention = 24 * time.Hour
)

func (p *logEventProvider) RefreshActiveUpkeeps(ids ...*big.Int) ([]*big.Int, error) {
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
			if err := p.UnregisterFilter(id); err != nil {
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

func (p *logEventProvider) RegisterFilter(opts FilterOptions) error {
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
		// removing filter so we can recreate it with updated values
		err := p.poller.UnregisterFilter(p.filterName(currentFilter.upkeepID))
		if err != nil {
			return fmt.Errorf("failed to unregister upkeep filter %s for update: %w", upkeepID.String(), err)
		}
		filter = *currentFilter
	} else { // new filter
		filter = upkeepFilter{
			upkeepID:     upkeepID,
			blockLimiter: rate.NewLimiter(p.opts.BlockRateLimit, p.opts.BlockLimitBurst),
		}
	}
	filter.lastPollBlock = 0
	filter.lastRePollBlock = 0
	filter.configUpdateBlock = opts.UpdateBlock
	filter.addr = lpFilter.Addresses[0].Bytes()
	filter.topics = make([]common.Hash, len(lpFilter.EventSigs))
	copy(filter.topics, lpFilter.EventSigs)

	if err := p.register(lpFilter, filter); err != nil {
		return fmt.Errorf("failed to register upkeep filter %s: %w", filter.upkeepID.String(), err)
	}

	return nil
}

// register registers the upkeep filter with the log poller and adds it to the filter store.
func (p *logEventProvider) register(lpFilter logpoller.Filter, ufilter upkeepFilter) error {
	if err := p.poller.RegisterFilter(lpFilter); err != nil {
		return err
	}
	p.filterStore.AddActiveUpkeeps(ufilter)
	p.poller.ReplayAsync(int64(ufilter.configUpdateBlock))

	return nil
}

func (p *logEventProvider) UnregisterFilter(upkeepID *big.Int) error {
	err := p.poller.UnregisterFilter(p.filterName(upkeepID))
	if err != nil {
		// TODO: mark as removed in filter store, so we'll
		// automatically retry on next refresh
		return fmt.Errorf("failed to unregister upkeep filter %s: %w", upkeepID.String(), err)
	}
	p.filterStore.RemoveActiveUpkeeps(upkeepFilter{
		upkeepID: upkeepID,
	})
	return nil
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (p *logEventProvider) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	topics := p.getFiltersBySelector(cfg.FilterSelector, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	topics = append([]common.Hash{common.BytesToHash(cfg.Topic0[:])}, topics...)
	return logpoller.Filter{
		Name:      p.filterName(upkeepID),
		EventSigs: topics,
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
	return nil
}

// getFiltersBySelector the filters based on the filterSelector
func (p *logEventProvider) getFiltersBySelector(filterSelector uint8, filters ...[]byte) []common.Hash {
	var sigs []common.Hash
	var zeroBytes [32]byte
	for i, f := range filters {
		// bitwise AND the filterSelector with the index to check if the filter is needed
		mask := uint8(1 << uint8(i))
		a := filterSelector & mask
		if a == uint8(0) {
			continue
		}
		if bytes.Equal(f, zeroBytes[:]) {
			continue
		}
		sigs = append(sigs, common.BytesToHash(common.LeftPadBytes(f, 32)))
	}
	return sigs
}

func (p *logEventProvider) filterName(upkeepID *big.Int) string {
	return logpoller.FilterName("KeepersRegistry LogUpkeep", upkeepID.String())
}
