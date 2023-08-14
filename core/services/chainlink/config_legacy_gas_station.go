package chainlink

import (
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	lgsconfig "github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types/config"
)

type legacyGasStationConfig struct {
	s v2.LegacyGasStationSecrets
}

func (l *legacyGasStationConfig) AuthConfig() *lgsconfig.AuthConfig {
	return &lgsconfig.AuthConfig{
		ClientCertificate: string(l.s.AuthConfig.ClientCertificate),
		ClientKey:         string(l.s.AuthConfig.ClientKey),
	}
}
