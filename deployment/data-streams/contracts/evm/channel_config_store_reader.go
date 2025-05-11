package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
)

// ConcreteChannelDefinitionIteratorAdapter adapts the concrete iterator to our interface
type ConcreteChannelDefinitionIteratorAdapter struct {
	*channel_config_store.ChannelConfigStoreNewChannelDefinitionIterator
}

func (a *ConcreteChannelDefinitionIteratorAdapter) GetEvent() *channel_config_store.ChannelConfigStoreNewChannelDefinition {
	return a.Event
}

// ChannelConfigStoreReader wraps the actual contract to return our interface
type ChannelConfigStoreReader struct {
	*channel_config_store.ChannelConfigStoreFilterer
	*channel_config_store.ChannelConfigStoreCaller
}

// FilterNewChannelDefinition wraps the original method to return our interface
func (w *ChannelConfigStoreReader) FilterNewChannelDefinition(
	opts *bind.FilterOpts,
	donID []*big.Int) (LogIterator[channel_config_store.ChannelConfigStoreNewChannelDefinition], error) {
	iter, err := w.ChannelConfigStoreFilterer.FilterNewChannelDefinition(opts, donID)
	if err != nil {
		return nil, err
	}

	return &ConcreteChannelDefinitionIteratorAdapter{iter}, nil
}

func NewChannelConfigStoreWrapper(contract *channel_config_store.ChannelConfigStore) *ChannelConfigStoreReader {
	return &ChannelConfigStoreReader{
		ChannelConfigStoreFilterer: &contract.ChannelConfigStoreFilterer,
		ChannelConfigStoreCaller:   &contract.ChannelConfigStoreCaller,
	}
}
