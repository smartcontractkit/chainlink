package aptos

import (
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

var _ deployment.ChangeSetV2[v1_6.SetOCR3OffRampConfig] = SetOCR3Offramp{}

// SetOCR3Offramp updates OCR3 Offramp configurations
type SetOCR3Offramp struct{}

func (cs SetOCR3Offramp) VerifyPreconditions(env deployment.Environment, config v1_6.SetOCR3OffRampConfig) error {
	for _, remoteSel := range config.RemoteChainSels {
		chainFamily, _ := chain_selectors.GetSelectorFamily(remoteSel)
		if chainFamily != chain_selectors.FamilyAptos {
			return fmt.Errorf("chain %d is not an Aptos chain", remoteSel)
		}
	}
	return nil
}

func (cs SetOCR3Offramp) Apply(env deployment.Environment, config v1_6.SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	seqReports := make([]operations.Report[any, any], 0)
	timeLockProposals := []mcms.TimelockProposal{}

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	for _, remoteSelector := range config.RemoteChainSels {
		deps := operation.AptosDeps{
			AptosChain:       env.AptosChains[remoteSelector],
			OnChainState:     state.AptosChains[remoteSelector],
			CCIPOnChainState: state,
		}
		in := seq.SetOCR3OfframpSeqInput{
			HomeChainSelector: config.HomeChainSel,
			ChainSelector:     remoteSelector,
		}
		setOCR3SeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.SetOCR3OfframpSequence, deps, in)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, setOCR3SeqReport.ExecutionReports...)

		// Generate MCMS proposals
		proposal, err := utils.GenerateProposal(
			deps.AptosChain.Client,
			state.AptosChains[remoteSelector].MCMSAddress,
			deps.AptosChain.Selector,
			[]types.BatchOperation{setOCR3SeqReport.Output},
			"Set OCR3 Configs",
			aptosmcms.TimelockRoleProposer,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", remoteSelector, err)
		}
		timeLockProposals = append(timeLockProposals, *proposal)
	}

	return deployment.ChangesetOutput{
		MCMSTimelockProposals: timeLockProposals,
		Reports:               seqReports,
	}, nil
}
