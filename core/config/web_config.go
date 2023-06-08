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
	CertFile() string
	KeyFile() string
	HTTPSPort() uint16
}

type RateLimit interface {
	Authenticated() int64
	AuthenticatedPeriod() time.Duration
	Unauthenticated() int64
	UnauthenticatedPeriod() time.Duration
}

type MFA interface {
	RPID() string
	RPOrigin() string
}

type WebServer interface {
	AllowOrigins() string
	BridgeCacheTTL() time.Duration
	BridgeResponseURL() *url.URL
	HTTPMaxSize() int64
	StartTimeout() time.Duration
	HTTPWriteTimeout() time.Duration
	HTTPPort() uint16
	SessionReaperExpiration() models.Duration
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionTimeout() models.Duration

	TLS() TLS
	RateLimit() RateLimit
	MFA() MFA
}
