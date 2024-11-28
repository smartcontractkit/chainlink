package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

var _ deployment.ChangeSet[ChainInboundChangesetConfig] = NewChainInboundChangeset

type ChainInboundChangesetConfig struct {
	HomeChainSelector    uint64
	NewChainSelector     uint64
	SourceChainSelectors []uint64
}

func (c ChainInboundChangesetConfig) Validate() error {
	if c.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}
	if c.NewChainSelector == 0 {
		return fmt.Errorf("NewChainSelector must be set")
	}
	if len(c.SourceChainSelectors) == 0 {
		return fmt.Errorf("SourceChainSelectors must be set")
	}
	return nil
}

// NewChainInboundChangeset generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundChangeset(
	e deployment.Environment,
	cfg ChainInboundChangesetConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	for _, source := range cfg.SourceChainSelectors {
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: cfg.NewChainSelector,
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
					DestChainSelector: cfg.NewChainSelector,
					DestChainConfig:   DefaultFeeQuoterDestChainConfig(),
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

	addChainOp, err := ApplyChainConfigUpdatesOp(e, state, cfg.HomeChainSelector, []uint64{cfg.NewChainSelector})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
		Batch: []mcms.Operation{
			addChainOp,
		},
	})

	var (
		timelocksPerChain = make(map[uint64]common.Address)
		proposerMCMSes    = make(map[uint64]*gethwrappers.ManyChainMultiSig)
	)
	for _, chain := range append(cfg.SourceChainSelectors, cfg.HomeChainSelector) {
		timelocksPerChain[chain] = state.Chains[chain].Timelock.Address()
		proposerMCMSes[chain] = state.Chains[chain].ProposerMcm
	}
	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to set new chains",
		0,
	)
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
	state CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig TokenConfig,
	pluginType types.PluginType,
) (deployment.ChangesetOutput, error) {
	ccipOCRParams := DefaultOCRParams(
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
	)
	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		nodes.NonBootstraps(),
		state.Chains[homeChainSel].RMNHome.Address(),
		ccipOCRParams.OCRParameters,
		ccipOCRParams.CommitOffChainConfig,
		ccipOCRParams.ExecuteOffChainConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	latestDon, err := internal.LatestCCIPDON(state.Chains[homeChainSel].CapabilityRegistry)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	commitConfig, ok := newDONArgs[pluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := NewDonWithCandidateOp(
		donID, commitConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var (
		timelocksPerChain = map[uint64]common.Address{
			homeChainSel: state.Chains[homeChainSel].Timelock.Address(),
		}
		proposerMCMSes = map[uint64]*gethwrappers.ManyChainMultiSig{
			homeChainSel: state.Chains[homeChainSel].ProposerMcm,
		}
	)
	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
			Batch:           []mcms.Operation{addDonOp},
		}},
		"setCandidate for commit and AddDon on new Chain",
		0, // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}
