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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	operation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.MigrateOnRampDestChainConfigsToV2Config] = MigrateOnRampDestChainConfigsToV2{}

type MigrateOnRampDestChainConfigsToV2 struct{}

func (cs MigrateOnRampDestChainConfigsToV2) VerifyPreconditions(env cldf.Environment, cfg config.MigrateOnRampDestChainConfigsToV2Config) error {
	var errs []error

	if cfg.MCMS == nil {
		errs = append(errs, errors.New("MCMS config is required"))
	}
	if len(cfg.DestChainSelectors) == 0 {
		errs = append(errs, errors.New("DestChainSelectors must not be empty"))
	}
	if len(cfg.DestChainSelectors) != len(cfg.RouterModuleAddresses) {
		errs = append(errs, fmt.Errorf("DestChainSelectors length (%d) must match RouterModuleAddresses length (%d)",
			len(cfg.DestChainSelectors), len(cfg.RouterModuleAddresses)))
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load onchain state: %w", err))
		return errors.Join(errs...)
	}

	supportedChains := state.SupportedChains()
	if _, ok := supportedChains[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}

	ccipAddress := state.AptosChains[cfg.ChainSelector].CCIPAddress
	if ccipAddress == (aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("CCIP package is not deployed on Aptos chain %d", cfg.ChainSelector))
	}

	client := env.BlockChains.AptosChains()[cfg.ChainSelector].Client
	hasOnramp, err := utils.IsModuleDeployed(client, ccipAddress, "onramp")
	if err != nil || !hasOnramp {
		errs = append(errs, fmt.Errorf("onRamp module is not deployed on Aptos chain %d: %w", cfg.ChainSelector, err))
	}

	return errors.Join(errs...)
}

func (cs MigrateOnRampDestChainConfigsToV2) Apply(env cldf.Environment, cfg config.MigrateOnRampDestChainConfigsToV2Config) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deps := dependency.AptosDeps{
		AptosChain:       env.BlockChains.AptosChains()[cfg.ChainSelector],
		CCIPOnChainState: state,
	}

	seqInput := operation.MigrateOnRampDestChainConfigsToV2Input{
		DestChainSelectors:    cfg.DestChainSelectors,
		RouterModuleAddresses: cfg.RouterModuleAddresses,
	}

	seqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.MigrateOnRampDestChainConfigsToV2Sequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute migrate sequence: %w", err)
	}

	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{seqReport.Output},
		"Migrate OnRamp destination chain configs to V2",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               seqReport.ExecutionReports,
	}, nil
}
