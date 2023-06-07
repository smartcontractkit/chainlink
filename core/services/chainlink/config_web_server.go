package chainlink

import (
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ config.WebServer = (*webServerConfig)(nil)

type webServerConfig struct {
	c       v2.WebServer
	rootDir func() string
}

func (w *webServerConfig) AllowOrigins() string {
	return *w.c.AllowOrigins
}

func (w *webServerConfig) AuthenticatedRateLimit() int64 {
	return *w.c.RateLimit.Authenticated
}

func (w *webServerConfig) AuthenticatedRateLimitPeriod() models.Duration {
	return *w.c.RateLimit.AuthenticatedPeriod
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

func (w *webServerConfig) CertFile() string {
	s := *w.c.TLS.CertPath
	if s == "" {
		s = filepath.Join(w.TLSDir(), "server.crt")
	}
	return s
}

func (w *webServerConfig) WebServerHTTPMaxSize() int64 {
	return int64(*w.c.HTTPMaxSize)
}

func (w *webServerConfig) WebServerStartTimeout() time.Duration {
	return w.c.StartTimeout.Duration()
}

func (w *webServerConfig) HTTPServerWriteTimeout() time.Duration {
	return w.c.HTTPWriteTimeout.Duration()
}

func (w *webServerConfig) Port() uint16 {
	return *w.c.HTTPPort
}

func (w *webServerConfig) RPID() string {
	return *w.c.MFA.RPID
}

func (w *webServerConfig) RPOrigin() string {
	return *w.c.MFA.RPOrigin
}

func (w *webServerConfig) ReaperExpiration() models.Duration {
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

func (w *webServerConfig) TLSCertPath() string {
	return *w.c.TLS.CertPath
}

func (w *webServerConfig) TLSDir() string {
	return filepath.Join(w.rootDir(), "tls")
}

func (w *webServerConfig) TLSHost() string {
	return *w.c.TLS.Host
}

func (w *webServerConfig) TLSKeyPath() string {
	return *w.c.TLS.KeyPath
}

func (w *webServerConfig) TLSPort() uint16 {
	return *w.c.TLS.HTTPSPort
}

func (w *webServerConfig) TLSRedirect() bool {
	return *w.c.TLS.ForceRedirect
}

func (w *webServerConfig) UnAuthenticatedRateLimit() int64 {
	return *w.c.RateLimit.Unauthenticated
}

func (w *webServerConfig) UnAuthenticatedRateLimitPeriod() models.Duration {
	return *w.c.RateLimit.UnauthenticatedPeriod
}

func (w *webServerConfig) KeyFile() string {
	if w.TLSKeyPath() == "" {
		return filepath.Join(w.TLSDir(), "server.key")
	}
	return w.TLSKeyPath()
}
