package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type Backup interface {
	Dir() string
	Frequency() time.Duration
	Mode() DatabaseBackupMode
	OnVersionUpgrade() bool
	URL() *url.URL
}

type Lock interface {
	LockingMode() string
	LeaseDuration() time.Duration
	LeaseRefreshInterval() time.Duration
}

type Listener interface {
	MaxReconnectDuration() time.Duration
	MinReconnectInterval() time.Duration
	FallbackPollInterval() time.Duration
}

type Database interface {
	Backup() Backup
	Listener() Listener
	Lock() Lock

	DefaultIdleInTxSessionTimeout() time.Duration
	DefaultLockTimeout() time.Duration
	DefaultQueryTimeout() time.Duration
	Dialect() dialects.DialectName
	LogSQL() bool
	MaxIdleConns() int
	MaxOpenConns() int
	MigrateDatabase() bool
	URL() url.URL
}
