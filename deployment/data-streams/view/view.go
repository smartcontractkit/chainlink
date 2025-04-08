package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
)

type ChainView struct {
	Configurator  map[string]v0_5.ConfiguratorView  `json:"configurator,omitempty"`
	RewardManager map[string]v0_5.RewardManagerView `json:"rewardManager,omitempty"`
	Verifier      map[string]v0_5.VerifierView      `json:"verifier,omitempty"`
	VerifierProxy map[string]v0_5.VerifierProxyView `json:"verifierProxy,omitempty"`
}

func NewChain() ChainView {
	return ChainView{
		Configurator:  make(map[string]v0_5.ConfiguratorView),
		RewardManager: make(map[string]v0_5.RewardManagerView),
		Verifier:      make(map[string]v0_5.VerifierView),
		VerifierProxy: make(map[string]v0_5.VerifierProxyView),
	}
}

type DataStreamsView struct {
	Chains map[string]ChainView `json:"chains,omitempty"`
}

func (v DataStreamsView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls.
	type Alias DataStreamsView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
