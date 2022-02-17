package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Logger = &zapLogger{}

// ZapLoggerConfig defines the struct that serves as config when spinning up a the zap logger
type ZapLoggerConfig struct {
	zap.Config
	local Config
	sinks []zapcore.WriteSyncer
}

type zapLogger struct {
	*zap.SugaredLogger
	config     ZapLoggerConfig
	name       string
	fields     []interface{}
	callerSkip int
}

func newZapLogger(cfg ZapLoggerConfig) (Logger, error) {
	cores := []zapcore.Core{
		newConsoleCore(cfg),
	}
	newCores, err := newCores(cfg)
	cores = append(cores, newCores...)

	core := zapcore.NewTee(cores...)

	return &zapLogger{
		config:        cfg,
		SugaredLogger: zap.New(core).Sugar(),
	}, err
}

func newCores(cfg ZapLoggerConfig) ([]zapcore.Core, error) {
	var err error
	var cores []zapcore.Core

	if cfg.local.ToDisk {
		core, diskErr := newDiskCore(cfg.local)
		if diskErr == nil {
			cores = append(cores, core)
		}
		err = diskErr
	}

	for _, sink := range cfg.sinks {
		cores = append(
			cores,
			zapcore.NewCore(
				zapcore.NewJSONEncoder(makeEncoderConfig(cfg.local)),
				sink,
				zap.LevelEnablerFunc(cfg.Level.Enabled),
			),
		)
	}

	return cores, err
}

func newDiskCore(cfg Config) (zapcore.Core, error) {
	diskUsage, err := disk.Usage(cfg.Dir)
	if err != nil {
		return nil, errors.Wrap(err, "error getting disk space available for logging")
	}
	diskSpaceAvailable := utils.FileSize(diskUsage.Free)
	requiredDiskSpace := utils.FileSize(cfg.DiskMaxSizeBeforeRotate * (cfg.DiskMaxBackupsBeforeDelete + 1))
	if diskSpaceAvailable < requiredDiskSpace {
		return nil, fmt.Errorf(
			"disk space is not enough to log into disk, Required disk space: %s, Available disk space: %s",
			requiredDiskSpace,
			diskSpaceAvailable,
		)
	}

	var (
		encoder = zapcore.NewConsoleEncoder(makeEncoderConfig(cfg))
		sink    = zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFileURI(cfg.Dir),
			MaxSize:    cfg.DiskMaxSizeBeforeRotate,
			MaxAge:     cfg.DiskMaxAgeBeforeDelete,
			MaxBackups: cfg.DiskMaxBackupsBeforeDelete,
			Compress:   true,
		})
		allLogLevels = zap.LevelEnablerFunc(zap.DebugLevel.Enabled)
	)

	return zapcore.NewCore(encoder, sink, allLogLevels), nil
}

func newConsoleCore(cfg ZapLoggerConfig) zapcore.Core {
	filteredLogLevels := zap.LevelEnablerFunc(cfg.Level.Enabled)

	encoder := zapcore.NewJSONEncoder(makeEncoderConfig(cfg.local))

	var sink zap.Sink
	if !cfg.local.JsonConsole {
		sink = PrettyConsole{os.Stderr}
	}

	return zapcore.NewCore(encoder, sink, filteredLogLevels)
}

func makeEncoderConfig(cfg Config) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	if !cfg.UnixTS {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	encoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		if l == zapcore.DPanicLevel {
			enc.AppendString("crit")
		} else {
			zapcore.LowercaseLevelEncoder(l, enc)
		}
	}

	return encoderConfig
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.config.Level.SetLevel(lvl)
}

func (l *zapLogger) With(args ...interface{}) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.With(args...)
	newLogger.fields = copyFields(l.fields, args...)
	return &newLogger
}

// copyFields returns a copy of fields with add appended.
func copyFields(fields []interface{}, add ...interface{}) []interface{} {
	f := make([]interface{}, 0, len(fields)+len(add))
	f = append(f, fields...)
	f = append(f, add...)
	return f
}

func joinName(old, new string) string {
	if old == "" {
		return new
	}
	return old + "." + new
}

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.name = joinName(l.name, name)
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	newLogger.Trace("Named logger created")
	return &newLogger
}

func (l *zapLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	newLogger := *l
	newLogger.config.Level = zap.NewAtomicLevelAt(lvl)

	cores := []zapcore.Core{
		// The console core is what we want to be unique per root, so we spin a new one here
		newConsoleCore(newLogger.config),
	}
	extraCores, err := newCores(newLogger.config)
	cores = append(cores, extraCores...)
	core := zap.New(zapcore.NewTee(cores...)).WithOptions(zap.AddCallerSkip(l.callerSkip))

	newLogger.SugaredLogger = core.Named(l.name).Sugar().With(l.fields...)

	return &newLogger, err
}

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.Helper(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.Helper(1).Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *zapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return nil
	}
	var msg string
	if uw := errors.Unwrap(err); uw != nil {
		msg = uw.Error()
	} else {
		msg = err.Error()
	}
	switch msg {
	case os.ErrInvalid.Error(), "bad file descriptor",
		"inappropriate ioctl for device":
		return nil
	}
	return err
}

func (l *zapLogger) Recover(panicErr interface{}) {
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}
