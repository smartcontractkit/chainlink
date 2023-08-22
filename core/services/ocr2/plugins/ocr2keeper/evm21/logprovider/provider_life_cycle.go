package logprovider

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

var (
	// LogRetention is the amount of time to retain logs for.
	LogRetention = 24 * time.Hour
)

func (p *logEventProvider) RegisterFilter(opts FilterOptions) error {
	upkeepID, cfg := opts.UpkeepID, opts.TriggerConfig
	if err := p.validateLogTriggerConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	lpFilter := p.newLogFilter(upkeepID, cfg)

	// using lock to facilitate multiple events causing filter registration
	// at the same time.
	// TODO: consider using a q to handle registration requests
	p.registerLock.Lock()
	defer p.registerLock.Unlock()

	var filter upkeepFilter
	currentFilter := p.filterStore.Get(upkeepID)
	if currentFilter != nil {
		if currentFilter.configUpdateBlock > opts.UpdateBlock {
			// already registered with a config from a higher block number
			return errors.Errorf("filter for upkeep with id %s already registered with newer config", upkeepID.String())
		} else if currentFilter.configUpdateBlock == opts.UpdateBlock {
			// already registered with the same config
			p.lggr.Debugf("filter for upkeep with id %s already registered with the same config", upkeepID.String())
			return nil
		}
		// removing filter so we can recreate it with updated values
		err := p.UnregisterFilter(currentFilter.upkeepID)
		if err != nil {
			return errors.Wrap(err, "failed to unregister upkeep filter for update")
		}
		filter = *currentFilter
	} else { // new filter
		filter = upkeepFilter{
			upkeepID:        upkeepID,
			blockLimiter:    rate.NewLimiter(p.opts.BlockRateLimit, p.opts.BlockLimitBurst),
			lastPollBlock:   0,
			lastRePollBlock: 0,
		}
	}
	filter.configUpdateBlock = opts.UpdateBlock
	filter.addr = lpFilter.Addresses[0].Bytes()
	filter.topics = make([]common.Hash, len(lpFilter.EventSigs))
	copy(filter.topics, lpFilter.EventSigs)

	if err := p.register(lpFilter, filter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}

	return nil
}

func (p *logEventProvider) register(lpFilter logpoller.Filter, ufilter upkeepFilter) error {
	if err := p.poller.RegisterFilter(lpFilter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}
	p.filterStore.AddActiveUpkeeps(ufilter)
	p.poller.ReplayAsync(int64(ufilter.configUpdateBlock))

	return nil
}

func (p *logEventProvider) UnregisterFilter(upkeepID *big.Int) error {
	err := p.poller.UnregisterFilter(p.filterName(upkeepID))
	if err == nil {
		p.filterStore.RemoveActiveUpkeeps(upkeepFilter{
			upkeepID: upkeepID,
		})
	}
	return errors.Wrap(err, "failed to unregister upkeep filter")
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
