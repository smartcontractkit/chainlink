package chainlink

import "net/url"

func (g *generalConfig) DatabaseURL() url.URL {
	if g.secrets.DatabaseURL == nil {
		return url.URL{}
	}
	return *(*url.URL)(g.secrets.DatabaseURL)
}

func (g *generalConfig) DatabaseBackupURL() *url.URL {
	return (*url.URL)(g.secrets.DatabaseBackupURL)
}

func (g *generalConfig) ExplorerAccessKey() string {
	if g.secrets.ExplorerAccessKey == nil {
		return ""
	}
	return *g.secrets.ExplorerAccessKey
}

func (g *generalConfig) ExplorerSecret() string {
	if g.secrets.ExplorerSecret == nil {
		return ""
	}
	return *g.secrets.ExplorerSecret
}

func (g *generalConfig) KeystorePassword() string {
	if g.secrets.KeystorePassword == nil {
		return ""
	}
	return *g.secrets.KeystorePassword
}

func (g *generalConfig) VRFPassword() string {
	if g.secrets.VRFPassword == nil {
		return ""
	}
	return *g.secrets.VRFPassword
}
