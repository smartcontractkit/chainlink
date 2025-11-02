package logger

import (
	"go.uber.org/zap/zapcore"
)

var NullLogger Logger = &nullLogger{}

type nullLogger struct{}

func (l *nullLogger) With(args ...any) Logger     { return l }
func (l *nullLogger) Named(name string) Logger    { return l }
func (l *nullLogger) SetLogLevel(_ zapcore.Level) {}

func (l *nullLogger) Trace(args ...any)    {}
func (l *nullLogger) Debug(args ...any)    {}
func (l *nullLogger) Info(args ...any)     {}
func (l *nullLogger) Warn(args ...any)     {}
func (l *nullLogger) Error(args ...any)    {}
func (l *nullLogger) Critical(args ...any) {}
func (l *nullLogger) Panic(args ...any)    {}
func (l *nullLogger) Fatal(args ...any)    {}

func (l *nullLogger) Tracef(format string, values ...any)    {}
func (l *nullLogger) Debugf(format string, values ...any)    {}
func (l *nullLogger) Infof(format string, values ...any)     {}
func (l *nullLogger) Warnf(format string, values ...any)     {}
func (l *nullLogger) Errorf(format string, values ...any)    {}
func (l *nullLogger) Criticalf(format string, values ...any) {}
func (l *nullLogger) Panicf(format string, values ...any)    {}
func (l *nullLogger) Fatalf(format string, values ...any)    {}

func (l *nullLogger) Tracew(msg string, keysAndValues ...any)    {}
func (l *nullLogger) Debugw(msg string, keysAndValues ...any)    {}
func (l *nullLogger) Infow(msg string, keysAndValues ...any)     {}
func (l *nullLogger) Warnw(msg string, keysAndValues ...any)     {}
func (l *nullLogger) Errorw(msg string, keysAndValues ...any)    {}
func (l *nullLogger) Criticalw(msg string, keysAndValues ...any) {}
func (l *nullLogger) Panicw(msg string, keysAndValues ...any)    {}
func (l *nullLogger) Fatalw(msg string, keysAndValues ...any)    {}

func (l *nullLogger) Sync() error            { return nil }
func (l *nullLogger) Helper(skip int) Logger { return l }
func (l *nullLogger) Name() string           { return "nullLogger" }

func (l *nullLogger) Recover(panicErr any) {}
