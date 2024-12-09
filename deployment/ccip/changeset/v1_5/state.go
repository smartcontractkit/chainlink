package v1_5

import (
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

// CCIPOnChainState state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
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
	Router             *router.Router
	Weth9              *weth9.WETH9
	MockRMN            *mock_rmn_contract.MockRMNContract
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token contract
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	// Not all tokens will be burn and mint tokens.
	BurnMintTokens677  map[TokenSymbol]*burn_mint_erc677.BurnMintERC677
	BurnMintTokenPools map[TokenSymbol]*burn_mint_token_pool.BurnMintTokenPool
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Test contracts
	Receiver               *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	TestRouter             *router.Router
	USDCTokenPool          *usdc_token_pool.USDCTokenPool
	MockUSDCTransmitter    *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	MockUSDCTokenMessenger *mock_usdc_token_messenger.MockE2EUSDCTokenMessenger
}
