package evm

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/burn_mint_erc677_helper"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/don_id_claimer"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/log_message_data_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool_factory"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	registry_module_owner_custom2 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/multicall3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	shared2 "github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	v1_1 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// CCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	state.MCMSWithTimelockState
	state.LinkTokenState
	state.StaticLinkTokenState
	ABIByAddress       map[string]string
	OnRamp             onramp.OnRampInterface
	OffRamp            offramp.OffRampInterface
	FeeQuoter          *fee_quoter.FeeQuoter
	RMNProxy           *rmn_proxy_contract.RMNProxy
	NonceManager       *nonce_manager.NonceManager
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	TokenPoolFactory   *token_pool_factory.TokenPoolFactory
	RegistryModules1_6 []*registry_module_owner_custom.RegistryModuleOwnerCustom
	// TODO change this to contract object for v1.5 RegistryModules once we have the wrapper available in chainlink-evm
	RegistryModules1_5 []*registry_module_owner_custom2.RegistryModuleOwnerCustom
	Router             *router.Router
	Weth9              *weth9.WETH9
	RMNRemote          *rmn_remote.RMNRemote
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token / token pool contract(s) (only one of which would be active on the registry).
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	ERC20Tokens                map[shared.TokenSymbol]*erc20.ERC20
	FactoryBurnMintERC20Token  *factory_burn_mint_erc20.FactoryBurnMintERC20
	ERC677Tokens               map[shared.TokenSymbol]*erc677.ERC677
	BurnMintTokens677          map[shared.TokenSymbol]*burn_mint_erc677.BurnMintERC677
	BurnMintTokens677Helper    map[shared.TokenSymbol]*burn_mint_erc677_helper.BurnMintERC677Helper
	BurnMintTokenPools         map[shared.TokenSymbol]map[semver.Version]*burn_mint_token_pool.BurnMintTokenPool
	BurnWithFromMintTokenPools map[shared.TokenSymbol]map[semver.Version]*burn_with_from_mint_token_pool.BurnWithFromMintTokenPool
	BurnFromMintTokenPools     map[shared.TokenSymbol]map[semver.Version]*burn_from_mint_token_pool.BurnFromMintTokenPool
	USDCTokenPools             map[semver.Version]*usdc_token_pool.USDCTokenPool
	LockReleaseTokenPools      map[shared.TokenSymbol]map[semver.Version]*lock_release_token_pool.LockReleaseTokenPool
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[shared.TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Note we only expect one of these (on the home chain)
	CapabilityRegistry *capabilities_registry.CapabilitiesRegistry
	CCIPHome           *ccip_home.CCIPHome
	RMNHome            *rmn_home.RMNHome
	DonIDClaimer       *don_id_claimer.DonIDClaimer

	// Test contracts
	Receiver               maybe_revert_message_receiver.MaybeRevertMessageReceiverInterface
	LogMessageDataReceiver *log_message_data_receiver.LogMessageDataReceiver
	TestRouter             *router.Router
	MockUSDCTransmitter    *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	MockUSDCTokenMessenger *mock_usdc_token_messenger.MockE2EUSDCTokenMessenger
	Multicall3             *multicall3.Multicall3

	// Legacy contracts
	EVM2EVMOnRamp  map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp   // mapping of dest chain selector -> EVM2EVMOnRamp
	CommitStore    map[uint64]*commit_store.CommitStore         // mapping of source chain selector -> CommitStore
	EVM2EVMOffRamp map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp // mapping of source chain selector -> EVM2EVMOffRamp
	MockRMN        *mock_rmn_contract.MockRMNContract
	PriceRegistry  *price_registry.PriceRegistry
	RMN            *rmn_contract.RMNContract

	// Treasury contracts
	FeeAggregator common.Address
}

// ValidateHomeChain validates the home chain contracts and their configurations after complete set up
// It cross-references the config across CCIPHome and OffRamps to ensure they are in sync
// This should be called after the complete deployment is done
func (c CCIPChainState) ValidateHomeChain(e cldf.Environment, nodes deployment.Nodes, offRampsByChain map[uint64]offramp.OffRampInterface) error {
	if c.RMNHome == nil {
		return errors.New("no RMNHome contract found in the state for home chain")
	}
	if c.CCIPHome == nil {
		return errors.New("no CCIPHome contract found in the state for home chain")
	}
	if c.CapabilityRegistry == nil {
		return errors.New("no CapabilityRegistry contract found in the state for home chain")
	}
	// get capReg from CCIPHome
	capReg, err := c.CCIPHome.GetCapabilityRegistry(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get capability registry from CCIPHome contract: %w", err)
	}
	if capReg != c.CapabilityRegistry.Address() {
		return fmt.Errorf("capability registry mismatch: expected %s, got %s", capReg.Hex(), c.CapabilityRegistry.Address().Hex())
	}
	ccipDons, err := shared.GetCCIPDonsFromCapRegistry(e.GetContext(), c.CapabilityRegistry)
	if err != nil {
		return fmt.Errorf("failed to get CCIP Dons from capability registry: %w", err)
	}
	if len(ccipDons) == 0 {
		return errors.New("no CCIP Dons found in capability registry")
	}
	// validate for all ccipDons
	for _, don := range ccipDons {
		if err := nodes.P2PIDsPresentInJD(don.NodeP2PIds); err != nil {
			return fmt.Errorf("failed to find Capability Registry p2pIDs in JD: %w", err)
		}
		commitConfig, err := c.CCIPHome.GetAllConfigs(&bind.CallOpts{
			Context: e.GetContext(),
		}, don.Id, uint8(types.PluginTypeCCIPCommit))
		if err != nil {
			return fmt.Errorf("failed to get commit config for don %d: %w", don.Id, err)
		}
		if err := c.validateCCIPHomeVersionedActiveConfig(e, nodes, commitConfig.ActiveConfig, offRampsByChain); err != nil {
			return fmt.Errorf("failed to validate active commit config for don %d: %w", don.Id, err)
		}
		execConfig, err := c.CCIPHome.GetAllConfigs(&bind.CallOpts{
			Context: e.GetContext(),
		}, don.Id, uint8(types.PluginTypeCCIPExec))
		if err != nil {
			return fmt.Errorf("failed to get exec config for don %d: %w", don.Id, err)
		}
		if err := c.validateCCIPHomeVersionedActiveConfig(e, nodes, execConfig.ActiveConfig, offRampsByChain); err != nil {
			return fmt.Errorf("failed to validate active exec config for don %d: %w", don.Id, err)
		}
	}
	return nil
}

