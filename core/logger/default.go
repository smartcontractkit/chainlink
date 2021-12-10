package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	jsonConsole, _ := strconv.ParseBool(os.Getenv("JSON_CONSOLE"))
	unixTS, _ := strconv.ParseBool(os.Getenv("LOG_UNIX_TS"))
	toDisk, _ := strconv.ParseBool(os.Getenv("LOG_TO_DISK"))

	logDir := os.Getenv("LOG_FILE_DIR")
	if logDir == "" {
		logDir = os.Getenv("ROOT")
	}
	l := newLogger(envLvl, logDir, jsonConsole, toDisk, unixTS)
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
	helper = newLogger.Helper(1)
}

// Warnw logs a debug message and any additional given information.
// Deprecated
func Warnw(msg string, keysAndValues ...interface{}) {
	helper.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
// Deprecated
func Errorw(msg string, keysAndValues ...interface{}) {
	helper.Errorw(msg, keysAndValues...)
}

// Warnf formats and then logs the message as Warn.
// Deprecated
func Warnf(format string, values ...interface{}) {
	helper.Warnf(format, values...)
}

// Warn logs a message at the warn level.
// Deprecated
func Warn(args ...interface{}) {
	helper.Warn(args...)
}

// Error logs an error message.
// Deprecated
func Error(args ...interface{}) {
	helper.Error(args...)
}

// Errorf logs a message at the error level using Sprintf.
// Deprecated
func Errorf(format string, values ...interface{}) {
	helper.Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
// Deprecated
func Fatalf(format string, values ...interface{}) {
	helper.Fatal(fmt.Sprintf(format, values...))
}
