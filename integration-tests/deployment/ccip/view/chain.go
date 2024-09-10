package view

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	DestinationChainSelectors []uint64                      `json:"destinationChainSelectors"`
	TokenAdminRegistry        map[string]TokenAdminRegistry `json:"tokenAdminRegistry"`
	NonceManager              map[string]NonceManager       `json:"nonceManager"`
	Router                    map[string]Router             `json:"router"`
}

func NewChain() Chain {
	return Chain{
		TokenAdminRegistry: make(map[string]TokenAdminRegistry),
		NonceManager:       make(map[string]NonceManager),
	}
}
