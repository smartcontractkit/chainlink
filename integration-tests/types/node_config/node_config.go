package node_config

type NodeConfig struct {
	RootDir   string `toml:"RootDir"`
	EVM       []EVM  `toml:"EVM"`
	EVMNodes  []Node `toml:"EVM.Nodes"`
	P2P       `toml:"P2P"`
	*OCR      `toml:"OCR"`
	*OCR2     `toml:"OCR2"`
	Feature   `toml:"Feature"`
	WebServer `toml:"WebServer"`
	Log       `toml:"Log"`
	Database  `toml:"Database"`
}

type NodeConfigOpt = func(c *NodeConfig)

func NewNodeConfig(opts ...NodeConfigOpt) NodeConfig {
	// Default options
	nodeConfOpts := NodeConfig{
		RootDir: "/home/chainlink",
		WebServer: WebServer{
			AllowOrigins:   "*",
			HTTPPort:       6688,
			SecureCookies:  false,
			SessionTimeout: "999h0m0s",
			WebServerTLS: &WebServerTLS{
				HTTPSPort: 0,
			},
			WebServerRateLimit: &WebServerRateLimit{
				Authenticated:   2000,
				Unauthenticated: 100,
			},
		},
		Log: Log{
			Level:       "debug",
			JSONConsole: true,
		},
		Database: Database{
			MaxIdleConns:     20,
			MaxOpenConns:     40,
			MigrateOnStartup: true,
		},
		Feature: Feature{
			LogPoller:    true,
			FeedsManager: true,
			UICSAKeys:    true,
		},
	}
	for _, opt := range opts {
		opt(&nodeConfOpts)
	}
	return nodeConfOpts
}
