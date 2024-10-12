package keystone

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"google.golang.org/protobuf/proto"
)

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig struct {
	Capability                                      kcr.CapabilitiesRegistryCapability
	kcr.CapabilitiesRegistryCapabilityConfiguration // embed bc otherwise with ~ config.config
}

type UpdateDonRequest struct {
	Registry *kcr.CapabilitiesRegistry
	Chain    deployment.Chain

	Name   string
	P2PIDs []p2pkey.PeerID // this is the unique identifier for the done
	// CapabilityId field is ignored. it is determined dynamically by the registry
	// If the Config is nil, a default config is used
	CapabilityConfigs []CapabilityConfig
}

func (req *UpdateDonRequest) Validate() error {
	if req.Registry == nil {
		return fmt.Errorf("registry is nil")
	}
	return nil
}

type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
}

func UpdateDon(lggr logger.Logger, req *UpdateDonRequest) (*UpdateDonResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	getDonsResp, err := req.Registry.GetDONs(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Dons: %w", err)
	}
	wantedDonID := sortedHash(peerIdsToBytes(req.P2PIDs))
	var don kcr.CapabilitiesRegistryDONInfo
	found := false
	for i, di := range getDonsResp {
		gotID := sortedHash(di.NodeP2PIds)
		if gotID == wantedDonID {
			don = getDonsResp[i]
			found = true
			break
		}
	}
	if !found {
		type debugDonInfo struct {
			OnchainID  uint32
			P2PIDsHash string
		}
		debugIds := make([]debugDonInfo, len(getDonsResp))
		for i, di := range getDonsResp {
			debugIds[i] = debugDonInfo{
				OnchainID:  di.Id,
				P2PIDsHash: sortedHash(di.NodeP2PIds),
			}
		}
		return nil, fmt.Errorf("don not found by p2pIDs %s in %v", wantedDonID, debugIds)
	}
	cfgs, err := computeConfigs(req.Registry, req.CapabilityConfigs, don)
	if err != nil {
		return nil, fmt.Errorf("failed to compute configs: %w", err)
	}
	tx, err := req.Registry.UpdateDON(req.Chain.DeployerKey, don.Id, don.NodeP2PIds, cfgs, don.IsPublic, don.F)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call UpdateDON: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm UpdateDON transaction %s: %w", tx.Hash().String(), err)
	}
	out := don
	out.CapabilityConfigurations = cfgs
	return &UpdateDonResponse{DonInfo: out}, nil
}

func peerIdsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	out := make([][32]byte, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}

func computeConfigs(registry *kcr.CapabilitiesRegistry, caps []CapabilityConfig, donInfo kcr.CapabilitiesRegistryDONInfo) ([]kcr.CapabilitiesRegistryCapabilityConfiguration, error) {
	out := make([]kcr.CapabilitiesRegistryCapabilityConfiguration, len(caps))
	for i, cap := range caps {
		out[i] = cap.CapabilitiesRegistryCapabilityConfiguration
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.Capability.LabelledName, cap.Capability.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id: %w", err)
		}
		out[i].CapabilityId = id
		if out[i].Config == nil {
			c := defaultCapConfig(cap.Capability.CapabilityType, int(donInfo.F))
			cb, err := proto.Marshal(c)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal capability config for %v: %w", c, err)
			}
			out[i].Config = cb
		}
	}
	return out, nil
}
