package v1_0

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// CapRegView denotes a view of the capabilities registry contract.
// Note that the contract itself is 1.0.0 versioned, but we're releasing it first
// as part of 1.6 for CCIP.
type CapRegView struct {
	types.ContractMetaData
	Capabilities []CapabilityView `json:"capabilities,omitempty"`
	Nodes        []NodeView       `json:"nodes,omitempty"`
	Nops         []NopView        `json:"nops,omitempty"`
	Dons         []DonView        `json:"dons,omitempty"`
}

type CapabilityView struct {
	CapabilityId          string // hex
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address `json:"omitempty"`
	IsDeprecated          bool           `json:"omitempty"`
}

func NewCapabilityView(capInfo capabilities_registry.CapabilitiesRegistryCapabilityInfo) CapabilityView {
	return CapabilityView{
		CapabilityId:          hex.EncodeToString(capInfo.HashedId[:]),
		LabelledName:          capInfo.LabelledName,
		Version:               capInfo.Version,
		CapabilityType:        capInfo.CapabilityType,
		ResponseType:          capInfo.ResponseType,
		ConfigurationContract: capInfo.ConfigurationContract,
		IsDeprecated:          capInfo.IsDeprecated,
	}
}

func (cv CapabilityView) Validate() error {
	id, err := hex.DecodeString(cv.CapabilityId)
	if err != nil {
		return err
	}
	if len(id) != 32 {
		return errors.New("capability id must be 32 bytes")
	}
	return nil
}

type DonView struct {
	Id                       uint32
	ConfigCount              uint32
	F                        uint8
	IsPublic                 bool
	AcceptsWorkflows         bool
	NodeP2PIds               []p2pkey.PeerID
	CapabilityConfigurations []CapabilitiesConfiguration
}

type CapabilitiesConfiguration struct {
	CapabilityId string // hex 32 bytes
	Config       string // hex
}

func NewDonView(d capabilities_registry.CapabilitiesRegistryDONInfo) DonView {
	return DonView{
		Id:                       d.Id,
		ConfigCount:              d.ConfigCount,
		F:                        d.F,
		IsPublic:                 d.IsPublic,
		AcceptsWorkflows:         d.AcceptsWorkflows,
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

func p2pIds(rawIds [][32]byte) []p2pkey.PeerID {
	var out []p2pkey.PeerID
	for _, id := range rawIds {
		out = append(out, p2pkey.PeerID(id))
	}
	return out
}

func NewCapabilityConfigurations(cfgs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration) []CapabilitiesConfiguration {
	var out []CapabilitiesConfiguration
	for _, cfg := range cfgs {
		out = append(out, CapabilitiesConfiguration{
			CapabilityId: hex.EncodeToString(cfg.CapabilityId[:]),
			Config:       hex.EncodeToString(cfg.Config),
		})
	}
	return out
}

func (cc CapabilitiesConfiguration) Validate() error {
	id, err := hex.DecodeString(cc.CapabilityId)
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

type NodeView struct {
	capabilities_registry.INodeInfoProviderNodeInfo
}

type NopView struct {
	capabilities_registry.CapabilitiesRegistryNodeOperator
}

func GenerateCapRegView(capReg *capabilities_registry.CapabilitiesRegistry) (CapRegView, error) {
	tv, err := types.NewContractMetaData(capReg, capReg.Address())
	if err != nil {
		return CapRegView{}, err
	}
	caps, err := capReg.GetCapabilities(nil)
	if err != nil {
		return CapRegView{}, err
	}
	var capViews []CapabilityView
	for _, capability := range caps {
		capViews = append(capViews, NewCapabilityView(capability))
	}
	donInfos, err := capReg.GetDONs(nil)
	if err != nil {
		return CapRegView{}, err
	}
	var donViews []DonView
	for _, donInfo := range donInfos {
		donViews = append(donViews, NewDonView(donInfo))
	}

	nodeInfos, err := capReg.GetNodes(nil)
	if err != nil {
		return CapRegView{}, err
	}
	var nodeViews []NodeView
	for _, nodeInfo := range nodeInfos {
		nodeViews = append(nodeViews, NodeView{nodeInfo})
	}

	nopInfos, err := capReg.GetNodeOperators(nil)
	if err != nil {
		return CapRegView{}, err
	}
	var nopViews []NopView
	for _, nopInfo := range nopInfos {
		nopViews = append(nopViews, NopView{nopInfo})
	}

	return CapRegView{
		ContractMetaData: tv,
		Capabilities:     capViews,
		Dons:             donViews,
		Nodes:            nodeViews,
		Nops:             nopViews,
	}, nil
}

type DonCapabilities struct {
	Don          DonView
	Nodes        []NodeView
	Capabilities []CapabilityView
}

func (v CapRegView) Denormalize() ([]DonCapabilities, error) {
	var out []DonCapabilities
	for _, don := range v.Dons {
		var nodes []NodeView
		for _, node := range v.Nodes {
			if nodeInDon(node, don) {
				nodes = append(nodes, node)
			}
		}
		var capabilities []CapabilityView
		for _, cap := range v.Capabilities {
			if capInDon(cap, don) {
				capabilities = append(capabilities, cap)
			}
		}
		out = append(out, DonCapabilities{
			Don:          don,
			Nodes:        nodes,
			Capabilities: capabilities,
		})
	}
	return out, nil
}

func nodeInDon(node NodeView, don DonView) bool {
	donId := big.NewInt(int64(don.Id))
	isMember := false
	for _, x := range node.CapabilitiesDONIds {
		if x.Cmp(donId) == 0 {
			isMember = true
			break
		}
	}
	return isMember
}

func capInDon(cap CapabilityView, don DonView) bool {
	isMember := false
	for _, cfg := range don.CapabilityConfigurations {
		if cfg.CapabilityId == cap.CapabilityId {
			isMember = true
			break
		}
	}
	return isMember
}
