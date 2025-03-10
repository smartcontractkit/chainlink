package types

type NodeEthKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
	ChainID  int    `toml:"ID"`
}

type NodeP2PKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
}

type NodeEthKeyWrapper struct {
	EthKeys []NodeEthKey `toml:"Keys"`
}

type NodeSecret struct {
	EthKeys NodeEthKeyWrapper `toml:"EVM"`
	P2PKey  NodeP2PKey        `toml:"P2PKey"`
	// Add more fields as needed to reflect 'Secrets' struct from /core/config/toml/types.go
	// We can't use the original struct, because it's using custom types that serlialize secrets to 'xxxxx'
}
