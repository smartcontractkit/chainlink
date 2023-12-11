package logger

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// Logger is a minimal subset of smartcontractkit/chainlink/core/logger.Logger implemented by go.uber.org/zap.SugaredLogger
type Logger interface {
	Name() string

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
	return NewWith(func(cfg *zap.Config) {
		cfg.Level.SetLevel(c.Level)
	})
}

// NewWith returns a new Logger from a modified [zap.Config].
func NewWith(cfgFn func(*zap.Config)) (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfgFn(&cfg)
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

// TestSugared returns a new test SugaredLogger.
func TestSugared(tb testing.TB) SugaredLogger {
	return Sugared(&logger{zaptest.NewLogger(tb).Sugar()})
}

// TestObserved returns a new test Logger for tb and ObservedLogs at the given Level.
func TestObserved(tb testing.TB, lvl zapcore.Level) (Logger, *observer.ObservedLogs) {
	sl, logs := testObserved(tb, lvl)
	return &logger{sl}, logs
}

func testObserved(tb testing.TB, lvl zapcore.Level) (*zap.SugaredLogger, *observer.ObservedLogs) {
	oCore, logs := observer.New(lvl)
	observe := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, oCore)
	})
	return zaptest.NewLogger(tb, zaptest.WrapOptions(observe, zap.AddCaller())).Sugar(), logs
}

// Nop returns a no-op Logger.
func Nop() Logger {
	return &logger{zap.New(zapcore.NewNopCore()).Sugar()}
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) with(args ...interface{}) Logger {
	return &logger{l.SugaredLogger.With(args...)}
}

func (l *logger) named(name string) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	return &newLogger
}

func (l *logger) Name() string {
	return l.Desugar().Name()
}

func (l *logger) helper(skip int) Logger {
	return &logger{l.sugaredHelper(skip)}
}

func (l *logger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.WithOptions(zap.AddCallerSkip(skip))
}

// With returns a Logger with keyvals, if 'l' has a method `With(...interface{}) L`, where L implements Logger, otherwise it returns l.
func With(l Logger, keyvals ...interface{}) Logger {
	switch t := l.(type) {
	case *logger:
		return t.with(keyvals...)
	}

	method := reflect.ValueOf(l).MethodByName("With")
	if method == (reflect.Value{}) {
		return l // not available
	}
	if ret := method.CallSlice([]reflect.Value{reflect.ValueOf(keyvals)}); len(ret) == 1 {
		nl, ok := ret[0].Interface().(Logger)
		if ok {
			return nl
		}
	}
	return l
}

// Named returns a logger with name 'n', if 'l' has a method `Named(string) L`, where L implements Logger, otherwise it returns l.
func Named(l Logger, n string) Logger {
	switch t := l.(type) {
	case *logger:
		return t.named(n)
	}

	method := reflect.ValueOf(l).MethodByName("Named")
	if method == (reflect.Value{}) {
		return l // not available
	}
	if ret := method.Call([]reflect.Value{reflect.ValueOf(n)}); len(ret) == 1 {
		nl, ok := ret[0].Interface().(Logger)
		if ok {
			return nl
		}
	}
	return l
}

// Helper returns a logger with 'skip' levels of callers skipped, if 'l' has a method `Helper(int) L`, where L implements Logger, otherwise it returns l.
// See [zap.AddCallerSkip]
func Helper(l Logger, skip int) Logger {
	switch t := l.(type) {
	case *logger:
		return t.helper(skip)
	}

	method := reflect.ValueOf(l).MethodByName("Helper")
	if method == (reflect.Value{}) {
		return l // not available
	}
	if ret := method.Call([]reflect.Value{reflect.ValueOf(skip)}); len(ret) == 1 {
		nl, ok := ret[0].Interface().(Logger)
		if ok {
			return nl
		}
	}
	return l
}

// Deprecated: instead use [SugaredLogger.Critical]:
//
//	Sugared(l).Critical(args...)
func Critical(l Logger, args ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Critical(args...)
}

// Deprecated: instead use [SugaredLogger.Criticalf]:
//
//	Sugared(l).Criticalf(args...)
func Criticalf(l Logger, format string, values ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Criticalf(format, values...)
}

// Deprecated: instead use [SugaredLogger.Criticalw]:
//
//	Sugared(l).Criticalw(args...)
func Criticalw(l Logger, msg string, keysAndValues ...interface{}) {
	s := &sugared{Logger: l, h: Helper(l, 2)}
	s.Criticalw(msg, keysAndValues...)
}
