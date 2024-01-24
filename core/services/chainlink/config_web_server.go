package chainlink

import (
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
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
	s       toml.WebServerSecrets
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

func (w *webServerConfig) LDAP() config.LDAP {
	return &ldapConfig{c: w.c.LDAP, s: w.s.LDAP}
}

func (w *webServerConfig) AuthenticationMethod() string {
	return *w.c.AuthenticationMethod
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

func (w *webServerConfig) SessionReaperExpiration() commonconfig.Duration {
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

func (w *webServerConfig) SessionTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(w.c.SessionTimeout.Duration())
}

func (w *webServerConfig) ListenIP() net.IP {
	return *w.c.ListenIP
}

type ldapConfig struct {
	c toml.WebServerLDAP
	s toml.WebServerLDAPSecrets
}

func (l *ldapConfig) ServerAddress() string {
	if l.s.ServerAddress == nil {
		return ""
	}
	return l.s.ServerAddress.URL().String()
}

func (l *ldapConfig) ReadOnlyUserLogin() string {
	if l.s.ReadOnlyUserLogin == nil {
		return ""
	}
	return string(*l.s.ReadOnlyUserLogin)
}

func (l *ldapConfig) ReadOnlyUserPass() string {
	if l.s.ReadOnlyUserPass == nil {
		return ""
	}
	return string(*l.s.ReadOnlyUserPass)
}

func (l *ldapConfig) ServerTLS() bool {
	if l.c.ServerTLS == nil {
		return false
	}
	return *l.c.ServerTLS
}

func (l *ldapConfig) SessionTimeout() commonconfig.Duration {
	return *l.c.SessionTimeout
}

func (l *ldapConfig) QueryTimeout() time.Duration {
	return l.c.QueryTimeout.Duration()
}

func (l *ldapConfig) UserAPITokenDuration() commonconfig.Duration {
	return *l.c.UserAPITokenDuration
}

func (l *ldapConfig) BaseUserAttr() string {
	if l.c.BaseUserAttr == nil {
		return ""
	}
	return *l.c.BaseUserAttr
}

func (l *ldapConfig) BaseDN() string {
	if l.c.BaseDN == nil {
		return ""
	}
	return *l.c.BaseDN
}

func (l *ldapConfig) UsersDN() string {
	if l.c.UsersDN == nil {
		return ""
	}
	return *l.c.UsersDN
}

func (l *ldapConfig) GroupsDN() string {
	if l.c.GroupsDN == nil {
		return ""
	}
	return *l.c.GroupsDN
}

func (l *ldapConfig) ActiveAttribute() string {
	if l.c.ActiveAttribute == nil {
		return ""
	}
	return *l.c.ActiveAttribute
}

func (l *ldapConfig) ActiveAttributeAllowedValue() string {
	if l.c.ActiveAttributeAllowedValue == nil {
		return ""
	}
	return *l.c.ActiveAttributeAllowedValue
}

func (l *ldapConfig) AdminUserGroupCN() string {
	if l.c.AdminUserGroupCN == nil {
		return ""
	}
	return *l.c.AdminUserGroupCN
}

func (l *ldapConfig) EditUserGroupCN() string {
	if l.c.EditUserGroupCN == nil {
		return ""
	}
	return *l.c.EditUserGroupCN
}

func (l *ldapConfig) RunUserGroupCN() string {
	if l.c.RunUserGroupCN == nil {
		return ""
	}
	return *l.c.RunUserGroupCN
}

func (l *ldapConfig) ReadUserGroupCN() string {
	if l.c.ReadUserGroupCN == nil {
		return ""
	}
	return *l.c.ReadUserGroupCN
}

func (l *ldapConfig) UserApiTokenEnabled() bool {
	if l.c.UserApiTokenEnabled == nil {
		return false
	}
	return *l.c.UserApiTokenEnabled
}

func (l *ldapConfig) UpstreamSyncInterval() commonconfig.Duration {
	if l.c.UpstreamSyncInterval == nil {
		return commonconfig.Duration{}
	}
	return *l.c.UpstreamSyncInterval
}

func (l *ldapConfig) UpstreamSyncRateLimit() commonconfig.Duration {
	if l.c.UpstreamSyncRateLimit == nil {
		return commonconfig.Duration{}
	}
	return *l.c.UpstreamSyncRateLimit
}
