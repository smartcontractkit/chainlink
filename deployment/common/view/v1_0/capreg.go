package v1_0

import (
	"bytes"
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
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
	capabilities_registry.CapabilitiesRegistryCapabilityInfo
}

type DonView struct {
	capabilities_registry.CapabilitiesRegistryDONInfo
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
		capViews = append(capViews, CapabilityView{capability})
	}
	donInfos, err := capReg.GetDONs(nil)
	if err != nil {
		return CapRegView{}, err
	}
	var donViews []DonView
	for _, donInfo := range donInfos {
		donViews = append(donViews, DonView{donInfo})
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
		if bytes.Equal(cfg.CapabilityId[:], cap.HashedId[:]) {
			isMember = true
			break
		}
	}
	return isMember
}
