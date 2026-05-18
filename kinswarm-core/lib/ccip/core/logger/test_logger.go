package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// TestLogger creates a logger that directs output to PrettyConsole configured
// for test output, and to the buffer testMemoryLog. t is optional.
// Log level is DEBUG by default.
//
// Note: It is not necessary to Sync().
func TestLogger(tb testing.TB) SugaredLogger {
	return testLogger(tb, nil)
}

// TestLoggerObserved creates a logger with an observer that can be used to
// test emitted logs at the given level or above
//
// Note: It is not necessary to Sync().
func TestLoggerObserved(tb testing.TB, lvl zapcore.Level) (Logger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(lvl)
	return testLogger(tb, observedZapCore), observedLogs
}

// testLogger returns a new SugaredLogger for tests. core is optional.
func testLogger(tb testing.TB, core zapcore.Core) SugaredLogger {
	a := zap.NewAtomicLevelAt(zap.DebugLevel)
	opts := []zaptest.LoggerOption{zaptest.Level(a)}
	zapOpts := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)}
	if core != nil {
		zapOpts = append(zapOpts, zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, core)
		}))
	}
	opts = append(opts, zaptest.WrapOptions(zapOpts...))
	l := &zapLogger{
		level:         a,
		SugaredLogger: zaptest.NewLogger(tb, opts...).Sugar(),
	}
	return Sugared(l.With("version", verShaNameStatic()))
}
