package changeset

import (
	"context"
	std_errors "errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	"github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	cciptypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/log_message_data_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool_factory"
	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc677"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	viewv1_0 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	viewv1_5 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	viewv1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
	viewv1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_6"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	registryModuleOwnerCustomv15 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	registryModuleOwnerCustomv16 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/multicall3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"
)

var (
	// Legacy
	CommitStore   deployment.ContractType = "CommitStore"
	PriceRegistry deployment.ContractType = "PriceRegistry"
	RMN           deployment.ContractType = "RMN"

	// Not legacy
	MockRMN              deployment.ContractType = "MockRMN"
	RMNRemote            deployment.ContractType = "RMNRemote"
	ARMProxy             deployment.ContractType = "ARMProxy"
	WETH9                deployment.ContractType = "WETH9"
	Router               deployment.ContractType = "Router"
	TokenAdminRegistry   deployment.ContractType = "TokenAdminRegistry"
	TokenPoolFactory     deployment.ContractType = "TokenPoolFactory"
	RegistryModule       deployment.ContractType = "RegistryModuleOwnerCustom"
	NonceManager         deployment.ContractType = "NonceManager"
	FeeQuoter            deployment.ContractType = "FeeQuoter"
	CCIPHome             deployment.ContractType = "CCIPHome"
	RMNHome              deployment.ContractType = "RMNHome"
	OnRamp               deployment.ContractType = "OnRamp"
	OffRamp              deployment.ContractType = "OffRamp"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	PriceFeed            deployment.ContractType = "PriceFeed"

	// Test contracts. Note test router maps to a regular router contract.
	TestRouter             deployment.ContractType = "TestRouter"
	Multicall3             deployment.ContractType = "Multicall3"
	CCIPReceiver           deployment.ContractType = "CCIPReceiver"
	LogMessageDataReceiver deployment.ContractType = "LogMessageDataReceiver"
	USDCMockTransmitter    deployment.ContractType = "USDCMockTransmitter"

	// Pools
	BurnMintToken                  deployment.ContractType = "BurnMintToken"
	FactoryBurnMintERC20Token      deployment.ContractType = "FactoryBurnMintERC20Token"
	ERC20Token                     deployment.ContractType = "ERC20Token"
	ERC677Token                    deployment.ContractType = "ERC677Token"
	BurnMintTokenPool              deployment.ContractType = "BurnMintTokenPool"
	BurnWithFromMintTokenPool      deployment.ContractType = "BurnWithFromMintTokenPool"
	BurnFromMintTokenPool          deployment.ContractType = "BurnFromMintTokenPool"
	LockReleaseTokenPool           deployment.ContractType = "LockReleaseTokenPool"
	USDCToken                      deployment.ContractType = "USDCToken"
	USDCTokenMessenger             deployment.ContractType = "USDCTokenMessenger"
	USDCTokenPool                  deployment.ContractType = "USDCTokenPool"
	HybridLockReleaseUSDCTokenPool deployment.ContractType = "HybridLockReleaseUSDCTokenPool"

	// Firedrill
	FiredrillEntrypointType deployment.ContractType = "FiredrillEntrypoint"

	// Treasury
	FeeAggregator deployment.ContractType = "FeeAggregator"
)

// CCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState struct {
	commonstate.MCMSWithTimelockState
	commonstate.LinkTokenState
	commonstate.StaticLinkTokenState
	ABIByAddress       map[string]string
	OnRamp             onramp.OnRampInterface
	OffRamp            offramp.OffRampInterface
	FeeQuoter          *fee_quoter.FeeQuoter
	RMNProxy           *rmn_proxy_contract.RMNProxy
	NonceManager       *nonce_manager.NonceManager
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	TokenPoolFactory   *token_pool_factory.TokenPoolFactory
	RegistryModules1_6 []*registryModuleOwnerCustomv16.RegistryModuleOwnerCustom
	// TODO change this to contract object for v1.5 RegistryModules once we have the wrapper available in chainlink-evm
	RegistryModules1_5 []*registryModuleOwnerCustomv15.RegistryModuleOwnerCustom
	Router             *router.Router
	Weth9              *weth9.WETH9
	RMNRemote          *rmn_remote.RMNRemote
	// Map between token Descriptor (e.g. LinkSymbol, WethSymbol)
	// and the respective token / token pool contract(s) (only one of which would be active on the registry).
	// This is more of an illustration of how we'll have tokens, and it might need some work later to work properly.
	ERC20Tokens                map[TokenSymbol]*erc20.ERC20
	FactoryBurnMintERC20Token  *factory_burn_mint_erc20.FactoryBurnMintERC20
	ERC677Tokens               map[TokenSymbol]*erc677.ERC677
	BurnMintTokens677          map[TokenSymbol]*burn_mint_erc677.BurnMintERC677
	BurnMintTokenPools         map[TokenSymbol]map[semver.Version]*burn_mint_token_pool.BurnMintTokenPool
	BurnWithFromMintTokenPools map[TokenSymbol]map[semver.Version]*burn_with_from_mint_token_pool.BurnWithFromMintTokenPool
	BurnFromMintTokenPools     map[TokenSymbol]map[semver.Version]*burn_from_mint_token_pool.BurnFromMintTokenPool
	USDCTokenPools             map[semver.Version]*usdc_token_pool.USDCTokenPool
	LockReleaseTokenPools      map[TokenSymbol]map[semver.Version]*lock_release_token_pool.LockReleaseTokenPool
	// Map between token Symbol (e.g. LinkSymbol, WethSymbol)
	// and the respective aggregator USD feed contract
	USDFeeds map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface

	// Note we only expect one of these (on the home chain)
	CapabilityRegistry *capabilities_registry.CapabilitiesRegistry
	CCIPHome           *ccip_home.CCIPHome
	RMNHome            *rmn_home.RMNHome

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
	PriceRegistry  *price_registry_1_2_0.PriceRegistry
	RMN            *rmn_contract.RMNContract

	// Treasury contracts
	FeeAggregator common.Address
}

