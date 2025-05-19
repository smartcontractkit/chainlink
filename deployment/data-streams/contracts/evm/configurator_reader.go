package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

// ProductionConfigSetIteratorAdapter adapts the concrete iterator to our interface
type ProductionConfigSetIteratorAdapter struct {
	*configurator.ConfiguratorProductionConfigSetIterator
}

// GetEvent returns the current event
func (a *ProductionConfigSetIteratorAdapter) GetEvent() *configurator.ConfiguratorProductionConfigSet {
	return a.Event
}

// StagingConfigSetIteratorAdapter adapts the concrete iterator to our interface
type StagingConfigSetIteratorAdapter struct {
	*configurator.ConfiguratorStagingConfigSetIterator
}

// GetEvent returns the current event
func (a *StagingConfigSetIteratorAdapter) GetEvent() *configurator.ConfiguratorStagingConfigSet {
	return a.Event
}

// PromoteStagingConfigIteratorAdapter adapts the concrete iterator to our interface
type PromoteStagingConfigIteratorAdapter struct {
	*configurator.ConfiguratorPromoteStagingConfigIterator
}

func (a *PromoteStagingConfigIteratorAdapter) GetEvent() *configurator.ConfiguratorPromoteStagingConfig {
	return a.Event
}

type ConfiguratorReader struct {
	*configurator.ConfiguratorFilterer
	*configurator.ConfiguratorCaller
}

func (r *ConfiguratorReader) FilterProductionConfigSet(
	opts *bind.FilterOpts,
	configID [][32]byte,
) (LogIterator[configurator.ConfiguratorProductionConfigSet], error) {
	iter, err := r.ConfiguratorFilterer.FilterProductionConfigSet(opts, configID)
	if err != nil {
		return nil, err
	}

	return &ProductionConfigSetIteratorAdapter{iter}, nil
}

func (r *ConfiguratorReader) FilterStagingConfigSet(
	opts *bind.FilterOpts,
	configID [][32]byte,
) (LogIterator[configurator.ConfiguratorStagingConfigSet], error) {
	iter, err := r.ConfiguratorFilterer.FilterStagingConfigSet(opts, configID)
	if err != nil {
		return nil, err
	}

	return &StagingConfigSetIteratorAdapter{iter}, nil
}

func (r *ConfiguratorReader) FilterPromoteStagingConfig(
	opts *bind.FilterOpts,
	configID [][32]byte,
	retiredConfigDigest [][32]byte,
) (LogIterator[configurator.ConfiguratorPromoteStagingConfig], error) {
	iter, err := r.ConfiguratorFilterer.FilterPromoteStagingConfig(opts, configID, retiredConfigDigest)
	if err != nil {
		return nil, err
	}

	return &PromoteStagingConfigIteratorAdapter{iter}, nil
}

// NewConfiguratorReader creates a new wrapper for the real contract
func NewConfiguratorReader(contract *configurator.Configurator) *ConfiguratorReader {
	return &ConfiguratorReader{
		ConfiguratorFilterer: &contract.ConfiguratorFilterer,
		ConfiguratorCaller:   &contract.ConfiguratorCaller,
	}
}
