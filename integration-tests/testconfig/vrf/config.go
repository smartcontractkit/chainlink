package vrf

type Config struct {
}

func (o *Config) ApplyOverrides(_ *Config) error {
	return nil
}

func (o *Config) Validate() error {
	return nil
}