// validateHomeChain validates the home chain contracts and their configurations after complete set up
// It cross-references the config across CCIPHome and OffRamps to ensure they are in sync
// This should be called after the complete deployment is done
func (c CCIPChainState) validateHomeChain(e deployment.Environment, nodes deployment.Nodes, offRampsByChain map[uint64]offramp.OffRampInterface) error {
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
	ccipDons, err := ccip.GetCCIPDonsFromCapRegistry(e.GetContext(), c.CapabilityRegistry)
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
		}, don.Id, uint8(cciptypes.PluginTypeCCIPCommit))
		if err != nil {
			return fmt.Errorf("failed to get commit config for don %d: %w", don.Id, err)
		}
		if err := c.validateCCIPHomeVersionedActiveConfig(e, nodes, commitConfig.ActiveConfig, offRampsByChain); err != nil {
			return fmt.Errorf("failed to validate active commit config for don %d: %w", don.Id, err)
		}
		execConfig, err := c.CCIPHome.GetAllConfigs(&bind.CallOpts{
			Context: e.GetContext(),
		}, don.Id, uint8(cciptypes.PluginTypeCCIPExec))
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
func (c CCIPChainState) validateCCIPHomeVersionedActiveConfig(e deployment.Environment, nodes deployment.Nodes, homeCfg ccip_home.CCIPHomeVersionedConfig, offRampsByChain map[uint64]offramp.OffRampInterface) error {
	if homeCfg.ConfigDigest == [32]byte{} {
		return errors.New("active config digest is empty")
	}
	chainSel := homeCfg.Config.ChainSelector
	if _, exists := e.SolChains[chainSel]; exists {
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
	case uint8(cciptypes.PluginTypeCCIPCommit):
		commitConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: e.GetContext(),
		}, uint8(cciptypes.PluginTypeCCIPCommit))
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
	case uint8(cciptypes.PluginTypeCCIPExec):
		execConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: e.GetContext(),
		}, uint8(cciptypes.PluginTypeCCIPExec))
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

