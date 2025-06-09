package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Billing = (*billingConfig)(nil)

type billingConfig struct {
	t toml.Billing
}

func (c *billingConfig) URL() string {
	return *c.t.URL
}
