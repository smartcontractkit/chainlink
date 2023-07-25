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

	// TODO: optimize locking, currently we lock the whole map while registering the filter
	p.lock.Lock()
	defer p.lock.Unlock()

	uid := upkeepID.String()
	if _, ok := p.active[uid]; ok {
		// TODO: check for updates
		return errors.Errorf("filter for upkeep with id %s already registered", uid)
	}
	if err := p.poller.RegisterFilter(filter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}
	p.active[uid] = upkeepFilterEntry{
		id:           upkeepID,
		filter:       filter,
		cfg:          cfg,
		blockLimiter: rate.NewLimiter(p.opts.BlockRateLimit, p.opts.BlockLimitBurst),
	}

	return nil
}

func (p *logEventProvider) UnregisterFilter(upkeepID *big.Int) error {
	err := p.poller.UnregisterFilter(p.filterName(upkeepID), nil)
	if err == nil {
		p.lock.Lock()
		delete(p.active, upkeepID.String())
		p.lock.Unlock()
	}
	return errors.Wrap(err, "failed to unregister upkeep filter")
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (p *logEventProvider) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	sigs := p.getFiltersBySelector(cfg.FilterSelector, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	sigs = append([]common.Hash{common.BytesToHash(cfg.Topic0[:])}, sigs...)
	return logpoller.Filter{
		Name:      p.filterName(upkeepID),
		EventSigs: sigs,
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
