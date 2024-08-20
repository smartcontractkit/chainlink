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

type CCIPChainState struct {
	EvmOnRampV160      *onramp.OnRamp
	EvmOffRampV160     *offramp.OffRamp
	PriceRegistry      *price_registry.PriceRegistry
	ArmProxy           *rmn_proxy_contract.RMNProxyContract
	NonceManager       *nonce_manager.NonceManager
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	Router             *router.Router
	Weth9              *weth9.WETH9
	MockRmn            *mock_rmn_contract.MockRMNContract
	// TODO: May need to support older link too
	LinkToken *burn_mint_erc677.BurnMintERC677
	// Note we only expect one of these (on the home chain)
	CapabilityRegistry *capabilities_registry.CapabilitiesRegistry
	CCIPConfig         *ccip_config.CCIPConfig
	Mcm                *owner_wrappers.ManyChainMultiSig
	// TODO: remove once we have Address() on wrappers
	McmsAddr common.Address
	Timelock *owner_wrappers.RBACTimelock

	// Test contracts
	Receiver *maybe_revert_message_receiver.MaybeRevertMessageReceiver
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
		if _, ok := s.Chains[chainSelector]; !ok {
			return snapshot, fmt.Errorf("chain not supported %d", chainSelector)
		}
		var c Chain
		ta := s.Chains[chainSelector].TokenAdminRegistry
		if ta != nil {
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
		nm := s.Chains[chainSelector].NonceManager
		if nm != nil {
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
	state, err := LoadOnchainState(e, ab)
	if err != nil {
		return CCIPSnapShot{}, err
	}
	return state.Snapshot(e.AllChainSelectors())
}

func LoadOnchainState(e deployment.Environment, ab deployment.AddressBook) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		Chains: make(map[uint64]CCIPChainState),
	}
	addresses, err := ab.Addresses()
	if err != nil {
		return state, errors.Wrap(err, "could not get addresses")
	}
	for chainSelector, addresses := range addresses {
		chainState, err := LoadChainState(e.Chains[chainSelector], addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = chainState
	}
	return state, nil
}

// Loads all state for a chain into state
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
		case deployment.NewTypeAndVersion(ManyChainMultisig, deployment.Version1_0_0).String():
			mcms, err := owner_wrappers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Mcm = mcms
			state.McmsAddr = common.HexToAddress(address)
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
			state.EvmOnRampV160 = onRampC
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev).String():
			offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.EvmOffRampV160 = offRamp
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.ArmProxy = armProxy
		case deployment.NewTypeAndVersion(MockARM, deployment.Version1_0_0).String():
			mockARM, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockRmn = mockARM
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
		case deployment.NewTypeAndVersion(PriceRegistry, deployment.Version1_6_0_dev).String():
			pr, err := price_registry.NewPriceRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.PriceRegistry = pr
		case deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0).String():
			lt, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.LinkToken = lt
		case deployment.NewTypeAndVersion(CCIPConfig, deployment.Version1_6_0_dev).String():
			cc, err := ccip_config.NewCCIPConfig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPConfig = cc
		case deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0).String():
			mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Receiver = mr
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
