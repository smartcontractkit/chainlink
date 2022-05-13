package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// Logger is a minimal subset of smartcontractkit/chainlink/core/logger.Logger implemented by go.uber.org/zap.SugaredLogger
type Logger interface {
	With(args ...interface{}) Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Panicf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Sync() error
}

type Config struct {
	Level zapcore.Level
}

var defaultConfig Config

// New returns a new Logger with the default configuration.
func New() (Logger, error) { return defaultConfig.New() }

// New returns a new Logger for Config.
func (c *Config) New() (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(c.Level)
	core, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &logger{core.Sugar()}, nil
}

// Test returns a new test Logger for tb.
func Test(tb testing.TB) Logger {
	return &logger{zaptest.NewLogger(tb).Sugar()}
}

// TestObserved returns a new test Logger for tb and ObservedLogs at the given Level.
func TestObserved(tb testing.TB, lvl zapcore.Level) (Logger, *observer.ObservedLogs) {
	oCore, logs := observer.New(lvl)
	observe := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, oCore)
	})
	return &logger{zaptest.NewLogger(tb, zaptest.WrapOptions(observe)).Sugar()}, logs
}

// Nop returns a no-op Logger.
func Nop() Logger {
	return &logger{zap.New(zapcore.NewNopCore()).Sugar()}
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) With(args ...interface{}) Logger {
	return &logger{l.SugaredLogger.With(args...)}
}
