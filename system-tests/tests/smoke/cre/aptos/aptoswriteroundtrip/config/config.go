package config

// Config for Aptos write->read roundtrip workflow.
// The workflow writes a benchmark report, then reads back get_feeds and validates the value.
type Config struct {
	ChainSelector      uint64 `yaml:"chainSelector"`
	WorkflowName       string `yaml:"workflowName"`
	ReceiverHex        string `yaml:"receiverHex"`
	RequiredSignatures int    `yaml:"requiredSignatures"`
	ReportPayloadHex   string `yaml:"reportPayloadHex"`
	MaxGasAmount       uint64 `yaml:"maxGasAmount"`
	GasUnitPrice       uint64 `yaml:"gasUnitPrice"`
	FeedIDHex          string `yaml:"feedIDHex"`
	ExpectedBenchmark  uint64 `yaml:"expectedBenchmark"`
}