// validateCCIPHomeVersionedActiveConfig validates the CCIPHomeVersionedConfig based on corresponding chain selector and its state
// The validation related to correctness of F and node length is omitted here as it is already validated in the contract
func (c CCIPChainState) validateCCIPHomeVersionedActiveConfig(e cldf.Environment, nodes deployment.Nodes, homeCfg ccip_home.CCIPHomeVersionedConfig, offRampsByChain map[uint64]offramp.OffRampInterface) error {
	if homeCfg.ConfigDigest == [32]byte{} {
		return errors.New("active config digest is empty")
	}
	chainSel := homeCfg.Config.ChainSelector
	if _, exists := e.BlockChains.SolanaChains()[chainSel]; exists {
		return nil
	}
	if _, exists := e.BlockChains.AptosChains()[chainSel]; exists {
		return nil
	}
	offRamp, ok := offRampsByChain[chainSel]
	if !ok {
		return fmt.Errorf("offRamp for chain %d not found in the state", chainSel)
	}
	// validate ChainConfig in CCIPHome
	homeChainConfig, err := c.CCIPHome.GetChainConfig(&bind.CallOpts{
		Context: e.GetContext(),
	}, chainSel)
	if err != nil {
		return fmt.Errorf("failed to get home chain config for chain %d: %w", chainSel, err)
	}
	// Node details should match with what we fetch from JD for CCIP Home Readers
	if err := nodes.P2PIDsPresentInJD(homeChainConfig.Readers); err != nil {
		return fmt.Errorf("failed to find homechain readers in JD for chain %d: %w",
			chainSel, err)
	}

	// Validate CCIPHome OCR3 Related Config
	if offRamp.Address() != common.BytesToAddress(homeCfg.Config.OfframpAddress) {
		return fmt.Errorf("offRamp address mismatch in active config for ccip home for chain %d: expected %s, got %s",
			chainSel, offRamp.Address().Hex(), homeCfg.Config.OfframpAddress)
	}
	if c.RMNHome.Address() != common.BytesToAddress(homeCfg.Config.RmnHomeAddress) {
		return fmt.Errorf("RMNHome address mismatch in active config for ccip home for chain %d: expected %s, got %s",
			chainSel, c.RMNHome.Address().Hex(), homeCfg.Config.RmnHomeAddress)
	}
	p2pIDs := make([][32]byte, 0)
	for _, node := range homeCfg.Config.Nodes {
		p2pIDs = append(p2pIDs, node.P2pId)
	}
	if err := nodes.P2PIDsPresentInJD(p2pIDs); err != nil {
		return fmt.Errorf("failed to find p2pIDs from CCIPHome config in JD for chain %d: %w", chainSel, err)
	}
	// cross-check with offRamp whether all in sync
	switch homeCfg.Config.PluginType {
	case uint8(types.PluginTypeCCIPCommit):
		commitConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: e.GetContext(),
		}, uint8(types.PluginTypeCCIPCommit))
		if err != nil {
			return fmt.Errorf("failed to get commit config for chain %d offRamp %s: %w", chainSel, c.OffRamp.Address().Hex(), err)
		}
		// the config digest should match with CCIP Home ActiveConfig
		if commitConfig.ConfigInfo.ConfigDigest != homeCfg.ConfigDigest {
			return fmt.Errorf("offRamp %s commit config digest mismatch with CCIPHome for chain %d: expected %x, got %x",
				offRamp.Address().Hex(), chainSel, homeCfg.ConfigDigest, commitConfig.ConfigInfo.ConfigDigest)
		}
		if !commitConfig.ConfigInfo.IsSignatureVerificationEnabled {
			return fmt.Errorf("offRamp %s for chain %d commit config signature verification is not enabled",
				offRamp.Address().Hex(), chainSel)
		}
		if err := validateLatestConfigOffRamp(offRamp, commitConfig, homeChainConfig); err != nil {
			return fmt.Errorf("offRamp %s for chain %d commit config validation error: %w",
				offRamp.Address().Hex(), chainSel, err)
		}
	case uint8(types.PluginTypeCCIPExec):
		execConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: e.GetContext(),
		}, uint8(types.PluginTypeCCIPExec))
		if err != nil {
			return fmt.Errorf("failed to get exec config for chain %d offRamp %s: %w", chainSel, offRamp.Address().Hex(), err)
		}
		// the config digest should match with CCIP Home ActiveConfig
		if execConfig.ConfigInfo.ConfigDigest != homeCfg.ConfigDigest {
			return fmt.Errorf("offRamp %s exec config digest mismatch with CCIPHome for chain %d: expected %x, got %x",
				offRamp.Address().Hex(), chainSel, homeCfg.ConfigDigest, execConfig.ConfigInfo.ConfigDigest)
		}
		if execConfig.ConfigInfo.IsSignatureVerificationEnabled {
			return fmt.Errorf("offRamp %s for chain %d exec config signature verification is enabled",
				offRamp.Address().Hex(), chainSel)
		}
		if err := validateLatestConfigOffRamp(offRamp, execConfig, homeChainConfig); err != nil {
			return fmt.Errorf("offRamp %s for chain %d exec config validation error: %w",
				offRamp.Address().Hex(), chainSel, err)
		}
	default:
		return fmt.Errorf("unsupported plugin type %d for chain %d", homeCfg.Config.PluginType, chainSel)
	}
	return nil
}

