package internal

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type AppendNodeCapabilitiesRequest struct {
	Chain                deployment.Chain
	CapabilitiesRegistry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	UseMCMS           bool
}

func (req *AppendNodeCapabilitiesRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return errors.New("p2pToCapabilities is empty")
	}
	if req.CapabilitiesRegistry == nil {
		return errors.New("registry is nil")
	}
	return nil
}

func AppendNodeCapabilitiesImpl(lggr logger.Logger, req *AppendNodeCapabilitiesRequest) (*UpdateNodesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	// for each node, merge the new capabilities with the existing ones and update the node
	updatesByPeer := make(map[p2pkey.PeerID]NodeUpdate)
	for p2pID, caps := range req.P2pToCapabilities {
		caps, err := AppendCapabilities(lggr, req.CapabilitiesRegistry, req.Chain, []p2pkey.PeerID{p2pID}, caps)
		if err != nil {
			return nil, fmt.Errorf("failed to append capabilities for p2p %s: %w", p2pID, err)
		}
		updatesByPeer[p2pID] = NodeUpdate{Capabilities: caps[p2pID]}
	}

	// collect all the capabilities and add them to the registry
	var capabilities []kcr.CapabilitiesRegistryCapability
	for _, cap := range req.P2pToCapabilities {
		capabilities = append(capabilities, cap...)
	}
	op, err := AddCapabilities(lggr, req.CapabilitiesRegistry, req.Chain, capabilities, req.UseMCMS)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	updateNodesReq := &UpdateNodesRequest{
		Chain:                req.Chain,
		CapabilitiesRegistry: req.CapabilitiesRegistry,
		P2pToUpdates:         updatesByPeer,
		UseMCMS:              req.UseMCMS,
		Ops:                  op,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
