package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.DynamicConfig = (*dynamicConfigConfig)(nil)

type dynamicConfigConfig struct {
	c toml.DynamicConfig
}

func (d *dynamicConfigConfig) Enabled() bool {
	if d.c.Enabled == nil {
		return false
	}

	return *d.c.Enabled
}

func (d *dynamicConfigConfig) Path() string {
	if d.c.Path == nil {
		return "/config/runtime.json"
	}

	return *d.c.Path
}

func (d *dynamicConfigConfig) ManagedJobPrefix() string {
	if d.c.ManagedJobPrefix == nil {
		return "dyncfg-"
	}

	return *d.c.ManagedJobPrefix
}
