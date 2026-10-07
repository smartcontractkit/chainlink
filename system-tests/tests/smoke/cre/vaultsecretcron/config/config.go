package config

// Config configures the cron-triggered vault secret reader: the schedule to run on
// and the single secret to fetch from the vault. See the module README for how this
// workflow differs from the vaultsecret verifier.
type Config struct {
	Schedule        string `yaml:"schedule,omitempty"`
	SecretNamespace string `yaml:"secretNamespace,omitempty"`
	SecretKey       string `yaml:"secretKey,omitempty"`
}
