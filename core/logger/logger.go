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
type Logger struct {
	*zap.SugaredLogger
}

// Write logs a message at the Info level and returns the length
// of the given bytes.
func (l *Logger) Write(b []byte) (int, error) {
	l.Info(string(b))
	return len(b), nil
}

// WarnIf logs the error if present.
func (l *Logger) WarnIf(err error) {
	if err != nil {
		l.Warn(err)
	}
}

// ErrorIf logs the error if present.
func (l *Logger) ErrorIf(err error, optionalMsg ...string) {
	if err != nil {
		if len(optionalMsg) > 0 {
			l.Error(errors.Wrap(err, optionalMsg[0]))
		} else {
			l.Error(err)
		}
	}
}

// ErrorIfCalling calls the given function and logs the error of it if there is.
func (l *Logger) ErrorIfCalling(f func() error, optionalMsg ...string) {
	err := f()
	if err != nil {
		e := errors.Wrap(err, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		if len(optionalMsg) > 0 {
			Default.Error(errors.Wrap(e, optionalMsg[0]))
		} else {
			Default.Error(e)
		}
	}
}

func (l *Logger) PanicIf(err error) {
	if err != nil {
		l.Panic(err)
	}
}

// CreateLogger dwisott
func CreateLogger(zl *zap.SugaredLogger) *Logger {
	return &Logger{
		SugaredLogger: zl,
	}
}

// CreateProductionLogger returns a log config for the passed directory
// with the given LogLevel and customizes stdout for pretty printing.
func CreateProductionLogger(
	dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) *Logger {
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

	zl, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{
		SugaredLogger: zl.Sugar(),
	}
}