// ValidateOnRamp validates whether the contract addresses configured in static and dynamic config are in sync with state
func (c CCIPChainState) ValidateOnRamp(
	e cldf.Environment,
	selector uint64,
	connectedChains []uint64,
) error {
	if c.OnRamp == nil {
		return errors.New("no OnRamp contract found in the state")
	}
	staticCfg, err := c.OnRamp.GetStaticConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return err
	}
	if staticCfg.ChainSelector != selector {
		return fmt.Errorf("onRamp %s chainSelector mismatch in static config: expected %d, got %d",
			c.OnRamp.Address().Hex(), selector, staticCfg.ChainSelector)
	}
	// it should be RMNProxy pointing to the RMNRemote
	if c.RMNProxy.Address() != staticCfg.RmnRemote {
		return fmt.Errorf("onRamp %s RMNRemote mismatch in static config: expected %s, got %s",
			c.OnRamp.Address().Hex(), c.RMNRemote.Address().Hex(), staticCfg.RmnRemote)
	}
	if c.NonceManager.Address() != staticCfg.NonceManager {
		return fmt.Errorf("onRamp %s NonceManager mismatch in static config: expected %s, got %s",
			c.OnRamp.Address().Hex(), c.NonceManager.Address().Hex(), staticCfg.NonceManager)
	}
	if c.TokenAdminRegistry.Address() != staticCfg.TokenAdminRegistry {
		return fmt.Errorf("onRamp %s TokenAdminRegistry mismatch in static config: expected %s, got %s",
			c.OnRamp.Address().Hex(), c.TokenAdminRegistry.Address().Hex(), staticCfg.TokenAdminRegistry)
	}
	dynamicCfg, err := c.OnRamp.GetDynamicConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get dynamic config for chain %d onRamp %s: %w", selector, c.OnRamp.Address().Hex(), err)
	}
	if dynamicCfg.FeeQuoter != c.FeeQuoter.Address() {
		return fmt.Errorf("onRamp %s feeQuoter mismatch in dynamic config: expected %s, got %s",
			c.OnRamp.Address().Hex(), c.FeeQuoter.Address().Hex(), dynamicCfg.FeeQuoter.Hex())
	}
	// if the fee aggregator is set, it should match the one in the dynamic config
	// otherwise the fee aggregator should be the timelock address
	if c.FeeAggregator != (common.Address{}) {
		if c.FeeAggregator != dynamicCfg.FeeAggregator {
			return fmt.Errorf("onRamp %s feeAggregator mismatch in dynamic config: expected %s, got %s",
				c.OnRamp.Address().Hex(), c.FeeAggregator.Hex(), dynamicCfg.FeeAggregator.Hex())
		}
	} else {
		if dynamicCfg.FeeAggregator != e.BlockChains.EVMChains()[selector].DeployerKey.From {
			return fmt.Errorf("onRamp %s feeAggregator mismatch in dynamic config: expected deployer key %s, got %s",
				c.OnRamp.Address().Hex(), e.BlockChains.EVMChains()[selector].DeployerKey.From.Hex(), dynamicCfg.FeeAggregator.Hex())
		}
	}

	for _, otherChainSel := range connectedChains {
		destChainCfg, err := c.OnRamp.GetDestChainConfig(&bind.CallOpts{
			Context: e.GetContext(),
		}, otherChainSel)
		if err != nil {
			return fmt.Errorf("failed to get dest chain config from source chain %d onRamp %s for dest chain %d: %w",
				selector, c.OnRamp.Address(), otherChainSel, err)
		}
		// if not blank, the dest chain config should be enabled
		if destChainCfg != (onramp.GetDestChainConfig{}) {
			if destChainCfg.Router != c.Router.Address() && destChainCfg.Router != c.TestRouter.Address() {
				return fmt.Errorf("onRamp %s router mismatch in dest chain config: expected router %s or test router %s, got %s",
					c.OnRamp.Address().Hex(), c.Router.Address().Hex(), c.TestRouter.Address().Hex(), destChainCfg.Router.Hex())
			}
		}
	}

	return nil
}

// ValidateFeeQuoter validates whether the fee quoter contract address configured in static config is in sync with state
func (c CCIPChainState) ValidateFeeQuoter(e cldf.Environment) error {
	if c.FeeQuoter == nil {
		return errors.New("no FeeQuoter contract found in the state")
	}
	staticConfig, err := c.FeeQuoter.GetStaticConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get static config for FeeQuoter %s: %w", c.FeeQuoter.Address().Hex(), err)
	}
	linktokenAddr, err := c.LinkTokenAddress()
	if err != nil {
		return fmt.Errorf("failed to get link token address for from state: %w", err)
	}
	if staticConfig.LinkToken != linktokenAddr {
		return fmt.Errorf("feeQuoter %s LinkToken mismatch: expected either linktoken %s or static link token %s, got %s",
			c.FeeQuoter.Address().Hex(), c.LinkToken.Address().Hex(), c.StaticLinkToken.Address(), staticConfig.LinkToken.Hex())
	}
	return nil
}

