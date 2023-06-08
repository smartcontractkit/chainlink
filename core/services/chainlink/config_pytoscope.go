package chainlink

import v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"

type PyroscopeConfig struct {
	c v2.Pyroscope
	s v2.PyroscopeSecrets
}

func (p *PyroscopeConfig) AuthToken() string {
	if p.s.AuthToken == nil {
		return ""
	}
	return string(*p.s.AuthToken)
}

func (p *PyroscopeConfig) ServerAddress() string {
	return *p.c.ServerAddress
}

func (p *PyroscopeConfig) Environment() string {
	return *p.c.Environment
}
