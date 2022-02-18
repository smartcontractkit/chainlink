package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/utils"
	utilsmocks "github.com/smartcontractkit/chainlink/core/utils/mocks"
	"github.com/test-go/testify/require"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLogger_OutOfDiskSpace(t *testing.T) {
	cfg := newTestConfig()
	ll, invalid := envvar.LogLevel.ParseLogLevel()
	assert.Empty(t, invalid)

	cfg.Level.SetLevel(ll)

	maxSize, invalid := envvar.LogFileMaxSize.ParseFileSize()
	assert.Empty(t, invalid)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer tmpFile.Close()

	diskMock := &utilsmocks.DiskStatsProvider{}

	t.Run("on logger creation", func(t *testing.T) {
		diskMock.On("AvailableSpace", logsDir).Return(maxSize, nil)
		defer diskMock.AssertExpectations(t)

		requiredDiskSpace := utils.FileSize(int(maxSize) * 2)
		_, err = newZapLogger(ZapLoggerConfig{
			Config: cfg,
			local: Config{
				Dir:                        logsDir,
				ToDisk:                     true,
				DiskMaxSizeBeforeRotate:    1,
				DiskMaxAgeBeforeDelete:     0,
				DiskMaxBackupsBeforeDelete: 1,
				RequiredDiskSpace:          requiredDiskSpace,
			},
			diskStats:    diskMock,
			diskLogLevel: zap.NewAtomicLevelAt(zapcore.DebugLevel),
		})

		expectedError := fmt.Sprintf(
			"disk space is not enough to log into disk, Required disk space: %s, Available disk space: %s",
			requiredDiskSpace,
			maxSize,
		)

		require.Error(t, err)
		require.Equal(t, err.Error(), expectedError)
	})
}
