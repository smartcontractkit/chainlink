package view

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	TokenAdminRegistry map[string]TokenAdminRegistry `json:"tokenAdminRegistry"`
	NonceManager       map[string]NonceManager       `json:"nonceManager"`
}

func NewChain() Chain {
	return Chain{
		TokenAdminRegistry: make(map[string]TokenAdminRegistry),
		NonceManager:       make(map[string]NonceManager),
	}
}
