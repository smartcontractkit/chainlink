package pkg

type VaultBootstrapJobsInput struct {
	ContractQualifierPrefix string        `json:"contractQualifierPrefix" yaml:"contractQualifierPrefix"`
	ChainSelector           ChainSelector `json:"chainSelector" yaml:"chainSelector"`
}
