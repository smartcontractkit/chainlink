package changeset

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	verifier "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

type DataStreamsEVMChainState struct {
	commonchangeset.MCMSWithTimelockState
	Configurators       map[view.Address]ContractAndMetadata[*configurator.Configurator, *metadata.GenericContractMetadata[v0_5.ConfiguratorView]]
	ChannelConfigStores map[view.Address]ContractAndMetadata[*channel_config_store.ChannelConfigStore, *metadata.GenericContractMetadata[v0_5.ChannelConfigStoreView]]
	FeeManagers         map[view.Address]ContractAndMetadata[*fee_manager.FeeManager, *metadata.GenericContractMetadata[v0_5.FeeManagerView]]
	RewardManagers      map[view.Address]ContractAndMetadata[*rewardManager.RewardManager, *metadata.GenericContractMetadata[v0_5.RewardManagerView]]
	Verifiers           map[view.Address]ContractAndMetadata[*verifier.Verifier, *metadata.GenericContractMetadata[v0_5.VerifierView]]
	VerifierProxys      map[view.Address]ContractAndMetadata[*verifier_proxy_v0_5_0.VerifierProxy, *metadata.GenericContractMetadata[v0_5.VerifierProxyView]]
}

type ContractAndMetadata[C any, M any] struct {
	Contract C
	Metadata M // retrieved from datastore contract metadata
}

type DataStreamsOnChainState struct {
	Chains map[uint64]DataStreamsEVMChainState
}

func LoadOnchainState(e cldf.Environment) (DataStreamsOnChainState, error) {
	state := DataStreamsOnChainState{
		Chains: make(map[uint64]DataStreamsEVMChainState),
	}
	envDatastore, err := datastore.FromDefault[metadata.SerializedContractMetadata, datastore.DefaultMetadata](e.DataStore)
	if err != nil {
		return state, fmt.Errorf("failed to convert datastore: %w", err)
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := LoadChainState(e.Logger, chain, addresses, envDatastore.ContractMetadata())
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = *chainState
	}
	return state, nil
}

