package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

// Note: this is a legacy interface. Any new fields should be added to the database
// interface defined below and accessed via cfg.Database().<FieldName>().
type DatabaseV1 interface {
	DatabaseDefaultIdleInTxSessionTimeout() time.Duration
	DatabaseDefaultLockTimeout() time.Duration
	DatabaseDefaultQueryTimeout() time.Duration
	DatabaseURL() url.URL
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	LogSQL() bool
}

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
	DatabaseDefaultIdleInTxSessionTimeout() time.Duration
	DatabaseDefaultLockTimeout() time.Duration
	DatabaseDefaultQueryTimeout() time.Duration
	DatabaseURL() url.URL
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	MigrateDatabase() bool
	MaxIdleConns() int
	MaxOpenConns() int
	LogSQL() bool
}
