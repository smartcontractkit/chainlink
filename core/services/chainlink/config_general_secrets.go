package chainlink

import (
	"fmt"
	"net/url"
	"strings"

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
	for _, creds := range g.secrets.Mercury.Credentials {
		if creds.URL != nil && creds.URL.URL().String() == url {
			if creds.Username == nil {
				return "", "", errors.Errorf("no Mercury username specified for server URL: %q", url)
			}
			if creds.Password == nil {
				return "", "", errors.Errorf("no Mercury password specified for server URL: %q", url)
			}
			return string(*creds.Username), string(*creds.Password), nil
		}
	}
	msg := fmt.Sprintf("no Mercury credentials specified for server URL: %q", url)
	if len(g.secrets.Mercury.Credentials) > 0 {
		urls := make([]string, len(g.secrets.Mercury.Credentials))
		for i, creds := range g.secrets.Mercury.Credentials {
			urls[i] = creds.URL.String()
		}
		msg += fmt.Sprintf(" (credentials available for these urls: %s)", strings.Join(urls, ","))
	}
	return "", "", errors.New(msg)
}
