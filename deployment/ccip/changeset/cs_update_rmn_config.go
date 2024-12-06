package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	mcmsWrappers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
)

type SetRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	RMNStaticConfig   rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig  rmn_home.RMNHomeDynamicConfig
	DigestToOverride  [32]byte
}

func (c SetRMNHomeCandidateConfig) Validate() error {
	if c.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}
	if len(c.RMNDynamicConfig.OffchainConfig) != 0 {
		return fmt.Errorf("RMNDynamicConfig.OffchainConfig must be empty")
	}
	if len(c.RMNStaticConfig.OffchainConfig) != 0 {
		return fmt.Errorf("RMNStaticConfig.OffchainConfig must be empty")
	}

	return nil
}

type PromoteRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
}

func (c PromoteRMNHomeCandidateConfig) Validate() error {
	if c.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}

	return nil
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

	prop, err := buildProposal(e, op, state, config.HomeChainSelector)

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

	prop, err := buildProposal(e, op, state, config.HomeChainSelector)

	if err != nil {
		lggr.Errorw("Failed to build proposal", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func buildTimelockPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*mcmsWrappers.RBACTimelock {
	timelocksPerChain := make(map[uint64]*mcmsWrappers.RBACTimelock)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = state.Chains[chain.Selector].Timelock
	}
	return timelocksPerChain
}

func buildTimelockAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]common.Address {
	timelocksPerChain := buildTimelockPerChain(e, state)
	timelockAddressPerChain := make(map[uint64]common.Address)
	for chain, timelock := range timelocksPerChain {
		timelockAddressPerChain[chain] = timelock.Address()
	}
	return timelockAddressPerChain
}

func buildProposerPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*gethwrappers.ManyChainMultiSig {
	proposerPerChain := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range e.Chains {
		proposerPerChain[chain.Selector] = state.Chains[chain.Selector].ProposerMcm
	}
	return proposerPerChain
}

func buildRemoteRemotePerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*rmn_remote.RMNRemote {
	timelocksPerChain := make(map[uint64]*rmn_remote.RMNRemote)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = state.Chains[chain.Selector].RMNRemote
	}
	return timelocksPerChain
}

func buildProposal(e deployment.Environment, op mcms.Operation, state CCIPOnChainState, homeChainSelector uint64) (*timelock.MCMSWithTimelockProposal, error) {
	batches := []timelock.BatchChainOperation{
		{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSelector),
			Batch:           []mcms.Operation{op},
		},
	}

	timelocksPerChain := buildTimelockAddressPerChain(e, state)

	proposerMCMSes := buildProposerPerChain(e, state)

	return proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to promote candidate config",
		0,
	)
}

type SetRMNRemoteConfig struct {
	HomeChainSelector uint64
	Signers           []rmn_remote.RMNRemoteSigner
	F                 uint64
}

func (c SetRMNRemoteConfig) Validate() error {
	if c.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}

	return nil
}

func NewSetRMNRemoteConfigChangeset(e deployment.Environment, config SetRMNRemoteConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	lggr := e.Logger
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	rmnHome := state.Chains[config.HomeChainSelector].RMNHome

	activeConfig, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		lggr.Errorw("Failed to get RMNHome active digest", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	rmnRemotePerChain := buildRemoteRemotePerChain(e, state)
	batches := make([]timelock.BatchChainOperation, 0)
	for chain, remote := range rmnRemotePerChain {
		if remote == nil {
			continue
		}

		tx, err := remote.SetConfig(deployment.SimTransactOpts(), rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeConfig,
			Signers:                     config.Signers,
			F:                           config.F,
		})

		if err != nil {
			lggr.Errorw("Failed to build call data to set RMNRemote config", "err", err)
			return deployment.ChangesetOutput{}, err
		}

		op := mcms.Operation{
			To:    remote.Address(),
			Data:  tx.Data(),
			Value: big.NewInt(0),
		}

		batch := timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain),
			Batch:           []mcms.Operation{op},
		}

		batches = append(batches, batch)
	}

	timelocksPerChain := buildTimelockAddressPerChain(e, state)

	proposerMCMSes := buildProposerPerChain(e, state)

	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to promote candidate config",
		0,
	)

	if err != nil {
		lggr.Errorw("Failed to build proposal", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}
