package periodicbackup

import (
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodicBackup_RunBackup(t *testing.T) {
	rawConfig := orm.NewConfig()
	periodicBackup := NewDatabaseBackup(time.Minute, rawConfig.DatabaseURL(), os.TempDir(), logger.Default).(*databaseBackup)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup()
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.Contains(t, result.path, "cl_backup")
}

func TestPeriodicBackup_FrequencyTooSmall(t *testing.T) {
	rawConfig := orm.NewConfig()
	periodicBackup := NewDatabaseBackup(time.Second, rawConfig.DatabaseURL(), os.TempDir(), logger.Default).(*databaseBackup)
	assert.True(t, periodicBackup.frequencyIsTooSmall())
}
