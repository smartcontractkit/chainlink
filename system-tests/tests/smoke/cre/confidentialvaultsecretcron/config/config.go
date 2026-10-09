package config

// Config is the JSON config the confidentialvaultsecretcron workflow parses.
// The struct lives in an untagged package so the module has a lintable file;
// the workflow binary itself only builds for wasip1.
type Config struct {
	Schedule            string `json:"schedule,omitempty"`
	SecretNamespace     string `json:"secret_namespace,omitempty"`
	SecretKey           string `json:"secret_key,omitempty"`
	ExpectedSecretValue string `json:"expected_secret_value,omitempty"`
}
