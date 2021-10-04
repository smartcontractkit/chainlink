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
	Default     Logger
	skipDefault Logger // Default.withCallerSkip(1) for helper funcs
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
	SetLogger(&zapLogger{SugaredLogger: zl.Sugar()})
}

// SetLogger sets the internal logger to the given input.
//
// DEPRECATED this method is deprecated because it leads to race conditions.
// Instead, you should fork the logger.Default instance to create a new logger
// for your module.
// Eg: logger.Default.Named("<my-package-name>")
func SetLogger(newLogger Logger) {
	if Default != nil {
		defer func(l Logger) {
			if err := l.Sync(); err != nil {
				if errors.Unwrap(err).Error() != os.ErrInvalid.Error() &&
					errors.Unwrap(err).Error() != "inappropriate ioctl for device" &&
					errors.Unwrap(err).Error() != "bad file descriptor" {
					// logger.Sync() will return 'invalid argument' error when closing file
					log.Fatalf("failed to sync logger %+v", err)
				}
			}
		}(Default)
	}
	Default = newLogger
	skipDefault = Default.withCallerSkip(1)
}

// Infow logs an info message and any additional given information.
func Infow(msg string, keysAndValues ...interface{}) {
	skipDefault.Infow(msg, keysAndValues...)
}

// Debugw logs a debug message and any additional given information.
func Debugw(msg string, keysAndValues ...interface{}) {
	skipDefault.Debugw(msg, keysAndValues...)
}

// Warnw logs a debug message and any additional given information.
func Warnw(msg string, keysAndValues ...interface{}) {
	skipDefault.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
func Errorw(msg string, keysAndValues ...interface{}) {
	skipDefault.Errorw(msg, keysAndValues...)
}

// Infof formats and then logs the message.
func Infof(format string, values ...interface{}) {
	skipDefault.Infof(format, values...)
}

// Debugf formats and then logs the message.
func Debugf(format string, values ...interface{}) {
	skipDefault.Debugf(format, values...)
}

// Tracef is a shim stand-in for when we have real trace-level logging support
func Tracef(format string, values ...interface{}) {
	skipDefault.Debug("TRACE: " + fmt.Sprintf(format, values...))
}

// Warnf formats and then logs the message as Warn.
func Warnf(format string, values ...interface{}) {
	skipDefault.Warnf(format, values...)
}

// Panicf formats and then logs the message before panicking.
func Panicf(format string, values ...interface{}) {
	skipDefault.Panic(fmt.Sprintf(format, values...))
}

// Info logs an info message.
func Info(args ...interface{}) {
	skipDefault.Info(args...)
}

// Debug logs a debug message.
func Debug(args ...interface{}) {
	skipDefault.Debug(args...)
}

// Warn logs a message at the warn level.
func Warn(args ...interface{}) {
	skipDefault.Warn(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	skipDefault.Error(args...)
}

func ErrorIf(err error, msg string) {
	skipDefault.ErrorIf(err, msg)
}

func ErrorIfCalling(f func() error) {
	skipDefault.ErrorIfCalling(f)
}

// Fatal logs a fatal message then exits the application.
func Fatal(args ...interface{}) {
	skipDefault.Fatal(args...)
}

// Errorf logs a message at the error level using Sprintf.
func Errorf(format string, values ...interface{}) {
	skipDefault.Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
func Fatalf(format string, values ...interface{}) {
	skipDefault.Fatal(fmt.Sprintf(format, values...))
}

// Fatalw logs a message and exits the application
func Fatalw(msg string, keysAndValues ...interface{}) {
	skipDefault.Fatalw(msg, keysAndValues...)
}

// Panic logs a panic message then panics.
func Panic(args ...interface{}) {
	skipDefault.Panic(args...)
}

// PanicIf logs the error if present.
func PanicIf(err error, msg string) {
	skipDefault.PanicIf(err, msg)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return Default.Sync()
}