// ValidateRouter validates the router contract to check if all wired contracts are synced with state
// and returns all connected chains with respect to the router
func (c CCIPChainState) ValidateRouter(e cldf.Environment, isTestRouter bool) ([]uint64, error) {
	if c.Router == nil && c.TestRouter == nil {
		return nil, errors.New("no Router or TestRouter contract found in the state")
	}
	routerC := c.Router
	if isTestRouter {
		routerC = c.TestRouter
	}
	armProxy, err := routerC.GetArmProxy(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get armProxy from router : %w", err)
	}
	if armProxy != c.RMNProxy.Address() {
		return nil, fmt.Errorf("armProxy %s mismatch in router %s: expected %s, got %s",
			armProxy.Hex(), routerC.Address().Hex(), c.RMNProxy.Address().Hex(), armProxy)
	}
	native, err := routerC.GetWrappedNative(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get wrapped native from router %s: %w", routerC.Address().Hex(), err)
	}
	if native != c.Weth9.Address() {
		return nil, fmt.Errorf("wrapped native %s mismatch in router %s: expected %s, got %s",
			native.Hex(), routerC.Address().Hex(), c.Weth9.Address().Hex(), native)
	}
	allConnectedChains := make([]uint64, 0)
	// get offRamps
	offRampDetails, err := routerC.GetOffRamps(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get offRamps from router %s: %w", routerC.Address().Hex(), err)
	}
	for _, d := range offRampDetails {
		// skip if solana - solana state is maintained in solana
		if _, exists := e.BlockChains.SolanaChains()[d.SourceChainSelector]; exists {
			continue
		}
		allConnectedChains = append(allConnectedChains, d.SourceChainSelector)
		// check if offRamp is valid
		if d.OffRamp != c.OffRamp.Address() {
			return nil, fmt.Errorf("offRamp %s mismatch for source %d in router %s: expected %s, got %s",
				d.OffRamp.Hex(), d.SourceChainSelector, routerC.Address().Hex(), c.OffRamp.Address().Hex(), d.OffRamp)
		}
	}
	// all lanes are bi-directional, if we have a lane from A to B, we also have a lane from B to A
	// source to offRamp should be same as dest to onRamp
	for _, dest := range allConnectedChains {
		onRamp, err := routerC.GetOnRamp(&bind.CallOpts{
			Context: context.Background(),
		}, dest)
		if err != nil {
			return nil, fmt.Errorf("failed to get onRamp for dest %d from router %s: %w", dest, routerC.Address().Hex(), err)
		}
		if onRamp != c.OnRamp.Address() {
			return nil, fmt.Errorf("onRamp %s mismatch for dest chain %d in router %s: expected %s, got %s",
				onRamp.Hex(), dest, routerC.Address().Hex(), c.OnRamp.Address().Hex(), onRamp)
		}
	}
	return allConnectedChains, nil
}

// ValidateRMNRemote validates the RMNRemote contract to check if all wired contracts are synced with state
// and returns whether RMN is enabled for the chain on the RMNRemote
// It validates whether RMNRemote is in sync with the RMNHome contract
func (c CCIPChainState) ValidateRMNRemote(
	e cldf.Environment,
	selector uint64,
	rmnHomeActiveDigest [32]byte,
) (bool, error) {
	if c.RMNRemote == nil {
		return false, errors.New("no RMNRemote contract found in the state")
	}
	chainSelector, err := c.RMNRemote.GetLocalChainSelector(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get local chain selector from RMNRemote %s: %w", c.RMNRemote.Address().Hex(), err)
	}
	if chainSelector != selector {
		return false, fmt.Errorf("RMNRemote %s chainSelector mismatch: expected %d, got %d",
			c.RMNRemote.Address().Hex(), selector, chainSelector)
	}
	versionedCfg, err := c.RMNRemote.GetVersionedConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get versioned config from RMNRemote %s: %w", c.RMNRemote.Address().Hex(), err)
	}
	if versionedCfg.Version == 0 {
		return false, errors.New("RMNRemote config is not set")
	}
	if versionedCfg.Config.RmnHomeContractConfigDigest != rmnHomeActiveDigest {
		return false, fmt.Errorf("RMNRemote %s config digest mismatch: expected %x, got %x",
			c.RMNRemote.Address().Hex(), rmnHomeActiveDigest, versionedCfg.Config.RmnHomeContractConfigDigest)
	}
	return versionedCfg.Config.FSign > 0, nil
}

