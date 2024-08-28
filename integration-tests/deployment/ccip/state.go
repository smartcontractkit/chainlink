package ccipdeployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	owner_wrappers "github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
)

// Onchain state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	EvmOnRampsV160       map[uint64]*onramp.OnRamp
	EvmOffRampsV160      map[uint64]*offramp.OffRamp
	PriceRegistries      map[uint64]*price_registry.PriceRegistry
	ArmProxies           map[uint64]*rmn_proxy_contract.RMNProxyContract
	NonceManagers        map[uint64]*nonce_manager.NonceManager
	TokenAdminRegistries map[uint64]*token_admin_registry.TokenAdminRegistry
	Routers              map[uint64]*router.Router
	Weth9s               map[uint64]*weth9.WETH9
	MockArms             map[uint64]*mock_rmn_contract.MockRMNContract
	// TODO: May need to support older link too
	LinkTokens map[uint64]*burn_mint_erc677.BurnMintERC677
	// Note we only expect one of these (on the home chain)
	CapabilityRegistry map[uint64]*capabilities_registry.CapabilitiesRegistry
	CCIPConfig         map[uint64]*ccip_config.CCIPConfig
	Mcms               map[uint64]*owner_wrappers.ManyChainMultiSig
	// TODO: remove once we have Address() on wrappers
	McmsAddrs map[uint64]common.Address
	Timelocks map[uint64]*owner_wrappers.RBACTimelock

	// Test contracts
	Receivers map[uint64]*maybe_revert_message_receiver.MaybeRevertMessageReceiver
}

type CCIPSnapShot struct {
	Chains map[string]Chain `json:"chains"`
}

type Contract struct {
	TypeAndVersion string         `json:"typeAndVersion"`
	Address        common.Address `json:"address"`
}

type TokenAdminRegistryView struct {
	Contract
	Tokens []common.Address `json:"tokens"`
}

type NonceManagerView struct {
	Contract
	AuthorizedCallers []common.Address `json:"authorizedCallers"`
}

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	TokenAdminRegistry TokenAdminRegistryView `json:"tokenAdminRegistry"`
	NonceManager       NonceManagerView       `json:"nonceManager"`
}

func (s CCIPOnChainState) Snapshot(chains []uint64) (CCIPSnapShot, error) {
	snapshot := CCIPSnapShot{
		Chains: make(map[string]Chain),
	}
	for _, chainSelector := range chains {
		// TODO: Need a utility for this
		chainid, err := chainsel.ChainIdFromSelector(chainSelector)
		if err != nil {
			return snapshot, err
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			return snapshot, err
		}
		var c Chain
		if ta, ok := s.TokenAdminRegistries[chainSelector]; ok {
			tokens, err := ta.GetAllConfiguredTokens(nil, 0, 10)
			if err != nil {
				return snapshot, err
			}
			tv, err := ta.TypeAndVersion(nil)
			if err != nil {
				return snapshot, err
			}
			c.TokenAdminRegistry = TokenAdminRegistryView{
				Contract: Contract{
					TypeAndVersion: tv,
					Address:        ta.Address(),
				},
				Tokens: tokens,
			}
		}
		if nm, ok := s.NonceManagers[chainSelector]; ok {
			authorizedCallers, err := nm.GetAllAuthorizedCallers(nil)
			if err != nil {
				return snapshot, err
			}
			tv, err := nm.TypeAndVersion(nil)
			if err != nil {
				return snapshot, err
			}
			c.NonceManager = NonceManagerView{
				Contract: Contract{
					TypeAndVersion: tv,
					Address:        nm.Address(),
				},
				// TODO: these can be resolved using an address book
				AuthorizedCallers: authorizedCallers,
			}
		}
		snapshot.Chains[chainName] = c
	}
	return snapshot, nil
}

func SnapshotState(e deployment.Environment, ab deployment.AddressBook) (CCIPSnapShot, error) {
	state, err := GenerateOnchainState(e, ab)
	if err != nil {
		return CCIPSnapShot{}, err
	}
	return state.Snapshot(e.AllChainSelectors())
}

