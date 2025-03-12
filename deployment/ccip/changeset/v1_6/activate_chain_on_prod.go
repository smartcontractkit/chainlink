package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"
)

var _ deployment.ChangeSet[ActivateChainOnProdConfig] = ActivateChainOnProd

// ActivateChainOnProdConfig is a configuration struct for ActivateChainOnProd.
type ActivateChainOnProdConfig struct {
	NewChainSelector uint64
	MCMSConfig       *changeset.MCMSConfig
}

func (c ActivateChainOnProdConfig) Validate(e deployment.Environment, state changeset.CCIPOnChainState) error {
	if c.MCMSConfig == nil {
		return fmt.Errorf("mcms config is required")
	}

	if _, ok := state.Chains[c.NewChainSelector]; !ok {
		return fmt.Errorf("chain with selector %d not found", c.NewChainSelector)
	}

	onRamp := state.Chains[c.NewChainSelector].OnRamp
	offRamp := state.Chains[c.NewChainSelector].OffRamp
	router := state.Chains[c.NewChainSelector].Router

	if onRamp == nil {
		return fmt.Errorf("onRamp contract not found for chain with selector %d", c.NewChainSelector)
	}
	if offRamp == nil {
		return fmt.Errorf("offRamp contract not found for chain with selector %d", c.NewChainSelector)
	}
	if router == nil {
		return fmt.Errorf("router contract not found for chain with selector %d", c.NewChainSelector)
	}

	return nil
}