// validateOnRamp validates whether the contract addresses configured in static and dynamic config are in sync with state
func (c CCIPChainState) validateOnRamp(
	e deployment.Environment,
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
		if dynamicCfg.FeeAggregator != e.Chains[selector].DeployerKey.From {
			return fmt.Errorf("onRamp %s feeAggregator mismatch in dynamic config: expected deployer key %s, got %s",
				c.OnRamp.Address().Hex(), e.Chains[selector].DeployerKey.From.Hex(), dynamicCfg.FeeAggregator.Hex())
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

// validateFeeQuoter validates whether the fee quoter contract address configured in static config is in sync with state
func (c CCIPChainState) validateFeeQuoter(e deployment.Environment) error {
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

// validateRouter validates the router contract to check if all wired contracts are synced with state
// and returns all connected chains with respect to the router
func (c CCIPChainState) validateRouter(e deployment.Environment, isTestRouter bool) ([]uint64, error) {
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
		if _, exists := e.SolChains[d.SourceChainSelector]; exists {
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

// validateRMNRemote validates the RMNRemote contract to check if all wired contracts are synced with state
// and returns whether RMN is enabled for the chain on the RMNRemote
// It validates whether RMNRemote is in sync with the RMNHome contract
func (c CCIPChainState) validateRMNRemote(
	e deployment.Environment,
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

// validateOffRamp validates the offRamp contract to check if all wired contracts are synced with state
func (c CCIPChainState) validateOffRamp(
	e deployment.Environment,
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

func (c CCIPChainState) TokenAddressBySymbol() (map[TokenSymbol]common.Address, error) {
	tokenAddresses := make(map[TokenSymbol]common.Address)
	if c.FactoryBurnMintERC20Token != nil {
		tokenAddresses[FactoryBurnMintERC20Symbol] = c.FactoryBurnMintERC20Token.Address()
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
	tokenAddresses[LinkSymbol], err = c.LinkTokenAddress()
	if err != nil {
		return nil, err
	}
	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenAddresses[WethSymbol] = c.Weth9.Address()
	return tokenAddresses, nil
}

// TokenDetailsBySymbol get token mapping from the state. It contains only tokens that we have in address book
func (c CCIPChainState) TokenDetailsBySymbol() (map[TokenSymbol]TokenDetails, error) {
	tokenDetails := make(map[TokenSymbol]TokenDetails)
	if c.FactoryBurnMintERC20Token != nil {
		tokenDetails[FactoryBurnMintERC20Symbol] = c.FactoryBurnMintERC20Token
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
		tokenDetails[LinkSymbol] = c.LinkToken
	}
	if c.StaticLinkToken != nil {
		tokenDetails[LinkSymbol] = c.StaticLinkToken
	}

	if _, ok := tokenDetails[LinkSymbol]; !ok {
		return nil, errors.New("no LINK contract found in the state")
	}

	if c.Weth9 == nil {
		return nil, errors.New("no WETH contract found in the state")
	}
	tokenDetails[WethSymbol] = c.Weth9
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
			taView, err := viewv1_5.GenerateTokenAdminRegistryView(c.TokenAdminRegistry)
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
			tpfView, err := viewv1_5_1.GenerateTokenPoolFactoryView(c.TokenPoolFactory)
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
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
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
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
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
				tokenPoolView, err := viewv1_5_1.GenerateTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
				if err != nil {
					return errors.Wrapf(err, "failed to generate burn mint token pool view for %s", tokenPool.Address().String())
				}
				chainView.UpdateTokenPool(tokenSymbol.String(), tokenPool.Address().Hex(), viewv1_5_1.PoolView{
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
				tokenPoolView, err := viewv1_5_1.GenerateLockReleaseTokenPoolView(tokenPool, c.usdFeedOrDefault(tokenSymbol))
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
			tokenPoolView, err := viewv1_5_1.GenerateUSDCTokenPoolView(pool)
			if err != nil {
				return errors.Wrapf(err, "failed to generate USDC token pool view for %s", pool.Address().String())
			}
			chainView.UpdateTokenPool(string(USDCSymbol), pool.Address().Hex(), tokenPoolView)
			lggr.Infow("generated USDC token pool view", "tokenPool", pool.Address().Hex(), "chain", chain)
			return nil
		})
	}
	if c.NonceManager != nil {
		grp.Go(func() error {
			nmView, err := viewv1_6.GenerateNonceManagerView(c.NonceManager)
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
			rmnView, err := viewv1_6.GenerateRMNRemoteView(c.RMNRemote)
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
			rmnHomeView, err := viewv1_6.GenerateRMNHomeView(c.RMNHome)
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
			fqView, err := viewv1_6.GenerateFeeQuoterView(c.FeeQuoter, c.Router, c.TestRouter, tokens)
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
			onRampView, err := viewv1_6.GenerateOnRampView(
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
			offRampView, err := viewv1_6.GenerateOffRampView(
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
			rmnProxyView, err := viewv1_0.GenerateRMNProxyView(c.RMNProxy)
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
			chView, err := viewv1_6.GenerateCCIPHomeView(c.CapabilityRegistry, c.CCIPHome)
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
			capRegView, err := common_v1_0.GenerateCapabilityRegistryView(c.CapabilityRegistry)
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
			registryModuleView, err := shared.GetRegistryModuleView(registryModule, c.TokenAdminRegistry.Address())
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
			registryModuleView, err := shared.GetRegistryModuleView(registryModule, c.TokenAdminRegistry.Address())
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
			rmnView, err := viewv1_5.GenerateRMNView(c.RMN)
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

func (c CCIPChainState) usdFeedOrDefault(symbol TokenSymbol) common.Address {
	if feed, ok := c.USDFeeds[symbol]; ok {
		return feed.Address()
	}
	return common.Address{}
}

// CCIPOnChainState state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	Chains      map[uint64]CCIPChainState
	SolChains   map[uint64]SolCCIPChainState
	AptosChains map[uint64]AptosCCIPChainState
}

// ValidatePostDeploymentState should be called after the deployment and configuration for all contracts
// in environment is complete.
// It validates the state of the contracts and ensures that they are correctly configured and wired with each other.
func (c CCIPOnChainState) ValidatePostDeploymentState(e deployment.Environment) error {
	onRampsBySelector := make(map[uint64]common.Address)
	offRampsBySelector := make(map[uint64]offramp.OffRampInterface)
	for selector, chainState := range c.Chains {
		if chainState.OnRamp == nil {
			return fmt.Errorf("onramp not found in the state for chain %d", selector)
		}
		onRampsBySelector[selector] = chainState.OnRamp.Address()
		offRampsBySelector[selector] = chainState.OffRamp
	}
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return fmt.Errorf("failed to get node info from env: %w", err)
	}
	homeChain, err := c.HomeChainSelector()
	if err != nil {
		return fmt.Errorf("failed to get home chain selector: %w", err)
	}
	homeChainState := c.Chains[homeChain]
	if err := homeChainState.validateHomeChain(e, nodes, offRampsBySelector); err != nil {
		return fmt.Errorf("failed to validate home chain %d: %w", homeChain, err)
	}
	rmnHomeActiveDigest, err := homeChainState.RMNHome.GetActiveDigest(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return fmt.Errorf("failed to get active digest for RMNHome %s at home chain %d: %w", homeChainState.RMNHome.Address().Hex(), homeChain, err)
	}
	isRMNEnabledInRMNHomeBySourceChain := make(map[uint64]bool)
	rmnHomeConfig, err := homeChainState.RMNHome.GetConfig(&bind.CallOpts{
		Context: e.GetContext(),
	}, rmnHomeActiveDigest)
	if err != nil {
		return fmt.Errorf("failed to get config for RMNHome %s at home chain %d: %w", homeChainState.RMNHome.Address().Hex(), homeChain, err)
	}
	// if Fobserve is greater than 0, RMN is enabled for the source chain in RMNHome
	for _, rmnHomeChain := range rmnHomeConfig.VersionedConfig.DynamicConfig.SourceChains {
		isRMNEnabledInRMNHomeBySourceChain[rmnHomeChain.ChainSelector] = rmnHomeChain.FObserve > 0
	}
	for selector, chainState := range c.Chains {
		isRMNEnabledInRmnRemote, err := chainState.validateRMNRemote(e, selector, rmnHomeActiveDigest)
		if err != nil {
			return fmt.Errorf("failed to validate RMNRemote %s for chain %d: %w", chainState.RMNRemote.Address().Hex(), selector, err)
		}
		// check whether RMNRemote and RMNHome are in sync in terms of RMNEnabled
		if isRMNEnabledInRmnRemote != isRMNEnabledInRMNHomeBySourceChain[selector] {
			return fmt.Errorf("RMNRemote %s rmnEnabled mismatch with RMNHome for chain %d: expected %v, got %v",
				chainState.RMNRemote.Address().Hex(), selector, isRMNEnabledInRMNHomeBySourceChain[selector], isRMNEnabledInRmnRemote)
		}
		otherOnRamps := make(map[uint64]common.Address)
		isTestRouter := true
		if chainState.Router != nil {
			isTestRouter = false
		}
		connectedChains, err := chainState.validateRouter(e, isTestRouter)
		if err != nil {
			return fmt.Errorf("failed to validate router %s for chain %d: %w", chainState.Router.Address().Hex(), selector, err)
		}
		for _, connectedChain := range connectedChains {
			if connectedChain == selector {
				continue
			}
			otherOnRamps[connectedChain] = c.Chains[connectedChain].OnRamp.Address()
		}
		if err := chainState.validateOffRamp(e, selector, otherOnRamps, isRMNEnabledInRMNHomeBySourceChain); err != nil {
			return fmt.Errorf("failed to validate offramp %s for chain %d: %w", chainState.OffRamp.Address().Hex(), selector, err)
		}
		if err := chainState.validateOnRamp(e, selector, connectedChains); err != nil {
			return fmt.Errorf("failed to validate onramp %s for chain %d: %w", chainState.OnRamp.Address().Hex(), selector, err)
		}
		if err := chainState.validateFeeQuoter(e); err != nil {
			return fmt.Errorf("failed to validate fee quoter %s for chain %d: %w", chainState.FeeQuoter.Address().Hex(), selector, err)
		}
	}
	return nil
}

// HomeChainSelector returns the selector of the home chain based on the presence of RMNHome, CapabilityRegistry and CCIPHome contracts.
func (c CCIPOnChainState) HomeChainSelector() (uint64, error) {
	for selector, chain := range c.Chains {
		if chain.RMNHome != nil && chain.CapabilityRegistry != nil && chain.CCIPHome != nil {
			return selector, nil
		}
	}
	return 0, errors.New("no home chain found")
}

func (c CCIPOnChainState) EVMMCMSStateByChain() map[uint64]commonstate.MCMSWithTimelockState {
	mcmsStateByChain := make(map[uint64]commonstate.MCMSWithTimelockState)
	for chainSelector, chain := range c.Chains {
		mcmsStateByChain[chainSelector] = commonstate.MCMSWithTimelockState{
			CancellerMcm: chain.CancellerMcm,
			BypasserMcm:  chain.BypasserMcm,
			ProposerMcm:  chain.ProposerMcm,
			Timelock:     chain.Timelock,
			CallProxy:    chain.CallProxy,
		}
	}
	return mcmsStateByChain
}

func (c CCIPOnChainState) OffRampPermissionLessExecutionThresholdSeconds(ctx context.Context, env deployment.Environment, selector uint64) (uint32, error) {
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return 0, err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		chain, ok := c.Chains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d not found in the state", selector)
		}
		offRamp := chain.OffRamp
		if offRamp == nil {
			return 0, fmt.Errorf("offramp not found in the state for chain %d", selector)
		}
		dCfg, err := offRamp.GetDynamicConfig(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return dCfg.PermissionLessExecutionThresholdSeconds, fmt.Errorf("fetch dynamic config from offRamp %s for chain %d: %w", offRamp.Address().String(), selector, err)
		}
		return dCfg.PermissionLessExecutionThresholdSeconds, nil
	case chain_selectors.FamilySolana:
		chainState, ok := c.SolChains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d not found in the state", selector)
		}
		chain, ok := env.SolChains[selector]
		if !ok {
			return 0, fmt.Errorf("solana chain %d not found in the environment", selector)
		}
		if chainState.OffRamp.IsZero() {
			return 0, fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", selector)
		}
		var offRampConfig solOffRamp.Config
		offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(chainState.OffRamp)
		err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
		if err != nil {
			return 0, fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
		}
		// #nosec G115
		return uint32(offRampConfig.EnableManualExecutionAfter), nil
	case chain_selectors.FamilyAptos:
		chainState, ok := c.AptosChains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d does not exist in state", selector)
		}
		chain, ok := env.AptosChains[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d does not exist in env", selector)
		}
		if chainState.CCIPAddress == (aptos.AccountAddress{}) {
			return 0, fmt.Errorf("ccip not found in existing state, deploy the ccip first for Aptos chain %d", selector)
		}
		offrampDynamicConfig, err := getOfframpDynamicConfig(chain, chainState.CCIPAddress)
		if err != nil {
			return 0, fmt.Errorf("failed to get offramp dynamic config for Aptos chain %d: %w", selector, err)
		}
		return offrampDynamicConfig.PermissionlessExecutionThresholdSeconds, nil
	}
	return 0, fmt.Errorf("unsupported chain family %s", family)
}

func (c CCIPOnChainState) Validate() error {
	for sel, chain := range c.Chains {
		// cannot have static link and link together
		if chain.LinkToken != nil && chain.StaticLinkToken != nil {
			return fmt.Errorf("cannot have both link and static link token on the same chain %d", sel)
		}
	}
	return nil
}

func (c CCIPOnChainState) GetAllProposerMCMSForChains(chains []uint64) (map[uint64]*gethwrappers.ManyChainMultiSig, error) {
	multiSigs := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range chains {
		chainState, ok := c.Chains[chain]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chain)
		}
		if chainState.ProposerMcm == nil {
			return nil, fmt.Errorf("proposer mcm not found for chain %d", chain)
		}
		multiSigs[chain] = chainState.ProposerMcm
	}
	return multiSigs, nil
}

func (c CCIPOnChainState) GetAllTimeLocksForChains(chains []uint64) (map[uint64]common.Address, error) {
	timelocks := make(map[uint64]common.Address)
	for _, chain := range chains {
		chainState, ok := c.Chains[chain]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chain)
		}
		if chainState.Timelock == nil {
			return nil, fmt.Errorf("timelock not found for chain %d", chain)
		}
		timelocks[chain] = chainState.Timelock.Address()
	}
	return timelocks, nil
}

func (c CCIPOnChainState) SupportedChains() map[uint64]struct{} {
	chains := make(map[uint64]struct{})
	for chain := range c.Chains {
		chains[chain] = struct{}{}
	}
	for chain := range c.SolChains {
		chains[chain] = struct{}{}
	}
	for chain := range c.AptosChains {
		chains[chain] = struct{}{}
	}
	return chains
}

// EnforceMCMSUsageIfProd determines if an MCMS config should be enforced for this particular environment.
// It checks if the CCIPHome and CapabilitiesRegistry contracts are owned by the Timelock because all other contracts should follow this precedent.
// If the home chain contracts are owned by the Timelock and no mcmsConfig is provided, this function will return an error.
func (c CCIPOnChainState) EnforceMCMSUsageIfProd(ctx context.Context, mcmsConfig *proposalutils.TimelockConfig) error {
	// Instead of accepting a homeChainSelector, we simply look for the CCIPHome and CapabilitiesRegistry in state.
	// This is because the home chain selector is not always available in the input to a changeset.
	// Also, if the underlying rules to EnforceMCMSUsageIfProd change (i.e. what determines "prod" changes),
	// we can simply update the function body without worrying about the function signature.
	var ccipHome *ccip_home.CCIPHome
	var capReg *capabilities_registry.CapabilitiesRegistry
	var homeChainSelector uint64
	for selector, chain := range c.Chains {
		if chain.CCIPHome == nil || chain.CapabilityRegistry == nil {
			continue
		}
		// This condition impacts the ability of this function to determine MCMS enforcement.
		// As such, we return an error if we find multiple chains with home chain contracts.
		if ccipHome != nil {
			return errors.New("multiple chains with CCIPHome and CapabilitiesRegistry contracts found")
		}
		ccipHome = chain.CCIPHome
		capReg = chain.CapabilityRegistry
		homeChainSelector = selector
	}
	// It is not the job of this function to enforce the existence of home chain contracts.
	// Some tests don't deploy these contracts, and we don't want to fail them.
	// We simply say that MCMS is not enforced in such environments.
	if ccipHome == nil {
		return nil
	}
	// If the timelock contract is not found on the home chain,
	// we know that MCMS is not enforced.
	timelock := c.Chains[homeChainSelector].Timelock
	if timelock == nil {
		return nil
	}
	ccipHomeOwner, err := ccipHome.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get CCIP home owner: %w", err)
	}
	capRegOwner, err := capReg.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get capabilities registry owner: %w", err)
	}
	if ccipHomeOwner != capRegOwner {
		return fmt.Errorf("CCIPHome and CapabilitiesRegistry owners do not match: %s != %s", ccipHomeOwner.String(), capRegOwner.String())
	}
	// If CCIPHome & CapabilitiesRegistry are owned by timelock, then MCMS is enforced.
	if ccipHomeOwner == timelock.Address() && mcmsConfig == nil {
		return errors.New("MCMS is enforced for environment (i.e. CCIPHome & CapReg are owned by timelock), but no MCMS config was provided")
	}

	return nil
}

// ValidateOwnershipOfChain validates the ownership of every CCIP contract on a chain.
// If mcmsConfig is nil, the expected owner of each contract is the chain's deployer key.
// If provided, the expected owner is the Timelock contract.
func (c CCIPOnChainState) ValidateOwnershipOfChain(e deployment.Environment, chainSel uint64, mcmsConfig *proposalutils.TimelockConfig) error {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain with selector %d not found in the environment", chainSel)
	}

	chainState, ok := c.Chains[chainSel]
	if !ok {
		return fmt.Errorf("%s not found in the state", chain)
	}
	if chainState.Timelock == nil {
		return fmt.Errorf("timelock not found on %s", chain)
	}

	ownedContracts := map[string]commoncs.Ownable{
		"router":             chainState.Router,
		"feeQuoter":          chainState.FeeQuoter,
		"offRamp":            chainState.OffRamp,
		"onRamp":             chainState.OnRamp,
		"nonceManager":       chainState.NonceManager,
		"rmnRemote":          chainState.RMNRemote,
		"rmnProxy":           chainState.RMNProxy,
		"tokenAdminRegistry": chainState.TokenAdminRegistry,
	}
	var wg sync.WaitGroup
	errs := make(chan error, len(ownedContracts))
	for contractName, contract := range ownedContracts {
		wg.Add(1)
		go func(name string, c commoncs.Ownable) {
			defer wg.Done()
			if c == nil {
				errs <- fmt.Errorf("missing %s contract on %s", name, chain)
				return
			}
			err := commoncs.ValidateOwnership(e.GetContext(), mcmsConfig != nil, chain.DeployerKey.From, chainState.Timelock.Address(), contract)
			if err != nil {
				errs <- fmt.Errorf("failed to validate ownership of %s contract on %s: %w", name, chain, err)
			}
		}(contractName, contract)
	}
	wg.Wait()
	close(errs)
	var multiErr error
	for err := range errs {
		multiErr = std_errors.Join(multiErr, err)
	}
	if multiErr != nil {
		return multiErr
	}

	return nil
}

