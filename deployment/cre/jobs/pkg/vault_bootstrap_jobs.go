package pkg

type VaultBootstrapJobsInput struct {
	ContractQualifierPrefix string        `json:"contractQualifierPrefix" yaml:"contractQualifierPrefix"`
	ChainSelector           Uint64 `json:"chainSelector" yaml:"chainSelector"`
}
