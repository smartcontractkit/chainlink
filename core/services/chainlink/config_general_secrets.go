package chainlink

import (
	"net/url"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
)

func (g *generalConfig) URL() url.URL {
	if g.secrets.Database.URL == nil {
		return url.URL{}
	}
	return *g.secrets.Database.URL.URL()
}

func (g *generalConfig) BackupURL() *url.URL {
	return g.secrets.Database.BackupURL.URL()
}

func (g *generalConfig) ExplorerAccessKey() string {
	if g.secrets.Explorer.AccessKey == nil {
		return ""
	}
	return string(*g.secrets.Explorer.AccessKey)
}

func (g *generalConfig) ExplorerSecret() string {
	if g.secrets.Explorer.Secret == nil {
		return ""
	}
	return string(*g.secrets.Explorer.Secret)
}
func (g *generalConfig) PyroscopeAuthToken() string {
	if g.secrets.Pyroscope.AuthToken == nil {
		return ""
	}
	return string(*g.secrets.Pyroscope.AuthToken)
}

func (g *generalConfig) PrometheusAuthToken() string {
	if g.secrets.Prometheus.AuthToken == nil {
		return ""
	}
	return string(*g.secrets.Prometheus.AuthToken)
}

func (g *generalConfig) MercuryCredentials(credName string) *models.MercuryCredentials {
	if mc, ok := g.secrets.Mercury.Credentials[credName]; ok {
		return &models.MercuryCredentials{
			URL:      mc.URL.URL().String(),
			Username: string(*mc.Username),
			Password: string(*mc.Password),
		}
	}
	return nil
}
