package chainlink

import (
	"net/url"
)

func (g *generalConfig) DatabaseURL() url.URL {
	if g.secrets.Database.URL == nil {
		return url.URL{}
	}
	return *g.secrets.Database.URL.URL()
}

func (g *generalConfig) DatabaseBackupURL() *url.URL {
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

func (g *generalConfig) MercuryID() string {
	if g.secrets.Mercury.ID == nil {
		return ""
	}
	return string(*g.secrets.Mercury.ID)
}

func (g *generalConfig) MercuryKey() string {
	if g.secrets.Mercury.Key == nil {
		return ""
	}
	return string(*g.secrets.Mercury.Key)
}

func (g *generalConfig) MercuryURL() *url.URL {
	if g.secrets.Mercury.URL == nil {
		return nil
	}
	return g.secrets.Mercury.URL.URL()
}
