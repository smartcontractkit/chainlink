package pkg

import (
	"encoding/json"
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

// Uint64 wraps uint64 so that YAML/JSON fields can be populated from either a numeric
// literal or a quoted string (e.g. after environment-variable substitution).
// Only unmarshal methods are provided so TOML/JSON output remains numeric.
type Uint64 uint64

func (u *Uint64) UnmarshalText(data []byte) error {
	ui, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return err
	}

	*u = Uint64(ui)
	return nil
}

func (u *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return json.Unmarshal(data, (*uint64)(u))
	}

	switch data[0] {
	case '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		return u.UnmarshalText([]byte(s))
	default:
		var n uint64
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*u = Uint64(n)
		return nil
	}
}

func (u *Uint64) UnmarshalYAML(node *yaml.Node) error {
	ui, err := strconv.ParseUint(node.Value, 10, 64)
	if err != nil {
		return err
	}

	*u = Uint64(ui)
	return nil
}

type ChainSelector = Uint64
