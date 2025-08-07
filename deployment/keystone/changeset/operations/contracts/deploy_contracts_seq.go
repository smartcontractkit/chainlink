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
	DeployVaultOCR3       bool
	DeployEVMOCR3         bool
	DeployConsensusOCR3   bool
}

type DeployKeystoneContractsSequenceOutput struct {
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

// DeployKeystoneContractsSequence is a sequence that deploys the Keystone contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder).
var DeployKeystoneContractsSequence = operations.NewSequence[DeployKeystoneContractsSequenceInput, DeployKeystoneContractsSequenceOutput, DeployKeystoneContractsSequenceDeps](
	"deploy-keystone-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Contracts (OCR3, Vault-OCR3, EVM-OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder)",
	func(b operations.Bundle, deps DeployKeystoneContractsSequenceDeps, input DeployKeystoneContractsSequenceInput) (output DeployKeystoneContractsSequenceOutput, err error) {
		ab := deployment.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// OCR3 Contract
		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_ocr3"})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), ocr3DeployReport.Output.Addresses, ab, ocr3DeployReport.Output.AddressBook)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		// Capabilities Registry contract
		capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, DeployCapabilityRegistryOpDeps(deps), DeployCapabilityRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), capabilitiesRegistryDeployReport.Output.Addresses, ab, capabilitiesRegistryDeployReport.Output.AddressBook)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		// Workflow Registry contract
		workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, DeployWorkflowRegistryOpDeps(deps), DeployWorkflowRegistryInput{ChainSelector: input.RegistryChainSelector})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), workflowRegistryDeployReport.Output.Addresses, ab, workflowRegistryDeployReport.Output.AddressBook)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		// Keystone Forwarder contract
		keystoneForwarderDeployReport, err := operations.ExecuteSequence(b, DeployKeystoneForwardersSequence, DeployKeystoneForwardersSequenceDeps(deps), DeployKeystoneForwardersInput{Targets: input.ForwardersSelectors})
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), keystoneForwarderDeployReport.Output.Addresses, ab, keystoneForwarderDeployReport.Output.AddressBook)
		if err != nil {
			return DeployKeystoneContractsSequenceOutput{}, err
		}

		if input.DeployVaultOCR3 {
			// Vault OCR3 Contract
			vaultOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_vault"})
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), vaultOCR3DeployReport.Output.Addresses, ab, vaultOCR3DeployReport.Output.AddressBook)
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
		}
		if input.DeployEVMOCR3 {
			// EVM cap OCR3 Contract
			evmOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_evm"})
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), evmOCR3DeployReport.Output.Addresses, ab, evmOCR3DeployReport.Output.AddressBook)
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
		}
		if input.DeployConsensusOCR3 {
			evmOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_consensus"})
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), evmOCR3DeployReport.Output.Addresses, ab, evmOCR3DeployReport.Output.AddressBook)
			if err != nil {
				return DeployKeystoneContractsSequenceOutput{}, err
			}
		}

		return DeployKeystoneContractsSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil
	},
)
