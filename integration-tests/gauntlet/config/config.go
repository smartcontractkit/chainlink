package config

type Config struct {
	ChainName        string
	ChainID          string
	RPCUrl           string
	WSUrl            string
	ProgramAddresses *ProgramAddresses
	PrivateKey       string
}

type ProgramAddresses struct {
	OCR2             string
	AccessController string
	Store            string
}

func DevnetConfig() *Config {
	return &Config{
		ChainName: "solana",
		ChainID:   "devnet",
		// Will be overridden if set in toml
		RPCUrl: "https://api.devnet.solana.com",
		WSUrl:  "wss://api.devnet.solana.com/",
	}
}

func LocalNetConfig() *Config {
	return &Config{
		ChainName: "solana",
		ChainID:   "localnet",
		// Will be overridden if set in toml
		RPCUrl: "http://sol:8899",
		WSUrl:  "ws://sol:8900",
		ProgramAddresses: &ProgramAddresses{
			OCR2:             "E3j24rx12SyVsG6quKuZPbQqZPkhAUCh8Uek4XrKYD2x",
			AccessController: "2ckhep7Mvy1dExenBqpcdevhRu7CLuuctMcx7G9mWEvo",
			Store:            "9kRNTZmoZSiTBuXC62dzK9E7gC7huYgcmRRhYv3i4osC",
		},
	}
}
