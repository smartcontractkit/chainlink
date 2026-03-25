package config

type Config struct {
	SecretKey       string `yaml:"secretKey"`
	SecretNamespace string `yaml:"secretNamespace"`
	ExpectNotFound  bool   `yaml:"expectNotFound"`
}