func (c CCIPOnChainState) View(e *deployment.Environment, chains []uint64) (map[string]view.ChainView, map[string]view.SolChainView, error) {
	m := sync.Map{}
	sm := sync.Map{}
	grp := errgroup.Group{}
	for _, chainSelector := range chains {
		var name string
		chainSelector := chainSelector
		grp.Go(func() error {
			family, err := chain_selectors.GetSelectorFamily(chainSelector)
			if err != nil {
				return err
			}
			chainInfo, err := deployment.ChainInfo(chainSelector)
			if err != nil {
				return err
			}
			name = chainInfo.ChainName
			if chainInfo.ChainName == "" {
				name = strconv.FormatUint(chainSelector, 10)
			}
			id, err := chain_selectors.GetChainIDFromSelector(chainSelector)
			if err != nil {
				return fmt.Errorf("failed to get chain id from selector %d: %w", chainSelector, err)
			}
			e.Logger.Infow("Generating view for", "chainSelector", chainSelector, "chainName", name, "chainID", id)
			switch family {
			case chain_selectors.FamilyEVM:
				if _, ok := c.Chains[chainSelector]; !ok {
					return fmt.Errorf("chain not supported %d", chainSelector)
				}
				chainState := c.Chains[chainSelector]
				chainView, err := chainState.GenerateView(e.Logger, name)
				if err != nil {
					return err
				}
				chainView.ChainSelector = chainSelector
				chainView.ChainID = id
				m.Store(name, chainView)
				e.Logger.Infow("Completed view for", "chainSelector", chainSelector, "chainName", name, "chainID", id)
			case chain_selectors.FamilySolana:
				if _, ok := c.SolChains[chainSelector]; !ok {
					return fmt.Errorf("chain not supported %d", chainSelector)
				}
				chainState := c.SolChains[chainSelector]
				chainView, err := chainState.GenerateView(e.SolChains[chainSelector])
				if err != nil {
					return err
				}
				chainView.ChainSelector = chainSelector
				chainView.ChainID = id
				sm.Store(name, chainView)
			default:
				return fmt.Errorf("unsupported chain family %s", family)
			}
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return nil, nil, err
	}
	finalEVMMap := make(map[string]view.ChainView)
	m.Range(func(key, value interface{}) bool {
		finalEVMMap[key.(string)] = value.(view.ChainView)
		return true
	})
	finalSolanaMap := make(map[string]view.SolChainView)
	sm.Range(func(key, value interface{}) bool {
		finalSolanaMap[key.(string)] = value.(view.SolChainView)
		return true
	})
	return finalEVMMap, finalSolanaMap, grp.Wait()
}

func (c CCIPOnChainState) GetOffRampAddressBytes(chainSelector uint64) ([]byte, error) {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, err
	}

	var offRampAddress []byte
	switch family {
	case chain_selectors.FamilyEVM:
		offRampAddress = c.Chains[chainSelector].OffRamp.Address().Bytes()
	case chain_selectors.FamilySolana:
		offRampAddress = c.SolChains[chainSelector].OffRamp.Bytes()
	case chain_selectors.FamilyAptos:
		ccipAddress := c.AptosChains[chainSelector].CCIPAddress
		offRampAddress = ccipAddress[:]
	default:
		return nil, fmt.Errorf("unsupported chain family %s", family)
	}

	return offRampAddress, nil
}

func (c CCIPOnChainState) GetOnRampAddressBytes(chainSelector uint64) ([]byte, error) {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, err
	}

	var onRampAddressBytes []byte
	switch family {
	case chain_selectors.FamilyEVM:
		if c.Chains[chainSelector].OnRamp == nil {
			return nil, fmt.Errorf("no onramp found in the state for chain %d", chainSelector)
		}
		onRampAddressBytes = c.Chains[chainSelector].OnRamp.Address().Bytes()
	case chain_selectors.FamilySolana:
		if c.SolChains[chainSelector].Router.IsZero() {
			return nil, fmt.Errorf("no router found in the state for chain %d", chainSelector)
		}
		onRampAddressBytes = c.SolChains[chainSelector].Router.Bytes()
	case chain_selectors.FamilyAptos:
		ccipAddress := c.AptosChains[chainSelector].CCIPAddress
		if ccipAddress == (aptos.AccountAddress{}) {
			return nil, fmt.Errorf("no ccip address found in the state for Aptos chain %d", chainSelector)
		}
		onRampAddressBytes = ccipAddress[:]
	default:
		return nil, fmt.Errorf("unsupported chain family %s", family)
	}

	return onRampAddressBytes, nil
}

