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

type zapDiskPollConfig struct {
	stop     func()
	pollChan <-chan time.Time
}

func newDiskPollConfig(interval time.Duration) zapDiskPollConfig {
	ticker := time.NewTicker(utils.WithJitter(interval))

	return zapDiskPollConfig{
		pollChan: ticker.C,
		stop:     ticker.Stop,
	}
}

func newDiskCore(cfg ZapLoggerConfig) (zapcore.Core, error) {
	availableSpace, err := cfg.diskStats.AvailableSpace(cfg.local.Dir)
	if err != nil {
		return nil, errors.Wrap(err, "error getting disk space available for logging")
	}
	if availableSpace < cfg.local.RequiredDiskSpace {
		return nil, fmt.Errorf(
			"disk space is not enough to log into disk, required disk space: %s, Available disk space: %s",
			cfg.local.RequiredDiskSpace,
			availableSpace,
		)
	}

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
	defer l.config.diskPollConfig.stop()

	for {
		select {
		case <-l.closeDiskPollChan:
			return
		case <-l.config.diskPollConfig.pollChan:
			lvl := zapcore.DebugLevel

			diskUsage, err := l.config.diskStats.AvailableSpace(l.config.local.Dir)
			if err != nil {
				// Will no longer log to disk
				lvl = zapcore.FatalLevel + 1
				l.Warnf("error getting disk space available for logging", "error", err)
			} else if diskUsage < l.config.local.RequiredDiskSpace {
				// Will no longer log to disk
				lvl = zapcore.FatalLevel + 1
				l.Warnf(
					"disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
					l.config.local.RequiredDiskSpace,
					diskUsage,
				)
			}

			l.config.diskLogLevel.SetLevel(lvl)

			if lvl == zapcore.DebugLevel {
				l.Info("resuming disk logs, disk has enough space")
			}

			if l.config.diskLogLvlChan != nil {
				l.config.diskLogLvlChan <- lvl
			}
		}
	}
}
