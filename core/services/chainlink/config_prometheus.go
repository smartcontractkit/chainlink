package chainlink

import (
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
)

type prometheusConfig struct {
	s v2.PrometheusSecrets
}

func (p *prometheusConfig) AuthToken() string {
	if p.s.AuthToken == nil {
		return ""
	}
	return string(*p.s.AuthToken)
}
