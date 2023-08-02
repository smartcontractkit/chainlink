package node_config

type Feature struct {
	LogPoller    bool `toml:"LogPoller"`
	FeedsManager bool `toml:"FeedsManager"`
	UICSAKeys    bool `toml:"UICSAKeys"`
}
