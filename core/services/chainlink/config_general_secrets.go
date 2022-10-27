package chainlink

import "net/url"

func (g *generalConfig) DatabaseURL() url.URL {
	if g.secrets.Database.URL == nil {
		return url.URL{}
	}
	return *(*url.URL)(g.secrets.Database.URL)
}

func (g *generalConfig) DatabaseBackupURL() *url.URL {
	return (*url.URL)(g.secrets.Database.BackupURL)
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
