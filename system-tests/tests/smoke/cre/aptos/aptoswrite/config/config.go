package config

// Config for Aptos write workflow (submits a report via the Aptos write capability).
type Config struct {
	ChainSelector uint64
	WorkflowName  string
	ReceiverHex   string
	ReportMessage string
	// When true, the workflow expects WriteReport to return TX_STATUS_FAILED and treats that as success.
	ExpectFailure bool
	// Number of OCR signatures to include in the submitted report (forwarder expects f+1).
	RequiredSignatures int
	// Optional hex-encoded payload to pass through OCR report generation.
	// If empty, ReportMessage bytes are used.
	ReportPayloadHex string
	MaxGasAmount     uint64
	GasUnitPrice     uint64
}
