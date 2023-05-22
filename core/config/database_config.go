package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

type Database interface {
	BackupDir() string
	BackupFrequency() time.Duration
	BackupMode() DatabaseBackupMode
	BackupOnVersionUpgrade() bool
	BackupURL() *url.URL
	DefaultIdleInTxSessionTimeout() time.Duration
	DefaultLockTimeout() time.Duration
	DefaultQueryTimeout() time.Duration
	ListenerMaxReconnectDuration() time.Duration
	ListenerMinReconnectInterval() time.Duration
	LockingMode() string
	URL() url.URL
	GetDialectConfiguredOrDefault() dialects.DialectName
	LeaseLockDuration() time.Duration
	LeaseLockRefreshInterval() time.Duration
	MigrateDatabase() bool
	MaxIdleConns() int
	ORMMaxOpenConns() int
	TriggerFallbackDBPollInterval() time.Duration
}
