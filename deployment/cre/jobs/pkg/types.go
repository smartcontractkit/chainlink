package pkg

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

type OracleFactory struct {
	Enabled                bool                   `yaml:"enabled"`
	BootstrapPeers         []string               `yaml:"bootstrapPeers"`
	OCRContractAddress     string                 `yaml:"ocrContractAddress"`
	OCRKeyBundleID         string                 `yaml:"ocrKeyBundleID"`
	ChainID                string                 `yaml:"chainID"`
	TransmitterID          string                 `yaml:"transmitterID"`
	OnchainSigningStrategy OnchainSigningStrategy `yaml:"onchainSigningStrategy"`
}

type OnchainSigningStrategy struct {
	StrategyName string            `yaml:"strategyName"`
	Config       map[string]string `yaml:"config"`
}

type OracleFactoryConfig struct {
	Enabled            bool     `toml:"enabled"`
	BootstrapPeers     []string `toml:"bootstrap_peers"`      // e.g.,["12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690"]
	OCRContractAddress string   `toml:"ocr_contract_address"` // e.g., 0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6
	ChainID            string   `toml:"chain_id"`             // e.g., "31337"
	Network            string   `toml:"network"`              // e.g., "evm"
}

// Int wraps int so that YAML fields can be populated from either a numeric
// literal or a quoted string (e.g. after environment-variable substitution).
type Int int

func (i *Int) UnmarshalYAML(node *yaml.Node) error {
	v, err := strconv.Atoi(node.Value)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}

func (i Int) MarshalYAML() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

type ChainSelector uint64

func (cs *ChainSelector) UnmarshalText(data []byte) error {
	ui, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return err
	}

	*cs = ChainSelector(ui)
	return nil
}

func (cs ChainSelector) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(cs), 10)), nil
}

func (cs *ChainSelector) UnmarshalYAML(node *yaml.Node) error {
	ui, err := strconv.ParseUint(node.Value, 10, 64)
	if err != nil {
		return err
	}

	*cs = ChainSelector(ui)
	return nil
}

func (cs ChainSelector) MarshalYAML() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(cs), 10)), nil
}
