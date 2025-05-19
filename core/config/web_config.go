package config

import (
	"net"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type TLS interface {
	Dir() string
	Host() string
	ForceRedirect() bool
	CertFile() string
	KeyFile() string
	HTTPSPort() uint16
	ListenIP() net.IP
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

type LDAP interface {
	ServerAddress() string
	ReadOnlyUserLogin() string
	ReadOnlyUserPass() string
	ServerTLS() bool
	SessionTimeout() commonconfig.Duration
	QueryTimeout() time.Duration
	BaseUserAttr() string
	BaseDN() string
	UsersDN() string
	GroupsDN() string
	ActiveAttribute() string
	ActiveAttributeAllowedValue() string
	AdminUserGroupCN() string
	EditUserGroupCN() string
	RunUserGroupCN() string
	ReadUserGroupCN() string
	UserApiTokenEnabled() bool
	UserAPITokenDuration() commonconfig.Duration
	UpstreamSyncInterval() commonconfig.Duration
	UpstreamSyncRateLimit() commonconfig.Duration
}

type WebServer interface {
	AuthenticationMethod() string
	AllowOrigins() string
	BridgeCacheTTL() time.Duration
	BridgeResponseURL() *url.URL
	HTTPMaxSize() int64
	StartTimeout() time.Duration
	HTTPWriteTimeout() time.Duration
	HTTPPort() uint16
	SessionReaperExpiration() commonconfig.Duration
	SecureCookies() bool
	SessionOptions() sessions.Options
	SessionTimeout() commonconfig.Duration
	ListenIP() net.IP

	TLS() TLS
	RateLimit() RateLimit
	MFA() MFA
	LDAP() LDAP
}
