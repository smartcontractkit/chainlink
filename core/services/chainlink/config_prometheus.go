package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type prometheusConfig struct {
	s toml.PrometheusSecrets
}

func (p *prometheusConfig) AuthToken() string {
	if p.s.AuthToken == nil {
		return ""
	}
	return string(*p.s.AuthToken)
}
