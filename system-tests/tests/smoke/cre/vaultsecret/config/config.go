package config

// Check is a single secret assertion evaluated by the workflow.
type Check struct {
	Name            string `yaml:"name,omitempty" json:"name,omitempty"`
	SecretKey       string `yaml:"secretKey" json:"secretKey"`
	SecretNamespace string `yaml:"secretNamespace" json:"secretNamespace"`
	ExpectedValue   string `yaml:"expectedValue,omitempty" json:"expectedValue,omitempty"`
	ExpectNotFound  bool   `yaml:"expectNotFound" json:"expectNotFound"`
}

type Phase struct {
	Name   string  `yaml:"name" json:"name"`
	Checks []Check `yaml:"checks" json:"checks"`
	// Batch fetches every check's secret in one GetSecrets call. Checks must not use ExpectNotFound.
	Batch bool `yaml:"batch,omitempty" json:"batch,omitempty"`
}

type Config struct {
	// AuthorizedKey is the EVM address allowed to sign HTTP trigger requests.
	AuthorizedKey string `yaml:"authorizedKey" json:"authorizedKey"`
}

type TriggerInput struct {
	// Phases are evaluated in order; the first whose checks all pass logs completion.
	Phases []Phase `json:"phases,omitempty"`

	ExpectInvalidIdentifier bool     `json:"expectInvalidIdentifier,omitempty"`
	SecretKey               string   `json:"secretKey,omitempty"`
	SecretNamespace         string   `json:"secretNamespace,omitempty"`
	SecretKey2              string   `json:"secretKey2,omitempty"`
	SecretNamespace2        string   `json:"secretNamespace2,omitempty"`
	ExpectBatchTooBig       bool     `json:"expectBatchTooBig,omitempty"`
	BatchSecretKeys         []string `json:"batchSecretKeys,omitempty"`
}
