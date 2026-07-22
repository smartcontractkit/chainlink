package config

// Config for the Stellar write smoke workflow.
//
// The workflow generates an OCR report via the Consensus capability (ed25519 /
// keccak256, "stellar" encoder), then submits it through the Stellar chain
// capability's WriteReport, which invokes the Soroban CRE forwarder's report(...)
// entrypoint. The forwarder verifies f+1 ed25519 signatures and dispatches to the
// receiver contract's on_report.
type Config struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	WorkflowName  string `yaml:"workflowName"`
	// ReceiverContractID is the C… StrKey of the Soroban receiver the forwarder
	// dispatches on_report to (passed to WriteReport as ContractId).
	ReceiverContractID string `yaml:"receiverContractID"`
	// ReportPayloadHex is an optional hex-encoded payload run through OCR report
	// generation. If empty, a small default payload is used.
	ReportPayloadHex string `yaml:"reportPayloadHex"`
	// RequiredSignatures is the number of OCR signatures to include in the
	// submitted report; the Soroban forwarder expects exactly f+1.
	RequiredSignatures int `yaml:"requiredSignatures"`
	// ExpectFailure, when true, treats a non-success tx status as the success
	// condition (negative test).
	ExpectFailure bool `yaml:"expectFailure"`
}
