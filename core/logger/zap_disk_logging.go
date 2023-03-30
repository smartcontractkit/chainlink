package logger

import (
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

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

var _ Logger = &zapDiskLogger{}

// zapDiskLoggerConfig defines the struct that serves as config when spinning up a the zap logger
type zapDiskLoggerConfig struct {
	local              Config
	diskSpaceAvailable diskSpaceAvailableFn
	diskPollConfig     zapDiskPollConfig

	// This is for tests only
	testDiskLogLvlChan chan zapcore.Level
}

func (cfg zapDiskLoggerConfig) newDiskCore(diskLogLevel zap.AtomicLevel) (zapcore.Core, error) {
	availableSpace, err := cfg.diskSpaceAvailable(cfg.local.Dir)
	if err != nil || availableSpace < cfg.local.RequiredDiskSpace() {
		// Won't log to disk if the directory is not found or there's not enough disk space
		diskLogLevel.SetLevel(disabledLevel)
	}

	var (
		encoder = zapcore.NewConsoleEncoder(makeEncoderConfig(cfg.local))
		sink    = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.local.logFileURI(),
			MaxSize:    cfg.local.FileMaxSizeMB,
			MaxAge:     cfg.local.FileMaxAgeDays,
			MaxBackups: cfg.local.FileMaxBackups,
			Compress:   true,
		})
		allLogLevels = zap.LevelEnablerFunc(diskLogLevel.Enabled)
	)

	return zapcore.NewCore(encoder, sink, allLogLevels), nil
}

type zapDiskLogger struct {
	zapLogger
	config            zapDiskLoggerConfig
	diskLogLevel      zap.AtomicLevel
	pollDiskSpaceStop chan struct{}
	pollDiskSpaceDone chan struct{}
}

func (cfg zapDiskLoggerConfig) newLogger(zcfg zap.Config, cores ...zapcore.Core) (Logger, func() error, error) {
	newCore, errWriter, err := cfg.newCore(zcfg)
	if err != nil {
		return nil, nil, err
	}
	cores = append(cores, newCore)
	diskLogLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if cfg.local.DebugLogsToDisk() {
		diskCore, diskErr := cfg.newDiskCore(diskLogLevel)
		if diskErr != nil {
			return nil, nil, diskErr
		}
		cores = append(cores, diskCore)
	}

	core := zapcore.NewTee(cores...)
	lggr := &zapDiskLogger{
		config:            cfg,
		pollDiskSpaceStop: make(chan struct{}),
		pollDiskSpaceDone: make(chan struct{}),
		zapLogger: zapLogger{
			level:         zcfg.Level,
			SugaredLogger: zap.New(core, zap.ErrorOutput(errWriter), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar(),
		},
		diskLogLevel: diskLogLevel,
	}

	if cfg.local.DebugLogsToDisk() {
		go lggr.pollDiskSpace()
	}

	var once sync.Once
	closeLogger := func() error {
		once.Do(func() {
			if cfg.local.DebugLogsToDisk() {
				close(lggr.pollDiskSpaceStop)
				<-lggr.pollDiskSpaceDone
			}
		})

		return lggr.Sync()
	}

	return lggr, closeLogger, err
}

func (cfg zapDiskLoggerConfig) newCore(zcfg zap.Config) (zapcore.Core, zapcore.WriteSyncer, error) {
	encoder := zapcore.NewJSONEncoder(makeEncoderConfig(cfg.local))

	sink, closeOut, err := zap.Open(zcfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}

	errSink, _, err := zap.Open(zcfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}

	if zcfg.Level == (zap.AtomicLevel{}) {
		return nil, nil, errors.New("missing Level")
	}

	filteredLogLevels := zap.LevelEnablerFunc(zcfg.Level.Enabled)

	return zapcore.NewCore(encoder, sink, filteredLogLevels), errSink, nil
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

			diskUsage, err := l.config.diskSpaceAvailable(l.config.local.Dir)
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
