package chainlink

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type backupConfig struct {
	c toml.DatabaseBackup
	s toml.DatabaseSecrets
}

func (b *backupConfig) Dir() string {
	return *b.c.Dir
}

func (b *backupConfig) Frequency() time.Duration {
	return b.c.Frequency.Duration()
}

func (b *backupConfig) Mode() config.DatabaseBackupMode {
	return *b.c.Mode
}

func (b *backupConfig) OnVersionUpgrade() bool {
	return *b.c.OnVersionUpgrade
}

func (b *backupConfig) URL() *url.URL {
	return b.s.BackupURL.URL()
}

type lockConfig struct {
	c toml.DatabaseLock
}

func (l *lockConfig) LockingMode() string {
	return l.c.Mode()
}

func (l *lockConfig) LeaseDuration() time.Duration {
	return l.c.LeaseDuration.Duration()
}

func (l *lockConfig) LeaseRefreshInterval() time.Duration {
	return l.c.LeaseRefreshInterval.Duration()
}

type listenerConfig struct {
	c toml.DatabaseListener
}

func (l *listenerConfig) MaxReconnectDuration() time.Duration {
	return l.c.MaxReconnectDuration.Duration()
}

func (l *listenerConfig) MinReconnectInterval() time.Duration {
	return l.c.MinReconnectInterval.Duration()
}

func (l *listenerConfig) FallbackPollInterval() time.Duration {
	return l.c.FallbackPollInterval.Duration()
}

var _ config.Database = (*databaseConfig)(nil)

type databaseConfig struct {
	c      toml.Database
	s      toml.DatabaseSecrets
	logSQL func() bool
}

func (d *databaseConfig) Backup() config.Backup {
	return &backupConfig{
		c: d.c.Backup,
		s: d.s,
	}
}

func (d *databaseConfig) Lock() config.Lock {
	return &lockConfig{
		d.c.Lock,
	}
}

func (d *databaseConfig) Listener() config.Listener {
	return &listenerConfig{
		c: d.c.Listener,
	}
}

func (d *databaseConfig) DefaultIdleInTxSessionTimeout() time.Duration {
	return d.c.DefaultIdleInTxSessionTimeout.Duration()
}

func (d *databaseConfig) DefaultLockTimeout() time.Duration {
	return d.c.DefaultLockTimeout.Duration()
}

func (d *databaseConfig) DefaultQueryTimeout() time.Duration {
	return d.c.DefaultQueryTimeout.Duration()
}

func (d *databaseConfig) URL() url.URL {
	if d.s.URL == nil {
		return url.URL{}
	}
	return *d.s.URL.URL()
}

func (d *databaseConfig) DriverName() string {
	return d.c.DriverName
}

func (d *databaseConfig) MigrateDatabase() bool {
	return *d.c.MigrateOnStartup
}

func (d *databaseConfig) MaxIdleConns() int {
	return int(*d.c.MaxIdleConns)
}

func (d *databaseConfig) MaxOpenConns() int {
	return int(*d.c.MaxOpenConns)
}

func (d *databaseConfig) LogSQL() (sql bool) {
	return d.logSQL()
}
