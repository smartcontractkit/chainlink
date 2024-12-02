package crib

import (
	"context"
	"errors"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can be started with the CR address in Capabilities.ExternalRegistry
func DeployHomeChainContracts(lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64) (deployment.CapabilityRegistryConfig, deployment.AddressBook, error) {
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}

	ab := deployment.NewMemoryAddressBook()

	capReg, err := ccipdeployment.DeployCapReg(lggr, ccipdeployment.CCIPOnChainState{}, ab, chains[homeChainSel])
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	return deployment.CapabilityRegistryConfig{
		NetworkType: relay.NetworkEVM,
		EVMChainID:  evmChainID,
		Contract:    capReg.Address,
	}, ab, nil
}

func DeployCCIPAndAddLanes(ctx context.Context, lggr logger.Logger, envCfg devenv.EnvironmentConfig, ab deployment.AddressBook) (DeployCCIPOutput, error) {
	env, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envCfg)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	if env == nil {
		return DeployCCIPOutput{}, errors.New("environment is nil")
	}
	selectors := env.AllChainSelectors()
	// deploy pre-requisites
	prerequisites, err := changeset.DeployPrerequisites(*env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: selectors,
	})
	err = env.ExistingAddresses.Merge(prerequisites.AddressBook)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	// Get new state after migration.
	state, err := ccipdeployment.LoadOnchainState(*env)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	// Apply the jobs.
	for nodeID, jobs := range prerequisites.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := env.Offchain.ProposeJob(context.Background(),
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			if err != nil {
				return DeployCCIPOutput{}, err
			}
		}
	}

	// Add all lanes
	err = ccipdeployment.AddLanesForAll(*env, state)
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	err = env.ExistingAddresses.Merge(prerequisites.AddressBook)
	addresses, err := ab.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, err
	}
	return DeployCCIPOutput{
		AddressesByChain: addresses,
		NodeIDs:          maps.Keys(prerequisites.JobSpecs),
	}, err
}
