package changeset

import (
	"fmt"
	"math/big"

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

// PromoteAllCandidatesChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
// TODO: make it conform to the ChangeSet interface.
func PromoteAllCandidatesChangeset(
	state CCIPOnChainState,
	homeChainSel, newChainSel uint64,
	nodes deployment.Nodes,
) (deployment.ChangesetOutput, error) {
	promoteCandidateOps, err := promoteAllCandidatesForChainOps(
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
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

// SetCandidateExecPluginProposal calls setCandidate on the CCIPHome for setting up OCR3 exec Plugin config for the new chain.
// TODO: make it conform to the ChangeSet interface.
func SetCandidatePluginChangeset(
	state CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig TokenConfig,
	pluginType cctypes.PluginType,
) (deployment.ChangesetOutput, error) {
	ccipOCRParams := DefaultOCRParams(
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
		nil,
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

	execConfig, ok := newDONArgs[pluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	setCandidateMCMSOps, err := setCandidateOnExistingDon(
		execConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
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
func promoteCandidateOp(donID uint32, pluginType uint8, capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome, nodes deployment.Nodes) (mcms.Operation, error) {

	allConfigs, err := ccipHome.GetAllConfigs(nil, donID, pluginType)
	if err != nil {
		return mcms.Operation{}, err
	}

	if allConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return mcms.Operation{}, fmt.Errorf("candidate digest is empty, expected nonempty")
	}
	fmt.Printf("commit candidate digest after setCandidate: %x\n", allConfigs.CandidateConfig.ConfigDigest)

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
		deployment.SimTransactOpts(),
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
	return mcms.Operation{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

// promoteAllCandidatesForChainOps promotes the candidate commit and exec configs to active by calling promoteCandidateAndRevokeActive on CCIPHome through the UpdateDON call on CapReg contract
func promoteAllCandidatesForChainOps(
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

	var mcmsOps []mcms.Operation
	updateCommitOp, err := promoteCandidateOp(donID, uint8(cctypes.PluginTypeCCIPCommit), capReg, ccipHome, nodes)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateCommitOp)

	updateExecOp, err := promoteCandidateOp(donID, uint8(cctypes.PluginTypeCCIPExec), capReg, ccipHome, nodes)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateExecOp)

	return mcmsOps, nil
}
