package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
)

type SetRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	RMNStaticConfig   rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig  rmn_home.RMNHomeDynamicConfig
	DigestToOverride  [32]byte
}

type PromoteRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
}

func NewSetRMNHomeCandidateConfigChangeset(e deployment.Environment, config SetRMNHomeCandidateConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	lggr := e.Logger
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	rmnHome := state.Chains[config.HomeChainSelector].RMNHome

	setCandidateTx, err := rmnHome.SetCandidate(deployment.SimTransactOpts(), config.RMNStaticConfig, config.RMNDynamicConfig, config.DigestToOverride)
	if err != nil {
		lggr.Errorw("Failed to build call data to set RMNHome candidate digest", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	op := mcms.Operation{
		To:    rmnHome.Address(),
		Data:  setCandidateTx.Data(),
		Value: big.NewInt(0),
	}

	prop, err := buildProposal(op, state, config.HomeChainSelector)

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func NewPromoteCandidateConfigChangeset(e deployment.Environment, config PromoteRMNHomeCandidateConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	lggr := e.Logger
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	rmnHome := state.Chains[config.HomeChainSelector].RMNHome

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		lggr.Errorw("Failed to get RMNHome candidate digest", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	currentActiveDigest, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		lggr.Errorw("Failed to get RMNHome active digest", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	promoteCandidateTx, err := rmnHome.PromoteCandidateAndRevokeActive(deployment.SimTransactOpts(), currentCandidateDigest, currentActiveDigest)
	if err != nil {
		lggr.Errorw("Failed to get call data to promote RMNHome candidate digest", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	op := mcms.Operation{
		To:    rmnHome.Address(),
		Data:  promoteCandidateTx.Data(),
		Value: big.NewInt(0),
	}

	prop, err := buildProposal(op, state, config.HomeChainSelector)

	if err != nil {
		lggr.Errorw("Failed to build proposal", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func buildProposal(op mcms.Operation, state CCIPOnChainState, homeChainSelector uint64) (*timelock.MCMSWithTimelockProposal, error) {
	batches := []timelock.BatchChainOperation{
		{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSelector),
			Batch:           []mcms.Operation{op},
		},
	}

	timelocksPerChain := map[uint64]common.Address{
		homeChainSelector: state.Chains[homeChainSelector].Timelock.Address(),
	}

	proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
		homeChainSelector: state.Chains[homeChainSelector].ProposerMcm,
	}

	return proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to promote candidate config",
		0,
	)
}
