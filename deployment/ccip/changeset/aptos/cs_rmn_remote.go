package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// CurseSubjects changeset: curses multiple subjects on Aptos RMN Remote.
var _ cldf.ChangeSetV2[config.RMNRemoteCurseInput] = CurseSubjects{}

type CurseSubjects struct{}

func (cs CurseSubjects) VerifyPreconditions(env cldf.Environment, cfg config.RMNRemoteCurseInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	if (state.AptosChains[cfg.ChainSelector].CCIPAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("CCIP is not deployed on Aptos chain %d", cfg.ChainSelector))
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if len(cfg.Subjects) == 0 {
		errs = append(errs, errors.New("subjects list cannot be empty"))
	}
	return errors.Join(errs...)
}

func (cs CurseSubjects) Apply(env cldf.Environment, cfg config.RMNRemoteCurseInput) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.BlockChains.AptosChains()[cfg.ChainSelector]
	ab := cldf.NewMemoryAddressBook()

	deps := operation.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	input := seq.CurseSubjectsSeqInput{
		CCIPAddress: state.AptosChains[cfg.ChainSelector].CCIPAddress,
		Subjects:    cfg.Subjects,
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.CurseSubjectsSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute CurseSubjectsSequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{report.Output},
		"Curse subjects on Aptos RMN Remote",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	ds, err := shared.PopulateDataStore(ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// UncurseSubjects changeset: uncurses multiple subjects on Aptos RMN Remote.
var _ cldf.ChangeSetV2[config.RMNRemoteCurseInput] = UncurseSubjects{}

type UncurseSubjects struct{}

func (cs UncurseSubjects) VerifyPreconditions(env cldf.Environment, cfg config.RMNRemoteCurseInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	if (state.AptosChains[cfg.ChainSelector].CCIPAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("CCIP is not deployed on Aptos chain %d", cfg.ChainSelector))
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if len(cfg.Subjects) == 0 {
		errs = append(errs, errors.New("subjects list cannot be empty"))
	}
	return errors.Join(errs...)
}

func (cs UncurseSubjects) Apply(env cldf.Environment, cfg config.RMNRemoteCurseInput) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.BlockChains.AptosChains()[cfg.ChainSelector]
	ab := cldf.NewMemoryAddressBook()

	deps := operation.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	input := seq.UncurseSubjectsSeqInput{
		CCIPAddress: state.AptosChains[cfg.ChainSelector].CCIPAddress,
		Subjects:    cfg.Subjects,
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.UncurseSubjectsSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute UncurseSubjectsSequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{report.Output},
		"Uncurse subjects on Aptos RMN Remote",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	ds, err := shared.PopulateDataStore(ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
