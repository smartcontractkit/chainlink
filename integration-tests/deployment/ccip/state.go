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

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
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
	EvmOnRampsV160       map[uint64]*evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp
	EvmOffRampsV160      map[uint64]*evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp
	PriceRegistries      map[uint64]*price_registry.PriceRegistry
	ArmProxies           map[uint64]*arm_proxy_contract.ARMProxyContract
	NonceManagers        map[uint64]*nonce_manager.NonceManager
	TokenAdminRegistries map[uint64]*token_admin_registry.TokenAdminRegistry
	Routers              map[uint64]*router.Router
	Weth9s               map[uint64]*weth9.WETH9
	MockArms             map[uint64]*mock_arm_contract.MockARMContract
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

type TokenAdminRegistry struct {
	Contract
	Tokens []common.Address `json:"tokens"`
}

type NonceManager struct {
	Contract
	AuthorizedCallers []common.Address `json:"authorizedCallers"`
}

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	TokenAdminRegistry TokenAdminRegistry `json:"tokenAdminRegistry"`
	NonceManager       NonceManager       `json:"nonceManager"`
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
			c.TokenAdminRegistry = TokenAdminRegistry{
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
			c.NonceManager = NonceManager{
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
		EvmOnRampsV160:       make(map[uint64]*evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp),
		EvmOffRampsV160:      make(map[uint64]*evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp),
		PriceRegistries:      make(map[uint64]*price_registry.PriceRegistry),
		ArmProxies:           make(map[uint64]*arm_proxy_contract.ARMProxyContract),
		NonceManagers:        make(map[uint64]*nonce_manager.NonceManager),
		TokenAdminRegistries: make(map[uint64]*token_admin_registry.TokenAdminRegistry),
		Routers:              make(map[uint64]*router.Router),
		MockArms:             make(map[uint64]*mock_arm_contract.MockARMContract),
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
			switch tvStr {
			case RBAC_Timelock_1_0_0:
				tl, err := owner_wrappers.NewRBACTimelock(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Timelocks[chainSelector] = tl
			case MCMS_1_0_0:
				mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Mcms[chainSelector] = mcms
				state.McmsAddrs[chainSelector] = common.HexToAddress(address)
			case CapabilitiesRegistry_1_0_0:
				cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.CapabilityRegistry[chainSelector] = cr
			case EVM2EVMMultiOnRamp_1_6_0:
				onRamp, err := evm_2_evm_multi_onramp.NewEVM2EVMMultiOnRamp(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.EvmOnRampsV160[chainSelector] = onRamp
			case EVM2EVMMultiOffRamp_1_6_0:
				offRamp, err := evm_2_evm_multi_offramp.NewEVM2EVMMultiOffRamp(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.EvmOffRampsV160[chainSelector] = offRamp
			case ARMProxy_1_1_0:
				armProxy, err := arm_proxy_contract.NewARMProxyContract(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.ArmProxies[chainSelector] = armProxy
			case MockARM_1_0_0:
				mockARM, err := mock_arm_contract.NewMockARMContract(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.MockArms[chainSelector] = mockARM
			case WETH9_1_0_0:
				weth9, err := weth9.NewWETH9(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Weth9s[chainSelector] = weth9
			case NonceManager_1_6_0:
				nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.NonceManagers[chainSelector] = nm
			case TokenAdminRegistry_1_5_0:
				tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.TokenAdminRegistries[chainSelector] = tm
			case Router_1_2_0:
				r, err := router.NewRouter(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.Routers[chainSelector] = r
			case PriceRegistry_1_6_0:
				pr, err := price_registry.NewPriceRegistry(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.PriceRegistries[chainSelector] = pr
			case LinkToken_1_0_0:
				lt, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.LinkTokens[chainSelector] = lt
			case CCIPConfig_1_6_0:
				cc, err := ccip_config.NewCCIPConfig(common.HexToAddress(address), e.Chains[chainSelector].Client)
				if err != nil {
					return state, err
				}
				state.CCIPConfig[chainSelector] = cc
			case CCIPReceiver_1_0_0:
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
