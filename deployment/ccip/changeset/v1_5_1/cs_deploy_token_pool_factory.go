package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipopsv1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

// DeployTokenPoolFactoryChangeset is a changeset that deploys the TokenPoolFactory contract on multiple chains.
// In most cases, running DeployPrerequisitesChangeset will be sufficient to deploy the TokenPoolFactory.
// However, if a chain has multiple registry modules with version 1.6.0 and you want to specify which one to use,
// you can use this changeset to do so.
var DeployTokenPoolFactoryChangeset = cldf.CreateChangeSet(deployTokenPoolFactoryLogic, deployTokenPoolFactoryPrecondition)

type DeployTokenPoolFactoryConfig struct {
	// Chains is the list of chains on which to deploy the token pool factory.
	Chains []uint64
	// RegistryModule1_6Addresses indicates which registry module to use for each chain.
	// If the chain only has one 1.6.0 registry module, you do not need to specify it here.
	RegistryModule1_6Addresses map[uint64]common.Address
}

func deployTokenPoolFactoryPrecondition(e cldf.Environment, config DeployTokenPoolFactoryConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chainSel := range config.Chains {
		err := stateview.ValidateChain(e, state, chainSel, nil)
		if err != nil {
			return fmt.Errorf("failed to validate chain with selector %d: %w", chainSel, err)
		}
		chain := e.BlockChains.EVMChains()[chainSel]
		state := state.Chains[chainSel]
		if state.TokenPoolFactory != nil {
			return fmt.Errorf("token pool factory already deployed on %s", chain.String())
		}
		if state.TokenAdminRegistry == nil {
			return fmt.Errorf("token admin registry does not exist on %s", chain.String())
		}
		if state.Router == nil {
			return fmt.Errorf("router does not exist on %s", chain.String())
		}
		if state.RMNProxy == nil {
			return fmt.Errorf("rmn proxy does not exist on %s", chain.String())
		}
		if len(state.RegistryModules1_6) == 0 {
			return fmt.Errorf("registry module with version 1.6.0 does not exist on %s", chain.String())
		}
		// There can be multiple registry modules with version 1.6.0 on a chain, but only one can be used for the token pool factory.
		// If the user has specified a registry module address, check that it exists on the chain.
		// If the user has not specified a registry module address, check that there is only one registry module with version 1.6.0 on the chain.
		// If there are multiple registry modules with version 1.6.0, the user MUST specify which one to use by providing the address.
		registryModuleAddress, ok := config.RegistryModule1_6Addresses[chainSel]
		if !ok && len(state.RegistryModules1_6) > 1 {
			return fmt.Errorf("multiple registry modules with version 1.6.0 exist on %s, must specify using RegistryModule1_6Addresses", chain.String())
		} else if ok {
			registryModuleExists := false
			for _, registryModule := range state.RegistryModules1_6 {
				if registryModuleAddress == registryModule.Address() {
					registryModuleExists = true
					break
				}
			}
			if !registryModuleExists {
				return fmt.Errorf("no registry module with version 1.6.0 and address %s found on %s", registryModuleAddress.String(), chain.String())
			}
		}
	}

	return nil
}

func deployTokenPoolFactoryLogic(e cldf.Environment, config DeployTokenPoolFactoryConfig) (cldf.ChangesetOutput, error) {
	addressBook := cldf.NewMemoryAddressBook()
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chainSel := range config.Chains {
		chain := e.BlockChains.EVMChains()[chainSel]
		chainState := state.Chains[chainSel]

		registryModuleAddress, ok := config.RegistryModule1_6Addresses[chainSel]
		if !ok {
			registryModuleAddress = chainState.RegistryModules1_6[0].Address()
		}

		tpfReport, err := operations.ExecuteOperation(e.OperationsBundle, ccipopsv1_5_1.DeployTokenPoolFactoryOp, chain, opsutil.EVMDeployInput[ccipopsv1_5_1.DeployTokenPoolFactoryInput]{
			ChainSelector: chain.ChainSelector(),
			DeployInput: ccipopsv1_5_1.DeployTokenPoolFactoryInput{
				ChainSelector:              chain.ChainSelector(),
				TokenAdminRegistry:         chainState.TokenAdminRegistry.Address(),
				RegistryModule1_6Addresses: registryModuleAddress,
				RMNProxy:                   chainState.RMNProxy.Address(),
				Router:                     chainState.Router.Address(),
			},
		})
		if err != nil {
			e.Logger.Errorw("Failed to deploy token pool factory", "chain", chain.String(), "err", err)
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool factory: %w", err)
		}

		err = addressBook.Save(chainSel, tpfReport.Output.Address.Hex(), cldf.MustTypeAndVersionFromString(tpfReport.Output.TypeAndVersion))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address %s for chain %d: %w", tpfReport.Output.Address.Hex(), chainSel, err)
		}
		e.Logger.Infof("Successfully deployed token pool factory %s on %s", tpfReport.Output.Address.Hex(), chain.String())
	}

	return cldf.ChangesetOutput{AddressBook: addressBook}, nil
}
