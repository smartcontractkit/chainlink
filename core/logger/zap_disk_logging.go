package logger

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func (c zapDiskPollConfig) isSet() bool {
	return c.stop != nil || c.pollChan != nil
}

func newDiskPollConfig(interval time.Duration) zapDiskPollConfig {
	ticker := time.NewTicker(utils.WithJitter(interval))

	return zapDiskPollConfig{
		pollChan: ticker.C,
		stop:     ticker.Stop,
	}
}

var _ Logger = &zapDiskLogger{}

type zapDiskLogger struct {
	zapLogger
	config            Config
	diskLogLevel      zap.AtomicLevel
	pollDiskSpaceStop chan struct{}
	pollDiskSpaceDone chan struct{}
}

func (l *zapDiskLogger) pollDiskSpace() {
	defer l.config.diskPollConfig.stop()
	defer close(l.pollDiskSpaceDone)

	for {
		select {
		case <-l.pollDiskSpaceStop:
			return
		case <-l.config.diskPollConfig.pollChan:
			lvl := zapcore.DebugLevel

			diskUsage, err := l.config.DiskSpaceAvailable(l.config.Dir)
			if err != nil {
				// Will no longer log to disk
				lvl = disabledLevel
				l.Warnw("Error getting disk space available for logging", "err", err)
			} else if diskUsage < l.config.RequiredDiskSpace() {
				// Will no longer log to disk
				lvl = disabledLevel
				l.Warnf(
					"Disk space is not enough to log into disk any longer, required disk space: %s, Available disk space: %s",
					l.config.RequiredDiskSpace(),
					diskUsage,
				)
			}

			lvlBefore := l.diskLogLevel.Level()

			l.diskLogLevel.SetLevel(lvl)

			if lvlBefore == disabledLevel && lvl == zapcore.DebugLevel {
				l.Info("Resuming disk logs, disk has enough space")
			}

			if l.config.testDiskLogLvlChan != nil {
				l.config.testDiskLogLvlChan <- lvl
			}
		}
	}
}

func newRotatingFileLogger(zcfg zap.Config, c Config, cores ...zapcore.Core) (*zapDiskLogger, func() error, error) {
	defaultCore, defaultCloseFn, err := newDefaultLoggingCore(zcfg, c.UnixTS)
	if err != nil {
		return nil, nil, err
	}
	cores = append(cores, defaultCore)

	diskLogLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	diskCore, diskErr := newDiskCore(diskLogLevel, c)
	if diskErr != nil {
		defaultCloseFn()
		return nil, nil, diskErr
	}
	cores = append(cores, diskCore)

	core := zapcore.NewTee(cores...)
	l, diskCloseFn, err := newLoggerForCore(zcfg, core)
	if err != nil {
		defaultCloseFn()
		return nil, nil, err
	}

	lggr := &zapDiskLogger{
		config: c,

		pollDiskSpaceStop: make(chan struct{}),
		pollDiskSpaceDone: make(chan struct{}),
		zapLogger:         *l,
		diskLogLevel:      diskLogLevel,
	}

	go lggr.pollDiskSpace()

	closeLogger := sync.OnceValue(func() error {
		defer defaultCloseFn()
		defer diskCloseFn()

		close(lggr.pollDiskSpaceStop)
		<-lggr.pollDiskSpaceDone

		return lggr.Sync()
	})

	return lggr, closeLogger, err
}
