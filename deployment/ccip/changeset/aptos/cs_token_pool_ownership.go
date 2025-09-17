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

var _ cldf.ChangeSetV2[config.TransferTokenPoolOwnershipInput] = TransferTokenPoolOwnership{}

type TransferTokenPoolOwnership struct{}

func (t TransferTokenPoolOwnership) VerifyPreconditions(env cldf.Environment, cfg config.TransferTokenPoolOwnershipInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	for i, transfer := range cfg.Transfers {
		if (transfer.TokenPoolAddress == aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("token pool address of transfer %d is empty", i))
		}
		if (transfer.To == aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("to address of transfer %d is empty", i))
		}
		switch transfer.TokenPoolType {
		case shared.AptosManagedTokenPoolType, shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		default:
			errs = append(errs, fmt.Errorf("token pool type %v of transfer %d is unsupported", transfer.TokenPoolType.String(), i))
		}
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if (state.AptosChains[cfg.ChainSelector].MCMSAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("MCMS is not deployed on Aptos chain %d", cfg.ChainSelector))
	}

	return errors.Join(errs...)
}

func (t TransferTokenPoolOwnership) Apply(env cldf.Environment, cfg config.TransferTokenPoolOwnershipInput) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	deps := operation.AptosDeps{
		AptosChain:       env.BlockChains.AptosChains()[cfg.ChainSelector],
		CCIPOnChainState: state,
	}
	input := seq.TransferTokenPoolOwnershipsSeqInput{
		Transfers: make([]seq.TokenPoolTransferInput, 0, len(cfg.Transfers)),
	}
	for _, transfer := range cfg.Transfers {
		input.Transfers = append(input.Transfers, seq.TokenPoolTransferInput{
			TokenPoolAddress: transfer.TokenPoolAddress,
			To:               transfer.To,
			TokenPoolType:    transfer.TokenPoolType,
		})
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.TransferTokenPoolOwnershipsSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute TransferTokenPoolOwnershipsSequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{report.Output},
		"Transfers ownership of one or multiple token pool instances to new addresses",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

var _ cldf.ChangeSetV2[config.AcceptTokenPoolOwnershipInput] = AcceptTokenPoolOwnership{}

type AcceptTokenPoolOwnership struct{}

func (a AcceptTokenPoolOwnership) VerifyPreconditions(env cldf.Environment, cfg config.AcceptTokenPoolOwnershipInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	for i, accept := range cfg.Accepts {
		if (accept.TokenPoolAddress == aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("token pool address of transfer %d is empty", i))
		}
		switch accept.TokenPoolType {
		case shared.AptosManagedTokenPoolType, shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		default:
			errs = append(errs, fmt.Errorf("token pool type %v of transfer %d is unsupported", accept.TokenPoolType.String(), i))
		}
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if (state.AptosChains[cfg.ChainSelector].MCMSAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("MCMS is not deployed on Aptos chain %d", cfg.ChainSelector))
	}

	return errors.Join(errs...)
}

func (a AcceptTokenPoolOwnership) Apply(env cldf.Environment, cfg config.AcceptTokenPoolOwnershipInput) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	deps := operation.AptosDeps{
		AptosChain:       env.BlockChains.AptosChains()[cfg.ChainSelector],
		CCIPOnChainState: state,
	}
	input := seq.AcceptTokenPoolOwnershipsSeqInput{
		Accepts: make([]seq.TokenPoolAcceptInput, 0, len(cfg.Accepts)),
	}
	for _, transfer := range cfg.Accepts {
		input.Accepts = append(input.Accepts, seq.TokenPoolAcceptInput{
			TokenPoolAddress: transfer.TokenPoolAddress,
			TokenPoolType:    transfer.TokenPoolType,
		})
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.AcceptTokenPoolOwnershipsSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute AcceptTokenPoolOwnershipsSequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{report.Output},
		"Accepts ownership of one or multiple token pool instances",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

var _ cldf.ChangeSetV2[config.ExecuteTokenPoolOwnershipTransferInput] = ExecuteTokenPoolOwnershipTransfer{}

type ExecuteTokenPoolOwnershipTransfer struct{}

func (e ExecuteTokenPoolOwnershipTransfer) VerifyPreconditions(env cldf.Environment, cfg config.ExecuteTokenPoolOwnershipTransferInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	for i, transfer := range cfg.Transfers {
		if (transfer.TokenPoolAddress == aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("token pool address of transfer %d is empty", i))
		}
		if (transfer.To == aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("to address of transfer %d is empty", i))
		}
		switch transfer.TokenPoolType {
		case shared.AptosManagedTokenPoolType, shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		default:
			errs = append(errs, fmt.Errorf("token pool type %v of transfer %d is unsupported", transfer.TokenPoolType.String(), i))
		}
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if (state.AptosChains[cfg.ChainSelector].MCMSAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("MCMS is not deployed on Aptos chain %d", cfg.ChainSelector))
	}

	return errors.Join(errs...)
}

func (e ExecuteTokenPoolOwnershipTransfer) Apply(env cldf.Environment, cfg config.ExecuteTokenPoolOwnershipTransferInput) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	deps := operation.AptosDeps{
		AptosChain:       env.BlockChains.AptosChains()[cfg.ChainSelector],
		CCIPOnChainState: state,
	}
	input := seq.ExecuteTokenPoolOwnershipTransfersSeqInput{
		Transfers: make([]seq.TokenPoolTransferInput, 0, len(cfg.Transfers)),
	}
	for _, transfer := range cfg.Transfers {
		input.Transfers = append(input.Transfers, seq.TokenPoolTransferInput{
			TokenPoolAddress: transfer.TokenPoolAddress,
			To:               transfer.To,
			TokenPoolType:    transfer.TokenPoolType,
		})
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.ExecuteTokenPoolOwnershipTransfersSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute ExecuteTokenPoolOwnershipTransfersSequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{report.Output},
		"Executed the pending ownership transfer of one or multiple token pool instances to new addresses",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
