package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type jobDistributorConfig struct {
	c toml.JobDistributor
}

func (s jobDistributorConfig) NopFriendlyName() string {
	if s.c.NopFriendlyName == nil {
		return ""
	}
	return *s.c.NopFriendlyName
}
