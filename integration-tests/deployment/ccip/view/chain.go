package view

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/v1_6"
)

type ChainView struct {
	// v1.0
	RMNProxy map[string]v1_0.RMNProxyView `json:"rmnProxy,omitempty"`
	// v1.2
	Router map[string]v1_2.RouterView `json:"router,omitempty"`
	// v1.5
	TokenAdminRegistry map[string]v1_5.TokenAdminRegistryView `json:"tokenAdminRegistry,omitempty"`
	CommitStore        map[string]v1_5.CommitStoreView        `json:"commitStore,omitempty"`
	// v1.6
	FeeQuoter          map[string]v1_6.FeeQuoterView    `json:"feeQuoter,omitempty"`
	NonceManager       map[string]v1_6.NonceManagerView `json:"nonceManager,omitempty"`
	RMN                map[string]v1_6.RMNRemoteView    `json:"rmn,omitempty"`
	OnRamp             map[string]v1_6.OnRampView       `json:"onRamp,omitempty"`
	OffRamp            map[string]v1_6.OffRampView      `json:"offRamp,omitempty"`
	CapabilityRegistry map[string]v1_6.CapRegView       `json:"capabilityRegistry,omitempty"`
}

func NewChain() ChainView {
	return ChainView{
		// v1.0
		RMNProxy: make(map[string]v1_0.RMNProxyView),
		// v1.2
		Router: make(map[string]v1_2.RouterView),
		// v1.5
		TokenAdminRegistry: make(map[string]v1_5.TokenAdminRegistryView),
		CommitStore:        make(map[string]v1_5.CommitStoreView),
		// v1.6
		FeeQuoter:          make(map[string]v1_6.FeeQuoterView),
		NonceManager:       make(map[string]v1_6.NonceManagerView),
		RMN:                make(map[string]v1_6.RMNRemoteView),
		OnRamp:             make(map[string]v1_6.OnRampView),
		OffRamp:            make(map[string]v1_6.OffRampView),
		CapabilityRegistry: make(map[string]v1_6.CapRegView),
	}
}
