package changeset

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var (
	_ deployment.ChangeSet[PromoteAllCandidatesChangesetConfig]  = PromoteAllCandidatesChangeset
	_ deployment.ChangeSet[AddDonAndSetCandidateChangesetConfig] = SetCandidatePluginChangeset
)

type PromoteAllCandidatesChangesetConfig struct {
	HomeChainSelector uint64
	// DONChainSelector is the chain selector of the DON that we want to promote the candidate config of.
	// Note that each (chain, ccip capability version) pair has a unique DON ID.
	DONChainSelector uint64
	NodeIDs          []string
	MCMS             *MCMSConfig
}

func (p PromoteAllCandidatesChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (deployment.Nodes, error) {
	if err := deployment.IsValidChainSelector(p.HomeChainSelector); err != nil {
		return nil, fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(p.DONChainSelector); err != nil {
		return nil, fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(p.NodeIDs) == 0 {
		return nil, fmt.Errorf("NodeIDs must be set")
	}
	if state.Chains[p.HomeChainSelector].CCIPHome == nil {
		return nil, fmt.Errorf("CCIPHome contract does not exist")
	}
	if state.Chains[p.HomeChainSelector].CapabilityRegistry == nil {
		return nil, fmt.Errorf("CapabilityRegistry contract does not exist")
	}

	nodes, err := deployment.NodeInfo(p.NodeIDs, e.Offchain)
	if err != nil {
		return nil, fmt.Errorf("fetch node info: %w", err)
	}

	donID, err := internal.DonIDForChain(
		state.Chains[p.HomeChainSelector].CapabilityRegistry,
		state.Chains[p.HomeChainSelector].CCIPHome,
		p.DONChainSelector,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return nil, fmt.Errorf("don doesn't exist in CR for chain %d", p.DONChainSelector)
	}

	// Check that candidate digest and active digest are not both zero - this is enforced onchain.
	commitConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
		Context: context.Background(),
	}, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("fetching commit configs from cciphome: %w", err)
	}

	execConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
		Context: context.Background(),
	}, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return nil, fmt.Errorf("fetching exec configs from cciphome: %w", err)
	}

	if commitConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
		commitConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return nil, fmt.Errorf("commit active and candidate config digests are both zero")
	}

	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
		execConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return nil, fmt.Errorf("exec active and candidate config digests are both zero")
	}

	return nodes, nil
}

// PromoteAllCandidatesChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
func PromoteAllCandidatesChangeset(
	e deployment.Environment,
	cfg PromoteAllCandidatesChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	nodes, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	homeChain := e.Chains[cfg.HomeChainSelector]

	promoteCandidateOps, err := promoteAllCandidatesForChainOps(
		homeChain,
		txOpts,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		cfg.DONChainSelector,
		nodes.NonBootstraps(),
		cfg.MCMS != nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("generating promote candidate ops: %w", err)
	}

	// Disabled MCMS means that we already executed the txes, so just return early w/out the proposals.
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	prop, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           promoteCandidateOps,
		}},
		"promoteCandidate for commit and execution",
		cfg.MCMS.MinDelay,
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
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	nodes, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		e.OCRSecrets,
		state.Chains[cfg.NewChainSelector].OffRamp,
		e.Chains[cfg.NewChainSelector],
		nodes.NonBootstraps(),
		state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
		cfg.CCIPOCRParams.OCRParameters,
		cfg.CCIPOCRParams.CommitOffChainConfig,
		cfg.CCIPOCRParams.ExecuteOffChainConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	config, ok := newDONArgs[cfg.PluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing %s plugin in ocr3Configs", cfg.PluginType.String())
	}

	setCandidateMCMSOps, err := setCandidateOnExistingDon(
		config,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		cfg.NewChainSelector,
		nodes.NonBootstraps(),
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
		fmt.Sprintf("SetCandidate for %s plugin", cfg.PluginType.String()),
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

// setCandidateOnExistingDon calls setCandidate on CCIPHome contract through the UpdateDON call on CapReg contract
// This proposes to set up OCR3 config for the provided plugin for the DON
func setCandidateOnExistingDon(
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSelector uint64,
	nodes deployment.Nodes,
) ([]mcms.Operation, error) {
	// fetch DON ID for the chain
	donID, err := internal.DonIDForChain(capReg, ccipHome, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
	}

	fmt.Printf("donID: %d", donID)
	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	// set candidate call
	updateDonTx, err := capReg.UpdateDON(
		deployment.SimTransactOpts(),
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ exec config: %w", err)
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}

// promoteCandidateOp will create the MCMS Operation for `promoteCandidateAndRevokeActive` directed towards the capabilityRegistry
func promoteCandidateOp(
	homeChain deployment.Chain,
	txOpts *bind.TransactOpts,
	donID uint32,
	pluginType uint8,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	mcmsEnabled bool,
) (mcms.Operation, error) {
	allConfigs, err := ccipHome.GetAllConfigs(nil, donID, pluginType)
	if err != nil {
		return mcms.Operation{}, err
	}

	encodedPromotionCall, err := internal.CCIPHomeABI.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		pluginType,
		allConfigs.CandidateConfig.ConfigDigest,
		allConfigs.ActiveConfig.ConfigDigest,
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack promotion call: %w", err)
	}

	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("error creating updateDon op for donID(%d) and plugin type (%d): %w", donID, pluginType, err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return mcms.Operation{}, fmt.Errorf("error confirming updateDon call for donID(%d) and plugin type (%d): %w", donID, pluginType, err)
		}
	}

	return mcms.Operation{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

// promoteAllCandidatesForChainOps promotes the candidate commit and exec configs to active by calling promoteCandidateAndRevokeActive on CCIPHome through the UpdateDON call on CapReg contract
func promoteAllCandidatesForChainOps(
	homeChain deployment.Chain,
	txOpts *bind.TransactOpts,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSelector uint64,
	nodes deployment.Nodes,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	// fetch DON ID for the chain
	donID, err := internal.DonIDForChain(capReg, ccipHome, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
	}

	var mcmsOps []mcms.Operation
	updateCommitOp, err := promoteCandidateOp(homeChain, txOpts, donID, uint8(cctypes.PluginTypeCCIPCommit), capReg, ccipHome, nodes, mcmsEnabled)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateCommitOp)

	updateExecOp, err := promoteCandidateOp(homeChain, txOpts, donID, uint8(cctypes.PluginTypeCCIPExec), capReg, ccipHome, nodes, mcmsEnabled)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateExecOp)

	return mcmsOps, nil
}
