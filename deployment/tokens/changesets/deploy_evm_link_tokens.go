package changesets

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/seqs"
)

// DeployLinkTokensInput contains the selectors of the chains to which the Link Token contract
// should be deployed.
//
// This input is shared between all changesets that deploy Link Token contracts on a specific
// chain.
type DeployLinkTokensInput struct {
	// Required: ChainSelectors are the chain selectors of the chains to which the Link Token contract
	// should be deployed.
	ChainSelectors []uint64
	// Optional: The Qualifier is a string that will be used to tag the deployed contracts in the datastore
	// to uniquely identify multiple versions of the same contract on the same chain.
	//
	// Default: "" (empty string)
	Qualifier string
	// Optional: The Labels are a list of labels that will be used to tag the deployed contracts in the
	// address book and datastore.
	//
	// Default: []string{}
	Labels []string
}

// DeployEVMLinkTokens deploys burn/mint Link Token contracts to the specified chains in the
// DeployLinkTokenInput and adds the addresses to the addressbook and datastore. The qualifiers and
// labels will be used to tag all deployed contracts in the address book and datastore.
//
// This changeset should be used to deploy Link Tokens for all new deployments.
var DeployEVMLinkTokens cldf.ChangeSetV2[DeployLinkTokensInput] = deployEVMLinkTokens{}

// deployEVMLinkTokens is the implementation of the DeployEVMLinkTokens changeset.
type deployEVMLinkTokens struct{}

// VerifyPreconditions ensures that all listed chain selectors are valid and available in the
// environment chains.
func (deployEVMLinkTokens) VerifyPreconditions(
	e cldf.Environment, input DeployLinkTokensInput,
) error {
	if err := validateNoDupeSelectors(input.ChainSelectors); err != nil {
		return err
	}

	if err := validateChainSelectorsFamily(
		input.ChainSelectors, chain_selectors.FamilyEVM,
	); err != nil {
		return err
	}

	if err := validateNoExistingLinkToken(
		input.ChainSelectors, input.Qualifier, e.ExistingAddresses, e.DataStore,
	); err != nil {
		return err
	}

	return validateSelectorsInEnvironment(e.BlockChains, input.ChainSelectors)
}

// Apply executes the SeqDeployEVMTokens sequence to deploy Link Token contracts to the specified
// chains.
func (deployEVMLinkTokens) Apply(
	e cldf.Environment, input DeployLinkTokensInput,
) (cldf.ChangesetOutput, error) {
	var (
		out = cldf.ChangesetOutput{
			AddressBook: cldf.NewMemoryAddressBook(),
			DataStore:   datastore.NewMemoryDataStore(),
		}

		seqDeps = seqs.SeqDeployEVMTokensDeps{
			EVMChains: e.BlockChains.EVMChains(),
			AddrBook:  out.AddressBook, //nolint:staticcheck // Will be removed once the address book is no longer required.
			Datastore: out.DataStore,
		}
		seqInput = seqs.SeqDeployEVMTokensInput{
			ChainSelectors: input.ChainSelectors,
			Qualifier:      input.Qualifier,
			Labels:         input.Labels,
		}
	)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle, seqs.SeqDeployEVMTokens, seqDeps, seqInput,
	)

	out.Reports = seqReport.ExecutionReports

	return out, err
}
