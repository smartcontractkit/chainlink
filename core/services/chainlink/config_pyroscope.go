package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type pyroscopeConfig struct {
	c toml.Pyroscope
	s toml.PyroscopeSecrets
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
