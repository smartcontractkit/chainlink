package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Transmission = (*transmissionConfig)(nil)

type transmissionConfig struct {
	s toml.Transmission
}

func (t transmissionConfig) TLS() config.TransmissionTLS {
	return transmissionTLSConfig{t.s.TLS}
}

type transmissionTLSConfig struct {
	s toml.TransmissionTLS
}

func (t transmissionTLSConfig) CertPath() string {
	return *t.s.CertPath
}
