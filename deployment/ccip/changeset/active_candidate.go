package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ deployment.ChangeSet[PromoteAllCandidatesChangesetConfig]  = PromoteAllCandidatesChangeset
	_ deployment.ChangeSet[AddDonAndSetCandidateChangesetConfig] = SetCandidatePluginChangeset
)

type PromoteAllCandidatesChangesetConfig struct {
	HomeChainSelector uint64
	NewChainSelector  uint64
	Nodes             deployment.Nodes
}

func (p PromoteAllCandidatesChangesetConfig) Validate() error {
	if p.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}
	if p.NewChainSelector == 0 {
		return fmt.Errorf("NewChainSelector must be set")
	}
	if len(p.Nodes) == 0 {
		return fmt.Errorf("Nodes must be set")
	}
	return nil
}

// PromoteAllCandidatesChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
func PromoteAllCandidatesChangeset(
	e deployment.Environment,
	cfg PromoteAllCandidatesChangesetConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	promoteCandidateOps, err := PromoteAllCandidatesForChainOps(
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		cfg.NewChainSelector,
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
			Batch:           promoteCandidateOps,
		}},
		"promoteCandidate for commit and execution",
		0, // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}

// SetCandidatePluginChangeset calls setCandidate on the CCIPHome for setting up OCR3 exec Plugin config for the new chain.
func SetCandidatePluginChangeset(
	e deployment.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	ccipOCRParams := DefaultOCRParams(
		cfg.FeedChainSelector,
		cfg.TokenConfig.GetTokenInfo(e.Logger, state.Chains[cfg.NewChainSelector].LinkToken, state.Chains[cfg.NewChainSelector].Weth9),
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

	execConfig, ok := newDONArgs[cfg.PluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	setCandidateMCMSOps, err := SetCandidateOnExistingDon(
		execConfig,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		cfg.NewChainSelector,
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
			Batch:           setCandidateMCMSOps,
		}},
		"SetCandidate for execution",
		0, // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}
