package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	config     zap.Config
	name       string
	fields     []interface{}
	callerSkip int
}

func newZapLogger(cfg zap.Config) (Logger, error) {
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &zapLogger{config: cfg, SugaredLogger: zl.Sugar()}, nil
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.config.Level.SetLevel(lvl)
}

func (l *zapLogger) With(args ...interface{}) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.With(args...)
	newLogger.fields = copyFields(l.fields, args...)
	return &newLogger
}

// copyFields returns a copy of fields with add appended.
func copyFields(fields []interface{}, add ...interface{}) []interface{} {
	f := make([]interface{}, 0, len(fields)+len(add))
	f = append(f, fields...)
	f = append(f, add...)
	return f
}

func joinName(old, new string) string {
	if old == "" {
		return new
	}
	return old + "." + new
}

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.name = joinName(l.name, name)
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	newLogger.Trace("Named logger created")
	return &newLogger
}

func (l *zapLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	newLogger := *l
	newLogger.config.Level = zap.NewAtomicLevelAt(lvl)
	zl, err := newLogger.config.Build()
	if err != nil {
		return nil, err
	}
	zl = zl.WithOptions(zap.AddCallerSkip(l.callerSkip))
	newLogger.SugaredLogger = zl.Named(l.name).Sugar().With(l.fields...)
	return &newLogger, nil
}

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.Helper(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.Helper(1).Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *zapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return nil
	}
	var msg string
	if uw := errors.Unwrap(err); uw != nil {
		msg = uw.Error()
	} else {
		msg = err.Error()
	}
	switch msg {
	case os.ErrInvalid.Error(), "bad file descriptor",
		"inappropriate ioctl for device":
		return nil
	}
	return err
}

func (l *zapLogger) Recover(panicErr interface{}) {
	l.CriticalW("Recovered goroutine panic", "panic", panicErr)
}
