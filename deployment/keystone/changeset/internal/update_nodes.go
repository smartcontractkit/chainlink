package internal

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

type UpdateNodesRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
}

func (req *UpdateNodesRequest) NodeParams() ([]kcr.CapabilitiesRegistryNodeParams, error) {
	return makeNodeParams(req.Registry, req.P2pToCapabilities)
}

// P2PSignerEnc represent the key fields in kcr.CapabilitiesRegistryNodeParams
// these values are obtain-able directly from the offchain node
type P2PSignerEnc struct {
	Signer              [32]byte
	P2PKey              p2pkey.PeerID
	EncryptionPublicKey [32]byte
}

func (req *UpdateNodesRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return errors.New("p2pToCapabilities is empty")
	}
	// no duplicate capabilities
	for peer, caps := range req.P2pToCapabilities {
		seen := make(map[string]struct{})
		for _, cap := range caps {
			id := kslib.CapabilityID(cap)
			if _, exists := seen[id]; exists {
				return fmt.Errorf("duplicate capability %s for %s", id, peer)
			}
			seen[id] = struct{}{}
		}
	}

	if req.Registry == nil {
		return errors.New("registry is nil")
	}

	return nil
}

type UpdateNodesResponse struct {
	NodeParams []kcr.CapabilitiesRegistryNodeParams
}

// UpdateNodes updates the nodes in the registry
// the update sets the signer and capabilities for each node. it does not append capabilities to the existing ones
func UpdateNodes(lggr logger.Logger, req *UpdateNodesRequest) (*UpdateNodesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	params, err := req.NodeParams()
	if err != nil {
		err = kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to make node params: %w", err)
	}
	tx, err := req.Registry.UpdateNodes(req.Chain.DeployerKey, params)
	if err != nil {
		err = kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call UpdateNodes: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm UpdateNodes confirm transaction %s: %w", tx.Hash().String(), err)
	}
	return &UpdateNodesResponse{NodeParams: params}, nil
}

// AppendCapabilities appends the capabilities to the existing capabilities of the nodes listed in p2pIds in the registry
func AppendCapabilities(lggr logger.Logger, registry *kcr.CapabilitiesRegistry, chain deployment.Chain, p2pIds []p2pkey.PeerID, capabilities []kcr.CapabilitiesRegistryCapability) (map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability, error) {
	out := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	allCapabilities, err := registry.GetCapabilities(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to GetCapabilities from registry: %w", err)
	}
	var capMap = make(map[[32]byte]kcr.CapabilitiesRegistryCapability)
	for _, cap := range allCapabilities {
		capMap[cap.HashedId] = kcr.CapabilitiesRegistryCapability{
			LabelledName:          cap.LabelledName,
			Version:               cap.Version,
			CapabilityType:        cap.CapabilityType,
			ResponseType:          cap.ResponseType,
			ConfigurationContract: cap.ConfigurationContract,
		}
	}

	for _, p2pID := range p2pIds {
		// read the existing capabilities for the node
		info, err := registry.GetNode(&bind.CallOpts{}, p2pID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node info for %s: %w", p2pID, err)
		}
		mergedCaps := make([]kcr.CapabilitiesRegistryCapability, 0)
		// we only have the id; need to fetch the capabilities details
		for _, capID := range info.HashedCapabilityIds {
			c, exists := capMap[capID]
			if !exists {
				return nil, fmt.Errorf("capability not found for %s", capID)
			}
			mergedCaps = append(mergedCaps, c)
		}
		// append the new capabilities and dedup
		mergedCaps = append(mergedCaps, capabilities...)
		var deduped []kcr.CapabilitiesRegistryCapability
		seen := make(map[string]struct{})
		for _, cap := range mergedCaps {
			if _, ok := seen[kslib.CapabilityID(cap)]; !ok {
				seen[kslib.CapabilityID(cap)] = struct{}{}
				deduped = append(deduped, cap)
			}
		}
		out[p2pID] = deduped
	}
	return out, nil
}

func makeNodeParams(registry *kcr.CapabilitiesRegistry,
	p2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability) ([]kcr.CapabilitiesRegistryNodeParams, error) {

	var out []kcr.CapabilitiesRegistryNodeParams
	var p2pIds []p2pkey.PeerID
	for p2pID := range p2pToCapabilities {
		p2pIds = append(p2pIds, p2pID)
	}

	nodes, err := registry.GetNodesByP2PIds(&bind.CallOpts{}, PeerIDsToBytes(p2pIds))
	if err != nil {
		err = kslib.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to get nodes by p2p ids: %w", err)
	}
	for _, node := range nodes {
		caps, ok := p2pToCapabilities[node.P2pId]
		if !ok {
			return nil, fmt.Errorf("capabilities not found for node %s", node.P2pId)
		}
		ids, err := capabilityIds(registry, caps)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability ids: %w", err)
		}
		out = append(out, kcr.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      node.NodeOperatorId,
			P2pId:               node.P2pId,
			HashedCapabilityIds: ids,
			EncryptionPublicKey: node.EncryptionPublicKey,
			Signer:              node.Signer,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].NodeOperatorId == out[j].NodeOperatorId {
			return bytes.Compare(out[i].P2pId[:], out[j].P2pId[:]) < 0
		}
		return out[i].NodeOperatorId < out[j].NodeOperatorId
	})

	return out, nil

}

// fetchkslib.CapabilityIDs fetches the capability ids for the given capabilities
func fetchCapabilityIDs(registry *kcr.CapabilitiesRegistry, caps []kcr.CapabilitiesRegistryCapability) (map[string][32]byte, error) {
	out := make(map[string][32]byte)
	for _, cap := range caps {
		name := kslib.CapabilityID(cap)
		if _, exists := out[name]; exists {
			continue
		}
		hashId, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id for %s: %w", name, err)
		}
		out[name] = hashId
	}
	return out, nil
}

func capabilityIds(registry *kcr.CapabilitiesRegistry, caps []kcr.CapabilitiesRegistryCapability) ([][32]byte, error) {
	out := make([][32]byte, len(caps))
	for i, cap := range caps {
		id, err := registry.GetHashedCapabilityId(&bind.CallOpts{}, cap.LabelledName, cap.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to get capability id: %w", err)
		}
		out[i] = id
	}
	return out, nil
}
