package internal

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type UpdateNodeCapabilitiesImplRequest struct {
	Chain                cldf_evm.Chain
	CapabilitiesRegistry *kcr.CapabilitiesRegistry
	P2pToCapabilities    map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability

	UseMCMS bool
}

func (req *UpdateNodeCapabilitiesImplRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return errors.New("p2pToCapabilities is empty")
	}
	if req.CapabilitiesRegistry == nil {
		return errors.New("registry is nil")
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
	op, err := AddCapabilities(lggr, req.CapabilitiesRegistry, req.Chain, capabilities, req.UseMCMS)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	p2pToUpdates := map[p2pkey.PeerID]NodeUpdate{}
	for id, caps := range req.P2pToCapabilities {
		p2pToUpdates[id] = NodeUpdate{Capabilities: caps}
	}

	updateNodesReq := &UpdateNodesRequest{
		Chain:                req.Chain,
		P2pToUpdates:         p2pToUpdates,
		CapabilitiesRegistry: req.CapabilitiesRegistry,
		Ops:                  op,
		UseMCMS:              req.UseMCMS,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
