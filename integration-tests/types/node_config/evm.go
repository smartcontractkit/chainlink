package node_config

type EVM struct {
	ChainID            string `toml:"ChainID"`
	AutoCreateKey      bool   `toml:"AutoCreateKey"`
	FinalityDepth      int    `toml:"finalityDepth"`
	MinContractPayment string `toml:"MinContractPayment"`
	Nodes              []Node `toml:"Nodes"`
}

type Node struct {
	WSURL    string `toml:"WSURL"`
	HTTPURL  string `toml:"HTTPURL"`
	Name     string `toml:"Name"`
	SendOnly bool   `toml:"SendOnly"`
}

func WithEvmNode(chainId, wsUrl, httpUrl string) NodeConfigOpt {
	return func(n *NodeConfig) {
		evm := EVM{
			ChainID:            chainId,
			AutoCreateKey:      true,
			FinalityDepth:      1,
			MinContractPayment: "0",
			Nodes: []Node{
				{
					WSURL:    wsUrl,
					HTTPURL:  httpUrl,
					Name:     "1337_primary_local_0",
					SendOnly: false,
				},
			},
		}
		n.EVM = append(n.EVM, evm)
	}
}
