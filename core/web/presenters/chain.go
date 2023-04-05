package presenters

type ChainResource struct {
	JAID
	Enabled bool   `json:"enabled"`
	Config  string `json:"config"` // TOML
}

type NodeResource struct {
	JAID
	ChainID string `json:"chainID"`
	Name    string `json:"name"`
	Config  string `json:"config"` // TOML
	State   string `json:"state"`
}
