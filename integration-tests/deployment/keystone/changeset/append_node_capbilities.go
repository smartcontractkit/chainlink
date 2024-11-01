package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

type AppendNodeCapabilitiesRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc
}

func (req *AppendNodeCapabilitiesRequest) Validate() error {
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

// AppendNodeCapabilibity adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(lggr logger.Logger, req *AppendNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	_, err := appendNodeCapabilitiesImpl(lggr, req)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

func appendNodeCapabilitiesImpl(lggr logger.Logger, req *AppendNodeCapabilitiesRequest) (*kslib.UpdateNodesResponse, error) {
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

	// for each node, merge the new capabilities with the existing ones and update the node
	capsByPeer := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for p2pID, caps := range req.P2pToCapabilities {
		caps, err := kslib.AppendCapabilities(lggr, req.Registry, req.Chain, []p2pkey.PeerID{p2pID}, caps)
		if err != nil {
			return nil, fmt.Errorf("failed to append capabilities for p2p %s: %w", p2pID, err)
		}
		capsByPeer[p2pID] = caps[p2pID]
	}

	updateNodesReq := &kslib.UpdateNodesRequest{
		Chain:             req.Chain,
		Registry:          req.Registry,
		P2pToCapabilities: capsByPeer,
		NopToNodes:        req.NopToNodes,
	}
	resp, err := kslib.UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
