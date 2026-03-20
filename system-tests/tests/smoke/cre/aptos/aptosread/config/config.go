package config

// Config for the Aptos read consensus workflow (reads 0x1::coin::name() on local devnet).
type Config struct {
	ChainSelector    uint64 `yaml:"chainSelector"`
	WorkflowName     string `yaml:"workflowName"`
	ExpectedCoinName string `yaml:"expectedCoinName"` // expected exact value in the View reply data (e.g. "Aptos Coin" for 0x1::coin::name())
}
