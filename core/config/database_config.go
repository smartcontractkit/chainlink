package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type DatabaseV1 interface {
	DatabaseDefaultLockTimeout() time.Duration
	DatabaseDefaultQueryTimeout() time.Duration
	DatabaseListenerMaxReconnectDuration() time.Duration
	DatabaseListenerMinReconnectInterval() time.Duration
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	MigrateDatabase() bool
	ORMMaxIdleConns() int
	ORMMaxOpenConns() int
	TriggerFallbackDBPollInterval() time.Duration
	LogSQL() bool
	SetLogSQL(bool)
}

type Backup interface {
	Dir() string
	Frequency() time.Duration
	Mode() DatabaseBackupMode
	OnVersionUpgrade() bool
	URL() *url.URL
}

type Database interface {
	Backup() Backup
	DefaultIdleInTxSessionTimeout() time.Duration
	DatabaseDefaultLockTimeout() time.Duration
	DatabaseDefaultQueryTimeout() time.Duration
	DatabaseListenerMaxReconnectDuration() time.Duration
	DatabaseListenerMinReconnectInterval() time.Duration
	DatabaseLockingMode() string
	DatabaseURL() url.URL
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	MigrateDatabase() bool
	ORMMaxIdleConns() int
	ORMMaxOpenConns() int
	TriggerFallbackDBPollInterval() time.Duration
	LogSQL() bool
	SetLogSQL(bool)
}
