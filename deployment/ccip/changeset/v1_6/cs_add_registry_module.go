package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
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
		return fmt.Errorf("no registry module addresses provided")
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

		// Create add operation for new 1.6 module
		addTx, err := chainState.TokenAdminRegistry.AddRegistryModule(nil, registryModuleAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create addRegistryModule transaction on chain %d: %w", chainSel, err)
		}
		op, err := proposalutils.BatchOperationForChain(
			chainSel, chainState.TokenAdminRegistry.Address().String(), addTx.Data(), big.NewInt(0), shared.TokenAdminRegistry.String(), nil)

		ops = append(ops, op)

		e.Logger.Infow("Added add operation to batch",
			"chain", chainSel,
			"newModule", registryModuleAddr.Hex())
	}

	// If no operations needed, return early
	if len(ops) == 0 {
		e.Logger.Info("No registry module operations needed")
		return cldf.ChangesetOutput{}, nil
	}

	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMSConfig, nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcm addresses per chain: %w", err)
	}
	// Generate MCMS proposal using proposalutils
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		ops,
		"PermaBless commit stores on RMN",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{
		*proposal,
	}}, nil
}

// isRegistryModule16 checks if the given registry module address is a 1.6 version
func isRegistryModule16(chainState evm.CCIPChainState, registryModuleAddr common.Address) bool {
	// Check if the address exists in the 1.6 registry modules map
	for _, module16 := range chainState.RegistryModules1_6 {
		if module16.Address() == registryModuleAddr {
			return true
		}
	}
	return false
}
