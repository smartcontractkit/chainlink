package node_config

type P2P struct {
	*P2Pv1 `toml:"V1"`
	*P2Pv2 `toml:"V2"`
}

type P2Pv1 struct {
	Enabled    bool   `toml:"Enabled"`
	ListenIP   string `toml:"ListenIP"`
	ListenPort string `toml:"ListenPort"`
}

type P2Pv2 struct {
	Enabled         bool     `toml:"Enabled"`
	ListenAddresses []string `toml:"ListenAddresses"`
}

func WithP2Pv1() NodeConfigOpt {
	return func(n *NodeConfig) {
		// default params
		n.P2P = P2P{
			P2Pv1: &P2Pv1{
				Enabled:    true,
				ListenIP:   "0.0.0.0",
				ListenPort: "6990",
			},
		}
	}
}

func WithP2Pv2() NodeConfigOpt {
	return func(n *NodeConfig) {
		// default params
		n.P2P = P2P{
			P2Pv2: &P2Pv2{
				Enabled:         true,
				ListenAddresses: []string{"0.0.0.0:6690"},
			},
		}
	}
}
