// Package logger is used to store details of events in the node.
// Events can be categorized by Debug, Info, Error, Fatal, and Panic.
package logger

import (
	stderr "errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *Logger
	mtx    sync.RWMutex
)

func init() {
	err := zap.RegisterSink("pretty", prettyConsoleSink(os.Stderr))
	if err != nil {
		log.Fatalf("failed to register pretty printer %+v", err)
	}
	err = registerOSSinks()
	if err != nil {
		log.Fatalf("failed to register os specific sinks %+v", err)
	}

	zl, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	SetLogger(zl)
}

func GetLogger() *Logger {
	mtx.RLock()
	defer mtx.RUnlock()
	return logger
}

func prettyConsoleSink(s zap.Sink) func(*url.URL) (zap.Sink, error) {
	return func(*url.URL) (zap.Sink, error) {
		return PrettyConsole{s}, nil
	}
}

// Logger holds a field for the logger interface.
type Logger struct {
	*zap.SugaredLogger
}

// Write logs a message at the Info level and returns the length
// of the given bytes.
func (l *Logger) Write(b []byte) (int, error) {
	l.Info(string(b))
	return len(b), nil
}

// SetLogger sets the internal logger to the given input.
func SetLogger(zl *zap.Logger) {
	mtx.Lock()
	defer mtx.Unlock()
	if logger != nil {
		defer func() {
			if err := logger.Sync(); err != nil {
				if stderr.Unwrap(err).Error() != os.ErrInvalid.Error() &&
					stderr.Unwrap(err).Error() != "inappropriate ioctl for device" &&
					stderr.Unwrap(err).Error() != "bad file descriptor" {
					// logger.Sync() will return 'invalid argument' error when closing file
					log.Fatalf("failed to sync logger %+v", err)
				}
			}
		}()
	}
	logger = &Logger{zl.Sugar()}
}

// CreateProductionLogger returns a log config for the passed directory
// with the given LogLevel and customizes stdout for pretty printing.
func CreateProductionLogger(
	dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) *zap.Logger {
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
	return zl
}

// Infow logs an info message and any additional given information.
func Infow(msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Infow(msg, keysAndValues...)
}

// Debugw logs a debug message and any additional given information.
func Debugw(msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Debugw(msg, keysAndValues...)
}

// Warnw logs a debug message and any additional given information.
func Warnw(msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
func Errorw(msg string, keysAndValues ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Errorw(msg, keysAndValues...)
}

// Infof formats and then logs the message.
func Infof(format string, values ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Info(fmt.Sprintf(format, values...))
}

// Debugf formats and then logs the message.
func Debugf(format string, values ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Debug(fmt.Sprintf(format, values...))
}

// Warnf formats and then logs the message as Warn.
func Warnf(format string, values ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Warn(fmt.Sprintf(format, values...))
}

// Panicf formats and then logs the message before panicking.
func Panicf(format string, values ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Panic(fmt.Sprintf(format, values...))
}

// Info logs an info message.
func Info(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Info(args...)
}

// Debug logs a debug message.
func Debug(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Debug(args...)
}

// Warn logs a message at the warn level.
func Warn(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Warn(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Error(args...)
}

// WarnIf logs the error if present.
func WarnIf(err error) {
	if err != nil {
		mtx.RLock()
		defer mtx.RUnlock()
		logger.Warn(err)
	}
}

// ErrorIf logs the error if present.
func ErrorIf(err error, optionalMsg ...string) {
	if err != nil {
		mtx.RLock()
		defer mtx.RUnlock()
		if len(optionalMsg) > 0 {
			logger.Error(errors.Wrap(err, optionalMsg[0]))
		} else {
			logger.Error(err)
		}
	}
}

// ErrorIfCalling calls the given function and logs the error of it if there is.
func ErrorIfCalling(f func() error, optionalMsg ...string) {
	err := f()
	if err != nil {
		mtx.RLock()
		defer mtx.RUnlock()
		e := errors.Wrap(err, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		if len(optionalMsg) > 0 {
			logger.Error(errors.Wrap(e, optionalMsg[0]))
		} else {
			logger.Error(e)
		}
	}
}

// PanicIf logs the error if present.
func PanicIf(err error) {
	if err != nil {
		mtx.RLock()
		defer mtx.RUnlock()
		logger.Panic(err)
	}
}

// Fatal logs a fatal message then exits the application.
func Fatal(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Fatal(args...)
}

// Errorf logs a message at the error level using Sprintf.
func Errorf(format string, values ...interface{}) {
	Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
func Fatalf(format string, values ...interface{}) {
	Fatal(fmt.Sprintf(format, values...))
}

// Panic logs a panic message then panics.
func Panic(args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	logger.Panic(args...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	mtx.RLock()
	defer mtx.RUnlock()
	return logger.Sync()
}
