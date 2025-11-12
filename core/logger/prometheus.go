package logger

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

var warnCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "log_warn_count",
	Help: "Number of warning messages in log",
})

var errorCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "log_error_count",
	Help: "Number of error messages in log",
})

var criticalCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "log_critical_count",
	Help: "Number of critical messages in log",
})

var panicCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "log_panic_count",
	Help: "Number of panic messages in log",
})

var fatalCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "log_fatal_count",
	Help: "Number of fatal messages in log",
})

type prometheusLogger struct {
	h           Logger
	warnCnt     prometheus.Counter
	errorCnt    prometheus.Counter
	criticalCnt prometheus.Counter
	panicCnt    prometheus.Counter
	fatalCnt    prometheus.Counter
}

func newPrometheusLoggerWithCounters(
	l Logger,
	warnCounter prometheus.Counter,
	errorCounter prometheus.Counter,
	criticalCounter prometheus.Counter,
	panicCounter prometheus.Counter,
	fatalCounter prometheus.Counter,
) Logger {
	return &prometheusLogger{
		h:           l.Helper(1),
		warnCnt:     warnCounter,
		errorCnt:    errorCounter,
		criticalCnt: criticalCounter,
		panicCnt:    panicCounter,
		fatalCnt:    fatalCounter,
	}
}

func newPrometheusLogger(l Logger) Logger {
	return newPrometheusLoggerWithCounters(l, warnCounter, errorCounter, criticalCounter, panicCounter, fatalCounter)
}

func (s *prometheusLogger) With(args ...any) Logger {
	return &prometheusLogger{
		h:           s.h.With(args...),
		warnCnt:     s.warnCnt,
		errorCnt:    s.errorCnt,
		criticalCnt: s.criticalCnt,
		panicCnt:    s.panicCnt,
		fatalCnt:    s.fatalCnt,
	}
}

func (s *prometheusLogger) Named(name string) Logger {
	return &prometheusLogger{
		h:           s.h.Named(name),
		warnCnt:     s.warnCnt,
		errorCnt:    s.errorCnt,
		criticalCnt: s.criticalCnt,
		panicCnt:    s.panicCnt,
		fatalCnt:    s.fatalCnt,
	}
}

func (s *prometheusLogger) Name() string {
	return s.h.Name()
}

func (s *prometheusLogger) SetLogLevel(level zapcore.Level) {
	s.h.SetLogLevel(level)
}

func (s *prometheusLogger) Trace(args ...any) {
	s.h.Trace(args...)
}

func (s *prometheusLogger) Debug(args ...any) {
	s.h.Debug(args...)
}

func (s *prometheusLogger) Info(args ...any) {
	s.h.Info(args...)
}

func (s *prometheusLogger) Warn(args ...any) {
	s.warnCnt.Inc()
	s.h.Warn(args...)
}

func (s *prometheusLogger) Error(args ...any) {
	s.errorCnt.Inc()
	s.h.Error(args...)
}

func (s *prometheusLogger) Critical(args ...any) {
	s.criticalCnt.Inc()
	s.h.Critical(args...)
}

func (s *prometheusLogger) Panic(args ...any) {
	s.panicCnt.Inc()
	s.h.Panic(args...)
}

func (s *prometheusLogger) Fatal(args ...any) {
	s.fatalCnt.Inc()
	s.h.Fatal(args...)
}

func (s *prometheusLogger) Tracef(format string, values ...any) {
	s.h.Tracef(format, values...)
}

func (s *prometheusLogger) Debugf(format string, values ...any) {
	s.h.Debugf(format, values...)
}

func (s *prometheusLogger) Infof(format string, values ...any) {
	s.h.Infof(format, values...)
}

func (s *prometheusLogger) Warnf(format string, values ...any) {
	s.warnCnt.Inc()
	s.h.Warnf(format, values...)
}

func (s *prometheusLogger) Errorf(format string, values ...any) {
	s.errorCnt.Inc()
	s.h.Errorf(format, values...)
}

func (s *prometheusLogger) Criticalf(format string, values ...any) {
	s.criticalCnt.Inc()
	s.h.Criticalf(format, values...)
}

func (s *prometheusLogger) Panicf(format string, values ...any) {
	s.panicCnt.Inc()
	s.h.Panicf(format, values...)
}

func (s *prometheusLogger) Fatalf(format string, values ...any) {
	s.fatalCnt.Inc()
	s.h.Fatalf(format, values...)
}

func (s *prometheusLogger) Tracew(msg string, keysAndValues ...any) {
	s.h.Tracew(msg, keysAndValues...)
}

func (s *prometheusLogger) Debugw(msg string, keysAndValues ...any) {
	s.h.Debugw(msg, keysAndValues...)
}

func (s *prometheusLogger) Infow(msg string, keysAndValues ...any) {
	s.h.Infow(msg, keysAndValues...)
}

func (s *prometheusLogger) Warnw(msg string, keysAndValues ...any) {
	s.warnCnt.Inc()
	s.h.Warnw(msg, keysAndValues...)
}

func (s *prometheusLogger) Errorw(msg string, keysAndValues ...any) {
	s.errorCnt.Inc()
	s.h.Errorw(msg, keysAndValues...)
}

func (s *prometheusLogger) Criticalw(msg string, keysAndValues ...any) {
	s.criticalCnt.Inc()
	s.h.Criticalw(msg, keysAndValues...)
}

func (s *prometheusLogger) Panicw(msg string, keysAndValues ...any) {
	s.panicCnt.Inc()
	s.h.Panicw(msg, keysAndValues...)
}

func (s *prometheusLogger) Fatalw(msg string, keysAndValues ...any) {
	s.fatalCnt.Inc()
	s.h.Fatalw(msg, keysAndValues...)
}

func (s *prometheusLogger) Sync() error {
	return s.h.Sync()
}

func (s *prometheusLogger) Helper(add int) Logger {
	return &prometheusLogger{
		s.h.Helper(add),
		s.warnCnt,
		s.errorCnt,
		s.criticalCnt,
		s.panicCnt,
		s.fatalCnt,
	}
}

func (s *prometheusLogger) Recover(panicErr any) {
	s.panicCnt.Inc()
	s.h.Recover(panicErr)
}
