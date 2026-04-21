package config

type Check struct {
	Name            string `yaml:"name,omitempty"`
	SecretKey       string `yaml:"secretKey"`
	SecretNamespace string `yaml:"secretNamespace"`
	ExpectedValue   string `yaml:"expectedValue,omitempty"`
	ExpectNotFound  bool   `yaml:"expectNotFound"`
}

type Phase struct {
	Name   string  `yaml:"name"`
	Checks []Check `yaml:"checks"`
}

type Config struct {
	Phases []Phase `yaml:"phases"`
	Checks []Check `yaml:"checks"`

	// Legacy single-check fields kept for compatibility with any older callers.
	SecretKey       string `yaml:"secretKey,omitempty"`
	SecretNamespace string `yaml:"secretNamespace,omitempty"`
	ExpectedValue   string `yaml:"expectedValue,omitempty"`
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
		ExpectedValue:   c.ExpectedValue,
		ExpectNotFound:  c.ExpectNotFound,
	}}
}

func (c Config) EffectivePhases() []Phase {
	if len(c.Phases) > 0 {
		return c.Phases
	}

	checks := c.EffectiveChecks()
	if len(checks) == 0 {
		return nil
	}

	return []Phase{{
		Name:   "default",
		Checks: checks,
	}}
}
