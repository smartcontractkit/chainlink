package changeset

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	verifier "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

type DataStreamsChainState struct {
	commonchangeset.MCMSWithTimelockState
	Configurators       map[common.Address]*configurator.Configurator
	ChannelConfigStores map[common.Address]*channel_config_store.ChannelConfigStore
	FeeManagers         map[common.Address]*fee_manager.FeeManager
	LinkTokens          map[common.Address]*link_token_interface.LinkToken
	RewardManagers      map[common.Address]*rewardManager.RewardManager
	Verifiers           map[common.Address]*verifier.Verifier
	VerifierProxys      map[common.Address]*verifier_proxy_v0_5_0.VerifierProxy
}

type DataStreamsOnChainState struct {
	Chains map[uint64]DataStreamsChainState
}

func LoadOnchainState(e deployment.Environment) (DataStreamsOnChainState, error) {
	state := DataStreamsOnChainState{
		Chains: make(map[uint64]DataStreamsChainState),
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := LoadChainState(e.Logger, chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = *chainState
	}
	return state, nil
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(logger logger.Logger, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*DataStreamsChainState, error) {
	var cc DataStreamsChainState

	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	cc.MCMSWithTimelockState = *mcmsWithTimelock

	cc.Configurators = make(map[common.Address]*configurator.Configurator)
	cc.ChannelConfigStores = make(map[common.Address]*channel_config_store.ChannelConfigStore)
	cc.FeeManagers = make(map[common.Address]*fee_manager.FeeManager)
	cc.LinkTokens = make(map[common.Address]*link_token_interface.LinkToken)
	cc.RewardManagers = make(map[common.Address]*rewardManager.RewardManager)
	cc.Verifiers = make(map[common.Address]*verifier.Verifier)
	cc.VerifierProxys = make(map[common.Address]*verifier_proxy_v0_5_0.VerifierProxy)

	for address, tv := range addresses {
		if belongsToMCMS(address, mcmsWithTimelock) {
			continue
		}

		switch tv.String() {
		case deployment.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0).String():
			ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.ChannelConfigStores[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0).String():
			ccs, err := fee_manager.NewFeeManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.FeeManagers[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(commonTypes.LinkToken, deployment.Version1_0_0).String():
			ccs, err := link_token_interface.NewLinkToken(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.LinkTokens[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0).String():
			ccs, err := configurator.NewConfigurator(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.Configurators[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0).String():
			ccs, err := rewardManager.NewRewardManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.RewardManagers[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0).String():
			ccs, err := verifier.NewVerifier(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.Verifiers[common.HexToAddress(address)] = ccs

		case deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0).String():
			css, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			cc.VerifierProxys[common.HexToAddress(address)] = css

		default:
			return &cc, fmt.Errorf("unknown contract %s", tv)
		}
	}
	return &cc, nil
}

func (s DataStreamsOnChainState) View(chains []uint64) (map[string]view.ChainView, error) {
	m := make(map[string]view.ChainView)
	for _, chainSelector := range chains {
		chainInfo, err := deployment.ChainInfo(chainSelector)
		if err != nil {
			return m, err
		}
		if _, ok := s.Chains[chainSelector]; !ok {
			return m, fmt.Errorf("chain not supported %d", chainSelector)
		}
		chainState := s.Chains[chainSelector]
		chainView, err := chainState.GenerateView()
		if err != nil {
			return m, err
		}
		name := chainInfo.ChainName
		if chainInfo.ChainName == "" {
			name = strconv.FormatUint(chainSelector, 10)
		}
		m[name] = chainView
	}
	return m, nil
}

func (c DataStreamsChainState) GenerateView() (view.ChainView, error) {
	chainView := view.NewChain()
	if c.Configurators != nil {
		for _, configurator := range c.Configurators {
			configuratorView, err := v0_5.GenerateConfiguratorView(configurator)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate configurator view %s", configurator.Address().String())
			}
			chainView.Configurator[configurator.Address().Hex()] = configuratorView
		}
	}
	if c.RewardManagers != nil {
		for _, rm := range c.RewardManagers {
			rmView, err := v0_5.GenerateRewardManagerView(rm)
			if err != nil {
				return chainView, errors.Wrapf(err, "failed to generate RewardManager view %s", rm.Address().String())
			}
			chainView.RewardManager[rm.Address().Hex()] = rmView
		}
	}
	return chainView, nil
}

// Helper function to determine if an address belongs to the MCMS contracts, and should be loaded in a separated way
func belongsToMCMS(addr string, mcmsWithTimelock *commonchangeset.MCMSWithTimelockState) bool {
	if mcmsWithTimelock == nil || mcmsWithTimelock.MCMSWithTimelockContracts == nil {
		return false
	}
	c := mcmsWithTimelock.MCMSWithTimelockContracts

	switch {
	case c.Timelock != nil && c.Timelock.Address().String() == addr:
		return true
	case c.CallProxy != nil && c.CallProxy.Address().String() == addr:
		return true
	case c.ProposerMcm != nil && c.ProposerMcm.Address().String() == addr:
		return true
	case c.CancellerMcm != nil && c.CancellerMcm.Address().String() == addr:
		return true
	case c.BypasserMcm != nil && c.BypasserMcm.Address().String() == addr:
		return true
	}
	return false
}
