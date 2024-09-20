package view

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_6"
)

type ChainView struct {
	TokenAdminRegistry map[string]v1_5.TokenAdminRegistryView `json:"tokenAdminRegistry,omitempty"`
	FeeQuoter          map[string]v1_6.FeeQuoterView          `json:"feeQuoter,omitempty"`
	NonceManager       map[string]v1_6.NonceManagerView       `json:"nonceManager,omitempty"`
	Router             map[string]v1_2.RouterView             `json:"router,omitempty"`
	RMN                map[string]v1_6.RMNRemoteView          `json:"rmn,omitempty"`
	OnRamp             map[string]v1_6.OnRampView             `json:"onRamp,omitempty"`
	OffRamp            map[string]v1_6.OffRampView            `json:"offRamp,omitempty"`
}

func NewChain() ChainView {
	return ChainView{
		TokenAdminRegistry: make(map[string]v1_5.TokenAdminRegistryView),
		NonceManager:       make(map[string]v1_6.NonceManagerView),
		Router:             make(map[string]v1_2.RouterView),
		RMN:                make(map[string]v1_6.RMNRemoteView),
		OnRamp:             make(map[string]v1_6.OnRampView),
		OffRamp:            make(map[string]v1_6.OffRampView),
		FeeQuoter:          make(map[string]v1_6.FeeQuoterView),
	}
}
