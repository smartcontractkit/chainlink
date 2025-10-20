package logger

import (
	"os"
	"sync/atomic"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AtomicCore provides thread-safe core swapping using atomic operations.
// It starts as a noop core and can be atomically swapped to include additional cores.
var _ zapcore.Core = &AtomicCore{}

type AtomicCore struct {
	atomic.Pointer[zapcore.Core]
}

// NewAtomicCore creates a new AtomicCore initialized with a noop core
func NewAtomicCore() *AtomicCore {
	ac := &AtomicCore{}
	noop := zapcore.NewNopCore()
	ac.Store(&noop)
	return ac
}

func (d *AtomicCore) load() zapcore.Core {
	p := d.Load()
	if p == nil {
		return zapcore.NewNopCore()
	}
	return *p
}

func (d *AtomicCore) Enabled(l zapcore.Level) bool { return d.load().Enabled(l) }

func (d *AtomicCore) With(fs []zapcore.Field) zapcore.Core { return d.load().With(fs) }

func (d *AtomicCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return d.load().Check(e, ce)
}

func (d *AtomicCore) Write(e zapcore.Entry, fs []zapcore.Field) error { return d.load().Write(e, fs) }

func (d *AtomicCore) Sync() error { return d.load().Sync() }

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	level      zap.AtomicLevel
	fields     []any
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

func (l *zapLogger) With(args ...any) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.With(args...)
	newLogger.fields = copyFields(l.fields, args...)
	return &newLogger
}

// copyFields returns a copy of fields with add appended.
func copyFields(fields []any, add ...any) []any {
	f := make([]any, 0, len(fields)+len(add))
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

func (l *zapLogger) Recover(panicErr any) {
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}
