package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/multicall3"
)

// CCIPOnChainState_Legacy state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState_Legacy struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	Chains map[uint64]CCIPChainState
}

// CCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	commoncs.MCMSWithTimelockState
	commoncs.LinkTokenState
	commoncs.StaticLinkTokenState
	OnRamp             map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp
	CommitStore        map[uint64]*commit_store.CommitStore
	OffRamp            map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp
	RMNProxy           *rmn_proxy_contract.RMNProxyContract
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	RegistryModule     *registry_module_owner_custom.RegistryModuleOwnerCustom
	FeeQuoter          *fee_quoter.FeeQuoter
	Router             *router.Router
	Weth9              *weth9.WETH9
	MockRMN            *mock_rmn_contract.MockRMNContract
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token contract
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	// Not all tokens will be burn and mint tokens.
	BurnMintTokens677  map[changeset.TokenSymbol]*burn_mint_erc677.BurnMintERC677
	BurnMintTokenPools map[changeset.TokenSymbol]*burn_mint_token_pool.BurnMintTokenPool
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[changeset.TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Test contracts
	Receiver               *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	TestRouter             *router.Router
	USDCTokenPool          *usdc_token_pool.USDCTokenPool
	MockUSDCTransmitter    *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	MockUSDCTokenMessenger *mock_usdc_token_messenger.MockE2EUSDCTokenMessenger
}

func LoadChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (CCIPChainState, error) {
	var state CCIPChainState
	mcmsWithTimelock, err := commoncs.MaybeLoadMCMSWithTimelockState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.MCMSWithTimelockState = *mcmsWithTimelock

	linkState, err := commoncs.MaybeLoadLinkTokenState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.LinkTokenState = *linkState
	staticLinkState, err := commoncs.MaybeLoadStaticLinkTokenState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.StaticLinkTokenState = *staticLinkState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.StaticLinkToken, deployment.Version1_0_0).String():
			// Skip common contracts, they are already loaded.
			continue
		case deployment.NewTypeAndVersion(changeset.OnRamp, deployment.Version1_5_0).String():
			onRampC, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := onRampC.GetStaticConfig(nil)
			if err != nil {
				return state, err
			}
			if state.OnRamp == nil {
				state.OnRamp = make(map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp)
			}
			state.OnRamp[sCfg.DestChainSelector] = onRampC
		case deployment.NewTypeAndVersion(changeset.OffRamp, deployment.Version1_5_0).String():
			offRamp, err := evm(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OffRamp = offRamp
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxyExisting = armProxy
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_6_0_dev).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxyNew = armProxy
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_6_0_dev).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxyContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxyNew = armProxy
		case deployment.NewTypeAndVersion(MockRMN, deployment.Version1_0_0).String():
			mockRMN, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockRMN = mockRMN
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
		case deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_5_0).String():
			rm, err := registry_module_owner_custom.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModule = rm
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
		case deployment.NewTypeAndVersion(USDCToken, deployment.Version1_0_0).String():
			ut, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.BurnMintTokens677 = map[TokenSymbol]*burn_mint_erc677.BurnMintERC677{
				USDCSymbol: ut,
			}
		case deployment.NewTypeAndVersion(USDCTokenPool, deployment.Version1_0_0).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.USDCTokenPool = utp
		case deployment.NewTypeAndVersion(USDCMockTransmitter, deployment.Version1_0_0).String():
			umt, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTransmitter = umt
		case deployment.NewTypeAndVersion(USDCTokenMessenger, deployment.Version1_0_0).String():
			utm, err := mock_usdc_token_messenger.NewMockE2EUSDCTokenMessenger(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTokenMessenger = utm
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
		case deployment.NewTypeAndVersion(Multicall3, deployment.Version1_0_0).String():
			mc, err := multicall3.NewMulticall3(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Multicall3 = mc
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
		case deployment.NewTypeAndVersion(BurnMintTokenPool, deployment.Version1_5_1).String():
			pool, err := burn_mint_token_pool.NewBurnMintTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.BurnMintTokenPools == nil {
				state.BurnMintTokenPools = make(map[TokenSymbol]*burn_mint_token_pool.BurnMintTokenPool)
			}
			tokAddress, err := pool.GetToken(nil)
			if err != nil {
				return state, err
			}
			tok, err := erc20.NewERC20(tokAddress, chain.Client)
			if err != nil {
				return state, err
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, err
			}
			state.BurnMintTokenPools[TokenSymbol(symbol)] = pool
		case deployment.NewTypeAndVersion(BurnMintToken, deployment.Version1_0_0).String():
			tok, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.BurnMintTokens677 == nil {
				state.BurnMintTokens677 = make(map[TokenSymbol]*burn_mint_erc677.BurnMintERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.BurnMintTokens677[TokenSymbol(symbol)] = tok
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
