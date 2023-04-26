package chainlink

import (
	"net/url"

	"github.com/pkg/errors"
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

func (g *generalConfig) MercurySecrets(credName string) (url, username, password string, err error) {
	for _, creds := range g.secrets.Mercury.Credentials {
		if *creds.Name == credName {
			return creds.URL.String(), string(*creds.Username), string(*creds.Password), nil
		}
	}
	return "", "", "", errors.Errorf("failed to find credentials for name: %s", credName)
}
