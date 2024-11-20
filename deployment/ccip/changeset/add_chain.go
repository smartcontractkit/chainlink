package changeset

import (
	"fmt"
	"math/big"

	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

// NewChainInboundChangeset generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundChangeset(
	e deployment.Environment,
	state ccipdeployment.CCIPOnChainState,
	homeChainSel uint64,
	newChainSel uint64,
	sources []uint64,
) (deployment.ChangesetOutput, error) {
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	for _, source := range sources {
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: newChainSel,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			deployment.SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: newChainSel,
					DestChainConfig:   ccipdeployment.DefaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(source),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enableFeeQuoterDest.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}

	addChainOp, err := ccipdeployment.ApplyChainConfigUpdatesOp(e, state, homeChainSel, []uint64{newChainSel})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch: []mcms.Operation{
			addChainOp,
		},
	})

	prop, err := ccipdeployment.BuildProposalFromBatches(state, batches, "proposal to set new chains", 0)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the commit plugin config as candidateConfig for the don.
func AddDonAndSetCandidateChangeset(
	state ccipdeployment.CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig ccipdeployment.TokenConfig,
	pluginType types.PluginType,
) (deployment.ChangesetOutput, error) {
	newDONArgs, err := ccipdeployment.BuildOCR3ConfigForCCIPHome(
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
	latestDon, err := ccipdeployment.LatestCCIPDON(state.Chains[homeChainSel].CapabilityRegistry)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	commitConfig, ok := newDONArgs[pluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := ccipdeployment.NewDonWithCandidateOp(
		donID, commitConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	prop, err := ccipdeployment.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           []mcms.Operation{addDonOp},
	}}, "setCandidate for commit and AddDon on new Chain", 0)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}
