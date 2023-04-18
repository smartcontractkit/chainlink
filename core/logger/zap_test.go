package logger

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (cfg zapDiskLoggerConfig) newTestLogger(t *testing.T, zcfg zap.Config, cores ...zapcore.Core) Logger {
	lggr, closeLggr, err := cfg.newLogger(zcfg, cores...)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, closeLggr())
	})
	return lggr
}

func TestZapLogger_OutOfDiskSpace(t *testing.T) {
	cfg := newZapConfigBase()
	maxSize := utils.FileSize(5 * utils.MB)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer func() { assert.NoError(t, tmpFile.Close()) }()

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	pollCfg := newDiskPollConfig(1 * time.Second)
	zapCfg := zapDiskLoggerConfig{
		local: Config{
			Dir:            logsDir,
			FileMaxAgeDays: 0,
			FileMaxBackups: 1,
			FileMaxSizeMB:  int(logFileSize / utils.MB),
		},
		diskPollConfig: pollCfg,
	}

	t.Run("on logger creation", func(t *testing.T) {
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.diskSpaceAvailable = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return maxSize, nil
		}
		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := zapCfg.newTestLogger(t, cfg)

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created")

		logFile := zapCfg.local.LogsFile()
		_, err = os.ReadFile(logFile)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("on logger creation generic error", func(t *testing.T) {
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		zapCfg.diskSpaceAvailable = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return 0, nil
		}
		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := zapCfg.newTestLogger(t, cfg)

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created - generic error")

		logFile := zapCfg.local.LogsFile()
		_, err = os.ReadFile(logFile)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("after logger is created", func(t *testing.T) {
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		available := maxSize * 10
		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskSpaceAvailable = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return available, nil
		}
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := zapCfg.newTestLogger(t, cfg)

		lggr.Debug("writing to disk on test")

		available = maxSize

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.SetLogLevel(zapcore.WarnLevel)
		lggr.Debug("writing to disk on test again")
		lggr.Warn("writing to disk on test again")

		logFile := zapCfg.local.LogsFile()
		b, err := os.ReadFile(logFile)
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
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		available := maxSize * 10

		zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
		zapCfg.diskSpaceAvailable = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return available, nil
		}
		zapCfg.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := zapCfg.newTestLogger(t, cfg)

		lggr.Debug("test")

		available = maxSize

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		available = maxSize * 12

		pollChan <- time.Now()
		<-zapCfg.testDiskLogLvlChan

		lggr.Debug("test again")

		logFile := zapCfg.local.LogsFile()
		b, err := os.ReadFile(logFile)
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

func TestZapLogger_LogCaller(t *testing.T) {
	cfg := newZapConfigBase()
	maxSize := utils.FileSize(5 * utils.MB)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer func() { assert.NoError(t, tmpFile.Close()) }()

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	pollCfg := newDiskPollConfig(1 * time.Second)
	zapCfg := zapDiskLoggerConfig{
		local: Config{
			Dir:            logsDir,
			FileMaxAgeDays: 1,
			FileMaxBackups: 1,
			FileMaxSizeMB:  int(logFileSize / utils.MB),
		},
		diskPollConfig: pollCfg,
	}

	pollChan := make(chan time.Time)
	stop := func() {
		close(pollChan)
	}

	zapCfg.testDiskLogLvlChan = make(chan zapcore.Level)
	zapCfg.diskSpaceAvailable = func(path string) (utils.FileSize, error) {
		assert.Equal(t, logsDir, path)
		return maxSize * 10, nil
	}
	zapCfg.diskPollConfig = zapDiskPollConfig{
		stop:     stop,
		pollChan: pollChan,
	}
	zapCfg.local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

	lggr := zapCfg.newTestLogger(t, cfg)

	lggr.Debug("test message with caller")

	pollChan <- time.Now()
	<-zapCfg.testDiskLogLvlChan

	logFile := zapCfg.local.LogsFile()
	b, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	logs := string(b)
	lines := strings.Split(logs, "\n")

	require.Contains(t, lines[0], "logger/zap_test.go:257")
}

func TestZapLogger_Name(t *testing.T) {
	cfg := newZapConfigBase()
	zapCfg := zapDiskLoggerConfig{}

	lggr := zapCfg.newTestLogger(t, cfg)
	require.Equal(t, "", lggr.Name())
	lggr1 := lggr.Named("Lggr1")
	require.Equal(t, "Lggr1", lggr1.Name())
	lggr2 := lggr1.Named("Lggr2")
	require.Equal(t, "Lggr1.Lggr2", lggr2.Name())
}
