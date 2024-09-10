package view

type Chain struct {
	// TODO: this will have to be versioned for getting state during upgrades.
	TokenAdminRegistry TokenAdminRegistry `json:"tokenAdminRegistry"`
	NonceManager       NonceManager       `json:"nonceManager"`
}
