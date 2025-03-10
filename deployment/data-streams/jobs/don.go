package jobs

import (
	"time"

	"github.com/pelletier/go-toml/v2"
)

type DonJobSpec struct {
	Base

	ContractID                        string          `toml:"contractID"`
	TransmitterID                     string          `toml:"transmitterID,omitempty"`
	ForwardingAllowed                 *bool           `toml:"forwardingAllowed,omitempty"`
	P2PV2Bootstrappers                []string        `toml:"p2pv2Bootstrappers,omitempty"`
	OCRKeyBundleID                    *string         `toml:"ocrKeyBundleID,omitempty"`
	MaxTaskDuration                   time.Duration   `toml:"maxTaskDuration,omitempty"`
	ContractConfigTrackerPollInterval time.Duration   `toml:"contractConfigTrackerPollInterval,omitempty"`
	Relay                             RelayType       `toml:"relay,omitempty"`
	PluginType                        string          `toml:"pluginType,omitempty"`
	RelayConfig                       RelayConfigDon  `toml:"relayConfig"`
	PluginConfig                      PluginConfigDon `toml:"pluginConfig"`
}

// RelayConfig is the configuration for the relay. This could change depending on the relay type.
type RelayConfigDon struct {
	ChainID       string `toml:"chainID"`
	FromBlock     uint64 `toml:"fromBlock,omitempty"`
	LLOConfigMode string `toml:"lloConfigMode,omitempty"`
	LLODonID      int64  `toml:"lloDonID,omitempty"`
}

type PluginConfigDon struct {
	ChannelDefinitionsContractAddress   string            `toml:"channelDefinitionsContractAddress"`
	ChannelDefinitionsContractFromBlock uint64            `toml:"channelDefinitionsContractFromBlock"`
	DonID                               int64             `toml:"donID"`
	Servers                             map[string]string `toml:"servers,inline"`
}

func (s *DonJobSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s)
}
