package chainlink

import v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"

type pyroscopeConfig struct {
	c v2.Pyroscope
	s v2.PyroscopeSecrets
}

func (p *pyroscopeConfig) AuthToken() string {
	if p.s.AuthToken == nil {
		return ""
	}
	return string(*p.s.AuthToken)
}

func (p *pyroscopeConfig) ServerAddress() string {
	return *p.c.ServerAddress
}

func (p *pyroscopeConfig) Environment() string {
	return *p.c.Environment
}
