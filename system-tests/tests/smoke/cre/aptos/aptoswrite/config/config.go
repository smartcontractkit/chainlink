package config

// Config for Aptos write workflow (submits a report via the Aptos write capability).
type Config struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	WorkflowName  string `yaml:"workflowName"`
	ReceiverHex   string `yaml:"receiverHex"`
	ReportMessage string `yaml:"reportMessage"`
	// When true, the workflow expects WriteReport to return a non-success tx status and treats that as success.
	ExpectFailure bool `yaml:"expectFailure"`
	// Number of OCR signatures to include in the submitted report (forwarder expects f+1).
	RequiredSignatures int `yaml:"requiredSignatures"`
	// Optional hex-encoded payload to pass through OCR report generation.
	// If empty, ReportMessage bytes are used.
	ReportPayloadHex string `yaml:"reportPayloadHex"`
	MaxGasAmount     uint64 `yaml:"maxGasAmount"`
	GasUnitPrice     uint64 `yaml:"gasUnitPrice"`
}
