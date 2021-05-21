package logger

import (
	"errors"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

var (
	// Default logger for use throughout the project.
	// All the package-level functions are calling Default.
	Default *Logger
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
	SetLogger(&Logger{
		SugaredLogger: zl.Sugar(),
	})
}

// SetLogger sets the internal logger to the given input.
//
// DEPRECATED this method is deprecated because it leads to race conditions.
// Instead, you should fork the logger.Default instance to create a new logger
// for your module.
// Eg: logger.Default.Named("<my-package-name>")
func SetLogger(newLogger *Logger) {
	if Default != nil {
		defer func() {
			if err := Default.Sync(); err != nil {
				if errors.Unwrap(err).Error() != os.ErrInvalid.Error() &&
					errors.Unwrap(err).Error() != "inappropriate ioctl for device" &&
					errors.Unwrap(err).Error() != "bad file descriptor" {
					// logger.Sync() will return 'invalid argument' error when closing file
					log.Fatalf("failed to sync logger %+v", err)
				}
			}
		}()
	}
	Default = newLogger
}

// Infow logs an info message and any additional given information.
func Infow(msg string, keysAndValues ...interface{}) {
	Default.Infow(msg, keysAndValues...)
}

// Debugw logs a debug message and any additional given information.
func Debugw(msg string, keysAndValues ...interface{}) {
	Default.Debugw(msg, keysAndValues...)
}

// Tracew is a shim stand-in for when we have real trace-level logging support
func Tracew(msg string, keysAndValues ...interface{}) {
	// Zap does not support trace logging just yet
	Default.Debugw("TRACE: "+msg, keysAndValues...)
}

// Warnw logs a debug message and any additional given information.
func Warnw(msg string, keysAndValues ...interface{}) {
	Default.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
func Errorw(msg string, keysAndValues ...interface{}) {
	Default.Errorw(msg, keysAndValues...)
}

// Logs and returns a new error
func NewErrorw(msg string, keysAndValues ...interface{}) error {
	Default.Errorw(msg, keysAndValues...)
	return errors.New(msg)
}

// Infof formats and then logs the message.
func Infof(format string, values ...interface{}) {
	Default.Info(fmt.Sprintf(format, values...))
}

// Debugf formats and then logs the message.
func Debugf(format string, values ...interface{}) {
	Default.Debug(fmt.Sprintf(format, values...))
}

// Tracef is a shim stand-in for when we have real trace-level logging support
func Tracef(format string, values ...interface{}) {
	Default.Debug("TRACE: " + fmt.Sprintf(format, values...))
}

// Warnf formats and then logs the message as Warn.
func Warnf(format string, values ...interface{}) {
	Default.Warn(fmt.Sprintf(format, values...))
}

// Panicf formats and then logs the message before panicking.
func Panicf(format string, values ...interface{}) {
	Default.Panic(fmt.Sprintf(format, values...))
}

// Info logs an info message.
func Info(args ...interface{}) {
	Default.Info(args...)
}

// Debug logs a debug message.
func Debug(args ...interface{}) {
	Default.Debug(args...)
}

// Trace is a shim stand-in for when we have real trace-level logging support
func Trace(args ...interface{}) {
	Default.Debug(append([]interface{}{"TRACE: "}, args...))
}

// Warn logs a message at the warn level.
func Warn(args ...interface{}) {
	Default.Warn(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	Default.Error(args...)
}

func WarnIf(err error) {
	Default.WarnIf(err)
}

func ErrorIf(err error, optionalMsg ...string) {
	Default.ErrorIf(err, optionalMsg...)
}

func ErrorIfCalling(f func() error, optionalMsg ...string) {
	Default.ErrorIfCalling(f, optionalMsg...)
}

// Fatal logs a fatal message then exits the application.
func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

// Errorf logs a message at the error level using Sprintf.
func Errorf(format string, values ...interface{}) {
	Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
func Fatalf(format string, values ...interface{}) {
	Fatal(fmt.Sprintf(format, values...))
}

// Fatalw logs a message and exits the application
func Fatalw(msg string, keysAndValues ...interface{}) {
	Default.Fatalw(msg, keysAndValues...)
}

// Panic logs a panic message then panics.
func Panic(args ...interface{}) {
	Default.Panic(args...)
}

// PanicIf logs the error if present.
func PanicIf(err error) {
	Default.PanicIf(err)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return Default.Sync()
}
