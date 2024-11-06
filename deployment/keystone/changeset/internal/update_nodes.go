package internal

import (
	"errors"
	"fmt"

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
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*P2PSignerEnc
}

func (req *UpdateNodesRequest) NodeParams() ([]kcr.CapabilitiesRegistryNodeParams, error) {
	return makeNodeParams(req.Registry, req.NopToNodes, req.P2pToCapabilities)
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
	if len(req.NopToNodes) == 0 {
		return errors.New("nopToNodes is empty")
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
	nopToNodes map[kcr.CapabilitiesRegistryNodeOperator][]*P2PSignerEnc,
	p2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability) ([]kcr.CapabilitiesRegistryNodeParams, error) {

	out := make([]kcr.CapabilitiesRegistryNodeParams, 0)
	// get all the node operators from chain
	registeredNops, err := registry.GetNodeOperators(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node operators: %w", err)
	}

	// make a cache of capability from chain
	var allCaps []kcr.CapabilitiesRegistryCapability
	for _, caps := range p2pToCapabilities {
		allCaps = append(allCaps, caps...)
	}
	capMap, err := fetchCapabilityIDs(registry, allCaps)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch capability ids: %w", err)
	}

	// flatten the onchain state to list of node params filtered by the input nops and nodes
	for idx, rnop := range registeredNops {
		// nop id is 1-indexed. no way to get value from chain. must infer from index
		nopID := uint32(idx + 1)
		nodes, ok := nopToNodes[rnop]
		if !ok {
			continue
		}
		for _, node := range nodes {
			caps, ok := p2pToCapabilities[node.P2PKey]
			if !ok {
				return nil, fmt.Errorf("capabilities not found for node %s", node.P2PKey)
			}
			hashedCaps := make([][32]byte, len(caps))
			for i, cap := range caps {
				hashedCap, exists := capMap[kslib.CapabilityID(cap)]
				if !exists {
					return nil, fmt.Errorf("capability id not found for %s", kslib.CapabilityID(cap))
				}
				hashedCaps[i] = hashedCap
			}
			out = append(out, kcr.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      nopID,
				P2pId:               node.P2PKey,
				HashedCapabilityIds: hashedCaps,
				EncryptionPublicKey: node.EncryptionPublicKey,
				Signer:              node.Signer,
			})
		}
	}

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
