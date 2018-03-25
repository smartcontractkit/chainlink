// Package logger is used to store details of events in the node.
// Events can be categorized by Debug, Info, Error, Fatal, and Panic.
package logger

import (
	"fmt"
	"log"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *Logger

func init() {
	zl, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	SetLogger(NewLogger(zl))
}

// Logger holds a field for the logger interface.
type Logger struct {
	*zap.SugaredLogger
}

// Write logs a message at the Info level and returns the length
// of the given bytes.
func (l *Logger) Write(b []byte) (n int, err error) {
	l.Info(string(b))
	return len(b), nil
}

// NewLogger returns the logger updated with the given Logger.
func NewLogger(zl *zap.Logger) *Logger {
	return &Logger{zl.Sugar()}
}

// SetLogger sets the internal logger to the given input.
func SetLogger(l *Logger) {
	if logger != nil {
		defer logger.Sync()
	}
	logger = l
}

// Reconfigure creates a new log file at the configured directory
// with the given LogLevel.
func Reconfigure(dir string, lvl zapcore.Level) {
	config := generateConfig(dir)
	config.Level.SetLevel(lvl)
	zl, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	SetLogger(NewLogger(zl))
}

func generateConfig(dir string) zap.Config {
	config := zap.NewProductionConfig()
	destination := path.Join(dir, "log.jsonl")
	config.OutputPaths = []string{"stderr", destination}
	config.ErrorOutputPaths = []string{"stderr", destination}
	return config
}

// Infow logs an info message and any additional given information.
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// Debugw logs a debug message and any additional given information.
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

// Warnw logs a debug message and any additional given information.
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message, any additional given information, and includes
// stack trace.
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

// Panicf formats and then logs the message before panicking.
func Panicf(format string, values ...interface{}) {
	logger.Panic(fmt.Sprintf(format, values...))
}

// Info logs an info message using Sprint.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug logs an debug message using Sprint.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Warn logs a message at the warn level using Sprint.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error logs an error message using Sprint.
func Error(args ...interface{}) {
	logger.Error(args...)
}

//WarnIf logs the error if present.
func WarnIf(err error) {
	if err != nil {
		logger.Warn(err)
	}
}

// Fatal logs a fatal message then exits the application using Sprint.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Panic logs a panic message then panics using Sprint.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return logger.Sync()
}
