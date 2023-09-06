package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
)

type mercuryConfig struct {
	s toml.MercurySecrets
}

func (m *mercuryConfig) Credentials(credName string) *models.MercuryCredentials {
	if mc, ok := m.s.Credentials[credName]; ok {
		c := &models.MercuryCredentials{
			URL:      mc.URL.URL().String(),
			Password: string(*mc.Password),
			Username: string(*mc.Username),
		}
		if mc.LegacyURL != nil && mc.LegacyURL.URL() != nil {
			c.LegacyURL = mc.LegacyURL.URL().String()
		}
		return c
	}
	return nil
}
