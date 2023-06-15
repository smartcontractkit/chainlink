package config

import (
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Web interface {
	AllowOrigins() string
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
	TLSHost() string
	TLSKeyPath() string
	TLSPort() uint16
	TLSRedirect() bool
	UnAuthenticatedRateLimit() int64
	UnAuthenticatedRateLimitPeriod() models.Duration
	ReaperExpiration() models.Duration
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionTimeout() models.Duration
}
