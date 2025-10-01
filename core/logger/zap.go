package logger

import (
	"os"

	pkgerrors "github.com/pkg/errors"
	otellog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger/otelzap"
)

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	level      zap.AtomicLevel
	fields     []interface{}
	callerSkip int
	opts       []zap.Option
	otelLogger otellog.Logger
}

func makeEncoderConfig(unixTS bool) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	if !unixTS {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	encoderConfig.EncodeLevel = encodeLevel

	return encoderConfig
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.level.SetLevel(lvl)
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

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	newLogger.Trace("Named logger created")
	return &newLogger
}

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) Name() string {
	return l.Desugar().Name()
}

func (l *zapLogger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.WithOptions(zap.AddCallerSkip(skip))
}

func (l *zapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return nil
	}
	var msg string
	if uw := pkgerrors.Unwrap(err); uw != nil {
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
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}

func (l *zapLogger) WithOtel(otelLogger otellog.Logger) (Logger, error) {
	if l.otelLogger != nil {
		return l, nil
	}
	l.otelLogger = otelLogger

	// Get the current core from the zap logger
	primaryCore := l.SugaredLogger.Desugar().Core()

	// Create OTel core with debug level to ensure all logs are captured
	otelCore := otelzap.NewCore(otelLogger, otelzap.WithLevel(zapcore.DebugLevel))

	// Create a new zap logger with both cores using Tee
	combinedCore := zapcore.NewTee(primaryCore, otelCore)
	newLogger := zap.New(combinedCore, l.opts...)

	// Update the zapLogger with the new core
	l.SugaredLogger = newLogger.Sugar()

	return l, nil
}
