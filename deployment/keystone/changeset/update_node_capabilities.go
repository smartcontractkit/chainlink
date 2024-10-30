package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type UpdateNodeCapabilitiesRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc
}

func (req *UpdateNodeCapabilitiesRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if len(req.NopToNodes) == 0 {
		return fmt.Errorf("nopToNodes is empty")
	}
	if req.Registry == nil {
		return fmt.Errorf("registry is nil")
	}
	return nil
}

// UpdateNodeCapabilibity sets the capabilities of the node to the new capabilities.
// New capabilities are added to the onchain registry and the node is updated to host the new capabilities.
func UpdateNodeCapabilities(lggr logger.Logger, req *UpdateNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	_, err := updateNodeCapabilitiesImpl(lggr, req)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

func updateNodeCapabilitiesImpl(lggr logger.Logger, req *UpdateNodeCapabilitiesRequest) (*kslib.UpdateNodesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	// collect all the capabilities and add them to the registry
	var capabilities []kcr.CapabilitiesRegistryCapability
	for _, cap := range req.P2pToCapabilities {
		capabilities = append(capabilities, cap...)
	}
	err := kslib.AddCapabilities(lggr, req.Registry, req.Chain, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	updateNodesReq := &kslib.UpdateNodesRequest{
		Chain:             req.Chain,
		Registry:          req.Registry,
		P2pToCapabilities: req.P2pToCapabilities,
		NopToNodes:        req.NopToNodes,
	}
	resp, err := kslib.UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
