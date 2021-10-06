// Package logger is used to store details of events in the node.
// Events can be categorized by Trace, Debug, Info, Error, Fatal, and Panic.
package logger

import (
	"fmt"
	"log"
	"reflect"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
type Logger interface {
	// With creates a new logger with the given arguments
	With(args ...interface{}) Logger
	// Named creates a new logger sub-scoped with name.
	// Names are inherited and dot-separated.
	//   a := l.Named("a") // logger=a
	//   b := a.Named("b") // logger=a.b
	Named(name string) Logger
	// NamedLevel creates a new Named logger with logLevel.
	NamedLevel(id string, logLevel string) (Logger, error)

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	// WarnIf logs the error if present.
	WarnIf(err error, msg string)
	// ErrorIf logs the error if present.
	ErrorIf(err error, msg string)
	PanicIf(err error, msg string)

	// ErrorIfCalling calls fn and logs any returned error along with func name.
	ErrorIfCalling(fn func() error)

	Sync() error

	// withCallerSkip creates a new logger with the number of callers skipped by
	// caller annotation increased by add. For wrappers to use internally.
	withCallerSkip(add int) Logger
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	dir         string
	jsonConsole bool
	toDisk      bool
	name        string
	fields      []interface{}
}

// Constants for service names for package specific logging configuration
const (
	HeadTracker = "head_tracker"
	FluxMonitor = "fluxmonitor"
	Keeper      = "keeper"
)

func GetLogServices() []string {
	return []string{
		HeadTracker,
		FluxMonitor,
		Keeper,
	}
}

func (l *zapLogger) Write(b []byte) (int, error) {
	l.Info(string(b))
	return len(b), nil
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
	return &newLogger
}

func (l *zapLogger) withCallerSkip(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return &newLogger
}

func (l *zapLogger) WarnIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Warn(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfCalling(fn func() error) {
	err := fn()
	if err != nil {
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		l.withCallerSkip(1).Errorw(fmt.Sprintf("Error calling %s", fnName), "err", err)
	}
}

func (l *zapLogger) PanicIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Panicw(msg, "err", err)
	}
}

// initLogConfig builds a zap.Config for a logger
func initLogConfig(dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) zap.Config {
	config := zap.NewProductionConfig()
	if !jsonConsole {
		config.OutputPaths = []string{"pretty://console"}
	}
	if toDisk {
		destination := logFileURI(dir)
		config.OutputPaths = append(config.OutputPaths, destination)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, destination)
	}
	config.Level.SetLevel(lvl)
	return config
}

// CreateProductionLogger returns a log config for the passed directory
// with the given LogLevel and customizes stdout for pretty printing.
func CreateProductionLogger(dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) Logger {
	zl, err := initLogConfig(dir, jsonConsole, lvl, toDisk).Build()
	if err != nil {
		log.Fatal(err)
	}
	return &zapLogger{
		SugaredLogger: zl.Sugar(),
		dir:           dir,
		jsonConsole:   jsonConsole,
		toDisk:        toDisk,
	}
}

func (l *zapLogger) NamedLevel(name string, logLevel string) (Logger, error) {
	var ll zapcore.Level
	if err := ll.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	zl, err := initLogConfig(l.dir, l.jsonConsole, ll, l.toDisk).Build()
	if err != nil {
		return nil, err
	}

	newLogger := *l
	newLogger.name = joinName(l.name, name)
	newLogger.SugaredLogger = zl.Named(newLogger.name).Sugar().With(l.fields...)
	return &newLogger, nil
}
