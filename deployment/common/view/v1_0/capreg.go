package v1_0

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// CapabilityRegistryView is a high-fidelity view of the capabilities registry contract.
type CapabilityRegistryView struct {
	types.ContractMetaData
	Capabilities []CapabilityView `json:"capabilities,omitempty"`
	Nodes        []NodeView       `json:"nodes,omitempty"`
	Nops         []NopView        `json:"nops,omitempty"`
	Dons         []DonView        `json:"dons,omitempty"`
}

// MarshalJSON marshals the CapabilityRegistryView to JSON. It includes the Capabilities, Nodes, Nops, and Dons
// and a denormalized summary of the Dons with their associated Nodes and Capabilities, which is useful for a high-level view
func (v *CapabilityRegistryView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias struct {
		types.ContractMetaData
		Capabilities    []CapabilityView      `json:"capabilities,omitempty"`
		Nodes           []NodeView            `json:"nodes,omitempty"`
		Nops            []NopView             `json:"nops,omitempty"`
		Dons            []DonView             `json:"dons,omitempty"`
		DonCapabilities []DonDenormalizedView `json:"don_capabilities_summary,omitempty"`
	}
	a := Alias{
		ContractMetaData: v.ContractMetaData,
		Capabilities:     v.Capabilities,
		Nodes:            v.Nodes,
		Nops:             v.Nops,
		Dons:             v.Dons,
	}
	dc, err := v.DonDenormalizedView()
	if err != nil {
		return nil, err
	}
	a.DonCapabilities = dc
	return json.MarshalIndent(&a, "", " ")
}

