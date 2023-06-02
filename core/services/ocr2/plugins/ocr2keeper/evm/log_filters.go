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
type LogTriggerConfig i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig

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
	err := validateLogTriggerConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	filter := newLogFilter(upkeepID, cfg)
	// TODO: TBD remove old filter (if exist) for this upkeep
	// _ = lfp.poller.UnregisterFilter(filterName(upkeepID), nil)
	err = lfm.poller.RegisterFilter(filter)
	return errors.Wrap(err, "failed to register upkeep filter")
}

// Unregister removes the filter for the given upkeepID
func (lfm *logFilterManager) Unregister(upkeepID *big.Int) error {
	err := lfm.poller.UnregisterFilter(filterName(upkeepID), nil)
	return errors.Wrap(err, "failed to unregister upkeep filter")
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
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

func validateLogTriggerConfig(cfg LogTriggerConfig) error {
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

// addFiltersBySelector adds the filters to the sigs slice based on the filterSelector
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

func filterName(upkeepID *big.Int) string {
	return logpoller.FilterName(upkeepID.String())
}