// ValidateOffRamp validates the offRamp contract to check if all wired contracts are synced with state
func (c CCIPChainState) ValidateOffRamp(
	e cldf.Environment,
	selector uint64,
	onRampsBySelector map[uint64]common.Address,
	isRMNEnabledBySource map[uint64]bool,
) error {
	if c.OffRamp == nil {
		return errors.New("no OffRamp contract found in the state")
	}
	// staticConfig chainSelector matches the selector key for the CCIPChainState
	staticConfig, err := c.OffRamp.GetStaticConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get static config for chain %d offRammp %s: %w", selector, c.OffRamp.Address().Hex(), err)
	}
	// staticConfig chainSelector should match the selector key for the CCIPChainState
	if staticConfig.ChainSelector != selector {
		return fmt.Errorf("offRamp %s chainSelector mismatch: expected %d, got %d",
			c.OffRamp.Address().Hex(), selector, staticConfig.ChainSelector)
	}
	// RMNProxy address for chain should be the same as the one in the static config for RMNRemote
	if c.RMNProxy.Address() != staticConfig.RmnRemote {
		return fmt.Errorf("offRamp %s RMNRemote mismatch: expected %s, got %s",
			c.OffRamp.Address().Hex(), c.RMNRemote.Address().Hex(), staticConfig.RmnRemote)
	}
	// NonceManager address for chain should be the same as the one in the static config
	if c.NonceManager.Address() != staticConfig.NonceManager {
		return fmt.Errorf("offRamp %s NonceManager mismatch: expected %s, got %s",
			c.OffRamp.Address().Hex(), c.NonceManager.Address().Hex(), staticConfig.NonceManager)
	}
	// TokenAdminRegistry address for chain should be the same as the one in the static config
	if c.TokenAdminRegistry.Address() != staticConfig.TokenAdminRegistry {
		return fmt.Errorf("offRamp %s TokenAdminRegistry mismatch: expected %s, got %s",
			c.OffRamp.Address().Hex(), c.TokenAdminRegistry.Address().Hex(), staticConfig.TokenAdminRegistry)
	}
	dynamicConfig, err := c.OffRamp.GetDynamicConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get dynamic config for chain %d offRamp %s: %w", selector, c.OffRamp.Address().Hex(), err)
	}
	// FeeQuoter address for chain should be the same as the one in the static config
	if dynamicConfig.FeeQuoter != c.FeeQuoter.Address() {
		return fmt.Errorf("offRamp %s feeQuoter mismatch: expected %s, got %s",
			c.OffRamp.Address().Hex(), c.FeeQuoter.Address().Hex(), dynamicConfig.FeeQuoter.Hex())
	}
	if dynamicConfig.PermissionLessExecutionThresholdSeconds != uint32(globals.PermissionLessExecutionThreshold.Seconds()) {
		return fmt.Errorf("offRamp %s permissionless execution threshold mismatch: expected %f, got %d",
			c.OffRamp.Address().Hex(), globals.PermissionLessExecutionThreshold.Seconds(), dynamicConfig.PermissionLessExecutionThresholdSeconds)
	}
	for chainSel, srcChainOnRamp := range onRampsBySelector {
		config, err := c.OffRamp.GetSourceChainConfig(&bind.CallOpts{
			Context: e.GetContext(),
		}, chainSel)
		if err != nil {
			return fmt.Errorf("failed to get source chain config for chain %d: %w", chainSel, err)
		}
		if config.IsEnabled {
			// For all configured sources, the address of configured onRamp for chain A must be the Address() of the onramp on chain A
			if srcChainOnRamp != common.BytesToAddress(config.OnRamp) {
				return fmt.Errorf("onRamp address mismatch for source chain %d on OffRamp %s : expected %s, got %x",
					chainSel, c.OffRamp.Address().Hex(), srcChainOnRamp.Hex(), config.OnRamp)
			}
			// The address of router should be accurate
			if c.Router.Address() != config.Router && c.TestRouter.Address() != config.Router {
				return fmt.Errorf("router address mismatch for source chain %d on OffRamp %s : expected either router %s or test router %s, got %s",
					chainSel, c.OffRamp.Address().Hex(), c.Router.Address().Hex(), c.TestRouter.Address().Hex(), config.Router.Hex())
			}
			// if RMN is enabled for the source chain, the RMNRemote and RMNHome should be configured to enable RMN
			// the reverse is not always true, as RMN verification can be disable at offRamp but enabled in RMNRemote and RMNHome
			if !config.IsRMNVerificationDisabled && !isRMNEnabledBySource[chainSel] {
				return fmt.Errorf("RMN verification is enabled in offRamp %s for source chain %d, "+
					"but RMN is not enabled in RMNHome and RMNRemote for the chain",
					c.OffRamp.Address().Hex(), chainSel)
			}
		}
	}
	return nil
}

func (c CCIPChainState) TokenAddressBySymbol() (map[shared.TokenSymbol]common.Address, error) {
	tokenAddresses := make(map[shared.TokenSymbol]common.Address)
	if c.FactoryBurnMintERC20Token != nil {
		tokenAddresses[shared.FactoryBurnMintERC20Symbol] = c.FactoryBurnMintERC20Token.Address()
	}
	for symbol, token := range c.ERC20Tokens {
		tokenAddresses[symbol] = token.Address()
	}
	for symbol, token := range c.ERC677Tokens {
		tokenAddresses[symbol] = token.Address()
	}
	for symbol, token := range c.BurnMintTokens677 {
		tokenAddresses[symbol] = token.Address()
	}
	var err error
	tokenAddresses[shared.LinkSymbol], err = c.LinkTokenAddress()
	if err != nil {
		return nil, err
	}
	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenAddresses[shared.WethSymbol] = c.Weth9.Address()
	return tokenAddresses, nil
}

