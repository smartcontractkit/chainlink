package v1_6

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSet[AddRegistryModuleConfig] = AddRegistryModuleChangeset

type AddRegistryModuleConfig struct {
	// Map of chain selector to registry module 1.6 address
	RegistryModuleAddrs map[uint64]common.Address
	// MCMS config
	MCMSConfig *proposalutils.TimelockConfig
}

func (c AddRegistryModuleConfig) Validate(e cldf.Environment) error {
	if len(c.RegistryModuleAddrs) == 0 {
		return errors.New("no registry module addresses provided")
	}

	// Load state to check TokenAdminRegistry exists
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSel, addr := range c.RegistryModuleAddrs {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSel, err)
		}

		if _, exists := e.BlockChains.EVMChains()[chainSel]; !exists {
			return fmt.Errorf("chain %d not found in environment", chainSel)
		}

		if addr == (common.Address{}) {
			return fmt.Errorf("registry module address for chain %d is zero address", chainSel)
		}

		// Check if TokenAdminRegistry exists on the chain
		chainState, exists := state.Chains[chainSel]
		if !exists {
			return fmt.Errorf("chain state not found for chain %d", chainSel)
		}

		if chainState.TokenAdminRegistry == nil {
			return fmt.Errorf("TokenAdminRegistry not found on chain %d", chainSel)
		}

		if c.MCMSConfig == nil {
			return errors.New("mcmsConfig is required for this changeset")
		}
	}
	return nil
}

func AddRegistryModuleChangeset(e cldf.Environment, cfg AddRegistryModuleConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Collect operations per chain
	ops := make([]mcmstypes.BatchOperation, 0)
	chainsNeedingOps := make([]uint64, 0) // Track which chains need operations

	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSel, registryModuleAddr := range cfg.RegistryModuleAddrs {
		chainState := state.Chains[chainSel]
		timelocks[chainSel] = chainState.Timelock.Address().Hex()

		inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
		}

		// Check if the 1.6 registry module we want to add is already registered
		isAlreadyModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if registry module is already added on chain %d: %w", chainSel, err)
		}

		if isAlreadyModule {
			e.Logger.Infow("RegistryModule 1.6 already added to TokenAdminRegistry, skipping",
				"chain", chainSel,
				"registryModule", registryModuleAddr.Hex())
			continue
		}
		// Track that this chain needs an operation
		chainsNeedingOps = append(chainsNeedingOps, chainSel)

		// Create add operation for new 1.6 module
		// Parse the ABI and encode the addRegistryModule call data
		parsedABI, err := abi.JSON(strings.NewReader(token_admin_registry.TokenAdminRegistryABI))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to parse TokenAdminRegistry ABI on chain %d: %w", chainSel, err)
		}

		callData, err := parsedABI.Pack("addRegistryModule", registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to encode addRegistryModule call data on chain %d: %w", chainSel, err)
		}

		op, err := proposalutils.BatchOperationForChain(
			chainSel, chainState.TokenAdminRegistry.Address().String(), callData, big.NewInt(0), shared.TokenAdminRegistry.String(), nil)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create batch operation for chain %d: %w", chainSel, err)
		}
		ops = append(ops, op)

		e.Logger.Infow("Added add operation to batch",
			"chain", chainSel,
			"newModule", registryModuleAddr.Hex())
	}

	// If no operations were needed, return early
	if len(chainsNeedingOps) == 0 {
		e.Logger.Info("No registry module operations needed - all modules already registered")
		return cldf.ChangesetOutput{}, nil
	}

	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMSConfig, nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcm addresses per chain: %w", err)
	}

	// Only iterate over chains that need operations
	for _, chainSel := range chainsNeedingOps {
		chainState := state.Chains[chainSel]
		// Safe to access Timelock here because validation already checked it exists
		timelocks[chainSel] = chainState.Timelock.Address().Hex()

		inspector, err := proposalutils.McmsInspectorForChain(e, chainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
		}
		inspectors[chainSel] = inspector
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		ops,
		"Add Registry Module 1.6 to TokenAdminRegistry",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
	}, nil
}
