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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.MintTokenInput] = MintToken{}

type MintToken struct{}

func (m MintToken) VerifyPreconditions(env cldf.Environment, cfg config.MintTokenInput) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	var errs []error
	if _, ok := state.SupportedChains()[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	if (cfg.TokenCodeObjectAddress == aptos.AccountAddress{}) {
		errs = append(errs, errors.New("managed token object address is empty"))
	}
	if cfg.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if (state.AptosChains[cfg.ChainSelector].MCMSAddress == aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("MCMS is not deployed on Aptos chain %d", cfg.ChainSelector))
	}
	if (cfg.To == aptos.AccountAddress{}) {
		errs = append(errs, errors.New("to address is empty"))
	}
	if cfg.Amount == 0 {
		errs = append(errs, errors.New("token amount is 0"))
	}

	return errors.Join(errs...)
}

func (m MintToken) Apply(env cldf.Environment, cfg config.MintTokenInput) (cldf.ChangesetOutput, error) {
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

	input := operation.MintTokensInput{
		TokenCodeObjectAddress: cfg.TokenCodeObjectAddress,
		To:                     cfg.To,
		Amount:                 cfg.Amount,
	}
	report, err := operations.ExecuteOperation(env.OperationsBundle, operation.MintTokensOp, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute MintTokensOp: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{
			{
				ChainSelector: mcmstypes.ChainSelector(cfg.ChainSelector),
				Transactions:  []mcmstypes.Transaction{report.Output},
			},
		},
		"Mint tokens",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
