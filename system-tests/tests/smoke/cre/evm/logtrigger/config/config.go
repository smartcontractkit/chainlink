package logtrigger

type Config struct {
	ChainSelector uint64
	Addresses     []string `yaml:"addresses"`
	Topics        []struct {
		Values []string `yaml:"values"`
	} `yaml:"topics"`
	Abi   string `yaml:"abi"`
	Event string `yaml:"event"`
}
