package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type EVMChainID uint64
type Selector uint64

// inputs and outputs have to be serializable, and must not contain sensitive data
type DeployContractsSequenceDeps struct {
	Env *deployment.Environment
}

type DeployRegistryContractsSequenceInput struct {
	RegistryChainSelector uint64
}
type DeployContractSequenceOutput struct {
	// Not sure if we can serialize the address book without modifications, but whatever is returned needs to be serializable.
	// This could also be the address datastore instead.
	AddressBook deployment.AddressBook
	Datastore   datastore.DataStore // Keeping the address store for backward compatibility, as not everything has been migrated to address book
}

func updateAddresses(addr datastore.MutableAddressRefStore, as datastore.AddressRefStore, sourceAB deployment.AddressBook, ab deployment.AddressBook) error {
	addresses, err := as.Fetch()
	if err != nil {
		return err
	}
	for _, a := range addresses {
		if err := addr.Add(a); err != nil {
			return err
		}
	}

	return sourceAB.Merge(ab)
}

// DeployRegistryContractsSequence is a sequence that deploys the the required registry contracts (Capabilities Registry, Workflow Registry).
var DeployRegistryContractsSequence = operations.NewSequence[DeployRegistryContractsSequenceInput, DeployContractSequenceOutput, DeployContractsSequenceDeps](
	// do not add optional contracts here (ocr, forwarder...), as this sequence is used to deploy the registry contracts that other sequences depend on
	"deploy-registry-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy registry Contracts (Capabilities Registry, Workflow Registry)",
	func(b operations.Bundle, deps DeployContractsSequenceDeps, input DeployRegistryContractsSequenceInput) (output DeployContractSequenceOutput, err error) {
		ab := deployment.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// Capabilities Registry contract
		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, DeployCapabilityRegistryOpDeps(deps), DeployCapabilityRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), capabilitiesRegistryDeployReport.Output.Addresses, ab, capabilitiesRegistryDeployReport.Output.AddressBook)
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}
		// Workflow Registry contract
		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, DeployWorkflowRegistryOpDeps(deps), DeployWorkflowRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), workflowRegistryDeployReport.Output.Addresses, ab, workflowRegistryDeployReport.Output.AddressBook)
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}
		return DeployContractSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil

	},
)

func CapabilityContractIdentifier(chainID uint64) string {
	return fmt.Sprintf("capability_evm_%d", chainID)
}
