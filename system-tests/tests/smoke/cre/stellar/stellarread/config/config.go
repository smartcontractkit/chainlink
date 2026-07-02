package config

// Config for the Stellar read smoke workflow.
//
// The workflow calls the Stellar chain capability's GetLatestLedger and asserts
// the returned ledger sequence is past MinLedgerSequence. NOTE: chain reads are
// aggregated by identical consensus across the DON, and the latest ledger
// sequence advances over time, so different nodes may observe different values.
// This workflow is therefore best used as a single-node/connectivity smoke check;
// a stable multi-node read test should ReadContract a static on-chain value.
type Config struct {
	ChainSelector     uint64 `yaml:"chainSelector"`
	WorkflowName      string `yaml:"workflowName"`
	MinLedgerSequence uint64 `yaml:"minLedgerSequence"`
}
