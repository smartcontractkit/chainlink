package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployOCR3OpDeps struct {
	Env *cldf.Environment
}

type DeployOCR3OpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployOCR3OpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook
}

// DeployOCR3Op is an operation that deploys the OCR3 contract.
var DeployOCR3Op = operations.NewOperation[DeployOCR3OpInput, DeployOCR3OpOutput, DeployOCR3OpDeps](
	"deploy-ocr3-op",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contract",
	func(b operations.Bundle, deps DeployOCR3OpDeps, input DeployOCR3OpInput) (DeployOCR3OpOutput, error) {
		ocr3Output, err := changeset.DeployOCR3V2(*deps.Env, &changeset.DeployRequestV2{
			ChainSel:  input.ChainSelector,
			Qualifier: input.Qualifier,
		})
		if err != nil {
			return DeployOCR3OpOutput{}, fmt.Errorf("DeployOCR3Op error: failed to deploy OCR3 contract: %w", err)
		}

		return DeployOCR3OpOutput{
			Addresses: ocr3Output.DataStore.Addresses(), AddressBook: ocr3Output.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)

type DeployOCR3ContractSequenceDeps struct {
	Env *cldf.Environment
}

type DeployOCR3ContractSequenceInput struct {
	RegistryChainSelector uint64
	Qualifier             string // qualifier for the OCR3 contract deployment
}

type DeployOCR3ContractSequenceOutput struct {
	// TODO: CRE-742 remove AddressBook
	AddressBook cldf.AddressBook // Keeping the address store for backward compatibility, as not everything has been migrated to datastore
	Datastore   datastore.DataStore
}

// DeployKeystoneContractsSequence is a sequence that deploys the Keystone contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder).
var DeployOCR3ContractsSequence = operations.NewSequence[DeployOCR3ContractSequenceInput, DeployOCR3ContractSequenceOutput, DeployOCR3ContractSequenceDeps](
	"deploy-registry-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy registry Contracts (Capabilities Registry, Workflow Registry)",
	func(b operations.Bundle, deps DeployOCR3ContractSequenceDeps, input DeployOCR3ContractSequenceInput) (output DeployOCR3ContractSequenceOutput, err error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// OCR3 Contract
		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: input.Qualifier})
		if err != nil {
			return DeployOCR3ContractSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), ocr3DeployReport.Output.Addresses, ab, ocr3DeployReport.Output.AddressBook)
		if err != nil {
			return DeployOCR3ContractSequenceOutput{}, err
		}
		return DeployOCR3ContractSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil

	},
)
