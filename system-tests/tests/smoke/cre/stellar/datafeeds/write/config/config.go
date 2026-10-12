package config

type Config struct {
	ChainSelector      uint64 `yaml:"chainSelector"`
	WorkflowName       string `yaml:"workflowName"`
	CacheContractID    string `yaml:"cacheContractID"`
	DataIDHex          string `yaml:"dataIDHex"`
	Answer             int64  `yaml:"answer"`
	RequiredSignatures int    `yaml:"requiredSignatures"`
}
