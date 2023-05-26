package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

func TestDatabaseConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	backup := cfg.Database().Backup()
	assert.Equal(t, backup.Dir(), "test/backup/dir")
	assert.Equal(t, backup.Frequency(), 1*time.Hour)
	assert.Equal(t, backup.Mode(), config.DatabaseBackupModeFull)
	assert.Equal(t, backup.OnVersionUpgrade(), true)
	assert.Nil(t, backup.URL())

	db := cfg.Database()
	assert.Equal(t, db.DatabaseDefaultIdleInTxSessionTimeout(), 1*time.Minute)
	assert.Equal(t, db.DatabaseDefaultLockTimeout(), 1*time.Hour)
	assert.Equal(t, db.DatabaseDefaultQueryTimeout(), 1*time.Second)
	assert.Equal(t, db.LogSQL(), true)
	assert.Equal(t, db.ORMMaxIdleConns(), 7)
	assert.Equal(t, db.ORMMaxOpenConns(), 13)
	assert.Equal(t, db.MigrateDatabase(), true)
	assert.Equal(t, db.DatabaseListenerMaxReconnectDuration(), 1*time.Minute)
	assert.Equal(t, db.DatabaseListenerMinReconnectInterval(), 5*time.Minute)
	assert.Equal(t, db.TriggerFallbackDBPollInterval(), 2*time.Minute)
	assert.Equal(t, db.GetDatabaseDialectConfiguredOrDefault(), dialects.Postgres)
	url := db.DatabaseURL()
	assert.NotEqual(t, url.String(), "")

	lock := db.Lock()
	assert.Equal(t, lock.LockingMode(), "none")
	assert.Equal(t, lock.LeaseDuration(), 1*time.Minute)
	assert.Equal(t, lock.LeaseRefreshInterval(), 1*time.Second)
}
