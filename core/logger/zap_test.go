package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/utils"
	utilsmocks "github.com/smartcontractkit/chainlink/core/utils/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLogger_OutOfDiskSpace(t *testing.T) {
	cfg := newZapConfigTest()
	ll, invalid := envvar.LogLevel.Parse()
	assert.Empty(t, invalid)

	cfg.Level.SetLevel(ll)

	maxSize, invalid := envvar.LogFileMaxSize.Parse()
	assert.Empty(t, invalid)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer tmpFile.Close()

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	pollCfg := newDiskPollConfig(1 * time.Second)
	zapCfg := zapLoggerConfig{
		Config: cfg,
		local: Config{
			Dir:            logsDir,
			FileMaxAgeDays: 0,
			FileMaxBackups: 1,
			FileMaxSizeMB:  int(logFileSize / utils.MB),
		},
		diskPollConfig: pollCfg,
		diskLogLevel:   zap.NewAtomicLevelAt(zapcore.DebugLevel),
	}

	t.Run("on logger creation", func(t *testing.T) {
		diskMock := &utilsmocks.DiskStatsProvider{}
		diskMock.On("AvailableSpace", logsDir).Return(maxSize, nil)
		defer diskMock.AssertExpectations(t)

		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.diskStats = diskMock
		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr, close, err := zapCfg.newLogger()
		assert.NoError(t, err)
		defer close()

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created")

		logFile := filepath.Join(zapCfg.local.Dir, LogsFile)
		_, err = ioutil.ReadFile(logFile)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("on logger creation generic error", func(t *testing.T) {
		diskMock := &utilsmocks.DiskStatsProvider{}
		diskMock.On("AvailableSpace", logsDir).Return(utils.FileSize(0), fmt.Errorf("custom error"))
		defer diskMock.AssertExpectations(t)

		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.diskStats = diskMock
		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr, close, err := zapCfg.newLogger()
		assert.NoError(t, err)
		defer close()

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created - generic error")

		logFile := filepath.Join(zapCfg.local.Dir, LogsFile)
		_, err = ioutil.ReadFile(logFile)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("after logger is created", func(t *testing.T) {
		diskMock := &utilsmocks.DiskStatsProvider{}
		diskMock.On("AvailableSpace", logsDir).Return(maxSize*10, nil).Once()
		defer diskMock.AssertExpectations(t)

		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskStats = diskMock
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr, close, err := zapCfg.newLogger()
		assert.NoError(t, err)
		defer close()

		lggr.Debug("writing to disk on test")

		diskMock.On("AvailableSpace", logsDir).Return(maxSize, nil)

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.SetLogLevel(zapcore.WarnLevel)
		lggr.Debug("writing to disk on test again")
		lggr.Warn("writing to disk on test again")

		logFile := filepath.Join(zapCfg.local.Dir, LogsFile)
		b, err := ioutil.ReadFile(logFile)
		assert.NoError(t, err)

		logs := string(b)
		lines := strings.Split(logs, "\n")
		// the last line is a blank line, hence why using len(lines) - 2 makes sense
		actualMessage := lines[len(lines)-2]
		expectedMessage := fmt.Sprintf(
			"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
			zapCfg.local.RequiredDiskSpace(),
			maxSize,
		)

		require.Contains(t, actualMessage, expectedMessage)
	})

	t.Run("after logger is created, recovers disk space", func(t *testing.T) {
		diskMock := &utilsmocks.DiskStatsProvider{}
		diskMock.On("AvailableSpace", logsDir).Return(maxSize*10, nil).Once()
		defer diskMock.AssertExpectations(t)

		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskStats = diskMock
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr, close, err := zapCfg.newLogger()
		assert.NoError(t, err)
		defer close()

		lggr.Debug("test")

		diskMock.On("AvailableSpace", logsDir).Return(maxSize, nil).Once()

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		diskMock.On("AvailableSpace", logsDir).Return(maxSize*12, nil).Once()

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("test again")

		logFile := filepath.Join(zapCfg.local.Dir, LogsFile)
		b, err := ioutil.ReadFile(logFile)
		assert.NoError(t, err)

		logs := string(b)
		lines := strings.Split(logs, "\n")
		expectedMessage := fmt.Sprintf(
			"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
			zapCfg.local.RequiredDiskSpace(),
			maxSize,
		)

		// the last line is a blank line, hence why using len(lines) - N makes sense
		require.Contains(t, lines[len(lines)-4], expectedMessage)
		require.Contains(t, lines[len(lines)-3], "Resuming disk logs, disk has enough space")
		require.Contains(t, lines[len(lines)-2], "test again")
	})
}
