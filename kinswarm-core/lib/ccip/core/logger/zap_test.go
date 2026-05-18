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
	"go.uber.org/zap/zapcore"
)

func newTestLogger(t *testing.T, cfg Config) Logger {
	lggr, closeFn := cfg.New()
	t.Cleanup(func() {
		assert.NoError(t, closeFn())
	})
	return lggr
}

func TestZapLogger_OutOfDiskSpace(t *testing.T) {
	maxSize := utils.FileSize(5 * utils.MB)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer func() { assert.NoError(t, tmpFile.Close()) }()

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	pollCfg := newDiskPollConfig(1 * time.Second)

	local := Config{
		Dir:            logsDir,
		FileMaxAgeDays: 0,
		FileMaxBackups: 1,
		FileMaxSizeMB:  int(logFileSize / utils.MB),

		diskPollConfig:     pollCfg,
		testDiskLogLvlChan: make(chan zapcore.Level),
	}

	t.Run("on logger creation", func(t *testing.T) {
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		local.diskSpaceAvailableFn = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return maxSize, nil
		}
		local.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := newTestLogger(t, local)

		pollChan <- time.Now()
		<-local.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created")

		logFile := local.LogsFile()
		_, err = os.ReadFile(logFile)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("on logger creation generic error", func(t *testing.T) {
		pollChan := make(chan time.Time)
		stop := func() {
			close(pollChan)
		}

		local.diskSpaceAvailableFn = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return 0, nil
		}
		local.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := newTestLogger(t, local)

		pollChan <- time.Now()
		<-local.testDiskLogLvlChan

		lggr.Debug("trying to write to disk when the disk logs should not be created - generic error")

		logFile := local.LogsFile()
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
		local.diskSpaceAvailableFn = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return available, nil
		}
		local.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := newTestLogger(t, local)

		lggr.Debug("writing to disk on test")

		available = maxSize

		pollChan <- time.Now()
		<-local.testDiskLogLvlChan

		lggr.SetLogLevel(zapcore.WarnLevel)
		lggr.Debug("writing to disk on test again")
		lggr.Warn("writing to disk on test again")

		logFile := local.LogsFile()
		b, err := os.ReadFile(logFile)
		assert.NoError(t, err)

		logs := string(b)
		lines := strings.Split(logs, "\n")
		// the last line is a blank line, hence why using len(lines) - 2 makes sense
		actualMessage := lines[len(lines)-2]
		expectedMessage := fmt.Sprintf(
			"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
			local.RequiredDiskSpace(),
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

		local.diskSpaceAvailableFn = func(path string) (utils.FileSize, error) {
			assert.Equal(t, logsDir, path)
			return available, nil
		}
		local.diskPollConfig = zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		}
		local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

		lggr := newTestLogger(t, local)

		lggr.Debug("test")

		available = maxSize

		pollChan <- time.Now()
		<-local.testDiskLogLvlChan

		available = maxSize * 12

		pollChan <- time.Now()
		<-local.testDiskLogLvlChan

		lggr.Debug("test again")

		logFile := local.LogsFile()
		b, err := os.ReadFile(logFile)
		assert.NoError(t, err)

		logs := string(b)
		lines := strings.Split(logs, "\n")
		expectedMessage := fmt.Sprintf(
			"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
			local.RequiredDiskSpace(),
			maxSize,
		)

		// the last line is a blank line, hence why using len(lines) - N makes sense
		require.Contains(t, lines[len(lines)-4], expectedMessage)
		require.Contains(t, lines[len(lines)-3], "Resuming disk logs, disk has enough space")
		require.Contains(t, lines[len(lines)-2], "test again")
	})
}

func TestZapLogger_LogCaller(t *testing.T) {
	maxSize := utils.FileSize(5 * utils.MB)

	logsDir := t.TempDir()
	tmpFile, err := os.CreateTemp(logsDir, "*")
	assert.NoError(t, err)
	defer func() { assert.NoError(t, tmpFile.Close()) }()

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	pollChan := make(chan time.Time)
	stop := func() {
		close(pollChan)
	}
	local := Config{
		Dir:            logsDir,
		FileMaxAgeDays: 1,
		FileMaxBackups: 1,
		FileMaxSizeMB:  int(logFileSize / utils.MB),

		diskPollConfig: zapDiskPollConfig{
			stop:     stop,
			pollChan: pollChan,
		},
		testDiskLogLvlChan: make(chan zapcore.Level),
	}

	local.diskSpaceAvailableFn = func(path string) (utils.FileSize, error) {
		assert.Equal(t, logsDir, path)
		return maxSize * 10, nil
	}
	local.FileMaxSizeMB = int(maxSize/utils.MB) * 2

	lggr := newTestLogger(t, local)

	lggr.Debug("test message with caller")

	pollChan <- time.Now()
	<-local.testDiskLogLvlChan

	logFile := local.LogsFile()
	b, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	logs := string(b)
	lines := strings.Split(logs, "\n")

	require.Contains(t, lines[0], "logger/zap_test.go:246")
}

func TestZapLogger_Name(t *testing.T) {
	cfg := Config{}
	lggr := newTestLogger(t, cfg)
	require.Equal(t, "", lggr.Name())
	lggr1 := lggr.Named("Lggr1")
	require.Equal(t, "Lggr1", lggr1.Name())
	lggr2 := lggr1.Named("Lggr2")
	require.Equal(t, "Lggr1.Lggr2", lggr2.Name())
}
