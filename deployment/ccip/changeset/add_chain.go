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

type AddDonAndSetCandidateChangesetConfig struct {
	HomeChainSelector uint64
	FeedChainSelector uint64
	NewChainSelector  uint64
	PluginType        types.PluginType
	TokenConfig       TokenConfig
	Nodes             deployment.Nodes
	OCRSecrets        deployment.OCRSecrets
}

func (a AddDonAndSetCandidateChangesetConfig) Validate() error {
	if a.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}
	if a.FeedChainSelector == 0 {
		return fmt.Errorf("FeedChainSelector must be set")
	}
	if a.NewChainSelector == 0 {
		return fmt.Errorf("NewChainSelector must be set")
	}
	if a.PluginType != types.PluginTypeCCIPCommit &&
		a.PluginType != types.PluginTypeCCIPExec {
		return fmt.Errorf("PluginType must be set to either CCIPCommit or CCIPExec")
	}
	// TODO: validate token config
	if len(a.Nodes) == 0 {
		return fmt.Errorf("Nodes must be set")
	}
	return nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the commit plugin config as candidateConfig for the don.
func AddDonAndSetCandidateChangeset(
	e deployment.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ccipOCRParams := DefaultOCRParams(
		cfg.FeedChainSelector,
		cfg.TokenConfig.GetTokenInfo(e.Logger,
			state.Chains[cfg.NewChainSelector].LinkToken,
			state.Chains[cfg.NewChainSelector].Weth9,
		),
	)
	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		cfg.OCRSecrets,
		state.Chains[cfg.NewChainSelector].OffRamp,
		e.Chains[cfg.NewChainSelector],
		cfg.Nodes.NonBootstraps(),
		state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
		ccipOCRParams.OCRParameters,
		ccipOCRParams.CommitOffChainConfig,
		ccipOCRParams.ExecuteOffChainConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	latestDon, err := internal.LatestCCIPDON(state.Chains[cfg.HomeChainSelector].CapabilityRegistry)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	commitConfig, ok := newDONArgs[cfg.PluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := NewDonWithCandidateOp(
		donID, commitConfig,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		cfg.Nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var (
		timelocksPerChain = map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		}
		proposerMCMSes = map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		}
	)
	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
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