func LoadChainState(logger logger.Logger,
	chain cldf.Chain,
	addresses map[string]cldf.TypeAndVersion,
	mdStore datastore.ContractMetadataStore[metadata.SerializedContractMetadata]) (*DataStreamsEVMChainState, error) {
	var cc DataStreamsEVMChainState

	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	cc.MCMSWithTimelockState = *mcmsWithTimelock

	cc.Configurators = make(map[view.Address]ContractAndMetadata[*configurator.Configurator, *metadata.GenericContractMetadata[v0_5.ConfiguratorView]])
	cc.ChannelConfigStores = make(map[view.Address]ContractAndMetadata[*channel_config_store.ChannelConfigStore, *metadata.GenericContractMetadata[v0_5.ChannelConfigStoreView]])
	cc.FeeManagers = make(map[view.Address]ContractAndMetadata[*fee_manager.FeeManager, *metadata.GenericContractMetadata[v0_5.FeeManagerView]])
	cc.RewardManagers = make(map[view.Address]ContractAndMetadata[*rewardManager.RewardManager, *metadata.GenericContractMetadata[v0_5.RewardManagerView]])
	cc.Verifiers = make(map[view.Address]ContractAndMetadata[*verifier.Verifier, *metadata.GenericContractMetadata[v0_5.VerifierView]])
	cc.VerifierProxys = make(map[view.Address]ContractAndMetadata[*verifier_proxy_v0_5_0.VerifierProxy, *metadata.GenericContractMetadata[v0_5.VerifierProxyView]])

	for address, tv := range addresses {
		if belongsToMCMS(address, mcmsWithTimelock) {
			continue
		}

		switch tv.String() {
		case cldf.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0).String():
			ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.ChannelConfigStoreView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.ChannelConfigStores[address] = ContractAndMetadata[
				*channel_config_store.ChannelConfigStore,
				*metadata.GenericContractMetadata[v0_5.ChannelConfigStoreView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		case cldf.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0).String():
			ccs, err := fee_manager.NewFeeManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.FeeManagerView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.FeeManagers[address] = ContractAndMetadata[
				*fee_manager.FeeManager,
				*metadata.GenericContractMetadata[v0_5.FeeManagerView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		case cldf.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0).String():
			ccs, err := configurator.NewConfigurator(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.ConfiguratorView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.Configurators[address] = ContractAndMetadata[
				*configurator.Configurator,
				*metadata.GenericContractMetadata[v0_5.ConfiguratorView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		case cldf.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0).String():
			ccs, err := rewardManager.NewRewardManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.RewardManagerView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.RewardManagers[address] = ContractAndMetadata[
				*rewardManager.RewardManager,
				*metadata.GenericContractMetadata[v0_5.RewardManagerView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		case cldf.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0).String():
			ccs, err := verifier.NewVerifier(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.VerifierView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.Verifiers[address] = ContractAndMetadata[
				*verifier.Verifier,
				*metadata.GenericContractMetadata[v0_5.VerifierView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		case cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0).String():
			ccs, err := verifier_proxy_v0_5_0.NewVerifierProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return &cc, err
			}
			previousMetadata, err := getPreviousMetadata[v0_5.VerifierProxyView](mdStore, chain.Selector, address)
			if err != nil {
				return &cc, fmt.Errorf("failed to get previous metadata: %w", err)
			}
			cc.VerifierProxys[address] = ContractAndMetadata[
				*verifier_proxy_v0_5_0.VerifierProxy,
				*metadata.GenericContractMetadata[v0_5.VerifierProxyView]]{
				Contract: ccs,
				Metadata: previousMetadata,
			}

		default:
			return &cc, fmt.Errorf("unknown contract %s", tv)
		}
	}
	return &cc, nil
}

func (s DataStreamsOnChainState) View(ctx context.Context, chains []uint64) (map[uint64]view.EVMChainView, error) {
	m := make(map[uint64]view.EVMChainView)
	for _, chainSelector := range chains {
		if _, ok := s.Chains[chainSelector]; !ok {
			return m, fmt.Errorf("chain not supported %d", chainSelector)
		}
		chainState := s.Chains[chainSelector]
		chainView, err := chainState.GenerateView(ctx)
		if err != nil {
			return m, err
		}
		m[chainSelector] = chainView
	}
	return m, nil
}

// GenerateView generates a view for the DataStreamsEVMChainState
func (s DataStreamsEVMChainState) GenerateView(ctx context.Context) (view.EVMChainView, error) {
	chainView := view.NewChain()
	// It would be a nice improvement to generate in goroutines & abstract this & use the ViewGenerator interface so that we don't have to repeat
	configuratorViews, err := s.GenerateConfiguratorViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate configurator views: %w", err)
	}
	for address, contractView := range configuratorViews {
		chainView.Configurator[address] = contractView
	}

	verifierViews, err := s.GenerateVerifierViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate verifier views: %w", err)
	}
	for address, contractView := range verifierViews {
		chainView.Verifier[address] = contractView
	}

	feeManagerViews, err := s.GenerateFeeManagerViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate fee manager views: %w", err)
	}
	for address, contractView := range feeManagerViews {
		chainView.FeeManager[address] = contractView
	}

	rewardManagerViews, err := s.GenerateRewardManagerViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate reward manager views: %w", err)
	}
	for address, contractView := range rewardManagerViews {
		chainView.RewardManager[address] = contractView
	}

	verifierProxyViews, err := s.GenerateVerifierProxyViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate verifier proxy views: %w", err)
	}
	for address, contractView := range verifierProxyViews {
		chainView.VerifierProxy[address] = contractView
	}

	channelConfigStoreViews, err := s.GenerateChannelConfigStoreViews(ctx)
	if err != nil {
		return chainView, fmt.Errorf("failed to generate channel config store views: %w", err)
	}
	for address, contractView := range channelConfigStoreViews {
		chainView.ChannelConfigStore[address] = contractView
	}

	return chainView, nil
}

