package config

import (
	"net/url"
	"time"
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
	DriverName() string
	LogSQL() bool
	MaxIdleConns() int
	MaxOpenConns() int
	MigrateDatabase() bool
	URL() url.URL
}
