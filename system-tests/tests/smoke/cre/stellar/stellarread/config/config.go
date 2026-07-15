package config

// Config for the Stellar read smoke workflow.
//
// The workflow calls the Stellar chain capability's GetLatestLedger and asserts
// the returned ledger sequence is past MinLedgerSequence.
type Config struct {
	ChainSelector     uint64 `yaml:"chainSelector"`
	WorkflowName      string `yaml:"workflowName"`
	MinLedgerSequence uint64 `yaml:"minLedgerSequence"`
}