// TokenDetailsBySymbol get token mapping from the state. It contains only tokens that we have in address book
func (c CCIPChainState) TokenDetailsBySymbol() (map[shared.TokenSymbol]shared.TokenDetails, error) {
	tokenDetails := make(map[shared.TokenSymbol]shared.TokenDetails)
	if c.FactoryBurnMintERC20Token != nil {
		tokenDetails[shared.FactoryBurnMintERC20Symbol] = c.FactoryBurnMintERC20Token
	}
	for symbol, token := range c.ERC20Tokens {
		tokenDetails[symbol] = token
	}
	for symbol, token := range c.ERC677Tokens {
		tokenDetails[symbol] = token
	}
	for symbol, token := range c.BurnMintTokens677 {
		tokenDetails[symbol] = token
	}
	if c.LinkToken != nil {
		tokenDetails[shared.LinkSymbol] = c.LinkToken
	}
	if c.StaticLinkToken != nil {
		tokenDetails[shared.LinkSymbol] = c.StaticLinkToken
	}

	if _, ok := tokenDetails[shared.LinkSymbol]; !ok {
		return nil, errors.New("no LINK contract found in the state")
	}

	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenDetails[shared.WethSymbol] = c.Weth9
	return tokenDetails, nil
}

func (c CCIPChainState) LinkTokenAddress() (common.Address, error) {
	if c.LinkToken != nil {
		return c.LinkToken.Address(), nil
	}
	if c.StaticLinkToken != nil {
		return c.StaticLinkToken.Address(), nil
	}
	return common.Address{}, errors.New("no link token found in the state")
}

