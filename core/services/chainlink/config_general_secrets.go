package chainlink

import "net/url"

func (g *generalConfig) DatabaseURL() url.URL {
	return *(*url.URL)(g.secrets.DatabaseURL)
}

func (g *generalConfig) DatabaseBackupURL() *url.URL {
	return (*url.URL)(g.secrets.DatabaseBackupURL)
}

func (g *generalConfig) ExplorerAccessKey() string {
	return *g.secrets.ExplorerAccessKey
}

func (g *generalConfig) ExplorerSecret() string {
	return *g.secrets.ExplorerSecret
}

func (g *generalConfig) KeystorePassword() string {
	return *g.secrets.KeystorePassword
}

func (g *generalConfig) VRFPassword() string {
	return *g.secrets.VRFPassword
}
