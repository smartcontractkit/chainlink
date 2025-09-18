package pkg

type OracleFactory struct {
	Enabled                bool
	BootstrapPeers         []string
	OCRContractAddress     string
	OCRKeyBundleID         string
	ChainID                string
	TransmitterID          string
	OnchainSigningStrategy OnchainSigningStrategy
}

type OnchainSigningStrategy struct {
	StrategyName string
	Config       map[string]string
}

type OracleFactoryConfig struct {
	Enabled            bool     `toml:"enabled"`
	BootstrapPeers     []string `toml:"bootstrap_peers"`      // e.g.,["12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690"]
	OCRContractAddress string   `toml:"ocr_contract_address"` // e.g., 0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6
	ChainID            string   `toml:"chain_id"`             // e.g., "31337"
	Network            string   `toml:"network"`              // e.g., "evm"
}
