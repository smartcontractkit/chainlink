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
		return &models.MercuryCredentials{
			URL:      mc.URL.URL().String(),
			Username: mc.Username.String(),
			Password: mc.Password.String(),
		}
	}
	return nil
}
