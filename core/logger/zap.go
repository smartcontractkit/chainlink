package logger

import (
	"os"
	"sync"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// UpdatableCore provides a thread-safe zap Core updatable setup.
// It provides a root zap Core (returned by Root()) and tracks all its derived With(*) cores. It features a cascading
// update where the root and all its derived cores are updated with the given zap Core.
type UpdatableCore struct {
	root     *trackedCore
	registry *WeakRegistry[trackedCore]
}

func NewUpdatableCore() *UpdatableCore {
	registry := NewWeakRegistry[trackedCore]()
	root := NewTrackedCore(zapcore.NewNopCore(), []zapcore.Field{}, registry)
	return &UpdatableCore{root, registry}
}

func (a *UpdatableCore) Root() zapcore.Core {
	return a.root
}

func (a *UpdatableCore) Update(core zapcore.Core) {
	a.registry.Update(func(tc *trackedCore) {
		tc.Store(core)
	})
}

func (a *UpdatableCore) Close() {
	a.registry.Close()
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

var _ zapcore.Core = &trackedCore{}

// trackedCore is a zapcore.Core wrapper that registers the wrapped core and its derived With(*) loggers in the
// provided registry. It also tracks used zapcore.Fields allowing the underlying core updates.
type trackedCore struct {
	core     atomicCore
	fields   []zapcore.Field
	registry *WeakRegistry[trackedCore]
}

func NewTrackedCore(core zapcore.Core, fields []zapcore.Field, registry *WeakRegistry[trackedCore]) *trackedCore {
	cw := &trackedCore{core: NewAtomicCore(core), fields: fields, registry: registry}
	registry.Add(cw)
	return cw
}

func (c *trackedCore) Store(core zapcore.Core) {
	c.core.Store(core.With(c.fields))
}

func (c *trackedCore) With(fields []zapcore.Field) zapcore.Core {
	combined := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	combined = append(combined, c.fields...)
	combined = append(combined, fields...)

	tc := NewTrackedCore(c.core.Load().With(fields), combined, c.registry)
	return tc
}

func (c *trackedCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return c.core.Load().Check(ent, ce)
}

func (c *trackedCore) Enabled(lvl zapcore.Level) bool {
	return c.core.Load().Enabled(lvl)
}

func (c *trackedCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.core.Load().Write(ent, fields)
}

func (c *trackedCore) Sync() error {
	return c.core.Load().Sync()
}

// atomicCore is a minimal wrapper around zapcore.Core providing thread-safe Load and Store operations.
type atomicCore struct {
	mu  sync.RWMutex
	val zapcore.Core
}

func NewAtomicCore(v zapcore.Core) atomicCore {
	return atomicCore{val: v}
}

func (a *atomicCore) Store(v zapcore.Core) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.val = v
}

func (a *atomicCore) Load() zapcore.Core {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.val
}
