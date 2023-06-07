package config

import (
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type TLS interface {
	Dir() string
	Host() string
	ForceRedirect() bool
	CertPath() string
	CertFile() string
	KeyPath() string
	KeyFile() string
	HTTPSPort() uint16
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
}
