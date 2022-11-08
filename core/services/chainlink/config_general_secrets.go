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

func (g *generalConfig) MercuryCredentials(url string) (username, password string, err error) {
	if g.secrets.Mercury.Credentials == nil {
		return "", "", errors.New("no Mercury credentials were specified in the config")
	}
	credentials, exists := g.secrets.Mercury.Credentials[url]
	if !exists {
		return "", "", errors.Errorf("no Mercury credentials specified for server URL: %q", url)
	}
	if credentials.Username == nil {
		return "", "", errors.Errorf("no Mercury username specified for server URL: %q", url)
	}
	if credentials.Password == nil {
		return "", "", errors.Errorf("no Mercury password specified for server URL: %q", url)
	}
	return string(*credentials.Username), string(*credentials.Password), nil

}
