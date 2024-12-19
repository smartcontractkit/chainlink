package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
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
	FeeQuoter    map[string]v1_6.FeeQuoterView    `json:"feeQuoter,omitempty"`
	NonceManager map[string]v1_6.NonceManagerView `json:"nonceManager,omitempty"`
	RMNHome      map[string]v1_6.RMNHomeView      `json:"rmnHome,omitempty"`
	RMN          map[string]v1_6.RMNRemoteView    `json:"rmn,omitempty"`
	OnRamp       map[string]v1_6.OnRampView       `json:"onRamp,omitempty"`
	OffRamp      map[string]v1_6.OffRampView      `json:"offRamp,omitempty"`
	// TODO: Perhaps restrict to one CCIPHome/CR? Shouldn't
	// be more than one per env.
	CCIPHome           map[string]v1_6.CCIPHomeView                  `json:"ccipHome,omitempty"`
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	MCMSWithTimelock   common_v1_0.MCMSWithTimelockView              `json:"mcmsWithTimelock,omitempty"`
	LinkToken          common_v1_0.LinkTokenView                     `json:"linkToken,omitempty"`
	StaticLinkToken    common_v1_0.StaticLinkTokenView               `json:"staticLinkToken,omitempty"`
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
		RMNHome:            make(map[string]v1_6.RMNHomeView),
		RMN:                make(map[string]v1_6.RMNRemoteView),
		OnRamp:             make(map[string]v1_6.OnRampView),
		OffRamp:            make(map[string]v1_6.OffRampView),
		CapabilityRegistry: make(map[string]common_v1_0.CapabilityRegistryView),
		CCIPHome:           make(map[string]v1_6.CCIPHomeView),
		MCMSWithTimelock:   common_v1_0.MCMSWithTimelockView{},
		LinkToken:          common_v1_0.LinkTokenView{},
		StaticLinkToken:    common_v1_0.StaticLinkTokenView{},
	}
}

type CCIPView struct {
	Chains map[string]ChainView    `json:"chains,omitempty"`
	Nops   map[string]view.NopView `json:"nops,omitempty"`
}

func (v CCIPView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias CCIPView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
