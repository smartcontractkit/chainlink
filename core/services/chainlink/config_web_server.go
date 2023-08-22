package chainlink

import (
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ config.WebServer = (*webServerConfig)(nil)

type tlsConfig struct {
	c       toml.WebServerTLS
	rootDir func() string
}

func (t *tlsConfig) Dir() string {
	return filepath.Join(t.rootDir(), "tls")
}

func (t *tlsConfig) Host() string {
	return *t.c.Host
}

func (t *tlsConfig) HTTPSPort() uint16 {
	return *t.c.HTTPSPort
}

func (t *tlsConfig) ForceRedirect() bool {
	return *t.c.ForceRedirect
}

func (t *tlsConfig) certPath() string {
	return *t.c.CertPath
}

func (t *tlsConfig) CertFile() string {
	s := t.certPath()
	if s == "" {
		s = filepath.Join(t.Dir(), "server.crt")
	}
	return s
}

func (t *tlsConfig) keyPath() string {
	return *t.c.KeyPath
}

func (t *tlsConfig) KeyFile() string {
	if t.keyPath() == "" {
		return filepath.Join(t.Dir(), "server.key")
	}
	return t.keyPath()
}

func (t *tlsConfig) ListenIP() net.IP {
	return *t.c.ListenIP
}

type rateLimitConfig struct {
	c toml.WebServerRateLimit
}

func (r *rateLimitConfig) Authenticated() int64 {
	return *r.c.Authenticated
}

func (r *rateLimitConfig) AuthenticatedPeriod() time.Duration {
	return r.c.AuthenticatedPeriod.Duration()
}

func (r *rateLimitConfig) Unauthenticated() int64 {
	return *r.c.Unauthenticated
}

func (r *rateLimitConfig) UnauthenticatedPeriod() time.Duration {
	return r.c.UnauthenticatedPeriod.Duration()
}

type mfaConfig struct {
	c toml.WebServerMFA
}

func (m *mfaConfig) RPID() string {
	return *m.c.RPID
}

func (m *mfaConfig) RPOrigin() string {
	return *m.c.RPOrigin
}

type webServerConfig struct {
	c       toml.WebServer
	rootDir func() string
}

func (w *webServerConfig) TLS() config.TLS {
	return &tlsConfig{c: w.c.TLS, rootDir: w.rootDir}
}

func (w *webServerConfig) RateLimit() config.RateLimit {
	return &rateLimitConfig{c: w.c.RateLimit}
}

func (w *webServerConfig) MFA() config.MFA {
	return &mfaConfig{c: w.c.MFA}
}

func (w *webServerConfig) AllowOrigins() string {
	return *w.c.AllowOrigins
}

func (w *webServerConfig) BridgeResponseURL() *url.URL {
	if w.c.BridgeResponseURL.IsZero() {
		return nil
	}
	return w.c.BridgeResponseURL.URL()
}

func (w *webServerConfig) BridgeCacheTTL() time.Duration {
	return w.c.BridgeCacheTTL.Duration()
}

func (w *webServerConfig) HTTPMaxSize() int64 {
	return int64(*w.c.HTTPMaxSize)
}

func (w *webServerConfig) StartTimeout() time.Duration {
	return w.c.StartTimeout.Duration()
}

func (w *webServerConfig) HTTPWriteTimeout() time.Duration {
	return w.c.HTTPWriteTimeout.Duration()
}

func (w *webServerConfig) HTTPPort() uint16 {
	return *w.c.HTTPPort
}

func (w *webServerConfig) SessionReaperExpiration() models.Duration {
	return *w.c.SessionReaperExpiration
}

func (w *webServerConfig) SecureCookies() bool {
	return *w.c.SecureCookies
}

func (w *webServerConfig) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   w.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteStrictMode,
	}
}

func (w *webServerConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(w.c.SessionTimeout.Duration())
}

func (w *webServerConfig) ListenIP() net.IP {
	return *w.c.ListenIP
}
