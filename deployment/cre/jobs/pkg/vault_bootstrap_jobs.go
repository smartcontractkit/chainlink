package pkg

type VaultBootstrapJobsInput struct {
	ContractQualifierPrefix string        `json:"contract_qualifier_prefix" yaml:"contract_qualifier_prefix"`
	ChainSelector           ChainSelector `json:"chain_selector" yaml:"chain_selector"`
}
