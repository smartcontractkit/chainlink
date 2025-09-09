package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type jobDistributorConfig struct {
	c toml.JobDistributor
}

func (s jobDistributorConfig) DisplayName() string {
	if s.c.DisplayName == nil {
		return ""
	}
	return *s.c.DisplayName
}
