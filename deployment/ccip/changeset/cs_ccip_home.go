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
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var (
	_ deployment.ChangeSet[AddDonAndSetCandidateChangesetConfig] = AddDonAndSetCandidateChangeset
	_ deployment.ChangeSet[PromoteAllCandidatesChangesetConfig]  = PromoteAllCandidatesChangeset
	_ deployment.ChangeSet[SetCandidateChangesetConfig]          = SetCandidateChangeset
	_ deployment.ChangeSet[RevokeCandidateChangesetConfig]       = RevokeCandidateChangeset
)

type PromoteAllCandidatesChangesetConfig struct {
	HomeChainSelector uint64

	// DONChainSelector is the chain selector of the DON that we want to promote the candidate config of.
	// Note that each (chain, ccip capability version) pair has a unique DON ID.
	DONChainSelector uint64

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (p PromoteAllCandidatesChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (donID uint32, err error) {
	if err := deployment.IsValidChainSelector(p.HomeChainSelector); err != nil {
		return 0, fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(p.DONChainSelector); err != nil {
		return 0, fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(e.NodeIDs) == 0 {
		return 0, fmt.Errorf("NodeIDs must be set")
	}
	if state.Chains[p.HomeChainSelector].CCIPHome == nil {
		return 0, fmt.Errorf("CCIPHome contract does not exist")
	}
	if state.Chains[p.HomeChainSelector].CapabilityRegistry == nil {
		return 0, fmt.Errorf("CapabilityRegistry contract does not exist")
	}
	if state.Chains[p.DONChainSelector].OffRamp == nil {
		// should not be possible, but a defensive check.
		return 0, fmt.Errorf("OffRamp contract does not exist")
	}

	donID, err = internal.DonIDForChain(
		state.Chains[p.HomeChainSelector].CapabilityRegistry,
		state.Chains[p.HomeChainSelector].CCIPHome,
		p.DONChainSelector,
	)
	if err != nil {
		return 0, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("don doesn't exist in CR for chain %d", p.DONChainSelector)
	}

	// Check that candidate digest and active digest are not both zero - this is enforced onchain.
	commitConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
		Context: context.Background(),
	}, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return 0, fmt.Errorf("fetching commit configs from cciphome: %w", err)
	}

	execConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
		Context: context.Background(),
	}, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return 0, fmt.Errorf("fetching exec configs from cciphome: %w", err)
	}

	if commitConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
		commitConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return 0, fmt.Errorf("commit active and candidate config digests are both zero")
	}

	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
		execConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return 0, fmt.Errorf("exec active and candidate config digests are both zero")
	}

	return donID, nil
}

// PromoteAllCandidatesChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// Note that a DON must exist prior to being able to use this changeset effectively,
// i.e AddDonAndSetCandidateChangeset must be called first.
// This can also be used to promote a 0x0 candidate config to be the active, effectively shutting down the DON.
// At that point you can call the RemoveDON changeset to remove the DON entirely from the capability registry.
func PromoteAllCandidatesChangeset(
	e deployment.Environment,
	cfg PromoteAllCandidatesChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	donID, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("fetch node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	homeChain := e.Chains[cfg.HomeChainSelector]

	promoteCandidateOps, err := promoteAllCandidatesForChainOps(
		txOpts,
		homeChain,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		nodes.NonBootstraps(),
		donID,
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

// SetCandidateConfigBase is a common base config struct for AddDonAndSetCandidateChangesetConfig and SetCandidateChangesetConfig.
// This is extracted to deduplicate most of the validation logic.
// Remaining validation logic is done in the specific config structs that inherit from this.
type SetCandidateConfigBase struct {
	HomeChainSelector uint64
	FeedChainSelector uint64

	// DONChainSelector is the chain selector of the chain where the DON will be added.
	DONChainSelector uint64

	PluginType types.PluginType
	// Note that the PluginType field is used to determine which field in CCIPOCRParams is used.
	CCIPOCRParams CCIPOCRParams

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (s SetCandidateConfigBase) Validate(e deployment.Environment, state CCIPOnChainState) error {
	if err := deployment.IsValidChainSelector(s.HomeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(s.FeedChainSelector); err != nil {
		return fmt.Errorf("feed chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(s.DONChainSelector); err != nil {
		return fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(e.NodeIDs) == 0 {
		return fmt.Errorf("nodeIDs must be set")
	}
	if state.Chains[s.HomeChainSelector].CCIPHome == nil {
		return fmt.Errorf("CCIPHome contract does not exist")
	}
	if state.Chains[s.HomeChainSelector].CapabilityRegistry == nil {
		return fmt.Errorf("CapabilityRegistry contract does not exist")
	}
	if state.Chains[s.DONChainSelector].OffRamp == nil {
		// should not be possible, but a defensive check.
		return fmt.Errorf("OffRamp contract does not exist on don chain selector %d", s.DONChainSelector)
	}
	if s.PluginType != types.PluginTypeCCIPCommit &&
		s.PluginType != types.PluginTypeCCIPExec {
		return fmt.Errorf("PluginType must be set to either CCIPCommit or CCIPExec")
	}

	// no donID check since this config is used for both adding a new DON and updating an existing one.
	// see AddDonAndSetCandidateChangesetConfig.Validate and SetCandidateChangesetConfig.Validate
	// for these checks.

	// check that chain config is set up for the new chain
	chainConfig, err := state.Chains[s.HomeChainSelector].CCIPHome.GetChainConfig(nil, s.DONChainSelector)
	if err != nil {
		return fmt.Errorf("get all chain configs: %w", err)
	}

	// FChain should never be zero if a chain config is set in CCIPHome
	if chainConfig.FChain == 0 {
		return fmt.Errorf("chain config not set up for new chain %d", s.DONChainSelector)
	}

	err = s.CCIPOCRParams.Validate()
	if err != nil {
		return fmt.Errorf("invalid ccip ocr params: %w", err)
	}

	// TODO: validate token config in the commit config, if commit is the plugin.
	// TODO: validate gas config in the chain config in cciphome for this DONChainSelector.

	if e.OCRSecrets.IsEmpty() {
		return fmt.Errorf("OCR secrets must be set")
	}

	return nil
}

// AddDonAndSetCandidateChangesetConfig is a separate config struct
// because the validation is slightly different from SetCandidateChangesetConfig.
// In particular, we check to make sure we don't already have a DON for the chain.
type AddDonAndSetCandidateChangesetConfig struct {
	SetCandidateConfigBase
}

func (a AddDonAndSetCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	err := a.SetCandidateConfigBase.Validate(e, state)
	if err != nil {
		return err
	}

	// check if a DON already exists for this chain
	donID, err := internal.DonIDForChain(
		state.Chains[a.HomeChainSelector].CapabilityRegistry,
		state.Chains[a.HomeChainSelector].CCIPHome,
		a.DONChainSelector,
	)
	if err != nil {
		return fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID != 0 {
		return fmt.Errorf("don already exists in CR for chain %d, it has id %d", a.DONChainSelector, donID)
	}

	return nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the plugin config as candidateConfig for the don.
//
// This is the first step to creating a CCIP DON and must be executed before any
// other changesets (SetCandidateChangeset, PromoteAllCandidatesChangeset)
// can be executed.
//
// Note that these operations must be done together because the createDON call
// in the capability registry calls the capability config contract, so we must
// provide suitable calldata for CCIPHome.
func AddDonAndSetCandidateChangeset(
	e deployment.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	err = cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		e.OCRSecrets,
		state.Chains[cfg.DONChainSelector].OffRamp,
		e.Chains[cfg.DONChainSelector],
		nodes.NonBootstraps(),
		state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
		cfg.CCIPOCRParams.OCRParameters,
		cfg.CCIPOCRParams.CommitOffChainConfig,
		cfg.CCIPOCRParams.ExecuteOffChainConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	latestDon, err := internal.LatestCCIPDON(state.Chains[cfg.HomeChainSelector].CapabilityRegistry)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	pluginOCR3Config, ok := newDONArgs[cfg.PluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing commit plugin in ocr3Configs")
	}

	expectedDonID := latestDon.Id + 1
	addDonOp, err := newDonWithCandidateOp(
		txOpts,
		e.Chains[cfg.HomeChainSelector],
		expectedDonID,
		pluginOCR3Config,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		nodes.NonBootstraps(),
		cfg.MCMS != nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
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
			Batch:           []mcms.Operation{addDonOp},
		}},
		fmt.Sprintf("addDON on new Chain && setCandidate for plugin %s", cfg.PluginType.String()),
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

// newDonWithCandidateOp sets the candidate commit config by calling setCandidate on CCIPHome contract through the AddDON call on CapReg contract
// This should be done first before calling any other UpdateDON calls
// This proposes to set up OCR3 config for the commit plugin for the DON
func newDonWithCandidateOp(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
	mcmsEnabled bool,
) (mcms.Operation, error) {
	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack set candidate call: %w", err)
	}

	addDonTx, err := capReg.AddDON(
		txOpts,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false, // isPublic
		false, // acceptsWorkflows
		nodes.DefaultF(),
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("could not generate add don tx w/ commit config: %w", err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, addDonTx, err)
		if err != nil {
			return mcms.Operation{}, fmt.Errorf("error confirming addDon call: %w", err)
		}
	}

	return mcms.Operation{
		To:    capReg.Address(),
		Data:  addDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

type SetCandidateChangesetConfig struct {
	SetCandidateConfigBase
}

func (s SetCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (donID uint32, err error) {
	err = s.SetCandidateConfigBase.Validate(e, state)
	if err != nil {
		return 0, err
	}

	donID, err = internal.DonIDForChain(
		state.Chains[s.HomeChainSelector].CapabilityRegistry,
		state.Chains[s.HomeChainSelector].CCIPHome,
		s.DONChainSelector,
	)
	if err != nil {
		return 0, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("don doesn't exist in CR for chain %d", s.DONChainSelector)
	}

	return donID, nil
}

// SetCandidateChangeset generates a proposal to call setCandidate on the CCIPHome through the capability registry.
// A DON must exist in order to use this changeset effectively, i.e AddDonAndSetCandidateChangeset must be called first.
func SetCandidateChangeset(
	e deployment.Environment,
	cfg SetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	donID, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		e.OCRSecrets,
		state.Chains[cfg.DONChainSelector].OffRamp,
		e.Chains[cfg.DONChainSelector],
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
		txOpts,
		e.Chains[cfg.HomeChainSelector],
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		nodes.NonBootstraps(),
		donID,
		config,
		cfg.MCMS != nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

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
			Batch:           setCandidateMCMSOps,
		}},
		fmt.Sprintf("SetCandidate for %s plugin", cfg.PluginType.String()),
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

// setCandidateOnExistingDon calls setCandidate on CCIPHome contract through the UpdateDON call on CapReg contract
// This proposes to set up OCR3 config for the provided plugin for the DON
func setCandidateOnExistingDon(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	if donID == 0 {
		return nil, fmt.Errorf("donID is zero")
	}

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
		txOpts,
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
		return nil, fmt.Errorf("update don w/ setCandidate call: %w", err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}

// promoteCandidateOp will create the MCMS Operation for `promoteCandidateAndRevokeActive` directed towards the capabilityRegistry
func promoteCandidateOp(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType uint8,
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
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	if donID == 0 {
		return nil, fmt.Errorf("donID is zero")
	}

	var mcmsOps []mcms.Operation
	updateCommitOp, err := promoteCandidateOp(
		txOpts,
		homeChain,
		capReg,
		ccipHome,
		nodes,
		donID,
		uint8(cctypes.PluginTypeCCIPCommit),
		mcmsEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateCommitOp)

	updateExecOp, err := promoteCandidateOp(
		txOpts,
		homeChain,
		capReg,
		ccipHome,
		nodes,
		donID,
		uint8(cctypes.PluginTypeCCIPExec),
		mcmsEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateExecOp)

	return mcmsOps, nil
}

type RevokeCandidateChangesetConfig struct {
	HomeChainSelector uint64

	// DONChainSelector is the chain selector whose candidate config we want to revoke.
	DONChainSelector uint64
	PluginType       types.PluginType

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (r RevokeCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (donID uint32, err error) {
	if err := deployment.IsValidChainSelector(r.HomeChainSelector); err != nil {
		return 0, fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(r.DONChainSelector); err != nil {
		return 0, fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(e.NodeIDs) == 0 {
		return 0, fmt.Errorf("NodeIDs must be set")
	}
	if state.Chains[r.HomeChainSelector].CCIPHome == nil {
		return 0, fmt.Errorf("CCIPHome contract does not exist")
	}
	if state.Chains[r.HomeChainSelector].CapabilityRegistry == nil {
		return 0, fmt.Errorf("CapabilityRegistry contract does not exist")
	}

	// check that the don exists for this chain
	donID, err = internal.DonIDForChain(
		state.Chains[r.HomeChainSelector].CapabilityRegistry,
		state.Chains[r.HomeChainSelector].CCIPHome,
		r.DONChainSelector,
	)
	if err != nil {
		return 0, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("don doesn't exist in CR for chain %d", r.DONChainSelector)
	}

	// check that candidate digest is not zero - this is enforced onchain.
	candidateDigest, err := state.Chains[r.HomeChainSelector].CCIPHome.GetCandidateDigest(nil, donID, uint8(r.PluginType))
	if err != nil {
		return 0, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}
	if candidateDigest == [32]byte{} {
		return 0, fmt.Errorf("candidate config digest is zero, can't revoke it")
	}

	return donID, nil
}

func RevokeCandidateChangeset(e deployment.Environment, cfg RevokeCandidateChangesetConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	donID, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("fetch nodes info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	homeChain := e.Chains[cfg.HomeChainSelector]
	ops, err := revokeCandidateOps(
		txOpts,
		homeChain,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		nodes.NonBootstraps(),
		donID,
		uint8(cfg.PluginType),
		cfg.MCMS != nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("revoke candidate ops: %w", err)
	}
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
			Batch:           ops,
		}},
		fmt.Sprintf("revokeCandidate for don %d", cfg.DONChainSelector),
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

func revokeCandidateOps(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType uint8,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	if donID == 0 {
		return nil, fmt.Errorf("donID is zero")
	}

	candidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, pluginType)
	if err != nil {
		return nil, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}

	encodedRevokeCandidateCall, err := internal.CCIPHomeABI.Pack(
		"revokeCandidate",
		donID,
		pluginType,
		candidateDigest,
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedRevokeCandidateCall,
			},
		},
		false, // isPublic
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ revokeCandidate call: %w", deployment.MaybeDataErr(err))
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}