func GenerateOnchainState(e deployment.Environment, ab deployment.AddressBook) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		EvmOnRampsV160:       make(map[uint64]*onramp.OnRamp),
		EvmOffRampsV160:      make(map[uint64]*offramp.OffRamp),
		PriceRegistries:      make(map[uint64]*price_registry.PriceRegistry),
		ArmProxies:           make(map[uint64]*rmn_proxy_contract.RMNProxyContract),
		NonceManagers:        make(map[uint64]*nonce_manager.NonceManager),
		TokenAdminRegistries: make(map[uint64]*token_admin_registry.TokenAdminRegistry),
		Routers:              make(map[uint64]*router.Router),
		MockArms:             make(map[uint64]*mock_rmn_contract.MockRMNContract),
		LinkTokens:           make(map[uint64]*burn_mint_erc677.BurnMintERC677),
		Weth9s:               make(map[uint64]*weth9.WETH9),
		Mcms:                 make(map[uint64]*owner_wrappers.ManyChainMultiSig),
		McmsAddrs:            make(map[uint64]common.Address),
		Timelocks:            make(map[uint64]*owner_wrappers.RBACTimelock),
		CapabilityRegistry:   make(map[uint64]*capabilities_registry.CapabilitiesRegistry),
		CCIPConfig:           make(map[uint64]*ccip_config.CCIPConfig),
		Receivers:            make(map[uint64]*maybe_revert_message_receiver.MaybeRevertMessageReceiver),
	}
	// Get all the onchain state
	addresses, err := ab.Addresses()
	if err != nil {
		return state, errors.Wrap(err, "could not get addresses")
	}
	for chainSelector, addresses := range addresses {
		for address, tvStr := range addresses {
			switch tvStr.String() {
			case deployment.NewTypeAndVersion(RBACTimelock, deployment.Version1_0_0).String():
				tl, err := owner_wrappers.NewRBACTimelock(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Timelocks[chainSelector] = tl
			case deployment.NewTypeAndVersion(ManyChainMultisig, deployment.Version1_0_0).String():
				mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Mcms[chainSelector] = mcms
				state.McmsAddrs[chainSelector] = common.HexToAddress(address)
			case deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0).String():
				cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.CapabilityRegistry[chainSelector] = cr
			case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0_dev).String():
				onRamp, err := onramp.NewOnRamp(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.EvmOnRampsV160[chainSelector] = onRamp
			case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev).String():
				offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.EvmOffRampsV160[chainSelector] = offRamp
			case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
				armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.ArmProxies[chainSelector] = armProxy
			case deployment.NewTypeAndVersion(MockARM, deployment.Version1_0_0).String():
				mockARM, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.MockArms[chainSelector] = mockARM
			case deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0).String():
				weth9, err := weth9.NewWETH9(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Weth9s[chainSelector] = weth9
			case deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0_dev).String():
				nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.NonceManagers[chainSelector] = nm
			case deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0).String():
				tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.TokenAdminRegistries[chainSelector] = tm
			case deployment.NewTypeAndVersion(Router, deployment.Version1_2_0).String():
				r, err := router.NewRouter(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Routers[chainSelector] = r
			case deployment.NewTypeAndVersion(PriceRegistry, deployment.Version1_6_0_dev).String():
				pr, err := price_registry.NewPriceRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.PriceRegistries[chainSelector] = pr
			case deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0).String():
				lt, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.LinkTokens[chainSelector] = lt
			case deployment.NewTypeAndVersion(CCIPConfig, deployment.Version1_6_0_dev).String():
				cc, err := ccip_config.NewCCIPConfig(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.CCIPConfig[chainSelector] = cc
			case deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0).String():
				mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Receivers[chainSelector] = mr
			default:
				return state, fmt.Errorf("unknown contract %s", tvStr)
			}
		}
	}
	return state, nil
}
