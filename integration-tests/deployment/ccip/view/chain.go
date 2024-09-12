package view

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_2"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_5"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_6"
)

type Chain struct {
	DestinationChainSelectors []uint64 `json:"destinationChainSelectors,omitempty"`
	// TODO - populate supportedTokensByDestination
	SupportedTokensByDestination map[uint64][]string                `json:"supportedTokensByDestination,omitempty"`
	TokenAdminRegistry           map[string]v1_5.TokenAdminRegistry `json:"tokenAdminRegistry,omitempty"`
	FeeQuoter                    map[string]v1_6.FeeQuoter          `json:"feeQuoter,omitempty"`
	NonceManager                 map[string]v1_6.NonceManager       `json:"nonceManager,omitempty"`
	Router                       map[string]v1_2.Router             `json:"router,omitempty"`
	RMN                          map[string]v1_6.RMN                `json:"rmn,omitempty"`
	OnRamp                       map[string]v1_6.OnRamp             `json:"onRamp,omitempty"`
}

func NewChain() Chain {
	return Chain{
		DestinationChainSelectors: make([]uint64, 0),
		TokenAdminRegistry:        make(map[string]v1_5.TokenAdminRegistry),
		NonceManager:              make(map[string]v1_6.NonceManager),
		Router:                    make(map[string]v1_2.Router),
		RMN:                       make(map[string]v1_6.RMN),
		OnRamp:                    make(map[string]v1_6.OnRamp),
	}
}
