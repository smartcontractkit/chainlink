package evm

import (
	"bytes"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
)

const (
	// logRetention is the amount of time to retain logs for
	logRetention = time.Minute * 5
)

// logFiltersManager manages log filters for upkeeps
type logFiltersManager struct {
	poller logpoller.LogPoller
}

func newLogFiltersManager(poller logpoller.LogPoller) *logFiltersManager {
	return &logFiltersManager{
		poller: poller,
	}
}

// Register takes an upkeep and register the corresponding filter if applicable
func (lfm *logFiltersManager) Register(upkeepID string, cfg keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig) error {
	err := validateLogTriggerConfig(cfg)
	if err != nil {
		return err
	}
	filter := newLogFilter(upkeepID, cfg)
	// TODO: TBD remove old filter (if exist) for this upkeep
	// _ = lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
	return lfm.poller.RegisterFilter(filter)
}

func (lfm *logFiltersManager) UnRegister(upkeepID string) error {
	return lfm.poller.UnregisterFilter(filterName(upkeepID), nil)
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func newLogFilter(upkeepID string, cfg keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig) logpoller.Filter {
	sigs := []common.Hash{
		common.BytesToHash(cfg.Topic0[:]),
	}
	sigs = addFiltersBySelector(cfg.FilterSelector, sigs, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	return logpoller.Filter{
		Name:      filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{cfg.ContractAddress},
		Retention: logRetention,
	}
}

func validateLogTriggerConfig(cfg keeper_registry_logic_b_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig) error {
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

func addFiltersBySelector(filterSelector uint8, sigs []common.Hash, filters ...[]byte) []common.Hash {
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

func filterName(upkeepID string) string {
	return logpoller.FilterName(upkeepID)
}
