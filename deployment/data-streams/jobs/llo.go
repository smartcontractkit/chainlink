package jobs

import (
	"fmt"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type LLOJobSpec struct {
	Base

	ContractID                        string          `toml:"contractID"`
	TransmitterID                     string          `toml:"transmitterID,omitempty"`
	ForwardingAllowed                 *bool           `toml:"forwardingAllowed,omitempty"`
	P2PV2Bootstrappers                []string        `toml:"p2pv2Bootstrappers,omitempty"`
	OCRKeyBundleID                    *string         `toml:"ocrKeyBundleID,omitempty"`
	MaxTaskDuration                   TOMLDuration    `toml:"maxTaskDuration,omitempty"`
	ContractConfigTrackerPollInterval TOMLDuration    `toml:"contractConfigTrackerPollInterval,omitempty"`
	Relay                             RelayType       `toml:"relay,omitempty"`
	PluginType                        PluginType      `toml:"pluginType,omitempty"`
	RelayConfig                       RelayConfigLLO  `toml:"relayConfig"`
	PluginConfig                      PluginConfigLLO `toml:"pluginConfig"`
}

func (s *LLOJobSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s)
}

// RelayConfig is the configuration for the relay. This could change depending on the relay type.
type RelayConfigLLO struct {
	ChainID       string `toml:"chainID"`
	FromBlock     uint64 `toml:"fromBlock,omitempty"`
	LLOConfigMode string `toml:"lloConfigMode,omitempty"`
	LLODonID      uint64 `toml:"lloDonID,omitempty"`
}

type PluginConfigLLO struct {
	ChannelDefinitionsContractAddress   string            `toml:"channelDefinitionsContractAddress"`
	ChannelDefinitionsContractFromBlock uint64            `toml:"channelDefinitionsContractFromBlock"`
	DonID                               uint64            `toml:"donID"`
	Servers                             map[string]string `toml:"servers,inline"`
}

type TOMLDuration time.Duration

func (d *TOMLDuration) UnmarshalText(text []byte) error {
	var err error
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = TOMLDuration(duration)
	return nil
}

func (d *TOMLDuration) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%ds", time.Duration(*d)/time.Second)), nil
}
