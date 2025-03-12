package data_streams

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
)

// ChainConfig holds a Go binding for all the currently deployed LLO contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type ChainConfig struct {
	// ChannelConfigStores is a map of chain selector to a list of all ChannelConfigStoreContract contracts on that chain.
	ChannelConfigStores map[uint64][]*channel_config_store.ChannelConfigStore
}

// LoadChainConfig Loads all state for a chain into state.
//
// Param addresses is a map of all known contract addresses on this chain.
func LoadChainConfig(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (ChainConfig, error) {
	cc := ChainConfig{
		ChannelConfigStores: make(map[uint64][]*channel_config_store.ChannelConfigStore),
	}
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0).String():
			ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.ChannelConfigStores[chain.Selector] = append(cc.ChannelConfigStores[chain.Selector], ccs)
		default:
			return cc, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return cc, nil
}

func (cc ChainConfig) Validate() error {
	// We want to ensure that the ChannelConfigStores map is not nil.
	if cc.ChannelConfigStores == nil {
		return errors.New("ChannelConfigStores is nil")
	}
	return nil
}
