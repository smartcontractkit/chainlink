package logger

import (
	"os"
	"sync"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AtomicCore provides thread-safe core swapping using atomic operations.
// It starts as a noop core and can be atomically swapped to include additional cores.
var _ zapcore.Core = &AtomicCore{}

type AtomicCore struct {
	mu           sync.RWMutex
	rootCore     *VersionedValue[zapcore.Core]
	localCore    zapcore.Core
	localVersion int64
	fields       []zapcore.Field
}

func NewAtomicCore() *AtomicCore {
	rootCore := &VersionedValue[zapcore.Core]{}
	rootCore.Store(zapcore.NewNopCore())
	return &AtomicCore{rootCore: rootCore}
}

func (a *AtomicCore) With(fields []zapcore.Field) zapcore.Core {
	combined := make([]zapcore.Field, 0, len(a.fields)+len(fields))
	combined = append(combined, a.fields...)
	combined = append(combined, fields...)

	return &AtomicCore{rootCore: a.rootCore, fields: combined}
}

func (a *AtomicCore) Store(core zapcore.Core) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.localVersion = a.rootCore.Store(core)
	a.localCore = core.With(a.fields)
}

func (a *AtomicCore) load() zapcore.Core {
	// Usual path: read localCore at the latest version with read lock only.
	a.mu.RLock()
	rootCore, version := a.rootCore.Load()
	localCore := a.localCore
	lastVerison := a.localVersion
	a.mu.RUnlock()
	if localCore != nil && lastVerison == version {
		return localCore
	}
	// Update path: need to read the latest version from the rootCore first and then update localCore.
	a.mu.Lock()
	defer a.mu.Unlock()
	rootCore, version = a.rootCore.Load()
	a.localCore = rootCore.With(a.fields)
	a.localVersion = version
	return a.localCore
}

func (a *AtomicCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return a.load().Check(ent, ce)
}

func (a *AtomicCore) Enabled(lvl zapcore.Level) bool {
	return a.load().Enabled(lvl)
}

func (a *AtomicCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return a.load().Write(ent, fields)
}

func (a *AtomicCore) Sync() error {
	return a.load().Sync()
}

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
	return l.WithOptions(zap.AddCallerSkip(skip))
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
