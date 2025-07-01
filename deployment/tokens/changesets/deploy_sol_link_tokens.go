package changesets

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/seqs"
)

// DeploySolLinkTokens deploys Link Token contracts to the specified Solana chains in the
// DeployLinkTokenInput. The qualifiers and labels will be used to tag all deployed contracts in
var DeploySolLinkTokens cldf.ChangeSetV2[DeployLinkTokensInput] = deploySolLinkTokens{}

// deploySolLinkTokens is the implementation of the DeploySolLinkTokens changeset.
type deploySolLinkTokens struct{}

// VerifyPreconditions ensures that all listed chain selectors are valid and available in the
// environment chains.
func (deploySolLinkTokens) VerifyPreconditions(
	e cldf.Environment, input DeployLinkTokensInput,
) error {
	if err := validateNoDupeSelectors(input.ChainSelectors); err != nil {
		return err
	}

	if err := validateChainSelectorsFamily(
		input.ChainSelectors, chain_selectors.FamilySolana,
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

// Apply executes the SeqDeploySolTokens sequence to deploy Link Token contracts to the specified
// chains.
func (deploySolLinkTokens) Apply(
	e cldf.Environment, input DeployLinkTokensInput,
) (cldf.ChangesetOutput, error) {
	var (
		out = cldf.ChangesetOutput{
			AddressBook: cldf.NewMemoryAddressBook(),
			DataStore:   datastore.NewMemoryDataStore(),
		}

		seqDeps = seqs.SeqDeploySolTokensDeps{
			SolChains: e.BlockChains.SolanaChains(),
			AddrBook:  out.AddressBook, //nolint:staticcheck // Will be removed once the address book is no longer required.
			Datastore: out.DataStore,
		}
		seqInput = seqs.SeqDeploySolTokensInput{
			ChainSelectors: input.ChainSelectors,
			Qualifier:      input.Qualifier,
			Labels:         input.Labels,
		}
	)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle, seqs.SeqDeploySolTokens, seqDeps, seqInput,
	)

	out.Reports = seqReport.ExecutionReports

	return out, err
}
