package presenters

type ChainResource struct {
	JAID
	Enabled bool   `json:"enabled"`
	Config  string `json:"config"` // TOML
}
