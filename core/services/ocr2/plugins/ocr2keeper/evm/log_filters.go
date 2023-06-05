package evm

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

const (
	// logRetention is the amount of time to retain logs for.
	// 5 minutes is the desired retention time for logs, but we add an extra 10 minutes buffer.
	// TODO: TBD make this configurable or based on block time.
	logRetention = (time.Minute * 5) + (time.Minute * 10)
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig = i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig

// logFilterManager manages log filters for upkeeps.
type logFilterManager struct {
	poller logpoller.LogPoller
}

func newLogFilterManager(poller logpoller.LogPoller) *logFilterManager {
	return &logFilterManager{
		poller: poller,
	}
}

// Register creates a filter from the given upkeep and calls log pollet to register it
func (lfm *logFilterManager) Register(upkeepID *big.Int, cfg LogTriggerConfig) error {
	if err := lfm.validateLogTriggerConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	filter := lfm.newLogFilter(upkeepID, cfg)
	// TODO: TBD remove old filter (if exist) for this upkeep
	// _ = lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
	return errors.Wrap(lfm.poller.RegisterFilter(filter), "failed to register upkeep filter")
}

// Unregister removes the filter for the given upkeepID
func (lfm *logFilterManager) Unregister(upkeepID *big.Int) error {
	err := lfm.poller.UnregisterFilter(lfm.filterName(upkeepID), nil)
	return errors.Wrap(err, "failed to unregister upkeep filter")
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (lfm *logFilterManager) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	sigs := lfm.getFiltersBySelector(cfg.FilterSelector, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	sigs = append([]common.Hash{common.BytesToHash(cfg.Topic0[:])}, sigs...)
	return logpoller.Filter{
		Name:      lfm.filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{cfg.ContractAddress},
		Retention: logRetention,
	}
}

func (lfm *logFilterManager) validateLogTriggerConfig(cfg LogTriggerConfig) error {
	var zeroAddr common.Address
	var zeroBytes [32]byte
	if bytes.Equal(cfg.ContractAddress[:], zeroAddr[:]) {
		return errors.New("invalid contract address: zeroed")
	}
	if bytes.Equal(cfg.Topic0[:], zeroBytes[:]) {
		return errors.New("invalid topic0: zeroed")
	}
	// TODO: TBD validate topic1, topic2, topic3
	return nil
}

// getFiltersBySelector the filters based on the filterSelector
func (lfm *logFilterManager) getFiltersBySelector(filterSelector uint8, filters ...[]byte) []common.Hash {
	var sigs []common.Hash
	var zeroBytes [32]byte
	for i, f := range filters {
		// bitwise AND the filterSelector with the index to check if the filter is needed
		mask := uint8(1 << uint8(i))
		a := filterSelector & mask
		if a == uint8(0) {
			continue
		}
		// TODO: TBD avoid adding zeroed filters
		if bytes.Equal(f, zeroBytes[:]) {
			continue
		}
		sigs = append(sigs, common.BytesToHash(common.LeftPadBytes(f, 32)))
	}
	return sigs
}

func (lfm *logFilterManager) filterName(upkeepID *big.Int) string {
	return logpoller.FilterName(upkeepID.String())
}
