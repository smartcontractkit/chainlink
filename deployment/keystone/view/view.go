package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapRegView `json:"capabilityRegistry,omitempty"`
	// TODO forwarders etc
}

type KeystoneView struct {
	Chains map[string]KeystoneChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView      `json:"nops,omitempty"`
}

func (v KeystoneView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}