func (c CCIPChainState) GenerateView(lggr logger.Logger, chain string) (view.ChainView, error) {
	chainView := view.NewChain()
	grp := errgroup.Group{}

	if c.Router != nil {
		grp.Go(func() error {
			routerView, err := v1_2.GenerateRouterView(c.Router, false)
			if err != nil {
				return errors.Wrapf(err, "failed to generate router view for router %s", c.Router.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.Router[c.Router.Address().Hex()] = routerView
			lggr.Infow("generated router view", "router", c.Router.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.TestRouter != nil {
		grp.Go(func() error {
			testRouterView, err := v1_2.GenerateRouterView(c.TestRouter, true)
			if err != nil {
				return errors.Wrapf(err, "failed to generate router view for test router %s", c.TestRouter.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.Router[c.TestRouter.Address().Hex()] = testRouterView
			lggr.Infow("generated test router view", "testRouter", c.TestRouter.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.TokenAdminRegistry != nil {
		grp.Go(func() error {
			lggr.Infow("generating token admin registry view, this might take a while based on number of tokens",
				"tokenAdminRegistry", c.TokenAdminRegistry.Address().Hex(), "chain", chain)
			taView, err := v1_5.GenerateTokenAdminRegistryView(c.TokenAdminRegistry)
			if err != nil {
				return errors.Wrapf(err, "failed to generate token admin registry view for token admin registry %s", c.TokenAdminRegistry.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.TokenAdminRegistry[c.TokenAdminRegistry.Address().Hex()] = taView
			lggr.Infow("generated token admin registry view", "tokenAdminRegistry", c.TokenAdminRegistry.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.TokenPoolFactory != nil {
		grp.Go(func() error {
			tpfView, err := v1_5_1.GenerateTokenPoolFactoryView(c.TokenPoolFactory)
			if err != nil {
				return errors.Wrapf(err, "failed to generate token pool factory view for token pool factory %s", c.TokenPoolFactory.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.TokenPoolFactory[c.TokenPoolFactory.Address().Hex()] = tpfView
			lggr.Infow("generated token pool factory view", "tokenPoolFactory", c.TokenPoolFactory.Address().Hex(), "chain", chain)
			return nil
		})
	}
	for tokenSymbol, versionToPool := range c.BurnMintTokenPools {
		for _, tokenPool := range versionToPool {
			grp.Go(func() error {
				tokenPoolView, err := v1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), v1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				lggr.Infow("generated burn mint token pool view", "tokenPool", tokenPool.Address().Hex(), "chain", chain)
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.BurnWithFromMintTokenPools {
		for _, tokenPool := range versionToPool {
			grp.Go(func() error {
				tokenPoolView, err := v1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), v1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				lggr.Infow("generated burn mint token pool view", "tokenPool", tokenPool.Address().Hex(), "chain", chain)
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.BurnFromMintTokenPools {
		for _, tokenPool := range versionToPool {
			grp.Go(func() error {
				tokenPoolView, err := v1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), v1_5_1.PoolView{
					TokenPoolView: tokenPoolView,
				})
				lggr.Infow("generated burn mint token pool view", "tokenPool", tokenPool.Address().Hex(), "chain", chain)
				return nil
			})
		}
	}
	for tokenSymbol, versionToPool := range c.LockReleaseTokenPools {
		for _, tokenPool := range versionToPool {
			grp.Go(func() error {
				tokenPoolView, err := v1_5_1.GenerateLockReleaseTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate lock release token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), tokenPoolView)
				lggr.Infow("generated lock release token pool view", "tokenPool", tokenPool.Address().Hex(), "chain", chain)
				return nil
			})
		}
	}
	for _, pool := range c.USDCTokenPools {
		grp.Go(func() error {
			tokenPoolView, err := v1_5_1.GenerateUSDCTokenPoolView(pool)
			if err != nil {
				return errors.Wrapf(err, "failed to generate USDC token pool view for %s", pool.Address().String())
			}
			chainView.UpdateTokenPool(string(shared.USDCSymbol), pool.Address().Hex(), tokenPoolView)
			lggr.Infow("generated USDC token pool view", "tokenPool", pool.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.NonceManager != nil {
		grp.Go(func() error {
			nmView, err := v1_6.GenerateNonceManagerView(c.NonceManager)
			if err != nil {
				return errors.Wrapf(err, "failed to generate nonce manager view for nonce manager %s", c.NonceManager.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.NonceManager[c.NonceManager.Address().Hex()] = nmView
			lggr.Infow("generated nonce manager view", "nonceManager", c.NonceManager.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.RMNRemote != nil {
		grp.Go(func() error {
			rmnView, err := v1_6.GenerateRMNRemoteView(c.RMNRemote)
			if err != nil {
				return errors.Wrapf(err, "failed to generate rmn remote view for rmn remote %s", c.RMNRemote.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.RMNRemote[c.RMNRemote.Address().Hex()] = rmnView
			lggr.Infow("generated rmn remote view", "rmnRemote", c.RMNRemote.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.RMNHome != nil {
		grp.Go(func() error {
			rmnHomeView, err := v1_6.GenerateRMNHomeView(c.RMNHome)
			if err != nil {
				return errors.Wrapf(err, "failed to generate rmn home view for rmn home %s", c.RMNHome.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.RMNHome[c.RMNHome.Address().Hex()] = rmnHomeView
			lggr.Infow("generated rmn home view", "rmnHome", c.RMNHome.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.FeeQuoter != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		grp.Go(func() error {
			// FeeQuoter knows only about tokens that managed by CCIP (i.e. imported from address book)
			tokenDetails, err := c.TokenDetailsBySymbol()
			if err != nil {
				return err
			}
			tokens := make([]common.Address, 0, len(tokenDetails))
			for _, tokenDetail := range tokenDetails {
				tokens = append(tokens, tokenDetail.Address())
			}
			fqView, err := v1_6.GenerateFeeQuoterView(c.FeeQuoter, c.Router, c.TestRouter, tokens)
			if err != nil {
				return errors.Wrapf(err, "failed to generate fee quoter view for fee quoter %s", c.FeeQuoter.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.FeeQuoter[c.FeeQuoter.Address().Hex()] = fqView
			lggr.Infow("generated fee quoter view", "feeQuoter", c.FeeQuoter.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.OnRamp != nil && c.Router != nil && c.TokenAdminRegistry != nil {
		grp.Go(func() error {
			onRampView, err := v1_6.GenerateOnRampView(
				c.OnRamp,
				c.Router,
				c.TestRouter,
				c.TokenAdminRegistry,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to generate on ramp view for on ramp %s", c.OnRamp.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.OnRamp[c.OnRamp.Address().Hex()] = onRampView
			lggr.Infow("generated on ramp view", "onRamp", c.OnRamp.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.OffRamp != nil && c.Router != nil {
		grp.Go(func() error {
			offRampView, err := v1_6.GenerateOffRampView(
				c.OffRamp,
				c.Router,
				c.TestRouter,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to generate off ramp view for off ramp %s", c.OffRamp.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.OffRamp[c.OffRamp.Address().Hex()] = offRampView
			lggr.Infow("generated off ramp view", "offRamp", c.OffRamp.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.RMNProxy != nil {
		grp.Go(func() error {
			rmnProxyView, err := v1_0.GenerateRMNProxyView(c.RMNProxy)
			if err != nil {
				return errors.Wrapf(err, "failed to generate rmn proxy view for rmn proxy %s", c.RMNProxy.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.RMNProxy[c.RMNProxy.Address().Hex()] = rmnProxyView
			lggr.Infow("generated rmn proxy view", "rmnProxy", c.RMNProxy.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.CCIPHome != nil && c.CapabilityRegistry != nil {
		grp.Go(func() error {
			chView, err := v1_6.GenerateCCIPHomeView(c.CapabilityRegistry, c.CCIPHome)
			if err != nil {
				return errors.Wrapf(err, "failed to generate CCIP home view for CCIP home %s", c.CCIPHome.Address())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.CCIPHome[c.CCIPHome.Address().Hex()] = chView
			lggr.Infow("generated CCIP home view", "CCIPHome", c.CCIPHome.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.CapabilityRegistry != nil {
		grp.Go(func() error {
			capRegView, err := v1_1.GenerateCapabilityRegistryView(c.CapabilityRegistry)
			if err != nil {
				return errors.Wrapf(err, "failed to generate capability registry view for capability registry %s", c.CapabilityRegistry.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.CapabilityRegistry[c.CapabilityRegistry.Address().Hex()] = capRegView
			lggr.Infow("generated capability registry view", "capabilityRegistry", c.CapabilityRegistry.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.MCMSWithTimelockState.Timelock != nil {
		grp.Go(func() error {
			mcmsView, err := c.MCMSWithTimelockState.GenerateMCMSWithTimelockView()
			if err != nil {
				return errors.Wrapf(err, "failed to generate MCMS with timelock view for MCMS with timelock %s", c.MCMSWithTimelockState.Timelock.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.MCMSWithTimelock = mcmsView
			lggr.Infow("generated MCMS with timelock view", "MCMSWithTimelock", c.MCMSWithTimelockState.Timelock.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.LinkToken != nil {
		grp.Go(func() error {
			linkTokenView, err := c.GenerateLinkView()
			if err != nil {
				return errors.Wrapf(err, "failed to generate link token view for link token %s", c.LinkToken.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.LinkToken = linkTokenView
			lggr.Infow("generated link token view", "linkToken", c.LinkToken.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.StaticLinkToken != nil {
		grp.Go(func() error {
			staticLinkTokenView, err := c.GenerateStaticLinkView()
			if err != nil {
				return err
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.StaticLinkToken = staticLinkTokenView
			lggr.Infow("generated static link token view", "staticLinkToken", c.StaticLinkToken.Address().Hex(), "chain", chain)
			return nil
		})
	}

	// Legacy contracts
	// OnRamp, OffRamp, CommitStore legacy contract related state generation is not done right now
	// considering the state of these contracts are not referred currently, and it's enormously expensive to generate
	// state for multiple lanes per chain
	for _, registryModule := range c.RegistryModules1_6 {
		grp.Go(func() error {
			registryModuleView, err := shared2.GetRegistryModuleView(registryModule, c.TokenAdminRegistry.Address())
			if err != nil {
				return errors.Wrapf(err, "failed to generate registry module view for registry module %s", registryModule.Address().Hex())
			}
			chainView.UpdateRegistryModuleView(registryModule.Address().Hex(), registryModuleView)
			lggr.Infow("generated registry module view", "registryModule", registryModule.Address().Hex(), "chain", chain)
			return nil
		})
	}

	for _, registryModule := range c.RegistryModules1_5 {
		grp.Go(func() error {
			registryModuleView, err := shared2.GetRegistryModuleView(registryModule, c.TokenAdminRegistry.Address())
			if err != nil {
				return errors.Wrapf(err, "failed to generate registry module view for registry module %s", registryModule.Address().Hex())
			}
			chainView.UpdateRegistryModuleView(registryModule.Address().Hex(), registryModuleView)
			lggr.Infow("generated registry module view", "registryModule", registryModule.Address().Hex(), "chain", chain)
			return nil
		})
	}

	if c.PriceRegistry != nil {
		grp.Go(func() error {
			priceRegistryView, err := v1_2.GeneratePriceRegistryView(c.PriceRegistry)
			if err != nil {
				return errors.Wrapf(err, "failed to generate price registry view for price registry %s", c.PriceRegistry.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.PriceRegistry[c.PriceRegistry.Address().String()] = priceRegistryView
			lggr.Infow("generated price registry view", "priceRegistry", c.PriceRegistry.Address().String(), "chain", chain)
			return nil
		})
	}

	if c.RMN != nil {
		grp.Go(func() error {
			rmnView, err := v1_5.GenerateRMNView(c.RMN)
			if err != nil {
				return errors.Wrapf(err, "failed to generate rmn view for rmn %s", c.RMN.Address().String())
			}
			chainView.UpdateMu.Lock()
			defer chainView.UpdateMu.Unlock()
			chainView.RMN[c.RMN.Address().Hex()] = rmnView
			lggr.Infow("generated rmn view", "rmn", c.RMN.Address().Hex(), "chain", chain)
			return nil
		})
	}

	return chainView, grp.Wait()
}

func (c CCIPChainState) usdFeedOrDefault(symbol shared.TokenSymbol) common.Address {
	if feed, ok := c.USDFeeds[symbol]; ok {
		return feed.Address()
	}
	return common.Address{}
}

func validateLatestConfigOffRamp(offRamp offramp.OffRampInterface, cfg offramp.MultiOCR3BaseOCRConfig, homeChainConfig ccip_home.CCIPHomeChainConfig) error {
	// check if number of signers are unique and greater than 3
	if cfg.ConfigInfo.IsSignatureVerificationEnabled {
		if len(cfg.Signers) < 3 {
			return fmt.Errorf("offRamp %s config signers count mismatch: expected at least 3, got %d",
				offRamp.Address().Hex(), len(cfg.Signers))
		}
		if !deployment.IsAddressListUnique(cfg.Signers) {
			return fmt.Errorf("offRamp %s config signers list %v is not unique", offRamp.Address().Hex(), cfg.Signers)
		}
		if deployment.AddressListContainsEmptyAddress(cfg.Signers) {
			return fmt.Errorf("offRamp %s config signers list %v contains empty address", offRamp.Address().Hex(), cfg.Signers)
		}
	} else if len(cfg.Signers) != 0 {
		return fmt.Errorf("offRamp %s config signers count mismatch: expected 0, got %d",
			offRamp.Address().Hex(), len(cfg.Signers))
	}
	if len(cfg.Transmitters) < 3 {
		return fmt.Errorf("offRamp %s config transmitters count mismatch: expected at least 3, got %d",
			offRamp.Address().Hex(), len(cfg.Transmitters))
	}
	if !deployment.IsAddressListUnique(cfg.Transmitters) {
		return fmt.Errorf("offRamp %s config transmitters list %v is not unique", offRamp.Address().Hex(), cfg.Transmitters)
	}
	if deployment.AddressListContainsEmptyAddress(cfg.Transmitters) {
		return fmt.Errorf("offRamp %s config transmitters list %v contains empty address", offRamp.Address().Hex(), cfg.Transmitters)
	}

	// FRoleDON >= fChain is a requirement
	if cfg.ConfigInfo.F < homeChainConfig.FChain {
		return fmt.Errorf("offRamp %s config fChain mismatch: expected at least %d, got %d",
			offRamp.Address().Hex(), homeChainConfig.FChain, cfg.ConfigInfo.F)
	}

	//  transmitters.length should be validated such that it meets the 3 * fChain + 1 requirement
	minTransmitterReq := 3*int(homeChainConfig.FChain) + 1
	if len(cfg.Transmitters) < minTransmitterReq {
		return fmt.Errorf("offRamp %s config transmitters count mismatch: expected at least %d, got %d",
			offRamp.Address().Hex(), minTransmitterReq, len(cfg.Transmitters))
	}
	return nil
}
