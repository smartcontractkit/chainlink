package jobs

import (
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
)

var _ JobSpec = &BootstrapSpec{}

type BootstrapSpec struct {
	Base

	ContractID  string      `toml:"contractID"`
	DonID       uint64      `toml:"donID,omitempty"`
	Relay       RelayType   `toml:"relay"`
	RelayConfig RelayConfig `toml:"relayConfig"`
}

// PluginType is the type of plugin set in the TOML.
type PluginType string

const (
	PluginTypeLLO PluginType = "llo"
)

// RelayType is the type of relay set in the TOML.
type RelayType string

const (
	RelayTypeEVM    RelayType = "evm"
	RelayTypeSolana RelayType = "solana"
)

// RelayConfig is the configuration for the relay. This could change depending on the relay type.
type RelayConfig struct {
	ChainID   string `toml:"chainID"`
	FromBlock uint64 `toml:"fromBlock,omitempty"`
}

func NewBootstrapSpec(contractID string, donID uint64, relay RelayType, relayConfig RelayConfig, externalJobID uuid.UUID) *BootstrapSpec {
	if externalJobID == uuid.Nil {
		externalJobID = uuid.New()
	}
	return &BootstrapSpec{
		Base: Base{
			Name:          "bootstrap",
			Type:          JobSpecTypeBootstrap,
			SchemaVersion: 1,
			ExternalJobID: externalJobID,
		},
		ContractID:  contractID,
		DonID:       donID,
		Relay:       relay,
		RelayConfig: relayConfig,
	}
}

func (b *BootstrapSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(b)
}
