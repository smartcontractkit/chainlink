// Package logger is used to store details of events in the node.
// Events can be categorized by Trace, Debug, Info, Error, Fatal, and Panic.
package logger

import (
	"log"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
type Logger interface {
	// With creates a new logger with the given arguments
	With(args ...interface{}) Logger
	// Named creates a new logger sub-scoped for id.
	// Names are inherited and dot-separated.
	//   a := l.Named("a") // id=a logger=a
	//   b := a.Named("b") // id=b logger=a.b
	Named(id string) Logger
	// NamedLevel creates a new Named logger with logLevel.
	NamedLevel(id string, logLevel string) (Logger, error)

	Info(args ...interface{})
	Infof(format string, values ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Debug(args ...interface{})
	Debugf(format string, values ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(args ...interface{})
	Warnf(format string, values ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	// WarnIf logs the error if present.
	WarnIf(err error)

	Error(args ...interface{})
	Errorf(format string, values ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	// ErrorIf logs the error if present.
	ErrorIf(err error, optionalMsg ...string)
	// ErrorIfCalling calls fn and logs any returned error along with func name.
	ErrorIfCalling(fn func() error)

	Fatal(values ...interface{})
	Fatalf(format string, values ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Panic(args ...interface{})
	PanicIf(err error)

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

func (l *zapLogger) Named(id string) Logger {
	newLogger := *l
	newLogger.name = joinName(l.name, id)
	newLogger.SugaredLogger = l.SugaredLogger.Named(id).With("id", id)
	newLogger.fields = copyFields(l.fields, "id", id)
	return &newLogger
}

func (l *zapLogger) withCallerSkip(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return &newLogger
}

func (l *zapLogger) WarnIf(err error) {
	if err != nil {
		l.withCallerSkip(1).Warn(err)
	}
}

func (l *zapLogger) ErrorIf(err error, optionalMsg ...string) {
	if err != nil {
		if len(optionalMsg) > 0 {
			err = errors.Wrap(err, optionalMsg[0])
		}
		l.withCallerSkip(1).Error(err)
	}
}

func (l *zapLogger) ErrorIfCalling(fn func() error) {
	err := fn()
	if err != nil {
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		e := errors.Wrap(err, fnName)
		l.withCallerSkip(1).Error(e)
	}
}

func (l *zapLogger) PanicIf(err error) {
	if err != nil {
		l.withCallerSkip(1).Panic(err)
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

func (l *zapLogger) NamedLevel(id string, logLevel string) (Logger, error) {
	var ll zapcore.Level
	if err := ll.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	zl, err := initLogConfig(l.dir, l.jsonConsole, ll, l.toDisk).Build()
	if err != nil {
		return nil, err
	}

	newLogger := *l
	newLogger.name = joinName(l.name, id)
	newLogger.fields = copyFields(l.fields, "id", id)
	newLogger.SugaredLogger = zl.Named(newLogger.name).Sugar().With(newLogger.fields...)
	return &newLogger, nil
}
