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
	logSql, _ := strconv.ParseBool(os.Getenv("LOG_SQL"))

	l := newLogger(envLvl, os.Getenv("ROOT"), jsonConsole, toDisk, unixTS, logSql)
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

// Errorf logs a message at the error level using Sprintf.
func Errorf(format string, values ...interface{}) {
	helper.Error(fmt.Sprintf(format, values...))
}

// Fatalf logs a message at the fatal level using Sprintf.
func Fatalf(format string, values ...interface{}) {
	helper.Fatal(fmt.Sprintf(format, values...))
}
