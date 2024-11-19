package internal

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type AppendNodeCapabilitiesRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
}

func (req *AppendNodeCapabilitiesRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if req.Registry == nil {
		return fmt.Errorf("registry is nil")
	}
	return nil
}

func AppendNodeCapabilitiesImpl(lggr logger.Logger, req *AppendNodeCapabilitiesRequest) (*UpdateNodesResponse, error) {
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
	updatesByPeer := make(map[p2pkey.PeerID]NodeUpdate)
	for p2pID, caps := range req.P2pToCapabilities {
		caps, err := AppendCapabilities(lggr, req.Registry, req.Chain, []p2pkey.PeerID{p2pID}, caps)
		if err != nil {
			return nil, fmt.Errorf("failed to append capabilities for p2p %s: %w", p2pID, err)
		}
		updatesByPeer[p2pID] = NodeUpdate{Capabilities: caps[p2pID]}
	}

	updateNodesReq := &UpdateNodesRequest{
		Chain:        req.Chain,
		Registry:     req.Registry,
		P2pToUpdates: updatesByPeer,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
