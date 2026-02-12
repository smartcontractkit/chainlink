package config

// Config for the log trigger + HTTP action load test workflow.
// Must produce YAML compatible with the logtrigger_config.Config used by the test harness.
type Config struct {
	ChainSelector uint64
	Addresses     []string `yaml:"addresses"`
	Topics        []struct {
		Values []string `yaml:"values"`
	} `yaml:"topics"`
	Abi   string `yaml:"abi"`
	Event string `yaml:"event"`
	URL   string `yaml:"url"` // Target URL for the HTTP action call
}
