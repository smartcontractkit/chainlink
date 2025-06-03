package stateview

import (
	"context"
	std_errors "errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	ccipshared "github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

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
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/burn_mint_erc677_helper"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/don_id_claimer"
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

// CCIPOnChainState state always derivable from an address book.
// Offchain state always derivable from a list of nodeIds.
// Note can translate this into Go struct needed for MCMS/Docs/UI.
type CCIPOnChainState struct {
	// Populated go bindings for the appropriate version for all contracts.
	// We would hold 2 versions of each contract here. Once we upgrade we can phase out the old one.
	// When generating bindings, make sure the package name corresponds to the version.
	Chains      map[uint64]evm.CCIPChainState
	SolChains   map[uint64]solana.CCIPChainState
	AptosChains map[uint64]aptosstate.CCIPChainState
	evmMu       *sync.RWMutex
}

func (c CCIPOnChainState) EVMChains() []uint64 {
	c.evmMu.RLock()
	defer c.evmMu.RUnlock()
	return maps.Keys(c.Chains)
}

func (c CCIPOnChainState) EVMChainState(selector uint64) (evm.CCIPChainState, bool) {
	c.evmMu.RLock()
	defer c.evmMu.RUnlock()
	chainState, ok := c.Chains[selector]
	return chainState, ok
}

func (c CCIPOnChainState) MustGetEVMChainState(selector uint64) evm.CCIPChainState {
	c.evmMu.RLock()
	defer c.evmMu.RUnlock()
	chainState, ok := c.Chains[selector]
	if !ok {
		panic("chain state not found for selector " + strconv.FormatUint(selector, 10))
	}
	return chainState
}

func (c CCIPOnChainState) WriteEVMChainState(selector uint64, chainState evm.CCIPChainState) {
	c.evmMu.Lock()
	defer c.evmMu.Unlock()
	c.Chains[selector] = chainState
}

// ValidatePostDeploymentState should be called after the deployment and configuration for all contracts
// in environment is complete.
// It validates the state of the contracts and ensures that they are correctly configured and wired with each other.
func (c CCIPOnChainState) ValidatePostDeploymentState(e cldf.Environment) error {
	onRampsBySelector := make(map[uint64]common.Address)
	offRampsBySelector := make(map[uint64]offramp.OffRampInterface)
	for _, selector := range c.EVMChains() {
		chainState := c.MustGetEVMChainState(selector)
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
	homeChainState := c.MustGetEVMChainState(homeChain)
	if err := homeChainState.ValidateHomeChain(e, nodes, offRampsBySelector); err != nil {
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
	for _, selector := range c.EVMChains() {
		chainState := c.MustGetEVMChainState(selector)
		isRMNEnabledInRmnRemote, err := chainState.ValidateRMNRemote(e, selector, rmnHomeActiveDigest)
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
		connectedChains, err := chainState.ValidateRouter(e, isTestRouter)
		if err != nil {
			return fmt.Errorf("failed to validate router %s for chain %d: %w", chainState.Router.Address().Hex(), selector, err)
		}
		for _, connectedChain := range connectedChains {
			if connectedChain == selector {
				continue
			}
			otherOnRamps[connectedChain] = c.MustGetEVMChainState(connectedChain).OnRamp.Address()
		}
		if err := chainState.ValidateOffRamp(e, selector, otherOnRamps, isRMNEnabledInRMNHomeBySourceChain); err != nil {
			return fmt.Errorf("failed to validate offramp %s for chain %d: %w", chainState.OffRamp.Address().Hex(), selector, err)
		}
		if err := chainState.ValidateOnRamp(e, selector, connectedChains); err != nil {
			return fmt.Errorf("failed to validate onramp %s for chain %d: %w", chainState.OnRamp.Address().Hex(), selector, err)
		}
		if err := chainState.ValidateFeeQuoter(e); err != nil {
			return fmt.Errorf("failed to validate fee quoter %s for chain %d: %w", chainState.FeeQuoter.Address().Hex(), selector, err)
		}
	}
	return nil
}

// HomeChainSelector returns the selector of the home chain based on the presence of RMNHome, CapabilityRegistry and CCIPHome contracts.
func (c CCIPOnChainState) HomeChainSelector() (uint64, error) {
	for _, selector := range c.EVMChains() {
		chain := c.MustGetEVMChainState(selector)
		if chain.RMNHome != nil && chain.CapabilityRegistry != nil && chain.CCIPHome != nil {
			return selector, nil
		}
	}
	return 0, errors.New("no home chain found")
}

func (c CCIPOnChainState) EVMMCMSStateByChain() map[uint64]commonstate.MCMSWithTimelockState {
	mcmsStateByChain := make(map[uint64]commonstate.MCMSWithTimelockState)
	for _, chainSelector := range c.EVMChains() {
		chain := c.MustGetEVMChainState(chainSelector)
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

func (c CCIPOnChainState) SolanaMCMSStateByChain(e cldf.Environment) map[uint64]commonstate.MCMSWithTimelockStateSolana {
	mcmsStateByChain := make(map[uint64]commonstate.MCMSWithTimelockStateSolana)
	for chainSelector := range e.BlockChains.SolanaChains() {
		addreses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return mcmsStateByChain
		}
		mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[chainSelector], addreses)
		if err != nil {
			return mcmsStateByChain
		}
		mcmsStateByChain[chainSelector] = *mcmState
	}
	return mcmsStateByChain
}

func (c CCIPOnChainState) OffRampPermissionLessExecutionThresholdSeconds(ctx context.Context, env cldf.Environment, selector uint64) (uint32, error) {
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return 0, err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		chain, ok := c.EVMChainState(selector)
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
		chain, ok := env.BlockChains.SolanaChains()[selector]
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
		chain, ok := env.BlockChains.AptosChains()[selector]
		if !ok {
			return 0, fmt.Errorf("chain %d does not exist in env", selector)
		}
		if chainState.CCIPAddress == (aptos.AccountAddress{}) {
			return 0, fmt.Errorf("ccip not found in existing state, deploy the ccip first for Aptos chain %d", selector)
		}
		offrampDynamicConfig, err := aptosstate.GetOfframpDynamicConfig(chain, chainState.CCIPAddress)
		if err != nil {
			return 0, fmt.Errorf("failed to get offramp dynamic config for Aptos chain %d: %w", selector, err)
		}
		return offrampDynamicConfig.PermissionlessExecutionThresholdSeconds, nil
	}
	return 0, fmt.Errorf("unsupported chain family %s", family)
}

func (c CCIPOnChainState) Validate() error {
	for _, sel := range c.EVMChains() {
		chain := c.MustGetEVMChainState(sel)
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
		chainState, ok := c.EVMChainState(chain)
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
		chainState, ok := c.EVMChainState(chain)
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
	for _, chain := range c.EVMChains() {
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
	for _, selector := range c.EVMChains() {
		chain := c.MustGetEVMChainState(selector)
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
	timelock := c.MustGetEVMChainState(homeChainSelector).Timelock
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
func (c CCIPOnChainState) ValidateOwnershipOfChain(e cldf.Environment, chainSel uint64, mcmsConfig *proposalutils.TimelockConfig) error {
	chain, ok := e.BlockChains.EVMChains()[chainSel]
	if !ok {
		return fmt.Errorf("chain with selector %d not found in the environment", chainSel)
	}

	chainState, ok := c.EVMChainState(chainSel)
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

func (c CCIPOnChainState) View(e *cldf.Environment, chains []uint64) (map[string]view.ChainView, map[string]view.SolChainView, error) {
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
			chainInfo, err := cldf_chain_utils.ChainInfo(chainSelector)
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
				if _, ok := c.EVMChainState(chainSelector); !ok {
					return fmt.Errorf("chain not supported %d", chainSelector)
				}
				chainState := c.MustGetEVMChainState(chainSelector)
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
				chainView, err := chainState.GenerateView(e, chainSelector)
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
		offRampAddress = c.MustGetEVMChainState(chainSelector).OffRamp.Address().Bytes()
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
		if c.MustGetEVMChainState(chainSelector).OnRamp == nil {
			return nil, fmt.Errorf("no onramp found in the state for chain %d", chainSelector)
		}
		onRampAddressBytes = c.MustGetEVMChainState(chainSelector).OnRamp.Address().Bytes()
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

func (c CCIPOnChainState) ValidateRamp(chainSelector uint64, rampType cldf.ContractType) error {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return err
	}
	switch family {
	case chain_selectors.FamilyEVM:
		chainState, exists := c.EVMChainState(chainSelector)
		if !exists {
			return fmt.Errorf("chain %d does not exist", chainSelector)
		}
		switch rampType {
		case ccipshared.OffRamp:
			if chainState.OffRamp == nil {
				return fmt.Errorf("offramp contract does not exist on evm chain %d", chainSelector)
			}
		case ccipshared.OnRamp:
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
		case ccipshared.OffRamp:
			if chainState.OffRamp.IsZero() {
				return fmt.Errorf("offramp contract does not exist on solana chain %d", chainSelector)
			}
		case ccipshared.OnRamp:
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

func LoadOnchainState(e cldf.Environment) (CCIPOnChainState, error) {
	solanaState, err := LoadOnchainStateSolana(e)
	if err != nil {
		return CCIPOnChainState{}, err
	}
	aptosChains, err := aptosstate.LoadOnchainStateAptos(e)
	if err != nil {
		return CCIPOnChainState{}, err
	}

	state := CCIPOnChainState{
		Chains:      make(map[uint64]evm.CCIPChainState),
		SolChains:   solanaState.SolChains,
		AptosChains: aptosChains,
		evmMu:       &sync.RWMutex{},
	}
	for chainSelector, chain := range e.BlockChains.EVMChains() {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return state, err
			}
			// Chain not found in address book, initialize empty
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := LoadChainState(e.GetContext(), chain, addresses)
		if err != nil {
			return state, err
		}
		state.WriteEVMChainState(chainSelector, chainState)
	}
	return state, state.Validate()
}

// LoadChainState Loads all state for a chain into state
func LoadChainState(ctx context.Context, chain cldf_evm.Chain, addresses map[string]cldf.TypeAndVersion) (evm.CCIPChainState, error) {
	var state evm.CCIPChainState
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
		case cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.RBACTimelockABI
		case cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.CallProxyABI
		case cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0).String(),
			cldf.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0).String(),
			cldf.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = gethwrappers.ManyChainMultiSigABI
		case cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token.LinkTokenABI
		case cldf.NewTypeAndVersion(commontypes.StaticLinkToken, deployment.Version1_0_0).String():
			state.ABIByAddress[address] = link_token_interface.LinkTokenABI
		case cldf.NewTypeAndVersion(ccipshared.CapabilitiesRegistry, deployment.Version1_0_0).String():
			cr, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CapabilityRegistry = cr
			state.ABIByAddress[address] = capabilities_registry.CapabilitiesRegistryABI
		case cldf.NewTypeAndVersion(ccipshared.OnRamp, deployment.Version1_6_0).String():
			onRampC, err := onramp.NewOnRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OnRamp = onRampC
			state.ABIByAddress[address] = onramp.OnRampABI
		case cldf.NewTypeAndVersion(ccipshared.OffRamp, deployment.Version1_6_0).String():
			offRamp, err := offramp.NewOffRamp(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.OffRamp = offRamp
			state.ABIByAddress[address] = offramp.OffRampABI
		case cldf.NewTypeAndVersion(ccipshared.ARMProxy, deployment.Version1_0_0).String():
			armProxy, err := rmn_proxy_contract.NewRMNProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNProxy = armProxy
			state.ABIByAddress[address] = rmn_proxy_contract.RMNProxyABI
		case cldf.NewTypeAndVersion(ccipshared.RMNRemote, deployment.Version1_6_0).String():
			rmnRemote, err := rmn_remote.NewRMNRemote(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNRemote = rmnRemote
			state.ABIByAddress[address] = rmn_remote.RMNRemoteABI
		case cldf.NewTypeAndVersion(ccipshared.RMNHome, deployment.Version1_6_0).String():
			rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMNHome = rmnHome
			state.ABIByAddress[address] = rmn_home.RMNHomeABI
		case cldf.NewTypeAndVersion(ccipshared.WETH9, deployment.Version1_0_0).String():
			_weth9, err := weth9.NewWETH9(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Weth9 = _weth9
			state.ABIByAddress[address] = weth9.WETH9ABI
		case cldf.NewTypeAndVersion(ccipshared.NonceManager, deployment.Version1_6_0).String():
			nm, err := nonce_manager.NewNonceManager(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.NonceManager = nm
			state.ABIByAddress[address] = nonce_manager.NonceManagerABI
		case cldf.NewTypeAndVersion(ccipshared.TokenAdminRegistry, deployment.Version1_5_0).String():
			tm, err := token_admin_registry.NewTokenAdminRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenAdminRegistry = tm
			state.ABIByAddress[address] = token_admin_registry.TokenAdminRegistryABI
		case cldf.NewTypeAndVersion(ccipshared.TokenPoolFactory, deployment.Version1_5_1).String():
			tpf, err := token_pool_factory.NewTokenPoolFactory(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TokenPoolFactory = tpf
			state.ABIByAddress[address] = token_pool_factory.TokenPoolFactoryABI
		case cldf.NewTypeAndVersion(ccipshared.RegistryModule, deployment.Version1_6_0).String():
			rm, err := registryModuleOwnerCustomv16.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModules1_6 = append(state.RegistryModules1_6, rm)
			state.ABIByAddress[address] = registryModuleOwnerCustomv16.RegistryModuleOwnerCustomABI
		case cldf.NewTypeAndVersion(ccipshared.RegistryModule, deployment.Version1_5_0).String():
			rm, err := registryModuleOwnerCustomv15.NewRegistryModuleOwnerCustom(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RegistryModules1_5 = append(state.RegistryModules1_5, rm)
			state.ABIByAddress[address] = registryModuleOwnerCustomv15.RegistryModuleOwnerCustomABI
		case cldf.NewTypeAndVersion(ccipshared.Router, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Router = r
			state.ABIByAddress[address] = router.RouterABI
		case cldf.NewTypeAndVersion(ccipshared.TestRouter, deployment.Version1_2_0).String():
			r, err := router.NewRouter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.TestRouter = r
			state.ABIByAddress[address] = router.RouterABI
		case cldf.NewTypeAndVersion(ccipshared.FeeQuoter, deployment.Version1_6_0).String():
			fq, err := fee_quoter.NewFeeQuoter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FeeQuoter = fq
			state.ABIByAddress[address] = fee_quoter.FeeQuoterABI
		case cldf.NewTypeAndVersion(ccipshared.USDCToken, deployment.Version1_0_0).String():
			ut, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.BurnMintTokens677 = map[ccipshared.TokenSymbol]*burn_mint_erc677.BurnMintERC677{
				ccipshared.USDCSymbol: ut,
			}
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case cldf.NewTypeAndVersion(ccipshared.USDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
		case cldf.NewTypeAndVersion(ccipshared.HybridLockReleaseUSDCTokenPool, deployment.Version1_5_1).String():
			utp, err := usdc_token_pool.NewUSDCTokenPool(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDCTokenPools == nil {
				state.USDCTokenPools = make(map[semver.Version]*usdc_token_pool.USDCTokenPool)
			}
			state.USDCTokenPools[deployment.Version1_5_1] = utp
			state.ABIByAddress[address] = usdc_token_pool.USDCTokenPoolABI
		case cldf.NewTypeAndVersion(ccipshared.USDCMockTransmitter, deployment.Version1_0_0).String():
			umt, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTransmitter = umt
			state.ABIByAddress[address] = mock_usdc_token_transmitter.MockE2EUSDCTransmitterABI
		case cldf.NewTypeAndVersion(ccipshared.USDCTokenMessenger, deployment.Version1_0_0).String():
			utm, err := mock_usdc_token_messenger.NewMockE2EUSDCTokenMessenger(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockUSDCTokenMessenger = utm
			state.ABIByAddress[address] = mock_usdc_token_messenger.MockE2EUSDCTokenMessengerABI
		case cldf.NewTypeAndVersion(ccipshared.CCIPHome, deployment.Version1_6_0).String():
			ccipHome, err := ccip_home.NewCCIPHome(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.CCIPHome = ccipHome
			state.ABIByAddress[address] = ccip_home.CCIPHomeABI
		case cldf.NewTypeAndVersion(ccipshared.CCIPReceiver, deployment.Version1_0_0).String():
			mr, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Receiver = mr
			state.ABIByAddress[address] = maybe_revert_message_receiver.MaybeRevertMessageReceiverABI
		case cldf.NewTypeAndVersion(ccipshared.LogMessageDataReceiver, deployment.Version1_0_0).String():
			mr, err := log_message_data_receiver.NewLogMessageDataReceiver(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.LogMessageDataReceiver = mr
			state.ABIByAddress[address] = log_message_data_receiver.LogMessageDataReceiverABI
		case cldf.NewTypeAndVersion(ccipshared.Multicall3, deployment.Version1_0_0).String():
			mc, err := multicall3.NewMulticall3(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.Multicall3 = mc
			state.ABIByAddress[address] = multicall3.Multicall3ABI
		case cldf.NewTypeAndVersion(ccipshared.PriceFeed, deployment.Version1_0_0).String():
			feed, err := aggregator_v3_interface.NewAggregatorV3Interface(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.USDFeeds == nil {
				state.USDFeeds = make(map[ccipshared.TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface)
			}
			desc, err := feed.Description(&bind.CallOpts{})
			if err != nil {
				return state, err
			}
			key, ok := ccipshared.GetSymbolFromDescription(desc)
			if !ok {
				return state, fmt.Errorf("unknown feed description %s", desc)
			}
			state.USDFeeds[key] = feed
			state.ABIByAddress[address] = aggregator_v3_interface.AggregatorV3InterfaceABI
		case cldf.NewTypeAndVersion(ccipshared.BurnMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := ccipshared.NewTokenPoolWithMetadata(ctx, burn_mint_token_pool.NewBurnMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnMintTokenPools = helpers.AddValueToNestedMap(state.BurnMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_mint_token_pool.BurnMintTokenPoolABI
		case cldf.NewTypeAndVersion(ccipshared.BurnWithFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := ccipshared.NewTokenPoolWithMetadata(ctx, burn_with_from_mint_token_pool.NewBurnWithFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnWithFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnWithFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolABI
		case cldf.NewTypeAndVersion(ccipshared.BurnFromMintTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := ccipshared.NewTokenPoolWithMetadata(ctx, burn_from_mint_token_pool.NewBurnFromMintTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.BurnFromMintTokenPools = helpers.AddValueToNestedMap(state.BurnFromMintTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = burn_from_mint_token_pool.BurnFromMintTokenPoolABI
		case cldf.NewTypeAndVersion(ccipshared.LockReleaseTokenPool, deployment.Version1_5_1).String():
			ethAddress := common.HexToAddress(address)
			pool, metadata, err := ccipshared.NewTokenPoolWithMetadata(ctx, lock_release_token_pool.NewLockReleaseTokenPool, ethAddress, chain.Client)
			if err != nil {
				return state, fmt.Errorf("failed to connect address %s with token pool bindings and get token symbol: %w", ethAddress, err)
			}
			state.LockReleaseTokenPools = helpers.AddValueToNestedMap(state.LockReleaseTokenPools, metadata.Symbol, metadata.Version, pool)
			state.ABIByAddress[address] = lock_release_token_pool.LockReleaseTokenPoolABI
		case cldf.NewTypeAndVersion(ccipshared.BurnMintToken, deployment.Version1_0_0).String():
			tok, err := burn_mint_erc677.NewBurnMintERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.BurnMintTokens677 == nil {
				state.BurnMintTokens677 = make(map[ccipshared.TokenSymbol]*burn_mint_erc677.BurnMintERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.BurnMintTokens677[ccipshared.TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = burn_mint_erc677.BurnMintERC677ABI
		case cldf.NewTypeAndVersion(ccipshared.ERC20Token, deployment.Version1_0_0).String():
			tok, err := erc20.NewERC20(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC20Tokens == nil {
				state.ERC20Tokens = make(map[ccipshared.TokenSymbol]*erc20.ERC20)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC20Tokens[ccipshared.TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc20.ERC20ABI
		case cldf.NewTypeAndVersion(ccipshared.FactoryBurnMintERC20Token, deployment.Version1_0_0).String():
			tok, err := factory_burn_mint_erc20.NewFactoryBurnMintERC20(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.FactoryBurnMintERC20Token = tok
			state.ABIByAddress[address] = factory_burn_mint_erc20.FactoryBurnMintERC20ABI
		case cldf.NewTypeAndVersion(ccipshared.ERC677Token, deployment.Version1_0_0).String():
			tok, err := erc677.NewERC677(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			if state.ERC677Tokens == nil {
				state.ERC677Tokens = make(map[ccipshared.TokenSymbol]*erc677.ERC677)
			}
			symbol, err := tok.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.ERC677Tokens[ccipshared.TokenSymbol(symbol)] = tok
			state.ABIByAddress[address] = erc677.ERC677ABI
		// legacy addresses below
		case cldf.NewTypeAndVersion(ccipshared.OnRamp, deployment.Version1_5_0).String():
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
		case cldf.NewTypeAndVersion(ccipshared.OffRamp, deployment.Version1_5_0).String():
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
		case cldf.NewTypeAndVersion(ccipshared.CommitStore, deployment.Version1_5_0).String():
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
		case cldf.NewTypeAndVersion(ccipshared.PriceRegistry, deployment.Version1_2_0).String():
			pr, err := price_registry_1_2_0.NewPriceRegistry(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.PriceRegistry = pr
			state.ABIByAddress[address] = price_registry_1_2_0.PriceRegistryABI
		case cldf.NewTypeAndVersion(ccipshared.RMN, deployment.Version1_5_0).String():
			rmnC, err := rmn_contract.NewRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.RMN = rmnC
			state.ABIByAddress[address] = rmn_contract.RMNContractABI
		case cldf.NewTypeAndVersion(ccipshared.MockRMN, deployment.Version1_0_0).String():
			mockRMN, err := mock_rmn_contract.NewMockRMNContract(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.MockRMN = mockRMN
			state.ABIByAddress[address] = mock_rmn_contract.MockRMNContractABI
		case cldf.NewTypeAndVersion(ccipshared.FeeAggregator, deployment.Version1_0_0).String():
			state.FeeAggregator = common.HexToAddress(address)
		case cldf.NewTypeAndVersion(ccipshared.FiredrillEntrypointType, deployment.Version1_5_0).String(),
			cldf.NewTypeAndVersion(ccipshared.FiredrillEntrypointType, deployment.Version1_6_0).String():
			// Ignore firedrill contracts
			// Firedrill contracts are unknown to core and their state is being loaded separately
		case cldf.NewTypeAndVersion(ccipshared.DonIDClaimer, deployment.Version1_6_1).String():
			donIDClaimer, err := don_id_claimer.NewDonIDClaimer(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}
			state.DonIDClaimer = donIDClaimer
			state.ABIByAddress[address] = don_id_claimer.DonIDClaimerABI
		case cldf.NewTypeAndVersion(ccipshared.ERC677TokenHelper, deployment.Version1_0_0).String():
			ERC677HelperToken, err := burn_mint_erc677_helper.NewBurnMintERC677Helper(common.HexToAddress(address), chain.Client)
			if err != nil {
				return state, err
			}

			if state.BurnMintTokens677Helper == nil {
				state.BurnMintTokens677Helper = make(map[ccipshared.TokenSymbol]*burn_mint_erc677_helper.BurnMintERC677Helper)
			}
			symbol, err := ERC677HelperToken.Symbol(nil)
			if err != nil {
				return state, fmt.Errorf("failed to get token symbol of token at %s: %w", address, err)
			}
			state.BurnMintTokens677Helper[ccipshared.TokenSymbol(symbol)] = ERC677HelperToken
			state.ABIByAddress[address] = burn_mint_erc677_helper.BurnMintERC677HelperABI
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

func ValidateChain(env cldf.Environment, state CCIPOnChainState, chainSel uint64, mcmsCfg *proposalutils.TimelockConfig) error {
	err := cldf.IsValidChainSelector(chainSel)
	if err != nil {
		return fmt.Errorf("is not valid chain selector %d: %w", chainSel, err)
	}
	chain, ok := env.BlockChains.EVMChains()[chainSel]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", chainSel)
	}
	chainState, ok := state.EVMChainState(chainSel)
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

func LoadOnchainStateSolana(e cldf.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		SolChains: make(map[uint64]solana.CCIPChainState),
	}
	for chainSelector, chain := range e.BlockChains.SolanaChains() {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !std_errors.Is(err, cldf.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := solana.LoadChainStateSolana(chain, addresses)
		if err != nil {
			return state, err
		}
		state.SolChains[chainSelector] = chainState
	}
	return state, nil
}
