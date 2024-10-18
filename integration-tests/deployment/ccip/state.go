package ccipdeployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_6"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	owner_wrappers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

// CCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	OnRamp             *onramp.OnRamp
	OffRamp            *offramp.OffRamp
	FeeQuoter          *fee_quoter.FeeQuoter
	RMNProxy           *rmn_proxy_contract.RMNProxyContract
	NonceManager       *nonce_manager.NonceManager
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	Router             *router.Router
	CommitStore        *commit_store.CommitStore
	Weth9              *weth9.WETH9
	RMNRemote          *rmn_remote.RMNRemote
	// TODO: May need to support older link too
	LinkToken *burn_mint_erc677.BurnMintERC677
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token contract
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	// Not all tokens will be burn and mint tokens.
	BurnMintTokens677 map[TokenSymbol]*burn_mint_erc677.BurnMintERC677
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Note we only expect one of these (on the home chain)
	CapabilityRegistry *capabilities_registry.CapabilitiesRegistry
	CCIPHome           *ccip_home.CCIPHome
	RMNHome            *rmn_home.RMNHome
	AdminMcm           *owner_wrappers.ManyChainMultiSig
	BypasserMcm        *owner_wrappers.ManyChainMultiSig
	CancellerMcm       *owner_wrappers.ManyChainMultiSig
	ProposerMcm        *owner_wrappers.ManyChainMultiSig
	Timelock           *owner_wrappers.RBACTimelock
	// TODO remove once staging upgraded.
	CCIPConfig *ccip_config.CCIPConfig

	// Test contracts
	Receiver   *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	TestRouter *router.Router
}

func (c CCIPChainState) GenerateView() (view.ChainView, error) {
	chainView := view.NewChain()
	if c.Router != nil {
		routerView, err := v1_2.GenerateRouterView(c.Router)
		if err != nil {
			return chainView, err
		}
		chainView.Router[c.Router.Address().Hex()] = routerView
	}
	if c.TokenAdminRegistry != nil {
		taView, err := v1_5.GenerateTokenAdminRegistryView(c.TokenAdminRegistry)
		if err != nil {
			return chainView, err
		}
		chainView.TokenAdminRegistry[c.TokenAdminRegistry.Address().Hex()] = taView
	}
	if c.NonceManager != nil {
		nmView, err := v1_6.GenerateNonceManagerView(c.NonceManager)
		if err != nil {
			return chainView, err
		}
		chainView.NonceManager[c.NonceManager.Address().Hex()] = nmView
	}
	if c.RMNRemote != nil {
		rmnView, err := v1_6.GenerateRMNRemoteView(c.RMNRemote)
		if err != nil {
			return chainView, err
		}
		chainView.RMN[c.RMNRemote.Address().Hex()] = rmnView
	}
	if c.FeeQuoter != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		fqView, err := v1_6.GenerateFeeQuoterView(c.FeeQuoter, c.Router, c.TokenAdminRegistry)
		if err != nil {
			return chainView, err
		}
		chainView.FeeQuoter[c.FeeQuoter.Address().Hex()] = fqView
	}

	if c.OnRamp != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		onRampView, err := v1_6.GenerateOnRampView(
			c.OnRamp,
			c.Router,
			c.TokenAdminRegistry,
		)
		if err != nil {
			return chainView, err
		}
		chainView.OnRamp[c.OnRamp.Address().Hex()] = onRampView
	}

	if c.OffRamp != nil && c.Router != nil {
		offRampView, err := v1_6.GenerateOffRampView(
			c.OffRamp,
			c.Router,
		)
		if err != nil {
			return chainView, err
		}
		chainView.OffRamp[c.OffRamp.Address().Hex()] = offRampView
	}

	if c.CommitStore != nil {
		commitStoreView, err := v1_5.GenerateCommitStoreView(c.CommitStore)
		if err != nil {
			return chainView, err
		}
		chainView.CommitStore[c.CommitStore.Address().Hex()] = commitStoreView
	}

	if c.RMNProxy != nil {
		rmnProxyView, err := v1_0.GenerateRMNProxyView(c.RMNProxy)
		if err != nil {
			return chainView, err
		}
		chainView.RMNProxy[c.RMNProxy.Address().Hex()] = rmnProxyView
	}
	if c.CapabilityRegistry != nil {
		capRegView, err := v1_6.GenerateCapRegView(c.CapabilityRegistry)
		if err != nil {
			return chainView, err
		}
		chainView.CapabilityRegistry[c.CapabilityRegistry.Address().Hex()] = capRegView
	}
	return chainView, nil
}

