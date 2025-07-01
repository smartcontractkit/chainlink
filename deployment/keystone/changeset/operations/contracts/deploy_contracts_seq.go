package contracts

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployKeystoneContractsSequenceDeps struct {
	Env *deployment.Environment
}

// inputs and outputs have to be serializable, and must not contain sensitive data

type DeployKeystoneContractsSequenceInput struct {
	RegistryChainSelector uint64
	ForwardersSelectors   []uint64
}

type DeployKeystoneContractsSequenceOutput struct {
	// Not sure if we can serialize the address book without modifications, but whatever is returned needs to be serializable.
	// This could also be the address datastore instead.
	AddressBook deployment.AddressBook
	Datastore   datastore.DataStore // Keeping the address store for backward compatibility, as not everything has been migrated to address book
}

// DeployKeystoneContractsSequence is a sequence that deploys the Keystone contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder).
var DeployKeystoneContractsSequence = operations.NewSequence[DeployKeystoneContractsSequenceInput, DeployKeystoneContractsSequenceOutput, DeployKeystoneContractsSequenceDeps](
	"deploy-keystone-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder)",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input DeployKeystoneContractsSequenceInput) (output DeployKeystoneContractsSequenceOutput, err error) {
		ab := deployment.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, DeployCapabilityRegistryOpDeps(deps), DeployCapabilityRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, DeployWorkflowRegistryOpDeps(deps), DeployWorkflowRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		keystoneForwarderDeployReport, err := operations.ExecuteSequence(b, DeployKeystoneForwardersSequence, DeployKeystoneForwardersSequenceDeps(deps), DeployKeystoneForwardersInput{Targets: input.ForwardersSelectors})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		// Merge the address book and datastore from the deployed contracts
		var allResultingAddresses []datastore.AddressRef
		ocr3Addrs, err := ocr3DeployReport.Output.Addresses.Fetch()
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		allResultingAddresses = append(allResultingAddresses, ocr3Addrs...)
		capabilitiesRegistryAddrs, err := capabilitiesRegistryDeployReport.Output.Addresses.Fetch()
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		allResultingAddresses = append(allResultingAddresses, capabilitiesRegistryAddrs...)
		workflowRegistryAddrs, err := workflowRegistryDeployReport.Output.Addresses.Fetch()
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		allResultingAddresses = append(allResultingAddresses, workflowRegistryAddrs...)
		keystoneForwarderAddrs, err := keystoneForwarderDeployReport.Output.Addresses.Fetch()
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		allResultingAddresses = append(allResultingAddresses, keystoneForwarderAddrs...)
		for _, addr := range allResultingAddresses {
			if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
				return DeployKeystoneContractsSequenceOutput{}, addrRefErr
			}
		}

		if err := ab.Merge(ocr3DeployReport.Output.AddressBook); err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		if err := ab.Merge(capabilitiesRegistryDeployReport.Output.AddressBook); err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		if err := ab.Merge(workflowRegistryDeployReport.Output.AddressBook); err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		if err := ab.Merge(keystoneForwarderDeployReport.Output.AddressBook); err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		return DeployKeystoneContractsSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil
	},
)
