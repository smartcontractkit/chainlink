package logger

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
)

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	level      zap.AtomicLevel
	name       string
	fields     []interface{}
	callerSkip int
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

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) Name() string {
	return l.name
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
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}

// loggerNameOverrideCore promotes any string field with key loggerName to the [zapcore.Entry.LoggerName] and removes it
// from the field set. If multiple matches are found, the last is used.
type loggerNameOverrideCore struct {
	zapcore.Core
	loggerName string
}

var (
	_ zapcore.Core = (*loggerNameOverrideCore)(nil)
)

func (c *loggerNameOverrideCore) Level() zapcore.Level {
	return zapcore.LevelOf(c.Core)
}

func (c *loggerNameOverrideCore) With(fields []zapcore.Field) zapcore.Core {
	return &loggerNameOverrideCore{c.Core.With(fields), c.loggerName}
}

func (c *loggerNameOverrideCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	fn := func(field zapcore.Field) bool {
		return field.Key == c.loggerName && field.String != ""
	}
	for i := slices.IndexFunc(fields, fn); i >= 0; i = slices.IndexFunc(fields, fn) {
		ent.LoggerName = fields[i].String
		fields = slices.Delete(fields, i, i+1)
	}
	return c.Core.Write(ent, fields)
}
