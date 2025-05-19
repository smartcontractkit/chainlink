package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

func TestDatabaseConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
		SecretsStrings: []string{`[Database]
URL = "postgresql://doesnotexist:justtopassvalidationtests@localhost:5432/chainlink_na_test"`},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	backup := cfg.Database().Backup()
	assert.Equal(t, "test/backup/dir", backup.Dir())
	assert.Equal(t, 1*time.Hour, backup.Frequency())
	assert.Equal(t, config.DatabaseBackupModeFull, backup.Mode())
	assert.True(t, backup.OnVersionUpgrade())
	assert.Nil(t, backup.URL())

	db := cfg.Database()
	assert.Equal(t, 1*time.Minute, db.DefaultIdleInTxSessionTimeout())
	assert.Equal(t, 1*time.Hour, db.DefaultLockTimeout())
	assert.Equal(t, 1*time.Second, db.DefaultQueryTimeout())
	assert.True(t, db.LogSQL())
	assert.Equal(t, 7, db.MaxIdleConns())
	assert.Equal(t, 13, db.MaxOpenConns())
	assert.True(t, db.MigrateDatabase())
	assert.Equal(t, pgcommon.DriverPostgres, db.DriverName())
	url := db.URL()
	assert.NotEqual(t, "", url.String())

	lock := db.Lock()
	assert.Equal(t, "none", lock.LockingMode())
	assert.Equal(t, 1*time.Minute, lock.LeaseDuration())
	assert.Equal(t, 1*time.Second, lock.LeaseRefreshInterval())

	l := db.Listener()
	assert.Equal(t, 1*time.Minute, l.MaxReconnectDuration())
	assert.Equal(t, 5*time.Minute, l.MinReconnectInterval())
	assert.Equal(t, 2*time.Minute, l.FallbackPollInterval())
}
