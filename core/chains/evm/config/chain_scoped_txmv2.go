package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type txmv2Config struct {
	c toml.TxmV2
}

func (t *txmv2Config) Enabled() bool {
	return *t.c.Enabled
}

func (t *txmv2Config) BlockTime() *time.Duration {
	d := t.c.BlockTime.Duration()
	return &d
}

func (t *txmv2Config) CustomURL() *url.URL {
	return t.c.CustomURL.URL()
}