// ActivateChainOnProd activates a new chain on production routers in CCIP.
// InitChainForTesting should be executed first, along with some E2E transfers using test routers.
// This changeset assumes that the onRamp and offRamp on the new chain are not yet owned by MCMS.
func ActivateChainOnProd(e deployment.Environment, config ActivateChainOnProdConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = config.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate ActivateChainOnProdConfig: %w", err)
	}

	readOpts := &bind.CallOpts{Context: e.GetContext()}

	/*
		OnRamp             onramp.OnRampInterface
		OffRamp            offramp.OffRampInterface
		FeeQuoter          *fee_quoter.FeeQuoter
		RMNProxy           *rmn_proxy_contract.RMNProxy
		NonceManager       *nonce_manager.NonceManager
		TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
		RegistryModule     *registry_module_owner_custom.RegistryModuleOwnerCustom
		Router             *router.Router
		Weth9              *weth9.WETH9
		RMNRemote          *rmn_remote.RMNRemote
	*/

	// Transfer ownership of any contracts that are still owned by the deployer key.
	allContracts := []commoncs.Ownable{
		state.Chains[config.NewChainSelector].OnRamp,
		state.Chains[config.NewChainSelector].OffRamp,
		state.Chains[config.NewChainSelector].FeeQuoter,
		state.Chains[config.NewChainSelector].RMNProxy,
		state.Chains[config.NewChainSelector].NonceManager,
		state.Chains[config.NewChainSelector].TokenAdminRegistry,
		state.Chains[config.NewChainSelector].Router,
		state.Chains[config.NewChainSelector].RMNRemote,
		// TODO: Are there any other contracts that should be added here? I think this covers all 1.6.0 contracts on non-home chains.
	}
	addressesToTransfer := make([]common.Address, 0, len(allContracts))
	for _, contract := range allContracts {
		if contract == nil {
			continue
		}
		owner, err := contract.Owner(readOpts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get owner of contract %s: %w", contract.Address().Hex(), err)
		}
		if owner == e.Chains[config.NewChainSelector].DeployerKey.From {
			addressesToTransfer = append(addressesToTransfer, contract.Address())
		}
	}
	out, err := commoncs.TransferToMCMSWithTimelock(e, commoncs.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{
			config.NewChainSelector: addressesToTransfer,
		},
		MinDelay: config.MCMSConfig.MinDelay,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run TransferToMCMSWithTimelock on chain with selector %d: %w", config.NewChainSelector, err)
	}
	ownershipTransferProposals := out.Proposals

	// Renounce the admin role on the Timelock from the deployer key.
	out, err = commoncs.RenounceTimelockDeployer(e, commoncs.RenounceTimelockDeployerConfig{
		ChainSel: config.NewChainSelector,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run RenounceTimelockDeployer on chain with selector %d: %w", config.NewChainSelector, err)
	}

	// Fetch the sourceChainSelectors from the offRamp on the new chain.
	sourceChainSelectors, sourceChainConfigs, err := state.Chains[config.NewChainSelector].OffRamp.GetAllSourceChainConfigs(readOpts)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get source chain configs from OffRamp on chain with selector %d: %w", config.NewChainSelector, err)
	}
	// We don't want to modify the router used by disabled chains.
	// If we decide to re-enable the disabled chain, we want it to re-enable on whichever router it was previously using to avoid unintended effects.
	enabledSourceChainSelectors := make([]uint64, 0, len(sourceChainSelectors))
	for i, srcChainSelector := range sourceChainSelectors {
		if sourceChainConfigs[i].IsEnabled {
			enabledSourceChainSelectors = append(enabledSourceChainSelectors, srcChainSelector)
		}
	}

	// Enable the production router on [new chain -> each remote chain] and [each remote chain -> new chain].
	var allEnablementProposals []mcmslib.TimelockProposal
	err = enableProductionRouter(e, config.NewChainSelector, enabledSourceChainSelectors, config.MCMSConfig, allEnablementProposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", config.NewChainSelector, err)
	}
	for _, srcChainSelector := range enabledSourceChainSelectors {
		err = enableProductionRouter(e, srcChainSelector, []uint64{config.NewChainSelector}, config.MCMSConfig, allEnablementProposals)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", srcChainSelector, err)
		}
	}

	// Aggregate all proposals.
	var batches []types.BatchOperation
	// Add ownership transfer proposals to the aggregate.
	// We know that ownership transfer will only touch the new chain, which the enablement proposals will also touch.
	// Therefore, we don't need to populate timelocks, proposers, and inspectors maps here.
	for _, proposal := range ownershipTransferProposals {
		for _, batchTransaction := range proposal.Transactions {
			for _, transaction := range batchTransaction.Batch {
				batchOperation, err := proposalutils.BatchOperationForChain(
					uint64(batchTransaction.ChainIdentifier),
					transaction.To.Hex(),
					transaction.Data,
					big.NewInt(0),
					transaction.ContractType,
					transaction.Tags,
				)
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to create batch operation on chain with selector %d: %w", batchTransaction.ChainIdentifier, err)
				}
				batches = append(batches, batchOperation)
			}
		}
	}
	// Add proposals that follow the new format to the aggregate.
	// Also, populate the timelocks, proposers, and inspectors maps.
	timelocks := make(map[uint64]string)
	proposers := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	for _, proposal := range allEnablementProposals {
		batches = append(batches, proposal.Operations...)
		for chainSelector, address := range proposal.TimelockAddresses {
			chainSel := uint64(chainSelector)
			timelocks[chainSel] = address
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
		}
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		fmt.Sprintf("Enable traffic to and from chain %d on production routers", config.NewChainSelector),
		config.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

// enableProductionRouter updates the onRamp and offRamp to point at the production router for the given remote chains.
// It also sets the onRamp and offRamp on the production router for the given remote chains.
// This function will add the proposals required to make these changes to the proposalAggregate slice.
func enableProductionRouter(
	e deployment.Environment,
	chainSelector uint64,
	remoteChainSelectors []uint64,
	mcmsConfig *changeset.MCMSConfig,
	proposalAggregate []mcmslib.TimelockProposal,
) error {
	// Update offRamp sources on the new chain.
	offRampUpdatesOnNew := make(map[uint64]OffRampSourceUpdate, len(remoteChainSelectors))
	for _, remoteChainSelector := range remoteChainSelectors {
		offRampUpdatesOnNew[remoteChainSelector] = OffRampSourceUpdate{
			TestRouter:                false,
			IsRMNVerificationDisabled: true, // TODO: We should eventually accept this as input.
			IsEnabled:                 true,
		}
	}
	out, err := UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
		UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
			chainSelector: offRampUpdatesOnNew,
		},
		MCMS: mcmsConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update onRamp destinations on the new chain.
	onRampUpdatesOnNew := make(map[uint64]OnRampDestinationUpdate, len(remoteChainSelectors))
	for _, remoteChainSelector := range remoteChainSelectors {
		onRampUpdatesOnNew[remoteChainSelector] = OnRampDestinationUpdate{
			TestRouter:       false,
			AllowListEnabled: false, // TODO: We should eventually accept this as input.
			IsEnabled:        true,
		}
	}
	out, err = UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
			chainSelector: onRampUpdatesOnNew,
		},
		MCMS: mcmsConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update router ramps on the new chain.
	offRampUpdates := make(map[uint64]bool, len(remoteChainSelectors))
	onRampUpdates := make(map[uint64]bool, len(remoteChainSelectors))
	for _, remoteChainSelector := range remoteChainSelectors {
		offRampUpdates[remoteChainSelector] = true
		onRampUpdates[remoteChainSelector] = true
	}
	out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
		TestRouter: false,
		UpdatesByChain: map[uint64]RouterUpdates{
			chainSelector: RouterUpdates{
				OnRampUpdates:  onRampUpdates,
				OffRampUpdates: offRampUpdates,
			},
		},
		MCMS: mcmsConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to run UpdateRouterRampsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	return nil
}
