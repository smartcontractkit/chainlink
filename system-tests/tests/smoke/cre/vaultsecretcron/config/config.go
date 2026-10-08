package config

type Config struct {
	Schedule        string `yaml:"schedule,omitempty"`
	SecretNamespace string `yaml:"secretNamespace,omitempty"`
	SecretKey       string `yaml:"secretKey,omitempty"`
}
