package view

type Chain struct {
	DestinationChainSelectors []uint64 `json:"destinationChainSelectors,omitempty"`
	// TODO - populate supportedTokensByDestination
	SupportedTokensByDestination map[uint64][]string           `json:"supportedTokensByDestination,omitempty"`
	TokenAdminRegistry           map[string]TokenAdminRegistry `json:"tokenAdminRegistry,omitempty"`
	FeeQuoter                    map[string]FeeQuoter          `json:"feeQuoter,omitempty"`
	NonceManager                 map[string]NonceManager       `json:"nonceManager,omitempty"`
	Router                       map[string]Router             `json:"router,omitempty"`
	RMN                          map[string]RMN                `json:"rmn,omitempty"`
	OnRamp                       map[string]OnRamp             `json:"onRamp,omitempty"`
}

func NewChain() Chain {
	return Chain{
		DestinationChainSelectors: make([]uint64, 0),
		TokenAdminRegistry:        make(map[string]TokenAdminRegistry),
		NonceManager:              make(map[string]NonceManager),
		Router:                    make(map[string]Router),
		RMN:                       make(map[string]RMN),
		OnRamp:                    make(map[string]OnRamp),
	}
}
