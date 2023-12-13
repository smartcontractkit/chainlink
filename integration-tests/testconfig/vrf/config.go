package vrf

type Config struct {
	Common struct {
		VRFishField *string `toml:"vrfish_field"`
	} `toml:"Common"`
}

func (o *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}

	if from.Common.VRFishField != nil {
		o.Common.VRFishField = from.Common.VRFishField
	}

	return nil
}

func (o *Config) Validate() error {
	return nil
}