func (c CCIPOnChainState) ValidateRamp(chainSelector uint64, rampType deployment.ContractType) error {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		chainState, exists := c.Chains[chainSelector]
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		switch rampType {
		case OffRamp:
			if chainState.OffRamp == nil {
				return fmt.Errorf("offramp contract does not exist on evm chain %d", chainSelector)
			}
		case OnRamp:
			if chainState.OnRamp == nil {
				return fmt.Errorf("onramp contract does not exist on evm chain %d", chainSelector)
			}
		default:
			return fmt.Errorf("unknown ramp type %s", rampType)
		}

	case chain_selectors.FamilySolana:
		chainState, exists := c.SolChains[chainSelector]
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		switch rampType {
		case OffRamp:
			if chainState.OffRamp.IsZero() {
				return fmt.Errorf("offramp contract does not exist on solana chain %d", chainSelector)
			}
		case OnRamp:
			if chainState.Router.IsZero() {
				return fmt.Errorf("router contract does not exist on solana chain %d", chainSelector)
			}
		default:
			return fmt.Errorf("unknown ramp type %s", rampType)
		}

	case chain_selectors.FamilyAptos:
		chainState, exists := c.AptosChains[chainSelector]
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		if chainState.CCIPAddress == (aptos.AccountAddress{}) {
			return fmt.Errorf("ccip package does not exist on aptos chain %d", chainSelector)
		}

	default:
		return fmt.Errorf("unknown chain family %s", family)
	}
	return nil
}

