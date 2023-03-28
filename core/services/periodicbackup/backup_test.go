package periodicbackup

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func mustNewDatabaseBackup(t *testing.T, config Config) *databaseBackup {
	testutils.SkipShortDB(t)
	b, err := NewDatabaseBackup(config, logger.TestLogger(t))
	require.NoError(t, err)
	return b.(*databaseBackup)
}

func must(t testing.TB, s string) *url.URL {
	v, err := url.Parse(s)
	require.NoError(t, err)
	return v
}

func TestPeriodicBackup_RunBackup(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(), "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_0.9.9")
	assert.NotContains(t, result.pgDumpArguments, "--exclude-table-data=pipeline_task_runs")
}

func TestPeriodicBackup_RunBackupInLiteMode(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(), "", config.DatabaseBackupModeLite)
	periodicBackup := mustNewDatabaseBackup(t, backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_0.9.9")
	assert.Contains(t, result.pgDumpArguments, "--exclude-table-data=pipeline_task_runs")
}

func TestPeriodicBackup_RunBackupWithoutVersion(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(), "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup(static.Unset)
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_unset")
}

func TestPeriodicBackup_RunBackupViaAltUrlAndMaskPassword(t *testing.T) {
	altUrl, _ := url.Parse("postgresql://invalid:some-pass@invalid")
	backupConfig := newTestConfig(time.Minute, altUrl, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(), "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	partialResult, err := periodicBackup.runBackup("")
	require.Error(t, err, "connection to database \"postgresql//invalid\" failed")
	assert.Contains(t, partialResult.maskedArguments, "postgresql://invalid:xxxxx@invalid")
}

func TestPeriodicBackup_FrequencyTooSmall(t *testing.T) {
	backupConfig := newTestConfig(time.Second, nil, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(), "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, backupConfig)
	assert.True(t, periodicBackup.frequencyIsTooSmall())
}

func TestPeriodicBackup_AlternativeOutputDir(t *testing.T) {
	backupConfig := newTestConfig(time.Second, nil, must(t, string(v2.EnvDatabaseURL.Get())), os.TempDir(),
		filepath.Join(os.TempDir(), "alternative"), config.DatabaseBackupModeFull)

	periodicBackup := mustNewDatabaseBackup(t, backupConfig)

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Contains(t, result.path, "/alternative/cl_backup_0.9.9.dump")

}

type testConfig struct {
	databaseBackupFrequency time.Duration
	databaseBackupMode      config.DatabaseBackupMode
	databaseBackupURL       *url.URL
	databaseBackupDir       string
	databaseURL             url.URL
	rootDir                 string
}

func (config testConfig) DatabaseBackupFrequency() time.Duration {
	return config.databaseBackupFrequency
}
func (config testConfig) DatabaseBackupMode() config.DatabaseBackupMode {
	return config.databaseBackupMode
}
func (config testConfig) DatabaseBackupURL() *url.URL {
	return config.databaseBackupURL
}
func (config testConfig) DatabaseBackupDir() string {
	return config.databaseBackupDir
}
func (config testConfig) DatabaseURL() url.URL {
	return config.databaseURL
}
func (config testConfig) RootDir() string {
	return config.rootDir
}

func newTestConfig(frequency time.Duration, databaseBackupURL *url.URL, databaseURL *url.URL, rootDir string, databaseBackupDir string, mode config.DatabaseBackupMode) testConfig {
	return testConfig{
		databaseBackupFrequency: frequency,
		databaseBackupMode:      mode,
		databaseBackupURL:       databaseBackupURL,
		databaseURL:             *databaseURL,
		rootDir:                 rootDir,
		databaseBackupDir:       databaseBackupDir,
	}
}
