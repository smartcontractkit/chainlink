package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type Database interface {
	DatabaseBackupDir() string
	DatabaseBackupFrequency() time.Duration
	DatabaseBackupMode() DatabaseBackupMode
	DatabaseBackupOnVersionUpgrade() bool
	DatabaseBackupURL() *url.URL
	DatabaseDefaultIdleInTxSessionTimeout() time.Duration
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
}
