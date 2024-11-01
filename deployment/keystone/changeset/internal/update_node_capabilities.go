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
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*P2PSignerEnc
}

func (req *UpdateNodeCapabilitiesImplRequest) Validate() error {
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
		NopToNodes:        req.NopToNodes,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}

/*
// AddCapabilities adds the capabilities to the registry
// it tries to add all capabilities in one go, if that fails, it falls back to adding them one by one
func AddCapabilities(lggr logger.Logger, registry *kcr.CapabilitiesRegistry, chain deployment.Chain, capabilities []kcr.CapabilitiesRegistryCapability) error {
	if len(capabilities) == 0 {
		return nil
	}
	// dedup capabilities
	var deduped []kcr.CapabilitiesRegistryCapability
	seen := make(map[string]struct{})
	for _, cap := range capabilities {
		if _, ok := seen[CapabilityID(cap)]; !ok {
			seen[CapabilityID(cap)] = struct{}{}
			deduped = append(deduped, cap)
		}
	}

	tx, err := registry.AddCapabilities(chain.DeployerKey, deduped)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		// no typed errors in the abi, so we have to do string matching
		// try to add all capabilities in one go, if that fails, fall back to 1-by-1
		if !strings.Contains(err.Error(), "CapabilityAlreadyExists") {
			return fmt.Errorf("failed to call AddCapabilities: %w", err)
		}
		lggr.Warnw("capabilities already exist, falling back to 1-by-1", "capabilities", deduped)
		for _, cap := range deduped {
			tx, err = registry.AddCapabilities(chain.DeployerKey, []kcr.CapabilitiesRegistryCapability{cap})
			if err != nil {
				err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
				if strings.Contains(err.Error(), "CapabilityAlreadyExists") {
					lggr.Warnw("capability already exists, skipping", "capability", cap)
					continue
				}
				return fmt.Errorf("failed to call AddCapabilities for capability %v: %w", cap, err)
			}
			// 1-by-1 tx is pending and we need to wait for it to be mined
			_, err = chain.Confirm(tx)
			if err != nil {
				return fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
			}
			lggr.Debugw("registered capability", "capability", cap)

		}
	} else {
		// the bulk add tx is pending and we need to wait for it to be mined
		_, err = chain.Confirm(tx)
		if err != nil {
			return fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
		}
		lggr.Info("registered capabilities", "capabilities", deduped)
	}
	return nil
}
*/
