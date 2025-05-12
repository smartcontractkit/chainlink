package view

import (
	"encoding/json"
	"sync"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

// ChainView is a json-persistable structure that represents chain state. Store all versions of CCIP contracts
// CCIP observability relies on ChainView. When making changes that makes final json backward incompatible, warn CCIP observability team
type ChainView struct {
	ChainSelector uint64 `json:"chainSelector,omitempty"`
	ChainID       string `json:"chainID,omitempty"`
	// v1.0
	RMNProxy map[string]v1_0.RMNProxyView `json:"rmnProxy,omitempty"`
	// v1.2
	Router map[string]v1_2.RouterView `json:"router,omitempty"`
	// v1.5
	TokenAdminRegistry map[string]v1_5.TokenAdminRegistryView `json:"tokenAdminRegistry,omitempty"`
	TokenPoolFactory   map[string]v1_5_1.TokenPoolFactoryView `json:"tokenPoolFactory,omitempty"`
	RegistryModules    map[string]shared.RegistryModulesView  `json:"registryModules,omitempty"`
	TokenPools         map[string]map[string]v1_5_1.PoolView  `json:"poolByTokens,omitempty"` // TokenSymbol => TokenPool Address => PoolView
	CommitStore        map[string]v1_5.CommitStoreView        `json:"commitStore,omitempty"`
	PriceRegistry      map[string]v1_2.PriceRegistryView      `json:"priceRegistry,omitempty"`
	EVM2EVMOnRamp      map[string]v1_5.OnRampView             `json:"evm2evmOnRamp,omitempty"`
	EVM2EVMOffRamp     map[string]v1_5.OffRampView            `json:"evm2evmOffRamp,omitempty"`
	RMN                map[string]v1_5.RMNView                `json:"rmn,omitempty"`

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

	UpdateMu *sync.Mutex `json:"-"`
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
		TokenPoolFactory:   make(map[string]v1_5_1.TokenPoolFactoryView),
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
		UpdateMu:           &sync.Mutex{},
	}
}

type SolChainView struct {
	ChainSelector uint64 `json:"chainSelector,omitempty"`
	ChainID       string `json:"chainID,omitempty"`
	// v1.6
	FeeQuoter map[string]solana.FeeQuoterView `json:"feeQuoter,omitempty"`
	Router    map[string]solana.RouterView    `json:"router,omitempty"`
	OffRamp   map[string]solana.OffRampView   `json:"offRamp,omitempty"`
	RMNRemote map[string]solana.RMNRemoteView `json:"rmnRemote,omitempty"`
	TokenPool map[string]solana.TokenPoolView `json:"tokenPool,omitempty"`
	LinkToken solana.TokenView                `json:"linkToken,omitempty"`
	Tokens    map[string]solana.TokenView     `json:"tokens,omitempty"`
}

func NewSolChain() SolChainView {
	return SolChainView{
		FeeQuoter: make(map[string]solana.FeeQuoterView),
		Router:    make(map[string]solana.RouterView),
		OffRamp:   make(map[string]solana.OffRampView),
		RMNRemote: make(map[string]solana.RMNRemoteView),
		TokenPool: make(map[string]solana.TokenPoolView),
		Tokens:    make(map[string]solana.TokenView),
	}
}

func (v *ChainView) UpdateTokenPool(tokenSymbol string, tokenPoolAddress string, poolView v1_5_1.PoolView) {
	v.UpdateMu.Lock()
	defer v.UpdateMu.Unlock()
	v.TokenPools = helpers.AddValueToNestedMap(v.TokenPools, tokenSymbol, tokenPoolAddress, poolView)
}

func (v *ChainView) UpdateRegistryModuleView(registryModuleAddress string, registryModuleView shared.RegistryModulesView) {
	v.UpdateMu.Lock()
	defer v.UpdateMu.Unlock()
	if v.RegistryModules == nil {
		v.RegistryModules = make(map[string]shared.RegistryModulesView)
	}
	v.RegistryModules[registryModuleAddress] = registryModuleView
}

type CCIPView struct {
	Chains    map[string]ChainView    `json:"chains,omitempty"`
	SolChains map[string]SolChainView `json:"solChains,omitempty"`
	Nops      map[string]view.NopView `json:"nops,omitempty"`
}

func (v CCIPView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias CCIPView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