// UnmarshalJSON unmarshals the CapabilityRegistryView from JSON. Since the CapabilityRegistryView doesn't hold a DonCapabilities field,
// it is not unmarshaled.
func (v *CapabilityRegistryView) UnmarshalJSON(data []byte) error {
	// Alias to avoid recursive calls
	type Alias struct {
		types.ContractMetaData
		Capabilities    []CapabilityView      `json:"capabilities,omitempty"`
		Nodes           []NodeView            `json:"nodes,omitempty"`
		Nops            []NopView             `json:"nops,omitempty"`
		Dons            []DonView             `json:"dons,omitempty"`
		DonCapabilities []DonDenormalizedView `json:"don_capabilities_summary,omitempty"`
	}
	a := Alias{
		ContractMetaData: v.ContractMetaData,
		Capabilities:     v.Capabilities,
		Nodes:            v.Nodes,
		Nops:             v.Nops,
		Dons:             v.Dons,
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	v.ContractMetaData = a.ContractMetaData
	v.Capabilities = a.Capabilities
	v.Nodes = a.Nodes
	v.Nops = a.Nops
	v.Dons = a.Dons
	return nil
}

// GenerateCapabilityRegistryView generates a CapRegView from a CapabilitiesRegistry contract.
func GenerateCapabilityRegistryView(capReg *capabilities_registry.CapabilitiesRegistry) (CapabilityRegistryView, error) {
	tv, err := types.NewContractMetaData(capReg, capReg.Address())
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	caps, err := capReg.GetCapabilities(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var capViews []CapabilityView
	for _, capability := range caps {
		capViews = append(capViews, NewCapabilityView(capability))
	}
	donInfos, err := capReg.GetDONs(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var donViews []DonView
	for _, donInfo := range donInfos {
		donViews = append(donViews, NewDonView(donInfo))
	}

	nodeInfos, err := capReg.GetNodes(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var nodeViews []NodeView
	for _, nodeInfo := range nodeInfos {
		nodeViews = append(nodeViews, NewNodeView(nodeInfo))
	}

	nopInfos, err := capReg.GetNodeOperators(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var nopViews []NopView
	for _, nopInfo := range nopInfos {
		nopViews = append(nopViews, NewNopView(nopInfo))
	}

	return CapabilityRegistryView{
		ContractMetaData: tv,
		Capabilities:     capViews,
		Dons:             donViews,
		Nodes:            nodeViews,
		Nops:             nopViews,
	}, nil
}

// DonDenormalizedView is a view of a Don with its associated Nodes and Capabilities.
type DonDenormalizedView struct {
	Don          DonUniversalMetadata   `json:"don"`
	Nodes        []NodeDenormalizedView `json:"nodes"`
	Capabilities []CapabilityView       `json:"capabilities"`
}

// DonDenormalizedView returns a list of DonDenormalizedView, which are Dons with their associated
// Nodes and Capabilities. This is a useful form of the CapabilityRegistryView, but it is not definitive.
// The full CapRegView should be used for the most accurate information as it can contain
// Capabilities and Nodes the are not associated with any Don.
func (v *CapabilityRegistryView) DonDenormalizedView() ([]DonDenormalizedView, error) {
	var out []DonDenormalizedView
	for _, don := range v.Dons {
		var nodes []NodeDenormalizedView
		for _, node := range v.Nodes {
			if don.hasNode(node) {
				ndv, err := v.nodeDenormalizedView(node)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, ndv)
			}
		}
		var capabilities []CapabilityView
		for _, cap := range v.Capabilities {
			if don.hasCapability(cap) {
				capabilities = append(capabilities, cap)
			}
		}
		out = append(out, DonDenormalizedView{
			Don:          don.DonUniversalMetadata,
			Nodes:        nodes,
			Capabilities: capabilities,
		})
	}
	return out, nil
}

func (v *CapabilityRegistryView) NodesToNodesParams() ([]capabilities_registry.CapabilitiesRegistryNodeParams, error) {
	var nodesParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, node := range v.Nodes {
		signer, err := hexTo32Bytes(node.Signer)
		if err != nil {
			return nil, err
		}
		encryptionPubKey, err := hexTo32Bytes(node.EncryptionPublicKey)
		if err != nil {
			return nil, err
		}
		capIDs := make([][32]byte, len(node.CapabilityIDs))
		for i, id := range node.CapabilityIDs {
			cid, err := hexTo32Bytes(id)
			if err != nil {
				return nil, err
			}
			capIDs[i] = cid
		}
		nodesParams = append(nodesParams, capabilities_registry.CapabilitiesRegistryNodeParams{
			Signer:              signer,
			P2pId:               node.P2pId,
			EncryptionPublicKey: encryptionPubKey,
			NodeOperatorId:      node.NodeOperatorID,
			HashedCapabilityIds: capIDs,
		})
	}

	return nodesParams, nil
}

func (v *CapabilityRegistryView) CapabilitiesToCapabilitiesParams() []capabilities_registry.CapabilitiesRegistryCapability {
	var capabilitiesParams []capabilities_registry.CapabilitiesRegistryCapability
	for _, capability := range v.Capabilities {
		capabilitiesParams = append(capabilitiesParams, capabilities_registry.CapabilitiesRegistryCapability{
			LabelledName:          capability.LabelledName,
			Version:               capability.Version,
			CapabilityType:        capability.CapabilityType,
			ResponseType:          capability.ResponseType,
			ConfigurationContract: capability.ConfigurationContract,
		})
	}
	return capabilitiesParams
}

func (v *CapabilityRegistryView) NopsToNopsParams() []capabilities_registry.CapabilitiesRegistryNodeOperator {
	var nopsParams []capabilities_registry.CapabilitiesRegistryNodeOperator
	for _, nop := range v.Nops {
		nopsParams = append(nopsParams, capabilities_registry.CapabilitiesRegistryNodeOperator{
			Admin: nop.Admin,
			Name:  nop.Name,
		})
	}
	return nopsParams
}

func (v *CapabilityRegistryView) CapabilityConfigToCapabilityConfigParams(don DonView) ([]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration, error) {
	var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
	for _, cfg := range don.CapabilityConfigurations {
		cid, err := hexTo32Bytes(cfg.ID)
		if err != nil {
			return nil, err
		}
		config, err := hex.DecodeString(cfg.Config)
		if err != nil {
			return nil, err
		}
		cfgs = append(cfgs, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: cid,
			Config:       config,
		})
	}
	return cfgs, nil
}

func hexTo32Bytes(val string) ([32]byte, error) {
	var out [32]byte
	b, err := hex.DecodeString(val)
	if err != nil {
		return out, err
	}
	copy(out[:], b)
	return out, nil
}

// CapabilityView is a serialization-friendly view of a capability in the capabilities registry.
type CapabilityView struct {
	ID                    string         `json:"id"` // hex 32 bytes
	LabelledName          string         `json:"labelled_name"`
	Version               string         `json:"version"`
	CapabilityType        uint8          `json:"capability_type"`
	ResponseType          uint8          `json:"response_type"`
	ConfigurationContract common.Address `json:"configuration_contract,omitempty"`
	IsDeprecated          bool           `json:"is_deprecated,omitempty"`
}

// NewCapabilityView creates a CapabilityView from a CapabilitiesRegistryCapabilityInfo.
func NewCapabilityView(capInfo capabilities_registry.CapabilitiesRegistryCapabilityInfo) CapabilityView {
	return CapabilityView{
		ID:                    hex.EncodeToString(capInfo.HashedId[:]),
		LabelledName:          capInfo.LabelledName,
		Version:               capInfo.Version,
		CapabilityType:        capInfo.CapabilityType,
		ResponseType:          capInfo.ResponseType,
		ConfigurationContract: capInfo.ConfigurationContract,
		IsDeprecated:          capInfo.IsDeprecated,
	}
}

// Validate checks that the CapabilityView is valid.
func (cv CapabilityView) Validate() error {
	id, err := hex.DecodeString(cv.ID)
	if err != nil {
		return err
	}
	if len(id) != 32 {
		return errors.New("capability id must be 32 bytes")
	}
	return nil
}

// DonView is a serialization-friendly view of a Don in the capabilities registry.
type DonView struct {
	DonUniversalMetadata
	NodeP2PIds               []p2pkey.PeerID             `json:"node_p2p_ids,omitempty"`
	CapabilityConfigurations []CapabilitiesConfiguration `json:"capability_configurations,omitempty"`
}

type DonUniversalMetadata struct {
	ID               uint32 `json:"id"`
	ConfigCount      uint32 `json:"config_count"`
	F                uint8  `json:"f"`
	IsPublic         bool   `json:"is_public,omitempty"`
	AcceptsWorkflows bool   `json:"accepts_workflows,omitempty"`
}

// NewDonView creates a DonView from a CapabilitiesRegistryDONInfo.
func NewDonView(d capabilities_registry.CapabilitiesRegistryDONInfo) DonView {
	return DonView{
		DonUniversalMetadata: DonUniversalMetadata{
			ID:               d.Id,
			ConfigCount:      d.ConfigCount,
			F:                d.F,
			IsPublic:         d.IsPublic,
			AcceptsWorkflows: d.AcceptsWorkflows,
		},
		NodeP2PIds:               p2pIds(d.NodeP2PIds),
		CapabilityConfigurations: NewCapabilityConfigurations(d.CapabilityConfigurations),
	}
}

func (dv DonView) Validate() error {
	for i, cfg := range dv.CapabilityConfigurations {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("capability configuration at index %d invalid:%w ", i, err)
		}
	}
	return nil
}

// CapabilitiesConfiguration is a serialization-friendly view of a capability configuration in the capabilities registry.
type CapabilitiesConfiguration struct {
	ID     string `json:"id"`     // hex 32 bytes
	Config string `json:"config"` // hex
}

// NewCapabilityConfigurations creates a list of CapabilitiesConfiguration from a list of CapabilitiesRegistryCapabilityConfiguration.
func NewCapabilityConfigurations(cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration) []CapabilitiesConfiguration {
	var out []CapabilitiesConfiguration
	for _, cfg := range cfgs {
		out = append(out, CapabilitiesConfiguration{
			ID:     hex.EncodeToString(cfg.CapabilityId[:]),
			Config: hex.EncodeToString(cfg.Config),
		})
	}
	return out
}

func (cc CapabilitiesConfiguration) Validate() error {
	id, err := hex.DecodeString(cc.ID)
	if err != nil {
		return errors.New("capability id must be hex encoded")
	}
	if len(id) != 32 {
		return errors.New("capability id must be 32 bytes")
	}
	_, err = hex.DecodeString(cc.Config)
	if err != nil {
		return errors.New("config must be hex encoded")
	}
	return nil
}

// NodeView is a serialization-friendly view of a node in the capabilities registry.
type NodeView struct {
	NodeUniversalMetadata
	NodeOperatorID   uint32     `json:"node_operator_id"`
	CapabilityIDs    []string   `json:"capability_ids,omitempty"` // hex 32 bytes
	CapabilityDONIDs []*big.Int `json:"capability_don_ids,omitempty"`
}

// NodeUniversalMetadata is a serialization-friendly view of the universal metadata of a node in the capabilities registry.
type NodeUniversalMetadata struct {
	ConfigCount         uint32        `json:"config_count"`
	WorkflowDONID       uint32        `json:"workflow_don_id"`
	Signer              string        `json:"signer"` // hex 32 bytes
	P2pId               p2pkey.PeerID `json:"p2p_id"`
	EncryptionPublicKey string        `json:"encryption_public_key"` // hex 32 bytes
}

// NewNodeView creates a NodeView from a CapabilitiesRegistryNodeInfoProviderNodeInfo.
func NewNodeView(n capabilities_registry.INodeInfoProviderNodeInfo) NodeView {
	return NodeView{
		NodeUniversalMetadata: NodeUniversalMetadata{
			ConfigCount:         n.ConfigCount,
			WorkflowDONID:       n.WorkflowDONId,
			Signer:              hex.EncodeToString(n.Signer[:]),
			P2pId:               n.P2pId,
			EncryptionPublicKey: hex.EncodeToString(n.EncryptionPublicKey[:]),
		},
		NodeOperatorID:   n.NodeOperatorId,
		CapabilityIDs:    hexIds(n.HashedCapabilityIds),
		CapabilityDONIDs: n.CapabilitiesDONIds,
	}
}

func (nv NodeView) Validate() error {
	s, err := hex.DecodeString(nv.Signer)
	if err != nil {
		return errors.New("signer must be hex encoded")
	}
	if len(s) != 32 {
		return errors.New("signer must be 32 bytes")
	}

	e, err := hex.DecodeString(nv.EncryptionPublicKey)
	if err != nil {
		return errors.New("encryption public key must be hex encoded")
	}
	if len(e) != 32 {
		return errors.New("encryption public key must be 32 bytes")
	}

	for _, id := range nv.CapabilityIDs {
		cid, err := hex.DecodeString(id)
		if err != nil {
			return errors.New("hashed capability id must be hex encoded")
		}
		if len(cid) != 32 {
			return errors.New("hashed capability id must be 32 bytes")
		}
	}
	return nil
}

// NodeDenormalizedView is a serialization-friendly view of a node in the capabilities registry with its associated NOP.
type NodeDenormalizedView struct {
	NodeUniversalMetadata
	Nop NopView `json:"nop"`
}

type NopView struct {
	Admin common.Address `json:"admin"`
	Name  string         `json:"name"`
}

func NewNopView(nop capabilities_registry.CapabilitiesRegistryNodeOperator) NopView {
	return NopView{
		Admin: nop.Admin,
		Name:  nop.Name,
	}
}

func (v *CapabilityRegistryView) nodeDenormalizedView(n NodeView) (NodeDenormalizedView, error) {
	nop, err := nodeNop(n, v.Nops)
	if err != nil {
		return NodeDenormalizedView{}, err
	}
	return NodeDenormalizedView{
		NodeUniversalMetadata: n.NodeUniversalMetadata,
		Nop:                   nop,
	}, nil
}

func nodeNop(n NodeView, nops []NopView) (NopView, error) {
	for i, nop := range nops {
		// nops are 1-indexed. there is no natural key to match on, so we use the index.
		idx := i + 1
		if n.NodeOperatorID == uint32(idx) {
			return nop, nil
		}
	}
	return NopView{}, fmt.Errorf("could not find nop for node %d", n.NodeOperatorID)
}

func p2pIds(rawIds [][32]byte) []p2pkey.PeerID {
	var out []p2pkey.PeerID
	for _, id := range rawIds {
		out = append(out, p2pkey.PeerID(id))
	}
	return out
}

func hexIds(ids [][32]byte) []string {
	var out []string
	for _, id := range ids {
		out = append(out, hex.EncodeToString(id[:]))
	}
	return out
}

func (v DonView) hasNode(node NodeView) bool {
	donId := big.NewInt(int64(v.ID))
	return slices.ContainsFunc(node.CapabilityDONIDs, func(elem *big.Int) bool { return elem.Cmp(donId) == 0 }) || node.WorkflowDONID == v.ID
}

func (v DonView) hasCapability(candidate CapabilityView) bool {
	return slices.ContainsFunc(v.CapabilityConfigurations, func(elem CapabilitiesConfiguration) bool { return elem.ID == candidate.ID })
}
