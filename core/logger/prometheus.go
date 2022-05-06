package logger

import (
	"fmt"
	"io"

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
	fatalCounter prometheus.Counter) Logger {
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

func (s *prometheusLogger) With(args ...interface{}) Logger {
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

func (s *prometheusLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	h, err := s.h.NewRootLogger(lvl)
	if err != nil {
		return nil, err
	}
	return &prometheusLogger{
		h:           h,
		warnCnt:     s.warnCnt,
		errorCnt:    s.errorCnt,
		criticalCnt: s.criticalCnt,
		panicCnt:    s.panicCnt,
		fatalCnt:    s.fatalCnt,
	}, nil
}

func (s *prometheusLogger) SetLogLevel(level zapcore.Level) {
	s.h.SetLogLevel(level)
}

func (s *prometheusLogger) Trace(args ...interface{}) {
	s.h.Trace(args...)
}

func (s *prometheusLogger) Debug(args ...interface{}) {
	s.h.Debug(args...)
}

func (s *prometheusLogger) Info(args ...interface{}) {
	s.h.Info(args...)
}

func (s *prometheusLogger) Warn(args ...interface{}) {
	s.warnCnt.Inc()
	s.h.Warn(args...)
}

func (s *prometheusLogger) Error(args ...interface{}) {
	s.errorCnt.Inc()
	s.h.Error(args...)
}

func (s *prometheusLogger) Critical(args ...interface{}) {
	s.criticalCnt.Inc()
	s.h.Critical(args...)
}

func (s *prometheusLogger) Panic(args ...interface{}) {
	s.panicCnt.Inc()
	s.h.Panic(args...)
}

func (s *prometheusLogger) Fatal(args ...interface{}) {
	s.fatalCnt.Inc()
	s.h.Fatal(args...)
}

func (s *prometheusLogger) Tracef(format string, values ...interface{}) {
	s.h.Tracef(format, values...)
}

func (s *prometheusLogger) Debugf(format string, values ...interface{}) {
	s.h.Debugf(format, values...)
}

func (s *prometheusLogger) Infof(format string, values ...interface{}) {
	s.h.Infof(format, values...)
}

func (s *prometheusLogger) Warnf(format string, values ...interface{}) {
	s.warnCnt.Inc()
	s.h.Warnf(format, values...)
}

func (s *prometheusLogger) Errorf(format string, values ...interface{}) {
	s.errorCnt.Inc()
	s.h.Errorf(format, values...)
}

func (s *prometheusLogger) Criticalf(format string, values ...interface{}) {
	s.criticalCnt.Inc()
	s.h.Criticalf(format, values...)
}

func (s *prometheusLogger) Panicf(format string, values ...interface{}) {
	s.panicCnt.Inc()
	s.h.Panicf(format, values...)
}

func (s *prometheusLogger) Fatalf(format string, values ...interface{}) {
	s.fatalCnt.Inc()
	s.h.Fatalf(format, values...)
}

func (s *prometheusLogger) Tracew(msg string, keysAndValues ...interface{}) {
	s.h.Tracew(msg, keysAndValues...)
}

func (s *prometheusLogger) Debugw(msg string, keysAndValues ...interface{}) {
	s.h.Debugw(msg, keysAndValues...)
}

func (s *prometheusLogger) Infow(msg string, keysAndValues ...interface{}) {
	s.h.Infow(msg, keysAndValues...)
}

func (s *prometheusLogger) Warnw(msg string, keysAndValues ...interface{}) {
	s.warnCnt.Inc()
	s.h.Warnw(msg, keysAndValues...)
}

func (s *prometheusLogger) Errorw(msg string, keysAndValues ...interface{}) {
	s.errorCnt.Inc()
	s.h.Errorw(msg, keysAndValues...)
}

func (s *prometheusLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	s.criticalCnt.Inc()
	s.h.Criticalw(msg, keysAndValues...)
}

func (s *prometheusLogger) Panicw(msg string, keysAndValues ...interface{}) {
	s.panicCnt.Inc()
	s.h.Panicw(msg, keysAndValues...)
}

func (s *prometheusLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	s.fatalCnt.Inc()
	s.h.Fatalw(msg, keysAndValues...)
}

func (s *prometheusLogger) ErrorIf(err error, msg string) {
	if err != nil {
		s.errorCnt.Inc()
		s.h.Errorw(msg, "err", err)
	}
}

func (s *prometheusLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		s.errorCnt.Inc()
		s.h.Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
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

func (s *prometheusLogger) Recover(panicErr interface{}) {
	s.panicCnt.Inc()
	s.h.Recover(panicErr)
}
