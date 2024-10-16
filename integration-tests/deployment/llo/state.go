package llo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
)

// LLOChainState holds a Go binding for all the currently deployed LLO contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type LLOChainState struct {
	ChannelConfigStore *channel_config_store.ChannelConfigStore
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (LLOChainState, error) {
	var state LLOChainState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(ChannelConfigStore, deployment.Version1_0_0).String():
			ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.ChannelConfigStore = ccs
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
