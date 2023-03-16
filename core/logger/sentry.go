package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

const (
	// SentryFlushDeadline indicates the maximum amount of time we allow sentry to
	// flush events on manual flush
	SentryFlushDeadline = 5 * time.Second

	loggerContextName = "Logger"
)

type sentryLogger struct {
	h Logger
}

func newSentryLogger(l Logger) Logger {
	return &sentryLogger{h: l.Helper(1)}
}

func (s *sentryLogger) With(args ...interface{}) Logger {
	return &sentryLogger{
		h: s.h.With(args...),
	}
}

func (s *sentryLogger) Named(name string) Logger {
	return &sentryLogger{
		h: s.h.Named(name),
	}
}

func (s *sentryLogger) Name() string {
	return s.h.Name()
}

func (s *sentryLogger) SetLogLevel(level zapcore.Level) {
	s.h.SetLogLevel(level)
}

func (s *sentryLogger) Trace(args ...interface{}) {
	s.h.Trace(args...)
}

func (s *sentryLogger) Debug(args ...interface{}) {
	s.h.Debug(args...)
}

func (s *sentryLogger) Info(args ...interface{}) {
	s.h.Info(args...)
}

func (s *sentryLogger) Warn(args ...interface{}) {
	s.h.Warn(args...)
}

func (s *sentryLogger) Error(args ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelError)
	})
	eid := hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.With("sentryEventID", eid).Error(args...)
}

func (s *sentryLogger) Critical(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.With("sentryEventID", eid).Critical(args...)
}

func (s *sentryLogger) Panic(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.With("sentryEventID", eid).Panic(args...)
}

func (s *sentryLogger) Fatal(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf("%v", args))
	s.h.With("sentryEventID", eid).Fatal(args...)
}

func (s *sentryLogger) Tracef(format string, values ...interface{}) {
	s.h.Tracef(format, values...)
}

func (s *sentryLogger) Debugf(format string, values ...interface{}) {
	s.h.Debugf(format, values...)
}

func (s *sentryLogger) Infof(format string, values ...interface{}) {
	s.h.Infof(format, values...)
}

func (s *sentryLogger) Warnf(format string, values ...interface{}) {
	s.h.Warnf(format, values...)
}

func (s *sentryLogger) Errorf(format string, values ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelError)
	})
	eid := hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.With("sentryEventID", eid).Errorf(format, values...)
}

func (s *sentryLogger) Criticalf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.With("sentryEventID", eid).Criticalf(format, values...)
}

func (s *sentryLogger) Panicf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.With("sentryEventID", eid).Panicf(format, values...)
}

func (s *sentryLogger) Fatalf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(fmt.Sprintf(format, values...))
	s.h.With("sentryEventID", eid).Fatalf(format, values...)
}

func (s *sentryLogger) Tracew(msg string, keysAndValues ...interface{}) {
	s.h.Tracew(msg, keysAndValues...)
}

func (s *sentryLogger) Debugw(msg string, keysAndValues ...interface{}) {
	s.h.Debugw(msg, keysAndValues...)
}

func (s *sentryLogger) Infow(msg string, keysAndValues ...interface{}) {
	s.h.Infow(msg, keysAndValues...)
}

func (s *sentryLogger) Warnw(msg string, keysAndValues ...interface{}) {
	s.h.Warnw(msg, keysAndValues...)
}

func (s *sentryLogger) Errorw(msg string, keysAndValues ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, toMap(keysAndValues))
		scope.SetLevel(sentry.LevelError)
	})
	eid := hub.CaptureMessage(msg)
	s.h.Errorw(msg, append(keysAndValues, "sentryEventID", eid)...)
}

func (s *sentryLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(msg)
	s.h.Criticalw(msg, append(keysAndValues, "sentryEventID", eid)...)
}

func (s *sentryLogger) Panicw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(msg)
	s.h.Panicw(msg, append(keysAndValues, "sentryEventID", eid)...)
}

func (s *sentryLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(loggerContextName, toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	eid := hub.CaptureMessage(msg)
	s.h.Fatalw(msg, append(keysAndValues, "sentryEventID", eid)...)
}

func (s *sentryLogger) ErrorIf(err error, msg string) {
	if err != nil {
		eid := sentry.CaptureException(err)
		s.h.Errorw(msg, "err", err, "sentryEventID", eid)
	}
}

func (s *sentryLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		eid := sentry.CaptureException(err)
		s.h.Errorw(fmt.Sprintf("Error closing %s", name), "err", err, "sentryEventID", eid)
	}
}

func (s *sentryLogger) Sync() error {
	return s.h.Sync()
}

func (s *sentryLogger) Helper(add int) Logger {
	return &sentryLogger{s.h.Helper(add)}
}

func toMap(args []interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); {
		// Make sure this element isn't a dangling key
		if i == len(args)-1 {
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string ignore it
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); ok {
			m[keyStr] = val
		}
		i += 2
	}
	return m
}

func (s *sentryLogger) Recover(panicErr interface{}) {
	eid := sentry.CurrentHub().Recover(panicErr)
	sentry.Flush(SentryFlushDeadline)

	s.h.With("sentryEventID", eid).Recover(panicErr)
}
