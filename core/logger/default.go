package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

// Logger used from package level helper functions.
var helper Logger

func init() {
	err := zap.RegisterSink("pretty", prettyConsoleSink(os.Stderr))
	if err != nil {
		log.Fatalf("failed to register pretty printer %+v", err)
	}
	err = registerOSSinks()
	if err != nil {
		log.Fatalf("failed to register os specific sinks %+v", err)
	}

	// HACK: This logic is a bit duplicated from newProductionConfig but the config object is not available in init
	// To be removed with https://app.shortcut.com/chainlinklabs/story/18500/logger-injection
	jsonStr := os.Getenv("JSON_CONSOLE")
	var jsonConsole bool
	if jsonStr == "true" {
		jsonConsole = true
	}
	unixTSStr := os.Getenv("LOG_UNIX_TS")
	var unixTS bool
	if unixTSStr == "true" {
		unixTS = true
	}
	toDiskStr := os.Getenv("LOG_TO_DISK")
	var toDisk bool
	if toDiskStr == "true" {
		toDisk = true
	}

	l, err := newZapLogger(newProductionConfig(os.Getenv("ROOT"), jsonConsole, toDisk, unixTS))
	if err != nil {
		log.Fatal(err)
	}
	InitLogger(l)
}

// InitLogger sets the helper Logger to newLogger. Not safe for concurrent use,
// so must be called from init() or the main goroutine during initialization.
//
// You probably don't want to use this. Loggers should be injected instead.
// Deprecated
func InitLogger(newLogger Logger) {
	if helper != nil {
		defer func(l Logger) {
			if err := l.Sync(); err != nil {
				// logger.Sync() will return 'invalid argument' error when closing file
				newLogger.Fatalf("failed to sync logger %+v", err)
			}
		}(helper)
	}
	helper = newLogger.withCallerSkip(1)
}

// Infow logs an info message and any additional given information.
func Infow(msg string, keysAndValues ...interface{}) {
	helper.Infow(msg, keysAndValues...)
}

// Debugw logs a debug message and any additional given information.
func Debugw(msg string, keysAndValues ...interface{}) {
	helper.Debugw(msg, keysAndValues...)
}

// Warnw logs a debug message and any additional given information.
func Warnw(msg string, keysAndValues ...interface{}) {
	helper.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
func Errorw(msg string, keysAndValues ...interface{}) {
	helper.Errorw(msg, keysAndValues...)
}

// Warnf formats and then logs the message as Warn.
func Warnf(format string, values ...interface{}) {
	helper.Warnf(format, values...)
}

// Warn logs a message at the warn level.
func Warn(args ...interface{}) {
	helper.Warn(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	helper.Error(args...)
}

func ErrorIfCalling(f func() error) {
	helper.ErrorIfCalling(f)
}

// Errorf logs a message at the error level using Sprintf.
func Errorf(format string, values ...interface{}) {
	helper.Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
func Fatalf(format string, values ...interface{}) {
	helper.Fatal(fmt.Sprintf(format, values...))
}

// Fatalw logs a message and exits the application
func Fatalw(msg string, keysAndValues ...interface{}) {
	helper.Fatalw(msg, keysAndValues...)
}
