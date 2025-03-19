package data_streams

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	verifier "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

// ChainConfig holds a Go binding for all the currently deployed LLO contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type ChainConfig struct {
	// ChannelConfigStores is a map of chain selector to a list of all ChannelConfigStoreContract contracts on that chain.
	ChannelConfigStores map[uint64][]*channel_config_store.ChannelConfigStore
	// FeeManagers is a map of chain selector to a list of all FeeManagerContract contracts on that chain.
	FeeManagers map[uint64][]*fee_manager.FeeManager
	// Configurators is a map of chain selector to a list of all ConfiguratorsContract contracts on that chain.
	Configurators map[uint64][]*configurator.Configurator

	// Verifiers is a map of chain selector to a list of all VerifiersContract contracts on that chain.
	Verifiers map[uint64][]*verifier.Verifier

	VerifierProxys map[uint64][]*verifier_proxy_v0_5_0.VerifierProxy
}

// LoadChainConfig Loads all state for a chain into state.
//
// Param addresses is a map of all known contract addresses on this chain.
func LoadChainConfig(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (ChainConfig, error) {
	cc := ChainConfig{
		ChannelConfigStores: make(map[uint64][]*channel_config_store.ChannelConfigStore),
		FeeManagers:         make(map[uint64][]*fee_manager.FeeManager),
		Configurators:       make(map[uint64][]*configurator.Configurator),
		Verifiers:           make(map[uint64][]*verifier.Verifier),
		VerifierProxys:      make(map[uint64][]*verifier_proxy_v0_5_0.VerifierProxy),
	}
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0).String():
			proxy, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.VerifierProxys[chain.Selector] = append(cc.VerifierProxys[chain.Selector], proxy)
		case deployment.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0).String():
			ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.ChannelConfigStores[chain.Selector] = append(cc.ChannelConfigStores[chain.Selector], ccs)

		case deployment.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0).String():
			ccs, err := fee_manager.NewFeeManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.FeeManagers[chain.Selector] = append(cc.FeeManagers[chain.Selector], ccs)

		case deployment.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0).String():
			ccs, err := configurator.NewConfigurator(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.Configurators[chain.Selector] = append(cc.Configurators[chain.Selector], ccs)

		case deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0).String():
			ccs, err := verifier.NewVerifier(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.Verifiers[chain.Selector] = append(cc.Verifiers[chain.Selector], ccs)

		case deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0).String():
			proxy, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return cc, err
			}
			cc.VerifierProxys[chain.Selector] = append(cc.VerifierProxys[chain.Selector], proxy)
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
