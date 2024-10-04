package v1_6

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

// CapRegView denotes a view of the capabilities registry contract.
// Note that the contract itself is 1.0.0 versioned, but we're releasing it first
// as part of 1.6.
type CapRegView struct {
	types.ContractMetaData
	Capabilities []CapabilityView `json:"capabilities,omitempty"`
}

type CapabilityView struct {
	LabelledName   string         `json:"labelledName"`
	Version        string         `json:"version"`
	ConfigContract common.Address `json:"configContract"`
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
		capViews = append(capViews, CapabilityView{
			LabelledName:   capability.LabelledName,
			Version:        capability.Version,
			ConfigContract: capability.ConfigurationContract,
		})
	}
	return CapRegView{
		ContractMetaData: tv,
		Capabilities:     capViews,
	}, nil
}
