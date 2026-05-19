package aptos

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/aptos-labs/aptos-go-sdk"
)

var _ cldf.ChangeSetV2[config.DeployRegulatedTokenConfig] = DeployRegulatedToken{}

// DeployRegulatedToken deploys and initializes a regulated token directly with the
// deployer signer (regulated_token cannot be deployed via MCMS due to DFA re-entrancy),
// transfers ownership and admin role to the MCMS registry owner, and returns a
// timelock proposal containing accept_ownership and accept_admin.
type DeployRegulatedToken struct{}

func (cs DeployRegulatedToken) VerifyPreconditions(env cldf.Environment, cfg config.DeployRegulatedTokenConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	aptosState, ok := state.AptosChains[cfg.ChainSelector]
	if !ok {
		errs = append(errs, fmt.Errorf("aptos chain %d not in state", cfg.ChainSelector))
	} else {
		if aptosState.CCIPAddress == (aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("CCIP is not deployed on Aptos chain %d", cfg.ChainSelector))
		}
		if aptosState.MCMSAddress == (aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("MCMS is not deployed on Aptos chain %d", cfg.ChainSelector))
		}
	}
	return errors.Join(errs...)
}

func (cs DeployRegulatedToken) Apply(env cldf.Environment, cfg config.DeployRegulatedTokenConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	aptosChain := env.BlockChains.AptosChains()[cfg.ChainSelector]
	mcmsAddress := state.AptosChains[cfg.ChainSelector].MCMSAddress

	ab := cldf.NewMemoryAddressBook()
	deps := dependency.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	seqReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		seq.DeployRegulatedTokenSequence,
		deps,
		seq.DeployRegulatedTokenSeqInput{
			MCMSAddress:          mcmsAddress,
			TokenParams:          cfg.TokenParams,
			TokenMint:            cfg.TokenMint,
			RegistrarPreregister: cfg.RegistrarPreregisterOrDefault(),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	out := seqReport.Output
	typeAndVersion := cldf.NewTypeAndVersion(shared.AptosRegulatedTokenType, deployment.Version1_6_0)
	typeAndVersion.AddLabel(string(cfg.TokenParams.Symbol))
	if err := ab.Save(aptosChain.Selector, out.TokenCodeObjectAddress.StringLong(), typeAndVersion); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("save regulated token code object: %w", err)
	}
	typeAndVersion = cldf.NewTypeAndVersion(cldf.ContractType(cfg.TokenParams.Symbol), deployment.Version1_6_0)
	if err := ab.Save(aptosChain.Selector, out.TokenMetadataAddress.StringLong(), typeAndVersion); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("save token metadata address: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		mcmsAddress,
		cfg.ChainSelector,
		out.MCMSOperations,
		"Accept regulated token ownership and admin role (after direct deploy)",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("generate MCMS proposal: %w", err)
	}

	ds, err := shared.PopulateDataStore(ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("populate datastore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               seqReport.ExecutionReports,
	}, nil
}
