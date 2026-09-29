package config

// Check describes a single secret state assertion evaluated inside the workflow.
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
}

// Config is the deploy-time workflow configuration. Everything the workflow
// evaluates is supplied per invocation through the HTTP trigger input (see
// TriggerInput), so a single deployment can serve every verification scenario.
type Config struct {
	// AuthorizedKey is the hex EOA address allowed to sign HTTP trigger
	// requests for this workflow.
	AuthorizedKey string `yaml:"authorizedKey" json:"authorizedKey"`
}

// TriggerInput is the per-invocation payload carried in the HTTP trigger
// request Input field.
type TriggerInput struct {
	// Phase verification mode: phases are evaluated in declaration order and
	// the first phase whose checks all succeed emits the completion log.
	Phases []Phase `json:"phases,omitempty"`

	// Identifier-validation mode: GetSecret is expected to reject the given
	// identifiers.
	ExpectInvalidIdentifier bool   `json:"expectInvalidIdentifier,omitempty"`
	SecretKey               string `json:"secretKey,omitempty"`
	SecretNamespace         string `json:"secretNamespace,omitempty"`
	SecretKey2              string `json:"secretKey2,omitempty"`
	SecretNamespace2        string `json:"secretNamespace2,omitempty"`

	// Batch-size mode: GetSecrets is expected to reject the oversized batch.
	ExpectBatchTooBig bool     `json:"expectBatchTooBig,omitempty"`
	BatchSecretKeys   []string `json:"batchSecretKeys,omitempty"`
}
