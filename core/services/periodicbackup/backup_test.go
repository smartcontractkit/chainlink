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
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func mustNewDatabaseBackup(t *testing.T, url url.URL, rootDir string, config BackupConfig) *databaseBackup {
	testutils.SkipShortDB(t)
	b, err := NewDatabaseBackup(url, rootDir, config, logger.TestLogger(t))
	require.NoError(t, err)
	return b.(*databaseBackup)
}

func must(t testing.TB, s string) *url.URL {
	v, err := url.Parse(s)
	require.NoError(t, err)
	return v
}

func TestPeriodicBackup_RunBackup(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Positive(t, file.Size())
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_0.9.9")
	assert.NotContains(t, result.pgDumpArguments, "--exclude-table-data=pipeline_task_runs")
}

func TestPeriodicBackup_RunBackupInLiteMode(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, "", config.DatabaseBackupModeLite)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Positive(t, file.Size())
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_0.9.9")
	assert.Contains(t, result.pgDumpArguments, "--exclude-table-data=pipeline_task_runs")
}

func TestPeriodicBackup_RunBackupWithoutVersion(t *testing.T) {
	backupConfig := newTestConfig(time.Minute, nil, "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup(static.Unset)
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Positive(t, file.Size())
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "backup/cl_backup_unset")
}

func TestPeriodicBackup_RunBackupViaAltUrlAndMaskPassword(t *testing.T) {
	altUrl, _ := url.Parse("postgresql://invalid:some-pass@invalid")
	backupConfig := newTestConfig(time.Minute, altUrl, "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	partialResult, err := periodicBackup.runBackup("")
	require.Error(t, err, "connection to database \"postgresql//invalid\" failed")
	assert.Contains(t, partialResult.maskedArguments, "postgresql://invalid:xxxxx@invalid")
}

func TestPeriodicBackup_FrequencyTooSmall(t *testing.T) {
	backupConfig := newTestConfig(time.Second, nil, "", config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)
	assert.True(t, periodicBackup.frequencyIsTooSmall())
}

func TestPeriodicBackup_AlternativeOutputDir(t *testing.T) {
	backupDir := filepath.Join(os.TempDir(), "alternative")
	backupConfig := newTestConfig(time.Second, nil, backupDir, config.DatabaseBackupModeFull)
	periodicBackup := mustNewDatabaseBackup(t, *(must(t, string(env.DatabaseURL.Get()))), os.TempDir(), backupConfig)

	result, err := periodicBackup.runBackup("0.9.9")
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Positive(t, file.Size())
	assert.Contains(t, result.path, "/alternative/cl_backup_0.9.9.dump")
}

type testConfig struct {
	frequency time.Duration
	mode      config.DatabaseBackupMode
	url       *url.URL
	dir       string
}

func (t *testConfig) Frequency() time.Duration {
	return t.frequency
}

func (t *testConfig) Mode() config.DatabaseBackupMode {
	return t.mode
}

func (t *testConfig) URL() *url.URL {
	return t.url
}

func (t *testConfig) Dir() string {
	return t.dir
}

func newTestConfig(frequency time.Duration, databaseBackupURL *url.URL, databaseBackupDir string, mode config.DatabaseBackupMode) *testConfig {
	return &testConfig{
		frequency: frequency,
		mode:      mode,
		url:       databaseBackupURL,
		dir:       databaseBackupDir,
	}
}
