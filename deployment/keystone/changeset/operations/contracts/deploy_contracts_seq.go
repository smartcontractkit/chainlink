package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cap_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	wf_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

type (
	EVMChainID uint64
	Selector   uint64
)

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
var DeployRegistryContractsSequence = operations.NewSequence(
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

// DeployV2RegistryContractsSequence is a sequence that deploys the the required registry contracts (Capabilities Registry, Workflow Registry).
var DeployV2RegistryContractsSequence = operations.NewSequence(
	// do not add optional contracts here (ocr, forwarder...), as this sequence is used to deploy the registry contracts that other sequences depend on
	"deploy-v2-registry-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy V2 registry Contracts (Capabilities Registry, Workflow Registry)",
	func(b operations.Bundle, deps DeployContractsSequenceDeps, input DeployRegistryContractsSequenceInput) (output DeployContractSequenceOutput, err error) {
		ab := deployment.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// Capabilities Registry contract
		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, cap_reg_v2.DeployCapabilitiesRegistry, cap_reg_v2.DeployCapabilitiesRegistryDeps(deps), cap_reg_v2.DeployCapabilitiesRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}

		v1Output, err := toV1Output(capabilitiesRegistryDeployReport.Output)
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}

		if err = updateAddresses(as.Addresses(), v1Output.Addresses, ab, v1Output.AddressBook); err != nil {
			return DeployContractSequenceOutput{}, err
		}

		// Workflow Registry contract
		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, wf_reg_v2.DeployWorkflowRegistryOp, wf_reg_v2.DeployWorkflowRegistryOpDeps(deps), wf_reg_v2.DeployWorkflowRegistryOpInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}

		v1Output, err = toV1Output(workflowRegistryDeployReport.Output)
		if err != nil {
			return DeployContractSequenceOutput{}, err
		}

		err = updateAddresses(as.Addresses(), v1Output.Addresses, ab, v1Output.AddressBook)
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

type DeprecatedOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook deployment.AddressBook
}

// toV1Output transforms a v2 output to a common output format that uses the deprecated
// address book.
func toV1Output(in any) (DeprecatedOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	labels := deployment.NewLabelSet()
	var r datastore.AddressRef

	switch v := in.(type) {
	case cap_reg_v2.DeployCapabilitiesRegistryOutput:
		r = datastore.AddressRef{
			ChainSelector: v.ChainSelector,
			Address:       v.Address,
			Type:          datastore.ContractType(v.Type),
			Version:       semver.MustParse(v.Version),
			Qualifier:     v.Qualifier,
			Labels:        datastore.NewLabelSet(v.Labels...),
		}
		for _, l := range v.Labels {
			labels.Add(l)
		}
	case wf_reg_v2.DeployWorkflowRegistryOpOutput:
		r = datastore.AddressRef{
			ChainSelector: v.ChainSelector,
			Address:       v.Address,
			Type:          datastore.ContractType(v.Type),
			Version:       semver.MustParse(v.Version),
			Qualifier:     v.Qualifier,
			Labels:        datastore.NewLabelSet(v.Labels...),
		}
		for _, l := range v.Labels {
			labels.Add(l)
		}
	default:
		return DeprecatedOutput{}, fmt.Errorf("unsupported input type for transform: %T", in)
	}

	if err := ds.Addresses().Add(r); err != nil {
		return DeprecatedOutput{}, fmt.Errorf("failed to add address ref: %w", err)
	}

	if err := ab.Save(r.ChainSelector, r.Address, deployment.TypeAndVersion{
		Type:    deployment.ContractType(r.Type),
		Version: *r.Version,
		Labels:  labels,
	}); err != nil {
		return DeprecatedOutput{}, fmt.Errorf("failed to save address to address book: %w", err)
	}

	return DeprecatedOutput{
		Addresses:   ds.Addresses(),
		AddressBook: ab,
	}, nil
}