// Onchain state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	Chains map[uint64]CCIPChainState
}

func (s CCIPOnChainState) View(chains []uint64) (view.CCIPView, error) {
	ccipView := view.NewCCIPView()
	for _, chainSelector := range chains {
		// TODO: Need a utility for this
		chainid, err := chainsel.ChainIdFromSelector(chainSelector)
		if err != nil {
			return ccipView, err
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			return ccipView, err
		}
		if _, ok := s.Chains[chainSelector]; !ok {
			return ccipView, fmt.Errorf("chain not supported %d", chainSelector)
		}
		chainState := s.Chains[chainSelector]
		chainView, err := chainState.GenerateView()
		if err != nil {
			return ccipView, err
		}
		ccipView.Chains[chainName] = chainView
	}
	return ccipView, nil
}

func StateView(e deployment.Environment, ab deployment.AddressBook) (view.CCIPView, error) {
	state, err := LoadOnchainState(e, ab)
	if err != nil {
		return view.CCIPView{}, err
	}
	ccipView, err := state.View(e.AllChainSelectors())
	if err != nil {
		return view.CCIPView{}, err
	}
	ccipView.NodeOperators, err = view.GenerateNopsView(e.NodeIDs, e.Offchain)
	if err != nil {
		return ccipView, err
	}
	return ccipView, nil
}

func LoadOnchainState(e deployment.Environment, ab deployment.AddressBook) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		Chains: make(map[uint64]CCIPChainState),
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := ab.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainState(chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = chainState
	}
	return state, nil
}

// LoadChainState Loads all state for a chain into state
// Modifies map in place
func LoadChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (CCIPChainState, error) {
	var state CCIPChainState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(RBACTimelock, deployment.Version1_0_0).String():
			tl, err := owner_wrappers.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Timelock = tl
		case deployment.NewTypeAndVersion(AdminManyChainMultisig, deployment.Version1_0_0).String():
			mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.AdminMcm = mcms
		case deployment.NewTypeAndVersion(ProposerManyChainMultisig, deployment.Version1_0_0).String():
			mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.ProposerMcm = mcms
		case deployment.NewTypeAndVersion(BypasserManyChainMultisig, deployment.Version1_0_0).String():
			mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.BypasserMcm = mcms
		case deployment.NewTypeAndVersion(CancellerManyChainMultisig, deployment.Version1_0_0).String():
			mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CancellerMcm = mcms
		case deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0).String():
			cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CapabilityRegistry = cr
		case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0_dev).String():
			onRampC, err := onramp.NewOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OnRamp = onRampC
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev).String():
			offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OffRamp = offRamp
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxy = armProxy
		case deployment.NewTypeAndVersion(RMNRemote, deployment.Version1_6_0_dev).String():
			rmnRemote, err := rmn_remote.NewRMNRemote(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNRemote = rmnRemote
		case deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0_dev).String():
			rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNHome = rmnHome
		case deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0).String():
			weth9, err := weth9.NewWETH9(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Weth9 = weth9
		case deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0_dev).String():
			nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.NonceManager = nm
		case deployment.NewTypeAndVersion(CommitStore, deployment.Version1_5_0).String():
			cs, err := commit_store.NewCommitStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CommitStore = cs
		case deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0).String():
			tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenAdminRegistry = tm
		case deployment.NewTypeAndVersion(Router, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Router = r
		case deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TestRouter = r
		case deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0_dev).String():
			fq, err := fee_quoter.NewFeeQuoter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FeeQuoter = fq
		case deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0).String():
			lt, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.LinkToken = lt
		case deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0_dev).String():
			ccipHome, err := ccip_home.NewCCIPHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPHome = ccipHome
		case deployment.NewTypeAndVersion(CCIPConfig, deployment.Version1_0_0).String():
			// TODO: Remove once staging upgraded.
			ccipConfig, err := ccip_config.NewCCIPConfig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPConfig = ccipConfig
		case deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0).String():
			mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Receiver = mr
		case deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0).String():
			feed, err := aggregator_v3_interface.NewAggregatorV3Interface(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDFeeds == nil {
				state.USDFeeds = make(map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface)
			}
			desc, err := feed.Description(&bind.CallOpts{})
			if err != nil {
				return state, err
			}
			key, ok := MockDescriptionToTokenSymbol[desc]
			if !ok {
				return state, fmt.Errorf("unknown feed description %s", desc)
			}
			state.USDFeeds[key] = feed
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
