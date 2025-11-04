package v2_0

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type CapabilityRegistryViewV2 struct {
	types.ContractMetaData
	Capabilities []CapabilityView `json:"capabilities,omitempty"`
	Nodes        []NodeView       `json:"nodes,omitempty"`
	Nops         []NopView        `json:"nops,omitempty"`
	Dons         []DonView        `json:"dons,omitempty"`
}

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

type unpagniatedCapabilityRegistry interface {
	GetCapabilitiesSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryCapabilityInfo, error)
	GetNodesSimple(opts *bind.CallOpts) ([]capabilities_registry.INodeInfoProviderNodeInfo, error)
	GetNodeOperatorsSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryNodeOperatorInfo, error)
	GetDONsSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryDONInfo, error)
}

var (
	MaxCapabilities = big.NewInt(128)
	MaxDONs         = big.NewInt(32)
	MaxNodes        = big.NewInt(256)
	MaxNOPs         = big.NewInt(128)
)

type ExtendedCapabilityRegistry struct {
	*capabilities_registry.CapabilitiesRegistry
}

var _ unpagniatedCapabilityRegistry = (*ExtendedCapabilityRegistry)(nil)

// GetCapabilitiesSimple implements unpagniatedCapabilityRegistry
func (e *ExtendedCapabilityRegistry) GetCapabilitiesSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryCapabilityInfo, error) {
	return e.GetCapabilities(opts, big.NewInt(0), MaxCapabilities)
}

// GetNodesSimple implements unpagniatedCapabilityRegistry
func (e *ExtendedCapabilityRegistry) GetNodesSimple(opts *bind.CallOpts) ([]capabilities_registry.INodeInfoProviderNodeInfo, error) {
	return e.GetNodes(opts, big.NewInt(0), MaxNodes)
}

// GetNodeOperatorsSimple implements unpagniatedCapabilityRegistry
func (e *ExtendedCapabilityRegistry) GetNodeOperatorsSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryNodeOperatorInfo, error) {
	return e.GetNodeOperators(opts, big.NewInt(0), MaxNOPs)
}

// GetDONsSimple implements unpagniatedCapabilityRegistry
func (e *ExtendedCapabilityRegistry) GetDONsSimple(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	return e.GetDONs(opts, big.NewInt(0), MaxDONs)
}

