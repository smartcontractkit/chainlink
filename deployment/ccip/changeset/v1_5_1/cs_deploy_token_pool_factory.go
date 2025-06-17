package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool_factory"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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

		tokenPoolFactory, err := cldf.DeployContract(e.Logger, chain, addressBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_pool_factory.TokenPoolFactory] {
				address, tx, tokenPoolFactory, err := token_pool_factory.DeployTokenPoolFactory(
					chain.DeployerKey,
					chain.Client,
					chainState.TokenAdminRegistry.Address(),
					registryModuleAddress,
					chainState.RMNProxy.Address(),
					chainState.Router.Address(),
				)

				return cldf.ContractDeploy[*token_pool_factory.TokenPoolFactory]{
					Address:  address,
					Contract: tokenPoolFactory,
					Tx:       tx,
					Tv:       cldf.NewTypeAndVersion(shared.TokenPoolFactory, deployment.Version1_5_1),
					Err:      err,
				}
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pool factory: %w", err)
		}
		e.Logger.Infof("Successfully deployed token pool factory %s on %s", tokenPoolFactory.Address.String(), chain.String())
	}

	return cldf.ChangesetOutput{AddressBook: addressBook}, nil
}
