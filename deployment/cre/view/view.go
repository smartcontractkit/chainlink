package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v2_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v2_0"
)

type CREChainView struct {
	CapabilityRegistry map[string]common_v2_0.CapabilityRegistryView `json:"capability_registry,omitempty"`
	OCRContracts       map[string]OCR3ConfigView                     `json:"ocr_contracts,omitempty"`
}

func NewCREChainView() CREChainView {
	return CREChainView{
		CapabilityRegistry: make(map[string]common_v2_0.CapabilityRegistryView),
		OCRContracts:       make(map[string]OCR3ConfigView),
	}
}

type CREView struct {
	Chains map[string]CREChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView `json:"nops,omitempty"`
}

func (v CREView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias CREView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