// GenerateCapabilityRegistryView generates a CapRegView from a CapabilitiesRegistry contract.
func GenerateCapabilityRegistryView(capReg *ExtendedCapabilityRegistry) (CapabilityRegistryView, error) {
	tv, err := types.NewContractMetaData(capReg, capReg.Address())
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	caps, err := capReg.GetCapabilitiesSimple(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var capViews []CapabilityView
	for _, capability := range caps {
		capView, capViewErr := NewCapabilityView(capability)
		if capViewErr != nil {
			return CapabilityRegistryView{}, fmt.Errorf("failed to create capability view for capability %s: %w", capability.CapabilityId, capViewErr)
		}
		capViews = append(capViews, capView)
	}
	donInfos, err := capReg.GetDONsSimple(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var donViews []DonView
	for _, donInfo := range donInfos {
		donView, donViewErr := NewDonView(donInfo)
		if donViewErr != nil {
			return CapabilityRegistryView{}, fmt.Errorf("failed to create don view for don %d: %w", donInfo.Id, donViewErr)
		}
		donViews = append(donViews, donView)
	}

	nodeInfos, err := capReg.GetNodesSimple(nil)
	if err != nil {
		return CapabilityRegistryView{}, err
	}
	var nodeViews []NodeView
	for _, nodeInfo := range nodeInfos {
		nodeViews = append(nodeViews, NewNodeView(nodeInfo))
	}

	nopInfos, err := capReg.GetNodeOperatorsSimple(nil)
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
		for _, capability := range v.Capabilities {
			if don.hasCapability(capability) {
				capabilities = append(capabilities, capability)
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
			return nil, fmt.Errorf("failed to decode signer: %w", err)
		}
		encryptionPubKey, err := hexTo32Bytes(node.EncryptionPublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode encryption public key: %w", err)
		}
		csaKey, err := hexTo32Bytes(node.CSAKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode csa key: %w", err)
		}

		nodesParams = append(nodesParams, capabilities_registry.CapabilitiesRegistryNodeParams{
			Signer:              signer,
			P2pId:               node.P2pID,
			EncryptionPublicKey: encryptionPubKey,
			CsaKey:              csaKey,
			NodeOperatorId:      node.NodeOperatorID,
			CapabilityIds:       node.CapabilityIDs,
		})
	}

	return nodesParams, nil
}

func (v *CapabilityRegistryView) CapabilitiesToCapabilitiesParams() ([]capabilities_registry.CapabilitiesRegistryCapability, error) {
	var capabilitiesParams []capabilities_registry.CapabilitiesRegistryCapability
	for _, capability := range v.Capabilities {
		fmt.Println("capInfo.Metadata:", capability.Metadata)
		metadataBytes, err := json.Marshal(capability.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal capability metadata for capability %s: %w", capability.ID, err)
		}
		capabilitiesParams = append(capabilitiesParams, capabilities_registry.CapabilitiesRegistryCapability{
			CapabilityId:          capability.ID,
			ConfigurationContract: capability.ConfigurationContract,
			Metadata:              metadataBytes,
		})
	}
	return capabilitiesParams, nil
}

func (v *CapabilityRegistryView) NopsToNopsParams() []capabilities_registry.CapabilitiesRegistryNodeOperatorInfo {
	var nopsParams []capabilities_registry.CapabilitiesRegistryNodeOperatorInfo
	for _, nop := range v.Nops {
		nopsParams = append(nopsParams, capabilities_registry.CapabilitiesRegistryNodeOperatorInfo{
			Admin: nop.Admin,
			Name:  nop.Name,
		})
	}
	return nopsParams
}

func (v *CapabilityRegistryView) CapabilityConfigToCapabilityConfigParams(don DonView) ([]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration, error) {
	var cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration
	for _, cfg := range don.CapabilityConfigurations {
		config := pkg.CapabilityConfig(cfg.Config)
		cfgBytes, err := config.MarshalProto()
		if err != nil {
			return nil, err
		}
		cfgs = append(cfgs, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: cfg.ID,
			Config:       cfgBytes,
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
	ConfigurationContract common.Address `json:"configuration_contract,omitempty"`
	Metadata              map[string]any `json:"metadata,omitempty"`
	IsDeprecated          bool           `json:"is_deprecated,omitempty"`
}

// NewCapabilityView creates a CapabilityView from a CapabilitiesRegistryCapabilityInfo.
func NewCapabilityView(capInfo capabilities_registry.CapabilitiesRegistryCapabilityInfo) (CapabilityView, error) {
	var metadata map[string]any
	// We have a weird case in which the metadata is just null chars (\x00) for a deprecated capability named `cap1` on eth sepolia.
	// First, find the first null byte.
	firstNull := bytes.IndexByte(capInfo.Metadata, 0)
	var cleanMetadata []byte
	if firstNull != -1 {
		cleanMetadata = capInfo.Metadata[:firstNull]
	} else {
		cleanMetadata = capInfo.Metadata
	}

	if len(cleanMetadata) > 0 {
		err := json.Unmarshal(cleanMetadata, &metadata)
		if err != nil {
			return CapabilityView{}, fmt.Errorf("failed to unmarshal capability metadata for capability %s: %w", capInfo.CapabilityId, err)
		}
	}

	return CapabilityView{
		ID:                    capInfo.CapabilityId,
		Metadata:              metadata,
		ConfigurationContract: capInfo.ConfigurationContract,
		IsDeprecated:          capInfo.IsDeprecated,
	}, nil
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
	ID               uint32         `json:"id"`
	Name             string         `json:"name"`
	ConfigCount      uint32         `json:"config_count"`
	F                uint8          `json:"f"`
	IsPublic         bool           `json:"is_public,omitempty"`
	AcceptsWorkflows bool           `json:"accepts_workflows,omitempty"`
	DONFamilies      []string       `json:"don_family,omitempty"`
	Config           map[string]any `json:"config,omitempty"`
}

// NewDonView creates a DonView from a CapabilitiesRegistryDONInfo.
func NewDonView(d capabilities_registry.CapabilitiesRegistryDONInfo) (DonView, error) {
	donCfg := pkg.CapabilityConfig{}
	err := donCfg.UnmarshalProto(d.Config)
	if err != nil {
		return DonView{}, fmt.Errorf("failed to unmarshal don config for don %d: %w", d.Id, err)
	}
	capCgfs, err := NewCapabilityConfigurations(d.CapabilityConfigurations)
	if err != nil {
		return DonView{}, fmt.Errorf("failed to create capability configurations for don %d: %w", d.Id, err)
	}
	return DonView{
		DonUniversalMetadata: DonUniversalMetadata{
			ID:               d.Id,
			Name:             d.Name,
			ConfigCount:      d.ConfigCount,
			F:                d.F,
			IsPublic:         d.IsPublic,
			AcceptsWorkflows: d.AcceptsWorkflows,
			DONFamilies:      d.DonFamilies,
			Config:           donCfg,
		},
		NodeP2PIds:               p2pIDs(d.NodeP2PIds),
		CapabilityConfigurations: capCgfs,
	}, nil
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
	ID     string         `json:"id"` // hex 32 bytes
	Config map[string]any `json:"config"`
}

// NewCapabilityConfigurations creates a list of CapabilitiesConfiguration from a list of CapabilitiesRegistryCapabilityConfiguration.
func NewCapabilityConfigurations(cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration) ([]CapabilitiesConfiguration, error) {
	var out []CapabilitiesConfiguration
	for _, cfg := range cfgs {
		capCfg := pkg.CapabilityConfig{}
		err := capCfg.UnmarshalProto(cfg.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal capability configuration for capability %s: %w", cfg.CapabilityId, err)
		}
		out = append(out, CapabilitiesConfiguration{
			ID:     cfg.CapabilityId,
			Config: capCfg,
		})
	}
	return out, nil
}

func (cc CapabilitiesConfiguration) Validate() error {
	id, err := hex.DecodeString(cc.ID)
	if err != nil {
		return errors.New("capability id must be hex encoded")
	}
	if len(id) != 32 {
		return errors.New("capability id must be 32 bytes")
	}
	x := pkg.CapabilityConfig(cc.Config)
	_, err = x.MarshalProto()
	if err != nil {
		return errors.New("config must be proto marshalable")
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
	P2pID               p2pkey.PeerID `json:"p2p_id"`
	CSAKey              string        `json:"csa_key"`               // hex 32 bytes
	EncryptionPublicKey string        `json:"encryption_public_key"` // hex 32 bytes

}

// NewNodeView creates a NodeView from a CapabilitiesRegistryNodeInfoProviderNodeInfo.
func NewNodeView(n capabilities_registry.INodeInfoProviderNodeInfo) NodeView {
	return NodeView{
		NodeUniversalMetadata: NodeUniversalMetadata{
			ConfigCount:         n.ConfigCount,
			WorkflowDONID:       n.WorkflowDONId,
			Signer:              hex.EncodeToString(n.Signer[:]),
			P2pID:               n.P2pId,
			EncryptionPublicKey: hex.EncodeToString(n.EncryptionPublicKey[:]),
			CSAKey:              hex.EncodeToString(n.CsaKey[:]),
		},
		NodeOperatorID:   n.NodeOperatorId,
		CapabilityIDs:    n.CapabilityIds,
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

func NewNopView(nop capabilities_registry.CapabilitiesRegistryNodeOperatorInfo) NopView {
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
		if n.NodeOperatorID == uint32(idx) { //nolint:gosec // G115
			return nop, nil
		}
	}
	return NopView{}, fmt.Errorf("could not find nop for node %d", n.NodeOperatorID)
}

func p2pIDs(rawIDs [][32]byte) []p2pkey.PeerID {
	var out []p2pkey.PeerID
	for _, id := range rawIDs {
		out = append(out, id)
	}
	return out
}

func (dv DonView) hasNode(node NodeView) bool {
	donID := big.NewInt(int64(dv.ID))
	return slices.ContainsFunc(node.CapabilityDONIDs, func(elem *big.Int) bool { return elem.Cmp(donID) == 0 }) || node.WorkflowDONID == dv.ID
}

func (dv DonView) hasCapability(candidate CapabilityView) bool {
	return slices.ContainsFunc(dv.CapabilityConfigurations, func(elem CapabilitiesConfiguration) bool { return elem.ID == candidate.ID })
}
