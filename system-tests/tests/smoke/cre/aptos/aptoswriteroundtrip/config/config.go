package config

// Config for Aptos write->read roundtrip workflow.
// The workflow writes a benchmark report, then reads back get_feeds and validates the value.
type Config struct {
	ChainSelector      uint64
	WorkflowName       string
	ReceiverHex        string
	RequiredSignatures int
	ReportPayloadHex   string
	MaxGasAmount       uint64
	GasUnitPrice       uint64
	FeedIDHex          string
	ExpectedBenchmark  uint64
}
