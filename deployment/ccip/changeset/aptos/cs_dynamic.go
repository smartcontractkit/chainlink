package aptos

import (
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.DynamicConfig] = DynamicCS{}

// TODO: add description
// DynamicCS ...
type DynamicCS struct{}

func (cs DynamicCS) VerifyPreconditions(env cldf.Environment, cfg config.DynamicConfig) error {
	// TODO: implement precondition checks
	return nil
}

func (cs DynamicCS) Apply(env cldf.Environment, cfg config.DynamicConfig) (cldf.ChangesetOutput, error) {
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
