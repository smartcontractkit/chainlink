package config

import (
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type TLS interface {
	TLSCertPath() string
	TLSDir() string
	TLSKeyPath() string
	TLSPort() uint16
	KeyFile() string
	CertFile() string

	Host() string
	ForceRedirect() bool
}

type WebServer interface {
	AllowOrigins() string
	AuthenticatedRateLimit() int64
	AuthenticatedRateLimitPeriod() models.Duration
	BridgeCacheTTL() time.Duration
	BridgeResponseURL() *url.URL
	HTTPServerWriteTimeout() time.Duration
	Port() uint16
	RPID() string
	RPOrigin() string
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
	ReaperExpiration() models.Duration
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionTimeout() models.Duration
	WebServerHTTPMaxSize() int64
	WebServerStartTimeout() time.Duration

	TLS() TLS
}

type WebV1 interface {
	AuthenticatedRateLimit() int64
	AuthenticatedRateLimitPeriod() models.Duration
	BridgeCacheTTL() time.Duration
	BridgeResponseURL() *url.URL
	CertFile() string
	HTTPServerWriteTimeout() time.Duration
	KeyFile() string
	Port() uint16
	RPID() string
	RPOrigin() string
	TLSCertPath() string
	TLSDir() string
	TLSKeyPath() string
	TLSPort() uint16
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
	ReaperExpiration() models.Duration
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionTimeout() models.Duration
	WebServerHTTPMaxSize() int64
	WebServerStartTimeout() time.Duration
}