func LoadOnchainState(e deployment.Environment) (CCIPOnChainState, error) {
	solanaState, err := LoadOnchainStateSolana(e)
	if err != nil {
		return CCIPOnChainState{}, err
	}
	aptosChains, err := LoadOnchainStateAptos(e)
	if err != nil {
		return CCIPOnChainState{}, err
	}
	state := CCIPOnChainState{
		Chains:      make(map[uint64]CCIPChainState),
		SolChains:   solanaState.SolChains,
		AptosChains: aptosChains,
	}
	for chainSelector, chain := range e.Chains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainState(e.GetContext(), chain, addresses)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = chainState
	}
	return state, state.Validate()
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(ctx context.Context, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (CCIPChainState, error) {
	var state CCIPChainState
	mcmsWithTimelock, err := commonstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.MCMSWithTimelockState = *mcmsWithTimelock

	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.LinkTokenState = *linkState
	staticLinkState, err := commonstate.MaybeLoadStaticLinkTokenState(chain, addresses)
	if err != nil {
		return state, err
	}
	state.StaticLinkTokenState = *staticLinkState
	state.ABIByAddress = make(map[string]string)
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.RBACTimelockABI
		case deployment.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.CallProxyABI
		case deployment.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String(),
			deployment.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
		case deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token.LinkTokenABI
		case deployment.NewTypeAndVersion(commontypes.StaticLinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token_interface.LinkTokenABI
		case deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0).String():
			cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CapabilityRegistry = cr
			state.ABIByAddress[address] = capabilities_registry.CapabilitiesRegistryABI
		case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0).String():
			onRampC, err := onramp.NewOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OnRamp = onRampC
			state.ABIByAddress[address] = onramp.OnRampABI
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0).String():
			offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OffRamp = offRamp
			state.ABIByAddress[address] = offramp.OffRampABI
		case deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxy = armProxy
			state.ABIByAddress[address] = rmn_proxy_contract.RMNProxyABI
		case deployment.NewTypeAndVersion(RMNRemote, deployment.Version1_6_0).String():
			rmnRemote, err := rmn_remote.NewRMNRemote(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNRemote = rmnRemote
			state.ABIByAddress[address] = rmn_remote.RMNRemoteABI
		case deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0).String():
			rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNHome = rmnHome
			state.ABIByAddress[address] = rmn_home.RMNHomeABI
		case deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0).String():
			_weth9, err := weth9.NewWETH9(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Weth9 = _weth9
			state.ABIByAddress[address] = weth9.WETH9ABI
		case deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0).String():
			nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.NonceManager = nm
			state.ABIByAddress[address] = nonce_manager.NonceManagerABI
		case deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0).String():
			tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenAdminRegistry = tm
			state.ABIByAddress[address] = token_admin_registry.TokenAdminRegistryABI
		case deployment.NewTypeAndVersion(TokenPoolFactory, deployment.Version1_5_1).String():
			tpf, err := token_pool_factory.NewTokenPoolFactory(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenPoolFactory = tpf
			state.ABIByAddress[address] = token_pool_factory.TokenPoolFactoryABI
		case deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_6_0).String():
			rm, err := registryModuleOwnerCustomv16.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModules1_6 = append(state.RegistryModules1_6, rm)
			state.ABIByAddress[address] = registryModuleOwnerCustomv16.RegistryModuleOwnerCustomABI
		case deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_5_0).String():
			rm, err := registryModuleOwnerCustomv15.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModules1_5 = append(state.RegistryModules1_5, rm)
			state.ABIByAddress[address] = registryModuleOwnerCustomv15.RegistryModuleOwnerCustomABI
		case deployment.NewTypeAndVersion(Router, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Router = r
			state.ABIByAddress[address] = router.RouterABI
		case deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TestRouter = r
			state.ABIByAddress[address] = router.RouterABI
		case deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0).String():
			fq, err := fee_quoter.NewFeeQuoter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FeeQuoter = fq
			state.ABIByAddress[address] = fee_quoter.FeeQuoterABI
		case deployment.NewTypeAndVersion(USDCToken, deployment.Version1_0_0).String():
			ut, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.BurnMintTokens677 = map[TokenSymbol]*burn_mint_erc677.BurnMintERC677{
				USDCSymbol: ut,
			}
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case deployment.NewTypeAndVersion(USDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
		case deployment.NewTypeAndVersion(HybridLockReleaseUSDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
			state.ABIByAddress[address] = usdc_token_pool.USDCTokenPoolABI
		case deployment.NewTypeAndVersion(USDCMockTransmitter, deployment.Version1_0_0).String():
			umt, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTransmitter = umt
			state.ABIByAddress[address] = mock_usdc_token_transmitter.MockE2EUSDCTransmitterABI
		case deployment.NewTypeAndVersion(USDCTokenMessenger, deployment.Version1_0_0).String():
			utm, err := mock_usdc_token_messenger.NewMockE2EUSDCTokenMessenger(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTokenMessenger = utm
			state.ABIByAddress[address] = mock_usdc_token_messenger.MockE2EUSDCTokenMessengerABI
		case deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0).String():
			ccipHome, err := ccip_home.NewCCIPHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPHome = ccipHome
			state.ABIByAddress[address] = ccip_home.CCIPHomeABI
		case deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0).String():
			mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Receiver = mr
			state.ABIByAddress[address] = maybe_revert_message_receiver.MaybeRevertMessageReceiverABI
		case deployment.NewTypeAndVersion(LogMessageDataReceiver, deployment.Version1_0_0).String():
			mr, err := log_message_data_receiver.NewLogMessageDataReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.LogMessageDataReceiver = mr
			state.ABIByAddress[address] = log_message_data_receiver.LogMessageDataReceiverABI
		case deployment.NewTypeAndVersion(Multicall3, deployment.Version1_0_0).String():
			mc, err := multicall3.NewMulticall3(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Multicall3 = mc
			state.ABIByAddress[address] = multicall3.Multicall3ABI
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
			key, ok := DescriptionToTokenSymbol[desc]
			if !ok {
				return state, fmt.Errorf("unknown feed description %s", desc)
			}
			state.USDFeeds[key] = feed
			state.ABIByAddress[address] = aggregator_v3_interface.AggregatorV3InterfaceABI
		case deployment.NewTypeAndVersion(BurnMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_mint_token_pool.NewBurnMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnMintTokenPools = helpers.AddValueToNestedMap(state.BurnMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_mint_token_pool.BurnMintTokenPoolABI
		case deployment.NewTypeAndVersion(BurnWithFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_with_from_mint_token_pool.NewBurnWithFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnWithFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnWithFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolABI
		case deployment.NewTypeAndVersion(BurnFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, burn_from_mint_token_pool.NewBurnFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_from_mint_token_pool.BurnFromMintTokenPoolABI
		case deployment.NewTypeAndVersion(LockReleaseTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := NewTokenPoolWithMetadata(ctx, lock_release_token_pool.NewLockReleaseTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.LockReleaseTokenPools = helpers.AddValueToNestedMap(state.LockReleaseTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = lock_release_token_pool.LockReleaseTokenPoolABI
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
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case deployment.NewTypeAndVersion(ERC20Token, deployment.Version1_0_0).String():
			tok, err := erc20.NewERC20(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC20Tokens == nil {
				state.ERC20Tokens = make(map[TokenSymbol]*erc20.ERC20)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC20Tokens[TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc20.ERC20ABI
		case deployment.NewTypeAndVersion(FactoryBurnMintERC20Token, deployment.Version1_0_0).String():
			tok, err := factory_burn_mint_erc20.NewFactoryBurnMintERC20(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FactoryBurnMintERC20Token = tok
			state.ABIByAddress[address] = factory_burn_mint_erc20.FactoryBurnMintERC20ABI
		case deployment.NewTypeAndVersion(ERC677Token, deployment.Version1_0_0).String():
			tok, err := erc677.NewERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC677Tokens == nil {
				state.ERC677Tokens = make(map[TokenSymbol]*erc677.ERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC677Tokens[TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc677.ERC677ABI
		// legacy addresses below
		case deployment.NewTypeAndVersion(OnRamp, deployment.Version1_5_0).String():
			onRampC, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := onRampC.GetStaticConfig(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get static config chain %s: %w", chain.String(), err)
			}
			if state.EVM2EVMOnRamp == nil {
				state.EVM2EVMOnRamp = make(map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp)
			}
			state.EVM2EVMOnRamp[sCfg.DestChainSelector] = onRampC
			state.ABIByAddress[address] = evm_2_evm_onramp.EVM2EVMOnRampABI
		case deployment.NewTypeAndVersion(OffRamp, deployment.Version1_5_0).String():
			offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := offRamp.GetStaticConfig(nil)
			if err != nil {
				return state, err
			}
			if state.EVM2EVMOffRamp == nil {
				state.EVM2EVMOffRamp = make(map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp)
			}
			state.EVM2EVMOffRamp[sCfg.SourceChainSelector] = offRamp
			state.ABIByAddress[address] = evm_2_evm_offramp.EVM2EVMOffRampABI
		case deployment.NewTypeAndVersion(CommitStore, deployment.Version1_5_0).String():
			commitStore, err := commit_store.NewCommitStore(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			sCfg, err := commitStore.GetStaticConfig(nil)
			if err != nil {
				return state, err
			}
			if state.CommitStore == nil {
				state.CommitStore = make(map[uint64]*commit_store.CommitStore)
			}
			state.CommitStore[sCfg.SourceChainSelector] = commitStore
			state.ABIByAddress[address] = commit_store.CommitStoreABI
		case deployment.NewTypeAndVersion(PriceRegistry, deployment.Version1_2_0).String():
			pr, err := price_registry_1_2_0.NewPriceRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.PriceRegistry = pr
			state.ABIByAddress[address] = price_registry_1_2_0.PriceRegistryABI
		case deployment.NewTypeAndVersion(RMN, deployment.Version1_5_0).String():
			rmnC, err := rmn_contract.NewRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMN = rmnC
			state.ABIByAddress[address] = rmn_contract.RMNContractABI
		case deployment.NewTypeAndVersion(MockRMN, deployment.Version1_0_0).String():
			mockRMN, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockRMN = mockRMN
			state.ABIByAddress[address] = mock_rmn_contract.MockRMNContractABI
		case deployment.NewTypeAndVersion(FeeAggregator, deployment.Version1_0_0).String():
			state.FeeAggregator = common.HexToAddress(address)
		case deployment.NewTypeAndVersion(FiredrillEntrypointType, deployment.Version1_5_0).String(),
			deployment.NewTypeAndVersion(FiredrillEntrypointType, deployment.Version1_6_0).String():
			// Ignore firedrill contracts
			// Firedrill contracts are unknown to core and their state is being loaded separately
		default:
			// ManyChainMultiSig 1.0.0 can have any of these labels, it can have either 1,2 or 3 of these -
			// bypasser, proposer and canceller
			// if you try to compare tvStr.String() you will have to compare all combinations of labels
			// so we will compare the type and version only
			if tvStr.Type == commontypes.ManyChainMultisig && tvStr.Version == deployment.Version1_0_0 {
				state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
				continue
			}
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}

func ValidateChain(env deployment.Environment, state CCIPOnChainState, chainSel uint64, mcmsCfg *proposalutils.TimelockConfig) error {
	err := deployment.IsValidChainSelector(chainSel)
	if err != nil {
		return fmt.Errorf("is not valid chain selector %d: %w", chainSel, err)
	}
	chain, ok := env.Chains[chainSel]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", chainSel)
	}
	chainState, ok := state.Chains[chainSel]
	if !ok {
		return fmt.Errorf("%s does not exist in state", chain)
	}
	if mcmsCfg != nil {
		err = mcmsCfg.Validate(chain, commonstate.MCMSWithTimelockState{
			CancellerMcm: chainState.CancellerMcm,
			ProposerMcm:  chainState.ProposerMcm,
			BypasserMcm:  chainState.BypasserMcm,
			Timelock:     chainState.Timelock,
			CallProxy:    chainState.CallProxy,
		})
		if err != nil {
			return err
		}
	}
	return nil
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
