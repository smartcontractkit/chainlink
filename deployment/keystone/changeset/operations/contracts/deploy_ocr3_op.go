package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type DeployOCR3ContractSequenceDeps struct {
	Env *cldf.Environment
}

type DeployOCR3ContractSequenceInput struct {
	ChainSelector uint64
	Qualifier     string // qualifier for the OCR3 contract deployment
}

type DeployOCR3ContractSequenceOutput struct {
	// TODO: CRE-742 remove AddressBook
	AddressBook cldf.AddressBook // Keeping the address store for backward compatibility, as not everything has been migrated to datastore
	Datastore   datastore.DataStore
}

// DeployOCR3ContractsSequence is a sequence that deploys the OCR3 contract.
// TODO dedup with sequence in ocr3/v2/changeset/sequences/deploy_ocr3.go CRE-803
var DeployOCR3ContractsSequence = operations.NewSequence[DeployOCR3ContractSequenceInput, DeployOCR3ContractSequenceOutput, DeployOCR3ContractSequenceDeps](
	"deploy-ocr3-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contracts",
	func(b operations.Bundle, deps DeployOCR3ContractSequenceDeps, input DeployOCR3ContractSequenceInput) (output DeployOCR3ContractSequenceOutput, err error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// OCR3 Contract
		deployInput := contracts.DeployOCR3Input{
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
		}
		ocr3DeployReport, err := operations.ExecuteOperation(b, contracts.DeployOCR3, contracts.DeployOCR3Deps(deps), deployInput)
		if err != nil {
			return DeployOCR3ContractSequenceOutput{}, fmt.Errorf("failed to execution operation DeployOCR3: %w", err)
		}
		err = updateAddresses(as.Addresses(), ocr3DeployReport.Output.Datastore.Addresses(), ab, ocr3DeployReport.Output.AddressBook)
		if err != nil {
			return DeployOCR3ContractSequenceOutput{}, fmt.Errorf("failed to update addresses after OCR3 deployment: %w", err)
		}
		return DeployOCR3ContractSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil

	},
)
