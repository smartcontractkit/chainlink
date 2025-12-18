package aptos

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.DynamicConfig] = DynamicCS{}

// DynamicCS enables dynamic execution of multiple Aptos operations at runtime
// without requiring dedicated changesets. It allows for flexible, configurable
// operation sequences to be applied based on runtime input, streamlining the
// process of managing Aptos chain changes.
type DynamicCS struct{}

func (cs DynamicCS) VerifyPreconditions(env cldf.Environment, cfg config.DynamicConfig) error {
	if len(cfg.Defs) != len(cfg.Inputs) {
		return fmt.Errorf("precondition failed: cfg.Defs and cfg.Inputs must have matching lengths (got %d and %d)", len(cfg.Defs), len(cfg.Inputs))
	}
	if cfg.MCMSConfig == nil {
		return errors.New("precondition failed: cfg.MCMSConfig must not be nil")
	}
	return nil
}

func (cs DynamicCS) Apply(env cldf.Environment, cfg config.DynamicConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.BlockChains.AptosChains()[cfg.ChainSelector]
	ab := cldf.NewMemoryAddressBook()

	deps := dependency.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	in := seq.DynamicSeqInput{
		Defs:          cfg.Defs,
		Inputs:        cfg.Inputs,
		ChainSelector: cfg.ChainSelector,
	}
	result, err := operations.ExecuteSequence(env.OperationsBundle, seq.DynamicSequence, deps, in)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
	}

	// Generate Aptos MCMS proposals
	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]types.BatchOperation{result.Output},
		cfg.Description,
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               result.ExecutionReports,
	}, nil
}
