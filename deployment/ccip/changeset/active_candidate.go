package changeset

import (
	"fmt"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// PromoteAllCandidatesChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
func PromoteAllCandidatesChangeset(
	state ccdeploy.CCIPOnChainState,
	homeChainSel, newChainSel uint64,
	nodes deployment.Nodes,
) (deployment.ChangesetOutput, error) {
	promoteCandidateOps, err := ccdeploy.PromoteAllCandidatesForChainOps(
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	prop, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           promoteCandidateOps,
	}}, "promoteCandidate for commit and execution", 0)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}

// SetCandidateExecPluginProposal calls setCandidate on the CCIPHome for setting up OCR3 exec Plugin config for the new chain.
func SetCandidatePluginChangeset(
	state ccdeploy.CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig ccdeploy.TokenConfig,
	pluginType cctypes.PluginType,
) (deployment.ChangesetOutput, error) {
	newDONArgs, err := ccdeploy.BuildOCR3ConfigForCCIPHome(
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
		nodes.NonBootstraps(),
		state.Chains[homeChainSel].RMNHome.Address(),
		nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	execConfig, ok := newDONArgs[pluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	setCandidateMCMSOps, err := ccdeploy.SetCandidateOnExistingDon(
		execConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	prop, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           setCandidateMCMSOps,
	}}, "SetCandidate for execution", 0)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil

}
