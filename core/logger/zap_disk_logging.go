package logger

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// `Fatal` is the max log level allowed, so log levels like `Panic` or `Critical` won't be logged to disk if this is set.
	disabledLevel = zapcore.FatalLevel + 1

	diskPollInterval = 1 * time.Minute
)

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

func (cfg zapLoggerConfig) newDiskCore() (zapcore.Core, error) {
	availableSpace, err := cfg.diskStats.AvailableSpace(cfg.local.Dir)
	if err != nil || availableSpace < cfg.local.RequiredDiskSpace() {
		// Won't log to disk if the directory is not found or there's not enough disk space
		cfg.diskLogLevel.SetLevel(disabledLevel)
	}

	var (
		encoder = zapcore.NewConsoleEncoder(makeEncoderConfig(cfg.local))
		sink    = zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFileURI(cfg.local.Dir),
			MaxSize:    cfg.local.FileMaxSizeMB,
			MaxAge:     cfg.local.FileMaxAgeDays,
			MaxBackups: cfg.local.FileMaxBackups,
			Compress:   true,
		})
		allLogLevels = zap.LevelEnablerFunc(cfg.diskLogLevel.Enabled)
	)

	return zapcore.NewCore(encoder, sink, allLogLevels), nil
}

func (l *zapLogger) pollDiskSpace() {
	defer l.config.diskPollConfig.stop()
	defer close(l.pollDiskSpaceDone)

	for {
		select {
		case <-l.pollDiskSpaceStop:
			return
		case <-l.config.diskPollConfig.pollChan:
			lvl := zapcore.DebugLevel

			diskUsage, err := l.config.diskStats.AvailableSpace(l.config.local.Dir)
			if err != nil {
				// Will no longer log to disk
				lvl = disabledLevel
				l.Warnw("Error getting disk space available for logging", "err", err)
			} else if diskUsage < l.config.local.RequiredDiskSpace() {
				// Will no longer log to disk
				lvl = disabledLevel
				l.Warnf(
					"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
					l.config.local.RequiredDiskSpace(),
					diskUsage,
				)
			}

			lvlBefore := l.config.diskLogLevel.Level()

			l.config.diskLogLevel.SetLevel(lvl)

			if lvlBefore == disabledLevel && lvl == zapcore.DebugLevel {
				l.Info("Resuming disk logs, disk has enough space")
			}

			if l.config.testDiskLogLvlChan != nil {
				l.config.testDiskLogLvlChan <- lvl
			}
		}
	}
}
