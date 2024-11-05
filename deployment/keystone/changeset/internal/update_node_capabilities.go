package internal

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type UpdateNodeCapabilitiesImplRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
}

func (req *UpdateNodeCapabilitiesImplRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if req.Registry == nil {
		return fmt.Errorf("registry is nil")
	}

	return nil
}

func UpdateNodeCapabilitiesImpl(lggr logger.Logger, req *UpdateNodeCapabilitiesImplRequest) (*UpdateNodesResponse, error) {
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

	updateNodesReq := &UpdateNodesRequest{
		Chain:             req.Chain,
		Registry:          req.Registry,
		P2pToCapabilities: req.P2pToCapabilities,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
