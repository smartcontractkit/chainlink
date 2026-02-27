package config

// Config for the Aptos read consensus workflow (reads 0x1::coin::name() on local devnet).
type Config struct {
	ChainSelector    uint64
	WorkflowName     string
	ExpectedCoinName string // expected substring in the View reply data (e.g. "Aptos" for 0x1::coin::name())
}
