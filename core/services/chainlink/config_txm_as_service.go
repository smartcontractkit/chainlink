package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.TxmAsService = (*txmAsServiceConfig)(nil)

type txmAsServiceConfig struct {
	c toml.TxmAsService
}

func (t *txmAsServiceConfig) Enabled() bool {
	return t.c.Enabled != nil && *t.c.Enabled
}
