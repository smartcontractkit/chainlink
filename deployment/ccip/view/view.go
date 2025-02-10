package view

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
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
	TokenAdminRegistry   map[string]v1_5.TokenAdminRegistryView                `json:"tokenAdminRegistry,omitempty"`
	BurnMintTokenPool    map[string]map[string]v1_5_1.TokenPoolView            `json:"burnMintTokenPool,omitempty"`
	LockReleaseTokenPool map[string]map[string]v1_5_1.LockReleaseTokenPoolView `json:"lockReleaseTokenPool,omitempty"`
	USDCTokenPool        map[string]map[string]v1_5_1.USDCTokenPoolView        `json:"usdcTokenPool,omitempty"`
	CommitStore          map[string]v1_5.CommitStoreView                       `json:"commitStore,omitempty"`
	PriceRegistry        map[string]v1_2.PriceRegistryView                     `json:"priceRegistry,omitempty"`
	EVM2EVMOnRamp        map[string]v1_5.OnRampView                            `json:"evm2evmOnRamp,omitempty"`
	EVM2EVMOffRamp       map[string]v1_5.OffRampView                           `json:"evm2evmOffRamp,omitempty"`
	RMN                  map[string]v1_5.RMNView                               `json:"rmn,omitempty"`

	// v1.6
	FeeQuoter    map[string]v1_6.FeeQuoterView    `json:"feeQuoter,omitempty"`
	NonceManager map[string]v1_6.NonceManagerView `json:"nonceManager,omitempty"`
	RMNRemote    map[string]v1_6.RMNRemoteView    `json:"rmnRemote,omitempty"`
	RMNHome      map[string]v1_6.RMNHomeView      `json:"rmnHome,omitempty"`
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
		Router:        make(map[string]v1_2.RouterView),
		PriceRegistry: make(map[string]v1_2.PriceRegistryView),
		// v1.5
		TokenAdminRegistry: make(map[string]v1_5.TokenAdminRegistryView),
		CommitStore:        make(map[string]v1_5.CommitStoreView),
		EVM2EVMOnRamp:      make(map[string]v1_5.OnRampView),
		EVM2EVMOffRamp:     make(map[string]v1_5.OffRampView),
		RMN:                make(map[string]v1_5.RMNView),
		// v1.6
		FeeQuoter:          make(map[string]v1_6.FeeQuoterView),
		NonceManager:       make(map[string]v1_6.NonceManagerView),
		RMNRemote:          make(map[string]v1_6.RMNRemoteView),
		RMNHome:            make(map[string]v1_6.RMNHomeView),
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
