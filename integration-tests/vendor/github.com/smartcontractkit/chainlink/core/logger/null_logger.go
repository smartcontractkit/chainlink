package logger

import (
	"io"

	"go.uber.org/zap/zapcore"
)

// nolint
var NullLogger Logger = &nullLogger{}

type nullLogger struct{}

func (l *nullLogger) With(args ...interface{}) Logger { return l }
func (l *nullLogger) Named(name string) Logger        { return l }
func (l *nullLogger) SetLogLevel(_ zapcore.Level)     {}

func (l *nullLogger) Trace(args ...interface{})    {}
func (l *nullLogger) Debug(args ...interface{})    {}
func (l *nullLogger) Info(args ...interface{})     {}
func (l *nullLogger) Warn(args ...interface{})     {}
func (l *nullLogger) Error(args ...interface{})    {}
func (l *nullLogger) Critical(args ...interface{}) {}
func (l *nullLogger) Panic(args ...interface{})    {}
func (l *nullLogger) Fatal(args ...interface{})    {}

func (l *nullLogger) Tracef(format string, values ...interface{})    {}
func (l *nullLogger) Debugf(format string, values ...interface{})    {}
func (l *nullLogger) Infof(format string, values ...interface{})     {}
func (l *nullLogger) Warnf(format string, values ...interface{})     {}
func (l *nullLogger) Errorf(format string, values ...interface{})    {}
func (l *nullLogger) Criticalf(format string, values ...interface{}) {}
func (l *nullLogger) Panicf(format string, values ...interface{})    {}
func (l *nullLogger) Fatalf(format string, values ...interface{})    {}

func (l *nullLogger) Tracew(msg string, keysAndValues ...interface{})    {}
func (l *nullLogger) Debugw(msg string, keysAndValues ...interface{})    {}
func (l *nullLogger) Infow(msg string, keysAndValues ...interface{})     {}
func (l *nullLogger) Warnw(msg string, keysAndValues ...interface{})     {}
func (l *nullLogger) Errorw(msg string, keysAndValues ...interface{})    {}
func (l *nullLogger) Criticalw(msg string, keysAndValues ...interface{}) {}
func (l *nullLogger) Panicw(msg string, keysAndValues ...interface{})    {}
func (l *nullLogger) Fatalw(msg string, keysAndValues ...interface{})    {}

func (l *nullLogger) ErrorIf(err error, msg string)    {}
func (l *nullLogger) ErrorIfClosing(io.Closer, string) {}
func (l *nullLogger) Sync() error                      { return nil }
func (l *nullLogger) Helper(skip int) Logger           { return l }
func (l *nullLogger) Name() string                     { return "nullLogger" }

func (l *nullLogger) Recover(panicErr interface{}) {}
