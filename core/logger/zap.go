package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

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
	sinksCores := []zapcore.Core{
		newConsoleCore(cfg),
	}

	if cfg.local.ToDisk {
		sinksCores = append(sinksCores, newDiskCore(cfg.local))
	}

	for _, sink := range cfg.sinks {
		sinksCores = append(
			sinksCores,
			zapcore.NewCore(
				zapcore.NewJSONEncoder(makeEncoderConfig(cfg.local)),
				sink,
				zap.LevelEnablerFunc(cfg.Level.Enabled),
			),
		)
	}

	core := zapcore.NewTee(sinksCores...)

	return &zapLogger{config: cfg, SugaredLogger: zap.New(core).Sugar()}, nil
}

func newDiskCore(cfg Config) zapcore.Core {
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
	return zapcore.NewCore(encoder, sink, allLogLevels)
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
	zl, err := newLogger.config.Build()
	if err != nil {
		return nil, err
	}
	zl = zl.WithOptions(zap.AddCallerSkip(l.callerSkip))
	newLogger.SugaredLogger = zl.Named(l.name).Sugar().With(l.fields...)
	return &newLogger, nil
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
