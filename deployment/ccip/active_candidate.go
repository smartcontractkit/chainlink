package ccipdeployment

import (
	"fmt"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"math/big"
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
	donID, err := DonIDForChain(capReg, ccipHome, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch don id for chain: %w", err)
	}
	fmt.Printf("donID: %d", donID)
	encodedSetCandidateCall, err := CCIPHomeABI.Pack(
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
				CapabilityId: CCIPCapabilityID,
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

	encodedPromotionCall, err := CCIPHomeABI.Pack(
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
				CapabilityId: CCIPCapabilityID,
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
	donID, err := DonIDForChain(capReg, ccipHome, chainSelector)
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
