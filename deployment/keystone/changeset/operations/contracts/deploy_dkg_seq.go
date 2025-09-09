package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	contracts "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type DeployDKGContractSequenceDeps struct {
	Env *cldf.Environment
}

type DeployDKGContractSequenceInput struct {
	ChainSelector uint64
	Qualifier     string // qualifier for the DKG contract deployment
}

type DeployDKGContractSequenceOutput struct {
	// TODO: CRE-742 remove AddressBook
	AddressBook cldf.AddressBook // Keeping the address store for backward compatibility, as not everything has been migrated to datastore
	Datastore   datastore.DataStore
}

// DeployDKGContractsSequence is a sequence that deploys the DKG contract.
// TODO dedup with sequence in ocr3/v2/changeset/sequences/deploy_ocr3.go CRE-803
var DeployDKGContractsSequence = operations.NewSequence[DeployDKGContractSequenceInput, DeployDKGContractSequenceOutput, DeployDKGContractSequenceDeps](
	"deploy-dkg-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy DKG Contracts",
	func(b operations.Bundle, deps DeployDKGContractSequenceDeps, input DeployDKGContractSequenceInput) (output DeployDKGContractSequenceOutput, err error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// DKG Contract
		dkgDeployReport, err := operations.ExecuteOperation(b, contracts.DeployDKG, contracts.DeployDKGDeps(deps), contracts.DeployDKGInput{ChainSelector: input.ChainSelector, Qualifier: input.Qualifier})
		if err != nil {
			return DeployDKGContractSequenceOutput{}, fmt.Errorf("failed to execution operation DeployDKG: %w", err)
		}
		err = updateAddresses(as.Addresses(), dkgDeployReport.Output.Datastore.Addresses(), ab, dkgDeployReport.Output.AddressBook)
		if err != nil {
			return DeployDKGContractSequenceOutput{}, fmt.Errorf("failed to update addresses after DKG deployment: %w", err)
		}
		return DeployDKGContractSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil

	},
)
