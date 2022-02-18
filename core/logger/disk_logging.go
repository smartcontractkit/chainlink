package logger

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var diskPollInterval = 1 * time.Minute

func newDiskCore(cfg ZapLoggerConfig) (zapcore.Core, error) {
	availableSpace, err := cfg.diskStats.AvailableSpace(cfg.local.Dir)
	if err != nil {
		return nil, errors.Wrap(err, "error getting disk space available for logging")
	}
	if availableSpace < cfg.local.RequiredDiskSpace {
		return nil, fmt.Errorf(
			"disk space is not enough to log into disk, Required disk space: %s, Available disk space: %s",
			cfg.local.RequiredDiskSpace,
			availableSpace,
		)
	}

	cfg.diskLogLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	var (
		encoder = zapcore.NewConsoleEncoder(makeEncoderConfig(cfg.local))
		sink    = zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFileURI(cfg.local.Dir),
			MaxSize:    cfg.local.DiskMaxSizeBeforeRotate,
			MaxAge:     cfg.local.DiskMaxAgeBeforeDelete,
			MaxBackups: cfg.local.DiskMaxBackupsBeforeDelete,
			Compress:   true,
		})
		allLogLevels = zap.LevelEnablerFunc(cfg.diskLogLevel.Enabled)
	)

	return zapcore.NewCore(encoder, sink, allLogLevels), nil
}

func (l *zapLogger) pollDiskSpace() {
	defer close(l.closeDiskPollChan)

	ticker := time.NewTicker(utils.WithJitter(diskPollInterval))
	defer ticker.Stop()

	for {
		select {
		case <-l.closeDiskPollChan:
			return
		case <-ticker.C:
			diskUsage, err := l.config.diskStats.AvailableSpace(l.config.local.Dir)
			if err != nil {
				l.Fatalw("error getting disk space available for logging", "error", err)
				// Will no longer log to disk
				l.config.diskLogLevel.SetLevel(zapcore.FatalLevel + 1)
			}

			if diskUsage < l.config.local.RequiredDiskSpace {
				l.Fatalf(
					"disk space is not enough to log into disk, Required disk space: %s, Available disk space: %s",
					l.config.local.RequiredDiskSpace,
					diskUsage,
				)
				// Will no longer log to disk
				l.config.diskLogLevel.SetLevel(zapcore.FatalLevel + 1)
			} else {
				// Will resume disk logs
				l.config.diskLogLevel.SetLevel(zapcore.DebugLevel)
			}
		}
	}
}
