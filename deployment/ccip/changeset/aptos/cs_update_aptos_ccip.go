package aptos

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

var _ cldf.ChangeSetV2[config.UpdateAptosChainConfig] = UpdateAptosChain{}

// DeployAptosChain deploys Aptos chain packages and modules
type UpdateAptosChain struct{}

func (cs UpdateAptosChain) VerifyPreconditions(env deployment.Environment, config config.UpdateAptosChainConfig) error {
	// TODO validate if required packages and modules are already deployed
	return nil
}

func (cs UpdateAptosChain) Apply(env deployment.Environment, cfg config.UpdateAptosChainConfig) (cldf.ChangesetOutput, error) {
	timeLockProposals := []mcms.TimelockProposal{}
	mcmsOperations := []types.BatchOperation{}
	seqReports := make([]operations.Report[any, any], 0)

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	deps := operation.AptosDeps{
		AptosChain:       env.AptosChains[cfg.ChainSelector],
		OnChainState:     state.AptosChains[cfg.ChainSelector],
		CCIPOnChainState: state,
	}

	// Execute the sequence
	updateSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.UpdateCCIPSequence, deps, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, updateSeqReport.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, updateSeqReport.Output...)

	// Generate MCMS proposals
	proposal, err := utils.GenerateProposal(
		deps.AptosChain.Client,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		deps.AptosChain.Selector,
		mcmsOperations,
		"Update chain contracts on Aptos chain",
		cfg.MCMSTimelockConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}
	timeLockProposals = append(timeLockProposals, *proposal)

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               seqReports,
	}, nil
}
