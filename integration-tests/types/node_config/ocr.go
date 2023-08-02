package node_config

type OCR struct {
	Enabled bool `toml:"Enabled"`
}

type OCR2 struct {
	Enabled bool `toml:"Enabled"`
}

func WithOCR1() NodeConfigOpt {
	return func(n *NodeConfig) {
		n.OCR = &OCR{Enabled: true}
	}
}

func WithOCR2() NodeConfigOpt {
	return func(n *NodeConfig) {
		n.OCR2 = &OCR2{Enabled: true}
	}
}