func (s DataStreamsEVMChainState) GenerateConfiguratorViews(ctx context.Context) (map[view.Address]v0_5.ConfiguratorView, error) {
	result := make(map[view.Address]v0_5.ConfiguratorView)
	for address, contractAndMeta := range s.Configurators {
		contractContext := v0_5.ConfiguratorViewParams{
			FromBlock: contractAndMeta.Metadata.Metadata.DeployBlock,
		}

		contractWrapper := evm.NewConfiguratorReader(contractAndMeta.Contract)
		generator := v0_5.NewConfiguratorViewGenerator(contractWrapper)
		configuratorView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = configuratorView
	}

	return result, nil
}

func (s DataStreamsEVMChainState) GenerateVerifierViews(ctx context.Context) (map[view.Address]v0_5.VerifierView, error) {
	result := make(map[view.Address]v0_5.VerifierView)
	for address, contractAndMeta := range s.Verifiers {
		contractContext := v0_5.VerifierViewParams{
			FromBlock: contractAndMeta.Metadata.Metadata.DeployBlock,
		}

		contractWrapper := evm.NewVerifierReader(contractAndMeta.Contract)
		generator := v0_5.NewVerifierViewGenerator(contractWrapper)
		contractView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = contractView
	}

	return result, nil
}

func (s DataStreamsEVMChainState) GenerateFeeManagerViews(ctx context.Context) (map[view.Address]v0_5.FeeManagerView, error) {
	result := make(map[view.Address]v0_5.FeeManagerView)
	for address, contractAndMeta := range s.FeeManagers {
		contractContext := v0_5.FeeManagerViewParams{
			FromBlock: contractAndMeta.Metadata.Metadata.DeployBlock,
		}

		contractWrapper := evm.NewFeeManagerReader(contractAndMeta.Contract)
		generator := v0_5.NewFeeManagerViewGenerator(contractWrapper)
		contractView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = contractView
	}

	return result, nil
}

func (s DataStreamsEVMChainState) GenerateRewardManagerViews(ctx context.Context) (map[view.Address]v0_5.RewardManagerView, error) {
	result := make(map[view.Address]v0_5.RewardManagerView)
	for address, contractAndMeta := range s.RewardManagers {
		contractContext := v0_5.RewardManagerViewParams{}

		contractWrapper := evm.NewRewardManagerReader(contractAndMeta.Contract)
		generator := v0_5.NewRewardManagerViewGenerator(contractWrapper)
		contractView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = contractView
	}

	return result, nil
}

func (s DataStreamsEVMChainState) GenerateVerifierProxyViews(ctx context.Context) (map[view.Address]v0_5.VerifierProxyView, error) {
	result := make(map[view.Address]v0_5.VerifierProxyView)
	for address, contractAndMeta := range s.VerifierProxys {
		contractContext := v0_5.VerifierProxyViewParams{}

		contractWrapper := evm.NewVerifierProxyReader(contractAndMeta.Contract)
		generator := v0_5.NewVerifierProxyViewGenerator(contractWrapper)
		contractView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = contractView
	}

	return result, nil
}

func (s DataStreamsEVMChainState) GenerateChannelConfigStoreViews(ctx context.Context) (map[view.Address]v0_5.ChannelConfigStoreView, error) {
	result := make(map[view.Address]v0_5.ChannelConfigStoreView)
	for address, contractAndMeta := range s.ChannelConfigStores {
		contractContext := v0_5.ChannelConfigStoreViewParams{}

		contractWrapper := evm.NewChannelConfigStoreWrapper(contractAndMeta.Contract)
		generator := v0_5.NewChannelConfigStoreViewGenerator(contractWrapper)
		contractView, err := generator.Generate(ctx, contractContext)
		if err != nil {
			return nil, fmt.Errorf("failed to build view for configurator %s: %w", address, err)
		}
		result[address] = contractView
	}

	return result, nil
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

func getPreviousMetadata[M any](
	mdStore datastore.ContractMetadataStore[metadata.SerializedContractMetadata],
	chainSelector uint64,
	address string,
) (*metadata.GenericContractMetadata[M], error) {
	cm, err := mdStore.Get(
		datastore.NewContractMetadataKey(chainSelector, address),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get contract metadata: %w", err)
	}
	contractMetadata, err := metadata.DeserializeMetadata[M](cm.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to convert contract metadata: %w", err)
	}
	return contractMetadata, nil
}
