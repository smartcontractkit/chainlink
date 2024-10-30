package changeset

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

// SetCandidateExecPluginOps calls setCandidate on CCIPHome contract through the UpdateDON call on CapReg contract
// This proposes to set up OCR3 config for the provided plugin for the DON
func SetCandidateOnExistingDon(
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSelector uint64,
	nodes deployment.Nodes,
) ([]mcms.Operation, error) {
	// fetch DON ID for the chain
	donID, err := ccdeploy.DonIDForChain(capReg, ccipHome, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}
	fmt.Printf("donID: %d", donID)
	encodedSetCandidateCall, err := ccdeploy.CCIPHomeABI.Pack(
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
				CapabilityId: ccdeploy.CCIPCapabilityID,
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

// PromoteCandidateOp will create the MCMS Operation for `promoteCandidateAndRevokeActive` directed towards the capabilityRegistry
func PromoteCandidateOp(donID uint32, pluginType uint8, capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome, nodes deployment.Nodes) (mcms.Operation, error) {

	allConfigs, err := ccipHome.GetAllConfigs(nil, donID, pluginType)
	if err != nil {
		return mcms.Operation{}, err
	}

	if allConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
		return mcms.Operation{}, fmt.Errorf("candidate digest is empty, expected nonempty")
	}
	fmt.Printf("commit candidate digest after setCandidate: %x\n", allConfigs.CandidateConfig.ConfigDigest)

	encodedPromotionCall, err := ccdeploy.CCIPHomeABI.Pack(
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
				CapabilityId: ccdeploy.CCIPCapabilityID,
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

// PromoteAllCandidatesForChainOps promotes the candidate commit and exec configs to active by calling promoteCandidateAndRevokeActive on CCIPHome through the UpdateDON call on CapReg contract
func PromoteAllCandidatesForChainOps(
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSelector uint64,
	nodes deployment.Nodes,
) ([]mcms.Operation, error) {
	// fetch DON ID for the chain
	donID, err := ccdeploy.DonIDForChain(capReg, ccipHome, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}

	var mcmsOps []mcms.Operation
	updateCommitOp, err := PromoteCandidateOp(donID, uint8(cctypes.PluginTypeCCIPCommit), capReg, ccipHome, nodes)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateCommitOp)

	updateExecOp, err := PromoteCandidateOp(donID, uint8(cctypes.PluginTypeCCIPExec), capReg, ccipHome, nodes)
	if err != nil {
		return nil, fmt.Errorf("promote candidate op: %w", err)
	}
	mcmsOps = append(mcmsOps, updateExecOp)

	return mcmsOps, nil
}

// PromoteAllCandidatesProposal generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
func PromoteAllCandidatesProposal(
	state ccdeploy.CCIPOnChainState,
	homeChainSel, newChainSel uint64,
	nodes deployment.Nodes,
) (*timelock.MCMSWithTimelockProposal, error) {
	promoteCandidateOps, err := PromoteAllCandidatesForChainOps(
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return nil, err
	}

	return ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           promoteCandidateOps,
	}}, "promoteCandidate for commit and execution", 0)
}

// SetCandidateExecPluginProposal calls setCandidate on the CCIPHome for setting up OCR3 exec Plugin config for the new chain.
func SetCandidatePluginProposal(
	state ccdeploy.CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig ccdeploy.TokenConfig,
	pluginType cctypes.PluginType,
) (*timelock.MCMSWithTimelockProposal, error) {
	newDONArgs, err := ccdeploy.BuildOCR3ConfigForCCIPHome(
		e.Logger,
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
		nodes.NonBootstraps(),
		state.Chains[homeChainSel].RMNHome.Address(),
	)
	if err != nil {
		return nil, err
	}

	execConfig, ok := newDONArgs[pluginType]
	if !ok {
		return nil, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	setCandidateMCMSOps, err := SetCandidateOnExistingDon(
		execConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return nil, err
	}

	return ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           setCandidateMCMSOps,
	}}, "SetCandidate for execution", 0)
}
