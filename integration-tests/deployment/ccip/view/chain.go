package view

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	DestinationChainSelectors []uint64                      `json:"destinationChainSelectors,omitempty"`
	SourceChainSelectors      []uint64                      `json:"sourceChainSelectors,omitempty"`
	TokenAdminRegistry        map[string]TokenAdminRegistry `json:"tokenAdminRegistry,omitempty"`
	NonceManager              map[string]NonceManager       `json:"nonceManager,omitempty"`
	Router                    map[string]Router             `json:"router,omitempty"`
	RMN                       map[string]RMN                `json:"rmn,omitempty"`
	OnRamp                    map[string]OnRamp             `json:"onRamp,omitempty"`
	OffRamp                   map[string]OffRamp            `json:"offRamp,omitempty"`
}

func NewChain() Chain {
	return Chain{
		DestinationChainSelectors: make([]uint64, 0),
		SourceChainSelectors:      make([]uint64, 0),
		TokenAdminRegistry:        make(map[string]TokenAdminRegistry),
		NonceManager:              make(map[string]NonceManager),
		Router:                    make(map[string]Router),
		RMN:                       make(map[string]RMN),
		OnRamp:                    make(map[string]OnRamp),
		OffRamp:                   make(map[string]OffRamp),
	}
}
