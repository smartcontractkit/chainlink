package vrf

import (
	pkg_errors "github.com/pkg/errors"
)

type Config struct {
	Common struct {
		VRFishField *string `toml:"vrfish_field"`
	} `toml:"Common"`
}

func (o *Config) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case *Config:
		if asCfg == nil {
			return nil
		}

		if asCfg.Common.VRFishField != nil {
			o.Common.VRFishField = asCfg.Common.VRFishField
		}

		return nil
	default:
		return pkg_errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (o *Config) Validate() error {
	return nil
}
