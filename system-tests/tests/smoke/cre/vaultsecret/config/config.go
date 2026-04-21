package config

type Check struct {
	Name            string `yaml:"name,omitempty"`
	SecretKey       string `yaml:"secretKey"`
	SecretNamespace string `yaml:"secretNamespace"`
	ExpectNotFound  bool   `yaml:"expectNotFound"`
}

type Config struct {
	Checks []Check `yaml:"checks"`

	// Legacy single-check fields kept for compatibility with any older callers.
	SecretKey       string `yaml:"secretKey,omitempty"`
	SecretNamespace string `yaml:"secretNamespace,omitempty"`
	ExpectNotFound  bool   `yaml:"expectNotFound,omitempty"`
}

func (c Config) EffectiveChecks() []Check {
	if len(c.Checks) > 0 {
		return c.Checks
	}

	if c.SecretKey == "" && c.SecretNamespace == "" {
		return nil
	}

	return []Check{{
		SecretKey:       c.SecretKey,
		SecretNamespace: c.SecretNamespace,
		ExpectNotFound:  c.ExpectNotFound,
	}}
}
