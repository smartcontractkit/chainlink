package logprovider

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

func (p *logEventProvider) RegisterFilter(upkeepID *big.Int, cfg LogTriggerConfig) error {
	if err := p.validateLogTriggerConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	filter := p.newLogFilter(upkeepID, cfg)

	if p.filterStore.Has(upkeepID) {
		// TODO: check for updates
		return errors.Errorf("filter for upkeep with id %s already registered", upkeepID.String())
	}
	if err := p.poller.RegisterFilter(filter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}
	topics := make([]common.Hash, len(filter.EventSigs))
	copy(topics, filter.EventSigs)
	p.filterStore.AddActiveUpkeeps(upkeepFilter{
		upkeepID:     upkeepID,
		addr:         filter.Addresses[0].Bytes(),
		topics:       topics,
		blockLimiter: rate.NewLimiter(p.opts.BlockRateLimit, p.opts.BlockLimitBurst),
	})

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
		Retention: p.opts.LogRetention,
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
